package scheduledtask

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	"gorm.io/gorm"
)

// RepeatType 重复类型枚举
type RepeatType string

const (
	RepeatTypeNone     RepeatType = "NONE"     // 不重复
	RepeatTypeDaily    RepeatType = "DAILY"    // 每天
	RepeatTypeWeekly   RepeatType = "WEEKLY"   // 每周
	RepeatTypeMonthly  RepeatType = "MONTHLY"  // 每月
	RepeatTypeWorkdays RepeatType = "WORKDAYS" // 工作日
	RepeatTypeCustom   RepeatType = "CUSTOM"   // 自定义
)

// ScheduleTaskStatus 任务状态枚举
type ScheduleTaskStatus string

const (
	ScheduleTaskStatusActive    ScheduleTaskStatus = "ACTIVE"    // 活跃
	ScheduleTaskStatusPaused    ScheduleTaskStatus = "PAUSED"    // 暂停
	ScheduleTaskStatusCompleted ScheduleTaskStatus = "COMPLETED" // 已完成
)

// RepeatConfig 重复配置值对象
type RepeatConfig struct {
	ExecuteDateTime *time.Time `json:"executeDateTime,omitempty"`
	Weekdays        []int      `json:"weekdays,omitempty"`
	MonthDay        *int       `json:"monthDay,omitempty"`
	Interval        *int       `json:"interval,omitempty"`
	TimeUnit        string     `json:"timeUnit,omitempty"`
	ExecuteTime     string     `json:"executeTime,omitempty"`
	EndDateTime     *time.Time `json:"endDateTime,omitempty"`
}

func (r RepeatConfig) Value() (driver.Value, error) {
	b, err := json.Marshal(r)
	return string(b), err
}

func (r *RepeatConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return nil
	}
	return json.Unmarshal(bytes, r)
}

// ScheduledTaskEntity 定时任务实体
type ScheduledTaskEntity struct {
	ID              string             `gorm:"column:id;primaryKey" json:"id"`
	UserID          string             `gorm:"column:user_id" json:"userId"`
	AgentID         string             `gorm:"column:agent_id" json:"agentId"`
	SessionID       string             `gorm:"column:session_id" json:"sessionId"`
	Content         string             `gorm:"column:content" json:"content"`
	RepeatType      RepeatType         `gorm:"column:repeat_type" json:"repeatType"`
	RepeatConfig    *RepeatConfig      `gorm:"column:repeat_config;type:jsonb" json:"repeatConfig"`
	Status          ScheduleTaskStatus `gorm:"column:status" json:"status"`
	LastExecuteTime *time.Time         `gorm:"column:last_execute_time" json:"lastExecuteTime"`
	NextExecuteTime *time.Time         `gorm:"column:next_execute_time" json:"nextExecuteTime"`

	entity.BaseEntity
}

func (ScheduledTaskEntity) TableName() string {
	return "scheduled_tasks"
}

// IsActive 检查任务是否活跃
func (e *ScheduledTaskEntity) IsActive() bool {
	return e.Status == ScheduleTaskStatusActive
}

// IsOneTime 检查是否为一次性任务
func (e *ScheduledTaskEntity) IsOneTime() bool {
	return e.RepeatType == RepeatTypeNone
}

// ---- 仓储接口 ----

type ScheduledTaskRepository interface {
	FindByID(id string) (*ScheduledTaskEntity, error)
	FindByIDAndUserID(id, userID string) (*ScheduledTaskEntity, error)
	FindByUserID(userID string) ([]ScheduledTaskEntity, error)
	FindByUserIDAndStatus(userID string, status ScheduleTaskStatus) ([]ScheduledTaskEntity, error)
	FindBySessionID(sessionID string) ([]ScheduledTaskEntity, error)
	FindByAgentID(agentID string) ([]ScheduledTaskEntity, error)
	FindActiveTasksToExecute() ([]ScheduledTaskEntity, error)
	Create(entity *ScheduledTaskEntity) error
	Update(entity *ScheduledTaskEntity) error
	UpdateStatus(id, userID string, status ScheduleTaskStatus) error
	DeleteByIDAndUserID(id, userID string) error
	DeleteBySessionIDAndUserID(sessionID, userID string) (int64, error)
	DeleteByAgentIDAndUserID(agentID, userID string) (int64, error)
	CountByUserID(userID string) (int64, error)
	CountByUserIDAndStatus(userID string, status ScheduleTaskStatus) (int64, error)
}

// ---- GORM 实现 ----

type ScheduledTaskRepositoryImpl struct {
	DB *gorm.DB
}

func NewScheduledTaskRepository(db *gorm.DB) ScheduledTaskRepository {
	return &ScheduledTaskRepositoryImpl{DB: db}
}

