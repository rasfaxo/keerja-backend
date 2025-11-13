package factory

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserFactory provides methods to create test users
type UserFactory struct {
	sequence int
}

// NewUserFactory creates a new user factory
func NewUserFactory() *UserFactory {
	return &UserFactory{
		sequence: 0,
	}
}

// UserBuilder provides a fluent interface for building users
type UserBuilder struct {
	ID                int64
	Email             string
	Phone             string
	Password          string
	PasswordHash      string
	FullName          string
	Role              string
	IsEmailVerified   bool
	IsPhoneVerified   bool
	EmailVerifiedAt   *time.Time
	PhoneVerifiedAt   *time.Time
	EmailVerifyToken  *string
	PhoneVerifyToken  *string
	LastLoginAt       *time.Time
	ProfilePictureURL *string
	Bio               *string
	Location          *string
	DateOfBirth       *time.Time
	Gender            *string
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Build builds the user
func (ub *UserBuilder) Build() *UserBuilder {
	// Generate password hash if password is set and hash is not
	if ub.Password != "" && ub.PasswordHash == "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(ub.Password), bcrypt.DefaultCost)
		ub.PasswordHash = string(hash)
	}

	// Set default timestamps if not set
	if ub.CreatedAt.IsZero() {
		ub.CreatedAt = time.Now()
	}
	if ub.UpdatedAt.IsZero() {
		ub.UpdatedAt = time.Now()
	}

	return ub
}

// WithID sets the user ID
func (ub *UserBuilder) WithID(id int64) *UserBuilder {
	ub.ID = id
	return ub
}

// WithEmail sets the email
func (ub *UserBuilder) WithEmail(email string) *UserBuilder {
	ub.Email = email
	return ub
}

// WithPhone sets the phone
func (ub *UserBuilder) WithPhone(phone string) *UserBuilder {
	ub.Phone = phone
	return ub
}

// WithPassword sets the password
func (ub *UserBuilder) WithPassword(password string) *UserBuilder {
	ub.Password = password
	return ub
}

// WithPasswordHash sets the password hash directly
func (ub *UserBuilder) WithPasswordHash(hash string) *UserBuilder {
	ub.PasswordHash = hash
	return ub
}

// WithFullName sets the full name
func (ub *UserBuilder) WithFullName(name string) *UserBuilder {
	ub.FullName = name
	return ub
}

// WithRole sets the role
func (ub *UserBuilder) WithRole(role string) *UserBuilder {
	ub.Role = role
	return ub
}

// AsJobSeeker sets the role as job seeker
func (ub *UserBuilder) AsJobSeeker() *UserBuilder {
	ub.Role = "job_seeker"
	return ub
}

// AsEmployer sets the role as employer
func (ub *UserBuilder) AsEmployer() *UserBuilder {
	ub.Role = "employer"
	return ub
}

// Verified marks email and phone as verified
func (ub *UserBuilder) Verified() *UserBuilder {
	now := time.Now()
	ub.IsEmailVerified = true
	ub.IsPhoneVerified = true
	ub.EmailVerifiedAt = &now
	ub.PhoneVerifiedAt = &now
	return ub
}

// EmailVerified marks email as verified
func (ub *UserBuilder) EmailVerified() *UserBuilder {
	now := time.Now()
	ub.IsEmailVerified = true
	ub.EmailVerifiedAt = &now
	return ub
}

// PhoneVerified marks phone as verified
func (ub *UserBuilder) PhoneVerified() *UserBuilder {
	now := time.Now()
	ub.IsPhoneVerified = true
	ub.PhoneVerifiedAt = &now
	return ub
}

// Unverified marks as unverified
func (ub *UserBuilder) Unverified() *UserBuilder {
	ub.IsEmailVerified = false
	ub.IsPhoneVerified = false
	ub.EmailVerifiedAt = nil
	ub.PhoneVerifiedAt = nil
	return ub
}

