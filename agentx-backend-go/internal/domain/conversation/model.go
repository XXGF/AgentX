package conversation

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
)

// SessionEntity 会话实体（对应 Java 的 SessionEntity）
type SessionEntity struct {
	ID          string `gorm:"column:id;primaryKey" json:"id"`
	Title       string `gorm:"column:title" json:"title"`
	UserID      string `gorm:"column:user_id" json:"userId"`
	AgentID     string `gorm:"column:agent_id" json:"agentId"`
	Description string `gorm:"column:description" json:"description"`
	IsArchived  bool   `gorm:"column:is_archived" json:"isArchived"`
	Metadata    string `gorm:"column:metadata" json:"metadata"`

	entity.BaseEntity
}

func (SessionEntity) TableName() string {
	return "sessions"
}

// CreateNewSession 创建新会话
func CreateNewSession(title, userID, agentID string) *SessionEntity {
	now := time.Now()
	return &SessionEntity{
		Title:      title,
		UserID:     userID,
		AgentID:    agentID,
		IsArchived: false,
		BaseEntity: entity.BaseEntity{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// Update 更新会话信息
func (e *SessionEntity) Update(title, description string) {
	e.Title = title
	e.Description = description
	e.UpdatedAt = time.Now()
}

// Archive 归档会话
func (e *SessionEntity) Archive() {
	e.IsArchived = true
	e.UpdatedAt = time.Now()
}

// Unarchive 恢复已归档会话
func (e *SessionEntity) Unarchive() {
	e.IsArchived = false
	e.UpdatedAt = time.Now()
}

// JSONStringList JSON字符串列表类型（用于 file_urls 等字段）
type JSONStringList []string

// Value 实现 driver.Valuer 接口
func (s JSONStringList) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	return string(b), err
}

// Scan 实现 sql.Scanner 接口
func (s *JSONStringList) Scan(value interface{}) error {
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

// MessageEntity 消息实体（对应 Java 的 MessageEntity）
type MessageEntity struct {
	ID             string         `gorm:"column:id;primaryKey" json:"id"`
	SessionID      string         `gorm:"column:session_id" json:"sessionId"`
	Role           Role           `gorm:"column:role" json:"role"`
	Content        string         `gorm:"column:content" json:"content"`
	MessageType    MessageType    `gorm:"column:message_type" json:"messageType"`
	TokenCount     int            `gorm:"column:token_count" json:"tokenCount"`
	BodyTokenCount int            `gorm:"column:body_token_count" json:"bodyTokenCount"`
	Provider       string         `gorm:"column:provider" json:"provider"`
	Model          string         `gorm:"column:model" json:"model"`
	Metadata       string         `gorm:"column:metadata" json:"metadata"`
	FileUrls       JSONStringList `gorm:"column:file_urls;type:jsonb" json:"fileUrls"`

	entity.BaseEntity
}

func (MessageEntity) TableName() string {
	return "messages"
}

// IsUserMessage 是否为用户消息
func (e *MessageEntity) IsUserMessage() bool {
	return e.Role == RoleUser
}

// IsAIMessage 是否为AI消息
func (e *MessageEntity) IsAIMessage() bool {
	return e.Role == RoleAssistant
}

// IsSystemMessage 是否为系统消息
func (e *MessageEntity) IsSystemMessage() bool {
	return e.Role == RoleSystem
}

// IsSummaryMessage 是否为摘要消息
func (e *MessageEntity) IsSummaryMessage() bool {
	return e.Role == RoleSummary
}

// ContextEntity 上下文实体（对应 Java 的 ContextEntity）
type ContextEntity struct {
	ID             string         `gorm:"column:id;primaryKey" json:"id"`
	SessionID      string         `gorm:"column:session_id" json:"sessionId"`
	ActiveMessages JSONStringList `gorm:"column:active_messages;type:jsonb" json:"activeMessages"`
	Summary        string         `gorm:"column:summary" json:"summary"`

	entity.BaseEntity
}

func (ContextEntity) TableName() string {
	return "context"
}
