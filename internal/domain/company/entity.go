package company

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Company represents the main company entity
type Company struct {
	ID                 int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID               uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex" json:"uuid"`
	CompanyName        string         `gorm:"type:varchar(200);not null" json:"company_name" validate:"required,min=2,max=200"`
	Slug               string         `gorm:"type:varchar(200);uniqueIndex;not null" json:"slug" validate:"required"`
	LegalName          *string        `gorm:"type:varchar(200)" json:"legal_name,omitempty"`
	RegistrationNumber *string        `gorm:"type:varchar(100)" json:"registration_number,omitempty"`
	Industry           *string        `gorm:"type:varchar(100)" json:"industry,omitempty"`
	CompanyType        *string        `gorm:"type:varchar(50);check:company_type IN ('private','public','startup','ngo','government')" json:"company_type,omitempty" validate:"omitempty,oneof=private public startup ngo government"`
	SizeCategory       *string        `gorm:"type:varchar(50);check:size_category IN ('1-10','11-50','51-200','201-1000','1000+')" json:"size_category,omitempty" validate:"omitempty,oneof=1-10 11-50 51-200 201-1000 1000+"`
	WebsiteURL         *string        `gorm:"type:text" json:"website_url,omitempty" validate:"omitempty,url"`
	EmailDomain        *string        `gorm:"type:varchar(100)" json:"email_domain,omitempty"`
	Phone              *string        `gorm:"type:varchar(30)" json:"phone,omitempty"`
	Address            *string        `gorm:"type:text" json:"address,omitempty"`
	City               *string        `gorm:"type:varchar(100)" json:"city,omitempty"`
	Province           *string        `gorm:"type:varchar(100)" json:"province,omitempty"`
	Country            string         `gorm:"type:varchar(100);default:'Indonesia'" json:"country"`
	PostalCode         *string        `gorm:"type:varchar(10)" json:"postal_code,omitempty"`
	Latitude           *float64       `gorm:"type:numeric(10,6)" json:"latitude,omitempty"`
	Longitude          *float64       `gorm:"type:numeric(10,6)" json:"longitude,omitempty"`
	LogoURL            *string        `gorm:"type:text" json:"logo_url,omitempty"`
	BannerURL          *string        `gorm:"type:text" json:"banner_url,omitempty"`
	About              *string        `gorm:"type:text" json:"about,omitempty"`
	Culture            *string        `gorm:"type:text" json:"culture,omitempty"`
	Benefits           pq.StringArray `gorm:"type:text[]" json:"benefits,omitempty"` // PostgreSQL array
	Verified           bool           `gorm:"default:false" json:"verified"`
	VerifiedAt         *time.Time     `gorm:"type:timestamp" json:"verified_at,omitempty"`
	VerifiedBy         *int64         `gorm:"type:bigint" json:"verified_by,omitempty"`
	IsActive           bool           `gorm:"default:true" json:"is_active"`
	CreatedAt          time.Time      `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	Profile       *CompanyProfile      `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"profile,omitempty"`
	Followers     []CompanyFollower    `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"followers,omitempty"`
	Reviews       []CompanyReview      `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"reviews,omitempty"`
	Documents     []CompanyDocument    `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"documents,omitempty"`
	Employees     []CompanyEmployee    `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"employees,omitempty"`
	EmployerUsers []EmployerUser       `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"employer_users,omitempty"`
	Verification  *CompanyVerification `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"verification,omitempty"`
}

// TableName specifies the table name for Company
func (Company) TableName() string {
	return "companies"
}

// IsVerified checks if company is verified
func (c *Company) IsVerified() bool {
	return c.Verified
}

// IsStartup checks if company is a startup
func (c *Company) IsStartup() bool {
	return c.CompanyType != nil && *c.CompanyType == "startup"
}

// CompanyProfile represents detailed company profile information
type CompanyProfile struct {
	ID               int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID        int64      `gorm:"not null;uniqueIndex" json:"company_id"`
	Tagline          *string    `gorm:"type:varchar(200)" json:"tagline,omitempty" validate:"omitempty,max=200"`
	ShortDescription *string    `gorm:"type:text" json:"short_description,omitempty"`
	LongDescription  *string    `gorm:"type:text" json:"long_description,omitempty"`
	Mission          *string    `gorm:"type:text" json:"mission,omitempty"`
	Vision           *string    `gorm:"type:text" json:"vision,omitempty"`
	Values           *string    `gorm:"type:text[]" json:"values,omitempty"` // PostgreSQL array
	Culture          *string    `gorm:"type:text" json:"culture,omitempty"`
	WorkEnvironment  *string    `gorm:"type:text" json:"work_environment,omitempty"`
	GalleryURLs      *string    `gorm:"type:text[]" json:"gallery_urls,omitempty"` // PostgreSQL array
	VideoURL         *string    `gorm:"type:text" json:"video_url,omitempty"`
	Awards           *string    `gorm:"type:text[]" json:"awards,omitempty"`      // PostgreSQL array
	SocialLinks      *string    `gorm:"type:jsonb" json:"social_links,omitempty"` // JSONB for social media links
	HiringTagline    *string    `gorm:"type:varchar(200)" json:"hiring_tagline,omitempty"`
	SEOTitle         *string    `gorm:"type:varchar(200)" json:"seo_title,omitempty"`
	SEOKeywords      *string    `gorm:"type:text[]" json:"seo_keywords,omitempty"` // PostgreSQL array
	SEODescription   *string    `gorm:"type:text" json:"seo_description,omitempty"`
	Status           string     `gorm:"type:varchar(20);default:'draft';check:status IN ('draft','published','suspended')" json:"status" validate:"oneof=draft published suspended"`
	Verified         bool       `gorm:"default:false" json:"verified"`
	VerifiedAt       *time.Time `gorm:"type:timestamp" json:"verified_at,omitempty"`
	VerifiedBy       *int64     `gorm:"type:bigint" json:"verified_by,omitempty"`
	CreatedAt        time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	Company *Company `gorm:"foreignKey:CompanyID" json:"-"`
}

// TableName specifies the table name for CompanyProfile
func (CompanyProfile) TableName() string {
	return "company_profiles"
}

// IsPublished checks if profile is published
func (cp *CompanyProfile) IsPublished() bool {
	return cp.Status == "published"
}

// CompanyIndustry represents industry classifications
type CompanyIndustry struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"code" validate:"required"`
	Name        string    `gorm:"type:varchar(150);uniqueIndex;not null" json:"name" validate:"required,max=150"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	ParentID    *int64    `gorm:"type:bigint" json:"parent_id,omitempty"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	Parent   *CompanyIndustry  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []CompanyIndustry `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName specifies the table name for CompanyIndustry
func (CompanyIndustry) TableName() string {
	return "company_industries"
}

// CompanyFollower represents users following companies
type CompanyFollower struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID    int64      `gorm:"not null;uniqueIndex:idx_company_user" json:"company_id"`
	UserID       int64      `gorm:"not null;uniqueIndex:idx_company_user" json:"user_id"`
	FollowedAt   time.Time  `gorm:"type:timestamp;default:now()" json:"followed_at"`
	UnfollowedAt *time.Time `gorm:"type:timestamp" json:"unfollowed_at,omitempty"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	Company *Company `gorm:"foreignKey:CompanyID" json:"-"`
}

