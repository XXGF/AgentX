package agent

import (
	"fmt"
	"regexp"
	"time"

	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/agent"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	"github.com/google/uuid"
)

// 版本号正则
var versionPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

// AppService Agent应用服务（对应 Java 的 AgentAppService）
type AppService struct {
	agentDomainService     *domain.DomainService
	workspaceDomainService *domain.WorkspaceDomainService
}

// NewAppService 创建Agent应用服务
func NewAppService(agentDomainService *domain.DomainService, workspaceDomainService *domain.WorkspaceDomainService) *AppService {
	return &AppService{
		agentDomainService:     agentDomainService,
		workspaceDomainService: workspaceDomainService,
	}
}

// CreateAgent 创建新Agent
func (s *AppService) CreateAgent(req *CreateAgentRequest, userID string) (*AgentDTO, error) {
	// 注意：计费检查可在此处集成 BillingService.CheckBalance()
	// 当前版本暂不强制计费检查，后续可通过中间件统一处理

	entity := CreateRequestToEntity(req, userID)
	agent, err := s.agentDomainService.CreateAgent(entity)
	if err != nil {
		return nil, err
	}

	// 创建工作区
	workspace := &domain.AgentWorkspaceEntity{
		ID:             uuid.New().String(),
		AgentID:        agent.ID,
		UserID:         userID,
		LLMModelConfig: domain.NewDefaultLLMModelConfig(),
	}
	if err := s.workspaceDomainService.Save(workspace); err != nil {
		return nil, err
	}

	return AgentToDTO(agent), nil
}

// GetAgent 获取Agent信息
func (s *AppService) GetAgent(agentID, userID string) (*AgentDTO, error) {
	agent, err := s.agentDomainService.GetAgent(agentID, userID)
	if err != nil {
		return nil, err
	}
	return AgentToDTO(agent), nil
}

// GetUserAgents 获取用户的Agent列表
func (s *AppService) GetUserAgents(userID string, req *SearchAgentsRequest) ([]*AgentDTO, error) {
	name := ""
	if req != nil {
		name = req.Name
	}
	agents, err := s.agentDomainService.GetUserAgents(userID, name)
	if err != nil {
		return nil, err
	}
	return AgentsToDTOs(agents), nil
}

// GetPublishedAgentsByName 获取已上架的Agent列表
func (s *AppService) GetPublishedAgentsByName(req *SearchAgentsRequest, userID string) ([]*AgentVersionDTO, error) {
	name := ""
	if req != nil {
		name = req.Name
	}
	versions, err := s.agentDomainService.GetPublishedAgentsByName(name)
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return []*AgentVersionDTO{}, nil
	}

	// 获取工作区信息
	agentIDs := make([]string, len(versions))
	for i, v := range versions {
		agentIDs[i] = v.AgentID
	}
	workspaces, err := s.workspaceDomainService.ListAgents(agentIDs, userID)
	if err != nil {
		return nil, err
	}
	workspaceSet := make(map[string]bool)
	for _, w := range workspaces {
		workspaceSet[w.AgentID] = true
	}

	dtos := VersionsToDTOs(versions)
	for _, dto := range dtos {
		added := workspaceSet[dto.AgentID]
		dto.IsAddWorkspace = &added
	}
	return dtos, nil
}

// UpdateAgent 更新Agent信息
func (s *AppService) UpdateAgent(req *UpdateAgentRequest, userID string) (*AgentDTO, error) {
	entity := UpdateRequestToEntity(req, userID)
	agent, err := s.agentDomainService.UpdateAgent(entity)
	if err != nil {
		return nil, err
	}
	return AgentToDTO(agent), nil
}

// ToggleAgentStatus 切换Agent的启用/禁用状态
func (s *AppService) ToggleAgentStatus(agentID string) (*AgentDTO, error) {
	agent, err := s.agentDomainService.ToggleAgentStatus(agentID)
	if err != nil {
		return nil, err
	}
	return AgentToDTO(agent), nil
}

// DeleteAgent 删除Agent
func (s *AppService) DeleteAgent(agentID, userID string) error {
	// 注意：删除关联的定时任务可在此处集成 ScheduledTaskExecutionService
	// 当前版本暂不自动清理，后续可通过事件驱动机制处理
	return s.agentDomainService.DeleteAgent(agentID, userID)
}

