package memory

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/tool"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
	"gorm.io/gorm"
)

// MemoryType 记忆类型
type MemoryType string

const (
	MemoryTypeProfile  MemoryType = "PROFILE"
	MemoryTypeTask     MemoryType = "TASK"
	MemoryTypeFact     MemoryType = "FACT"
	MemoryTypeEpisodic MemoryType = "EPISODIC"
)

// SafeOf 安全转换记忆类型
func SafeOf(name string) MemoryType {
	if name == "" {
		return MemoryTypeFact
	}
	upper := strings.ToUpper(strings.TrimSpace(name))
	switch upper {
	case "PROFILE":
		return MemoryTypeProfile
	case "TASK":
		return MemoryTypeTask
	case "FACT":
		return MemoryTypeFact
	case "EPISODIC":
		return MemoryTypeEpisodic
	default:
		return MemoryTypeFact
	}
}

const MemoryStatusActive = 1

// MemoryItemEntity 记忆条目实体
type MemoryItemEntity struct {
	ID              string       `gorm:"column:id;primaryKey" json:"id"`
	UserID          string       `gorm:"column:user_id" json:"userId"`
	Type            string       `gorm:"column:type" json:"type"`
	Text            string       `gorm:"column:text" json:"text"`
	Data            tool.JSONMap `gorm:"column:data;type:jsonb" json:"data"`
	Importance      *float32     `gorm:"column:importance" json:"importance"`
	Tags            tool.JSONStringList `gorm:"column:tags;type:jsonb" json:"tags"`
	SourceSessionID string       `gorm:"column:source_session_id" json:"sourceSessionId"`
	DedupeHash      string       `gorm:"column:dedupe_hash" json:"dedupeHash"`
	Status          int          `gorm:"column:status" json:"status"`

	entity.BaseEntity
}

func (MemoryItemEntity) TableName() string {
	return "memory_items"
}

