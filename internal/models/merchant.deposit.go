package models

import "github.com/shopspring/decimal"

// MerchantDeposit 充值记录表
type MerchantDeposit struct {
	ID        uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TrxID     string `json:"trx_id" gorm:"column:trx_id;type:varchar(64);uniqueIndex"`
	Mid       string `json:"mid" gorm:"column:mid;type:varchar(64);index"`
	AccountID string `json:"account_id" gorm:"column:account_id;type:varchar(64);index"`
	*MerchantDepositValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type MerchantDepositValues struct {
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

func (MerchantDeposit) TableName() string {
	return "t_merchant_deposits"
}

// Getter methods for MerchantDepositValues
func (mdv *MerchantDepositValues) GetStatus() string {
	if mdv.Status == nil {
		return ""
	}
	return *mdv.Status
}

func (mdv *MerchantDepositValues) GetCcy() string {
	if mdv.Ccy == nil {
		return ""
	}
	return *mdv.Ccy
}

func (mdv *MerchantDepositValues) GetAmount() decimal.Decimal {
	if mdv.Amount == nil {
		return decimal.Zero
	}
	return *mdv.Amount
}

func (mdv *MerchantDepositValues) GetFee() decimal.Decimal {
	if mdv.Fee == nil {
		return decimal.Zero
	}
	return *mdv.Fee
}

func (mdv *MerchantDepositValues) GetChannelCode() string {
	if mdv.ChannelCode == nil {
		return ""
	}
	return *mdv.ChannelCode
}

func (mdv *MerchantDepositValues) GetNotifyURL() string {
	if mdv.NotifyURL == nil {
		return ""
	}
	return *mdv.NotifyURL
}

func (mdv *MerchantDepositValues) GetCountry() string {
	if mdv.Country == nil {
		return ""
	}
	return *mdv.Country
}

func (mdv *MerchantDepositValues) GetCanceledAt() int64 {
	if mdv.CanceledAt == nil {
		return 0
	}
	return *mdv.CanceledAt
}

func (mdv *MerchantDepositValues) GetCompletedAt() int64 {
	if mdv.CompletedAt == nil {
		return 0
	}
	return *mdv.CompletedAt
}

func (mdv *MerchantDepositValues) GetExpiredAt() int64 {
	if mdv.ExpiredAt == nil {
		return 0
	}
	return *mdv.ExpiredAt
}

func (mdv *MerchantDepositValues) GetConfirmedAt() int64 {
	if mdv.ConfirmedAt == nil {
		return 0
	}
	return *mdv.ConfirmedAt
}

// Setter methods for MerchantDepositValues (支持链式调用)
func (mdv *MerchantDepositValues) SetStatus(status string) *MerchantDepositValues {
	mdv.Status = &status
	return mdv
}

func (mdv *MerchantDepositValues) SetCcy(ccy string) *MerchantDepositValues {
	mdv.Ccy = &ccy
	return mdv
}

func (mdv *MerchantDepositValues) SetAmount(amount decimal.Decimal) *MerchantDepositValues {
	mdv.Amount = &amount
	return mdv
}

func (mdv *MerchantDepositValues) SetFee(fee decimal.Decimal) *MerchantDepositValues {
	mdv.Fee = &fee
	return mdv
}

func (mdv *MerchantDepositValues) SetChannelCode(channelCode string) *MerchantDepositValues {
	mdv.ChannelCode = &channelCode
	return mdv
}

func (mdv *MerchantDepositValues) SetNotifyURL(notifyURL string) *MerchantDepositValues {
	mdv.NotifyURL = &notifyURL
	return mdv
}

func (mdv *MerchantDepositValues) SetCountry(country string) *MerchantDepositValues {
	mdv.Country = &country
	return mdv
}

func (mdv *MerchantDepositValues) SetCanceledAt(canceledAt int64) *MerchantDepositValues {
	mdv.CanceledAt = &canceledAt
	return mdv
}

func (mdv *MerchantDepositValues) SetCompletedAt(completedAt int64) *MerchantDepositValues {
	mdv.CompletedAt = &completedAt
	return mdv
}

func (mdv *MerchantDepositValues) SetExpiredAt(expiredAt int64) *MerchantDepositValues {
	mdv.ExpiredAt = &expiredAt
	return mdv
}

func (mdv *MerchantDepositValues) SetConfirmedAt(confirmedAt int64) *MerchantDepositValues {
	mdv.ConfirmedAt = &confirmedAt
	return mdv
}

// SetValues 为MerchantDeposit设置MerchantDepositValues
func (md *MerchantDeposit) SetValues(values *MerchantDepositValues) *MerchantDeposit {
	if values == nil {
		return md
	}

	if md.MerchantDepositValues == nil {
		md.MerchantDepositValues = &MerchantDepositValues{}
	}

	if values.Status != nil {
		md.MerchantDepositValues.SetStatus(*values.Status)
	}
	if values.Ccy != nil {
		md.MerchantDepositValues.SetCcy(*values.Ccy)
	}
	if values.Amount != nil {
		md.MerchantDepositValues.SetAmount(*values.Amount)
	}
	if values.Fee != nil {
		md.MerchantDepositValues.SetFee(*values.Fee)
	}
	if values.ChannelCode != nil {
		md.MerchantDepositValues.SetChannelCode(*values.ChannelCode)
	}
	if values.NotifyURL != nil {
		md.MerchantDepositValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.Country != nil {
		md.MerchantDepositValues.SetCountry(*values.Country)
	}
	if values.CanceledAt != nil {
		md.MerchantDepositValues.SetCanceledAt(*values.CanceledAt)
	}
	if values.CompletedAt != nil {
		md.MerchantDepositValues.SetCompletedAt(*values.CompletedAt)
	}
	if values.ExpiredAt != nil {
		md.MerchantDepositValues.SetExpiredAt(*values.ExpiredAt)
	}
	if values.ConfirmedAt != nil {
		md.MerchantDepositValues.SetConfirmedAt(*values.ConfirmedAt)
	}

	return md
}
