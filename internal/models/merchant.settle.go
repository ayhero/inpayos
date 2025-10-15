package models

import (
	"fmt"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type MerchantSettleLog struct {
	ID                       int64            `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	SettleID                 string           `json:"settle_id" gorm:"column:settle_id;index"`       // SettleID 结算ID
	MID                      int64            `json:"mid" gorm:"column:mid;index"`                   // MerchantID 商户ID
	PeriodType               string           `json:"period_type" gorm:"column:period_type;index"`   // SettlePeriodType 结算周期类型: D0, D1, D7, M1
	Period                   int64            `json:"period" gorm:"column:period;index"`             // SettlePeriod 结算周期
	TrxStartAt               int64            `json:"trx_start_at" gorm:"column:trx_start_at"`       // TrxStartAt 交易开始时间
	TrxEndAt                 int64            `json:"trx_end_at" gorm:"column:trx_end_at"`           // TrxEndAt 交易结束时间
	SettleStartAt            int64            `json:"settle_start_at" gorm:"column:settle_start_at"` // SettleStartAt 结算开始时间
	SettleEndAt              int64            `json:"settle_end_at" gorm:"column:settle_end_at"`     // SettleEndAt 结算结束时间
	TrxTotal                 int64            `json:"trx_total" gorm:"column:trx_total"`             // TrxTotal 交易笔数
	TrxCcy                   string           `json:"trx_ccy" gorm:"column:trx_ccy"`                 // TrxCcy 交易币种
	TrxAmount                *decimal.Decimal `json:"trx_amount" gorm:"column:trx_amount"`           // TrxAmount 交易金额
	TrxUsdAmount             *decimal.Decimal `json:"trx_usd_amount" gorm:"column:trx_usd_amount"`   // TrxUsdAmount 交易美元金额
	SettleCcy                string           `json:"settle_ccy" gorm:"column:settle_ccy"`           // SettleCcy 结算币种
	TrxType                  string           `json:"trx_type" gorm:"column:trx_type;index"`         // TrxType 交易类型
	*MerchantSettleLogValues `gorm:"embedded"`
	CreatedAt                int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt                int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

type MerchantSettleLogValues struct {
	Query           *TrxQuery        `json:"query" gorm:"column:query;query;type:json;serializer:json"`             // Query 查询条件
	Status          *string          `json:"status" gorm:"column:status;index"`                                     // Status 状态
	SettleAmount    *decimal.Decimal `json:"settle_amount" gorm:"column:settle_amount"`                             // TrxSettleAmount 结算金额
	SettleUsdAmount *decimal.Decimal `json:"settle_usd_amount" gorm:"column:settle_usd_amount"`                     // TrxSettleUsdAmount 结算美元金额
	FeeCcy          *string          `json:"fee_ccy" gorm:"column:fee_ccy"`                                         // FeeCcy 手续费币种
	Fee             *decimal.Decimal `json:"fee" gorm:"column:fee"`                                                 // TrxFee 交易手续费
	UsdFee          *decimal.Decimal `json:"usd_fee" gorm:"column:usd_fee"`                                         // TrxUsdFee 交易美元手续费
	FixedFee        *decimal.Decimal `json:"fixed_fee" gorm:"column:fixed_fee"`                                     // TrxFixedFee 交易固定手续费
	FixedUsdFee     *decimal.Decimal `json:"fixed_usd_fee" gorm:"column:fixed_usd_fee"`                             // TrxFixedUsdFee 交易固定美元手续费
	Rate            *decimal.Decimal `json:"rate" gorm:"column:rate"`                                               // TrxRate 交易费率
	UsdRate         *decimal.Decimal `json:"usd_rate" gorm:"column:usd_rate"`                                       // TrxUsdRate 交易美元费率
	SettleRules     SettleRules      `json:"settle_rules" gorm:"column:settle_rules;type:json;serializer:json"`     // SettleRules 结算规则
	StrategyCodes   []string         `json:"strategy_codes" gorm:"column:strategy_codes;type:json;serializer:json"` // StrategyCodes 结算策略代码列表
	SettleFile      *string          `json:"settle_file" gorm:"column:settle_file"`                                 // SettleFile 结算文件
	CompletedAt     *int64           `json:"completed_at" gorm:"column:completed_at"`                               // CompletedAt 结算完成时间
}

func (t MerchantSettleLog) TableName() string {
	return "t_merchant_settle_log"
}

// Getters for MerchantSettleLog
func (t *MerchantSettleLog) GetTrxAmount() decimal.Decimal {
	if t.TrxAmount == nil {
		return decimal.Zero
	}
	return *t.TrxAmount
}
func (t *MerchantSettleLog) GetTrxUsdAmount() decimal.Decimal {
	if t.TrxUsdAmount == nil {
		return decimal.Zero
	}
	return *t.TrxUsdAmount
}

func (t *MerchantSettleLog) SetTrxAmount(amount decimal.Decimal) *MerchantSettleLog {
	t.TrxAmount = &amount
	return t
}
func (t *MerchantSettleLog) SetTrxUsdAmount(amount decimal.Decimal) *MerchantSettleLog {
	t.TrxUsdAmount = &amount
	return t
}

// Getters
func (v *MerchantSettleLogValues) GetStatus() string {
	return *v.Status
}
func (v *MerchantSettleLogValues) GetSettleAmount() decimal.Decimal {
	if v.SettleAmount == nil {
		return decimal.Zero
	}
	return *v.SettleAmount
}
func (v *MerchantSettleLogValues) GetSettleUsdAmount() decimal.Decimal {
	if v.SettleUsdAmount == nil {
		return decimal.Zero
	}
	return *v.SettleUsdAmount
}
func (v *MerchantSettleLogValues) GetFeeCcy() string {
	if v.FeeCcy == nil {
		return ""
	}
	return *v.FeeCcy
}
func (v *MerchantSettleLogValues) GetFee() decimal.Decimal {
	if v.Fee == nil {
		return decimal.Zero
	}
	return *v.Fee
}
func (v *MerchantSettleLogValues) GetUsdFee() decimal.Decimal {
	if v.UsdFee == nil {
		return decimal.Zero
	}
	return *v.UsdFee
}
func (v *MerchantSettleLogValues) GetFixedFee() decimal.Decimal {
	if v.FixedFee == nil {
		return decimal.Zero
	}
	return *v.FixedFee
}
func (v *MerchantSettleLogValues) GetFixedUsdFee() decimal.Decimal {
	if v.FixedUsdFee == nil {
		return decimal.Zero
	}
	return *v.FixedUsdFee
}
func (v *MerchantSettleLogValues) GetRate() decimal.Decimal {
	if v.Rate == nil {
		return decimal.Zero
	}
	return *v.Rate
}
func (v *MerchantSettleLogValues) GetUsdRate() decimal.Decimal {
	if v.UsdRate == nil {
		return decimal.Zero
	}
	return *v.UsdRate
}
func (v *MerchantSettleLogValues) GetSettleRules() SettleRules {
	if v.SettleRules == nil {
		return SettleRules{}
	}
	return v.SettleRules
}
func (v *MerchantSettleLogValues) GetSettleFile() string {
	return *v.SettleFile
}
func (v *MerchantSettleLogValues) GetCompletedAt() int64 {
	return *v.CompletedAt
}
func (v *MerchantSettleLogValues) GetStrategyCodes() []string {
	return v.StrategyCodes
}

// Setters
func (v *MerchantSettleLogValues) SetStatus(status string) *MerchantSettleLogValues {
	v.Status = &status
	return v
}
func (v *MerchantSettleLogValues) SetSettleAmount(amount decimal.Decimal) *MerchantSettleLogValues {
	v.SettleAmount = &amount
	return v
}
func (v *MerchantSettleLogValues) SetSettleUsdAmount(amount decimal.Decimal) *MerchantSettleLogValues {
	v.SettleUsdAmount = &amount
	return v
}
func (v *MerchantSettleLogValues) SetFeeCcy(feeCcy string) *MerchantSettleLogValues {
	v.FeeCcy = &feeCcy
	return v
}
func (v *MerchantSettleLogValues) SetFee(fee decimal.Decimal) *MerchantSettleLogValues {
	v.Fee = &fee
	return v
}
func (v *MerchantSettleLogValues) SetUsdFee(usdFee decimal.Decimal) *MerchantSettleLogValues {
	v.UsdFee = &usdFee
	return v
}
func (v *MerchantSettleLogValues) SetFixedFee(fixedFee decimal.Decimal) *MerchantSettleLogValues {
	v.FixedFee = &fixedFee
	return v
}
func (v *MerchantSettleLogValues) SetFixedUsdFee(fixedUsdFee decimal.Decimal) *MerchantSettleLogValues {
	v.FixedUsdFee = &fixedUsdFee
	return v
}
func (v *MerchantSettleLogValues) SetRate(rate decimal.Decimal) *MerchantSettleLogValues {
	v.Rate = &rate
	return v
}
func (v *MerchantSettleLogValues) SetUsdRate(usdRate decimal.Decimal) *MerchantSettleLogValues {
	v.UsdRate = &usdRate
	return v
}
func (v *MerchantSettleLogValues) SetSettleRules(rules SettleRules) *MerchantSettleLogValues {
	v.SettleRules = rules
	return v
}
func (v *MerchantSettleLogValues) SetSettleFile(file string) *MerchantSettleLogValues {
	v.SettleFile = &file
	return v
}
func (v *MerchantSettleLogValues) SetCompletedAt(completedAt int64) *MerchantSettleLogValues {
	v.CompletedAt = &completedAt
	return v
}
func (v *MerchantSettleLogValues) SetStrategyCodes(codes []string) *MerchantSettleLogValues {
	v.StrategyCodes = codes
	return v
}

// SetValues updates the MerchantSettleLogValues of the settle log
func (t *MerchantSettleLog) SetValues(values *MerchantSettleLogValues) {
	if values == nil {
		return
	}

	if t.MerchantSettleLogValues == nil {
		t.MerchantSettleLogValues = &MerchantSettleLogValues{}
	}

	if values.Status != nil && *values.Status != "" {
		t.Status = values.Status
	}
	if values.SettleAmount != nil {
		t.SettleAmount = values.SettleAmount
	}
	if values.SettleUsdAmount != nil {
		t.SettleUsdAmount = values.SettleUsdAmount
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
	if values.SettleRules != nil {
		t.SettleRules = values.SettleRules
	}
	if values.SettleFile != nil && *values.SettleFile != "" {
		t.SettleFile = values.SettleFile
	}
	if values.CompletedAt != nil && *values.CompletedAt != 0 {
		t.CompletedAt = values.CompletedAt
	}
}

// GetOrCreateSettleLog 获取或创建结算周期记录,根据执行结算的当时
func GetOrCreateSettleLog(mid string, completedAt, settleAt int64, settleConfig *ContractSettleSetting) *MerchantSettleLog {
	settleLog := GetMerchantSettleLogByMIDAndPeriod(mid, settleConfig.TrxType, settleAt, settleConfig.Type)
	if settleLog != nil {
		return settleLog
	}
	period, startAt, endAt := settleConfig.GetSettlePeriodByTime(completedAt, settleAt)
	// 创建新的结算周期记录
	settleLog = &MerchantSettleLog{
		SettleID:      utils.GenerateSettleLogID(),
		MID:           cast.ToInt64(mid),
		PeriodType:    settleConfig.Type,
		Period:        period,
		SettleStartAt: startAt,
		SettleEndAt:   endAt,
		TrxTotal:      0,
		TrxCcy:        "",
		SettleCcy:     settleConfig.Ccy,
		MerchantSettleLogValues: &MerchantSettleLogValues{
			StrategyCodes: settleConfig.Strategies, // 存储策略代码
		},
	}
	settleLog.SetStatus(protocol.StatusPending)
	// 保存到数据库
	if err := WriteDB.Create(settleLog).Error; err != nil {
		return nil
	}
	return settleLog
}

// RefreshSettleLog 更新结算周期记录的交易统计信息
func RefreshSettleLog(tx *gorm.DB, settleLogID string) error {
	return nil
}

// GetSettleLogByID 根据ID获取结算周期记录
func GetSettleLogByID(settleID string) (*MerchantSettleLog, error) {
	if settleID == "" {
		return nil, fmt.Errorf("settleID cannot be empty")
	}

	var settleLog MerchantSettleLog
	err := ReadDB.Where("settle_id = ?", settleID).First(&settleLog).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &settleLog, nil
}

func GetMerchantSettleLogByMIDAndPeriod(mid, trx_type string, period int64, periodType string) *MerchantSettleLog {
	var settleLog MerchantSettleLog
	err := ReadDB.Where("mid = ? AND period = ? AND period_type = ? AND trx_type = ?", mid, period, periodType, trx_type).First(&settleLog).Error
	if err != nil {
		return nil
	}
	return &settleLog
}

// NewSettleTransaction 创建新的结算交易记录，挂靠到结算周期记录
func NewSettleTransaction(trx *Transaction, settleLogID string) *MerchantSettleTransaction {
	settleTransaction := &MerchantSettleTransaction{
		TrxID:        trx.TrxID,
		SettleID:     utils.GenerateSettleTrxID(),
		MID:          trx.Mid,
		TrxType:      trx.TrxType,
		TrxCcy:       trx.Ccy,
		TrxAmount:    trx.Amount,
		TrxUsdAmount: trx.UsdAmount,
		TrxAt:        trx.CreatedAt,
		SettleCcy:    trx.Ccy,
		MerchantSettleTransactionValues: &MerchantSettleTransactionValues{
			SettleLogID: &settleLogID, // 直接关联结算周期记录
		},
	}
	settleTransaction.SetStatus(protocol.StatusPending)
	return settleTransaction
}

// NewSettleTransactionRecord 创建包含完整信息的结算交易记录
func NewSettleTransactionRecord(trx *Transaction, settleLogID string, strategy *protocol.SettleStrategy, rule *protocol.SettleRule, result *protocol.SettlementResult) *MerchantSettleTransaction {
	settleTransaction := NewSettleTransaction(trx, settleLogID)

	settleTransaction.MerchantSettleTransactionValues = &MerchantSettleTransactionValues{
		SettleAmount:    &result.SettleAmount,
		SettleUsdAmount: &result.SettleUsdAmount,
		FeeCcy:          &result.FeeCcy,
		Fee:             &result.Fee,
		UsdFee:          &result.UsdFee,
		FixedFee:        &result.FixedFee,
		FixedUsdFee:     &result.FixedUsdFee,
		Rate:            &result.Rate,
		UsdRate:         &result.UsdRate,
		SettleStrategy:  strategy,
		SettleRule:      rule,
	}

	return settleTransaction
}

// GetPendingAccountingSettleLogs 获取待记账的结算记录
// 返回所有 settle_end_at 已过期且状态为成功但尚未记账的结算记录
func GetPendingAccountingSettleLogs(currentTime int64) ([]*MerchantSettleLog, error) {
	var settleLogs []*MerchantSettleLog

	// 查询条件：
	// 1. settle_end_at < currentTime (结算周期已结束)
	// 2. status = 'success' (结算成功)
	// 3. completed_at IS NULL 或者需要增加一个专门的记账状态字段
	err := ReadDB.Where(`
		settle_end_at < ? AND 
		status = ? AND 
		(completed_at IS NULL OR completed_at = 0)
	`, currentTime, protocol.StatusSuccess).Find(&settleLogs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query pending accounting settle logs: %v", err)
	}

	return settleLogs, nil
}
