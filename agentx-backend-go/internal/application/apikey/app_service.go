package apikey

import (
	"time"

	domainAgent "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/agent"
	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/apikey"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// ApiKeyDTO API密钥DTO
type ApiKeyDTO struct {
	ID         string     `json:"id"`
	ApiKey     string     `json:"apiKey"`
	AgentID    string     `json:"agentId"`
	AgentName  string     `json:"agentName,omitempty"`
	UserID     string     `json:"userId"`
	Name       string     `json:"name"`
	Status     *bool      `json:"status"`
	UsageCount int        `json:"usageCount"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
	ExpiresAt  *time.Time `json:"expiresAt"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

// ApiKeyValidationResult API Key验证结果
type ApiKeyValidationResult struct {
	Valid   bool   `json:"valid"`
	UserID  string `json:"userId,omitempty"`
	AgentID string `json:"agentId,omitempty"`
	Message string `json:"message,omitempty"`
}

// CreateApiKeyRequest 创建API密钥请求
type CreateApiKeyRequest struct {
	AgentID string `json:"agentId" binding:"required"`
	Name    string `json:"name" binding:"required"`
}

// QueryApiKeyRequest 查询API密钥请求
type QueryApiKeyRequest struct {
	Name    string `form:"name" json:"name"`
	Status  *bool  `form:"status" json:"status"`
	AgentID string `form:"agentId" json:"agentId"`
}

// UpdateApiKeyStatusRequest 更新API密钥状态请求
type UpdateApiKeyStatusRequest struct {
	Status *bool `json:"status" binding:"required"`
}

// EntityToDTO 实体转DTO
func EntityToDTO(e *domain.ApiKeyEntity) *ApiKeyDTO {
	if e == nil {
		return nil
	}
	return &ApiKeyDTO{
		ID:         e.ID,
		ApiKey:     e.ApiKey,
		AgentID:    e.AgentID,
		UserID:     e.UserID,
		Name:       e.Name,
		Status:     e.Status,
		UsageCount: e.UsageCount,
		LastUsedAt: e.LastUsedAt,
		ExpiresAt:  e.ExpiresAt,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

// EntitiesToDTOs 实体列表转DTO列表
func EntitiesToDTOs(entities []domain.ApiKeyEntity) []*ApiKeyDTO {
	dtos := make([]*ApiKeyDTO, len(entities))
	for i, e := range entities {
		dtos[i] = EntityToDTO(&e)
	}
	return dtos
}

// AppService API密钥应用服务
type AppService struct {
	apiKeyDomainService *domain.DomainService
	agentDomainService  *domainAgent.DomainService
}

func NewAppService(apiKeyDomainService *domain.DomainService, agentDomainService *domainAgent.DomainService) *AppService {
	return &AppService{
		apiKeyDomainService: apiKeyDomainService,
		agentDomainService:  agentDomainService,
	}
}

// CreateApiKey 创建API密钥
func (s *AppService) CreateApiKey(agentID, name, userID string) (*ApiKeyDTO, error) {
	agent, err := s.agentDomainService.GetAgent(agentID, userID)
	if err != nil || agent == nil {
		return nil, exception.NewBusinessException("Agent不存在或无权限访问")
	}

	e := &domain.ApiKeyEntity{
		AgentID: agentID,
		UserID:  userID,
		Name:    name,
	}
	status := true
	e.Status = &status

	created, err := s.apiKeyDomainService.CreateApiKey(e)
	if err != nil {
		return nil, err
	}

	dto := EntityToDTO(created)
	dto.AgentName = agent.Name
	return dto, nil
}

// GetUserApiKeys 获取用户的API密钥列表
func (s *AppService) GetUserApiKeys(userID string, req *QueryApiKeyRequest) ([]*ApiKeyDTO, error) {
	var name *string
	var status *bool
	var agentID *string
	if req != nil {
		if req.Name != "" {
			name = &req.Name
		}
		status = req.Status
		if req.AgentID != "" {
			agentID = &req.AgentID
		}
	}

	apiKeys, err := s.apiKeyDomainService.GetUserApiKeys(userID, name, status, agentID)
	if err != nil {
		return nil, err
	}
	dtos := EntitiesToDTOs(apiKeys)

	// 批量获取Agent信息
	agentIDs := make(map[string]bool)
	for _, k := range apiKeys {
		agentIDs[k.AgentID] = true
	}
	ids := make([]string, 0, len(agentIDs))
	for id := range agentIDs {
		ids = append(ids, id)
	}
	if len(ids) > 0 {
		agents, _ := s.agentDomainService.GetAgentsByIDs(ids)
		agentMap := make(map[string]string)
		for _, a := range agents {
			agentMap[a.ID] = a.Name
		}
		for _, dto := range dtos {
			if n, ok := agentMap[dto.AgentID]; ok {
				dto.AgentName = n
			}
		}
	}

	return dtos, nil
}

// GetAgentApiKeys 获取Agent的API密钥列表
func (s *AppService) GetAgentApiKeys(agentID, userID string) ([]*ApiKeyDTO, error) {
	agent, err := s.agentDomainService.GetAgent(agentID, userID)
	if err != nil || agent == nil {
		return nil, exception.NewBusinessException("Agent不存在或无权限访问")
	}

	apiKeys, err := s.apiKeyDomainService.GetAgentApiKeys(agentID, userID)
	if err != nil {
		return nil, err
	}
	dtos := EntitiesToDTOs(apiKeys)
	for _, dto := range dtos {
		dto.AgentName = agent.Name
	}
	return dtos, nil
}

// GetApiKey 获取API密钥详情
func (s *AppService) GetApiKey(apiKeyID, userID string) (*ApiKeyDTO, error) {
	e, err := s.apiKeyDomainService.GetApiKey(apiKeyID, userID)
	if err != nil {
		return nil, err
	}
	dto := EntityToDTO(e)
	agent, _ := s.agentDomainService.GetAgent(e.AgentID, userID)
	if agent != nil {
		dto.AgentName = agent.Name
	}
	return dto, nil
}

// UpdateApiKeyStatus 更新API密钥状态
func (s *AppService) UpdateApiKeyStatus(apiKeyID string, status bool, userID string) error {
	return s.apiKeyDomainService.UpdateStatus(apiKeyID, userID, status)
}

// DeleteApiKey 删除API密钥
func (s *AppService) DeleteApiKey(apiKeyID, userID string) error {
	return s.apiKeyDomainService.DeleteApiKey(apiKeyID, userID)
}

// ResetApiKey 重置API密钥
func (s *AppService) ResetApiKey(apiKeyID, userID string) (*ApiKeyDTO, error) {
	e, err := s.apiKeyDomainService.ResetApiKey(apiKeyID, userID)
	if err != nil {
		return nil, err
	}
	dto := EntityToDTO(e)
	agent, _ := s.agentDomainService.GetAgent(e.AgentID, userID)
	if agent != nil {
		dto.AgentName = agent.Name
	}
	return dto, nil
}

// ValidateExternalApiKey 验证外部API Key
func (s *AppService) ValidateExternalApiKey(apiKey string) *ApiKeyValidationResult {
	e, err := s.apiKeyDomainService.ValidateApiKey(apiKey)
	if err != nil {
		return &ApiKeyValidationResult{Valid: false, Message: err.Error()}
	}
	_ = s.apiKeyDomainService.UpdateUsage(apiKey)
	return &ApiKeyValidationResult{Valid: true, UserID: e.UserID, AgentID: e.AgentID}
}
