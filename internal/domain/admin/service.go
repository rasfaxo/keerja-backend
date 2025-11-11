package admin

import (
	"context"
	"time"
)

// AdminAuthService defines business logic for admin authentication
type AdminAuthService interface {
	// Authentication
	Login(ctx context.Context, email, password string) (*AdminUser, string, string, error) // returns: user, accessToken, refreshToken, error
	Logout(ctx context.Context, adminID int64) error
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error) // returns: newAccessToken, newRefreshToken, error

	// Profile Management
	GetCurrentProfile(ctx context.Context, adminID int64) (*AdminUser, error)
	ChangePassword(ctx context.Context, adminID int64, currentPassword, newPassword string) error

	// Session Management
	ValidateSession(ctx context.Context, adminID int64) (bool, error)
	InvalidateAllSessions(ctx context.Context, adminID int64) error
}

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

// AdminJobService defines business logic for admin job moderation
// Note: Actual implementation is in internal/service/admin_job_service.go
type AdminJobService interface {
	// Job approval/rejection (moderation)
	// ApproveJob changes job status from pending_review to published
	ApproveJob(ctx context.Context, jobID int64) (interface{}, error)
	// RejectJob changes job status from pending_review back to draft
	RejectJob(ctx context.Context, jobID int64, reason string) (interface{}, error)

	// Job list for approval
	// GetPendingJobs retrieves all jobs pending review (status = pending_review)
	GetPendingJobs(ctx context.Context, page, limit int) ([]interface{}, int64, error)
	// GetJobsForReview retrieves jobs for review with specific status
	GetJobsForReview(ctx context.Context, status string, page, limit int) ([]interface{}, int64, error)
}

// AdminCompanyService defines business logic for admin company management
// This service is used by admins to moderate, manage, and oversee companies
type AdminCompanyService interface {
	// Company Moderation Queue
	// Task 2.1: List companies with filters, search, and pagination
	ListCompanies(ctx context.Context, req *AdminCompanyListRequest) (*AdminCompanyListResponse, error)

	// Task 2.2: Get full company details for moderation
	GetCompanyDetail(ctx context.Context, companyID int64) (*AdminCompanyDetailResponse, error)

	// Task 2.3: Update company verification status (Approve/Reject/Suspend)
	UpdateCompanyStatus(ctx context.Context, companyID int64, req *AdminCompanyStatusRequest, adminID int64) error

	// Task 2.4: Edit company details (admin support)
	UpdateCompany(ctx context.Context, companyID int64, req *AdminUpdateCompanyRequest, adminID int64) error

	// Task 2.5: Delete company (with validation)
	DeleteCompany(ctx context.Context, companyID int64, req *AdminDeleteCompanyRequest, adminID int64) error

	// Additional operations
	GetCompanyStats(ctx context.Context, companyID int64) (*AdminCompanyStatsResponse, error)
	GetDashboardStats(ctx context.Context) (*AdminDashboardStatsResponse, error)
	BulkUpdateStatus(ctx context.Context, companyIDs []int64, status string, adminID int64) (*BulkOperationResult, error)
	GetAuditLogs(ctx context.Context, companyID int64, page, limit int) (*AuditLogListResponse, error)
}

// Request DTOs

// CreateRoleRequest represents a request to create an admin role
type CreateRoleRequest struct {
	RoleName        string `json:"role_name" validate:"required,min=3,max=100"`
	RoleDescription string `json:"role_description,omitempty"`
	AccessLevel     int16  `json:"access_level" validate:"required,min=1,max=10"`
	IsSystemRole    bool   `json:"is_system_role"`
	CreatedBy       *int64 `json:"created_by,omitempty"`
}

// UpdateRoleRequest represents a request to update an admin role
type UpdateRoleRequest struct {
	RoleName        string `json:"role_name,omitempty" validate:"omitempty,min=3,max=100"`
	RoleDescription string `json:"role_description,omitempty"`
	AccessLevel     *int16 `json:"access_level,omitempty" validate:"omitempty,min=1,max=10"`
}

