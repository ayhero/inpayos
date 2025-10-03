package services

import "sync"

type WebhookService struct {
}

var (
	webhookService     *WebhookService
	webhookServiceOnce sync.Once
)

func SetupWebhookService() {
	webhookServiceOnce.Do(func() {
		webhookService = &WebhookService{}
	})
}

// GetWebhookService 获取Webhook服务单例
func GetWebhookService() *WebhookService {
	if webhookService == nil {
		SetupWebhookService()
	}
	return webhookService
}
