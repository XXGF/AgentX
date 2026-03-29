package portal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appScheduledTask "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/scheduledtask"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// ScheduledTaskController 定时任务控制器
type ScheduledTaskController struct {
	appService *appScheduledTask.AppService
}

func NewScheduledTaskController(appService *appScheduledTask.AppService) *ScheduledTaskController {
	return &ScheduledTaskController{appService: appService}
}

// CreateScheduledTask 创建定时任务（POST /scheduled-tasks）
func (ctrl *ScheduledTaskController) CreateScheduledTask(c *gin.Context) {
	var req appScheduledTask.CreateScheduledTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.CreateScheduledTask(&req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// UpdateScheduledTask 更新定时任务（PUT /scheduled-tasks/:taskId）
func (ctrl *ScheduledTaskController) UpdateScheduledTask(c *gin.Context) {
	taskID := c.Param("taskId")
	var req appScheduledTask.UpdateScheduledTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}
	req.ID = taskID
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.UpdateScheduledTask(&req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// DeleteScheduledTask 删除定时任务（DELETE /scheduled-tasks/:taskId）
func (ctrl *ScheduledTaskController) DeleteScheduledTask(c *gin.Context) {
	taskID := c.Param("taskId")
	userID := auth.GetCurrentUserID(c)
	if err := ctrl.appService.DeleteTask(taskID, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// GetScheduledTasks 获取定时任务列表（GET /scheduled-tasks）
func (ctrl *ScheduledTaskController) GetScheduledTasks(c *gin.Context) {
	userID := auth.GetCurrentUserID(c)
	sessionID := c.Query("sessionId")
	agentID := c.Query("agentId")

	var result []*appScheduledTask.ScheduledTaskDTO
	var err error

	if agentID != "" {
		result, err = ctrl.appService.GetTasksByAgentID(agentID, userID)
	} else if sessionID != "" {
		result, err = ctrl.appService.GetTasksBySessionID(sessionID, userID)
	} else {
		result, err = ctrl.appService.GetUserTasks(userID)
	}

	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetScheduledTasksByAgent 根据Agent ID获取定时任务列表（GET /scheduled-tasks/agent/:agentId）
func (ctrl *ScheduledTaskController) GetScheduledTasksByAgent(c *gin.Context) {
	agentID := c.Param("agentId")
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.GetTasksByAgentID(agentID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetScheduledTask 获取单个定时任务详情（GET /scheduled-tasks/:taskId）
func (ctrl *ScheduledTaskController) GetScheduledTask(c *gin.Context) {
	taskID := c.Param("taskId")
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.GetTask(taskID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// PauseTask 暂停定时任务（POST /scheduled-tasks/:taskId/pause）
func (ctrl *ScheduledTaskController) PauseTask(c *gin.Context) {
	taskID := c.Param("taskId")
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.PauseTask(taskID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// ResumeTask 恢复定时任务（POST /scheduled-tasks/:taskId/resume）
func (ctrl *ScheduledTaskController) ResumeTask(c *gin.Context) {
	taskID := c.Param("taskId")
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.ResumeTask(taskID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}
