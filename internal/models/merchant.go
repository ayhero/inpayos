package models

import (
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
)

type Merchant struct {
	ID  int64  `json:"id" gorm:"column:id;primaryKey;AUTO_INCREMENT"`
	MID string `json:"mid" gorm:"column:mid;type:varchar(64);uniqueIndex"`
	*MerchantValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

type MerchantValues struct {
	AuthID    *string `json:"auth_id" gorm:"column:auth_id;type:varchar(32);uniqueIndex"`
	Name      *string `json:"name" gorm:"column:name;type:varchar(64)"`
	Type      *string `json:"type" gorm:"column:type;type:varchar(32)"`
	Email     *string `json:"email" gorm:"column:email;type:varchar(128);uniqueIndex"`
	Phone     *string `json:"phone" gorm:"column:phone;type:varchar(20)"`
	Status    *string `json:"status" gorm:"column:status;type:varchar(32)"`
	Password  *string `json:"password" gorm:"column:password;type:varchar(128);not null"`
	Salt      *string `json:"salt" gorm:"column:salt;type:varchar(128);not null"`
	Region    *string `json:"region" gorm:"column:region;type:varchar(32)"`
	Avatar    *string `json:"avatar" gorm:"column:avatar;type:varchar(255)"`
	G2FA      *string `json:"g2fa" gorm:"column:g2fa;type:varchar(256)"`
	NotifyURL *string `json:"notify_url" gorm:"column:notify_url;type:varchar(1024)"`
	RegIP     *string `json:"reg_ip" gorm:"column:reg_ip;type:varchar(64)"` // 注册IP
}

func (t *Merchant) TableName() string {
	return "t_merchant"
}
func (m *MerchantValues) GetPassword() string {
	if m.Password == nil {
		return ""
	}
	return *m.Password
}
func (m *MerchantValues) SetPassword(password string) *MerchantValues {
	m.Password = &password
	return m
}
func (m *MerchantValues) GetSalt() string {
	if m.Salt == nil {
		return ""
	}
	return *m.Salt
}
func (m *MerchantValues) SetSalt(salt string) *MerchantValues {
	m.Salt = &salt
	return m
}
func (m *MerchantValues) GetRegion() string {
	if m.Region == nil {
		return ""
	}
	return *m.Region
}
func (m *MerchantValues) SetRegion(region string) *MerchantValues {
	m.Region = &region
	return m
}
func (m *MerchantValues) GetAvatar() string {
	if m.Avatar == nil {
		return ""
	}
	return *m.Avatar
}
func (m *MerchantValues) SetAvatar(avatar string) *MerchantValues {
	m.Avatar = &avatar
	return m
}
func (m *MerchantValues) GetG2FA() string {
	if m.G2FA == nil {
		return ""
	}
	return *m.G2FA
}
func (m *MerchantValues) SetG2FA(g2fa string) *MerchantValues {
	m.G2FA = &g2fa
	return m
}

func (m *MerchantValues) GetAuthID() string {
	if m.AuthID == nil {
		return ""
	}
	return *m.AuthID
}
func (m *MerchantValues) SetAuthID(authID string) *MerchantValues {
	m.AuthID = &authID
	return m
}

func (m *MerchantValues) GetName() string {
	if m.Name == nil {
		return ""
	}
	return *m.Name
}

// SetName sets the merchant name and returns the instance
func (m *MerchantValues) SetName(name string) *MerchantValues {
	m.Name = &name
	return m
}

// GetType returns the merchant type
func (m *MerchantValues) GetType() string {
	if m.Type == nil {
		return ""
	}
	return *m.Type
}

// SetType sets the merchant type and returns the instance
func (m *MerchantValues) SetType(t string) *MerchantValues {
	m.Type = &t
	return m
}

// GetEmail returns the merchant email
func (m *MerchantValues) GetEmail() string {
	if m.Email == nil {
		return ""
	}
	return *m.Email
}

// SetEmail sets the merchant email and returns the instance
func (m *MerchantValues) SetEmail(email string) *MerchantValues {
	m.Email = &email
	return m
}

// GetPhone returns the merchant phone
func (m *MerchantValues) GetPhone() string {
	if m.Phone == nil {
		return ""
	}
	return *m.Phone
}

// SetPhone sets the merchant phone and returns the instance
func (m *MerchantValues) SetPhone(phone string) *MerchantValues {
	m.Phone = &phone
	return m
}

// GetStatus returns the merchant status
func (m *MerchantValues) GetStatus() string {
	if m.Status == nil {
		return ""
	}
	return *m.Status
}

// SetStatus sets the merchant status and returns the instance
func (m *MerchantValues) SetStatus(status string) *MerchantValues {
	m.Status = &status
	return m
}

// GetNotifyURL returns the merchant notify URL
func (m *MerchantValues) GetNotifyURL() string {
	if m.NotifyURL == nil {
		return ""
	}
	return *m.NotifyURL
}

// SetNotifyURL sets the merchant notify URL and returns the instance
func (m *MerchantValues) SetNotifyURL(url string) *MerchantValues {
	m.NotifyURL = &url
	return m
}

// GetRegIP returns the merchant registration IP
func (m *MerchantValues) GetRegIP() string {
	if m.RegIP == nil {
		return ""
	}
	return *m.RegIP
}

// SetRegIP sets the merchant registration IP and returns the instance
func (m *MerchantValues) SetRegIP(ip string) *MerchantValues {
	m.RegIP = &ip
	return m
}

// Protocol converts model.Merchant to protocol.MerchantInfo
func (t *Merchant) Protocol() *protocol.Merchant {
	var region, avatar string
	if t.MerchantValues.Region != nil {
		region = *t.MerchantValues.Region
	}
	if t.MerchantValues.Avatar != nil {
		avatar = *t.MerchantValues.Avatar
	}

	return &protocol.Merchant{
		Mid:     t.MID,
		Name:    t.MerchantValues.GetName(),
		Type:    t.MerchantValues.GetType(),
		Email:   t.MerchantValues.GetEmail(),
		Phone:   t.MerchantValues.GetPhone(),
		Status:  t.MerchantValues.GetStatus(),
		Region:  region,
		Avatar:  avatar,
		HasG2FA: t.GetG2FA() != "",
	}
}

func GetMerchantByAuthID(authID string) (m *Merchant) {
	if err := ReadDB.Where("auth_id = ?", authID).First(&m).Error; err != nil {
		//log.Get().Errorf("GetMerchantByAuthID error: %v", err)
		return nil
	}
	return
}

func GetMerchantByEmail(email string) (data *Merchant) {
	err := ReadDB.Where("email = ?", email).First(&data).Error
	if err != nil {
		//log.Get().Errorf("GetMerchantByEmail error: %v", err)
		return nil
	}
	return
}
func GetMerchantByMID(mid string) (data *Merchant) {
	err := ReadDB.Where("mid = ?", mid).First(&data).Error
	if err != nil {
		//log.Get().Errorf("GetMerchantByMID error: %v", err)
		return nil
	}
	return
}

func GetMerchantByAppID(appID string) (data *Merchant) {
	err := ReadDB.Where("app_id = ?", appID).First(&data).Error
	if err != nil {
		//log.Get().Errorf("GetMerchantByAppID error: %v", err)
		return nil
	}
	return
}

func GetActiveMerchants() ([]*Merchant, error) {
	return GetMerchantsByStatus(protocol.StatusActive)
}

// GetMerchantsByStatus 获取指定状态的商户列表
func GetMerchantsByStatus(status string) ([]*Merchant, error) {
	var merchants []*Merchant
	err := ReadDB.Where("status = ?", status).Find(&merchants).Error
	if err != nil {
		return nil, err
	}
	return merchants, nil
}

func (u *Merchant) Decrypt() {
	salt := u.GetSalt()
	pwd, err := utils.Decrypt(u.GetPassword(), []byte(salt))
	if err == nil {
		u.SetPassword(pwd)
	}
}

func (u *Merchant) Encrypt() {
	salt := u.GetSalt()
	pwd, err := utils.Encrypt([]byte(u.GetPassword()), []byte(salt))
	if err == nil {
		u.SetPassword(pwd)
	}
}

// IsPasswordValid 验证密码是否正确
func (u *Merchant) IsPasswordValid(password string) bool {
	return u.GetPassword() == password
}

func CheckMerchantEmail(email string) bool {
	var count int64
	ReadDB.Model(&Merchant{}).Where("email = ?", email).Count(&count)
	return count > 0
}
