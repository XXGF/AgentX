package portal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appUser "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/user"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
	dto "github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/dto/user"
)

// UserController 用户控制器（对应 Java 的 PortalUserController）
type UserController struct {
	userAppService *appUser.AppService
}

// NewUserController 创建用户控制器
func NewUserController(userAppService *appUser.AppService) *UserController {
	return &UserController{userAppService: userAppService}
}

// GetUserInfo 获取用户信息（GET /users）
func (ctrl *UserController) GetUserInfo(c *gin.Context) {
	userID := auth.GetCurrentUserID(c)
	userDTO, err := ctrl.userAppService.GetUserInfo(userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessWithData(userDTO))
}

// UpdateUserInfo 修改用户信息（POST /users）
func (ctrl *UserController) UpdateUserInfo(c *gin.Context) {
	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}

	userID := auth.GetCurrentUserID(c)
	if err := ctrl.userAppService.UpdateUserInfo(req.Nickname, userID); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.Success())
}

// ChangePassword 修改密码（PUT /users/password）
func (ctrl *UserController) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}

	userID := auth.GetCurrentUserID(c)
	if err := ctrl.userAppService.ChangePassword(req.CurrentPassword, req.NewPassword, req.ConfirmPassword, userID); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.Success().SetMessage("密码修改成功"))
}
