package order

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	"gorm.io/gorm"
)

// ---- 枚举 ----

// OrderStatus 订单状态枚举
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"   // 待支付
	OrderStatusPaid      OrderStatus = "PAID"       // 已支付
	OrderStatusCancelled OrderStatus = "CANCELLED"  // 已取消
	OrderStatusRefunded  OrderStatus = "REFUNDED"   // 已退款
	OrderStatusExpired   OrderStatus = "EXPIRED"    // 已过期
)

// IsFinished 检查订单是否为结束状态
func (s OrderStatus) IsFinished() bool {
	return s == OrderStatusPaid || s == OrderStatusCancelled || s == OrderStatusRefunded || s == OrderStatusExpired
}

// CanPay 检查订单是否可以支付
func (s OrderStatus) CanPay() bool { return s == OrderStatusPending }

// CanCancel 检查订单是否可以取消
func (s OrderStatus) CanCancel() bool { return s == OrderStatusPending }

// CanRefund 检查订单是否可以退款
func (s OrderStatus) CanRefund() bool { return s == OrderStatusPaid }

// OrderType 订单类型枚举
type OrderType string

const (
	OrderTypeRecharge     OrderType = "RECHARGE"     // 充值订单
	OrderTypePurchase     OrderType = "PURCHASE"     // 购买订单
	OrderTypeSubscription OrderType = "SUBSCRIPTION" // 订阅订单
	OrderTypeRenewal      OrderType = "RENEWAL"      // 续费订单
)

// PaymentPlatform 支付平台枚举
type PaymentPlatform string

const (
	PaymentPlatformAlipay PaymentPlatform = "alipay" // 支付宝
	PaymentPlatformWechat PaymentPlatform = "wechat" // 微信支付
	PaymentPlatformStripe PaymentPlatform = "stripe" // Stripe
)

// PaymentType 支付类型枚举
type PaymentType string

const (
	PaymentTypeWeb         PaymentType = "WEB"          // 网页支付
	PaymentTypeQRCode      PaymentType = "QR_CODE"      // 二维码支付
	PaymentTypeMobile      PaymentType = "MOBILE"       // 移动端支付
	PaymentTypeH5          PaymentType = "H5"           // H5支付
	PaymentTypeMiniProgram PaymentType = "mini_program" // 小程序支付
)

// ---- 实体 ----

// OrderEntity 订单实体
type OrderEntity struct {
	ID              string            `gorm:"column:id;primaryKey" json:"id"`
	UserID          string            `gorm:"column:user_id" json:"userId"`
	OrderNo         string            `gorm:"column:order_no" json:"orderNo"`
	OrderType       OrderType         `gorm:"column:order_type" json:"orderType"`
	Title           string            `gorm:"column:title" json:"title"`
	Description     string            `gorm:"column:description" json:"description"`
	Amount          float64           `gorm:"column:amount" json:"amount"`
	Currency        string            `gorm:"column:currency" json:"currency"`
	Status          OrderStatus       `gorm:"column:status" json:"status"`
	ExpiredAt       *time.Time        `gorm:"column:expired_at" json:"expiredAt"`
	PaidAt          *time.Time        `gorm:"column:paid_at" json:"paidAt"`
	CancelledAt     *time.Time        `gorm:"column:cancelled_at" json:"cancelledAt"`
	RefundedAt      *time.Time        `gorm:"column:refunded_at" json:"refundedAt"`
	RefundAmount    float64           `gorm:"column:refund_amount" json:"refundAmount"`
	PaymentPlatform PaymentPlatform   `gorm:"column:payment_platform" json:"paymentPlatform"`
	PaymentType     PaymentType       `gorm:"column:payment_type" json:"paymentType"`
	ProviderOrderID string            `gorm:"column:provider_order_id" json:"providerOrderId"`
	Metadata        map[string]interface{} `gorm:"column:metadata;serializer:json" json:"metadata"`

	entity.BaseEntity
}

func (OrderEntity) TableName() string {
	return "orders"
}

// GenerateOrderNo 生成订单号
func (e *OrderEntity) GenerateOrderNo() {
	if e.OrderNo == "" {
		e.OrderNo = fmt.Sprintf("ORD%d%08X", time.Now().UnixMilli(), rand.Uint32())
	}
}

// SetDefaultExpiration 设置默认过期时间（30分钟）
func (e *OrderEntity) SetDefaultExpiration() {
	if e.ExpiredAt == nil {
		t := time.Now().Add(30 * time.Minute)
		e.ExpiredAt = &t
	}
}

// IsExpired 检查订单是否过期
func (e *OrderEntity) IsExpired() bool {
	return e.ExpiredAt != nil && time.Now().After(*e.ExpiredAt)
}

