package auth

import (
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// DomainService 认证配置领域服务（对应 Java 的 AuthSettingDomainService）
type DomainService struct {
	repo AuthSettingRepository
}

// NewDomainService 创建认证配置领域服务
func NewDomainService(repo AuthSettingRepository) *DomainService {
	return &DomainService{repo: repo}
}

// GetEnabledFeatures 获取指定类型的启用功能列表
func (s *DomainService) GetEnabledFeatures(featureType FeatureType) ([]AuthSettingEntity, error) {
	return s.repo.FindEnabledByFeatureType(string(featureType))
}

// GetAllFeatures 获取指定类型的所有功能列表
func (s *DomainService) GetAllFeatures(featureType FeatureType) ([]AuthSettingEntity, error) {
	return s.repo.FindByFeatureType(string(featureType))
}

// GetAllAuthSettings 获取所有认证配置
func (s *DomainService) GetAllAuthSettings() ([]AuthSettingEntity, error) {
	return s.repo.FindAll()
}

// IsFeatureEnabled 检查指定功能是否启用
func (s *DomainService) IsFeatureEnabled(featureKey AuthFeatureKey) (bool, error) {
	count, err := s.repo.CountByFeatureKeyAndEnabled(string(featureKey))
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetByFeatureKey 根据功能键获取认证配置
func (s *DomainService) GetByFeatureKey(featureKey AuthFeatureKey) (*AuthSettingEntity, error) {
	return s.repo.FindByFeatureKey(string(featureKey))
}

// GetByID 根据ID获取认证配置
func (s *DomainService) GetByID(id string) (*AuthSettingEntity, error) {
	entity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, exception.NewBusinessException("认证配置不存在")
	}
	return entity, nil
}

// ToggleEnabled 切换功能启用状态
func (s *DomainService) ToggleEnabled(id string) (*AuthSettingEntity, error) {
	entity, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	entity.ToggleEnabled()
	if err := s.repo.UpdateEnabled(id, entity.IsEnabled()); err != nil {
		return nil, err
	}

	return entity, nil
}

// UpdateAuthSetting 更新认证配置
func (s *DomainService) UpdateAuthSetting(entity *AuthSettingEntity) (*AuthSettingEntity, error) {
	// 确认存在
	if _, err := s.GetByID(entity.ID); err != nil {
		return nil, err
	}

	if err := s.repo.Update(entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// CreateAuthSetting 创建认证配置
func (s *DomainService) CreateAuthSetting(entity *AuthSettingEntity) (*AuthSettingEntity, error) {
	// 检查功能键是否已存在
	count, err := s.repo.CountByFeatureKey(entity.FeatureKey)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, exception.NewBusinessException("功能键已存在: " + entity.FeatureKey)
	}

	if err := s.repo.Create(entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// DeleteAuthSetting 删除认证配置
func (s *DomainService) DeleteAuthSetting(id string) error {
	// 确认存在
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
