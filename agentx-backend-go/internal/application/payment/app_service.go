package payment

import (
	"fmt"
	"math/rand"
	"time"

	domainOrder "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/order"
	domainUser "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/user"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	infraPayment "github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/payment"
)

// RechargeRequest 充值请求
type RechargeRequest struct {
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	PaymentPlatform string  `json:"paymentPlatform" binding:"required"`
	PaymentType     string  `json:"paymentType" binding:"required"`
	Remark          string  `json:"remark"`
}

// PaymentResponseDTO 支付响应
type PaymentResponseDTO struct {
	OrderID   string  `json:"orderId"`
	OrderNo   string  `json:"orderNo"`
	Amount    float64 `json:"amount"`
	PayURL    string  `json:"payUrl,omitempty"`
	QRCodeURL string  `json:"qrCodeUrl,omitempty"`
}

// OrderStatusResponseDTO 订单状态响应
type OrderStatusResponseDTO struct {
	OrderID         string  `json:"orderId"`
	OrderNo         string  `json:"orderNo"`
	Status          string  `json:"status"`
	PaymentPlatform string  `json:"paymentPlatform"`
	PaymentType     string  `json:"paymentType"`
	Amount          float64 `json:"amount"`
	Title           string  `json:"title"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
	ExpiredAt       string  `json:"expiredAt,omitempty"`
}

// PaymentMethodDTO 支付方法DTO
type PaymentMethodDTO struct {
	PlatformCode string           `json:"platformCode"`
	PlatformName string           `json:"platformName"`
	Available    bool             `json:"available"`
	Description  string           `json:"description"`
	PaymentTypes []PaymentTypeDTO `json:"paymentTypes"`
}

// PaymentTypeDTO 支付类型DTO
type PaymentTypeDTO struct {
	TypeCode    string `json:"typeCode"`
	TypeName    string `json:"typeName"`
	Default     bool   `json:"default"`
	Description string `json:"description"`
}

// AppService 支付应用服务
type AppService struct {
	orderDomainService   *domainOrder.DomainService
	accountDomainService *domainUser.AccountDomainService
	paymentFactory       *infraPayment.PaymentProviderFactory
}

func NewAppService(
	orderDomainService *domainOrder.DomainService,
	accountDomainService *domainUser.AccountDomainService,
	paymentFactory *infraPayment.PaymentProviderFactory,
) *AppService {
	return &AppService{
		orderDomainService:   orderDomainService,
		accountDomainService: accountDomainService,
		paymentFactory:       paymentFactory,
	}
}

// CreateRechargePayment 创建充值支付
func (s *AppService) CreateRechargePayment(userID string, req *RechargeRequest) (*PaymentResponseDTO, error) {
	if req.Amount <= 0 {
		return nil, exception.NewBusinessException("充值金额必须大于0")
	}

	// 创建充值订单
	order := &domainOrder.OrderEntity{
		UserID:          userID,
		OrderNo:         generateOrderNo(),
		OrderType:       domainOrder.OrderTypeRecharge,
		Title:           "账户充值",
		Description:     fmt.Sprintf("账户余额充值 ¥%.2f", req.Amount),
		Amount:          req.Amount,
		Currency:        "CNY",
		Status:          domainOrder.OrderStatusPending,
		PaymentPlatform: domainOrder.PaymentPlatform(req.PaymentPlatform),
		PaymentType:     domainOrder.PaymentType(req.PaymentType),
	}

	createdOrder, err := s.orderDomainService.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	// 调用支付提供商创建支付
	provider, err := s.paymentFactory.GetProvider(req.PaymentPlatform)
	if err != nil {
		// 如果找不到对应的支付提供商，使用默认的
		provider = s.paymentFactory.GetDefaultProvider()
	}

	payResult, err := provider.CreatePayment(createdOrder.OrderNo, createdOrder.Amount, createdOrder.Title)
	if err != nil {
		return nil, exception.NewBusinessExceptionWithCause("创建支付失败", err)
	}

	return &PaymentResponseDTO{
		OrderID: createdOrder.ID,
		OrderNo: createdOrder.OrderNo,
		Amount:  createdOrder.Amount,
		PayURL:  payResult.PaymentURL,
		QRCodeURL: payResult.QRCode,
	}, nil
}

// QueryOrderStatus 查询订单状态
func (s *AppService) QueryOrderStatus(orderNo string) (*OrderStatusResponseDTO, error) {
	order, err := s.orderDomainService.FindOrderByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, exception.NewBusinessException("订单不存在: " + orderNo)
	}

	return &OrderStatusResponseDTO{
		OrderID:         order.ID,
		OrderNo:         order.OrderNo,
		Status:          string(order.Status),
		PaymentPlatform: string(order.PaymentPlatform),
		PaymentType:     string(order.PaymentType),
		Amount:          order.Amount,
		Title:           order.Title,
		CreatedAt:       order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       order.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// GetAvailablePaymentMethods 获取可用的支付方法列表
func (s *AppService) GetAvailablePaymentMethods() []*PaymentMethodDTO {
	return []*PaymentMethodDTO{
		{
			PlatformCode: "alipay",
			PlatformName: "支付宝",
			Available:    false,
			Description:  "支持扫码支付",
			PaymentTypes: []PaymentTypeDTO{
				{TypeCode: "QR_CODE", TypeName: "扫码支付", Description: "扫描二维码完成支付"},
			},
		},
		{
			PlatformCode: "wechat",
			PlatformName: "微信支付",
			Available:    false,
			Description:  "支持扫码支付",
			PaymentTypes: []PaymentTypeDTO{
				{TypeCode: "QR_CODE", TypeName: "扫码支付", Description: "扫描二维码完成支付"},
			},
		},
	}
}

func generateOrderNo() string {
	return fmt.Sprintf("RCH%d%04d", time.Now().UnixMilli(), rand.Intn(10000))
}
