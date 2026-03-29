package llm

import (
	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/llm"
)

// AdminAppService 管理员LLM应用服务（对应 Java 的 AdminLLMAppService）
type AdminAppService struct {
	llmDomainService *domain.DomainService
}

// NewAdminAppService 创建管理员LLM应用服务
func NewAdminAppService(llmDomainService *domain.DomainService) *AdminAppService {
	return &AdminAppService{llmDomainService: llmDomainService}
}

// CreateProvider 创建官方服务商
func (s *AdminAppService) CreateProvider(req *ProviderCreateRequest, userID string) (*ProviderDTO, error) {
	entity := ProviderCreateRequestToEntity(req, userID)
	isOfficial := true
	entity.IsOfficial = &isOfficial
	provider, err := s.llmDomainService.CreateProvider(entity)
	if err != nil {
		return nil, err
	}
	return ProviderToDTO(provider), nil
}

// UpdateProvider 修改服务商
func (s *AdminAppService) UpdateProvider(req *ProviderUpdateRequest, userID string) (*ProviderDTO, error) {
	// 先获取当前服务商数据
	existing, err := s.llmDomainService.GetProviderByID(req.ID)
	if err != nil {
		return nil, err
	}

	// 如果传入的是掩码，使用原有的密钥
	if req.Config != nil && req.Config.ApiKey != "" && IsApiKeyMasked(req.Config.ApiKey) {
		req.Config.ApiKey = existing.Config.ApiKey
	}

	entity := ProviderUpdateRequestToEntity(req, userID)
	entity.SetAdmin()
	if err := s.llmDomainService.UpdateProvider(entity); err != nil {
		return nil, err
	}

	agg, err := s.llmDomainService.GetProviderAggregate(entity.ID, userID)
	if err != nil {
		return nil, err
	}
	return ProviderAggregateToDTO(agg), nil
}

// DeleteProvider 删除服务商
func (s *AdminAppService) DeleteProvider(providerID, userID string) error {
	return s.llmDomainService.DeleteProvider(providerID, userID, false)
}

// CreateModel 创建模型
func (s *AdminAppService) CreateModel(req *ModelCreateRequest, userID string) (*ModelDTO, error) {
	entity := ModelCreateRequestToEntity(req, userID)
	entity.SetAdmin()
	isOfficial := true
	entity.IsOfficial = &isOfficial
	if err := s.llmDomainService.CreateModel(entity); err != nil {
		return nil, err
	}
	return ModelToDTO(entity), nil
}

// UpdateModel 更新模型
func (s *AdminAppService) UpdateModel(req *ModelUpdateRequest, userID string) (*ModelDTO, error) {
	entity := ModelUpdateRequestToEntity(req, userID)
	entity.SetAdmin()
	if err := s.llmDomainService.UpdateModel(entity); err != nil {
		return nil, err
	}
	return ModelToDTO(entity), nil
}

// DeleteModel 删除模型
func (s *AdminAppService) DeleteModel(modelID, userID string) error {
	return s.llmDomainService.DeleteModel(modelID, userID, false)
}

// GetOfficialProviders 获取官方服务商列表（管理员需要看到所有模型）
func (s *AdminAppService) GetOfficialProviders(userID string, page, pageSize int) ([]*ProviderDTO, error) {
	aggregates, err := s.llmDomainService.GetOfficialProvidersWithAllModels()
	if err != nil {
		return nil, err
	}
	dtos := make([]*ProviderDTO, len(aggregates))
	for i, agg := range aggregates {
		dtos[i] = ProviderAggregateToDTO(agg)
	}
	return dtos, nil
}

// GetProviderDetail 获取服务商详情
func (s *AdminAppService) GetProviderDetail(providerID, userID string) (*ProviderDTO, error) {
	agg, err := s.llmDomainService.GetProviderAggregate(providerID, userID)
	if err != nil {
		return nil, err
	}
	return ProviderAggregateToDTO(agg), nil
}

// ToggleProviderStatus 切换服务商状态
func (s *AdminAppService) ToggleProviderStatus(providerID, userID string) error {
	return s.llmDomainService.UpdateProviderStatus(providerID, userID)
}

// GetProviderProtocols 获取支持的协议列表
func (s *AdminAppService) GetProviderProtocols() []domain.ProviderProtocol {
	return s.llmDomainService.GetProviderProtocols()
}

// GetOfficialModels 获取官方模型列表
func (s *AdminAppService) GetOfficialModels(userID string, providerID *string, modelType *domain.ModelType, page, pageSize int) ([]*ModelDTO, error) {
	aggregates, err := s.llmDomainService.GetOfficialProvidersWithAllModels()
	if err != nil {
		return nil, err
	}

	var dtos []*ModelDTO
	for _, agg := range aggregates {
		for _, model := range agg.Models {
			// 按类型过滤
			if modelType != nil && model.Type != string(*modelType) {
				continue
			}
			// 按服务商过滤
			if providerID != nil && agg.GetID() != *providerID {
				continue
			}
			dtos = append(dtos, ModelToDTO(&model))
		}
	}
	return dtos, nil
}

// ToggleModelStatus 切换模型状态
func (s *AdminAppService) ToggleModelStatus(modelID, userID string) error {
	return s.llmDomainService.UpdateModelStatus(modelID, userID)
}

// GetModelTypes 获取模型类型列表
func (s *AdminAppService) GetModelTypes() []domain.ModelType {
	return domain.AllModelTypes
}
