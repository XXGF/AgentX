package agent

import (
	"gorm.io/gorm"
)

// AgentRepository Agent仓储接口
type AgentRepository interface {
	FindByID(id string) (*AgentEntity, error)
	FindByIDAndUserID(id, userID string) (*AgentEntity, error)
	FindByUserID(userID string, name string) ([]AgentEntity, error)
	FindByIDs(ids []string) ([]AgentEntity, error)
	Create(entity *AgentEntity) error
	Update(entity *AgentEntity) error
	UpdateByIDAndUserID(entity *AgentEntity) error
	DeleteByIDAndUserID(id, userID string) error
	Count(wrapper map[string]interface{}) (int64, error)
	FindPage(page, pageSize int, keyword string, enabled *bool) ([]AgentEntity, int64, error)
}

// AgentVersionRepository Agent版本仓储接口
type AgentVersionRepository interface {
	FindByID(id string) (*AgentVersionEntity, error)
	FindByAgentID(agentID string) ([]AgentVersionEntity, error)
	FindByAgentIDOrderByCreatedAtDesc(agentID string) ([]AgentVersionEntity, error)
	FindLatestByAgentID(agentID string) (*AgentVersionEntity, error)
	FindLatestPublishedByAgentID(agentID string) (*AgentVersionEntity, error)
	FindLatestVersionsByStatus(status *int) ([]AgentVersionEntity, error)
	FindLatestVersionsByNameAndStatus(name string, status int) ([]AgentVersionEntity, error)
	FindByAgentIDs(agentIDs []string) ([]AgentVersionEntity, error)
	Create(entity *AgentVersionEntity) error
	Update(entity *AgentVersionEntity) error
	DeleteByAgentIDAndUserID(agentID, userID string) error
	CountByStatus(status int) (int64, error)
}

// AgentWorkspaceRepository Agent工作区仓储接口
type AgentWorkspaceRepository interface {
	FindByAgentIDAndUserID(agentID, userID string) (*AgentWorkspaceEntity, error)
	FindByUserID(userID string) ([]AgentWorkspaceEntity, error)
	Exists(agentID, userID string) (bool, error)
	Create(entity *AgentWorkspaceEntity) error
	Update(entity *AgentWorkspaceEntity) error
	DeleteByAgentIDAndUserID(agentID, userID string) error
	FindByAgentIDsAndUserID(agentIDs []string, userID string) ([]AgentWorkspaceEntity, error)
}

// ---- GORM 实现 ----

// AgentRepositoryImpl GORM实现
type AgentRepositoryImpl struct {
	DB *gorm.DB
}

func NewAgentRepository(db *gorm.DB) AgentRepository {
	return &AgentRepositoryImpl{DB: db}
}

