package repository

import (
	"errors"
	"fmt"

	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	"gorm.io/gorm"
)

// BaseRepository 基础仓储（对应 Java 的 MyBatisPlusExtRepository）
// 提供带检查的 CRUD 操作，失败时抛出 BusinessException
type BaseRepository[T any] struct {
	DB *gorm.DB
}

// NewBaseRepository 创建基础仓储实例
func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{DB: db}
}

// CheckedCreate 创建记录，失败时返回 BusinessException（对应 Java 的 checkInsert）
func (r *BaseRepository[T]) CheckedCreate(entity *T) error {
	result := r.DB.Create(entity)
	if result.Error != nil {
		return exception.NewBusinessExceptionWithCause("数据创建失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return exception.NewBusinessException("数据创建失败")
	}
	return nil
}

// CheckedUpdate 更新记录，失败时返回 BusinessException（对应 Java 的 checkedUpdate）
func (r *BaseRepository[T]) CheckedUpdate(db *gorm.DB) error {
	result := db.Updates(new(T))
	if result.Error != nil {
		return exception.NewBusinessExceptionWithCause("数据更新失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return exception.NewBusinessException("数据更新失败")
	}
	return nil
}

// CheckedUpdateByID 根据ID更新记录（对应 Java 的 checkedUpdateById）
func (r *BaseRepository[T]) CheckedUpdateByID(entity *T) error {
	result := r.DB.Save(entity)
	if result.Error != nil {
		return exception.NewBusinessExceptionWithCause("数据更新失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return exception.NewBusinessException("数据更新失败")
	}
	return nil
}

// CheckedDelete 删除记录，失败时返回 BusinessException（对应 Java 的 checkedDelete）
func (r *BaseRepository[T]) CheckedDelete(db *gorm.DB) error {
	result := db.Delete(new(T))
	if result.Error != nil {
		return exception.NewBusinessExceptionWithCause("数据删除失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return exception.NewBusinessException("数据删除失败")
	}
	return nil
}

// FindByID 根据ID查询（对应 Java 的 selectById）
func (r *BaseRepository[T]) FindByID(id string) (*T, error) {
	var entity T
	result := r.DB.First(&entity, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, nil
}

// GetByID 根据ID查询，不存在则返回异常（对应 Java 的 getXxx 命名规范）
func (r *BaseRepository[T]) GetByID(id string) (*T, error) {
	entity, err := r.FindByID(id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, exception.NewEntityNotFoundException(fmt.Sprintf("ID为 %s 的记录不存在", id))
	}
	return entity, nil
}

// FindAll 查询所有记录
func (r *BaseRepository[T]) FindAll() ([]T, error) {
	var entities []T
	result := r.DB.Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	return entities, nil
}

// Exists 检查记录是否存在
func (r *BaseRepository[T]) Exists(db *gorm.DB) (bool, error) {
	var count int64
	result := db.Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// Page 分页查询
func (r *BaseRepository[T]) Page(db *gorm.DB, page, pageSize int) ([]T, int64, error) {
	var entities []T
	var total int64

	// 先查总数
	countDB := db.Session(&gorm.Session{})
	if err := countDB.Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 再查分页数据
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
