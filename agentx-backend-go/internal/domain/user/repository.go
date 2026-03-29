package user

import (
	"gorm.io/gorm"
)

// UserRepository 用户仓储接口（对应 Java 的 UserRepository）
type UserRepository interface {
	FindByID(id string) (*UserEntity, error)
	FindByEmail(email string) (*UserEntity, error)
	FindByPhone(phone string) (*UserEntity, error)
	FindByAccount(account string) (*UserEntity, error)
	FindByGithubID(githubID string) (*UserEntity, error)
	ExistsByEmail(email string) (bool, error)
	Create(user *UserEntity) error
	Update(user *UserEntity) error
	FindByIDs(ids []string) ([]UserEntity, error)
	Page(keyword string, page, pageSize int) ([]UserEntity, int64, error)
}

// UserSettingsRepository 用户设置仓储接口（对应 Java 的 UserSettingsRepository）
type UserSettingsRepository interface {
	FindByUserID(userID string) (*UserSettingsEntity, error)
	Create(settings *UserSettingsEntity) error
	Update(settings *UserSettingsEntity) error
	UpdateByUserID(settings *UserSettingsEntity) error
}

// AccountRepository 账户仓储接口（对应 Java 的 AccountRepository）
type AccountRepository interface {
	FindByUserID(userID string) (*AccountEntity, error)
	Create(account *AccountEntity) error
	Update(account *AccountEntity) error
}

// ---- GORM 实现 ----

// UserRepositoryImpl 用户仓储 GORM 实现
type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

func (r *UserRepositoryImpl) FindByID(id string) (*UserEntity, error) {
	var user UserEntity
	result := r.DB.First(&user, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*UserEntity, error) {
	var user UserEntity
	result := r.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByPhone(phone string) (*UserEntity, error) {
	var user UserEntity
	result := r.DB.Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindByAccount 根据邮箱或手机号查找用户（对应 Java 的 findUserByAccount）
func (r *UserRepositoryImpl) FindByAccount(account string) (*UserEntity, error) {
	var user UserEntity
	result := r.DB.Where("email = ? OR phone = ?", account, account).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByGithubID(githubID string) (*UserEntity, error) {
	var user UserEntity
	result := r.DB.Where("github_id = ?", githubID).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepositoryImpl) ExistsByEmail(email string) (bool, error) {
	var count int64
	result := r.DB.Model(&UserEntity{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r *UserRepositoryImpl) Create(user *UserEntity) error {
	return r.DB.Create(user).Error
}

func (r *UserRepositoryImpl) Update(user *UserEntity) error {
	result := r.DB.Save(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *UserRepositoryImpl) FindByIDs(ids []string) ([]UserEntity, error) {
	var users []UserEntity
	if len(ids) == 0 {
		return users, nil
	}
	result := r.DB.Where("id IN ?", ids).Find(&users)
	return users, result.Error
}

// Page 分页查询（对应 Java 的 selectPage）
func (r *UserRepositoryImpl) Page(keyword string, page, pageSize int) ([]UserEntity, int64, error) {
	var users []UserEntity
	var total int64

	query := r.DB.Model(&UserEntity{})

	// 关键词搜索：昵称、邮箱、手机号
	if keyword != "" {
		query = query.Where("nickname LIKE ? OR email LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 先查总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按创建时间倒序
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ---- UserSettings GORM 实现 ----

type UserSettingsRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserSettingsRepository(db *gorm.DB) UserSettingsRepository {
	return &UserSettingsRepositoryImpl{DB: db}
}

func (r *UserSettingsRepositoryImpl) FindByUserID(userID string) (*UserSettingsEntity, error) {
	var settings UserSettingsEntity
	result := r.DB.Where("user_id = ?", userID).First(&settings)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &settings, nil
}

func (r *UserSettingsRepositoryImpl) Create(settings *UserSettingsEntity) error {
	return r.DB.Create(settings).Error
}

func (r *UserSettingsRepositoryImpl) Update(settings *UserSettingsEntity) error {
	result := r.DB.Save(settings)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *UserSettingsRepositoryImpl) UpdateByUserID(settings *UserSettingsEntity) error {
	result := r.DB.Where("user_id = ?", settings.UserID).Updates(settings)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
