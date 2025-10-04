package models

import (
	"inpayos/internal/log"
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
)

type Contract struct {
	ID            uint64 `gorm:"primaryKey;autoIncrement;<-:create" json:"id"`
	ContractID    string `gorm:"type:varchar(32);uniqueIndex;not null;<-:create" json:"contract_id"` // 合同ID(唯一)
	OriContractID string `gorm:"type:varchar(32);not null;<-:create" json:"ori_contract_id"`         // 原合同ID(用于更新时的原始合同ID)
	Mid           string `gorm:"type:varchar(32);not null;<-:create" json:"mid"`                     // 商户ID
	*ContractValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

func (t Contract) TableName() string {
	return "t_merchant_contract"
}

type ContractValues struct {
	StartAt      int64                 `json:"start_at" gorm:"index;comment:生效时间"`          // StartAt 生效时间
	ExpiredAt    int64                 `json:"expired_at" gorm:"index;comment:过期时间"`        // ExpiredAt 过期时间
	Status       int64                 `json:"status" gorm:"index;comment:状态"`              // Status 状态
	Payin        *PayinSetting         `json:"payin" gorm:"column:payin;serializer:json"`   // PayinSetting 充值配置
	Payout       *PayinSetting         `json:"payout" gorm:"column:payout;serializer:json"` // PayoutSetting 提现配置
	SettleConfig *ContractSettleConfig `json:"settle_config" gorm:"column:settle_config"`   // SettleConfig 结算配置
}

type PayinSetting struct {
	Status      int64                 `json:"status" gorm:"column:status"`                             // 状态
	SettingList []*ContractTrxSetting `json:"setting_list" gorm:"column:setting_list;serializer:json"` // 充值设置列表
}

type ContractTrxConfig struct {
	Payin  []*ContractTrxSetting `json:"payin" gorm:"column:payin"`
	Payout []*ContractTrxSetting `json:"payout" gorm:"column:payout"`
}

type ContractTrxSetting struct {
	Pkg          string           `json:"pkg" gorm:"column:pkg"`
	TrxType      string           `json:"trx_type" gorm:"column:trx_type"`
	TrxSubType   string           `json:"trx_sub_type" gorm:"column:trx_sub_type"`
	TrxMethod    string           `json:"trx_method" gorm:"column:trx_method"`
	Ccy          string           `json:"ccy" gorm:"column:ccy"`
	Country      string           `json:"country" gorm:"column:country"` // 国家代码
	MinAmount    *decimal.Decimal `json:"min_amount" gorm:"column:min_amount"`
	MaxAmount    *decimal.Decimal `json:"max_amount" gorm:"column:max_amount"`
	MinUsdAmount *decimal.Decimal `json:"min_usd_amount" gorm:"column:min_usd_amount"`
	MaxUsdAmount *decimal.Decimal `json:"max_usd_amount" gorm:"column:max_usd_amount"`
}

type ContractSettleConfig struct {
	Type   string                   `json:"type" gorm:"column:type"` // 结算类型，如：TRX, TRX_APP, TRX_PKG
	Ccy    string                   `json:"ccy" gorm:"column:ccy"`   // 结算币种，如：CNY, USD
	Payin  []*ContractSettleSetting `json:"payin" gorm:"column:payin"`
	Payout []*ContractSettleSetting `json:"payout" gorm:"column:payout"`
}

type ContractSettleSetting struct {
	Type       string   `json:"type" gorm:"column:type"`             // 结算周期，如：D0, D1, D7, M1
	Ccy        string   `json:"ccy" gorm:"column:ccy"`               // 结算币种，如：CNY, USD
	Strategies []string `json:"strategies" gorm:"column:strategies"` // 结算策略
}

func ListMerchantContractByMid(mid string) []*Contract {
	var contracts []*Contract
	if err := ReadDB.Where("mid = ? and status=?", mid, protocol.StatusActive).Find(&contracts).Error; err != nil {
		log.Get().Errorf("ListMerchantContractByMid failed, mid: %s, err: %v", mid, err)
		return nil
	}
	return contracts
}

// GetValidContractsAtTime 获取在指定时间有效的合同
func GetValidContractsAtTime(mid string, trxTime int64) *Contract {
	if mid == "" {
		return nil
	}

	var contract *Contract
	err := ReadDB.Where("mid = ? AND status = ? AND start_at <= ? AND (expired_at = 0 OR expired_at >= ?)",
		mid, protocol.StatusActive, trxTime, trxTime).Order("id desc").First(&contract).Error
	if err != nil {
		log.Get().Errorf("getValidContractsAtTime failed, mid: %s, trxTime: %d, err: %v", mid, trxTime, err)
		return nil
	}

	return contract
}
