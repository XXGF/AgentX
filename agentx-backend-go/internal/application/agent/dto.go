package agent

import (
	"time"

	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/agent"
)

// AgentDTO Agent数据传输对象（对应 Java 的 AgentDTO）
type AgentDTO struct {
	ID               string                  `json:"id"`
	Name             string                  `json:"name"`
	Avatar           string                  `json:"avatar"`
	Description      string                  `json:"description"`
	SystemPrompt     string                  `json:"systemPrompt"`
	WelcomeMessage   string                  `json:"welcomeMessage"`
	ToolIDs          []string                `json:"toolIds"`
	KnowledgeBaseIDs []string                `json:"knowledgeBaseIds"`
	PublishedVersion string                  `json:"publishedVersion"`
	Enabled          *bool                   `json:"enabled"`
	UserID           string                  `json:"userId"`
	UserNickname     string                  `json:"userNickname,omitempty"`
	ToolPresetParams domain.ToolPresetParams `json:"toolPresetParams,omitempty"`
	MultiModal       *bool                   `json:"multiModal"`
	CreatedAt        time.Time               `json:"createdAt"`
	UpdatedAt        time.Time               `json:"updatedAt"`
}

// AgentVersionDTO Agent版本DTO（对应 Java 的 AgentVersionDTO）
type AgentVersionDTO struct {
	ID               string     `json:"id"`
	AgentID          string     `json:"agentId"`
	Name             string     `json:"name"`
	Avatar           string     `json:"avatar"`
	Description      string     `json:"description"`
	VersionNumber    string     `json:"versionNumber"`
	SystemPrompt     string     `json:"systemPrompt"`
	WelcomeMessage   string     `json:"welcomeMessage"`
	ToolIDs          []string   `json:"toolIds"`
	KnowledgeBaseIDs []string   `json:"knowledgeBaseIds"`
	ChangeLog        string     `json:"changeLog"`
	PublishStatus    int        `json:"publishStatus"`
	RejectReason     string     `json:"rejectReason"`
	ReviewTime       *time.Time `json:"reviewTime"`
	PublishedAt      *time.Time `json:"publishedAt"`
	UserID           string     `json:"userId"`
	IsAddWorkspace   *bool      `json:"isAddWorkspace,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

// AgentWorkspaceDTO Agent工作区DTO
type AgentWorkspaceDTO struct {
	ID             string                `json:"id"`
	AgentID        string                `json:"agentId"`
	UserID         string                `json:"userId"`
	LLMModelConfig *domain.LLMModelConfig `json:"llmModelConfig"`
}

// CreateAgentRequest 创建Agent请求（对应 Java 的 CreateAgentRequest）
type CreateAgentRequest struct {
	Name             string                  `json:"name" binding:"required"`
	Description      string                  `json:"description"`
	Avatar           string                  `json:"avatar"`
	SystemPrompt     string                  `json:"systemPrompt"`
	WelcomeMessage   string                  `json:"welcomeMessage"`
	ToolIDs          []string                `json:"toolIds"`
	KnowledgeBaseIDs []string                `json:"knowledgeBaseIds"`
	ToolPresetParams domain.ToolPresetParams `json:"toolPresetParams"`
	MultiModal       *bool                   `json:"multiModal"`
}

// UpdateAgentRequest 更新Agent请求（对应 Java 的 UpdateAgentRequest）
type UpdateAgentRequest struct {
	ID               string                  `json:"id"`
	Name             string                  `json:"name" binding:"required"`
	Avatar           string                  `json:"avatar"`
	Description      string                  `json:"description"`
	Enabled          *bool                   `json:"enabled"`
	SystemPrompt     string                  `json:"systemPrompt"`
	WelcomeMessage   string                  `json:"welcomeMessage"`
	ToolIDs          []string                `json:"toolIds"`
	KnowledgeBaseIDs []string                `json:"knowledgeBaseIds"`
	ToolPresetParams domain.ToolPresetParams `json:"toolPresetParams"`
	MultiModal       *bool                   `json:"multiModal"`
}

// SearchAgentsRequest 搜索Agent请求
type SearchAgentsRequest struct {
	Name string `form:"name" json:"name"`
}

// PublishAgentVersionRequest 发布Agent版本请求（对应 Java 的 PublishAgentVersionRequest）
type PublishAgentVersionRequest struct {
	VersionNumber string `json:"versionNumber" binding:"required"`
	ChangeLog     string `json:"changeLog" binding:"required"`
}

// ReviewAgentVersionRequest 审核Agent版本请求
type ReviewAgentVersionRequest struct {
	Status       domain.PublishStatus `json:"status"`
	RejectReason string               `json:"rejectReason"`
}

// QueryAgentRequest 管理员查询Agent请求
type QueryAgentRequest struct {
	Keyword  string `form:"keyword" json:"keyword"`
	Enabled  *bool  `form:"enabled" json:"enabled"`
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"pageSize" json:"pageSize"`
}
