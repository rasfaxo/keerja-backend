package company

import (
	"context"
	"mime/multipart"
)

// CompanyService defines the business logic interface for company operations
type CompanyService interface {
	// Company registration and management
	RegisterCompany(ctx context.Context, req *RegisterCompanyRequest, userID int64) (*Company, error)
	GetCompany(ctx context.Context, id int64) (*Company, error)
	GetCompanyBySlug(ctx context.Context, slug string) (*Company, error)
	UpdateCompany(ctx context.Context, companyID int64, req *UpdateCompanyRequest) error
	DeleteCompany(ctx context.Context, companyID int64) error
	ListCompanies(ctx context.Context, filter *CompanyFilter) ([]Company, int64, error)
	SearchCompanies(ctx context.Context, query string, filter *CompanyFilter) ([]Company, int64, error)

	// Profile management
	CreateProfile(ctx context.Context, companyID int64, req *CreateProfileRequest) error
	UpdateProfile(ctx context.Context, companyID int64, req *UpdateProfileRequest) error
	GetProfile(ctx context.Context, companyID int64) (*CompanyProfile, error)
	PublishProfile(ctx context.Context, companyID int64) error
	UnpublishProfile(ctx context.Context, companyID int64) error

	// Logo and banner
	UploadLogo(ctx context.Context, companyID int64, file *multipart.FileHeader) (string, error)
	UploadBanner(ctx context.Context, companyID int64, file *multipart.FileHeader) (string, error)
	DeleteLogo(ctx context.Context, companyID int64) error
	DeleteBanner(ctx context.Context, companyID int64) error

	// Follower management
	FollowCompany(ctx context.Context, companyID, userID int64) error
	UnfollowCompany(ctx context.Context, companyID, userID int64) error
	IsFollowing(ctx context.Context, companyID, userID int64) (bool, error)
	GetFollowers(ctx context.Context, companyID int64, page, limit int) ([]CompanyFollower, int64, error)
	GetFollowedCompanies(ctx context.Context, userID int64, page, limit int) ([]Company, int64, error)
	GetFollowerCount(ctx context.Context, companyID int64) (int64, error)

	// Review management
	AddReview(ctx context.Context, req *AddReviewRequest) (*CompanyReview, error)
	UpdateReview(ctx context.Context, reviewID int64, userID int64, req *UpdateReviewRequest) error
	DeleteReview(ctx context.Context, reviewID, userID int64) error
	GetReview(ctx context.Context, reviewID int64) (*CompanyReview, error)
	GetCompanyReviews(ctx context.Context, companyID int64, filter *ReviewFilter) ([]CompanyReview, int64, error)
	GetUserReviews(ctx context.Context, userID int64) ([]CompanyReview, error)
	GetAverageRatings(ctx context.Context, companyID int64) (*AverageRatings, error)

	// Review moderation (admin only)
	ApproveReview(ctx context.Context, reviewID, moderatedBy int64) error
	RejectReview(ctx context.Context, reviewID, moderatedBy int64) error
	HideReview(ctx context.Context, reviewID, moderatedBy int64) error
	GetPendingReviews(ctx context.Context, page, limit int) ([]CompanyReview, int64, error)

	// Document management
	UploadDocument(ctx context.Context, companyID int64, file *multipart.FileHeader, req *UploadDocumentRequest) (*CompanyDocument, error)
	UpdateDocument(ctx context.Context, documentID int64, req *UpdateDocumentRequest) error
	DeleteDocument(ctx context.Context, documentID, companyID int64) error
	GetDocuments(ctx context.Context, companyID int64) ([]CompanyDocument, error)

	// Document verification (admin only)
	ApproveDocument(ctx context.Context, documentID, verifiedBy int64) error
	RejectDocument(ctx context.Context, documentID, verifiedBy int64, reason string) error
	CheckExpiredDocuments(ctx context.Context) error

	// Employee management
	AddEmployee(ctx context.Context, companyID int64, req *AddEmployeeRequest) (*CompanyEmployee, error)
	UpdateEmployee(ctx context.Context, employeeID, companyID int64, req *UpdateEmployeeRequest) error
	RemoveEmployee(ctx context.Context, employeeID, companyID int64) error
	GetEmployees(ctx context.Context, companyID int64, includeInactive bool) ([]CompanyEmployee, error)
	GetEmployeeCount(ctx context.Context, companyID int64) (int64, error)

	// Employer user management
	InviteEmployer(ctx context.Context, req *InviteEmployerRequest) error
	AcceptInvitation(ctx context.Context, userID, companyID int64) error
	UpdateEmployerRole(ctx context.Context, employerUserID int64, newRole string) error
	RemoveEmployerUser(ctx context.Context, employerUserID, companyID int64) error
	GetEmployerUsers(ctx context.Context, companyID int64) ([]EmployerUser, error)
	GetUserCompanies(ctx context.Context, userID int64) ([]Company, error)
	CheckEmployerPermission(ctx context.Context, userID, companyID int64, requiredRole string) (bool, error)

	// Verification management
	RequestVerification(ctx context.Context, companyID, requestedBy int64) error
	GetVerificationStatus(ctx context.Context, companyID int64) (*CompanyVerification, error)
	ApproveVerification(ctx context.Context, companyID, reviewedBy int64, notes string) error
	RejectVerification(ctx context.Context, companyID, reviewedBy int64, reason string) error
	GetPendingVerifications(ctx context.Context, page, limit int) ([]CompanyVerification, int64, error)
	RenewVerification(ctx context.Context, companyID int64) error
	CheckVerificationExpiry(ctx context.Context) error

	// Industry management (admin only)
	CreateIndustry(ctx context.Context, req *CreateIndustryRequest) (*CompanyIndustry, error)
	UpdateIndustry(ctx context.Context, industryID int64, req *UpdateIndustryRequest) error
	DeleteIndustry(ctx context.Context, industryID int64) error
	GetIndustry(ctx context.Context, industryID int64) (*CompanyIndustry, error)
	GetAllIndustries(ctx context.Context) ([]CompanyIndustry, error)
	GetIndustryTree(ctx context.Context) ([]CompanyIndustry, error)

	// Analytics and stats
	GetCompanyStats(ctx context.Context, companyID int64) (*CompanyStats, error)
	GetTopRatedCompanies(ctx context.Context, limit int) ([]Company, error)
	GetVerifiedCompanies(ctx context.Context, page, limit int) ([]Company, int64, error)
	GetCompanyEngagement(ctx context.Context, companyID int64) (*EngagementStats, error)
}