func (r *ScheduledTaskRepositoryImpl) FindByID(id string) (*ScheduledTaskEntity, error) {
	var e ScheduledTaskEntity
	result := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *ScheduledTaskRepositoryImpl) FindByIDAndUserID(id, userID string) (*ScheduledTaskEntity, error) {
	var e ScheduledTaskEntity
	result := r.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *ScheduledTaskRepositoryImpl) FindByUserID(userID string) ([]ScheduledTaskEntity, error) {
	var entities []ScheduledTaskEntity
	result := r.DB.Where("user_id = ? AND deleted_at IS NULL", userID).Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *ScheduledTaskRepositoryImpl) FindByUserIDAndStatus(userID string, status ScheduleTaskStatus) ([]ScheduledTaskEntity, error) {
	var entities []ScheduledTaskEntity
	result := r.DB.Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID, status).Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *ScheduledTaskRepositoryImpl) FindBySessionID(sessionID string) ([]ScheduledTaskEntity, error) {
	var entities []ScheduledTaskEntity
	result := r.DB.Where("session_id = ? AND deleted_at IS NULL", sessionID).Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *ScheduledTaskRepositoryImpl) FindByAgentID(agentID string) ([]ScheduledTaskEntity, error) {
	var entities []ScheduledTaskEntity
	result := r.DB.Where("agent_id = ? AND deleted_at IS NULL", agentID).Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *ScheduledTaskRepositoryImpl) FindActiveTasksToExecute() ([]ScheduledTaskEntity, error) {
	var entities []ScheduledTaskEntity
	result := r.DB.Where("status = ? AND deleted_at IS NULL", ScheduleTaskStatusActive).Order("created_at ASC").Find(&entities)
	return entities, result.Error
}

func (r *ScheduledTaskRepositoryImpl) Create(e *ScheduledTaskEntity) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return r.DB.Create(e).Error
}

func (r *ScheduledTaskRepositoryImpl) Update(e *ScheduledTaskEntity) error {
	return r.DB.Save(e).Error
}

func (r *ScheduledTaskRepositoryImpl) UpdateStatus(id, userID string, status ScheduleTaskStatus) error {
	result := r.DB.Model(&ScheduledTaskEntity{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).
		Update("status", string(status))
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *ScheduledTaskRepositoryImpl) DeleteByIDAndUserID(id, userID string) error {
	now := time.Now()
	result := r.DB.Model(&ScheduledTaskEntity{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).
		Update("deleted_at", now)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *ScheduledTaskRepositoryImpl) DeleteBySessionIDAndUserID(sessionID, userID string) (int64, error) {
	now := time.Now()
	result := r.DB.Model(&ScheduledTaskEntity{}).
		Where("session_id = ? AND user_id = ? AND deleted_at IS NULL", sessionID, userID).
		Update("deleted_at", now)
	return result.RowsAffected, result.Error
}

func (r *ScheduledTaskRepositoryImpl) DeleteByAgentIDAndUserID(agentID, userID string) (int64, error) {
	now := time.Now()
	result := r.DB.Model(&ScheduledTaskEntity{}).
		Where("agent_id = ? AND user_id = ? AND deleted_at IS NULL", agentID, userID).
		Update("deleted_at", now)
	return result.RowsAffected, result.Error
}

func (r *ScheduledTaskRepositoryImpl) CountByUserID(userID string) (int64, error) {
	var count int64
	result := r.DB.Model(&ScheduledTaskEntity{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&count)
	return count, result.Error
}

func (r *ScheduledTaskRepositoryImpl) CountByUserIDAndStatus(userID string, status ScheduleTaskStatus) (int64, error) {
	var count int64
	result := r.DB.Model(&ScheduledTaskEntity{}).Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID, status).Count(&count)
	return count, result.Error
}

// ---- 领域服务 ----

type DomainService struct {
	repo ScheduledTaskRepository
}

func NewDomainService(repo ScheduledTaskRepository) *DomainService {
	return &DomainService{repo: repo}
}

func (s *DomainService) CreateTask(task *ScheduledTaskEntity) (*ScheduledTaskEntity, error) {
	if err := s.repo.Create(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *DomainService) GetTasksByUserID(userID string) ([]ScheduledTaskEntity, error) {
	return s.repo.FindByUserID(userID)
}

func (s *DomainService) GetTasksByUserIDAndStatus(userID string, status ScheduleTaskStatus) ([]ScheduledTaskEntity, error) {
	return s.repo.FindByUserIDAndStatus(userID, status)
}

func (s *DomainService) GetTasksBySessionID(sessionID string) ([]ScheduledTaskEntity, error) {
	return s.repo.FindBySessionID(sessionID)
}

func (s *DomainService) GetTasksByAgentID(agentID string) ([]ScheduledTaskEntity, error) {
	return s.repo.FindByAgentID(agentID)
}

func (s *DomainService) GetActiveTasksToExecute() ([]ScheduledTaskEntity, error) {
	return s.repo.FindActiveTasksToExecute()
}

func (s *DomainService) UpdateTask(task *ScheduledTaskEntity) error {
	return s.repo.Update(task)
}

func (s *DomainService) DeleteTask(taskID, userID string) error {
	return s.repo.DeleteByIDAndUserID(taskID, userID)
}

func (s *DomainService) PauseTask(taskID, userID string) error {
	return s.repo.UpdateStatus(taskID, userID, ScheduleTaskStatusPaused)
}

func (s *DomainService) ResumeTask(taskID, userID string) error {
	return s.repo.UpdateStatus(taskID, userID, ScheduleTaskStatusActive)
}

func (s *DomainService) CompleteTask(taskID, userID string) error {
	return s.repo.UpdateStatus(taskID, userID, ScheduleTaskStatusCompleted)
}

func (s *DomainService) GetTask(taskID, userID string) (*ScheduledTaskEntity, error) {
	task, err := s.repo.FindByIDAndUserID(taskID, userID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, exception.NewBusinessException("定时任务不存在")
	}
	return task, nil
}

func (s *DomainService) DeleteTasksBySessionID(sessionID, userID string) (int64, error) {
	return s.repo.DeleteBySessionIDAndUserID(sessionID, userID)
}

func (s *DomainService) DeleteTasksByAgentID(agentID, userID string) (int64, error) {
	return s.repo.DeleteByAgentIDAndUserID(agentID, userID)
}
