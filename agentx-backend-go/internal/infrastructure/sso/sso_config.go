package sso

import (
	"fmt"

	domainAuth "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/auth"
	domainSso "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/sso"
)

// GitHubSsoConfig GitHub SSO配置（对应 Java 的 SsoConfigProvider.GitHubSsoConfig）
type GitHubSsoConfig struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RedirectURI  string `json:"redirectUri"`
	AuthorizeURL string `json:"authorizeUrl"`
	TokenURL     string `json:"tokenUrl"`
	UserInfoURL  string `json:"userInfoUrl"`
	UserEmailURL string `json:"userEmailUrl"`
}

// CommunitySsoConfig Community SSO配置（对应 Java 的 SsoConfigProvider.CommunitySsoConfig）
type CommunitySsoConfig struct {
	BaseURL     string `json:"baseUrl"`
	AppKey      string `json:"appKey"`
	AppSecret   string `json:"appSecret"`
	CallbackURL string `json:"callbackUrl"`
}

// ConfigProvider SSO配置提供者（对应 Java 的 SsoConfigProvider）
type ConfigProvider struct {
	authDomainService *domainAuth.DomainService
}

// NewConfigProvider 创建SSO配置提供者
func NewConfigProvider(authDomainService *domainAuth.DomainService) *ConfigProvider {
	return &ConfigProvider{authDomainService: authDomainService}
}

// GetGitHubConfig 获取GitHub SSO配置
func (p *ConfigProvider) GetGitHubConfig() *GitHubSsoConfig {
	entity, err := p.authDomainService.GetByFeatureKey(domainAuth.AuthFeatureKeyGithubLogin)
	if err != nil || entity == nil || entity.ConfigData == nil {
		return &GitHubSsoConfig{
			AuthorizeURL: "https://github.com/login/oauth/authorize",
			TokenURL:     "https://github.com/login/oauth/access_token",
			UserInfoURL:  "https://api.github.com/user",
			UserEmailURL: "https://api.github.com/user/emails",
		}
	}

	config := &GitHubSsoConfig{
		AuthorizeURL: "https://github.com/login/oauth/authorize",
		TokenURL:     "https://github.com/login/oauth/access_token",
		UserInfoURL:  "https://api.github.com/user",
		UserEmailURL: "https://api.github.com/user/emails",
	}

	if v, ok := entity.ConfigData["clientId"].(string); ok {
		config.ClientID = v
	}
	if v, ok := entity.ConfigData["clientSecret"].(string); ok {
		config.ClientSecret = v
	}
	if v, ok := entity.ConfigData["redirectUri"].(string); ok {
		config.RedirectURI = v
	}

	return config
}

// GetCommunityConfig 获取Community SSO配置
func (p *ConfigProvider) GetCommunityConfig() *CommunitySsoConfig {
	entity, err := p.authDomainService.GetByFeatureKey(domainAuth.AuthFeatureKeyCommunityLogin)
	if err != nil || entity == nil || entity.ConfigData == nil {
		return &CommunitySsoConfig{}
	}

	config := &CommunitySsoConfig{}
	if v, ok := entity.ConfigData["baseUrl"].(string); ok {
		config.BaseURL = v
	}
	if v, ok := entity.ConfigData["appKey"].(string); ok {
		config.AppKey = v
	}
	if v, ok := entity.ConfigData["appSecret"].(string); ok {
		config.AppSecret = v
	}
	if v, ok := entity.ConfigData["callbackUrl"].(string); ok {
		config.CallbackURL = v
	}

	return config
}

// ServiceFactory SSO服务工厂（对应 Java 的 SsoServiceFactory）
type ServiceFactory struct {
	services          map[domainSso.SsoProvider]domainSso.SsoService
	authDomainService *domainAuth.DomainService
}

// NewServiceFactory 创建SSO服务工厂
func NewServiceFactory(authDomainService *domainAuth.DomainService, services ...domainSso.SsoService) *ServiceFactory {
	serviceMap := make(map[domainSso.SsoProvider]domainSso.SsoService)
	for _, svc := range services {
		serviceMap[svc.GetProvider()] = svc
	}
	return &ServiceFactory{
		services:          serviceMap,
		authDomainService: authDomainService,
	}
}

// GetSsoService 根据提供商获取SSO服务
func (f *ServiceFactory) GetSsoService(providerCode string) (domainSso.SsoService, error) {
	provider, ok := domainSso.SsoProviderFromCode(providerCode)
	if !ok {
		return nil, fmt.Errorf("不支持的SSO提供商: %s", providerCode)
	}

	svc, ok := f.services[provider]
	if !ok {
		return nil, fmt.Errorf("不支持的SSO提供商: %s", provider.GetName())
	}

	// 检查SSO提供商是否启用
	featureKey := getAuthFeatureKeyByProvider(provider)
	if featureKey != "" {
		enabled, err := f.authDomainService.IsFeatureEnabled(domainAuth.AuthFeatureKey(featureKey))
		if err != nil {
			return nil, err
		}
		if !enabled {
			return nil, fmt.Errorf("SSO提供商已禁用: %s", provider.GetName())
		}
	}

	return svc, nil
}

// getAuthFeatureKeyByProvider 根据SSO提供商获取对应的认证功能键
func getAuthFeatureKeyByProvider(provider domainSso.SsoProvider) domainAuth.AuthFeatureKey {
	switch provider {
	case domainSso.SsoProviderGitHub:
		return domainAuth.AuthFeatureKeyGithubLogin
	case domainSso.SsoProviderCommunity:
		return domainAuth.AuthFeatureKeyCommunityLogin
	default:
		return ""
	}
}
