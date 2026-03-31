package chat

import (
	domainAgent "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/agent"
	domainConv "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/conversation"
	domainLLM "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/llm"
)

// ChatContext 聊天上下文（对应 Java 的 ChatContext）
// 包含对话所需的所有信息
type ChatContext struct {
	// SessionID 会话ID
	SessionID string `json:"sessionId"`

	// UserID 用户ID
	UserID string `json:"userId"`

	// UserMessage 用户消息
	UserMessage string `json:"userMessage"`

	// Agent 智能体实体
	Agent *domainAgent.AgentEntity `json:"-"`

	// Model 模型实体
	Model *domainLLM.ModelEntity `json:"-"`

	// OriginalModel 原始模型实体（用于追踪模型切换）
	OriginalModel *domainLLM.ModelEntity `json:"-"`

	// Provider 服务商实体
	Provider *domainLLM.ProviderEntity `json:"-"`

	// OriginalProvider 原始服务商实体（用于追踪服务商切换）
	OriginalProvider *domainLLM.ProviderEntity `json:"-"`

	// LLMModelConfig 大模型配置
	LLMModelConfig *domainAgent.LLMModelConfig `json:"-"`

	// ContextEntity 上下文实体
	ContextEntity *domainConv.ContextEntity `json:"-"`

	// MessageHistory 历史消息列表
	MessageHistory []domainConv.MessageEntity `json:"-"`

	// McpServerNames 使用的 MCP server name 列表
	McpServerNames []string `json:"mcpServerNames"`

	// FileUrls 多模态的文件URL列表
	FileUrls []string `json:"fileUrls"`

	// InstanceID 高可用实例ID
	InstanceID string `json:"instanceId"`

	// Streaming 是否流式响应
	Streaming bool `json:"streaming"`

	// PublicAccess 是否为公开访问（嵌入模式）
	PublicAccess bool `json:"publicAccess"`

	// PublicID 公开访问ID（嵌入模式使用）
	PublicID string `json:"publicId"`
}

// NewChatContext 创建默认的聊天上下文
func NewChatContext() *ChatContext {
	return &ChatContext{
		Streaming:      true,
		McpServerNames: make([]string, 0),
		FileUrls:       make([]string, 0),
	}
}

// ChatRequest 聊天请求DTO（对应 Java 的 ChatRequest）
type ChatRequest struct {
	Message   string   `json:"message" binding:"required"`
	SessionID string   `json:"sessionId" binding:"required"`
	FileUrls  []string `json:"fileUrls"`
}

// ChatResponse 同步聊天响应DTO（对应 Java 的 ChatResponse）
type ChatResponse struct {
	Content   string `json:"content"`
	SessionID string `json:"sessionId"`
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	Timestamp int64  `json:"timestamp"`
}

// AgentChatResponse SSE流式响应数据（对应 Java 的 AgentChatResponse）
type AgentChatResponse struct {
	// Type 消息类型：text/tool_call/tool_result/error/done/thinking
	Type string `json:"type"`

	// Content 消息内容
	Content string `json:"content,omitempty"`

	// ToolName 工具名称（type=tool_call时）
	ToolName string `json:"toolName,omitempty"`

	// ToolArgs 工具参数（type=tool_call时）
	ToolArgs string `json:"toolArgs,omitempty"`

	// ToolResult 工具执行结果（type=tool_result时）
	ToolResult string `json:"toolResult,omitempty"`

	// SessionID 会话ID
	SessionID string `json:"sessionId,omitempty"`

	// MessageID 消息ID
	MessageID string `json:"messageId,omitempty"`

	// Provider 使用的服务商
	Provider string `json:"provider,omitempty"`

	// Model 使用的模型
	Model string `json:"model,omitempty"`

	// Error 错误信息（type=error时）
	Error string `json:"error,omitempty"`

	// TokenUsage Token使用量
	TokenUsage *TokenUsage `json:"tokenUsage,omitempty"`
}

// TokenUsage Token使用量统计
type TokenUsage struct {
	InputTokens  int `json:"inputTokens"`
	OutputTokens int `json:"outputTokens"`
	TotalTokens  int `json:"totalTokens"`
}

// AgentPreviewRequest Agent预览请求（对应 Java 的 AgentPreviewRequest）
type AgentPreviewRequest struct {
	UserMessage      string                            `json:"userMessage" binding:"required"`
	SystemPrompt     string                            `json:"systemPrompt"`
	ModelID          string                            `json:"modelId"`
	ToolIDs          []string                          `json:"toolIds"`
	ToolPresetParams map[string]map[string]interface{} `json:"toolPresetParams"`
	KnowledgeBaseIDs []string                          `json:"knowledgeBaseIds"`
	MessageHistory   []MessageHistoryItem              `json:"messageHistory"`
	FileUrls         []string                          `json:"fileUrls"`
}

// MessageHistoryItem 历史消息项
type MessageHistoryItem struct {
	ID       string `json:"id"`
	Role     string `json:"role"`
	Content  string `json:"content"`
	FileUrls []string `json:"fileUrls"`
}

// StopChatRequest 停止聊天请求
type StopChatRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
}
