package job

import (
	"time"

	"keerja-backend/internal/domain/master"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Job represents a job posting entity
type Job struct {
	ID               int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UUID             uuid.UUID `gorm:"column:uuid;type:uuid;default:gen_random_uuid();unique" json:"uuid"`
	CompanyID        int64     `gorm:"column:company_id;not null;index" json:"company_id" validate:"required"`
	EmployerUserID   *int64    `gorm:"column:employer_user_id;index" json:"employer_user_id,omitempty"`
	CategoryID       *int64    `gorm:"column:category_id;index" json:"category_id,omitempty"`
	Title            string    `gorm:"column:title;type:varchar(200);not null" json:"title" validate:"required,max=200"`
	Slug             string    `gorm:"column:slug;type:varchar(220);unique" json:"slug"`
	Description      string    `gorm:"column:description;type:text;not null" json:"description" validate:"required"`
	RequirementsText string    `gorm:"column:requirements;type:text" json:"requirements_text,omitempty"`
	Responsibilities string    `gorm:"column:responsibilities;type:text" json:"responsibilities,omitempty"`

	// Master Data Relations - Job Master Data (New FK columns)
	JobTitleID         *int64 `gorm:"column:job_title_id;index" json:"job_title_id,omitempty"`
	JobTypeID          *int64 `gorm:"column:job_type_id;index" json:"job_type_id,omitempty"`
	WorkPolicyID       *int64 `gorm:"column:work_policy_id;index" json:"work_policy_id,omitempty"`
	EducationLevelID   *int64 `gorm:"column:education_level_id;index" json:"education_level_id,omitempty"`
	ExperienceLevelID  *int64 `gorm:"column:experience_level_id;index" json:"experience_level_id,omitempty"`
	GenderPreferenceID *int64 `gorm:"column:gender_preference_id;index" json:"gender_preference_id,omitempty"`

	// Legacy Location Fields (for backward compatibility)
	Location     string `gorm:"column:location;type:varchar(150)" json:"location,omitempty"`
	City         string `gorm:"column:city;type:varchar(100)" json:"city,omitempty"`
	Province     string `gorm:"column:province;type:varchar(100)" json:"province,omitempty"`
	RemoteOption bool   `gorm:"column:remote_option;default:false" json:"remote_option"`

	SalaryMin     *float64 `gorm:"column:salary_min;type:numeric(12,2)" json:"salary_min,omitempty"`
	SalaryMax     *float64 `gorm:"column:salary_max;type:numeric(12,2)" json:"salary_max,omitempty"`
	SalaryDisplay string   `gorm:"column:salary_display;type:varchar(20);default:'range'" json:"salary_display"`
	Currency      string   `gorm:"column:currency;type:varchar(10);default:'IDR'" json:"currency" validate:"omitempty,len=3"`

	// Age Requirements
	MinAge        *int   `gorm:"column:min_age" json:"min_age,omitempty" validate:"omitempty,min=17,max=100"`
	MaxAge        *int   `gorm:"column:max_age" json:"max_age,omitempty" validate:"omitempty,min=17,max=100"`
	ExperienceMin *int16 `gorm:"column:experience_min" json:"experience_min,omitempty"`
	ExperienceMax *int16 `gorm:"column:experience_max" json:"experience_max,omitempty"`
	TotalHires    int16  `gorm:"column:total_hires;default:1" json:"total_hires"`

	// Company Address (Optional - for work location)
	CompanyAddressID *int64 `gorm:"column:company_address_id" json:"company_address_id,omitempty"`

	Status            string     `gorm:"column:status;type:varchar(20);default:'draft';index" json:"status" validate:"omitempty,oneof='draft' 'pending_review' 'published' 'closed' 'expired' 'suspended' 'rejected'"`
	ViewsCount        int64      `gorm:"column:views_count;default:0" json:"views_count"`
	ApplicationsCount int64      `gorm:"column:applications_count;default:0" json:"applications_count"`
	PublishedAt       *time.Time `gorm:"column:published_at" json:"published_at,omitempty"`
	ExpiredAt         *time.Time `gorm:"column:expired_at" json:"expired_at,omitempty"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relationships
	Category        *JobCategory     `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:SET NULL" json:"category,omitempty"`
	Locations       []JobLocation    `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE" json:"locations,omitempty"`
	Benefits        []JobBenefit     `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE" json:"benefits,omitempty"`
	Skills          []JobSkill       `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE" json:"skills,omitempty"`
	JobRequirements []JobRequirement `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE" json:"job_requirements,omitempty"`

	// Master Data Relations - Job Master Data (New Relations)
	JobTitle         *master.JobTitle         `gorm:"foreignKey:JobTitleID;references:ID;constraint:OnDelete:SET NULL" json:"job_title,omitempty"`
	JobType          *master.JobType          `gorm:"foreignKey:JobTypeID;references:ID;constraint:OnDelete:SET NULL" json:"job_type,omitempty"`
	WorkPolicy       *master.WorkPolicy       `gorm:"foreignKey:WorkPolicyID;references:ID;constraint:OnDelete:SET NULL" json:"work_policy,omitempty"`
	EducationLevelM  *master.EducationLevel   `gorm:"foreignKey:EducationLevelID;references:ID;constraint:OnDelete:SET NULL" json:"education_level_m,omitempty"`
	ExperienceLevelM *master.ExperienceLevel  `gorm:"foreignKey:ExperienceLevelID;references:ID;constraint:OnDelete:SET NULL" json:"experience_level_m,omitempty"`
	GenderPreference *master.GenderPreference `gorm:"foreignKey:GenderPreferenceID;references:ID;constraint:OnDelete:SET NULL" json:"gender_preference,omitempty"`
}

// TableName specifies the table name for Job
func (Job) TableName() string {
	return "jobs"
}

// BeforeCreate hook for Job
func (j *Job) BeforeCreate(tx *gorm.DB) error {
	if j.UUID == uuid.Nil {
		j.UUID = uuid.New()
	}
	return nil
}

// IsPublished checks if job is published
func (j *Job) IsPublished() bool {
	return j.Status == "published"
}

// IsPendingReview checks if job is pending review
func (j *Job) IsPendingReview() bool {
	return j.Status == "pending_review"
}

// IsDraft checks if job is draft
func (j *Job) IsDraft() bool {
	return j.Status == "draft"
}

// IsRejected checks if job is rejected
func (j *Job) IsRejected() bool {
	return j.Status == "rejected"
}

// IsClosed checks if job is closed
func (j *Job) IsClosed() bool {
	return j.Status == "closed" || j.Status == "expired"
}

// IsExpired checks if job has expired
func (j *Job) IsExpired() bool {
	if j.ExpiredAt == nil {
		return false
	}
	return time.Now().After(*j.ExpiredAt)
}

// IsActive checks if job is active and accepting applications
func (j *Job) IsActive() bool {
	return j.Status == "published" && !j.IsExpired()
}

// CanApply checks if job accepts applications
func (j *Job) CanApply() bool {
	return j.IsActive()
}

// ==========================================
// MASTER DATA HELPER METHODS
// ==========================================

// HasMasterDataRelations checks if job has master data relations loaded
func (j *Job) HasMasterDataRelations() bool {
	return j.JobTitle != nil || j.JobType != nil || j.WorkPolicy != nil ||
		j.EducationLevelM != nil || j.ExperienceLevelM != nil || j.GenderPreference != nil
}

// HasJobMasterDataRelations checks if job has job master data relations loaded
func (j *Job) HasJobMasterDataRelations() bool {
	return j.JobTitle != nil || j.JobType != nil || j.WorkPolicy != nil ||
		j.EducationLevelM != nil || j.ExperienceLevelM != nil || j.GenderPreference != nil
}

// ==========================================
// JOB MASTER DATA HELPERS (NEW)
// ==========================================

// GetJobTitle returns job title master data if loaded
func (j *Job) GetJobTitle() *master.JobTitle {
	return j.JobTitle
}

// GetJobTitleName returns job title name from master data
func (j *Job) GetJobTitleName() string {
	if j.JobTitle != nil {
		return j.JobTitle.Name
	}
	return ""
}

// GetJobType returns job type master data if loaded
func (j *Job) GetJobType() *master.JobType {
	return j.JobType
}

// GetJobTypeName returns job type name from master data
func (j *Job) GetJobTypeName() string {
	if j.JobType != nil {
		return j.JobType.Name
	}
	return ""
}

// GetWorkPolicy returns work policy master data if loaded
func (j *Job) GetWorkPolicy() *master.WorkPolicy {
	return j.WorkPolicy
}

// GetWorkPolicyName returns work policy name from master data
func (j *Job) GetWorkPolicyName() string {
	if j.WorkPolicy != nil {
		return j.WorkPolicy.Name
	}
	return ""
}

// GetEducationLevel returns education level master data if loaded
func (j *Job) GetEducationLevel() *master.EducationLevel {
	return j.EducationLevelM
}

// GetEducationLevelName returns education level name from master data
func (j *Job) GetEducationLevelName() string {
	if j.EducationLevelM != nil {
		return j.EducationLevelM.Name
	}
	return ""
}

// GetExperienceLevel returns experience level master data if loaded
func (j *Job) GetExperienceLevel() *master.ExperienceLevel {
	return j.ExperienceLevelM
}

// GetExperienceLevelName returns experience level name from master data
func (j *Job) GetExperienceLevelName() string {
	if j.ExperienceLevelM != nil {
		return j.ExperienceLevelM.Name
	}
	return ""
}

// GetGenderPreference returns gender preference master data if loaded
func (j *Job) GetGenderPreference() *master.GenderPreference {
	return j.GenderPreference
}

// GetGenderPreferenceName returns gender preference name from master data or returns default
func (j *Job) GetGenderPreferenceName() string {
	if j.GenderPreference != nil {
		return j.GenderPreference.Name
	}
	// Default to "Any" if not specified
	return "Any"
}

// JobCategory represents job category entity
type JobCategory struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ParentID    *int64    `gorm:"column:parent_id;index" json:"parent_id,omitempty"`
	Code        string    `gorm:"column:code;type:varchar(30);not null;unique" json:"code" validate:"required,max=30"`
	Name        string    `gorm:"column:name;type:varchar(150);not null;unique" json:"name" validate:"required,max=150"`
	Description string    `gorm:"column:description;type:text" json:"description,omitempty"`
	IsActive    bool      `gorm:"column:is_active;default:true;index" json:"is_active"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relationships
	Parent        *JobCategory     `gorm:"foreignKey:ParentID;references:ID;constraint:OnDelete:SET NULL" json:"parent,omitempty"`
	Children      []JobCategory    `gorm:"foreignKey:ParentID;references:ID" json:"children,omitempty"`
	Subcategories []JobSubcategory `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:CASCADE" json:"subcategories,omitempty"`
	Jobs          []Job            `gorm:"foreignKey:CategoryID;references:ID" json:"jobs,omitempty"`
}

