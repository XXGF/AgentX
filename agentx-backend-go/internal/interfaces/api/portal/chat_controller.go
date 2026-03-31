package portal

import (
	"github.com/gin-gonic/gin"
	appChat "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/chat"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
	"go.uber.org/zap"
)

// ChatController 聊天控制器（对应 Java 的 ConversationAppService 中的 Chat 功能）
type ChatController struct {
	chatAppService *appChat.ChatAppService
	logger         *zap.Logger
}

func NewChatController(chatAppService *appChat.ChatAppService, logger *zap.Logger) *ChatController {
	return &ChatController{
		chatAppService: chatAppService,
		logger:         logger,
	}
}

// StreamChat 流式聊天（SSE）
// POST /chat/stream
func (ctrl *ChatController) StreamChat(c *gin.Context) {
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

// Chat 同步聊天
// POST /chat
func (ctrl *ChatController) Chat(c *gin.Context) {
	var req appChat.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorJSON(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := auth.GetCurrentUserID(c)

	response, err := ctrl.chatAppService.Chat(&req, userID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}

	common.SuccessJSON(c, response)
}

// StopChat 停止聊天
// POST /chat/stop
func (ctrl *ChatController) StopChat(c *gin.Context) {
	var req appChat.StopChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorJSON(c, 400, "参数错误: "+err.Error())
		return
	}

	stopped := ctrl.chatAppService.StopChat(req.SessionID)
	common.SuccessJSON(c, gin.H{
		"stopped": stopped,
	})
}

// WidgetChatController 小组件聊天控制器（对应 Java 的 WidgetChatController）
type WidgetChatController struct {
	chatAppService *appChat.ChatAppService
	logger         *zap.Logger
}

func NewWidgetChatController(chatAppService *appChat.ChatAppService, logger *zap.Logger) *WidgetChatController {
	return &WidgetChatController{
		chatAppService: chatAppService,
		logger:         logger,
	}
}

// WidgetStreamChat 小组件流式聊天（公开API，无需登录）
// POST /public/widget/chat/stream
func (ctrl *WidgetChatController) WidgetStreamChat(c *gin.Context) {
	// TODO: 实现小组件流式聊天（需要Widget配置验证）
	common.ErrorJSON(c, 501, "小组件聊天功能待实现")
}