// PublishAgentVersion 发布Agent版本
func (s *AppService) PublishAgentVersion(agentID string, req *PublishAgentVersionRequest, userID string) (*AgentVersionDTO, error) {
	// 验证版本号格式
	if !versionPattern.MatchString(req.VersionNumber) {
		return nil, exception.NewBusinessException("版本号必须遵循 x.y.z 格式")
	}

	// 获取当前Agent
	agent, err := s.agentDomainService.GetAgent(agentID, userID)
	if err != nil {
		return nil, err
	}

	// 获取最新版本，检查版本号大小
	latestVersion, err := s.agentDomainService.GetLatestAgentVersion(agentID)
	if err != nil {
		return nil, err
	}
	if latestVersion != nil {
		if req.VersionNumber == latestVersion.VersionNumber {
			return nil, exception.NewBusinessException("版本号已存在: " + req.VersionNumber)
		}
	}

	// 创建版本实体
	versionEntity := CreateVersionEntity(agent, req)
	versionEntity.UserID = userID

	// 注意：验证Agent依赖的工具和知识库权限可在此处集成
	// 当前版本暂不强制验证，后续可通过权限中间件统一处理

	version, err := s.agentDomainService.PublishAgentVersion(agentID, versionEntity)
	if err != nil {
		return nil, err
	}
	return VersionToDTO(version), nil
}

// GetAgentVersions 获取Agent的所有版本
func (s *AppService) GetAgentVersions(agentID string, userID *string) ([]*AgentVersionDTO, error) {
	versions, err := s.agentDomainService.GetAgentVersions(agentID, userID)
	if err != nil {
		return nil, err
	}
	return VersionsToDTOs(versions), nil
}

// GetAgentVersion 获取Agent的特定版本
func (s *AppService) GetAgentVersion(agentID, versionNumber string) (*AgentVersionDTO, error) {
	version, err := s.agentDomainService.GetAgentVersion(agentID, versionNumber)
	if err != nil {
		return nil, err
	}
	return VersionToDTO(version), nil
}

// GetLatestAgentVersion 获取Agent的最新版本
func (s *AppService) GetLatestAgentVersion(agentID string) (*AgentVersionDTO, error) {
	version, err := s.agentDomainService.GetLatestAgentVersion(agentID)
	if err != nil {
		return nil, err
	}
	return VersionToDTO(version), nil
}

// ReviewAgentVersion 审核Agent版本
func (s *AppService) ReviewAgentVersion(versionID string, req *ReviewAgentVersionRequest) (*AgentVersionDTO, error) {
	var version *domain.AgentVersionEntity
	var err error

	if req.Status == domain.PublishStatusRejected {
		version, err = s.agentDomainService.RejectVersion(versionID, req.RejectReason)
	} else {
		version, err = s.agentDomainService.UpdateVersionPublishStatus(versionID, req.Status)
	}
	if err != nil {
		return nil, err
	}
	return VersionToDTO(version), nil
}

// GetVersionsByStatus 根据发布状态获取版本列表
func (s *AppService) GetVersionsByStatus(status *domain.PublishStatus) ([]*AgentVersionDTO, error) {
	versions, err := s.agentDomainService.GetVersionsByStatus(status)
	if err != nil {
		return nil, err
	}
	return VersionsToDTOs(versions), nil
}

// GetAgentStatistics 获取Agent统计信息
func (s *AppService) GetAgentStatistics() (*domain.AgentStatisticsDTO, error) {
	return s.agentDomainService.GetAgentStatistics()
}

// GetAgentsPage 分页查询Agent列表（管理员使用）
func (s *AppService) GetAgentsPage(req *QueryAgentRequest) ([]AgentDTO, int64, error) {
	page := 1
	pageSize := 15
	if req.Page > 0 {
		page = req.Page
	}
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}

	agents, total, err := s.agentDomainService.GetAgentsPage(page, pageSize, req.Keyword, req.Enabled)
	if err != nil {
		return nil, 0, err
	}

	dtos := make([]AgentDTO, len(agents))
	for i, a := range agents {
		dto := AgentToDTO(&a)
		dtos[i] = *dto
	}
	return dtos, total, nil
}

// generateRequestID 生成请求ID
func generateRequestID(userID, action string) string {
	return fmt.Sprintf("agent_%s_%s_%d", action, userID, time.Now().UnixMilli())
}