// TableName specifies the table name for JobCategory
func (JobCategory) TableName() string {
	return "job_categories"
}

// JobSubcategory represents job subcategory entity
type JobSubcategory struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CategoryID  int64     `gorm:"column:category_id;not null;index" json:"category_id" validate:"required"`
	Code        string    `gorm:"column:code;type:varchar(50);not null;unique" json:"code" validate:"required,max=50"`
	Name        string    `gorm:"column:name;type:varchar(150);not null;unique" json:"name" validate:"required,max=150"`
	Description string    `gorm:"column:description;type:text" json:"description,omitempty"`
	IsActive    bool      `gorm:"column:is_active;default:true;index" json:"is_active"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relationships
	Category *JobCategory `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:CASCADE" json:"category,omitempty"`
}

// TableName specifies the table name for JobSubcategory
func (JobSubcategory) TableName() string {
	return "job_subcategories"
}

// JobLocation represents job location entity
type JobLocation struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	JobID         int64     `gorm:"column:job_id;not null;index" json:"job_id" validate:"required"`
	CompanyID     *int64    `gorm:"column:company_id;index" json:"company_id,omitempty"`
	LocationType  string    `gorm:"column:location_type;type:varchar(20);default:'onsite'" json:"location_type" validate:"omitempty,oneof='onsite' 'hybrid' 'remote'"`
	Address       string    `gorm:"column:address;type:text" json:"address,omitempty"`
	City          string    `gorm:"column:city;type:varchar(100)" json:"city,omitempty"`
	Province      string    `gorm:"column:province;type:varchar(100)" json:"province,omitempty"`
	PostalCode    string    `gorm:"column:postal_code;type:varchar(20)" json:"postal_code,omitempty"`
	Country       string    `gorm:"column:country;type:varchar(100);default:'Indonesia'" json:"country"`
	Latitude      *float64  `gorm:"column:latitude;type:numeric(10,6)" json:"latitude,omitempty"`
	Longitude     *float64  `gorm:"column:longitude;type:numeric(10,6)" json:"longitude,omitempty"`
	GooglePlaceID string    `gorm:"column:google_place_id;type:varchar(100)" json:"google_place_id,omitempty"`
	MapURL        string    `gorm:"column:map_url;type:text" json:"map_url,omitempty"`
	IsPrimary     bool      `gorm:"column:is_primary;default:false" json:"is_primary"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relationships
	Job *Job `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE" json:"job,omitempty"`
}

