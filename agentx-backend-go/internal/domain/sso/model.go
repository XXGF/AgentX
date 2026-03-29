package sso

// SsoProvider SSO提供商枚举（对应 Java 的 SsoProvider）
type SsoProvider string

const (
	SsoProviderCommunity SsoProvider = "community"
	SsoProviderGitHub    SsoProvider = "github"
	SsoProviderGoogle    SsoProvider = "google"
	SsoProviderWechat    SsoProvider = "wechat"
)

// GetName 获取提供商名称
func (p SsoProvider) GetName() string {
	switch p {
	case SsoProviderCommunity:
		return "敲鸭"
	case SsoProviderGitHub:
		return "GitHub"
	case SsoProviderGoogle:
		return "Google"
	case SsoProviderWechat:
		return "微信"
	default:
		return string(p)
	}
}

// FromCode 根据code获取SsoProvider
func SsoProviderFromCode(code string) (SsoProvider, bool) {
	switch code {
	case "community":
		return SsoProviderCommunity, true
	case "github":
		return SsoProviderGitHub, true
	case "google":
		return SsoProviderGoogle, true
	case "wechat":
		return SsoProviderWechat, true
	default:
		return "", false
	}
}

// SsoUserInfo SSO用户信息（对应 Java 的 SsoUserInfo）
type SsoUserInfo struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Avatar   string      `json:"avatar"`
	Desc     string      `json:"desc"`
	Provider SsoProvider `json:"provider"`
}

// SsoService SSO服务接口（对应 Java 的 SsoService）
type SsoService interface {
	// GetLoginUrl 获取SSO登录重定向URL
	GetLoginUrl(redirectUrl string) (string, error)
	// GetUserInfo 通过授权码获取用户信息
	GetUserInfo(authCode string) (*SsoUserInfo, error)
	// GetProvider 获取支持的SSO提供商类型
	GetProvider() SsoProvider
}
