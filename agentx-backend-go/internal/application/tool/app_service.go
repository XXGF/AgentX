package tool

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/tool"
	domainUser "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/user"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// AppService 工具应用服务（对应 Java 的 ToolAppService）
type AppService struct {
	toolDomainService        *domain.ToolDomainService
	userToolDomainService    *domain.UserToolDomainService
	toolVersionDomainService *domain.ToolVersionDomainService
	userDomainService        *domainUser.DomainService
}

func NewAppService(
	toolDomainService *domain.ToolDomainService,
	userToolDomainService *domain.UserToolDomainService,
	toolVersionDomainService *domain.ToolVersionDomainService,
	userDomainService *domainUser.DomainService,
) *AppService {
	return &AppService{
		toolDomainService:        toolDomainService,
		userToolDomainService:    userToolDomainService,
		toolVersionDomainService: toolVersionDomainService,
		userDomainService:        userDomainService,
	}
}

// UploadTool 上传工具
func (s *AppService) UploadTool(req *CreateToolRequest, userID string) (*ToolDTO, error) {
	toolEntity := CreateToolRequestToEntity(req, userID)
	toolEntity.Status = domain.ToolStatusWaitingReview

	result, err := s.toolDomainService.CreateTool(toolEntity)
	if err != nil {
		return nil, err
	}

	// TODO: 如果需要状态转换，调用 toolStateStateMachine.submitToolForProcessing
	_ = result.NeedStateTransition

	return ToolEntityToDTO(result.Tool), nil
}

// GetToolDetail 获取工具详情
func (s *AppService) GetToolDetail(toolID, userID string) (*ToolDTO, error) {
	tool, err := s.toolDomainService.GetTool(toolID, userID)
	if err != nil {
		return nil, err
	}
	return ToolEntityToDTO(tool), nil
}

// GetUserTools 获取用户工具列表
func (s *AppService) GetUserTools(userID string) ([]*ToolDTO, error) {
	tools, err := s.toolDomainService.GetUserTools(userID)
	if err != nil {
		return nil, err
	}
	return ToolEntitiesToDTOs(tools), nil
}

// UpdateTool 更新工具
func (s *AppService) UpdateTool(toolID string, req *UpdateToolRequest, userID string) (*ToolDTO, error) {
	toolEntity := UpdateToolRequestToEntity(req, userID)
	toolEntity.ID = toolID

	result, err := s.toolDomainService.UpdateTool(toolEntity)
	if err != nil {
		return nil, err
	}

	// TODO: 如果需要状态转换，调用 toolStateStateMachine.submitToolForProcessing
	_ = result.NeedStateTransition

	return ToolEntityToDTO(result.Tool), nil
}

// DeleteTool 删除工具
func (s *AppService) DeleteTool(toolID, userID string) error {
	return s.toolDomainService.DeleteTool(toolID, userID)
}

// MarketTool 上架工具
func (s *AppService) MarketTool(req *MarketToolRequest, userID string) error {
	tool, err := s.toolDomainService.GetTool(req.ToolID, userID)
	if err != nil {
		return err
	}
	if tool.Status != domain.ToolStatusApproved {
		return exception.NewBusinessException("工具未审核通过，不能上架")
	}

	latestVersion, err := s.toolVersionDomainService.FindLatestToolVersion(req.ToolID)
	if err != nil {
		return err
	}
	if latestVersion != nil {
		if !isVersionGreaterThan(req.Version, latestVersion.Version) {
			return exception.NewBusinessException(
				fmt.Sprintf("新版本号(%s)必须大于当前最新版本号(%s)", req.Version, latestVersion.Version))
		}
	}

	// 创建工具版本
	versionEntity := &domain.ToolVersionEntity{
		Name:          tool.Name,
		Icon:          tool.Icon,
		Subtitle:      tool.Subtitle,
		Description:   tool.Description,
		UserID:        tool.UserID,
		Version:       req.Version,
		ToolID:        req.ToolID,
		UploadType:    tool.UploadType,
		UploadUrl:     tool.UploadUrl,
		ToolList:      tool.ToolList,
		Labels:        domain.JSONStringList(tool.Labels),
		IsOffice:      tool.IsOffice,
		ChangeLog:     req.ChangeLog,
		McpServerName: tool.McpServerName,
	}
	publicStatus := true
	versionEntity.PublicStatus = &publicStatus

	return s.toolVersionDomainService.AddToolVersion(versionEntity)
}

