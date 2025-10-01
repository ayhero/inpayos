package services

// InitializeServices 初始化所有服务单例
func InitializeServices() error {
	// 确保各个服务单例正确初始化
	GetPayinService()
	GetPayoutService()
	GetAccountService()
	GetCashierService()
	GetCheckoutService()
	GetTransactionService()
	return nil
}
