package product

import (
	"time"

	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/product"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// ProductDTO 商品DTO
type ProductDTO struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Type          string                 `json:"type"`
	ServiceID     string                 `json:"serviceId"`
	RuleID        string                 `json:"ruleId"`
	PricingConfig map[string]interface{} `json:"pricingConfig"`
	Status        int                    `json:"status"`
	ModelName     string                 `json:"modelName,omitempty"`
	ModelID       string                 `json:"modelId,omitempty"`
	ProviderName  string                 `json:"providerName,omitempty"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
}

func entityToDTO(e *domain.ProductEntity) *ProductDTO {
	if e == nil {
		return nil
	}
	return &ProductDTO{
		ID:            e.ID,
		Name:          e.Name,
		Type:          string(e.Type),
		ServiceID:     e.ServiceID,
		RuleID:        e.RuleID,
		PricingConfig: e.PricingConfig,
		Status:        e.Status,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

func entitiesToDTOs(entities []domain.ProductEntity) []*ProductDTO {
	dtos := make([]*ProductDTO, len(entities))
	for i, e := range entities {
		dtos[i] = entityToDTO(&e)
	}
	return dtos
}

// AppService 商品应用服务
type AppService struct {
	productDomainService *domain.DomainService
}

func NewAppService(productDomainService *domain.DomainService) *AppService {
	return &AppService{productDomainService: productDomainService}
}

// GetProductByID 根据ID获取商品
func (s *AppService) GetProductByID(productID string) (*ProductDTO, error) {
	entity, err := s.productDomainService.GetProductByID(productID)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, exception.NewBusinessException("商品不存在")
	}
	return entityToDTO(entity), nil
}

// GetProductByBusinessKey 根据业务主键获取商品
func (s *AppService) GetProductByBusinessKey(billingType, serviceID string) (*ProductDTO, error) {
	entity, err := s.productDomainService.FindProductByBusinessKey(domain.BillingType(billingType), serviceID)
	if err != nil {
		return nil, err
	}
	return entityToDTO(entity), nil
}

// GetActiveProducts 获取激活的商品列表
func (s *AppService) GetActiveProducts(billingType string) ([]*ProductDTO, error) {
	var bt *domain.BillingType
	if billingType != "" {
		t := domain.BillingType(billingType)
		bt = &t
	}
	entities, err := s.productDomainService.GetActiveProducts(bt)
	if err != nil {
		return nil, err
	}
	return entitiesToDTOs(entities), nil
}

// IsProductActive 检查商品是否存在且激活
func (s *AppService) IsProductActive(billingType, serviceID string) (bool, error) {
	return s.productDomainService.IsProductActive(domain.BillingType(billingType), serviceID)
}

// GetProductsPaged 分页获取商品列表
func (s *AppService) GetProductsPaged(page, pageSize int, keyword string) ([]*ProductDTO, int64, error) {
	entities, total, err := s.productDomainService.ProductRepo().FindPaged(page, pageSize, keyword)
	if err != nil {
		return nil, 0, err
	}
	return entitiesToDTOs(entities), total, nil
}

// GetAllProducts 获取所有商品列表
func (s *AppService) GetAllProducts() ([]*ProductDTO, error) {
	entities, err := s.productDomainService.ProductRepo().FindAll()
	if err != nil {
		return nil, err
	}
	return entitiesToDTOs(entities), nil
}

// DeleteProduct 删除商品
func (s *AppService) DeleteProduct(productID string) error {
	return s.productDomainService.DeleteProduct(productID)
}

// ExistsProduct 检查商品是否存在
func (s *AppService) ExistsProduct(productID string) (bool, error) {
	entity, err := s.productDomainService.GetProductByID(productID)
	if err != nil {
		return false, err
	}
	return entity != nil, nil
}

// ExistsByBusinessKey 检查业务标识是否存在
func (s *AppService) ExistsByBusinessKey(billingType, serviceID string) (bool, error) {
	entity, err := s.productDomainService.FindProductByBusinessKey(domain.BillingType(billingType), serviceID)
	if err != nil {
		return false, err
	}
	return entity != nil, nil
}

// EnableProduct 启用商品
func (s *AppService) EnableProduct(productID string) (*ProductDTO, error) {
	entity, err := s.productDomainService.GetProductByID(productID)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, exception.NewBusinessException("商品不存在")
	}
	entity.Activate()
	updated, err := s.productDomainService.UpdateProduct(entity)
	if err != nil {
		return nil, err
	}
	return entityToDTO(updated), nil
}

// DisableProduct 禁用商品
func (s *AppService) DisableProduct(productID string) (*ProductDTO, error) {
	entity, err := s.productDomainService.GetProductByID(productID)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, exception.NewBusinessException("商品不存在")
	}
	entity.Deactivate()
	updated, err := s.productDomainService.UpdateProduct(entity)
	if err != nil {
		return nil, err
	}
	return entityToDTO(updated), nil
}
