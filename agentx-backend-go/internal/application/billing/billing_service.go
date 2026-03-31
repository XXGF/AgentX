package billing

import (
	"github.com/google/uuid"
	infraBilling "github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/billing"
	domainProduct "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/product"
	domainRule "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/rule"
	domainUser "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/user"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	"time"
)

// RuleContext 计费上下文（对应 Java 的 RuleContext）
type RuleContext struct {
	Type      string                 `json:"type"`      // 规则类型（BillingType）
	ServiceID string                 `json:"serviceId"` // 服务ID
	UsageData map[string]interface{} `json:"usageData"` // 用量数据
	RequestID string                 `json:"requestId"` // 请求ID（幂等性）
	UserID    string                 `json:"userId"`    // 用户ID
}

// IsValid 验证上下文是否有效
func (c *RuleContext) IsValid() bool {
	return c.Type != "" && c.ServiceID != "" &&
		c.UsageData != nil && len(c.UsageData) > 0 && c.UserID != ""
}

// BillingService 计费服务（对应 Java 的 BillingService，完整实现）
type BillingService struct {
	productDomainService     *domainProduct.DomainService
	ruleDomainService        *domainRule.DomainService
	accountDomainService     *domainUser.AccountDomainService
	usageRecordDomainService *domainUser.UsageRecordDomainService
	billingStrategyFactory   *infraBilling.BillingStrategyFactory
}

func NewBillingService(
	productDomainService *domainProduct.DomainService,
	ruleDomainService *domainRule.DomainService,
	accountDomainService *domainUser.AccountDomainService,
	usageRecordDomainService *domainUser.UsageRecordDomainService,
) *BillingService {
	return &BillingService{
		productDomainService:     productDomainService,
		ruleDomainService:        ruleDomainService,
		accountDomainService:     accountDomainService,
		usageRecordDomainService: usageRecordDomainService,
		billingStrategyFactory:   infraBilling.NewBillingStrategyFactory(),
	}
}

// Charge 执行计费（对应 Java 的 BillingService.charge）
func (s *BillingService) Charge(ctx *RuleContext) error {
	// 1. 验证上下文
	if !ctx.IsValid() {
		return exception.NewBusinessException("无效的计费上下文")
	}

	// 2. 查找商品
	product, err := s.productDomainService.FindProductByBusinessKey(domainProduct.BillingType(ctx.Type), ctx.ServiceID)
	if err != nil {
		return err
	}
	if product == nil {
		return nil // 没有配置计费规则，直接放行
	}
	if !product.IsActive() {
		return exception.NewBusinessException("商品已被禁用，无法计费")
	}

	// 3. 检查幂等性
	if ctx.RequestID != "" {
		exists, err := s.usageRecordDomainService.ExistsByRequestID(ctx.RequestID)
		if err != nil {
			return err
		}
		if exists {
			return nil // 请求已处理
		}
	}

	// 4. 获取规则和策略
	rule, err := s.ruleDomainService.GetRuleByID(product.RuleID)
	if err != nil {
		return err
	}
	if rule == nil {
		return exception.NewBusinessException("关联的计费规则不存在")
	}

	strategy, err := s.billingStrategyFactory.GetStrategy(string(rule.HandlerKey))
	if err != nil {
		return err
	}

	// 5. 计算费用
	cost, err := strategy.Process(ctx.UsageData, product.PricingConfig)
	if err != nil {
		return err
	}

	if cost < 0 {
		return exception.NewBusinessException("计算出的费用不能为负数")
	}

	// 实现最低计费0.01元逻辑：如果费用大于0但小于0.01，则按0.01计算
	if cost > 0 && cost < 0.01 {
		cost = 0.01
	}

	// 如果费用为0，也需要记录用量，但不扣费
	if cost == 0 {
		s.recordUsage(ctx, product, cost)
		return nil
	}

	// 6. 检查余额并扣费
	if err := s.accountDomainService.DeductBalance(ctx.UserID, cost); err != nil {
		return err
	}

	// 7. 记录用量
	s.recordUsage(ctx, product, cost)

	return nil
}

// CheckBalance 检查余额是否充足（不实际扣费）
func (s *BillingService) CheckBalance(ctx *RuleContext) (bool, error) {
	// 查找商品
	product, err := s.productDomainService.FindProductByBusinessKey(domainProduct.BillingType(ctx.Type), ctx.ServiceID)
	if err != nil {
		return false, err
	}
	if product == nil || !product.IsActive() {
		return true, nil // 无需计费
	}

	// 获取规则和策略
	rule, err := s.ruleDomainService.GetRuleByID(product.RuleID)
	if err != nil {
		return false, err
	}
	if rule == nil {
		return false, nil
	}

	strategy, err := s.billingStrategyFactory.GetStrategy(string(rule.HandlerKey))
	if err != nil {
		return false, nil
	}

	// 计算费用
	cost, err := strategy.Process(ctx.UsageData, product.PricingConfig)
	if err != nil {
		return false, nil
	}

	// 实现最低计费0.01元逻辑
	if cost > 0 && cost < 0.01 {
		cost = 0.01
	}

	if cost <= 0 {
		return true, nil // 无需扣费
	}

	// 检查余额
	return s.accountDomainService.CheckSufficientBalance(ctx.UserID, cost)
}

// recordUsage 记录用量
func (s *BillingService) recordUsage(ctx *RuleContext, product *domainProduct.ProductEntity, cost float64) {
	now := time.Now()
	usageRecord := &domainUser.UsageRecordEntity{
		ID:          uuid.New().String(),
		UserID:      ctx.UserID,
		ProductID:   product.ID,
		QuantityData: ctx.UsageData,
		Cost:        cost,
		RequestID:   ctx.RequestID,
		BilledAt:    &now,
		ServiceName: product.Name,
		ServiceType: string(product.Type),
	}
	// 忽略记录用量的错误，不影响主流程
	_ = s.usageRecordDomainService.CreateUsageRecord(usageRecord)
}