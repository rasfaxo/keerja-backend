package request

// AdminGetCompaniesRequest represents the request for getting companies list (admin)
type AdminGetCompaniesRequest struct {
	// Pagination
	Page  int `query:"page" validate:"omitempty,min=1"`
	Limit int `query:"limit" validate:"omitempty,min=1,max=100"`

	// Search
	Search string `query:"search"` // Search by company name, email, legal name

	// Filters
	Status            string `query:"status"`             // pending_verification, verified, rejected, suspended
	VerificationStatus string `query:"verification_status"` // For company_verifications.status
	IndustryID        *int64 `query:"industry_id"`
	CompanySizeID     *int64 `query:"company_size_id"`
	ProvinceID        *int64 `query:"province_id"`
	CityID            *int64 `query:"city_id"`
	Verified          *bool  `query:"verified"` // true/false
	IsActive          *bool  `query:"is_active"` // true/false

	// Date range
	CreatedFrom string `query:"created_from"` // Format: 2024-01-01
	CreatedTo   string `query:"created_to"`   // Format: 2024-12-31

	// Sorting
	SortBy    string `query:"sort_by" validate:"omitempty,oneof=company_name created_at verified_at updated_at"` // Default: created_at
	SortOrder string `query:"sort_order" validate:"omitempty,oneof=asc desc"`                                   // Default: desc
}

// AdminUpdateCompanyStatusRequest represents the request for updating company verification status
type AdminUpdateCompanyStatusRequest struct {
	Status          string  `json:"status" validate:"required,oneof=pending_verification verified rejected suspended blacklisted"`
	RejectionReason *string `json:"rejection_reason" validate:"omitempty,max=1000"`
	Notes           *string `json:"notes" validate:"omitempty,max=2000"`
	GrantBadge      *bool   `json:"grant_badge"` // Grant verification badge
}

// AdminUpdateCompanyRequest represents the request for admin to edit company details
type AdminUpdateCompanyRequest struct {
	CompanyName        *string  `json:"company_name" validate:"omitempty,min=2,max=200"`
	LegalName          *string  `json:"legal_name" validate:"omitempty,max=200"`
	RegistrationNumber *string  `json:"registration_number" validate:"omitempty,max=100"`
	IndustryID         *int64   `json:"industry_id"`
	CompanySizeID      *int64   `json:"company_size_id"`
	DistrictID         *int64   `json:"district_id"`
	FullAddress        *string  `json:"full_address" validate:"omitempty,max=500"`
	Description        *string  `json:"description" validate:"omitempty,max=2000"`
	WebsiteURL         *string  `json:"website_url" validate:"omitempty,url"`
	EmailDomain        *string  `json:"email_domain" validate:"omitempty,max=100"`
	Phone              *string  `json:"phone" validate:"omitempty,min=10,max=30"`
	PostalCode         *string  `json:"postal_code" validate:"omitempty,max=10"`
	Latitude           *float64 `json:"latitude" validate:"omitempty,min=-90,max=90"`
	Longitude          *float64 `json:"longitude" validate:"omitempty,min=-180,max=180"`
	About              *string  `json:"about" validate:"omitempty,max=5000"`
	Culture            *string  `json:"culture" validate:"omitempty,max=5000"`
	Benefits           []string `json:"benefits"`
	IsActive           *bool    `json:"is_active"`
	Verified           *bool    `json:"verified"`
}

// AdminDeleteCompanyRequest represents validation for company deletion
type AdminDeleteCompanyRequest struct {
	Force  bool   `query:"force"` // Force delete even with active jobs
	Reason string `json:"reason" validate:"required,min=10,max=500"` // Reason for deletion
}