// MarkAsPaid 标记订单为已支付
func (e *OrderEntity) MarkAsPaid() error {
	if !e.Status.CanPay() {
		return exception.NewBusinessException("订单状态不允许支付操作")
	}
	e.Status = OrderStatusPaid
	now := time.Now()
	e.PaidAt = &now
	return nil
}

// Cancel 取消订单
func (e *OrderEntity) Cancel() error {
	if !e.Status.CanCancel() {
		return exception.NewBusinessException("订单状态不允许取消操作")
	}
	e.Status = OrderStatusCancelled
	now := time.Now()
	e.CancelledAt = &now
	return nil
}

// Refund 退款订单
func (e *OrderEntity) Refund(refundAmount float64) error {
	if !e.Status.CanRefund() {
		return exception.NewBusinessException("订单状态不允许退款操作")
	}
	if refundAmount <= 0 {
		return exception.NewBusinessException("退款金额必须大于0")
	}
	if refundAmount > e.Amount {
		return exception.NewBusinessException("退款金额不能超过订单金额")
	}
	e.Status = OrderStatusRefunded
	e.RefundAmount = refundAmount
	now := time.Now()
	e.RefundedAt = &now
	return nil
}

// MarkAsExpired 标记订单为已过期
func (e *OrderEntity) MarkAsExpired() {
	if e.Status.IsFinished() {
		return
	}
	e.Status = OrderStatusExpired
}

// Validate 验证订单数据
func (e *OrderEntity) Validate() error {
	if e.UserID == "" {
		return exception.NewBusinessException("用户ID不能为空")
	}
	if e.OrderType == "" {
		return exception.NewBusinessException("订单类型不能为空")
	}
	if e.Title == "" {
		return exception.NewBusinessException("订单标题不能为空")
	}
	if e.Amount <= 0 {
		return exception.NewBusinessException("订单金额必须大于0")
	}
	if e.Currency == "" {
		e.Currency = "CNY"
	}
	if e.Status == "" {
		e.Status = OrderStatusPending
	}
	return nil
}

// CreateNewOrder 创建新订单
func CreateNewOrder(userID string, orderType OrderType, title, description string, amount float64) (*OrderEntity, error) {
	order := &OrderEntity{
		ID:           uuid.New().String(),
		UserID:       userID,
		OrderType:    orderType,
		Title:        title,
		Description:  description,
		Amount:       amount,
		Currency:     "CNY",
		Status:       OrderStatusPending,
		RefundAmount: 0,
	}
	order.GenerateOrderNo()
	order.SetDefaultExpiration()
	if err := order.Validate(); err != nil {
		return nil, err
	}
	return order, nil
}

// CreateRechargeOrder 创建充值订单
func CreateRechargeOrder(userID string, amount float64, remark string) (*OrderEntity, error) {
	title := "账户充值"
	description := remark
	if description == "" {
		description = fmt.Sprintf("充值 ¥%.2f 到账户余额", amount)
	}
	return CreateNewOrder(userID, OrderTypeRecharge, title, description, amount)
}

// ---- 仓储接口 ----

type OrderRepository interface {
	FindByID(id string) (*OrderEntity, error)
	FindByOrderNo(orderNo string) (*OrderEntity, error)
	FindByUserID(userID string, orderType *OrderType, status *OrderStatus) ([]OrderEntity, error)
	FindByUserIDPaged(userID string, page, pageSize int, orderType *OrderType, status *OrderStatus) ([]OrderEntity, int64, error)
	FindAllPaged(page, pageSize int, keyword string) ([]OrderEntity, int64, error)
	FindExpiredPendingOrders(limit int) ([]OrderEntity, error)
	Create(entity *OrderEntity) error
	Update(entity *OrderEntity) error
	UpdateStatus(id string, status OrderStatus) error
}

// ---- GORM 实现 ----

type OrderRepositoryImpl struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{DB: db}
}

