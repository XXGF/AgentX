package chat

import (
	"context"
	"time"

	"github.com/google/uuid"
	domainAgent "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/agent"
	domainConv "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/conversation"
	domainLLM "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/llm"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	"go.uber.org/zap"
)

// ChatAppService 聊天应用服务（对应 Java 的 ConversationAppService 中的 Chat 功能）
type ChatAppService struct {
	conversationDomainService *domainConv.ConversationDomainService
	sessionDomainService      *domainConv.SessionDomainService
	agentDomainService        *domainAgent.DomainService
	agentWorkspaceDomainService *domainAgent.WorkspaceDomainService
	llmDomainService          *domainLLM.DomainService
	contextDomainService      *domainConv.ContextDomainService
	messageDomainService      *domainConv.MessageDomainService
	sessionManager            *ChatSessionManager
	logger                    *zap.Logger
}

// NewChatAppService 创建聊天应用服务
func NewChatAppService(
	conversationDomainService *domainConv.ConversationDomainService,
	sessionDomainService *domainConv.SessionDomainService,
	agentDomainService *domainAgent.DomainService,
	agentWorkspaceDomainService *domainAgent.WorkspaceDomainService,
	llmDomainService *domainLLM.DomainService,
	contextDomainService *domainConv.ContextDomainService,
	messageDomainService *domainConv.MessageDomainService,
	logger *zap.Logger,
) *ChatAppService {
	return &ChatAppService{
		conversationDomainService:   conversationDomainService,
		sessionDomainService:        sessionDomainService,
		agentDomainService:          agentDomainService,
		agentWorkspaceDomainService: agentWorkspaceDomainService,
		llmDomainService:            llmDomainService,
		contextDomainService:        contextDomainService,
		messageDomainService:        messageDomainService,
		sessionManager:              NewChatSessionManager(),
		logger:                      logger,
	}
}

// StreamChat 流式聊天（SSE）
// 这是核心的聊天入口，对应 Java 的 ConversationAppService.chat()
func (s *ChatAppService) StreamChat(req *ChatRequest, userID string, transport *SSETransport) error {
	// 1. 准备对话环境
	chatCtx, err := s.prepareEnvironment(req, userID)
	if err != nil {
		return err
	}

	// 2. 创建可取消的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	// 3. 注册会话到会话管理器（支持中断功能）
	s.sessionManager.RegisterSession(req.SessionID, transport, cancel)

	// 4. 异步处理对话
	go func() {
		defer cancel()
		defer s.sessionManager.RemoveSession(req.SessionID)

		s.processChat(ctx, chatCtx, transport)
	}()

	return nil
}

