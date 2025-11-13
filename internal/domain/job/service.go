package job

import (
	"context"
	"time"
)

// JobService defines the interface for job business logic
type JobService interface {
	// Job management (Employer)
	CreateJob(ctx context.Context, req *CreateJobRequest) (*Job, error)
	UpdateJob(ctx context.Context, jobID int64, req *UpdateJobRequest) (*Job, error)
	DeleteJob(ctx context.Context, jobID int64, employerUserID int64) error
	GetJob(ctx context.Context, jobID int64) (*Job, error)
	GetJobBySlug(ctx context.Context, slug string) (*Job, error)
	GetJobByUUID(ctx context.Context, uuid string) (*Job, error)
	GetMyJobs(ctx context.Context, employerUserID int64, filter JobFilter, page, limit int) ([]Job, int64, error)
	GetCompanyJobs(ctx context.Context, companyID int64, filter JobFilter, page, limit int) ([]Job, int64, error)

	// Phase 6: Job draft workflow
	SaveJobDraft(ctx context.Context, companyID int64, req *SaveJobDraftRequest) (*Job, error)

	// Job status management
	PublishJob(ctx context.Context, jobID int64, employerUserID int64) error
	UnpublishJob(ctx context.Context, jobID int64, employerUserID int64) error
	CloseJob(ctx context.Context, jobID int64, employerUserID int64) error
	ReopenJob(ctx context.Context, jobID int64, employerUserID int64) error
	SuspendJob(ctx context.Context, jobID int64, employerUserID int64, reason string) error
	SetJobExpiry(ctx context.Context, jobID int64, expiryDate time.Time) error
	ExtendJobExpiry(ctx context.Context, jobID int64, days int) error
	AutoExpireJobs(ctx context.Context) error
	UpdateStatus(ctx context.Context, jobID int64, status string) error

	// Job search and discovery (Public)
	ListJobs(ctx context.Context, filter JobFilter, page, limit int) ([]Job, int64, error)
	SearchJobs(ctx context.Context, filter JobSearchFilter, page, limit int) (*JobSearchResponse, error)
	SearchJobsByLocation(ctx context.Context, latitude, longitude, radius float64, filter JobFilter, page, limit int) ([]Job, int64, error)
	GetFeaturedJobs(ctx context.Context, limit int) ([]Job, error)
	GetLatestJobs(ctx context.Context, limit int) ([]Job, error)
	GetTrendingJobs(ctx context.Context, limit int) ([]Job, error)
	GetRecommendedJobs(ctx context.Context, userID int64, limit int) ([]Job, error)
	GetSimilarJobs(ctx context.Context, jobID int64, limit int) ([]Job, error)

	// Job matching
	CalculateMatchScore(ctx context.Context, jobID, userID int64) (*MatchScore, error)
	GetMatchingJobs(ctx context.Context, userID int64, filter JobFilter, page, limit int) (*MatchResponse, error)

	// Job views and interactions
	IncrementView(ctx context.Context, jobID int64, userID *int64) error
	GetJobStats(ctx context.Context, jobID int64) (*JobStats, error)
	GetCompanyJobStats(ctx context.Context, companyID int64) (*CompanyJobStats, error)

	// Job details management
	AddLocation(ctx context.Context, jobID int64, req *AddLocationRequest) (*JobLocation, error)
	UpdateLocation(ctx context.Context, locationID int64, req *UpdateLocationRequest) (*JobLocation, error)
	DeleteLocation(ctx context.Context, locationID int64) error
	SetPrimaryLocation(ctx context.Context, jobID, locationID int64) error

	AddBenefit(ctx context.Context, jobID int64, req *AddBenefitRequest) (*JobBenefit, error)
	UpdateBenefit(ctx context.Context, benefitID int64, req *UpdateBenefitRequest) (*JobBenefit, error)
	DeleteBenefit(ctx context.Context, benefitID int64) error
	BulkAddBenefits(ctx context.Context, jobID int64, benefits []AddBenefitRequest) error

	AddSkill(ctx context.Context, jobID int64, req *AddSkillRequest) (*JobSkill, error)
	UpdateSkill(ctx context.Context, jobSkillID int64, req *UpdateSkillRequest) (*JobSkill, error)
	DeleteSkill(ctx context.Context, jobSkillID int64) error
	BulkAddSkills(ctx context.Context, jobID int64, skills []AddSkillRequest) error

	AddRequirement(ctx context.Context, jobID int64, req *AddRequirementRequest) (*JobRequirement, error)
	UpdateRequirement(ctx context.Context, requirementID int64, req *UpdateRequirementRequest) (*JobRequirement, error)
	DeleteRequirement(ctx context.Context, requirementID int64) error
	BulkAddRequirements(ctx context.Context, jobID int64, requirements []AddRequirementRequest) error

	// Category management (Admin)
	CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*JobCategory, error)
	UpdateCategory(ctx context.Context, categoryID int64, req *UpdateCategoryRequest) (*JobCategory, error)
	DeleteCategory(ctx context.Context, categoryID int64) error
	GetCategory(ctx context.Context, categoryID int64) (*JobCategory, error)
	GetCategoryByCode(ctx context.Context, code string) (*JobCategory, error)
	ListCategories(ctx context.Context, filter CategoryFilter, page, limit int) ([]JobCategory, int64, error)
	GetCategoryTree(ctx context.Context) ([]JobCategory, error)
	GetActiveCategories(ctx context.Context) ([]JobCategory, error)

	// Subcategory management (Admin)
	CreateSubcategory(ctx context.Context, req *CreateSubcategoryRequest) (*JobSubcategory, error)
	UpdateSubcategory(ctx context.Context, subcategoryID int64, req *UpdateSubcategoryRequest) (*JobSubcategory, error)
	DeleteSubcategory(ctx context.Context, subcategoryID int64) error
	GetSubcategory(ctx context.Context, subcategoryID int64) (*JobSubcategory, error)
	ListSubcategories(ctx context.Context, categoryID int64) ([]JobSubcategory, error)
	GetActiveSubcategories(ctx context.Context, categoryID int64) ([]JobSubcategory, error)

	// Analytics and reporting
	GetJobAnalytics(ctx context.Context, jobID int64, startDate, endDate time.Time) (*JobAnalytics, error)
	GetCompanyAnalytics(ctx context.Context, companyID int64, startDate, endDate time.Time) (*CompanyAnalytics, error)
	GetCategoryAnalytics(ctx context.Context, categoryID int64, startDate, endDate time.Time) (*CategoryAnalytics, error)
	GetPopularCategories(ctx context.Context, limit int) ([]CategoryStats, error)
	GetTopCompanies(ctx context.Context, limit int) ([]CompanyStats, error)

	// Bulk operations
	BulkPublishJobs(ctx context.Context, jobIDs []int64) error
	BulkCloseJobs(ctx context.Context, jobIDs []int64) error
	BulkDeleteJobs(ctx context.Context, jobIDs []int64) error

	// Validation
	ValidateJob(ctx context.Context, job *Job) error
	CheckJobOwnership(ctx context.Context, jobID, employerUserID int64) error
	CheckJobStatus(ctx context.Context, jobID int64) (string, error)

	// Master Data Validation
	ValidateMasterDataIDs(ctx context.Context, industryID, districtID *int64) error
	GetJobWithMasterData(ctx context.Context, jobID int64) (*Job, error)
}

