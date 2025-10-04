package models

import (
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
)

// Admin 管理员表
type Admin struct {
	ID      int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	AdminID string `json:"admin_id" gorm:"column:admin_id;type:varchar(64);uniqueIndex"`
	Salt    string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*AdminValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type AdminValues struct {
	Username     *string `json:"username" gorm:"column:username;type:varchar(50);uniqueIndex"`
	Email        *string `json:"email" gorm:"column:email;type:varchar(255);uniqueIndex"`
	Role         *string `json:"role" gorm:"column:role;type:varchar(50);index"`
	Status       *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'active'"`
	ActiveStatus *string `json:"active_status" gorm:"column:active_status;type:varchar(32);default:'offline'"`
	UpdatedAt    int64   `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Admin) TableName() string {
	return "t_admins"
}

// 创建新管理员
func NewAdmin() *Admin {
	return &Admin{
		AdminID: utils.GenerateAdminID(),
		Salt:    utils.GenerateSalt(),
		AdminValues: &AdminValues{
			Status:       utils.StringPtr(protocol.AdminStatusActive),
			ActiveStatus: utils.StringPtr(protocol.AdminActiveStatusOffline),
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
		return protocol.AdminStatusActive
	}
	return *av.Status
}

func (av *AdminValues) GetActiveStatus() string {
	if av.ActiveStatus == nil {
		return protocol.AdminActiveStatusOffline
	}
	return *av.ActiveStatus
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

func (av *AdminValues) SetActiveStatus(value string) *AdminValues {
	av.ActiveStatus = &value
	return av
}

// 状态检查方法
func (a *AdminValues) IsActive() bool {
	return a.GetStatus() == protocol.AdminStatusActive
}

func (a *AdminValues) IsInactive() bool {
	return a.GetStatus() == protocol.AdminStatusInactive
}

func (a *AdminValues) IsSuspended() bool {
	return a.GetStatus() == protocol.AdminStatusSuspended
}

func (a *AdminValues) IsLocked() bool {
	return a.GetStatus() == protocol.AdminStatusLocked
}

func (a *AdminValues) IsOnline() bool {
	return a.GetActiveStatus() == protocol.AdminActiveStatusOnline
}

func (a *AdminValues) IsOffline() bool {
	return a.GetActiveStatus() == protocol.AdminActiveStatusOffline
}

func (a *AdminValues) IsBusy() bool {
	return a.GetActiveStatus() == protocol.AdminActiveStatusBusy
}

func (a *AdminValues) IsSuperAdmin() bool {
	return a.GetRole() == protocol.AdminRoleSuperAdmin
}

func (a *AdminValues) IsAdmin() bool {
	return a.GetRole() == protocol.AdminRoleAdmin
}

// 业务操作方法
func (a *AdminValues) Logout() *AdminValues {
	a.SetActiveStatus(protocol.AdminActiveStatusOffline)
	return a
}

func (a *AdminValues) Activate(adminID string) *AdminValues {
	a.SetStatus(protocol.AdminStatusActive)
	return a
}

func (a *AdminValues) Deactivate(adminID string) *AdminValues {
	a.SetStatus(protocol.AdminStatusInactive)
	return a
}

func (a *AdminValues) Suspend(adminID, reason string) *AdminValues {
	a.SetStatus(protocol.AdminStatusSuspended)
	return a
}

func (a *AdminValues) Unsuspend(adminID string) *AdminValues {
	a.SetStatus(protocol.AdminStatusActive)
	return a
}

func (a *AdminValues) Lock(adminID string) *AdminValues {
	a.SetStatus(protocol.AdminStatusLocked)
	return a
}

func (a *AdminValues) Approve(adminID string) *AdminValues {
	a.SetStatus(protocol.AdminStatusActive)
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

// 转换方法
func (a *Admin) ToProtocol() *protocol.Admin {
	return &protocol.Admin{
		AdminID:      a.AdminID,
		Username:     a.GetUsername(),
		Email:        a.GetEmail(),
		Role:         a.GetRole(),
		Status:       a.GetStatus(),
		ActiveStatus: a.GetActiveStatus(),
		CreatedAt:    a.CreatedAt,
	}
}
