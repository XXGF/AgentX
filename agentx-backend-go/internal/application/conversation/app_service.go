package conversation

import (
	domainAgent "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/agent"
	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/conversation"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// AgentSessionAppService Agent会话应用服务（对应 Java 的 AgentSessionAppService）
type AgentSessionAppService struct {
	agentWorkspaceDomainService *domainAgent.WorkspaceDomainService
	agentDomainService          *domainAgent.DomainService
	sessionDomainService        *domain.SessionDomainService
	conversationDomainService   *domain.ConversationDomainService
}

// NewAgentSessionAppService 创建Agent会话应用服务
func NewAgentSessionAppService(
	agentWorkspaceDomainService *domainAgent.WorkspaceDomainService,
	agentDomainService *domainAgent.DomainService,
	sessionDomainService *domain.SessionDomainService,
	conversationDomainService *domain.ConversationDomainService,
) *AgentSessionAppService {
	return &AgentSessionAppService{
		agentWorkspaceDomainService: agentWorkspaceDomainService,
		agentDomainService:          agentDomainService,
		sessionDomainService:        sessionDomainService,
		conversationDomainService:   conversationDomainService,
	}
}

// GetAgentSessionList 获取助理下的会话列表
func (s *AgentSessionAppService) GetAgentSessionList(userID, agentID string) ([]*SessionDTO, error) {
	// 校验该 agent 是否被添加了工作区
	isOwner, err := s.agentDomainService.Exist(agentID, userID)
	if err != nil {
		return nil, err
	}
	isInWorkspace, err := s.agentWorkspaceDomainService.Exist(agentID, userID)
	if err != nil {
		return nil, err
	}

	if !isOwner && !isInWorkspace {
		return nil, exception.NewBusinessException("助理不存在")
	}

	// 获取对应的会话列表
	sessions, err := s.sessionDomainService.GetSessionsByAgentID(agentID, userID)
	if err != nil {
		return nil, err
	}

	if len(sessions) == 0 {
		// 如果会话列表为空，则新创建一个并且返回
		session, err := s.sessionDomainService.CreateSession(agentID, userID)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, *session)
	}

	// 获取Agent信息以判断multiModal
	agent, err := s.agentDomainService.GetAgentByID(agentID)
	if err != nil {
		return nil, err
	}
	multiModal := agent.MultiModal
	if agent.UserID != userID {
		latestVersion, err := s.agentDomainService.GetLatestAgentVersion(agentID)
		if err == nil && latestVersion != nil {
			multiModal = latestVersion.MultiModal
		}
	}

	dtos := SessionsToDTOs(sessions)
	for _, dto := range dtos {
		dto.MultiModal = multiModal
	}
	return dtos, nil
}

// CreateSession 创建会话
func (s *AgentSessionAppService) CreateSession(userID, agentID string) (*SessionDTO, error) {
	session, err := s.sessionDomainService.CreateSession(agentID, userID)
	if err != nil {
		return nil, err
	}

	// 获取Agent的欢迎消息
	agent, err := s.agentDomainService.GetAgentByID(agentID)
	if err != nil {
		return nil, err
	}

	if agent.WelcomeMessage != "" {
		welcomeMsg := &domain.MessageEntity{
			SessionID:   session.ID,
			Role:        domain.RoleSystem,
			Content:     agent.WelcomeMessage,
			MessageType: domain.MessageTypeText,
		}
		if err := s.conversationDomainService.SaveMessage(welcomeMsg); err != nil {
			return nil, err
		}
	}

	return SessionToDTO(session), nil
}

// UpdateSession 更新会话
func (s *AgentSessionAppService) UpdateSession(id, userID, title string) error {
	return s.sessionDomainService.UpdateSession(id, userID, title)
}

// DeleteSession 删除会话
func (s *AgentSessionAppService) DeleteSession(id, userID string) error {
	if err := s.sessionDomainService.DeleteSession(id, userID); err != nil {
		return err
	}
	// 删除会话下的消息
	if err := s.conversationDomainService.DeleteConversationMessages(id); err != nil {
		return err
	}
	// 注意：删除关联的定时任务可通过事件驱动机制处理，当前版本暂不自动清理
	return nil
}

// ConversationAppService 对话应用服务（对应 Java 的 ConversationAppService 的基础部分）
type ConversationAppService struct {
	conversationDomainService *domain.ConversationDomainService
	sessionDomainService      *domain.SessionDomainService
}

// NewConversationAppService 创建对话应用服务
func NewConversationAppService(
	conversationDomainService *domain.ConversationDomainService,
	sessionDomainService *domain.SessionDomainService,
) *ConversationAppService {
	return &ConversationAppService{
		conversationDomainService: conversationDomainService,
		sessionDomainService:      sessionDomainService,
	}
}

// GetConversationMessages 获取会话中的消息列表
func (s *ConversationAppService) GetConversationMessages(sessionID, userID string) ([]*MessageDTO, error) {
	// 检查会话是否存在
	if err := s.sessionDomainService.CheckSessionExist(sessionID, userID); err != nil {
		return nil, err
	}

	messages, err := s.conversationDomainService.GetConversationMessages(sessionID)
	if err != nil {
		return nil, err
	}
	return MessagesToDTOs(messages), nil
}

// 注意：Chat 和 PreviewAgent 方法已迁移到 ChatAppService（internal/application/chat）
// 通过 ChatController 和 SessionController 提供 SSE 流式响应、Agent工作流、MCP工具调用等功能
