package master

import (
	"database/sql"
	"time"
)

// JobTitle represents master data for job titles with smart recommendations
// Maps to: job_titles table
type JobTitle struct {
	ID                    int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                  string    `gorm:"type:varchar(200);not null;uniqueIndex" json:"name" validate:"required,min=2,max=200"`
	NormalizedName        string    `gorm:"type:varchar(200);index:idx_job_titles_normalized" json:"normalized_name"`
	RecommendedCategoryID *int64    `gorm:"index" json:"recommended_category_id,omitempty"`
	PopularityScore       float64   `gorm:"type:numeric(5,2);default:0.00;index:idx_job_titles_popularity,sort:desc" json:"popularity_score" validate:"min=0,max=100"`
	SearchCount           int64     `gorm:"default:0" json:"search_count"`
	IsActive              bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt             time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for JobTitle
func (JobTitle) TableName() string {
	return "job_titles"
}

// IncrementSearchCount increments search count for analytics
func (jt *JobTitle) IncrementSearchCount() {
	jt.SearchCount++
}

// IncrementPopularity increases popularity score
func (jt *JobTitle) IncrementPopularity(amount float64) {
	jt.PopularityScore += amount
	if jt.PopularityScore > 100.0 {
		jt.PopularityScore = 100.0
	}
}

// JobType represents job type options (full-time, part-time, etc)
type JobType struct {
	ID    int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Code  string `gorm:"type:varchar(30);not null;uniqueIndex" json:"code" validate:"required"`
	Name  string `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	Order int    `gorm:"default:0" json:"order"`
}

// TableName specifies the table name for JobType
func (JobType) TableName() string {
	return "job_types"
}

// WorkPolicy represents work location policies (onsite, remote, hybrid)
type WorkPolicy struct {
	ID    int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Code  string `gorm:"type:varchar(30);not null;uniqueIndex" json:"code" validate:"required"`
	Name  string `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	Order int    `gorm:"default:0" json:"order"`
}

// TableName specifies the table name for WorkPolicy
func (WorkPolicy) TableName() string {
	return "work_policies"
}

// EducationLevel represents education level requirements
type EducationLevel struct {
	ID    int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Code  string `gorm:"type:varchar(30);not null;uniqueIndex" json:"code" validate:"required"`
	Name  string `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	Order int    `gorm:"default:0" json:"order"`
}

// TableName specifies the table name for EducationLevel
func (EducationLevel) TableName() string {
	return "education_levels"
}

// ExperienceLevel represents experience level requirements
type ExperienceLevel struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Code     string `gorm:"type:varchar(30);not null;uniqueIndex" json:"code" validate:"required"`
	Name     string `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	MinYears int    `gorm:"default:0" json:"min_years"`
	MaxYears *int   `gorm:"" json:"max_years,omitempty"`
	Order    int    `gorm:"default:0" json:"order"`
}

// TableName specifies the table name for ExperienceLevel
func (ExperienceLevel) TableName() string {
	return "experience_levels"
}

// GenderPreference represents gender preference options
type GenderPreference struct {
	ID    int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Code  string `gorm:"type:varchar(30);not null;uniqueIndex" json:"code" validate:"required"`
	Name  string `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	Order int    `gorm:"default:0" json:"order"`
}

// TableName specifies the table name for GenderPreference
func (GenderPreference) TableName() string {
	return "gender_preferences"
}

// JobOptionsResponse represents combined response for GetJobOptions endpoint
type JobOptionsResponse struct {
	JobTypes          []JobType          `json:"job_types"`
	WorkPolicies      []WorkPolicy       `json:"work_policies"`
	EducationLevels   []EducationLevel   `json:"education_levels"`
	ExperienceLevels  []ExperienceLevel  `json:"experience_levels"`
	GenderPreferences []GenderPreference `json:"genders"`
}

// BenefitsMaster represents a master data entry for job benefits
// Maps to: benefits_master table
type BenefitsMaster struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Code            string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"code" validate:"required,min=2,max=50"`
	Name            string    `gorm:"type:varchar(150);not null;uniqueIndex" json:"name" validate:"required,min=2,max=150"`
	Category        string    `gorm:"type:varchar(50);default:'other'" json:"category" validate:"oneof=financial health career lifestyle flexibility other"`
	Description     string    `gorm:"type:text" json:"description,omitempty"`
	Icon            string    `gorm:"type:varchar(100)" json:"icon,omitempty"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	PopularityScore float64   `gorm:"type:numeric(5,2);default:0.00;index:idx_benefits_master_popularity,sort:desc" json:"popularity_score" validate:"min=0,max=100"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for BenefitsMaster
func (BenefitsMaster) TableName() string {
	return "benefits_master"
}

// IsFinancial checks if the benefit is in the financial category
func (b *BenefitsMaster) IsFinancial() bool {
	return b.Category == "financial"
}

// IsHealth checks if the benefit is in the health category
func (b *BenefitsMaster) IsHealth() bool {
	return b.Category == "health"
}

// IsCareer checks if the benefit is in the career category
func (b *BenefitsMaster) IsCareer() bool {
	return b.Category == "career"
}

// IsLifestyle checks if the benefit is in the lifestyle category
func (b *BenefitsMaster) IsLifestyle() bool {
	return b.Category == "lifestyle"
}

// IsFlexibility checks if the benefit is in the flexibility category
func (b *BenefitsMaster) IsFlexibility() bool {
	return b.Category == "flexibility"
}

// IsPopular checks if the benefit has high popularity score (>= 70)
func (b *BenefitsMaster) IsPopular() bool {
	return b.PopularityScore >= 70.0
}

// IncrementPopularity increases the popularity score
func (b *BenefitsMaster) IncrementPopularity(amount float64) {
	b.PopularityScore += amount
	if b.PopularityScore > 100.0 {
		b.PopularityScore = 100.0
	}
}

// SkillsMaster represents a master data entry for skills
// Maps to: skills_master table
type SkillsMaster struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Code            string    `gorm:"type:varchar(50);uniqueIndex" json:"code,omitempty" validate:"omitempty,min=2,max=50"`
	Name            string    `gorm:"type:varchar(150);not null;uniqueIndex" json:"name" validate:"required,min=2,max=150"`
	NormalizedName  string    `gorm:"type:varchar(150)" json:"normalized_name,omitempty"`
	CategoryID      *int64    `gorm:"index" json:"category_id,omitempty"`
	Description     string    `gorm:"type:text" json:"description,omitempty"`
	SkillType       string    `gorm:"type:varchar(30);default:'technical'" json:"skill_type" validate:"oneof=technical soft language tool"`
	DifficultyLevel string    `gorm:"type:varchar(20);default:'intermediate'" json:"difficulty_level" validate:"oneof=beginner intermediate advanced"`
	PopularityScore float64   `gorm:"type:numeric(5,2);default:0.00;index:idx_skills_master_popularity,sort:desc" json:"popularity_score" validate:"min=0,max=100"`
	Aliases         []string  `gorm:"type:text[]" json:"aliases,omitempty"`
	ParentID        *int64    `gorm:"index" json:"parent_id,omitempty"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	// Note: Category references job_categories table (not included here to avoid circular dependency)
	// In implementation, use: Category *job.JobCategory `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL"`
	Parent   *SkillsMaster  `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL" json:"parent,omitempty"`
	Children []SkillsMaster `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL" json:"children,omitempty"`
}

// TableName specifies the table name for SkillsMaster
func (SkillsMaster) TableName() string {
	return "skills_master"
}

// IsTechnical checks if the skill is a technical skill
func (s *SkillsMaster) IsTechnical() bool {
	return s.SkillType == "technical"
}

// IsSoft checks if the skill is a soft skill
func (s *SkillsMaster) IsSoft() bool {
	return s.SkillType == "soft"
}

// IsLanguage checks if the skill is a language skill
func (s *SkillsMaster) IsLanguage() bool {
	return s.SkillType == "language"
}

// IsTool checks if the skill is a tool/software skill
func (s *SkillsMaster) IsTool() bool {
	return s.SkillType == "tool"
}

// IsBeginner checks if the skill has beginner difficulty level
func (s *SkillsMaster) IsBeginner() bool {
	return s.DifficultyLevel == "beginner"
}

// IsIntermediate checks if the skill has intermediate difficulty level
func (s *SkillsMaster) IsIntermediate() bool {
	return s.DifficultyLevel == "intermediate"
}

// IsAdvanced checks if the skill has advanced difficulty level
func (s *SkillsMaster) IsAdvanced() bool {
	return s.DifficultyLevel == "advanced"
}

// IsPopular checks if the skill has high popularity score (>= 70)
func (s *SkillsMaster) IsPopular() bool {
	return s.PopularityScore >= 70.0
}

// HasParent checks if the skill has a parent skill
func (s *SkillsMaster) HasParent() bool {
	return s.ParentID != nil
}

// HasChildren checks if the skill has child skills
func (s *SkillsMaster) HasChildren() bool {
	return len(s.Children) > 0
}

// IncrementPopularity increases the popularity score
func (s *SkillsMaster) IncrementPopularity(amount float64) {
	s.PopularityScore += amount
	if s.PopularityScore > 100.0 {
		s.PopularityScore = 100.0
	}
}

// AddAlias adds a new alias to the skill
func (s *SkillsMaster) AddAlias(alias string) {
	for _, existing := range s.Aliases {
		if existing == alias {
			return
		}
	}
	s.Aliases = append(s.Aliases, alias)
}

// RemoveAlias removes an alias from the skill
func (s *SkillsMaster) RemoveAlias(alias string) {
	for i, existing := range s.Aliases {
		if existing == alias {
			s.Aliases = append(s.Aliases[:i], s.Aliases[i+1:]...)
			return
		}
	}
}

// HasAlias checks if the skill has a specific alias
func (s *SkillsMaster) HasAlias(alias string) bool {
	for _, existing := range s.Aliases {
		if existing == alias {
			return true
		}
	}
	return false
}

// ========================================
// Company Refactor Master Data Entities
// ========================================

// Industry represents the master data for company industries
// This table stores predefined industry categories that companies can select
type Industry struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string         `gorm:"type:varchar(100);not null;uniqueIndex:idx_industries_name" json:"name" validate:"required,min=2,max=100"`
	Slug         string         `gorm:"type:varchar(100);not null;uniqueIndex:idx_industries_slug" json:"slug" validate:"required,slug,min=2,max=100"`
	Description  sql.NullString `gorm:"type:text" json:"description,omitempty"`
	IconURL      sql.NullString `gorm:"type:text" json:"icon_url,omitempty" validate:"omitempty,url"`
	IsActive     bool           `gorm:"default:true;not null;index:idx_industries_active_deleted" json:"is_active"`
	DisplayOrder int            `gorm:"default:0;not null;index:idx_industries_display_order" json:"display_order"`
	CreatedAt    time.Time      `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime;not null" json:"updated_at"`
	DeletedAt    sql.NullTime   `gorm:"index:idx_industries_active_deleted" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Industry entity
func (Industry) TableName() string {
	return "industries"
}

// IsDeleted checks if the industry is soft deleted
func (i *Industry) IsDeleted() bool {
	return i.DeletedAt.Valid
}

// GetDescription returns the description if valid, otherwise empty string
func (i *Industry) GetDescription() string {
	if i.Description.Valid {
		return i.Description.String
	}
	return ""
}

// GetIconURL returns the icon URL if valid, otherwise empty string
func (i *Industry) GetIconURL() string {
	if i.IconURL.Valid {
		return i.IconURL.String
	}
	return ""
}

// CompanySize represents the master data for company size categories
// This table stores predefined company size ranges based on employee count
type CompanySize struct {
	ID           int64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Label        string        `gorm:"type:varchar(50);not null;uniqueIndex:idx_company_sizes_label" json:"label" validate:"required,min=2,max=50"`
	MinEmployees int           `gorm:"not null" json:"min_employees" validate:"required,gte=1"`
	MaxEmployees sql.NullInt32 `gorm:"" json:"max_employees,omitempty" validate:"omitempty,gtefield=MinEmployees"`
	IsActive     bool          `gorm:"default:true;not null;index:idx_company_sizes_active" json:"is_active"`
	DisplayOrder int           `gorm:"default:0;not null;index:idx_company_sizes_display_order" json:"display_order"`
	CreatedAt    time.Time     `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt    time.Time     `gorm:"autoUpdateTime;not null" json:"updated_at"`
}

