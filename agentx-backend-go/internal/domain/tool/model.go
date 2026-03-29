package tool

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
)

// ToolDefinition 工具定义（对应 Java 的 ToolDefinition）
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Enabled     *bool                  `json:"enabled"`
}

// ToolDefinitionList 工具定义列表（JSONB类型）
type ToolDefinitionList []ToolDefinition

func (t ToolDefinitionList) Value() (driver.Value, error) {
	if t == nil {
		return "[]", nil
	}
	b, err := json.Marshal(t)
	return string(b), err
}

func (t *ToolDefinitionList) Scan(value interface{}) error {
	if value == nil {
		*t = []ToolDefinition{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		*t = []ToolDefinition{}
		return nil
	}
	return json.Unmarshal(bytes, t)
}

// JSONStringList JSON字符串列表类型
type JSONStringList []string

func (s JSONStringList) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	return string(b), err
}

func (s *JSONStringList) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		*s = []string{}
		return nil
	}
	return json.Unmarshal(bytes, s)
}

// JSONMap JSON Map类型（用于 install_command 等字段）
type JSONMap map[string]interface{}

func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return "{}", nil
	}
	b, err := json.Marshal(m)
	return string(b), err
}

func (m *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*m = map[string]interface{}{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		*m = map[string]interface{}{}
		return nil
	}
	return json.Unmarshal(bytes, m)
}

// ToolEntity 工具实体（对应 Java 的 ToolEntity）
type ToolEntity struct {
	ID               string             `gorm:"column:id;primaryKey" json:"id"`
	Name             string             `gorm:"column:name" json:"name"`
	Icon             string             `gorm:"column:icon" json:"icon"`
	Subtitle         string             `gorm:"column:subtitle" json:"subtitle"`
	Description      string             `gorm:"column:description" json:"description"`
	UserID           string             `gorm:"column:user_id" json:"userId"`
	Labels           JSONStringList     `gorm:"column:labels;type:jsonb" json:"labels"`
	ToolType         ToolType           `gorm:"column:tool_type" json:"toolType"`
	UploadType       UploadType         `gorm:"column:upload_type" json:"uploadType"`
	UploadUrl        string             `gorm:"column:upload_url" json:"uploadUrl"`
	InstallCommand   JSONMap            `gorm:"column:install_command;type:jsonb" json:"installCommand"`
	ToolList         ToolDefinitionList `gorm:"column:tool_list;type:jsonb" json:"toolList"`
	Status           ToolStatus         `gorm:"column:status" json:"status"`
	IsOffice         *bool              `gorm:"column:is_office" json:"isOffice"`
	RejectReason     string             `gorm:"column:reject_reason" json:"rejectReason"`
	FailedStepStatus ToolStatus         `gorm:"column:failed_step_status" json:"failedStepStatus"`
	McpServerName    string             `gorm:"column:mcp_server_name" json:"mcpServerName"`
	IsGlobal         *bool              `gorm:"column:is_global" json:"isGlobal"`

	entity.BaseEntity
}

func (ToolEntity) TableName() string {
	return "tools"
}

// IsGlobalTool 是否为全局工具
func (e *ToolEntity) IsGlobalTool() bool {
	return e.IsGlobal != nil && *e.IsGlobal
}

// RequiresUserContainer 是否需要用户容器
func (e *ToolEntity) RequiresUserContainer() bool {
	return !e.IsGlobalTool()
}

// ToolVersionEntity 工具版本实体（对应 Java 的 ToolVersionEntity）
type ToolVersionEntity struct {
	ID            string             `gorm:"column:id;primaryKey" json:"id"`
	Name          string             `gorm:"column:name" json:"name"`
	Icon          string             `gorm:"column:icon" json:"icon"`
	Subtitle      string             `gorm:"column:subtitle" json:"subtitle"`
	Description   string             `gorm:"column:description" json:"description"`
	UserID        string             `gorm:"column:user_id" json:"userId"`
	Version       string             `gorm:"column:version" json:"version"`
	ToolID        string             `gorm:"column:tool_id" json:"toolId"`
	UploadType    UploadType         `gorm:"column:upload_type" json:"uploadType"`
	UploadUrl     string             `gorm:"column:upload_url" json:"uploadUrl"`
	ToolList      ToolDefinitionList `gorm:"column:tool_list;type:jsonb" json:"toolList"`
	Labels        JSONStringList     `gorm:"column:labels;type:jsonb" json:"labels"`
	IsOffice      *bool              `gorm:"column:is_office" json:"isOffice"`
	PublicStatus  *bool              `gorm:"column:public_status" json:"publicStatus"`
	ChangeLog     string             `gorm:"column:change_log" json:"changeLog"`
	McpServerName string             `gorm:"column:mcp_server_name" json:"mcpServerName"`

	entity.BaseEntity
}

func (ToolVersionEntity) TableName() string {
	return "tool_versions"
}

// UserToolEntity 用户工具关联实体（对应 Java 的 UserToolEntity）
type UserToolEntity struct {
	ID            string             `gorm:"column:id;primaryKey" json:"id"`
	UserID        string             `gorm:"column:user_id" json:"userId"`
	Name          string             `gorm:"column:name" json:"name"`
	Description   string             `gorm:"column:description" json:"description"`
	Icon          string             `gorm:"column:icon" json:"icon"`
	Subtitle      string             `gorm:"column:subtitle" json:"subtitle"`
	ToolID        string             `gorm:"column:tool_id" json:"toolId"`
	Version       string             `gorm:"column:version" json:"version"`
	ToolList      ToolDefinitionList `gorm:"column:tool_list;type:jsonb" json:"toolList"`
	Labels        JSONStringList     `gorm:"column:labels;type:jsonb" json:"labels"`
	IsOffice      *bool              `gorm:"column:is_office" json:"isOffice"`
	PublicState   *bool              `gorm:"column:public_state" json:"publicState"`
	McpServerName string             `gorm:"column:mcp_server_name" json:"mcpServerName"`
	IsGlobal      *bool              `gorm:"column:is_global" json:"isGlobal"`

	entity.BaseEntity
}

func (UserToolEntity) TableName() string {
	return "user_tools"
}

// ToolOperationResult 工具操作结果（对应 Java 的 ToolOperationResult）
type ToolOperationResult struct {
	Tool                *ToolEntity
	NeedStateTransition bool
}

func NewToolOperationResult(tool *ToolEntity, needStateTransition bool) *ToolOperationResult {
	return &ToolOperationResult{Tool: tool, NeedStateTransition: needStateTransition}
}
