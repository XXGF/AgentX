package trace

import (
	"strconv"
	"time"

	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/trace"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// AgentExecutionSummaryDTO 执行汇总DTO
type AgentExecutionSummaryDTO struct {
	TraceID                string     `json:"traceId"`
	UserID                 string     `json:"userId"`
	SessionID              string     `json:"sessionId"`
	AgentID                string     `json:"agentId"`
	AgentName              string     `json:"agentName,omitempty"`
	ExecutionStartTime     *time.Time `json:"executionStartTime"`
	ExecutionEndTime       *time.Time `json:"executionEndTime"`
	TotalExecutionTime     *int       `json:"totalExecutionTime"`
	TotalInputTokens       int        `json:"totalInputTokens"`
	TotalOutputTokens      int        `json:"totalOutputTokens"`
	TotalTokens            int        `json:"totalTokens"`
	ToolCallCount          int        `json:"toolCallCount"`
	TotalToolExecutionTime int        `json:"totalToolExecutionTime"`
	ExecutionSuccess       *bool      `json:"executionSuccess"`
	ErrorPhase             string     `json:"errorPhase,omitempty"`
	ErrorMessage           string     `json:"errorMessage,omitempty"`
	CreatedTime            time.Time  `json:"createdTime"`
}

// AgentExecutionDetailDTO 执行详情DTO
type AgentExecutionDetailDTO struct {
	TraceID              string     `json:"traceId"`
	MessageContent       string     `json:"messageContent"`
	MessageType          string     `json:"messageType"`
	ModelEndpoint        string     `json:"modelEndpoint,omitempty"`
	ProviderName         string     `json:"providerName,omitempty"`
	MessageTokens        *int       `json:"messageTokens,omitempty"`
	ModelCallTime        *int       `json:"modelCallTime,omitempty"`
	ToolName             string     `json:"toolName,omitempty"`
	ToolRequestArgs      string     `json:"toolRequestArgs,omitempty"`
	ToolResponseData     string     `json:"toolResponseData,omitempty"`
	ToolExecutionTime    *int       `json:"toolExecutionTime,omitempty"`
	ToolSuccess          *bool      `json:"toolSuccess,omitempty"`
	IsFallbackUsed       *bool      `json:"isFallbackUsed"`
	FallbackReason       string     `json:"fallbackReason,omitempty"`
	FallbackFromEndpoint string     `json:"fallbackFromEndpoint,omitempty"`
	FallbackToEndpoint   string     `json:"fallbackToEndpoint,omitempty"`
	FallbackFromProvider string     `json:"fallbackFromProvider,omitempty"`
	FallbackToProvider   string     `json:"fallbackToProvider,omitempty"`
	StepSuccess          *bool      `json:"stepSuccess"`
	StepErrorMessage     string     `json:"stepErrorMessage,omitempty"`
	CreatedTime          time.Time  `json:"createdTime"`
}

// TraceDetailResponse 追踪详情响应
type TraceDetailResponse struct {
	Summary *AgentExecutionSummaryDTO  `json:"summary"`
	Details []*AgentExecutionDetailDTO `json:"details"`
}

// ExecutionStatisticsResponse 执行统计响应
type ExecutionStatisticsResponse struct {
	TotalExecutions      int     `json:"totalExecutions"`
	SuccessfulExecutions int     `json:"successfulExecutions"`
	FailedExecutions     int     `json:"failedExecutions"`
	SuccessRate          float64 `json:"successRate"`
	TotalTokens          int64   `json:"totalTokens"`
}

// QueryExecutionHistoryRequest 查询执行历史请求
type QueryExecutionHistoryRequest struct {
	Page     *int `form:"page" json:"page"`
	PageSize *int `form:"pageSize" json:"pageSize"`
}

// ---- Assembler ----

func summaryEntityToDTO(e *domain.AgentExecutionSummaryEntity) *AgentExecutionSummaryDTO {
	if e == nil {
		return nil
	}
	return &AgentExecutionSummaryDTO{
		TraceID:                strconv.FormatInt(e.ID, 10),
		UserID:                 e.UserID,
		SessionID:              e.SessionID,
		AgentID:                e.AgentID,
		ExecutionStartTime:     e.ExecutionStartTime,
		ExecutionEndTime:       e.ExecutionEndTime,
		TotalExecutionTime:     e.TotalExecutionTime,
		TotalInputTokens:       e.TotalInputTokens,
		TotalOutputTokens:      e.TotalOutputTokens,
		TotalTokens:            e.TotalTokens,
		ToolCallCount:          e.ToolCallCount,
		TotalToolExecutionTime: e.TotalToolExecutionTime,
		ExecutionSuccess:       e.ExecutionSuccess,
		ErrorPhase:             e.ErrorPhase,
		ErrorMessage:           e.ErrorMessage,
		CreatedTime:            e.CreatedAt,
	}
}

func summaryEntitiesToDTOs(entities []domain.AgentExecutionSummaryEntity) []*AgentExecutionSummaryDTO {
	dtos := make([]*AgentExecutionSummaryDTO, len(entities))
	for i, e := range entities {
		dtos[i] = summaryEntityToDTO(&e)
	}
	return dtos
}

func detailEntityToDTO(e *domain.AgentExecutionDetailEntity) *AgentExecutionDetailDTO {
	if e == nil {
		return nil
	}
	return &AgentExecutionDetailDTO{
		TraceID:              strconv.FormatInt(e.ID, 10),
		MessageContent:       e.MessageContent,
		MessageType:          e.MessageType,
		ModelEndpoint:        e.ModelEndpoint,
		ProviderName:         e.ProviderName,
		MessageTokens:        e.MessageTokens,
		ModelCallTime:        e.ModelCallTime,
		ToolName:             e.ToolName,
		ToolRequestArgs:      e.ToolRequestArgs,
		ToolResponseData:     e.ToolResponseData,
		ToolExecutionTime:    e.ToolExecutionTime,
		ToolSuccess:          e.ToolSuccess,
		IsFallbackUsed:       e.IsFallbackUsed,
		FallbackReason:       e.FallbackReason,
		FallbackFromEndpoint: e.FallbackFromEndpoint,
		FallbackToEndpoint:   e.FallbackToEndpoint,
		FallbackFromProvider: e.FallbackFromProvider,
		FallbackToProvider:   e.FallbackToProvider,
		StepSuccess:          e.StepSuccess,
		StepErrorMessage:     e.StepErrorMessage,
		CreatedTime:          e.CreatedAt,
	}
}

func detailEntitiesToDTOs(entities []domain.AgentExecutionDetailEntity) []*AgentExecutionDetailDTO {
	dtos := make([]*AgentExecutionDetailDTO, len(entities))
	for i, e := range entities {
		dtos[i] = detailEntityToDTO(&e)
	}
	return dtos
}

// AppService 执行链路追踪应用服务
type AppService struct {
	traceDomainService *domain.DomainService
}

func NewAppService(traceDomainService *domain.DomainService) *AppService {
	return &AppService{traceDomainService: traceDomainService}
}

// GetExecutionHistory 分页查询用户执行历史
func (s *AppService) GetExecutionHistory(req *QueryExecutionHistoryRequest, userID string) ([]*AgentExecutionSummaryDTO, int64, error) {
	page := 1
	pageSize := 15
	if req.Page != nil {
		page = *req.Page
	}
	if req.PageSize != nil {
		pageSize = *req.PageSize
	}
	entities, total, err := s.traceDomainService.GetUserExecutionHistory(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return summaryEntitiesToDTOs(entities), total, nil
}

// GetTraceDetail 获取追踪详情
func (s *AppService) GetTraceDetail(traceID, userID string) (*TraceDetailResponse, error) {
	summary, err := s.traceDomainService.GetExecutionSummary(traceID, userID)
	if err != nil {
		return nil, err
	}
	if summary == nil {
		return nil, exception.NewBusinessException("追踪记录不存在")
	}
	details, err := s.traceDomainService.GetExecutionDetails(traceID, userID)
	if err != nil {
		return nil, err
	}
	return &TraceDetailResponse{
		Summary: summaryEntityToDTO(summary),
		Details: detailEntitiesToDTOs(details),
	}, nil
}

// GetExecutionDetails 获取执行详情列表
func (s *AppService) GetExecutionDetails(traceID, userID string) ([]*AgentExecutionDetailDTO, error) {
	entities, err := s.traceDomainService.GetExecutionDetails(traceID, userID)
	if err != nil {
		return nil, err
	}
	return detailEntitiesToDTOs(entities), nil
}

// GetSessionExecutionHistory 查询会话执行历史
func (s *AppService) GetSessionExecutionHistory(sessionID, userID string) ([]*AgentExecutionSummaryDTO, error) {
	entities, err := s.traceDomainService.GetSessionExecutionHistory(sessionID, userID)
	if err != nil {
		return nil, err
	}
	return summaryEntitiesToDTOs(entities), nil
}

// GetUserExecutionStatistics 获取用户执行统计
func (s *AppService) GetUserExecutionStatistics(userID string) (*ExecutionStatisticsResponse, error) {
	stats, err := s.traceDomainService.GetUserExecutionStatistics(userID)
	if err != nil {
		return nil, err
	}
	return &ExecutionStatisticsResponse{
		TotalExecutions:      stats.TotalExecutions,
		SuccessfulExecutions: stats.SuccessfulExecutions,
		FailedExecutions:     stats.FailedExecutions,
		SuccessRate:          stats.SuccessRate,
		TotalTokens:          stats.TotalTokens,
	}, nil
}
