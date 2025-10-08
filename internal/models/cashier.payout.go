package models

import (
	"github.com/shopspring/decimal"
)

// CashierPayout 代付记录表
type CashierPayout struct {
	ID          uint64           `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Tid         string           `json:"tid" gorm:"column:tid"`
	ReqID       string           `json:"req_id" gorm:"column:req_id;type:varchar(64);index"`
	TrxID       string           `json:"transaction_id" gorm:"column:transaction_id;type:varchar(64);uniqueIndex"`
	TrxType     string           `json:"trx_type" gorm:"column:trx_type;type:varchar(16);index;default:'payin'"`
	OriTrxID    string           `json:"ori_trx_id" gorm:"column:ori_trx_id;index;<-:create"`
	OriReqID    string           `json:"ori_req_id" gorm:"column:ori_req_id;index;<-:create"`
	OriFlowNo   string           `json:"ori_flow_no" gorm:"column:ori_flow_no"`
	TrxMethod   string           `json:"trx_method" gorm:"column:trx_method;<-:create"`
	Ccy         string           `json:"ccy" gorm:"column:ccy;<-:create"`
	Amount      *decimal.Decimal `json:"amount" gorm:"column:amount;<-:create"`
	UsdAmount   *decimal.Decimal `json:"usd_amount" gorm:"column:usd_amount;<-:create"`
	AccountNo   string           `json:"account_no" gorm:"column:account_no;<-:create"`
	AccountName string           `json:"account_name" gorm:"column:account_name;<-:create"`
	AccountType string           `json:"account_type" gorm:"column:account_type;<-:create"`
	BankCode    string           `json:"bank_code" gorm:"column:bank_code;<-:create"`
	BankName    string           `json:"bank_name" gorm:"column:bank_name;<-:create"`
	*CashierPayoutValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type CashierPayoutValues struct {
	FlowNo    string           `json:"flow_no" gorm:"column:flow_no"`
	CashierID *string          `json:"cashier_id" gorm:"column:cashier_id;type:varchar(32);index"`
	Status    *string          `json:"status" gorm:"column:status;type:varchar(16);index;default:'pending'"`
	Fee       *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(36,18);default:0"`
	NotifyURL *string          `json:"notify_url" gorm:"column:notify_url;type:varchar(512)"`
	ReturnURL *string          `json:"return_url" gorm:"column:return_url;type:varchar(512)"`
	OriTrxID  *string          `json:"ori_trx_id" gorm:"column:ori_trx_id;type:varchar(64)"` // 原交易ID(退款使用)
	Metadata  *string          `json:"metadata" gorm:"column:metadata;type:json"`
	Remark    *string          `json:"remark" gorm:"column:remark;type:varchar(512)"`

	SettleID           *string `json:"settle_id" gorm:"column:settle_id;type:varchar(64)"`         // 结算ID
	SettleStatus       *string `json:"settle_status" gorm:"column:settle_status;type:varchar(16)"` // 结算状态
	SettledAt          *int64  `json:"settled_at" gorm:"column:settled_at"`                        // 结算时间
	ExpiredAt          *int64  `json:"expired_at" gorm:"column:expired_at"`
	ConfirmedAt        *int64  `json:"confirmed_at" gorm:"column:confirmed_at"`
	CanceledAt         *int64  `json:"canceled_at" gorm:"column:canceled_at"`
	CancelFailedResult *string `json:"cancel_failed_result" gorm:"column:cancel_failed_result;type:varchar(512)"`
}

func (CashierPayout) TableName() string {
	return "t_cashier_payouts"
}

// GetStatus returns the Status value
func (pov *CashierPayoutValues) GetStatus() string {
	if pov.Status == nil {
		return ""
	}
	return *pov.Status
}

// GetCashierID returns the CashierID value
func (pov *CashierPayoutValues) GetCashierID() string {
	if pov.CashierID == nil {
		return ""
	}
	return *pov.CashierID
}

// GetFee returns the Fee value
func (pov *CashierPayoutValues) GetFee() decimal.Decimal {
	if pov.Fee == nil {
		return decimal.Zero
	}
	return *pov.Fee
}

