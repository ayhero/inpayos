package services

import "sync"

type CashierService struct {
}

var (
	cashierService     *CashierService
	cashierServiceOnce sync.Once
)

func SetupCashierService() {
	cashierServiceOnce.Do(func() {
		cashierService = &CashierService{}
	})
}

// GetCashierService 获取Cashier服务单例
func GetCashierService() *CashierService {
	if cashierService == nil {
		SetupCashierService()
	}
	return cashierService
}
