package tool

import (
	"time"

	"gorm.io/gorm"
)

// ToolRepository 工具仓储接口
type ToolRepository interface {
	FindByID(id string) (*ToolEntity, error)
	FindByIDAndUserID(id, userID string) (*ToolEntity, error)
	FindByUserID(userID string) ([]ToolEntity, error)
	FindByIDs(ids []string) ([]ToolEntity, error)
	Create(entity *ToolEntity) error
	Update(entity *ToolEntity) error
	UpdateFields(id string, fields map[string]interface{}) error
	DeleteByIDAndUserID(id, userID string) error
	Count(conditions map[string]interface{}) (int64, error)
	FindPage(page, pageSize int, conditions map[string]interface{}, keyword string) ([]ToolEntity, int64, error)
}

// ToolVersionRepository 工具版本仓储接口
type ToolVersionRepository interface {
	FindByToolIDAndVersion(toolID, version string) (*ToolVersionEntity, error)
	FindByToolID(toolID string) ([]ToolVersionEntity, error)
	FindLatestByToolID(toolID string) (*ToolVersionEntity, error)
	FindPublicVersions(toolName string) ([]ToolVersionEntity, error)
	Create(entity *ToolVersionEntity) error
	Update(entity *ToolVersionEntity) error
	UpdatePublicStatus(toolID, version, userID string, publicStatus bool) error
}

// UserToolRepository 用户工具仓储接口
type UserToolRepository interface {
	FindByToolIDAndUserID(toolID, userID string) (*UserToolEntity, error)
	FindByUserID(userID string, page, pageSize int) ([]UserToolEntity, int64, error)
	FindByUserIDAndToolIDs(userID string, toolIDs []string) ([]UserToolEntity, error)
	FindByMcpServerNameAndUserID(serverName, userID string) (*UserToolEntity, error)
	CountByToolIDs(toolIDs []string) (map[string]int64, error)
	Create(entity *UserToolEntity) error
	Update(entity *UserToolEntity) error
	DeleteByToolIDAndUserID(toolID, userID string) error
	UpdateGlobalStatusByToolID(toolID string, isGlobal bool) error
	CountByMcpServerNameAndUserID(serverName, userID string, excludeToolID *string) (int64, error)
}

// ---- GORM 实现 ----

// ToolRepositoryImpl GORM实现
type ToolRepositoryImpl struct {
	DB *gorm.DB
}

func NewToolRepository(db *gorm.DB) ToolRepository {
	return &ToolRepositoryImpl{DB: db}
}

