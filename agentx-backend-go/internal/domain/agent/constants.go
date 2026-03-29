package agent

// AgentStatus Agent状态枚举（对应 Java 的 AgentStatus）
type AgentStatus int

const (
	AgentStatusDraft         AgentStatus = 0 // 草稿
	AgentStatusPendingReview AgentStatus = 1 // 待审核
	AgentStatusPublished     AgentStatus = 2 // 已上架
	AgentStatusUnpublished   AgentStatus = 3 // 已下架
	AgentStatusRejected      AgentStatus = 4 // 审核拒绝
)

// PublishStatus 版本发布状态枚举（对应 Java 的 PublishStatus）
type PublishStatus int

const (
	PublishStatusReviewing PublishStatus = 1 // 审核中
	PublishStatusPublished PublishStatus = 2 // 已发布
	PublishStatusRejected  PublishStatus = 3 // 拒绝
	PublishStatusRemoved   PublishStatus = 4 // 已下架
)

// PublishStatusDescription 获取发布状态描述
func PublishStatusDescription(status PublishStatus) string {
	switch status {
	case PublishStatusReviewing:
		return "审核中"
	case PublishStatusPublished:
		return "已发布"
	case PublishStatusRejected:
		return "拒绝"
	case PublishStatusRemoved:
		return "已下架"
	default:
		return "未知"
	}
}

// TokenOverflowStrategy Token溢出策略枚举（对应 Java 的 TokenOverflowStrategyEnum）
type TokenOverflowStrategy string

const (
	TokenOverflowStrategyNone          TokenOverflowStrategy = "NONE"           // 不处理
	TokenOverflowStrategySlidingWindow TokenOverflowStrategy = "SLIDING_WINDOW" // 滑动窗口
	TokenOverflowStrategySummary       TokenOverflowStrategy = "SUMMARY"        // 摘要
)

// LLMModelConfig Agent模型配置（对应 Java 的 LLMModelConfig）
type LLMModelConfig struct {
	ModelID          string                `json:"modelId,omitempty"`
	Temperature      *float64              `json:"temperature,omitempty"`
	TopP             *float64              `json:"topP,omitempty"`
	TopK             *int                  `json:"topK,omitempty"`
	MaxTokens        *int                  `json:"maxTokens,omitempty"`
	StrategyType     TokenOverflowStrategy `json:"strategyType,omitempty"`
	ReserveRatio     *float64              `json:"reserveRatio,omitempty"`
	SummaryThreshold *int                  `json:"summaryThreshold,omitempty"`
}

// NewDefaultLLMModelConfig 创建默认模型配置
func NewDefaultLLMModelConfig() *LLMModelConfig {
	temp := 0.7
	topP := 0.7
	topK := 50
	return &LLMModelConfig{
		Temperature:  &temp,
		TopP:         &topP,
		TopK:         &topK,
		StrategyType: TokenOverflowStrategyNone,
	}
}
