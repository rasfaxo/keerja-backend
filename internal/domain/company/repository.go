package company

import (
	"context"
)

// CompanyRepository defines the interface for company data access
type CompanyRepository interface {
	// Company CRUD
	Create(ctx context.Context, company *Company) error
	FindByID(ctx context.Context, id int64) (*Company, error)
	FindByUUID(ctx context.Context, uuid string) (*Company, error)
	FindBySlug(ctx context.Context, slug string) (*Company, error)
	Update(ctx context.Context, company *Company) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter *CompanyFilter) ([]Company, int64, error)

	// Company CRUD with Master Data Preloading
	FindByIDWithMasterData(ctx context.Context, id int64) (*Company, error)
	FindByUUIDWithMasterData(ctx context.Context, uuid string) (*Company, error)
	FindBySlugWithMasterData(ctx context.Context, slug string) (*Company, error)
	ListWithMasterData(ctx context.Context, filter *CompanyFilter) ([]Company, int64, error)

	// Profile operations
	CreateProfile(ctx context.Context, profile *CompanyProfile) error
	FindProfileByCompanyID(ctx context.Context, companyID int64) (*CompanyProfile, error)
	UpdateProfile(ctx context.Context, profile *CompanyProfile) error

	// Follower operations
	FollowCompany(ctx context.Context, companyID, userID int64) error
	UnfollowCompany(ctx context.Context, companyID, userID int64) error
	IsFollowing(ctx context.Context, companyID, userID int64) (bool, error)
	GetFollowers(ctx context.Context, companyID int64, page, limit int) ([]CompanyFollower, int64, error)
	GetFollowedCompanies(ctx context.Context, userID int64, page, limit int) ([]Company, int64, error)
	CountFollowers(ctx context.Context, companyID int64) (int64, error)

	// Review operations
	CreateReview(ctx context.Context, review *CompanyReview) error
	UpdateReview(ctx context.Context, review *CompanyReview) error
	DeleteReview(ctx context.Context, id int64) error
	FindReviewByID(ctx context.Context, id int64) (*CompanyReview, error)
	GetReviewsByCompanyID(ctx context.Context, companyID int64, filter *ReviewFilter) ([]CompanyReview, int64, error)
	GetReviewsByUserID(ctx context.Context, userID int64) ([]CompanyReview, error)
	ApproveReview(ctx context.Context, id, moderatedBy int64) error
	RejectReview(ctx context.Context, id, moderatedBy int64) error
	CalculateAverageRatings(ctx context.Context, companyID int64) (*AverageRatings, error)

	// Document operations
	CreateDocument(ctx context.Context, doc *CompanyDocument) error
	UpdateDocument(ctx context.Context, doc *CompanyDocument) error
	DeleteDocument(ctx context.Context, id int64) error
	FindDocumentByID(ctx context.Context, id int64) (*CompanyDocument, error)
	GetDocumentsByCompanyID(ctx context.Context, companyID int64) ([]CompanyDocument, error)
	ApproveDocument(ctx context.Context, id, verifiedBy int64) error
	RejectDocument(ctx context.Context, id, verifiedBy int64, reason string) error

	// Employee operations
	AddEmployee(ctx context.Context, employee *CompanyEmployee) error
	UpdateEmployee(ctx context.Context, employee *CompanyEmployee) error
	DeleteEmployee(ctx context.Context, id int64) error
	GetEmployeesByCompanyID(ctx context.Context, companyID int64, includeInactive bool) ([]CompanyEmployee, error)
	CountEmployees(ctx context.Context, companyID int64, activeOnly bool) (int64, error)

	// Employer User operations
	CreateEmployerUser(ctx context.Context, employerUser *EmployerUser) error
	UpdateEmployerUser(ctx context.Context, employerUser *EmployerUser) error
	DeleteEmployerUser(ctx context.Context, id int64) error
	FindEmployerUserByID(ctx context.Context, id int64) (*EmployerUser, error)
	FindEmployerUserByUserAndCompany(ctx context.Context, userID, companyID int64) (*EmployerUser, error)
	GetEmployerUsersByCompanyID(ctx context.Context, companyID int64) ([]EmployerUser, error)
	GetCompaniesByUserID(ctx context.Context, userID int64) ([]Company, error)

	// Company Invitation operations
	CreateInvitation(ctx context.Context, invitation *CompanyInvitation) error
	FindInvitationByToken(ctx context.Context, token string) (*CompanyInvitation, error)
	FindInvitationByID(ctx context.Context, id int64) (*CompanyInvitation, error)
	UpdateInvitation(ctx context.Context, invitation *CompanyInvitation) error
	GetPendingInvitationsByCompany(ctx context.Context, companyID int64) ([]CompanyInvitation, error)
	GetPendingInvitationsByEmail(ctx context.Context, email string) ([]CompanyInvitation, error)
	ExpireOldInvitations(ctx context.Context) error
	DeleteInvitation(ctx context.Context, id int64) error

	// Verification operations
	CreateVerification(ctx context.Context, verification *CompanyVerification) error
	UpdateVerification(ctx context.Context, verification *CompanyVerification) error
	FindVerificationByCompanyID(ctx context.Context, companyID int64) (*CompanyVerification, error)
	RequestVerification(ctx context.Context, companyID, requestedBy int64) error
	ApproveVerification(ctx context.Context, companyID, reviewedBy int64, notes string) error
	RejectVerification(ctx context.Context, companyID, reviewedBy int64, reason string) error
	GetPendingVerifications(ctx context.Context, page, limit int) ([]CompanyVerification, int64, error)

	// Industry operations
	CreateIndustry(ctx context.Context, industry *CompanyIndustry) error
	UpdateIndustry(ctx context.Context, industry *CompanyIndustry) error
	DeleteIndustry(ctx context.Context, id int64) error
	FindIndustryByID(ctx context.Context, id int64) (*CompanyIndustry, error)
	FindIndustryByCode(ctx context.Context, code string) (*CompanyIndustry, error)
	GetAllIndustries(ctx context.Context, activeOnly bool) ([]CompanyIndustry, error)
	GetIndustryTree(ctx context.Context) ([]CompanyIndustry, error)

	// Search and analytics
	SearchCompanies(ctx context.Context, query string, filter *CompanyFilter) ([]Company, int64, error)
	GetVerifiedCompanies(ctx context.Context, page, limit int) ([]Company, int64, error)
	GetTopRatedCompanies(ctx context.Context, limit int) ([]Company, error)
	GetCompaniesNeedingVerificationRenewal(ctx context.Context) ([]Company, error)

	// Full data with relationships
	GetFullCompanyProfile(ctx context.Context, companyID int64) (*Company, error)
}

// CompanyFilter represents filters for querying companies
type CompanyFilter struct {
	// Master Data Filters
	IndustryID    *int64
	CompanySizeID *int64
	ProvinceID    *int64
	CityID        *int64
	DistrictID    *int64

	// Legacy Filters
	Industry     *string
	CompanyType  *string
	SizeCategory *string
	City         *string
	Province     *string
	Verified     *bool
	IsActive     *bool
	SearchQuery  *string
	Page         int
	Limit        int
	SortBy       string
	SortOrder    string
}

// ReviewFilter represents filters for querying reviews
type ReviewFilter struct {
	ReviewerType *string
	Status       *string
	MinRating    *float64
	MaxRating    *float64
	Page         int
	Limit        int
	SortBy       string
	SortOrder    string
}

// AverageRatings represents average ratings for a company
type AverageRatings struct {
	Overall      float64
	Culture      float64
	WorkLife     float64
	Salary       float64
	Management   float64
	TotalReviews int64
}
