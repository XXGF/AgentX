package task

import (
	"time"

	"github.com/google/uuid"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
	"gorm.io/gorm"
)

// TaskStatus 任务状态枚举
type TaskStatus string

const (
	TaskStatusWaiting    TaskStatus = "WAITING"     // 等待中
	TaskStatusInProgress TaskStatus = "IN_PROGRESS" // 进行中
	TaskStatusCompleted  TaskStatus = "COMPLETED"   // 已完成
	TaskStatusFailed     TaskStatus = "FAILED"      // 失败
)

// TaskEntity 任务实体
type TaskEntity struct {
	ID           string     `gorm:"column:id;primaryKey" json:"id"`
	SessionID    string     `gorm:"column:session_id" json:"sessionId"`
	UserID       string     `gorm:"column:user_id" json:"userId"`
	ParentTaskID string     `gorm:"column:parent_task_id" json:"parentTaskId"`
	TaskName     string     `gorm:"column:task_name" json:"taskName"`
	Description  string     `gorm:"column:description" json:"description"`
	Status       TaskStatus `gorm:"column:status" json:"status"`
	Progress     int        `gorm:"column:progress" json:"progress"`
	StartTime    *time.Time `gorm:"column:start_time" json:"startTime"`
	EndTime      *time.Time `gorm:"column:end_time" json:"endTime"`
	TaskResult   string     `gorm:"column:task_result" json:"taskResult"`

	entity.BaseEntity
}

func (TaskEntity) TableName() string {
	return "agent_tasks"
}

// UpdateStatus 更新任务状态
func (e *TaskEntity) UpdateStatus(status TaskStatus) {
	e.Status = status
	now := time.Now()
	if status == TaskStatusInProgress && e.StartTime == nil {
		e.StartTime = &now
	} else if (status == TaskStatusCompleted || status == TaskStatusFailed) && e.EndTime == nil {
		e.EndTime = &now
	}
	if status == TaskStatusCompleted {
		e.Progress = 100
	}
}

// UpdateProgress 更新进度
func (e *TaskEntity) UpdateProgress(progress int) {
	if progress < 0 {
		e.Progress = 0
	} else if progress > 100 {
		e.Progress = 100
	} else {
		e.Progress = progress
	}
}

// TaskAggregate 任务聚合（父任务+子任务）
type TaskAggregate struct {
	Task     *TaskEntity  `json:"task"`
	SubTasks []TaskEntity `json:"subTasks"`
}

// ---- 仓储接口 ----

type TaskRepository interface {
	FindByID(id string) (*TaskEntity, error)
	FindBySessionIDAndUserID(sessionID, userID string) ([]TaskEntity, error)
	FindLatestParentTask(sessionID, userID string) (*TaskEntity, error)
	FindSubTasks(parentTaskID string) ([]TaskEntity, error)
	Create(entity *TaskEntity) error
	Update(entity *TaskEntity) error
}

// ---- GORM 实现 ----

type TaskRepositoryImpl struct {
	DB *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &TaskRepositoryImpl{DB: db}
}

func (r *TaskRepositoryImpl) FindByID(id string) (*TaskEntity, error) {
	var e TaskEntity
	result := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *TaskRepositoryImpl) FindBySessionIDAndUserID(sessionID, userID string) ([]TaskEntity, error) {
	var entities []TaskEntity
	result := r.DB.Where("session_id = ? AND user_id = ? AND deleted_at IS NULL", sessionID, userID).
		Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *TaskRepositoryImpl) FindLatestParentTask(sessionID, userID string) (*TaskEntity, error) {
	var e TaskEntity
	result := r.DB.Where("session_id = ? AND user_id = ? AND parent_task_id = '0' AND deleted_at IS NULL", sessionID, userID).
		Order("created_at DESC").Limit(1).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *TaskRepositoryImpl) FindSubTasks(parentTaskID string) ([]TaskEntity, error) {
	var entities []TaskEntity
	result := r.DB.Where("parent_task_id = ? AND deleted_at IS NULL", parentTaskID).
		Order("created_at ASC").Find(&entities)
	return entities, result.Error
}

func (r *TaskRepositoryImpl) Create(e *TaskEntity) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return r.DB.Create(e).Error
}

func (r *TaskRepositoryImpl) Update(e *TaskEntity) error {
	return r.DB.Save(e).Error
}

// ---- 领域服务 ----

type DomainService struct {
	repo TaskRepository
}

func NewDomainService(repo TaskRepository) *DomainService {
	return &DomainService{repo: repo}
}

// AddTask 添加任务
func (s *DomainService) AddTask(task *TaskEntity) (*TaskEntity, error) {
	now := time.Now()
	task.StartTime = &now
	if err := s.repo.Create(task); err != nil {
		return nil, err
	}
	return task, nil
}

// UpdateTask 更新任务
func (s *DomainService) UpdateTask(task *TaskEntity) (*TaskEntity, error) {
	if err := s.repo.Update(task); err != nil {
		return nil, err
	}
	return task, nil
}

// GetCurrentSessionTask 获取当前会话的最新任务聚合
func (s *DomainService) GetCurrentSessionTask(sessionID, userID string) (*TaskAggregate, error) {
	parentTask, err := s.repo.FindLatestParentTask(sessionID, userID)
	if err != nil {
		return nil, err
	}
	if parentTask == nil {
		return nil, nil
	}
	subTasks, err := s.repo.FindSubTasks(parentTask.ID)
	if err != nil {
		return nil, err
	}
	return &TaskAggregate{Task: parentTask, SubTasks: subTasks}, nil
}

// GetSubTasks 获取子任务列表
func (s *DomainService) GetSubTasks(parentTaskID string) ([]TaskEntity, error) {
	return s.repo.FindSubTasks(parentTaskID)
}
