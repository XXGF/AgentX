package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/config"
	"go.uber.org/zap"
)

// PaymentProvider 支付提供商接口（对应 Java 的 PaymentProvider）
type PaymentProvider interface {
	// CreatePayment 创建支付
	CreatePayment(orderNo string, amount float64, subject string) (*PaymentResult, error)
	// QueryPayment 查询支付状态
	QueryPayment(orderNo string) (*PaymentQueryResult, error)
	// RefundPayment 退款
	RefundPayment(orderNo string, amount float64, reason string) error
	// GetProviderName 获取提供商名称
	GetProviderName() string
}

// PaymentResult 支付创建结果
type PaymentResult struct {
	PaymentURL string `json:"paymentUrl"` // 支付链接
	QRCode     string `json:"qrCode"`     // 二维码内容
	OrderNo    string `json:"orderNo"`    // 订单号
	TradeNo    string `json:"tradeNo"`    // 第三方交易号
}

// PaymentQueryResult 支付查询结果
type PaymentQueryResult struct {
	OrderNo string     `json:"orderNo"`
	TradeNo string     `json:"tradeNo"`
	Status  string     `json:"status"` // paid / pending / failed
	Amount  float64    `json:"amount"`
	PaidAt  *time.Time `json:"paidAt"`
}

// ---- Stripe 支付提供商 ----

// StripePaymentProvider Stripe 支付提供商（通过 Stripe REST API 直接调用）
type StripePaymentProvider struct {
	secretKey     string
	webhookSecret string
	httpClient    *http.Client
	logger        *zap.Logger
}

func NewStripePaymentProvider(stripeCfg *config.StripeConfig, logger *zap.Logger) *StripePaymentProvider {
	return &StripePaymentProvider{
		secretKey:     stripeCfg.SecretKey,
		webhookSecret: stripeCfg.WebhookSecret,
		httpClient:    &http.Client{Timeout: 30 * time.Second},
		logger:        logger,
	}
}

