package models

import (
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
)

type CashierTeam struct {
	ID  int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Tid string `json:"tid" gorm:"column:tid"`
	*CashierTeamValues
	CreatedAt int64 `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli" json:"updated_at"`
}

type CashierTeamValues struct {
	Salt        *string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	Description *string `gorm:"type:varchar(255)" json:"description"`
	AuthID      *string `json:"auth_id" gorm:"column:auth_id;type:varchar(32);uniqueIndex"`
	Name        *string `json:"name" gorm:"column:name;type:varchar(64)"`
	Type        *string `json:"type" gorm:"column:type;type:varchar(32)"`
	Email       *string `json:"email" gorm:"column:email;type:varchar(128);uniqueIndex"`
	Phone       *string `json:"phone" gorm:"column:phone;type:varchar(20)"`
	Status      *string `json:"status" gorm:"column:status;type:varchar(32)"`
	Password    *string `json:"password" gorm:"column:password;type:varchar(128);not null"`
	Region      *string `json:"region" gorm:"column:region;type:varchar(32)"`
	Avatar      *string `json:"avatar" gorm:"column:avatar;type:varchar(255)"`
	G2FA        *string `json:"g2fa" gorm:"column:g2fa;type:varchar(256)"`
	NotifyURL   *string `json:"notify_url" gorm:"column:notify_url;type:varchar(1024)"`
	RegIP       *string `json:"reg_ip" gorm:"column:reg_ip;type:varchar(64)"` // 注册IP
}

func (CashierTeam) TableName() string {
	return "t_cashier_teams"
}

