package portal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appConv "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/conversation"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// SessionController Agent会话控制器（对应 Java 的 PortalAgentSessionController）
type SessionController struct {
	agentSessionAppService  *appConv.AgentSessionAppService
	conversationAppService  *appConv.ConversationAppService
}

// NewSessionController 创建会话控制器
func NewSessionController(
	agentSessionAppService *appConv.AgentSessionAppService,
	conversationAppService *appConv.ConversationAppService,
) *SessionController {
	return &SessionController{
		agentSessionAppService: agentSessionAppService,
		conversationAppService: conversationAppService,
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
// TODO: 涉及复杂的SSE流式响应、Agent工作流、MCP工具调用等，后续迁移
func (ctrl *SessionController) Chat(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, common.BadRequest("聊天功能正在迁移中"))
}

// InterruptSession 中断对话会话（POST /agents/sessions/:sessionId/interrupt）
// TODO: 涉及ChatSessionManager，后续迁移
func (ctrl *SessionController) InterruptSession(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, common.BadRequest("中断功能正在迁移中"))
}
