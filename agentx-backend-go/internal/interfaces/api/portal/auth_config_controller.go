package portal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appAuth "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// AuthConfigController 认证配置控制器（用户端）（对应 Java 的 AuthConfigController）
type AuthConfigController struct {
	authAppService *appAuth.AppService
}

// NewAuthConfigController 创建认证配置控制器
func NewAuthConfigController(authAppService *appAuth.AppService) *AuthConfigController {
	return &AuthConfigController{authAppService: authAppService}
}

// GetAuthConfig 获取可用的认证配置（GET /auth/config）
func (ctrl *AuthConfigController) GetAuthConfig(c *gin.Context) {
	config, err := ctrl.authAppService.GetAuthConfig()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessWithData(config))
}
