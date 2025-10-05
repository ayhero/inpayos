package services

import "sync"

type MerchantFeeConfigService struct {
}

var (
	merchantfeeConfigService     *MerchantFeeConfigService
	merchantfeeConfigServiceOnce sync.Once
)

func SetupMerchantFeeConfigService() {
	merchantfeeConfigServiceOnce.Do(func() {
		merchantfeeConfigService = &MerchantFeeConfigService{}
	})
}

// GetFeeConfigService 获取FeeConfig服务单例
func GetFeeConfigService() *MerchantFeeConfigService {
	if merchantfeeConfigService == nil {
		SetupMerchantFeeConfigService()
	}
	return merchantfeeConfigService
}
