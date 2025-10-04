package models

import (
	"fmt"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
	"slices"
	"strings"
)

// CashierAdmin 管理员表 - 后台管理员账户、权限和登录管理
type CashierAdmin struct {
	ID      int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	AdminID string `json:"admin_id" gorm:"column:admin_id;type:varchar(64);uniqueIndex"`
	Salt    string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*CashierAdminValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type CashierAdminValues struct {
	// 基本信息
	Username *string `json:"username" gorm:"column:username;type:varchar(50);uniqueIndex"`
	Email    *string `json:"email" gorm:"column:email;type:varchar(255);uniqueIndex"`
	Phone    *string `json:"phone" gorm:"column:phone;type:varchar(20);index"`

	// 个人信息
	FirstName *string `json:"first_name" gorm:"column:first_name;type:varchar(100)"`
	LastName  *string `json:"last_name" gorm:"column:last_name;type:varchar(100)"`
	FullName  *string `json:"full_name" gorm:"column:full_name;type:varchar(200)"`
	Avatar    *string `json:"avatar" gorm:"column:avatar;type:varchar(500)"`

	// 认证信息
	PasswordHash     *string `json:"-" gorm:"column:password_hash;type:varchar(255)"`
	PasswordSalt     *string `json:"-" gorm:"column:password_salt;type:varchar(255)"`
	EmailVerified    *bool   `json:"email_verified" gorm:"column:email_verified;default:false"`
	PhoneVerified    *bool   `json:"phone_verified" gorm:"column:phone_verified;default:false"`
	TwoFactorEnabled *bool   `json:"two_factor_enabled" gorm:"column:two_factor_enabled;default:false"`
	TwoFactorSecret  *string `json:"-" gorm:"column:two_factor_secret;type:varchar(255)"`

	// 角色和权限
	Role        *string `json:"role" gorm:"column:role;type:varchar(50);index"`        // super_admin, admin, moderator, support, analyst
	Permissions *string `json:"permissions" gorm:"column:permissions;type:json"`       // JSON数组存储权限列表
	Department  *string `json:"department" gorm:"column:department;type:varchar(100)"` // 部门
	JobTitle    *string `json:"job_title" gorm:"column:job_title;type:varchar(100)"`   // 职位

	// 状态管理
	Status       *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'active'"`          // active, inactive, suspended, locked
	ActiveStatus *string `json:"active_status" gorm:"column:active_status;type:varchar(32);default:'offline'"` // online, offline, busy

	// 登录相关
	LastLoginAt    *int64  `json:"last_login_at" gorm:"column:last_login_at"`
	LastLoginIP    *string `json:"last_login_ip" gorm:"column:last_login_ip;type:varchar(45)"`
	LoginCount     *int    `json:"login_count" gorm:"column:login_count;default:0"`
	FailedAttempts *int    `json:"failed_attempts" gorm:"column:failed_attempts;default:0"`
	LastFailedAt   *int64  `json:"last_failed_at" gorm:"column:last_failed_at"`
	LockedUntil    *int64  `json:"locked_until" gorm:"column:locked_until"`

	// 会话管理
	CurrentSessionID      *string `json:"current_session_id" gorm:"column:current_session_id;type:varchar(255)"`
	SessionCount          *int    `json:"session_count" gorm:"column:session_count;default:0"`
	MaxConcurrentSessions *int    `json:"max_concurrent_sessions" gorm:"column:max_concurrent_sessions;default:3"`

	// 安全设置
	PasswordChangedAt  *int64  `json:"password_changed_at" gorm:"column:password_changed_at"`
	MustChangePassword *bool   `json:"must_change_password" gorm:"column:must_change_password;default:false"`
	AllowedIPs         *string `json:"allowed_ips" gorm:"column:allowed_ips;type:text"` // 允许的IP地址列表
	SecurityQuestion   *string `json:"security_question" gorm:"column:security_question;type:varchar(255)"`
	SecurityAnswerHash *string `json:"-" gorm:"column:security_answer_hash;type:varchar(255)"`

	// 审计信息
	CreatedBy        *string `json:"created_by" gorm:"column:created_by;type:varchar(64)"`
	LastUpdatedBy    *string `json:"last_updated_by" gorm:"column:last_updated_by;type:varchar(64)"`
	ApprovedBy       *string `json:"approved_by" gorm:"column:approved_by;type:varchar(64)"`
	ApprovedAt       *int64  `json:"approved_at" gorm:"column:approved_at"`
	SuspendedBy      *string `json:"suspended_by" gorm:"column:suspended_by;type:varchar(64)"`
	SuspendedAt      *int64  `json:"suspended_at" gorm:"column:suspended_at"`
	SuspensionReason *string `json:"suspension_reason" gorm:"column:suspension_reason;type:varchar(500)"`

	// 工作时间和区域
	WorkingHours    *string `json:"working_hours" gorm:"column:working_hours;type:json"` // JSON对象存储工作时间配置
	WorkingTimeZone *string `json:"working_timezone" gorm:"column:working_timezone;type:varchar(50)"`
	WorkingAreas    *string `json:"working_areas" gorm:"column:working_areas;type:json"` // JSON数组存储负责区域

	// 联系信息
	OfficePhone      *string `json:"office_phone" gorm:"column:office_phone;type:varchar(20)"`
	OfficeAddress    *string `json:"office_address" gorm:"column:office_address;type:varchar(500)"`
	EmergencyContact *string `json:"emergency_contact" gorm:"column:emergency_contact;type:varchar(255)"`

	// 偏好设置
	Language    *string `json:"language" gorm:"column:language;type:varchar(10);default:'en'"`
	Timezone    *string `json:"timezone" gorm:"column:timezone;type:varchar(50);default:'UTC'"`
	DateFormat  *string `json:"date_format" gorm:"column:date_format;type:varchar(20);default:'YYYY-MM-DD'"`
	TimeFormat  *string `json:"time_format" gorm:"column:time_format;type:varchar(20);default:'24h'"`
	Preferences *string `json:"preferences" gorm:"column:preferences;type:json"` // JSON对象存储个人偏好

	// 统计数据
	ActionsToday     *int `json:"actions_today" gorm:"column:actions_today;default:0"`
	ActionsThisWeek  *int `json:"actions_this_week" gorm:"column:actions_this_week;default:0"`
	ActionsThisMonth *int `json:"actions_this_month" gorm:"column:actions_this_month;default:0"`
	TotalActions     *int `json:"total_actions" gorm:"column:total_actions;default:0"`

	// 元数据
	Metadata *string `json:"metadata" gorm:"column:metadata;type:json"`
	Notes    *string `json:"notes" gorm:"column:notes;type:text"`
	Tags     *string `json:"tags" gorm:"column:tags;type:varchar(500)"` // 逗号分隔的标签

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (CashierAdmin) TableName() string {
	return "t_admins"
}

// 创建新的管理员对象
func NewCashierAdminV2() *CashierAdmin {
	return &CashierAdmin{
		AdminID: utils.GenerateAdminID(),
		Salt:    utils.GenerateSalt(),
		CashierAdminValues: &CashierAdminValues{
			Status:                utils.StringPtr(protocol.AdminStatusActive),
			ActiveStatus:          utils.StringPtr(protocol.AdminActiveStatusOffline),
			EmailVerified:         utils.BoolPtr(false),
			PhoneVerified:         utils.BoolPtr(false),
			TwoFactorEnabled:      utils.BoolPtr(false),
			LoginCount:            utils.IntPtr(0),
			FailedAttempts:        utils.IntPtr(0),
			SessionCount:          utils.IntPtr(0),
			MaxConcurrentSessions: utils.IntPtr(3),
			MustChangePassword:    utils.BoolPtr(false),
			Language:              utils.StringPtr("en"),
			Timezone:              utils.StringPtr("UTC"),
			DateFormat:            utils.StringPtr("YYYY-MM-DD"),
			TimeFormat:            utils.StringPtr("24h"),
			ActionsToday:          utils.IntPtr(0),
			ActionsThisWeek:       utils.IntPtr(0),
			ActionsThisMonth:      utils.IntPtr(0),
			TotalActions:          utils.IntPtr(0),
		},
	}
}

// SetValues 更新CashierAdminV2Values中的非nil值
func (a *CashierAdminValues) SetValues(values *CashierAdminValues) {
	if values == nil {
		return
	}

	if values.Username != nil {
		a.Username = values.Username
	}
	if values.Email != nil {
		a.Email = values.Email
	}
	if values.Phone != nil {
		a.Phone = values.Phone
	}
	if values.FirstName != nil {
		a.FirstName = values.FirstName
	}
	if values.LastName != nil {
		a.LastName = values.LastName
	}
	if values.FullName != nil {
		a.FullName = values.FullName
	}
	if values.Role != nil {
		a.Role = values.Role
	}
	if values.Status != nil {
		a.Status = values.Status
	}
	if values.Department != nil {
		a.Department = values.Department
	}
	if values.JobTitle != nil {
		a.JobTitle = values.JobTitle
	}
	if values.Permissions != nil {
		a.Permissions = values.Permissions
	}
	if values.Notes != nil {
		a.Notes = values.Notes
	}
	if values.UpdatedAt > 0 {
		a.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (a *CashierAdminValues) GetUsername() string {
	if a.Username == nil {
		return ""
	}
	return *a.Username
}

func (a *CashierAdminValues) GetEmail() string {
	if a.Email == nil {
		return ""
	}
	return *a.Email
}

func (a *CashierAdminValues) GetFullName() string {
	if a.FullName == nil {
		return ""
	}
	return *a.FullName
}

func (a *CashierAdminValues) GetRole() string {
	if a.Role == nil {
		return ""
	}
	return *a.Role
}

func (a *CashierAdminValues) GetStatus() string {
	if a.Status == nil {
		return ""
	}
	return *a.Status
}

func (a *CashierAdminValues) GetActiveStatus() string {
	if a.ActiveStatus == nil {
		return ""
	}
	return *a.ActiveStatus
}

func (a *CashierAdminValues) GetLoginCount() int {
	if a.LoginCount == nil {
		return 0
	}
	return *a.LoginCount
}

func (a *CashierAdminValues) GetFailedAttempts() int {
	if a.FailedAttempts == nil {
		return 0
	}
	return *a.FailedAttempts
}

func (a *CashierAdminValues) GetSessionCount() int {
	if a.SessionCount == nil {
		return 0
	}
	return *a.SessionCount
}

func (a *CashierAdminValues) GetMaxConcurrentSessions() int {
	if a.MaxConcurrentSessions == nil {
		return 3
	}
	return *a.MaxConcurrentSessions
}

func (a *CashierAdminValues) GetEmailVerified() bool {
	if a.EmailVerified == nil {
		return false
	}
	return *a.EmailVerified
}

func (a *CashierAdminValues) GetTwoFactorEnabled() bool {
	if a.TwoFactorEnabled == nil {
		return false
	}
	return *a.TwoFactorEnabled
}

func (a *CashierAdminValues) GetMustChangePassword() bool {
	if a.MustChangePassword == nil {
		return false
	}
	return *a.MustChangePassword
}

func (a *CashierAdminValues) GetActionsToday() int {
	if a.ActionsToday == nil {
		return 0
	}
	return *a.ActionsToday
}

func (a *CashierAdminValues) GetActionsThisWeek() int {
	if a.ActionsThisWeek == nil {
		return 0
	}
	return *a.ActionsThisWeek
}

func (a *CashierAdminValues) GetActionsThisMonth() int {
	if a.ActionsThisMonth == nil {
		return 0
	}
	return *a.ActionsThisMonth
}

func (a *CashierAdminValues) GetTotalActions() int {
	if a.TotalActions == nil {
		return 0
	}
	return *a.TotalActions
}

// Setter 方法
func (a *CashierAdminValues) SetUsername(username string) *CashierAdminValues {
	a.Username = &username
	return a
}

func (a *CashierAdminValues) SetEmail(email string) *CashierAdminValues {
	a.Email = &email
	return a
}

func (a *CashierAdminValues) SetPhone(phone string) *CashierAdminValues {
	a.Phone = &phone
	return a
}

func (a *CashierAdminValues) SetFullName(firstName, lastName string) *CashierAdminValues {
	a.FirstName = &firstName
	a.LastName = &lastName
	fullName := strings.TrimSpace(firstName + " " + lastName)
	a.FullName = &fullName
	return a
}

func (a *CashierAdminValues) SetRole(role string) *CashierAdminValues {
	a.Role = &role
	return a
}

func (a *CashierAdminValues) SetStatus(status string) *CashierAdminValues {
	a.Status = &status
	return a
}

func (a *CashierAdminValues) SetActiveStatus(status string) *CashierAdminValues {
	a.ActiveStatus = &status
	return a
}

func (a *CashierAdminValues) SetDepartment(department string) *CashierAdminValues {
	a.Department = &department
	return a
}

func (a *CashierAdminValues) SetJobTitle(title string) *CashierAdminValues {
	a.JobTitle = &title
	return a
}

func (a *CashierAdminValues) SetPasswordHash(hash, salt string) *CashierAdminValues {
	a.PasswordHash = &hash
	a.PasswordSalt = &salt
	now := utils.TimeNowMilli()
	a.PasswordChangedAt = &now
	return a
}

func (a *CashierAdminValues) SetEmailVerified(verified bool) *CashierAdminValues {
	a.EmailVerified = &verified
	return a
}

func (a *CashierAdminValues) SetTwoFactorEnabled(enabled bool) *CashierAdminValues {
	a.TwoFactorEnabled = &enabled
	return a
}

func (a *CashierAdminValues) SetTwoFactorSecret(secret string) *CashierAdminValues {
	a.TwoFactorSecret = &secret
	return a
}

func (a *CashierAdminValues) SetMustChangePassword(must bool) *CashierAdminValues {
	a.MustChangePassword = &must
	return a
}

// 业务方法
func (a *CashierAdmin) IsActive() bool {
	return a.GetStatus() == protocol.StatusActive
}

func (a *CashierAdmin) IsInactive() bool {
	return a.GetStatus() == protocol.StatusInactive
}

func (a *CashierAdmin) IsSuspended() bool {
	return a.GetStatus() == protocol.StatusSuspended
}

func (a *CashierAdmin) IsLocked() bool {
	return a.GetStatus() == protocol.StatusLocked
}

func (a *CashierAdmin) IsOnline() bool {
	return a.GetActiveStatus() == protocol.StatusOnline
}

func (a *CashierAdmin) IsOffline() bool {
	return a.GetActiveStatus() == protocol.StatusOffline
}

func (a *CashierAdmin) IsBusy() bool {
	return a.GetActiveStatus() == protocol.StatusBusy
}

func (a *CashierAdmin) IsSuperAdmin() bool {
	return a.GetRole() == protocol.AdminRoleSuperAdmin
}

func (a *CashierAdmin) IsCashierAdmin() bool {
	return a.GetRole() == protocol.AdminRoleAdmin
}

func (a *CashierAdmin) CanLogin() bool {
	return a.IsActive() && !a.IsAccountLocked()
}

func (a *CashierAdmin) IsAccountLocked() bool {
	if a.LockedUntil == nil {
		return false
	}
	return utils.TimeNowMilli() < *a.LockedUntil
}

func (a *CashierAdmin) HasExceededSessionLimit() bool {
	return a.GetSessionCount() >= a.GetMaxConcurrentSessions()
}

func (a *CashierAdmin) ShouldForcePasswordChange() bool {
	return a.GetMustChangePassword()
}

// 权限相关方法
func (a *CashierAdminValues) HasPermission(permission string) bool {
	if a.Permissions == nil {
		return false
	}

	var permissions []string
	if err := utils.FromJSON(*a.Permissions, &permissions); err != nil {
		return false
	}

	return slices.Contains(permissions, permission)
}

func (a *CashierAdminValues) AddPermission(permission string) error {
	var permissions []string
	if a.Permissions != nil {
		if err := utils.FromJSON(*a.Permissions, &permissions); err != nil {
			return fmt.Errorf("failed to parse existing permissions: %v", err)
		}
	}

	// 避免重复添加
	for _, p := range permissions {
		if p == permission {
			return nil
		}
	}

	permissions = append(permissions, permission)
	permissionsJSON, err := utils.ToJSON(permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal permissions: %v", err)
	}

	a.Permissions = &permissionsJSON
	return nil
}

func (a *CashierAdminValues) RemovePermission(permission string) error {
	if a.Permissions == nil {
		return nil
	}

	var permissions []string
	if err := utils.FromJSON(*a.Permissions, &permissions); err != nil {
		return fmt.Errorf("failed to parse existing permissions: %v", err)
	}

	var newPermissions []string
	for _, p := range permissions {
		if p != permission {
			newPermissions = append(newPermissions, p)
		}
	}

	permissionsJSON, err := utils.ToJSON(newPermissions)
	if err != nil {
		return fmt.Errorf("failed to marshal permissions: %v", err)
	}

	a.Permissions = &permissionsJSON
	return nil
}

func (a *CashierAdminValues) SetPermissions(permissions []string) error {
	permissionsJSON, err := utils.ToJSON(permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal permissions: %v", err)
	}

	a.Permissions = &permissionsJSON
	return nil
}

func (a *CashierAdminValues) GetPermissions() []string {
	if a.Permissions == nil {
		return []string{}
	}

	var permissions []string
	if err := utils.FromJSON(*a.Permissions, &permissions); err != nil {
		return []string{}
	}

	return permissions
}

// 登录相关方法
func (a *CashierAdminValues) RecordLogin(ip string, sessionID string) *CashierAdminValues {
	now := utils.TimeNowMilli()
	a.LastLoginAt = &now
	a.LastLoginIP = &ip
	a.CurrentSessionID = &sessionID

	loginCount := a.GetLoginCount() + 1
	a.LoginCount = &loginCount

	sessionCount := a.GetSessionCount() + 1
	a.SessionCount = &sessionCount

	// 重置失败次数
	a.FailedAttempts = utils.IntPtr(0)

	return a
}

func (a *CashierAdminValues) RecordFailedLogin() *CashierAdminValues {
	now := utils.TimeNowMilli()
	a.LastFailedAt = &now

	attempts := a.GetFailedAttempts() + 1
	a.FailedAttempts = &attempts

	// 如果失败次数过多，锁定账户
	if attempts >= 5 {
		lockUntil := now + (30 * 60 * 1000) // 锁定30分钟
		a.LockedUntil = &lockUntil
	}

	return a
}

func (a *CashierAdminValues) Logout() *CashierAdminValues {
	a.CurrentSessionID = nil
	sessionCount := a.GetSessionCount() - 1
	if sessionCount < 0 {
		sessionCount = 0
	}
	a.SessionCount = &sessionCount
	a.SetActiveStatus(protocol.StatusOffline)

	return a
}

func (a *CashierAdminValues) UnlockAccount() *CashierAdminValues {
	a.LockedUntil = nil
	a.FailedAttempts = utils.IntPtr(0)
	return a
}

// 状态管理方法
func (a *CashierAdminValues) Activate(adminID string) *CashierAdminValues {
	a.SetStatus(protocol.StatusActive)
	a.LastUpdatedBy = &adminID
	return a
}

func (a *CashierAdminValues) Deactivate(adminID string) *CashierAdminValues {
	a.SetStatus(protocol.StatusInactive)
	a.LastUpdatedBy = &adminID
	return a
}

func (a *CashierAdminValues) Suspend(adminID, reason string) *CashierAdminValues {
	a.SetStatus(protocol.StatusSuspended)
	a.SuspendedBy = &adminID
	a.SuspensionReason = &reason
	now := utils.TimeNowMilli()
	a.SuspendedAt = &now
	a.LastUpdatedBy = &adminID
	return a
}

func (a *CashierAdminValues) Unsuspend(adminID string) *CashierAdminValues {
	a.SetStatus(protocol.StatusActive)
	a.SuspendedBy = nil
	a.SuspendedAt = nil
	a.SuspensionReason = nil
	a.LastUpdatedBy = &adminID
	return a
}

func (a *CashierAdminValues) Lock(adminID string) *CashierAdminValues {
	a.SetStatus(protocol.StatusLocked)
	now := utils.TimeNowMilli()
	a.LockedUntil = &now
	a.LastUpdatedBy = &adminID
	return a
}

// 审批方法
func (a *CashierAdminValues) Approve(adminID string) *CashierAdminValues {
	a.ApprovedBy = &adminID
	now := utils.TimeNowMilli()
	a.ApprovedAt = &now
	a.SetStatus(protocol.StatusActive)
	return a
}

// 统计更新方法
func (a *CashierAdminValues) IncrementActions() *CashierAdminValues {
	today := a.GetActionsToday() + 1
	week := a.GetActionsThisWeek() + 1
	month := a.GetActionsThisMonth() + 1
	total := a.GetTotalActions() + 1

	a.ActionsToday = &today
	a.ActionsThisWeek = &week
	a.ActionsThisMonth = &month
	a.TotalActions = &total

	return a
}

// 工作时间设置
func (a *CashierAdminValues) SetWorkingHours(workingHours map[string]interface{}) error {
	workingHoursJSON, err := utils.ToJSON(workingHours)
	if err != nil {
		return fmt.Errorf("failed to marshal working hours: %v", err)
	}

	a.WorkingHours = &workingHoursJSON
	return nil
}

func (a *CashierAdminValues) SetWorkingAreas(areas []string) error {
	areasJSON, err := utils.ToJSON(areas)
	if err != nil {
		return fmt.Errorf("failed to marshal working areas: %v", err)
	}

	a.WorkingAreas = &areasJSON
	return nil
}

// 偏好设置
func (a *CashierAdminValues) SetPreferences(preferences map[string]interface{}) error {
	preferencesJSON, err := utils.ToJSON(preferences)
	if err != nil {
		return fmt.Errorf("failed to marshal preferences: %v", err)
	}

	a.Preferences = &preferencesJSON
	return nil
}

// IP白名单管理
func (a *CashierAdminValues) SetAllowedIPs(ips []string) *CashierAdminValues {
	allowedIPs := strings.Join(ips, ",")
	a.AllowedIPs = &allowedIPs
	return a
}

func (a *CashierAdminValues) IsIPAllowed(ip string) bool {
	if a.AllowedIPs == nil || *a.AllowedIPs == "" {
		return true // 如果没有设置白名单，允许所有IP
	}

	allowedIPs := strings.Split(*a.AllowedIPs, ",")
	for _, allowedIP := range allowedIPs {
		if strings.TrimSpace(allowedIP) == ip {
			return true
		}
	}

	return false
}

// 创建具有指定角色的管理员
func NewCashierAdminV2WithRole(username, email, role string) *CashierAdmin {
	admin := NewCashierAdminV2()
	admin.SetUsername(username).
		SetEmail(email).
		SetRole(role)

	// 设置角色对应的权限 - 简化实现
	// admin.SetPermissions([]string{}) // 可以根据角色设置具体权限

	return admin
}

func GetCashierAdminByID(adminID string) *CashierAdmin {
	var admin CashierAdmin
	if err := GetDB().Where("id = ?", adminID).First(&admin).Error; err != nil {
		return nil
	}
	return &admin
}

// NewCashierAdminInfoV2FromModel 从模型创建管理员信息
func (admin *CashierAdmin) Protocol() protocol.Admin {
	if admin == nil {
		return protocol.Admin{}
	}

	adminInfo := protocol.Admin{
		AdminID:   admin.AdminID,
		Username:  admin.GetUsername(),
		Email:     admin.GetEmail(),
		FullName:  admin.GetFullName(),
		CreatedAt: admin.CreatedAt,
	}

	// 处理可能为空的指针字段
	if admin.Role != nil {
		adminInfo.Role = *admin.Role
	}

	if admin.Department != nil {
		adminInfo.Department = *admin.Department
	}

	if admin.Status != nil {
		adminInfo.Status = *admin.Status
	}

	if admin.ActiveStatus != nil {
		adminInfo.ActiveStatus = *admin.ActiveStatus
	}

	if admin.LastLoginAt != nil {
		adminInfo.LastLoginAt = admin.LastLoginAt
	}

	return adminInfo
}

func (a *CashierAdminValues) GetFirstName() string {
	if a.FirstName == nil {
		return ""
	}
	return *a.FirstName
}

func (a *CashierAdminValues) GetLastName() string {
	if a.LastName == nil {
		return ""
	}
	return *a.LastName
}
