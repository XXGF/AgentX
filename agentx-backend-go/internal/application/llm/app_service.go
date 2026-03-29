package llm

import (
	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/llm"
	domainUser "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/user"
)

// AppService LLM应用服务（对应 Java 的 LLMAppService）
type AppService struct {
	llmDomainService  *domain.DomainService
	userDomainService *domainUser.DomainService
}

// NewAppService 创建LLM应用服务
func NewAppService(llmDomainService *domain.DomainService, userDomainService *domainUser.DomainService) *AppService {
	return &AppService{
		llmDomainService:  llmDomainService,
		userDomainService: userDomainService,
	}
}

// GetProviderDetail 获取服务商详情
func (s *AppService) GetProviderDetail(providerID, userID string) (*ProviderDTO, error) {
	agg, err := s.llmDomainService.GetProviderAggregate(providerID, userID)
	if err != nil {
		return nil, err
	}
	return ProviderAggregateToDTO(agg), nil
}

// CreateProvider 创建服务商
func (s *AppService) CreateProvider(req *ProviderCreateRequest, userID string) (*ProviderDTO, error) {
	entity := ProviderCreateRequestToEntity(req, userID)
	isOfficial := false
	entity.IsOfficial = &isOfficial
	provider, err := s.llmDomainService.CreateProvider(entity)
	if err != nil {
		return nil, err
	}
	return ProviderToDTO(provider), nil
}

// UpdateProvider 更新服务商
func (s *AppService) UpdateProvider(req *ProviderUpdateRequest, userID string) (*ProviderDTO, error) {
	// 先获取当前服务商数据
	existing, err := s.llmDomainService.GetProvider(req.ID, userID)
	if err != nil {
		return nil, err
	}

	// 如果传入的是掩码，使用原有的密钥
	if req.Config != nil && req.Config.ApiKey != "" && IsApiKeyMasked(req.Config.ApiKey) {
		req.Config.ApiKey = existing.Config.ApiKey
	}

	entity := ProviderUpdateRequestToEntity(req, userID)
	if err := s.llmDomainService.UpdateProvider(entity); err != nil {
		return nil, err
	}
	return ProviderToDTO(entity), nil
}

// GetProvider 获取服务商
func (s *AppService) GetProvider(providerID, userID string) (*ProviderDTO, error) {
	provider, err := s.llmDomainService.GetProvider(providerID, userID)
	if err != nil {
		return nil, err
	}
	return ProviderToDTO(provider), nil
}

// DeleteProvider 删除服务商
func (s *AppService) DeleteProvider(providerID, userID string) error {
	return s.llmDomainService.DeleteProvider(providerID, userID, true)
}

// GetUserProviders 获取用户自己的服务商
func (s *AppService) GetUserProviders(userID string) ([]*ProviderDTO, error) {
	aggregates, err := s.llmDomainService.GetUserProviders(userID)
	if err != nil {
		return nil, err
	}
	dtos := make([]*ProviderDTO, len(aggregates))
	for i, agg := range aggregates {
		dtos[i] = ProviderAggregateToDTO(agg)
	}
	return dtos, nil
}

// GetAllProviders 获取所有服务商
func (s *AppService) GetAllProviders(userID string) ([]*ProviderDTO, error) {
	aggregates, err := s.llmDomainService.GetAllProviders(userID)
	if err != nil {
		return nil, err
	}
	dtos := make([]*ProviderDTO, len(aggregates))
	for i, agg := range aggregates {
		dtos[i] = ProviderAggregateToDTO(agg)
	}
	return dtos, nil
}

// GetOfficialProviders 获取官方服务商
func (s *AppService) GetOfficialProviders() ([]*ProviderDTO, error) {
	aggregates, err := s.llmDomainService.GetOfficialProviders()
	if err != nil {
		return nil, err
	}
	dtos := make([]*ProviderDTO, len(aggregates))
	for i, agg := range aggregates {
		dtos[i] = ProviderAggregateToDTO(agg)
	}
	return dtos, nil
}

// GetCustomProviders 获取用户自定义服务商
func (s *AppService) GetCustomProviders(userID string) ([]*ProviderDTO, error) {
	aggregates, err := s.llmDomainService.GetCustomProviders(userID)
	if err != nil {
		return nil, err
	}
	dtos := make([]*ProviderDTO, len(aggregates))
	for i, agg := range aggregates {
		dtos[i] = ProviderAggregateToDTO(agg)
	}
	return dtos, nil
}

