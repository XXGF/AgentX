package conversation

import (
	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/conversation"
)

// SessionToDTO Session实体转DTO
func SessionToDTO(entity *domain.SessionEntity) *SessionDTO {
	if entity == nil {
		return nil
	}
	return &SessionDTO{
		ID:          entity.ID,
		Title:       entity.Title,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		IsArchived:  entity.IsArchived,
		AgentID:     entity.AgentID,
	}
}

// SessionsToDTOs Session实体列表转DTO列表
func SessionsToDTOs(entities []domain.SessionEntity) []*SessionDTO {
	if len(entities) == 0 {
		return []*SessionDTO{}
	}
	dtos := make([]*SessionDTO, len(entities))
	for i, e := range entities {
		dtos[i] = SessionToDTO(&e)
	}
	return dtos
}

// MessageToDTO Message实体转DTO
func MessageToDTO(entity *domain.MessageEntity) *MessageDTO {
	if entity == nil {
		return nil
	}
	return &MessageDTO{
		ID:          entity.ID,
		Role:        entity.Role,
		Content:     entity.Content,
		CreatedAt:   entity.CreatedAt,
		Provider:    entity.Provider,
		Model:       entity.Model,
		MessageType: entity.MessageType,
		FileUrls:    entity.FileUrls,
	}
}

// MessagesToDTOs Message实体列表转DTO列表
func MessagesToDTOs(entities []domain.MessageEntity) []*MessageDTO {
	if len(entities) == 0 {
		return []*MessageDTO{}
	}
	dtos := make([]*MessageDTO, len(entities))
	for i, e := range entities {
		dtos[i] = MessageToDTO(&e)
	}
	return dtos
}
