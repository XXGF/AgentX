package rule

import (
	"github.com/google/uuid"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	"gorm.io/gorm"
)

// ---- 枚举 ----

// RuleHandlerKey 规则处理器标识枚举
type RuleHandlerKey string

const (
	RuleHandlerKeyModelTokenStrategy  RuleHandlerKey = "MODEL_TOKEN_STRATEGY"  // 模型Token计费策略
	RuleHandlerKeyPerUnitStrategy     RuleHandlerKey = "PER_UNIT_STRATEGY"     // 按次计费策略
	RuleHandlerKeyTieredStrategy      RuleHandlerKey = "TIERED_STRATEGY"       // 分层计费策略
	RuleHandlerKeyVolumeTieredStrategy RuleHandlerKey = "VOLUME_TIERED_STRATEGY" // 按量阶梯计费策略
)

// ---- 实体 ----

// RuleEntity 规则实体
type RuleEntity struct {
	ID          string         `gorm:"column:id;primaryKey" json:"id"`
	Name        string         `gorm:"column:name" json:"name"`
	HandlerKey  RuleHandlerKey `gorm:"column:handler_key" json:"handlerKey"`
	Description string         `gorm:"column:description" json:"description"`

	entity.BaseEntity
}

func (RuleEntity) TableName() string {
	return "rules"
}

// Validate 验证规则信息
func (e *RuleEntity) Validate() error {
	if e.Name == "" {
		return exception.NewBusinessException("规则名称不能为空")
	}
	if e.HandlerKey == "" {
		return exception.NewBusinessException("规则处理器标识不能为空")
	}
	return nil
}

// ---- 仓储接口 ----

type RuleRepository interface {
	FindByID(id string) (*RuleEntity, error)
	FindByHandlerKey(handlerKey RuleHandlerKey) (*RuleEntity, error)
	FindAll() ([]RuleEntity, error)
	FindPaged(page, pageSize int) ([]RuleEntity, int64, error)
	Create(entity *RuleEntity) error
	Update(entity *RuleEntity) error
	Delete(id string) error
}

// ---- GORM 实现 ----

type RuleRepositoryImpl struct {
	DB *gorm.DB
}

func NewRuleRepository(db *gorm.DB) RuleRepository {
	return &RuleRepositoryImpl{DB: db}
}

func (r *RuleRepositoryImpl) FindByID(id string) (*RuleEntity, error) {
	var e RuleEntity
	result := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *RuleRepositoryImpl) FindByHandlerKey(handlerKey RuleHandlerKey) (*RuleEntity, error) {
	var e RuleEntity
	result := r.DB.Where("handler_key = ? AND deleted_at IS NULL", handlerKey).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *RuleRepositoryImpl) FindAll() ([]RuleEntity, error) {
	var entities []RuleEntity
	result := r.DB.Where("deleted_at IS NULL").Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *RuleRepositoryImpl) FindPaged(page, pageSize int) ([]RuleEntity, int64, error) {
	var entities []RuleEntity
	var total int64
	query := r.DB.Model(&RuleEntity{}).Where("deleted_at IS NULL")
	query.Count(&total)
	result := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities)
	return entities, total, result.Error
}

func (r *RuleRepositoryImpl) Create(e *RuleEntity) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return r.DB.Create(e).Error
}

func (r *RuleRepositoryImpl) Update(e *RuleEntity) error {
	return r.DB.Save(e).Error
}

func (r *RuleRepositoryImpl) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&RuleEntity{}).Error
}

// ---- 领域服务 ----

type DomainService struct {
	repo RuleRepository
}

func NewDomainService(repo RuleRepository) *DomainService {
	return &DomainService{repo: repo}
}

// Repo 获取仓储（供应用层使用）
func (s *DomainService) Repo() RuleRepository {
	return s.repo
}

// GetRuleByID 根据ID获取规则
func (s *DomainService) GetRuleByID(ruleID string) (*RuleEntity, error) {
	return s.repo.FindByID(ruleID)
}

// GetRuleByHandlerKey 根据处理器标识获取规则
func (s *DomainService) GetRuleByHandlerKey(handlerKey RuleHandlerKey) (*RuleEntity, error) {
	return s.repo.FindByHandlerKey(handlerKey)
}

// CreateRule 创建规则
func (s *DomainService) CreateRule(rule *RuleEntity) (*RuleEntity, error) {
	if err := rule.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.repo.FindByHandlerKey(rule.HandlerKey)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, exception.NewBusinessException("该处理器标识的规则已存在")
	}
	if err := s.repo.Create(rule); err != nil {
		return nil, err
	}
	return rule, nil
}

// UpdateRule 更新规则
func (s *DomainService) UpdateRule(rule *RuleEntity) (*RuleEntity, error) {
	if rule.ID == "" {
		return nil, exception.NewBusinessException("规则ID不能为空")
	}
	if err := rule.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.repo.FindByID(rule.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, exception.NewBusinessException("规则不存在")
	}
	if err := s.repo.Update(rule); err != nil {
		return nil, err
	}
	return rule, nil
}

// DeleteRule 删除规则
func (s *DomainService) DeleteRule(ruleID string) error {
	existing, err := s.repo.FindByID(ruleID)
	if err != nil {
		return err
	}
	if existing == nil {
		return exception.NewBusinessException("规则不存在")
	}
	return s.repo.Delete(ruleID)
}

// GetAllRules 获取所有规则
func (s *DomainService) GetAllRules() ([]RuleEntity, error) {
	return s.repo.FindAll()
}

// ExistsRule 检查规则是否存在
func (s *DomainService) ExistsRule(ruleID string) (bool, error) {
	rule, err := s.repo.FindByID(ruleID)
	if err != nil {
		return false, err
	}
	return rule != nil, nil
}
