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

// CashierMemberValues Getter Methods
// GetStatus returns the Status value
func (cmv *CashierMemberValues) GetStatus() string {
	if cmv.Status == nil {
		return ""
	}
	return *cmv.Status
}

// GetPayinStatus returns the PayinStatus value
func (cmv *CashierMemberValues) GetPayinStatus() string {
	if cmv.PayinStatus == nil {
		return ""
	}
	return *cmv.PayinStatus
}

// GetPayinConfig returns the PayinConfig value
func (cmv *CashierMemberValues) GetPayinConfig() protocol.MapData {
	if cmv.PayinConfig == nil {
		return nil
	}
	return *cmv.PayinConfig
}

// GetPayoutStatus returns the PayoutStatus value
func (cmv *CashierMemberValues) GetPayoutStatus() string {
	if cmv.PayoutStatus == nil {
		return ""
	}
	return *cmv.PayoutStatus
}

// GetPayoutConfig returns the PayoutConfig value
func (cmv *CashierMemberValues) GetPayoutConfig() protocol.MapData {
	if cmv.PayoutConfig == nil {
		return nil
	}
	return *cmv.PayoutConfig
}

// GetRemark returns the Remark value
func (cmv *CashierMemberValues) GetRemark() string {
	if cmv.Remark == nil {
		return ""
	}
	return *cmv.Remark
}

// CashierMemberValues Setter Methods (support method chaining)
// SetStatus sets the Status value
func (cmv *CashierMemberValues) SetStatus(value string) *CashierMemberValues {
	cmv.Status = &value
	return cmv
}

// SetPayinStatus sets the PayinStatus value
func (cmv *CashierMemberValues) SetPayinStatus(value string) *CashierMemberValues {
	cmv.PayinStatus = &value
	return cmv
}

// SetPayinConfig sets the PayinConfig value
func (cmv *CashierMemberValues) SetPayinConfig(value protocol.MapData) *CashierMemberValues {
	cmv.PayinConfig = &value
	return cmv
}

// SetPayoutStatus sets the PayoutStatus value
func (cmv *CashierMemberValues) SetPayoutStatus(value string) *CashierMemberValues {
	cmv.PayoutStatus = &value
	return cmv
}

// SetPayoutConfig sets the PayoutConfig value
func (cmv *CashierMemberValues) SetPayoutConfig(value protocol.MapData) *CashierMemberValues {
	cmv.PayoutConfig = &value
	return cmv
}

// SetRemark sets the Remark value
func (cmv *CashierMemberValues) SetRemark(value string) *CashierMemberValues {
	cmv.Remark = &value
	return cmv
}

// SetValues sets multiple CashierMemberValues fields at once
func (cm *CashierMember) SetValues(values *CashierMemberValues) *CashierMember {
	if values == nil {
		return cm
	}

	if cm.CashierMemberValues == nil {
		cm.CashierMemberValues = &CashierMemberValues{}
	}

	// Set all fields from the provided values
	if values.Status != nil {
		cm.CashierMemberValues.SetStatus(*values.Status)
	}
	if values.PayinStatus != nil {
		cm.CashierMemberValues.SetPayinStatus(*values.PayinStatus)
	}
	if values.PayinConfig != nil {
		cm.CashierMemberValues.SetPayinConfig(*values.PayinConfig)
	}
	if values.PayoutStatus != nil {
		cm.CashierMemberValues.SetPayoutStatus(*values.PayoutStatus)
	}
	if values.PayoutConfig != nil {
		cm.CashierMemberValues.SetPayoutConfig(*values.PayoutConfig)
	}
	if values.Remark != nil {
		cm.CashierMemberValues.SetRemark(*values.Remark)
	}

	return cm
}
