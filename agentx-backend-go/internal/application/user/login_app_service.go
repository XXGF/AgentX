package user

import (
	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/user"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	infraEmail "github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/email"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// LoginAppService 登录应用服务（对应 Java 的 LoginAppService）
type LoginAppService struct {
	domainService       *domain.DomainService
	jwtUtils            *auth.JWTUtils
	emailService        infraEmail.EmailService
	verificationService *infraEmail.VerificationCodeService
	captchaService      *infraEmail.CaptchaService
}

// NewLoginAppService 创建登录应用服务
func NewLoginAppService(
	domainService *domain.DomainService,
	jwtUtils *auth.JWTUtils,
	emailService infraEmail.EmailService,
	verificationService *infraEmail.VerificationCodeService,
	captchaService *infraEmail.CaptchaService,
) *LoginAppService {
	return &LoginAppService{
		domainService:       domainService,
		jwtUtils:            jwtUtils,
		emailService:        emailService,
		verificationService: verificationService,
		captchaService:      captchaService,
	}
}

// Login 登录（对应 Java 的 login）
func (s *LoginAppService) Login(account, password string) (string, error) {
	user, err := s.domainService.Login(account, password)
	if err != nil {
		return "", err
	}

	token, err := s.jwtUtils.GenerateToken(user.ID)
	if err != nil {
		return "", exception.NewBusinessExceptionWithCause("生成Token失败", err)
	}

	return token, nil
}

// Register 注册（对应 Java 的 register）
func (s *LoginAppService) Register(email, phone, password, code string) error {
	// 如果是邮箱注册，需要验证码
	if email != "" && phone == "" {
		if code == "" {
			return exception.NewBusinessException("邮箱注册需要验证码")
		}
		if !s.verificationService.VerifyCode(email, code, "register") {
			return exception.NewBusinessException("验证码无效或已过期")
		}
	}

	_, err := s.domainService.Register(email, phone, password)
	return err
}

// SendEmailVerificationCode 发送注册邮箱验证码（对应 Java 的 sendEmailVerificationCode）
func (s *LoginAppService) SendEmailVerificationCode(email, captchaUuid, captchaCode, ip string) error {
	// 1. 验证图形验证码
	if !s.captchaService.VerifyCaptcha(captchaUuid, captchaCode) {
		return exception.NewBusinessException("图形验证码错误")
	}

	// 2. 检查邮箱是否已注册
	existingUser, err := s.domainService.FindUserByAccount(email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return exception.NewBusinessException("该邮箱已注册")
	}

	// 3. 生成并发送验证码
	code := s.verificationService.GenerateCode(email, "register")
	if err := s.emailService.SendVerificationCode(email, code, "注册验证码"); err != nil {
		return exception.NewBusinessExceptionWithCause("发送验证码失败", err)
	}

	return nil
}

// SendResetPasswordCode 发送重置密码邮箱验证码（对应 Java 的 sendResetPasswordCode）
func (s *LoginAppService) SendResetPasswordCode(email, captchaUuid, captchaCode, ip string) error {
	// 1. 验证图形验证码
	if !s.captchaService.VerifyCaptcha(captchaUuid, captchaCode) {
		return exception.NewBusinessException("图形验证码错误")
	}

	// 2. 检查邮箱是否存在
	existingUser, err := s.domainService.FindUserByAccount(email)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return exception.NewBusinessException("该邮箱未注册")
	}

	// 3. 生成并发送验证码
	code := s.verificationService.GenerateCode(email, "reset_password")
	if err := s.emailService.SendResetPasswordCode(email, code); err != nil {
		return exception.NewBusinessExceptionWithCause("发送验证码失败", err)
	}

	return nil
}

// VerifyEmailCode 验证邮箱验证码（对应 Java 的 verifyEmailCode）
func (s *LoginAppService) VerifyEmailCode(email, code string) bool {
	return s.verificationService.VerifyCode(email, code, "register")
}

// VerifyResetPasswordCode 验证重置密码邮箱验证码
func (s *LoginAppService) VerifyResetPasswordCode(email, code string) bool {
	return s.verificationService.VerifyCode(email, code, "reset_password")
}

// ResetPassword 重置密码（对应 Java 的 resetPassword）
func (s *LoginAppService) ResetPassword(email, newPassword, code string) error {
	// 验证重置密码验证码
	if !s.verificationService.VerifyCode(email, code, "reset_password") {
		return exception.NewBusinessException("验证码无效或已过期")
	}

	user, err := s.domainService.FindUserByAccount(email)
	if err != nil {
		return err
	}
	if user == nil {
		return exception.NewBusinessException("用户不存在")
	}

	return s.domainService.UpdatePassword(user.ID, newPassword)
}
