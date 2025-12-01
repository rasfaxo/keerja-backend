package request

// CreateJobRequest represents job creation request
// Only uses master data IDs - no legacy format support
type CreateJobRequest struct {
	// Required Fields
	CompanyID   int64  `json:"company_id" validate:"required,min=1"`
	Description string `json:"description" validate:"required,min=50"`

	// Category/Subcategory (optional but recommended)
	JobCategoryID    int64 `json:"job_category_id" validate:"omitempty,min=1"`
	JobSubcategoryID int64 `json:"job_subcategory_id" validate:"omitempty,min=1"`

	// Master Data IDs (All Required)
	JobTitleID         *int64 `json:"job_title_id" validate:"omitempty,min=1"`
	JobTypeID          int64  `json:"job_type_id" validate:"required,min=1"`
	WorkPolicyID       int64  `json:"work_policy_id" validate:"required,min=1"`
	EducationLevelID   int64  `json:"education_level_id" validate:"required,min=1"`
	ExperienceLevelID  int64  `json:"experience_level_id" validate:"required,min=1"`
	GenderPreferenceID int64  `json:"gender_preference_id" validate:"required,min=1"`

	// Salary Range (Required - validation depends on salary_display mode)
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
	SkillID         int64  `json:"skill_id" validate:"required"`
	ImportanceLevel string `json:"importance_level" validate:"omitempty,oneof='required' 'preferred' 'optional'"`
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

// UpdateJobRequest represents job update request (only master data format)
type UpdateJobRequest struct {
	// Master Data IDs (Optional for updates)
	JobTitleID         *int64 `json:"job_title_id" validate:"omitempty,min=1"`
	JobTypeID          *int64 `json:"job_type_id" validate:"omitempty,min=1"`
	WorkPolicyID       *int64 `json:"work_policy_id" validate:"omitempty,min=1"`
	EducationLevelID   *int64 `json:"education_level_id" validate:"omitempty,min=1"`
	ExperienceLevelID  *int64 `json:"experience_level_id" validate:"omitempty,min=1"`
	GenderPreferenceID *int64 `json:"gender_preference_id" validate:"omitempty,min=1"`

	// Other Fields (Optional)
	Description      *string           `json:"description" validate:"omitempty,min=50"`
	SalaryMin        *float64          `json:"salary_min" validate:"omitempty,min=0"`
	SalaryMax        *float64          `json:"salary_max" validate:"omitempty,min=0"`
	SalaryDisplay    *string           `json:"salary_display" validate:"omitempty,oneof='range' 'min_only' 'max_only' 'negotiable' 'competitive' 'hidden'"`
	MinAge           *int              `json:"min_age" validate:"omitempty,min=17,max=65"`
	MaxAge           *int              `json:"max_age" validate:"omitempty,min=17,max=65,gtefield=MinAge"`
	CompanyAddressID *int64            `json:"company_address_id" validate:"omitempty,min=1"`
	Skills           []AddSkillRequest `json:"skills,omitempty" validate:"omitempty,dive"`
	// NOTE: Status should NOT be updated by users - it's controlled by workflow
	// - draft: initial state (automatic)
	// - pending_approval: submitted for review (automatic when published)
	// - published: approved by admin (admin action)
	// - rejected: rejected by admin (admin action)
	// - closed/expired/suspended: lifecycle management (admin/system action)
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

// SaveJobDraftRequest represents save job draft request (Phase 6)
type SaveJobDraftRequest struct {
	DraftID          *int64  `json:"draft_id" validate:"omitempty,min=1"`                           // Optional: for updating existing draft
	CompanyID        int64   `json:"company_id" validate:"required,min=1"`                          // Company ID for the job draft
	JobTitleID       *int64  `json:"job_title_id" validate:"omitempty,min=1"`                       // Master data: job title ID
	JobCategoryID    int64   `json:"job_category_id" validate:"required,min=1"`                     // Master data: job category ID
	JobSubcategoryID int64   `json:"job_subcategory_id" validate:"omitempty,min=1"`                 // Master data: job subcategory ID (optional)
	JobTypeID        int64   `json:"job_type_id" validate:"required,min=1"`                         // Master data: job type ID (full-time, part-time, etc.)
	WorkPolicyID     int64   `json:"work_policy_id" validate:"required,min=1"`                      // Master data: work policy ID (onsite, remote, hybrid)
	GajiMin          int     `json:"gaji_min" validate:"required,min=0"`                            // Minimum salary
	GajiMaks         int     `json:"gaji_maks" validate:"required,min=0,gtefield=GajiMin"`          // Maximum salary (must be >= GajiMin)
	AdaBonus         bool    `json:"ada_bonus"`                                                     // Has bonus?
	GenderPreference string  `json:"gender_preference" validate:"required,oneof=male female any"`   // Gender preference
	UmurMin          *int    `json:"umur_min" validate:"omitempty,min=17,max=65"`                   // Minimum age (nullable)
	UmurMaks         *int    `json:"umur_maks" validate:"omitempty,min=17,max=65,gtefield=UmurMin"` // Maximum age (nullable, must be >= UmurMin)
	SkillIDs         []int64 `json:"skill_ids" validate:"required,min=1,dive,min=1"`                // Array of skill IDs (at least 1)
	PendidikanID     int64   `json:"pendidikan_id" validate:"required,min=1"`                       // Master data: education level ID
	PengalamanID     int64   `json:"pengalaman_id" validate:"required,min=1"`                       // Master data: experience level ID
	Deskripsi        string  `json:"deskripsi" validate:"required,min=50,max=5000"`                 // Job description (will be sanitized for XSS)
}
