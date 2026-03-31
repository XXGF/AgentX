package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	appOrder "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/order"
	appProduct "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/product"
	appRule "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/rule"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// ---- AdminProductController 管理员商品控制器 ----

type AdminProductController struct {
	appService *appProduct.AppService
}

func NewAdminProductController(appService *appProduct.AppService) *AdminProductController {
	return &AdminProductController{appService: appService}
}

// GetProducts 分页获取商品列表
// GET /admin/products
func (ctrl *AdminProductController) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "15"))
	keyword := c.Query("keyword")
	products, total, err := ctrl.appService.GetProductsPaged(page, pageSize, keyword)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, gin.H{"records": products, "total": total})
}

// GetAllProducts 获取所有商品列表
// GET /admin/products/all
func (ctrl *AdminProductController) GetAllProducts(c *gin.Context) {
	products, err := ctrl.appService.GetAllProducts()
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, products)
}

// GetProductByID 根据ID获取商品详情
// GET /admin/products/:productId
func (ctrl *AdminProductController) GetProductByID(c *gin.Context) {
	productID := c.Param("productId")
	product, err := ctrl.appService.GetProductByID(productID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, product)
}

// DeleteProduct 删除商品
// DELETE /admin/products/:productId
func (ctrl *AdminProductController) DeleteProduct(c *gin.Context) {
	productID := c.Param("productId")
	if err := ctrl.appService.DeleteProduct(productID); err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, nil)
}

// ---- AdminOrderController 管理员订单控制器 ----

type AdminOrderController struct {
	appService *appOrder.AppService
}

func NewAdminOrderController(appService *appOrder.AppService) *AdminOrderController {
	return &AdminOrderController{appService: appService}
}

// GetAllOrders 分页获取所有订单列表
// GET /admin/orders
func (ctrl *AdminOrderController) GetAllOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "15"))
	keyword := c.Query("keyword")
	orders, total, err := ctrl.appService.GetAllOrdersPaged(page, pageSize, keyword)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, gin.H{"records": orders, "total": total})
}

// GetOrderDetail 获取订单详情
// GET /admin/orders/:orderId
func (ctrl *AdminOrderController) GetOrderDetail(c *gin.Context) {
	orderID := c.Param("orderId")
	order, err := ctrl.appService.GetOrderDetail(orderID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, order)
}

// ---- AdminRuleController 管理员规则控制器 ----

type AdminRuleController struct {
	appService *appRule.AppService
}

func NewAdminRuleController(appService *appRule.AppService) *AdminRuleController {
	return &AdminRuleController{appService: appService}
}

// CreateRule 创建计费规则
// POST /admin/rules
func (ctrl *AdminRuleController) CreateRule(c *gin.Context) {
	var req appRule.CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorJSON(c, 400, "参数错误: "+err.Error())
		return
	}
	rule, err := ctrl.appService.CreateRule(&req)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, rule)
}

// UpdateRule 更新计费规则
// PUT /admin/rules/:ruleId
func (ctrl *AdminRuleController) UpdateRule(c *gin.Context) {
	ruleID := c.Param("ruleId")
	var req appRule.UpdateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorJSON(c, 400, "参数错误: "+err.Error())
		return
	}
	rule, err := ctrl.appService.UpdateRule(&req, ruleID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, rule)
}

// GetRuleByID 根据ID获取规则
// GET /admin/rules/:ruleId
func (ctrl *AdminRuleController) GetRuleByID(c *gin.Context) {
	ruleID := c.Param("ruleId")
	rule, err := ctrl.appService.GetRuleByID(ruleID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, rule)
}

// GetRules 分页查询规则
// GET /admin/rules
func (ctrl *AdminRuleController) GetRules(c *gin.Context) {
	var req appRule.QueryRuleRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ErrorJSON(c, 400, "参数错误: "+err.Error())
		return
	}
	rules, total, err := ctrl.appService.GetRules(&req)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, gin.H{"records": rules, "total": total})
}

// GetAllRules 获取所有规则
// GET /admin/rules/all
func (ctrl *AdminRuleController) GetAllRules(c *gin.Context) {
	rules, err := ctrl.appService.GetAllRules()
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, rules)
}

// DeleteRule 删除规则
// DELETE /admin/rules/:ruleId
func (ctrl *AdminRuleController) DeleteRule(c *gin.Context) {
	ruleID := c.Param("ruleId")
	if err := ctrl.appService.DeleteRule(ruleID); err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, nil)
}
