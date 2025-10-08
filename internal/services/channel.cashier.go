package services

import "sync"

type CashierChannelService struct {
}

var (
	cashierChannelService     *CashierChannelService
	cashierChannelServiceOnce sync.Once
)

func SetupCashierChannelService() {
	cashierChannelServiceOnce.Do(func() {
		cashierChannelService = &CashierChannelService{}
	})
}

// GetCashierChannelService 获取CashierChannel服务单例
func GetCashierChannelService() *CashierChannelService {
	if cashierChannelService == nil {
		SetupCashierChannelService()
	}
	return cashierChannelService
}
