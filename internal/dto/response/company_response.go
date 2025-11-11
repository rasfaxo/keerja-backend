package response

import "time"

// CompanyResponse represents company public response
type CompanyResponse struct {
	ID          int64  `json:"id"`
	UUID        string `json:"uuid"`
	CompanyName string `json:"company_name"`
	Slug        string `json:"slug"`

	// Master Data Relations (new structure)
	IndustryDetail    *MasterIndustryResponse    `json:"industry_detail,omitempty"`
	CompanySizeDetail *MasterCompanySizeResponse `json:"company_size_detail,omitempty"`
	LocationDetail    *CompanyLocationResponse   `json:"location_detail,omitempty"`
	FullAddress       string                     `json:"full_address,omitempty"`
	Description       string                     `json:"description,omitempty"`

	// Legacy Fields (for backward compatibility)
	Industry     string `json:"industry,omitempty"`
	CompanyType  string `json:"company_type,omitempty"`
	SizeCategory string `json:"size_category,omitempty"`
	City         string `json:"city,omitempty"`
	Province     string `json:"province,omitempty"`

	WebsiteURL string     `json:"website_url,omitempty"`
	Phone      string     `json:"phone,omitempty"`
	Country    string     `json:"country"`
	LogoURL    string     `json:"logo_url,omitempty"`
	BannerURL  string     `json:"banner_url,omitempty"`
	About      string     `json:"about,omitempty"`
	Verified   bool       `json:"verified"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`

	// Verification Status Details
	Status       string `json:"status,omitempty"`        // "verified", "pending", "under_review", "rejected", "not_requested"
	BadgeGranted bool   `json:"badge_granted,omitempty"` // Whether company has verification badge

	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Stats
	FollowersCount int64   `json:"followers_count"`
	JobsCount      int64   `json:"jobs_count"`
	AverageRating  float64 `json:"average_rating"`
	ReviewsCount   int64   `json:"reviews_count"`
}

// CompanyDetailResponse represents detailed company response
type CompanyDetailResponse struct {
	ID                 int64  `json:"id"`
	UUID               string `json:"uuid"`
	CompanyName        string `json:"company_name"`
	Slug               string `json:"slug"`
	LegalName          string `json:"legal_name,omitempty"`
	RegistrationNumber string `json:"registration_number,omitempty"`

	// Master Data Relations
	IndustryDetail    *MasterIndustryResponse    `json:"industry_detail,omitempty"`
	CompanySizeDetail *MasterCompanySizeResponse `json:"company_size_detail,omitempty"`
	LocationDetail    *CompanyLocationResponse   `json:"location_detail,omitempty"`
	FullAddress       string                     `json:"full_address,omitempty"`
	Description       string                     `json:"description,omitempty"`
	ShortDescription  string                     `json:"short_description,omitempty"` // From about field

	// Legacy Fields (for backward compatibility)
	Industry     string `json:"industry,omitempty"`
	SizeCategory string `json:"size_category,omitempty"`
	City         string `json:"city,omitempty"`
	Province     string `json:"province,omitempty"`
	Address      string `json:"address,omitempty"`

	// Other Fields
	CompanyType string `json:"company_type,omitempty"`
	WebsiteURL  string `json:"website_url,omitempty"`
	EmailDomain string `json:"email_domain,omitempty"`
	Phone       string `json:"phone,omitempty"`

	// Social Media URLs
	InstagramURL string `json:"instagram_url,omitempty"`
	FacebookURL  string `json:"facebook_url,omitempty"`
	LinkedinURL  string `json:"linkedin_url,omitempty"`
	TwitterURL   string `json:"twitter_url,omitempty"`

	Country    string     `json:"country"`
	PostalCode string     `json:"postal_code,omitempty"`
	Latitude   *float64   `json:"latitude,omitempty"`
	Longitude  *float64   `json:"longitude,omitempty"`
	LogoURL    string     `json:"logo_url,omitempty"`
	BannerURL  string     `json:"banner_url,omitempty"`
	Culture    string     `json:"culture,omitempty"`
	Benefits   []string   `json:"benefits,omitempty"`
	Verified   bool       `json:"verified"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`

	// Verification Status Details
	Status       string `json:"status"`                // "verified", "pending", "under_review", "rejected", "not_requested"
	BadgeGranted bool   `json:"badge_granted"`         // Whether company has verification badge
	NPWPNumber   string `json:"npwp_number,omitempty"` // NPWP number if verified
	NIBNumber    string `json:"nib_number,omitempty"`  // NIB number if provided

	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Profile    *CompanyProfileResponse   `json:"profile,omitempty"`
	Industries []CompanyIndustryResponse `json:"industries,omitempty"`
	Reviews    []CompanyReviewResponse   `json:"reviews,omitempty"`

	// Stats
	FollowersCount int64   `json:"followers_count"`
	JobsCount      int64   `json:"jobs_count"`
	AverageRating  float64 `json:"average_rating"`
	ReviewsCount   int64   `json:"reviews_count"`
	IsFollowing    bool    `json:"is_following,omitempty"` // For authenticated users
}

