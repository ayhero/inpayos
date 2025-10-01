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
