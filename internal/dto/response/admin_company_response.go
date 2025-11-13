package response

import "time"

// AdminCompanyListItemResponse represents a company item in admin list view
type AdminCompanyListItemResponse struct {
	ID                 int64      `json:"id"`
	UUID               string     `json:"uuid"`
	CompanyName        string     `json:"company_name"`
	Slug               string     `json:"slug"`
	LegalName          string     `json:"legal_name,omitempty"`
	RegistrationNumber string     `json:"registration_number,omitempty"`
	Email              string     `json:"email,omitempty"` // From owner/creator
	Phone              string     `json:"phone,omitempty"`
	Industry           string     `json:"industry,omitempty"`
	CompanySize        string     `json:"company_size,omitempty"`
	Location           string     `json:"location,omitempty"` // City, Province
	Verified           bool       `json:"verified"`
	VerifiedAt         *time.Time `json:"verified_at,omitempty"`
	IsActive           bool       `json:"is_active"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// Verification info
	VerificationStatus string     `json:"verification_status"` // pending, verified, rejected, suspended
	VerifiedBy         *int64     `json:"verified_by,omitempty"`
	VerifierName       string     `json:"verifier_name,omitempty"`
	ReviewedAt         *time.Time `json:"reviewed_at,omitempty"`

	// Stats (quick view)
	TotalJobs         int64 `json:"total_jobs"`
	ActiveJobs        int64 `json:"active_jobs"`
	TotalApplications int64 `json:"total_applications"`
	TotalFollowers    int64 `json:"total_followers"`
}

// AdminCompanyDetailResponse represents full company details for admin
type AdminCompanyDetailResponse struct {
	// Basic Info
	ID                 int64  `json:"id"`
	UUID               string `json:"uuid"`
	CompanyName        string `json:"company_name"`
	Slug               string `json:"slug"`
	LegalName          string `json:"legal_name,omitempty"`
	RegistrationNumber string `json:"registration_number,omitempty"`

	// Master Data Relations
	IndustryID    *int64                     `json:"industry_id,omitempty"`
	IndustryInfo  *MasterIndustryResponse    `json:"industry_info,omitempty"`
	CompanySizeID *int64                     `json:"company_size_id,omitempty"`
	CompanySize   *MasterCompanySizeResponse `json:"company_size,omitempty"`
	LocationInfo  *CompanyLocationResponse   `json:"location_info,omitempty"`

	// Legacy fields
	Industry     string `json:"industry,omitempty"`
	CompanyType  string `json:"company_type,omitempty"`
	SizeCategory string `json:"size_category,omitempty"`

	// Contact & Address
	WebsiteURL  string   `json:"website_url,omitempty"`
	EmailDomain string   `json:"email_domain,omitempty"`
	Phone       string   `json:"phone,omitempty"`
	FullAddress string   `json:"full_address,omitempty"`
	Country     string   `json:"country"`
	PostalCode  string   `json:"postal_code,omitempty"`
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`

	// Media
	LogoURL   string `json:"logo_url,omitempty"`
	BannerURL string `json:"banner_url,omitempty"`

	// Content
	Description string   `json:"description,omitempty"`
	About       string   `json:"about,omitempty"`
	Culture     string   `json:"culture,omitempty"`
	Benefits    []string `json:"benefits,omitempty"`

	// Verification & Status
	Verified     bool                   `json:"verified"`
	VerifiedAt   *time.Time             `json:"verified_at,omitempty"`
	VerifiedBy   *int64                 `json:"verified_by,omitempty"`
	VerifierInfo *AdminUserInfoResponse `json:"verifier_info,omitempty"`
	IsActive     bool                   `json:"is_active"`

	// Verification Detail
	VerificationDetail *AdminCompanyVerificationDetailResponse `json:"verification_detail,omitempty"`

	// Profile
	Profile *CompanyProfileResponse `json:"profile,omitempty"`

	// Documents
	Documents []AdminCompanyDocumentResponse `json:"documents,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Owner/Creator Info
	CreatorInfo *AdminCompanyCreatorResponse `json:"creator_info,omitempty"`

	// Employer Users (Team members)
	EmployerUsers []AdminEmployerUserResponse `json:"employer_users,omitempty"`

	// Statistics
	Stats AdminCompanyStatsResponse `json:"stats"`
}

// AdminCompanyVerificationDetailResponse represents verification details
type AdminCompanyVerificationDetailResponse struct {
	ID                 int64      `json:"id"`
	Status             string     `json:"status"` // pending, under_review, verified, rejected, blacklisted, expired
	VerificationScore  float64    `json:"verification_score"`
	VerificationNotes  string     `json:"verification_notes,omitempty"`
	RejectionReason    string     `json:"rejection_reason,omitempty"`
	VerificationExpiry *time.Time `json:"verification_expiry,omitempty"`
	BadgeGranted       bool       `json:"badge_granted"`
	AutoExpired        bool       `json:"auto_expired"`
	LastChecked        *time.Time `json:"last_checked,omitempty"`
	RequestedBy        *int64     `json:"requested_by,omitempty"`
	RequestedByName    string     `json:"requested_by_name,omitempty"`
	ReviewedBy         *int64     `json:"reviewed_by,omitempty"`
	ReviewedByName     string     `json:"reviewed_by_name,omitempty"`
	ReviewedAt         *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// AdminCompanyDocumentResponse represents company document for admin view
type AdminCompanyDocumentResponse struct {
	ID              int64      `json:"id"`
	DocumentType    string     `json:"document_type"`
	DocumentNumber  string     `json:"document_number,omitempty"`
	DocumentName    string     `json:"document_name,omitempty"`
	FilePath        string     `json:"file_path"`
	FileURL         string     `json:"file_url,omitempty"`
	IssueDate       *time.Time `json:"issue_date,omitempty"`
	ExpiryDate      *time.Time `json:"expiry_date,omitempty"`
	Status          string     `json:"status"` // pending, approved, rejected, expired
	VerifiedBy      *int64     `json:"verified_by,omitempty"`
	VerifiedByName  string     `json:"verified_by_name,omitempty"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty"`
	RejectionReason string     `json:"rejection_reason,omitempty"`
	IsActive        bool       `json:"is_active"`
	IsExpired       bool       `json:"is_expired"`
	UploadedBy      *int64     `json:"uploaded_by,omitempty"`
	UploadedByName  string     `json:"uploaded_by_name,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// AdminCompanyCreatorResponse represents company creator/owner info
type AdminCompanyCreatorResponse struct {
	UserID    int64     `json:"user_id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	Role      string    `json:"role"` // owner, admin, recruiter
	CreatedAt time.Time `json:"created_at"`
}

