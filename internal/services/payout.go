package services

import (
	"inpayos/internal/protocol"
	"sync"
)

type PayoutService struct {
}

var (
	payoutService     *PayoutService
	payoutServiceOnce sync.Once
)

func SetupPayoutService() {
	payoutServiceOnce.Do(func() {
		payoutService = &PayoutService{}
	})
}

// GetPayoutService 获取Payout服务单例
func GetPayoutService() *PayoutService {
	if payoutService == nil {
		SetupPayoutService()
	}
	return payoutService
}

func (s *PayoutService) Create() (trx *protocol.Transaction, code protocol.ErrorCode) {
	return &protocol.Transaction{}, protocol.Success
}
