package models

import "inpayos/internal/protocol"

// Checkout 收银台会话表
type Checkout struct {
	ID            uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	SessionID     string           `gorm:"column:session_id;type:varchar(64);uniqueIndex;not null" json:"session_id"`
	MerchantID    string           `gorm:"column:merchant_id;type:varchar(64);not null;index" json:"merchant_id"`
	BillID        string           `gorm:"column:bill_id;type:varchar(64);index" json:"bill_id"`
	Amount        string           `gorm:"column:amount;type:decimal(20,8);not null" json:"amount"`
	Currency      string           `gorm:"column:currency;type:varchar(10);not null" json:"currency"`
	Country       string           `gorm:"column:country;type:varchar(3)" json:"country"`
	PaymentMethod string           `gorm:"column:payment_method;type:varchar(32)" json:"payment_method"`
	ReturnURL     string           `gorm:"column:return_url;type:varchar(1024)" json:"return_url"`
	CancelURL     string           `gorm:"column:cancel_url;type:varchar(1024)" json:"cancel_url"`
	NotifyURL     string           `gorm:"column:notify_url;type:varchar(1024)" json:"notify_url"`
	Status        string           `gorm:"column:status;type:varchar(32);not null;default:'created'" json:"status"` // created, pending, completed, cancelled, expired
	TransactionID string           `gorm:"column:transaction_id;type:varchar(64);index" json:"transaction_id"`      // 关联的交易ID
	ChannelCode   string           `gorm:"column:channel_code;type:varchar(32)" json:"channel_code"`
	Metadata      protocol.MapData `gorm:"column:metadata;type:text" json:"metadata"` // JSON格式的元数据
	ErrorCode     string           `gorm:"column:error_code;type:varchar(32)" json:"error_code"`
	ErrorMsg      string           `gorm:"column:error_msg;type:varchar(512)" json:"error_msg"`
	ExpiredAt     int64            `gorm:"column:expired_at;type:bigint" json:"expired_at"`
	CompletedAt   int64            `gorm:"column:completed_at;type:bigint" json:"completed_at"`
	CreatedAt     int64            `gorm:"column:created_at;type:bigint;autoCreateTime:milli" json:"created_at"`
	UpdatedAt     int64            `gorm:"column:updated_at;type:bigint;autoUpdateTime:milli" json:"updated_at"`
}

// TableName 返回表名
func (Checkout) TableName() string {
	return "t_checkouts"
}