// MarketTools 工具市场列表
func (s *AppService) MarketTools(req *QueryToolRequest) ([]*ToolVersionDTO, int64, error) {
	page := req.GetPageOrDefault()
	pageSize := req.GetPageSizeOrDefault()

	versions, total, err := s.toolVersionDomainService.ListToolVersion(page, pageSize, req.ToolName)
	if err != nil {
		return nil, 0, err
	}

	toolIDs := make([]string, len(versions))
	for i, v := range versions {
		toolIDs[i] = v.ToolID
	}

	installMap, _ := s.userToolDomainService.GetToolsInstall(toolIDs)

	dtos := make([]*ToolVersionDTO, len(versions))
	for i, v := range versions {
		dto := ToolVersionEntityToDTO(&v)
		if count, ok := installMap[v.ToolID]; ok {
			dto.InstallCount = &count
		}
		dtos[i] = dto
	}

	// 获取用户昵称
	userIDs := make([]string, len(dtos))
	for i, dto := range dtos {
		userIDs[i] = dto.UserID
	}
	users, _ := s.userDomainService.GetByIDs(userIDs)
	userNicknameMap := make(map[string]string)
	for _, u := range users {
		userNicknameMap[u.ID] = u.Nickname
	}
	for _, dto := range dtos {
		if name, ok := userNicknameMap[dto.UserID]; ok {
			dto.UserName = name
		}
	}

	return dtos, total, nil
}

// GetToolVersionDetail 获取工具版本详情
func (s *AppService) GetToolVersionDetail(toolID, version, userID string) (*ToolVersionDTO, error) {
	entity, err := s.toolVersionDomainService.GetToolVersion(toolID, version, userID)
	if err != nil {
		return nil, err
	}
	dto := ToolVersionEntityToDTO(entity)

	// 设置创建者昵称
	user, err := s.userDomainService.GetUserInfo(dto.UserID)
	if err == nil && user != nil {
		dto.UserName = user.Nickname
	}

	// 设置历史版本
	allVersions, _ := s.toolVersionDomainService.GetToolVersions(toolID, userID)
	versionDTOs := make([]*ToolVersionDTO, len(allVersions))
	for i, v := range allVersions {
		versionDTOs[i] = ToolVersionEntityToDTO(&v)
	}
	dto.Versions = versionDTOs

	installMap, _ := s.userToolDomainService.GetToolsInstall([]string{toolID})
	if count, ok := installMap[toolID]; ok {
		dto.InstallCount = &count
	}

	return dto, nil
}

// InstallTool 安装工具
func (s *AppService) InstallTool(toolID, version, userID string) error {
	userTool, err := s.userToolDomainService.FindByToolIDAndUserID(toolID, userID)
	if err != nil {
		return err
	}

	toolVersion, err := s.toolVersionDomainService.GetToolVersion(toolID, version, userID)
	if err != nil {
		return err
	}

	if userTool == nil {
		userTool = &domain.UserToolEntity{
			ToolID: toolVersion.ToolID,
		}
	}

	existingID := userTool.ID
	userTool.Name = toolVersion.Name
	userTool.Icon = toolVersion.Icon
	userTool.Subtitle = toolVersion.Subtitle
	userTool.Description = toolVersion.Description
	userTool.UserID = userID
	userTool.Version = toolVersion.Version
	userTool.ToolList = toolVersion.ToolList
	userTool.Labels = domain.JSONStringList(toolVersion.Labels)
	userTool.IsOffice = toolVersion.IsOffice
	userTool.McpServerName = toolVersion.McpServerName
	userTool.ID = existingID

	if userTool.ID == "" {
		return s.userToolDomainService.Add(userTool)
	}
	return s.userToolDomainService.Update(userTool)
}

