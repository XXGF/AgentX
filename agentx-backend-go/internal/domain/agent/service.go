package agent

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// DomainService Agent领域服务（对应 Java 的 AgentDomainService）
type DomainService struct {
	agentRepo     AgentRepository
	versionRepo   AgentVersionRepository
	workspaceRepo AgentWorkspaceRepository
}

// NewDomainService 创建Agent领域服务
func NewDomainService(agentRepo AgentRepository, versionRepo AgentVersionRepository, workspaceRepo AgentWorkspaceRepository) *DomainService {
	return &DomainService{
		agentRepo:     agentRepo,
		versionRepo:   versionRepo,
		workspaceRepo: workspaceRepo,
	}
}

// CreateAgent 创建新Agent
func (s *DomainService) CreateAgent(agent *AgentEntity) (*AgentEntity, error) {
	if err := s.agentRepo.Create(agent); err != nil {
		return nil, err
	}
	return agent, nil
}

// GetAgent 获取单个Agent信息
func (s *DomainService) GetAgent(agentID, userID string) (*AgentEntity, error) {
	agent, err := s.agentRepo.FindByIDAndUserID(agentID, userID)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, exception.NewBusinessException("Agent不存在: " + agentID)
	}
	return agent, nil
}

// GetUserAgents 获取用户的Agent列表
func (s *DomainService) GetUserAgents(userID, name string) ([]AgentEntity, error) {
	return s.agentRepo.FindByUserID(userID, name)
}

// GetPublishedAgentsByName 获取已上架的Agent列表
func (s *DomainService) GetPublishedAgentsByName(name string) ([]AgentVersionEntity, error) {
	latestVersions, err := s.versionRepo.FindLatestVersionsByNameAndStatus(name, int(PublishStatusPublished))
	if err != nil {
		return nil, err
	}
	return s.combineAgentsWithVersions(latestVersions)
}

// UpdateAgent 更新Agent信息
func (s *DomainService) UpdateAgent(updateEntity *AgentEntity) (*AgentEntity, error) {
	if err := s.agentRepo.UpdateByIDAndUserID(updateEntity); err != nil {
		return nil, err
	}
	return updateEntity, nil
}

// ToggleAgentStatus 切换Agent的启用/禁用状态
func (s *DomainService) ToggleAgentStatus(agentID string) (*AgentEntity, error) {
	agent, err := s.agentRepo.FindByID(agentID)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, exception.NewBusinessException("Agent不存在: " + agentID)
	}

	if agent.IsEnabled() {
		agent.Disable()
	} else {
		agent.Enable()
	}

	if err := s.agentRepo.Update(agent); err != nil {
		return nil, err
	}
	return agent, nil
}

// DeleteAgent 删除Agent
func (s *DomainService) DeleteAgent(agentID, userID string) error {
	if err := s.agentRepo.DeleteByIDAndUserID(agentID, userID); err != nil {
		return err
	}
	// 删除版本
	return s.versionRepo.DeleteByAgentIDAndUserID(agentID, userID)
}

// PublishAgentVersion 发布Agent版本
func (s *DomainService) PublishAgentVersion(agentID string, versionEntity *AgentVersionEntity) (*AgentVersionEntity, error) {
	agent, err := s.agentRepo.FindByID(agentID)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, exception.NewBusinessException("Agent不存在: " + agentID)
	}

	// 查询最新版本号进行比较
	latestVersion, err := s.versionRepo.FindLatestByAgentID(agentID)
	if err != nil {
		return nil, err
	}

	if latestVersion != nil {
		newVersion := versionEntity.VersionNumber
		oldVersion := latestVersion.VersionNumber

		if newVersion == oldVersion {
			return nil, exception.NewBusinessException("版本号已存在: " + newVersion)
		}

		if !isVersionGreaterThan(newVersion, oldVersion) {
			return nil, exception.NewBusinessException(
				fmt.Sprintf("新版本号(%s)必须大于当前最新版本号(%s)", newVersion, oldVersion))
		}
	}

	versionEntity.AgentID = agentID
	versionEntity.PublishStatus = int(PublishStatusReviewing)

	if err := s.versionRepo.Create(versionEntity); err != nil {
		return nil, err
	}
	return versionEntity, nil
}