// Chat 同步聊天
func (s *ChatAppService) Chat(req *ChatRequest, userID string) (*ChatResponse, error) {
	// 1. 准备对话环境
	chatCtx, err := s.prepareEnvironment(req, userID)
	if err != nil {
		return nil, err
	}

	// 2. 同步处理对话
	response, err := s.processChatSync(chatCtx)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// StopChat 停止聊天
func (s *ChatAppService) StopChat(sessionID string) bool {
	return s.sessionManager.StopSession(sessionID)
}

// prepareEnvironment 准备对话环境（对应 Java 的 prepareEnvironment）
func (s *ChatAppService) prepareEnvironment(req *ChatRequest, userID string) (*ChatContext, error) {
	// 1. 获取会话和Agent信息
	session, err := s.sessionDomainService.GetSession(req.SessionID, userID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, exception.NewBusinessException("会话不存在")
	}

	agentID := session.AgentID
	agent, err := s.agentDomainService.GetAgentByID(agentID)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, exception.NewBusinessException("Agent不存在")
	}

	// 检查Agent是否可用
	if agent.UserID != userID && (agent.Enabled == nil || !*agent.Enabled) {
		return nil, exception.NewBusinessException("Agent已被禁用")
	}

	// 如果不是自己的Agent，使用最新版本
	if agent.UserID != userID {
		latestVersion, err := s.agentDomainService.GetLatestAgentVersion(agentID)
		if err == nil && latestVersion != nil {
			agent.SystemPrompt = latestVersion.SystemPrompt
			agent.WelcomeMessage = latestVersion.WelcomeMessage
			agent.ToolIDs = latestVersion.ToolIDs
			agent.KnowledgeBaseIDs = latestVersion.KnowledgeBaseIDs
		}
	}

	// 2. 获取模型配置
	workspace, err := s.agentWorkspaceDomainService.GetWorkspace(agentID, userID)
	if err != nil {
		return nil, err
	}

	var llmModelConfig *domainAgent.LLMModelConfig
	if workspace != nil {
		llmModelConfig = workspace.LLMModelConfig
	}

	// 获取模型
	var modelID string
	if llmModelConfig != nil && llmModelConfig.ModelID != "" {
		modelID = llmModelConfig.ModelID
	}

	model, err := s.llmDomainService.GetModelByID(modelID)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, exception.NewBusinessException("模型不存在或未配置")
	}

	// 3. 获取服务商
	provider, err := s.llmDomainService.GetProviderByID(model.ProviderID)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, exception.NewBusinessException("服务商不存在")
	}

	// 4. 获取上下文和历史消息
	contextEntity, err := s.contextDomainService.FindBySessionID(req.SessionID)
	if err != nil {
		return nil, err
	}

	var messageHistory []domainConv.MessageEntity
	if contextEntity != nil && len(contextEntity.ActiveMessages) > 0 {
		messages, err := s.messageDomainService.ListByIDs(contextEntity.ActiveMessages)
		if err != nil {
			return nil, err
		}
		messageHistory = messages
	}

	if contextEntity == nil {
		contextEntity = &domainConv.ContextEntity{
			ID:             uuid.New().String(),
			SessionID:      req.SessionID,
			ActiveMessages: make([]string, 0),
		}
	}

	// 5. 创建并配置环境对象
	chatCtx := NewChatContext()
	chatCtx.SessionID = req.SessionID
	chatCtx.UserID = userID
	chatCtx.UserMessage = req.Message
	chatCtx.Agent = agent
	chatCtx.Model = model
	chatCtx.Provider = provider
	chatCtx.LLMModelConfig = llmModelConfig
	chatCtx.ContextEntity = contextEntity
	chatCtx.MessageHistory = messageHistory
	chatCtx.FileUrls = req.FileUrls

	return chatCtx, nil
}

