package portal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appTrace "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/trace"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// TraceController 执行链路追踪控制器
type TraceController struct {
	appService *appTrace.AppService
}

func NewTraceController(appService *appTrace.AppService) *TraceController {
	return &TraceController{appService: appService}
}

// GetExecutionHistory 分页查询用户执行历史（GET /traces/history）
func (ctrl *TraceController) GetExecutionHistory(c *gin.Context) {
	var req appTrace.QueryExecutionHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)
	dtos, total, err := ctrl.appService.GetExecutionHistory(&req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	page := 1
	pageSize := 15
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

// GetTraceDetail 获取单个追踪详情（GET /traces/:traceId）
func (ctrl *TraceController) GetTraceDetail(c *gin.Context) {
	traceID := c.Param("traceId")
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.GetTraceDetail(traceID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetExecutionDetails 获取执行详情列表（GET /traces/:traceId/details）
func (ctrl *TraceController) GetExecutionDetails(c *gin.Context) {
	traceID := c.Param("traceId")
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.GetExecutionDetails(traceID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetSessionExecutionHistory 查询会话执行记录（GET /traces/sessions/:sessionId）
func (ctrl *TraceController) GetSessionExecutionHistory(c *gin.Context) {
	sessionID := c.Param("sessionId")
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.GetSessionExecutionHistory(sessionID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetUserExecutionStatistics 获取用户执行统计（GET /traces/statistics）
func (ctrl *TraceController) GetUserExecutionStatistics(c *gin.Context) {
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.GetUserExecutionStatistics(userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}
