package models

import (
	"github.com/shopspring/decimal"
)

// MerchantWithdraw 提现记录表
type MerchantWithdraw struct {
	ID        uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TrxID     string `json:"trx_id" gorm:"column:trx_id;type:varchar(64);uniqueIndex"`
	Mid       string `json:"mid" gorm:"column:mid;type:varchar(64);index"`
	AccountID string `json:"account_id" gorm:"column:account_id;type:varchar(64);index"`
	*MerchantWithdrawValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type MerchantWithdrawValues struct {
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

func (MerchantWithdraw) TableName() string {
	return "t_merchant_withdraws"
}

// Getter methods for MerchantWithdrawValues
func (mwv *MerchantWithdrawValues) GetStatus() string {
	if mwv.Status == nil {
		return ""
	}
	return *mwv.Status
}

func (mwv *MerchantWithdrawValues) GetCcy() string {
	if mwv.Ccy == nil {
		return ""
	}
	return *mwv.Ccy
}

func (mwv *MerchantWithdrawValues) GetAmount() decimal.Decimal {
	if mwv.Amount == nil {
		return decimal.Zero
	}
	return *mwv.Amount
}

func (mwv *MerchantWithdrawValues) GetFee() decimal.Decimal {
	if mwv.Fee == nil {
		return decimal.Zero
	}
	return *mwv.Fee
}

func (mwv *MerchantWithdrawValues) GetChannelCode() string {
	if mwv.ChannelCode == nil {
		return ""
	}
	return *mwv.ChannelCode
}

func (mwv *MerchantWithdrawValues) GetNotifyURL() string {
	if mwv.NotifyURL == nil {
		return ""
	}
	return *mwv.NotifyURL
}

func (mwv *MerchantWithdrawValues) GetCountry() string {
	if mwv.Country == nil {
		return ""
	}
	return *mwv.Country
}

func (mwv *MerchantWithdrawValues) GetCanceledAt() int64 {
	if mwv.CanceledAt == nil {
		return 0
	}
	return *mwv.CanceledAt
}

func (mwv *MerchantWithdrawValues) GetCompletedAt() int64 {
	if mwv.CompletedAt == nil {
		return 0
	}
	return *mwv.CompletedAt
}

func (mwv *MerchantWithdrawValues) GetExpiredAt() int64 {
	if mwv.ExpiredAt == nil {
		return 0
	}
	return *mwv.ExpiredAt
}

func (mwv *MerchantWithdrawValues) GetConfirmedAt() int64 {
	if mwv.ConfirmedAt == nil {
		return 0
	}
	return *mwv.ConfirmedAt
}

// Setter methods for MerchantWithdrawValues (支持链式调用)
func (mwv *MerchantWithdrawValues) SetStatus(status string) *MerchantWithdrawValues {
	mwv.Status = &status
	return mwv
}

func (mwv *MerchantWithdrawValues) SetCcy(ccy string) *MerchantWithdrawValues {
	mwv.Ccy = &ccy
	return mwv
}

func (mwv *MerchantWithdrawValues) SetAmount(amount decimal.Decimal) *MerchantWithdrawValues {
	mwv.Amount = &amount
	return mwv
}

func (mwv *MerchantWithdrawValues) SetFee(fee decimal.Decimal) *MerchantWithdrawValues {
	mwv.Fee = &fee
	return mwv
}

func (mwv *MerchantWithdrawValues) SetChannelCode(channelCode string) *MerchantWithdrawValues {
	mwv.ChannelCode = &channelCode
	return mwv
}

func (mwv *MerchantWithdrawValues) SetNotifyURL(notifyURL string) *MerchantWithdrawValues {
	mwv.NotifyURL = &notifyURL
	return mwv
}

func (mwv *MerchantWithdrawValues) SetCountry(country string) *MerchantWithdrawValues {
	mwv.Country = &country
	return mwv
}

func (mwv *MerchantWithdrawValues) SetCanceledAt(canceledAt int64) *MerchantWithdrawValues {
	mwv.CanceledAt = &canceledAt
	return mwv
}

func (mwv *MerchantWithdrawValues) SetCompletedAt(completedAt int64) *MerchantWithdrawValues {
	mwv.CompletedAt = &completedAt
	return mwv
}

func (mwv *MerchantWithdrawValues) SetExpiredAt(expiredAt int64) *MerchantWithdrawValues {
	mwv.ExpiredAt = &expiredAt
	return mwv
}

func (mwv *MerchantWithdrawValues) SetConfirmedAt(confirmedAt int64) *MerchantWithdrawValues {
	mwv.ConfirmedAt = &confirmedAt
	return mwv
}

// SetValues 为MerchantWithdraw设置MerchantWithdrawValues
func (mw *MerchantWithdraw) SetValues(values *MerchantWithdrawValues) *MerchantWithdraw {
	if values == nil {
		return mw
	}

	if mw.MerchantWithdrawValues == nil {
		mw.MerchantWithdrawValues = &MerchantWithdrawValues{}
	}

	if values.Status != nil {
		mw.MerchantWithdrawValues.SetStatus(*values.Status)
	}
	if values.Ccy != nil {
		mw.MerchantWithdrawValues.SetCcy(*values.Ccy)
	}
	if values.Amount != nil {
		mw.MerchantWithdrawValues.SetAmount(*values.Amount)
	}
	if values.Fee != nil {
		mw.MerchantWithdrawValues.SetFee(*values.Fee)
	}
	if values.ChannelCode != nil {
		mw.MerchantWithdrawValues.SetChannelCode(*values.ChannelCode)
	}
	if values.NotifyURL != nil {
		mw.MerchantWithdrawValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.Country != nil {
		mw.MerchantWithdrawValues.SetCountry(*values.Country)
	}
	if values.CanceledAt != nil {
		mw.MerchantWithdrawValues.SetCanceledAt(*values.CanceledAt)
	}
	if values.CompletedAt != nil {
		mw.MerchantWithdrawValues.SetCompletedAt(*values.CompletedAt)
	}
	if values.ExpiredAt != nil {
		mw.MerchantWithdrawValues.SetExpiredAt(*values.ExpiredAt)
	}
	if values.ConfirmedAt != nil {
		mw.MerchantWithdrawValues.SetConfirmedAt(*values.ConfirmedAt)
	}

	return mw
}
