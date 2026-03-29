package agent

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// ToolPresetParams 预先设置的工具参数类型
type ToolPresetParams map[string]map[string]map[string]string

// Value 实现 driver.Valuer 接口
func (t ToolPresetParams) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	b, err := json.Marshal(t)
	return string(b), err
}

// Scan 实现 sql.Scanner 接口
func (t *ToolPresetParams) Scan(value interface{}) error {
	if value == nil {
		*t = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return nil
	}
	return json.Unmarshal(bytes, t)
}

// StringList JSON字符串列表类型
type StringList []string

// Value 实现 driver.Valuer 接口
func (s StringList) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	return string(b), err
}

// Scan 实现 sql.Scanner 接口
func (s *StringList) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		*s = []string{}
		return nil
	}
	return json.Unmarshal(bytes, s)
}

// AgentEntity Agent实体（对应 Java 的 AgentEntity）
type AgentEntity struct {
	ID               string           `gorm:"column:id;primaryKey" json:"id"`
	Name             string           `gorm:"column:name" json:"name"`
	Avatar           string           `gorm:"column:avatar" json:"avatar"`
	Description      string           `gorm:"column:description" json:"description"`
	SystemPrompt     string           `gorm:"column:system_prompt" json:"systemPrompt"`
	WelcomeMessage   string           `gorm:"column:welcome_message" json:"welcomeMessage"`
	ToolIDs          StringList       `gorm:"column:tool_ids;type:jsonb" json:"toolIds"`
	KnowledgeBaseIDs StringList       `gorm:"column:knowledge_base_ids;type:jsonb" json:"knowledgeBaseIds"`
	PublishedVersion string           `gorm:"column:published_version" json:"publishedVersion"`
	Enabled          *bool            `gorm:"column:enabled" json:"enabled"`
	UserID           string           `gorm:"column:user_id" json:"userId"`
	ToolPresetParams ToolPresetParams `gorm:"column:tool_preset_params;type:jsonb" json:"toolPresetParams"`
	MultiModal       *bool            `gorm:"column:multi_modal" json:"multiModal"`

	entity.BaseEntity
}

func (AgentEntity) TableName() string {
	return "agents"
}

// IsEnabled 检查Agent是否启用
func (e *AgentEntity) IsEnabled() bool {
	return e.Enabled != nil && *e.Enabled
}

// CheckEnabled 检查Agent是否启用，未启用则返回错误
func (e *AgentEntity) CheckEnabled() error {
	if !e.IsEnabled() {
		return exception.NewBusinessException("助理未激活")
	}
	return nil
}

// Enable 启用Agent
func (e *AgentEntity) Enable() {
	enabled := true
	e.Enabled = &enabled
	e.UpdatedAt = time.Now()
}

// Disable 禁用Agent
func (e *AgentEntity) Disable() {
	enabled := false
	e.Enabled = &enabled
	e.UpdatedAt = time.Now()
}

// PublishVersion 发布新版本
func (e *AgentEntity) PublishVersion(versionID string) {
	e.PublishedVersion = versionID
	e.UpdatedAt = time.Now()
}

// UpdateBasicInfo 更新基本信息
func (e *AgentEntity) UpdateBasicInfo(name, avatar, description string) {
	e.Name = name
	e.Avatar = avatar
	e.Description = description
	e.UpdatedAt = time.Now()
}