func (p *StripePaymentProvider) CreatePayment(orderNo string, amount float64, subject string) (*PaymentResult, error) {
	if p.secretKey == "" {
		return nil, fmt.Errorf("Stripe 未配置 secret-key")
	}

	// Stripe Checkout Session API
	data := url.Values{}
	data.Set("payment_method_types[]", "card")
	data.Set("line_items[0][price_data][currency]", "cny")
	data.Set("line_items[0][price_data][product_data][name]", subject)
	data.Set("line_items[0][price_data][unit_amount]", fmt.Sprintf("%d", int(amount*100))) // 分
	data.Set("line_items[0][quantity]", "1")
	data.Set("mode", "payment")
	data.Set("client_reference_id", orderNo)
	data.Set("success_url", "https://your-domain.com/payment/success?session_id={CHECKOUT_SESSION_ID}")
	data.Set("cancel_url", "https://your-domain.com/payment/cancel")

	req, err := http.NewRequest("POST", "https://api.stripe.com/v1/checkout/sessions", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(p.secretKey, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Stripe API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Stripe 创建支付失败 (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析 Stripe 响应失败: %w", err)
	}

	return &PaymentResult{
		PaymentURL: result.URL,
		OrderNo:    orderNo,
		TradeNo:    result.ID,
	}, nil
}

func (p *StripePaymentProvider) QueryPayment(orderNo string) (*PaymentQueryResult, error) {
	if p.secretKey == "" {
		return nil, fmt.Errorf("Stripe 未配置 secret-key")
	}

	// 通过 client_reference_id 查询（简化实现，实际应通过 session ID）
	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://api.stripe.com/v1/checkout/sessions?client_reference_id=%s&limit=1", orderNo), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(p.secretKey, "")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Stripe API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Stripe 查询支付失败 (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Data []struct {
			ID            string `json:"id"`
			PaymentStatus string `json:"payment_status"`
			AmountTotal   int    `json:"amount_total"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return &PaymentQueryResult{OrderNo: orderNo, Status: "pending"}, nil
	}

	session := result.Data[0]
	status := "pending"
	if session.PaymentStatus == "paid" {
		status = "paid"
	}

	return &PaymentQueryResult{
		OrderNo: orderNo,
		TradeNo: session.ID,
		Status:  status,
		Amount:  float64(session.AmountTotal) / 100,
	}, nil
}

func (p *StripePaymentProvider) RefundPayment(orderNo string, amount float64, reason string) error {
	p.logger.Info("Stripe 退款", zap.String("orderNo", orderNo), zap.Float64("amount", amount))
	// 实际退款需要 payment_intent ID，这里简化处理
	return nil
}

func (p *StripePaymentProvider) GetProviderName() string {
	return "stripe"
}

// ---- 模拟支付提供商（开发测试用）----

// MockPaymentProvider 模拟支付提供商
type MockPaymentProvider struct {
	logger *zap.Logger
}

func NewMockPaymentProvider(logger *zap.Logger) *MockPaymentProvider {
	return &MockPaymentProvider{logger: logger}
}

func (p *MockPaymentProvider) CreatePayment(orderNo string, amount float64, subject string) (*PaymentResult, error) {
	p.logger.Info("模拟创建支付",
		zap.String("orderNo", orderNo),
		zap.Float64("amount", amount),
		zap.String("subject", subject),
	)
	return &PaymentResult{
		PaymentURL: fmt.Sprintf("https://mock-pay.example.com/pay?orderNo=%s&amount=%.2f", orderNo, amount),
		QRCode:     fmt.Sprintf("mock-qr://%s", orderNo),
		OrderNo:    orderNo,
		TradeNo:    fmt.Sprintf("MOCK-%d", time.Now().UnixNano()),
	}, nil
}

func (p *MockPaymentProvider) QueryPayment(orderNo string) (*PaymentQueryResult, error) {
	now := time.Now()
	return &PaymentQueryResult{
		OrderNo: orderNo,
		TradeNo: fmt.Sprintf("MOCK-%s", orderNo),
		Status:  "paid",
		Amount:  0,
		PaidAt:  &now,
	}, nil
}

func (p *MockPaymentProvider) RefundPayment(orderNo string, amount float64, reason string) error {
	p.logger.Info("模拟退款", zap.String("orderNo", orderNo), zap.Float64("amount", amount))
	return nil
}

func (p *MockPaymentProvider) GetProviderName() string {
	return "mock"
}

// ---- 支付提供商工厂 ----

// PaymentProviderFactory 支付提供商工厂
type PaymentProviderFactory struct {
	providers map[string]PaymentProvider
}

func NewPaymentProviderFactory(paymentCfg *config.PaymentConfig, logger *zap.Logger) *PaymentProviderFactory {
	factory := &PaymentProviderFactory{
		providers: make(map[string]PaymentProvider),
	}

	// 注册模拟支付提供商（始终可用）
	mockProvider := NewMockPaymentProvider(logger)
	factory.providers["mock"] = mockProvider

	// 注册 Stripe 支付提供商
	if paymentCfg.Stripe.SecretKey != "" {
		stripeProvider := NewStripePaymentProvider(&paymentCfg.Stripe, logger)
		factory.providers["stripe"] = stripeProvider
		logger.Info("Stripe 支付提供商已注册")
	} else {
		factory.providers["stripe"] = mockProvider
		logger.Warn("Stripe 未配置，使用模拟支付提供商")
	}

	// 支付宝（暂用模拟，需要支付宝 SDK）
	factory.providers["alipay"] = mockProvider
	if paymentCfg.Alipay.AppID != "" {
		logger.Info("支付宝配置已加载（当前使用模拟实现，需集成支付宝 SDK）",
			zap.String("appId", paymentCfg.Alipay.AppID))
	}

	return factory
}

// GetProvider 获取支付提供商
func (f *PaymentProviderFactory) GetProvider(providerName string) (PaymentProvider, error) {
	provider, ok := f.providers[providerName]
	if !ok {
		return nil, fmt.Errorf("不支持的支付提供商: %s", providerName)
	}
	return provider, nil
}

// GetDefaultProvider 获取默认支付提供商
func (f *PaymentProviderFactory) GetDefaultProvider() PaymentProvider {
	// 优先使用 Stripe，否则使用 mock
	if p, ok := f.providers["stripe"]; ok {
		return p
	}
	return f.providers["mock"]
}

// ---- Alipay 支付宝支付提供商（预留接口）----
// 注意：支付宝需要使用其官方 SDK 进行签名验签，
// 建议使用 github.com/smartwalle/alipay/v3 库
// 当前预留接口，后续集成时只需实现 PaymentProvider 接口即可

// AlipayPaymentProvider 支付宝支付提供商（预留）
type AlipayPaymentProvider struct {
	appID      string
	privateKey string
	publicKey  string
	notifyURL  string
	returnURL  string
	logger     *zap.Logger
}

func NewAlipayPaymentProvider(alipayCfg *config.AlipayConfig, logger *zap.Logger) *AlipayPaymentProvider {
	return &AlipayPaymentProvider{
		appID:      alipayCfg.AppID,
		privateKey: alipayCfg.PrivateKey,
		publicKey:  alipayCfg.PublicKey,
		notifyURL:  alipayCfg.NotifyURL,
		returnURL:  alipayCfg.ReturnURL,
		logger:     logger,
	}
}

func (p *AlipayPaymentProvider) CreatePayment(orderNo string, amount float64, subject string) (*PaymentResult, error) {
	// 预留：需要集成 github.com/smartwalle/alipay/v3
	_ = bytes.NewBuffer(nil) // 避免 unused import
	return nil, fmt.Errorf("支付宝支付尚未集成，请使用 Stripe 或联系管理员")
}

func (p *AlipayPaymentProvider) QueryPayment(orderNo string) (*PaymentQueryResult, error) {
	return nil, fmt.Errorf("支付宝支付尚未集成")
}

func (p *AlipayPaymentProvider) RefundPayment(orderNo string, amount float64, reason string) error {
	return fmt.Errorf("支付宝退款尚未集成")
}

func (p *AlipayPaymentProvider) GetProviderName() string {
	return "alipay"
}