package user

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	"gorm.io/gorm"
)

// ---- AccountEntity 完整版（补充 credit, totalConsumed, lastTransactionAt 字段）----

// AccountEntity 账户实体（对应 Java 的 AccountEntity）
type AccountEntity struct {
	ID                string     `gorm:"column:id;primaryKey" json:"id"`
	UserID            string     `gorm:"column:user_id" json:"userId"`
	Balance           float64    `gorm:"column:balance" json:"balance"`
	Credit            float64    `gorm:"column:credit" json:"credit"`
	TotalConsumed     float64    `gorm:"column:total_consumed" json:"totalConsumed"`
	LastTransactionAt *time.Time `gorm:"column:last_transaction_at" json:"lastTransactionAt"`

	entity.BaseEntity
}

func (AccountEntity) TableName() string {
	return "accounts"
}

// GetAvailableBalance 获取可用余额（余额+信用额度）
func (e *AccountEntity) GetAvailableBalance() float64 {
	return e.Balance + e.Credit
}

// CheckSufficientBalance 检查余额是否充足
func (e *AccountEntity) CheckSufficientBalance(amount float64) bool {
	return e.GetAvailableBalance() >= amount
}

// Deduct 扣除余额
func (e *AccountEntity) Deduct(amount float64) error {
	if amount <= 0 {
		return exception.NewBusinessException("扣费金额必须大于0")
	}
	if !e.CheckSufficientBalance(amount) {
		return exception.NewBusinessException("账户余额不足")
	}
	if e.Balance >= amount {
		e.Balance -= amount
	} else {
		remaining := amount - e.Balance
		e.Balance = 0
		e.Credit -= remaining
	}
	e.TotalConsumed += amount
	now := time.Now()
	e.LastTransactionAt = &now
	return nil
}

// Recharge 充值
func (e *AccountEntity) Recharge(amount float64) error {
	if amount <= 0 {
		return exception.NewBusinessException("充值金额必须大于0")
	}
	e.Balance += amount
	now := time.Now()
	e.LastTransactionAt = &now
	return nil
}

// AddCredit 增加信用额度
func (e *AccountEntity) AddCredit(amount float64) error {
	if amount <= 0 {
		return exception.NewBusinessException("信用额度必须大于0")
	}
	e.Credit += amount
	now := time.Now()
	e.LastTransactionAt = &now
	return nil
}

// CreateNewAccount 创建新账户
func CreateNewAccount(userID string) *AccountEntity {
	return &AccountEntity{
		ID:     uuid.New().String(),
		UserID: userID,
	}
}

// ---- UsageRecordEntity 用量记录实体 ----

type UsageRecordEntity struct {
	ID                 string                 `gorm:"column:id;primaryKey" json:"id"`
	UserID             string                 `gorm:"column:user_id" json:"userId"`
	ProductID          string                 `gorm:"column:product_id" json:"productId"`
	QuantityData       map[string]interface{} `gorm:"column:quantity_data;serializer:json" json:"quantityData"`
	Cost               float64                `gorm:"column:cost" json:"cost"`
	RequestID          string                 `gorm:"column:request_id" json:"requestId"`
	BilledAt           *time.Time             `gorm:"column:billed_at" json:"billedAt"`
	ServiceName        string                 `gorm:"column:service_name" json:"serviceName"`
	ServiceType        string                 `gorm:"column:service_type" json:"serviceType"`
	ServiceDescription string                 `gorm:"column:service_description" json:"serviceDescription"`
	PricingRule        string                 `gorm:"column:pricing_rule" json:"pricingRule"`
	RelatedEntityName  string                 `gorm:"column:related_entity_name" json:"relatedEntityName"`

	entity.BaseEntity
}

func (UsageRecordEntity) TableName() string {
	return "usage_records"
}

// ---- 仓储接口 ----

type AccountRepository interface {
	FindByID(id string) (*AccountEntity, error)
	FindByUserID(userID string) (*AccountEntity, error)
	Create(account *AccountEntity) error
	Update(account *AccountEntity) error
}

type UsageRecordRepository interface {
	FindByID(id string) (*UsageRecordEntity, error)
	FindByUserIDPaged(userID string, page, pageSize int) ([]UsageRecordEntity, int64, error)
	FindByUserIDAndProductIDPaged(userID, productID string, page, pageSize int) ([]UsageRecordEntity, int64, error)
	FindByUserIDAndTimeRange(userID string, startTime, endTime time.Time) ([]UsageRecordEntity, error)
	FindByProductIDPaged(productID string, page, pageSize int) ([]UsageRecordEntity, int64, error)
	ExistsByRequestID(requestID string) (bool, error)
	QueryPaged(userID, productID, requestID string, startTime, endTime *time.Time, page, pageSize int) ([]UsageRecordEntity, int64, error)
	Create(record *UsageRecordEntity) error
}