func (r *OrderRepositoryImpl) FindByID(id string) (*OrderEntity, error) {
	var e OrderEntity
	result := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *OrderRepositoryImpl) FindByOrderNo(orderNo string) (*OrderEntity, error) {
	var e OrderEntity
	result := r.DB.Where("order_no = ? AND deleted_at IS NULL", orderNo).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *OrderRepositoryImpl) FindByUserID(userID string, orderType *OrderType, status *OrderStatus) ([]OrderEntity, error) {
	var entities []OrderEntity
	query := r.DB.Where("user_id = ? AND deleted_at IS NULL", userID)
	if orderType != nil {
		query = query.Where("order_type = ?", *orderType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	result := query.Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *OrderRepositoryImpl) FindByUserIDPaged(userID string, page, pageSize int, orderType *OrderType, status *OrderStatus) ([]OrderEntity, int64, error) {
	var entities []OrderEntity
	var total int64
	query := r.DB.Model(&OrderEntity{}).Where("user_id = ? AND deleted_at IS NULL", userID)
	if orderType != nil {
		query = query.Where("order_type = ?", *orderType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	query.Count(&total)
	result := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities)
	return entities, total, result.Error
}

func (r *OrderRepositoryImpl) FindExpiredPendingOrders(limit int) ([]OrderEntity, error) {
	var entities []OrderEntity
	result := r.DB.Where("status = ? AND expired_at <= ? AND deleted_at IS NULL", OrderStatusPending, time.Now()).
		Order("expired_at ASC").Limit(limit).Find(&entities)
	return entities, result.Error
}

func (r *OrderRepositoryImpl) Create(e *OrderEntity) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return r.DB.Create(e).Error
}

func (r *OrderRepositoryImpl) Update(e *OrderEntity) error {
	return r.DB.Save(e).Error
}

func (r *OrderRepositoryImpl) UpdateStatus(id string, status OrderStatus) error {
	return r.DB.Model(&OrderEntity{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderRepositoryImpl) FindAllPaged(page, pageSize int, keyword string) ([]OrderEntity, int64, error) {
	var entities []OrderEntity
	var total int64
	query := r.DB.Model(&OrderEntity{}).Where("deleted_at IS NULL")
	if keyword != "" {
		query = query.Where("order_no LIKE ? OR title LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)
	result := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities)
	return entities, total, result.Error
}

// ---- 领域服务 ----

type DomainService struct {
	repo OrderRepository
}

func NewDomainService(repo OrderRepository) *DomainService {
	return &DomainService{repo: repo}
}

// OrderRepo 获取仓储（供应用层使用）
func (s *DomainService) OrderRepo() OrderRepository {
	return s.repo
}

// CreateOrder 创建订单
func (s *DomainService) CreateOrder(order *OrderEntity) (*OrderEntity, error) {
	if err := order.Validate(); err != nil {
		return nil, err
	}
	if err := s.repo.Create(order); err != nil {
		return nil, err
	}
	return order, nil
}

// GetOrderByID 根据ID获取订单
func (s *DomainService) GetOrderByID(orderID string) (*OrderEntity, error) {
	if orderID == "" {
		return nil, exception.NewBusinessException("订单ID不能为空")
	}
	order, err := s.repo.FindByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, exception.NewBusinessException("订单不存在: " + orderID)
	}
	return order, nil
}

// FindOrderByOrderNo 根据订单号查找订单
func (s *DomainService) FindOrderByOrderNo(orderNo string) (*OrderEntity, error) {
	return s.repo.FindByOrderNo(orderNo)
}

// GetOrdersByUserID 获取用户订单列表
func (s *DomainService) GetOrdersByUserID(userID string, orderType *OrderType, status *OrderStatus) ([]OrderEntity, error) {
	if userID == "" {
		return nil, exception.NewBusinessException("用户ID不能为空")
	}
	return s.repo.FindByUserID(userID, orderType, status)
}

// GetOrdersByUserIDPaged 分页查询用户订单
func (s *DomainService) GetOrdersByUserIDPaged(userID string, page, pageSize int, orderType *OrderType, status *OrderStatus) ([]OrderEntity, int64, error) {
	if userID == "" {
		return nil, 0, exception.NewBusinessException("用户ID不能为空")
	}
	return s.repo.FindByUserIDPaged(userID, page, pageSize, orderType, status)
}

// MarkOrderAsPaid 标记订单为已支付
func (s *DomainService) MarkOrderAsPaid(orderID string) error {
	order, err := s.GetOrderByID(orderID)
	if err != nil {
		return err
	}
	if err := order.MarkAsPaid(); err != nil {
		return err
	}
	return s.repo.Update(order)
}

// CancelOrder 取消订单
func (s *DomainService) CancelOrder(orderID string) error {
	order, err := s.GetOrderByID(orderID)
	if err != nil {
		return err
	}
	if err := order.Cancel(); err != nil {
		return err
	}
	return s.repo.Update(order)
}

// RefundOrder 退款订单
func (s *DomainService) RefundOrder(orderID string, refundAmount float64) error {
	order, err := s.GetOrderByID(orderID)
	if err != nil {
		return err
	}
	if err := order.Refund(refundAmount); err != nil {
		return err
	}
	return s.repo.Update(order)
}

// FindExpiredPendingOrders 查找过期的未支付订单
func (s *DomainService) FindExpiredPendingOrders(limit int) ([]OrderEntity, error) {
	return s.repo.FindExpiredPendingOrders(limit)
}

// UpdateOrderStatus 更新订单状态
func (s *DomainService) UpdateOrderStatus(orderID string, newStatus OrderStatus) error {
	if orderID == "" {
		return exception.NewBusinessException("订单ID不能为空")
	}
	return s.repo.UpdateStatus(orderID, newStatus)
}
