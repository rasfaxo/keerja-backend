package job

import (
	"context"
	"time"
)

// JobRepository defines the interface for job data access
type JobRepository interface {
	UpdateStatusByCompany(ctx context.Context, companyID int64, fromStatus, toStatus string) error
	// Job CRUD operations
	Create(ctx context.Context, job *Job) error
	FindByID(ctx context.Context, id int64) (*Job, error)
	FindByUUID(ctx context.Context, uuid string) (*Job, error)
	FindBySlug(ctx context.Context, slug string) (*Job, error)
	Update(ctx context.Context, job *Job) error
	Delete(ctx context.Context, id int64) error
	SoftDelete(ctx context.Context, id int64) error

	// Job listing and search
	List(ctx context.Context, filter JobFilter, page, limit int) ([]Job, int64, error)
	ListByCompany(ctx context.Context, companyID int64, filter JobFilter, page, limit int) ([]Job, int64, error)
	ListByEmployer(ctx context.Context, employerUserID int64, filter JobFilter, page, limit int) ([]Job, int64, error)
	SearchJobs(ctx context.Context, filter JobSearchFilter, page, limit int) ([]Job, int64, error)

	// Job status operations
	UpdateStatus(ctx context.Context, id int64, status string) error
	UpdateStatusWithExpiry(ctx context.Context, id int64, status string, publishedAt *time.Time, expiredAt *time.Time) error
	PublishJob(ctx context.Context, id int64) error
	CloseJob(ctx context.Context, id int64) error
	ExpireJob(ctx context.Context, id int64) error
	SuspendJob(ctx context.Context, id int64) error
	GetExpiredJobs(ctx context.Context) ([]Job, error)
	GetExpiringJobs(ctx context.Context, days int) ([]Job, error)

	// Job statistics
	IncrementViews(ctx context.Context, id int64) error
	IncrementApplications(ctx context.Context, id int64) error
	GetJobStats(ctx context.Context, jobID int64) (*JobStats, error)
	GetCompanyJobStats(ctx context.Context, companyID int64) (*CompanyJobStats, error)

	// Recommendation and matching
	GetRecommendedJobs(ctx context.Context, userID int64, limit int) ([]Job, error)
	GetSimilarJobs(ctx context.Context, jobID int64, limit int) ([]Job, error)
	GetMatchingJobs(ctx context.Context, userID int64, filter JobFilter, page, limit int) ([]Job, int64, error)

	// Advanced search
	SearchByLocation(ctx context.Context, latitude, longitude, radius float64, filter JobFilter, page, limit int) ([]Job, int64, error)
	SearchBySkills(ctx context.Context, skillIDs []int64, filter JobFilter, page, limit int) ([]Job, int64, error)
	SearchBySalaryRange(ctx context.Context, minSalary, maxSalary float64, filter JobFilter, page, limit int) ([]Job, int64, error)

	// JobCategory CRUD
	CreateCategory(ctx context.Context, category *JobCategory) error
	FindCategoryByID(ctx context.Context, id int64) (*JobCategory, error)
	FindCategoryByCode(ctx context.Context, code string) (*JobCategory, error)
	UpdateCategory(ctx context.Context, category *JobCategory) error
	DeleteCategory(ctx context.Context, id int64) error
	ListCategories(ctx context.Context, filter CategoryFilter, page, limit int) ([]JobCategory, int64, error)
	GetCategoryTree(ctx context.Context) ([]JobCategory, error)
	GetActiveCategories(ctx context.Context) ([]JobCategory, error)

	// Get jobs by specific status
	GetJobsByStatus(ctx context.Context, userID int64, status string, page, limit int) ([]Job, int64, error)

	// JobSubcategory CRUD
	CreateSubcategory(ctx context.Context, subcategory *JobSubcategory) error
	FindSubcategoryByID(ctx context.Context, id int64) (*JobSubcategory, error)
	FindSubcategoryByCode(ctx context.Context, code string) (*JobSubcategory, error)
	UpdateSubcategory(ctx context.Context, subcategory *JobSubcategory) error
	DeleteSubcategory(ctx context.Context, id int64) error
	ListSubcategories(ctx context.Context, categoryID int64) ([]JobSubcategory, error)
	GetActiveSubcategories(ctx context.Context, categoryID int64) ([]JobSubcategory, error)

	// JobLocation operations
	CreateLocation(ctx context.Context, location *JobLocation) error
	FindLocationByID(ctx context.Context, id int64) (*JobLocation, error)
	UpdateLocation(ctx context.Context, location *JobLocation) error
	DeleteLocation(ctx context.Context, id int64) error
	ListLocationsByJob(ctx context.Context, jobID int64) ([]JobLocation, error)
	GetPrimaryLocation(ctx context.Context, jobID int64) (*JobLocation, error)
	SetPrimaryLocation(ctx context.Context, jobID, locationID int64) error

	// JobBenefit operations
	CreateBenefit(ctx context.Context, benefit *JobBenefit) error
	FindBenefitByID(ctx context.Context, id int64) (*JobBenefit, error)
	UpdateBenefit(ctx context.Context, benefit *JobBenefit) error
	DeleteBenefit(ctx context.Context, id int64) error
	ListBenefitsByJob(ctx context.Context, jobID int64) ([]JobBenefit, error)
	GetHighlightedBenefits(ctx context.Context, jobID int64) ([]JobBenefit, error)
	BulkCreateBenefits(ctx context.Context, benefits []JobBenefit) error
	BulkDeleteBenefits(ctx context.Context, jobID int64) error

	// JobSkill operations
	CreateSkill(ctx context.Context, skill *JobSkill) error
	FindSkillByID(ctx context.Context, id int64) (*JobSkill, error)
	UpdateSkill(ctx context.Context, skill *JobSkill) error
	DeleteSkill(ctx context.Context, id int64) error
	ListSkillsByJob(ctx context.Context, jobID int64) ([]JobSkill, error)
	GetRequiredSkills(ctx context.Context, jobID int64) ([]JobSkill, error)
	GetPreferredSkills(ctx context.Context, jobID int64) ([]JobSkill, error)
	BulkCreateSkills(ctx context.Context, skills []JobSkill) error
	BulkDeleteSkills(ctx context.Context, jobID int64) error

	// JobRequirement operations
	CreateRequirement(ctx context.Context, requirement *JobRequirement) error
	FindRequirementByID(ctx context.Context, id int64) (*JobRequirement, error)
	UpdateRequirement(ctx context.Context, requirement *JobRequirement) error
	DeleteRequirement(ctx context.Context, id int64) error
	ListRequirementsByJob(ctx context.Context, jobID int64) ([]JobRequirement, error)
	GetMandatoryRequirements(ctx context.Context, jobID int64) ([]JobRequirement, error)
	BulkCreateRequirements(ctx context.Context, requirements []JobRequirement) error
	BulkDeleteRequirements(ctx context.Context, jobID int64) error

	// Analytics
	GetTrendingJobs(ctx context.Context, limit int) ([]Job, error)
	GetPopularCategories(ctx context.Context, limit int) ([]CategoryStats, error)
	GetJobsByDateRange(ctx context.Context, startDate, endDate time.Time, filter JobFilter) ([]Job, error)

	// Master Data Preload
	PreloadMasterData(ctx context.Context, job *Job) error
	FindByIDWithMasterData(ctx context.Context, id int64) (*Job, error)
}

