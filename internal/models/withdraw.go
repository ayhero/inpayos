package models

import (
	"github.com/shopspring/decimal"
)

// Withdraw 提现记录表
type Withdraw struct {
	ID    uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TrxID string `json:"transaction_id" gorm:"column:transaction_id;type:varchar(64);uniqueIndex"`
	Salt  string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*WithdrawValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type WithdrawValues struct {
	UserID      *string          `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	ReqID       *string          `json:"req_id" gorm:"column:req_id;type:varchar(64);index"`
	Status      *string          `json:"status" gorm:"column:status;type:varchar(16);index;default:'pending'"`
	Ccy         *string          `json:"ccy" gorm:"column:ccy;type:varchar(16)"`
	Amount      *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Fee         *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(36,18);default:0"`
	ChannelCode *string          `json:"channel_code" gorm:"column:channel_code;type:varchar(32)"`
	NotifyURL   *string          `json:"notify_url" gorm:"column:notify_url;type:varchar(512)"`
	Country     *string          `json:"country" gorm:"column:country;type:varchar(8)"`
	CanceledAt  *int64           `json:"canceled_at" gorm:"column:canceled_at"`
	UpdatedAt   int64            `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Withdraw) TableName() string {
	return "t_withdraws"
}
