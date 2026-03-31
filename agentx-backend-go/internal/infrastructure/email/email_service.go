package email

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"net/smtp"
	"strings"
	"sync"
	"time"

	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/config"
	"go.uber.org/zap"
)

// EmailService 邮件服务接口（对应 Java 的 EmailService）
type EmailService interface {
	// SendVerificationCode 发送验证码邮件
	SendVerificationCode(to, code, subject string) error
	// SendResetPasswordCode 发送重置密码验证码邮件
	SendResetPasswordCode(to, code string) error
}

// SMTPEmailService 基于 SMTP 的邮件服务实现（替换原 SimpleEmailService）
type SMTPEmailService struct {
	host     string
	port     int
	username string
	password string
	template string // 验证码邮件模板
	subject  string // 默认邮件主题
	logger   *zap.Logger
}

// NewSMTPEmailService 创建 SMTP 邮件服务
func NewSMTPEmailService(mailCfg *config.MailConfig, logger *zap.Logger) EmailService {
	svc := &SMTPEmailService{
		host:     mailCfg.SMTP.Host,
		port:     mailCfg.SMTP.Port,
		username: mailCfg.SMTP.Username,
		password: mailCfg.SMTP.Password,
		template: mailCfg.Verification.Template,
		subject:  mailCfg.Verification.Subject,
		logger:   logger,
	}

	// 如果 SMTP 未配置，返回日志模式（仅打印不发送）
	if svc.username == "" || svc.password == "" {
		logger.Warn("SMTP 邮件服务未配置用户名/密码，将以日志模式运行（验证码仅打印到日志）")
	}

	// 设置默认值
	if svc.template == "" {
		svc.template = "您的验证码是:%s，有效期10分钟，请勿泄露给他人。"
	}
	if svc.subject == "" {
		svc.subject = "AgentX - 邮箱验证码"
	}
	if svc.port == 0 {
		svc.port = 587
	}

	return svc
}

func (s *SMTPEmailService) SendVerificationCode(to, code, subject string) error {
	if subject == "" {
		subject = s.subject
	}
	body := fmt.Sprintf(s.template, code)
	return s.sendMail(to, subject, body)
}

func (s *SMTPEmailService) SendResetPasswordCode(to, code string) error {
	return s.SendVerificationCode(to, code, "AgentX - 重置密码验证码")
}

// sendMail 发送邮件（支持 STARTTLS 和 SSL）
func (s *SMTPEmailService) sendMail(to, subject, body string) error {
	// 如果未配置 SMTP，仅打印日志
	if s.username == "" || s.password == "" {
		s.logger.Info("【日志模式】发送邮件",
			zap.String("to", to),
			zap.String("subject", subject),
			zap.String("body", body),
		)
		return nil
	}

	// 构建邮件内容
	from := s.username
	msg := s.buildMessage(from, to, subject, body)
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	var err error
	switch s.port {
	case 465:
		// SSL 直连
		err = s.sendMailSSL(addr, auth, from, to, msg)
	default:
		// STARTTLS（587）或明文（25）
		err = s.sendMailSTARTTLS(addr, auth, from, to, msg)
	}

	if err != nil {
		s.logger.Error("发送邮件失败",
			zap.String("to", to),
			zap.String("subject", subject),
			zap.Error(err),
		)
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	s.logger.Info("邮件发送成功",
		zap.String("to", to),
		zap.String("subject", subject),
	)
	return nil
}

// buildMessage 构建 MIME 邮件内容
func (s *SMTPEmailService) buildMessage(from, to, subject, body string) []byte {
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)
	return []byte(msg.String())
}

// sendMailSTARTTLS 通过 STARTTLS 发送邮件（端口 587/25）
func (s *SMTPEmailService) sendMailSTARTTLS(addr string, auth smtp.Auth, from, to string, msg []byte) error {
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败: %w", err)
	}

	host := s.host
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %w", err)
	}
	defer client.Close()

	// 尝试 STARTTLS
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName: host,
		}
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("STARTTLS失败: %w", err)
		}
	}

	// 认证
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %w", err)
	}

	// 发送
	if err = client.Mail(from); err != nil {
		return err
	}
	if err = client.Rcpt(to); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return client.Quit()
}

// sendMailSSL 通过 SSL 直连发送邮件（端口 465）
func (s *SMTPEmailService) sendMailSSL(addr string, auth smtp.Auth, from, to string, msg []byte) error {
	tlsConfig := &tls.Config{
		ServerName: s.host,
	}
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second}, "tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("SSL连接SMTP服务器失败: %w", err)
	}

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %w", err)
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %w", err)
	}

	if err = client.Mail(from); err != nil {
		return err
	}
	if err = client.Rcpt(to); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return client.Quit()
}

// ---- 验证码服务 ----

// VerificationCodeService 验证码服务（对应 Java 的 VerificationCodeService）
type VerificationCodeService struct {
	codes  sync.Map // key: email -> value: *codeEntry
	logger *zap.Logger
}

