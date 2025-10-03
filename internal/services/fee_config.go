package services

import "sync"

type FeeConfigService struct {
}

var (
	feeConfigService     *FeeConfigService
	feeConfigServiceOnce sync.Once
)

func SetupFeeConfigService() {
	feeConfigServiceOnce.Do(func() {
		feeConfigService = &FeeConfigService{}
	})
}

// GetFeeConfigService 获取FeeConfig服务单例
func GetFeeConfigService() *FeeConfigService {
	if feeConfigService == nil {
		SetupFeeConfigService()
	}
	return feeConfigService
}
