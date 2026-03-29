package auth

import (
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
)

// AuthSettingEntity 认证配置实体（对应 Java 的 AuthSettingEntity）
type AuthSettingEntity struct {
	ID           string                 `gorm:"column:id;primaryKey" json:"id"`
	FeatureType  string                 `gorm:"column:feature_type" json:"featureType"`
	FeatureKey   string                 `gorm:"column:feature_key" json:"featureKey"`
	FeatureName  string                 `gorm:"column:feature_name" json:"featureName"`
	Enabled      *bool                  `gorm:"column:enabled" json:"enabled"`
	ConfigData   map[string]interface{} `gorm:"column:config_data;serializer:json" json:"configData"`
	DisplayOrder *int                   `gorm:"column:display_order" json:"displayOrder"`
	Description  string                 `gorm:"column:description" json:"description"`

	entity.BaseEntity
}

// TableName 指定表名
func (AuthSettingEntity) TableName() string {
	return "auth_settings"
}

// IsEnabled 判断是否启用
func (e *AuthSettingEntity) IsEnabled() bool {
	return e.Enabled != nil && *e.Enabled
}

// SetEnabled 设置启用状态
func (e *AuthSettingEntity) SetEnabled(enabled bool) {
	e.Enabled = &enabled
}

// ToggleEnabled 切换启用状态
func (e *AuthSettingEntity) ToggleEnabled() {
	if e.Enabled == nil {
		enabled := true
		e.Enabled = &enabled
	} else {
		toggled := !*e.Enabled
		e.Enabled = &toggled
	}
}
