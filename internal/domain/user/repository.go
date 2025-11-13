package user

import (
	"context"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// User CRUD
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByUUID(ctx context.Context, uuid string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter *UserFilter) ([]User, int64, error)

	// Profile operations
	CreateProfile(ctx context.Context, profile *UserProfile) error
	FindProfileByUserID(ctx context.Context, userID int64) (*UserProfile, error)
	FindProfileBySlug(ctx context.Context, slug string) (*UserProfile, error)
	UpdateProfile(ctx context.Context, profile *UserProfile) error

	// Preference operations
	CreatePreference(ctx context.Context, preference *UserPreference) error
	FindPreferenceByUserID(ctx context.Context, userID int64) (*UserPreference, error)
	UpdatePreference(ctx context.Context, preference *UserPreference) error

	// Education operations
	AddEducation(ctx context.Context, education *UserEducation) error
	UpdateEducation(ctx context.Context, education *UserEducation) error
	DeleteEducation(ctx context.Context, id int64) error
	GetEducationsByUserID(ctx context.Context, userID int64) ([]UserEducation, error)

	// Experience operations
	AddExperience(ctx context.Context, experience *UserExperience) error
	UpdateExperience(ctx context.Context, experience *UserExperience) error
	DeleteExperience(ctx context.Context, id int64) error
	GetExperiencesByUserID(ctx context.Context, userID int64) ([]UserExperience, error)

	// Skill operations
	AddSkill(ctx context.Context, skill *UserSkill) error
	UpdateSkill(ctx context.Context, skill *UserSkill) error
	DeleteSkill(ctx context.Context, id int64) error
	GetSkillsByUserID(ctx context.Context, userID int64) ([]UserSkill, error)

	// Certification operations
	AddCertification(ctx context.Context, cert *UserCertification) error
	UpdateCertification(ctx context.Context, cert *UserCertification) error
	DeleteCertification(ctx context.Context, id int64) error
	GetCertificationsByUserID(ctx context.Context, userID int64) ([]UserCertification, error)

	// Language operations
	AddLanguage(ctx context.Context, lang *UserLanguage) error
	UpdateLanguage(ctx context.Context, lang *UserLanguage) error
	DeleteLanguage(ctx context.Context, id int64) error
	GetLanguagesByUserID(ctx context.Context, userID int64) ([]UserLanguage, error)

	// Project operations
	AddProject(ctx context.Context, project *UserProject) error
	UpdateProject(ctx context.Context, project *UserProject) error
	DeleteProject(ctx context.Context, id int64) error
	GetProjectsByUserID(ctx context.Context, userID int64) ([]UserProject, error)

	// Document operations
	AddDocument(ctx context.Context, doc *UserDocument) error
	UpdateDocument(ctx context.Context, doc *UserDocument) error
	DeleteDocument(ctx context.Context, id int64) error
	GetDocumentsByUserID(ctx context.Context, userID int64) ([]UserDocument, error)

	// Full profile with relationships
	GetFullProfile(ctx context.Context, userID int64) (*User, error)
}

// UserFilter represents filters for querying users
type UserFilter struct {
	UserType         *string
	Status           *string
	IsVerified       *bool
	LocationCity     *string
	ExperienceLevel  *string
	IndustryInterest *string
	SkillNames       []string
	SearchQuery      *string
	Page             int
	Limit            int
	SortBy           string
	SortOrder        string
}