// WithEmailVerifyToken sets the email verification token
func (ub *UserBuilder) WithEmailVerifyToken(token string) *UserBuilder {
	ub.EmailVerifyToken = &token
	return ub
}

// WithPhoneVerifyToken sets the phone verification token
func (ub *UserBuilder) WithPhoneVerifyToken(token string) *UserBuilder {
	ub.PhoneVerifyToken = &token
	return ub
}

// WithLastLoginAt sets the last login time
func (ub *UserBuilder) WithLastLoginAt(t time.Time) *UserBuilder {
	ub.LastLoginAt = &t
	return ub
}

// WithProfilePicture sets the profile picture URL
func (ub *UserBuilder) WithProfilePicture(url string) *UserBuilder {
	ub.ProfilePictureURL = &url
	return ub
}

// WithBio sets the bio
func (ub *UserBuilder) WithBio(bio string) *UserBuilder {
	ub.Bio = &bio
	return ub
}

// WithLocation sets the location
func (ub *UserBuilder) WithLocation(location string) *UserBuilder {
	ub.Location = &location
	return ub
}

// WithDateOfBirth sets the date of birth
func (ub *UserBuilder) WithDateOfBirth(dob time.Time) *UserBuilder {
	ub.DateOfBirth = &dob
	return ub
}

// WithGender sets the gender
func (ub *UserBuilder) WithGender(gender string) *UserBuilder {
	ub.Gender = &gender
	return ub
}

// Active marks user as active
func (ub *UserBuilder) Active() *UserBuilder {
	ub.IsActive = true
	return ub
}

// Inactive marks user as inactive
func (ub *UserBuilder) Inactive() *UserBuilder {
	ub.IsActive = false
	return ub
}

// WithCreatedAt sets the created at time
func (ub *UserBuilder) WithCreatedAt(t time.Time) *UserBuilder {
	ub.CreatedAt = t
	return ub
}

// WithUpdatedAt sets the updated at time
func (ub *UserBuilder) WithUpdatedAt(t time.Time) *UserBuilder {
	ub.UpdatedAt = t
	return ub
}

// CreateUser creates a user builder with default values
func (f *UserFactory) CreateUser() *UserBuilder {
	f.sequence++
	return &UserBuilder{
		ID:              int64(f.sequence),
		Email:           fmt.Sprintf("user%d@test.com", f.sequence),
		Phone:           fmt.Sprintf("+62812345%05d", f.sequence),
		Password:        "Password123!",
		FullName:        fmt.Sprintf("Test User %d", f.sequence),
		Role:            "job_seeker",
		IsEmailVerified: true,
		IsPhoneVerified: true,
		IsActive:        true,
	}
}

// CreateJobSeeker creates a job seeker user
func (f *UserFactory) CreateJobSeeker() *UserBuilder {
	return f.CreateUser().AsJobSeeker().Verified()
}

// CreateEmployer creates an employer user
func (f *UserFactory) CreateEmployer() *UserBuilder {
	return f.CreateUser().AsEmployer().Verified()
}

// CreateUnverifiedUser creates an unverified user
func (f *UserFactory) CreateUnverifiedUser() *UserBuilder {
	token := uuid.New().String()
	return f.CreateUser().Unverified().WithEmailVerifyToken(token)
}

// CreateUserWithEmail creates a user with specific email
func (f *UserFactory) CreateUserWithEmail(email string) *UserBuilder {
	return f.CreateUser().WithEmail(email)
}

// CreateUserWithPhone creates a user with specific phone
func (f *UserFactory) CreateUserWithPhone(phone string) *UserBuilder {
	return f.CreateUser().WithPhone(phone)
}

// CreateMultipleUsers creates multiple users
func (f *UserFactory) CreateMultipleUsers(count int) []*UserBuilder {
	users := make([]*UserBuilder, count)
	for i := 0; i < count; i++ {
		users[i] = f.CreateUser()
	}
	return users
}