// NewCashierTeam 创建新的出纳员团队
func NewCashierTeam() *CashierTeam {
	salt := utils.GenerateSalt()
	return &CashierTeam{
		Tid: utils.GenerateCashierTeamID(),
		CashierTeamValues: &CashierTeamValues{
			Salt: &salt,
		},
	}
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

// GetAuthID returns the AuthID value
func (ctv *CashierTeamValues) GetAuthID() string {
	if ctv.AuthID == nil {
		return ""
	}
	return *ctv.AuthID
}

// GetType returns the Type value
func (ctv *CashierTeamValues) GetType() string {
	if ctv.Type == nil {
		return ""
	}
	return *ctv.Type
}

// GetEmail returns the Email value
func (ctv *CashierTeamValues) GetEmail() string {
	if ctv.Email == nil {
		return ""
	}
	return *ctv.Email
}

// GetPhone returns the Phone value
func (ctv *CashierTeamValues) GetPhone() string {
	if ctv.Phone == nil {
		return ""
	}
	return *ctv.Phone
}

// GetPassword returns the Password value
func (ctv *CashierTeamValues) GetPassword() string {
	if ctv.Password == nil {
		return ""
	}
	return *ctv.Password
}

// GetRegion returns the Region value
func (ctv *CashierTeamValues) GetRegion() string {
	if ctv.Region == nil {
		return ""
	}
	return *ctv.Region
}

// GetAvatar returns the Avatar value
func (ctv *CashierTeamValues) GetAvatar() string {
	if ctv.Avatar == nil {
		return ""
	}
	return *ctv.Avatar
}

// GetG2FA returns the G2FA value
func (ctv *CashierTeamValues) GetG2FA() string {
	if ctv.G2FA == nil {
		return ""
	}
	return *ctv.G2FA
}

// GetNotifyURL returns the NotifyURL value
func (ctv *CashierTeamValues) GetNotifyURL() string {
	if ctv.NotifyURL == nil {
		return ""
	}
	return *ctv.NotifyURL
}

// GetRegIP returns the RegIP value
func (ctv *CashierTeamValues) GetRegIP() string {
	if ctv.RegIP == nil {
		return ""
	}
	return *ctv.RegIP
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

// SetAuthID sets the AuthID value
func (ctv *CashierTeamValues) SetAuthID(value string) *CashierTeamValues {
	ctv.AuthID = &value
	return ctv
}

// SetType sets the Type value
func (ctv *CashierTeamValues) SetType(value string) *CashierTeamValues {
	ctv.Type = &value
	return ctv
}

// SetEmail sets the Email value
func (ctv *CashierTeamValues) SetEmail(value string) *CashierTeamValues {
	ctv.Email = &value
	return ctv
}

// SetPhone sets the Phone value
func (ctv *CashierTeamValues) SetPhone(value string) *CashierTeamValues {
	ctv.Phone = &value
	return ctv
}

// SetPassword sets the Password value
func (ctv *CashierTeamValues) SetPassword(value string) *CashierTeamValues {
	ctv.Password = &value
	return ctv
}

// SetRegion sets the Region value
func (ctv *CashierTeamValues) SetRegion(value string) *CashierTeamValues {
	ctv.Region = &value
	return ctv
}

// SetAvatar sets the Avatar value
func (ctv *CashierTeamValues) SetAvatar(value string) *CashierTeamValues {
	ctv.Avatar = &value
	return ctv
}

// SetG2FA sets the G2FA value
func (ctv *CashierTeamValues) SetG2FA(value string) *CashierTeamValues {
	ctv.G2FA = &value
	return ctv
}

// SetNotifyURL sets the NotifyURL value
func (ctv *CashierTeamValues) SetNotifyURL(value string) *CashierTeamValues {
	ctv.NotifyURL = &value
	return ctv
}

// SetRegIP sets the RegIP value
func (ctv *CashierTeamValues) SetRegIP(value string) *CashierTeamValues {
	ctv.RegIP = &value
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
	if values.Description != nil {
		ct.CashierTeamValues.SetDescription(*values.Description)
	}
	if values.AuthID != nil {
		ct.CashierTeamValues.SetAuthID(*values.AuthID)
	}
	if values.Name != nil {
		ct.CashierTeamValues.SetName(*values.Name)
	}
	if values.Type != nil {
		ct.CashierTeamValues.SetType(*values.Type)
	}
	if values.Email != nil {
		ct.CashierTeamValues.SetEmail(*values.Email)
	}
	if values.Phone != nil {
		ct.CashierTeamValues.SetPhone(*values.Phone)
	}
	if values.Status != nil {
		ct.CashierTeamValues.SetStatus(*values.Status)
	}
	if values.Password != nil {
		ct.CashierTeamValues.SetPassword(*values.Password)
	}
	if values.Region != nil {
		ct.CashierTeamValues.SetRegion(*values.Region)
	}
	if values.Avatar != nil {
		ct.CashierTeamValues.SetAvatar(*values.Avatar)
	}
	if values.G2FA != nil {
		ct.CashierTeamValues.SetG2FA(*values.G2FA)
	}
	if values.NotifyURL != nil {
		ct.CashierTeamValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.RegIP != nil {
		ct.CashierTeamValues.SetRegIP(*values.RegIP)
	}

	return ct
}

func (u *CashierTeam) Decrypt() {
	if u.Salt == nil {
		return
	}
	salt := *u.Salt
	pwd, err := utils.Decrypt(u.GetPassword(), []byte(salt))
	if err == nil {
		u.SetPassword(pwd)
	}
}

func (u *CashierTeam) Encrypt() {
	if u.Salt == nil {
		return
	}
	salt := *u.Salt
	pwd, err := utils.Encrypt([]byte(u.GetPassword()), []byte(salt))
	if err == nil {
		u.SetPassword(pwd)
	}
}

// IsPasswordValid 验证密码是否正确
func (u *CashierTeam) IsPasswordValid(password string) bool {
	return u.GetPassword() == password
}
func GetCashierTeamByTid(tid string) *CashierTeam {
	var ct CashierTeam
	if err := ReadDB.Where("tid = ?", tid).First(&ct).Error; err != nil {
		return nil
	}
	return &ct
}

func GetCashierTeamByEmail(email string) *CashierTeam {
	var ct CashierTeam
	if err := ReadDB.Where("email = ?", email).First(&ct).Error; err != nil {
		return nil
	}
	return &ct
}

func (t *CashierTeam) Protocol() *protocol.CashierTeam {
	info := &protocol.CashierTeam{}
	return info
}

func CheckCashierTeamEmail(email string) bool {
	var count int64
	ReadDB.Model(&CashierTeam{}).Where("email = ?", email).Count(&count)
	return count > 0
}
