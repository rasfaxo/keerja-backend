package request

// UpdateProfileRequest represents user profile update request
type UpdateProfileRequest struct {
	FullName     *string `json:"full_name" validate:"omitempty,min=3,max=150"`
	Phone        *string `json:"phone" validate:"omitempty,min=10,max=20"`
	Headline     *string `json:"headline" validate:"omitempty,max=150"`
	Summary      *string `json:"summary" validate:"omitempty"`
	DateOfBirth  *string `json:"date_of_birth" validate:"omitempty"`
	Gender       *string `json:"gender" validate:"omitempty,oneof=male female other"`
	Nationality  *string `json:"nationality" validate:"omitempty,max=100"`
	Address      *string `json:"address" validate:"omitempty"`
	City         *string `json:"city" validate:"omitempty,max=100"`
	Province     *string `json:"province" validate:"omitempty,max=100"`
	Country      *string `json:"country" validate:"omitempty,max=100"`
	PostalCode   *string `json:"postal_code" validate:"omitempty,max=10"`
	LinkedinURL  *string `json:"linkedin_url" validate:"omitempty,url"`
	PortfolioURL *string `json:"portfolio_url" validate:"omitempty,url"`
	GithubURL    *string `json:"github_url" validate:"omitempty,url"`
}

// AddEducationRequest represents add education request
type AddEducationRequest struct {
	InstitutionName     string  `json:"institution_name" validate:"required,max=200"`
	Degree              string  `json:"degree" validate:"required,max=100"`
	FieldOfStudy        string  `json:"field_of_study" validate:"required,max=150"`
	StartDate           string  `json:"start_date" validate:"required"`
	EndDate             *string `json:"end_date" validate:"omitempty"`
	Grade               *string `json:"grade" validate:"omitempty,max=20"`
	Description         *string `json:"description" validate:"omitempty"`
	IsCurrentlyStudying bool    `json:"is_currently_studying"`
}

// UpdateEducationRequest represents update education request
type UpdateEducationRequest struct {
	InstitutionName     *string `json:"institution_name" validate:"omitempty,max=200"`
	Degree              *string `json:"degree" validate:"omitempty,max=100"`
	FieldOfStudy        *string `json:"field_of_study" validate:"omitempty,max=150"`
	StartDate           *string `json:"start_date" validate:"omitempty"`
	EndDate             *string `json:"end_date" validate:"omitempty"`
	Grade               *string `json:"grade" validate:"omitempty,max=20"`
	Description         *string `json:"description" validate:"omitempty"`
	IsCurrentlyStudying *bool   `json:"is_currently_studying"`
}

// AddExperienceRequest represents add experience request
type AddExperienceRequest struct {
	CompanyName        string  `json:"company_name" validate:"required,max=200"`
	JobTitle           string  `json:"job_title" validate:"required,max=150"`
	EmploymentType     string  `json:"employment_type" validate:"required,oneof='Full-Time' 'Part-Time' 'Contract' 'Internship' 'Freelance'"`
	Location           *string `json:"location" validate:"omitempty,max=150"`
	StartDate          string  `json:"start_date" validate:"required"`
	EndDate            *string `json:"end_date" validate:"omitempty"`
	IsCurrentlyWorking bool    `json:"is_currently_working"`
	Description        *string `json:"description" validate:"omitempty"`
}

// UpdateExperienceRequest represents update experience request
type UpdateExperienceRequest struct {
	CompanyName        *string `json:"company_name" validate:"omitempty,max=200"`
	JobTitle           *string `json:"job_title" validate:"omitempty,max=150"`
	EmploymentType     *string `json:"employment_type" validate:"omitempty,oneof='Full-Time' 'Part-Time' 'Contract' 'Internship' 'Freelance'"`
	Location           *string `json:"location" validate:"omitempty,max=150"`
	StartDate          *string `json:"start_date" validate:"omitempty"`
	EndDate            *string `json:"end_date" validate:"omitempty"`
	IsCurrentlyWorking *bool   `json:"is_currently_working"`
	Description        *string `json:"description" validate:"omitempty"`
}

// AddSkillRequest represents add skill request
type AddSkillRequest struct {
	SkillID           int64  `json:"skill_id" validate:"required,min=1"`
	ProficiencyLevel  string `json:"proficiency_level" validate:"required,oneof=beginner intermediate advanced expert"`
	YearsOfExperience *int16 `json:"years_of_experience" validate:"omitempty,min=0"`
}

// UpdateSkillRequest represents update skill request
type UpdateSkillRequest struct {
	ProficiencyLevel  *string `json:"proficiency_level" validate:"omitempty,oneof=beginner intermediate advanced expert"`
	YearsOfExperience *int16  `json:"years_of_experience" validate:"omitempty,min=0"`
}