// CompanyProfileResponse represents company profile response
type CompanyProfileResponse struct {
	ID             int64     `json:"id"`
	FoundedYear    int16     `json:"founded_year,omitempty"`
	EmployeeCount  int32     `json:"employee_count,omitempty"`
	Description    string    `json:"description,omitempty"`
	Mission        string    `json:"mission,omitempty"`
	Vision         string    `json:"vision,omitempty"`
	CoreValues     string    `json:"core_values,omitempty"`
	FacebookURL    string    `json:"facebook_url,omitempty"`
	TwitterURL     string    `json:"twitter_url,omitempty"`
	LinkedinURL    string    `json:"linkedin_url,omitempty"`
	InstagramURL   string    `json:"instagram_url,omitempty"`
	YoutubeURL     string    `json:"youtube_url,omitempty"`
	Awards         string    `json:"awards,omitempty"`
	Certifications string    `json:"certifications,omitempty"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CompanyIndustryResponse represents company industry response
type CompanyIndustryResponse struct {
	ID           int64     `json:"id"`
	IndustryName string    `json:"industry_name"`
	CreatedAt    time.Time `json:"created_at"`
}

// CompanyReviewResponse represents company review response
type CompanyReviewResponse struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id,omitempty"`   // Hidden if anonymous
	UserName       string    `json:"user_name,omitempty"` // Hidden if anonymous
	Rating         int16     `json:"rating"`
	ReviewText     string    `json:"review_text"`
	ReviewTitle    string    `json:"review_title,omitempty"`
	Position       string    `json:"position"`
	EmploymentType string    `json:"employment_type"`
	WorkDuration   string    `json:"work_duration,omitempty"`
	IsAnonymous    bool      `json:"is_anonymous"`
	Pros           string    `json:"pros,omitempty"`
	Cons           string    `json:"cons,omitempty"`
	HelpfulCount   int32     `json:"helpful_count"`
	IsVerified     bool      `json:"is_verified"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CompanyFollowerResponse represents company follower response
type CompanyFollowerResponse struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	UserName   string    `json:"user_name"`
	FollowedAt time.Time `json:"followed_at"`
}

// CompanyEmployeeResponse represents company employee response
type CompanyEmployeeResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	UserName  string    `json:"user_name"`
	Position  string    `json:"position"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	JoinedAt  time.Time `json:"joined_at"`
	CreatedAt time.Time `json:"created_at"`
}

// CompanyVerificationResponse represents company verification response
type CompanyVerificationResponse struct {
	ID              int64      `json:"id"`
	DocumentType    string     `json:"document_type"`
	DocumentURL     string     `json:"document_url"`
	Status          string     `json:"status"`
	Notes           string     `json:"notes,omitempty"`
	VerifiedBy      *int64     `json:"verified_by,omitempty"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty"`
	RejectionReason string     `json:"rejection_reason,omitempty"`
	RequestedAt     time.Time  `json:"requested_at"`
}

// CompanyDocumentResponse represents company document response
type CompanyDocumentResponse struct {
	ID           int64     `json:"id"`
	DocumentType string    `json:"document_type"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	FileURL      string    `json:"file_url"`
	FileName     string    `json:"file_name"`
	FileSize     int64     `json:"file_size"`
	IsVerified   bool      `json:"is_verified"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

// CompanyListResponse represents list of companies response
type CompanyListResponse struct {
	Companies []CompanyResponse `json:"companies"`
}

// CompanyStatsResponse represents company statistics response
type CompanyStatsResponse struct {
	TotalJobs         int64   `json:"total_jobs"`
	ActiveJobs        int64   `json:"active_jobs"`
	TotalApplications int64   `json:"total_applications"`
	TotalFollowers    int64   `json:"total_followers"`
	AverageRating     float64 `json:"average_rating"`
	TotalReviews      int64   `json:"total_reviews"`
	TotalEmployees    int64   `json:"total_employees"`
}

// =============================================================================
// Master Data Responses
// =============================================================================

// MasterIndustryResponse represents industry master data in company response
type MasterIndustryResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"icon_url,omitempty"`
}

// MasterCompanySizeResponse represents company size master data in company response
type MasterCompanySizeResponse struct {
	ID           int64  `json:"id"`
	Label        string `json:"label"`
	Code         string `json:"code"`
	MinEmployees int    `json:"min_employees"`
	MaxEmployees *int   `json:"max_employees,omitempty"`
	Description  string `json:"description,omitempty"`
}

// CompanyLocationResponse represents complete location hierarchy in company response
type CompanyLocationResponse struct {
	Province     ProvinceResponse `json:"province"`
	City         CityResponse     `json:"city"`
	District     DistrictResponse `json:"district"`
	FullLocation string           `json:"full_location"` // e.g., "Batujajar, Kabupaten Bandung Barat, Jawa Barat"
}

// ProvinceResponse represents province in location response
type ProvinceResponse struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// CityResponse represents city in location response
type CityResponse struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Type        string `json:"type"`      // "Kota" or "Kabupaten"
	FullName    string `json:"full_name"` // e.g., "Kabupaten Bandung Barat"
	ProvinceID  int64  `json:"province_id"`
	PostalCodes string `json:"postal_codes,omitempty"`
}

// DistrictResponse represents district in location response
type DistrictResponse struct {
	ID     int64  `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	CityID int64  `json:"city_id"`
}

// CompanyAddressResponse represents company address for job posting
type CompanyAddressResponse struct {
	ID            int64   `json:"id"`
	AlamatLengkap string  `json:"alamat_lengkap"` // Full address text
	Latitude      float64 `json:"latitude,omitempty"`
	Longitude     float64 `json:"longitude,omitempty"`
}
