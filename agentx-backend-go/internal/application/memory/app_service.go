package memory

import (
	"time"

	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/memory"
)

// MemoryItemDTO 记忆条目DTO
type MemoryItemDTO struct {
	ID              string                 `json:"id"`
	UserID          string                 `json:"userId"`
	Type            string                 `json:"type"`
	Text            string                 `json:"text"`
	Data            map[string]interface{} `json:"data,omitempty"`
	Importance      *float32               `json:"importance"`
	Tags            []string               `json:"tags"`
	SourceSessionID string                 `json:"sourceSessionId,omitempty"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
}

// CreateMemoryRequest 创建记忆请求
type CreateMemoryRequest struct {
	Text       string                 `json:"text" binding:"required"`
	Type       string                 `json:"type"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Importance *float32               `json:"importance,omitempty"`
	Tags       []string               `json:"tags,omitempty"`
}

// QueryMemoryRequest 查询记忆请求
type QueryMemoryRequest struct {
	Type     string `form:"type" json:"type"`
	Page     *int   `form:"page" json:"page"`
	PageSize *int   `form:"pageSize" json:"pageSize"`
}

// EntityToDTO 实体转DTO
func EntityToDTO(e *domain.MemoryItemEntity) *MemoryItemDTO {
	if e == nil {
		return nil
	}
	return &MemoryItemDTO{
		ID:              e.ID,
		UserID:          e.UserID,
		Type:            e.Type,
		Text:            e.Text,
		Data:            e.Data,
		Importance:      e.Importance,
		Tags:            e.Tags,
		SourceSessionID: e.SourceSessionID,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

// EntitiesToDTOs 实体列表转DTO列表
func EntitiesToDTOs(entities []domain.MemoryItemEntity) []*MemoryItemDTO {
	dtos := make([]*MemoryItemDTO, len(entities))
	for i, e := range entities {
		dtos[i] = EntityToDTO(&e)
	}
	return dtos
}

// AppService 记忆应用服务
type AppService struct {
	domainService *domain.DomainService
}

func NewAppService(domainService *domain.DomainService) *AppService {
	return &AppService{domainService: domainService}
}

// ListUserMemories 分页列出用户记忆
func (s *AppService) ListUserMemories(userID string, req *QueryMemoryRequest) ([]*MemoryItemDTO, int64, error) {
	page := 1
	pageSize := 20
	if req.Page != nil {
		page = *req.Page
	}
	if req.PageSize != nil {
		pageSize = *req.PageSize
	}
	var memType *string
	if req.Type != "" {
		memType = &req.Type
	}
	entities, total, err := s.domainService.PageMemories(userID, memType, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return EntitiesToDTOs(entities), total, nil
}

// CreateMemory 手动创建记忆
func (s *AppService) CreateMemory(userID string, req *CreateMemoryRequest) ([]string, error) {
	memType := domain.SafeOf(req.Type)
	candidate := domain.CandidateMemory{
		Text:       req.Text,
		Type:       &memType,
		Data:       req.Data,
		Importance: req.Importance,
		Tags:       req.Tags,
	}
	return s.domainService.SaveMemories(userID, "", []domain.CandidateMemory{candidate})
}

// DeleteMemory 归档（软删除）记忆
func (s *AppService) DeleteMemory(userID, itemID string) error {
	return s.domainService.Delete(userID, itemID)
}
