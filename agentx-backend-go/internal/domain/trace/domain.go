package trace

import (
	"time"

	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
	"gorm.io/gorm"
)

// ExecutionPhase 执行阶段枚举
type ExecutionPhase string

const (
	ExecutionPhaseInitialization          ExecutionPhase = "INITIALIZATION"
	ExecutionPhaseEnvironmentPreparation  ExecutionPhase = "ENVIRONMENT_PREPARATION"
	ExecutionPhaseBalanceCheck            ExecutionPhase = "BALANCE_CHECK"
	ExecutionPhaseMemoryInitialization    ExecutionPhase = "MEMORY_INITIALIZATION"
	ExecutionPhaseModelCall               ExecutionPhase = "MODEL_CALL"
	ExecutionPhaseToolExecution           ExecutionPhase = "TOOL_EXECUTION"
	ExecutionPhaseBilling                 ExecutionPhase = "BILLING"
	ExecutionPhaseResultProcessing        ExecutionPhase = "RESULT_PROCESSING"
)

// ExecutionStepType 执行步骤类型枚举
type ExecutionStepType string

const (
	ExecutionStepTypeUserMessage  ExecutionStepType = "USER_MESSAGE"
	ExecutionStepTypeAIResponse   ExecutionStepType = "AI_RESPONSE"
	ExecutionStepTypeToolCall     ExecutionStepType = "TOOL_CALL"
	ExecutionStepTypeErrorMessage ExecutionStepType = "ERROR_MESSAGE"
)

// AgentExecutionSummaryEntity Agent执行链路汇总实体
type AgentExecutionSummaryEntity struct {
	ID                     int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID                 string     `gorm:"column:user_id" json:"userId"`
	SessionID              string     `gorm:"column:session_id" json:"sessionId"`
	AgentID                string     `gorm:"column:agent_id" json:"agentId"`
	ExecutionStartTime     *time.Time `gorm:"column:execution_start_time" json:"executionStartTime"`
	ExecutionEndTime       *time.Time `gorm:"column:execution_end_time" json:"executionEndTime"`
	TotalExecutionTime     *int       `gorm:"column:total_execution_time" json:"totalExecutionTime"`
	TotalInputTokens       int        `gorm:"column:total_input_tokens" json:"totalInputTokens"`
	TotalOutputTokens      int        `gorm:"column:total_output_tokens" json:"totalOutputTokens"`
	TotalTokens            int        `gorm:"column:total_tokens" json:"totalTokens"`
	ToolCallCount          int        `gorm:"column:tool_call_count" json:"toolCallCount"`
	TotalToolExecutionTime int        `gorm:"column:total_tool_execution_time" json:"totalToolExecutionTime"`
	ExecutionSuccess       *bool      `gorm:"column:execution_success" json:"executionSuccess"`
	ErrorPhase             string     `gorm:"column:error_phase" json:"errorPhase"`
	ErrorMessage           string     `gorm:"column:error_message" json:"errorMessage"`

	entity.BaseEntity
}

func (AgentExecutionSummaryEntity) TableName() string {
	return "agent_execution_summary"
}

// Create 创建新的执行追踪汇总
func CreateSummary(userID, sessionID, agentID string) *AgentExecutionSummaryEntity {
	now := time.Now()
	f := false
	return &AgentExecutionSummaryEntity{
		UserID:             userID,
		SessionID:          sessionID,
		AgentID:            agentID,
		ExecutionStartTime: &now,
		ExecutionSuccess:   &f,
	}
}

// MarkCompleted 标记执行完成
func (e *AgentExecutionSummaryEntity) MarkCompleted(success bool, errorPhase, errorMessage string) {
	now := time.Now()
	e.ExecutionEndTime = &now
	e.ExecutionSuccess = &success
	e.ErrorPhase = errorPhase
	e.ErrorMessage = errorMessage
	if e.ExecutionStartTime != nil {
		duration := int(now.Sub(*e.ExecutionStartTime).Milliseconds())
		e.TotalExecutionTime = &duration
	}
}

// AddTokens 添加Token统计
func (e *AgentExecutionSummaryEntity) AddTokens(inputTokens, outputTokens *int) {
	if inputTokens != nil {
		e.TotalInputTokens += *inputTokens
	}
	if outputTokens != nil {
		e.TotalOutputTokens += *outputTokens
	}
	e.TotalTokens = e.TotalInputTokens + e.TotalOutputTokens
}

