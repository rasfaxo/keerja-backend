package request

// CreateJobRequest represents job creation request
type CreateJobRequest struct {
	CompanyID        int64                      `json:"company_id" validate:"required,min=1"`
	EmployerUserID   int64                      `json:"employer_user_id" validate:"required,min=1"`
	CategoryID       *int64                     `json:"category_id" validate:"omitempty,min=1"`
	Title            string                     `json:"title" validate:"required,min=3,max=200"`
	JobLevel         string                     `json:"job_level" validate:"omitempty,oneof='Internship' 'Entry Level' 'Mid Level' 'Senior Level' 'Manager' 'Director'"`
	EmploymentType   string                     `json:"employment_type" validate:"required,oneof='Full-Time' 'Part-Time' 'Contract' 'Internship' 'Freelance'"`
	Description      string                     `json:"description" validate:"required,min=50"`
	RequirementsText string                     `json:"requirements_text" validate:"omitempty"`
	Responsibilities string                     `json:"responsibilities" validate:"omitempty"`
	RemoteOption     bool                       `json:"remote_option"`
	SalaryMin        *float64                   `json:"salary_min" validate:"omitempty,min=0"`
	SalaryMax        *float64                   `json:"salary_max" validate:"omitempty,min=0,gtefield=SalaryMin"`
	Currency         string                     `json:"currency" validate:"omitempty,len=3"`
	ExperienceMin    *int16                     `json:"experience_min" validate:"omitempty,min=0"`
	ExperienceMax    *int16                     `json:"experience_max" validate:"omitempty,min=0,gtefield=ExperienceMin"`
	EducationLevel   string                     `json:"education_level" validate:"omitempty"`
	TotalHires       int16                      `json:"total_hires" validate:"omitempty,min=1"`
	ExpiredAt        *string                    `json:"expired_at" validate:"omitempty"`
	Locations        []CreateJobLocationRequest `json:"locations" validate:"omitempty"`
	Skills           []AddSkillRequest          `json:"skills" validate:"omitempty"`
	Benefits         []AddBenefitRequest        `json:"benefits" validate:"omitempty"`
	JobRequirements  []AddRequirementRequest    `json:"job_requirements" validate:"omitempty"`
}

// CreateJobLocationRequest represents additional job location
type CreateJobLocationRequest struct {
	LocationType  string   `json:"location_type" validate:"omitempty,oneof='onsite' 'hybrid' 'remote'"`
	City          string   `json:"city" validate:"required,max=100"`
	Province      string   `json:"province" validate:"required,max=100"`
	Address       string   `json:"address" validate:"omitempty"`
	PostalCode    string   `json:"postal_code" validate:"omitempty,max=10"`
	Country       string   `json:"country" validate:"omitempty,max=100"`
	Latitude      *float64 `json:"latitude" validate:"omitempty"`
	Longitude     *float64 `json:"longitude" validate:"omitempty"`
	GooglePlaceID string   `json:"google_place_id" validate:"omitempty"`
	MapURL        string   `json:"map_url" validate:"omitempty,url"`
	IsPrimary     bool     `json:"is_primary"`
}

// AddSkillRequest represents add skill to job request
type AddSkillRequest struct {
	SkillID         int64   `json:"skill_id" validate:"required"`
	ImportanceLevel string  `json:"importance_level" validate:"omitempty,oneof='required' 'preferred' 'optional'"`
	Weight          float64 `json:"weight" validate:"omitempty,min=0,max=1"`
}

// AddBenefitRequest represents add benefit to job request
type AddBenefitRequest struct {
	BenefitID   *int64 `json:"benefit_id" validate:"omitempty"`
	BenefitName string `json:"benefit_name" validate:"required,max=150"`
	Description string `json:"description" validate:"omitempty"`
	IsHighlight bool   `json:"is_highlight"`
}

// AddRequirementRequest represents add requirement to job request
type AddRequirementRequest struct {
	RequirementType string `json:"requirement_type" validate:"omitempty,oneof='education' 'experience' 'skill' 'language' 'certification' 'other'"`
	RequirementText string `json:"requirement_text" validate:"required"`
	SkillID         *int64 `json:"skill_id" validate:"omitempty"`
	MinExperience   *int16 `json:"min_experience" validate:"omitempty"`
	MaxExperience   *int16 `json:"max_experience" validate:"omitempty"`
	EducationLevel  string `json:"education_level" validate:"omitempty"`
	Language        string `json:"language" validate:"omitempty"`
	IsMandatory     bool   `json:"is_mandatory"`
	Priority        int16  `json:"priority" validate:"omitempty,min=1"`
}

