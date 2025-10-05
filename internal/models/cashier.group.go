package models

import "inpayos/internal/protocol"

type CashierGroup struct {
	ID   int64  `json:"id" gorm:"column:id;primaryKey;AUTO_INCREMENT"`
	Tid  string `json:"tid" gorm:"column:tid"`
	Code string `json:"code" gorm:"column:code;type:varchar(50);uniqueIndex"`
	*CashierGroupValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

type CashierGroupValues struct {
	Name    *string       `json:"name" gorm:"column:name;type:varchar(255)"`
	Status  *string       `json:"status" gorm:"column:status;type:varchar(50);default:0"`
	Setting *GroupSetting `json:"setting" gorm:"column:setting;serializer:json"`
	Members GroupMembers  `json:"members" gorm:"column:members;serializer:json"`
}

func (t *CashierGroup) TableName() string {
	return "t_cashier_groups"
}

func GetActiveCashierGroupByCode(code string) *CashierGroup {
	group := &CashierGroup{}
	err := ReadDB.Where("code = ?", code).First(group).Error
	if err != nil {
		return nil
	}
	return group
}

// GetActiveCashierGroups returns all active channel groups
func GetActiveCashierGroups() ([]*CashierGroup, error) {
	var groups []*CashierGroup
	err := ReadDB.Where("status = ?", protocol.StatusActive).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// CashierGroupValues Getter Methods
// GetName returns the Name value
func (cgv *CashierGroupValues) GetName() string {
	if cgv.Name == nil {
		return ""
	}
	return *cgv.Name
}

// GetStatus returns the Status value
func (cgv *CashierGroupValues) GetStatus() string {
	if cgv.Status == nil {
		return ""
	}
	return *cgv.Status
}

// GetSetting returns the Setting value
func (cgv *CashierGroupValues) GetSetting() *GroupSetting {
	return cgv.Setting
}

// GetMembers returns the Members value
func (cgv *CashierGroupValues) GetMembers() GroupMembers {
	return cgv.Members
}

// CashierGroupValues Setter Methods (support method chaining)
// SetName sets the Name value
func (cgv *CashierGroupValues) SetName(value string) *CashierGroupValues {
	cgv.Name = &value
	return cgv
}

// SetStatus sets the Status value
func (cgv *CashierGroupValues) SetStatus(value string) *CashierGroupValues {
	cgv.Status = &value
	return cgv
}

// SetSetting sets the Setting value
func (cgv *CashierGroupValues) SetSetting(value *GroupSetting) *CashierGroupValues {
	cgv.Setting = value
	return cgv
}

// SetMembers sets the Members value
func (cgv *CashierGroupValues) SetMembers(value GroupMembers) *CashierGroupValues {
	cgv.Members = value
	return cgv
}

// SetValues sets multiple CashierGroupValues fields at once
func (cg *CashierGroup) SetValues(values *CashierGroupValues) *CashierGroup {
	if values == nil {
		return cg
	}

	if cg.CashierGroupValues == nil {
		cg.CashierGroupValues = &CashierGroupValues{}
	}

	// Set all fields from the provided values
	if values.Name != nil {
		cg.CashierGroupValues.SetName(*values.Name)
	}
	if values.Status != nil {
		cg.CashierGroupValues.SetStatus(*values.Status)
	}
	if values.Setting != nil {
		cg.CashierGroupValues.SetSetting(values.Setting)
	}
	// Members is not a pointer, so we always set it
	cg.CashierGroupValues.SetMembers(values.Members)

	return cg
}
