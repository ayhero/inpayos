package models

import "github.com/shopspring/decimal"

// CashierDeposit 充值记录表
type CashierDeposit struct {
	ID        uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TrxID     string `json:"trx_id" gorm:"column:trx_id;type:varchar(64);uniqueIndex"`
	Tid       string `json:"tid" gorm:"column:tid;type:varchar(32);index"`
	AccountID string `json:"account_id" gorm:"column:account_id;type:varchar(64);index"`
	*CashierDepositValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type CashierDepositValues struct {
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

func (CashierDeposit) TableName() string {
	return "t_cashier_deposits"
}

// Getter methods for CashierDepositValues
func (cdv *CashierDepositValues) GetStatus() string {
	if cdv.Status == nil {
		return ""
	}
	return *cdv.Status
}

func (cdv *CashierDepositValues) GetCcy() string {
	if cdv.Ccy == nil {
		return ""
	}
	return *cdv.Ccy
}

func (cdv *CashierDepositValues) GetAmount() decimal.Decimal {
	if cdv.Amount == nil {
		return decimal.Zero
	}
	return *cdv.Amount
}

func (cdv *CashierDepositValues) GetFee() decimal.Decimal {
	if cdv.Fee == nil {
		return decimal.Zero
	}
	return *cdv.Fee
}

func (cdv *CashierDepositValues) GetChannelCode() string {
	if cdv.ChannelCode == nil {
		return ""
	}
	return *cdv.ChannelCode
}

func (cdv *CashierDepositValues) GetNotifyURL() string {
	if cdv.NotifyURL == nil {
		return ""
	}
	return *cdv.NotifyURL
}

func (cdv *CashierDepositValues) GetCountry() string {
	if cdv.Country == nil {
		return ""
	}
	return *cdv.Country
}

func (cdv *CashierDepositValues) GetCanceledAt() int64 {
	if cdv.CanceledAt == nil {
		return 0
	}
	return *cdv.CanceledAt
}

func (cdv *CashierDepositValues) GetCompletedAt() int64 {
	if cdv.CompletedAt == nil {
		return 0
	}
	return *cdv.CompletedAt
}

func (cdv *CashierDepositValues) GetExpiredAt() int64 {
	if cdv.ExpiredAt == nil {
		return 0
	}
	return *cdv.ExpiredAt
}

func (cdv *CashierDepositValues) GetConfirmedAt() int64 {
	if cdv.ConfirmedAt == nil {
		return 0
	}
	return *cdv.ConfirmedAt
}

// Setter methods for CashierDepositValues (支持链式调用)
func (cdv *CashierDepositValues) SetStatus(status string) *CashierDepositValues {
	cdv.Status = &status
	return cdv
}

func (cdv *CashierDepositValues) SetCcy(ccy string) *CashierDepositValues {
	cdv.Ccy = &ccy
	return cdv
}

func (cdv *CashierDepositValues) SetAmount(amount decimal.Decimal) *CashierDepositValues {
	cdv.Amount = &amount
	return cdv
}

func (cdv *CashierDepositValues) SetFee(fee decimal.Decimal) *CashierDepositValues {
	cdv.Fee = &fee
	return cdv
}

func (cdv *CashierDepositValues) SetChannelCode(channelCode string) *CashierDepositValues {
	cdv.ChannelCode = &channelCode
	return cdv
}

func (cdv *CashierDepositValues) SetNotifyURL(notifyURL string) *CashierDepositValues {
	cdv.NotifyURL = &notifyURL
	return cdv
}

func (cdv *CashierDepositValues) SetCountry(country string) *CashierDepositValues {
	cdv.Country = &country
	return cdv
}

func (cdv *CashierDepositValues) SetCanceledAt(canceledAt int64) *CashierDepositValues {
	cdv.CanceledAt = &canceledAt
	return cdv
}

func (cdv *CashierDepositValues) SetCompletedAt(completedAt int64) *CashierDepositValues {
	cdv.CompletedAt = &completedAt
	return cdv
}

func (cdv *CashierDepositValues) SetExpiredAt(expiredAt int64) *CashierDepositValues {
	cdv.ExpiredAt = &expiredAt
	return cdv
}

func (cdv *CashierDepositValues) SetConfirmedAt(confirmedAt int64) *CashierDepositValues {
	cdv.ConfirmedAt = &confirmedAt
	return cdv
}

// SetValues 为CashierDeposit设置CashierDepositValues
func (cd *CashierDeposit) SetValues(values *CashierDepositValues) *CashierDeposit {
	if values == nil {
		return cd
	}

	if cd.CashierDepositValues == nil {
		cd.CashierDepositValues = &CashierDepositValues{}
	}

	if values.Status != nil {
		cd.CashierDepositValues.SetStatus(*values.Status)
	}
	if values.Ccy != nil {
		cd.CashierDepositValues.SetCcy(*values.Ccy)
	}
	if values.Amount != nil {
		cd.CashierDepositValues.SetAmount(*values.Amount)
	}
	if values.Fee != nil {
		cd.CashierDepositValues.SetFee(*values.Fee)
	}
	if values.ChannelCode != nil {
		cd.CashierDepositValues.SetChannelCode(*values.ChannelCode)
	}
	if values.NotifyURL != nil {
		cd.CashierDepositValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.Country != nil {
		cd.CashierDepositValues.SetCountry(*values.Country)
	}
	if values.CanceledAt != nil {
		cd.CashierDepositValues.SetCanceledAt(*values.CanceledAt)
	}
	if values.CompletedAt != nil {
		cd.CashierDepositValues.SetCompletedAt(*values.CompletedAt)
	}
	if values.ExpiredAt != nil {
		cd.CashierDepositValues.SetExpiredAt(*values.ExpiredAt)
	}
	if values.ConfirmedAt != nil {
		cd.CashierDepositValues.SetConfirmedAt(*values.ConfirmedAt)
	}

	return cd
}
