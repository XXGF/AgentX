package entity

import (
	"time"

	"gorm.io/gorm"
)

// Operator 操作者类型（对应 Java 的 org.xhy.infrastructure.entity.Operator 枚举）
type Operator int

const (
	// OperatorUser 普通用户操作
	OperatorUser Operator = iota
	// OperatorAdmin 管理员操作
	OperatorAdmin
)

// NeedCheckUserId 是否需要检查用户ID（普通用户需要，管理员不需要）
func (o Operator) NeedCheckUserId() bool {
	return o == OperatorUser
}

// BaseEntity 基础实体（对应 Java 的 org.xhy.infrastructure.entity.BaseEntity）
// 所有实体都应嵌入此结构体，提供统一的审计字段和软删除支持
type BaseEntity struct {
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`

	// OperatedBy 不映射到数据库，仅用于运行时权限判断（对应 Java 的 @TableField(exist = false)）
	OperatedBy Operator `gorm:"-" json:"-"`
}

// SetAdmin 设置为管理员操作（对应 Java 的 entity.setAdmin()）
func (b *BaseEntity) SetAdmin() {
	b.OperatedBy = OperatorAdmin
}

// NeedCheckUserId 是否需要检查用户ID（对应 Java 的 entity.needCheckUserId()）
func (b *BaseEntity) NeedCheckUserId() bool {
	return b.OperatedBy.NeedCheckUserId()
}