// TableName specifies the table name for JobLocation
func (JobLocation) TableName() string {
	return "job_locations"
}

// IsRemote checks if location is remote
func (jl *JobLocation) IsRemote() bool {
	return jl.LocationType == "remote"
}

// IsHybrid checks if location is hybrid
func (jl *JobLocation) IsHybrid() bool {
	return jl.LocationType == "hybrid"
}

// JobBenefit represents job benefit entity
type JobBenefit struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	JobID       int64     `gorm:"column:job_id;not null;index:idx_job_benefit_unique,unique" json:"job_id" validate:"required"`
	BenefitID   *int64    `gorm:"column:benefit_id;index" json:"benefit_id,omitempty"`
	BenefitName string    `gorm:"column:benefit_name;type:varchar(150);not null;index:idx_job_benefit_unique,unique" json:"benefit_name" validate:"required,max=150"`
	Description string    `gorm:"column:description;type:text" json:"description,omitempty"`
	IsHighlight bool      `gorm:"column:is_highlight;default:false;index" json:"is_highlight"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relationships
	Job *Job `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE" json:"job,omitempty"`
}

// TableName specifies the table name for JobBenefit
func (JobBenefit) TableName() string {
	return "job_benefits"
}

// JobSkill represents job skill entity (many-to-many relationship between jobs and skills)
type JobSkill struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	JobID           int64     `gorm:"column:job_id;not null;index:idx_job_skill_unique,unique" json:"job_id" validate:"required"`
	SkillID         int64     `gorm:"column:skill_id;not null;index:idx_job_skill_unique,unique" json:"skill_id" validate:"required"`
	ImportanceLevel string    `gorm:"column:importance_level;type:varchar(20);default:'required';index" json:"importance_level" validate:"omitempty,oneof='required' 'preferred' 'optional'"`
	Weight          float64   `gorm:"column:weight;type:numeric(3,2);default:1.00" json:"weight" validate:"omitempty,min=0,max=1"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relationships
	Job   *Job                 `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE" json:"job,omitempty"`
	Skill *master.SkillsMaster `gorm:"foreignKey:SkillID;references:ID;constraint:OnDelete:CASCADE" json:"skill,omitempty"`
}