// AddToolExecution 添加工具调用信息
func (e *AgentExecutionSummaryEntity) AddToolExecution(executionTime *int) {
	e.ToolCallCount++
	if executionTime != nil {
		e.TotalToolExecutionTime += *executionTime
	}
}

// AgentExecutionDetailEntity Agent执行链路详细记录实体
type AgentExecutionDetailEntity struct {
	ID                   int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SessionID            string `gorm:"column:session_id" json:"sessionId"`
	MessageContent       string `gorm:"column:message_content" json:"messageContent"`
	MessageType          string `gorm:"column:message_type" json:"messageType"`
	ModelEndpoint        string `gorm:"column:model_endpoint" json:"modelEndpoint"`
	ProviderName         string `gorm:"column:provider_name" json:"providerName"`
	MessageTokens        *int   `gorm:"column:message_tokens" json:"messageTokens"`
	ModelCallTime        *int   `gorm:"column:model_call_time" json:"modelCallTime"`
	ToolName             string `gorm:"column:tool_name" json:"toolName"`
	ToolRequestArgs      string `gorm:"column:tool_request_args" json:"toolRequestArgs"`
	ToolResponseData     string `gorm:"column:tool_response_data" json:"toolResponseData"`
	ToolExecutionTime    *int   `gorm:"column:tool_execution_time" json:"toolExecutionTime"`
	ToolSuccess          *bool  `gorm:"column:tool_success" json:"toolSuccess"`
	IsFallbackUsed       *bool  `gorm:"column:is_fallback_used" json:"isFallbackUsed"`
	FallbackReason       string `gorm:"column:fallback_reason" json:"fallbackReason"`
	FallbackFromEndpoint string `gorm:"column:fallback_from_endpoint" json:"fallbackFromEndpoint"`
	FallbackToEndpoint   string `gorm:"column:fallback_to_endpoint" json:"fallbackToEndpoint"`
	FallbackFromProvider string `gorm:"column:fallback_from_provider" json:"fallbackFromProvider"`
	FallbackToProvider   string `gorm:"column:fallback_to_provider" json:"fallbackToProvider"`
	StepSuccess          *bool  `gorm:"column:step_success" json:"stepSuccess"`
	StepErrorMessage     string `gorm:"column:step_error_message" json:"stepErrorMessage"`

	entity.BaseEntity
}

func (AgentExecutionDetailEntity) TableName() string {
	return "agent_execution_details"
}

// CreateUserMessageStep 创建用户消息步骤
func CreateUserMessageStep(sessionID, userMessage string, messageTokens *int) *AgentExecutionDetailEntity {
	t := true
	f := false
	return &AgentExecutionDetailEntity{
		SessionID:      sessionID,
		MessageContent: userMessage,
		MessageType:    string(ExecutionStepTypeUserMessage),
		MessageTokens:  messageTokens,
		IsFallbackUsed: &f,
		StepSuccess:    &t,
	}
}

// CreateAIResponseStep 创建AI响应步骤
func CreateAIResponseStep(sessionID, aiResponse, modelEndpoint, providerName string, messageTokens, modelCallTime *int) *AgentExecutionDetailEntity {
	t := true
	f := false
	return &AgentExecutionDetailEntity{
		SessionID:      sessionID,
		MessageContent: aiResponse,
		MessageType:    string(ExecutionStepTypeAIResponse),
		ModelEndpoint:  modelEndpoint,
		ProviderName:   providerName,
		MessageTokens:  messageTokens,
		ModelCallTime:  modelCallTime,
		IsFallbackUsed: &f,
		StepSuccess:    &t,
	}
}

// CreateToolCallStep 创建工具调用步骤
func CreateToolCallStep(sessionID, toolName, requestArgs, responseData string, executionTime *int, success bool) *AgentExecutionDetailEntity {
	f := false
	return &AgentExecutionDetailEntity{
		SessionID:         sessionID,
		MessageContent:    "执行工具：" + toolName,
		MessageType:       string(ExecutionStepTypeToolCall),
		ToolName:          toolName,
		ToolRequestArgs:   requestArgs,
		ToolResponseData:  responseData,
		ToolExecutionTime: executionTime,
		ToolSuccess:       &success,
		IsFallbackUsed:    &f,
		StepSuccess:       &success,
	}
}