// CreateNew 创建新的Agent
func CreateNewAgent(name, description, avatar, userID string) *AgentEntity {
	now := time.Now()
	enabled := true
	return &AgentEntity{
		Name:        name,
		Description: description,
		Avatar:      avatar,
		UserID:      userID,
		Enabled:     &enabled,
		ToolIDs:     StringList{},
		KnowledgeBaseIDs: StringList{},
		BaseEntity: entity.BaseEntity{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// AgentVersionEntity Agent版本实体（对应 Java 的 AgentVersionEntity）
type AgentVersionEntity struct {
	ID               string           `gorm:"column:id;primaryKey" json:"id"`
	AgentID          string           `gorm:"column:agent_id" json:"agentId"`
	Name             string           `gorm:"column:name" json:"name"`
	Avatar           string           `gorm:"column:avatar" json:"avatar"`
	Description      string           `gorm:"column:description" json:"description"`
	VersionNumber    string           `gorm:"column:version_number" json:"versionNumber"`
	SystemPrompt     string           `gorm:"column:system_prompt" json:"systemPrompt"`
	WelcomeMessage   string           `gorm:"column:welcome_message" json:"welcomeMessage"`
	ToolIDs          StringList       `gorm:"column:tool_ids;type:jsonb" json:"toolIds"`
	KnowledgeBaseIDs StringList       `gorm:"column:knowledge_base_ids;type:jsonb" json:"knowledgeBaseIds"`
	ChangeLog        string           `gorm:"column:change_log" json:"changeLog"`
	PublishStatus    int              `gorm:"column:publish_status" json:"publishStatus"`
	RejectReason     string           `gorm:"column:reject_reason" json:"rejectReason"`
	ReviewTime       *time.Time       `gorm:"column:review_time" json:"reviewTime"`
	PublishedAt      *time.Time       `gorm:"column:published_at" json:"publishedAt"`
	UserID           string           `gorm:"column:user_id" json:"userId"`
	ToolPresetParams ToolPresetParams `gorm:"column:tool_preset_params;type:jsonb" json:"toolPresetParams"`
	MultiModal       *bool            `gorm:"column:multi_modal" json:"multiModal"`

	entity.BaseEntity
}

func (AgentVersionEntity) TableName() string {
	return "agent_versions"
}

// UpdatePublishStatus 更新发布状态
func (e *AgentVersionEntity) UpdatePublishStatus(status PublishStatus) {
	e.PublishStatus = int(status)
	now := time.Now()
	e.ReviewTime = &now
}

// Reject 拒绝发布
func (e *AgentVersionEntity) Reject(reason string) {
	e.PublishStatus = int(PublishStatusRejected)
	e.RejectReason = reason
	now := time.Now()
	e.ReviewTime = &now
}

// CreateVersionFromAgent 从Agent实体创建版本
func CreateVersionFromAgent(agent *AgentEntity, versionNumber, changeLog string) *AgentVersionEntity {
	now := time.Now()
	return &AgentVersionEntity{
		AgentID:          agent.ID,
		Name:             agent.Name,
		Avatar:           agent.Avatar,
		Description:      agent.Description,
		VersionNumber:    versionNumber,
		SystemPrompt:     agent.SystemPrompt,
		WelcomeMessage:   agent.WelcomeMessage,
		ToolIDs:          agent.ToolIDs,
		KnowledgeBaseIDs: agent.KnowledgeBaseIDs,
		ChangeLog:        changeLog,
		UserID:           agent.UserID,
		PublishStatus:    int(PublishStatusReviewing),
		ReviewTime:       &now,
		PublishedAt:      &now,
		ToolPresetParams: agent.ToolPresetParams,
		MultiModal:       agent.MultiModal,
		BaseEntity: entity.BaseEntity{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// AgentWorkspaceEntity Agent工作区实体（对应 Java 的 AgentWorkspaceEntity）
type AgentWorkspaceEntity struct {
	ID             string         `gorm:"column:id;primaryKey" json:"id"`
	AgentID        string         `gorm:"column:agent_id" json:"agentId"`
	UserID         string         `gorm:"column:user_id" json:"userId"`
	LLMModelConfig *LLMModelConfig `gorm:"column:llm_model_config;serializer:json" json:"llmModelConfig"`

	entity.BaseEntity
}

func (AgentWorkspaceEntity) TableName() string {
	return "agent_workspace"
}

// GetLLMModelConfig 获取模型配置（带兜底）
func (e *AgentWorkspaceEntity) GetLLMModelConfig() *LLMModelConfig {
	if e.LLMModelConfig == nil {
		e.LLMModelConfig = NewDefaultLLMModelConfig()
	}
	return e.LLMModelConfig
}
