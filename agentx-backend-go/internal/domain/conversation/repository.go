package conversation

import (
	"time"

	"gorm.io/gorm"
)

// SessionRepository 会话仓储接口
type SessionRepository interface {
	FindByID(id string) (*SessionEntity, error)
	FindByIDAndUserID(id, userID string) (*SessionEntity, error)
	FindByAgentIDAndUserID(agentID, userID string) ([]SessionEntity, error)
	Create(entity *SessionEntity) error
	Update(entity *SessionEntity) error
	DeleteByIDAndUserID(id, userID string) error
	DeleteByIDs(ids []string) error
}

// MessageRepository 消息仓储接口
type MessageRepository interface {
	FindByID(id string) (*MessageEntity, error)
	FindByIDs(ids []string) ([]MessageEntity, error)
	FindBySessionID(sessionID string) ([]MessageEntity, error)
	FindBySessionIDExcludeRole(sessionID string, excludeRole Role) ([]MessageEntity, error)
	CountBySessionID(sessionID string) (int64, error)
	Create(entity *MessageEntity) error
	CreateBatch(entities []MessageEntity) error
	Update(entity *MessageEntity) error
	DeleteBySessionID(sessionID string) error
	DeleteBySessionIDs(sessionIDs []string) error
}

// ContextRepository 上下文仓储接口
type ContextRepository interface {
	FindBySessionID(sessionID string) (*ContextEntity, error)
	InsertOrUpdate(entity *ContextEntity) error
}

// ---- GORM 实现 ----

// SessionRepositoryImpl GORM实现
type SessionRepositoryImpl struct {
	DB *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &SessionRepositoryImpl{DB: db}
}

func (r *SessionRepositoryImpl) FindByID(id string) (*SessionEntity, error) {
	var entity SessionEntity
	result := r.DB.First(&entity, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *SessionRepositoryImpl) FindByIDAndUserID(id, userID string) (*SessionEntity, error) {
	var entity SessionEntity
	result := r.DB.Where("id = ? AND user_id = ?", id, userID).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *SessionRepositoryImpl) FindByAgentIDAndUserID(agentID, userID string) ([]SessionEntity, error) {
	var entities []SessionEntity
	result := r.DB.Where("agent_id = ? AND user_id = ?", agentID, userID).
		Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *SessionRepositoryImpl) Create(entity *SessionEntity) error {
	return r.DB.Create(entity).Error
}

func (r *SessionRepositoryImpl) Update(entity *SessionEntity) error {
	return r.DB.Save(entity).Error
}

func (r *SessionRepositoryImpl) DeleteByIDAndUserID(id, userID string) error {
	now := time.Now()
	result := r.DB.Model(&SessionEntity{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("deleted_at", now)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *SessionRepositoryImpl) DeleteByIDs(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	now := time.Now()
	return r.DB.Model(&SessionEntity{}).Where("id IN ?", ids).Update("deleted_at", now).Error
}

// MessageRepositoryImpl GORM实现
type MessageRepositoryImpl struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &MessageRepositoryImpl{DB: db}
}

func (r *MessageRepositoryImpl) FindByID(id string) (*MessageEntity, error) {
	var entity MessageEntity
	result := r.DB.First(&entity, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *MessageRepositoryImpl) FindByIDs(ids []string) ([]MessageEntity, error) {
	var entities []MessageEntity
	if len(ids) == 0 {
		return entities, nil
	}
	result := r.DB.Where("id IN ?", ids).Find(&entities)
	return entities, result.Error
}

func (r *MessageRepositoryImpl) FindBySessionID(sessionID string) ([]MessageEntity, error) {
	var entities []MessageEntity
	result := r.DB.Where("session_id = ?", sessionID).Order("created_at ASC").Find(&entities)
	return entities, result.Error
}

func (r *MessageRepositoryImpl) FindBySessionIDExcludeRole(sessionID string, excludeRole Role) ([]MessageEntity, error) {
	var entities []MessageEntity
	result := r.DB.Where("session_id = ? AND role != ?", sessionID, string(excludeRole)).
		Order("created_at ASC").Find(&entities)
	return entities, result.Error
}

func (r *MessageRepositoryImpl) CountBySessionID(sessionID string) (int64, error) {
	var count int64
	result := r.DB.Model(&MessageEntity{}).Where("session_id = ?", sessionID).Count(&count)
	return count, result.Error
}

func (r *MessageRepositoryImpl) Create(entity *MessageEntity) error {
	return r.DB.Create(entity).Error
}

func (r *MessageRepositoryImpl) CreateBatch(entities []MessageEntity) error {
	if len(entities) == 0 {
		return nil
	}
	return r.DB.Create(&entities).Error
}

func (r *MessageRepositoryImpl) Update(entity *MessageEntity) error {
	return r.DB.Save(entity).Error
}

func (r *MessageRepositoryImpl) DeleteBySessionID(sessionID string) error {
	now := time.Now()
	return r.DB.Model(&MessageEntity{}).Where("session_id = ?", sessionID).Update("deleted_at", now).Error
}

func (r *MessageRepositoryImpl) DeleteBySessionIDs(sessionIDs []string) error {
	if len(sessionIDs) == 0 {
		return nil
	}
	now := time.Now()
	return r.DB.Model(&MessageEntity{}).Where("session_id IN ?", sessionIDs).Update("deleted_at", now).Error
}

// ContextRepositoryImpl GORM实现
type ContextRepositoryImpl struct {
	DB *gorm.DB
}

func NewContextRepository(db *gorm.DB) ContextRepository {
	return &ContextRepositoryImpl{DB: db}
}

func (r *ContextRepositoryImpl) FindBySessionID(sessionID string) (*ContextEntity, error) {
	var entity ContextEntity
	result := r.DB.Where("session_id = ?", sessionID).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *ContextRepositoryImpl) InsertOrUpdate(entity *ContextEntity) error {
	// 先查找是否存在
	existing, err := r.FindBySessionID(entity.SessionID)
	if err != nil {
		return err
	}
	if existing != nil {
		entity.ID = existing.ID
		return r.DB.Save(entity).Error
	}
	return r.DB.Create(entity).Error
}
