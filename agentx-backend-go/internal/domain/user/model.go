package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// UserEntity 用户实体（对应 Java 的 org.xhy.domain.user.model.UserEntity）
type UserEntity struct {
	ID            string `gorm:"column:id;primaryKey" json:"id"`
	Nickname      string `gorm:"column:nickname" json:"nickname"`
	Email         string `gorm:"column:email" json:"email"`
	Phone         string `gorm:"column:phone" json:"phone"`
	Password      string `gorm:"column:password" json:"-"` // json中隐藏密码
	GithubID      string `gorm:"column:github_id" json:"githubId"`
	GithubLogin   string `gorm:"column:github_login" json:"githubLogin"`
	AvatarURL     string `gorm:"column:avatar_url" json:"avatarUrl"`
	LoginPlatform string `gorm:"column:login_platform" json:"loginPlatform"`
	IsAdmin       *bool  `gorm:"column:is_admin" json:"isAdmin"`

	entity.BaseEntity
}

// TableName 指定表名
func (UserEntity) TableName() string {
	return "users"
}

// BeforeCreate GORM 钩子：创建前自动生成UUID
func (u *UserEntity) BeforeCreate(tx interface{}) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// Valid 校验实体有效性（对应 Java 的 UserEntity.valid()）
func (u *UserEntity) Valid() error {
	if u.Email == "" && u.Phone == "" && u.GithubID == "" {
		return exception.NewBusinessException("必须使用邮箱、手机号或GitHub账号来作为账号")
	}
	return nil
}

// IsAdminUser 判断是否为管理员（对应 Java 的 UserEntity.isAdmin()）
func (u *UserEntity) IsAdminUser() bool {
	return u.IsAdmin != nil && *u.IsAdmin
}

// UserSettingsConfig 用户设置配置（对应 Java 的 UserSettingsConfig）
type UserSettingsConfig struct {
	DefaultModel          string          `json:"defaultModel,omitempty"`
	DefaultOcrModel       string          `json:"defaultOcrModel,omitempty"`
	DefaultEmbeddingModel string          `json:"defaultEmbeddingModel,omitempty"`
	FallbackConfig        *FallbackConfig `json:"fallbackConfig,omitempty"`
}

// FallbackConfig 降级配置
type FallbackConfig struct {
	// TODO: 后续迁移时补充字段
}

// UserSettingsEntity 用户设置实体（对应 Java 的 UserSettingsEntity）
type UserSettingsEntity struct {
	ID            string              `gorm:"column:id;primaryKey" json:"id"`
	UserID        string              `gorm:"column:user_id" json:"userId"`
	SettingConfig *UserSettingsConfig `gorm:"column:setting_config;serializer:json" json:"settingConfig"`

	entity.BaseEntity
}

// TableName 指定表名
func (UserSettingsEntity) TableName() string {
	return "user_settings"
}

// BeforeCreate GORM 钩子：创建前自动生成UUID
func (u *UserSettingsEntity) BeforeCreate(tx interface{}) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// GetDefaultModelID 获取默认模型ID
func (u *UserSettingsEntity) GetDefaultModelID() string {
	if u.SettingConfig == nil {
		return ""
	}
	return u.SettingConfig.DefaultModel
}

// SetDefaultModelID 设置默认模型ID
func (u *UserSettingsEntity) SetDefaultModelID(modelID string) {
	if u.SettingConfig == nil {
		u.SettingConfig = &UserSettingsConfig{}
	}
	u.SettingConfig.DefaultModel = modelID
}

// AccountEntity 账户实体（对应 Java 的 AccountEntity）
type AccountEntity struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	UserID    string    `gorm:"column:user_id" json:"userId"`
	Balance   float64   `gorm:"column:balance" json:"balance"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (AccountEntity) TableName() string {
	return "accounts"
}
