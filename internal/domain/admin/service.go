package admin

import (
	"context"
	"time"
)

// AdminRoleService defines business logic for admin role management
type AdminRoleService interface {
	// Role Management
	CreateRole(ctx context.Context, req *CreateRoleRequest) (*AdminRole, error)
	UpdateRole(ctx context.Context, id int64, req *UpdateRoleRequest) (*AdminRole, error)
	DeleteRole(ctx context.Context, id int64) error
	GetRole(ctx context.Context, id int64) (*AdminRoleResponse, error)
	GetRoleByName(ctx context.Context, roleName string) (*AdminRoleResponse, error)
	GetRoles(ctx context.Context, filter *AdminRoleFilter) (*RoleListResponse, error)

	// System Roles
	GetSystemRoles(ctx context.Context) ([]AdminRoleResponse, error)
	GetNonSystemRoles(ctx context.Context) ([]AdminRoleResponse, error)

	// Access Level Operations
	GetRolesByAccessLevel(ctx context.Context, minLevel, maxLevel int16) ([]AdminRoleResponse, error)
	PromoteRole(ctx context.Context, id int64, newLevel int16) error
	DemoteRole(ctx context.Context, id int64, newLevel int16) error

	// Statistics
	GetRoleStats(ctx context.Context) (*RoleStatsResponse, error)
	GetRoleUsage(ctx context.Context, roleID int64) (*RoleUsageResponse, error)

	// Validation
	ValidateRole(ctx context.Context, req *CreateRoleRequest) error
	CheckRolePermissions(ctx context.Context, roleID int64, requiredLevel int16) (bool, error)
	CanModifyRole(ctx context.Context, actorRoleID, targetRoleID int64) (bool, error)
}

// AdminUserService defines business logic for admin user management
type AdminUserService interface {
	// User Management
	CreateUser(ctx context.Context, req *CreateUserRequest) (*AdminUser, error)
	UpdateUser(ctx context.Context, id int64, req *UpdateUserRequest) (*AdminUser, error)
	DeleteUser(ctx context.Context, id int64) error
	GetUser(ctx context.Context, id int64) (*AdminUserResponse, error)
	GetUserByUUID(ctx context.Context, uuid string) (*AdminUserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*AdminUserResponse, error)
	GetUsers(ctx context.Context, filter *AdminUserFilter) (*UserListResponse, error)
	SearchUsers(ctx context.Context, query string, page, pageSize int) (*UserListResponse, error)

	// Authentication
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context, userID int64) error
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	ChangePassword(ctx context.Context, userID int64, req *ChangePasswordRequest) error
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
	RequestPasswordReset(ctx context.Context, email string) error

	// Two-Factor Authentication
	Enable2FA(ctx context.Context, userID int64) (*Enable2FAResponse, error)
	Verify2FA(ctx context.Context, userID int64, code string) error
	Disable2FA(ctx context.Context, userID int64, password string) error
	Generate2FAQRCode(ctx context.Context, userID int64) (string, error)

	// Status Management
	ActivateUser(ctx context.Context, id int64) error
	DeactivateUser(ctx context.Context, id int64) error
	SuspendUser(ctx context.Context, id int64, reason string) error
	UnsuspendUser(ctx context.Context, id int64) error
	GetActiveUsers(ctx context.Context) ([]AdminUserResponse, error)
	GetInactiveUsers(ctx context.Context) ([]AdminUserResponse, error)
	GetSuspendedUsers(ctx context.Context) ([]AdminUserResponse, error)

	// Role Management
	AssignRole(ctx context.Context, userID, roleID int64) error
	RemoveRole(ctx context.Context, userID int64) error
	GetUsersByRole(ctx context.Context, roleID int64) ([]AdminUserResponse, error)

	// Profile Management
	UpdateProfile(ctx context.Context, userID int64, req *UpdateProfileRequest) (*AdminUser, error)
	UpdateProfileImage(ctx context.Context, userID int64, imageURL string) error
	GetProfile(ctx context.Context, userID int64) (*ProfileResponse, error)

	// Statistics & Analytics
	GetUserStats(ctx context.Context) (*UserStatsResponse, error)
	GetActivityStats(ctx context.Context, startDate, endDate time.Time) (*ActivityStatsResponse, error)
	GetRecentLogins(ctx context.Context, limit int) ([]AdminUserResponse, error)
	GetUserActivity(ctx context.Context, userID int64, startDate, endDate time.Time) (*UserActivityResponse, error)

	// Audit & Tracking
	GetCreatedUsers(ctx context.Context, creatorID int64) ([]AdminUserResponse, error)
	TrackLogin(ctx context.Context, userID int64) error
	GetLoginHistory(ctx context.Context, userID int64, limit int) ([]LoginRecord, error)

	// Validation & Authorization
	ValidateUser(ctx context.Context, req *CreateUserRequest) error
	CheckUserPermissions(ctx context.Context, userID int64, requiredLevel int16) (bool, error)
	CanModifyUser(ctx context.Context, actorID, targetID int64) (bool, error)
	IsUserActive(ctx context.Context, userID int64) (bool, error)
}

