package user

import (
	"context"
	"mime/multipart"
)

// UserService defines the business logic interface for user operations
type UserService interface {
	// Registration and verification
	Register(ctx context.Context, req *RegisterRequest) (*User, error)
	VerifyEmail(ctx context.Context, token string) error
	ResendVerificationEmail(ctx context.Context, email string) error

	// Profile management
	GetProfile(ctx context.Context, userID int64) (*User, error)
	GetProfileBySlug(ctx context.Context, slug string) (*User, error)
	UpdateProfile(ctx context.Context, userID int64, req *UpdateProfileRequest) error
	UploadAvatar(ctx context.Context, userID int64, file *multipart.FileHeader) (string, error)
	UploadCover(ctx context.Context, userID int64, file *multipart.FileHeader) (string, error)
	DeleteAvatar(ctx context.Context, userID int64) error
	DeleteCover(ctx context.Context, userID int64) error

	// Preference management
	GetPreferences(ctx context.Context, userID int64) (*UserPreference, error)
	UpdatePreferences(ctx context.Context, userID int64, req *UpdatePreferenceRequest) error

	// Education management
	AddEducation(ctx context.Context, userID int64, req *AddEducationRequest) (*UserEducation, error)
	UpdateEducation(ctx context.Context, userID int64, educationID int64, req *UpdateEducationRequest) error
	DeleteEducation(ctx context.Context, userID int64, educationID int64) error
	GetEducations(ctx context.Context, userID int64) ([]UserEducation, error)

	// Experience management
	AddExperience(ctx context.Context, userID int64, req *AddExperienceRequest) (*UserExperience, error)
	UpdateExperience(ctx context.Context, userID int64, experienceID int64, req *UpdateExperienceRequest) error
	DeleteExperience(ctx context.Context, userID int64, experienceID int64) error
	GetExperiences(ctx context.Context, userID int64) ([]UserExperience, error)

	// Skill management
	AddSkills(ctx context.Context, userID int64, skillNames []string) error
	UpdateSkill(ctx context.Context, userID int64, skillID int64, req *UpdateSkillRequest) error
	DeleteSkill(ctx context.Context, userID int64, skillID int64) error
	GetSkills(ctx context.Context, userID int64) ([]UserSkill, error)

	// Certification management
	AddCertification(ctx context.Context, userID int64, req *AddCertificationRequest) (*UserCertification, error)
	UpdateCertification(ctx context.Context, userID int64, certID int64, req *UpdateCertificationRequest) error
	DeleteCertification(ctx context.Context, userID int64, certID int64) error
	GetCertifications(ctx context.Context, userID int64) ([]UserCertification, error)

	// Language management
	AddLanguage(ctx context.Context, userID int64, req *AddLanguageRequest) (*UserLanguage, error)
	UpdateLanguage(ctx context.Context, userID int64, langID int64, req *UpdateLanguageRequest) error
	DeleteLanguage(ctx context.Context, userID int64, langID int64) error
	GetLanguages(ctx context.Context, userID int64) ([]UserLanguage, error)

	// Project management
	AddProject(ctx context.Context, userID int64, req *AddProjectRequest) (*UserProject, error)
	UpdateProject(ctx context.Context, userID int64, projectID int64, req *UpdateProjectRequest) error
	DeleteProject(ctx context.Context, userID int64, projectID int64) error
	GetProjects(ctx context.Context, userID int64) ([]UserProject, error)

	// Document management
	UploadDocument(ctx context.Context, userID int64, file *multipart.FileHeader, req *UploadDocumentRequest) (*UserDocument, error)
	DeleteDocument(ctx context.Context, userID int64, documentID int64) error
	GetDocuments(ctx context.Context, userID int64) ([]UserDocument, error)

	// Search and discovery
	SearchUsers(ctx context.Context, filter *UserFilter) ([]User, int64, error)
	GetUsersBySkills(ctx context.Context, skillNames []string) ([]User, error)

	// Profile completion and analytics
	GetProfileCompletionPercentage(ctx context.Context, userID int64) (int, error)
	UpdateLastLogin(ctx context.Context, userID int64) error

	// Account management
	UpdateStatus(ctx context.Context, userID int64, status string) error
	SuspendAccount(ctx context.Context, userID int64, reason string) error
	DeactivateAccount(ctx context.Context, userID int64) error
	DeleteAccount(ctx context.Context, userID int64) error
}