// GetNotifyURL returns the NotifyURL value
func (pov *CashierPayoutValues) GetNotifyURL() string {
	if pov.NotifyURL == nil {
		return ""
	}
	return *pov.NotifyURL
}

// GetReturnURL returns the ReturnURL value
func (pov *CashierPayoutValues) GetReturnURL() string {
	if pov.ReturnURL == nil {
		return ""
	}
	return *pov.ReturnURL
}

// GetOriTrxID returns the OriTrxID value
func (pov *CashierPayoutValues) GetOriTrxID() string {
	if pov.OriTrxID == nil {
		return ""
	}
	return *pov.OriTrxID
}

// GetMetadata returns the Metadata value
func (pov *CashierPayoutValues) GetMetadata() string {
	if pov.Metadata == nil {
		return ""
	}
	return *pov.Metadata
}

// GetRemark returns the Remark value
func (pov *CashierPayoutValues) GetRemark() string {
	if pov.Remark == nil {
		return ""
	}
	return *pov.Remark
}

// GetSettleID returns the SettleID value
func (pov *CashierPayoutValues) GetSettleID() string {
	if pov.SettleID == nil {
		return ""
	}
	return *pov.SettleID
}

// GetSettleStatus returns the SettleStatus value
func (pov *CashierPayoutValues) GetSettleStatus() string {
	if pov.SettleStatus == nil {
		return ""
	}
	return *pov.SettleStatus
}

