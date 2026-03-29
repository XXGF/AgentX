package user

import (
	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/user"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// AppService 用户应用服务（对应 Java 的 UserAppService）
type AppService struct {
	domainService *domain.DomainService
}

// NewAppService 创建用户应用服务
func NewAppService(domainService *domain.DomainService) *AppService {
	return &AppService{domainService: domainService}
}

// GetUserInfo 获取用户信息（对应 Java 的 getUserInfo）
func (s *AppService) GetUserInfo(id string) (*UserDTO, error) {
	user, err := s.domainService.GetUserInfo(id)
	if err != nil {
		return nil, err
	}
	return ToDTO(user), nil
}

// UpdateUserInfo 修改用户信息（对应 Java 的 updateUserInfo）
func (s *AppService) UpdateUserInfo(nickname string, userID string) error {
	user := UpdateRequestToEntity(nickname, userID)
	return s.domainService.UpdateUserInfo(user)
}

// ChangePassword 修改用户密码（对应 Java 的 changePassword）
func (s *AppService) ChangePassword(currentPassword, newPassword, confirmPassword, userID string) error {
	// 验证确认密码
	if newPassword != confirmPassword {
		return exception.NewBusinessException("新密码和确认密码不一致")
	}

	return s.domainService.ChangePassword(userID, currentPassword, newPassword)
}

// GetUsers 分页获取用户列表（对应 Java 的 getUsers）
func (s *AppService) GetUsers(keyword string, page, pageSize int) (*common.PageResponse[*UserDTO], error) {
	users, total, err := s.domainService.GetUsers(keyword, page, pageSize)
	if err != nil {
		return nil, err
	}

	dtos := ToDTOList(users)
	return common.NewPageResponse(dtos, total, page, pageSize), nil
}