// TableName specifies the table name for CompanySize entity
func (CompanySize) TableName() string {
	return "company_sizes"
}

// GetMaxEmployees returns the max employees if valid, otherwise -1 (unlimited)
func (cs *CompanySize) GetMaxEmployees() int {
	if cs.MaxEmployees.Valid {
		return int(cs.MaxEmployees.Int32)
	}
	return -1 // Indicates unlimited
}

// IsUnlimited checks if this size category has no upper limit
func (cs *CompanySize) IsUnlimited() bool {
	return !cs.MaxEmployees.Valid
}

// GetRange returns the employee range as a formatted string
func (cs *CompanySize) GetRange() string {
	return cs.Label
}

// Province represents the first-level location (Provinsi)
// This table stores all Indonesian provinces with official BPS codes
type Province struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;index:idx_provinces_name" json:"name" validate:"required,min=2,max=100"`
	Code      string    `gorm:"type:varchar(10);not null;uniqueIndex:idx_provinces_code" json:"code" validate:"required,min=2,max=10"`
	IsActive  bool      `gorm:"default:true;not null;index:idx_provinces_active" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null" json:"updated_at"`

	// Relations
	Cities []City `gorm:"foreignKey:ProvinceID;constraint:OnDelete:RESTRICT" json:"cities,omitempty"`
}

