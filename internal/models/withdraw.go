package models

import (
	"inpayos/internal/utils"

	"github.com/shopspring/decimal"
)

// Withdraw 提现记录表
type Withdraw struct {
	ID        uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TrxID     string `json:"trx_id" gorm:"column:trx_id;type:varchar(64);uniqueIndex"`
	Sid       string `json:"sid" gorm:"column:sid;type:varchar(32);index"`
	SType     string `json:"s_type" gorm:"column:s_type;type:varchar(32);index"` // service类型，如 "merchant", "cashier"
	AccountID string `json:"account_id" gorm:"column:account_id;type:varchar(64);index"`
	*WithdrawValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type WithdrawValues struct {
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

func (Withdraw) TableName() string {
	return "t_withdraws"
}

// NewWithdraw 创建新的提现记录
func NewWithdraw() *Withdraw {
	return &Withdraw{
		TrxID:          utils.GenerateWithdrawID(),
		WithdrawValues: &WithdrawValues{},
	}
}

// Getter methods for WithdrawValues
func (wv *WithdrawValues) GetStatus() string {
	if wv.Status == nil {
		return ""
	}
	return *wv.Status
}

func (wv *WithdrawValues) GetCcy() string {
	if wv.Ccy == nil {
		return ""
	}
	return *wv.Ccy
}

func (wv *WithdrawValues) GetAmount() decimal.Decimal {
	if wv.Amount == nil {
		return decimal.Zero
	}
	return *wv.Amount
}

func (wv *WithdrawValues) GetFee() decimal.Decimal {
	if wv.Fee == nil {
		return decimal.Zero
	}
	return *wv.Fee
}

func (wv *WithdrawValues) GetChannelCode() string {
	if wv.ChannelCode == nil {
		return ""
	}
	return *wv.ChannelCode
}

func (wv *WithdrawValues) GetNotifyURL() string {
	if wv.NotifyURL == nil {
		return ""
	}
	return *wv.NotifyURL
}

func (wv *WithdrawValues) GetCountry() string {
	if wv.Country == nil {
		return ""
	}
	return *wv.Country
}

func (wv *WithdrawValues) GetCanceledAt() int64 {
	if wv.CanceledAt == nil {
		return 0
	}
	return *wv.CanceledAt
}

func (wv *WithdrawValues) GetCompletedAt() int64 {
	if wv.CompletedAt == nil {
		return 0
	}
	return *wv.CompletedAt
}

func (wv *WithdrawValues) GetExpiredAt() int64 {
	if wv.ExpiredAt == nil {
		return 0
	}
	return *wv.ExpiredAt
}

func (wv *WithdrawValues) GetConfirmedAt() int64 {
	if wv.ConfirmedAt == nil {
		return 0
	}
	return *wv.ConfirmedAt
}

// Setter methods for WithdrawValues (支持链式调用)
func (wv *WithdrawValues) SetStatus(status string) *WithdrawValues {
	wv.Status = &status
	return wv
}

func (wv *WithdrawValues) SetCcy(ccy string) *WithdrawValues {
	wv.Ccy = &ccy
	return wv
}

func (wv *WithdrawValues) SetAmount(amount decimal.Decimal) *WithdrawValues {
	wv.Amount = &amount
	return wv
}

func (wv *WithdrawValues) SetFee(fee decimal.Decimal) *WithdrawValues {
	wv.Fee = &fee
	return wv
}

func (wv *WithdrawValues) SetChannelCode(channelCode string) *WithdrawValues {
	wv.ChannelCode = &channelCode
	return wv
}

func (wv *WithdrawValues) SetNotifyURL(notifyURL string) *WithdrawValues {
	wv.NotifyURL = &notifyURL
	return wv
}

func (wv *WithdrawValues) SetCountry(country string) *WithdrawValues {
	wv.Country = &country
	return wv
}

func (wv *WithdrawValues) SetCanceledAt(canceledAt int64) *WithdrawValues {
	wv.CanceledAt = &canceledAt
	return wv
}

func (wv *WithdrawValues) SetCompletedAt(completedAt int64) *WithdrawValues {
	wv.CompletedAt = &completedAt
	return wv
}

func (wv *WithdrawValues) SetExpiredAt(expiredAt int64) *WithdrawValues {
	wv.ExpiredAt = &expiredAt
	return wv
}

func (wv *WithdrawValues) SetConfirmedAt(confirmedAt int64) *WithdrawValues {
	wv.ConfirmedAt = &confirmedAt
	return wv
}

// SetValues 为Withdraw设置WithdrawValues
func (w *Withdraw) SetValues(values *WithdrawValues) *Withdraw {
	if values == nil {
		return w
	}

	if w.WithdrawValues == nil {
		w.WithdrawValues = &WithdrawValues{}
	}

	if values.Status != nil {
		w.WithdrawValues.SetStatus(*values.Status)
	}
	if values.Ccy != nil {
		w.WithdrawValues.SetCcy(*values.Ccy)
	}
	if values.Amount != nil {
		w.WithdrawValues.SetAmount(*values.Amount)
	}
	if values.Fee != nil {
		w.WithdrawValues.SetFee(*values.Fee)
	}
	if values.ChannelCode != nil {
		w.WithdrawValues.SetChannelCode(*values.ChannelCode)
	}
	if values.NotifyURL != nil {
		w.WithdrawValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.Country != nil {
		w.WithdrawValues.SetCountry(*values.Country)
	}
	if values.CanceledAt != nil {
		w.WithdrawValues.SetCanceledAt(*values.CanceledAt)
	}
	if values.CompletedAt != nil {
		w.WithdrawValues.SetCompletedAt(*values.CompletedAt)
	}
	if values.ExpiredAt != nil {
		w.WithdrawValues.SetExpiredAt(*values.ExpiredAt)
	}
	if values.ConfirmedAt != nil {
		w.WithdrawValues.SetConfirmedAt(*values.ConfirmedAt)
	}

	return w
}