func (r *ToolRepositoryImpl) FindByID(id string) (*ToolEntity, error) {
	var entity ToolEntity
	result := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *ToolRepositoryImpl) FindByIDAndUserID(id, userID string) (*ToolEntity, error) {
	var entity ToolEntity
	result := r.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *ToolRepositoryImpl) FindByUserID(userID string) ([]ToolEntity, error) {
	var entities []ToolEntity
	result := r.DB.Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("updated_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *ToolRepositoryImpl) FindByIDs(ids []string) ([]ToolEntity, error) {
	if len(ids) == 0 {
		return []ToolEntity{}, nil
	}
	var entities []ToolEntity
	result := r.DB.Where("id IN ? AND deleted_at IS NULL", ids).Find(&entities)
	return entities, result.Error
}

func (r *ToolRepositoryImpl) Create(entity *ToolEntity) error {
	return r.DB.Create(entity).Error
}

func (r *ToolRepositoryImpl) Update(entity *ToolEntity) error {
	return r.DB.Save(entity).Error
}

func (r *ToolRepositoryImpl) UpdateFields(id string, fields map[string]interface{}) error {
	return r.DB.Model(&ToolEntity{}).Where("id = ?", id).Updates(fields).Error
}

func (r *ToolRepositoryImpl) DeleteByIDAndUserID(id, userID string) error {
	now := time.Now()
	result := r.DB.Model(&ToolEntity{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).
		Update("deleted_at", now)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *ToolRepositoryImpl) Count(conditions map[string]interface{}) (int64, error) {
	var count int64
	query := r.DB.Model(&ToolEntity{}).Where("deleted_at IS NULL")
	for k, v := range conditions {
		query = query.Where(k+" = ?", v)
	}
	result := query.Count(&count)
	return count, result.Error
}

func (r *ToolRepositoryImpl) FindPage(page, pageSize int, conditions map[string]interface{}, keyword string) ([]ToolEntity, int64, error) {
	var entities []ToolEntity
	var total int64

	query := r.DB.Model(&ToolEntity{}).Where("deleted_at IS NULL")
	for k, v := range conditions {
		query = query.Where(k+" = ?", v)
	}
	if keyword != "" {
		query = query.Where("(name LIKE ? OR description LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	result := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&entities)
	return entities, total, result.Error
}

// ToolVersionRepositoryImpl GORM实现
type ToolVersionRepositoryImpl struct {
	DB *gorm.DB
}

func NewToolVersionRepository(db *gorm.DB) ToolVersionRepository {
	return &ToolVersionRepositoryImpl{DB: db}
}

func (r *ToolVersionRepositoryImpl) FindByToolIDAndVersion(toolID, version string) (*ToolVersionEntity, error) {
	var entity ToolVersionEntity
	result := r.DB.Where("tool_id = ? AND version = ? AND deleted_at IS NULL", toolID, version).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *ToolVersionRepositoryImpl) FindByToolID(toolID string) ([]ToolVersionEntity, error) {
	var entities []ToolVersionEntity
	result := r.DB.Where("tool_id = ? AND deleted_at IS NULL", toolID).
		Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *ToolVersionRepositoryImpl) FindLatestByToolID(toolID string) (*ToolVersionEntity, error) {
	var entity ToolVersionEntity
	result := r.DB.Where("tool_id = ? AND deleted_at IS NULL", toolID).
		Order("created_at DESC").Limit(1).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *ToolVersionRepositoryImpl) FindPublicVersions(toolName string) ([]ToolVersionEntity, error) {
	var entities []ToolVersionEntity
	query := r.DB.Where("public_status = true AND deleted_at IS NULL")
	if toolName != "" {
		query = query.Where("name LIKE ?", "%"+toolName+"%")
	}
	result := query.Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *ToolVersionRepositoryImpl) Create(entity *ToolVersionEntity) error {
	return r.DB.Create(entity).Error
}

func (r *ToolVersionRepositoryImpl) Update(entity *ToolVersionEntity) error {
	return r.DB.Save(entity).Error
}

func (r *ToolVersionRepositoryImpl) UpdatePublicStatus(toolID, version, userID string, publicStatus bool) error {
	result := r.DB.Model(&ToolVersionEntity{}).
		Where("tool_id = ? AND version = ? AND user_id = ? AND deleted_at IS NULL", toolID, version, userID).
		Update("public_status", publicStatus)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// UserToolRepositoryImpl GORM实现
type UserToolRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserToolRepository(db *gorm.DB) UserToolRepository {
	return &UserToolRepositoryImpl{DB: db}
}

func (r *UserToolRepositoryImpl) FindByToolIDAndUserID(toolID, userID string) (*UserToolEntity, error) {
	var entity UserToolEntity
	result := r.DB.Where("tool_id = ? AND user_id = ? AND deleted_at IS NULL", toolID, userID).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *UserToolRepositoryImpl) FindByUserID(userID string, page, pageSize int) ([]UserToolEntity, int64, error) {
	var entities []UserToolEntity
	var total int64

	query := r.DB.Model(&UserToolEntity{}).Where("user_id = ? AND deleted_at IS NULL", userID)
	query.Count(&total)
	result := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities)
	return entities, total, result.Error
}

func (r *UserToolRepositoryImpl) FindByUserIDAndToolIDs(userID string, toolIDs []string) ([]UserToolEntity, error) {
	if len(toolIDs) == 0 {
		return []UserToolEntity{}, nil
	}
	var entities []UserToolEntity
	result := r.DB.Where("user_id = ? AND tool_id IN ? AND deleted_at IS NULL", userID, toolIDs).Find(&entities)
	return entities, result.Error
}

func (r *UserToolRepositoryImpl) FindByMcpServerNameAndUserID(serverName, userID string) (*UserToolEntity, error) {
	var entity UserToolEntity
	result := r.DB.Where("mcp_server_name = ? AND user_id = ? AND deleted_at IS NULL", serverName, userID).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *UserToolRepositoryImpl) CountByToolIDs(toolIDs []string) (map[string]int64, error) {
	if len(toolIDs) == 0 {
		return map[string]int64{}, nil
	}
	var results []struct {
		ToolID string `gorm:"column:tool_id"`
		Count  int64  `gorm:"column:count"`
	}
	err := r.DB.Model(&UserToolEntity{}).
		Select("tool_id, COUNT(*) as count").
		Where("tool_id IN ? AND deleted_at IS NULL", toolIDs).
		Group("tool_id").
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	m := make(map[string]int64)
	for _, r := range results {
		m[r.ToolID] = r.Count
	}
	return m, nil
}

func (r *UserToolRepositoryImpl) Create(entity *UserToolEntity) error {
	return r.DB.Create(entity).Error
}

func (r *UserToolRepositoryImpl) Update(entity *UserToolEntity) error {
	return r.DB.Save(entity).Error
}

func (r *UserToolRepositoryImpl) DeleteByToolIDAndUserID(toolID, userID string) error {
	now := time.Now()
	result := r.DB.Model(&UserToolEntity{}).
		Where("tool_id = ? AND user_id = ? AND deleted_at IS NULL", toolID, userID).
		Update("deleted_at", now)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *UserToolRepositoryImpl) UpdateGlobalStatusByToolID(toolID string, isGlobal bool) error {
	return r.DB.Model(&UserToolEntity{}).
		Where("tool_id = ? AND deleted_at IS NULL", toolID).
		Update("is_global", isGlobal).Error
}

func (r *UserToolRepositoryImpl) CountByMcpServerNameAndUserID(serverName, userID string, excludeToolID *string) (int64, error) {
	var count int64
	query := r.DB.Model(&UserToolEntity{}).
		Where("mcp_server_name = ? AND user_id = ? AND deleted_at IS NULL", serverName, userID)
	if excludeToolID != nil {
		query = query.Where("tool_id != ?", *excludeToolID)
	}
	result := query.Count(&count)
	return count, result.Error
}
