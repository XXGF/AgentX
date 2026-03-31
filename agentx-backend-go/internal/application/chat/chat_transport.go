package chat

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"go.uber.org/zap"
)

// ---- SSE Transport（对应 Java 的 SseMessageTransport）----

// SSEWriter SSE写入接口（由 gin.Context 的 ResponseWriter 实现）
type SSEWriter interface {
	io.Writer
	Flush()
}

// SSETransport SSE消息传输
type SSETransport struct {
	writer  SSEWriter
	flusher func()
	closed  bool
	mu      sync.Mutex
	logger  *zap.Logger
}

// NewSSETransport 创建SSE传输
func NewSSETransport(writer SSEWriter, logger *zap.Logger) *SSETransport {
	return &SSETransport{
		writer: writer,
		logger: logger,
	}
}

// SendMessage 发送SSE消息
func (t *SSETransport) SendMessage(response *AgentChatResponse) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return fmt.Errorf("SSE connection already closed")
	}

	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	// SSE格式: data: {json}\n\n
	_, err = fmt.Fprintf(t.writer, "data: %s\n\n", string(data))
	if err != nil {
		return err
	}

	t.writer.Flush()
	return nil
}

// SendEndMessage 发送结束消息
func (t *SSETransport) SendEndMessage(response *AgentChatResponse) error {
	response.Type = "done"
	return t.SendMessage(response)
}

// SendError 发送错误消息
func (t *SSETransport) SendError(errMsg string) error {
	return t.SendMessage(&AgentChatResponse{
		Type:  "error",
		Error: errMsg,
	})
}

// Close 关闭连接
func (t *SSETransport) Close() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.closed = true
}

// IsClosed 检查连接是否已关闭
func (t *SSETransport) IsClosed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.closed
}

// ---- ChatSessionManager 会话管理器（对应 Java 的 ChatSessionManager）----

// ChatSessionManager 管理活跃的聊天会话，支持中断功能
type ChatSessionManager struct {
	sessions sync.Map // sessionID -> *ChatSession
}

// ChatSession 聊天会话
type ChatSession struct {
	SessionID string
	Transport *SSETransport
	Cancel    func() // 取消函数
}

// NewChatSessionManager 创建会话管理器
func NewChatSessionManager() *ChatSessionManager {
	return &ChatSessionManager{}
}

// RegisterSession 注册会话
func (m *ChatSessionManager) RegisterSession(sessionID string, transport *SSETransport, cancel func()) {
	m.sessions.Store(sessionID, &ChatSession{
		SessionID: sessionID,
		Transport: transport,
		Cancel:    cancel,
	})
}

// StopSession 停止会话（中断聊天）
func (m *ChatSessionManager) StopSession(sessionID string) bool {
	if val, ok := m.sessions.LoadAndDelete(sessionID); ok {
		session := val.(*ChatSession)
		if session.Cancel != nil {
			session.Cancel()
		}
		if session.Transport != nil {
			session.Transport.Close()
		}
		return true
	}
	return false
}

// RemoveSession 移除会话
func (m *ChatSessionManager) RemoveSession(sessionID string) {
	m.sessions.Delete(sessionID)
}

// HasSession 检查会话是否存在
func (m *ChatSessionManager) HasSession(sessionID string) bool {
	_, ok := m.sessions.Load(sessionID)
	return ok
}