// TableName specifies the table name for Province entity
func (Province) TableName() string {
	return "provinces"
}

// City represents the second-level location (Kota/Kabupaten)
// This table stores cities and regencies with their parent province
type City struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProvinceID int64     `gorm:"not null;index:idx_cities_province_active" json:"province_id" validate:"required,gt=0"`
	Name       string    `gorm:"type:varchar(100);not null;index:idx_cities_name" json:"name" validate:"required,min=2,max=100"`
	Type       string    `gorm:"type:varchar(20);not null;index:idx_cities_type" json:"type" validate:"required,oneof=Kota Kabupaten"`
	Code       string    `gorm:"type:varchar(10);not null;uniqueIndex:idx_cities_code" json:"code" validate:"required,min=2,max=10"`
	IsActive   bool      `gorm:"default:true;not null;index:idx_cities_province_active" json:"is_active"`
	CreatedAt  time.Time `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime;not null" json:"updated_at"`

	// Relations
	Province  *Province  `gorm:"foreignKey:ProvinceID;constraint:OnDelete:RESTRICT" json:"province,omitempty"`
	Districts []District `gorm:"foreignKey:CityID;constraint:OnDelete:RESTRICT" json:"districts,omitempty"`
}

// TableName specifies the table name for City entity
func (City) TableName() string {
	return "cities"
}

