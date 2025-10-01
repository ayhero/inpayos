package services

import (
	"inpayos/internal/protocol"
	"sync"
)

type PayinService struct {
}

var (
	payinService     *PayinService
	payinServiceOnce sync.Once
)

func SetupPayinService() {
	payinServiceOnce.Do(func() {
		payinService = &PayinService{}
	})
}

// GetPayinService 获取Payin服务单例
func GetPayinService() *PayinService {
	if payinService == nil {
		SetupPayinService()
	}
	return payinService
}

func (s *PayinService) Create() (trx *protocol.Transaction, code protocol.ErrorCode) {
	return &protocol.Transaction{}, protocol.Success
}
