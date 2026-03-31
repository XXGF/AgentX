package portal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appChat "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/chat"
	appConv "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/conversation"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
	"go.uber.org/zap"
)

// SessionController Agent会话控制器（对应 Java 的 PortalAgentSessionController）
type SessionController struct {
	agentSessionAppService  *appConv.AgentSessionAppService
	conversationAppService  *appConv.ConversationAppService
	chatAppService          *appChat.ChatAppService
	logger                  *zap.Logger
}

// NewSessionController 创建会话控制器
func NewSessionController(
	agentSessionAppService *appConv.AgentSessionAppService,
	conversationAppService *appConv.ConversationAppService,
	chatAppService *appChat.ChatAppService,
	logger *zap.Logger,
) *SessionController {
	return &SessionController{
		agentSessionAppService: agentSessionAppService,
		conversationAppService: conversationAppService,
		chatAppService:         chatAppService,
		logger:                 logger,
	}
}

// GetConversationMessages 获取会话中的消息列表（GET /agents/sessions/:sessionId/messages）
func (ctrl *SessionController) GetConversationMessages(c *gin.Context) {
	sessionID := c.Param("sessionId")
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.conversationAppService.GetConversationMessages(sessionID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetAgentSessionList 获取助理会话列表（GET /agents/sessions/:agentId）
func (ctrl *SessionController) GetAgentSessionList(c *gin.Context) {
	agentID := c.Param("agentId")
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.agentSessionAppService.GetAgentSessionList(userID, agentID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// CreateSession 创建会话（POST /agents/sessions/:agentId）
func (ctrl *SessionController) CreateSession(c *gin.Context) {
	agentID := c.Param("agentId")
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.agentSessionAppService.CreateSession(userID, agentID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// UpdateSession 更新会话（PUT /agents/sessions/:id）
func (ctrl *SessionController) UpdateSession(c *gin.Context) {
	id := c.Param("id")
	title := c.Query("title")
	userID := auth.GetCurrentUserID(c)

	if err := ctrl.agentSessionAppService.UpdateSession(id, userID, title); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// DeleteSession 删除会话（DELETE /agents/sessions/:id）
func (ctrl *SessionController) DeleteSession(c *gin.Context) {
	id := c.Param("id")
	userID := auth.GetCurrentUserID(c)

	if err := ctrl.agentSessionAppService.DeleteSession(id, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// Chat 发送消息（POST /agents/sessions/chat）
// 集成 ChatAppService 实现 SSE 流式聊天
func (ctrl *SessionController) Chat(c *gin.Context) {
	var req appChat.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorJSON(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := auth.GetCurrentUserID(c)

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("X-Accel-Buffering", "no")

	// 创建SSE传输
	writer := c.Writer
	transport := appChat.NewSSETransport(writer, ctrl.logger)

	// 启动流式聊天
	if err := ctrl.chatAppService.StreamChat(&req, userID, transport); err != nil {
		transport.SendError(err.Error())
		return
	}

	// 保持连接直到客户端断开
	<-c.Request.Context().Done()
	transport.Close()
}

// InterruptSession 中断对话会话（POST /agents/sessions/:sessionId/interrupt）
func (ctrl *SessionController) InterruptSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	stopped := ctrl.chatAppService.StopChat(sessionID)
	c.JSON(http.StatusOK, common.SuccessWithData(gin.H{
		"stopped": stopped,
	}))
}
