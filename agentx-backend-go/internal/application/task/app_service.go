package task

import (
	"time"

	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/task"
)

// TaskDTO 任务DTO
type TaskDTO struct {
	ID           string     `json:"id"`
	SessionID    string     `json:"sessionId"`
	UserID       string     `json:"userId"`
	ParentTaskID string     `json:"parentTaskId"`
	TaskName     string     `json:"taskName"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	Progress     int        `json:"progress"`
	StartTime    *time.Time `json:"startTime"`
	EndTime      *time.Time `json:"endTime"`
	TaskResult   string     `json:"taskResult"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

// EntityToDTO 实体转DTO
func EntityToDTO(e *domain.TaskEntity) *TaskDTO {
	if e == nil {
		return nil
	}
	return &TaskDTO{
		ID:           e.ID,
		SessionID:    e.SessionID,
		UserID:       e.UserID,
		ParentTaskID: e.ParentTaskID,
		TaskName:     e.TaskName,
		Description:  e.Description,
		Status:       string(e.Status),
		Progress:     e.Progress,
		StartTime:    e.StartTime,
		EndTime:      e.EndTime,
		TaskResult:   e.TaskResult,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

// AppService 任务应用服务
type AppService struct {
	domainService *domain.DomainService
}

func NewAppService(domainService *domain.DomainService) *AppService {
	return &AppService{domainService: domainService}
}

// GetCurrentSessionTask 获取当前会话的最新任务
func (s *AppService) GetCurrentSessionTask(sessionID, userID string) (*domain.TaskAggregate, error) {
	return s.domainService.GetCurrentSessionTask(sessionID, userID)
}
