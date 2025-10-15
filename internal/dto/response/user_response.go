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
	ID                  int64      `json:"id"`
	Headline            string     `json:"headline,omitempty"`
	Summary             string     `json:"summary,omitempty"`
	DateOfBirth         *time.Time `json:"date_of_birth,omitempty"`
	Gender              string     `json:"gender,omitempty"`
	Nationality         string     `json:"nationality,omitempty"`
	Address             string     `json:"address,omitempty"`
	City                string     `json:"city,omitempty"`
	Province            string     `json:"province,omitempty"`
	Country             string     `json:"country,omitempty"`
	PostalCode          string     `json:"postal_code,omitempty"`
	ProfilePictureURL   string     `json:"profile_picture_url,omitempty"`
	CoverImageURL       string     `json:"cover_image_url,omitempty"`
	LinkedinURL         string     `json:"linkedin_url,omitempty"`
	PortfolioURL        string     `json:"portfolio_url,omitempty"`
	GithubURL           string     `json:"github_url,omitempty"`
	ResumeURL           string     `json:"resume_url,omitempty"`
	TotalExperience     int16      `json:"total_experience"`
	ProfileCompleteness int16      `json:"profile_completeness"`
	UpdatedAt           time.Time  `json:"updated_at"`
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
	ID                  int64      `json:"id"`
	InstitutionName     string     `json:"institution_name"`
	Degree              string     `json:"degree"`
	FieldOfStudy        string     `json:"field_of_study"`
	StartDate           time.Time  `json:"start_date"`
	EndDate             *time.Time `json:"end_date,omitempty"`
	Grade               string     `json:"grade,omitempty"`
	Description         string     `json:"description,omitempty"`
	IsCurrentlyStudying bool       `json:"is_currently_studying"`
	CreatedAt           time.Time  `json:"created_at"`
}

// UserExperienceResponse represents experience response
type UserExperienceResponse struct {
	ID                 int64      `json:"id"`
	CompanyName        string     `json:"company_name"`
	JobTitle           string     `json:"job_title"`
	EmploymentType     string     `json:"employment_type"`
	Location           string     `json:"location,omitempty"`
	StartDate          time.Time  `json:"start_date"`
	EndDate            *time.Time `json:"end_date,omitempty"`
	IsCurrentlyWorking bool       `json:"is_currently_working"`
	Description        string     `json:"description,omitempty"`
	Duration           string     `json:"duration,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
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