// UpdateJobRequest represents job update request
type UpdateJobRequest struct {
	CategoryID       *int64   `json:"category_id" validate:"omitempty,min=1"`
	Title            *string  `json:"title" validate:"omitempty,min=3,max=200"`
	JobLevel         *string  `json:"job_level" validate:"omitempty,oneof='Internship' 'Entry Level' 'Mid Level' 'Senior Level' 'Manager' 'Director'"`
	EmploymentType   *string  `json:"employment_type" validate:"omitempty,oneof='Full-Time' 'Part-Time' 'Contract' 'Internship' 'Freelance'"`
	Description      *string  `json:"description" validate:"omitempty,min=50"`
	RequirementsText *string  `json:"requirements_text" validate:"omitempty"`
	Responsibilities *string  `json:"responsibilities" validate:"omitempty"`
	RemoteOption     *bool    `json:"remote_option"`
	SalaryMin        *float64 `json:"salary_min" validate:"omitempty,min=0"`
	SalaryMax        *float64 `json:"salary_max" validate:"omitempty,min=0"`
	Currency         *string  `json:"currency" validate:"omitempty,len=3"`
	ExperienceMin    *int16   `json:"experience_min" validate:"omitempty,min=0"`
	ExperienceMax    *int16   `json:"experience_max" validate:"omitempty,min=0"`
	EducationLevel   *string  `json:"education_level" validate:"omitempty"`
	TotalHires       *int16   `json:"total_hires" validate:"omitempty,min=1"`
	Status           *string  `json:"status" validate:"omitempty,oneof=draft published closed expired suspended"`
	ExpiredAt        *string  `json:"expired_at" validate:"omitempty"`
}

// UpdateJobSkillsRequest represents update job skills request
type UpdateJobSkillsRequest struct {
	SkillIDs []int64 `json:"skill_ids" validate:"required"`
}

// UpdateJobBenefitsRequest represents update job benefits request
type UpdateJobBenefitsRequest struct {
	BenefitIDs []int64 `json:"benefit_ids" validate:"required"`
}

// JobSearchRequest represents job search request
type JobSearchRequest struct {
	Query          string   `json:"query" query:"q" validate:"omitempty"`
	Keywords       []string `json:"keywords" query:"keywords" validate:"omitempty"`
	CategoryID     *int64   `json:"category_id" query:"category_id" validate:"omitempty"`
	Location       string   `json:"location" query:"location" validate:"omitempty"`
	City           string   `json:"city" query:"city" validate:"omitempty"`
	Province       string   `json:"province" query:"province" validate:"omitempty"`
	RemoteOnly     bool     `json:"remote_only" query:"remote_only"`
	EmploymentType string   `json:"employment_type" query:"employment_type" validate:"omitempty"`
	JobLevel       string   `json:"job_level" query:"job_level" validate:"omitempty"`
	SalaryMin      *float64 `json:"salary_min" query:"salary_min" validate:"omitempty"`
	SalaryMax      *float64 `json:"salary_max" query:"salary_max" validate:"omitempty"`
	ExperienceMin  *int16   `json:"experience_min" query:"experience_min" validate:"omitempty"`
	ExperienceMax  *int16   `json:"experience_max" query:"experience_max" validate:"omitempty"`
	EducationLevel string   `json:"education_level" query:"education_level" validate:"omitempty"`
	CompanyID      *int64   `json:"company_id" query:"company_id" validate:"omitempty"`
	SkillIDs       []int64  `json:"skill_ids" query:"skill_ids" validate:"omitempty"`
	BenefitIDs     []int64  `json:"benefit_ids" query:"benefit_ids" validate:"omitempty"`
	PostedWithin   *int     `json:"posted_within" query:"posted_within" validate:"omitempty,min=1"` // in days
	Page           int      `json:"page" query:"page" validate:"omitempty,min=1"`
	Limit          int      `json:"limit" query:"limit" validate:"omitempty,min=1,max=100"`
	SortBy         string   `json:"sort_by" query:"sort_by" validate:"omitempty,oneof=relevance posted_date salary views applications"`
	SortOrder      string   `json:"sort_order" query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// JobFilterRequest represents job filter request
type JobFilterRequest struct {
	Status         string `json:"status" query:"status" validate:"omitempty,oneof=draft published closed expired suspended"`
	IsExpired      *bool  `json:"is_expired" query:"is_expired"`
	CompanyID      *int64 `json:"company_id" query:"company_id" validate:"omitempty"`
	EmployerUserID *int64 `json:"employer_user_id" query:"employer_user_id" validate:"omitempty"`
	CategoryID     *int64 `json:"category_id" query:"category_id" validate:"omitempty"`
	Page           int    `json:"page" query:"page" validate:"omitempty,min=1"`
	Limit          int    `json:"limit" query:"limit" validate:"omitempty,min=1,max=100"`
	SortBy         string `json:"sort_by" query:"sort_by" validate:"omitempty"`
	SortOrder      string `json:"sort_order" query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// PublishJobRequest represents publish job request
type PublishJobRequest struct {
	PublishAt *string `json:"publish_at" validate:"omitempty"` // Optional: schedule publish
	ExpiredAt *string `json:"expired_at" validate:"omitempty"` // Optional: set expiration
}

// CloseJobRequest represents close job request
type CloseJobRequest struct {
	Reason string `json:"reason" validate:"omitempty,max=500"`
}

// JobRecommendationRequest represents job recommendation request for user
type JobRecommendationRequest struct {
	UserID int64 `json:"user_id" validate:"required,min=1"`
	Limit  int   `json:"limit" validate:"omitempty,min=1,max=50"`
}