// ===== Request DTOs =====

// CreateJobRequest represents request to create a new job (master data only)
type CreateJobRequest struct {
	// Required Fields
	CompanyID      int64  `json:"company_id" validate:"required"`
	EmployerUserID int64  `json:"employer_user_id" validate:"required"`
	Description    string `json:"description" validate:"required"`

	// Master Data IDs (All Required)
	JobTitleID         int64 `json:"job_title_id" validate:"required"`
	JobTypeID          int64 `json:"job_type_id" validate:"required"`
	WorkPolicyID       int64 `json:"work_policy_id" validate:"required"`
	EducationLevelID   int64 `json:"education_level_id" validate:"required"`
	ExperienceLevelID  int64 `json:"experience_level_id" validate:"required"`
	GenderPreferenceID int64 `json:"gender_preference_id" validate:"required"`

	// Salary Range (Required)
	SalaryMin     *float64 `json:"salary_min" validate:"required,min=0"`
	SalaryMax     *float64 `json:"salary_max" validate:"required,min=0"`
	SalaryDisplay string   `json:"salary_display" validate:"required,oneof='range' 'min_only' 'max_only' 'negotiable' 'competitive' 'hidden'"`

	// Age Requirements (Optional)
	MinAge *int `json:"min_age" validate:"omitempty,min=17,max=65"`
	MaxAge *int `json:"max_age" validate:"omitempty,min=17,max=65,gtefield=MinAge"`

	// Company Address (Optional - for work location)
	CompanyAddressID *int64 `json:"company_address_id" validate:"omitempty,min=1"`

	// Skills (Required)
	Skills []AddSkillRequest `json:"skills" validate:"required,min=1"`
}

