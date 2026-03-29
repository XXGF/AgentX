package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appAuth "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// AdminAuthSettingController 管理员认证配置控制器（对应 Java 的 AdminAuthSettingController）
type AdminAuthSettingController struct {
	authAppService *appAuth.AppService
}

// NewAdminAuthSettingController 创建管理员认证配置控制器
func NewAdminAuthSettingController(authAppService *appAuth.AppService) *AdminAuthSettingController {
	return &AdminAuthSettingController{authAppService: authAppService}
}

// GetAllAuthSettings 获取所有认证配置（GET /admin/auth-settings）
func (ctrl *AdminAuthSettingController) GetAllAuthSettings(c *gin.Context) {
	settings, err := ctrl.authAppService.GetAllAuthSettings()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessWithData(settings))
}

// GetAuthSettingByID 根据ID获取认证配置（GET /admin/auth-settings/:id）
func (ctrl *AdminAuthSettingController) GetAuthSettingByID(c *gin.Context) {
	id := c.Param("id")
	setting, err := ctrl.authAppService.GetAuthSettingByID(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessWithData(setting))
}

// ToggleAuthSetting 切换认证配置启用状态（PUT /admin/auth-settings/:id/toggle）
func (ctrl *AdminAuthSettingController) ToggleAuthSetting(c *gin.Context) {
	id := c.Param("id")
	setting, err := ctrl.authAppService.ToggleAuthSetting(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessWithData(setting))
}

// UpdateAuthSetting 更新认证配置（PUT /admin/auth-settings/:id）
func (ctrl *AdminAuthSettingController) UpdateAuthSetting(c *gin.Context) {
	id := c.Param("id")
	var req appAuth.UpdateAuthSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}

	setting, err := ctrl.authAppService.UpdateAuthSetting(id, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessWithData(setting))
}

// DeleteAuthSetting 删除认证配置（DELETE /admin/auth-settings/:id）
func (ctrl *AdminAuthSettingController) DeleteAuthSetting(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.authAppService.DeleteAuthSetting(id); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.Success())
}