// Request DTOs (simplified - will be detailed in DTO layer)

type RegisterRequest struct {
	FullName string
	Email    string
	Phone    *string
	Password string
	UserType string
}

type UpdateProfileRequest struct {
	Headline           *string
	Bio                *string
	Gender             *string
	BirthDate          *string
	LocationCity       *string
	LocationCountry    *string
	DesiredPosition    *string
	DesiredSalaryMin   *float64
	DesiredSalaryMax   *float64
	ExperienceLevel    *string
	IndustryInterest   *string
	AvailabilityStatus *string
}

type UpdatePreferenceRequest struct {
	LanguagePreference  *string
	ThemePreference     *string
	PreferredJobType    *string
	PreferredIndustry   *string
	PreferredLocation   *string
	PreferredSalaryMin  *float64
	PreferredSalaryMax  *float64
	EmailNotifications  *bool
	SMSNotifications    *bool
	PushNotifications   *bool
	EmailMarketing      *bool
	ProfileVisibility   *string
	ShowOnlineStatus    *bool
	AllowDirectMessages *bool
	DataSharingConsent  *bool
}

type AddEducationRequest struct {
	InstitutionName string
	Major           *string
	DegreeLevel     *string
	StartYear       *int
	EndYear         *int
	GPA             *float64
	Activities      *string
	Description     *string
	IsCurrent       bool
}

type UpdateEducationRequest struct {
	InstitutionName *string
	Major           *string
	DegreeLevel     *string
	StartYear       *int
	EndYear         *int
	GPA             *float64
	Activities      *string
	Description     *string
	IsCurrent       *bool
}

type AddExperienceRequest struct {
	CompanyName     string
	PositionTitle   string
	Industry        *string
	EmploymentType  *string
	StartDate       string
	EndDate         *string
	IsCurrent       bool
	Description     *string
	Achievements    *string
	LocationCity    *string
	LocationCountry *string
}

type UpdateExperienceRequest struct {
	CompanyName     *string
	PositionTitle   *string
	Industry        *string
	EmploymentType  *string
	StartDate       *string
	EndDate         *string
	IsCurrent       *bool
	Description     *string
	Achievements    *string
	LocationCity    *string
	LocationCountry *string
}

type UpdateSkillRequest struct {
	SkillLevel      *string
	YearsExperience *int
	LastUsedAt      *string
}

type AddCertificationRequest struct {
	CertificationName   string
	IssuingOrganization string
	IssueDate           *string
	ExpirationDate      *string
	CredentialID        *string
	CredentialURL       *string
	Description         *string
	FileURL             *string
}

type UpdateCertificationRequest struct {
	CertificationName   *string
	IssuingOrganization *string
	IssueDate           *string
	ExpirationDate      *string
	CredentialID        *string
	CredentialURL       *string
	Description         *string
	FileURL             *string
}

type AddLanguageRequest struct {
	LanguageName       string
	ProficiencyLevel   *string
	CertificationName  *string
	CertificationScore *string
	CertificationDate  *string
	Notes              *string
}

type UpdateLanguageRequest struct {
	LanguageName       *string
	ProficiencyLevel   *string
	CertificationName  *string
	CertificationScore *string
	CertificationDate  *string
	Notes              *string
}

type AddProjectRequest struct {
	ProjectTitle  string
	RoleInProject *string
	ProjectType   *string
	Description   *string
	Industry      *string
	StartDate     *string
	EndDate       *string
	IsCurrent     bool
	ProjectURL    *string
	RepoURL       *string
	MediaURLs     []string
	SkillsUsed    []string
	Collaborators []string
	Visibility    *string
}

type UpdateProjectRequest struct {
	ProjectTitle  *string
	RoleInProject *string
	ProjectType   *string
	Description   *string
	Industry      *string
	StartDate     *string
	EndDate       *string
	IsCurrent     *bool
	ProjectURL    *string
	RepoURL       *string
	MediaURLs     []string
	SkillsUsed    []string
	Collaborators []string
	Visibility    *string
}

type UploadDocumentRequest struct {
	DocumentType *string
	DocumentName string
	Description  *string
}