// CreateErrorMessageStep 创建异常消息步骤
func CreateErrorMessageStep(sessionID, errorMessage string) *AgentExecutionDetailEntity {
	f := false
	msg := errorMessage
	if msg == "" {
		msg = "未知错误"
	}
	return &AgentExecutionDetailEntity{
		SessionID:        sessionID,
		MessageContent:   msg,
		MessageType:      string(ExecutionStepTypeErrorMessage),
		IsFallbackUsed:   &f,
		StepSuccess:      &f,
		StepErrorMessage: errorMessage,
	}
}

// SetFallbackInfo 设置模型降级信息
func (e *AgentExecutionDetailEntity) SetFallbackInfo(reason, fromEndpoint, toEndpoint, fromProvider, toProvider string) {
	t := true
	e.IsFallbackUsed = &t
	e.FallbackReason = reason
	e.FallbackFromEndpoint = fromEndpoint
	e.FallbackToEndpoint = toEndpoint
	e.FallbackFromProvider = fromProvider
	e.FallbackToProvider = toProvider
}

// MarkStepFailed 标记步骤失败
func (e *AgentExecutionDetailEntity) MarkStepFailed(errorMessage string) {
	f := false
	e.StepSuccess = &f
	e.StepErrorMessage = errorMessage
}

// ExecutionStatistics 执行统计信息
type ExecutionStatistics struct {
	TotalExecutions      int     `json:"totalExecutions"`
	SuccessfulExecutions int     `json:"successfulExecutions"`
	FailedExecutions     int     `json:"failedExecutions"`
	SuccessRate          float64 `json:"successRate"`
	TotalTokens          int64   `json:"totalTokens"`
}

// ---- 仓储接口 ----

type ExecutionSummaryRepository interface {
	FindByID(id int64) (*AgentExecutionSummaryEntity, error)
	FindByIDAndUserID(id int64, userID string) (*AgentExecutionSummaryEntity, error)
	FindBySessionIDAndUserID(sessionID, userID string) ([]AgentExecutionSummaryEntity, error)
	FindByUserIDPaged(userID string, page, pageSize int) ([]AgentExecutionSummaryEntity, int64, error)
	FindByUserIDAndTimeRange(userID string, startTime, endTime time.Time) ([]AgentExecutionSummaryEntity, error)
	FindFailedByUserID(userID string) ([]AgentExecutionSummaryEntity, error)
	GetUserStatistics(userID string) (*ExecutionStatistics, error)
	Create(entity *AgentExecutionSummaryEntity) error
	Update(entity *AgentExecutionSummaryEntity) error
}

type ExecutionDetailRepository interface {
	FindBySessionID(sessionID string) ([]AgentExecutionDetailEntity, error)
	FindToolCallsBySessionID(sessionID string) ([]AgentExecutionDetailEntity, error)
	FindModelCallsBySessionID(sessionID string) ([]AgentExecutionDetailEntity, error)
	FindFallbackCallsBySessionID(sessionID string) ([]AgentExecutionDetailEntity, error)
	Create(entity *AgentExecutionDetailEntity) error
	BatchCreate(entities []AgentExecutionDetailEntity) error
}

// ---- GORM 实现 ----

type ExecutionSummaryRepositoryImpl struct {
	DB *gorm.DB
}

func NewExecutionSummaryRepository(db *gorm.DB) ExecutionSummaryRepository {
	return &ExecutionSummaryRepositoryImpl{DB: db}
}

func (r *ExecutionSummaryRepositoryImpl) FindByID(id int64) (*AgentExecutionSummaryEntity, error) {
	var e AgentExecutionSummaryEntity
	result := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *ExecutionSummaryRepositoryImpl) FindByIDAndUserID(id int64, userID string) (*AgentExecutionSummaryEntity, error) {
	var e AgentExecutionSummaryEntity
	result := r.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *ExecutionSummaryRepositoryImpl) FindBySessionIDAndUserID(sessionID, userID string) ([]AgentExecutionSummaryEntity, error) {
	var entities []AgentExecutionSummaryEntity
	result := r.DB.Where("session_id = ? AND user_id = ? AND deleted_at IS NULL", sessionID, userID).
		Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *ExecutionSummaryRepositoryImpl) FindByUserIDPaged(userID string, page, pageSize int) ([]AgentExecutionSummaryEntity, int64, error) {
	var entities []AgentExecutionSummaryEntity
	var total int64
	query := r.DB.Model(&AgentExecutionSummaryEntity{}).Where("user_id = ? AND deleted_at IS NULL", userID)
	query.Count(&total)
	result := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities)
	return entities, total, result.Error
}

