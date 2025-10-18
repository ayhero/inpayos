package models

import (
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
)

// Admin 管理员表
type Admin struct {
	ID     int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID string `json:"user_id" gorm:"column:user_id;type:varchar(64);uniqueIndex"`
	*AdminValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type AdminValues struct {
	Salt     *string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	Username *string `json:"username" gorm:"column:username;type:varchar(50);uniqueIndex"`
	Email    *string `json:"email" gorm:"column:email;type:varchar(255);uniqueIndex"`
	Role     *string `json:"role" gorm:"column:role;type:varchar(50);index"`
	Status   *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'active'"`
	Password *string `json:"password" gorm:"column:password;type:varchar(128);not null"`
}

func (Admin) TableName() string {
	return "t_admins"
}

// 创建新管理员
func NewAdmin() *Admin {
	salt := utils.GenerateSalt()
	return &Admin{
		UserID: utils.GenerateAdminID(),
		AdminValues: &AdminValues{
			Salt:   &salt,
			Status: utils.StringPtr(protocol.StatusActive),
		},
		CreatedAt: utils.TimeNowMilli(),
	}
}

// Getter方法
func (av *AdminValues) GetUsername() string {
	if av.Username == nil {
		return ""
	}
	return *av.Username
}

func (av *AdminValues) GetEmail() string {
	if av.Email == nil {
		return ""
	}
	return *av.Email
}

func (av *AdminValues) GetRole() string {
	if av.Role == nil {
		return protocol.AdminRoleSupport
	}
	return *av.Role
}

func (av *AdminValues) GetStatus() string {
	if av.Status == nil {
		return protocol.StatusActive
	}
	return *av.Status
}

// Setter方法
func (av *AdminValues) SetUsername(value string) *AdminValues {
	av.Username = &value
	return av
}

func (av *AdminValues) SetEmail(value string) *AdminValues {
	av.Email = &value
	return av
}

func (av *AdminValues) SetRole(value string) *AdminValues {
	av.Role = &value
	return av
}

func (av *AdminValues) SetStatus(value string) *AdminValues {
	av.Status = &value
	return av
}

// 状态检查方法
func (a *AdminValues) IsActive() bool {
	return a.GetStatus() == protocol.StatusActive
}

func (a *AdminValues) IsInactive() bool {
	return a.GetStatus() == protocol.StatusInactive
}

func (a *AdminValues) IsSuspended() bool {
	return a.GetStatus() == protocol.StatusSuspended
}

func (a *AdminValues) IsLocked() bool {
	return a.GetStatus() == protocol.StatusLocked
}

func (a *AdminValues) IsSuperAdmin() bool {
	return a.GetRole() == protocol.AdminRoleSuperAdmin
}

func (a *AdminValues) IsAdmin() bool {
	return a.GetRole() == protocol.AdminRoleAdmin
}

func (a *AdminValues) Activate(adminID string) *AdminValues {
	a.SetStatus(protocol.StatusActive)
	return a
}

func (a *AdminValues) Deactivate(adminID string) *AdminValues {
	a.SetStatus(protocol.StatusInactive)
	return a
}

func (a *AdminValues) Suspend(adminID, reason string) *AdminValues {
	a.SetStatus(protocol.StatusSuspended)
	return a
}

func (a *AdminValues) Unsuspend(adminID string) *AdminValues {
	a.SetStatus(protocol.StatusActive)
	return a
}

func (a *AdminValues) Lock(adminID string) *AdminValues {
	a.SetStatus(protocol.StatusLocked)
	return a
}

func (a *AdminValues) Approve(adminID string) *AdminValues {
	a.SetStatus(protocol.StatusActive)
	return a
}

// 权限相关方法
func (a *AdminValues) GetPermissionList() []string {
	switch a.GetRole() {
	case protocol.AdminRoleSuperAdmin:
		return []string{
			protocol.PermissionUserManagement,
			protocol.PermissionOrderManagement,
			protocol.PermissionPaymentManagement,
			protocol.PermissionFinancialManagement,
			protocol.PermissionSystemConfig,
			protocol.PermissionAnalytics,
			protocol.PermissionCustomerSupport,
			protocol.PermissionAdminManagement,
			protocol.PermissionAuditLogs,
			protocol.PermissionEmergencyActions,
		}
	case protocol.AdminRoleAdmin:
		return []string{
			protocol.PermissionUserManagement,
			protocol.PermissionOrderManagement,
			protocol.PermissionPaymentManagement,
			protocol.PermissionAnalytics,
			protocol.PermissionCustomerSupport,
			protocol.PermissionAuditLogs,
		}
	case protocol.AdminRoleModerator:
		return []string{
			protocol.PermissionUserManagement,
			protocol.PermissionOrderManagement,
			protocol.PermissionCustomerSupport,
		}
	case protocol.AdminRoleSupport:
		return []string{
			protocol.PermissionCustomerSupport,
		}
	default:
		return []string{}
	}
}

// SetValues 设置AdminValues
func (a *Admin) SetValues(values *AdminValues) *Admin {
	if values == nil {
		return a
	}

	if a.AdminValues == nil {
		a.AdminValues = &AdminValues{}
	}

	if values.Username != nil {
		a.AdminValues.SetUsername(*values.Username)
	}
	if values.Email != nil {
		a.AdminValues.SetEmail(*values.Email)
	}
	if values.Role != nil {
		a.AdminValues.SetRole(*values.Role)
	}
	if values.Status != nil {
		a.AdminValues.SetStatus(*values.Status)
	}

	return a
}

// 转换方法
func (a *Admin) ToProtocol() *protocol.Admin {
	return &protocol.Admin{
		AdminID:   a.UserID,
		Username:  a.GetUsername(),
		Email:     a.GetEmail(),
		Role:      a.GetRole(),
		Status:    a.GetStatus(),
		CreatedAt: a.CreatedAt,
	}
}

func GetAdminUserByID(userID string) *Admin {
	var admin Admin
	result := ReadDB.Where("user_id = ?", userID).First(&admin)
	if result.Error != nil {
		return nil
	}
	return &admin
}
