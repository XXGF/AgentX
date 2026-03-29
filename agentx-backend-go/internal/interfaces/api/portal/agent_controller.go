package portal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appAgent "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/agent"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// AgentController Agent控制器（对应 Java 的 PortalAgentController）
type AgentController struct {
	agentAppService *appAgent.AppService
}

// NewAgentController 创建Agent控制器
func NewAgentController(agentAppService *appAgent.AppService) *AgentController {
	return &AgentController{agentAppService: agentAppService}
}

// CreateAgent 创建新Agent（POST /agents）
func (ctrl *AgentController) CreateAgent(c *gin.Context) {
	var req appAgent.CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.agentAppService.CreateAgent(&req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetAgent 获取Agent详情（GET /agents/:agentId）
func (ctrl *AgentController) GetAgent(c *gin.Context) {
	agentID := c.Param("agentId")
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.agentAppService.GetAgent(agentID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetUserAgents 获取用户的Agent列表（GET /agents/user）
func (ctrl *AgentController) GetUserAgents(c *gin.Context) {
	userID := auth.GetCurrentUserID(c)
	var req appAgent.SearchAgentsRequest
	_ = c.ShouldBindQuery(&req)

	result, err := ctrl.agentAppService.GetUserAgents(userID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetPublishedAgents 获取已上架的Agent列表（GET /agents/published）
func (ctrl *AgentController) GetPublishedAgents(c *gin.Context) {
	userID := auth.GetCurrentUserID(c)
	var req appAgent.SearchAgentsRequest
	_ = c.ShouldBindQuery(&req)

	result, err := ctrl.agentAppService.GetPublishedAgentsByName(&req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// UpdateAgent 更新Agent信息（PUT /agents/:agentId）
func (ctrl *AgentController) UpdateAgent(c *gin.Context) {
	agentID := c.Param("agentId")
	var req appAgent.UpdateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}
	req.ID = agentID
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.agentAppService.UpdateAgent(&req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// ToggleAgentStatus 切换Agent的启用/禁用状态（PUT /agents/:agentId/toggle-status）
func (ctrl *AgentController) ToggleAgentStatus(c *gin.Context) {
	agentID := c.Param("agentId")

	result, err := ctrl.agentAppService.ToggleAgentStatus(agentID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// DeleteAgent 删除Agent（DELETE /agents/:agentId）
func (ctrl *AgentController) DeleteAgent(c *gin.Context) {
	agentID := c.Param("agentId")
	userID := auth.GetCurrentUserID(c)

	if err := ctrl.agentAppService.DeleteAgent(agentID, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// PublishAgentVersion 发布Agent版本（POST /agents/:agentId/publish）
func (ctrl *AgentController) PublishAgentVersion(c *gin.Context) {
	agentID := c.Param("agentId")
	var req appAgent.PublishAgentVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.agentAppService.PublishAgentVersion(agentID, &req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetAgentVersions 获取Agent的所有版本（GET /agents/:agentId/versions）
func (ctrl *AgentController) GetAgentVersions(c *gin.Context) {
	agentID := c.Param("agentId")
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.agentAppService.GetAgentVersions(agentID, &userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetAgentVersion 获取Agent的特定版本（GET /agents/:agentId/versions/:versionNumber）
func (ctrl *AgentController) GetAgentVersion(c *gin.Context) {
	agentID := c.Param("agentId")
	versionNumber := c.Param("versionNumber")

	result, err := ctrl.agentAppService.GetAgentVersion(agentID, versionNumber)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetLatestAgentVersion 获取Agent的最新版本（GET /agents/:agentId/versions/latest）
func (ctrl *AgentController) GetLatestAgentVersion(c *gin.Context) {
	agentID := c.Param("agentId")

	result, err := ctrl.agentAppService.GetLatestAgentVersion(agentID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}
