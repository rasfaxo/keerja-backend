package fixtures

import (
	"time"

	"github.com/google/uuid"
)

// CommonTestData provides common test data used across tests
type CommonTestData struct {
	TestEmail         string
	TestPhone         string
	TestPassword      string
	TestHashedPass    string
	TestToken         string
	TestUserID        int64
	TestCompanyID     int64
	TestJobID         int64
	TestApplicationID int64
}

// DefaultTestData returns default test data
func DefaultTestData() *CommonTestData {
	return &CommonTestData{
		TestEmail:         "test@example.com",
		TestPhone:         "+628123456789",
		TestPassword:      "Test123!@#",
		TestHashedPass:    "$2a$10$abcdefghijklmnopqrstuv", // Placeholder bcrypt hash
		TestToken:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.token",
		TestUserID:        1,
		TestCompanyID:     1,
		TestJobID:         1,
		TestApplicationID: 1,
	}
}

// UserFixture represents a test user
type UserFixture struct {
	ID              int64
	Email           string
	Phone           string
	Password        string
	FullName        string
	Role            string
	IsEmailVerified bool
	IsPhoneVerified bool
	EmailVerifiedAt *time.Time
	PhoneVerifiedAt *time.Time
	LastLoginAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewJobSeekerUser creates a job seeker user fixture
func NewJobSeekerUser() *UserFixture {
	now := time.Now()
	return &UserFixture{
		ID:              1,
		Email:           "jobseeker@test.com",
		Phone:           "+628111111111",
		Password:        "$2a$10$abcdefghijklmnopqrstuv",
		FullName:        "Job Seeker Test",
		Role:            "job_seeker",
		IsEmailVerified: true,
		IsPhoneVerified: true,
		EmailVerifiedAt: &now,
		PhoneVerifiedAt: &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// NewEmployerUser creates an employer user fixture
func NewEmployerUser() *UserFixture {
	now := time.Now()
	return &UserFixture{
		ID:              2,
		Email:           "employer@test.com",
		Phone:           "+628222222222",
		Password:        "$2a$10$abcdefghijklmnopqrstuv",
		FullName:        "Employer Test",
		Role:            "employer",
		IsEmailVerified: true,
		IsPhoneVerified: true,
		EmailVerifiedAt: &now,
		PhoneVerifiedAt: &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// NewUnverifiedUser creates an unverified user fixture
func NewUnverifiedUser() *UserFixture {
	now := time.Now()
	return &UserFixture{
		ID:              3,
		Email:           "unverified@test.com",
		Phone:           "+628333333333",
		Password:        "$2a$10$abcdefghijklmnopqrstuv",
		FullName:        "Unverified Test",
		Role:            "job_seeker",
		IsEmailVerified: false,
		IsPhoneVerified: false,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// CompanyFixture represents a test company
type CompanyFixture struct {
	ID          int64
	Name        string
	Description string
	Industry    string
	Location    string
	Website     string
	Email       string
	Phone       string
	LogoURL     string
	IsVerified  bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewCompany creates a company fixture
func NewCompany() *CompanyFixture {
	now := time.Now()
	return &CompanyFixture{
		ID:          1,
		Name:        "Test Company Inc",
		Description: "A test company for testing purposes",
		Industry:    "Technology",
		Location:    "Jakarta, Indonesia",
		Website:     "https://testcompany.com",
		Email:       "info@testcompany.com",
		Phone:       "+628987654321",
		LogoURL:     "https://testcompany.com/logo.png",
		IsVerified:  true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// JobFixture represents a test job
type JobFixture struct {
	ID              int64
	CompanyID       int64
	Title           string
	Description     string
	Requirements    string
	Location        string
	EmploymentType  string
	ExperienceLevel string
	SalaryMin       *int64
	SalaryMax       *int64
	Skills          []string
	Status          string
	PostedAt        time.Time
	ClosingDate     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewJob creates a job fixture
func NewJob() *JobFixture {
	now := time.Now()
	closingDate := now.AddDate(0, 1, 0) // 1 month from now
	salaryMin := int64(5000000)
	salaryMax := int64(10000000)

	return &JobFixture{
		ID:              1,
		CompanyID:       1,
		Title:           "Software Engineer",
		Description:     "We are looking for a talented software engineer",
		Requirements:    "Bachelor's degree in Computer Science or related field",
		Location:        "Jakarta, Indonesia",
		EmploymentType:  "full_time",
		ExperienceLevel: "mid_level",
		SalaryMin:       &salaryMin,
		SalaryMax:       &salaryMax,
		Skills:          []string{"Go", "PostgreSQL", "Docker"},
		Status:          "open",
		PostedAt:        now,
		ClosingDate:     &closingDate,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// NewJobWithCustomStatus creates a job with custom status
func NewJobWithCustomStatus(status string) *JobFixture {
	job := NewJob()
	job.Status = status
	return job
}

// ApplicationFixture represents a test application
type ApplicationFixture struct {
	ID          int64
	JobID       int64
	UserID      int64
	Status      string
	CoverLetter string
	ResumeURL   string
	AppliedAt   time.Time
	UpdatedAt   time.Time
}

// NewApplication creates an application fixture
func NewApplication() *ApplicationFixture {
	now := time.Now()
	return &ApplicationFixture{
		ID:          1,
		JobID:       1,
		UserID:      1,
		Status:      "pending",
		CoverLetter: "I am very interested in this position...",
		ResumeURL:   "https://storage.example.com/resumes/user1_resume.pdf",
		AppliedAt:   now,
		UpdatedAt:   now,
	}
}

// NewApplicationWithStatus creates an application with custom status
func NewApplicationWithStatus(status string) *ApplicationFixture {
	app := NewApplication()
	app.Status = status
	return app
}

// NotificationFixture represents a test notification
type NotificationFixture struct {
	ID        int64
	UserID    int64
	Type      string
	Title     string
	Message   string
	Data      map[string]interface{}
	IsRead    bool
	ReadAt    *time.Time
	CreatedAt time.Time
}

// NewNotification creates a notification fixture
func NewNotification() *NotificationFixture {
	now := time.Now()
	return &NotificationFixture{
		ID:      1,
		UserID:  1,
		Type:    "application_update",
		Title:   "Application Status Update",
		Message: "Your application status has been updated",
		Data: map[string]interface{}{
			"application_id": 1,
			"new_status":     "reviewed",
		},
		IsRead:    false,
		CreatedAt: now,
	}
}

// DeviceTokenFixture represents a test device token
type DeviceTokenFixture struct {
	ID           int64
	UserID       int64
	Token        string
	Platform     string
	DeviceModel  string
	IsActive     bool
	LastUsedAt   time.Time
	FailureCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewDeviceToken creates a device token fixture
func NewDeviceToken() *DeviceTokenFixture {
	now := time.Now()
	return &DeviceTokenFixture{
		ID:           1,
		UserID:       1,
		Token:        "fcm_token_" + uuid.New().String(),
		Platform:     "android",
		DeviceModel:  "Samsung Galaxy S21",
		IsActive:     true,
		LastUsedAt:   now,
		FailureCount: 0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NewIOSDeviceToken creates an iOS device token fixture
func NewIOSDeviceToken() *DeviceTokenFixture {
	token := NewDeviceToken()
	token.Platform = "ios"
	token.DeviceModel = "iPhone 13"
	return token
}

// TestSkills provides a list of test skills
var TestSkills = []string{
	"Go",
	"Python",
	"JavaScript",
	"TypeScript",
	"React",
	"Node.js",
	"PostgreSQL",
	"MySQL",
	"MongoDB",
	"Docker",
	"Kubernetes",
	"AWS",
	"Azure",
	"GCP",
	"Git",
	"CI/CD",
	"REST API",
	"GraphQL",
	"Microservices",
	"Agile",
}

// TestIndustries provides a list of test industries
var TestIndustries = []string{
	"Technology",
	"Finance",
	"Healthcare",
	"Education",
	"E-commerce",
	"Manufacturing",
	"Retail",
	"Consulting",
	"Media",
	"Real Estate",
}

// TestLocations provides a list of test locations
var TestLocations = []string{
	"Jakarta, Indonesia",
	"Bandung, Indonesia",
	"Surabaya, Indonesia",
	"Yogyakarta, Indonesia",
	"Bali, Indonesia",
}

// TestEmploymentTypes provides employment types
var TestEmploymentTypes = []string{
	"full_time",
	"part_time",
	"contract",
	"freelance",
	"internship",
}

// TestExperienceLevels provides experience levels
var TestExperienceLevels = []string{
	"entry_level",
	"junior",
	"mid_level",
	"senior",
	"lead",
	"manager",
}

// TestApplicationStatuses provides application statuses
var TestApplicationStatuses = []string{
	"pending",
	"reviewing",
	"shortlisted",
	"interview",
	"offered",
	"rejected",
	"withdrawn",
	"accepted",
}

// TestJobStatuses provides job statuses
var TestJobStatuses = []string{
	"draft",
	"open",
	"closed",
	"filled",
	"cancelled",
}
