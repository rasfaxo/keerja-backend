package response

import "time"

// JobResponse represents job public response (simplified - master data only)
type JobResponse struct {
	ID              int64  `json:"id"`
	UUID            string `json:"uuid"`
	CompanyID       int64  `json:"company_id"`
	CompanyName     string `json:"company_name"`
	CompanyLogoURL  string `json:"company_logo_url,omitempty"`
	CompanyVerified bool   `json:"company_verified"`
	Title           string `json:"title"`
	Slug            string `json:"slug"`

	// Master Data
	JobType    *JobMasterDataItem `json:"job_type,omitempty"`
	WorkPolicy *JobMasterDataItem `json:"work_policy,omitempty"`

	SalaryMin         *float64   `json:"salary_min,omitempty"`
	SalaryMax         *float64   `json:"salary_max,omitempty"`
	Currency          string     `json:"currency"`
	Status            string     `json:"status"`
	ViewsCount        int64      `json:"views_count"`
	ApplicationsCount int64      `json:"applications_count"`
	PublishedAt       *time.Time `json:"published_at,omitempty"`
	ExpiredAt         *time.Time `json:"expired_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	IsExpired         bool       `json:"is_expired"`
	DaysRemaining     *int       `json:"days_remaining,omitempty"`
}

// JobDetailResponse represents detailed job response (master data only)
type JobDetailResponse struct {
	ID              int64  `json:"id"`
	UUID            string `json:"uuid"`
	CompanyID       int64  `json:"company_id"`
	CompanyName     string `json:"company_name"`
	CompanyLogoURL  string `json:"company_logo_url,omitempty"`
	CompanyVerified bool   `json:"company_verified"`
	CompanySlug     string `json:"company_slug"`
	EmployerUserID  *int64 `json:"employer_user_id,omitempty"`
	Title           string `json:"title"`
	Slug            string `json:"slug"`

	// Master Data IDs
	JobTitleID         int64 `json:"job_title_id"`
	JobTypeID          int64 `json:"job_type_id"`
	WorkPolicyID       int64 `json:"work_policy_id"`
	EducationLevelID   int64 `json:"education_level_id"`
	ExperienceLevelID  int64 `json:"experience_level_id"`
	GenderPreferenceID int64 `json:"gender_preference_id"`

	// Master Data Details
	JobTitle         *JobMasterDataItem `json:"job_title,omitempty"`
	JobType          *JobMasterDataItem `json:"job_type,omitempty"`
	WorkPolicy       *JobMasterDataItem `json:"work_policy,omitempty"`
	EducationLevel   *JobMasterDataItem `json:"education_level,omitempty"`
	ExperienceLevel  *JobMasterDataItem `json:"experience_level,omitempty"`
	GenderPreference *JobMasterDataItem `json:"gender_preference,omitempty"`

	Description string `json:"description"`

	SalaryMin         *float64                 `json:"salary_min,omitempty"`
	SalaryMax         *float64                 `json:"salary_max,omitempty"`
	Currency          string                   `json:"currency"`
	Status            string                   `json:"status"`
	ViewsCount        int64                    `json:"views_count"`
	ApplicationsCount int64                    `json:"applications_count"`
	PublishedAt       *time.Time               `json:"published_at,omitempty"`
	ExpiredAt         *time.Time               `json:"expired_at,omitempty"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
	IsExpired         bool                     `json:"is_expired"`
	DaysRemaining     *int                     `json:"days_remaining,omitempty"`
	Skills            []JobSkillResponse       `json:"skills,omitempty"`
	Benefits          []JobBenefitResponse     `json:"benefits,omitempty"`
	Locations         []JobLocationResponse    `json:"locations,omitempty"`
	JobRequirements   []JobRequirementResponse `json:"job_requirements,omitempty"`
	HasApplied        bool                     `json:"has_applied,omitempty"`
	IsSaved           bool                     `json:"is_saved,omitempty"`
}

