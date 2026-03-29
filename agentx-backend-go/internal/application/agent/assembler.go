package agent

import (
	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/agent"
	"github.com/google/uuid"
)

// === Agent Assembler ===

// AgentToDTO Agent实体转DTO
func AgentToDTO(entity *domain.AgentEntity) *AgentDTO {
	if entity == nil {
		return nil
	}
	return &AgentDTO{
		ID:               entity.ID,
		Name:             entity.Name,
		Avatar:           entity.Avatar,
		Description:      entity.Description,
		SystemPrompt:     entity.SystemPrompt,
		WelcomeMessage:   entity.WelcomeMessage,
		ToolIDs:          entity.ToolIDs,
		KnowledgeBaseIDs: entity.KnowledgeBaseIDs,
		PublishedVersion: entity.PublishedVersion,
		Enabled:          entity.Enabled,
		UserID:           entity.UserID,
		ToolPresetParams: entity.ToolPresetParams,
		MultiModal:       entity.MultiModal,
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
	}
}

// AgentsToDTOs Agent实体列表转DTO列表
func AgentsToDTOs(entities []domain.AgentEntity) []*AgentDTO {
	if len(entities) == 0 {
		return []*AgentDTO{}
	}
	dtos := make([]*AgentDTO, len(entities))
	for i, e := range entities {
		dtos[i] = AgentToDTO(&e)
	}
	return dtos
}

// CreateRequestToEntity 创建请求转实体
func CreateRequestToEntity(req *CreateAgentRequest, userID string) *domain.AgentEntity {
	enabled := true
	toolIDs := domain.StringList(req.ToolIDs)
	if toolIDs == nil {
		toolIDs = domain.StringList{}
	}
	kbIDs := domain.StringList(req.KnowledgeBaseIDs)
	if kbIDs == nil {
		kbIDs = domain.StringList{}
	}
	return &domain.AgentEntity{
		ID:               uuid.New().String(),
		Name:             req.Name,
		Description:      req.Description,
		Avatar:           req.Avatar,
		SystemPrompt:     req.SystemPrompt,
		WelcomeMessage:   req.WelcomeMessage,
		ToolIDs:          toolIDs,
		KnowledgeBaseIDs: kbIDs,
		UserID:           userID,
		Enabled:          &enabled,
		ToolPresetParams: req.ToolPresetParams,
		MultiModal:       req.MultiModal,
	}
}

// UpdateRequestToEntity 更新请求转实体
func UpdateRequestToEntity(req *UpdateAgentRequest, userID string) *domain.AgentEntity {
	return &domain.AgentEntity{
		ID:               req.ID,
		Name:             req.Name,
		Avatar:           req.Avatar,
		Description:      req.Description,
		Enabled:          req.Enabled,
		SystemPrompt:     req.SystemPrompt,
		WelcomeMessage:   req.WelcomeMessage,
		ToolIDs:          domain.StringList(req.ToolIDs),
		KnowledgeBaseIDs: domain.StringList(req.KnowledgeBaseIDs),
		UserID:           userID,
		ToolPresetParams: req.ToolPresetParams,
		MultiModal:       req.MultiModal,
	}
}

// === AgentVersion Assembler ===

// VersionToDTO 版本实体转DTO
func VersionToDTO(entity *domain.AgentVersionEntity) *AgentVersionDTO {
	if entity == nil {
		return nil
	}
	return &AgentVersionDTO{
		ID:               entity.ID,
		AgentID:          entity.AgentID,
		Name:             entity.Name,
		Avatar:           entity.Avatar,
		Description:      entity.Description,
		VersionNumber:    entity.VersionNumber,
		SystemPrompt:     entity.SystemPrompt,
		WelcomeMessage:   entity.WelcomeMessage,
		ToolIDs:          entity.ToolIDs,
		KnowledgeBaseIDs: entity.KnowledgeBaseIDs,
		ChangeLog:        entity.ChangeLog,
		PublishStatus:    entity.PublishStatus,
		RejectReason:     entity.RejectReason,
		ReviewTime:       entity.ReviewTime,
		PublishedAt:      entity.PublishedAt,
		UserID:           entity.UserID,
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
	}
}

// VersionsToDTOs 版本实体列表转DTO列表
func VersionsToDTOs(entities []domain.AgentVersionEntity) []*AgentVersionDTO {
	if len(entities) == 0 {
		return []*AgentVersionDTO{}
	}
	dtos := make([]*AgentVersionDTO, len(entities))
	for i, e := range entities {
		dtos[i] = VersionToDTO(&e)
	}
	return dtos
}

// CreateVersionEntity 创建版本实体
func CreateVersionEntity(agent *domain.AgentEntity, req *PublishAgentVersionRequest) *domain.AgentVersionEntity {
	return domain.CreateVersionFromAgent(agent, req.VersionNumber, req.ChangeLog)
}
