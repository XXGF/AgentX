package portal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appTool "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/tool"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// ToolController 工具控制器（对应 Java 的 PortalToolController）
type ToolController struct {
	toolAppService *appTool.AppService
}

func NewToolController(toolAppService *appTool.AppService) *ToolController {
	return &ToolController{toolAppService: toolAppService}
}

// CreateTool 上传工具（POST /tools）
func (ctrl *ToolController) CreateTool(c *gin.Context) {
	var req appTool.CreateToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.toolAppService.UploadTool(&req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetToolDetail 获取工具详情（GET /tools/:toolId）
func (ctrl *ToolController) GetToolDetail(c *gin.Context) {
	toolID := c.Param("toolId")
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.toolAppService.GetToolDetail(toolID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetUserTools 获取用户工具列表（GET /tools/user）
func (ctrl *ToolController) GetUserTools(c *gin.Context) {
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.toolAppService.GetUserTools(userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// UpdateTool 编辑工具（PUT /tools/:toolId）
func (ctrl *ToolController) UpdateTool(c *gin.Context) {
	toolID := c.Param("toolId")
	var req appTool.UpdateToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.toolAppService.UpdateTool(toolID, &req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// DeleteTool 删除工具（DELETE /tools/:toolId）
func (ctrl *ToolController) DeleteTool(c *gin.Context) {
	toolID := c.Param("toolId")
	userID := auth.GetCurrentUserID(c)
	if err := ctrl.toolAppService.DeleteTool(toolID, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// MarketTool 上架工具（POST /tools/market）
func (ctrl *ToolController) MarketTool(c *gin.Context) {
	var req appTool.MarketToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)
	if err := ctrl.toolAppService.MarketTool(&req, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithMessage("上架成功", nil))
}

// Market 工具市场列表（GET /tools/market）
func (ctrl *ToolController) Market(c *gin.Context) {
	var req appTool.QueryToolRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}
	dtos, total, err := ctrl.toolAppService.MarketTools(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(gin.H{
		"records": dtos,
		"total":   total,
		"current": req.GetPageOrDefault(),
		"size":    req.GetPageSizeOrDefault(),
	}))
}

// GetToolVersionDetail 获取工具版本详情（GET /tools/market/:toolId/:version）
func (ctrl *ToolController) GetToolVersionDetail(c *gin.Context) {
	toolID := c.Param("toolId")
	version := c.Param("version")
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.toolAppService.GetToolVersionDetail(toolID, version, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// InstallTool 安装工具（POST /tools/install/:toolId/:version）
func (ctrl *ToolController) InstallTool(c *gin.Context) {
	toolID := c.Param("toolId")
	version := c.Param("version")
	userID := auth.GetCurrentUserID(c)
	if err := ctrl.toolAppService.InstallTool(toolID, version, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithMessage("安装成功", nil))
}

// UninstallTool 卸载工具（POST /tools/uninstall/:toolId）
func (ctrl *ToolController) UninstallTool(c *gin.Context) {
	toolID := c.Param("toolId")
	userID := auth.GetCurrentUserID(c)
	if err := ctrl.toolAppService.UninstallTool(toolID, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithMessage("卸载成功", nil))
}

// GetInstalledTools 获取已安装的工具列表（GET /tools/installed）
func (ctrl *ToolController) GetInstalledTools(c *gin.Context) {
	var req appTool.QueryToolRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)
	dtos, total, err := ctrl.toolAppService.GetInstalledTools(userID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(gin.H{
		"records": dtos,
		"total":   total,
		"current": req.GetPageOrDefault(),
		"size":    req.GetPageSizeOrDefault(),
	}))
}

// GetToolVersions 获取工具已发布的所有版本（GET /tools/market/:toolId/versions）
func (ctrl *ToolController) GetToolVersions(c *gin.Context) {
	toolID := c.Param("toolId")
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.toolAppService.GetToolVersions(toolID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetRecommendTools 推荐工具（GET /tools/recommend）
func (ctrl *ToolController) GetRecommendTools(c *gin.Context) {
	result, err := ctrl.toolAppService.GetRecommendTools()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// UpdateToolVersionStatus 修改工具版本发布状态（POST /tools/user/:toolId/:version/status）
func (ctrl *ToolController) UpdateToolVersionStatus(c *gin.Context) {
	toolID := c.Param("toolId")
	version := c.Param("version")
	publishStatus := c.Query("publishStatus") == "true"
	userID := auth.GetCurrentUserID(c)
	if err := ctrl.toolAppService.UpdateUserToolVersionStatus(toolID, version, publishStatus, userID); err != nil {
		_ = c.Error(err)
		return
	}
	msg := "发布成功"
	if !publishStatus {
		msg = "下架成功"
	}
	c.JSON(http.StatusOK, common.SuccessWithMessage(msg, nil))
}

// GetLatestToolVersion 获取工具最新版本（GET /tools/:toolId/latest）
func (ctrl *ToolController) GetLatestToolVersion(c *gin.Context) {
	toolID := c.Param("toolId")
	result, err := ctrl.toolAppService.GetLatestToolVersion(toolID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(gin.H{"version": result.Version}))
}