// JobMasterDataItem represents a generic master data item for job details
type JobMasterDataItem struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// JobSkillResponse represents job skill response
type JobSkillResponse struct {
	ID              int64   `json:"id"`
	SkillID         int64   `json:"skill_id"`
	SkillName       string  `json:"skill_name"`
	ImportanceLevel string  `json:"importance_level"`
	Weight          float64 `json:"weight"`
}

// JobBenefitResponse represents job benefit response
type JobBenefitResponse struct {
	ID          int64  `json:"id"`
	BenefitID   *int64 `json:"benefit_id,omitempty"`
	BenefitName string `json:"benefit_name"`
	Description string `json:"description,omitempty"`
	IsHighlight bool   `json:"is_highlight"`
}

// JobLocationResponse represents job location response
type JobLocationResponse struct {
	ID            int64    `json:"id"`
	LocationType  string   `json:"location_type"`
	Address       string   `json:"address,omitempty"`
	City          string   `json:"city,omitempty"`
	Province      string   `json:"province,omitempty"`
	PostalCode    string   `json:"postal_code,omitempty"`
	Country       string   `json:"country"`
	Latitude      *float64 `json:"latitude,omitempty"`
	Longitude     *float64 `json:"longitude,omitempty"`
	GooglePlaceID string   `json:"google_place_id,omitempty"`
	MapURL        string   `json:"map_url,omitempty"`
	IsPrimary     bool     `json:"is_primary"`
}

// JobRequirementResponse represents job requirement response
type JobRequirementResponse struct {
	ID              int64  `json:"id"`
	RequirementType string `json:"requirement_type"`
	RequirementText string `json:"requirement_text"`
	SkillID         *int64 `json:"skill_id,omitempty"`
	MinExperience   *int16 `json:"min_experience,omitempty"`
	MaxExperience   *int16 `json:"max_experience,omitempty"`
	EducationLevel  string `json:"education_level,omitempty"`
	Language        string `json:"language,omitempty"`
	IsMandatory     bool   `json:"is_mandatory"`
	Priority        int16  `json:"priority"`
}

// JobCategoryResponse represents job category response
type JobCategoryResponse struct {
	ID            int64                    `json:"id"`
	CategoryName  string                   `json:"category_name"`
	Description   string                   `json:"description,omitempty"`
	IconURL       string                   `json:"icon_url,omitempty"`
	JobsCount     int64                    `json:"jobs_count"`
	Subcategories []JobSubcategoryResponse `json:"subcategories,omitempty"`
}

// JobSubcategoryResponse represents job subcategory response
type JobSubcategoryResponse struct {
	ID              int64  `json:"id"`
	SubcategoryName string `json:"subcategory_name"`
	Description     string `json:"description,omitempty"`
	JobsCount       int64  `json:"jobs_count"`
}

// JobListResponse represents list of jobs response
type JobListResponse struct {
	Jobs []JobResponse `json:"jobs"`
}

// JobStatsResponse represents job statistics response
type JobStatsResponse struct {
	TotalViews         int64   `json:"total_views"`
	TotalApplications  int64   `json:"total_applications"`
	NewApplications    int64   `json:"new_applications"`
	ViewsToday         int64   `json:"views_today"`
	ApplicationsToday  int64   `json:"applications_today"`
	AverageTimeToApply float64 `json:"average_time_to_apply"` // in hours
	ApplicationRate    float64 `json:"application_rate"`      // percentage
}

// JobRecommendationResponse represents recommended jobs response
type JobRecommendationResponse struct {
	Job        JobResponse `json:"job"`
	MatchScore float64     `json:"match_score"`
	Reasons    []string    `json:"reasons,omitempty"`
}

// SimilarJobsResponse represents similar jobs response
type SimilarJobsResponse struct {
	Jobs []JobResponse `json:"jobs"`
}

// JobLocationDetail represents hierarchical location information from master data
type JobLocationDetail struct {
	District *DistrictResponse `json:"district,omitempty"`
	City     *CityResponse     `json:"city,omitempty"`
	Province *ProvinceResponse `json:"province,omitempty"`
}
