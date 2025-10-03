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
