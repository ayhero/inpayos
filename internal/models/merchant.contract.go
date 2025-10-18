package models

import (
	"fmt"
	"inpayos/internal/log"
	"inpayos/internal/protocol"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

type MerchantContract struct {
	ID            uint64 `gorm:"primaryKey;autoIncrement;<-:create" json:"id"`
	ContractID    string `gorm:"type:varchar(32);uniqueIndex;not null;<-:create" json:"contract_id"` // 合同ID(唯一)
	OriContractID string `gorm:"type:varchar(32);not null;<-:create" json:"ori_contract_id"`         // 原合同ID(用于更新时的原始合同ID)
	Mid           string `gorm:"type:varchar(64);not null;<-:create" json:"mid"`                     // 商户ID
	*MerchantContractValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

func (t MerchantContract) TableName() string {
	return "t_merchant_contracts"
}

type MerchantContractValues struct {
	StartAt   *int64                  `json:"start_at" gorm:"index;comment:生效时间"`                    // StartAt 生效时间
	ExpiredAt *int64                  `json:"expired_at" gorm:"index;comment:过期时间"`                  // ExpiredAt 过期时间
	Status    *string                 `json:"status" gorm:"index;comment:状态"`                        // Status 状态
	Payin     *MerchantContractConfig `json:"payin" gorm:"column:payin;type:json;serializer:json"`   // PayinSetting 充值配置
	Payout    *MerchantContractConfig `json:"payout" gorm:"column:payout;type:json;serializer:json"` // PayoutSetting 提现配置
}

type MerchantContractConfig struct {
	TrxType string                `json:"trx_type" gorm:"column:trx_type"`               // 交易类型
	Status  string                `json:"status" gorm:"column:status"`                   // 状态
	Configs []*ContractTrxConfig  `json:"configs" gorm:"column:configs;serializer:json"` // 配置列表
	Settle  *ContractSettleConfig `json:"settle" gorm:"column:settle;serializer:json"`   // 结算配置
}

type ContractTrxConfig struct {
	Pkg          string           `json:"pkg" gorm:"column:pkg"`
	TrxType      string           `json:"trx_type" gorm:"column:trx_type"`
	TrxMethod    string           `json:"trx_method" gorm:"column:trx_method"`
	Ccy          string           `json:"ccy" gorm:"column:ccy"`
	Country      string           `json:"country" gorm:"column:country"` // 国家代码
	MinAmount    *decimal.Decimal `json:"min_amount" gorm:"column:min_amount"`
	MaxAmount    *decimal.Decimal `json:"max_amount" gorm:"column:max_amount"`
	MinUsdAmount *decimal.Decimal `json:"min_usd_amount" gorm:"column:min_usd_amount"`
	MaxUsdAmount *decimal.Decimal `json:"max_usd_amount" gorm:"column:max_usd_amount"`
}

// 合同结算配置
type ContractSettleConfig struct {
	Type         string           `json:"type" gorm:"column:type"` // 结算类型，如 T0、T1、T2、T3、W1、M1
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
	Strategies   []string         `json:"strategies" gorm:"column:strategies;serializer:json;type:json"` // 结算策略
}

func (c *MerchantContract) GetSettleConfig(trxType string) *ContractSettleConfig {
	switch trxType {
	case protocol.TrxTypePayin:
		if _v := c.Payin; _v != nil {
			return _v.Settle
		}
	case protocol.TrxTypePayout:
		if _v := c.Payout; _v != nil {
			return _v.Settle
		}
	}
	return nil
}

// GetSettlePeriodByTime 根据交易完成时间和结算周期类型计算结算周期及其开始和结束时间
// 参数 completedAt: 交易完成时间（毫秒时间戳）
// 参数 executeAt: 定时任务执行时间（毫秒时间戳），如果理论结算时间已过，则使用执行时间计算周期
func (c *ContractSettleConfig) GetSettlePeriodByTime(completedAt, executeAt int64) (period, startAt, endAt int64) {
	completedTime := time.UnixMilli(completedAt)
	executeTime := time.UnixMilli(executeAt)
	periodStr := ""
	var periodStart, periodEnd time.Time

	switch c.Type {
	case protocol.T0:
		// T+0 当天结算 - 交易完成当天结算
		theoreticalSettleDay := time.Date(completedTime.Year(), completedTime.Month(), completedTime.Day(), 0, 0, 0, 0, completedTime.Location())

		// 如果理论结算日期已经过去，使用执行时间的日期
		if executeTime.After(theoreticalSettleDay.Add(24*time.Hour - time.Millisecond)) {
			// 理论结算时间已过，使用执行时间所在日期
			actualSettleDay := time.Date(executeTime.Year(), executeTime.Month(), executeTime.Day(), 0, 0, 0, 0, executeTime.Location())
			periodStart = actualSettleDay
			periodEnd = actualSettleDay.Add(24*time.Hour - time.Millisecond)
			periodStr = actualSettleDay.Format("20060102")
		} else {
			// 理论结算时间未过，使用原计划时间
			periodStart = theoreticalSettleDay
			periodEnd = theoreticalSettleDay.Add(24*time.Hour - time.Millisecond)
			periodStr = theoreticalSettleDay.Format("20060102")
		}
	case protocol.T1:
		// T+1 次日结算 - 交易完成次日结算
		theoreticalSettleDay := time.Date(completedTime.Year(), completedTime.Month(), completedTime.Day()+1, 0, 0, 0, 0, completedTime.Location())

		// 如果理论结算日期已经过去，使用执行时间的日期
		if executeTime.After(theoreticalSettleDay.Add(24*time.Hour - time.Millisecond)) {
			actualSettleDay := time.Date(executeTime.Year(), executeTime.Month(), executeTime.Day(), 0, 0, 0, 0, executeTime.Location())
			periodStart = actualSettleDay
			periodEnd = actualSettleDay.Add(24*time.Hour - time.Millisecond)
			periodStr = actualSettleDay.Format("20060102")
		} else {
			periodStart = theoreticalSettleDay
			periodEnd = theoreticalSettleDay.Add(24*time.Hour - time.Millisecond)
			periodStr = theoreticalSettleDay.Format("20060102")
		}
	case protocol.T2:
		// T+2 两日后结算 - 交易完成两日后结算
		theoreticalSettleDay := time.Date(completedTime.Year(), completedTime.Month(), completedTime.Day()+2, 0, 0, 0, 0, completedTime.Location())

		// 如果理论结算日期已经过去，使用执行时间的日期
		if executeTime.After(theoreticalSettleDay.Add(24*time.Hour - time.Millisecond)) {
			actualSettleDay := time.Date(executeTime.Year(), executeTime.Month(), executeTime.Day(), 0, 0, 0, 0, executeTime.Location())
			periodStart = actualSettleDay
			periodEnd = actualSettleDay.Add(24*time.Hour - time.Millisecond)
			periodStr = actualSettleDay.Format("20060102")
		} else {
			periodStart = theoreticalSettleDay
			periodEnd = theoreticalSettleDay.Add(24*time.Hour - time.Millisecond)
			periodStr = theoreticalSettleDay.Format("20060102")
		}
	case protocol.T3:
		// T+3 三日后结算 - 交易完成三日后结算
		theoreticalSettleDay := time.Date(completedTime.Year(), completedTime.Month(), completedTime.Day()+3, 0, 0, 0, 0, completedTime.Location())

		// 如果理论结算日期已经过去，使用执行时间的日期
		if executeTime.After(theoreticalSettleDay.Add(24*time.Hour - time.Millisecond)) {
			actualSettleDay := time.Date(executeTime.Year(), executeTime.Month(), executeTime.Day(), 0, 0, 0, 0, executeTime.Location())
			periodStart = actualSettleDay
			periodEnd = actualSettleDay.Add(24*time.Hour - time.Millisecond)
			periodStr = actualSettleDay.Format("20060102")
		} else {
			periodStart = theoreticalSettleDay
			periodEnd = theoreticalSettleDay.Add(24*time.Hour - time.Millisecond)
			periodStr = theoreticalSettleDay.Format("20060102")
		}
	case protocol.W1:
		// W+1 下周结算 - 基于交易完成时间计算下周
		// 获取交易完成日期是周几（周一为1，周日为7）
		weekday := int(completedTime.Weekday())
		if weekday == 0 {
			weekday = 7 // 将周日从0改为7
		}

		// 计算理论下周一的日期
		daysToNextMonday := 7 - weekday + 1
		theoreticalNextMonday := completedTime.AddDate(0, 0, daysToNextMonday)
		theoreticalMondayStart := time.Date(theoreticalNextMonday.Year(), theoreticalNextMonday.Month(), theoreticalNextMonday.Day(), 0, 0, 0, 0, theoreticalNextMonday.Location())
		theoreticalSundayEnd := theoreticalMondayStart.AddDate(0, 0, 6).Add(24*time.Hour - time.Millisecond)

		// 如果理论结算周期已经过去，使用执行时间计算当前周期
		if executeTime.After(theoreticalSundayEnd) {
			// 基于执行时间计算当前周
			executeWeekday := int(executeTime.Weekday())
			if executeWeekday == 0 {
				executeWeekday = 7
			}
			daysToCurrentMonday := executeWeekday - 1
			currentMonday := executeTime.AddDate(0, 0, -daysToCurrentMonday)
			currentMondayStart := time.Date(currentMonday.Year(), currentMonday.Month(), currentMonday.Day(), 0, 0, 0, 0, currentMonday.Location())
			currentSundayEnd := currentMondayStart.AddDate(0, 0, 6).Add(24*time.Hour - time.Millisecond)

			periodStart = currentMondayStart
			periodEnd = currentSundayEnd

			// 计算当前周所在月的第几周
			firstDay := time.Date(currentMondayStart.Year(), currentMondayStart.Month(), 1, 0, 0, 0, 0, currentMondayStart.Location())
			firstWeekday := int(firstDay.Weekday())
			if firstWeekday == 0 {
				firstWeekday = 7
			}
			weekNum := (currentMondayStart.Day()-1+firstWeekday-1)/7 + 1
			periodStr = fmt.Sprintf("%s%02d", currentMondayStart.Format("200601"), weekNum)
		} else {
			// 使用理论结算周期
			periodStart = theoreticalMondayStart
			periodEnd = theoreticalSundayEnd

			firstDay := time.Date(theoreticalMondayStart.Year(), theoreticalMondayStart.Month(), 1, 0, 0, 0, 0, theoreticalMondayStart.Location())
			firstWeekday := int(firstDay.Weekday())
			if firstWeekday == 0 {
				firstWeekday = 7
			}
			weekNum := (theoreticalMondayStart.Day()-1+firstWeekday-1)/7 + 1
			periodStr = fmt.Sprintf("%s%02d", theoreticalMondayStart.Format("200601"), weekNum)
		}
	case protocol.M1:
		// M+1 下月结算 - 基于交易完成时间计算下月
		// 理论下个月1号00:00:00
		theoreticalNextMonth := time.Date(completedTime.Year(), completedTime.Month()+1, 1, 0, 0, 0, 0, completedTime.Location())
		theoreticalMonthEnd := theoreticalNextMonth.AddDate(0, 1, 0).Add(-time.Millisecond)

		// 如果理论结算月份已经过去，使用执行时间所在月份
		if executeTime.After(theoreticalMonthEnd) {
			actualMonth := time.Date(executeTime.Year(), executeTime.Month(), 1, 0, 0, 0, 0, executeTime.Location())
			actualMonthEnd := actualMonth.AddDate(0, 1, 0).Add(-time.Millisecond)

			periodStart = actualMonth
			periodEnd = actualMonthEnd
			periodStr = actualMonth.Format("200601")
		} else {
			// 使用理论结算月份
			periodStart = theoreticalNextMonth
			periodEnd = theoreticalMonthEnd
			periodStr = theoreticalNextMonth.Format("200601")
		}
	default:
		// 默认按天结算（T+0当天结算）
		theoreticalSettleDay := time.Date(completedTime.Year(), completedTime.Month(), completedTime.Day(), 0, 0, 0, 0, completedTime.Location())

		// 如果理论结算日期已经过去，使用执行时间的日期
		if executeTime.After(theoreticalSettleDay.Add(24*time.Hour - time.Millisecond)) {
			actualSettleDay := time.Date(executeTime.Year(), executeTime.Month(), executeTime.Day(), 0, 0, 0, 0, executeTime.Location())
			periodStart = actualSettleDay
			periodEnd = actualSettleDay.Add(24*time.Hour - time.Millisecond)
			periodStr = actualSettleDay.Format("20060102")
		} else {
			periodStart = theoreticalSettleDay
			periodEnd = theoreticalSettleDay.Add(24*time.Hour - time.Millisecond)
			periodStr = theoreticalSettleDay.Format("20060102")
		}
	}

	period = cast.ToInt64(periodStr)
	startAt = periodStart.UnixMilli()
	endAt = periodEnd.UnixMilli()

	return
}
func ListMerchantContractByMid(mid string) []*MerchantContract {
	var contracts []*MerchantContract
	if err := ReadDB.Where("mid = ? and status=?", mid, protocol.StatusActive).Find(&contracts).Error; err != nil {
		log.Get().Errorf("ListMerchantContractByMid failed, mid: %s, err: %v", mid, err)
		return nil
	}
	return contracts
}

// GetValidContractsAtTime 获取在指定时间有效的合同
func GetValidContractsAtTime(mid string, trxTime int64) *MerchantContract {
	if mid == "" {
		return nil
	}

	var contract *MerchantContract
	err := ReadDB.Where("mid = ? AND status = ? AND start_at <= ? AND (expired_at = 0 OR expired_at >= ?)",
		mid, protocol.StatusActive, trxTime, trxTime).Order("id desc").First(&contract).Error
	if err != nil {
		log.Get().Errorf("getValidContractsAtTime failed, mid: %s, trxTime: %d, err: %v", mid, trxTime, err)
		return nil
	}

	return contract
}
