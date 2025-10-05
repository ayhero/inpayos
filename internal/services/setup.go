package services

// InitializeMerchantServices 初始化所有服务单例
func InitializeMerchantServices() error {
	// 确保各个服务单例正确初始化
	GetMerchantPayinService()
	GetMerchantPayoutService()
	GetAccountService()
	GetCashierService()
	GetCheckoutService()
	GetMerchantTransactionService()
	return nil
}
