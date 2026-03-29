package user

import (
	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/user"
)

// ToDTO 将 UserEntity 转换为 UserDTO（对应 Java 的 UserAssembler.toDTO）
func ToDTO(entity *domain.UserEntity) *UserDTO {
	if entity == nil {
		return nil
	}
	return &UserDTO{
		ID:            entity.ID,
		Nickname:      entity.Nickname,
		Email:         entity.Email,
		Phone:         entity.Phone,
		GithubID:      entity.GithubID,
		GithubLogin:   entity.GithubLogin,
		AvatarURL:     entity.AvatarURL,
		LoginPlatform: entity.LoginPlatform,
		IsAdmin:       entity.IsAdminUser(),
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
}

// ToDTOList 批量转换
func ToDTOList(entities []domain.UserEntity) []*UserDTO {
	dtos := make([]*UserDTO, 0, len(entities))
	for i := range entities {
		dtos = append(dtos, ToDTO(&entities[i]))
	}
	return dtos
}

// ToSettingsDTO 将 UserSettingsEntity 转换为 UserSettingsDTO
func ToSettingsDTO(entity *domain.UserSettingsEntity) *UserSettingsDTO {
	if entity == nil {
		return nil
	}
	dto := &UserSettingsDTO{
		ID:     entity.ID,
		UserID: entity.UserID,
	}
	if entity.SettingConfig != nil {
		dto.DefaultModel = entity.SettingConfig.DefaultModel
		dto.DefaultOcrModel = entity.SettingConfig.DefaultOcrModel
		dto.DefaultEmbeddingModel = entity.SettingConfig.DefaultEmbeddingModel
	}
	return dto
}

// UpdateRequestToEntity 将更新请求转换为实体（对应 Java 的 UserAssembler.toEntity(UserUpdateRequest, userId)）
func UpdateRequestToEntity(nickname string, userID string) *domain.UserEntity {
	return &domain.UserEntity{
		ID:       userID,
		Nickname: nickname,
	}
}