// CandidateMemory 候选记忆（用于保存）
type CandidateMemory struct {
	Text       string                 `json:"text"`
	Type       *MemoryType            `json:"type,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Importance *float32               `json:"importance,omitempty"`
	Tags       []string               `json:"tags,omitempty"`
}

// MemoryResult 记忆检索结果
type MemoryResult struct {
	ItemID     string     `json:"itemId"`
	Type       MemoryType `json:"type"`
	Text       string     `json:"text"`
	Importance *float32   `json:"importance"`
	Tags       []string   `json:"tags"`
	Score      float64    `json:"score"`
}

// ---- 仓储接口 ----

type MemoryItemRepository interface {
	FindByID(id string) (*MemoryItemEntity, error)
	FindByUserIDAndDedupeHash(userID, hash string) (*MemoryItemEntity, error)
	FindByUserID(userID string, memType *string, limit *int) ([]MemoryItemEntity, error)
	FindByUserIDPaged(userID string, memType *string, page, pageSize int) ([]MemoryItemEntity, int64, error)
	FindByIDs(ids []string) ([]MemoryItemEntity, error)
	Create(entity *MemoryItemEntity) error
	Update(entity *MemoryItemEntity) error
	DeleteByUserIDAndID(userID, id string) error
}

// ---- GORM 实现 ----

type MemoryItemRepositoryImpl struct {
	DB *gorm.DB
}

func NewMemoryItemRepository(db *gorm.DB) MemoryItemRepository {
	return &MemoryItemRepositoryImpl{DB: db}
}

func (r *MemoryItemRepositoryImpl) FindByID(id string) (*MemoryItemEntity, error) {
	var e MemoryItemEntity
	result := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *MemoryItemRepositoryImpl) FindByUserIDAndDedupeHash(userID, hash string) (*MemoryItemEntity, error) {
	var e MemoryItemEntity
	result := r.DB.Where("user_id = ? AND dedupe_hash = ? AND deleted_at IS NULL", userID, hash).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *MemoryItemRepositoryImpl) FindByUserID(userID string, memType *string, limit *int) ([]MemoryItemEntity, error) {
	var entities []MemoryItemEntity
	query := r.DB.Where("user_id = ? AND deleted_at IS NULL", userID)
	if memType != nil && *memType != "" {
		query = query.Where("type = ?", strings.ToUpper(*memType))
	}
	query = query.Order("updated_at DESC")
	if limit != nil && *limit > 0 {
		query = query.Limit(*limit)
	}
	result := query.Find(&entities)
	return entities, result.Error
}

func (r *MemoryItemRepositoryImpl) FindByUserIDPaged(userID string, memType *string, page, pageSize int) ([]MemoryItemEntity, int64, error) {
	var entities []MemoryItemEntity
	var total int64
	query := r.DB.Model(&MemoryItemEntity{}).Where("user_id = ? AND deleted_at IS NULL", userID)
	if memType != nil && *memType != "" {
		query = query.Where("type = ?", strings.ToUpper(*memType))
	}
	query.Count(&total)
	result := query.Order("updated_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities)
	return entities, total, result.Error
}

func (r *MemoryItemRepositoryImpl) FindByIDs(ids []string) ([]MemoryItemEntity, error) {
	if len(ids) == 0 {
		return []MemoryItemEntity{}, nil
	}
	var entities []MemoryItemEntity
	result := r.DB.Where("id IN ? AND deleted_at IS NULL", ids).Find(&entities)
	return entities, result.Error
}

func (r *MemoryItemRepositoryImpl) Create(e *MemoryItemEntity) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return r.DB.Create(e).Error
}

func (r *MemoryItemRepositoryImpl) Update(e *MemoryItemEntity) error {
	return r.DB.Save(e).Error
}

func (r *MemoryItemRepositoryImpl) DeleteByUserIDAndID(userID, id string) error {
	now := time.Now()
	return r.DB.Model(&MemoryItemEntity{}).
		Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, id).
		Update("deleted_at", now).Error
}

// ---- 领域服务 ----

type DomainService struct {
	repo MemoryItemRepository
}

func NewDomainService(repo MemoryItemRepository) *DomainService {
	return &DomainService{repo: repo}
}

// SaveMemories 保存记忆（去重/合并）
// 注意：向量入库部分需要后续集成 embedding 服务，此处仅做数据库持久化
func (s *DomainService) SaveMemories(userID, sessionID string, candidates []CandidateMemory) ([]string, error) {
	if len(candidates) == 0 {
		return []string{}, nil
	}

	var itemIDs []string
	for _, c := range candidates {
		if c.Text == "" {
			continue
		}

		memType := MemoryTypeFact
		if c.Type != nil {
			memType = *c.Type
		}
		normalized := normalizeText(c.Text)
		hash := sha256Hash(normalized)

		existed, err := s.repo.FindByUserIDAndDedupeHash(userID, hash)
		if err != nil {
			return nil, err
		}

		if existed == nil {
			// 新增
			newItem := &MemoryItemEntity{
				UserID:          userID,
				Type:            string(memType),
				Text:            strings.TrimSpace(c.Text),
				Data:            tool.JSONMap(c.Data),
				Importance:      safeImportance(c.Importance),
				Tags:            tool.JSONStringList(c.Tags),
				SourceSessionID: sessionID,
				DedupeHash:      hash,
				Status:          MemoryStatusActive,
			}
			if err := s.repo.Create(newItem); err != nil {
				return nil, err
			}
			itemIDs = append(itemIDs, newItem.ID)
		} else {
			// 合并
			existed.Importance = maxFloat(existed.Importance, c.Importance)
			existed.Tags = mergeTags(existed.Tags, c.Tags)
			existed.Text = pickRichText(existed.Text, c.Text)
			if err := s.repo.Update(existed); err != nil {
				return nil, err
			}
			itemIDs = append(itemIDs, existed.ID)
		}
		// TODO: 向量入库（需要集成 embedding 服务）
	}

	return itemIDs, nil
}

// ListMemories 列出用户记忆
func (s *DomainService) ListMemories(userID string, memType *string, limit *int) ([]MemoryItemEntity, error) {
	return s.repo.FindByUserID(userID, memType, limit)
}

// PageMemories 分页列出用户记忆
func (s *DomainService) PageMemories(userID string, memType *string, page, pageSize int) ([]MemoryItemEntity, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	return s.repo.FindByUserIDPaged(userID, memType, page, pageSize)
}

// Delete 归档（软删除）记忆
func (s *DomainService) Delete(userID, itemID string) error {
	return s.repo.DeleteByUserIDAndID(userID, itemID)
}

// ---- 辅助函数 ----

func normalizeText(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(s))
}

func sha256Hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", h)
}

func safeImportance(f *float32) *float32 {
	if f == nil {
		v := float32(0.5)
		return &v
	}
	v := *f
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	return &v
}

func maxFloat(a, b *float32) *float32 {
	if a == nil && b == nil {
		v := float32(0.5)
		return &v
	}
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if *a > *b {
		return a
	}
	return b
}

func mergeTags(a tool.JSONStringList, b []string) tool.JSONStringList {
	seen := make(map[string]bool)
	var result []string
	for _, t := range a {
		if !seen[t] {
			seen[t] = true
			result = append(result, t)
		}
	}
	for _, t := range b {
		if !seen[t] {
			seen[t] = true
			result = append(result, t)
		}
	}
	return tool.JSONStringList(result)
}

func pickRichText(oldText, newText string) string {
	if newText == "" {
		return oldText
	}
	if oldText == "" {
		return newText
	}
	if len(newText) >= len(oldText) {
		return newText
	}
	return oldText
}