// CreateUserRequest represents a request to create an admin user
type CreateUserRequest struct {
	FullName        string `json:"full_name" validate:"required,min=2,max=100"`
	Email           string `json:"email" validate:"required,email,max=150"`
	Phone           string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	Password        string `json:"password" validate:"required,min=8"`
	RoleID          *int64 `json:"role_id,omitempty"`
	Status          string `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
	CreatedBy       *int64 `json:"created_by,omitempty"`
}

// UpdateUserRequest represents a request to update an admin user
type UpdateUserRequest struct {
	FullName        string `json:"full_name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone           string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	RoleID          *int64 `json:"role_id,omitempty"`
	Status          string `json:"status,omitempty" validate:"omitempty,oneof=active inactive suspended"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
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
	Secret      string   `json:"secret"`
	QRCode      string   `json:"qr_code"`
	BackupCodes []string `json:"backup_codes"`
}

// ProfileResponse represents a profile response with full details
type ProfileResponse struct {
	User         *AdminUserResponse `json:"user"`
	CreatedUsers int64              `json:"created_users"`
	CreatedRoles int64              `json:"created_roles"`
	TotalLogins  int64              `json:"total_logins"`
	RecentLogins []time.Time        `json:"recent_logins"`
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
	Role        *AdminRoleResponse  `json:"role"`
	UserCount   int64               `json:"user_count"`
	ActiveUsers int64               `json:"active_users"`
	Users       []AdminUserResponse `json:"users"`
}

// UserStatsResponse represents statistics about admin users
type UserStatsResponse struct {
	TotalUsers     int64           `json:"total_users"`
	ActiveUsers    int64           `json:"active_users"`
	InactiveUsers  int64           `json:"inactive_users"`
	SuspendedUsers int64           `json:"suspended_users"`
	Users2FA       int64           `json:"users_2fa"`
	ByRole         map[int64]int64 `json:"by_role"`
	SuperAdmins    int64           `json:"super_admins"`
	Admins         int64           `json:"admins"`
	Moderators     int64           `json:"moderators"`
	RecentLogins   int64           `json:"recent_logins"`
}

// ActivityStatsResponse represents activity statistics
type ActivityStatsResponse struct {
	Period         string              `json:"period"`
	TotalLogins    int64               `json:"total_logins"`
	UniqueUsers    int64               `json:"unique_users"`
	AverageLogins  float64             `json:"average_logins"`
	LoginsByDate   map[string]int64    `json:"logins_by_date"`
	LoginsByUser   map[int64]int64     `json:"logins_by_user"`
	TopActiveUsers []AdminUserResponse `json:"top_active_users"`
}

// UserActivityResponse represents activity information for a user
type UserActivityResponse struct {
	User         *AdminUserResponse `json:"user"`
	TotalLogins  int64              `json:"total_logins"`
	LastLogin    *time.Time         `json:"last_login,omitempty"`
	CreatedUsers int64              `json:"created_users"`
	CreatedRoles int64              `json:"created_roles"`
	LoginHistory []LoginRecord      `json:"login_history"`
}

// LoginRecord represents a single login record
type LoginRecord struct {
	UserID    int64     `json:"user_id"`
	LoginTime time.Time `json:"login_time"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
}

// =============================================================================
// Admin Company Management DTOs
// =============================================================================

// AdminCompanyListRequest represents request for listing companies (Task 2.1)
type AdminCompanyListRequest struct {
	Page               int
	Limit              int
	Search             string // Search by company name, email, legal name
	Status             string // verification status filter
	VerificationStatus string // From company_verifications table
	IndustryID         *int64
	CompanySizeID      *int64
	ProvinceID         *int64
	CityID             *int64
	Verified           *bool
	IsActive           *bool
	CreatedFrom        string // Date filter
	CreatedTo          string // Date filter
	SortBy             string // company_name, created_at, verified_at, updated_at
	SortOrder          string // asc, desc
}

