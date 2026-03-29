package conversation

import (
	"time"

	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/conversation"
)

// SessionDTO 会话DTO（对应 Java 的 SessionDTO）
type SessionDTO struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	IsArchived  bool      `json:"isArchived"`
	AgentID     string    `json:"agentId"`
	MultiModal  *bool     `json:"multiModal,omitempty"`
}

// MessageDTO 消息DTO（对应 Java 的 MessageDTO）
type MessageDTO struct {
	ID          string             `json:"id"`
	Role        domain.Role        `json:"role"`
	Content     string             `json:"content"`
	CreatedAt   time.Time          `json:"createdAt"`
	Provider    string             `json:"provider"`
	Model       string             `json:"model"`
	MessageType domain.MessageType `json:"messageType"`
	FileUrls    []string           `json:"fileUrls"`
}

// ChatRequest 聊天请求DTO（对应 Java 的 ChatRequest）
type ChatRequest struct {
	Message   string   `json:"message" binding:"required"`
	SessionID string   `json:"sessionId" binding:"required"`
	FileUrls  []string `json:"fileUrls"`
}

// ConversationRequest 对话请求DTO
type ConversationRequest struct {
	Message string `json:"message"`
}
