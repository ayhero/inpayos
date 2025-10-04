package protocol

// Admin 管理员信息V2
type Admin struct {
	AdminID      string `json:"admin_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	Role         string `json:"role"`
	Department   string `json:"department"`
	Status       string `json:"status"`
	ActiveStatus string `json:"active_status"`
	CreatedAt    int64  `json:"created_at"`
	LastLoginAt  *int64 `json:"last_login_at,omitempty"`
}

// 管理员角色常量
const (
	AdminRoleSuperAdmin = "super_admin"
	AdminRoleAdmin      = "admin"
	AdminRoleModerator  = "moderator"
	AdminRoleSupport    = "support"
	AdminRoleAnalyst    = "analyst"
)

// 权限常量
const (
	PermissionUserManagement      = "user_management"
	PermissionOrderManagement     = "order_management"
	PermissionPaymentManagement   = "payment_management"
	PermissionFinancialManagement = "financial_management"
	PermissionSystemConfig        = "system_config"
	PermissionAnalytics           = "analytics"
	PermissionCustomerSupport     = "customer_support"
	PermissionAdminManagement     = "admin_management"
	PermissionAuditLogs           = "audit_logs"
	PermissionEmergencyActions    = "emergency_actions"
)
