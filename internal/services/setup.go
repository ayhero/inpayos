package services

// SetupService 初始化所有服务
func SetupService() {
	// 初始化各种服务单例
	// 这里可以添加服务初始化逻辑

	// 确保各个服务单例正确初始化
	GetAccountService()
	GetTransactionService()
	GetWebhookService()
	GetChannelService()
}

// 全局服务实例
var (
	channelService *ChannelService
)

// GetChannelService 获取渠道服务单例
func GetChannelService() *ChannelService {
	if channelService == nil {
		channelService = NewChannelService()
	}
	return channelService
}
