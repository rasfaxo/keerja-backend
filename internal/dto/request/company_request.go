package request

// RegisterCompanyRequest represents company registration request
type RegisterCompanyRequest struct {
	CompanyName        string `json:"company_name" validate:"required,min=2,max=200"`
	LegalName          string `json:"legal_name" validate:"omitempty,max=200"`
	RegistrationNumber string `json:"registration_number" validate:"omitempty,max=100"`
	Industry           string `json:"industry" validate:"omitempty,max=100"`
	CompanyType        string `json:"company_type" validate:"omitempty,oneof=private public startup ngo government"`
	SizeCategory       string `json:"size_category" validate:"omitempty,oneof='1-10' '11-50' '51-200' '201-1000' '1000+'"`
	WebsiteURL         string `json:"website_url" validate:"omitempty,url"`
	EmailDomain        string `json:"email_domain" validate:"omitempty,max=100"`
	Phone              string `json:"phone" validate:"omitempty,max=30"`
	Address            string `json:"address" validate:"omitempty"`
	City               string `json:"city" validate:"omitempty,max=100"`
	Province           string `json:"province" validate:"omitempty,max=100"`
	Country            string `json:"country" validate:"omitempty,max=100"`
	PostalCode         string `json:"postal_code" validate:"omitempty,max=10"`
	About              string `json:"about" validate:"omitempty"`
}

// UpdateCompanyRequest represents company update request
type UpdateCompanyRequest struct {
	CompanyName        *string  `json:"company_name" validate:"omitempty,min=2,max=200"`
	LegalName          *string  `json:"legal_name" validate:"omitempty,max=200"`
	RegistrationNumber *string  `json:"registration_number" validate:"omitempty,max=100"`
	Industry           *string  `json:"industry" validate:"omitempty,max=100"`
	CompanyType        *string  `json:"company_type" validate:"omitempty,oneof=private public startup ngo government"`
	SizeCategory       *string  `json:"size_category" validate:"omitempty,oneof='1-10' '11-50' '51-200' '201-1000' '1000+'"`
	WebsiteURL         *string  `json:"website_url" validate:"omitempty,url"`
	EmailDomain        *string  `json:"email_domain" validate:"omitempty,max=100"`
	Phone              *string  `json:"phone" validate:"omitempty,max=30"`
	Address            *string  `json:"address" validate:"omitempty"`
	City               *string  `json:"city" validate:"omitempty,max=100"`
	Province           *string  `json:"province" validate:"omitempty,max=100"`
	Country            *string  `json:"country" validate:"omitempty,max=100"`
	PostalCode         *string  `json:"postal_code" validate:"omitempty,max=10"`
	Latitude           *float64 `json:"latitude" validate:"omitempty"`
	Longitude          *float64 `json:"longitude" validate:"omitempty"`
	About              *string  `json:"about" validate:"omitempty"`
	Culture            *string  `json:"culture" validate:"omitempty"`
	Benefits           []string `json:"benefits" validate:"omitempty"`
}

// UpdateCompanyProfileRequest represents company profile update request
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
	Rating         int16   `json:"rating" validate:"required,min=1,max=5"`
	ReviewText     string  `json:"review_text" validate:"required,min=10"`
	ReviewTitle    *string `json:"review_title" validate:"omitempty,max=200"`
	Position       string  `json:"position" validate:"required,max=150"`
	EmploymentType string  `json:"employment_type" validate:"required,oneof='Full-Time' 'Part-Time' 'Contract' 'Internship' 'Freelance'"`
	WorkDuration   *string `json:"work_duration" validate:"omitempty,max=50"`
	IsAnonymous    bool    `json:"is_anonymous"`
	Pros           *string `json:"pros" validate:"omitempty"`
	Cons           *string `json:"cons" validate:"omitempty"`
}

// UpdateReviewRequest represents update review request
type UpdateReviewRequest struct {
	Rating      *int16  `json:"rating" validate:"omitempty,min=1,max=5"`
	ReviewText  *string `json:"review_text" validate:"omitempty,min=10"`
	ReviewTitle *string `json:"review_title" validate:"omitempty,max=200"`
	Pros        *string `json:"pros" validate:"omitempty"`
	Cons        *string `json:"cons" validate:"omitempty"`
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
	DocumentType string `json:"document_type" validate:"required,oneof='business_license' 'tax_id' 'incorporation_certificate' 'other'"`
	Notes        string `json:"notes" validate:"omitempty"`
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
