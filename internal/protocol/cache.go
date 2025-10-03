package protocol

import "time"

// Cache Keys
const (
	// FundFlow cache
	FundFlowCacheKeyPrefix = "fund_flow:"
	FundFlowByFlowNoTpl    = FundFlowCacheKeyPrefix + "flow_no:%s"
	FundFlowByBizInfoTpl   = FundFlowCacheKeyPrefix + "biz:%s:%s"
	FundFlowByOriFlowNoTpl = FundFlowCacheKeyPrefix + "ori_flow_no:%s"

	// G2FA cache
	G2FACacheKeyPrefix = "g2fa:"
	G2FABindingTpl     = G2FACacheKeyPrefix + "binding:%s"
)

// Cache Expirations
const (
	FundFlowCacheExpiration = time.Hour
	G2FACacheExpiration     = 10 * time.Minute // G2FA binding should expire in 10 minutes
)