// Request DTOs

// CreateRoleRequest represents a request to create an admin role
type CreateRoleRequest struct {
	RoleName        string  `json:"role_name" validate:"required,min=3,max=100"`
	RoleDescription string  `json:"role_description,omitempty"`
	AccessLevel     int16   `json:"access_level" validate:"required,min=1,max=10"`
	IsSystemRole    bool    `json:"is_system_role"`
	CreatedBy       *int64  `json:"created_by,omitempty"`
}

// UpdateRoleRequest represents a request to update an admin role
type UpdateRoleRequest struct {
	RoleName        string  `json:"role_name,omitempty" validate:"omitempty,min=3,max=100"`
	RoleDescription string  `json:"role_description,omitempty"`
	AccessLevel     *int16  `json:"access_level,omitempty" validate:"omitempty,min=1,max=10"`
}

// CreateUserRequest represents a request to create an admin user
type CreateUserRequest struct {
	FullName        string  `json:"full_name" validate:"required,min=2,max=100"`
	Email           string  `json:"email" validate:"required,email,max=150"`
	Phone           string  `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	Password        string  `json:"password" validate:"required,min=8"`
	RoleID          *int64  `json:"role_id,omitempty"`
	Status          string  `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
	ProfileImageURL string  `json:"profile_image_url,omitempty"`
	CreatedBy       *int64  `json:"created_by,omitempty"`
}

