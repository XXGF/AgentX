package billing

import (
	"fmt"
	"math"

	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// ---- 常量定义（对应 Java 的 UsageDataKeys / PricingConfigKeys）----

const (
	// UsageDataKeys 用量数据键
	UsageDataKeyInputTokens  = "input_tokens"
	UsageDataKeyOutputTokens = "output_tokens"
	UsageDataKeyQuantity     = "quantity"

	// PricingConfigKeys 价格配置键
	PricingConfigKeyInputCostPerMillion  = "input_cost_per_million"
	PricingConfigKeyOutputCostPerMillion = "output_cost_per_million"
	PricingConfigKeyCostPerUnit          = "cost_per_unit"
)

// ---- RuleStrategy 接口（对应 Java 的 RuleStrategy）----

// RuleStrategy 规则策略接口
type RuleStrategy interface {
	// Process 计算费用
	// usageData: 用量数据 (如：{"input_tokens": 1000, "output_tokens": 500}、{"quantity": 1})
	// pricingConfig: 价格配置 (如：{"input_cost_per_million": 5.0, "output_cost_per_million": 15.0})
	Process(usageData map[string]interface{}, pricingConfig map[string]interface{}) (float64, error)

	// GetStrategyName 获取策略名称
	GetStrategyName() string

	// ValidateUsageData 验证用量数据是否有效
	ValidateUsageData(usageData map[string]interface{}) bool

	// ValidatePricingConfig 验证价格配置是否有效
	ValidatePricingConfig(pricingConfig map[string]interface{}) bool
}

// ---- ModelTokenStrategy 模型Token计费策略（对应 Java 的 ModelTokenStrategy）----

// ModelTokenStrategy 基于输入输出Token数量分别计费
type ModelTokenStrategy struct{}

func NewModelTokenStrategy() *ModelTokenStrategy {
	return &ModelTokenStrategy{}
}

func (s *ModelTokenStrategy) Process(usageData map[string]interface{}, pricingConfig map[string]interface{}) (float64, error) {
	if !s.ValidateUsageData(usageData) || !s.ValidatePricingConfig(pricingConfig) {
		return 0, exception.NewBusinessException("无效的用量数据或价格配置")
	}

	// 获取Token数量
	inputTokens := toFloat64(usageData[UsageDataKeyInputTokens])
	outputTokens := toFloat64(usageData[UsageDataKeyOutputTokens])

	// 获取价格配置
	inputCostPerMillion := toFloat64(pricingConfig[PricingConfigKeyInputCostPerMillion])
	outputCostPerMillion := toFloat64(pricingConfig[PricingConfigKeyOutputCostPerMillion])

	// 计算输入Token费用：(inputTokens / 1000000) * inputCostPerMillion
	inputCost := (inputTokens / 1000000.0) * inputCostPerMillion

	// 计算输出Token费用：(outputTokens / 1000000) * outputCostPerMillion
	outputCost := (outputTokens / 1000000.0) * outputCostPerMillion

	// 总费用，保留8位小数
	total := inputCost + outputCost
	return math.Round(total*1e8) / 1e8, nil
}

func (s *ModelTokenStrategy) GetStrategyName() string {
	return "MODEL_TOKEN_STRATEGY"
}

func (s *ModelTokenStrategy) ValidateUsageData(usageData map[string]interface{}) bool {
	if usageData == nil || len(usageData) == 0 {
		return false
	}
	inputTokens, ok1 := usageData[UsageDataKeyInputTokens]
	outputTokens, ok2 := usageData[UsageDataKeyOutputTokens]
	if !ok1 || !ok2 {
		return false
	}
	return toFloat64(inputTokens) >= 0 && toFloat64(outputTokens) >= 0
}

func (s *ModelTokenStrategy) ValidatePricingConfig(pricingConfig map[string]interface{}) bool {
	if pricingConfig == nil || len(pricingConfig) == 0 {
		return false
	}
	inputCost, ok1 := pricingConfig[PricingConfigKeyInputCostPerMillion]
	outputCost, ok2 := pricingConfig[PricingConfigKeyOutputCostPerMillion]
	if !ok1 || !ok2 {
		return false
	}
	return toFloat64(inputCost) >= 0 && toFloat64(outputCost) >= 0
}

// ---- PerUnitStrategy 按次计费策略（对应 Java 的 PerUnitStrategy）----

// PerUnitStrategy 按使用次数进行固定计费
type PerUnitStrategy struct{}

func NewPerUnitStrategy() *PerUnitStrategy {
	return &PerUnitStrategy{}
}

func (s *PerUnitStrategy) Process(usageData map[string]interface{}, pricingConfig map[string]interface{}) (float64, error) {
	if !s.ValidateUsageData(usageData) || !s.ValidatePricingConfig(pricingConfig) {
		return 0, exception.NewBusinessException("无效的用量数据或价格配置")
	}

	// 获取使用数量
	quantity := toFloat64(usageData[UsageDataKeyQuantity])

	// 获取单价
	costPerUnit := toFloat64(pricingConfig[PricingConfigKeyCostPerUnit])

	// 计算总费用：quantity * costPerUnit
	total := quantity * costPerUnit
	return math.Round(total*1e8) / 1e8, nil
}

func (s *PerUnitStrategy) GetStrategyName() string {
	return "PER_UNIT_STRATEGY"
}

func (s *PerUnitStrategy) ValidateUsageData(usageData map[string]interface{}) bool {
	if usageData == nil || len(usageData) == 0 {
		return false
	}
	quantity, ok := usageData[UsageDataKeyQuantity]
	if !ok {
		return false
	}
	return toFloat64(quantity) > 0
}

func (s *PerUnitStrategy) ValidatePricingConfig(pricingConfig map[string]interface{}) bool {
	if pricingConfig == nil || len(pricingConfig) == 0 {
		return false
	}
	costPerUnit, ok := pricingConfig[PricingConfigKeyCostPerUnit]
	if !ok {
		return false
	}
	return toFloat64(costPerUnit) >= 0
}

// ---- BillingStrategyFactory 计费策略工厂（对应 Java 的 BillingStrategyFactory）----

// BillingStrategyFactory 管理所有计费策略实例
type BillingStrategyFactory struct {
	strategyMap map[string]RuleStrategy
}

// NewBillingStrategyFactory 创建计费策略工厂，自动注册所有内置策略
func NewBillingStrategyFactory() *BillingStrategyFactory {
	factory := &BillingStrategyFactory{
		strategyMap: make(map[string]RuleStrategy),
	}

	// 注册内置策略
	modelTokenStrategy := NewModelTokenStrategy()
	perUnitStrategy := NewPerUnitStrategy()

	factory.strategyMap[modelTokenStrategy.GetStrategyName()] = modelTokenStrategy
	factory.strategyMap[perUnitStrategy.GetStrategyName()] = perUnitStrategy

	return factory
}

// GetStrategy 根据handler_key获取对应的策略实例
func (f *BillingStrategyFactory) GetStrategy(handlerKey string) (RuleStrategy, error) {
	strategy, ok := f.strategyMap[handlerKey]
	if !ok {
		return nil, exception.NewBusinessException(fmt.Sprintf("未找到对应的计费策略: %s", handlerKey))
	}
	return strategy, nil
}

// HasStrategy 检查策略是否存在
func (f *BillingStrategyFactory) HasStrategy(handlerKey string) bool {
	_, ok := f.strategyMap[handlerKey]
	return ok
}

// GetAllStrategyNames 获取所有已注册的策略名称
func (f *BillingStrategyFactory) GetAllStrategyNames() []string {
	names := make([]string, 0, len(f.strategyMap))
	for name := range f.strategyMap {
		names = append(names, name)
	}
	return names
}

// RegisterStrategy 注册自定义策略
func (f *BillingStrategyFactory) RegisterStrategy(strategy RuleStrategy) {
	f.strategyMap[strategy.GetStrategyName()] = strategy
}

// ---- 工具函数 ----

// toFloat64 将 interface{} 转换为 float64
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		return 0
	default:
		return 0
	}
}
