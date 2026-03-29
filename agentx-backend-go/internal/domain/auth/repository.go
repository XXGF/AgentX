package auth

import (
	"gorm.io/gorm"
)

// AuthSettingRepository 认证配置仓储接口（对应 Java 的 AuthSettingRepository）
type AuthSettingRepository interface {
	FindByID(id string) (*AuthSettingEntity, error)
	FindByFeatureKey(featureKey string) (*AuthSettingEntity, error)
	FindByFeatureType(featureType string) ([]AuthSettingEntity, error)
	FindEnabledByFeatureType(featureType string) ([]AuthSettingEntity, error)
	FindAll() ([]AuthSettingEntity, error)
	CountByFeatureKeyAndEnabled(featureKey string) (int64, error)
	CountByFeatureKey(featureKey string) (int64, error)
	Create(entity *AuthSettingEntity) error
	Update(entity *AuthSettingEntity) error
	UpdateEnabled(id string, enabled bool) error
	Delete(id string) error
}

// ---- GORM 实现 ----

type AuthSettingRepositoryImpl struct {
	DB *gorm.DB
}

func NewAuthSettingRepository(db *gorm.DB) AuthSettingRepository {
	return &AuthSettingRepositoryImpl{DB: db}
}

func (r *AuthSettingRepositoryImpl) FindByID(id string) (*AuthSettingEntity, error) {
	var entity AuthSettingEntity
	result := r.DB.First(&entity, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *AuthSettingRepositoryImpl) FindByFeatureKey(featureKey string) (*AuthSettingEntity, error) {
	var entity AuthSettingEntity
	result := r.DB.Where("feature_key = ?", featureKey).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *AuthSettingRepositoryImpl) FindByFeatureType(featureType string) ([]AuthSettingEntity, error) {
	var entities []AuthSettingEntity
	result := r.DB.Where("feature_type = ?", featureType).Order("display_order ASC").Find(&entities)
	return entities, result.Error
}

func (r *AuthSettingRepositoryImpl) FindEnabledByFeatureType(featureType string) ([]AuthSettingEntity, error) {
	var entities []AuthSettingEntity
	result := r.DB.Where("feature_type = ? AND enabled = ?", featureType, true).
		Order("display_order ASC").Find(&entities)
	return entities, result.Error
}

func (r *AuthSettingRepositoryImpl) FindAll() ([]AuthSettingEntity, error) {
	var entities []AuthSettingEntity
	result := r.DB.Order("feature_type ASC, display_order ASC").Find(&entities)
	return entities, result.Error
}

func (r *AuthSettingRepositoryImpl) CountByFeatureKeyAndEnabled(featureKey string) (int64, error) {
	var count int64
	result := r.DB.Model(&AuthSettingEntity{}).
		Where("feature_key = ? AND enabled = ?", featureKey, true).Count(&count)
	return count, result.Error
}

func (r *AuthSettingRepositoryImpl) CountByFeatureKey(featureKey string) (int64, error) {
	var count int64
	result := r.DB.Model(&AuthSettingEntity{}).Where("feature_key = ?", featureKey).Count(&count)
	return count, result.Error
}

func (r *AuthSettingRepositoryImpl) Create(entity *AuthSettingEntity) error {
	return r.DB.Create(entity).Error
}

func (r *AuthSettingRepositoryImpl) Update(entity *AuthSettingEntity) error {
	return r.DB.Save(entity).Error
}

func (r *AuthSettingRepositoryImpl) UpdateEnabled(id string, enabled bool) error {
	return r.DB.Model(&AuthSettingEntity{}).Where("id = ?", id).Update("enabled", enabled).Error
}

func (r *AuthSettingRepositoryImpl) Delete(id string) error {
	return r.DB.Delete(&AuthSettingEntity{}, "id = ?", id).Error
}