// AdminEmployerUserResponse represents employer user in company
type AdminEmployerUserResponse struct {
	ID            int64      `json:"id"`
	UserID        int64      `json:"user_id"`
	FullName      string     `json:"full_name"`
	Email         string     `json:"email"`
	Role          string     `json:"role"` // owner, admin, recruiter, viewer
	PositionTitle string     `json:"position_title,omitempty"`
	Department    string     `json:"department,omitempty"`
	IsVerified    bool       `json:"is_verified"`
	IsActive      bool       `json:"is_active"`
	LastLogin     *time.Time `json:"last_login,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

// AdminCompanyStatsResponse represents company statistics for admin
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
	ApprovedReviews     int64   `json:"approved_reviews"`
	PendingReviews      int64   `json:"pending_reviews"`
}

// AdminUserInfoResponse represents admin user who performed an action
type AdminUserInfoResponse struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role,omitempty"`
}

// AdminCompaniesListResponse represents paginated list of companies for admin
type AdminCompaniesListResponse struct {
	Companies []AdminCompanyListItemResponse `json:"companies"`
	Meta      PaginationMeta                 `json:"meta"`
}

// AdminCompanyStatusUpdateResponse represents response after status update
type AdminCompanyStatusUpdateResponse struct {
	ID                 int64      `json:"id"`
	CompanyName        string     `json:"company_name"`
	Status             string     `json:"status"`
	VerificationStatus string     `json:"verification_status"`
	Verified           bool       `json:"verified"`
	VerifiedAt         *time.Time `json:"verified_at,omitempty"`
	ReviewedBy         int64      `json:"reviewed_by"`
	ReviewedByName     string     `json:"reviewed_by_name"`
	ReviewedAt         time.Time  `json:"reviewed_at"`
	Notes              string     `json:"notes,omitempty"`
	RejectionReason    string     `json:"rejection_reason,omitempty"`
	BadgeGranted       bool       `json:"badge_granted"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// AdminCompanyDeleteResponse represents response after company deletion
type AdminCompanyDeleteResponse struct {
	ID          int64     `json:"id"`
	CompanyName string    `json:"company_name"`
	Reason      string    `json:"reason"`
	DeletedBy   int64     `json:"deleted_by"`
	DeletedAt   time.Time `json:"deleted_at"`
	Message     string    `json:"message"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

// AdminDashboardStatsResponse represents dashboard statistics
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