// AddCertificationRequest represents add certification request
type AddCertificationRequest struct {
	Name                string  `json:"name" validate:"required,max=200"`
	IssuingOrganization string  `json:"issuing_organization" validate:"required,max=200"`
	IssueDate           string  `json:"issue_date" validate:"required"`
	ExpiryDate          *string `json:"expiry_date" validate:"omitempty"`
	CredentialID        *string `json:"credential_id" validate:"omitempty,max=100"`
	CredentialURL       *string `json:"credential_url" validate:"omitempty,url"`
	DoesNotExpire       bool    `json:"does_not_expire"`
}

// UpdateCertificationRequest represents update certification request
type UpdateCertificationRequest struct {
	Name                *string `json:"name" validate:"omitempty,max=200"`
	IssuingOrganization *string `json:"issuing_organization" validate:"omitempty,max=200"`
	IssueDate           *string `json:"issue_date" validate:"omitempty"`
	ExpiryDate          *string `json:"expiry_date" validate:"omitempty"`
	CredentialID        *string `json:"credential_id" validate:"omitempty,max=100"`
	CredentialURL       *string `json:"credential_url" validate:"omitempty,url"`
	DoesNotExpire       *bool   `json:"does_not_expire"`
}

// AddLanguageRequest represents add language request
type AddLanguageRequest struct {
	LanguageName     string `json:"language_name" validate:"required,max=100"`
	ProficiencyLevel string `json:"proficiency_level" validate:"required,oneof=elementary limited professional 'full professional' native"`
}

// UpdateLanguageRequest represents update language request
type UpdateLanguageRequest struct {
	ProficiencyLevel string `json:"proficiency_level" validate:"required,oneof=elementary limited professional 'full professional' native"`
}

// AddProjectRequest represents add project request
type AddProjectRequest struct {
	ProjectName  string  `json:"project_name" validate:"required,max=200"`
	Description  string  `json:"description" validate:"required"`
	Role         *string `json:"role" validate:"omitempty,max=100"`
	StartDate    string  `json:"start_date" validate:"required"`
	EndDate      *string `json:"end_date" validate:"omitempty"`
	ProjectURL   *string `json:"project_url" validate:"omitempty,url"`
	IsOngoing    bool    `json:"is_ongoing"`
	Technologies *string `json:"technologies" validate:"omitempty"`
}

// UpdateProjectRequest represents update project request
type UpdateProjectRequest struct {
	ProjectName  *string `json:"project_name" validate:"omitempty,max=200"`
	Description  *string `json:"description" validate:"omitempty"`
	Role         *string `json:"role" validate:"omitempty,max=100"`
	StartDate    *string `json:"start_date" validate:"omitempty"`
	EndDate      *string `json:"end_date" validate:"omitempty"`
	ProjectURL   *string `json:"project_url" validate:"omitempty,url"`
	IsOngoing    *bool   `json:"is_ongoing"`
	Technologies *string `json:"technologies" validate:"omitempty"`
}

// UploadDocumentRequest represents document upload request
type UploadDocumentRequest struct {
	DocumentType string `form:"document_type" validate:"required,oneof=resume cover_letter portfolio certificate other"`
	Title        string `form:"title" validate:"required,max=200"`
	Description  string `form:"description" validate:"omitempty"`
}

// UpdatePreferenceRequest represents user preferences update request
type UpdatePreferenceRequest struct {
	JobTypes            []string `json:"job_types" validate:"omitempty"`
	PreferredLocations  []string `json:"preferred_locations" validate:"omitempty"`
	ExpectedSalaryMin   *int64   `json:"expected_salary_min" validate:"omitempty,min=0"`
	ExpectedSalaryMax   *int64   `json:"expected_salary_max" validate:"omitempty,min=0,gtefield=ExpectedSalaryMin"`
	Currency            *string  `json:"currency" validate:"omitempty,len=3"`
	WillingToRelocate   *bool    `json:"willing_to_relocate"`
	AvailableForRemote  *bool    `json:"available_for_remote"`
	NoticePeriodInDays  *int16   `json:"notice_period_in_days" validate:"omitempty,min=0"`
	IsOpenToWork        *bool    `json:"is_open_to_work"`
	PreferredIndustries *string  `json:"preferred_industries" validate:"omitempty"`
	JobAlertFrequency   *string  `json:"job_alert_frequency" validate:"omitempty,oneof=daily weekly never"`
	ProfileVisibility   *string  `json:"profile_visibility" validate:"omitempty,oneof=public private recruiter_only"`
}

// UserSearchRequest represents user search request
type UserSearchRequest struct {
	Query      string  `json:"query" query:"q" validate:"omitempty"`
	Skills     []int64 `json:"skills" query:"skills" validate:"omitempty"`
	Location   string  `json:"location" query:"location" validate:"omitempty"`
	Experience string  `json:"experience" query:"experience" validate:"omitempty"`
	Page       int     `json:"page" query:"page" validate:"omitempty,min=1"`
	Limit      int     `json:"limit" query:"limit" validate:"omitempty,min=1,max=100"`
}