// ---- GORM 实现 ----

type AccountRepositoryImpl struct {
	DB *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &AccountRepositoryImpl{DB: db}
}

func (r *AccountRepositoryImpl) FindByID(id string) (*AccountEntity, error) {
	var e AccountEntity
	result := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *AccountRepositoryImpl) FindByUserID(userID string) (*AccountEntity, error) {
	var e AccountEntity
	result := r.DB.Where("user_id = ? AND deleted_at IS NULL", userID).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *AccountRepositoryImpl) Create(e *AccountEntity) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return r.DB.Create(e).Error
}

func (r *AccountRepositoryImpl) Update(e *AccountEntity) error {
	return r.DB.Save(e).Error
}

// ---- UsageRecord GORM 实现 ----

type UsageRecordRepositoryImpl struct {
	DB *gorm.DB
}

func NewUsageRecordRepository(db *gorm.DB) UsageRecordRepository {
	return &UsageRecordRepositoryImpl{DB: db}
}

func (r *UsageRecordRepositoryImpl) FindByID(id string) (*UsageRecordEntity, error) {
	var e UsageRecordEntity
	result := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *UsageRecordRepositoryImpl) FindByUserIDPaged(userID string, page, pageSize int) ([]UsageRecordEntity, int64, error) {
	var entities []UsageRecordEntity
	var total int64
	query := r.DB.Model(&UsageRecordEntity{}).Where("user_id = ? AND deleted_at IS NULL", userID)
	query.Count(&total)
	result := query.Order("billed_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities)
	return entities, total, result.Error
}

func (r *UsageRecordRepositoryImpl) FindByUserIDAndProductIDPaged(userID, productID string, page, pageSize int) ([]UsageRecordEntity, int64, error) {
	var entities []UsageRecordEntity
	var total int64
	query := r.DB.Model(&UsageRecordEntity{}).Where("user_id = ? AND product_id = ? AND deleted_at IS NULL", userID, productID)
	query.Count(&total)
	result := query.Order("billed_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities)
	return entities, total, result.Error
}

func (r *UsageRecordRepositoryImpl) FindByUserIDAndTimeRange(userID string, startTime, endTime time.Time) ([]UsageRecordEntity, error) {
	var entities []UsageRecordEntity
	result := r.DB.Where("user_id = ? AND billed_at >= ? AND billed_at <= ? AND deleted_at IS NULL", userID, startTime, endTime).
		Order("billed_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *UsageRecordRepositoryImpl) FindByProductIDPaged(productID string, page, pageSize int) ([]UsageRecordEntity, int64, error) {
	var entities []UsageRecordEntity
	var total int64
	query := r.DB.Model(&UsageRecordEntity{}).Where("product_id = ? AND deleted_at IS NULL", productID)
	query.Count(&total)
	result := query.Order("billed_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities)
	return entities, total, result.Error
}

func (r *UsageRecordRepositoryImpl) ExistsByRequestID(requestID string) (bool, error) {
	var count int64
	result := r.DB.Model(&UsageRecordEntity{}).Where("request_id = ? AND deleted_at IS NULL", requestID).Count(&count)
	return count > 0, result.Error
}

func (r *UsageRecordRepositoryImpl) QueryPaged(userID, productID, requestID string, startTime, endTime *time.Time, page, pageSize int) ([]UsageRecordEntity, int64, error) {
	var entities []UsageRecordEntity
	var total int64
	query := r.DB.Model(&UsageRecordEntity{}).Where("deleted_at IS NULL")
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if productID != "" {
		query = query.Where("product_id = ?", productID)
	}
	if requestID != "" {
		query = query.Where("request_id = ?", requestID)
	}
	if startTime != nil {
		query = query.Where("billed_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("billed_at <= ?", *endTime)
	}
	query.Count(&total)
	result := query.Order("billed_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities)
	return entities, total, result.Error
}

func (r *UsageRecordRepositoryImpl) Create(e *UsageRecordEntity) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return r.DB.Create(e).Error
}

// ---- Account 领域服务 ----

type AccountDomainService struct {
	repo  AccountRepository
	locks sync.Map // 用户级别的锁
}

func NewAccountDomainService(repo AccountRepository) *AccountDomainService {
	return &AccountDomainService{repo: repo}
}

func (s *AccountDomainService) getUserLock(userID string) *sync.Mutex {
	lock, _ := s.locks.LoadOrStore(userID, &sync.Mutex{})
	return lock.(*sync.Mutex)
}

// GetOrCreateAccount 获取或创建用户账户
func (s *AccountDomainService) GetOrCreateAccount(userID string) (*AccountEntity, error) {
	account, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if account != nil {
		return account, nil
	}
	lock := s.getUserLock(userID)
	lock.Lock()
	defer lock.Unlock()
	account, err = s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if account != nil {
		return account, nil
	}
	account = CreateNewAccount(userID)
	if err := s.repo.Create(account); err != nil {
		return nil, err
	}
	return account, nil
}

// GetAccountByUserID 根据用户ID获取账户
func (s *AccountDomainService) GetAccountByUserID(userID string) (*AccountEntity, error) {
	return s.repo.FindByUserID(userID)
}

// GetAccountByID 根据ID获取账户
func (s *AccountDomainService) GetAccountByID(accountID string) (*AccountEntity, error) {
	return s.repo.FindByID(accountID)
}

// DeductBalance 扣除余额
func (s *AccountDomainService) DeductBalance(userID string, amount float64) error {
	lock := s.getUserLock(userID)
	lock.Lock()
	defer lock.Unlock()
	account, err := s.GetOrCreateAccount(userID)
	if err != nil {
		return err
	}
	if err := account.Deduct(amount); err != nil {
		return err
	}
	return s.repo.Update(account)
}

// RechargeBalance 充值
func (s *AccountDomainService) RechargeBalance(userID string, amount float64) error {
	lock := s.getUserLock(userID)
	lock.Lock()
	defer lock.Unlock()
	account, err := s.GetOrCreateAccount(userID)
	if err != nil {
		return err
	}
	if err := account.Recharge(amount); err != nil {
		return err
	}
	return s.repo.Update(account)
}

// AddCredit 增加信用额度
func (s *AccountDomainService) AddCredit(userID string, amount float64) error {
	lock := s.getUserLock(userID)
	lock.Lock()
	defer lock.Unlock()
	account, err := s.GetOrCreateAccount(userID)
	if err != nil {
		return err
	}
	if err := account.AddCredit(amount); err != nil {
		return err
	}
	return s.repo.Update(account)
}

// CheckSufficientBalance 检查余额是否充足
func (s *AccountDomainService) CheckSufficientBalance(userID string, amount float64) (bool, error) {
	account, err := s.repo.FindByUserID(userID)
	if err != nil {
		return false, err
	}
	if account == nil {
		return false, nil
	}
	return account.CheckSufficientBalance(amount), nil
}

// GetAvailableBalance 获取可用余额
func (s *AccountDomainService) GetAvailableBalance(userID string) (float64, error) {
	account, err := s.repo.FindByUserID(userID)
	if err != nil {
		return 0, err
	}
	if account == nil {
		return 0, nil
	}
	return account.GetAvailableBalance(), nil
}

// ExistsAccount 检查账户是否存在
func (s *AccountDomainService) ExistsAccount(userID string) (bool, error) {
	account, err := s.repo.FindByUserID(userID)
	if err != nil {
		return false, err
	}
	return account != nil, nil
}

// ---- UsageRecord 领域服务 ----

type UsageRecordDomainService struct {
	repo UsageRecordRepository
}

func NewUsageRecordDomainService(repo UsageRecordRepository) *UsageRecordDomainService {
	return &UsageRecordDomainService{repo: repo}
}

// GetUsageRecordByID 根据ID获取用量记录
func (s *UsageRecordDomainService) GetUsageRecordByID(recordID string) (*UsageRecordEntity, error) {
	return s.repo.FindByID(recordID)
}

// GetUserUsageHistory 获取用户用量历史（分页）
func (s *UsageRecordDomainService) GetUserUsageHistory(userID string, page, pageSize int) ([]UsageRecordEntity, int64, error) {
	return s.repo.FindByUserIDPaged(userID, page, pageSize)
}

// ExistsByRequestID 检查请求ID是否存在
func (s *UsageRecordDomainService) ExistsByRequestID(requestID string) (bool, error) {
	return s.repo.ExistsByRequestID(requestID)
}

// CreateUsageRecord 创建用量记录
func (s *UsageRecordDomainService) CreateUsageRecord(record *UsageRecordEntity) error {
	return s.repo.Create(record)
}

// QueryUsageRecords 按条件查询用量记录
func (s *UsageRecordDomainService) QueryUsageRecords(userID, productID, requestID string, startTime, endTime *time.Time, page, pageSize int) ([]UsageRecordEntity, int64, error) {
	return s.repo.QueryPaged(userID, productID, requestID, startTime, endTime, page, pageSize)
}

// GetUserTotalCost 获取用户总消费
func (s *UsageRecordDomainService) GetUserTotalCost(userID string) (float64, error) {
	records, _, err := s.repo.FindByUserIDPaged(userID, 1, 10000)
	if err != nil {
		return 0, err
	}
	var total float64
	for _, r := range records {
		total += r.Cost
	}
	return total, nil
}
