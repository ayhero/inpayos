package models

import (
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
)

// CashierPayin 代收记录表
type CashierPayin struct {
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
	*CashierPayinValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type CashierPayinValues struct {
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

// GetCashierID returns the CashierID value
func (pv *CashierPayinValues) GetCashierID() string {
	if pv.CashierID == nil {
		return ""
	}
	return *pv.CashierID
}

// GetFee returns the Fee value
func (pv *CashierPayinValues) GetFee() decimal.Decimal {
	if pv.Fee == nil {
		return decimal.Zero
	}
	return *pv.Fee
}

// GetOriTrxID returns the OriTrxID value
func (pv *CashierPayinValues) GetOriTrxID() string {
	if pv.OriTrxID == nil {
		return ""
	}
	return *pv.OriTrxID
}

// GetMetadata returns the Metadata value
func (pv *CashierPayinValues) GetMetadata() string {
	if pv.Metadata == nil {
		return ""
	}
	return *pv.Metadata
}

// GetRemark returns the Remark value
func (pv *CashierPayinValues) GetRemark() string {
	if pv.Remark == nil {
		return ""
	}
	return *pv.Remark
}

// GetSettleID returns the SettleID value
func (pv *CashierPayinValues) GetSettleID() string {
	if pv.SettleID == nil {
		return ""
	}
	return *pv.SettleID
}

// GetSettleStatus returns the SettleStatus value
func (pv *CashierPayinValues) GetSettleStatus() string {
	if pv.SettleStatus == nil {
		return ""
	}
	return *pv.SettleStatus
}

// GetSettledAt returns the SettledAt value
func (pv *CashierPayinValues) GetSettledAt() int64 {
	if pv.SettledAt == nil {
		return 0
	}
	return *pv.SettledAt
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

// GetCancelFailedResult returns the CancelFailedResult value
func (pv *CashierPayinValues) GetCancelFailedResult() string {
	if pv.CancelFailedResult == nil {
		return ""
	}
	return *pv.CancelFailedResult
}

// SetStatus sets the Status value
func (pv *CashierPayinValues) SetStatus(value string) *CashierPayinValues {
	pv.Status = &value
	return pv
}

// SetCashierID sets the CashierID value
func (pv *CashierPayinValues) SetCashierID(value string) *CashierPayinValues {
	pv.CashierID = &value
	return pv
}

// SetFee sets the Fee value
func (pv *CashierPayinValues) SetFee(value decimal.Decimal) *CashierPayinValues {
	pv.Fee = &value
	return pv
}

// SetOriTrxID sets the OriTrxID value
func (pv *CashierPayinValues) SetOriTrxID(value string) *CashierPayinValues {
	pv.OriTrxID = &value
	return pv
}

// SetMetadata sets the Metadata value
func (pv *CashierPayinValues) SetMetadata(value string) *CashierPayinValues {
	pv.Metadata = &value
	return pv
}

// SetRemark sets the Remark value
func (pv *CashierPayinValues) SetRemark(value string) *CashierPayinValues {
	pv.Remark = &value
	return pv
}

// SetSettleID sets the SettleID value
func (pv *CashierPayinValues) SetSettleID(value string) *CashierPayinValues {
	pv.SettleID = &value
	return pv
}

// SetSettleStatus sets the SettleStatus value
func (pv *CashierPayinValues) SetSettleStatus(value string) *CashierPayinValues {
	pv.SettleStatus = &value
	return pv
}

// SetSettledAt sets the SettledAt value
func (pv *CashierPayinValues) SetSettledAt(value int64) *CashierPayinValues {
	pv.SettledAt = &value
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

// SetCancelFailedResult sets the CancelFailedResult value
func (pv *CashierPayinValues) SetCancelFailedResult(value string) *CashierPayinValues {
	pv.CancelFailedResult = &value
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
	if values.CashierID != nil {
		p.CashierPayinValues.SetCashierID(*values.CashierID)
	}
	if values.Status != nil {
		p.CashierPayinValues.SetStatus(*values.Status)
	}
	if values.Fee != nil {
		p.CashierPayinValues.SetFee(*values.Fee)
	}
	if values.NotifyURL != nil {
		p.CashierPayinValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.ReturnURL != nil {
		p.CashierPayinValues.SetReturnURL(*values.ReturnURL)
	}
	if values.OriTrxID != nil {
		p.CashierPayinValues.SetOriTrxID(*values.OriTrxID)
	}
	if values.Metadata != nil {
		p.CashierPayinValues.SetMetadata(*values.Metadata)
	}
	if values.Remark != nil {
		p.CashierPayinValues.SetRemark(*values.Remark)
	}
	if values.SettleID != nil {
		p.CashierPayinValues.SetSettleID(*values.SettleID)
	}
	if values.SettleStatus != nil {
		p.CashierPayinValues.SetSettleStatus(*values.SettleStatus)
	}
	if values.SettledAt != nil {
		p.CashierPayinValues.SetSettledAt(*values.SettledAt)
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
	if values.CancelFailedResult != nil {
		p.CashierPayinValues.SetCancelFailedResult(*values.CancelFailedResult)
	}

	return p
}

// ToTransaction converts CashierPayin to Transaction
func (p *CashierPayin) ToTransaction() *Transaction {
	if p == nil {
		return nil
	}

	transaction := &Transaction{
		ID:          int64(p.ID), // Convert uint64 to int64
		Tid:         p.Tid,
		TrxType:     protocol.TrxTypePayin,
		CashierID:   p.CashierPayinValues.GetCashierID(),
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
		ReturnURL:   "", // CashierPayin doesn't have ReturnURL in main struct
		TransactionValues: &TransactionValues{
			Status:             p.CashierPayinValues.Status,
			NotifyURL:          p.CashierPayinValues.NotifyURL,
			Remark:             p.CashierPayinValues.Remark,
			FlowNo:             &p.CashierPayinValues.FlowNo,
			SettleID:           p.CashierPayinValues.SettleID,
			SettleStatus:       p.CashierPayinValues.SettleStatus,
			SettledAt:          p.CashierPayinValues.SettledAt,
			ExpiredAt:          p.CashierPayinValues.ExpiredAt,
			ConfirmedAt:        p.CashierPayinValues.ConfirmedAt,
			CanceledAt:         p.CashierPayinValues.CanceledAt,
			CancelFailedResult: p.CashierPayinValues.CancelFailedResult,
			FeeAmount:          p.CashierPayinValues.Fee,
		},
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	return transaction
}
