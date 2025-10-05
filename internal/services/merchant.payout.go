package services

import (
	"inpayos/internal/protocol"
	"sync"
)

type MerchantPayoutService struct {
}

var (
	merchantPayoutService     *MerchantPayoutService
	merchantPayoutServiceOnce sync.Once
)

func SetupMerchantPayoutService() {
	merchantPayoutServiceOnce.Do(func() {
		merchantPayoutService = &MerchantPayoutService{}
	})
}

// GetMerchantPayoutService 获取Payout服务单例
func GetMerchantPayoutService() *MerchantPayoutService {
	if merchantPayoutService == nil {
		SetupMerchantPayoutService()
	}
	return merchantPayoutService
}

func (s *MerchantPayoutService) Create() (trx *protocol.Transaction, code protocol.ErrorCode) {
	return &protocol.Transaction{}, protocol.Success
}