// UpdateJobRequest represents request to update job (master data only)
type UpdateJobRequest struct {
	// Internal - for ownership verification (not from API)
	EmployerUserID int64 `json:"employer_user_id,omitempty"` // Will be set by handler
	CompanyID      int64 `json:"company_id,omitempty"`       // Will be set by handler

	// Master Data IDs (Optional for updates)
	JobTitleID         *int64 `json:"job_title_id,omitempty"`
	JobTypeID          *int64 `json:"job_type_id,omitempty"`
	WorkPolicyID       *int64 `json:"work_policy_id,omitempty"`
	EducationLevelID   *int64 `json:"education_level_id,omitempty"`
	ExperienceLevelID  *int64 `json:"experience_level_id,omitempty"`
	GenderPreferenceID *int64 `json:"gender_preference_id,omitempty"`

	// Other Fields (Optional)
	Description      string   `json:"description,omitempty"`
	SalaryMin        *float64 `json:"salary_min,omitempty"`
	SalaryMax        *float64 `json:"salary_max,omitempty"`
	SalaryDisplay    *string  `json:"salary_display,omitempty" validate:"omitempty,oneof='range' 'min_only' 'max_only' 'negotiable' 'competitive' 'hidden'"`
	MinAge           *int     `json:"min_age,omitempty" validate:"omitempty,min=17,max=65"`
	MaxAge           *int     `json:"max_age,omitempty" validate:"omitempty,min=17,max=65,gtefield=MinAge"`
	CompanyAddressID *int64   `json:"company_address_id,omitempty" validate:"omitempty,min=1"`
}

// SaveJobDraftRequest represents request to save job draft (Phase 6)
type SaveJobDraftRequest struct {
	DraftID          *int64  `json:"draft_id"`          // Optional: for updating existing draft
	JobTitleID       int64   `json:"job_title_id"`      // Master data: job title ID
	JobCategoryID    int64   `json:"job_category_id"`   // Master data: job category ID
	JobTypeID        int64   `json:"job_type_id"`       // Master data: job type ID (full-time, part-time, etc.)
	WorkPolicyID     int64   `json:"work_policy_id"`    // Master data: work policy ID (onsite, remote, hybrid)
	GajiMin          int     `json:"gaji_min"`          // Minimum salary
	GajiMaks         int     `json:"gaji_maks"`         // Maximum salary
	AdaBonus         bool    `json:"ada_bonus"`         // Has bonus?
	GenderPreference string  `json:"gender_preference"` // Gender preference
	UmurMin          *int    `json:"umur_min"`          // Minimum age (nullable)
	UmurMaks         *int    `json:"umur_maks"`         // Maximum age (nullable)
	SkillIDs         []int64 `json:"skill_ids"`         // Array of skill IDs
	PendidikanID     int64   `json:"pendidikan_id"`     // Master data: education level ID
	PengalamanID     int64   `json:"pengalaman_id"`     // Master data: experience level ID
	Deskripsi        string  `json:"deskripsi"`         // Job description (will be sanitized for XSS)
}

// ApproveJobRequest represents request to approve a job (admin only)
type ApproveJobRequest struct {
	Notes string `json:"notes,omitempty" validate:"omitempty,max=500"`
}

// RejectJobRequest represents request to reject a job (admin only)
type RejectJobRequest struct {
	Reason string `json:"reason" validate:"required,max=500"`
}

// AddLocationRequest represents request to add job location
type AddLocationRequest struct {
	LocationType  string   `json:"location_type" validate:"omitempty,oneof='onsite' 'hybrid' 'remote'"`
	Address       string   `json:"address,omitempty"`
	City          string   `json:"city,omitempty"`
	Province      string   `json:"province,omitempty"`
	PostalCode    string   `json:"postal_code,omitempty"`
	Country       string   `json:"country,omitempty"`
	Latitude      *float64 `json:"latitude,omitempty"`
	Longitude     *float64 `json:"longitude,omitempty"`
	GooglePlaceID string   `json:"google_place_id,omitempty"`
	MapURL        string   `json:"map_url,omitempty"`
	IsPrimary     bool     `json:"is_primary"`
}

