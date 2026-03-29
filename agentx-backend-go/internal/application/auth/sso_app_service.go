package auth

import (
	"strings"

	domainSso "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/sso"
	domainUser "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/user"
	jwtAuth "github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	infraSso "github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/sso"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/utils"
	"github.com/google/uuid"
)

// SsoAppService SSO应用服务（对应 Java 的 SsoAppService）
type SsoAppService struct {
	ssoFactory        *infraSso.ServiceFactory
	userDomainService *domainUser.DomainService
	jwtUtils          *jwtAuth.JWTUtils
}

// NewSsoAppService 创建SSO应用服务
func NewSsoAppService(ssoFactory *infraSso.ServiceFactory, userDomainService *domainUser.DomainService, jwtUtils *jwtAuth.JWTUtils) *SsoAppService {
	return &SsoAppService{
		ssoFactory:        ssoFactory,
		userDomainService: userDomainService,
		jwtUtils:          jwtUtils,
	}
}

// GetSsoLoginUrl 获取SSO登录URL
func (s *SsoAppService) GetSsoLoginUrl(provider, redirectUrl string) (string, error) {
	ssoService, err := s.ssoFactory.GetSsoService(provider)
	if err != nil {
		return "", err
	}
	return ssoService.GetLoginUrl(redirectUrl)
}

// HandleSsoCallback 处理SSO登录回调
func (s *SsoAppService) HandleSsoCallback(provider, authCode string) (string, error) {
	ssoService, err := s.ssoFactory.GetSsoService(provider)
	if err != nil {
		return "", err
	}

	ssoUserInfo, err := ssoService.GetUserInfo(authCode)
	if err != nil {
		return "", err
	}

	// 根据SSO用户信息创建或获取本地用户
	userEntity, err := s.findOrCreateUser(ssoUserInfo)
	if err != nil {
		return "", err
	}

	// 生成JWT token
	token, err := s.jwtUtils.GenerateToken(userEntity.ID)
	if err != nil {
		return "", exception.NewBusinessExceptionWithCause("生成Token失败", err)
	}

	return token, nil
}

// findOrCreateUser 根据SSO用户信息查找或创建本地用户
func (s *SsoAppService) findOrCreateUser(ssoUserInfo *domainSso.SsoUserInfo) (*domainUser.UserEntity, error) {
	var existingUser *domainUser.UserEntity
	var err error

	// GitHub用户优先通过GitHub ID查找
	if ssoUserInfo.Provider == domainSso.SsoProviderGitHub {
		existingUser, err = s.userDomainService.FindUserByGithubID(ssoUserInfo.ID)
		if err != nil {
			return nil, err
		}
	}

	// 如果通过特定ID没找到，再通过邮箱查找
	if existingUser == nil && ssoUserInfo.Email != "" {
		existingUser, err = s.userDomainService.FindUserByAccount(ssoUserInfo.Email)
		if err != nil {
			return nil, err
		}
	}

	if existingUser != nil {
		// 用户已存在，更新用户信息和登录平台
		s.updateUserFromSso(existingUser, ssoUserInfo)
		return existingUser, nil
	}

	// 用户不存在，创建新用户
	return s.createUserFromSso(ssoUserInfo)
}

// updateUserFromSso 从SSO信息更新用户
func (s *SsoAppService) updateUserFromSso(user *domainUser.UserEntity, ssoUserInfo *domainSso.SsoUserInfo) {
	needUpdate := false

	if ssoUserInfo.Name != "" && ssoUserInfo.Name != user.Nickname {
		user.Nickname = ssoUserInfo.Name
		needUpdate = true
	}
	if ssoUserInfo.Avatar != "" && ssoUserInfo.Avatar != user.AvatarURL {
		user.AvatarURL = ssoUserInfo.Avatar
		needUpdate = true
	}

	// GitHub用户需要更新GitHub相关信息
	if ssoUserInfo.Provider == domainSso.SsoProviderGitHub {
		if ssoUserInfo.ID != user.GithubID {
			user.GithubID = ssoUserInfo.ID
			needUpdate = true
		}
	}

	// 更新登录平台
	currentPlatform := string(ssoUserInfo.Provider)
	if currentPlatform != user.LoginPlatform {
		user.LoginPlatform = currentPlatform
		needUpdate = true
	}

	if needUpdate {
		_ = s.userDomainService.UpdateUserInfo(user)
	}
}

// createUserFromSso 从SSO信息创建新用户
func (s *SsoAppService) createUserFromSso(ssoUserInfo *domainSso.SsoUserInfo) (*domainUser.UserEntity, error) {
	nickname := ssoUserInfo.Name
	if nickname == "" {
		nickname = "sso-user-" + uuid.New().String()[:8]
	}

	user := &domainUser.UserEntity{
		ID:            uuid.New().String(),
		Email:         ssoUserInfo.Email,
		Nickname:      nickname,
		AvatarURL:     ssoUserInfo.Avatar,
		LoginPlatform: string(ssoUserInfo.Provider),
	}

	// 设置提供商特定的信息
	if ssoUserInfo.Provider == domainSso.SsoProviderGitHub {
		user.GithubID = ssoUserInfo.ID
		if ssoUserInfo.Desc != "" && strings.HasPrefix(ssoUserInfo.Desc, "GitHub用户: ") {
			user.GithubLogin = strings.TrimPrefix(ssoUserInfo.Desc, "GitHub用户: ")
		}
	}

	// SSO用户生成一个随机密码（用户无法知道这个密码，只能通过SSO登录）
	randomPassword := "SSO_" + uuid.New().String()
	encodedPassword, err := utils.EncodePassword(randomPassword)
	if err != nil {
		return nil, err
	}
	user.Password = encodedPassword

	if err := s.userDomainService.CreateDefaultUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
