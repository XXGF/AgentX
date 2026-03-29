package llm

import (
	"time"

	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/llm"
)

// ProviderDTO 服务商DTO（对应 Java 的 ProviderDTO）
type ProviderDTO struct {
	ID          string                `json:"id"`
	Protocol    string                `json:"protocol"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Config      *domain.ProviderConfig `json:"config,omitempty"`
	IsOfficial  *bool                 `json:"isOfficial"`
	Status      *bool                 `json:"status"`
	CreatedAt   time.Time             `json:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt"`
	Models      []*ModelDTO           `json:"models"`
}

// MaskSensitiveInfo 脱敏处理
func (dto *ProviderDTO) MaskSensitiveInfo() {
	if dto.Config != nil {
		dto.Config.MaskSensitiveInfo()
	}
}

// ModelDTO 模型DTO（对应 Java 的 ModelDTO）
type ModelDTO struct {
	ID            string    `json:"id"`
	UserID        string    `json:"userId"`
	ProviderID    string    `json:"providerId"`
	ProviderName  string    `json:"providerName,omitempty"`
	ModelID       string    `json:"modelId"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Type          string    `json:"type"`
	ModelEndpoint string    `json:"modelEndpoint"`
	IsOfficial    *bool     `json:"isOfficial"`
	Status        *bool     `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// ProviderCreateRequest 创建服务商请求（对应 Java 的 ProviderCreateRequest）
type ProviderCreateRequest struct {
	Protocol    string                `json:"protocol" binding:"required"`
	Name        string                `json:"name" binding:"required"`
	Description string                `json:"description"`
	Config      *domain.ProviderConfig `json:"config"`
	Status      *bool                 `json:"status"`
}

// ProviderUpdateRequest 更新服务商请求（对应 Java 的 ProviderUpdateRequest）
type ProviderUpdateRequest struct {
	ID          string                `json:"id"`
	Protocol    string                `json:"protocol" binding:"required"`
	Name        string                `json:"name" binding:"required"`
	Description string                `json:"description"`
	Config      *domain.ProviderConfig `json:"config"`
	Status      *bool                 `json:"status"`
}

// ModelCreateRequest 创建模型请求（对应 Java 的 ModelCreateRequest）
type ModelCreateRequest struct {
	ProviderID    string `json:"providerId"`
	ModelID       string `json:"modelId" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	Type          string `json:"type"`
	ModelEndpoint string `json:"modelEndpoint"`
}

// ModelUpdateRequest 更新模型请求（对应 Java 的 ModelUpdateRequest）
type ModelUpdateRequest struct {
	ID            string `json:"id"`
	ModelID       string `json:"modelId" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	ModelEndpoint string `json:"modelEndpoint"`
}