// UpdateLocationRequest represents request to update job location
type UpdateLocationRequest struct {
	LocationType  string   `json:"location_type,omitempty"`
	Address       string   `json:"address,omitempty"`
	City          string   `json:"city,omitempty"`
	Province      string   `json:"province,omitempty"`
	PostalCode    string   `json:"postal_code,omitempty"`
	Country       string   `json:"country,omitempty"`
	Latitude      *float64 `json:"latitude,omitempty"`
	Longitude     *float64 `json:"longitude,omitempty"`
	GooglePlaceID string   `json:"google_place_id,omitempty"`
	MapURL        string   `json:"map_url,omitempty"`
	IsPrimary     *bool    `json:"is_primary,omitempty"`
}

// AddBenefitRequest represents request to add job benefit
type AddBenefitRequest struct {
	BenefitID   *int64 `json:"benefit_id,omitempty"`
	BenefitName string `json:"benefit_name" validate:"required,max=150"`
	Description string `json:"description,omitempty"`
	IsHighlight bool   `json:"is_highlight"`
}

// UpdateBenefitRequest represents request to update job benefit
type UpdateBenefitRequest struct {
	BenefitName string `json:"benefit_name,omitempty"`
	Description string `json:"description,omitempty"`
	IsHighlight *bool  `json:"is_highlight,omitempty"`
}

// AddSkillRequest represents request to add job skill
type AddSkillRequest struct {
	SkillID         int64  `json:"skill_id" validate:"required"`
	ImportanceLevel string `json:"importance_level" validate:"omitempty,oneof='required' 'preferred' 'optional'"`
}

// UpdateSkillRequest represents request to update job skill
type UpdateSkillRequest struct {
	ImportanceLevel string   `json:"importance_level,omitempty"`
	Weight          *float64 `json:"weight,omitempty"`
}

// AddRequirementRequest represents request to add job requirement
type AddRequirementRequest struct {
	RequirementType string `json:"requirement_type" validate:"omitempty,oneof='education' 'experience' 'skill' 'language' 'certification' 'other'"`
	RequirementText string `json:"requirement_text" validate:"required"`
	SkillID         *int64 `json:"skill_id,omitempty"`
	MinExperience   *int16 `json:"min_experience,omitempty"`
	MaxExperience   *int16 `json:"max_experience,omitempty"`
	EducationLevel  string `json:"education_level,omitempty"`
	Language        string `json:"language,omitempty"`
	IsMandatory     bool   `json:"is_mandatory"`
	Priority        int16  `json:"priority" validate:"omitempty,min=1"`
}

// UpdateRequirementRequest represents request to update job requirement
type UpdateRequirementRequest struct {
	RequirementType string `json:"requirement_type,omitempty"`
	RequirementText string `json:"requirement_text,omitempty"`
	SkillID         *int64 `json:"skill_id,omitempty"`
	MinExperience   *int16 `json:"min_experience,omitempty"`
	MaxExperience   *int16 `json:"max_experience,omitempty"`
	EducationLevel  string `json:"education_level,omitempty"`
	Language        string `json:"language,omitempty"`
	IsMandatory     *bool  `json:"is_mandatory,omitempty"`
	Priority        *int16 `json:"priority,omitempty"`
}