// GetUserProviderProtocols 获取用户服务商协议
func (s *AppService) GetUserProviderProtocols() []domain.ProviderProtocol {
	return s.llmDomainService.GetProviderProtocols()
}

// GetProvidersByType 根据类型获取服务商
func (s *AppService) GetProvidersByType(providerType domain.ProviderType, userID string) ([]*ProviderDTO, error) {
	aggregates, err := s.llmDomainService.GetProvidersByType(providerType, userID)
	if err != nil {
		return nil, err
	}
	dtos := make([]*ProviderDTO, len(aggregates))
	for i, agg := range aggregates {
		dtos[i] = ProviderAggregateToDTO(agg)
	}
	return dtos, nil
}

// UpdateProviderStatus 修改服务商状态
func (s *AppService) UpdateProviderStatus(providerID, userID string) error {
	return s.llmDomainService.UpdateProviderStatus(providerID, userID)
}

// CreateModel 创建模型
func (s *AppService) CreateModel(req *ModelCreateRequest, userID string) (*ModelDTO, error) {
	entity := ModelCreateRequestToEntity(req, userID)
	isOfficial := false
	entity.IsOfficial = &isOfficial
	if err := s.llmDomainService.CheckProviderExists(req.ProviderID, userID); err != nil {
		return nil, err
	}
	if err := s.llmDomainService.CreateModel(entity); err != nil {
		return nil, err
	}

	// 如果用户没有默认模型则设置当前模型
	defaultModelID := s.userDomainService.GetUserDefaultModelID(userID)
	if defaultModelID == "" {
		_ = s.userDomainService.SetUserDefaultModelID(userID, entity.ID)
	}

	return ModelToDTO(entity), nil
}

// UpdateModel 修改模型
func (s *AppService) UpdateModel(req *ModelUpdateRequest, userID string) (*ModelDTO, error) {
	entity := ModelUpdateRequestToEntity(req, userID)
	if err := s.llmDomainService.UpdateModel(entity); err != nil {
		return nil, err
	}
	return ModelToDTO(entity), nil
}

// DeleteModel 删除模型
func (s *AppService) DeleteModel(modelID, userID string) error {
	return s.llmDomainService.DeleteModel(modelID, userID, false)
}

// UpdateModelStatus 修改模型状态
func (s *AppService) UpdateModelStatus(modelID, userID string) error {
	return s.llmDomainService.UpdateModelStatus(modelID, userID)
}

// GetActiveModelsByType 获取所有激活模型
func (s *AppService) GetActiveModelsByType(providerType domain.ProviderType, userID string, modelType *domain.ModelType) ([]*ModelDTO, error) {
	aggregates, err := s.llmDomainService.GetProvidersByType(providerType, userID)
	if err != nil {
		return nil, err
	}

	var dtos []*ModelDTO
	for _, agg := range aggregates {
		if !agg.GetStatus() {
			continue
		}
		for _, model := range agg.Models {
			if modelType != nil && model.Type != string(*modelType) {
				continue
			}
			dto := ModelToDTO(&model, agg.GetName())
			dtos = append(dtos, dto)
		}
	}
	return dtos, nil
}

// GetDefaultModel 获取用户默认模型
func (s *AppService) GetDefaultModel(userID string) (*ModelDTO, error) {
	defaultModelID := s.userDomainService.GetUserDefaultModelID(userID)
	if defaultModelID == "" {
		return nil, nil
	}
	model, err := s.llmDomainService.FindModelByID(defaultModelID)
	if err != nil {
		return nil, err
	}
	return ModelToDTO(model), nil
}

// CanUserUseModel 检查用户是否可以使用指定模型
func (s *AppService) CanUserUseModel(modelID, userID string) bool {
	model, err := s.llmDomainService.GetModelByID(modelID)
	if err != nil {
		return false
	}
	// 官方模型 + 激活状态 = 可用
	if model.GetIsOfficial() && model.GetStatus() {
		return true
	}
	// 用户自己的模型 + 激活状态 = 可用
	if userID == model.UserID && model.GetStatus() {
		return true
	}
	return false
}

// GetAvailableModelsForUser 获取用户可用的模型列表
func (s *AppService) GetAvailableModelsForUser(userID string) ([]*ModelDTO, error) {
	chatType := domain.ModelTypeChat
	return s.GetActiveModelsByType(domain.ProviderTypeAll, userID, &chatType)
}