// GetFullName returns the city name with type prefix (e.g., "Kota Bandung", "Kabupaten Bandung Barat")
func (c *City) GetFullName() string {
	return c.Type + " " + c.Name
}

// IsKota checks if this is a city (Kota) as opposed to regency (Kabupaten)
func (c *City) IsKota() bool {
	return c.Type == "Kota"
}

// IsKabupaten checks if this is a regency (Kabupaten) as opposed to city (Kota)
func (c *City) IsKabupaten() bool {
	return c.Type == "Kabupaten"
}

// District represents the third-level location (Kecamatan)
// This table stores districts with their parent city and postal codes
type District struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CityID     int64          `gorm:"not null;index:idx_districts_city_active" json:"city_id" validate:"required,gt=0"`
	Name       string         `gorm:"type:varchar(100);not null;index:idx_districts_name" json:"name" validate:"required,min=2,max=100"`
	Code       string         `gorm:"type:varchar(10);not null;uniqueIndex:idx_districts_code" json:"code" validate:"required,min=2,max=10"`
	PostalCode sql.NullString `gorm:"type:varchar(10);index:idx_districts_postal" json:"postal_code,omitempty" validate:"omitempty,len=5"`
	IsActive   bool           `gorm:"default:true;not null;index:idx_districts_city_active" json:"is_active"`
	CreatedAt  time.Time      `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime;not null" json:"updated_at"`

	// Relations
	City *City `gorm:"foreignKey:CityID;constraint:OnDelete:RESTRICT" json:"city,omitempty"`
}

// TableName specifies the table name for District entity
func (District) TableName() string {
	return "districts"
}

// GetPostalCode returns the postal code if valid, otherwise empty string
func (d *District) GetPostalCode() string {
	if d.PostalCode.Valid {
		return d.PostalCode.String
	}
	return ""
}

// GetFullLocationPath returns the complete location hierarchy
// Example: "Batujajar, Kabupaten Bandung Barat, Jawa Barat"
func (d *District) GetFullLocationPath() string {
	if d.City == nil {
		return d.Name
	}

	result := d.Name + ", " + d.City.GetFullName()

	if d.City.Province != nil {
		result += ", " + d.City.Province.Name
	}

	return result
}

// HasPostalCode checks if this district has a postal code
func (d *District) HasPostalCode() bool {
	return d.PostalCode.Valid && d.PostalCode.String != ""
}