// JobFilter defines filter criteria for job listing
type JobFilter struct {
	Status         string
	CompanyID      int64
	CategoryID     int64
	SubcategoryID  int64
	City           string
	Province       string
	JobLevel       string
	EmploymentType string
	RemoteOption   *bool
	MinSalary      *float64
	MaxSalary      *float64
	MinExperience  *int16
	MaxExperience  *int16
	EducationLevel string
	IsActive       *bool
	PublishedAfter *time.Time
	SortBy         string // "latest", "salary_asc", "salary_desc", "views", "applications"
}

// JobSearchFilter defines advanced search criteria
type JobSearchFilter struct {
	Keyword         string
	Location        string
	CategoryIDs     []int64
	SkillIDs        []int64
	JobLevels       []string
	EmploymentTypes []string
	RemoteOnly      bool
	MinSalary       *float64
	MaxSalary       *float64
	MinExperience   *int16
	MaxExperience   *int16
	EducationLevels []string
	CompanyIDs      []int64
	PostedWithin    *int // days
}

// CategoryFilter defines filter criteria for category listing
type CategoryFilter struct {
	ParentID *int64
	IsActive *bool
	Keyword  string
}

// JobStats represents job statistics
type JobStats struct {
	JobID                 int64
	ViewsCount            int64
	ApplicationsCount     int64
	ViewsToday            int64
	ApplicationsToday     int64
	ViewsThisWeek         int64
	ApplicationsThisWeek  int64
	ViewsThisMonth        int64
	ApplicationsThisMonth int64
	ConversionRate        float64
}

// CompanyJobStats represents company's job statistics
type CompanyJobStats struct {
	CompanyID                 int64
	TotalJobs                 int64
	ActiveJobs                int64
	DraftJobs                 int64
	ClosedJobs                int64
	ExpiredJobs               int64
	TotalViews                int64
	TotalApplications         int64
	AverageViewsPerJob        float64
	AverageApplicationsPerJob float64
}

// CategoryStats represents category statistics
type CategoryStats struct {
	CategoryID        int64
	CategoryName      string
	JobCount          int64
	ActiveJobCount    int64
	TotalViews        int64
	TotalApplications int64
}