// AdminCompanyListResponse represents paginated company list response
type AdminCompanyListResponse struct {
	Companies  []AdminCompanyListItem `json:"companies"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
	HasNext    bool                   `json:"has_next"`
	HasPrev    bool                   `json:"has_prev"`
}

// AdminCompanyListItem represents a company item in the list
type AdminCompanyListItem struct {
	ID                 int64      `json:"id"`
	UUID               string     `json:"uuid"`
	CompanyName        string     `json:"company_name"`
	Slug               string     `json:"slug"`
	LegalName          string     `json:"legal_name,omitempty"`
	RegistrationNumber string     `json:"registration_number,omitempty"`
	Industry           string     `json:"industry,omitempty"`
	CompanySize        string     `json:"company_size,omitempty"`
	Location           string     `json:"location,omitempty"`
	Verified           bool       `json:"verified"`
	VerifiedAt         *time.Time `json:"verified_at,omitempty"`
	IsActive           bool       `json:"is_active"`
	VerificationStatus string     `json:"verification_status"`
	TotalJobs          int64      `json:"total_jobs"`
	ActiveJobs         int64      `json:"active_jobs"`
	TotalApplications  int64      `json:"total_applications"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// AdminCompanyDetailResponse represents full company details (Task 2.2)
type AdminCompanyDetailResponse struct {
	ID                 int64      `json:"id"`
	UUID               string     `json:"uuid"`
	CompanyName        string     `json:"company_name"`
	Slug               string     `json:"slug"`
	LegalName          string     `json:"legal_name,omitempty"`
	RegistrationNumber string     `json:"registration_number,omitempty"`
	IndustryID         *int64     `json:"industry_id,omitempty"`
	CompanySizeID      *int64     `json:"company_size_id,omitempty"`
	WebsiteURL         string     `json:"website_url,omitempty"`
	EmailDomain        string     `json:"email_domain,omitempty"`
	Phone              string     `json:"phone,omitempty"`
	FullAddress        string     `json:"full_address,omitempty"`
	Description        string     `json:"description,omitempty"`
	About              string     `json:"about,omitempty"`
	Culture            string     `json:"culture,omitempty"`
	Verified           bool       `json:"verified"`
	VerifiedAt         *time.Time `json:"verified_at,omitempty"`
	IsActive           bool       `json:"is_active"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// Extended info
	VerificationDetail *CompanyVerificationDetail `json:"verification_detail,omitempty"`
	Documents          []CompanyDocumentDetail    `json:"documents,omitempty"`
	OwnerInfo          *CompanyOwnerInfo          `json:"owner_info,omitempty"`
	Stats              *AdminCompanyStatsResponse `json:"stats,omitempty"`
}

// CompanyVerificationDetail represents verification details
type CompanyVerificationDetail struct {
	ID                int64      `json:"id"`
	Status            string     `json:"status"`
	VerificationScore float64    `json:"verification_score"`
	VerificationNotes string     `json:"verification_notes,omitempty"`
	RejectionReason   string     `json:"rejection_reason,omitempty"`
	ReviewedBy        *int64     `json:"reviewed_by,omitempty"`
	ReviewedByName    string     `json:"reviewed_by_name,omitempty"`
	ReviewedAt        *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

// CompanyDocumentDetail represents document details
type CompanyDocumentDetail struct {
	ID              int64      `json:"id"`
	DocumentType    string     `json:"document_type"`
	DocumentNumber  string     `json:"document_number,omitempty"`
	FilePath        string     `json:"file_path"`
	Status          string     `json:"status"`
	VerifiedBy      *int64     `json:"verified_by,omitempty"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty"`
	RejectionReason string     `json:"rejection_reason,omitempty"`
}

// CompanyOwnerInfo represents company owner information
type CompanyOwnerInfo struct {
	UserID   int64  `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone,omitempty"`
	Role     string `json:"role"`
}

// AdminCompanyStatusRequest represents request to update company status (Task 2.3)
type AdminCompanyStatusRequest struct {
	Status          string  `json:"status" validate:"required,oneof=pending_verification verified rejected suspended blacklisted"`
	RejectionReason *string `json:"rejection_reason" validate:"omitempty,max=1000"`
	Notes           *string `json:"notes" validate:"omitempty,max=2000"`
	GrantBadge      *bool   `json:"grant_badge"`
}

// AdminUpdateCompanyRequest represents admin request to update company (Task 2.4)
type AdminUpdateCompanyRequest struct {
	CompanyName        *string `json:"company_name" validate:"omitempty,min=2,max=200"`
	LegalName          *string `json:"legal_name" validate:"omitempty,max=200"`
	RegistrationNumber *string `json:"registration_number" validate:"omitempty,max=100"`
	IndustryID         *int64  `json:"industry_id"`
	CompanySizeID      *int64  `json:"company_size_id"`
	DistrictID         *int64  `json:"district_id"`
	FullAddress        *string `json:"full_address" validate:"omitempty,max=500"`
	Description        *string `json:"description" validate:"omitempty,max=2000"`
	WebsiteURL         *string `json:"website_url" validate:"omitempty,url"`
	EmailDomain        *string `json:"email_domain" validate:"omitempty,max=100"`
	Phone              *string `json:"phone" validate:"omitempty,min=10,max=30"`
	About              *string `json:"about" validate:"omitempty,max=5000"`
	Culture            *string `json:"culture" validate:"omitempty,max=5000"`
	IsActive           *bool   `json:"is_active"`
	Verified           *bool   `json:"verified"`
}

// AdminDeleteCompanyRequest represents request to delete company (Task 2.5)
type AdminDeleteCompanyRequest struct {
	Force  bool   `json:"force"` // Force delete even with active jobs
	Reason string `json:"reason" validate:"required,min=10,max=500"`
}

// AdminCompanyStatsResponse represents company statistics
type AdminCompanyStatsResponse struct {
	TotalJobs           int64   `json:"total_jobs"`
	ActiveJobs          int64   `json:"active_jobs"`
	ClosedJobs          int64   `json:"closed_jobs"`
	DraftJobs           int64   `json:"draft_jobs"`
	TotalApplications   int64   `json:"total_applications"`
	PendingApplications int64   `json:"pending_applications"`
	TotalFollowers      int64   `json:"total_followers"`
	TotalEmployees      int64   `json:"total_employees"`
	TotalReviews        int64   `json:"total_reviews"`
	AverageRating       float64 `json:"average_rating"`
}

// AdminDashboardStatsResponse represents overall dashboard statistics
type AdminDashboardStatsResponse struct {
	TotalCompanies        int64 `json:"total_companies"`
	VerifiedCompanies     int64 `json:"verified_companies"`
	PendingVerification   int64 `json:"pending_verification"`
	RejectedCompanies     int64 `json:"rejected_companies"`
	SuspendedCompanies    int64 `json:"suspended_companies"`
	NewCompaniesThisMonth int64 `json:"new_companies_this_month"`
	NewCompaniesToday     int64 `json:"new_companies_today"`
	TotalJobs             int64 `json:"total_jobs"`
	ActiveJobs            int64 `json:"active_jobs"`
	TotalApplications     int64 `json:"total_applications"`
}

// BulkOperationResult represents result of bulk operations
type BulkOperationResult struct {
	SuccessCount int64       `json:"success_count"`
	FailedCount  int64       `json:"failed_count"`
	Errors       []BulkError `json:"errors,omitempty"`
}

// BulkError represents error in bulk operation
type BulkError struct {
	CompanyID int64  `json:"company_id"`
	Error     string `json:"error"`
}

// AuditLogListResponse represents paginated audit log list
type AuditLogListResponse struct {
	Logs       []AuditLogEntry `json:"logs"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// AuditLogEntry represents single audit log entry
type AuditLogEntry struct {
	ID          int64     `json:"id"`
	CompanyID   int64     `json:"company_id"`
	CompanyName string    `json:"company_name"`
	AdminID     int64     `json:"admin_id"`
	AdminName   string    `json:"admin_name"`
	Action      string    `json:"action"` // created, updated, verified, rejected, suspended, deleted
	Description string    `json:"description"`
	OldValue    string    `json:"old_value,omitempty"`
	NewValue    string    `json:"new_value,omitempty"`
	IPAddress   string    `json:"ip_address,omitempty"`
	UserAgent   string    `json:"user_agent,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
