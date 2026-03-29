package auth

import (
	"strings"

	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/auth"
	domainSso "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/sso"
)

// AppService 认证配置应用服务（对应 Java 的 AuthSettingAppService）
type AppService struct {
	domainService *domain.DomainService
}

// NewAppService 创建认证配置应用服务
func NewAppService(domainService *domain.DomainService) *AppService {
	return &AppService{domainService: domainService}
}

// GetAuthConfig 获取前端认证配置（对应 Java 的 getAuthConfig）
func (s *AppService) GetAuthConfig() (*AuthConfigDTO, error) {
	// 获取启用的登录方式
	loginSettings, err := s.domainService.GetEnabledFeatures(domain.FeatureTypeLogin)
	if err != nil {
		return nil, err
	}

	loginMethods := make(map[string]*LoginMethodDTO)
	for _, setting := range loginSettings {
		method := &LoginMethodDTO{
			Enabled: setting.IsEnabled(),
			Name:    setting.FeatureName,
		}

		// 根据功能键设置provider
		providerCode := getProviderCodeByFeatureKey(setting.FeatureKey)
		if providerCode != "" {
			method.Provider = providerCode
		}

		loginMethods[setting.FeatureKey] = method
	}

	// 检查注册是否启用
	registerEnabled, err := s.domainService.IsFeatureEnabled(domain.AuthFeatureKeyUserRegister)
	if err != nil {
		return nil, err
	}

	return &AuthConfigDTO{
		LoginMethods:    loginMethods,
		RegisterEnabled: registerEnabled,
	}, nil
}

// GetAllAuthSettings 获取所有认证配置
func (s *AppService) GetAllAuthSettings() ([]*AuthSettingDTO, error) {
	entities, err := s.domainService.GetAllAuthSettings()
	if err != nil {
		return nil, err
	}
	return ToDTOList(entities), nil
}

// GetAuthSettingByID 根据ID获取认证配置
func (s *AppService) GetAuthSettingByID(id string) (*AuthSettingDTO, error) {
	entity, err := s.domainService.GetByID(id)
	if err != nil {
		return nil, err
	}
	return ToDTO(entity), nil
}

// ToggleAuthSetting 切换认证配置启用状态
func (s *AppService) ToggleAuthSetting(id string) (*AuthSettingDTO, error) {
	entity, err := s.domainService.ToggleEnabled(id)
	if err != nil {
		return nil, err
	}
	return ToDTO(entity), nil
}

// UpdateAuthSetting 更新认证配置
func (s *AppService) UpdateAuthSetting(id string, req *UpdateAuthSettingRequest) (*AuthSettingDTO, error) {
	entity, err := s.domainService.GetByID(id)
	if err != nil {
		return nil, err
	}

	updatedEntity := UpdateEntityFromRequest(entity, req)
	savedEntity, err := s.domainService.UpdateAuthSetting(updatedEntity)
	if err != nil {
		return nil, err
	}
	return ToDTO(savedEntity), nil
}

// DeleteAuthSetting 删除认证配置
func (s *AppService) DeleteAuthSetting(id string) error {
	return s.domainService.DeleteAuthSetting(id)
}

// getProviderCodeByFeatureKey 根据认证功能键获取对应的SSO提供商代码
func getProviderCodeByFeatureKey(featureKey string) string {
	switch domain.AuthFeatureKey(featureKey) {
	case domain.AuthFeatureKeyGithubLogin:
		return strings.ToUpper(string(domainSso.SsoProviderGitHub))
	case domain.AuthFeatureKeyCommunityLogin:
		return strings.ToUpper(string(domainSso.SsoProviderCommunity))
	default:
		return ""
	}
}
