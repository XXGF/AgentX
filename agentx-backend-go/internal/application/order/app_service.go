package order

import (
	"time"

	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/order"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// OrderDTO 订单DTO
type OrderDTO struct {
	ID              string  `json:"id"`
	UserID          string  `json:"userId"`
	UserNickname    string  `json:"userNickname,omitempty"`
	OrderNo         string  `json:"orderNo"`
	OrderType       string  `json:"orderType"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	Status          string  `json:"status"`
	ExpiredAt       *time.Time `json:"expiredAt"`
	PaidAt          *time.Time `json:"paidAt"`
	CancelledAt     *time.Time `json:"cancelledAt"`
	RefundedAt      *time.Time `json:"refundedAt"`
	RefundAmount    float64 `json:"refundAmount"`
	PaymentPlatform string  `json:"paymentPlatform"`
	PaymentType     string  `json:"paymentType"`
	ProviderOrderID string  `json:"providerOrderId"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// QueryUserOrderRequest 查询用户订单请求
type QueryUserOrderRequest struct {
	Page      int    `form:"page" json:"page"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
	OrderType string `form:"orderType" json:"orderType"`
}

func entityToDTO(e *domain.OrderEntity) *OrderDTO {
	if e == nil {
		return nil
	}
	return &OrderDTO{
		ID:              e.ID,
		UserID:          e.UserID,
		OrderNo:         e.OrderNo,
		OrderType:       string(e.OrderType),
		Title:           e.Title,
		Description:     e.Description,
		Amount:          e.Amount,
		Currency:        e.Currency,
		Status:          string(e.Status),
		ExpiredAt:       e.ExpiredAt,
		PaidAt:          e.PaidAt,
		CancelledAt:     e.CancelledAt,
		RefundedAt:      e.RefundedAt,
		RefundAmount:    e.RefundAmount,
		PaymentPlatform: string(e.PaymentPlatform),
		PaymentType:     string(e.PaymentType),
		ProviderOrderID: e.ProviderOrderID,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

func entitiesToDTOs(entities []domain.OrderEntity) []*OrderDTO {
	dtos := make([]*OrderDTO, len(entities))
	for i, e := range entities {
		dtos[i] = entityToDTO(&e)
	}
	return dtos
}

// AppService 订单应用服务
type AppService struct {
	orderDomainService *domain.DomainService
}

func NewAppService(orderDomainService *domain.DomainService) *AppService {
	return &AppService{orderDomainService: orderDomainService}
}

// GetUserPaidOrders 获取用户已支付订单列表（分页）
func (s *AppService) GetUserPaidOrders(req *QueryUserOrderRequest, userID string) ([]*OrderDTO, int64, error) {
	if userID == "" {
		return nil, 0, exception.NewBusinessException("用户ID不能为空")
	}
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var orderType *domain.OrderType
	if req.OrderType != "" {
		ot := domain.OrderType(req.OrderType)
		orderType = &ot
	}

	// 用户只能查看已支付的订单
	status := domain.OrderStatusPaid

	entities, total, err := s.orderDomainService.GetOrdersByUserIDPaged(userID, page, pageSize, orderType, &status)
	if err != nil {
		return nil, 0, err
	}
	return entitiesToDTOs(entities), total, nil
}

// GetUserOrderDetail 获取用户订单详情
func (s *AppService) GetUserOrderDetail(orderID, userID string) (*OrderDTO, error) {
	if orderID == "" {
		return nil, exception.NewBusinessException("订单ID不能为空")
	}
	order, err := s.orderDomainService.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		return nil, exception.NewBusinessException("无权限访问此订单")
	}
	if order.Status != domain.OrderStatusPaid {
		return nil, exception.NewBusinessException("订单状态不允许查看")
	}
	return entityToDTO(order), nil
}

// GetAllOrdersPaged 管理员分页获取所有订单
func (s *AppService) GetAllOrdersPaged(page, pageSize int, keyword string) ([]*OrderDTO, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 15
	}
	entities, total, err := s.orderDomainService.OrderRepo().FindAllPaged(page, pageSize, keyword)
	if err != nil {
		return nil, 0, err
	}
	return entitiesToDTOs(entities), total, nil
}

// GetOrderDetail 管理员获取订单详情
func (s *AppService) GetOrderDetail(orderID string) (*OrderDTO, error) {
	if orderID == "" {
		return nil, exception.NewBusinessException("订单ID不能为空")
	}
	order, err := s.orderDomainService.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	return entityToDTO(order), nil
}
