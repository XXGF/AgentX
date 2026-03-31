package portal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appTask "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/task"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// TaskController 任务控制器
type TaskController struct {
	appService *appTask.AppService
}

func NewTaskController(appService *appTask.AppService) *TaskController {
	return &TaskController{appService: appService}
}

// GetSessionTasks 获取当前会话的任务（GET /tasks/session/:sessionId/latest）
func (ctrl *TaskController) GetSessionTasks(c *gin.Context) {
	sessionID := c.Param("sessionId")
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.GetCurrentSessionTask(sessionID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}
