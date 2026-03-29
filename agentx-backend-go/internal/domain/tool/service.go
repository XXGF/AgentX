package tool

import (
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// ToolDomainService 工具领域服务（对应 Java 的 ToolDomainService）
type ToolDomainService struct {
	toolRepo     ToolRepository
	versionRepo  ToolVersionRepository
	userToolRepo UserToolRepository
}

func NewToolDomainService(toolRepo ToolRepository, versionRepo ToolVersionRepository, userToolRepo UserToolRepository) *ToolDomainService {
	return &ToolDomainService{
		toolRepo:     toolRepo,
		versionRepo:  versionRepo,
		userToolRepo: userToolRepo,
	}
}

// CreateTool 创建工具
func (s *ToolDomainService) CreateTool(tool *ToolEntity) (*ToolOperationResult, error) {
	tool.Status = ToolStatusWaitingReview
	tool.ID = uuid.New().String()

	mcpServerName, err := s.getMcpServerName(tool)
	if err != nil {
		return nil, err
	}
	tool.McpServerName = mcpServerName

	// 校验MCP名称唯一性
	if err := s.validateMcpServerNameUnique(mcpServerName, tool.UserID, nil); err != nil {
		return nil, err
	}

	if err := s.toolRepo.Create(tool); err != nil {
		return nil, err
	}
	return NewToolOperationResult(tool, true), nil
}

// GetTool 获取工具（带用户校验）
func (s *ToolDomainService) GetTool(toolID, userID string) (*ToolEntity, error) {
	tool, err := s.toolRepo.FindByIDAndUserID(toolID, userID)
	if err != nil {
		return nil, err
	}
	if tool == nil {
		return nil, exception.NewBusinessException("工具不存在: " + toolID)
	}
	return tool, nil
}

// GetToolByID 获取工具（不校验用户）
func (s *ToolDomainService) GetToolByID(toolID string) (*ToolEntity, error) {
	tool, err := s.toolRepo.FindByID(toolID)
	if err != nil {
		return nil, err
	}
	if tool == nil {
		return nil, exception.NewBusinessException("工具不存在: " + toolID)
	}
	return tool, nil
}

// GetUserTools 获取用户的工具列表
func (s *ToolDomainService) GetUserTools(userID string) ([]ToolEntity, error) {
	return s.toolRepo.FindByUserID(userID)
}

// UpdateTool 更新工具
func (s *ToolDomainService) UpdateTool(tool *ToolEntity) (*ToolOperationResult, error) {
	oldTool, err := s.toolRepo.FindByID(tool.ID)
	if err != nil {
		return nil, err
	}
	if oldTool == nil {
		return nil, exception.NewBusinessException("工具不存在: " + tool.ID)
	}

	needStateTransition := false
	if (tool.UploadUrl != "" && tool.UploadUrl != oldTool.UploadUrl) ||
		(tool.InstallCommand != nil && len(tool.InstallCommand) > 0) {
		needStateTransition = true
		mcpServerName, err := s.getMcpServerName(tool)
		if err != nil {
			return nil, err
		}
		tool.McpServerName = mcpServerName
		if err := s.validateMcpServerNameUnique(mcpServerName, tool.UserID, &tool.ID); err != nil {
			return nil, err
		}
		tool.Status = ToolStatusWaitingReview
	} else {
		tool.Status = ToolStatusManualReview
	}

	if err := s.toolRepo.Update(tool); err != nil {
		return nil, err
	}
	return NewToolOperationResult(tool, needStateTransition), nil
}

// DeleteTool 删除工具
func (s *ToolDomainService) DeleteTool(toolID, userID string) error {
	if err := s.toolRepo.DeleteByIDAndUserID(toolID, userID); err != nil {
		return err
	}
	// 删除用户安装的该工具
	_ = s.userToolRepo.DeleteByToolIDAndUserID(toolID, userID)
	return nil
}

// UpdateToolEntity 更新工具实体
func (s *ToolDomainService) UpdateToolEntity(tool *ToolEntity) error {
	if tool == nil || tool.ID == "" {
		return exception.NewBusinessException("工具实体或工具ID不能为空")
	}
	return s.toolRepo.Update(tool)
}

// UpdateApprovedToolStatus 更新已审核工具状态
func (s *ToolDomainService) UpdateApprovedToolStatus(toolID string, status ToolStatus) (*ToolEntity, error) {
	if err := s.toolRepo.UpdateFields(toolID, map[string]interface{}{"status": string(status)}); err != nil {
		return nil, err
	}
	return s.toolRepo.FindByID(toolID)
}

// UpdateFailedToolStatus 更新失败工具状态
func (s *ToolDomainService) UpdateFailedToolStatus(toolID string, failedStepStatus ToolStatus, rejectReason string) (*ToolEntity, error) {
	fields := map[string]interface{}{
		"failed_step_status": string(failedStepStatus),
		"reject_reason":      rejectReason,
		"status":             string(ToolStatusFailed),
	}
	if err := s.toolRepo.UpdateFields(toolID, fields); err != nil {
		return nil, err
	}
	return s.toolRepo.FindByID(toolID)
}

// ManualReviewComplete 人工审核完成
func (s *ToolDomainService) ManualReviewComplete(tool *ToolEntity, approved bool) (string, error) {
	if tool == nil {
		return "", exception.NewBusinessException("工具不存在")
	}
	if tool.Status != ToolStatusManualReview {
		return "", exception.NewBusinessException("工具当前状态不是MANUAL_REVIEW，无法进行人工审核操作")
	}
	if approved {
		tool.Status = ToolStatusApproved
	} else {
		tool.Status = ToolStatusFailed
	}
	if err := s.UpdateToolEntity(tool); err != nil {
		return "", err
	}
	return tool.ID, nil
}

// TransitionToStatus 转换工具状态
func (s *ToolDomainService) TransitionToStatus(tool *ToolEntity, targetStatus ToolStatus) error {
	if tool == nil {
		return exception.NewBusinessException("工具不存在")
	}
	if tool.Status == targetStatus {
		return nil
	}
	tool.Status = targetStatus
	return s.UpdateToolEntity(tool)
}

// GetByIDs 根据ID列表获取工具
func (s *ToolDomainService) GetByIDs(toolIDs []string) ([]ToolEntity, error) {
	return s.toolRepo.FindByIDs(toolIDs)
}

// GetTools 分页查询工具列表
func (s *ToolDomainService) GetTools(page, pageSize int, keyword string, status *ToolStatus, isOffice *bool) ([]ToolEntity, int64, error) {
	conditions := make(map[string]interface{})
	if status != nil {
		conditions["status"] = string(*status)
	}
	if isOffice != nil {
		conditions["is_office"] = *isOffice
	}
	return s.toolRepo.FindPage(page, pageSize, conditions, keyword)
}

// GetToolStatistics 获取工具统计信息
func (s *ToolDomainService) GetToolStatistics() (*ToolStatisticsDTO, error) {
	stats := &ToolStatisticsDTO{}

	total, err := s.toolRepo.Count(nil)
	if err != nil {
		return nil, err
	}
	stats.TotalTools = total

	pending, _ := s.toolRepo.Count(map[string]interface{}{"status": string(ToolStatusWaitingReview)})
	stats.PendingReviewTools = pending

	manual, _ := s.toolRepo.Count(map[string]interface{}{"status": string(ToolStatusManualReview)})
	stats.ManualReviewTools = manual

	approved, _ := s.toolRepo.Count(map[string]interface{}{"status": string(ToolStatusApproved)})
	stats.ApprovedTools = approved

	failed, _ := s.toolRepo.Count(map[string]interface{}{"status": string(ToolStatusFailed)})
	stats.FailedTools = failed

	official, _ := s.toolRepo.Count(map[string]interface{}{"is_office": true})
	stats.OfficialTools = official

	return stats, nil
}

// UpdateToolGlobalStatus 更新工具全局状态
func (s *ToolDomainService) UpdateToolGlobalStatus(toolID string, isGlobal bool) error {
	if err := s.toolRepo.UpdateFields(toolID, map[string]interface{}{"is_global": isGlobal}); err != nil {
		return err
	}
	return s.userToolRepo.UpdateGlobalStatusByToolID(toolID, isGlobal)
}

func (s *ToolDomainService) getMcpServerName(tool *ToolEntity) (string, error) {
	if tool == nil || tool.InstallCommand == nil {
		return "", exception.NewBusinessException("安装命令不能为空")
	}
	mcpServers, ok := tool.InstallCommand["mcpServers"]
	if !ok {
		return "", exception.NewBusinessException("安装命令中mcpServers为空")
	}
	mcpServersMap, ok := mcpServers.(map[string]interface{})
	if !ok || len(mcpServersMap) == 0 {
		return "", exception.NewBusinessException("安装命令中mcpServers为空")
	}
	for name := range mcpServersMap {
		return name, nil
	}
	return "", exception.NewBusinessException("无法从安装命令中获取工具名称")
}

func (s *ToolDomainService) validateMcpServerNameUnique(mcpServerName, userID string, excludeToolID *string) error {
	count, err := s.userToolRepo.CountByMcpServerNameAndUserID(mcpServerName, userID, excludeToolID)
	if err != nil {
		return err
	}
	if count > 0 {
		return exception.NewBusinessException("MCP服务器名称 '" + mcpServerName + "' 与已安装工具冲突，请使用其他名称")
	}
	return nil
}

// ToolStatisticsDTO 工具统计信息（对应 Java 的 ToolStatisticsDTO）
type ToolStatisticsDTO struct {
	TotalTools        int64 `json:"totalTools"`
	PendingReviewTools int64 `json:"pendingReviewTools"`
	ManualReviewTools  int64 `json:"manualReviewTools"`
	ApprovedTools      int64 `json:"approvedTools"`
	FailedTools        int64 `json:"failedTools"`
	OfficialTools      int64 `json:"officialTools"`
}

// UserToolDomainService 用户已安装工具领域服务（对应 Java 的 UserToolDomainService）
type UserToolDomainService struct {
	userToolRepo UserToolRepository
}

func NewUserToolDomainService(userToolRepo UserToolRepository) *UserToolDomainService {
	return &UserToolDomainService{userToolRepo: userToolRepo}
}

func (s *UserToolDomainService) Add(entity *UserToolEntity) error {
	entity.ID = uuid.New().String()
	return s.userToolRepo.Create(entity)
}

func (s *UserToolDomainService) FindByToolIDAndUserID(toolID, userID string) (*UserToolEntity, error) {
	return s.userToolRepo.FindByToolIDAndUserID(toolID, userID)
}

func (s *UserToolDomainService) Update(entity *UserToolEntity) error {
	return s.userToolRepo.Update(entity)
}

func (s *UserToolDomainService) Delete(toolID, userID string) error {
	return s.userToolRepo.DeleteByToolIDAndUserID(toolID, userID)
}

func (s *UserToolDomainService) ListByUserID(userID string, page, pageSize int) ([]UserToolEntity, int64, error) {
	return s.userToolRepo.FindByUserID(userID, page, pageSize)
}

func (s *UserToolDomainService) GetToolsInstall(toolIDs []string) (map[string]int64, error) {
	return s.userToolRepo.CountByToolIDs(toolIDs)
}

func (s *UserToolDomainService) GetInstallTool(toolIDs []string, userID string) ([]UserToolEntity, error) {
	return s.userToolRepo.FindByUserIDAndToolIDs(userID, toolIDs)
}

func (s *UserToolDomainService) GetUserToolsByIDs(userID string, toolIDs []string) ([]UserToolEntity, error) {
	return s.userToolRepo.FindByUserIDAndToolIDs(userID, toolIDs)
}

// ToolVersionDomainService 工具版本领域服务（对应 Java 的 ToolVersionDomainService）
type ToolVersionDomainService struct {
	versionRepo ToolVersionRepository
}

func NewToolVersionDomainService(versionRepo ToolVersionRepository) *ToolVersionDomainService {
	return &ToolVersionDomainService{versionRepo: versionRepo}
}

// ListToolVersion 获取工具市场列表（按tool_id分组取最新公开版本）
func (s *ToolVersionDomainService) ListToolVersion(page, pageSize int, toolName string) ([]ToolVersionEntity, int64, error) {
	allPublic, err := s.versionRepo.FindPublicVersions(toolName)
	if err != nil {
		return nil, 0, err
	}

	// 按tool_id分组，取每组created_at最大的一条
	latestMap := make(map[string]ToolVersionEntity)
	for _, v := range allPublic {
		existing, ok := latestMap[v.ToolID]
		if !ok || v.CreatedAt.After(existing.CreatedAt) {
			latestMap[v.ToolID] = v
		}
	}

	latestList := make([]ToolVersionEntity, 0, len(latestMap))
	for _, v := range latestMap {
		latestList = append(latestList, v)
	}

	// 按创建时间倒序
	sort.Slice(latestList, func(i, j int) bool {
		return latestList[i].CreatedAt.After(latestList[j].CreatedAt)
	})

	total := int64(len(latestList))

	// 手动分页
	fromIndex := (page - 1) * pageSize
	toIndex := fromIndex + pageSize
	if fromIndex >= len(latestList) {
		return []ToolVersionEntity{}, total, nil
	}
	if toIndex > len(latestList) {
		toIndex = len(latestList)
	}

	return latestList[fromIndex:toIndex], total, nil
}

// GetToolVersion 获取工具版本（带权限验证）
func (s *ToolVersionDomainService) GetToolVersion(toolID, version, userID string) (*ToolVersionEntity, error) {
	entity, err := s.versionRepo.FindByToolIDAndVersion(toolID, version)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, exception.NewBusinessException("工具版本不存在: " + toolID + " " + version)
	}
	if entity.UserID != userID && (entity.PublicStatus == nil || !*entity.PublicStatus) {
		return nil, exception.NewBusinessException("该工具版本未公开，无权访问")
	}
	return entity, nil
}

