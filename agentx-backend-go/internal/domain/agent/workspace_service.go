package agent

import (
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// WorkspaceDomainService Agent工作区领域服务（对应 Java 的 AgentWorkspaceDomainService）
type WorkspaceDomainService struct {
	workspaceRepo AgentWorkspaceRepository
	agentRepo     AgentRepository
}

// NewWorkspaceDomainService 创建工作区领域服务
func NewWorkspaceDomainService(workspaceRepo AgentWorkspaceRepository, agentRepo AgentRepository) *WorkspaceDomainService {
	return &WorkspaceDomainService{
		workspaceRepo: workspaceRepo,
		agentRepo:     agentRepo,
	}
}

// GetWorkspaceAgents 获取工作区中的Agent列表
func (s *WorkspaceDomainService) GetWorkspaceAgents(userID string) ([]AgentEntity, error) {
	workspaces, err := s.workspaceRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	if len(workspaces) == 0 {
		return []AgentEntity{}, nil
	}

	agentIDs := make([]string, len(workspaces))
	for i, w := range workspaces {
		agentIDs[i] = w.AgentID
	}

	return s.agentRepo.FindByIDs(agentIDs)
}

// Exist 检查工作区是否存在
func (s *WorkspaceDomainService) Exist(agentID, userID string) (bool, error) {
	return s.workspaceRepo.Exists(agentID, userID)
}

// DeleteAgent 从工作区删除Agent
func (s *WorkspaceDomainService) DeleteAgent(agentID, userID string) error {
	return s.workspaceRepo.DeleteByAgentIDAndUserID(agentID, userID)
}

// GetWorkspace 获取工作区
func (s *WorkspaceDomainService) GetWorkspace(agentID, userID string) (*AgentWorkspaceEntity, error) {
	workspace, err := s.workspaceRepo.FindByAgentIDAndUserID(agentID, userID)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, exception.NewBusinessException("助理不存在")
	}
	return workspace, nil
}

// FindWorkspace 查找工作区（不抛异常）
func (s *WorkspaceDomainService) FindWorkspace(agentID, userID string) (*AgentWorkspaceEntity, error) {
	return s.workspaceRepo.FindByAgentIDAndUserID(agentID, userID)
}

// Save 保存工作区
func (s *WorkspaceDomainService) Save(workspace *AgentWorkspaceEntity) error {
	return s.workspaceRepo.Create(workspace)
}

// Update 更新工作区
func (s *WorkspaceDomainService) Update(workspace *AgentWorkspaceEntity) error {
	return s.workspaceRepo.Update(workspace)
}

// ListAgents 获取工作区中的Agent列表
func (s *WorkspaceDomainService) ListAgents(agentIDs []string, userID string) ([]AgentWorkspaceEntity, error) {
	return s.workspaceRepo.FindByAgentIDsAndUserID(agentIDs, userID)
}
