package product

import (
	"github.com/google/uuid"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	"gorm.io/gorm"
)

// ---- 枚举 ----

// BillingType 计费类型枚举
type BillingType string

const (
	BillingTypeModelUsage   BillingType = "MODEL_USAGE"   // 模型调用计费
	BillingTypeAgentCreation BillingType = "AGENT_CREATION" // Agent创建计费
	BillingTypeAgentUsage   BillingType = "AGENT_USAGE"   // Agent使用计费
	BillingTypeAPICall      BillingType = "API_CALL"      // API调用计费
	BillingTypeStorageUsage BillingType = "STORAGE_USAGE" // 存储使用计费
)

// ---- 实体 ----

// ProductEntity 商品实体
type ProductEntity struct {
	ID            string                 `gorm:"column:id;primaryKey" json:"id"`
	Name          string                 `gorm:"column:name" json:"name"`
	Type          BillingType            `gorm:"column:type" json:"type"`
	ServiceID     string                 `gorm:"column:service_id" json:"serviceId"`
	RuleID        string                 `gorm:"column:rule_id" json:"ruleId"`
	PricingConfig map[string]interface{} `gorm:"column:pricing_config;serializer:json" json:"pricingConfig"`
	Status        int                    `gorm:"column:status" json:"status"`

	entity.BaseEntity
}

func (ProductEntity) TableName() string {
	return "products"
}

// IsActive 检查商品是否激活
func (e *ProductEntity) IsActive() bool {
	return e.Status == 1
}

// Activate 激活商品
func (e *ProductEntity) Activate() {
	e.Status = 1
}

// Deactivate 禁用商品
func (e *ProductEntity) Deactivate() {
	e.Status = 0
}

// Validate 验证商品基本信息
func (e *ProductEntity) Validate() error {
	if e.Name == "" {
		return exception.NewBusinessException("商品名称不能为空")
	}
	if e.Type == "" {
		return exception.NewBusinessException("商品类型不能为空")
	}
	if e.ServiceID == "" {
		return exception.NewBusinessException("业务ID不能为空")
	}
	if e.RuleID == "" {
		return exception.NewBusinessException("规则ID不能为空")
	}
	if e.PricingConfig == nil || len(e.PricingConfig) == 0 {
		return exception.NewBusinessException("商品价格配置不能为空")
	}
	return nil
}

// ---- 仓储接口 ----

type ProductRepository interface {
	FindByID(id string) (*ProductEntity, error)
	FindByBusinessKey(billingType BillingType, serviceID string) (*ProductEntity, error)
	FindActiveProducts(billingType *BillingType) ([]ProductEntity, error)
	FindAll() ([]ProductEntity, error)
	FindPaged(page, pageSize int, keyword string) ([]ProductEntity, int64, error)
	FindByIDs(ids []string) ([]ProductEntity, error)
	Create(entity *ProductEntity) error
	Update(entity *ProductEntity) error
	Delete(id string) error
}

// ---- GORM 实现 ----

type ProductRepositoryImpl struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{DB: db}
}

func (r *ProductRepositoryImpl) FindByID(id string) (*ProductEntity, error) {
	var e ProductEntity
	result := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *ProductRepositoryImpl) FindByBusinessKey(billingType BillingType, serviceID string) (*ProductEntity, error) {
	var e ProductEntity
	result := r.DB.Where("type = ? AND service_id = ? AND deleted_at IS NULL", billingType, serviceID).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *ProductRepositoryImpl) FindActiveProducts(billingType *BillingType) ([]ProductEntity, error) {
	var entities []ProductEntity
	query := r.DB.Where("status = 1 AND deleted_at IS NULL")
	if billingType != nil {
		query = query.Where("type = ?", *billingType)
	}
	result := query.Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *ProductRepositoryImpl) FindAll() ([]ProductEntity, error) {
	var entities []ProductEntity
	result := r.DB.Where("deleted_at IS NULL").Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *ProductRepositoryImpl) FindByIDs(ids []string) ([]ProductEntity, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var entities []ProductEntity
	result := r.DB.Where("id IN ? AND deleted_at IS NULL", ids).Find(&entities)
	return entities, result.Error
}

func (r *ProductRepositoryImpl) Create(e *ProductEntity) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return r.DB.Create(e).Error
}

func (r *ProductRepositoryImpl) Update(e *ProductEntity) error {
	return r.DB.Save(e).Error
}

func (r *ProductRepositoryImpl) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&ProductEntity{}).Error
}

func (r *ProductRepositoryImpl) FindPaged(page, pageSize int, keyword string) ([]ProductEntity, int64, error) {
	var entities []ProductEntity
	var total int64
	query := r.DB.Model(&ProductEntity{}).Where("deleted_at IS NULL")
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	query.Count(&total)
	result := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities)
	return entities, total, result.Error
}

// ---- 领域服务 ----

type DomainService struct {
	repo ProductRepository
}

func NewDomainService(repo ProductRepository) *DomainService {
	return &DomainService{repo: repo}
}

// ProductRepo 获取仓储（供应用层使用）
func (s *DomainService) ProductRepo() ProductRepository {
	return s.repo
}

// FindProductByBusinessKey 根据业务主键查找商品
func (s *DomainService) FindProductByBusinessKey(billingType BillingType, serviceID string) (*ProductEntity, error) {
	return s.repo.FindByBusinessKey(billingType, serviceID)
}

// CreateProduct 创建商品
func (s *DomainService) CreateProduct(product *ProductEntity) (*ProductEntity, error) {
	if err := product.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.repo.FindByBusinessKey(product.Type, product.ServiceID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, exception.NewBusinessException("该业务类型和业务ID的商品已存在")
	}
	if product.Status == 0 {
		product.Activate()
	}
	if err := s.repo.Create(product); err != nil {
		return nil, err
	}
	return product, nil
}

// UpdateProduct 更新商品
func (s *DomainService) UpdateProduct(product *ProductEntity) (*ProductEntity, error) {
	if product.ID == "" {
		return nil, exception.NewBusinessException("商品ID不能为空")
	}
	existing, err := s.repo.FindByID(product.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, exception.NewBusinessException("商品不存在")
	}
	if err := s.repo.Update(product); err != nil {
		return nil, err
	}
	return product, nil
}

// GetProductByID 根据ID获取商品
func (s *DomainService) GetProductByID(productID string) (*ProductEntity, error) {
	return s.repo.FindByID(productID)
}

// GetActiveProducts 获取激活的商品列表
func (s *DomainService) GetActiveProducts(billingType *BillingType) ([]ProductEntity, error) {
	return s.repo.FindActiveProducts(billingType)
}

// GetAllProducts 获取所有商品
func (s *DomainService) GetAllProducts() ([]ProductEntity, error) {
	return s.repo.FindAll()
}

// DeleteProduct 删除商品
func (s *DomainService) DeleteProduct(productID string) error {
	existing, err := s.repo.FindByID(productID)
	if err != nil {
		return err
	}
	if existing == nil {
		return exception.NewBusinessException("商品不存在")
	}
	return s.repo.Delete(productID)
}

// IsProductActive 检查商品是否存在且激活
func (s *DomainService) IsProductActive(billingType BillingType, serviceID string) (bool, error) {
	product, err := s.repo.FindByBusinessKey(billingType, serviceID)
	if err != nil {
		return false, err
	}
	return product != nil && product.IsActive(), nil
}