// CreateCategoryRequest represents request to create job category
type CreateCategoryRequest struct {
	ParentID    *int64 `json:"parent_id,omitempty"`
	Code        string `json:"code" validate:"required,max=30"`
	Name        string `json:"name" validate:"required,max=150"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
}

// UpdateCategoryRequest represents request to update job category
type UpdateCategoryRequest struct {
	ParentID    *int64 `json:"parent_id,omitempty"`
	Code        string `json:"code,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// CreateSubcategoryRequest represents request to create job subcategory
type CreateSubcategoryRequest struct {
	CategoryID  int64  `json:"category_id" validate:"required"`
	Code        string `json:"code" validate:"required,max=50"`
	Name        string `json:"name" validate:"required,max=150"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
}

// UpdateSubcategoryRequest represents request to update job subcategory
type UpdateSubcategoryRequest struct {
	CategoryID  *int64 `json:"category_id,omitempty"`
	Code        string `json:"code,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// ===== Response DTOs =====

// JobSearchResponse represents job search results with metadata
type JobSearchResponse struct {
	Jobs        []Job         `json:"jobs"`
	Total       int64         `json:"total"`
	Page        int           `json:"page"`
	Limit       int           `json:"limit"`
	TotalPages  int           `json:"total_pages"`
	Facets      *SearchFacets `json:"facets,omitempty"`
	Suggestions []string      `json:"suggestions,omitempty"`
}

// SearchFacets represents search facets/filters
type SearchFacets struct {
	Categories      []FacetItem `json:"categories"`
	Locations       []FacetItem `json:"locations"`
	JobLevels       []FacetItem `json:"job_levels"`
	EmploymentTypes []FacetItem `json:"employment_types"`
	SalaryRanges    []FacetItem `json:"salary_ranges"`
}

// FacetItem represents a facet item with count
type FacetItem struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

// MatchScore represents job-user match score
type MatchScore struct {
	JobID           int64    `json:"job_id"`
	UserID          int64    `json:"user_id"`
	OverallScore    float64  `json:"overall_score"`
	SkillScore      float64  `json:"skill_score"`
	ExperienceScore float64  `json:"experience_score"`
	EducationScore  float64  `json:"education_score"`
	LocationScore   float64  `json:"location_score"`
	MatchedSkills   []string `json:"matched_skills"`
	MissingSkills   []string `json:"missing_skills"`
	Recommendation  string   `json:"recommendation"`
}

// MatchResponse represents matching jobs response
type MatchResponse struct {
	Jobs       []JobWithScore `json:"jobs"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// JobWithScore represents job with match score
type JobWithScore struct {
	Job        Job     `json:"job"`
	MatchScore float64 `json:"match_score"`
}

// JobAnalytics represents job analytics data
type JobAnalytics struct {
	JobID              int64            `json:"job_id"`
	Period             string           `json:"period"`
	ViewsData          []TimeSeriesData `json:"views_data"`
	ApplicationsData   []TimeSeriesData `json:"applications_data"`
	TotalViews         int64            `json:"total_views"`
	TotalApplications  int64            `json:"total_applications"`
	UniqueViewers      int64            `json:"unique_viewers"`
	ConversionRate     float64          `json:"conversion_rate"`
	AverageTimeToApply float64          `json:"average_time_to_apply"`
	TopSources         []SourceStats    `json:"top_sources"`
}

// CompanyAnalytics represents company job analytics
type CompanyAnalytics struct {
	CompanyID         int64            `json:"company_id"`
	Period            string           `json:"period"`
	TotalJobs         int64            `json:"total_jobs"`
	ActiveJobs        int64            `json:"active_jobs"`
	TotalViews        int64            `json:"total_views"`
	TotalApplications int64            `json:"total_applications"`
	TopJobs           []JobPerformance `json:"top_jobs"`
	CategoryBreakdown []CategoryStats  `json:"category_breakdown"`
}

// CategoryAnalytics represents category analytics
type CategoryAnalytics struct {
	CategoryID        int64          `json:"category_id"`
	CategoryName      string         `json:"category_name"`
	Period            string         `json:"period"`
	TotalJobs         int64          `json:"total_jobs"`
	ActiveJobs        int64          `json:"active_jobs"`
	TotalViews        int64          `json:"total_views"`
	TotalApplications int64          `json:"total_applications"`
	TopCompanies      []CompanyStats `json:"top_companies"`
}

// TimeSeriesData represents time-series data point
type TimeSeriesData struct {
	Date  time.Time `json:"date"`
	Value int64     `json:"value"`
}

// SourceStats represents traffic source statistics
type SourceStats struct {
	Source string `json:"source"`
	Count  int64  `json:"count"`
}

// JobPerformance represents job performance metrics
type JobPerformance struct {
	JobID          int64   `json:"job_id"`
	JobTitle       string  `json:"job_title"`
	Views          int64   `json:"views"`
	Applications   int64   `json:"applications"`
	ConversionRate float64 `json:"conversion_rate"`
}

// CompanyStats represents company statistics
type CompanyStats struct {
	CompanyID         int64  `json:"company_id"`
	CompanyName       string `json:"company_name"`
	TotalJobs         int64  `json:"total_jobs"`
	ActiveJobs        int64  `json:"active_jobs"`
	TotalViews        int64  `json:"total_views"`
	TotalApplications int64  `json:"total_applications"`
}