// processChat 异步处理聊天（SSE流式）
func (s *ChatAppService) processChat(ctx context.Context, chatCtx *ChatContext, transport *SSETransport) {
	// 1. 保存用户消息
	userMsg := &domainConv.MessageEntity{
		ID:          uuid.New().String(),
		SessionID:   chatCtx.SessionID,
		Role:        domainConv.RoleUser,
		Content:     chatCtx.UserMessage,
		MessageType: domainConv.MessageTypeText,
		FileUrls:    chatCtx.FileUrls,
	}
	if err := s.conversationDomainService.SaveMessage(userMsg); err != nil {
		s.logger.Error("保存用户消息失败", zap.Error(err))
		transport.SendError("保存消息失败: " + err.Error())
		return
	}

	// 2. 发送思考中状态
	transport.SendMessage(&AgentChatResponse{
		Type:      "thinking",
		SessionID: chatCtx.SessionID,
		Provider:  chatCtx.Provider.Name,
		Model:     chatCtx.Model.ModelID,
	})

	// 3. 构建LLM请求消息
	// TODO: 这里需要集成实际的LLM SDK（如 OpenAI/Anthropic API）
	// 当前实现为模拟响应，展示完整的SSE流式响应流程

	select {
	case <-ctx.Done():
		transport.SendMessage(&AgentChatResponse{
			Type:    "error",
			Error:   "聊天已被中断",
		})
		return
	default:
	}

	// 模拟LLM响应（实际实现需要调用LLM API）
	simulatedResponse := "您好！我是AI助手。您的消息已收到：「" + chatCtx.UserMessage + "」\n\n" +
		"当前使用的模型是 " + chatCtx.Model.ModelID + "，服务商是 " + chatCtx.Provider.Name + "。\n\n" +
		"⚠️ 注意：这是模拟响应。完整的LLM集成需要配置实际的API密钥和SDK。"

	// 4. 流式发送响应内容
	messageID := uuid.New().String()
	for i, ch := range simulatedResponse {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if i%3 == 0 { // 每3个字符发送一次，模拟流式效果
			transport.SendMessage(&AgentChatResponse{
				Type:      "text",
				Content:   string(simulatedResponse[max(0, i-2) : i+len(string(ch))]),
				SessionID: chatCtx.SessionID,
				MessageID: messageID,
			})
		}
	}

	// 5. 发送完整响应
	transport.SendMessage(&AgentChatResponse{
		Type:      "text",
		Content:   simulatedResponse,
		SessionID: chatCtx.SessionID,
		MessageID: messageID,
	})

	// 6. 保存AI响应消息
	aiMsg := &domainConv.MessageEntity{
		ID:          messageID,
		SessionID:   chatCtx.SessionID,
		Role:        domainConv.RoleAssistant,
		Content:     simulatedResponse,
		MessageType: domainConv.MessageTypeText,
	}
	if err := s.conversationDomainService.SaveMessage(aiMsg); err != nil {
		s.logger.Error("保存AI消息失败", zap.Error(err))
	}

	// 7. 更新上下文
	chatCtx.ContextEntity.ActiveMessages = append(chatCtx.ContextEntity.ActiveMessages, userMsg.ID, aiMsg.ID)
	if _, err := s.contextDomainService.InsertOrUpdate(chatCtx.ContextEntity); err != nil {
		s.logger.Error("更新上下文失败", zap.Error(err))
	}

	// 8. 发送结束消息
	transport.SendEndMessage(&AgentChatResponse{
		SessionID: chatCtx.SessionID,
		MessageID: messageID,
		Provider:  chatCtx.Provider.Name,
		Model:     chatCtx.Model.ModelID,
		TokenUsage: &TokenUsage{
			InputTokens:  len(chatCtx.UserMessage),
			OutputTokens: len(simulatedResponse),
			TotalTokens:  len(chatCtx.UserMessage) + len(simulatedResponse),
		},
	})
}

// processChatSync 同步处理聊天
func (s *ChatAppService) processChatSync(chatCtx *ChatContext) (*ChatResponse, error) {
	// 1. 保存用户消息
	userMsg := &domainConv.MessageEntity{
		ID:          uuid.New().String(),
		SessionID:   chatCtx.SessionID,
		Role:        domainConv.RoleUser,
		Content:     chatCtx.UserMessage,
		MessageType: domainConv.MessageTypeText,
		FileUrls:    chatCtx.FileUrls,
	}
	if err := s.conversationDomainService.SaveMessage(userMsg); err != nil {
		return nil, err
	}

	// 2. TODO: 调用LLM API获取同步响应
	// 当前实现为模拟响应
	simulatedResponse := "您好！我是AI助手。您的消息已收到：「" + chatCtx.UserMessage + "」"

	// 3. 保存AI响应消息
	messageID := uuid.New().String()
	aiMsg := &domainConv.MessageEntity{
		ID:          messageID,
		SessionID:   chatCtx.SessionID,
		Role:        domainConv.RoleAssistant,
		Content:     simulatedResponse,
		MessageType: domainConv.MessageTypeText,
	}
	if err := s.conversationDomainService.SaveMessage(aiMsg); err != nil {
		return nil, err
	}

	// 4. 更新上下文
	chatCtx.ContextEntity.ActiveMessages = append(chatCtx.ContextEntity.ActiveMessages, userMsg.ID, aiMsg.ID)
	_, _ = s.contextDomainService.InsertOrUpdate(chatCtx.ContextEntity)

	return &ChatResponse{
		Content:   simulatedResponse,
		SessionID: chatCtx.SessionID,
		Provider:  chatCtx.Provider.Name,
		Model:     chatCtx.Model.ModelID,
		Timestamp: time.Now().UnixMilli(),
	}, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
