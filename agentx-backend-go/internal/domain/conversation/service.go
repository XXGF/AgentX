package conversation

import (
	"github.com/google/uuid"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// SessionDomainService 会话领域服务（对应 Java 的 SessionDomainService）
type SessionDomainService struct {
	sessionRepo SessionRepository
}

// NewSessionDomainService 创建会话领域服务
func NewSessionDomainService(sessionRepo SessionRepository) *SessionDomainService {
	return &SessionDomainService{sessionRepo: sessionRepo}
}

// GetSessionsByAgentID 根据agentId获取会话列表
func (s *SessionDomainService) GetSessionsByAgentID(agentID, userID string) ([]SessionEntity, error) {
	return s.sessionRepo.FindByAgentIDAndUserID(agentID, userID)
}

// DeleteSession 删除会话
func (s *SessionDomainService) DeleteSession(sessionID, userID string) error {
	return s.sessionRepo.DeleteByIDAndUserID(sessionID, userID)
}

// UpdateSession 更新会话标题
func (s *SessionDomainService) UpdateSession(sessionID, userID, title string) error {
	session, err := s.sessionRepo.FindByIDAndUserID(sessionID, userID)
	if err != nil {
		return err
	}
	if session == nil {
		return exception.NewBusinessException("会话不存在")
	}
	session.Title = title
	return s.sessionRepo.Update(session)
}

// CreateSession 创建会话
func (s *SessionDomainService) CreateSession(agentID, userID string) (*SessionEntity, error) {
	session := CreateNewSession("新会话", userID, agentID)
	session.ID = uuid.New().String()
	if err := s.sessionRepo.Create(session); err != nil {
		return nil, err
	}
	return session, nil
}

// CheckSessionExist 检查会话是否存在
func (s *SessionDomainService) CheckSessionExist(sessionID, userID string) error {
	session, err := s.sessionRepo.FindByIDAndUserID(sessionID, userID)
	if err != nil {
		return err
	}
	if session == nil {
		return exception.NewBusinessException("会话不存在")
	}
	return nil
}

// Find 查找会话（不抛异常）
func (s *SessionDomainService) Find(sessionID, userID string) (*SessionEntity, error) {
	return s.sessionRepo.FindByIDAndUserID(sessionID, userID)
}

// DeleteSessions 批量删除会话
func (s *SessionDomainService) DeleteSessions(sessionIDs []string) error {
	return s.sessionRepo.DeleteByIDs(sessionIDs)
}

// GetSession 获取会话（不存在则抛异常）
func (s *SessionDomainService) GetSession(sessionID, userID string) (*SessionEntity, error) {
	session, err := s.sessionRepo.FindByIDAndUserID(sessionID, userID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, exception.NewBusinessException("会话不存在")
	}
	return session, nil
}

// MessageDomainService 消息领域服务（对应 Java 的 MessageDomainService）
type MessageDomainService struct {
	messageRepo MessageRepository
	contextRepo ContextRepository
}

// NewMessageDomainService 创建消息领域服务
func NewMessageDomainService(messageRepo MessageRepository, contextRepo ContextRepository) *MessageDomainService {
	return &MessageDomainService{
		messageRepo: messageRepo,
		contextRepo: contextRepo,
	}
}

// ListByIDs 根据ID列表获取消息
func (s *MessageDomainService) ListByIDs(ids []string) ([]MessageEntity, error) {
	return s.messageRepo.FindByIDs(ids)
}

// SaveMessageAndUpdateContext 保存消息并更新上下文
func (s *MessageDomainService) SaveMessageAndUpdateContext(messages []MessageEntity, context *ContextEntity) error {
	if len(messages) == 0 {
		return nil
	}
	// 设置ID
	for i := range messages {
		if messages[i].ID == "" {
			messages[i].ID = uuid.New().String()
		}
	}
	if err := s.messageRepo.CreateBatch(messages); err != nil {
		return err
	}
	// 更新上下文
	for _, msg := range messages {
		context.ActiveMessages = append(context.ActiveMessages, msg.ID)
	}
	return s.contextRepo.InsertOrUpdate(context)
}

// SaveMessage 保存消息
func (s *MessageDomainService) SaveMessages(messages []MessageEntity) error {
	return s.messageRepo.CreateBatch(messages)
}

// UpdateMessage 更新消息
func (s *MessageDomainService) UpdateMessage(message *MessageEntity) error {
	return s.messageRepo.Update(message)
}

// IsFirstConversation 是否为首次对话
func (s *MessageDomainService) IsFirstConversation(sessionID string) (bool, error) {
	count, err := s.messageRepo.CountBySessionID(sessionID)
	if err != nil {
		return false, err
	}
	return count <= 3, nil
}

// ContextDomainService 上下文领域服务（对应 Java 的 ContextDomainService）
type ContextDomainService struct {
	contextRepo ContextRepository
}

// NewContextDomainService 创建上下文领域服务
func NewContextDomainService(contextRepo ContextRepository) *ContextDomainService {
	return &ContextDomainService{contextRepo: contextRepo}
}

// GetBySessionID 根据sessionId获取上下文
func (s *ContextDomainService) GetBySessionID(sessionID string) (*ContextEntity, error) {
	context, err := s.contextRepo.FindBySessionID(sessionID)
	if err != nil {
		return nil, err
	}
	if context == nil {
		return nil, exception.NewBusinessException("消息上下文不存在")
	}
	return context, nil
}

// FindBySessionID 查找上下文（不抛异常）
func (s *ContextDomainService) FindBySessionID(sessionID string) (*ContextEntity, error) {
	return s.contextRepo.FindBySessionID(sessionID)
}

// InsertOrUpdate 插入或更新上下文
func (s *ContextDomainService) InsertOrUpdate(context *ContextEntity) (*ContextEntity, error) {
	if err := s.contextRepo.InsertOrUpdate(context); err != nil {
		return nil, err
	}
	return context, nil
}

// ConversationDomainService 对话领域服务（对应 Java 的 ConversationDomainService）
type ConversationDomainService struct {
	messageRepo MessageRepository
}

// NewConversationDomainService 创建对话领域服务
func NewConversationDomainService(messageRepo MessageRepository) *ConversationDomainService {
	return &ConversationDomainService{messageRepo: messageRepo}
}

// GetConversationMessages 获取会话中的消息列表（排除SUMMARY角色）
func (s *ConversationDomainService) GetConversationMessages(sessionID string) ([]MessageEntity, error) {
	return s.messageRepo.FindBySessionIDExcludeRole(sessionID, RoleSummary)
}

// InsertBatchMessage 批量插入消息
func (s *ConversationDomainService) InsertBatchMessage(messages []MessageEntity) error {
	return s.messageRepo.CreateBatch(messages)
}

// SaveMessage 保存单条消息
func (s *ConversationDomainService) SaveMessage(message *MessageEntity) error {
	if message.ID == "" {
		message.ID = uuid.New().String()
	}
	return s.messageRepo.Create(message)
}

// DeleteConversationMessages 删除会话下的消息
func (s *ConversationDomainService) DeleteConversationMessages(sessionID string) error {
	return s.messageRepo.DeleteBySessionID(sessionID)
}

// DeleteConversationMessagesBySessionIDs 批量删除会话下的消息
func (s *ConversationDomainService) DeleteConversationMessagesBySessionIDs(sessionIDs []string) error {
	return s.messageRepo.DeleteBySessionIDs(sessionIDs)
}

// UpdateMessageTokenCount 更新消息的token数量
func (s *ConversationDomainService) UpdateMessageTokenCount(message *MessageEntity) error {
	return s.messageRepo.Update(message)
}
