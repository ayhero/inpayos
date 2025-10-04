package models

import "inpayos/internal/protocol"

type ChannelAccount struct {
	ID          int64  `json:"id" gorm:"column:id;primaryKey;AUTO_INCREMENT"`
	MID         string `json:"mid" gorm:"column:mid"`
	ChannelCode string `json:"channel_code" gorm:"column:channel_code"`
	*ChannelAccountValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

type ChannelAccountValues struct {
	AccountID *string          `json:"account_id" gorm:"column:account_id;uniqueIndex"`
	Secret    *string          `json:"secret" gorm:"column:secret"`
	Detail    protocol.MapData `json:"detail" gorm:"column:detail;type:json;serializer:json"`
	Pkgs      *[]string        `json:"pkgs" gorm:"column:pkgs;type:json;serializer:json"`
	Status    *string          `json:"status" gorm:"column:status"`
	Settings  protocol.MapData `json:"settings" gorm:"column:settings;type:json;serializer:json"`
	Groups    *[]string        `json:"groups" gorm:"column:groups;type:json;serializer:json"`
}

func (t *ChannelAccount) TableName() string {
	return "t_channel_accounts"
}

func GetChannelAccountsByAccountID(accountID string) *ChannelAccount {
	var account ChannelAccount
	err := ReadDB.Where("account_id = ?", accountID).First(&account).Error
	if err != nil {
		return nil
	}
	return &account
}
func GetActiveChannelAccountByCode(mid, channelCode string) *ChannelAccount {
	var account ChannelAccount
	err := ReadDB.Where("mid=? and channel_code = ?", mid, channelCode).First(&account).Error
	if err != nil {
		return nil
	}
	return &account
}

// GetChannelAccounts 获取渠道账户列表
func GetChannelAccounts() []*ChannelAccount {
	var accounts []*ChannelAccount
	err := ReadDB.Find(&accounts).Error
	if err != nil {
		return nil
	}
	return accounts
}

// ChannelAccountValues Getter Methods
// GetAccountID returns the AccountID value
func (cav *ChannelAccountValues) GetAccountID() string {
	if cav.AccountID == nil {
		return ""
	}
	return *cav.AccountID
}

// GetSecret returns the Secret value
func (cav *ChannelAccountValues) GetSecret() string {
	if cav.Secret == nil {
		return ""
	}
	return *cav.Secret
}

// GetDetail returns the Detail value
func (cav *ChannelAccountValues) GetDetail() protocol.MapData {
	return cav.Detail
}

// GetPkgs returns the Pkgs value
func (cav *ChannelAccountValues) GetPkgs() []string {
	if cav.Pkgs == nil {
		return nil
	}
	return *cav.Pkgs
}

// GetStatus returns the Status value
func (cav *ChannelAccountValues) GetStatus() string {
	if cav.Status == nil {
		return ""
	}
	return *cav.Status
}

// GetSettings returns the Settings value
func (cav *ChannelAccountValues) GetSettings() protocol.MapData {
	return cav.Settings
}

// GetGroups returns the Groups value
func (cav *ChannelAccountValues) GetGroups() []string {
	if cav.Groups == nil {
		return nil
	}
	return *cav.Groups
}

// ChannelAccountValues Setter Methods (support method chaining)
// SetAccountID sets the AccountID value
func (cav *ChannelAccountValues) SetAccountID(value string) *ChannelAccountValues {
	cav.AccountID = &value
	return cav
}

// SetSecret sets the Secret value
func (cav *ChannelAccountValues) SetSecret(value string) *ChannelAccountValues {
	cav.Secret = &value
	return cav
}

// SetDetail sets the Detail value
func (cav *ChannelAccountValues) SetDetail(value protocol.MapData) *ChannelAccountValues {
	cav.Detail = value
	return cav
}

// SetPkgs sets the Pkgs value
func (cav *ChannelAccountValues) SetPkgs(value []string) *ChannelAccountValues {
	cav.Pkgs = &value
	return cav
}

// SetStatus sets the Status value
func (cav *ChannelAccountValues) SetStatus(value string) *ChannelAccountValues {
	cav.Status = &value
	return cav
}

// SetSettings sets the Settings value
func (cav *ChannelAccountValues) SetSettings(value protocol.MapData) *ChannelAccountValues {
	cav.Settings = value
	return cav
}

// SetGroups sets the Groups value
func (cav *ChannelAccountValues) SetGroups(value []string) *ChannelAccountValues {
	cav.Groups = &value
	return cav
}

// SetValues sets multiple ChannelAccountValues fields at once
func (ca *ChannelAccount) SetValues(values *ChannelAccountValues) *ChannelAccount {
	if values == nil {
		return ca
	}

	if ca.ChannelAccountValues == nil {
		ca.ChannelAccountValues = &ChannelAccountValues{}
	}

	// Set all fields from the provided values
	if values.AccountID != nil {
		ca.ChannelAccountValues.SetAccountID(*values.AccountID)
	}
	if values.Secret != nil {
		ca.ChannelAccountValues.SetSecret(*values.Secret)
	}
	// Detail is not a pointer, so we always set it
	ca.ChannelAccountValues.SetDetail(values.Detail)
	if values.Pkgs != nil {
		ca.ChannelAccountValues.SetPkgs(*values.Pkgs)
	}
	if values.Status != nil {
		ca.ChannelAccountValues.SetStatus(*values.Status)
	}
	// Settings is not a pointer, so we always set it
	ca.ChannelAccountValues.SetSettings(values.Settings)
	if values.Groups != nil {
		ca.ChannelAccountValues.SetGroups(*values.Groups)
	}

	return ca
}
