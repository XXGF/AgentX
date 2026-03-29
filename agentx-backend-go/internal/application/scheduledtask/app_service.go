package scheduledtask

import (
	"time"

	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/scheduledtask"
)

// ScheduledTaskDTO 定时任务DTO
type ScheduledTaskDTO struct {
	ID              string                    `json:"id"`
	UserID          string                    `json:"userId"`
	AgentID         string                    `json:"agentId"`
	SessionID       string                    `json:"sessionId"`
	Content         string                    `json:"content"`
	RepeatType      domain.RepeatType         `json:"repeatType"`
	RepeatConfig    *domain.RepeatConfig      `json:"repeatConfig"`
	Status          domain.ScheduleTaskStatus `json:"status"`
	LastExecuteTime *time.Time                `json:"lastExecuteTime"`
	NextExecuteTime *time.Time                `json:"nextExecuteTime"`
	CreatedAt       time.Time                 `json:"createdAt"`
	UpdatedAt       time.Time                 `json:"updatedAt"`
}

// CreateScheduledTaskRequest 创建定时任务请求
type CreateScheduledTaskRequest struct {
	AgentID      string               `json:"agentId" binding:"required"`
	SessionID    string               `json:"sessionId" binding:"required"`
	Content      string               `json:"content" binding:"required"`
	RepeatType   domain.RepeatType    `json:"repeatType" binding:"required"`
	RepeatConfig *domain.RepeatConfig `json:"repeatConfig" binding:"required"`
}

// UpdateScheduledTaskRequest 更新定时任务请求
type UpdateScheduledTaskRequest struct {
	ID           string               `json:"id"`
	Content      string               `json:"content"`
	RepeatType   domain.RepeatType    `json:"repeatType"`
	RepeatConfig *domain.RepeatConfig `json:"repeatConfig"`
}

// EntityToDTO 实体转DTO
func EntityToDTO(e *domain.ScheduledTaskEntity) *ScheduledTaskDTO {
	if e == nil {
		return nil
	}
	return &ScheduledTaskDTO{
		ID:              e.ID,
		UserID:          e.UserID,
		AgentID:         e.AgentID,
		SessionID:       e.SessionID,
		Content:         e.Content,
		RepeatType:      e.RepeatType,
		RepeatConfig:    e.RepeatConfig,
		Status:          e.Status,
		LastExecuteTime: e.LastExecuteTime,
		NextExecuteTime: e.NextExecuteTime,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

// EntitiesToDTOs 实体列表转DTO列表
func EntitiesToDTOs(entities []domain.ScheduledTaskEntity) []*ScheduledTaskDTO {
	dtos := make([]*ScheduledTaskDTO, len(entities))
	for i, e := range entities {
		dtos[i] = EntityToDTO(&e)
	}
	return dtos
}

// AppService 定时任务应用服务
type AppService struct {
	domainService *domain.DomainService
}

func NewAppService(domainService *domain.DomainService) *AppService {
	return &AppService{domainService: domainService}
}

// CreateScheduledTask 创建定时任务
func (s *AppService) CreateScheduledTask(req *CreateScheduledTaskRequest, userID string) (*ScheduledTaskDTO, error) {
	entity := &domain.ScheduledTaskEntity{
		UserID:       userID,
		AgentID:      req.AgentID,
		SessionID:    req.SessionID,
		Content:      req.Content,
		RepeatType:   req.RepeatType,
		RepeatConfig: req.RepeatConfig,
		Status:       domain.ScheduleTaskStatusActive,
	}

	// TODO: 计算下次执行时间（需要 TaskScheduleService）

	saved, err := s.domainService.CreateTask(entity)
	if err != nil {
		return nil, err
	}

	// TODO: 调度任务执行（需要 ScheduledTaskExecutionService）

	return EntityToDTO(saved), nil
}

// UpdateScheduledTask 更新定时任务
func (s *AppService) UpdateScheduledTask(req *UpdateScheduledTaskRequest, userID string) (*ScheduledTaskDTO, error) {
	task, err := s.domainService.GetTask(req.ID, userID)
	if err != nil {
		return nil, err
	}

	if req.Content != "" {
		task.Content = req.Content
	}
	if req.RepeatType != "" {
		task.RepeatType = req.RepeatType
	}
	if req.RepeatConfig != nil {
		task.RepeatConfig = req.RepeatConfig
	}

	if err := s.domainService.UpdateTask(task); err != nil {
		return nil, err
	}

	return EntityToDTO(task), nil
}

// DeleteTask 删除定时任务
func (s *AppService) DeleteTask(taskID, userID string) error {
	return s.domainService.DeleteTask(taskID, userID)
}

// GetTask 获取单个定时任务
func (s *AppService) GetTask(taskID, userID string) (*ScheduledTaskDTO, error) {
	entity, err := s.domainService.GetTask(taskID, userID)
	if err != nil {
		return nil, err
	}
	return EntityToDTO(entity), nil
}

// GetUserTasks 获取用户的定时任务列表
func (s *AppService) GetUserTasks(userID string) ([]*ScheduledTaskDTO, error) {
	entities, err := s.domainService.GetTasksByUserID(userID)
	if err != nil {
		return nil, err
	}
	return EntitiesToDTOs(entities), nil
}

// GetTasksBySessionID 根据会话ID获取定时任务列表
func (s *AppService) GetTasksBySessionID(sessionID, userID string) ([]*ScheduledTaskDTO, error) {
	entities, err := s.domainService.GetTasksBySessionID(sessionID)
	if err != nil {
		return nil, err
	}
	// 过滤出属于当前用户的任务
	var filtered []domain.ScheduledTaskEntity
	for _, e := range entities {
		if e.UserID == userID {
			filtered = append(filtered, e)
		}
	}
	return EntitiesToDTOs(filtered), nil
}

// GetTasksByAgentID 根据Agent ID获取定时任务列表
func (s *AppService) GetTasksByAgentID(agentID, userID string) ([]*ScheduledTaskDTO, error) {
	entities, err := s.domainService.GetTasksByAgentID(agentID)
	if err != nil {
		return nil, err
	}
	var filtered []domain.ScheduledTaskEntity
	for _, e := range entities {
		if e.UserID == userID {
			filtered = append(filtered, e)
		}
	}
	return EntitiesToDTOs(filtered), nil
}

// PauseTask 暂停定时任务
func (s *AppService) PauseTask(taskID, userID string) (*ScheduledTaskDTO, error) {
	if err := s.domainService.PauseTask(taskID, userID); err != nil {
		return nil, err
	}
	entity, err := s.domainService.GetTask(taskID, userID)
	if err != nil {
		return nil, err
	}
	return EntityToDTO(entity), nil
}

// ResumeTask 恢复定时任务
func (s *AppService) ResumeTask(taskID, userID string) (*ScheduledTaskDTO, error) {
	if err := s.domainService.ResumeTask(taskID, userID); err != nil {
		return nil, err
	}
	entity, err := s.domainService.GetTask(taskID, userID)
	if err != nil {
		return nil, err
	}
	return EntityToDTO(entity), nil
}
