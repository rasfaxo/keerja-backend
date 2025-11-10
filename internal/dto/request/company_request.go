package request

// RegisterCompanyRequest represents company registration request
type RegisterCompanyRequest struct {
	CompanyName        string  `json:"company_name" validate:"required,min=2,max=200"`
	LegalName          *string `json:"legal_name" validate:"omitempty,max=200"`
	RegistrationNumber *string `json:"registration_number" validate:"omitempty,max=100"`

	// Master Data Relations (ID-based - backward compatibility)
	IndustryID    *int64 `json:"industry_id" validate:"omitempty"`
	CompanySizeID *int64 `json:"company_size_id" validate:"omitempty"`
	DistrictID    *int64 `json:"district_id" validate:"omitempty"`

	// Master Data Relations (Name-based - for mobile UI dropdown)
	IndustryName    *string `json:"industry_name" validate:"omitempty,max=100"`
	CompanySizeName *string `json:"company_size_name" validate:"omitempty,max=50"`
	ProvinceName    *string `json:"province_name" validate:"omitempty,max=100"`
	CityName        *string `json:"city_name" validate:"omitempty,max=100"`
	DistrictName    *string `json:"district_name" validate:"omitempty,max=100"`

	// Location Details
	FullAddress string  `json:"full_address" validate:"omitempty,max=500"`
	Description *string `json:"description" validate:"omitempty"`

	// Legacy Fields (kept for backward compatibility)
	Industry     *string `json:"industry" validate:"omitempty,max=100"`
	SizeCategory *string `json:"size_category" validate:"omitempty,oneof='1-10' '11-50' '51-200' '201-1000' '1000+'"`
	City         *string `json:"city" validate:"omitempty,max=100"`
	Province     *string `json:"province" validate:"omitempty,max=100"`
	Address      *string `json:"address" validate:"omitempty"`

	// Other Fields
	CompanyType *string `json:"company_type" validate:"omitempty,oneof=private public startup ngo government"`
	WebsiteURL  *string `json:"website_url" validate:"omitempty,url"`
	EmailDomain *string `json:"email_domain" validate:"omitempty,max=100"`
	Phone       *string `json:"phone" validate:"omitempty,max=30"`
	Country     *string `json:"country" validate:"omitempty,max=100"`
	PostalCode  *string `json:"postal_code" validate:"omitempty,max=10"`
	About       *string `json:"about" validate:"omitempty"`
}

// UpdateCompanyRequest represents company update request (Edit Profil Perusahaan)
type UpdateCompanyRequest struct {
	// NOTE: company_name, country, province, city, employee_count (company_size), industry
	// akan di-get dari data company yang sudah ada saat create company
	// Fields tersebut tidak perlu di-update via endpoint ini

	// Full Address (dari data company saat create, bisa di-edit)
	FullAddress *string `form:"full_address" json:"full_address" validate:"omitempty,max=500"`

	// Deskripsi Singkat - Visi dan Misi Perusahaan (required)
	ShortDescription *string `form:"short_description" json:"short_description" validate:"required,max=1000"`

	// Website & Social Media
	WebsiteURL   *string `form:"website_url" json:"website_url" validate:"omitempty,url"`
	InstagramURL *string `form:"instagram_url" json:"instagram_url" validate:"omitempty"`
	FacebookURL  *string `form:"facebook_url" json:"facebook_url" validate:"omitempty"`
	LinkedinURL  *string `form:"linkedin_url" json:"linkedin_url" validate:"omitempty"`
	TwitterURL   *string `form:"twitter_url" json:"twitter_url" validate:"omitempty"`

	// Rich Text Descriptions
	CompanyDescription *string `form:"company_description" json:"company_description" validate:"required"` // Deskripsi Perusahaan (required)
	CompanyCulture     *string `form:"company_culture" json:"company_culture" validate:"omitempty"`        // Budaya Perusahaan (optional)
} // UpdateCompanyProfileRequest represents company profile update request
type UpdateCompanyProfileRequest struct {
	FoundedYear    *int16  `json:"founded_year" validate:"omitempty,min=1800,max=2100"`
	EmployeeCount  *int32  `json:"employee_count" validate:"omitempty,min=0"`
	Description    *string `json:"description" validate:"omitempty"`
	Mission        *string `json:"mission" validate:"omitempty"`
	Vision         *string `json:"vision" validate:"omitempty"`
	CoreValues     *string `json:"core_values" validate:"omitempty"`
	FacebookURL    *string `json:"facebook_url" validate:"omitempty,url"`
	TwitterURL     *string `json:"twitter_url" validate:"omitempty,url"`
	LinkedinURL    *string `json:"linkedin_url" validate:"omitempty,url"`
	InstagramURL   *string `json:"instagram_url" validate:"omitempty,url"`
	YoutubeURL     *string `json:"youtube_url" validate:"omitempty,url"`
	Awards         *string `json:"awards" validate:"omitempty"`
	Certifications *string `json:"certifications" validate:"omitempty"`
}

