package tool

import (
	"encoding/json"
	"time"

	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/tool"
)

// ToolDTO 工具DTO（对应 Java 的 ToolDTO）
type ToolDTO struct {
	ID               string                    `json:"id"`
	Name             string                    `json:"name"`
	Icon             string                    `json:"icon"`
	Subtitle         string                    `json:"subtitle"`
	Description      string                    `json:"description"`
	UserID           string                    `json:"userId"`
	UserName         string                    `json:"userName,omitempty"`
	Labels           []string                  `json:"labels"`
	ToolType         domain.ToolType           `json:"toolType"`
	UploadType       domain.UploadType         `json:"uploadType"`
	UploadUrl        string                    `json:"uploadUrl"`
	ToolList         domain.ToolDefinitionList `json:"toolList"`
	Status           domain.ToolStatus         `json:"status"`
	IsOffice         *bool                     `json:"isOffice"`
	InstallCount     int                       `json:"installCount,omitempty"`
	CurrentVersion   string                    `json:"currentVersion,omitempty"`
	InstallCommand   string                    `json:"installCommand,omitempty"`
	CreatedAt        time.Time                 `json:"createdAt"`
	UpdatedAt        time.Time                 `json:"updatedAt"`
	RejectReason     string                    `json:"rejectReason,omitempty"`
	FailedStepStatus domain.ToolStatus         `json:"failedStepStatus,omitempty"`
	McpServerName    string                    `json:"mcpServerName,omitempty"`
	IsGlobal         *bool                     `json:"isGlobal"`
}

// ToolVersionDTO 工具版本DTO（对应 Java 的 ToolVersionDTO）
type ToolVersionDTO struct {
	ID            string                    `json:"id"`
	Name          string                    `json:"name"`
	Icon          string                    `json:"icon"`
	Subtitle      string                    `json:"subtitle"`
	Description   string                    `json:"description"`
	UserID        string                    `json:"userId"`
	Version       string                    `json:"version"`
	ToolID        string                    `json:"toolId"`
	UploadType    string                    `json:"uploadType"`
	UploadUrl     string                    `json:"uploadUrl"`
	ToolList      domain.ToolDefinitionList `json:"toolList"`
	Labels        []string                  `json:"labels"`
	IsOffice      *bool                     `json:"isOffice"`
	PublicStatus  *bool                     `json:"publicStatus"`
	ChangeLog     string                    `json:"changeLog,omitempty"`
	CreatedAt     time.Time                 `json:"createdAt"`
	UpdatedAt     time.Time                 `json:"updatedAt"`
	UserName      string                    `json:"userName,omitempty"`
	Versions      []*ToolVersionDTO         `json:"versions,omitempty"`
	InstallCount  *int64                    `json:"installCount,omitempty"`
	McpServerName string                    `json:"mcpServerName,omitempty"`
	IsDelete      bool                      `json:"isDelete,omitempty"`
}

// ToolWithUserDTO 包含用户信息的工具DTO（对应 Java 的 ToolWithUserDTO）
type ToolWithUserDTO struct {
	ToolDTO
	UserNickname  string `json:"userNickname,omitempty"`
	UserEmail     string `json:"userEmail,omitempty"`
	UserAvatarUrl string `json:"userAvatarUrl,omitempty"`
}

// CreateToolRequest 创建工具请求（对应 Java 的 CreateToolRequest）
type CreateToolRequest struct {
	Name           string                 `json:"name" binding:"required"`
	Icon           string                 `json:"icon"`
	Subtitle       string                 `json:"subtitle" binding:"required"`
	Description    string                 `json:"description" binding:"required"`
	Labels         []string               `json:"labels" binding:"required"`
	UploadUrl      string                 `json:"uploadUrl" binding:"required"`
	InstallCommand map[string]interface{} `json:"installCommand" binding:"required"`
	IsGlobal       *bool                  `json:"isGlobal"`
}

// UpdateToolRequest 更新工具请求（对应 Java 的 UpdateToolRequest）
type UpdateToolRequest struct {
	Name           string                 `json:"name" binding:"required"`
	Icon           string                 `json:"icon"`
	Subtitle       string                 `json:"subtitle" binding:"required"`
	Description    string                 `json:"description" binding:"required"`
	Labels         []string               `json:"labels" binding:"required"`
	UploadUrl      string                 `json:"uploadUrl" binding:"required"`
	InstallCommand map[string]interface{} `json:"installCommand" binding:"required"`
	IsGlobal       *bool                  `json:"isGlobal"`
}

// MarketToolRequest 上架工具请求（对应 Java 的 MarketToolRequest）
type MarketToolRequest struct {
	ToolID    string `json:"toolId" binding:"required"`
	Version   string `json:"version" binding:"required"`
	ChangeLog string `json:"changeLog"`
}

// QueryToolRequest 工具查询请求（对应 Java 的 QueryToolRequest）
type QueryToolRequest struct {
	Page     int               `form:"page" json:"page"`
	PageSize int               `form:"pageSize" json:"pageSize"`
	Keyword  string            `form:"keyword" json:"keyword"`
	Status   *domain.ToolStatus `form:"status" json:"status"`
	IsOffice *bool             `form:"isOffice" json:"isOffice"`
	ToolName string            `form:"toolName" json:"toolName"`
}

