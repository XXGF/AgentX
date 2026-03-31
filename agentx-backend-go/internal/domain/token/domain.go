package token

import (
	"time"
)

// TokenOverflowStrategy Token超限策略枚举
type TokenOverflowStrategy string

const (
	TokenOverflowStrategyNone          TokenOverflowStrategy = "NONE"           // 不处理
	TokenOverflowStrategySlidingWindow TokenOverflowStrategy = "SLIDING_WINDOW" // 滑动窗口
	TokenOverflowStrategySummarize     TokenOverflowStrategy = "SUMMARIZE"      // 摘要
)

// TokenMessage Token领域的消息模型
type TokenMessage struct {
	ID             string    `json:"id"`
	Content        string    `json:"content"`
	Role           string    `json:"role"`
	TokenCount     *int      `json:"tokenCount"`
	BodyTokenCount *int      `json:"bodyTokenCount"`
	CreatedAt      time.Time `json:"createdAt"`
}

func NewTokenMessage(id, content, role string, tokenCount *int) *TokenMessage {
	return &TokenMessage{
		ID:         id,
		Content:    content,
		Role:       role,
		TokenCount: tokenCount,
		CreatedAt:  time.Now(),
	}
}

// TokenOverflowConfig Token超限处理配置
type TokenOverflowConfig struct {
	StrategyType     TokenOverflowStrategy `json:"strategyType"`
	MaxTokens        *int                  `json:"maxTokens"`
	ReserveRatio     *float64              `json:"reserveRatio"`
	SummaryThreshold *int                  `json:"summaryThreshold"`
}

// CreateDefault 创建默认配置（不处理）
func CreateDefaultConfig() *TokenOverflowConfig {
	return &TokenOverflowConfig{StrategyType: TokenOverflowStrategyNone}
}

// CreateSlidingWindowConfig 创建滑动窗口策略配置
func CreateSlidingWindowConfig(maxTokens int, reserveRatio *float64) *TokenOverflowConfig {
	ratio := 0.1
	if reserveRatio != nil {
		ratio = *reserveRatio
	}
	return &TokenOverflowConfig{
		StrategyType: TokenOverflowStrategySlidingWindow,
		MaxTokens:    &maxTokens,
		ReserveRatio: &ratio,
	}
}

// CreateSummaryConfig 创建摘要策略配置
func CreateSummaryConfig(maxTokens int, summaryThreshold *int) *TokenOverflowConfig {
	threshold := 20
	if summaryThreshold != nil {
		threshold = *summaryThreshold
	}
	return &TokenOverflowConfig{
		StrategyType:     TokenOverflowStrategySummarize,
		MaxTokens:        &maxTokens,
		SummaryThreshold: &threshold,
	}
}

// TokenProcessResult Token处理结果
type TokenProcessResult struct {
	RetainedMessages []*TokenMessage `json:"retainedMessages"`
	Summary          string          `json:"summary"`
	TotalTokens      int             `json:"totalTokens"`
	StrategyName     string          `json:"strategyName"`
	Processed        bool            `json:"processed"`
}

// TokenResult Token处理结果（简化版）
type TokenResult struct {
	RetainedMessages []*TokenMessage `json:"retainedMessages"`
	Summary          string          `json:"summary"`
	StrategyName     string          `json:"strategyName"`
	TotalTokens      int             `json:"totalTokens"`
}

// GetRetainedMessageIDs 获取保留消息的ID列表
func (r *TokenResult) GetRetainedMessageIDs() []string {
	var ids []string
	for _, msg := range r.RetainedMessages {
		ids = append(ids, msg.ID)
	}
	return ids
}

// ---- Token处理策略接口 ----

type TokenOverflowProcessor interface {
	Process(messages []*TokenMessage, config *TokenOverflowConfig) *TokenProcessResult
}

// NoTokenOverflowProcessor 不处理策略
type NoTokenOverflowProcessor struct{}

func (p *NoTokenOverflowProcessor) Process(messages []*TokenMessage, config *TokenOverflowConfig) *TokenProcessResult {
	total := calculateTotalTokens(messages)
	return &TokenProcessResult{
		RetainedMessages: messages,
		TotalTokens:      total,
		StrategyName:     "NONE",
		Processed:        false,
	}
}

// SlidingWindowProcessor 滑动窗口策略
type SlidingWindowProcessor struct{}

func (p *SlidingWindowProcessor) Process(messages []*TokenMessage, config *TokenOverflowConfig) *TokenProcessResult {
	total := calculateTotalTokens(messages)
	maxTokens := 0
	if config.MaxTokens != nil {
		maxTokens = *config.MaxTokens
	}
	reserveRatio := 0.1
	if config.ReserveRatio != nil {
		reserveRatio = *config.ReserveRatio
	}

	effectiveMax := int(float64(maxTokens) * (1 - reserveRatio))
	if total <= effectiveMax {
		return &TokenProcessResult{
			RetainedMessages: messages,
			TotalTokens:      total,
			StrategyName:     "SLIDING_WINDOW",
			Processed:        false,
		}
	}

	// 从最新消息开始保留，直到超过限制
	var retained []*TokenMessage
	currentTokens := 0
	for i := len(messages) - 1; i >= 0; i-- {
		msgTokens := 0
		if messages[i].TokenCount != nil {
			msgTokens = *messages[i].TokenCount
		}
		if currentTokens+msgTokens > effectiveMax && len(retained) > 0 {
			break
		}
		retained = append([]*TokenMessage{messages[i]}, retained...)
		currentTokens += msgTokens
	}

	return &TokenProcessResult{
		RetainedMessages: retained,
		TotalTokens:      currentTokens,
		StrategyName:     "SLIDING_WINDOW",
		Processed:        true,
	}
}

// ---- 领域服务 ----

type DomainService struct{}

func NewDomainService() *DomainService {
	return &DomainService{}
}

// ProcessMessages 处理消息列表
func (s *DomainService) ProcessMessages(messages []*TokenMessage, config *TokenOverflowConfig) *TokenProcessResult {
	processor := createProcessor(config)
	return processor.Process(messages, config)
}

func createProcessor(config *TokenOverflowConfig) TokenOverflowProcessor {
	switch config.StrategyType {
	case TokenOverflowStrategySlidingWindow:
		return &SlidingWindowProcessor{}
	case TokenOverflowStrategySummarize:
		// TODO: 摘要策略需要集成LLM服务
		return &NoTokenOverflowProcessor{}
	default:
		return &NoTokenOverflowProcessor{}
	}
}

func calculateTotalTokens(messages []*TokenMessage) int {
	total := 0
	for _, msg := range messages {
		if msg.TokenCount != nil {
			total += *msg.TokenCount
		}
	}
	return total
}
