package models

import (
	"github.com/shopspring/decimal"
)

// Transaction 通用交易记录表（作为所有业务交易的抽象层）
// 每个具体业务表（Payment, Receipt, Refund等）通过 ToTransaction() 方法转换为此通用模型
type Transaction struct {
	ID    uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TrxID string `json:"transaction_id" gorm:"column:transaction_id;type:varchar(64);uniqueIndex"`
	*TransactionValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type TransactionValues struct {
	UserID        *string          `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	UserType      *string          `json:"user_type" gorm:"column:user_type;type:varchar(16);index"`
	ReqID         *string          `json:"req_id" gorm:"column:req_id;type:varchar(64);index"`
	TrxType       *string          `json:"trx_type" gorm:"column:trx_type;type:varchar(16);index"` // receipt, payment, refund, transfer
	Status        *string          `json:"status" gorm:"column:status;type:varchar(16);index"`     // pending, processing, success, failed
	Amount        *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Fee           *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(36,18);default:0"`
	Ccy           *string          `json:"ccy" gorm:"column:ccy;type:varchar(16)"`
	ChannelCode   *string          `json:"channel_code" gorm:"column:channel_code;type:varchar(32)"`
	PaymentMethod *string          `json:"payment_method" gorm:"column:payment_method;type:varchar(32)"`
	NotifyURL     *string          `json:"notify_url" gorm:"column:notify_url;type:varchar(512)"`
	ReturnURL     *string          `json:"return_url" gorm:"column:return_url;type:varchar(512)"`
	NotifyStatus  *string          `json:"notify_status" gorm:"column:notify_status;type:varchar(16);default:'pending'"`
	NotifyTimes   *int             `json:"notify_times" gorm:"column:notify_times;type:int;default:0"`
	OriTrxID      *string          `json:"ori_trx_id" gorm:"column:ori_trx_id;type:varchar(64)"` // 原交易ID(退款使用)
	Metadata      *string          `json:"metadata" gorm:"column:metadata;type:json"`
	Remark        *string          `json:"remark" gorm:"column:remark;type:varchar(512)"`
	ExpiredAt     *int64           `json:"expired_at" gorm:"column:expired_at"`
	ConfirmedAt   *int64           `json:"confirmed_at" gorm:"column:confirmed_at"`
	CanceledAt    *int64           `json:"canceled_at" gorm:"column:canceled_at"`
	UpdatedAt     int64            `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Transaction) TableName() string {
	return "t_transactions"
}