// GetPageOrDefault 获取页码（默认1）
func (q *QueryToolRequest) GetPageOrDefault() int {
	if q.Page <= 0 {
		return 1
	}
	return q.Page
}

// GetPageSizeOrDefault 获取每页大小（默认15）
func (q *QueryToolRequest) GetPageSizeOrDefault() int {
	if q.PageSize <= 0 {
		return 15
	}
	return q.PageSize
}

// ---- Assembler 转换函数 ----

// ToolEntityToDTO 工具实体转DTO
func ToolEntityToDTO(entity *domain.ToolEntity) *ToolDTO {
	if entity == nil {
		return nil
	}
	installCommandJSON, _ := json.Marshal(entity.InstallCommand)
	return &ToolDTO{
		ID:               entity.ID,
		Name:             entity.Name,
		Icon:             entity.Icon,
		Subtitle:         entity.Subtitle,
		Description:      entity.Description,
		UserID:           entity.UserID,
		Labels:           entity.Labels,
		ToolType:         entity.ToolType,
		UploadType:       entity.UploadType,
		UploadUrl:        entity.UploadUrl,
		ToolList:         entity.ToolList,
		Status:           entity.Status,
		IsOffice:         entity.IsOffice,
		InstallCommand:   string(installCommandJSON),
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
		RejectReason:     entity.RejectReason,
		FailedStepStatus: entity.FailedStepStatus,
		McpServerName:    entity.McpServerName,
		IsGlobal:         entity.IsGlobal,
	}
}

// ToolEntitiesToDTOs 工具实体列表转DTO列表
func ToolEntitiesToDTOs(entities []domain.ToolEntity) []*ToolDTO {
	if len(entities) == 0 {
		return []*ToolDTO{}
	}
	dtos := make([]*ToolDTO, len(entities))
	for i, e := range entities {
		dtos[i] = ToolEntityToDTO(&e)
	}
	return dtos
}

// ToolVersionEntityToDTO 工具版本实体转DTO
func ToolVersionEntityToDTO(entity *domain.ToolVersionEntity) *ToolVersionDTO {
	if entity == nil {
		return nil
	}
	return &ToolVersionDTO{
		ID:            entity.ID,
		Name:          entity.Name,
		Icon:          entity.Icon,
		Subtitle:      entity.Subtitle,
		Description:   entity.Description,
		UserID:        entity.UserID,
		Version:       entity.Version,
		ToolID:        entity.ToolID,
		UploadType:    string(entity.UploadType),
		UploadUrl:     entity.UploadUrl,
		ToolList:      entity.ToolList,
		Labels:        entity.Labels,
		IsOffice:      entity.IsOffice,
		PublicStatus:  entity.PublicStatus,
		ChangeLog:     entity.ChangeLog,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		McpServerName: entity.McpServerName,
	}
}

// UserToolEntityToDTO 用户工具实体转DTO
func UserToolEntityToDTO(entity *domain.UserToolEntity) *ToolVersionDTO {
	if entity == nil {
		return nil
	}
	return &ToolVersionDTO{
		ID:            entity.ID,
		Name:          entity.Name,
		Icon:          entity.Icon,
		Subtitle:      entity.Subtitle,
		Description:   entity.Description,
		UserID:        entity.UserID,
		Version:       entity.Version,
		ToolID:        entity.ToolID,
		ToolList:      entity.ToolList,
		Labels:        entity.Labels,
		IsOffice:      entity.IsOffice,
		PublicStatus:  entity.PublicState,
		McpServerName: entity.McpServerName,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
}

// CreateToolRequestToEntity 创建请求转实体
func CreateToolRequestToEntity(req *CreateToolRequest, userID string) *domain.ToolEntity {
	return &domain.ToolEntity{
		Name:           req.Name,
		Icon:           req.Icon,
		Subtitle:       req.Subtitle,
		Description:    req.Description,
		UserID:         userID,
		Labels:         domain.JSONStringList(req.Labels),
		UploadUrl:      req.UploadUrl,
		InstallCommand: domain.JSONMap(req.InstallCommand),
		ToolType:       domain.ToolTypeMCP,
		UploadType:     domain.UploadTypeGitHub,
		IsGlobal:       req.IsGlobal,
	}
}

// UpdateToolRequestToEntity 更新请求转实体
func UpdateToolRequestToEntity(req *UpdateToolRequest, userID string) *domain.ToolEntity {
	return &domain.ToolEntity{
		Name:           req.Name,
		Icon:           req.Icon,
		Subtitle:       req.Subtitle,
		Description:    req.Description,
		UserID:         userID,
		Labels:         domain.JSONStringList(req.Labels),
		UploadUrl:      req.UploadUrl,
		InstallCommand: domain.JSONMap(req.InstallCommand),
		IsGlobal:       req.IsGlobal,
	}
}
