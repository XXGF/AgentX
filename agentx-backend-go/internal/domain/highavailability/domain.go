package highavailability

import (
	"go.uber.org/zap"
)

// AffinityType 亲和性类型
type AffinityType string

const (
	AffinityTypeSession AffinityType = "SESSION" // 会话亲和性
)

// HighAvailabilityResult 高可用选择结果
type HighAvailabilityResult struct {
	ProviderID string `json:"providerId"` // 选择的Provider ID
	ModelID    string `json:"modelId"`    // 选择的Model ID
	InstanceID string `json:"instanceId"` // 实例ID（用于结果上报）
	Switched   bool   `json:"switched"`   // 模型是否被切换（降级到备用模型）
}

// ApiInstanceDTO API实例信息
type ApiInstanceDTO struct {
	ID         string `json:"id"`
	BusinessID string `json:"businessId"`
	Type       string `json:"type"`
}

// SelectInstanceRequest 选择实例请求
type SelectInstanceRequest struct {
	UserID        string       `json:"userId"`
	ModelID       string       `json:"modelId"`
	Type          string       `json:"type"`
	AffinityKey   string       `json:"affinityKey,omitempty"`
	AffinityType  AffinityType `json:"affinityType,omitempty"`
	FallbackChain []string     `json:"fallbackChain,omitempty"`
}

// ApiInstanceCreateRequest 创建API实例请求
type ApiInstanceCreateRequest struct {
	UserID     string `json:"userId"`
	ModelID    string `json:"modelId"`
	Type       string `json:"type"`
	BusinessID string `json:"businessId"`
}

