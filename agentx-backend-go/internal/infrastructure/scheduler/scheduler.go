package scheduler

import (
	"fmt"
	"sync"
	"time"

	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/scheduledtask"
	"go.uber.org/zap"
)

// TaskScheduleService 任务调度计算服务（对应 Java 的 TaskScheduleService）
// 负责计算定时任务的下次执行时间
type TaskScheduleService struct{}

func NewTaskScheduleService() *TaskScheduleService {
	return &TaskScheduleService{}
}

// CalculateNextExecuteTime 计算下次执行时间
func (s *TaskScheduleService) CalculateNextExecuteTime(task *domain.ScheduledTaskEntity) *time.Time {
	if task.RepeatConfig == nil {
		return nil
	}

	now := time.Now()
	var next time.Time

	switch task.RepeatType {
	case domain.RepeatTypeNone:
		// 一次性任务：使用配置的执行时间
		if task.RepeatConfig.ExecuteDateTime != nil {
			next = *task.RepeatConfig.ExecuteDateTime
			if next.Before(now) {
				return nil // 已过期
			}
		}

	case domain.RepeatTypeDaily:
		// 每日任务：使用 ExecuteTime（格式 "HH:MM"）
		hour, minute := s.parseExecuteTime(task.RepeatConfig.ExecuteTime)
		if hour >= 0 {
			next = time.Date(now.Year(), now.Month(), now.Day(),
				hour, minute, 0, 0, now.Location())
			if next.Before(now) {
				next = next.Add(24 * time.Hour)
			}
		}

	case domain.RepeatTypeWeekly:
		// 每周任务：计算下一个匹配的星期几
		hour, minute := s.parseExecuteTime(task.RepeatConfig.ExecuteTime)
		if task.RepeatConfig.Weekdays != nil && hour >= 0 {
			next = s.findNextWeekday(now, task.RepeatConfig.Weekdays, hour, minute)
		}

	case domain.RepeatTypeMonthly:
		// 每月任务：计算下一个匹配的日期
		hour, minute := s.parseExecuteTime(task.RepeatConfig.ExecuteTime)
		if task.RepeatConfig.MonthDay != nil && hour >= 0 {
			next = time.Date(now.Year(), now.Month(), *task.RepeatConfig.MonthDay,
				hour, minute, 0, 0, now.Location())
			if next.Before(now) {
				next = next.AddDate(0, 1, 0)
			}
		}

	case domain.RepeatTypeWorkdays:
		// 工作日任务：类似每周，但只在周一到周五
		hour, minute := s.parseExecuteTime(task.RepeatConfig.ExecuteTime)
		if hour >= 0 {
			workdays := []int{1, 2, 3, 4, 5} // 周一到周五
			next = s.findNextWeekday(now, workdays, hour, minute)
		}

	case domain.RepeatTypeCustom:
		// 自定义间隔执行
		if task.RepeatConfig.Interval != nil {
			interval := s.calculateInterval(task.RepeatConfig)
			if task.LastExecuteTime != nil {
				next = task.LastExecuteTime.Add(interval)
			} else {
				next = now.Add(interval)
			}
		}

	default:
		return nil
	}

	// 检查是否超过结束时间
	if task.RepeatConfig.EndDateTime != nil && !next.IsZero() && next.After(*task.RepeatConfig.EndDateTime) {
		return nil
	}

	if next.IsZero() {
		return nil
	}
	return &next
}

// parseExecuteTime 解析执行时间字符串（格式 "HH:MM"）
func (s *TaskScheduleService) parseExecuteTime(executeTime string) (hour, minute int) {
	if executeTime == "" {
		return -1, 0
	}
	var h, m int
	n, _ := fmt.Sscanf(executeTime, "%d:%d", &h, &m)
	if n != 2 {
		return -1, 0
	}
	return h, m
}

// calculateInterval 根据配置计算间隔时间
func (s *TaskScheduleService) calculateInterval(config *domain.RepeatConfig) time.Duration {
	if config.Interval == nil {
		return time.Hour
	}
	interval := *config.Interval
	switch config.TimeUnit {
	case "second":
		return time.Duration(interval) * time.Second
	case "minute":
		return time.Duration(interval) * time.Minute
	case "hour":
		return time.Duration(interval) * time.Hour
	case "day":
		return time.Duration(interval) * 24 * time.Hour
	default:
		return time.Duration(interval) * time.Minute
	}
}

