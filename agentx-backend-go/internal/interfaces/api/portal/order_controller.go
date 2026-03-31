package portal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appOrder "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/order"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// OrderController 订单控制器
type OrderController struct {
	appService *appOrder.AppService
}

func NewOrderController(appService *appOrder.AppService) *OrderController {
	return &OrderController{appService: appService}
}

// GetUserOrders 获取当前用户的已支付订单列表（GET /orders）
func (ctrl *OrderController) GetUserOrders(c *gin.Context) {
	var req appOrder.QueryUserOrderRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)
	dtos, total, err := ctrl.appService.GetUserPaidOrders(&req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	c.JSON(http.StatusOK, common.SuccessWithData(gin.H{
		"records": dtos,
		"total":   total,
		"current": page,
		"size":    pageSize,
	}))
}

// GetOrderDetail 获取订单详情（GET /orders/:orderId）
func (ctrl *OrderController) GetOrderDetail(c *gin.Context) {
	orderID := c.Param("orderId")
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.GetUserOrderDetail(orderID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}