func (r *ExecutionSummaryRepositoryImpl) FindByUserIDAndTimeRange(userID string, startTime, endTime time.Time) ([]AgentExecutionSummaryEntity, error) {
	var entities []AgentExecutionSummaryEntity
	result := r.DB.Where("user_id = ? AND execution_start_time >= ? AND execution_start_time <= ? AND deleted_at IS NULL", userID, startTime, endTime).
		Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *ExecutionSummaryRepositoryImpl) FindFailedByUserID(userID string) ([]AgentExecutionSummaryEntity, error) {
	var entities []AgentExecutionSummaryEntity
	result := r.DB.Where("user_id = ? AND execution_success = false AND deleted_at IS NULL", userID).
		Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *ExecutionSummaryRepositoryImpl) GetUserStatistics(userID string) (*ExecutionStatistics, error) {
	var total, success int64
	var totalTokens int64

	r.DB.Model(&AgentExecutionSummaryEntity{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&total)
	r.DB.Model(&AgentExecutionSummaryEntity{}).Where("user_id = ? AND execution_success = true AND deleted_at IS NULL", userID).Count(&success)

	row := r.DB.Model(&AgentExecutionSummaryEntity{}).Where("user_id = ? AND deleted_at IS NULL", userID).
		Select("COALESCE(SUM(total_tokens), 0)").Row()
	_ = row.Scan(&totalTokens)

	failed := total - success
	var rate float64
	if total > 0 {
		rate = float64(success) / float64(total) * 100
	}

	return &ExecutionStatistics{
		TotalExecutions:      int(total),
		SuccessfulExecutions: int(success),
		FailedExecutions:     int(failed),
		SuccessRate:          rate,
		TotalTokens:          totalTokens,
	}, nil
}

func (r *ExecutionSummaryRepositoryImpl) Create(e *AgentExecutionSummaryEntity) error {
	return r.DB.Create(e).Error
}

func (r *ExecutionSummaryRepositoryImpl) Update(e *AgentExecutionSummaryEntity) error {
	return r.DB.Save(e).Error
}

// ---- Detail GORM 实现 ----

type ExecutionDetailRepositoryImpl struct {
	DB *gorm.DB
}

func NewExecutionDetailRepository(db *gorm.DB) ExecutionDetailRepository {
	return &ExecutionDetailRepositoryImpl{DB: db}
}

func (r *ExecutionDetailRepositoryImpl) FindBySessionID(sessionID string) ([]AgentExecutionDetailEntity, error) {
	var entities []AgentExecutionDetailEntity
	result := r.DB.Where("session_id = ? AND deleted_at IS NULL", sessionID).
		Order("created_at ASC").Find(&entities)
	return entities, result.Error
}

func (r *ExecutionDetailRepositoryImpl) FindToolCallsBySessionID(sessionID string) ([]AgentExecutionDetailEntity, error) {
	var entities []AgentExecutionDetailEntity
	result := r.DB.Where("session_id = ? AND message_type = ? AND deleted_at IS NULL", sessionID, string(ExecutionStepTypeToolCall)).
		Order("created_at ASC").Find(&entities)
	return entities, result.Error
}

func (r *ExecutionDetailRepositoryImpl) FindModelCallsBySessionID(sessionID string) ([]AgentExecutionDetailEntity, error) {
	var entities []AgentExecutionDetailEntity
	result := r.DB.Where("session_id = ? AND message_type = ? AND deleted_at IS NULL", sessionID, string(ExecutionStepTypeAIResponse)).
		Order("created_at ASC").Find(&entities)
	return entities, result.Error
}

func (r *ExecutionDetailRepositoryImpl) FindFallbackCallsBySessionID(sessionID string) ([]AgentExecutionDetailEntity, error) {
	var entities []AgentExecutionDetailEntity
	result := r.DB.Where("session_id = ? AND is_fallback_used = true AND deleted_at IS NULL", sessionID).
		Order("created_at ASC").Find(&entities)
	return entities, result.Error
}

func (r *ExecutionDetailRepositoryImpl) Create(e *AgentExecutionDetailEntity) error {
	return r.DB.Create(e).Error
}

func (r *ExecutionDetailRepositoryImpl) BatchCreate(entities []AgentExecutionDetailEntity) error {
	if len(entities) == 0 {
		return nil
	}
	return r.DB.Create(&entities).Error
}

// ---- 领域服务 ----

type DomainService struct {
	summaryRepo ExecutionSummaryRepository
	detailRepo  ExecutionDetailRepository
}

func NewDomainService(summaryRepo ExecutionSummaryRepository, detailRepo ExecutionDetailRepository) *DomainService {
	return &DomainService{summaryRepo: summaryRepo, detailRepo: detailRepo}
}

// GetExecutionSummary 获取执行汇总
func (s *DomainService) GetExecutionSummary(traceID string, userID string) (*AgentExecutionSummaryEntity, error) {
	// traceID 实际上是 session_id
	entities, err := s.summaryRepo.FindBySessionIDAndUserID(traceID, userID)
	if err != nil {
		return nil, err
	}
	if len(entities) == 0 {
		return nil, nil
	}
	return &entities[0], nil
}

// GetExecutionDetails 获取执行详情
func (s *DomainService) GetExecutionDetails(traceID, userID string) ([]AgentExecutionDetailEntity, error) {
	return s.detailRepo.FindBySessionID(traceID)
}

// GetSessionExecutionHistory 获取会话执行历史
func (s *DomainService) GetSessionExecutionHistory(sessionID, userID string) ([]AgentExecutionSummaryEntity, error) {
	return s.summaryRepo.FindBySessionIDAndUserID(sessionID, userID)
}

// GetUserExecutionHistory 分页获取用户执行历史
func (s *DomainService) GetUserExecutionHistory(userID string, page, pageSize int) ([]AgentExecutionSummaryEntity, int64, error) {
	return s.summaryRepo.FindByUserIDPaged(userID, page, pageSize)
}

// GetUserExecutionsByTimeRange 按时间范围查询
func (s *DomainService) GetUserExecutionsByTimeRange(userID string, startTime, endTime time.Time) ([]AgentExecutionSummaryEntity, error) {
	return s.summaryRepo.FindByUserIDAndTimeRange(userID, startTime, endTime)
}

// GetUserFailedExecutions 获取用户失败的执行记录
func (s *DomainService) GetUserFailedExecutions(userID string) ([]AgentExecutionSummaryEntity, error) {
	return s.summaryRepo.FindFailedByUserID(userID)
}

// GetUserExecutionStatistics 获取用户执行统计
func (s *DomainService) GetUserExecutionStatistics(userID string) (*ExecutionStatistics, error) {
	return s.summaryRepo.GetUserStatistics(userID)
}

// GetToolCallsBySessionID 获取工具调用记录
func (s *DomainService) GetToolCallsBySessionID(sessionID string) ([]AgentExecutionDetailEntity, error) {
	return s.detailRepo.FindToolCallsBySessionID(sessionID)
}

// GetModelCallsBySessionID 获取模型调用记录
func (s *DomainService) GetModelCallsBySessionID(sessionID string) ([]AgentExecutionDetailEntity, error) {
	return s.detailRepo.FindModelCallsBySessionID(sessionID)
}

// GetFallbackCallsBySessionID 获取降级调用记录
func (s *DomainService) GetFallbackCallsBySessionID(sessionID string) ([]AgentExecutionDetailEntity, error) {
	return s.detailRepo.FindFallbackCallsBySessionID(sessionID)
}

// CreateSummaryRecord 创建汇总记录
func (s *DomainService) CreateSummaryRecord(e *AgentExecutionSummaryEntity) error {
	return s.summaryRepo.Create(e)
}

// UpdateSummaryRecord 更新汇总记录
func (s *DomainService) UpdateSummaryRecord(e *AgentExecutionSummaryEntity) error {
	return s.summaryRepo.Update(e)
}

// CreateDetailRecord 创建详情记录
func (s *DomainService) CreateDetailRecord(e *AgentExecutionDetailEntity) error {
	return s.detailRepo.Create(e)
}

// BatchCreateDetailRecords 批量创建详情记录
func (s *DomainService) BatchCreateDetailRecords(entities []AgentExecutionDetailEntity) error {
	return s.detailRepo.BatchCreate(entities)
}
