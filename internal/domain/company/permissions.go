package company

// Permission represents a specific action that can be performed
type Permission string

const (
	// Company Management Permissions
	PermissionUpdateCompany Permission = "company:update"
	PermissionDeleteCompany Permission = "company:delete"
	PermissionVerifyCompany Permission = "company:verify"

	// Employee Management Permissions
	PermissionInviteEmployee     Permission = "employee:invite"
	PermissionRemoveEmployee     Permission = "employee:remove"
	PermissionUpdateEmployeeRole Permission = "employee:update_role"
	PermissionViewEmployees      Permission = "employee:view"

	// Job Management Permissions
	PermissionCreateJob  Permission = "job:create"
	PermissionUpdateJob  Permission = "job:update"
	PermissionDeleteJob  Permission = "job:delete"
	PermissionPublishJob Permission = "job:publish"
	PermissionCloseJob   Permission = "job:close"
	PermissionViewJobs   Permission = "job:view"

	// Application Management Permissions
	PermissionViewApplications        Permission = "application:view"
	PermissionUpdateApplicationStatus Permission = "application:update_status"
	PermissionRejectApplication       Permission = "application:reject"

	// Review Management Permissions
	PermissionViewReviews  Permission = "review:view"
	PermissionReplyReview  Permission = "review:reply"
	PermissionDeleteReview Permission = "review:delete"

	// Statistics Permissions
	PermissionViewStatistics Permission = "statistics:view"
	PermissionViewAnalytics  Permission = "analytics:view"
)

// RolePermissions maps roles to their allowed permissions
var RolePermissions = map[string][]Permission{
	"admin": {
		// Company Management
		PermissionUpdateCompany,
		PermissionDeleteCompany,

		// Employee Management
		PermissionInviteEmployee,
		PermissionRemoveEmployee,
		PermissionUpdateEmployeeRole,
		PermissionViewEmployees,

		// Job Management
		PermissionCreateJob,
		PermissionUpdateJob,
		PermissionDeleteJob,
		PermissionPublishJob,
		PermissionCloseJob,
		PermissionViewJobs,

		// Application Management
		PermissionViewApplications,
		PermissionUpdateApplicationStatus,
		PermissionRejectApplication,

		// Review Management
		PermissionViewReviews,
		PermissionReplyReview,
		PermissionDeleteReview,

		// Statistics
		PermissionViewStatistics,
		PermissionViewAnalytics,
	},

	"recruiter": {
		// Job Management
		PermissionCreateJob,
		PermissionUpdateJob,
		PermissionDeleteJob,
		PermissionPublishJob,
		PermissionCloseJob,
		PermissionViewJobs,

		// Application Management
		PermissionViewApplications,
		PermissionUpdateApplicationStatus,
		PermissionRejectApplication,

		// Review Management
		PermissionViewReviews,
		PermissionReplyReview,

		// Statistics
		PermissionViewStatistics,

		// Limited Employee View
		PermissionViewEmployees,
	},

	"viewer": {
		// Read-only permissions
		PermissionViewJobs,
		PermissionViewApplications,
		PermissionViewReviews,
		PermissionViewStatistics,
		PermissionViewEmployees,
	},
}

// HasPermission checks if a role has a specific permission
func HasPermission(role string, permission Permission) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}

	return false
}

// GetRolePermissions returns all permissions for a role
func GetRolePermissions(role string) []Permission {
	permissions, exists := RolePermissions[role]
	if !exists {
		return []Permission{}
	}

	return permissions
}

// CanManageCompany checks if role can perform company management actions
func CanManageCompany(role string) bool {
	return HasPermission(role, PermissionUpdateCompany) ||
		HasPermission(role, PermissionDeleteCompany)
}

// CanManageEmployees checks if role can manage employees
func CanManageEmployees(role string) bool {
	return HasPermission(role, PermissionInviteEmployee) ||
		HasPermission(role, PermissionRemoveEmployee) ||
		HasPermission(role, PermissionUpdateEmployeeRole)
}

// CanManageJobs checks if role can manage jobs
func CanManageJobs(role string) bool {
	return HasPermission(role, PermissionCreateJob) ||
		HasPermission(role, PermissionUpdateJob) ||
		HasPermission(role, PermissionDeleteJob) ||
		HasPermission(role, PermissionPublishJob) ||
		HasPermission(role, PermissionCloseJob)
}

// CanManageApplications checks if role can manage applications
func CanManageApplications(role string) bool {
	return HasPermission(role, PermissionUpdateApplicationStatus) ||
		HasPermission(role, PermissionRejectApplication)
}

// CanViewOnly checks if role is view-only
func CanViewOnly(role string) bool {
	return role == "viewer"
}

// IsAdmin checks if role is admin
func IsAdmin(role string) bool {
	return role == "admin"
}

// IsRecruiter checks if role is recruiter
func IsRecruiter(role string) bool {
	return role == "recruiter"
}

// ValidateRole checks if a role is valid
func ValidateRole(role string) bool {
	validRoles := []string{"admin", "recruiter", "viewer"}
	for _, r := range validRoles {
		if r == role {
			return true
		}
	}
	return false
}

// GetRoleDescription returns a description for each role
func GetRoleDescription(role string) string {
	descriptions := map[string]string{
		"owner":     "Company owner with full control including ability to delete company and transfer ownership",
		"admin":     "Full access to all company features including employee management, job postings, applications, and company settings",
		"recruiter": "Can create and manage job postings, review applications, and respond to reviews. Cannot manage company settings or employees",
		"viewer":    "Read-only access to company data. Can view jobs, applications, reviews, and statistics but cannot make changes",
	}

	if desc, exists := descriptions[role]; exists {
		return desc
	}

	return "Unknown role"
}

// RoleHierarchy defines role hierarchy (higher number = more permissions)
var RoleHierarchy = map[string]int{
	"viewer":    1,
	"recruiter": 2,
	"admin":     3,
	"owner":     4, // Highest privilege
}

// HasHigherRole checks if role1 has higher or equal hierarchy than role2
func HasHigherRole(role1, role2 string) bool {
	level1, exists1 := RoleHierarchy[role1]
	level2, exists2 := RoleHierarchy[role2]

	if !exists1 || !exists2 {
		return false
	}

	return level1 >= level2
}
