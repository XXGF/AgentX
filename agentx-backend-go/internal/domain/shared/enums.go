package shared

// TokenOverflowStrategy Token超限处理策略枚举
// 注意：此枚举与 domain/token 包中的 TokenOverflowStrategy 保持一致
// 在 shared 包中提供是为了跨模块引用
type TokenOverflowStrategy string

const (
	// TokenOverflowStrategyNone 无策略 - 不做任何处理，可能导致超限错误
	TokenOverflowStrategyNone TokenOverflowStrategy = "NONE"
	// TokenOverflowStrategySlidingWindow 滑动窗口 - 自动移除旧消息，保留最新内容
	TokenOverflowStrategySlidingWindow TokenOverflowStrategy = "SLIDING_WINDOW"
	// TokenOverflowStrategySummarize 摘要策略 - 将旧消息转换为摘要，保留关键信息
	TokenOverflowStrategySummarize TokenOverflowStrategy = "SUMMARIZE"
)

// IsValid 判断给定字符串是否为有效的枚举值
func (s TokenOverflowStrategy) IsValid() bool {
	switch s {
	case TokenOverflowStrategyNone, TokenOverflowStrategySlidingWindow, TokenOverflowStrategySummarize:
		return true
	default:
		return false
	}
}

// TokenOverflowStrategyFromString 从字符串转换为枚举值，如果不存在则返回默认值NONE
func TokenOverflowStrategyFromString(value string) TokenOverflowStrategy {
	s := TokenOverflowStrategy(value)
	if s.IsValid() {
		return s
	}
	return TokenOverflowStrategyNone
}
