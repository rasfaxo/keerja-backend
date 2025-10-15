package response

import "time"

// JobResponse represents job public response
type JobResponse struct {
	ID                int64      `json:"id"`
	UUID              string     `json:"uuid"`
	CompanyID         int64      `json:"company_id"`
	CompanyName       string     `json:"company_name"`
	CompanyLogoURL    string     `json:"company_logo_url,omitempty"`
	CompanyVerified   bool       `json:"company_verified"`
	Title             string     `json:"title"`
	Slug              string     `json:"slug"`
	JobLevel          string     `json:"job_level,omitempty"`
	EmploymentType    string     `json:"employment_type"`
	Location          string     `json:"location,omitempty"`
	City              string     `json:"city,omitempty"`
	Province          string     `json:"province,omitempty"`
	RemoteOption      bool       `json:"remote_option"`
	SalaryMin         *float64   `json:"salary_min,omitempty"`
	SalaryMax         *float64   `json:"salary_max,omitempty"`
	Currency          string     `json:"currency"`
	ExperienceMin     *int16     `json:"experience_min,omitempty"`
	ExperienceMax     *int16     `json:"experience_max,omitempty"`
	EducationLevel    string     `json:"education_level,omitempty"`
	Status            string     `json:"status"`
	ViewsCount        int64      `json:"views_count"`
	ApplicationsCount int64      `json:"applications_count"`
	PublishedAt       *time.Time `json:"published_at,omitempty"`
	ExpiredAt         *time.Time `json:"expired_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	IsExpired         bool       `json:"is_expired"`
	DaysRemaining     *int       `json:"days_remaining,omitempty"`
}

// JobDetailResponse represents detailed job response
type JobDetailResponse struct {
	ID                int64                    `json:"id"`
	UUID              string                   `json:"uuid"`
	CompanyID         int64                    `json:"company_id"`
	CompanyName       string                   `json:"company_name"`
	CompanyLogoURL    string                   `json:"company_logo_url,omitempty"`
	CompanyVerified   bool                     `json:"company_verified"`
	CompanySlug       string                   `json:"company_slug"`
	EmployerUserID    *int64                   `json:"employer_user_id,omitempty"`
	CategoryID        *int64                   `json:"category_id,omitempty"`
	CategoryName      string                   `json:"category_name,omitempty"`
	Title             string                   `json:"title"`
	Slug              string                   `json:"slug"`
	JobLevel          string                   `json:"job_level,omitempty"`
	EmploymentType    string                   `json:"employment_type"`
	Description       string                   `json:"description"`
	Requirements      string                   `json:"requirements,omitempty"`
	Responsibilities  string                   `json:"responsibilities,omitempty"`
	Location          string                   `json:"location,omitempty"`
	City              string                   `json:"city,omitempty"`
	Province          string                   `json:"province,omitempty"`
	RemoteOption      bool                     `json:"remote_option"`
	SalaryMin         *float64                 `json:"salary_min,omitempty"`
	SalaryMax         *float64                 `json:"salary_max,omitempty"`
	Currency          string                   `json:"currency"`
	ExperienceMin     *int16                   `json:"experience_min,omitempty"`
	ExperienceMax     *int16                   `json:"experience_max,omitempty"`
	EducationLevel    string                   `json:"education_level,omitempty"`
	TotalHires        int16                    `json:"total_hires"`
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
	Requirements_     []JobRequirementResponse `json:"job_requirements,omitempty"`
	HasApplied        bool                     `json:"has_applied,omitempty"` // For authenticated users
	IsSaved           bool                     `json:"is_saved,omitempty"`    // For authenticated users
}

// JobSkillResponse represents job skill response
type JobSkillResponse struct {
	ID          int64  `json:"id"`
	SkillID     int64  `json:"skill_id"`
	SkillName   string `json:"skill_name"`
	IsRequired  bool   `json:"is_required"`
	MinYearsExp *int16 `json:"min_years_exp,omitempty"`
}

// JobBenefitResponse represents job benefit response
type JobBenefitResponse struct {
	ID          int64  `json:"id"`
	BenefitID   int64  `json:"benefit_id"`
	BenefitName string `json:"benefit_name"`
	Description string `json:"description,omitempty"`
}

// JobLocationResponse represents job location response
type JobLocationResponse struct {
	ID        int64    `json:"id"`
	City      string   `json:"city"`
	Province  string   `json:"province"`
	Address   string   `json:"address,omitempty"`
	IsRemote  bool     `json:"is_remote"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

// JobRequirementResponse represents job requirement response
type JobRequirementResponse struct {
	ID              int64  `json:"id"`
	RequirementType string `json:"requirement_type"`
	Description     string `json:"description"`
	IsRequired      bool   `json:"is_required"`
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