// GetInstalledTools 获取已安装的工具列表
func (s *AppService) GetInstalledTools(userID string, req *QueryToolRequest) ([]*ToolVersionDTO, int64, error) {
	page := req.GetPageOrDefault()
	pageSize := req.GetPageSizeOrDefault()

	userTools, total, err := s.userToolDomainService.ListByUserID(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	toolIDs := make([]string, len(userTools))
	for i, ut := range userTools {
		toolIDs[i] = ut.ToolID
	}

	toolMap := make(map[string]bool)
	tools, _ := s.toolDomainService.GetByIDs(toolIDs)
	for _, t := range tools {
		toolMap[t.ID] = true
	}

	dtos := make([]*ToolVersionDTO, len(userTools))
	for i, ut := range userTools {
		dto := UserToolEntityToDTO(&ut)
		if !toolMap[ut.ToolID] {
			dto.IsDelete = true
		}
		dtos[i] = dto
	}

	return dtos, total, nil
}

// GetToolVersions 获取工具已发布的所有版本
func (s *AppService) GetToolVersions(toolID, userID string) ([]*ToolVersionDTO, error) {
	versions, err := s.toolVersionDomainService.GetToolVersions(toolID, userID)
	if err != nil {
		return nil, err
	}
	dtos := make([]*ToolVersionDTO, len(versions))
	for i, v := range versions {
		dtos[i] = ToolVersionEntityToDTO(&v)
	}
	return dtos, nil
}

// UninstallTool 卸载工具
func (s *AppService) UninstallTool(toolID, userID string) error {
	tool, err := s.toolDomainService.GetToolByID(toolID)
	if err == nil && tool != nil && tool.UserID == userID {
		return exception.NewBusinessException("不允许卸载自己创建的工具")
	}
	return s.userToolDomainService.Delete(toolID, userID)
}

// GetRecommendTools 获取推荐工具
func (s *AppService) GetRecommendTools() ([]*ToolVersionDTO, error) {
	versions, _, err := s.toolVersionDomainService.ListToolVersion(1, 9999, "")
	if err != nil {
		return nil, err
	}

	toolIDs := make([]string, len(versions))
	for i, v := range versions {
		toolIDs[i] = v.ToolID
	}
	installMap, _ := s.userToolDomainService.GetToolsInstall(toolIDs)

	dtos := make([]*ToolVersionDTO, len(versions))
	for i, v := range versions {
		dto := ToolVersionEntityToDTO(&v)
		if count, ok := installMap[v.ToolID]; ok {
			dto.InstallCount = &count
		}
		dtos[i] = dto
	}

	// 随机选取最多10条
	if len(dtos) > 10 {
		rand.Shuffle(len(dtos), func(i, j int) {
			dtos[i], dtos[j] = dtos[j], dtos[i]
		})
		dtos = dtos[:10]
	}

	// 获取用户昵称
	userIDs := make([]string, len(dtos))
	for i, dto := range dtos {
		userIDs[i] = dto.UserID
	}
	users, _ := s.userDomainService.GetByIDs(userIDs)
	userNicknameMap := make(map[string]string)
	for _, u := range users {
		userNicknameMap[u.ID] = u.Nickname
	}
	for _, dto := range dtos {
		if name, ok := userNicknameMap[dto.UserID]; ok {
			dto.UserName = name
		}
	}

	return dtos, nil
}

// UpdateUserToolVersionStatus 修改工具版本发布状态
func (s *AppService) UpdateUserToolVersionStatus(toolID, version string, publishStatus bool, userID string) error {
	return s.toolVersionDomainService.UpdateToolVersionStatus(toolID, version, userID, publishStatus)
}

// GetLatestToolVersion 获取工具最新版本
func (s *AppService) GetLatestToolVersion(toolID string) (*ToolVersionDTO, error) {
	entity, err := s.toolVersionDomainService.FindLatestToolVersion(toolID)
	if err != nil {
		return nil, err
	}
	return ToolVersionEntityToDTO(entity), nil
}

// AutoInstallApprovedTool 为工具创建者自动安装审核通过的工具
func (s *AppService) AutoInstallApprovedTool(toolID string) error {
	tool, err := s.toolDomainService.GetToolByID(toolID)
	if err != nil || tool == nil || tool.Status != domain.ToolStatusApproved {
		return nil
	}

	ownerID := tool.UserID
	existing, _ := s.userToolDomainService.FindByToolIDAndUserID(toolID, ownerID)
	if existing != nil {
		return nil
	}

	versionToInstall, _ := s.toolVersionDomainService.FindLatestToolVersion(toolID)
	if versionToInstall == nil {
		// 创建内部基础版本
		baseVersion := &domain.ToolVersionEntity{
			Name:          tool.Name,
			Icon:          tool.Icon,
			Subtitle:      tool.Subtitle,
			Description:   tool.Description,
			UserID:        ownerID,
			ToolID:        toolID,
			Version:       "0.0.0",
			ChangeLog:     "Base configuration for owner auto-installation.",
			McpServerName: tool.McpServerName,
			UploadType:    tool.UploadType,
			UploadUrl:     tool.UploadUrl,
			ToolList:      tool.ToolList,
			Labels:        domain.JSONStringList(tool.Labels),
		}
		publicStatus := false
		baseVersion.PublicStatus = &publicStatus
		baseVersion.CreatedAt = time.Now()
		baseVersion.UpdatedAt = time.Now()

		if err := s.toolVersionDomainService.AddToolVersion(baseVersion); err != nil {
			return err
		}
		versionToInstall = baseVersion
	}

	return s.InstallTool(toolID, versionToInstall.Version, ownerID)
}

// ---- 辅助函数 ----

var versionPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

func isVersionGreaterThan(current, last string) bool {
	if last == "" {
		return true
	}
	if !versionPattern.MatchString(current) || !versionPattern.MatchString(last) {
		return false
	}
	currentParts := strings.Split(current, ".")
	lastParts := strings.Split(last, ".")

	for i := 0; i < 3; i++ {
		c, _ := strconv.Atoi(currentParts[i])
		l, _ := strconv.Atoi(lastParts[i])
		if c > l {
			return true
		}
		if c < l {
			return false
		}
	}
	return false
}
