package llm

import (
	"gorm.io/gorm"
)

// ProviderRepository 服务商仓储接口
type ProviderRepository interface {
	FindByID(id string) (*ProviderEntity, error)
	FindByIDAndUserID(id, userID string) (*ProviderEntity, error)
	FindByUserID(userID string) ([]ProviderEntity, error)
	FindByUserIDOrOfficial(userID string) ([]ProviderEntity, error)
	FindOfficial() ([]ProviderEntity, error)
	FindCustomByUserID(userID string) ([]ProviderEntity, error)
	Create(entity *ProviderEntity) error
	Update(entity *ProviderEntity) error
	ToggleStatus(id, userID string) error
	Delete(id string, userID string, checkUser bool) error
}

// ModelRepository 模型仓储接口
type ModelRepository interface {
	FindByID(id string) (*ModelEntity, error)
	FindByProviderIDs(providerIDs []string) ([]ModelEntity, error)
	FindActiveByProviderIDs(providerIDs []string) ([]ModelEntity, error)
	FindByProviderIDAndUserID(providerID, userID string) ([]ModelEntity, error)
	FindActiveByProviderIDAndUserID(providerID, userID string) ([]ModelEntity, error)
	FindAllActive() ([]ModelEntity, error)
	FindByIDs(ids []string) ([]ModelEntity, error)
	Create(entity *ModelEntity) error
	Update(entity *ModelEntity) error
	ToggleStatus(id, userID string) error
	DeleteByProviderID(providerID string) (int64, error)
	Delete(id string, userID string, checkUser bool) error
}

// ---- GORM 实现 ----

// ProviderRepositoryImpl GORM实现
type ProviderRepositoryImpl struct {
	DB *gorm.DB
}

func NewProviderRepository(db *gorm.DB) ProviderRepository {
	return &ProviderRepositoryImpl{DB: db}
}

func (r *ProviderRepositoryImpl) FindByID(id string) (*ProviderEntity, error) {
	var entity ProviderEntity
	result := r.DB.First(&entity, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *ProviderRepositoryImpl) FindByIDAndUserID(id, userID string) (*ProviderEntity, error) {
	var entity ProviderEntity
	result := r.DB.Where("id = ? AND user_id = ?", id, userID).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *ProviderRepositoryImpl) FindByUserID(userID string) ([]ProviderEntity, error) {
	var entities []ProviderEntity
	result := r.DB.Where("user_id = ?", userID).Find(&entities)
	return entities, result.Error
}

func (r *ProviderRepositoryImpl) FindByUserIDOrOfficial(userID string) ([]ProviderEntity, error) {
	var entities []ProviderEntity
	result := r.DB.Where("user_id = ? OR is_official = ?", userID, true).Find(&entities)
	return entities, result.Error
}

func (r *ProviderRepositoryImpl) FindOfficial() ([]ProviderEntity, error) {
	var entities []ProviderEntity
	result := r.DB.Where("is_official = ?", true).Find(&entities)
	return entities, result.Error
}

func (r *ProviderRepositoryImpl) FindCustomByUserID(userID string) ([]ProviderEntity, error) {
	var entities []ProviderEntity
	result := r.DB.Where("user_id = ? AND is_official = ?", userID, false).Find(&entities)
	return entities, result.Error
}

func (r *ProviderRepositoryImpl) Create(entity *ProviderEntity) error {
	return r.DB.Create(entity).Error
}

func (r *ProviderRepositoryImpl) Update(entity *ProviderEntity) error {
	query := r.DB.Where("id = ?", entity.ID)
	if entity.NeedCheckUserId() {
		query = query.Where("user_id = ?", entity.UserID)
	}
	result := query.Updates(entity)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *ProviderRepositoryImpl) ToggleStatus(id, userID string) error {
	result := r.DB.Model(&ProviderEntity{}).
		Where("id = ? AND user_id = ?", id, userID).
		UpdateColumn("status", gorm.Expr("NOT status"))
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *ProviderRepositoryImpl) Delete(id string, userID string, checkUser bool) error {
	query := r.DB.Where("id = ?", id)
	if checkUser {
		query = query.Where("user_id = ?", userID)
	}
	result := query.Delete(&ProviderEntity{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// ModelRepositoryImpl GORM实现
type ModelRepositoryImpl struct {
	DB *gorm.DB
}

func NewModelRepository(db *gorm.DB) ModelRepository {
	return &ModelRepositoryImpl{DB: db}
}

func (r *ModelRepositoryImpl) FindByID(id string) (*ModelEntity, error) {
	var entity ModelEntity
	result := r.DB.First(&entity, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *ModelRepositoryImpl) FindByProviderIDs(providerIDs []string) ([]ModelEntity, error) {
	var entities []ModelEntity
	if len(providerIDs) == 0 {
		return entities, nil
	}
	result := r.DB.Where("provider_id IN ?", providerIDs).Find(&entities)
	return entities, result.Error
}

func (r *ModelRepositoryImpl) FindActiveByProviderIDs(providerIDs []string) ([]ModelEntity, error) {
	var entities []ModelEntity
	if len(providerIDs) == 0 {
		return entities, nil
	}
	result := r.DB.Where("provider_id IN ? AND status = ?", providerIDs, true).Find(&entities)
	return entities, result.Error
}

func (r *ModelRepositoryImpl) FindByProviderIDAndUserID(providerID, userID string) ([]ModelEntity, error) {
	var entities []ModelEntity
	result := r.DB.Where("provider_id = ? AND user_id = ?", providerID, userID).Find(&entities)
	return entities, result.Error
}

func (r *ModelRepositoryImpl) FindActiveByProviderIDAndUserID(providerID, userID string) ([]ModelEntity, error) {
	var entities []ModelEntity
	result := r.DB.Where("provider_id = ? AND user_id = ? AND status = ?", providerID, userID, true).Find(&entities)
	return entities, result.Error
}

func (r *ModelRepositoryImpl) FindAllActive() ([]ModelEntity, error) {
	var entities []ModelEntity
	result := r.DB.Where("status = ?", true).Find(&entities)
	return entities, result.Error
}

func (r *ModelRepositoryImpl) FindByIDs(ids []string) ([]ModelEntity, error) {
	var entities []ModelEntity
	if len(ids) == 0 {
		return entities, nil
	}
	result := r.DB.Where("id IN ?", ids).Find(&entities)
	return entities, result.Error
}

func (r *ModelRepositoryImpl) Create(entity *ModelEntity) error {
	return r.DB.Create(entity).Error
}

func (r *ModelRepositoryImpl) Update(entity *ModelEntity) error {
	query := r.DB.Where("id = ?", entity.ID)
	if entity.UserID != "" {
		query = query.Where("user_id = ?", entity.UserID)
	}
	result := query.Updates(entity)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *ModelRepositoryImpl) ToggleStatus(id, userID string) error {
	result := r.DB.Model(&ModelEntity{}).
		Where("id = ? AND user_id = ?", id, userID).
		UpdateColumn("status", gorm.Expr("NOT status"))
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *ModelRepositoryImpl) DeleteByProviderID(providerID string) (int64, error) {
	result := r.DB.Where("provider_id = ?", providerID).Delete(&ModelEntity{})
	return result.RowsAffected, result.Error
}

func (r *ModelRepositoryImpl) Delete(id string, userID string, checkUser bool) error {
	query := r.DB.Where("id = ?", id)
	if checkUser {
		query = query.Where("user_id = ?", userID)
	}
	result := query.Delete(&ModelEntity{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
