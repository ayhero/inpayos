package models

import (
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
)

// CashierPayin 代收记录表
type CashierPayin struct {
	ID        uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Tid       string `json:"tid" gorm:"column:tid"`
	CashierID string `json:"cashier_id" gorm:"column:cashier_id;type:varchar(32);index"`
	ReqID     string `json:"req_id" gorm:"column:req_id;type:varchar(64);index"`
	TrxID     string `json:"transaction_id" gorm:"column:transaction_id;type:varchar(64);uniqueIndex"`
	*CashierPayinValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type CashierPayinValues struct {
	Status        *string          `json:"status" gorm:"column:status;type:varchar(16);index;default:'pending'"`
	Country       *string          `json:"country" gorm:"column:country;type:varchar(8)"`
	Ccy           *string          `json:"ccy" gorm:"column:ccy;type:varchar(16)"`
	Amount        *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Fee           *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(36,18);default:0"`
	ChannelCode   *string          `json:"channel_code" gorm:"column:channel_code;type:varchar(32)"`
	PaymentMethod *string          `json:"payment_method" gorm:"column:payment_method;type:varchar(32)"`
	NotifyURL     *string          `json:"notify_url" gorm:"column:notify_url;type:varchar(512)"`
	ReturnURL     *string          `json:"return_url" gorm:"column:return_url;type:varchar(512)"`
	OriTrxID      *string          `json:"ori_trx_id" gorm:"column:ori_trx_id;type:varchar(64)"` // 原交易ID(退款使用)
	Metadata      *string          `json:"metadata" gorm:"column:metadata;type:json"`
	Remark        *string          `json:"remark" gorm:"column:remark;type:varchar(512)"`
	UsdAmount     *decimal.Decimal `json:"usd_amount" gorm:"column:usd_amount;type:decimal(36,18)"`    // USD金额
	SettleID      *string          `json:"settle_id" gorm:"column:settle_id;type:varchar(64)"`         // 结算ID
	SettleStatus  *string          `json:"settle_status" gorm:"column:settle_status;type:varchar(16)"` // 结算状态
	SettledAt     *int64           `json:"settled_at" gorm:"column:settled_at"`                        // 结算时间
	ExpiredAt     *int64           `json:"expired_at" gorm:"column:expired_at"`
	ConfirmedAt   *int64           `json:"confirmed_at" gorm:"column:confirmed_at"`
	CanceledAt    *int64           `json:"canceled_at" gorm:"column:canceled_at"`
	UpdatedAt     int64            `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (CashierPayin) TableName() string {
	return "t_cashier_payins"
}

// GetStatus returns the Status value
func (pv *CashierPayinValues) GetStatus() string {
	if pv.Status == nil {
		return ""
	}
	return *pv.Status
}

// GetCountry returns the Country value
func (pv *CashierPayinValues) GetCountry() string {
	if pv.Country == nil {
		return ""
	}
	return *pv.Country
}

// GetCcy returns the Ccy value
func (pv *CashierPayinValues) GetCcy() string {
	if pv.Ccy == nil {
		return ""
	}
	return *pv.Ccy
}

// GetAmount returns the Amount value
func (pv *CashierPayinValues) GetAmount() decimal.Decimal {
	if pv.Amount == nil {
		return decimal.Zero
	}
	return *pv.Amount
}

// GetFee returns the Fee value
func (pv *CashierPayinValues) GetFee() decimal.Decimal {
	if pv.Fee == nil {
		return decimal.Zero
	}
	return *pv.Fee
}

// GetChannelCode returns the ChannelCode value
func (pv *CashierPayinValues) GetChannelCode() string {
	if pv.ChannelCode == nil {
		return ""
	}
	return *pv.ChannelCode
}

// GetPaymentMethod returns the PaymentMethod value
func (pv *CashierPayinValues) GetPaymentMethod() string {
	if pv.PaymentMethod == nil {
		return ""
	}
	return *pv.PaymentMethod
}

// GetNotifyURL returns the NotifyURL value
func (pv *CashierPayinValues) GetNotifyURL() string {
	if pv.NotifyURL == nil {
		return ""
	}
	return *pv.NotifyURL
}

// GetReturnURL returns the ReturnURL value
func (pv *CashierPayinValues) GetReturnURL() string {
	if pv.ReturnURL == nil {
		return ""
	}
	return *pv.ReturnURL
}

// GetExpiredAt returns the ExpiredAt value
func (pv *CashierPayinValues) GetExpiredAt() int64 {
	if pv.ExpiredAt == nil {
		return 0
	}
	return *pv.ExpiredAt
}

// GetConfirmedAt returns the ConfirmedAt value
func (pv *CashierPayinValues) GetConfirmedAt() int64 {
	if pv.ConfirmedAt == nil {
		return 0
	}
	return *pv.ConfirmedAt
}

// GetCanceledAt returns the CanceledAt value
func (pv *CashierPayinValues) GetCanceledAt() int64 {
	if pv.CanceledAt == nil {
		return 0
	}
	return *pv.CanceledAt
}

// GetUpdatedAt returns the UpdatedAt value
func (pv *CashierPayinValues) GetUpdatedAt() int64 {
	return pv.UpdatedAt
}

// SetStatus sets the Status value
func (pv *CashierPayinValues) SetStatus(value string) *CashierPayinValues {
	pv.Status = &value
	return pv
}

// SetCountry sets the Country value
func (pv *CashierPayinValues) SetCountry(value string) *CashierPayinValues {
	pv.Country = &value
	return pv
}

// SetCcy sets the Ccy value
func (pv *CashierPayinValues) SetCcy(value string) *CashierPayinValues {
	pv.Ccy = &value
	return pv
}

// SetAmount sets the Amount value
func (pv *CashierPayinValues) SetAmount(value decimal.Decimal) *CashierPayinValues {
	pv.Amount = &value
	return pv
}

// SetFee sets the Fee value
func (pv *CashierPayinValues) SetFee(value decimal.Decimal) *CashierPayinValues {
	pv.Fee = &value
	return pv
}

// SetChannelCode sets the ChannelCode value
func (pv *CashierPayinValues) SetChannelCode(value string) *CashierPayinValues {
	pv.ChannelCode = &value
	return pv
}

// SetPaymentMethod sets the PaymentMethod value
func (pv *CashierPayinValues) SetPaymentMethod(value string) *CashierPayinValues {
	pv.PaymentMethod = &value
	return pv
}

// SetNotifyURL sets the NotifyURL value
func (pv *CashierPayinValues) SetNotifyURL(value string) *CashierPayinValues {
	pv.NotifyURL = &value
	return pv
}

// SetReturnURL sets the ReturnURL value
func (pv *CashierPayinValues) SetReturnURL(value string) *CashierPayinValues {
	pv.ReturnURL = &value
	return pv
}

// SetExpiredAt sets the ExpiredAt value
func (pv *CashierPayinValues) SetExpiredAt(value int64) *CashierPayinValues {
	pv.ExpiredAt = &value
	return pv
}

// SetConfirmedAt sets the ConfirmedAt value
func (pv *CashierPayinValues) SetConfirmedAt(value int64) *CashierPayinValues {
	pv.ConfirmedAt = &value
	return pv
}

// SetCanceledAt sets the CanceledAt value
func (pv *CashierPayinValues) SetCanceledAt(value int64) *CashierPayinValues {
	pv.CanceledAt = &value
	return pv
}

// SetUpdatedAt sets the UpdatedAt value
func (pv *CashierPayinValues) SetUpdatedAt(value int64) *CashierPayinValues {
	pv.UpdatedAt = value
	return pv
}

// SetValues sets multiple CashierPayinValues fields at once
func (p *CashierPayin) SetValues(values *CashierPayinValues) *CashierPayin {
	if values == nil {
		return p
	}

	if p.CashierPayinValues == nil {
		p.CashierPayinValues = &CashierPayinValues{}
	}
	if values.Status != nil {
		p.CashierPayinValues.SetStatus(*values.Status)
	}
	if values.Country != nil {
		p.CashierPayinValues.SetCountry(*values.Country)
	}
	if values.Ccy != nil {
		p.CashierPayinValues.SetCcy(*values.Ccy)
	}
	if values.Amount != nil {
		p.CashierPayinValues.SetAmount(*values.Amount)
	}
	if values.Fee != nil {
		p.CashierPayinValues.SetFee(*values.Fee)
	}
	if values.ChannelCode != nil {
		p.CashierPayinValues.SetChannelCode(*values.ChannelCode)
	}
	if values.PaymentMethod != nil {
		p.CashierPayinValues.SetPaymentMethod(*values.PaymentMethod)
	}
	if values.NotifyURL != nil {
		p.CashierPayinValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.ReturnURL != nil {
		p.CashierPayinValues.SetReturnURL(*values.ReturnURL)
	}
	if values.ExpiredAt != nil {
		p.CashierPayinValues.SetExpiredAt(*values.ExpiredAt)
	}
	if values.ConfirmedAt != nil {
		p.CashierPayinValues.SetConfirmedAt(*values.ConfirmedAt)
	}
	if values.CanceledAt != nil {
		p.CashierPayinValues.SetCanceledAt(*values.CanceledAt)
	}
	// UpdatedAt is not a pointer, so we always set it
	p.CashierPayinValues.SetUpdatedAt(values.UpdatedAt)

	return p
}

// ToTransaction converts CashierPayin to Transaction
func (p *CashierPayin) ToTransaction() *Transaction {
	if p == nil {
		return nil
	}

	transaction := &Transaction{
		ID:        p.ID,
		Tid:       p.Tid,
		CashierID: p.CashierID,
		TrxID:     p.TrxID,
		ReqID:     p.ReqID,
		TrxType:   protocol.TrxTypePayin, // Set transaction type to payin
		TransactionValues: &TransactionValues{
			Status:        p.CashierPayinValues.Status,
			Amount:        p.CashierPayinValues.Amount,
			Fee:           p.CashierPayinValues.Fee,
			Ccy:           p.CashierPayinValues.Ccy,
			ChannelCode:   p.CashierPayinValues.ChannelCode,
			PaymentMethod: p.CashierPayinValues.PaymentMethod,
			NotifyURL:     p.CashierPayinValues.NotifyURL,
			ReturnURL:     p.CashierPayinValues.ReturnURL,
			NotifyStatus:  nil, // CashierPayin doesn't have NotifyStatus
			NotifyTimes:   nil, // CashierPayin doesn't have NotifyTimes
			OriTrxID:      nil, // CashierPayin doesn't have OriTrxID
			Metadata:      nil, // CashierPayin doesn't have Metadata
			Remark:        nil, // CashierPayin doesn't have Remark
			ExpiredAt:     p.CashierPayinValues.ExpiredAt,
			ConfirmedAt:   p.CashierPayinValues.ConfirmedAt,
			CanceledAt:    p.CashierPayinValues.CanceledAt,
			UpdatedAt:     p.CashierPayinValues.UpdatedAt,
		},
		CreatedAt: p.CreatedAt,
	}

	return transaction
}