// ApiInstanceUpdateRequest 更新API实例请求
type ApiInstanceUpdateRequest struct {
	UserID        string                 `json:"userId"`
	ModelID       string                 `json:"modelId"`
	RoutingParams map[string]interface{} `json:"routingParams,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// ReportResultRequest 结果上报请求
type ReportResultRequest struct {
	InstanceID    string `json:"instanceId"`
	BusinessID    string `json:"businessId"`
	Success       bool   `json:"success"`
	LatencyMs     int64  `json:"latencyMs"`
	ErrorMessage  string `json:"errorMessage,omitempty"`
	CallTimestamp int64  `json:"callTimestamp"`
}

// ProjectCreateRequest 项目创建请求
type ProjectCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ApiKey      string `json:"apiKey"`
}

// ModelDeleteItem 模型删除项
type ModelDeleteItem struct {
	Type       string `json:"type"`
	BusinessID string `json:"businessId"`
}

// ---- Gateway 接口（基础设施层实现）----

// Gateway 高可用网关接口
type Gateway interface {
	// SelectBestInstance 选择最佳实例
	SelectBestInstance(request *SelectInstanceRequest) (*ApiInstanceDTO, error)
	// CreateApiInstance 创建API实例
	CreateApiInstance(request *ApiInstanceCreateRequest) error
	// DeleteApiInstance 删除API实例
	DeleteApiInstance(instanceType, businessID string) error
	// UpdateApiInstance 更新API实例
	UpdateApiInstance(instanceType, businessID string, request *ApiInstanceUpdateRequest) error
	// ReportResult 上报调用结果
	ReportResult(request *ReportResultRequest) error
	// CreateProject 创建项目
	CreateProject(request *ProjectCreateRequest) error
	// BatchCreateApiInstances 批量创建API实例
	BatchCreateApiInstances(requests []*ApiInstanceCreateRequest) error
	// ActivateApiInstance 激活API实例
	ActivateApiInstance(instanceType, businessID string) error
	// DeactivateApiInstance 停用API实例
	DeactivateApiInstance(instanceType, businessID string) error
	// BatchDeleteApiInstances 批量删除API实例
	BatchDeleteApiInstances(items []*ModelDeleteItem) error
}

// ---- 高可用领域服务接口 ----

// DomainService 高可用领域服务接口
type DomainService interface {
	// SyncModelToGateway 同步模型到高可用网关
	SyncModelToGateway(modelID, userID, modelEndpoint string) error
	// RemoveModelFromGateway 从高可用网关删除模型
	RemoveModelFromGateway(modelID, userID string)
	// UpdateModelInGateway 更新高可用网关中的模型
	UpdateModelInGateway(modelID, userID, modelEndpoint string)
	// SelectBestProvider 通过高可用网关选择最佳Provider
	SelectBestProvider(modelID, userID, providerID string) *HighAvailabilityResult
	// SelectBestProviderWithSession 支持会话亲和性
	SelectBestProviderWithSession(modelID, userID, providerID, sessionID string) *HighAvailabilityResult
	// SelectBestProviderWithFallback 支持会话亲和性和降级链
	SelectBestProviderWithFallback(modelID, userID, providerID, sessionID string, fallbackChain []string) *HighAvailabilityResult
	// ReportCallResult 上报调用结果
	ReportCallResult(instanceID, modelID string, success bool, latencyMs int64, errorMessage string)
	// InitializeProject 初始化项目
	InitializeProject()
	// SyncAllModelsToGateway 批量同步所有模型
	SyncAllModelsToGateway()
	// ChangeModelStatusInGateway 变更模型状态
	ChangeModelStatusInGateway(modelID string, enabled bool, reason string)
	// BatchRemoveModelsFromGateway 批量删除模型
	BatchRemoveModelsFromGateway(deleteItems []*ModelDeleteItem, userID string)
}

// ---- 高可用配置 ----

// Config 高可用配置
type Config struct {
	Enabled bool   `mapstructure:"enabled" json:"enabled"`
	ApiKey  string `mapstructure:"api-key" json:"apiKey"`
	BaseURL string `mapstructure:"base-url" json:"baseUrl"`
}

// ---- 默认实现（高可用未启用时的降级实现）----

// DefaultDomainService 默认高可用领域服务（不启用高可用时使用）
type DefaultDomainService struct {
	config *Config
	logger *zap.Logger
}

func NewDefaultDomainService(config *Config, logger *zap.Logger) DomainService {
	return &DefaultDomainService{config: config, logger: logger}
}

func (s *DefaultDomainService) SyncModelToGateway(modelID, userID, modelEndpoint string) error {
	if !s.config.Enabled {
		s.logger.Debug("高可用功能未启用，跳过模型同步", zap.String("modelId", modelID))
		return nil
	}
	return nil
}

func (s *DefaultDomainService) RemoveModelFromGateway(modelID, userID string) {
	if !s.config.Enabled {
		s.logger.Debug("高可用功能未启用，跳过模型删除", zap.String("modelId", modelID))
	}
}

func (s *DefaultDomainService) UpdateModelInGateway(modelID, userID, modelEndpoint string) {
	if !s.config.Enabled {
		s.logger.Debug("高可用功能未启用，跳过模型更新", zap.String("modelId", modelID))
	}
}

func (s *DefaultDomainService) SelectBestProvider(modelID, userID, providerID string) *HighAvailabilityResult {
	// 高可用未启用，直接返回原始Provider
	return &HighAvailabilityResult{
		ProviderID: providerID,
		ModelID:    modelID,
		InstanceID: "",
		Switched:   false,
	}
}

func (s *DefaultDomainService) SelectBestProviderWithSession(modelID, userID, providerID, sessionID string) *HighAvailabilityResult {
	return s.SelectBestProvider(modelID, userID, providerID)
}

func (s *DefaultDomainService) SelectBestProviderWithFallback(modelID, userID, providerID, sessionID string, fallbackChain []string) *HighAvailabilityResult {
	return s.SelectBestProvider(modelID, userID, providerID)
}

func (s *DefaultDomainService) ReportCallResult(instanceID, modelID string, success bool, latencyMs int64, errorMessage string) {
	if !s.config.Enabled {
		return
	}
}

func (s *DefaultDomainService) InitializeProject() {
	if !s.config.Enabled {
		s.logger.Info("高可用功能未启用，跳过项目初始化")
	}
}

func (s *DefaultDomainService) SyncAllModelsToGateway() {
	if !s.config.Enabled {
		s.logger.Info("高可用功能未启用，跳过模型批量同步")
	}
}

func (s *DefaultDomainService) ChangeModelStatusInGateway(modelID string, enabled bool, reason string) {
	if !s.config.Enabled {
		s.logger.Debug("高可用功能未启用，跳过模型状态变更", zap.String("modelId", modelID))
	}
}

func (s *DefaultDomainService) BatchRemoveModelsFromGateway(deleteItems []*ModelDeleteItem, userID string) {
	if !s.config.Enabled {
		s.logger.Debug("高可用功能未启用，跳过批量模型删除",
			zap.String("userId", userID),
			zap.Int("count", len(deleteItems)))
	}
}
