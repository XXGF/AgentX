package auth

import "time"

// AuthSettingDTO 认证配置DTO（对应 Java 的 AuthSettingDTO）
type AuthSettingDTO struct {
	ID           string                 `json:"id"`
	FeatureType  string                 `json:"featureType"`
	FeatureKey   string                 `json:"featureKey"`
	FeatureName  string                 `json:"featureName"`
	Enabled      *bool                  `json:"enabled"`
	ConfigData   map[string]interface{} `json:"configData,omitempty"`
	DisplayOrder *int                   `json:"displayOrder"`
	Description  string                 `json:"description"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

// AuthConfigDTO 认证配置响应DTO（对应 Java 的 AuthConfigDTO）
type AuthConfigDTO struct {
	LoginMethods    map[string]*LoginMethodDTO `json:"loginMethods"`
	RegisterEnabled bool                       `json:"registerEnabled"`
}

// LoginMethodDTO 登录方式DTO（对应 Java 的 LoginMethodDTO）
type LoginMethodDTO struct {
	Enabled  bool   `json:"enabled"`
	Name     string `json:"name"`
	Provider string `json:"provider,omitempty"`
}

// UpdateAuthSettingRequest 更新认证配置请求（对应 Java 的 UpdateAuthSettingRequest）
type UpdateAuthSettingRequest struct {
	FeatureName  *string                `json:"featureName,omitempty"`
	Enabled      *bool                  `json:"enabled,omitempty"`
	ConfigData   map[string]interface{} `json:"configData,omitempty"`
	DisplayOrder *int                   `json:"displayOrder,omitempty"`
	Description  *string                `json:"description,omitempty"`
}