// UpdateVersionPublishStatus 更新版本发布状态
func (s *DomainService) UpdateVersionPublishStatus(versionID string, status PublishStatus) (*AgentVersionEntity, error) {
	version, err := s.versionRepo.FindByID(versionID)
	if err != nil {
		return nil, err
	}
	if version == nil {
		return nil, exception.NewBusinessException("版本不存在: " + versionID)
	}

	version.RejectReason = ""
	version.UpdatePublishStatus(status)

	if err := s.versionRepo.Update(version); err != nil {
		return nil, err
	}

	// 如果状态更新为已发布，则绑定为Agent的publishedVersion
	if status == PublishStatusPublished {
		agent, err := s.agentRepo.FindByID(version.AgentID)
		if err != nil {
			return nil, err
		}
		if agent != nil {
			agent.PublishVersion(versionID)
			if err := s.agentRepo.Update(agent); err != nil {
				return nil, err
			}
		}
	}

	return version, nil
}

// RejectVersion 拒绝版本发布
func (s *DomainService) RejectVersion(versionID, reason string) (*AgentVersionEntity, error) {
	version, err := s.versionRepo.FindByID(versionID)
	if err != nil {
		return nil, err
	}
	if version == nil {
		return nil, exception.NewBusinessException("版本不存在: " + versionID)
	}

	version.Reject(reason)
	if err := s.versionRepo.Update(version); err != nil {
		return nil, err
	}
	return version, nil
}

// GetAgentVersions 获取Agent的所有版本
func (s *DomainService) GetAgentVersions(agentID string, userID *string) ([]AgentVersionEntity, error) {
	agent, err := s.agentRepo.FindByID(agentID)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, exception.NewBusinessException("Agent不存在")
	}

	// 如果userID不为空，需要检查权限
	if userID != nil && agent.UserID != *userID {
		return nil, exception.NewBusinessException("Agent不存在或无权访问")
	}

	return s.versionRepo.FindByAgentIDOrderByCreatedAtDesc(agentID)
}

// GetAgentVersion 获取Agent的特定版本
func (s *DomainService) GetAgentVersion(agentID, versionNumber string) (*AgentVersionEntity, error) {
	versions, err := s.versionRepo.FindByAgentID(agentID)
	if err != nil {
		return nil, err
	}
	for _, v := range versions {
		if v.VersionNumber == versionNumber {
			return &v, nil
		}
	}
	return nil, exception.NewBusinessException("Agent版本不存在: " + versionNumber)
}

// GetLatestAgentVersion 获取Agent的最新版本
func (s *DomainService) GetLatestAgentVersion(agentID string) (*AgentVersionEntity, error) {
	return s.versionRepo.FindLatestByAgentID(agentID)
}

// GetPublishedAgentVersion 获取指定Agent的已发布版本
func (s *DomainService) GetPublishedAgentVersion(agentID string) (*AgentVersionEntity, error) {
	return s.versionRepo.FindLatestPublishedByAgentID(agentID)
}

// GetVersionsByStatus 获取指定状态的所有版本
func (s *DomainService) GetVersionsByStatus(status *PublishStatus) ([]AgentVersionEntity, error) {
	var statusInt *int
	if status != nil {
		v := int(*status)
		statusInt = &v
	}
	return s.versionRepo.FindLatestVersionsByStatus(statusInt)
}

// Exist 校验agent是否存在
func (s *DomainService) Exist(agentID, userID string) (bool, error) {
	agent, err := s.agentRepo.FindByIDAndUserID(agentID, userID)
	if err != nil {
		return false, err
	}
	return agent != nil, nil
}

