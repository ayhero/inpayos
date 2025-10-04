package models

import (
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
)

// Payout 代付记录表
type Payout struct {
	ID     uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TrxID  string `json:"transaction_id" gorm:"column:transaction_id;type:varchar(64);uniqueIndex"`
	Mid    string `json:"mid" gorm:"column:mid;type:varchar(32);index"`
	UserID string `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	ReqID  string `json:"req_id" gorm:"column:req_id;type:varchar(64);index"`
	Salt   string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*PayoutValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type PayoutValues struct {
	Status        *string          `json:"status" gorm:"column:status;type:varchar(16);index;default:'pending'"`
	Amount        *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Fee           *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(36,18);default:0"`
	Ccy           *string          `json:"ccy" gorm:"column:ccy;type:varchar(16)"`
	Country       *string          `json:"country" gorm:"column:country;type:varchar(8)"`
	ChannelCode   *string          `json:"channel_code" gorm:"column:channel_code;type:varchar(32)"`
	PaymentMethod *string          `json:"payment_method" gorm:"column:payment_method;type:varchar(32)"`
	RecipientInfo *string          `json:"recipient_info" gorm:"column:recipient_info;type:json"` // 收款方信息
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

func (Payout) TableName() string {
	return "t_payouts"
}

// GetStatus returns the Status value
func (pov *PayoutValues) GetStatus() string {
	if pov.Status == nil {
		return ""
	}
	return *pov.Status
}

// GetAmount returns the Amount value
func (pov *PayoutValues) GetAmount() decimal.Decimal {
	if pov.Amount == nil {
		return decimal.Zero
	}
	return *pov.Amount
}

// GetFee returns the Fee value
func (pov *PayoutValues) GetFee() decimal.Decimal {
	if pov.Fee == nil {
		return decimal.Zero
	}
	return *pov.Fee
}

// GetCcy returns the Ccy value
func (pov *PayoutValues) GetCcy() string {
	if pov.Ccy == nil {
		return ""
	}
	return *pov.Ccy
}

// GetChannelCode returns the ChannelCode value
func (pov *PayoutValues) GetChannelCode() string {
	if pov.ChannelCode == nil {
		return ""
	}
	return *pov.ChannelCode
}

// GetPaymentMethod returns the PaymentMethod value
func (pov *PayoutValues) GetPaymentMethod() string {
	if pov.PaymentMethod == nil {
		return ""
	}
	return *pov.PaymentMethod
}

// GetRecipientInfo returns the RecipientInfo value
func (pov *PayoutValues) GetRecipientInfo() string {
	if pov.RecipientInfo == nil {
		return ""
	}
	return *pov.RecipientInfo
}

// GetNotifyURL returns the NotifyURL value
func (pov *PayoutValues) GetNotifyURL() string {
	if pov.NotifyURL == nil {
		return ""
	}
	return *pov.NotifyURL
}

// GetCountry returns the Country value
func (pov *PayoutValues) GetCountry() string {
	if pov.Country == nil {
		return ""
	}
	return *pov.Country
}

// GetExpiredAt returns the ExpiredAt value
func (pov *PayoutValues) GetExpiredAt() int64 {
	if pov.ExpiredAt == nil {
		return 0
	}
	return *pov.ExpiredAt
}

// GetCanceledAt returns the CanceledAt value
func (pov *PayoutValues) GetCanceledAt() int64 {
	if pov.CanceledAt == nil {
		return 0
	}
	return *pov.CanceledAt
}

// SetStatus sets the Status value
func (pov *PayoutValues) SetStatus(value string) *PayoutValues {
	pov.Status = &value
	return pov
}

// SetAmount sets the Amount value
func (pov *PayoutValues) SetAmount(value decimal.Decimal) *PayoutValues {
	pov.Amount = &value
	return pov
}

// SetFee sets the Fee value
func (pov *PayoutValues) SetFee(value decimal.Decimal) *PayoutValues {
	pov.Fee = &value
	return pov
}

// SetCcy sets the Ccy value
func (pov *PayoutValues) SetCcy(value string) *PayoutValues {
	pov.Ccy = &value
	return pov
}

// SetChannelCode sets the ChannelCode value
func (pov *PayoutValues) SetChannelCode(value string) *PayoutValues {
	pov.ChannelCode = &value
	return pov
}

// SetPaymentMethod sets the PaymentMethod value
func (pov *PayoutValues) SetPaymentMethod(value string) *PayoutValues {
	pov.PaymentMethod = &value
	return pov
}

// SetRecipientInfo sets the RecipientInfo value
func (pov *PayoutValues) SetRecipientInfo(value string) *PayoutValues {
	pov.RecipientInfo = &value
	return pov
}

// SetNotifyURL sets the NotifyURL value
func (pov *PayoutValues) SetNotifyURL(value string) *PayoutValues {
	pov.NotifyURL = &value
	return pov
}

// SetCountry sets the Country value
func (pov *PayoutValues) SetCountry(value string) *PayoutValues {
	pov.Country = &value
	return pov
}

// SetExpiredAt sets the ExpiredAt value
func (pov *PayoutValues) SetExpiredAt(value int64) *PayoutValues {
	pov.ExpiredAt = &value
	return pov
}

// SetCanceledAt sets the CanceledAt value
func (pov *PayoutValues) SetCanceledAt(value int64) *PayoutValues {
	pov.CanceledAt = &value
	return pov
}

// SetValues sets multiple PayoutValues fields at once
func (p *Payout) SetValues(values *PayoutValues) *Payout {
	if values == nil {
		return p
	}

	if p.PayoutValues == nil {
		p.PayoutValues = &PayoutValues{}
	}

	if values.Status != nil {
		p.PayoutValues.SetStatus(*values.Status)
	}
	if values.Amount != nil {
		p.PayoutValues.SetAmount(*values.Amount)
	}
	if values.Fee != nil {
		p.PayoutValues.SetFee(*values.Fee)
	}
	if values.Ccy != nil {
		p.PayoutValues.SetCcy(*values.Ccy)
	}
	if values.ChannelCode != nil {
		p.PayoutValues.SetChannelCode(*values.ChannelCode)
	}
	if values.PaymentMethod != nil {
		p.PayoutValues.SetPaymentMethod(*values.PaymentMethod)
	}
	if values.RecipientInfo != nil {
		p.PayoutValues.SetRecipientInfo(*values.RecipientInfo)
	}
	if values.NotifyURL != nil {
		p.PayoutValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.Country != nil {
		p.PayoutValues.SetCountry(*values.Country)
	}
	if values.ExpiredAt != nil {
		p.PayoutValues.SetExpiredAt(*values.ExpiredAt)
	}
	if values.CanceledAt != nil {
		p.PayoutValues.SetCanceledAt(*values.CanceledAt)
	}

	return p
}

// ToTransaction converts Payout to Transaction
func (p *Payout) ToTransaction() *Transaction {
	if p == nil {
		return nil
	}

	transaction := &Transaction{
		ID:      p.ID,
		TrxID:   p.TrxID,
		Mid:     p.Mid,
		UserID:  p.UserID,
		ReqID:   p.ReqID,
		TrxType: protocol.TrxTypePayout, // Set transaction type to payout
		TransactionValues: &TransactionValues{
			Status:        p.PayoutValues.Status,
			Amount:        p.PayoutValues.Amount,
			Fee:           p.PayoutValues.Fee,
			Ccy:           p.PayoutValues.Ccy,
			ChannelCode:   p.PayoutValues.ChannelCode,
			PaymentMethod: p.PayoutValues.PaymentMethod,
			NotifyURL:     p.PayoutValues.NotifyURL,
			ReturnURL:     nil,                          // Payout doesn't have ReturnURL
			NotifyStatus:  nil,                          // Payout doesn't have NotifyStatus
			NotifyTimes:   nil,                          // Payout doesn't have NotifyTimes
			OriTrxID:      nil,                          // Payout doesn't have OriTrxID
			Metadata:      p.PayoutValues.RecipientInfo, // Map RecipientInfo to Metadata
			Remark:        nil,                          // Payout doesn't have Remark
			ExpiredAt:     p.PayoutValues.ExpiredAt,
			ConfirmedAt:   nil, // Payout doesn't have ConfirmedAt
			CanceledAt:    p.PayoutValues.CanceledAt,
			UpdatedAt:     p.UpdatedAt, // Payout has UpdatedAt in main struct
		},
		CreatedAt: p.CreatedAt,
	}

	return transaction
}
