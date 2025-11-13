package response

import "time"

// UserResponse represents user public response
type UserResponse struct {
	ID         int64                `json:"id"`
	UUID       string               `json:"uuid"`
	FullName   string               `json:"full_name"`
	Email      string               `json:"email"`
	Phone      string               `json:"phone,omitempty"`
	UserType   string               `json:"user_type"`
	IsVerified bool                 `json:"is_verified"`
	Status     string               `json:"status"`
	LastLogin  *time.Time           `json:"last_login,omitempty"`
	CreatedAt  time.Time            `json:"created_at"`
	Profile    *UserProfileResponse `json:"profile,omitempty"`
}

// UserProfileResponse represents user profile response
type UserProfileResponse struct {
	ID                 int64      `json:"id"`
	UserID             int64      `json:"user_id"`
	Headline           *string    `json:"headline,omitempty"`
	Bio                *string    `json:"bio,omitempty"`
	Gender             *string    `json:"gender,omitempty"`
	BirthDate          *time.Time `json:"birth_date,omitempty"`
	Nationality        *string    `json:"nationality,omitempty"`
	Address            *string    `json:"address,omitempty"`
	LocationCity       *string    `json:"location_city,omitempty"`
	LocationState      *string    `json:"location_state,omitempty"`
	LocationCountry    *string    `json:"location_country,omitempty"`
	PostalCode         *string    `json:"postal_code,omitempty"`
	LinkedInURL        *string    `json:"linkedin_url,omitempty"`
	PortfolioURL       *string    `json:"portfolio_url,omitempty"`
	GithubURL          *string    `json:"github_url,omitempty"`
	DesiredPosition    *string    `json:"desired_position,omitempty"`
	DesiredSalaryMin   *float64   `json:"desired_salary_min,omitempty"`
	DesiredSalaryMax   *float64   `json:"desired_salary_max,omitempty"`
	ExperienceLevel    *string    `json:"experience_level,omitempty"`
	IndustryInterest   *string    `json:"industry_interest,omitempty"`
	AvailabilityStatus string     `json:"availability_status"`
	ProfileVisibility  bool       `json:"profile_visibility"`
	Slug               *string    `json:"slug,omitempty"`
	AvatarURL          *string    `json:"avatar_url,omitempty"`
	CoverURL           *string    `json:"cover_url,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// UserDetailResponse represents detailed user response with all relations
type UserDetailResponse struct {
	ID             int64                       `json:"id"`
	UUID           string                      `json:"uuid"`
	FullName       string                      `json:"full_name"`
	Email          string                      `json:"email"`
	Phone          string                      `json:"phone,omitempty"`
	UserType       string                      `json:"user_type"`
	IsVerified     bool                        `json:"is_verified"`
	Status         string                      `json:"status"`
	LastLogin      *time.Time                  `json:"last_login,omitempty"`
	CreatedAt      time.Time                   `json:"created_at"`
	Profile        *UserProfileResponse        `json:"profile,omitempty"`
	Educations     []UserEducationResponse     `json:"educations,omitempty"`
	Experiences    []UserExperienceResponse    `json:"experiences,omitempty"`
	Skills         []UserSkillResponse         `json:"skills,omitempty"`
	Certifications []UserCertificationResponse `json:"certifications,omitempty"`
	Languages      []UserLanguageResponse      `json:"languages,omitempty"`
	Projects       []UserProjectResponse       `json:"projects,omitempty"`
	Documents      []UserDocumentResponse      `json:"documents,omitempty"`
	Preference     *UserPreferenceResponse     `json:"preference,omitempty"`
}

// UserEducationResponse represents education response
type UserEducationResponse struct {
	ID              int64    `json:"id"`
	InstitutionName string   `json:"institution_name"`
	Major           *string  `json:"major,omitempty"`
	DegreeLevel     *string  `json:"degree_level,omitempty"`
	StartYear       *int     `json:"start_year,omitempty"`
	EndYear         *int     `json:"end_year,omitempty"`
	GPA             *float64 `json:"gpa,omitempty"`
	Activities      *string  `json:"activities,omitempty"`
	Description     *string  `json:"description,omitempty"`
	IsCurrent       bool     `json:"is_current"`
	CreatedAt       string   `json:"created_at"`
}

// UserExperienceResponse represents experience response
type UserExperienceResponse struct {
	ID              int64   `json:"id"`
	CompanyName     string  `json:"company_name"`
	PositionTitle   string  `json:"position_title"`
	Industry        *string `json:"industry,omitempty"`
	EmploymentType  *string `json:"employment_type,omitempty"`
	StartDate       string  `json:"start_date"`
	EndDate         *string `json:"end_date,omitempty"`
	IsCurrent       bool    `json:"is_current"`
	Description     *string `json:"description,omitempty"`
	Achievements    *string `json:"achievements,omitempty"`
	LocationCity    *string `json:"location_city,omitempty"`
	LocationCountry string  `json:"location_country"`
	CreatedAt       string  `json:"created_at"`
}

// UserSkillResponse represents skill response
type UserSkillResponse struct {
	ID                int64     `json:"id"`
	SkillID           int64     `json:"skill_id"`
	SkillName         string    `json:"skill_name"`
	ProficiencyLevel  string    `json:"proficiency_level"`
	YearsOfExperience int16     `json:"years_of_experience,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

// UserCertificationResponse represents certification response
type UserCertificationResponse struct {
	ID                  int64      `json:"id"`
	Name                string     `json:"name"`
	IssuingOrganization string     `json:"issuing_organization"`
	IssueDate           time.Time  `json:"issue_date"`
	ExpiryDate          *time.Time `json:"expiry_date,omitempty"`
	CredentialID        string     `json:"credential_id,omitempty"`
	CredentialURL       string     `json:"credential_url,omitempty"`
	DoesNotExpire       bool       `json:"does_not_expire"`
	IsExpired           bool       `json:"is_expired"`
	CreatedAt           time.Time  `json:"created_at"`
}

// UserLanguageResponse represents language response
type UserLanguageResponse struct {
	ID               int64     `json:"id"`
	LanguageName     string    `json:"language_name"`
	ProficiencyLevel string    `json:"proficiency_level"`
	CreatedAt        time.Time `json:"created_at"`
}

// UserProjectResponse represents project response
type UserProjectResponse struct {
	ID           int64      `json:"id"`
	ProjectName  string     `json:"project_name"`
	Description  string     `json:"description"`
	Role         string     `json:"role,omitempty"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	ProjectURL   string     `json:"project_url,omitempty"`
	IsOngoing    bool       `json:"is_ongoing"`
	Technologies string     `json:"technologies,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// UserDocumentResponse represents document response
type UserDocumentResponse struct {
	ID           int64     `json:"id"`
	DocumentType string    `json:"document_type"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	FileURL      string    `json:"file_url"`
	FileName     string    `json:"file_name"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	IsVerified   bool      `json:"is_verified"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

// UserPreferenceResponse represents user preference response
type UserPreferenceResponse struct {
	ID                  int64     `json:"id"`
	JobTypes            []string  `json:"job_types,omitempty"`
	PreferredLocations  []string  `json:"preferred_locations,omitempty"`
	ExpectedSalaryMin   *int64    `json:"expected_salary_min,omitempty"`
	ExpectedSalaryMax   *int64    `json:"expected_salary_max,omitempty"`
	Currency            string    `json:"currency"`
	WillingToRelocate   bool      `json:"willing_to_relocate"`
	AvailableForRemote  bool      `json:"available_for_remote"`
	NoticePeriodInDays  int16     `json:"notice_period_in_days,omitempty"`
	IsOpenToWork        bool      `json:"is_open_to_work"`
	PreferredIndustries string    `json:"preferred_industries,omitempty"`
	JobAlertFrequency   string    `json:"job_alert_frequency"`
	ProfileVisibility   string    `json:"profile_visibility"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// UserListResponse represents list of users response
type UserListResponse struct {
	Users []UserResponse `json:"users"`
}
