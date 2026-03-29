package llm

import (
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// DomainService LLM领域服务（对应 Java 的 LLMDomainService）
type DomainService struct {
	providerRepo ProviderRepository
	modelRepo    ModelRepository
}

// NewDomainService 创建LLM领域服务
func NewDomainService(providerRepo ProviderRepository, modelRepo ModelRepository) *DomainService {
	return &DomainService{
		providerRepo: providerRepo,
		modelRepo:    modelRepo,
	}
}

// === 服务商相关 ===

// CreateProvider 创建服务商
func (s *DomainService) CreateProvider(provider *ProviderEntity) (*ProviderEntity, error) {
	if err := s.validateProviderProtocol(provider.Protocol); err != nil {
		return nil, err
	}
	if err := s.providerRepo.Create(provider); err != nil {
		return nil, err
	}
	return provider, nil
}

// UpdateProvider 更新服务商
func (s *DomainService) UpdateProvider(provider *ProviderEntity) error {
	if err := s.validateProviderProtocol(provider.Protocol); err != nil {
		return err
	}
	return s.providerRepo.Update(provider)
}

// GetUserProviders 获取用户自己的服务商
func (s *DomainService) GetUserProviders(userID string) ([]*ProviderAggregate, error) {
	providers, err := s.providerRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return s.buildProviderAggregatesWithActiveModels(providers)
}

// GetAllProviders 获取所有服务商（包含官方和用户自定义）
func (s *DomainService) GetAllProviders(userID string) ([]*ProviderAggregate, error) {
	providers, err := s.providerRepo.FindByUserIDOrOfficial(userID)
	if err != nil {
		return nil, err
	}
	return s.buildProviderAggregatesWithActiveModels(providers)
}

// GetOfficialProviders 获取官方服务商
func (s *DomainService) GetOfficialProviders() ([]*ProviderAggregate, error) {
	providers, err := s.providerRepo.FindOfficial()
	if err != nil {
		return nil, err
	}
	return s.buildProviderAggregatesWithActiveModels(providers)
}

// GetCustomProviders 获取用户自定义服务商
func (s *DomainService) GetCustomProviders(userID string) ([]*ProviderAggregate, error) {
	providers, err := s.providerRepo.FindCustomByUserID(userID)
	if err != nil {
		return nil, err
	}
	return s.buildProviderAggregatesWithActiveModels(providers)
}

// GetProvidersByType 根据类型获取服务商
func (s *DomainService) GetProvidersByType(providerType ProviderType, userID string) ([]*ProviderAggregate, error) {
	switch providerType {
	case ProviderTypeOfficial:
		return s.GetOfficialProviders()
	case ProviderTypeCustom:
		return s.GetCustomProviders(userID)
	default: // ALL
		return s.GetAllProviders(userID)
	}
}

// GetProvider 获取服务商（带用户校验）
func (s *DomainService) GetProvider(providerID, userID string) (*ProviderEntity, error) {
	provider, err := s.providerRepo.FindByIDAndUserID(providerID, userID)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, exception.NewBusinessException("服务商不存在")
	}
	return provider, nil
}

// GetProviderByID 获取服务商（不校验用户）
func (s *DomainService) GetProviderByID(providerID string) (*ProviderEntity, error) {
	provider, err := s.providerRepo.FindByID(providerID)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, exception.NewBusinessException("服务商不存在")
	}
	return provider, nil
}

// FindProviderByID 查找服务商（不抛异常）
func (s *DomainService) FindProviderByID(providerID string) (*ProviderEntity, error) {
	return s.providerRepo.FindByID(providerID)
}

// CheckProviderExists 检查服务商是否存在
func (s *DomainService) CheckProviderExists(providerID, userID string) error {
	provider, err := s.providerRepo.FindByIDAndUserID(providerID, userID)
	if err != nil {
		return err
	}
	if provider == nil {
		return exception.NewBusinessException("服务商不存在")
	}
	return nil
}

// GetProviderAggregate 获取服务商聚合根
func (s *DomainService) GetProviderAggregate(providerID, userID string) (*ProviderAggregate, error) {
	provider, err := s.GetProvider(providerID, userID)
	if err != nil {
		return nil, err
	}
	models, err := s.GetActiveModelList(providerID, userID)
	if err != nil {
		return nil, err
	}
	return NewProviderAggregate(provider, models), nil
}

// UpdateProviderStatus 修改服务商状态
func (s *DomainService) UpdateProviderStatus(providerID, userID string) error {
	return s.providerRepo.ToggleStatus(providerID, userID)
}

// DeleteProvider 删除服务商
func (s *DomainService) DeleteProvider(providerID, userID string, checkUser bool) error {
	if err := s.providerRepo.Delete(providerID, userID, checkUser); err != nil {
		return err
	}
	// 删除服务商下的所有模型
	_, err := s.modelRepo.DeleteByProviderID(providerID)
	return err
}