// TableName specifies the table name for CompanyFollower
func (CompanyFollower) TableName() string {
	return "company_followers"
}

// CompanyReview represents company reviews from employees/ex-employees
type CompanyReview struct {
	ID                 int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID          int64      `gorm:"not null;index" json:"company_id"`
	UserID             *int64     `gorm:"type:bigint" json:"user_id,omitempty"`
	ReviewerType       *string    `gorm:"type:varchar(30);check:reviewer_type IN ('employee','ex-employee','applicant')" json:"reviewer_type,omitempty" validate:"omitempty,oneof=employee ex-employee applicant"`
	PositionTitle      *string    `gorm:"type:varchar(100)" json:"position_title,omitempty"`
	EmploymentPeriod   *string    `gorm:"type:varchar(50)" json:"employment_period,omitempty"`
	RatingOverall      *float64   `gorm:"type:numeric(2,1);check:rating_overall >= 0 AND rating_overall <= 5" json:"rating_overall,omitempty" validate:"omitempty,min=0,max=5"`
	RatingCulture      *float64   `gorm:"type:numeric(2,1)" json:"rating_culture,omitempty" validate:"omitempty,min=0,max=5"`
	RatingWorkLife     *float64   `gorm:"type:numeric(2,1)" json:"rating_worklife,omitempty" validate:"omitempty,min=0,max=5"`
	RatingSalary       *float64   `gorm:"type:numeric(2,1)" json:"rating_salary,omitempty" validate:"omitempty,min=0,max=5"`
	RatingManagement   *float64   `gorm:"type:numeric(2,1)" json:"rating_management,omitempty" validate:"omitempty,min=0,max=5"`
	Pros               *string    `gorm:"type:text" json:"pros,omitempty"`
	Cons               *string    `gorm:"type:text" json:"cons,omitempty"`
	AdviceToManagement *string    `gorm:"type:text" json:"advice_to_management,omitempty"`
	IsAnonymous        bool       `gorm:"default:true" json:"is_anonymous"`
	RecommendToFriend  bool       `gorm:"default:true" json:"recommend_to_friend"`
	Status             string     `gorm:"type:varchar(20);default:'pending';check:status IN ('pending','approved','rejected','hidden')" json:"status"`
	ModeratedBy        *int64     `gorm:"type:bigint" json:"moderated_by,omitempty"`
	ModeratedAt        *time.Time `gorm:"type:timestamp" json:"moderated_at,omitempty"`
	CreatedAt          time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	Company *Company `gorm:"foreignKey:CompanyID" json:"-"`
}

