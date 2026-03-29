package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appTool "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/tool"
	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/tool"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// AdminToolController 管理员工具控制器（对应 Java 的 AdminToolController）
type AdminToolController struct {
	toolAppService    *appTool.AppService
	toolDomainService *domain.ToolDomainService
}

func NewAdminToolController(toolAppService *appTool.AppService, toolDomainService *domain.ToolDomainService) *AdminToolController {
	return &AdminToolController{
		toolAppService:    toolAppService,
		toolDomainService: toolDomainService,
	}
}

// GetTools 分页获取工具列表（GET /admin/tools）
func (ctrl *AdminToolController) GetTools(c *gin.Context) {
	var req appTool.QueryToolRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}

	page := req.GetPageOrDefault()
	pageSize := req.GetPageSizeOrDefault()

	tools, total, err := ctrl.toolDomainService.GetTools(page, pageSize, req.Keyword, req.Status, req.IsOffice)
	if err != nil {
		_ = c.Error(err)
		return
	}

	dtos := appTool.ToolEntitiesToDTOs(tools)
	c.JSON(http.StatusOK, common.SuccessWithData(gin.H{
		"records": dtos,
		"total":   total,
		"current": page,
		"size":    pageSize,
	}))
}

// GetToolStatistics 获取工具统计信息（GET /admin/tools/statistics）
func (ctrl *AdminToolController) GetToolStatistics(c *gin.Context) {
	stats, err := ctrl.toolDomainService.GetToolStatistics()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(stats))
}

// CreateOfficialTool 创建官方工具（POST /admin/tools/official）
func (ctrl *AdminToolController) CreateOfficialTool(c *gin.Context) {
	var req appTool.CreateToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)
	isOffice := true
	req.IsGlobal = req.IsGlobal // 保留原值

	toolEntity := appTool.CreateToolRequestToEntity(&req, userID)
	toolEntity.IsOffice = &isOffice

	result, err := ctrl.toolDomainService.CreateTool(toolEntity)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result.Tool.ID))
}

// UpdateStatus 修改工具状态（POST /admin/tools/:toolId/status）
func (ctrl *AdminToolController) UpdateStatus(c *gin.Context) {
	toolID := c.Param("toolId")
	statusStr := c.Query("status")
	reason := c.Query("reason")

	status := domain.ToolStatus(statusStr)
	if status == domain.ToolStatusFailed && reason == "" {
		c.JSON(http.StatusBadRequest, common.BadRequest("拒绝操作需要提供原因"))
		return
	}

	tool, err := ctrl.toolDomainService.GetToolByID(toolID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if status == domain.ToolStatusApproved {
		_, err = ctrl.toolDomainService.UpdateApprovedToolStatus(toolID, status)
	} else if status == domain.ToolStatusFailed {
		_, err = ctrl.toolDomainService.UpdateFailedToolStatus(toolID, tool.Status, reason)
	} else {
		err = ctrl.toolDomainService.TransitionToStatus(tool, status)
	}

	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// UpdateGlobalStatus 更新工具全局状态（PUT /admin/tools/:toolId/global-status）
func (ctrl *AdminToolController) UpdateGlobalStatus(c *gin.Context) {
	toolID := c.Param("toolId")
	var req struct {
		IsGlobal *bool `json:"isGlobal" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}
	if err := ctrl.toolDomainService.UpdateToolGlobalStatus(toolID, *req.IsGlobal); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}