// GetAgentsByIDs 根据agentIds获取agents
func (s *DomainService) GetAgentsByIDs(agentIDs []string) ([]AgentEntity, error) {
	return s.agentRepo.FindByIDs(agentIDs)
}

// GetAgentByID 根据agentId获取agent
func (s *DomainService) GetAgentByID(agentID string) (*AgentEntity, error) {
	agent, err := s.agentRepo.FindByID(agentID)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, exception.NewBusinessException("Agent不存在: " + agentID)
	}
	return agent, nil
}

// GetAgentVersionByID 根据versionId获取版本
func (s *DomainService) GetAgentVersionByID(versionID string) (*AgentVersionEntity, error) {
	return s.versionRepo.FindByID(versionID)
}

// GetAgentStatistics 获取Agent统计信息
func (s *DomainService) GetAgentStatistics() (*AgentStatisticsDTO, error) {
	stats := &AgentStatisticsDTO{}

	totalAgents, err := s.agentRepo.Count(nil)
	if err != nil {
		return nil, err
	}
	stats.TotalAgents = totalAgents

	enabledAgents, err := s.agentRepo.Count(map[string]interface{}{"enabled": true})
	if err != nil {
		return nil, err
	}
	stats.EnabledAgents = enabledAgents

	disabledAgents, err := s.agentRepo.Count(map[string]interface{}{"enabled": false})
	if err != nil {
		return nil, err
	}
	stats.DisabledAgents = disabledAgents

	pendingVersions, err := s.versionRepo.CountByStatus(int(PublishStatusReviewing))
	if err != nil {
		return nil, err
	}
	stats.PendingVersions = pendingVersions

	return stats, nil
}

// GetAgentsPage 分页查询Agent列表
func (s *DomainService) GetAgentsPage(page, pageSize int, keyword string, enabled *bool) ([]AgentEntity, int64, error) {
	return s.agentRepo.FindPage(page, pageSize, keyword, enabled)
}

// === 内部方法 ===

// combineAgentsWithVersions 组合助理和版本信息
func (s *DomainService) combineAgentsWithVersions(versionEntities []AgentVersionEntity) ([]AgentVersionEntity, error) {
	if len(versionEntities) == 0 {
		return []AgentVersionEntity{}, nil
	}

	agentIDs := make([]string, len(versionEntities))
	for i, v := range versionEntities {
		agentIDs[i] = v.AgentID
	}

	agents, err := s.agentRepo.FindByIDs(agentIDs)
	if err != nil {
		return nil, err
	}

	// 过滤出启用的agent
	enabledAgentIDs := make(map[string]bool)
	for _, a := range agents {
		if a.IsEnabled() {
			enabledAgentIDs[a.ID] = true
		}
	}

	// 版本转为map
	versionMap := make(map[string]AgentVersionEntity)
	for _, v := range versionEntities {
		versionMap[v.AgentID] = v
	}

	var result []AgentVersionEntity
	for _, a := range agents {
		if enabledAgentIDs[a.ID] {
			if v, ok := versionMap[a.ID]; ok {
				result = append(result, v)
			}
		}
	}
	return result, nil
}

// AgentStatisticsDTO Agent统计DTO（放在领域层因为领域服务直接返回）
type AgentStatisticsDTO struct {
	TotalAgents     int64 `json:"totalAgents"`
	EnabledAgents   int64 `json:"enabledAgents"`
	DisabledAgents  int64 `json:"disabledAgents"`
	PendingVersions int64 `json:"pendingVersions"`
}

// isVersionGreaterThan 比较版本号大小
func isVersionGreaterThan(newVersion, oldVersion string) bool {
	if oldVersion == "" {
		return true
	}

	current := strings.Split(newVersion, ".")
	last := strings.Split(oldVersion, ".")

	if len(current) != 3 || len(last) != 3 {
		return false
	}

	for i := 0; i < 3; i++ {
		c, _ := strconv.Atoi(current[i])
		l, _ := strconv.Atoi(last[i])
		if c > l {
			return true
		}
		if c < l {
			return false
		}
	}
	return false
}
