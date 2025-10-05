package services

import (
	"inpayos/internal/protocol"
	"sync"
)

type MerchantPayinService struct {
}

var (
	merchantPayinService     *MerchantPayinService
	merchantPayinServiceOnce sync.Once
)

func SetupMerchantPayinService() {
	merchantPayinServiceOnce.Do(func() {
		merchantPayinService = &MerchantPayinService{}
	})
}

// GetMerchantPayinService 获取Payin服务单例
func GetMerchantPayinService() *MerchantPayinService {
	if merchantPayinService == nil {
		SetupMerchantPayinService()
	}
	return merchantPayinService
}

func (s *MerchantPayinService) Create() (trx *protocol.Transaction, code protocol.ErrorCode) {
	return &protocol.Transaction{}, protocol.Success
}