// AddCompanyIndustryRequest represents add industry request
type AddCompanyIndustryRequest struct {
	IndustryName string `json:"industry_name" validate:"required,max=100"`
}

// AddReviewRequest represents add company review request
type AddReviewRequest struct {
	RatingOverall      float64  `json:"rating_overall" validate:"required,min=1,max=5"`
	RatingCulture      *float64 `json:"rating_culture" validate:"omitempty,min=1,max=5"`
	RatingWorkLife     *float64 `json:"rating_work_life" validate:"omitempty,min=1,max=5"`
	RatingSalary       *float64 `json:"rating_salary" validate:"omitempty,min=1,max=5"`
	RatingManagement   *float64 `json:"rating_management" validate:"omitempty,min=1,max=5"`
	ReviewerType       *string  `json:"reviewer_type" validate:"omitempty,max=50"`
	PositionTitle      *string  `json:"position_title" validate:"omitempty,max=150"`
	EmploymentPeriod   *string  `json:"employment_period" validate:"omitempty,max=50"`
	Pros               *string  `json:"pros" validate:"omitempty"`
	Cons               *string  `json:"cons" validate:"omitempty"`
	AdviceToManagement *string  `json:"advice_to_management" validate:"omitempty"`
	IsAnonymous        bool     `json:"is_anonymous"`
	RecommendToFriend  bool     `json:"recommend_to_friend"`
}

// UpdateReviewRequest represents update review request
type UpdateReviewRequest struct {
	ReviewerType       *string  `json:"reviewer_type" validate:"omitempty,max=50"`
	PositionTitle      *string  `json:"position_title" validate:"omitempty,max=150"`
	EmploymentPeriod   *string  `json:"employment_period" validate:"omitempty,max=50"`
	RatingOverall      *float64 `json:"rating_overall" validate:"omitempty,min=1,max=5"`
	RatingCulture      *float64 `json:"rating_culture" validate:"omitempty,min=1,max=5"`
	RatingWorkLife     *float64 `json:"rating_work_life" validate:"omitempty,min=1,max=5"`
	RatingSalary       *float64 `json:"rating_salary" validate:"omitempty,min=1,max=5"`
	RatingManagement   *float64 `json:"rating_management" validate:"omitempty,min=1,max=5"`
	Pros               *string  `json:"pros" validate:"omitempty"`
	Cons               *string  `json:"cons" validate:"omitempty"`
	AdviceToManagement *string  `json:"advice_to_management" validate:"omitempty"`
	RecommendToFriend  *bool    `json:"recommend_to_friend"`
}

// InviteEmployeeRequest represents employee invitation request
type InviteEmployeeRequest struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required,min=3,max=150"`
	Position string `json:"position" validate:"required,max=150"`
	Role     string `json:"role" validate:"required,oneof=admin hr recruiter manager viewer"`
}

// UpdateEmployeeRequest represents update employee request
type UpdateEmployeeRequest struct {
	Position *string `json:"position" validate:"omitempty,max=150"`
	Role     *string `json:"role" validate:"omitempty,oneof=admin hr recruiter manager viewer"`
	IsActive *bool   `json:"is_active"`
}

// RequestVerificationRequest represents company verification request
type RequestVerificationRequest struct {
	NPWPNumber *string `form:"npwp_number" json:"npwp_number" validate:"required"`
	NIBNumber  *string `form:"nib_number" json:"nib_number" validate:"omitempty"`
	Notes      *string `form:"notes" json:"notes" validate:"omitempty"`
}

// UploadCompanyDocumentRequest represents company document upload request
type UploadCompanyDocumentRequest struct {
	DocumentType string `form:"document_type" validate:"required,oneof='business_license' 'tax_id' 'certificate' 'other'"`
	Title        string `form:"title" validate:"required,max=200"`
	Description  string `form:"description" validate:"omitempty"`
}

// CompanySearchRequest represents company search request
type CompanySearchRequest struct {
	Query        string `json:"query" query:"q" validate:"omitempty"`
	Industry     string `json:"industry" query:"industry" validate:"omitempty"`
	CompanyType  string `json:"company_type" query:"company_type" validate:"omitempty,oneof=private public startup ngo government"`
	SizeCategory string `json:"size_category" query:"size_category" validate:"omitempty"`
	Location     string `json:"location" query:"location" validate:"omitempty"`
	IsVerified   *bool  `json:"is_verified" query:"is_verified" validate:"omitempty"`
	Page         int    `json:"page" query:"page" validate:"omitempty,min=1"`
	Limit        int    `json:"limit" query:"limit" validate:"omitempty,min=1,max=100"`
	SortBy       string `json:"sort_by" query:"sort_by" validate:"omitempty,oneof=name created_at followers rating"`
	SortOrder    string `json:"sort_order" query:"sort_order" validate:"omitempty,oneof=asc desc"`
}
