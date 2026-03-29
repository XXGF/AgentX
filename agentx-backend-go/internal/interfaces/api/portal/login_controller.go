package portal

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	appUser "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/user"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
	dto "github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/dto/user"
)

// LoginController 登录注册控制器（对应 Java 的 LoginController）
type LoginController struct {
	loginAppService *appUser.LoginAppService
}

// NewLoginController 创建登录控制器
func NewLoginController(loginAppService *appUser.LoginAppService) *LoginController {
	return &LoginController{loginAppService: loginAppService}
}

// Login 登录（POST /login）
func (ctrl *LoginController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}

	token, err := ctrl.loginAppService.Login(req.Account, req.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessWithMessage("登录成功", gin.H{"token": token}))
}

// Register 注册（POST /register）
func (ctrl *LoginController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}

	if err := ctrl.loginAppService.Register(req.Email, req.Phone, req.Password, req.Code); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.Success().SetMessage("注册成功"))
}

// GetCaptcha 获取图形验证码（POST /get-captcha）
func (ctrl *LoginController) GetCaptcha(c *gin.Context) {
	// TODO: 实现图形验证码生成（需要 CaptchaUtils）
	c.JSON(http.StatusOK, common.SuccessWithData(gin.H{
		"uuid":        "",
		"imageBase64": "",
	}))
}

// SendEmailCode 发送邮箱验证码（POST /send-email-code）
func (ctrl *LoginController) SendEmailCode(c *gin.Context) {
	var req dto.SendEmailCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}

	clientIP := getClientIP(c)
	if err := ctrl.loginAppService.SendEmailVerificationCode(req.Email, req.CaptchaUuid, req.CaptchaCode, clientIP); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.Success().SetMessage("验证码已发送，请查收邮件"))
}

// SendResetPasswordCode 发送重置密码验证码（POST /send-reset-password-code）
func (ctrl *LoginController) SendResetPasswordCode(c *gin.Context) {
	var req dto.SendResetPasswordCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}

	clientIP := getClientIP(c)
	if err := ctrl.loginAppService.SendResetPasswordCode(req.Email, req.CaptchaUuid, req.CaptchaCode, clientIP); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.Success().SetMessage("验证码已发送，请查收邮件"))
}

// VerifyEmailCode 验证邮箱验证码（POST /verify-email-code）
func (ctrl *LoginController) VerifyEmailCode(c *gin.Context) {
	var req dto.VerifyEmailCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}

	isValid := ctrl.loginAppService.VerifyEmailCode(req.Email, req.Code)
	if isValid {
		c.JSON(http.StatusOK, common.SuccessWithMessage("验证码验证成功", true))
	} else {
		c.JSON(http.StatusForbidden, common.Error(403, "验证码无效或已过期"))
	}
}

// ResetPassword 重置密码（POST /reset-password）
func (ctrl *LoginController) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}

	if err := ctrl.loginAppService.ResetPassword(req.Email, req.NewPassword, req.Code); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.Success().SetMessage("密码重置成功"))
}

// getClientIP 获取客户端IP（对应 Java 的 getClientIp）
func getClientIP(c *gin.Context) string {
	// 检查 X-Forwarded-For
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" && ip != "unknown" {
		// 多个代理时取第一个
		if idx := strings.Index(ip, ","); idx != -1 {
			ip = strings.TrimSpace(ip[:idx])
		}
		return ip
	}

	// 检查 X-Real-IP
	ip = c.GetHeader("X-Real-IP")
	if ip != "" && ip != "unknown" {
		return ip
	}

	// 检查 Proxy-Client-IP
	ip = c.GetHeader("Proxy-Client-IP")
	if ip != "" && ip != "unknown" {
		return ip
	}

	// 使用 RemoteAddr
	return c.ClientIP()
}