// UpdateUserRequest represents a request to update an admin user
type UpdateUserRequest struct {
	FullName        string  `json:"full_name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone           string  `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	RoleID          *int64  `json:"role_id,omitempty"`
	Status          string  `json:"status,omitempty" validate:"omitempty,oneof=active inactive suspended"`
	ProfileImageURL string  `json:"profile_image_url,omitempty"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	TwoFACode string `json:"two_fa_code,omitempty"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// ResetPasswordRequest represents a password reset request
type ResetPasswordRequest struct {
	Email           string `json:"email" validate:"required,email"`
	ResetToken      string `json:"reset_token" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	FullName        string `json:"full_name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone           string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
}

// Response DTOs

// AdminRoleResponse represents an admin role response
type AdminRoleResponse struct {
	ID              int64     `json:"id"`
	RoleName        string    `json:"role_name"`
	RoleDescription string    `json:"role_description,omitempty"`
	AccessLevel     int16     `json:"access_level"`
	IsSystemRole    bool      `json:"is_system_role"`
	UserCount       int64     `json:"user_count"`
	CreatedBy       *int64    `json:"created_by,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// RoleListResponse represents a paginated list of admin roles
type RoleListResponse struct {
	Roles      []AdminRoleResponse `json:"roles"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

// AdminUserResponse represents an admin user response
type AdminUserResponse struct {
	ID              int64              `json:"id"`
	UUID            string             `json:"uuid"`
	FullName        string             `json:"full_name"`
	Email           string             `json:"email"`
	Phone           string             `json:"phone,omitempty"`
	Role            *AdminRoleResponse `json:"role,omitempty"`
	Status          string             `json:"status"`
	LastLogin       *time.Time         `json:"last_login,omitempty"`
	Has2FA          bool               `json:"has_2fa"`
	ProfileImageURL string             `json:"profile_image_url,omitempty"`
	CreatedBy       *int64             `json:"created_by,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

// UserListResponse represents a paginated list of admin users
type UserListResponse struct {
	Users      []AdminUserResponse `json:"users"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

// LoginResponse represents a login response with tokens
type LoginResponse struct {
	User         *AdminUserResponse `json:"user"`
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	ExpiresIn    int64              `json:"expires_in"`
	TokenType    string             `json:"token_type"`
	Requires2FA  bool               `json:"requires_2fa"`
}

// TokenResponse represents a token refresh response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// Enable2FAResponse represents a 2FA setup response
type Enable2FAResponse struct {
	Secret    string `json:"secret"`
	QRCode    string `json:"qr_code"`
	BackupCodes []string `json:"backup_codes"`
}

// ProfileResponse represents a profile response with full details
type ProfileResponse struct {
	User            *AdminUserResponse `json:"user"`
	CreatedUsers    int64              `json:"created_users"`
	CreatedRoles    int64              `json:"created_roles"`
	TotalLogins     int64              `json:"total_logins"`
	RecentLogins    []time.Time        `json:"recent_logins"`
}

// RoleStatsResponse represents statistics about admin roles
type RoleStatsResponse struct {
	TotalRoles        int64              `json:"total_roles"`
	SystemRoles       int64              `json:"system_roles"`
	CustomRoles       int64              `json:"custom_roles"`
	ByAccessLevel     map[int16]int64    `json:"by_access_level"`
	MostUsedRole      *AdminRoleResponse `json:"most_used_role,omitempty"`
	MostUsedRoleCount int64              `json:"most_used_role_count"`
}

// RoleUsageResponse represents usage information for a role
type RoleUsageResponse struct {
	Role       *AdminRoleResponse   `json:"role"`
	UserCount  int64                `json:"user_count"`
	ActiveUsers int64               `json:"active_users"`
	Users      []AdminUserResponse  `json:"users"`
}

// UserStatsResponse represents statistics about admin users
type UserStatsResponse struct {
	TotalUsers      int64           `json:"total_users"`
	ActiveUsers     int64           `json:"active_users"`
	InactiveUsers   int64           `json:"inactive_users"`
	SuspendedUsers  int64           `json:"suspended_users"`
	Users2FA        int64           `json:"users_2fa"`
	ByRole          map[int64]int64 `json:"by_role"`
	SuperAdmins     int64           `json:"super_admins"`
	Admins          int64           `json:"admins"`
	Moderators      int64           `json:"moderators"`
	RecentLogins    int64           `json:"recent_logins"`
}

// ActivityStatsResponse represents activity statistics
type ActivityStatsResponse struct {
	Period          string                 `json:"period"`
	TotalLogins     int64                  `json:"total_logins"`
	UniqueUsers     int64                  `json:"unique_users"`
	AverageLogins   float64                `json:"average_logins"`
	LoginsByDate    map[string]int64       `json:"logins_by_date"`
	LoginsByUser    map[int64]int64        `json:"logins_by_user"`
	TopActiveUsers  []AdminUserResponse    `json:"top_active_users"`
}

// UserActivityResponse represents activity information for a user
type UserActivityResponse struct {
	User            *AdminUserResponse `json:"user"`
	TotalLogins     int64              `json:"total_logins"`
	LastLogin       *time.Time         `json:"last_login,omitempty"`
	CreatedUsers    int64              `json:"created_users"`
	CreatedRoles    int64              `json:"created_roles"`
	LoginHistory    []LoginRecord      `json:"login_history"`
}

// LoginRecord represents a single login record
type LoginRecord struct {
	UserID    int64     `json:"user_id"`
	LoginTime time.Time `json:"login_time"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
}
