package services

import "sync"

type ChannelService struct {
}

var (
	channelService     *ChannelService
	channelServiceOnce sync.Once
)

func SetupChannelService() {
	channelServiceOnce.Do(func() {
		channelService = &ChannelService{}
	})
}

// GetChannelService 获取Channel服务单例
func GetChannelService() *ChannelService {
	if channelService == nil {
		SetupChannelService()
	}
	return channelService
}