// CreateMultipleJobSeekers creates multiple job seekers
func (f *UserFactory) CreateMultipleJobSeekers(count int) []*UserBuilder {
	users := make([]*UserBuilder, count)
	for i := 0; i < count; i++ {
		users[i] = f.CreateJobSeeker()
	}
	return users
}

// CreateMultipleEmployers creates multiple employers
func (f *UserFactory) CreateMultipleEmployers(count int) []*UserBuilder {
	users := make([]*UserBuilder, count)
	for i := 0; i < count; i++ {
		users[i] = f.CreateEmployer()
	}
	return users
}

// RandomUser creates a random user with random data
func (f *UserFactory) RandomUser() *UserBuilder {
	return f.CreateUser().
		WithEmail(fmt.Sprintf("random_%s@test.com", uuid.New().String()[:8])).
		WithPhone(fmt.Sprintf("+62812%08d", time.Now().UnixNano()%100000000)).
		WithFullName(fmt.Sprintf("Random User %s", uuid.New().String()[:8]))
}

// UserProfileBuilder extends UserBuilder with profile data
type UserProfileBuilder struct {
	*UserBuilder
	Skills             []string
	Experience         []ExperienceData
	Education          []EducationData
	Certifications     []string
	Languages          []string
	PreferredJobTypes  []string
	PreferredLocations []string
	ExpectedSalaryMin  *int64
	ExpectedSalaryMax  *int64
	ResumeURL          *string
	PortfolioURL       *string
	LinkedInURL        *string
	GitHubURL          *string
}

// ExperienceData represents work experience
type ExperienceData struct {
	Company     string
	Position    string
	StartDate   time.Time
	EndDate     *time.Time
	Description string
	IsCurrent   bool
}

// EducationData represents education
type EducationData struct {
	Institution string
	Degree      string
	Field       string
	StartDate   time.Time
	EndDate     *time.Time
	GPA         *float64
	IsCurrent   bool
}

// WithProfile converts UserBuilder to UserProfileBuilder with profile data
func (ub *UserBuilder) WithProfile() *UserProfileBuilder {
	return &UserProfileBuilder{
		UserBuilder:        ub,
		Skills:             []string{"Go", "PostgreSQL", "Docker"},
		Experience:         []ExperienceData{},
		Education:          []EducationData{},
		Languages:          []string{"Indonesian", "English"},
		PreferredJobTypes:  []string{"full_time"},
		PreferredLocations: []string{"Jakarta, Indonesia"},
	}
}

// WithSkills sets the skills
func (upb *UserProfileBuilder) WithSkills(skills []string) *UserProfileBuilder {
	upb.Skills = skills
	return upb
}

// WithExperience adds work experience
func (upb *UserProfileBuilder) WithExperience(exp ExperienceData) *UserProfileBuilder {
	upb.Experience = append(upb.Experience, exp)
	return upb
}

// WithEducation adds education
func (upb *UserProfileBuilder) WithEducation(edu EducationData) *UserProfileBuilder {
	upb.Education = append(upb.Education, edu)
	return upb
}

// WithResume sets the resume URL
func (upb *UserProfileBuilder) WithResume(url string) *UserProfileBuilder {
	upb.ResumeURL = &url
	return upb
}

// WithPortfolio sets the portfolio URL
func (upb *UserProfileBuilder) WithPortfolio(url string) *UserProfileBuilder {
	upb.PortfolioURL = &url
	return upb
}

// WithLinkedIn sets the LinkedIn URL
func (upb *UserProfileBuilder) WithLinkedIn(url string) *UserProfileBuilder {
	upb.LinkedInURL = &url
	return upb
}

// WithGitHub sets the GitHub URL
func (upb *UserProfileBuilder) WithGitHub(url string) *UserProfileBuilder {
	upb.GitHubURL = &url
	return upb
}

// WithExpectedSalary sets the expected salary range
func (upb *UserProfileBuilder) WithExpectedSalary(min, max int64) *UserProfileBuilder {
	upb.ExpectedSalaryMin = &min
	upb.ExpectedSalaryMax = &max
	return upb
}