// FindLatestToolVersion 获取最新工具版本
func (s *ToolVersionDomainService) FindLatestToolVersion(toolID string) (*ToolVersionEntity, error) {
	return s.versionRepo.FindLatestByToolID(toolID)
}

// AddToolVersion 添加工具版本
func (s *ToolVersionDomainService) AddToolVersion(entity *ToolVersionEntity) error {
	entity.ID = uuid.New().String()
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()
	return s.versionRepo.Create(entity)
}

// GetToolVersions 获取工具的所有版本
func (s *ToolVersionDomainService) GetToolVersions(toolID, userID string) ([]ToolVersionEntity, error) {
	versions, err := s.versionRepo.FindByToolID(toolID)
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, exception.NewBusinessException("工具版本不存在")
	}

	// 如果不是创建者，只返回公开版本
	if versions[0].UserID != userID {
		var publicVersions []ToolVersionEntity
		for _, v := range versions {
			if v.PublicStatus != nil && *v.PublicStatus {
				publicVersions = append(publicVersions, v)
			}
		}
		return publicVersions, nil
	}
	return versions, nil
}

// UpdateToolVersionStatus 更新工具版本发布状态
func (s *ToolVersionDomainService) UpdateToolVersionStatus(toolID, version, userID string, publishStatus bool) error {
	return s.versionRepo.UpdatePublicStatus(toolID, version, userID, publishStatus)
}
