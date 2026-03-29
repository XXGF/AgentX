package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	appAgent "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/agent"
	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/agent"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// AdminAgentController 管理员Agent控制器（对应 Java 的 AdminAgentController）
type AdminAgentController struct {
	agentAppService *appAgent.AppService
}

// NewAdminAgentController 创建管理员Agent控制器
func NewAdminAgentController(agentAppService *appAgent.AppService) *AdminAgentController {
	return &AdminAgentController{agentAppService: agentAppService}
}

// GetAgents 分页获取Agent列表（GET /admin/agents）
func (ctrl *AdminAgentController) GetAgents(c *gin.Context) {
	var req appAgent.QueryAgentRequest
	_ = c.ShouldBindQuery(&req)

	agents, total, err := ctrl.agentAppService.GetAgentsPage(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessWithData(gin.H{
		"records": agents,
		"total":   total,
		"page":    req.Page,
		"size":    req.PageSize,
	}))
}

// GetAgentStatistics 获取Agent统计信息（GET /admin/agents/statistics）
func (ctrl *AdminAgentController) GetAgentStatistics(c *gin.Context) {
	result, err := ctrl.agentAppService.GetAgentStatistics()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetVersions 获取版本列表（GET /admin/agents/versions）
func (ctrl *AdminAgentController) GetVersions(c *gin.Context) {
	statusStr := c.Query("status")
	agentID := c.Query("agentId")

	if agentID != "" {
		// 获取指定Agent的所有版本
		result, err := ctrl.agentAppService.GetAgentVersions(agentID, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, common.SuccessWithData(result))
		return
	}

	// 根据状态参数获取对应的版本列表
	var status *domain.PublishStatus
	if statusStr != "" {
		statusInt, err := strconv.Atoi(statusStr)
		if err == nil {
			s := domain.PublishStatus(statusInt)
			status = &s
		}
	}

	result, err := ctrl.agentAppService.GetVersionsByStatus(status)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// UpdateVersionStatus 更新版本状态（POST /admin/agents/versions/:versionId/status）
func (ctrl *AdminAgentController) UpdateVersionStatus(c *gin.Context) {
	versionID := c.Param("versionId")
	statusStr := c.Query("status")
	reason := c.Query("reason")

	statusInt, err := strconv.Atoi(statusStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("无效的状态码"))
		return
	}

	publishStatus := domain.PublishStatus(statusInt)

	// 如果是拒绝操作，需要检查原因
	if publishStatus == domain.PublishStatusRejected && reason == "" {
		c.JSON(http.StatusBadRequest, common.BadRequest("拒绝操作需要提供原因"))
		return
	}

	req := &appAgent.ReviewAgentVersionRequest{
		Status: publishStatus,
	}
	if publishStatus == domain.PublishStatusRejected {
		req.RejectReason = reason
	}

	result, err2 := ctrl.agentAppService.ReviewAgentVersion(versionID, req)
	if err2 != nil {
		_ = c.Error(err2)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}
