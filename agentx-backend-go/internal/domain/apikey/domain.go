package apikey

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	"gorm.io/gorm"
)

// ApiKeyEntity API密钥实体
type ApiKeyEntity struct {
	ID         string     `gorm:"column:id;primaryKey" json:"id"`
	ApiKey     string     `gorm:"column:api_key" json:"apiKey"`
	AgentID    string     `gorm:"column:agent_id" json:"agentId"`
	UserID     string     `gorm:"column:user_id" json:"userId"`
	Name       string     `gorm:"column:name" json:"name"`
	Status     *bool      `gorm:"column:status" json:"status"`
	UsageCount int        `gorm:"column:usage_count" json:"usageCount"`
	LastUsedAt *time.Time `gorm:"column:last_used_at" json:"lastUsedAt"`
	ExpiresAt  *time.Time `gorm:"column:expires_at" json:"expiresAt"`

	entity.BaseEntity
}

func (ApiKeyEntity) TableName() string {
	return "api_keys"
}

// GenerateApiKey 生成API Key，格式：ak_{agentId}_{随机字符串}
func (e *ApiKeyEntity) GenerateApiKey() {
	randomStr := strings.ReplaceAll(uuid.New().String(), "-", "")[:16]
	e.ApiKey = "ak_" + e.AgentID + "_" + randomStr
}

// IsExpired 检查API Key是否过期
func (e *ApiKeyEntity) IsExpired() bool {
	return e.ExpiresAt != nil && time.Now().After(*e.ExpiresAt)
}

// IsAvailable 检查API Key是否可用
func (e *ApiKeyEntity) IsAvailable() bool {
	return e.Status != nil && *e.Status && !e.IsExpired()
}

// IncrementUsage 增加使用次数
func (e *ApiKeyEntity) IncrementUsage() {
	e.UsageCount++
	now := time.Now()
	e.LastUsedAt = &now
}

// ---- 仓储接口 ----

type ApiKeyRepository interface {
	FindByID(id string) (*ApiKeyEntity, error)
	FindByIDAndUserID(id, userID string) (*ApiKeyEntity, error)
	FindByApiKey(apiKey string) (*ApiKeyEntity, error)
	FindByUserID(userID string, name *string, status *bool, agentID *string) ([]ApiKeyEntity, error)
	FindByAgentIDAndUserID(agentID, userID string) ([]ApiKeyEntity, error)
	Create(entity *ApiKeyEntity) error
	Update(entity *ApiKeyEntity) error
	UpdateStatus(id, userID string, status bool) error
	UpdateUsage(apiKey string) error
	DeleteByIDAndUserID(id, userID string) error
}

// ---- GORM 实现 ----

type ApiKeyRepositoryImpl struct {
	DB *gorm.DB
}

func NewApiKeyRepository(db *gorm.DB) ApiKeyRepository {
	return &ApiKeyRepositoryImpl{DB: db}
}

