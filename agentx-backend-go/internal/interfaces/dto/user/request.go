package dto

// LoginRequest 登录请求（对应 Java 的 LoginRequest）
type LoginRequest struct {
	Account  string `json:"account" binding:"required"`  // 邮箱或手机号
	Password string `json:"password" binding:"required"` // 密码
}

// RegisterRequest 注册请求（对应 Java 的 RegisterRequest）
type RegisterRequest struct {
	Email    string `json:"email" binding:"omitempty,email"` // 邮箱
	Phone    string `json:"phone"`                           // 手机号
	Password string `json:"password" binding:"required"`     // 密码
	Code     string `json:"code"`                            // 验证码（邮箱注册时必填）
}

// UserUpdateRequest 修改用户信息请求（对应 Java 的 UserUpdateRequest）
type UserUpdateRequest struct {
	Nickname string `json:"nickname" binding:"required"` // 昵称
}

// ChangePasswordRequest 修改密码请求（对应 Java 的 ChangePasswordRequest）
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`          // 当前密码
	NewPassword     string `json:"newPassword" binding:"required,min=6,max=20"` // 新密码
	ConfirmPassword string `json:"confirmPassword" binding:"required"`          // 确认密码
}

// QueryUserRequest 查询用户请求（对应 Java 的 QueryUserRequest）
type QueryUserRequest struct {
	Keyword  string `form:"keyword" json:"keyword"`   // 搜索关键词
	Page     int    `form:"page" json:"page"`         // 页码
	PageSize int    `form:"pageSize" json:"pageSize"` // 每页大小
}

// GetPage 获取页码（默认1）
func (q *QueryUserRequest) GetPage() int {
	if q.Page <= 0 {
		return 1
	}
	return q.Page
}

// GetPageSize 获取每页大小（默认15）
func (q *QueryUserRequest) GetPageSize() int {
	if q.PageSize <= 0 {
		return 15
	}
	if q.PageSize > 100 {
		return 100
	}
	return q.PageSize
}

// SendEmailCodeRequest 发送邮箱验证码请求（对应 Java 的 SendEmailCodeRequest）
type SendEmailCodeRequest struct {
	Email       string `json:"email" binding:"required,email"` // 邮箱
	CaptchaUuid string `json:"captchaUuid" binding:"required"` // 验证码UUID
	CaptchaCode string `json:"captchaCode" binding:"required"` // 图形验证码
}

// SendResetPasswordCodeRequest 发送重置密码验证码请求
type SendResetPasswordCodeRequest struct {
	Email       string `json:"email" binding:"required,email"` // 邮箱
	CaptchaUuid string `json:"captchaUuid" binding:"required"` // 验证码UUID
	CaptchaCode string `json:"captchaCode" binding:"required"` // 图形验证码
}

// VerifyEmailCodeRequest 验证邮箱验证码请求
type VerifyEmailCodeRequest struct {
	Email string `json:"email" binding:"required,email"` // 邮箱
	Code  string `json:"code" binding:"required"`        // 验证码
}

// VerifyResetPasswordCodeRequest 验证重置密码验证码请求
type VerifyResetPasswordCodeRequest struct {
	Email string `json:"email" binding:"required,email"` // 邮箱
	Code  string `json:"code" binding:"required"`        // 验证码
}

// ResetPasswordRequest 重置密码请求（对应 Java 的 ResetPasswordRequest）
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`              // 邮箱
	NewPassword string `json:"newPassword" binding:"required,min=6,max=20"` // 新密码
	Code        string `json:"code" binding:"required"`                     // 验证码
}

// GetCaptchaRequest 获取图形验证码请求
type GetCaptchaRequest struct {
	// 暂无字段
}

// UserSettingsUpdateRequest 更新用户设置请求
type UserSettingsUpdateRequest struct {
	DefaultModel          string `json:"defaultModel"`
	DefaultOcrModel       string `json:"defaultOcrModel"`
	DefaultEmbeddingModel string `json:"defaultEmbeddingModel"`
}
