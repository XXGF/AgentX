package auth

// FeatureType 功能类型枚举（对应 Java 的 FeatureType）
type FeatureType string

const (
	FeatureTypeLogin    FeatureType = "LOGIN"    // 登录功能
	FeatureTypeRegister FeatureType = "REGISTER" // 注册功能
)

// AuthFeatureKey 认证功能键枚举（对应 Java 的 AuthFeatureKey）
type AuthFeatureKey string

const (
	AuthFeatureKeyNormalLogin    AuthFeatureKey = "NORMAL_LOGIN"    // 普通登录
	AuthFeatureKeyGithubLogin   AuthFeatureKey = "GITHUB_LOGIN"    // GitHub登录
	AuthFeatureKeyCommunityLogin AuthFeatureKey = "COMMUNITY_LOGIN" // 敲鸭登录
	AuthFeatureKeyUserRegister  AuthFeatureKey = "USER_REGISTER"   // 用户注册
)

// GetName 获取功能键名称
func (k AuthFeatureKey) GetName() string {
	switch k {
	case AuthFeatureKeyNormalLogin:
		return "普通登录"
	case AuthFeatureKeyGithubLogin:
		return "GitHub登录"
	case AuthFeatureKeyCommunityLogin:
		return "敲鸭登录"
	case AuthFeatureKeyUserRegister:
		return "用户注册"
	default:
		return string(k)
	}
}
