package auth

import (
	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/auth"
)

// ToDTO 实体转DTO（对应 Java 的 AuthSettingAssembler.toDTO）
func ToDTO(entity *domain.AuthSettingEntity) *AuthSettingDTO {
	if entity == nil {
		return nil
	}
	return &AuthSettingDTO{
		ID:           entity.ID,
		FeatureType:  entity.FeatureType,
		FeatureKey:   entity.FeatureKey,
		FeatureName:  entity.FeatureName,
		Enabled:      entity.Enabled,
		ConfigData:   entity.ConfigData,
		DisplayOrder: entity.DisplayOrder,
		Description:  entity.Description,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}
}

// ToDTOList 实体列表转DTO列表
func ToDTOList(entities []domain.AuthSettingEntity) []*AuthSettingDTO {
	dtos := make([]*AuthSettingDTO, 0, len(entities))
	for i := range entities {
		dtos = append(dtos, ToDTO(&entities[i]))
	}
	return dtos
}

// UpdateEntityFromRequest 根据更新请求更新实体（对应 Java 的 AuthSettingAssembler.updateEntity）
func UpdateEntityFromRequest(entity *domain.AuthSettingEntity, req *UpdateAuthSettingRequest) *domain.AuthSettingEntity {
	if req.FeatureName != nil {
		entity.FeatureName = *req.FeatureName
	}
	if req.Enabled != nil {
		entity.SetEnabled(*req.Enabled)
	}
	if req.DisplayOrder != nil {
		entity.DisplayOrder = req.DisplayOrder
	}
	if req.Description != nil {
		entity.Description = *req.Description
	}
	if req.ConfigData != nil {
		entity.ConfigData = req.ConfigData
	}
	return entity
}