// Request DTOs

type RegisterCompanyRequest struct {
	CompanyName        string
	LegalName          *string
	RegistrationNumber *string
	Industry           *string
	CompanyType        *string
	SizeCategory       *string
	WebsiteURL         *string
	EmailDomain        *string
	Phone              *string
	Address            *string
	City               *string
	Province           *string
	Country            *string
	PostalCode         *string
	About              *string
}

type UpdateCompanyRequest struct {
	CompanyName        *string
	LegalName          *string
	RegistrationNumber *string
	Industry           *string
	CompanyType        *string
	SizeCategory       *string
	WebsiteURL         *string
	EmailDomain        *string
	Phone              *string
	Address            *string
	City               *string
	Province           *string
	PostalCode         *string
	Latitude           *float64
	Longitude          *float64
	About              *string
	Culture            *string
	Benefits           []string
}

type CreateProfileRequest struct {
	Tagline          *string
	ShortDescription *string
	LongDescription  *string
	Mission          *string
	Vision           *string
	Values           []string
	Culture          *string
	WorkEnvironment  *string
	VideoURL         *string
	Awards           []string
	SocialLinks      map[string]string
	HiringTagline    *string
	SEOTitle         *string
	SEOKeywords      []string
	SEODescription   *string
}

type UpdateProfileRequest struct {
	Tagline          *string
	ShortDescription *string
	LongDescription  *string
	Mission          *string
	Vision           *string
	Values           []string
	Culture          *string
	WorkEnvironment  *string
	VideoURL         *string
	Awards           []string
	SocialLinks      map[string]string
	HiringTagline    *string
	SEOTitle         *string
	SEOKeywords      []string
	SEODescription   *string
}

type AddReviewRequest struct {
	CompanyID          int64
	UserID             int64
	ReviewerType       *string
	PositionTitle      *string
	EmploymentPeriod   *string
	RatingOverall      float64
	RatingCulture      *float64
	RatingWorkLife     *float64
	RatingSalary       *float64
	RatingManagement   *float64
	Pros               *string
	Cons               *string
	AdviceToManagement *string
	IsAnonymous        bool
	RecommendToFriend  bool
}

type UpdateReviewRequest struct {
	ReviewerType       *string
	PositionTitle      *string
	EmploymentPeriod   *string
	RatingOverall      *float64
	RatingCulture      *float64
	RatingWorkLife     *float64
	RatingSalary       *float64
	RatingManagement   *float64
	Pros               *string
	Cons               *string
	AdviceToManagement *string
	RecommendToFriend  *bool
}

type UploadDocumentRequest struct {
	DocumentType   string
	DocumentNumber *string
	DocumentName   *string
	IssueDate      *string
	ExpiryDate     *string
}

type UpdateDocumentRequest struct {
	DocumentNumber *string
	DocumentName   *string
	IssueDate      *string
	ExpiryDate     *string
}

type AddEmployeeRequest struct {
	UserID           *int64
	FullName         *string
	JobTitle         *string
	Department       *string
	EmploymentType   string
	EmploymentStatus string
	JoinDate         *string
	SalaryRangeMin   *float64
	SalaryRangeMax   *float64
	Note             *string
	IsVisiblePublic  bool
}

type UpdateEmployeeRequest struct {
	FullName         *string
	JobTitle         *string
	Department       *string
	EmploymentType   *string
	EmploymentStatus *string
	JoinDate         *string
	EndDate          *string
	SalaryRangeMin   *float64
	SalaryRangeMax   *float64
	Note             *string
	IsVisiblePublic  *bool
}

type InviteEmployerRequest struct {
	CompanyID     int64
	Email         string
	Role          string
	PositionTitle *string
	Department    *string
}

type CreateIndustryRequest struct {
	Code        string
	Name        string
	Description *string
	ParentID    *int64
}

type UpdateIndustryRequest struct {
	Code        *string
	Name        *string
	Description *string
	ParentID    *int64
	IsActive    *bool
}

// Response DTOs

type CompanyStats struct {
	TotalJobs           int64
	ActiveJobs          int64
	TotalApplications   int64
	TotalFollowers      int64
	TotalEmployees      int64
	AverageRating       float64
	TotalReviews        int64
	VerificationStatus  string
	ProfileCompleteness int
}

type EngagementStats struct {
	TotalViews     int64
	TotalFollowers int64
	FollowerGrowth int64
	TotalReviews   int64
	AverageRating  float64
	ResponseRate   float64
}
