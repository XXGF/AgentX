package tool

// ToolType 工具类型枚举（对应 Java 的 ToolType）
type ToolType string

const (
	ToolTypeMCP ToolType = "MCP"
)

// UploadType 上传方式枚举（对应 Java 的 UploadType）
type UploadType string

const (
	UploadTypeGitHub UploadType = "GITHUB"
	UploadTypeZip    UploadType = "ZIP"
)

// ToolStatus 工具审核状态枚举（对应 Java 的 ToolStatus）
type ToolStatus string

const (
	ToolStatusWaitingReview    ToolStatus = "WAITING_REVIEW"     // 等待审核
	ToolStatusGithubURLValidate ToolStatus = "GITHUB_URL_VALIDATE" // GitHub URL 验证中
	ToolStatusDeploying        ToolStatus = "DEPLOYING"          // 部署中
	ToolStatusFetchingTools    ToolStatus = "FETCHING_TOOLS"     // 获取工具中
	ToolStatusManualReview     ToolStatus = "MANUAL_REVIEW"      // 人工审核
	ToolStatusApproved         ToolStatus = "APPROVED"           // 已通过
	ToolStatusFailed           ToolStatus = "FAILED"             // 通用失败状态
)
