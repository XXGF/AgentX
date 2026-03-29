package portal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appMemory "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/memory"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// MemoryController 记忆管理控制器
type MemoryController struct {
	appService *appMemory.AppService
}

func NewMemoryController(appService *appMemory.AppService) *MemoryController {
	return &MemoryController{appService: appService}
}

// ListMemories 分页列出当前用户的记忆（GET /portal/memory/items）
func (ctrl *MemoryController) ListMemories(c *gin.Context) {
	var req appMemory.QueryMemoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)
	dtos, total, err := ctrl.appService.ListUserMemories(userID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	page := 1
	pageSize := 20
	if req.Page != nil {
		page = *req.Page
	}
	if req.PageSize != nil {
		pageSize = *req.PageSize
	}
	c.JSON(http.StatusOK, common.SuccessWithData(gin.H{
		"records": dtos,
		"total":   total,
		"current": page,
		"size":    pageSize,
	}))
}

// CreateMemory 手动新增记忆（POST /portal/memory/items）
func (ctrl *MemoryController) CreateMemory(c *gin.Context) {
	var req appMemory.CreateMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)
	_, err := ctrl.appService.CreateMemory(userID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// DeleteMemory 归档（软删除）记忆（DELETE /portal/memory/items/:itemId）
func (ctrl *MemoryController) DeleteMemory(c *gin.Context) {
	itemID := c.Param("itemId")
	userID := auth.GetCurrentUserID(c)
	if err := ctrl.appService.DeleteMemory(userID, itemID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}
