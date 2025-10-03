package models

import "inpayos/internal/protocol"

type CashierMember struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	TeamID    string `gorm:"type:varchar(64);not null;index" json:"team_id"`
	CashierID string `gorm:"type:varchar(64);not null;index" json:"cashier_id"`
	*CashierMemberValues
	CreatedAt int64 `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli" json:"updated_at"`
}

type CashierMemberValues struct {
	Status       *string           `gorm:"type:varchar(32);not null;default:'active'" json:"status"`                    // active, inactive
	PayinStatus  *string           `json:"payin_status" gorm:"column:payin_status;type:varchar(16);default:'active'"`   // 收款状态：active, inactive, frozen, suspended
	PayinConfig  *protocol.MapData `json:"payin_config" gorm:"column:payin_config;type:text"`                           // 收款配置
	PayoutStatus *string           `json:"payout_status" gorm:"column:payout_status;type:varchar(16);default:'active'"` // 付款状态：active, inactive, frozen, suspended
	PayoutConfig *protocol.MapData `json:"payout_config" gorm:"column:payout_config;type:text"`                         // 付款配置
	Remark       *string           `gorm:"type:varchar(255)" json:"remark"`
}

func (CashierMember) TableName() string {
	return "t_cashier_members"
}
