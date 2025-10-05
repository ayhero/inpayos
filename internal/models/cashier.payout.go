package models

import (
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
)

// CashierPayout 代付记录表
type CashierPayout struct {
	ID        uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Tid       string `json:"tid" gorm:"column:tid"`
	CashierID string `json:"cashier_id" gorm:"column:cashier_id;type:varchar(32);index"`
	ReqID     string `json:"req_id" gorm:"column:req_id;type:varchar(64);index"`
	TrxID     string `json:"transaction_id" gorm:"column:transaction_id;type:varchar(64);uniqueIndex"`
	*CashierPayoutValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type CashierPayoutValues struct {
	Status        *string          `json:"status" gorm:"column:status;type:varchar(16);index;default:'pending'"`
	Amount        *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Fee           *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(36,18);default:0"`
	Ccy           *string          `json:"ccy" gorm:"column:ccy;type:varchar(16)"`
	Country       *string          `json:"country" gorm:"column:country;type:varchar(8)"`
	ChannelCode   *string          `json:"channel_code" gorm:"column:channel_code;type:varchar(32)"`
	PaymentMethod *string          `json:"payment_method" gorm:"column:payment_method;type:varchar(32)"`
	NotifyURL     *string          `json:"notify_url" gorm:"column:notify_url;type:varchar(512)"`
	OriTrxID      *string          `json:"ori_trx_id" gorm:"column:ori_trx_id;type:varchar(64)"` // 原交易ID(退款使用)
	Metadata      *string          `json:"metadata" gorm:"column:metadata;type:json"`
	Remark        *string          `json:"remark" gorm:"column:remark;type:varchar(512)"`
	UsdAmount     *decimal.Decimal `json:"usd_amount" gorm:"column:usd_amount;type:decimal(36,18)"`    // USD金额
	SettleID      *string          `json:"settle_id" gorm:"column:settle_id;type:varchar(64)"`         // 结算ID
	SettleStatus  *string          `json:"settle_status" gorm:"column:settle_status;type:varchar(16)"` // 结算状态
	SettledAt     *int64           `json:"settled_at" gorm:"column:settled_at"`                        // 结算时间
	ExpiredAt     *int64           `json:"expired_at" gorm:"column:expired_at"`
	CanceledAt    *int64           `json:"canceled_at" gorm:"column:canceled_at"`
}

func (CashierPayout) TableName() string {
	return "t_payouts"
}

// GetStatus returns the Status value
func (pov *CashierPayoutValues) GetStatus() string {
	if pov.Status == nil {
		return ""
	}
	return *pov.Status
}

// GetAmount returns the Amount value
func (pov *CashierPayoutValues) GetAmount() decimal.Decimal {
	if pov.Amount == nil {
		return decimal.Zero
	}
	return *pov.Amount
}

// GetFee returns the Fee value
func (pov *CashierPayoutValues) GetFee() decimal.Decimal {
	if pov.Fee == nil {
		return decimal.Zero
	}
	return *pov.Fee
}

// GetCcy returns the Ccy value
func (pov *CashierPayoutValues) GetCcy() string {
	if pov.Ccy == nil {
		return ""
	}
	return *pov.Ccy
}

// GetChannelCode returns the ChannelCode value
func (pov *CashierPayoutValues) GetChannelCode() string {
	if pov.ChannelCode == nil {
		return ""
	}
	return *pov.ChannelCode
}

// GetPaymentMethod returns the PaymentMethod value
func (pov *CashierPayoutValues) GetPaymentMethod() string {
	if pov.PaymentMethod == nil {
		return ""
	}
	return *pov.PaymentMethod
}

// GetNotifyURL returns the NotifyURL value
func (pov *CashierPayoutValues) GetNotifyURL() string {
	if pov.NotifyURL == nil {
		return ""
	}
	return *pov.NotifyURL
}

// GetCountry returns the Country value
func (pov *CashierPayoutValues) GetCountry() string {
	if pov.Country == nil {
		return ""
	}
	return *pov.Country
}

// GetExpiredAt returns the ExpiredAt value
func (pov *CashierPayoutValues) GetExpiredAt() int64 {
	if pov.ExpiredAt == nil {
		return 0
	}
	return *pov.ExpiredAt
}

// GetCanceledAt returns the CanceledAt value
func (pov *CashierPayoutValues) GetCanceledAt() int64 {
	if pov.CanceledAt == nil {
		return 0
	}
	return *pov.CanceledAt
}

// SetStatus sets the Status value
func (pov *CashierPayoutValues) SetStatus(value string) *CashierPayoutValues {
	pov.Status = &value
	return pov
}

// SetAmount sets the Amount value
func (pov *CashierPayoutValues) SetAmount(value decimal.Decimal) *CashierPayoutValues {
	pov.Amount = &value
	return pov
}

// SetFee sets the Fee value
func (pov *CashierPayoutValues) SetFee(value decimal.Decimal) *CashierPayoutValues {
	pov.Fee = &value
	return pov
}

// SetCcy sets the Ccy value
func (pov *CashierPayoutValues) SetCcy(value string) *CashierPayoutValues {
	pov.Ccy = &value
	return pov
}

// SetChannelCode sets the ChannelCode value
func (pov *CashierPayoutValues) SetChannelCode(value string) *CashierPayoutValues {
	pov.ChannelCode = &value
	return pov
}

// SetPaymentMethod sets the PaymentMethod value
func (pov *CashierPayoutValues) SetPaymentMethod(value string) *CashierPayoutValues {
	pov.PaymentMethod = &value
	return pov
}

// SetNotifyURL sets the NotifyURL value
func (pov *CashierPayoutValues) SetNotifyURL(value string) *CashierPayoutValues {
	pov.NotifyURL = &value
	return pov
}

// SetCountry sets the Country value
func (pov *CashierPayoutValues) SetCountry(value string) *CashierPayoutValues {
	pov.Country = &value
	return pov
}

// SetExpiredAt sets the ExpiredAt value
func (pov *CashierPayoutValues) SetExpiredAt(value int64) *CashierPayoutValues {
	pov.ExpiredAt = &value
	return pov
}

// SetCanceledAt sets the CanceledAt value
func (pov *CashierPayoutValues) SetCanceledAt(value int64) *CashierPayoutValues {
	pov.CanceledAt = &value
	return pov
}

// SetValues sets multiple CashierPayoutValues fields at once
func (p *CashierPayout) SetValues(values *CashierPayoutValues) *CashierPayout {
	if values == nil {
		return p
	}

	if p.CashierPayoutValues == nil {
		p.CashierPayoutValues = &CashierPayoutValues{}
	}

	if values.Status != nil {
		p.CashierPayoutValues.SetStatus(*values.Status)
	}
	if values.Amount != nil {
		p.CashierPayoutValues.SetAmount(*values.Amount)
	}
	if values.Fee != nil {
		p.CashierPayoutValues.SetFee(*values.Fee)
	}
	if values.Ccy != nil {
		p.CashierPayoutValues.SetCcy(*values.Ccy)
	}
	if values.ChannelCode != nil {
		p.CashierPayoutValues.SetChannelCode(*values.ChannelCode)
	}
	if values.PaymentMethod != nil {
		p.CashierPayoutValues.SetPaymentMethod(*values.PaymentMethod)
	}
	if values.NotifyURL != nil {
		p.CashierPayoutValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.Country != nil {
		p.CashierPayoutValues.SetCountry(*values.Country)
	}
	if values.ExpiredAt != nil {
		p.CashierPayoutValues.SetExpiredAt(*values.ExpiredAt)
	}
	if values.CanceledAt != nil {
		p.CashierPayoutValues.SetCanceledAt(*values.CanceledAt)
	}

	return p
}

// ToTransaction converts CashierPayout to Transaction
func (p *CashierPayout) ToTransaction() *Transaction {
	if p == nil {
		return nil
	}

	transaction := &Transaction{
		ID:        p.ID,
		Tid:       p.Tid,
		CashierID: p.CashierID,
		TrxID:     p.TrxID,
		ReqID:     p.ReqID,
		TrxType:   protocol.TrxTypePayout, // Set transaction type to payout
		TransactionValues: &TransactionValues{
			Status:        p.CashierPayoutValues.Status,
			Amount:        p.CashierPayoutValues.Amount,
			Fee:           p.CashierPayoutValues.Fee,
			Ccy:           p.CashierPayoutValues.Ccy,
			ChannelCode:   p.CashierPayoutValues.ChannelCode,
			PaymentMethod: p.CashierPayoutValues.PaymentMethod,
			NotifyURL:     p.CashierPayoutValues.NotifyURL,
			ReturnURL:     nil, // CashierPayout doesn't have ReturnURL
			NotifyStatus:  nil, // CashierPayout doesn't have NotifyStatus
			NotifyTimes:   nil, // CashierPayout doesn't have NotifyTimes
			OriTrxID:      nil, // CashierPayout doesn't have OriTrxID
			Remark:        nil, // CashierPayout doesn't have Remark
			ExpiredAt:     p.CashierPayoutValues.ExpiredAt,
			ConfirmedAt:   nil, // CashierPayout doesn't have ConfirmedAt
			CanceledAt:    p.CashierPayoutValues.CanceledAt,
			UpdatedAt:     p.UpdatedAt, // CashierPayout has UpdatedAt in main struct
		},
		CreatedAt: p.CreatedAt,
	}

	return transaction
}
