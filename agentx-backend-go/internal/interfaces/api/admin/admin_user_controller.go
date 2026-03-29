package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appUser "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/user"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
	dto "github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/dto/user"
)

// AdminUserController 管理员用户管理控制器（对应 Java 的 AdminUserController）
type AdminUserController struct {
	userAppService *appUser.AppService
}

// NewAdminUserController 创建管理员用户控制器
func NewAdminUserController(userAppService *appUser.AppService) *AdminUserController {
	return &AdminUserController{userAppService: userAppService}
}

// GetUsers 分页获取用户列表（GET /admin/users）
func (ctrl *AdminUserController) GetUsers(c *gin.Context) {
	var req dto.QueryUserRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}

	result, err := ctrl.userAppService.GetUsers(req.Keyword, req.GetPage(), req.GetPageSize())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessWithData(result))
}
