package llm

import (
	"regexp"

	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/llm"
	"github.com/google/uuid"
)

// === Provider Assembler ===

// ProviderToDTO 服务商实体转DTO
func ProviderToDTO(entity *domain.ProviderEntity) *ProviderDTO {
	if entity == nil {
		return nil
	}
	dto := &ProviderDTO{
		ID:          entity.ID,
		Protocol:    entity.Protocol,
		Name:        entity.Name,
		Description: entity.Description,
		Config:      entity.Config,
		IsOfficial:  entity.IsOfficial,
		Status:      entity.Status,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Models:      []*ModelDTO{},
	}
	// 脱敏处理
	dto.MaskSensitiveInfo()
	return dto
}

// ProviderAggregateToDTO 服务商聚合根转DTO
func ProviderAggregateToDTO(agg *domain.ProviderAggregate) *ProviderDTO {
	if agg == nil {
		return nil
	}
	dto := ProviderToDTO(agg.Entity)
	if dto == nil {
		return nil
	}
	for _, model := range agg.Models {
		modelDTO := ModelToDTO(&model, agg.GetName())
		modelDTO.IsOfficial = agg.Entity.IsOfficial
		dto.Models = append(dto.Models, modelDTO)
	}
	return dto
}

// ProviderCreateRequestToEntity 创建请求转实体
func ProviderCreateRequestToEntity(req *ProviderCreateRequest, userID string) *domain.ProviderEntity {
	status := true
	if req.Status != nil {
		status = *req.Status
	}
	return &domain.ProviderEntity{
		ID:          uuid.New().String(),
		UserID:      userID,
		Protocol:    req.Protocol,
		Name:        req.Name,
		Description: req.Description,
		Config:      req.Config,
		Status:      &status,
	}
}

// ProviderUpdateRequestToEntity 更新请求转实体
func ProviderUpdateRequestToEntity(req *ProviderUpdateRequest, userID string) *domain.ProviderEntity {
	return &domain.ProviderEntity{
		ID:          req.ID,
		UserID:      userID,
		Protocol:    req.Protocol,
		Name:        req.Name,
		Description: req.Description,
		Config:      req.Config,
		Status:      req.Status,
	}
}

// IsApiKeyMasked 检查API Key是否是掩码
func IsApiKeyMasked(apiKey string) bool {
	if apiKey == "" {
		return false
	}
	matched, _ := regexp.MatchString(`^\*+$`, apiKey)
	return matched
}

// === Model Assembler ===

// ModelToDTO 模型实体转DTO
func ModelToDTO(entity *domain.ModelEntity, providerName ...string) *ModelDTO {
	if entity == nil {
		return nil
	}
	dto := &ModelDTO{
		ID:            entity.ID,
		UserID:        entity.UserID,
		ProviderID:    entity.ProviderID,
		ModelID:       entity.ModelID,
		Name:          entity.Name,
		Description:   entity.Description,
		Type:          entity.Type,
		ModelEndpoint: entity.ModelEndpoint,
		IsOfficial:    entity.IsOfficial,
		Status:        entity.Status,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
	if len(providerName) > 0 {
		dto.ProviderName = providerName[0]
	}
	return dto
}

// ModelCreateRequestToEntity 创建请求转实体
func ModelCreateRequestToEntity(req *ModelCreateRequest, userID string) *domain.ModelEntity {
	endpoint := req.ModelEndpoint
	if endpoint == "" {
		endpoint = req.ModelID
	}
	return &domain.ModelEntity{
		ID:            uuid.New().String(),
		UserID:        userID,
		ProviderID:    req.ProviderID,
		ModelID:       req.ModelID,
		Name:          req.Name,
		Description:   req.Description,
		Type:          req.Type,
		ModelEndpoint: endpoint,
	}
}

// ModelUpdateRequestToEntity 更新请求转实体
func ModelUpdateRequestToEntity(req *ModelUpdateRequest, userID string) *domain.ModelEntity {
	endpoint := req.ModelEndpoint
	if endpoint == "" {
		endpoint = req.ModelID
	}
	return &domain.ModelEntity{
		ID:            req.ID,
		UserID:        userID,
		ModelID:       req.ModelID,
		Name:          req.Name,
		Description:   req.Description,
		ModelEndpoint: endpoint,
	}
}
