package portal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appProduct "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/product"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// ProductController 商品控制器
type ProductController struct {
	appService *appProduct.AppService
}

func NewProductController(appService *appProduct.AppService) *ProductController {
	return &ProductController{appService: appService}
}

// GetProductByID 根据ID获取商品详情（GET /products/:productId）
func (ctrl *ProductController) GetProductByID(c *gin.Context) {
	productID := c.Param("productId")
	result, err := ctrl.appService.GetProductByID(productID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetProductByBusinessKey 根据业务标识获取商品（GET /products/business）
func (ctrl *ProductController) GetProductByBusinessKey(c *gin.Context) {
	billingType := c.Query("type")
	serviceID := c.Query("serviceId")
	result, err := ctrl.appService.GetProductByBusinessKey(billingType, serviceID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetActiveProducts 获取指定类型的活跃商品列表（GET /products/active）
func (ctrl *ProductController) GetActiveProducts(c *gin.Context) {
	billingType := c.Query("type")
	result, err := ctrl.appService.GetActiveProducts(billingType)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// IsProductActive 检查商品是否存在且激活（GET /products/business/active）
func (ctrl *ProductController) IsProductActive(c *gin.Context) {
	billingType := c.Query("type")
	serviceID := c.Query("serviceId")
	result, err := ctrl.appService.IsProductActive(billingType, serviceID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}