// findNextWeekday 查找下一个匹配的星期几
func (s *TaskScheduleService) findNextWeekday(now time.Time, daysOfWeek []int, hour, minute int) time.Time {
	for i := 0; i < 8; i++ {
		candidate := now.AddDate(0, 0, i)
		candidateTime := time.Date(candidate.Year(), candidate.Month(), candidate.Day(),
			hour, minute, 0, 0, now.Location())

		if i == 0 && candidateTime.Before(now) {
			continue
		}

		weekday := int(candidateTime.Weekday())
		for _, d := range daysOfWeek {
			if d == weekday {
				return candidateTime
			}
		}
	}
	return time.Time{}
}

// ---- ScheduledTaskExecutionService 任务执行服务 ----

// TaskExecutor 任务执行器接口
type TaskExecutor interface {
	Execute(task *domain.ScheduledTaskEntity) error
}

// ScheduledTaskExecutionService 定时任务执行服务（对应 Java 的 ScheduledTaskExecutionService）
type ScheduledTaskExecutionService struct {
	domainService   *domain.DomainService
	scheduleService *TaskScheduleService
	executor        TaskExecutor
	logger          *zap.Logger
	stopCh          chan struct{}
	wg              sync.WaitGroup
}

func NewScheduledTaskExecutionService(
	domainService *domain.DomainService,
	scheduleService *TaskScheduleService,
	executor TaskExecutor,
	logger *zap.Logger,
) *ScheduledTaskExecutionService {
	return &ScheduledTaskExecutionService{
		domainService:   domainService,
		scheduleService: scheduleService,
		executor:        executor,
		logger:          logger,
		stopCh:          make(chan struct{}),
	}
}

// Start 启动调度服务
func (s *ScheduledTaskExecutionService) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.runScheduleLoop()
	}()
	s.logger.Info("定时任务调度服务已启动")
}

// Stop 停止调度服务
func (s *ScheduledTaskExecutionService) Stop() {
	close(s.stopCh)
	s.wg.Wait()
	s.logger.Info("定时任务调度服务已停止")
}

// runScheduleLoop 调度循环
func (s *ScheduledTaskExecutionService) runScheduleLoop() {
	ticker := time.NewTicker(30 * time.Second) // 每30秒检查一次
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.checkAndExecuteTasks()
		}
	}
}

// checkAndExecuteTasks 检查并执行到期的任务
func (s *ScheduledTaskExecutionService) checkAndExecuteTasks() {
	now := time.Now()

	// 获取所有活跃的任务
	tasks, err := s.domainService.GetActiveTasksToExecute()
	if err != nil {
		s.logger.Error("获取活跃任务失败", zap.Error(err))
		return
	}

	for _, task := range tasks {
		if task.NextExecuteTime != nil && task.NextExecuteTime.Before(now) {
			s.executeTask(&task)
		}
	}
}

// executeTask 执行单个任务
func (s *ScheduledTaskExecutionService) executeTask(task *domain.ScheduledTaskEntity) {
	s.logger.Info("执行定时任务",
		zap.String("taskId", task.ID),
		zap.String("content", task.Content),
	)

	// 执行任务
	if err := s.executor.Execute(task); err != nil {
		s.logger.Error("执行定时任务失败",
			zap.String("taskId", task.ID),
			zap.Error(err),
		)
	}

	// 更新执行时间
	now := time.Now()
	task.LastExecuteTime = &now

	// 计算下次执行时间
	nextTime := s.scheduleService.CalculateNextExecuteTime(task)
	task.NextExecuteTime = nextTime

	// 如果是一次性任务且已执行，标记为完成
	if task.RepeatType == domain.RepeatTypeNone {
		task.Status = domain.ScheduleTaskStatusCompleted
	}

	// 如果没有下次执行时间，标记为完成
	if nextTime == nil && task.RepeatType != domain.RepeatTypeNone {
		task.Status = domain.ScheduleTaskStatusCompleted
	}

	if err := s.domainService.UpdateTask(task); err != nil {
		s.logger.Error("更新任务状态失败",
			zap.String("taskId", task.ID),
			zap.Error(err),
		)
	}
}

// ---- ChatTaskExecutor 聊天任务执行器 ----

// ChatTaskExecutor 通过聊天接口执行定时任务
// TODO: 集成 ChatAppService 实现实际的聊天执行
type ChatTaskExecutor struct {
	logger *zap.Logger
}

func NewChatTaskExecutor(logger *zap.Logger) *ChatTaskExecutor {
	return &ChatTaskExecutor{logger: logger}
}

func (e *ChatTaskExecutor) Execute(task *domain.ScheduledTaskEntity) error {
	// TODO: 调用 ChatAppService 发送消息到对应的会话
	e.logger.Info("模拟执行聊天任务",
		zap.String("taskId", task.ID),
		zap.String("sessionId", task.SessionID),
		zap.String("content", task.Content),
	)
	return nil
}