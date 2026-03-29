package llm

import (
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// ProviderEntity 服务商实体（对应 Java 的 ProviderEntity）
type ProviderEntity struct {
	ID          string          `gorm:"column:id;primaryKey" json:"id"`
	UserID      string          `gorm:"column:user_id" json:"userId"`
	Protocol    string          `gorm:"column:protocol" json:"protocol"`
	Name        string          `gorm:"column:name" json:"name"`
	Description string          `gorm:"column:description" json:"description"`
	Config      *ProviderConfig `gorm:"column:config;serializer:json" json:"config"`
	IsOfficial  *bool           `gorm:"column:is_official" json:"isOfficial"`
	Status      *bool           `gorm:"column:status" json:"status"`

	entity.BaseEntity
}

func (ProviderEntity) TableName() string {
	return "providers"
}

// IsActive 检查服务商是否激活
func (e *ProviderEntity) IsActive() error {
	if e.Status == nil || !*e.Status {
		return exception.NewBusinessException("服务商未激活")
	}
	return nil
}

// GetStatus 获取状态
func (e *ProviderEntity) GetStatus() bool {
	return e.Status != nil && *e.Status
}

// GetIsOfficial 获取是否官方
func (e *ProviderEntity) GetIsOfficial() bool {
	return e.IsOfficial != nil && *e.IsOfficial
}

// SetAdmin 设置为管理员操作（不检查userId）
func (e *ProviderEntity) SetAdmin() {
	e.UserID = ""
}

// NeedCheckUserId 是否需要检查userId
func (e *ProviderEntity) NeedCheckUserId() bool {
	return e.UserID != ""
}

// ModelEntity 模型实体（对应 Java 的 ModelEntity）
type ModelEntity struct {
	ID            string `gorm:"column:id;primaryKey" json:"id"`
	UserID        string `gorm:"column:user_id" json:"userId"`
	ProviderID    string `gorm:"column:provider_id" json:"providerId"`
	ModelID       string `gorm:"column:model_id" json:"modelId"`
	Name          string `gorm:"column:name" json:"name"`
	Description   string `gorm:"column:description" json:"description"`
	ModelEndpoint string `gorm:"column:model_endpoint" json:"modelEndpoint"`
	IsOfficial    *bool  `gorm:"column:is_official" json:"isOfficial"`
	Type          string `gorm:"column:type" json:"type"`
	Status        *bool  `gorm:"column:status" json:"status"`

	entity.BaseEntity
}

func (ModelEntity) TableName() string {
	return "models"
}

// GetStatus 获取状态
func (e *ModelEntity) GetStatus() bool {
	return e.Status != nil && *e.Status
}

// GetIsOfficial 获取是否官方
func (e *ModelEntity) GetIsOfficial() bool {
	return e.IsOfficial != nil && *e.IsOfficial
}

// SetAdmin 设置为管理员操作
func (e *ModelEntity) SetAdmin() {
	e.UserID = ""
}

// IsChatType 是否是对话类型
func (e *ModelEntity) IsChatType() bool {
	return e.Type == string(ModelTypeChat)
}

// ProviderAggregate 服务商聚合根（对应 Java 的 ProviderAggregate）
type ProviderAggregate struct {
	Entity *ProviderEntity
	Models []ModelEntity
}

// NewProviderAggregate 创建服务商聚合根
func NewProviderAggregate(entity *ProviderEntity, models []ModelEntity) *ProviderAggregate {
	if models == nil {
		models = []ModelEntity{}
	}
	return &ProviderAggregate{
		Entity: entity,
		Models: models,
	}
}

// GetStatus 获取服务商状态
func (a *ProviderAggregate) GetStatus() bool {
	return a.Entity.GetStatus()
}

// GetIsOfficial 获取是否官方
func (a *ProviderAggregate) GetIsOfficial() bool {
	return a.Entity.GetIsOfficial()
}

// GetID 获取服务商ID
func (a *ProviderAggregate) GetID() string {
	return a.Entity.ID
}

// GetName 获取服务商名称
func (a *ProviderAggregate) GetName() string {
	return a.Entity.Name
}
