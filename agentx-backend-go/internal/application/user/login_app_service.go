package user

import (
	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/user"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// LoginAppService 登录应用服务（对应 Java 的 LoginAppService）
type LoginAppService struct {
	domainService *domain.DomainService
	jwtUtils      *auth.JWTUtils
	// TODO: 后续迁移时添加
	// emailService          EmailService
	// verificationService   VerificationCodeService
	// authSettingService    AuthSettingDomainService
}

// NewLoginAppService 创建登录应用服务
func NewLoginAppService(domainService *domain.DomainService, jwtUtils *auth.JWTUtils) *LoginAppService {
	return &LoginAppService{
		domainService: domainService,
		jwtUtils:      jwtUtils,
	}
}

// Login 登录（对应 Java 的 login）
func (s *LoginAppService) Login(account, password string) (string, error) {
	// TODO: 检查普通登录是否启用（需要 AuthSettingDomainService）
	// if !s.authSettingService.IsFeatureEnabled(AuthFeatureKeyNormalLogin) {
	//     return "", exception.NewBusinessException("普通登录已禁用")
	// }

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
	// TODO: 检查用户注册是否启用（需要 AuthSettingDomainService）
	// if !s.authSettingService.IsFeatureEnabled(AuthFeatureKeyUserRegister) {
	//     return exception.NewBusinessException("用户注册已禁用")
	// }

	// TODO: 如果是邮箱注册，需要验证码（需要 VerificationCodeService）
	// if email != "" && phone == "" {
	//     if code == "" {
	//         return exception.NewBusinessException("邮箱注册需要验证码")
	//     }
	//     if !s.verificationService.VerifyCode(email, code) {
	//         return exception.NewBusinessException("验证码无效或已过期")
	//     }
	// }

	_, err := s.domainService.Register(email, phone, password)
	return err
}

// SendEmailVerificationCode 发送注册邮箱验证码（对应 Java 的 sendEmailVerificationCode）
func (s *LoginAppService) SendEmailVerificationCode(email, captchaUuid, captchaCode, ip string) error {
	// TODO: 实现邮件验证码发送（需要 EmailService + VerificationCodeService）
	return exception.NewBusinessException("邮件验证码功能尚未实现")
}

// SendResetPasswordCode 发送重置密码邮箱验证码（对应 Java 的 sendResetPasswordCode）
func (s *LoginAppService) SendResetPasswordCode(email, captchaUuid, captchaCode, ip string) error {
	// TODO: 实现重置密码验证码发送
	return exception.NewBusinessException("重置密码验证码功能尚未实现")
}

// VerifyEmailCode 验证邮箱验证码（对应 Java 的 verifyEmailCode）
func (s *LoginAppService) VerifyEmailCode(email, code string) bool {
	// TODO: 实现验证码验证
	return false
}

// VerifyResetPasswordCode 验证重置密码邮箱验证码
func (s *LoginAppService) VerifyResetPasswordCode(email, code string) bool {
	// TODO: 实现验证码验证
	return false
}

// ResetPassword 重置密码（对应 Java 的 resetPassword）
func (s *LoginAppService) ResetPassword(email, newPassword, code string) error {
	// TODO: 验证重置密码验证码
	// if !s.VerifyResetPasswordCode(email, code) {
	//     return exception.NewBusinessException("验证码无效或已过期")
	// }

	user, err := s.domainService.FindUserByAccount(email)
	if err != nil {
		return err
	}
	if user == nil {
		return exception.NewBusinessException("用户不存在")
	}

	return s.domainService.UpdatePassword(user.ID, newPassword)
}