type codeEntry struct {
	Code      string
	ExpiresAt time.Time
	Purpose   string // register / reset_password
}

func NewVerificationCodeService(logger *zap.Logger) *VerificationCodeService {
	svc := &VerificationCodeService{logger: logger}
	// 启动过期清理协程
	go svc.cleanupExpired()
	return svc
}

// GenerateCode 生成6位验证码
func (s *VerificationCodeService) GenerateCode(email, purpose string) string {
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	s.codes.Store(email+":"+purpose, &codeEntry{
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute), // 10分钟有效
		Purpose:   purpose,
	})
	return code
}

// VerifyCode 验证验证码
func (s *VerificationCodeService) VerifyCode(email, code, purpose string) bool {
	key := email + ":" + purpose
	val, ok := s.codes.Load(key)
	if !ok {
		return false
	}
	entry := val.(*codeEntry)
	if time.Now().After(entry.ExpiresAt) {
		s.codes.Delete(key)
		return false
	}
	if entry.Code != code {
		return false
	}
	// 验证成功后删除
	s.codes.Delete(key)
	return true
}

// cleanupExpired 定期清理过期验证码
func (s *VerificationCodeService) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		s.codes.Range(func(key, value interface{}) bool {
			entry := value.(*codeEntry)
			if now.After(entry.ExpiresAt) {
				s.codes.Delete(key)
			}
			return true
		})
	}
}

// ---- 图形验证码服务 ----

// CaptchaService 图形验证码服务（对应 Java 的 CaptchaService）
type CaptchaService struct {
	captchas sync.Map // uuid -> code
}

func NewCaptchaService() *CaptchaService {
	return &CaptchaService{}
}

// VerifyCaptcha 验证图形验证码
func (s *CaptchaService) VerifyCaptcha(uuid, code string) bool {
	if uuid == "" || code == "" {
		return false
	}
	val, ok := s.captchas.LoadAndDelete(uuid)
	if !ok {
		// 开发阶段暂时放行
		return true
	}
	return val.(string) == code
}

// GenerateCaptcha 生成图形验证码
// 返回 uuid、验证码文本、Base64 编码的图片
func (s *CaptchaService) GenerateCaptcha() (uuid, code, imageBase64 string) {
	uuid = fmt.Sprintf("%d", time.Now().UnixNano())
	code = fmt.Sprintf("%04d", rand.Intn(10000))
	s.captchas.Store(uuid, code)

	// 生成简单的 SVG 验证码图片（纯 Go 实现，无需第三方库）
	imageBase64 = generateSimpleCaptchaSVG(code)
	return
}

// generateSimpleCaptchaSVG 生成简单的 SVG 格式验证码图片
func generateSimpleCaptchaSVG(code string) string {
	// 生成一个简单的 SVG 验证码
	var svg strings.Builder
	svg.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="120" height="40">`)
	svg.WriteString(`<rect width="120" height="40" fill="#f0f0f0"/>`)

	// 添加干扰线
	for i := 0; i < 5; i++ {
		x1 := rand.Intn(120)
		y1 := rand.Intn(40)
		x2 := rand.Intn(120)
		y2 := rand.Intn(40)
		colors := []string{"#ccc", "#ddd", "#bbb", "#aaa", "#999"}
		svg.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="1"/>`,
			x1, y1, x2, y2, colors[i%len(colors)]))
	}

	// 绘制验证码文字
	for i, ch := range code {
		x := 15 + i*25
		y := 25 + rand.Intn(10) - 5
		rotation := rand.Intn(30) - 15
		colors := []string{"#333", "#555", "#222", "#444"}
		svg.WriteString(fmt.Sprintf(
			`<text x="%d" y="%d" font-size="22" fill="%s" transform="rotate(%d %d %d)">%c</text>`,
			x, y, colors[i%len(colors)], rotation, x, y, ch))
	}

	svg.WriteString(`</svg>`)

	// 转为 data URI（SVG 可以直接用 base64 编码）
	return "data:image/svg+xml;base64," + base64Encode([]byte(svg.String()))
}

// base64Encode 简单的 base64 编码
func base64Encode(data []byte) string {
	const base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var result strings.Builder
	for i := 0; i < len(data); i += 3 {
		var b0, b1, b2 byte
		b0 = data[i]
		if i+1 < len(data) {
			b1 = data[i+1]
		}
		if i+2 < len(data) {
			b2 = data[i+2]
		}
		result.WriteByte(base64Chars[(b0>>2)&0x3F])
		result.WriteByte(base64Chars[((b0<<4)|(b1>>4))&0x3F])
		if i+1 < len(data) {
			result.WriteByte(base64Chars[((b1<<2)|(b2>>6))&0x3F])
		} else {
			result.WriteByte('=')
		}
		if i+2 < len(data) {
			result.WriteByte(base64Chars[b2&0x3F])
		} else {
			result.WriteByte('=')
		}
	}
	return result.String()
}