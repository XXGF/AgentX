package rule

import (
	domainRule "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/rule"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// RuleDTO 规则DTO
type RuleDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	HandlerKey  string `json:"handlerKey"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// CreateRuleRequest 创建规则请求
type CreateRuleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	HandlerKey  string `json:"handlerKey" binding:"required"`
}

// UpdateRuleRequest 更新规则请求
type UpdateRuleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HandlerKey  string `json:"handlerKey"`
}

// QueryRuleRequest 查询规则请求
type QueryRuleRequest struct {
	Keyword  string `form:"keyword" json:"keyword"`
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"pageSize" json:"pageSize"`
}

func entityToDTO(e *domainRule.RuleEntity) *RuleDTO {
	if e == nil {
		return nil
	}
	return &RuleDTO{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		HandlerKey:  string(e.HandlerKey),
		CreatedAt:   e.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   e.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// AppService 规则应用服务
type AppService struct {
	domainService *domainRule.DomainService
}

func NewAppService(domainService *domainRule.DomainService) *AppService {
	return &AppService{domainService: domainService}
}

// CreateRule 创建规则
func (s *AppService) CreateRule(req *CreateRuleRequest) (*RuleDTO, error) {
	entity := &domainRule.RuleEntity{
		Name:        req.Name,
		Description: req.Description,
		HandlerKey:  domainRule.RuleHandlerKey(req.HandlerKey),
	}
	created, err := s.domainService.CreateRule(entity)
	if err != nil {
		return nil, err
	}
	return entityToDTO(created), nil
}

// UpdateRule 更新规则
func (s *AppService) UpdateRule(req *UpdateRuleRequest, ruleID string) (*RuleDTO, error) {
	existing, err := s.domainService.GetRuleByID(ruleID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, exception.NewBusinessException("规则不存在")
	}
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.HandlerKey != "" {
		existing.HandlerKey = domainRule.RuleHandlerKey(req.HandlerKey)
	}
	updated, err := s.domainService.UpdateRule(existing)
	if err != nil {
		return nil, err
	}
	return entityToDTO(updated), nil
}

// GetRuleByID 根据ID获取规则
func (s *AppService) GetRuleByID(ruleID string) (*RuleDTO, error) {
	entity, err := s.domainService.GetRuleByID(ruleID)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, exception.NewBusinessException("规则不存在")
	}
	return entityToDTO(entity), nil
}

// GetRuleByHandlerKey 根据处理器标识获取规则
func (s *AppService) GetRuleByHandlerKey(handlerKey string) (*RuleDTO, error) {
	entity, err := s.domainService.GetRuleByHandlerKey(domainRule.RuleHandlerKey(handlerKey))
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, exception.NewBusinessException("规则不存在")
	}
	return entityToDTO(entity), nil
}

// GetRules 分页查询规则
func (s *AppService) GetRules(req *QueryRuleRequest) ([]*RuleDTO, int64, error) {
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 15
	}
	entities, total, err := s.domainService.Repo().FindPaged(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	dtos := make([]*RuleDTO, len(entities))
	for i, e := range entities {
		dtos[i] = entityToDTO(&e)
	}
	return dtos, total, nil
}

// GetAllRules 获取所有规则
func (s *AppService) GetAllRules() ([]*RuleDTO, error) {
	entities, err := s.domainService.GetAllRules()
	if err != nil {
		return nil, err
	}
	dtos := make([]*RuleDTO, len(entities))
	for i, e := range entities {
		dtos[i] = entityToDTO(&e)
	}
	return dtos, nil
}

// DeleteRule 删除规则
func (s *AppService) DeleteRule(ruleID string) error {
	return s.domainService.DeleteRule(ruleID)
}

// ExistsRule 检查规则是否存在
func (s *AppService) ExistsRule(ruleID string) (bool, error) {
	entity, err := s.domainService.GetRuleByID(ruleID)
	if err != nil {
		return false, err
	}
	return entity != nil, nil
}

// ExistsByHandlerKey 检查处理器标识是否存在
func (s *AppService) ExistsByHandlerKey(handlerKey string) (bool, error) {
	entity, err := s.domainService.GetRuleByHandlerKey(domainRule.RuleHandlerKey(handlerKey))
	if err != nil {
		return false, err
	}
	return entity != nil, nil
}