func (r *AgentRepositoryImpl) FindByID(id string) (*AgentEntity, error) {
	var entity AgentEntity
	result := r.DB.First(&entity, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *AgentRepositoryImpl) FindByIDAndUserID(id, userID string) (*AgentEntity, error) {
	var entity AgentEntity
	result := r.DB.Where("id = ? AND user_id = ?", id, userID).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *AgentRepositoryImpl) FindByUserID(userID string, name string) ([]AgentEntity, error) {
	var entities []AgentEntity
	query := r.DB.Where("user_id = ?", userID)
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	result := query.Order("updated_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *AgentRepositoryImpl) FindByIDs(ids []string) ([]AgentEntity, error) {
	var entities []AgentEntity
	if len(ids) == 0 {
		return entities, nil
	}
	result := r.DB.Where("id IN ?", ids).Find(&entities)
	return entities, result.Error
}

func (r *AgentRepositoryImpl) Create(entity *AgentEntity) error {
	return r.DB.Create(entity).Error
}

func (r *AgentRepositoryImpl) Update(entity *AgentEntity) error {
	return r.DB.Save(entity).Error
}

func (r *AgentRepositoryImpl) UpdateByIDAndUserID(entity *AgentEntity) error {
	result := r.DB.Where("id = ? AND user_id = ?", entity.ID, entity.UserID).Updates(entity)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *AgentRepositoryImpl) DeleteByIDAndUserID(id, userID string) error {
	result := r.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&AgentEntity{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *AgentRepositoryImpl) Count(wrapper map[string]interface{}) (int64, error) {
	var count int64
	query := r.DB.Model(&AgentEntity{})
	for k, v := range wrapper {
		query = query.Where(k+" = ?", v)
	}
	result := query.Count(&count)
	return count, result.Error
}

func (r *AgentRepositoryImpl) FindPage(page, pageSize int, keyword string, enabled *bool) ([]AgentEntity, int64, error) {
	var entities []AgentEntity
	var total int64

	query := r.DB.Model(&AgentEntity{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	result := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&entities)
	return entities, total, result.Error
}

// AgentVersionRepositoryImpl GORM实现
type AgentVersionRepositoryImpl struct {
	DB *gorm.DB
}

func NewAgentVersionRepository(db *gorm.DB) AgentVersionRepository {
	return &AgentVersionRepositoryImpl{DB: db}
}

func (r *AgentVersionRepositoryImpl) FindByID(id string) (*AgentVersionEntity, error) {
	var entity AgentVersionEntity
	result := r.DB.First(&entity, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *AgentVersionRepositoryImpl) FindByAgentID(agentID string) ([]AgentVersionEntity, error) {
	var entities []AgentVersionEntity
	result := r.DB.Where("agent_id = ?", agentID).Find(&entities)
	return entities, result.Error
}

func (r *AgentVersionRepositoryImpl) FindByAgentIDOrderByCreatedAtDesc(agentID string) ([]AgentVersionEntity, error) {
	var entities []AgentVersionEntity
	result := r.DB.Where("agent_id = ?", agentID).Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *AgentVersionRepositoryImpl) FindLatestByAgentID(agentID string) (*AgentVersionEntity, error) {
	var entity AgentVersionEntity
	result := r.DB.Where("agent_id = ?", agentID).Order("published_at DESC").Limit(1).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *AgentVersionRepositoryImpl) FindLatestPublishedByAgentID(agentID string) (*AgentVersionEntity, error) {
	var entity AgentVersionEntity
	result := r.DB.Where("agent_id = ? AND publish_status = ?", agentID, int(PublishStatusPublished)).
		Order("published_at DESC").Limit(1).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *AgentVersionRepositoryImpl) FindLatestVersionsByStatus(status *int) ([]AgentVersionEntity, error) {
	var entities []AgentVersionEntity
	// 子查询：每个agent_id的最新版本
	subQuery := r.DB.Model(&AgentVersionEntity{}).
		Select("agent_id, MAX(published_at) as max_published_at").
		Group("agent_id")

	if status != nil {
		subQuery = subQuery.Where("publish_status = ?", *status)
	}

	query := r.DB.Table("agent_versions AS av").
		Joins("INNER JOIN (?) AS latest ON av.agent_id = latest.agent_id AND av.published_at = latest.max_published_at", subQuery)

	if status != nil {
		query = query.Where("av.publish_status = ?", *status)
	}

	result := query.Find(&entities)
	return entities, result.Error
}

func (r *AgentVersionRepositoryImpl) FindLatestVersionsByNameAndStatus(name string, status int) ([]AgentVersionEntity, error) {
	var entities []AgentVersionEntity
	subQuery := r.DB.Model(&AgentVersionEntity{}).
		Select("agent_id, MAX(published_at) as max_published_at").
		Where("publish_status = ?", status).
		Group("agent_id")

	query := r.DB.Table("agent_versions AS av").
		Joins("INNER JOIN (?) AS latest ON av.agent_id = latest.agent_id AND av.published_at = latest.max_published_at", subQuery).
		Where("av.publish_status = ?", status)

	if name != "" {
		query = query.Where("av.name LIKE ?", "%"+name+"%")
	}

	result := query.Find(&entities)
	return entities, result.Error
}

func (r *AgentVersionRepositoryImpl) FindByAgentIDs(agentIDs []string) ([]AgentVersionEntity, error) {
	var entities []AgentVersionEntity
	if len(agentIDs) == 0 {
		return entities, nil
	}
	result := r.DB.Where("agent_id IN ?", agentIDs).Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *AgentVersionRepositoryImpl) Create(entity *AgentVersionEntity) error {
	return r.DB.Create(entity).Error
}

func (r *AgentVersionRepositoryImpl) Update(entity *AgentVersionEntity) error {
	return r.DB.Save(entity).Error
}

func (r *AgentVersionRepositoryImpl) DeleteByAgentIDAndUserID(agentID, userID string) error {
	result := r.DB.Where("agent_id = ? AND user_id = ?", agentID, userID).Delete(&AgentVersionEntity{})
	return result.Error
}

func (r *AgentVersionRepositoryImpl) CountByStatus(status int) (int64, error) {
	var count int64
	result := r.DB.Model(&AgentVersionEntity{}).Where("publish_status = ?", status).Count(&count)
	return count, result.Error
}

// AgentWorkspaceRepositoryImpl GORM实现
type AgentWorkspaceRepositoryImpl struct {
	DB *gorm.DB
}

func NewAgentWorkspaceRepository(db *gorm.DB) AgentWorkspaceRepository {
	return &AgentWorkspaceRepositoryImpl{DB: db}
}

func (r *AgentWorkspaceRepositoryImpl) FindByAgentIDAndUserID(agentID, userID string) (*AgentWorkspaceEntity, error) {
	var entity AgentWorkspaceEntity
	result := r.DB.Where("agent_id = ? AND user_id = ?", agentID, userID).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *AgentWorkspaceRepositoryImpl) FindByUserID(userID string) ([]AgentWorkspaceEntity, error) {
	var entities []AgentWorkspaceEntity
	result := r.DB.Where("user_id = ?", userID).Find(&entities)
	return entities, result.Error
}

func (r *AgentWorkspaceRepositoryImpl) Exists(agentID, userID string) (bool, error) {
	var count int64
	result := r.DB.Model(&AgentWorkspaceEntity{}).Where("agent_id = ? AND user_id = ?", agentID, userID).Count(&count)
	return count > 0, result.Error
}

func (r *AgentWorkspaceRepositoryImpl) Create(entity *AgentWorkspaceEntity) error {
	return r.DB.Create(entity).Error
}

func (r *AgentWorkspaceRepositoryImpl) Update(entity *AgentWorkspaceEntity) error {
	result := r.DB.Where("agent_id = ? AND user_id = ?", entity.AgentID, entity.UserID).Updates(entity)
	return result.Error
}

func (r *AgentWorkspaceRepositoryImpl) DeleteByAgentIDAndUserID(agentID, userID string) error {
	result := r.DB.Where("agent_id = ? AND user_id = ?", agentID, userID).Delete(&AgentWorkspaceEntity{})
	return result.Error
}

func (r *AgentWorkspaceRepositoryImpl) FindByAgentIDsAndUserID(agentIDs []string, userID string) ([]AgentWorkspaceEntity, error) {
	var entities []AgentWorkspaceEntity
	if len(agentIDs) == 0 {
		return entities, nil
	}
	result := r.DB.Where("user_id = ? AND agent_id IN ?", userID, agentIDs).Find(&entities)
	return entities, result.Error
}
