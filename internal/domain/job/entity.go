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
	JobLevel         string    `gorm:"column:job_level;type:varchar(50)" json:"job_level,omitempty" validate:"omitempty,oneof='Internship' 'Entry Level' 'Mid Level' 'Senior Level' 'Manager' 'Director'"`
	EmploymentType   string    `gorm:"column:employment_type;type:varchar(30)" json:"employment_type,omitempty" validate:"omitempty,oneof='Full-Time' 'Part-Time' 'Contract' 'Internship' 'Freelance'"`
	Description      string    `gorm:"column:description;type:text;not null" json:"description" validate:"required"`
	RequirementsText string    `gorm:"column:requirements;type:text" json:"requirements_text,omitempty"`
	Responsibilities string    `gorm:"column:responsibilities;type:text" json:"responsibilities,omitempty"`

	// Master Data Relations
	IndustryID *int64 `gorm:"column:industry_id;index" json:"industry_id,omitempty"`
	DistrictID *int64 `gorm:"column:district_id;index" json:"district_id,omitempty"`
	CityID     *int64 `gorm:"column:city_id;index" json:"city_id,omitempty"`
	ProvinceID *int64 `gorm:"column:province_id;index" json:"province_id,omitempty"`

	// Legacy Location Fields (for backward compatibility)
	Location     string `gorm:"column:location;type:varchar(150)" json:"location,omitempty"`
	City         string `gorm:"column:city;type:varchar(100)" json:"city,omitempty"`
	Province     string `gorm:"column:province;type:varchar(100)" json:"province,omitempty"`
	RemoteOption bool   `gorm:"column:remote_option;default:false" json:"remote_option"`

	SalaryMin         *float64   `gorm:"column:salary_min;type:numeric(12,2)" json:"salary_min,omitempty"`
	SalaryMax         *float64   `gorm:"column:salary_max;type:numeric(12,2)" json:"salary_max,omitempty"`
	Currency          string     `gorm:"column:currency;type:varchar(10);default:'IDR'" json:"currency" validate:"omitempty,len=3"`
	ExperienceMin     *int16     `gorm:"column:experience_min" json:"experience_min,omitempty"`
	ExperienceMax     *int16     `gorm:"column:experience_max" json:"experience_max,omitempty"`
	EducationLevel    string     `gorm:"column:education_level;type:varchar(50)" json:"education_level,omitempty"`
	TotalHires        int16      `gorm:"column:total_hires;default:1" json:"total_hires"`
	Status            string     `gorm:"column:status;type:varchar(20);default:'draft';index" json:"status" validate:"omitempty,oneof='draft' 'published' 'closed' 'expired' 'suspended'"`
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

	// Master Data Relations
	Industry  *master.Industry `gorm:"foreignKey:IndustryID;references:ID;constraint:OnDelete:SET NULL" json:"industry,omitempty"`
	District  *master.District `gorm:"foreignKey:DistrictID;references:ID;constraint:OnDelete:SET NULL" json:"district,omitempty"`
	MCity     *master.City     `gorm:"foreignKey:CityID;references:ID;constraint:OnDelete:SET NULL" json:"m_city,omitempty"`
	MProvince *master.Province `gorm:"foreignKey:ProvinceID;references:ID;constraint:OnDelete:SET NULL" json:"m_province,omitempty"`
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
	return j.Industry != nil || j.District != nil || j.MCity != nil || j.MProvince != nil
}

// GetIndustry returns the industry relation if loaded
func (j *Job) GetIndustry() *master.Industry {
	return j.Industry
}

// GetIndustryName returns industry name from master data or falls back to legacy field
// This provides smart fallback for backward compatibility
func (j *Job) GetIndustryName() string {
	if j.Industry != nil {
		return j.Industry.Name
	}
	// No legacy industry field in Job, return empty
	return ""
}

// GetDistrict returns the district relation if loaded
func (j *Job) GetDistrict() *master.District {
	return j.District
}

// GetCity returns the city relation if loaded
func (j *Job) GetCity() *master.City {
	return j.MCity
}

// GetProvince returns the province relation if loaded
func (j *Job) GetProvince() *master.Province {
	return j.MProvince
}

// GetCityName returns city name from master data or falls back to legacy field
func (j *Job) GetCityName() string {
	if j.MCity != nil {
		return j.MCity.Name
	}
	// Fallback to legacy field
	return j.City
}

// GetProvinceName returns province name from master data or falls back to legacy field
func (j *Job) GetProvinceName() string {
	if j.MProvince != nil {
		return j.MProvince.Name
	}
	// Fallback to legacy field
	return j.Province
}

// GetFullLocation returns full location string with master data or legacy fields
// Returns format: "District, City, Province" or legacy location string
func (j *Job) GetFullLocation() string {
	// Try master data first
	if j.District != nil && j.MCity != nil && j.MProvince != nil {
		districtName := j.District.Name
		cityName := j.MCity.GetFullName() // e.g., "Kota Bandung"
		provinceName := j.MProvince.Name
		return districtName + ", " + cityName + ", " + provinceName
	}

	// Fallback to legacy fields
	if j.City != "" && j.Province != "" {
		return j.City + ", " + j.Province
	}

	// Last resort: use Location field
	return j.Location
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
	Job *Job `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE" json:"job,omitempty"`
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