// TableName specifies the table name for JobSkill
func (JobSkill) TableName() string {
	return "job_skills"
}

// IsRequired checks if skill is required
func (js *JobSkill) IsRequired() bool {
	return js.ImportanceLevel == "required"
}

// IsPreferred checks if skill is preferred
func (js *JobSkill) IsPreferred() bool {
	return js.ImportanceLevel == "preferred"
}

// JobRequirement represents job requirement entity
type JobRequirement struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	JobID           int64     `gorm:"column:job_id;not null;index" json:"job_id" validate:"required"`
	RequirementType string    `gorm:"column:requirement_type;type:varchar(50);default:'other'" json:"requirement_type" validate:"omitempty,oneof='education' 'experience' 'skill' 'language' 'certification' 'other'"`
	RequirementText string    `gorm:"column:requirement_text;type:text;not null" json:"requirement_text" validate:"required"`
	SkillID         *int64    `gorm:"column:skill_id;index" json:"skill_id,omitempty"`
	MinExperience   *int16    `gorm:"column:min_experience" json:"min_experience,omitempty"`
	MaxExperience   *int16    `gorm:"column:max_experience" json:"max_experience,omitempty"`
	EducationLevel  string    `gorm:"column:education_level;type:varchar(50)" json:"education_level,omitempty"`
	Language        string    `gorm:"column:language;type:varchar(50)" json:"language,omitempty"`
	IsMandatory     bool      `gorm:"column:is_mandatory;default:true" json:"is_mandatory"`
	Priority        int16     `gorm:"column:priority;default:1" json:"priority"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relationships
	Job *Job `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE" json:"job,omitempty"`
}

// TableName specifies the table name for JobRequirement
func (JobRequirement) TableName() string {
	return "job_requirements"
}

// IsEducationRequirement checks if requirement is education type
func (jr *JobRequirement) IsEducationRequirement() bool {
	return jr.RequirementType == "education"
}

// IsExperienceRequirement checks if requirement is experience type
func (jr *JobRequirement) IsExperienceRequirement() bool {
	return jr.RequirementType == "experience"
}

// IsSkillRequirement checks if requirement is skill type
func (jr *JobRequirement) IsSkillRequirement() bool {
	return jr.RequirementType == "skill"
}
