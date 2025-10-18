package models

import (
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
)

type MerchantSettleTransaction struct {
	ID                               int64            `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TrxID                            string           `json:"trx_id" gorm:"column:trx_id;index;uniqueIndex:idx_trx_settle"`       // TrxID 交易ID
	SettleID                         string           `json:"settle_id" gorm:"column:settle_id;index;uniqueIndex:idx_trx_settle"` // SettleID 结算ID
	Mid                              string           `json:"mid" gorm:"column:mid;type:varchar(64);index"`                       // MerchantID 商户ID
	TrxType                          string           `json:"trx_type" gorm:"column:trx_type"`                                    // TrxType 交易类型
	TrxCcy                           string           `json:"trx_ccy" gorm:"column:trx_ccy"`                                      // TrxCcy 交易币种
	TrxAmount                        *decimal.Decimal `json:"trx_amount" gorm:"column:trx_amount"`                                // TrxAmount 交易金额
	TrxUsdAmount                     *decimal.Decimal `json:"trx_usd_amount" gorm:"column:trx_usd_amount"`                        // TrxUsdAmount 交易美元金额
	TrxAt                            int64            `json:"trx_at" gorm:"column:trx_at"`                                        // TrxAt 交易时间
	SettleCcy                        string           `json:"settle_ccy" gorm:"column:settle_ccy"`                                // SettleCcy 结算币种
	*MerchantSettleTransactionValues `gorm:"embedded"`
	CreatedAt                        int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt                        int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

type MerchantSettleTransactionValues struct {
	SettleLogID     *string                  `json:"settle_log_id" gorm:"column:settle_log_id;index"`                         // SettleLogID 结算周期记录ID
	FlowNo          *string                  `json:"flow_no" gorm:"column:flow_no;index"`                                     // FlowNo 流水号
	SettleAmount    *decimal.Decimal         `json:"settle_amount" gorm:"column:settle_amount"`                               // TrxSettleAmount 结算金额
	SettleUsdAmount *decimal.Decimal         `json:"settle_usd_amount" gorm:"column:settle_usd_amount"`                       // TrxSettleUsdAmount 结算美元金额
	SettledAt       *int64                   `json:"settled_at" gorm:"column:settled_at"`                                     // SettledAt 结算时间
	FeeCcy          *string                  `json:"fee_ccy" gorm:"column:fee_ccy"`                                           // FeeCcy 手续费币种
	Fee             *decimal.Decimal         `json:"fee" gorm:"column:fee"`                                                   // TrxFee 交易手续费
	UsdFee          *decimal.Decimal         `json:"usd_fee" gorm:"column:usd_fee"`                                           // TrxUsdFee 交易美元手续费
	FixedFee        *decimal.Decimal         `json:"fixed_fee" gorm:"column:fixed_fee"`                                       // TrxFixedFee 交易固定手续费
	FixedUsdFee     *decimal.Decimal         `json:"fixed_usd_fee" gorm:"column:fixed_usd_fee"`                               // TrxFixedUsdFee 交易固定美元手续费
	Rate            *decimal.Decimal         `json:"rate" gorm:"column:rate"`                                                 // TrxRate 交易费率
	UsdRate         *decimal.Decimal         `json:"usd_rate" gorm:"column:usd_rate"`                                         // TrxUsdRate 交易美元费率
	SettleStrategy  *protocol.SettleStrategy `json:"settle_strategy" gorm:"column:settle_strategy;type:json;serializer:json"` // SettleStrategy 结算策略
	SettleRule      *protocol.SettleRule     `json:"settle_rule" gorm:"column:settle_rule;type:json;serializer:json"`         // SettleRule 结算规则
	Status          *string                  `json:"status" gorm:"column:status;index"`                                       // Status 状态
}

func (t MerchantSettleTransaction) TableName() string {
	return "t_merchant_settle_transaction"
}

// Getters for MerchantSettleTransaction
func (t *MerchantSettleTransaction) GetTrxAmount() decimal.Decimal {
	if t.TrxAmount == nil {
		return decimal.Zero
	}
	return *t.TrxAmount
}
func (t *MerchantSettleTransaction) GetTrxUsdAmount() decimal.Decimal {
	if t.TrxUsdAmount == nil {
		return decimal.Zero
	}
	return *t.TrxUsdAmount
}

func (t *MerchantSettleTransaction) SetTrxAmount(amount decimal.Decimal) *MerchantSettleTransaction {
	t.TrxAmount = &amount
	return t
}
func (t *MerchantSettleTransaction) SetTrxUsdAmount(amount decimal.Decimal) *MerchantSettleTransaction {
	t.TrxUsdAmount = &amount
	return t
}

// Getters for MerchantSettleTransactionValues
func (v *MerchantSettleTransactionValues) GetSettleAmount() decimal.Decimal {
	if v.SettleAmount == nil {
		return decimal.Zero
	}
	return *v.SettleAmount
}
func (v *MerchantSettleTransactionValues) GetSettleUsdAmount() decimal.Decimal {
	if v.SettleUsdAmount == nil {
		return decimal.Zero
	}
	return *v.SettleUsdAmount
}
func (v *MerchantSettleTransactionValues) GetSettledAt() int64 {
	if v.SettledAt == nil {
		return 0
	}
	return *v.SettledAt
}
func (v *MerchantSettleTransactionValues) GetFeeCcy() string {
	if v.FeeCcy == nil {
		return ""
	}
	return *v.FeeCcy
}
func (v *MerchantSettleTransactionValues) GetFee() decimal.Decimal {
	if v.Fee == nil {
		return decimal.Zero
	}
	return *v.Fee
}
func (v *MerchantSettleTransactionValues) GetUsdFee() decimal.Decimal {
	if v.UsdFee == nil {
		return decimal.Zero
	}
	return *v.UsdFee
}
func (v *MerchantSettleTransactionValues) GetFixedFee() decimal.Decimal {
	if v.FixedFee == nil {
		return decimal.Zero
	}
	return *v.FixedFee
}
func (v *MerchantSettleTransactionValues) GetFixedUsdFee() decimal.Decimal {
	if v.FixedUsdFee == nil {
		return decimal.Zero
	}
	return *v.FixedUsdFee
}
func (v *MerchantSettleTransactionValues) GetRate() decimal.Decimal {
	if v.Rate == nil {
		return decimal.Zero
	}
	return *v.Rate
}
func (v *MerchantSettleTransactionValues) GetUsdRate() decimal.Decimal {
	if v.UsdRate == nil {
		return decimal.Zero
	}
	return *v.UsdRate
}
func (v *MerchantSettleTransactionValues) GetStatus() string {
	if v.Status == nil {
		return ""
	}
	return *v.Status
}
func (v *MerchantSettleTransactionValues) GetFlowNo() string {
	if v.FlowNo == nil {
		return ""
	}
	return *v.FlowNo
}
func (v *MerchantSettleTransactionValues) GetSettleStrategy() *protocol.SettleStrategy {
	return v.SettleStrategy
}
func (v *MerchantSettleTransactionValues) GetSettleRule() *protocol.SettleRule {
	return v.SettleRule
}

// Setters for MerchantSettleTransactionValues
func (v *MerchantSettleTransactionValues) SetSettleAmount(amount decimal.Decimal) *MerchantSettleTransactionValues {
	v.SettleAmount = &amount
	return v
}
func (v *MerchantSettleTransactionValues) SetSettleUsdAmount(amount decimal.Decimal) *MerchantSettleTransactionValues {
	v.SettleUsdAmount = &amount
	return v
}
func (v *MerchantSettleTransactionValues) SetSettledAt(settledAt int64) *MerchantSettleTransactionValues {
	v.SettledAt = &settledAt
	return v
}
func (v *MerchantSettleTransactionValues) SetFeeCcy(feeCcy string) *MerchantSettleTransactionValues {
	v.FeeCcy = &feeCcy
	return v
}
func (v *MerchantSettleTransactionValues) SetFee(fee decimal.Decimal) *MerchantSettleTransactionValues {
	v.Fee = &fee
	return v
}
func (v *MerchantSettleTransactionValues) SetUsdFee(usdFee decimal.Decimal) *MerchantSettleTransactionValues {
	v.UsdFee = &usdFee
	return v
}
func (v *MerchantSettleTransactionValues) SetFixedFee(fixedFee decimal.Decimal) *MerchantSettleTransactionValues {
	v.FixedFee = &fixedFee
	return v
}
func (v *MerchantSettleTransactionValues) SetFixedUsdFee(fixedUsdFee decimal.Decimal) *MerchantSettleTransactionValues {
	v.FixedUsdFee = &fixedUsdFee
	return v
}
func (v *MerchantSettleTransactionValues) SetRate(rate decimal.Decimal) *MerchantSettleTransactionValues {
	v.Rate = &rate
	return v
}
func (v *MerchantSettleTransactionValues) SetUsdRate(usdRate decimal.Decimal) *MerchantSettleTransactionValues {
	v.UsdRate = &usdRate
	return v
}
func (v *MerchantSettleTransactionValues) SetStatus(status string) *MerchantSettleTransactionValues {
	v.Status = &status
	return v
}
func (v *MerchantSettleTransactionValues) SetFlowNo(flowNo string) *MerchantSettleTransactionValues {
	v.FlowNo = &flowNo
	return v
}
func (v *MerchantSettleTransactionValues) SetSettleStrategy(strategy *protocol.SettleStrategy) *MerchantSettleTransactionValues {
	v.SettleStrategy = strategy
	return v
}
func (v *MerchantSettleTransactionValues) SetSettleRule(rule *protocol.SettleRule) *MerchantSettleTransactionValues {
	v.SettleRule = rule
	return v
}

// SetSettleLogID sets the settle log ID
func (t *MerchantSettleTransactionValues) SetSettleLogID(settleLogID string) *MerchantSettleTransactionValues {
	t.SettleLogID = &settleLogID
	return t
}

// GetSettleLogID gets the settle log ID
func (t *MerchantSettleTransactionValues) GetSettleLogID() string {
	if t.SettleLogID == nil {
		return ""
	}
	return *t.SettleLogID
}

// SetValues updates the MerchantSettleTransactionValues
func (t *MerchantSettleTransaction) SetValues(values *MerchantSettleTransactionValues) {
	if values == nil {
		return
	}

	if t.MerchantSettleTransactionValues == nil {
		t.MerchantSettleTransactionValues = &MerchantSettleTransactionValues{}
	}

	if values.FlowNo != nil {
		t.FlowNo = values.FlowNo
	}
	if values.SettleAmount != nil {
		t.SettleAmount = values.SettleAmount
	}
	if values.SettleUsdAmount != nil {
		t.SettleUsdAmount = values.SettleUsdAmount
	}
	if values.SettledAt != nil {
		t.SettledAt = values.SettledAt
	}
	if values.FeeCcy != nil {
		t.FeeCcy = values.FeeCcy
	}
	if values.Fee != nil {
		t.Fee = values.Fee
	}
	if values.UsdFee != nil {
		t.UsdFee = values.UsdFee
	}
	if values.FixedFee != nil {
		t.FixedFee = values.FixedFee
	}
	if values.FixedUsdFee != nil {
		t.FixedUsdFee = values.FixedUsdFee
	}
	if values.Rate != nil {
		t.Rate = values.Rate
	}
	if values.UsdRate != nil {
		t.UsdRate = values.UsdRate
	}
	if values.SettleStrategy != nil {
		t.SettleStrategy = values.SettleStrategy
	}
	if values.SettleRule != nil {
		t.SettleRule = values.SettleRule
	}
	if values.Status != nil {
		t.Status = values.Status
	}
}

// GetExistingSettleRecord 获取现有的结算记录
func GetExistingSettleRecord(trxID string) *MerchantSettleTransaction {
	if trxID == "" {
		return nil
	}

	var settleRecord MerchantSettleTransaction
	err := ReadDB.Where("trx_id = ?", trxID).First(&settleRecord).Error
	if err != nil {
		return nil
	}

	return &settleRecord
}

func CountSettleTransactionsByTimeRange(startTime, endTime int64) int64 {
	var count int64
	err := ReadDB.Model(&MerchantSettleTransaction{}).Where(`
		settle_log_id IS NOT NULL AND settle_log_id != '' AND 
		created_at BETWEEN ? AND ?
	`, startTime, endTime).Count(&count).Error
	if err != nil {
		return 0
	}
	return count
}

func ListSettleTransactionsByTimeRange(startTime, endTime int64) []*MerchantSettleTransaction {
	var settleTransactions []*MerchantSettleTransaction
	err := ReadDB.Where(`
		settle_log_id IS NOT NULL AND settle_log_id != '' AND 
		created_at BETWEEN ? AND ?
	`, startTime, endTime).Find(&settleTransactions).Error
	if err != nil {
		return nil
	}
	return settleTransactions
}

// ListSettleTransactionsWithLimitedFields 查询结算交易记录（限制字段，支持分页）
// 只查询必要字段：settle_id, settle_log_id, trx_id, trx_type
func ListSettleTransactionsWithLimitedFields(startTime, endTime int64, offset, limit int) []*MerchantSettleTransaction {
	var settleTransactions []*MerchantSettleTransaction
	err := ReadDB.Select("settle_id, settle_log_id, trx_id, trx_type, created_at").Where(`
		settle_log_id IS NOT NULL AND settle_log_id != '' AND 
		created_at BETWEEN ? AND ?
	`, startTime, endTime).Offset(offset).Limit(limit).Find(&settleTransactions).Error

	if err != nil {
		return nil
	}
	return settleTransactions
}
