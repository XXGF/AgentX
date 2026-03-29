package user

import (
	"github.com/google/uuid"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/utils"
)

// DomainService 用户领域服务（对应 Java 的 UserDomainService）
type DomainService struct {
	userRepo     UserRepository
	settingsRepo UserSettingsRepository
}

// NewDomainService 创建用户领域服务
func NewDomainService(userRepo UserRepository, settingsRepo UserSettingsRepository) *DomainService {
	return &DomainService{
		userRepo:     userRepo,
		settingsRepo: settingsRepo,
	}
}

// GetUserInfo 获取用户信息（对应 Java 的 getUserInfo）
func (s *DomainService) GetUserInfo(id string) (*UserEntity, error) {
	return s.userRepo.FindByID(id)
}

// FindUserByAccount 根据邮箱或手机号查找用户（对应 Java 的 findUserByAccount）
func (s *DomainService) FindUserByAccount(account string) (*UserEntity, error) {
	return s.userRepo.FindByAccount(account)
}

// FindUserByGithubID 根据GitHub ID查找用户（对应 Java 的 findUserByGithubId）
func (s *DomainService) FindUserByGithubID(githubID string) (*UserEntity, error) {
	return s.userRepo.FindByGithubID(githubID)
}

// Register 注册用户（对应 Java 的 register）
func (s *DomainService) Register(email, phone, password string) (*UserEntity, error) {
	// 加密密码
	encodedPassword, err := utils.EncodePassword(password)
	if err != nil {
		return nil, exception.NewBusinessExceptionWithCause("密码加密失败", err)
	}

	user := &UserEntity{
		ID:            uuid.New().String(),
		Email:         email,
		Phone:         phone,
		Password:      encodedPassword,
		Nickname:      generateNickname(),
		LoginPlatform: "normal",
	}

	// 校验实体
	if err := user.Valid(); err != nil {
		return nil, err
	}

	// 检查账号是否已存在
	if err := s.CheckAccountExist(email); err != nil {
		return nil, err
	}

	// 创建用户
	if err := s.userRepo.Create(user); err != nil {
		return nil, exception.NewBusinessExceptionWithCause("创建用户失败", err)
	}

	// 创建用户设置
	settings := &UserSettingsEntity{
		ID:     uuid.New().String(),
		UserID: user.ID,
	}
	if err := s.settingsRepo.Create(settings); err != nil {
		return nil, exception.NewBusinessExceptionWithCause("创建用户设置失败", err)
	}

	return user, nil
}

// Login 登录（对应 Java 的 login）
func (s *DomainService) Login(account, password string) (*UserEntity, error) {
	user, err := s.userRepo.FindByAccount(account)
	if err != nil {
		return nil, exception.NewBusinessExceptionWithCause("查询用户失败", err)
	}

	if user == nil || !utils.MatchesPassword(password, user.Password) {
		return nil, exception.NewBusinessException("账号密码错误")
	}

	return user, nil
}

// CheckAccountExist 检查账号是否已存在（对应 Java 的 checkAccountExist）
func (s *DomainService) CheckAccountExist(email string) error {
	if email == "" {
		return nil
	}
	exists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return exception.NewBusinessExceptionWithCause("检查账号失败", err)
	}
	if exists {
		return exception.NewBusinessException("账号已存在,不可重复注册")
	}
	return nil
}

// UpdateUserInfo 更新用户信息（对应 Java 的 updateUserInfo）
func (s *DomainService) UpdateUserInfo(user *UserEntity) error {
	return s.userRepo.Update(user)
}

// UpdatePassword 更新用户密码（对应 Java 的 updatePassword）
func (s *DomainService) UpdatePassword(userID, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return exception.NewBusinessExceptionWithCause("查询用户失败", err)
	}
	if user == nil {
		return exception.NewBusinessException("用户不存在")
	}

	encodedPassword, err := utils.EncodePassword(newPassword)
	if err != nil {
		return exception.NewBusinessExceptionWithCause("密码加密失败", err)
	}

	user.Password = encodedPassword
	return s.userRepo.Update(user)
}

// ChangePassword 修改密码（需要验证当前密码）（对应 Java 的 changePassword）
func (s *DomainService) ChangePassword(userID, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return exception.NewBusinessExceptionWithCause("查询用户失败", err)
	}
	if user == nil {
		return exception.NewBusinessException("用户不存在")
	}

	// 验证当前密码
	if !utils.MatchesPassword(currentPassword, user.Password) {
		return exception.NewBusinessException("当前密码不正确")
	}

	// 检查新密码是否与当前密码相同
	if utils.MatchesPassword(newPassword, user.Password) {
		return exception.NewBusinessException("新密码不能与当前密码相同")
	}

	// 加密新密码并更新
	encodedPassword, err := utils.EncodePassword(newPassword)
	if err != nil {
		return exception.NewBusinessExceptionWithCause("密码加密失败", err)
	}

	user.Password = encodedPassword
	return s.userRepo.Update(user)
}

// GetByIDs 批量获取用户（对应 Java 的 getByIds）
func (s *DomainService) GetByIDs(userIDs []string) ([]UserEntity, error) {
	if len(userIDs) == 0 {
		return []UserEntity{}, nil
	}
	return s.userRepo.FindByIDs(userIDs)
}

// CreateDefaultUser 创建默认用户（对应 Java 的 createDefaultUser）
func (s *DomainService) CreateDefaultUser(user *UserEntity) error {
	if err := user.Valid(); err != nil {
		return err
	}

	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	settings := &UserSettingsEntity{
		ID:     uuid.New().String(),
		UserID: user.ID,
	}
	return s.settingsRepo.Create(settings)
}

// GetUsers 分页查询用户列表（对应 Java 的 getUsers）
func (s *DomainService) GetUsers(keyword string, page, pageSize int) ([]UserEntity, int64, error) {
	return s.userRepo.Page(keyword, page, pageSize)
}

// EncryptPassword 加密密码（对应 Java 的 encryptPassword）
func (s *DomainService) EncryptPassword(password string) (string, error) {
	return utils.EncodePassword(password)
}

// generateNickname 随机生成用户昵称（对应 Java 的 generateNickname）
func generateNickname() string {
	return "agent-x" + uuid.New().String()[:6]
}

// GetUserDefaultModelID 获取用户默认模型ID（对应 Java 的 getUserDefaultModelId）
func (s *DomainService) GetUserDefaultModelID(userID string) string {
	settings, err := s.settingsRepo.FindByUserID(userID)
	if err != nil || settings == nil {
		return ""
	}
	return settings.GetDefaultModelID()
}

// SetUserDefaultModelID 设置用户默认模型ID（对应 Java 的 setUserDefaultModelId）
func (s *DomainService) SetUserDefaultModelID(userID, modelID string) error {
	settings, err := s.settingsRepo.FindByUserID(userID)
	if err != nil {
		return err
	}
	if settings == nil {
		// 创建新的用户设置
		settings = &UserSettingsEntity{
			ID:     uuid.New().String(),
			UserID: userID,
		}
		settings.SetDefaultModelID(modelID)
		return s.settingsRepo.Create(settings)
	}
	// 更新现有设置
	settings.SetDefaultModelID(modelID)
	return s.settingsRepo.UpdateByUserID(settings)
}
