package models

import (
	"github.com/shopspring/decimal"
)

// CashierWithdraw 提现记录表
type CashierWithdraw struct {
	ID        uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Tid       string `json:"tid" gorm:"column:tid;type:varchar(32);index"`
	TrxID     string `json:"trx_id" gorm:"column:trx_id;type:varchar(64);uniqueIndex"`
	AccountID string `json:"account_id" gorm:"column:account_id;type:varchar(64);index"`
	*CashierWithdrawValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type CashierWithdrawValues struct {
	Status      *string          `json:"status" gorm:"column:status;type:varchar(16);index;default:'pending'"`
	Ccy         *string          `json:"ccy" gorm:"column:ccy;type:varchar(16)"`
	Amount      *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Fee         *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(36,18);default:0"`
	ChannelCode *string          `json:"channel_code" gorm:"column:channel_code;type:varchar(32)"`
	NotifyURL   *string          `json:"notify_url" gorm:"column:notify_url;type:varchar(512)"`
	Country     *string          `json:"country" gorm:"column:country;type:varchar(8)"`
	CanceledAt  *int64           `json:"canceled_at" gorm:"column:canceled_at"`
	CompletedAt *int64           `json:"completed_at" gorm:"column:completed_at"`
	ExpiredAt   *int64           `json:"expired_at" gorm:"column:expired_at"`
	ConfirmedAt *int64           `json:"confirmed_at" gorm:"column:confirmed_at"`
}

func (CashierWithdraw) TableName() string {
	return "t_cashier_withdraws"
}

// Getter methods for CashierWithdrawValues
func (cwv *CashierWithdrawValues) GetStatus() string {
	if cwv.Status == nil {
		return ""
	}
	return *cwv.Status
}

func (cwv *CashierWithdrawValues) GetCcy() string {
	if cwv.Ccy == nil {
		return ""
	}
	return *cwv.Ccy
}

func (cwv *CashierWithdrawValues) GetAmount() decimal.Decimal {
	if cwv.Amount == nil {
		return decimal.Zero
	}
	return *cwv.Amount
}

func (cwv *CashierWithdrawValues) GetFee() decimal.Decimal {
	if cwv.Fee == nil {
		return decimal.Zero
	}
	return *cwv.Fee
}

func (cwv *CashierWithdrawValues) GetChannelCode() string {
	if cwv.ChannelCode == nil {
		return ""
	}
	return *cwv.ChannelCode
}

func (cwv *CashierWithdrawValues) GetNotifyURL() string {
	if cwv.NotifyURL == nil {
		return ""
	}
	return *cwv.NotifyURL
}

func (cwv *CashierWithdrawValues) GetCountry() string {
	if cwv.Country == nil {
		return ""
	}
	return *cwv.Country
}

func (cwv *CashierWithdrawValues) GetCanceledAt() int64 {
	if cwv.CanceledAt == nil {
		return 0
	}
	return *cwv.CanceledAt
}

func (cwv *CashierWithdrawValues) GetCompletedAt() int64 {
	if cwv.CompletedAt == nil {
		return 0
	}
	return *cwv.CompletedAt
}

func (cwv *CashierWithdrawValues) GetExpiredAt() int64 {
	if cwv.ExpiredAt == nil {
		return 0
	}
	return *cwv.ExpiredAt
}

func (cwv *CashierWithdrawValues) GetConfirmedAt() int64 {
	if cwv.ConfirmedAt == nil {
		return 0
	}
	return *cwv.ConfirmedAt
}

// Setter methods for CashierWithdrawValues (支持链式调用)
func (cwv *CashierWithdrawValues) SetStatus(status string) *CashierWithdrawValues {
	cwv.Status = &status
	return cwv
}

func (cwv *CashierWithdrawValues) SetCcy(ccy string) *CashierWithdrawValues {
	cwv.Ccy = &ccy
	return cwv
}

func (cwv *CashierWithdrawValues) SetAmount(amount decimal.Decimal) *CashierWithdrawValues {
	cwv.Amount = &amount
	return cwv
}

func (cwv *CashierWithdrawValues) SetFee(fee decimal.Decimal) *CashierWithdrawValues {
	cwv.Fee = &fee
	return cwv
}

func (cwv *CashierWithdrawValues) SetChannelCode(channelCode string) *CashierWithdrawValues {
	cwv.ChannelCode = &channelCode
	return cwv
}

func (cwv *CashierWithdrawValues) SetNotifyURL(notifyURL string) *CashierWithdrawValues {
	cwv.NotifyURL = &notifyURL
	return cwv
}

func (cwv *CashierWithdrawValues) SetCountry(country string) *CashierWithdrawValues {
	cwv.Country = &country
	return cwv
}

func (cwv *CashierWithdrawValues) SetCanceledAt(canceledAt int64) *CashierWithdrawValues {
	cwv.CanceledAt = &canceledAt
	return cwv
}

func (cwv *CashierWithdrawValues) SetCompletedAt(completedAt int64) *CashierWithdrawValues {
	cwv.CompletedAt = &completedAt
	return cwv
}

func (cwv *CashierWithdrawValues) SetExpiredAt(expiredAt int64) *CashierWithdrawValues {
	cwv.ExpiredAt = &expiredAt
	return cwv
}

func (cwv *CashierWithdrawValues) SetConfirmedAt(confirmedAt int64) *CashierWithdrawValues {
	cwv.ConfirmedAt = &confirmedAt
	return cwv
}

// SetValues 为CashierWithdraw设置CashierWithdrawValues
func (cw *CashierWithdraw) SetValues(values *CashierWithdrawValues) *CashierWithdraw {
	if values == nil {
		return cw
	}

	if cw.CashierWithdrawValues == nil {
		cw.CashierWithdrawValues = &CashierWithdrawValues{}
	}

	if values.Status != nil {
		cw.CashierWithdrawValues.SetStatus(*values.Status)
	}
	if values.Ccy != nil {
		cw.CashierWithdrawValues.SetCcy(*values.Ccy)
	}
	if values.Amount != nil {
		cw.CashierWithdrawValues.SetAmount(*values.Amount)
	}
	if values.Fee != nil {
		cw.CashierWithdrawValues.SetFee(*values.Fee)
	}
	if values.ChannelCode != nil {
		cw.CashierWithdrawValues.SetChannelCode(*values.ChannelCode)
	}
	if values.NotifyURL != nil {
		cw.CashierWithdrawValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.Country != nil {
		cw.CashierWithdrawValues.SetCountry(*values.Country)
	}
	if values.CanceledAt != nil {
		cw.CashierWithdrawValues.SetCanceledAt(*values.CanceledAt)
	}
	if values.CompletedAt != nil {
		cw.CashierWithdrawValues.SetCompletedAt(*values.CompletedAt)
	}
	if values.ExpiredAt != nil {
		cw.CashierWithdrawValues.SetExpiredAt(*values.ExpiredAt)
	}
	if values.ConfirmedAt != nil {
		cw.CashierWithdrawValues.SetConfirmedAt(*values.ConfirmedAt)
	}

	return cw
}
