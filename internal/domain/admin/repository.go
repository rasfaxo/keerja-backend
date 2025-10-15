package admin

import (
	"context"
	"time"
)

// AdminRoleRepository defines data access methods for AdminRole
type AdminRoleRepository interface {
	// Basic CRUD
	Create(ctx context.Context, role *AdminRole) error
	FindByID(ctx context.Context, id int64) (*AdminRole, error)
	FindByName(ctx context.Context, roleName string) (*AdminRole, error)
	Update(ctx context.Context, role *AdminRole) error
	Delete(ctx context.Context, id int64) error

	// Listing & Search
	List(ctx context.Context, filter *AdminRoleFilter) ([]AdminRole, int64, error)
	ListActive(ctx context.Context) ([]AdminRole, error)
	GetSystemRoles(ctx context.Context) ([]AdminRole, error)
	GetNonSystemRoles(ctx context.Context) ([]AdminRole, error)

	// Access Level Operations
	GetRolesByAccessLevel(ctx context.Context, minLevel, maxLevel int16) ([]AdminRole, error)
	GetRolesByMinAccessLevel(ctx context.Context, minLevel int16) ([]AdminRole, error)

	// Statistics
	Count(ctx context.Context) (int64, error)
	CountByAccessLevel(ctx context.Context, accessLevel int16) (int64, error)
	GetRoleStats(ctx context.Context) (*AdminRoleStats, error)
}

// AdminUserRepository defines data access methods for AdminUser
type AdminUserRepository interface {
	// Basic CRUD
	Create(ctx context.Context, user *AdminUser) error
	FindByID(ctx context.Context, id int64) (*AdminUser, error)
	FindByUUID(ctx context.Context, uuid string) (*AdminUser, error)
	FindByEmail(ctx context.Context, email string) (*AdminUser, error)
	Update(ctx context.Context, user *AdminUser) error
	Delete(ctx context.Context, id int64) error

	// Listing & Search
	List(ctx context.Context, filter *AdminUserFilter) ([]AdminUser, int64, error)
	ListByRole(ctx context.Context, roleID int64, page, pageSize int) ([]AdminUser, int64, error)
	ListByStatus(ctx context.Context, status string, page, pageSize int) ([]AdminUser, int64, error)
	SearchUsers(ctx context.Context, query string, page, pageSize int) ([]AdminUser, int64, error)

	// Status Operations
	UpdateStatus(ctx context.Context, id int64, status string) error
	ActivateUser(ctx context.Context, id int64) error
	DeactivateUser(ctx context.Context, id int64) error
	SuspendUser(ctx context.Context, id int64) error
	GetActiveUsers(ctx context.Context) ([]AdminUser, error)
	GetInactiveUsers(ctx context.Context) ([]AdminUser, error)
	GetSuspendedUsers(ctx context.Context) ([]AdminUser, error)

	// Role Operations
	UpdateRole(ctx context.Context, userID, roleID int64) error
	GetUsersByRole(ctx context.Context, roleID int64) ([]AdminUser, error)

	// Authentication & Security
	UpdatePassword(ctx context.Context, id int64, passwordHash string) error
	UpdateLastLogin(ctx context.Context, id int64) error
	Enable2FA(ctx context.Context, id int64, secret string) error
	Disable2FA(ctx context.Context, id int64) error
	Get2FAUsers(ctx context.Context) ([]AdminUser, error)

	// Profile Operations
	UpdateProfile(ctx context.Context, id int64, fullName, phone, profileImageURL string) error
	UpdateProfileImage(ctx context.Context, id int64, imageURL string) error

	// Statistics
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status string) (int64, error)
	CountByRole(ctx context.Context, roleID int64) (int64, error)
	GetUserStats(ctx context.Context) (*AdminUserStats, error)
	GetActivityStats(ctx context.Context, startDate, endDate time.Time) (*AdminActivityStats, error)
	GetRecentLogins(ctx context.Context, limit int) ([]AdminUser, error)

	// Audit & Tracking
	GetCreatedUsers(ctx context.Context, creatorID int64) ([]AdminUser, error)
	GetUserActivity(ctx context.Context, userID int64, startDate, endDate time.Time) (*UserActivity, error)
}

// AdminRoleFilter defines filter options for admin role queries
type AdminRoleFilter struct {
	Search       string
	AccessLevel  *int16
	MinLevel     *int16
	MaxLevel     *int16
	IsSystemRole *bool
	CreatedBy    *int64
	Page         int
	PageSize     int
	SortBy       string
	SortOrder    string
}

// AdminUserFilter defines filter options for admin user queries
type AdminUserFilter struct {
	Search           string
	Status           string
	RoleID           *int64
	MinAccessLevel   *int16
	Has2FA           *bool
	CreatedBy        *int64
	LastLoginAfter   *time.Time
	LastLoginBefore  *time.Time
	CreatedAfter     *time.Time
	CreatedBefore    *time.Time
	Page             int
	PageSize         int
	SortBy           string
	SortOrder        string
}

// AdminRoleStats contains statistics about admin roles
type AdminRoleStats struct {
	TotalRoles       int64
	SystemRoles      int64
	CustomRoles      int64
	ByAccessLevel    map[int16]int64
	MostUsedRole     *AdminRole
	MostUsedRoleCount int64
}

// AdminUserStats contains statistics about admin users
type AdminUserStats struct {
	TotalUsers      int64
	ActiveUsers     int64
	InactiveUsers   int64
	SuspendedUsers  int64
	Users2FA        int64
	ByRole          map[int64]int64
	SuperAdmins     int64
	Admins          int64
	Moderators      int64
}

// AdminActivityStats contains activity statistics for admin users
type AdminActivityStats struct {
	Period          string
	TotalLogins     int64
	UniqueUsers     int64
	AverageLogins   float64
	LoginsByDate    map[string]int64
	LoginsByUser    map[int64]int64
	TopActiveUsers  []AdminUser
}

// UserActivity contains activity information for a specific admin user
type UserActivity struct {
	UserID          int64
	TotalLogins     int64
	LastLogin       *time.Time
	CreatedUsers    int64
	CreatedRoles    int64
	LoginHistory    []time.Time
}
