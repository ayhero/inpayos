package models

type CashierTeam struct {
	ID     int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	TeamID string `gorm:"type:varchar(64);not null;uniqueIndex" json:"team_id"`
	*CashierTeamValues
	CreatedAt int64 `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli" json:"updated_at"`
}

type CashierTeamValues struct {
	Name        *string `gorm:"type:varchar(64);not null;uniqueIndex" json:"name"`
	Description *string `gorm:"type:varchar(255)" json:"description"`
	Status      *string `gorm:"type:varchar(32);not null;default:'active'" json:"status"` // active, inactive
}

func (CashierTeam) TableName() string {
	return "t_cashier_teams"
}

// CashierTeamValues Getter Methods
// GetName returns the Name value
func (ctv *CashierTeamValues) GetName() string {
	if ctv.Name == nil {
		return ""
	}
	return *ctv.Name
}

// GetDescription returns the Description value
func (ctv *CashierTeamValues) GetDescription() string {
	if ctv.Description == nil {
		return ""
	}
	return *ctv.Description
}

// GetStatus returns the Status value
func (ctv *CashierTeamValues) GetStatus() string {
	if ctv.Status == nil {
		return ""
	}
	return *ctv.Status
}

// CashierTeamValues Setter Methods (support method chaining)
// SetName sets the Name value
func (ctv *CashierTeamValues) SetName(value string) *CashierTeamValues {
	ctv.Name = &value
	return ctv
}

// SetDescription sets the Description value
func (ctv *CashierTeamValues) SetDescription(value string) *CashierTeamValues {
	ctv.Description = &value
	return ctv
}

// SetStatus sets the Status value
func (ctv *CashierTeamValues) SetStatus(value string) *CashierTeamValues {
	ctv.Status = &value
	return ctv
}

// SetValues sets multiple CashierTeamValues fields at once
func (ct *CashierTeam) SetValues(values *CashierTeamValues) *CashierTeam {
	if values == nil {
		return ct
	}

	if ct.CashierTeamValues == nil {
		ct.CashierTeamValues = &CashierTeamValues{}
	}

	// Set all fields from the provided values
	if values.Name != nil {
		ct.CashierTeamValues.SetName(*values.Name)
	}
	if values.Description != nil {
		ct.CashierTeamValues.SetDescription(*values.Description)
	}
	if values.Status != nil {
		ct.CashierTeamValues.SetStatus(*values.Status)
	}

	return ct
}