// GetSettledAt returns the SettledAt value
func (pov *CashierPayoutValues) GetSettledAt() int64 {
	if pov.SettledAt == nil {
		return 0
	}
	return *pov.SettledAt
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

// GetCancelFailedResult returns the CancelFailedResult value
func (pov *CashierPayoutValues) GetCancelFailedResult() string {
	if pov.CancelFailedResult == nil {
		return ""
	}
	return *pov.CancelFailedResult
}

// SetStatus sets the Status value
func (pov *CashierPayoutValues) SetStatus(value string) *CashierPayoutValues {
	pov.Status = &value
	return pov
}

// SetCashierID sets the CashierID value
func (pov *CashierPayoutValues) SetCashierID(value string) *CashierPayoutValues {
	pov.CashierID = &value
	return pov
}

// SetFee sets the Fee value
func (pov *CashierPayoutValues) SetFee(value decimal.Decimal) *CashierPayoutValues {
	pov.Fee = &value
	return pov
}

// SetReturnURL sets the ReturnURL value
func (pov *CashierPayoutValues) SetReturnURL(value string) *CashierPayoutValues {
	pov.ReturnURL = &value
	return pov
}

// SetOriTrxID sets the OriTrxID value
func (pov *CashierPayoutValues) SetOriTrxID(value string) *CashierPayoutValues {
	pov.OriTrxID = &value
	return pov
}

// SetMetadata sets the Metadata value
func (pov *CashierPayoutValues) SetMetadata(value string) *CashierPayoutValues {
	pov.Metadata = &value
	return pov
}

// SetRemark sets the Remark value
func (pov *CashierPayoutValues) SetRemark(value string) *CashierPayoutValues {
	pov.Remark = &value
	return pov
}

// SetSettleID sets the SettleID value
func (pov *CashierPayoutValues) SetSettleID(value string) *CashierPayoutValues {
	pov.SettleID = &value
	return pov
}

// SetSettleStatus sets the SettleStatus value
func (pov *CashierPayoutValues) SetSettleStatus(value string) *CashierPayoutValues {
	pov.SettleStatus = &value
	return pov
}

// SetSettledAt sets the SettledAt value
func (pov *CashierPayoutValues) SetSettledAt(value int64) *CashierPayoutValues {
	pov.SettledAt = &value
	return pov
}

// SetNotifyURL sets the NotifyURL value
func (pov *CashierPayoutValues) SetNotifyURL(value string) *CashierPayoutValues {
	pov.NotifyURL = &value
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

// SetCancelFailedResult sets the CancelFailedResult value
func (pov *CashierPayoutValues) SetCancelFailedResult(value string) *CashierPayoutValues {
	pov.CancelFailedResult = &value
	return pov
}

// SetConfirmedAt sets the ConfirmedAt value
func (pov *CashierPayoutValues) SetConfirmedAt(value int64) *CashierPayoutValues {
	pov.ConfirmedAt = &value
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

	if values.CashierID != nil {
		p.CashierPayoutValues.SetCashierID(*values.CashierID)
	}
	if values.Status != nil {
		p.CashierPayoutValues.SetStatus(*values.Status)
	}
	if values.Fee != nil {
		p.CashierPayoutValues.SetFee(*values.Fee)
	}
	if values.NotifyURL != nil {
		p.CashierPayoutValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.ReturnURL != nil {
		p.CashierPayoutValues.SetReturnURL(*values.ReturnURL)
	}
	if values.OriTrxID != nil {
		p.CashierPayoutValues.SetOriTrxID(*values.OriTrxID)
	}
	if values.Metadata != nil {
		p.CashierPayoutValues.SetMetadata(*values.Metadata)
	}
	if values.Remark != nil {
		p.CashierPayoutValues.SetRemark(*values.Remark)
	}
	if values.SettleID != nil {
		p.CashierPayoutValues.SetSettleID(*values.SettleID)
	}
	if values.SettleStatus != nil {
		p.CashierPayoutValues.SetSettleStatus(*values.SettleStatus)
	}
	if values.SettledAt != nil {
		p.CashierPayoutValues.SetSettledAt(*values.SettledAt)
	}
	if values.ExpiredAt != nil {
		p.CashierPayoutValues.SetExpiredAt(*values.ExpiredAt)
	}
	if values.ConfirmedAt != nil {
		p.CashierPayoutValues.SetConfirmedAt(*values.ConfirmedAt)
	}
	if values.CanceledAt != nil {
		p.CashierPayoutValues.SetCanceledAt(*values.CanceledAt)
	}
	if values.CancelFailedResult != nil {
		p.CashierPayoutValues.SetCancelFailedResult(*values.CancelFailedResult)
	}

	return p
}

// ToTransaction converts CashierPayout to Transaction
func (p *CashierPayout) ToTransaction() *Transaction {
	if p == nil {
		return nil
	}

	transaction := &Transaction{
		ID:          int64(p.ID), // Convert uint64 to int64
		Tid:         p.Tid,
		CashierID:   p.CashierPayoutValues.GetCashierID(),
		TrxID:       p.TrxID,
		ReqID:       p.ReqID,
		OriTrxID:    p.OriTrxID,
		OriReqID:    p.OriReqID,
		OriFlowNo:   p.OriFlowNo,
		TrxMethod:   p.TrxMethod,
		Ccy:         p.Ccy,
		Amount:      p.Amount,
		UsdAmount:   p.UsdAmount,
		AccountNo:   p.AccountNo,
		AccountName: p.AccountName,
		AccountType: p.AccountType,
		BankCode:    p.BankCode,
		BankName:    p.BankName,
		ReturnURL:   "", // CashierPayout doesn't have ReturnURL in main struct
		TransactionValues: &TransactionValues{
			Status:             p.CashierPayoutValues.Status,
			NotifyURL:          p.CashierPayoutValues.NotifyURL,
			Remark:             p.CashierPayoutValues.Remark,
			FlowNo:             &p.CashierPayoutValues.FlowNo,
			SettleID:           p.CashierPayoutValues.SettleID,
			SettleStatus:       p.CashierPayoutValues.SettleStatus,
			SettledAt:          p.CashierPayoutValues.SettledAt,
			ExpiredAt:          p.CashierPayoutValues.ExpiredAt,
			ConfirmedAt:        p.CashierPayoutValues.ConfirmedAt,
			CanceledAt:         p.CashierPayoutValues.CanceledAt,
			CancelFailedResult: p.CashierPayoutValues.CancelFailedResult,
			FeeAmount:          p.CashierPayoutValues.Fee,
		},
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	return transaction
}
