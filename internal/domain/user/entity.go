package user

import (
	"time"

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
	ID                 int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID             int64      `gorm:"not null;uniqueIndex" json:"user_id"`
	Headline           *string    `gorm:"type:varchar(150)" json:"headline,omitempty" validate:"omitempty,max=150"`
	Bio                *string    `gorm:"type:text" json:"bio,omitempty"`
	Gender             *string    `gorm:"type:varchar(10);check:gender IN ('male','female','other')" json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	BirthDate          *time.Time `gorm:"type:date" json:"birth_date,omitempty"`
	LocationCity       *string    `gorm:"type:varchar(100)" json:"location_city,omitempty"`
	LocationCountry    *string    `gorm:"type:varchar(100)" json:"location_country,omitempty"`
	DesiredPosition    *string    `gorm:"type:varchar(150)" json:"desired_position,omitempty"`
	DesiredSalaryMin   *float64   `gorm:"type:numeric(12,2)" json:"desired_salary_min,omitempty"`
	DesiredSalaryMax   *float64   `gorm:"type:numeric(12,2)" json:"desired_salary_max,omitempty"`
	ExperienceLevel    *string    `gorm:"type:varchar(50);check:experience_level IN ('internship','junior','mid','senior','lead')" json:"experience_level,omitempty" validate:"omitempty,oneof=internship junior mid senior lead"`
	IndustryInterest   *string    `gorm:"type:varchar(100)" json:"industry_interest,omitempty"`
	AvailabilityStatus string     `gorm:"type:varchar(50);default:'open'" json:"availability_status" validate:"omitempty,oneof=open looking_actively not_looking"`
	ProfileVisibility  bool       `gorm:"default:true" json:"profile_visibility"`
	Slug               *string    `gorm:"type:varchar(100);uniqueIndex" json:"slug,omitempty"`
	AvatarURL          *string    `gorm:"type:text" json:"avatar_url,omitempty"`
	CoverURL           *string    `gorm:"type:text" json:"cover_url,omitempty"`
	CreatedAt          time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for UserProfile
func (UserProfile) TableName() string {
	return "user_profiles"
}

// UserPreference represents user preferences and settings
type UserPreference struct {
	ID                  int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID              int64     `gorm:"not null;uniqueIndex" json:"user_id"`
	LanguagePreference  string    `gorm:"type:varchar(10);default:'id'" json:"language_preference"`
	ThemePreference     string    `gorm:"type:varchar(10);default:'light'" json:"theme_preference"`
	PreferredJobType    string    `gorm:"type:varchar(50);default:'onsite'" json:"preferred_job_type"`
	PreferredIndustry   *string   `gorm:"type:varchar(100)" json:"preferred_industry,omitempty"`
	PreferredLocation   *string   `gorm:"type:varchar(100)" json:"preferred_location,omitempty"`
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
}

// TableName specifies the table name for UserPreference
func (UserPreference) TableName() string {
	return "user_preferences"
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