func (r *ApiKeyRepositoryImpl) FindByID(id string) (*ApiKeyEntity, error) {
	var e ApiKeyEntity
	result := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *ApiKeyRepositoryImpl) FindByIDAndUserID(id, userID string) (*ApiKeyEntity, error) {
	var e ApiKeyEntity
	result := r.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *ApiKeyRepositoryImpl) FindByApiKey(apiKey string) (*ApiKeyEntity, error) {
	var e ApiKeyEntity
	result := r.DB.Where("api_key = ? AND deleted_at IS NULL", apiKey).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *ApiKeyRepositoryImpl) FindByUserID(userID string, name *string, status *bool, agentID *string) ([]ApiKeyEntity, error) {
	var entities []ApiKeyEntity
	query := r.DB.Where("user_id = ? AND deleted_at IS NULL", userID)
	if name != nil && *name != "" {
		query = query.Where("name LIKE ?", "%"+*name+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if agentID != nil && *agentID != "" {
		query = query.Where("agent_id = ?", *agentID)
	}
	result := query.Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *ApiKeyRepositoryImpl) FindByAgentIDAndUserID(agentID, userID string) ([]ApiKeyEntity, error) {
	var entities []ApiKeyEntity
	result := r.DB.Where("agent_id = ? AND user_id = ? AND deleted_at IS NULL", agentID, userID).
		Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *ApiKeyRepositoryImpl) Create(e *ApiKeyEntity) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return r.DB.Create(e).Error
}

func (r *ApiKeyRepositoryImpl) Update(e *ApiKeyEntity) error {
	return r.DB.Save(e).Error
}

func (r *ApiKeyRepositoryImpl) UpdateStatus(id, userID string, status bool) error {
	result := r.DB.Model(&ApiKeyEntity{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).
		Update("status", status)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *ApiKeyRepositoryImpl) UpdateUsage(apiKey string) error {
	return r.DB.Model(&ApiKeyEntity{}).
		Where("api_key = ? AND deleted_at IS NULL", apiKey).
		Updates(map[string]interface{}{
			"usage_count": gorm.Expr("usage_count + 1"),
			"last_used_at": time.Now(),
		}).Error
}

func (r *ApiKeyRepositoryImpl) DeleteByIDAndUserID(id, userID string) error {
	now := time.Now()
	result := r.DB.Model(&ApiKeyEntity{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).
		Update("deleted_at", now)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// ---- 领域服务 ----

type DomainService struct {
	repo ApiKeyRepository
}

func NewDomainService(repo ApiKeyRepository) *DomainService {
	return &DomainService{repo: repo}
}

func (s *DomainService) CreateApiKey(e *ApiKeyEntity) (*ApiKeyEntity, error) {
	e.GenerateApiKey()
	if err := s.repo.Create(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *DomainService) FindByApiKey(apiKey string) (*ApiKeyEntity, error) {
	return s.repo.FindByApiKey(apiKey)
}

func (s *DomainService) ValidateApiKey(apiKey string) (*ApiKeyEntity, error) {
	e, err := s.repo.FindByApiKey(apiKey)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, exception.NewBusinessException("无效的API Key")
	}
	if !e.IsAvailable() {
		return nil, exception.NewBusinessException("API Key已禁用或过期")
	}
	return e, nil
}

func (s *DomainService) UpdateUsage(apiKey string) error {
	return s.repo.UpdateUsage(apiKey)
}

func (s *DomainService) GetUserApiKeys(userID string, name *string, status *bool, agentID *string) ([]ApiKeyEntity, error) {
	return s.repo.FindByUserID(userID, name, status, agentID)
}

func (s *DomainService) GetAgentApiKeys(agentID, userID string) ([]ApiKeyEntity, error) {
	return s.repo.FindByAgentIDAndUserID(agentID, userID)
}

func (s *DomainService) GetApiKey(apiKeyID, userID string) (*ApiKeyEntity, error) {
	e, err := s.repo.FindByIDAndUserID(apiKeyID, userID)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, exception.NewBusinessException("API密钥不存在: " + apiKeyID)
	}
	return e, nil
}

func (s *DomainService) UpdateStatus(apiKeyID, userID string, status bool) error {
	return s.repo.UpdateStatus(apiKeyID, userID, status)
}

func (s *DomainService) DeleteApiKey(apiKeyID, userID string) error {
	return s.repo.DeleteByIDAndUserID(apiKeyID, userID)
}

func (s *DomainService) ResetApiKey(apiKeyID, userID string) (*ApiKeyEntity, error) {
	e, err := s.GetApiKey(apiKeyID, userID)
	if err != nil {
		return nil, err
	}
	e.GenerateApiKey()
	e.UsageCount = 0
	e.LastUsedAt = nil
	if err := s.repo.Update(e); err != nil {
		return nil, err
	}
	return e, nil
}