// GetProviderProtocols 获取所有支持的服务商协议
func (s *DomainService) GetProviderProtocols() []ProviderProtocol {
	return AllProviderProtocols
}

// === 模型相关 ===

// CreateModel 创建模型
func (s *DomainService) CreateModel(model *ModelEntity) error {
	return s.modelRepo.Create(model)
}

// UpdateModel 修改模型
func (s *DomainService) UpdateModel(model *ModelEntity) error {
	return s.modelRepo.Update(model)
}

// DeleteModel 删除模型
func (s *DomainService) DeleteModel(modelID, userID string, checkUser bool) error {
	return s.modelRepo.Delete(modelID, userID, checkUser)
}

// UpdateModelStatus 修改模型状态
func (s *DomainService) UpdateModelStatus(modelID, userID string) error {
	return s.modelRepo.ToggleStatus(modelID, userID)
}

// GetModelByID 获取模型
func (s *DomainService) GetModelByID(modelID string) (*ModelEntity, error) {
	model, err := s.modelRepo.FindByID(modelID)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, exception.NewBusinessException("模型不存在")
	}
	return model, nil
}

// FindModelByID 查找模型（不抛异常）
func (s *DomainService) FindModelByID(modelID string) (*ModelEntity, error) {
	return s.modelRepo.FindByID(modelID)
}

// GetModelList 获取模型列表
func (s *DomainService) GetModelList(providerID, userID string) ([]ModelEntity, error) {
	return s.modelRepo.FindByProviderIDAndUserID(providerID, userID)
}

// GetActiveModelList 获取激活的模型列表
func (s *DomainService) GetActiveModelList(providerID, userID string) ([]ModelEntity, error) {
	return s.modelRepo.FindActiveByProviderIDAndUserID(providerID, userID)
}

// GetAllActiveModels 获取所有激活的模型
func (s *DomainService) GetAllActiveModels() ([]ModelEntity, error) {
	return s.modelRepo.FindAllActive()
}

// GetModelsByIDs 批量获取模型信息
func (s *DomainService) GetModelsByIDs(modelIDs []string) ([]ModelEntity, error) {
	return s.modelRepo.FindByIDs(modelIDs)
}

// GetOfficialProvidersWithAllModels 获取官方服务商（包含所有模型）- 管理员功能
func (s *DomainService) GetOfficialProvidersWithAllModels() ([]*ProviderAggregate, error) {
	providers, err := s.providerRepo.FindOfficial()
	if err != nil {
		return nil, err
	}
	return s.buildProviderAggregatesWithAllModels(providers)
}

// === 内部方法 ===

// buildProviderAggregatesWithActiveModels 构建服务商聚合根（只包含激活模型）
func (s *DomainService) buildProviderAggregatesWithActiveModels(providers []ProviderEntity) ([]*ProviderAggregate, error) {
	if len(providers) == 0 {
		return []*ProviderAggregate{}, nil
	}

	providerIDs := make([]string, len(providers))
	for i, p := range providers {
		providerIDs[i] = p.ID
	}

	activeModels, err := s.modelRepo.FindActiveByProviderIDs(providerIDs)
	if err != nil {
		return nil, err
	}

	modelMap := make(map[string][]ModelEntity)
	for _, m := range activeModels {
		modelMap[m.ProviderID] = append(modelMap[m.ProviderID], m)
	}

	aggregates := make([]*ProviderAggregate, len(providers))
	for i, p := range providers {
		pCopy := p
		aggregates[i] = NewProviderAggregate(&pCopy, modelMap[p.ID])
	}
	return aggregates, nil
}

// buildProviderAggregatesWithAllModels 构建服务商聚合根（包含所有模型）- 管理员功能
func (s *DomainService) buildProviderAggregatesWithAllModels(providers []ProviderEntity) ([]*ProviderAggregate, error) {
	if len(providers) == 0 {
		return []*ProviderAggregate{}, nil
	}

	providerIDs := make([]string, len(providers))
	for i, p := range providers {
		providerIDs[i] = p.ID
	}

	allModels, err := s.modelRepo.FindByProviderIDs(providerIDs)
	if err != nil {
		return nil, err
	}

	modelMap := make(map[string][]ModelEntity)
	for _, m := range allModels {
		modelMap[m.ProviderID] = append(modelMap[m.ProviderID], m)
	}

	aggregates := make([]*ProviderAggregate, len(providers))
	for i, p := range providers {
		pCopy := p
		aggregates[i] = NewProviderAggregate(&pCopy, modelMap[p.ID])
	}
	return aggregates, nil
}

// validateProviderProtocol 验证服务商协议
func (s *DomainService) validateProviderProtocol(protocol string) error {
	_, ok := ProviderProtocolFromCode(protocol)
	if !ok {
		return exception.NewBusinessException("不支持的服务商协议类型: " + protocol)
	}
	return nil
}
