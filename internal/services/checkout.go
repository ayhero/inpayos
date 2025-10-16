package services

import (
	"inpayos/internal/protocol"
	"sync"
)

type CheckoutService struct {
}

var (
	checkoutService     *CheckoutService
	checkoutServiceOnce sync.Once
)

func SetupCheckoutService() {
	checkoutServiceOnce.Do(func() {
		checkoutService = &CheckoutService{}
	})
}

// GetCheckoutService 获取Checkout服务单例
func GetCheckoutService() *CheckoutService {
	if checkoutService == nil {
		SetupCheckoutService()
	}
	return checkoutService
}

func (s *CheckoutService) Create(req *protocol.CreateCheckoutRequest) (trx *protocol.Checkout, code protocol.ErrorCode) {
	return &protocol.Checkout{}, protocol.Success
}
func (s *CheckoutService) Confirm(req *protocol.ConfirmCheckoutRequest) (trx *protocol.Checkout, code protocol.ErrorCode) {
	return &protocol.Checkout{}, protocol.Success
}

func (s *CheckoutService) Info(checkoutID string) (trx *protocol.Checkout, code protocol.ErrorCode) {
	return &protocol.Checkout{}, protocol.Success
}

func (s *CheckoutService) Cancel(checkoutID string) (trx *protocol.Checkout, code protocol.ErrorCode) {
	return &protocol.Checkout{}, protocol.Success
}

func (s *CheckoutService) Submit(req *protocol.SubmitCheckoutRequest) (trx *protocol.Checkout, code protocol.ErrorCode) {
	// TODO: 实现提交收银台支付逻辑
	return &protocol.Checkout{
		CheckoutID: req.CheckoutID,
		Status:     "processing",
	}, protocol.Success
}

func (s *CheckoutService) ListServices(checkoutID string) (services []protocol.CheckoutService, code protocol.ErrorCode) {
	// TODO: 根据checkoutID从数据库获取收银台信息，然后返回该收银台会话可用的支付服务列表
	// 目前返回默认的支付服务列表，实际应该根据checkoutID对应的收银台配置来过滤

	// 这里可以根据checkoutID获取收银台信息，然后基于以下条件过滤服务：
	// 1. 收银台设置的币种
	// 2. 收银台设置的国家/地区
	// 3. 支付金额范围
	// 4. 商户的支付通道配置

	allServices := []protocol.CheckoutService{
		{
			ID:          "alipay",
			Name:        "支付宝",
			Description: "支付宝在线支付",
			LogoURL:     "/images/alipay-logo.png",
			MinAmount:   "0.01",
			MaxAmount:   "50000.00",
			Currencies:  []string{"CNY", "USD"},
			Countries:   []string{"CN", "US"},
			Status:      "active",
		},
		{
			ID:          "wechatpay",
			Name:        "微信支付",
			Description: "微信在线支付",
			LogoURL:     "/images/wechatpay-logo.png",
			MinAmount:   "0.01",
			MaxAmount:   "50000.00",
			Currencies:  []string{"CNY"},
			Countries:   []string{"CN"},
			Status:      "active",
		},
		{
			ID:          "card",
			Name:        "银行卡支付",
			Description: "信用卡/借记卡支付",
			LogoURL:     "/images/card-logo.png",
			MinAmount:   "1.00",
			MaxAmount:   "100000.00",
			Currencies:  []string{"USD", "EUR", "CNY"},
			Countries:   []string{"US", "EU", "CN"},
			Status:      "active",
		},
		{
			ID:          "paypal",
			Name:        "PayPal",
			Description: "PayPal在线支付",
			LogoURL:     "/images/paypal-logo.png",
			MinAmount:   "0.01",
			MaxAmount:   "10000.00",
			Currencies:  []string{"USD", "EUR"},
			Countries:   []string{"US", "EU"},
			Status:      "active",
		},
	}

	// TODO: 实际实现中应该根据checkoutID获取收银台信息，并过滤出适用的服务
	// 这里简单返回所有服务作为示例
	return allServices, protocol.Success
}
