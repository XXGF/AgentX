package portal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appAuth "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// SsoController SSO控制器（对应 Java 的 SsoController）
type SsoController struct {
	ssoAppService *appAuth.SsoAppService
}

// NewSsoController 创建SSO控制器
func NewSsoController(ssoAppService *appAuth.SsoAppService) *SsoController {
	return &SsoController{ssoAppService: ssoAppService}
}

// GetSsoLoginUrl 获取SSO登录URL（GET /sso/:provider/login）
func (ctrl *SsoController) GetSsoLoginUrl(c *gin.Context) {
	provider := c.Param("provider")
	redirectUrl := c.Query("redirectUrl")

	loginUrl, err := ctrl.ssoAppService.GetSsoLoginUrl(provider, redirectUrl)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessWithData(gin.H{"loginUrl": loginUrl}))
}

// HandleSsoCallback SSO登录回调处理（GET /sso/:provider/callback）
func (ctrl *SsoController) HandleSsoCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")

	token, err := ctrl.ssoAppService.HandleSsoCallback(provider, code)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessWithMessage("登录成功", gin.H{"token": token}))
}
