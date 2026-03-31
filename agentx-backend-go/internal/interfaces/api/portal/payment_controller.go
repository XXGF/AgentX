package portal

import (
	"github.com/gin-gonic/gin"
	appPayment "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/payment"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// PaymentController 支付控制器（对应 Java 的 PaymentController）
type PaymentController struct {
	appService *appPayment.AppService
}

func NewPaymentController(appService *appPayment.AppService) *PaymentController {
	return &PaymentController{appService: appService}
}

// CreateRechargePayment 创建充值支付
// POST /payments/recharge
func (ctrl *PaymentController) CreateRechargePayment(c *gin.Context) {
	var req appPayment.RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorJSON(c, 400, "参数错误: "+err.Error())
		return
	}
	userID := auth.GetCurrentUserID(c)
	resp, err := ctrl.appService.CreateRechargePayment(userID, &req)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, resp)
}

// QueryOrderStatus 查询订单状态
// GET /payments/orders/:orderNo/status
func (ctrl *PaymentController) QueryOrderStatus(c *gin.Context) {
	orderNo := c.Param("orderNo")
	resp, err := ctrl.appService.QueryOrderStatus(orderNo)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, resp)
}

// GetAvailablePaymentMethods 获取可用的支付方法列表
// GET /payments/methods
func (ctrl *PaymentController) GetAvailablePaymentMethods(c *gin.Context) {
	methods := ctrl.appService.GetAvailablePaymentMethods()
	common.SuccessJSON(c, methods)
}
