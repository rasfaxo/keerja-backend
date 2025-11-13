package user

import (
	"time"

	"keerja-backend/internal/domain/master"

	"github.com/google/uuid"
)

// User represents the main user entity
type User struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex" json:"uuid"`
	FullName     string     `gorm:"type:varchar(150);not null" json:"full_name" validate:"required,min=3,max=150"`
	Email        string     `gorm:"type:varchar(150);uniqueIndex;not null" json:"email" validate:"required,email,max=150"`
	Phone        *string    `gorm:"type:varchar(20)" json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	PasswordHash string     `gorm:"type:text;not null" json:"-"`
	UserType     string     `gorm:"type:varchar(20);check:user_type IN ('jobseeker','employer','admin')" json:"user_type" validate:"required,oneof=jobseeker employer admin"`
	IsVerified   bool       `gorm:"default:false" json:"is_verified"`
	Status       string     `gorm:"type:varchar(20);default:'active';check:status IN ('active','inactive','suspended')" json:"status"`
	LastLogin    *time.Time `gorm:"type:timestamp" json:"last_login,omitempty"`
	CreatedAt    time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	Profile        *UserProfile        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"profile,omitempty"`
	Preference     *UserPreference     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"preference,omitempty"`
	Educations     []UserEducation     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"educations,omitempty"`
	Experiences    []UserExperience    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"experiences,omitempty"`
	Skills         []UserSkill         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"skills,omitempty"`
	Certifications []UserCertification `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"certifications,omitempty"`
	Languages      []UserLanguage      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"languages,omitempty"`
	Projects       []UserProject       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"projects,omitempty"`
	Documents      []UserDocument      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"documents,omitempty"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// IsActive checks if user is active
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// IsJobseeker checks if user is a jobseeker
func (u *User) IsJobseeker() bool {
	return u.UserType == "jobseeker"
}

// IsEmployer checks if user is an employer
func (u *User) IsEmployer() bool {
	return u.UserType == "employer"
}

// UserProfile represents user profile information
type UserProfile struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64      `gorm:"not null;uniqueIndex" json:"user_id"`
	Headline    *string    `gorm:"type:varchar(150)" json:"headline,omitempty" validate:"omitempty,max=150"`
	Bio         *string    `gorm:"type:text" json:"bio,omitempty"`
	Gender      *string    `gorm:"type:varchar(10);check:gender IN ('male','female','other')" json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	BirthDate   *time.Time `gorm:"type:date;column:birth_date" json:"birth_date,omitempty"`
	Nationality *string    `gorm:"type:varchar(100);column:nationality" json:"nationality,omitempty"`
	Address     *string    `gorm:"type:text;column:address" json:"address,omitempty"`

	// Master Data Location Fields
	DistrictID *int64 `gorm:"column:district_id;index" json:"district_id,omitempty"`
	CityID     *int64 `gorm:"column:city_id;index" json:"city_id,omitempty"`
	ProvinceID *int64 `gorm:"column:province_id;index" json:"province_id,omitempty"`

	// Legacy Location Fields (backward compatibility)
	LocationCity    *string `gorm:"type:varchar(100);column:location_city" json:"location_city,omitempty"`
	LocationState   *string `gorm:"type:varchar(100);column:location_state" json:"location_state,omitempty"`
	LocationCountry *string `gorm:"type:varchar(100);column:location_country" json:"location_country,omitempty"`

	PostalCode       *string  `gorm:"type:varchar(10);column:postal_code" json:"postal_code,omitempty"`
	LinkedInURL      *string  `gorm:"type:varchar(255);column:linkedin_url" json:"linkedin_url,omitempty"`
	PortfolioURL     *string  `gorm:"type:varchar(255);column:portfolio_url" json:"portfolio_url,omitempty"`
	GithubURL        *string  `gorm:"type:varchar(255);column:github_url" json:"github_url,omitempty"`
	DesiredPosition  *string  `gorm:"type:varchar(150);column:desired_position" json:"desired_position,omitempty"`
	DesiredSalaryMin *float64 `gorm:"type:numeric(12,2)" json:"desired_salary_min,omitempty"`
	DesiredSalaryMax *float64 `gorm:"type:numeric(12,2)" json:"desired_salary_max,omitempty"`
	ExperienceLevel  *string  `gorm:"type:varchar(50);check:experience_level IN ('internship','junior','mid','senior','lead')" json:"experience_level,omitempty" validate:"omitempty,oneof=internship junior mid senior lead"`

	// Master Data Industry Fields - Support multiple industries
	IndustryIDs []int64 `gorm:"-" json:"industry_ids,omitempty"` // Virtual field, stored in junction table

	// Legacy Industry Field (backward compatibility)
	IndustryInterest *string `gorm:"type:varchar(100)" json:"industry_interest,omitempty"`

	AvailabilityStatus string    `gorm:"type:varchar(50);default:'open'" json:"availability_status" validate:"omitempty,oneof=open looking_actively not_looking"`
	ProfileVisibility  bool      `gorm:"default:true" json:"profile_visibility"`
	Slug               *string   `gorm:"type:varchar(100);uniqueIndex" json:"slug,omitempty"`
	AvatarURL          *string   `gorm:"type:text" json:"avatar_url,omitempty"`
	CoverURL           *string   `gorm:"type:text" json:"cover_url,omitempty"`
	CreatedAt          time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt          time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"-"`

	// Master Data Relations
	District   *master.District  `gorm:"foreignKey:DistrictID;references:ID;constraint:OnDelete:SET NULL" json:"district,omitempty"`
	MCity      *master.City      `gorm:"foreignKey:CityID;references:ID;constraint:OnDelete:SET NULL" json:"m_city,omitempty"`
	MProvince  *master.Province  `gorm:"foreignKey:ProvinceID;references:ID;constraint:OnDelete:SET NULL" json:"m_province,omitempty"`
	Industries []master.Industry `gorm:"many2many:user_profile_industries;constraint:OnDelete:CASCADE" json:"industries,omitempty"`
}

