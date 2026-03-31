package portal

import (
	"github.com/gin-gonic/gin"
	appAccount "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/account"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// AccountController 账户控制器（对应 Java 的 AccountController）
type AccountController struct {
	appService *appAccount.AppService
}

func NewAccountController(appService *appAccount.AppService) *AccountController {
	return &AccountController{appService: appService}
}

// GetCurrentUserAccount 获取当前用户账户信息
// GET /accounts/current
func (ctrl *AccountController) GetCurrentUserAccount(c *gin.Context) {
	userID := auth.GetCurrentUserID(c)
	account, err := ctrl.appService.GetUserAccount(userID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, account)
}