// TableName specifies the table name for CompanyReview
func (CompanyReview) TableName() string {
	return "company_reviews"
}

// IsApproved checks if review is approved
func (cr *CompanyReview) IsApproved() bool {
	return cr.Status == "approved"
}

// CompanyDocument represents company legal documents
type CompanyDocument struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID       int64      `gorm:"not null;uniqueIndex:idx_company_doc_type_num" json:"company_id"`
	UploadedBy      *int64     `gorm:"type:bigint" json:"uploaded_by,omitempty"`
	DocumentType    string     `gorm:"type:varchar(50);not null;uniqueIndex:idx_company_doc_type_num;check:document_type IN ('SIUP','NPWP','NIB','AKTA','TDP','ISO','SERTIFIKAT','LAINNYA')" json:"document_type" validate:"required,oneof=SIUP NPWP NIB AKTA TDP ISO SERTIFIKAT LAINNYA"`
	DocumentNumber  *string    `gorm:"type:varchar(100);uniqueIndex:idx_company_doc_type_num" json:"document_number,omitempty"`
	DocumentName    *string    `gorm:"type:varchar(150)" json:"document_name,omitempty"`
	FilePath        string     `gorm:"type:text;not null" json:"file_path" validate:"required"`
	IssueDate       *time.Time `gorm:"type:date" json:"issue_date,omitempty"`
	ExpiryDate      *time.Time `gorm:"type:date" json:"expiry_date,omitempty"`
	Status          string     `gorm:"type:varchar(20);default:'pending';check:status IN ('pending','approved','rejected','expired')" json:"status"`
	VerifiedBy      *int64     `gorm:"type:bigint" json:"verified_by,omitempty"`
	VerifiedAt      *time.Time `gorm:"type:timestamp" json:"verified_at,omitempty"`
	RejectionReason *string    `gorm:"type:text" json:"rejection_reason,omitempty"`
	IsActive        bool       `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	Company *Company `gorm:"foreignKey:CompanyID" json:"-"`
}

// TableName specifies the table name for CompanyDocument
func (CompanyDocument) TableName() string {
	return "company_documents"
}

// IsApproved checks if document is approved
func (cd *CompanyDocument) IsApproved() bool {
	return cd.Status == "approved"
}

// IsExpired checks if document is expired
func (cd *CompanyDocument) IsExpired() bool {
	if cd.ExpiryDate == nil {
		return false
	}
	return cd.ExpiryDate.Before(time.Now())
}

// CompanyEmployee represents company employees
type CompanyEmployee struct {
	ID               int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID        int64      `gorm:"not null;index" json:"company_id"`
	UserID           *int64     `gorm:"type:bigint" json:"user_id,omitempty"`
	FullName         *string    `gorm:"type:varchar(150)" json:"full_name,omitempty"`
	JobTitle         *string    `gorm:"type:varchar(100)" json:"job_title,omitempty"`
	Department       *string    `gorm:"type:varchar(100)" json:"department,omitempty"`
	EmploymentType   string     `gorm:"type:varchar(30);default:'permanent';check:employment_type IN ('permanent','contract','intern','freelance')" json:"employment_type"`
	EmploymentStatus string     `gorm:"type:varchar(30);default:'active';check:employment_status IN ('active','resigned','terminated','on_leave')" json:"employment_status"`
	JoinDate         *time.Time `gorm:"type:date" json:"join_date,omitempty"`
	EndDate          *time.Time `gorm:"type:date" json:"end_date,omitempty"`
	SalaryRangeMin   *float64   `gorm:"type:numeric(12,2)" json:"salary_range_min,omitempty"`
	SalaryRangeMax   *float64   `gorm:"type:numeric(12,2)" json:"salary_range_max,omitempty"`
	AddedBy          *int64     `gorm:"type:bigint" json:"added_by,omitempty"`
	Note             *string    `gorm:"type:text" json:"note,omitempty"`
	IsVisiblePublic  bool       `gorm:"default:false" json:"is_visible_public"`
	Verified         bool       `gorm:"default:false" json:"verified"`
	VerifiedAt       *time.Time `gorm:"type:timestamp" json:"verified_at,omitempty"`
	VerifiedBy       *int64     `gorm:"type:bigint" json:"verified_by,omitempty"`
	CreatedAt        time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	Company *Company `gorm:"foreignKey:CompanyID" json:"-"`
}

// TableName specifies the table name for CompanyEmployee
func (CompanyEmployee) TableName() string {
	return "company_employees"
}

// IsActive checks if employee is currently active
func (ce *CompanyEmployee) IsActive() bool {
	return ce.EmploymentStatus == "active"
}

// CompanyVerification represents company verification status
type CompanyVerification struct {
	ID                 int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID          int64      `gorm:"not null;uniqueIndex" json:"company_id"`
	RequestedBy        *int64     `gorm:"type:bigint" json:"requested_by,omitempty"`
	ReviewedBy         *int64     `gorm:"type:bigint" json:"reviewed_by,omitempty"`
	ReviewedAt         *time.Time `gorm:"type:timestamp" json:"reviewed_at,omitempty"`
	Status             string     `gorm:"type:varchar(20);default:'pending';check:status IN ('pending','under_review','verified','rejected','blacklisted','expired')" json:"status"`
	VerificationScore  float64    `gorm:"type:numeric(5,2);default:0.00" json:"verification_score"`
	VerificationNotes  *string    `gorm:"type:text" json:"verification_notes,omitempty"`
	RejectionReason    *string    `gorm:"type:text" json:"rejection_reason,omitempty"`
	VerificationExpiry *time.Time `gorm:"type:date" json:"verification_expiry,omitempty"`
	BadgeGranted       bool       `gorm:"default:false" json:"badge_granted"`
	AutoExpired        bool       `gorm:"default:false" json:"auto_expired"`
	LastChecked        *time.Time `gorm:"type:timestamp" json:"last_checked,omitempty"`
	CreatedAt          time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	Company *Company `gorm:"foreignKey:CompanyID" json:"-"`
}

// TableName specifies the table name for CompanyVerification
func (CompanyVerification) TableName() string {
	return "company_verifications"
}

// IsVerified checks if verification is approved
func (cv *CompanyVerification) IsVerified() bool {
	return cv.Status == "verified"
}

// IsExpired checks if verification is expired
func (cv *CompanyVerification) IsExpired() bool {
	if cv.VerificationExpiry == nil {
		return false
	}
	return cv.VerificationExpiry.Before(time.Now())
}

// EmployerUser represents users with employer privileges
type EmployerUser struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64      `gorm:"not null;uniqueIndex:idx_user_company" json:"user_id"`
	CompanyID     int64      `gorm:"not null;uniqueIndex:idx_user_company" json:"company_id"`
	Role          string     `gorm:"type:varchar(30);default:'recruiter';check:role IN ('owner','admin','recruiter','viewer')" json:"role" validate:"oneof=owner admin recruiter viewer"`
	PositionTitle *string    `gorm:"type:varchar(100)" json:"position_title,omitempty"`
	Department    *string    `gorm:"type:varchar(100)" json:"department,omitempty"`
	EmailCompany  *string    `gorm:"type:varchar(150)" json:"email_company,omitempty" validate:"omitempty,email"`
	PhoneCompany  *string    `gorm:"type:varchar(30)" json:"phone_company,omitempty"`
	IsVerified    bool       `gorm:"default:false" json:"is_verified"`
	VerifiedAt    *time.Time `gorm:"type:timestamp" json:"verified_at,omitempty"`
	VerifiedBy    *int64     `gorm:"type:bigint" json:"verified_by,omitempty"`
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	LastLogin     *time.Time `gorm:"type:timestamp" json:"last_login,omitempty"`
	CreatedAt     time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	Company *Company `gorm:"foreignKey:CompanyID" json:"-"`
}

// TableName specifies the table name for EmployerUser
func (EmployerUser) TableName() string {
	return "employer_users"
}

// IsOwner checks if user is company owner
func (eu *EmployerUser) IsOwner() bool {
	return eu.Role == "owner"
}

// IsAdmin checks if user is company admin
func (eu *EmployerUser) IsAdmin() bool {
	return eu.Role == "admin" || eu.Role == "owner"
}

// CanManageJobs checks if user can manage jobs
func (eu *EmployerUser) CanManageJobs() bool {
	return eu.Role == "owner" || eu.Role == "admin" || eu.Role == "recruiter"
}