// TableName specifies the table name for UserProfile
func (UserProfile) TableName() string {
	return "user_profiles"
}

// UserPreference represents user preferences and settings
type UserPreference struct {
	ID                 int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID             int64  `gorm:"not null;uniqueIndex" json:"user_id"`
	LanguagePreference string `gorm:"type:varchar(10);default:'id'" json:"language_preference"`
	ThemePreference    string `gorm:"type:varchar(10);default:'light'" json:"theme_preference"`
	PreferredJobType   string `gorm:"type:varchar(50);default:'onsite'" json:"preferred_job_type"`

	// Master Data Preferred Industries - Support multiple industries
	PreferredIndustryIDs []int64 `gorm:"-" json:"preferred_industry_ids,omitempty"` // Virtual field, stored in junction table

	// Master Data Preferred Location
	PreferredDistrictID *int64 `gorm:"column:preferred_district_id;index" json:"preferred_district_id,omitempty"`
	PreferredCityID     *int64 `gorm:"column:preferred_city_id;index" json:"preferred_city_id,omitempty"`
	PreferredProvinceID *int64 `gorm:"column:preferred_province_id;index" json:"preferred_province_id,omitempty"`

	// Legacy Fields (backward compatibility)
	PreferredIndustry *string `gorm:"type:varchar(100)" json:"preferred_industry,omitempty"`
	PreferredLocation *string `gorm:"type:varchar(100)" json:"preferred_location,omitempty"`

	PreferredSalaryMin  *float64  `gorm:"type:numeric(12,2)" json:"preferred_salary_min,omitempty"`
	PreferredSalaryMax  *float64  `gorm:"type:numeric(12,2)" json:"preferred_salary_max,omitempty"`
	EmailNotifications  bool      `gorm:"default:true" json:"email_notifications"`
	SMSNotifications    bool      `gorm:"default:false" json:"sms_notifications"`
	PushNotifications   bool      `gorm:"default:true" json:"push_notifications"`
	EmailMarketing      bool      `gorm:"default:false" json:"email_marketing"`
	ProfileVisibility   string    `gorm:"type:varchar(20);default:'public';check:profile_visibility IN ('public','private','recruiter-only')" json:"profile_visibility"`
	ShowOnlineStatus    bool      `gorm:"default:true" json:"show_online_status"`
	AllowDirectMessages bool      `gorm:"default:true" json:"allow_direct_messages"`
	DataSharingConsent  bool      `gorm:"default:true" json:"data_sharing_consent"`
	CreatedAt           time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt           time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"-"`

	// Master Data Relations
	PreferredDistrict   *master.District  `gorm:"foreignKey:PreferredDistrictID;references:ID;constraint:OnDelete:SET NULL" json:"preferred_district,omitempty"`
	PreferredMCity      *master.City      `gorm:"foreignKey:PreferredCityID;references:ID;constraint:OnDelete:SET NULL" json:"preferred_m_city,omitempty"`
	PreferredMProvince  *master.Province  `gorm:"foreignKey:PreferredProvinceID;references:ID;constraint:OnDelete:SET NULL" json:"preferred_m_province,omitempty"`
	PreferredIndustries []master.Industry `gorm:"many2many:user_preference_industries;constraint:OnDelete:CASCADE" json:"preferred_industries,omitempty"`
}

// TableName specifies the table name for UserPreference
func (UserPreference) TableName() string {
	return "user_preferences"
}

// ===========================================
// UserProfile Helper Methods
// ===========================================

// HasMasterDataRelations checks if user profile has master data relations loaded
func (up *UserProfile) HasMasterDataRelations() bool {
	return up.District != nil || up.MCity != nil || up.MProvince != nil || len(up.Industries) > 0
}

// GetDistrict returns the district relation if loaded
func (up *UserProfile) GetDistrict() *master.District {
	return up.District
}

// GetCity returns the city relation if loaded
func (up *UserProfile) GetCity() *master.City {
	return up.MCity
}

// GetProvince returns the province relation if loaded
func (up *UserProfile) GetProvince() *master.Province {
	return up.MProvince
}

// GetCityName returns city name from master data or falls back to legacy field
func (up *UserProfile) GetCityName() string {
	if up.MCity != nil {
		return up.MCity.Name
	}
	if up.LocationCity != nil {
		return *up.LocationCity
	}
	return ""
}

// GetProvinceName returns province name from master data or falls back to legacy field
func (up *UserProfile) GetProvinceName() string {
	if up.MProvince != nil {
		return up.MProvince.Name
	}
	if up.LocationState != nil {
		return *up.LocationState
	}
	return ""
}

// GetFullLocation returns formatted location string with smart fallback
func (up *UserProfile) GetFullLocation() string {
	// Try master data first
	if up.District != nil && up.MCity != nil && up.MProvince != nil {
		districtName := up.District.Name
		cityName := up.MCity.GetFullName()
		provinceName := up.MProvince.Name
		return districtName + ", " + cityName + ", " + provinceName
	}

	// Fallback to legacy fields
	if up.LocationCity != nil && up.LocationState != nil {
		return *up.LocationCity + ", " + *up.LocationState
	}

	if up.LocationCity != nil {
		return *up.LocationCity
	}

	return ""
}

// GetIndustries returns array of industry names
func (up *UserProfile) GetIndustries() []string {
	if len(up.Industries) == 0 {
		// Fallback to legacy field
		if up.IndustryInterest != nil && *up.IndustryInterest != "" {
			return []string{*up.IndustryInterest}
		}
		return []string{}
	}

	industries := make([]string, len(up.Industries))
	for i, ind := range up.Industries {
		industries[i] = ind.Name
	}
	return industries
}

// ===========================================
// UserPreference Helper Methods
// ===========================================

// HasMasterDataRelations checks if user preference has master data relations loaded
func (upref *UserPreference) HasMasterDataRelations() bool {
	return upref.PreferredDistrict != nil || upref.PreferredMCity != nil ||
		upref.PreferredMProvince != nil || len(upref.PreferredIndustries) > 0
}

// GetPreferredCityName returns preferred city name from master data or falls back to legacy field
func (upref *UserPreference) GetPreferredCityName() string {
	if upref.PreferredMCity != nil {
		return upref.PreferredMCity.Name
	}
	if upref.PreferredLocation != nil {
		return *upref.PreferredLocation
	}
	return ""
}

// GetPreferredProvinceName returns preferred province name from master data
func (upref *UserPreference) GetPreferredProvinceName() string {
	if upref.PreferredMProvince != nil {
		return upref.PreferredMProvince.Name
	}
	return ""
}

// GetPreferredFullLocation returns formatted preferred location string
func (upref *UserPreference) GetPreferredFullLocation() string {
	if upref.PreferredDistrict != nil && upref.PreferredMCity != nil && upref.PreferredMProvince != nil {
		districtName := upref.PreferredDistrict.Name
		cityName := upref.PreferredMCity.GetFullName()
		provinceName := upref.PreferredMProvince.Name
		return districtName + ", " + cityName + ", " + provinceName
	}

	// Fallback to legacy field
	if upref.PreferredLocation != nil {
		return *upref.PreferredLocation
	}

	return ""
}

// GetPreferredIndustries returns array of preferred industry names
func (upref *UserPreference) GetPreferredIndustries() []string {
	if len(upref.PreferredIndustries) == 0 {
		// Fallback to legacy field
		if upref.PreferredIndustry != nil && *upref.PreferredIndustry != "" {
			return []string{*upref.PreferredIndustry}
		}
		return []string{}
	}

	industries := make([]string, len(upref.PreferredIndustries))
	for i, ind := range upref.PreferredIndustries {
		industries[i] = ind.Name
	}
	return industries
}

// UserEducation represents user's education history
type UserEducation struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int64     `gorm:"not null;index" json:"user_id"`
	InstitutionName string    `gorm:"type:varchar(150);not null" json:"institution_name" validate:"required,max=150"`
	Major           *string   `gorm:"type:varchar(100)" json:"major,omitempty"`
	DegreeLevel     *string   `gorm:"type:varchar(50);check:degree_level IN ('SMA','D1','D2','D3','S1','S2','S3','Other')" json:"degree_level,omitempty" validate:"omitempty,oneof=SMA D1 D2 D3 S1 S2 S3 Other"`
	StartYear       *int      `gorm:"type:int;check:start_year >= 1950" json:"start_year,omitempty"`
	EndYear         *int      `gorm:"type:int;check:end_year >= 1950" json:"end_year,omitempty"`
	GPA             *float64  `gorm:"type:numeric(3,2)" json:"gpa,omitempty"`
	Activities      *string   `gorm:"type:text" json:"activities,omitempty"`
	Description     *string   `gorm:"type:text" json:"description,omitempty"`
	IsCurrent       bool      `gorm:"default:false" json:"is_current"`
	Verified        bool      `gorm:"default:false" json:"verified"`
	CreatedAt       time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for UserEducation
func (UserEducation) TableName() string {
	return "user_educations"
}

// UserExperience represents user's work experience
type UserExperience struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int64      `gorm:"not null;index" json:"user_id"`
	CompanyName     string     `gorm:"type:varchar(150);not null" json:"company_name" validate:"required,max=150"`
	PositionTitle   string     `gorm:"type:varchar(150);not null" json:"position_title" validate:"required,max=150"`
	Industry        *string    `gorm:"type:varchar(100)" json:"industry,omitempty"`
	EmploymentType  *string    `gorm:"type:varchar(30);check:employment_type IN ('full-time','part-time','internship','freelance','contract')" json:"employment_type,omitempty" validate:"omitempty,oneof=full-time part-time internship freelance contract"`
	StartDate       time.Time  `gorm:"type:date;not null" json:"start_date" validate:"required"`
	EndDate         *time.Time `gorm:"type:date" json:"end_date,omitempty"`
	IsCurrent       bool       `gorm:"default:false" json:"is_current"`
	Description     *string    `gorm:"type:text" json:"description,omitempty"`
	Achievements    *string    `gorm:"type:text" json:"achievements,omitempty"`
	LocationCity    *string    `gorm:"type:varchar(100)" json:"location_city,omitempty"`
	LocationCountry string     `gorm:"type:varchar(100);default:'Indonesia'" json:"location_country"`
	Verified        bool       `gorm:"default:false" json:"verified"`
	CreatedAt       time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for UserExperience
func (UserExperience) TableName() string {
	return "user_experiences"
}

// UserSkill represents user's skills
type UserSkill struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int64      `gorm:"not null;index" json:"user_id"`
	SkillName       string     `gorm:"type:varchar(100);not null" json:"skill_name" validate:"required,max=100"`
	SkillLevel      *string    `gorm:"type:varchar(20);check:skill_level IN ('beginner','intermediate','advanced','expert')" json:"skill_level,omitempty" validate:"omitempty,oneof=beginner intermediate advanced expert"`
	YearsExperience *int       `gorm:"type:int;check:years_experience >= 0" json:"years_experience,omitempty"`
	LastUsedAt      *time.Time `gorm:"type:date" json:"last_used_at,omitempty"`
	Verified        bool       `gorm:"default:false" json:"verified"`
	CreatedAt       time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for UserSkill
func (UserSkill) TableName() string {
	return "user_skills"
}

// UserCertification represents user's certifications
type UserCertification struct {
	ID                  int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID              int64      `gorm:"not null;index" json:"user_id"`
	CertificationName   string     `gorm:"type:varchar(150);not null" json:"certification_name" validate:"required,max=150"`
	IssuingOrganization string     `gorm:"type:varchar(150);not null" json:"issuing_organization" validate:"required,max=150"`
	IssueDate           *time.Time `gorm:"type:date" json:"issue_date,omitempty"`
	ExpirationDate      *time.Time `gorm:"type:date" json:"expiration_date,omitempty"`
	CredentialID        *string    `gorm:"type:varchar(100)" json:"credential_id,omitempty"`
	CredentialURL       *string    `gorm:"type:text" json:"credential_url,omitempty"`
	Description         *string    `gorm:"type:text" json:"description,omitempty"`
	Verified            bool       `gorm:"default:false" json:"verified"`
	VerificationDate    *time.Time `gorm:"type:timestamp" json:"verification_date,omitempty"`
	FileURL             *string    `gorm:"type:text" json:"file_url,omitempty"`
	IsActive            bool       `gorm:"default:true" json:"is_active"`
	CreatedAt           time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for UserCertification
func (UserCertification) TableName() string {
	return "user_certifications"
}

// UserLanguage represents user's language proficiency
type UserLanguage struct {
	ID                 int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID             int64      `gorm:"not null;index" json:"user_id"`
	LanguageName       string     `gorm:"type:varchar(100);not null" json:"language_name" validate:"required,max=100"`
	ProficiencyLevel   *string    `gorm:"type:varchar(50);check:proficiency_level IN ('basic','intermediate','advanced','fluent','native')" json:"proficiency_level,omitempty" validate:"omitempty,oneof=basic intermediate advanced fluent native"`
	CertificationName  *string    `gorm:"type:varchar(100)" json:"certification_name,omitempty"`
	CertificationScore *string    `gorm:"type:varchar(50)" json:"certification_score,omitempty"`
	CertificationDate  *time.Time `gorm:"type:date" json:"certification_date,omitempty"`
	Verified           bool       `gorm:"default:false" json:"verified"`
	IsActive           bool       `gorm:"default:true" json:"is_active"`
	Notes              *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt          time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for UserLanguage
func (UserLanguage) TableName() string {
	return "user_languages"
}

// UserProject represents user's projects
type UserProject struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64      `gorm:"not null;index" json:"user_id"`
	ProjectTitle  string     `gorm:"type:varchar(150);not null" json:"project_title" validate:"required,max=150"`
	RoleInProject *string    `gorm:"type:varchar(100)" json:"role_in_project,omitempty"`
	ProjectType   *string    `gorm:"type:varchar(50);check:project_type IN ('personal','freelance','company','academic','community')" json:"project_type,omitempty" validate:"omitempty,oneof=personal freelance company academic community"`
	Description   *string    `gorm:"type:text" json:"description,omitempty"`
	Industry      *string    `gorm:"type:varchar(100)" json:"industry,omitempty"`
	StartDate     *time.Time `gorm:"type:date" json:"start_date,omitempty"`
	EndDate       *time.Time `gorm:"type:date" json:"end_date,omitempty"`
	IsCurrent     bool       `gorm:"default:false" json:"is_current"`
	ProjectURL    *string    `gorm:"type:text" json:"project_url,omitempty"`
	RepoURL       *string    `gorm:"type:text" json:"repo_url,omitempty"`
	MediaURLs     *string    `gorm:"type:text[]" json:"media_urls,omitempty"`    // PostgreSQL array
	SkillsUsed    *string    `gorm:"type:text[]" json:"skills_used,omitempty"`   // PostgreSQL array
	Collaborators *string    `gorm:"type:text[]" json:"collaborators,omitempty"` // PostgreSQL array
	Verified      bool       `gorm:"default:false" json:"verified"`
	Featured      bool       `gorm:"default:false" json:"featured"`
	Visibility    string     `gorm:"type:varchar(20);default:'public';check:visibility IN ('public','private','limited')" json:"visibility"`
	CreatedAt     time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for UserProject
func (UserProject) TableName() string {
	return "user_projects"
}

// UserDocument represents user's documents
type UserDocument struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int64      `gorm:"not null;index" json:"user_id"`
	DocumentType *string    `gorm:"type:varchar(50);check:document_type IN ('resume','id_card','certificate','portfolio','transcript','others')" json:"document_type,omitempty" validate:"omitempty,oneof=resume id_card certificate portfolio transcript others"`
	DocumentName string     `gorm:"type:varchar(150);not null" json:"document_name" validate:"required,max=150"`
	FileURL      string     `gorm:"type:text;not null" json:"file_url" validate:"required"`
	FileSize     *int64     `gorm:"type:bigint" json:"file_size,omitempty"`
	MimeType     *string    `gorm:"type:varchar(100)" json:"mime_type,omitempty"`
	Description  *string    `gorm:"type:text" json:"description,omitempty"`
	UploadedAt   time.Time  `gorm:"type:timestamp;default:now()" json:"uploaded_at"`
	Verified     bool       `gorm:"default:false" json:"verified"`
	VerifiedAt   *time.Time `gorm:"type:timestamp" json:"verified_at,omitempty"`
	VerifiedBy   *int64     `gorm:"type:bigint" json:"verified_by,omitempty"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	Checksum     *string    `gorm:"type:varchar(100)" json:"checksum,omitempty"`
	CreatedAt    time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for UserDocument
func (UserDocument) TableName() string {
	return "user_documents"
}
