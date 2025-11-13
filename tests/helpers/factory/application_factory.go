package factory

import (
	"fmt"
	"time"
)

// ApplicationFactory provides methods to create test applications
type ApplicationFactory struct {
	sequence int
}

// NewApplicationFactory creates a new application factory
func NewApplicationFactory() *ApplicationFactory {
	return &ApplicationFactory{
		sequence: 0,
	}
}

// ApplicationBuilder provides a fluent interface for building applications
type ApplicationBuilder struct {
	ID            int64
	JobID         int64
	UserID        int64
	Status        string
	CoverLetter   string
	ResumeURL     *string
	PortfolioURL  *string
	LinkedInURL   *string
	GitHubURL     *string
	AppliedAt     time.Time
	ReviewedAt    *time.Time
	ReviewedBy    *int64
	Notes         string
	Rating        *int
	IsShortlisted bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Build builds the application
func (ab *ApplicationBuilder) Build() *ApplicationBuilder {
	// Set default timestamps if not set
	if ab.CreatedAt.IsZero() {
		ab.CreatedAt = time.Now()
	}
	if ab.UpdatedAt.IsZero() {
		ab.UpdatedAt = time.Now()
	}
	if ab.AppliedAt.IsZero() {
		ab.AppliedAt = time.Now()
	}

	// Set default status
	if ab.Status == "" {
		ab.Status = "pending"
	}

	return ab
}

// WithID sets the application ID
func (ab *ApplicationBuilder) WithID(id int64) *ApplicationBuilder {
	ab.ID = id
	return ab
}

// WithJobID sets the job ID
func (ab *ApplicationBuilder) WithJobID(jobID int64) *ApplicationBuilder {
	ab.JobID = jobID
	return ab
}

// WithUserID sets the user ID
func (ab *ApplicationBuilder) WithUserID(userID int64) *ApplicationBuilder {
	ab.UserID = userID
	return ab
}

// WithStatus sets the application status
func (ab *ApplicationBuilder) WithStatus(status string) *ApplicationBuilder {
	ab.Status = status
	return ab
}

// Pending sets the status as pending
func (ab *ApplicationBuilder) Pending() *ApplicationBuilder {
	ab.Status = "pending"
	return ab
}

// Reviewing sets the status as reviewing
func (ab *ApplicationBuilder) Reviewing() *ApplicationBuilder {
	ab.Status = "reviewing"
	now := time.Now()
	ab.Status = "reviewing"
	ab.ReviewedAt = &now
	return ab
}

// Shortlisted sets the status as shortlisted
func (ab *ApplicationBuilder) Shortlisted() *ApplicationBuilder {
	now := time.Now()
	ab.Status = "shortlisted"
	ab.IsShortlisted = true
	ab.ReviewedAt = &now
	return ab
}

// Interview sets the status as interview
func (ab *ApplicationBuilder) Interview() *ApplicationBuilder {
	now := time.Now()
	ab.Status = "interview"
	ab.ReviewedAt = &now
	return ab
}

// Offered sets the status as offered
func (ab *ApplicationBuilder) Offered() *ApplicationBuilder {
	now := time.Now()
	ab.Status = "offered"
	ab.ReviewedAt = &now
	return ab
}

// Rejected sets the status as rejected
func (ab *ApplicationBuilder) Rejected() *ApplicationBuilder {
	now := time.Now()
	ab.Status = "rejected"
	ab.ReviewedAt = &now
	return ab
}

// Withdrawn sets the status as withdrawn
func (ab *ApplicationBuilder) Withdrawn() *ApplicationBuilder {
	ab.Status = "withdrawn"
	return ab
}

// Accepted sets the status as accepted
func (ab *ApplicationBuilder) Accepted() *ApplicationBuilder {
	now := time.Now()
	ab.Status = "accepted"
	ab.ReviewedAt = &now
	return ab
}

// WithCoverLetter sets the cover letter
func (ab *ApplicationBuilder) WithCoverLetter(letter string) *ApplicationBuilder {
	ab.CoverLetter = letter
	return ab
}

// WithResume sets the resume URL
func (ab *ApplicationBuilder) WithResume(url string) *ApplicationBuilder {
	ab.ResumeURL = &url
	return ab
}

// WithPortfolio sets the portfolio URL
func (ab *ApplicationBuilder) WithPortfolio(url string) *ApplicationBuilder {
	ab.PortfolioURL = &url
	return ab
}

// WithLinkedIn sets the LinkedIn URL
func (ab *ApplicationBuilder) WithLinkedIn(url string) *ApplicationBuilder {
	ab.LinkedInURL = &url
	return ab
}

// WithGitHub sets the GitHub URL
func (ab *ApplicationBuilder) WithGitHub(url string) *ApplicationBuilder {
	ab.GitHubURL = &url
	return ab
}

// WithReviewedBy sets who reviewed the application
func (ab *ApplicationBuilder) WithReviewedBy(userID int64) *ApplicationBuilder {
	now := time.Now()
	ab.ReviewedBy = &userID
	ab.ReviewedAt = &now
	return ab
}

// WithNotes sets the review notes
func (ab *ApplicationBuilder) WithNotes(notes string) *ApplicationBuilder {
	ab.Notes = notes
	return ab
}

// WithRating sets the rating
func (ab *ApplicationBuilder) WithRating(rating int) *ApplicationBuilder {
	ab.Rating = &rating
	return ab
}

// WithAppliedAt sets the applied at time
func (ab *ApplicationBuilder) WithAppliedAt(t time.Time) *ApplicationBuilder {
	ab.AppliedAt = t
	return ab
}

// WithCreatedAt sets the created at time
func (ab *ApplicationBuilder) WithCreatedAt(t time.Time) *ApplicationBuilder {
	ab.CreatedAt = t
	return ab
}

// WithUpdatedAt sets the updated at time
func (ab *ApplicationBuilder) WithUpdatedAt(t time.Time) *ApplicationBuilder {
	ab.UpdatedAt = t
	return ab
}

// CreateApplication creates an application builder with default values
func (f *ApplicationFactory) CreateApplication() *ApplicationBuilder {
	f.sequence++
	resumeURL := fmt.Sprintf("https://storage.example.com/resumes/resume_%d.pdf", f.sequence)

	return &ApplicationBuilder{
		ID:          int64(f.sequence),
		JobID:       1,
		UserID:      1,
		Status:      "pending",
		CoverLetter: "I am very interested in this position and believe my skills and experience make me a great fit.",
		ResumeURL:   &resumeURL,
	}
}

// CreateApplicationForJob creates an application for a specific job
func (f *ApplicationFactory) CreateApplicationForJob(jobID int64) *ApplicationBuilder {
	return f.CreateApplication().WithJobID(jobID)
}

// CreateApplicationByUser creates an application by a specific user
func (f *ApplicationFactory) CreateApplicationByUser(userID int64) *ApplicationBuilder {
	return f.CreateApplication().WithUserID(userID)
}

// CreateApplicationForJobByUser creates an application for a job by a user
func (f *ApplicationFactory) CreateApplicationForJobByUser(jobID, userID int64) *ApplicationBuilder {
	return f.CreateApplication().WithJobID(jobID).WithUserID(userID)
}

// CreatePendingApplication creates a pending application
func (f *ApplicationFactory) CreatePendingApplication() *ApplicationBuilder {
	return f.CreateApplication().Pending()
}

// CreateReviewingApplication creates a reviewing application
func (f *ApplicationFactory) CreateReviewingApplication() *ApplicationBuilder {
	return f.CreateApplication().Reviewing()
}

// CreateShortlistedApplication creates a shortlisted application
func (f *ApplicationFactory) CreateShortlistedApplication() *ApplicationBuilder {
	return f.CreateApplication().Shortlisted()
}

// CreateRejectedApplication creates a rejected application
func (f *ApplicationFactory) CreateRejectedApplication() *ApplicationBuilder {
	return f.CreateApplication().Rejected()
}

// CreateAcceptedApplication creates an accepted application
func (f *ApplicationFactory) CreateAcceptedApplication() *ApplicationBuilder {
	return f.CreateApplication().Accepted()
}

// CreateMultipleApplications creates multiple applications
func (f *ApplicationFactory) CreateMultipleApplications(count int) []*ApplicationBuilder {
	applications := make([]*ApplicationBuilder, count)
	for i := 0; i < count; i++ {
		applications[i] = f.CreateApplication()
	}
	return applications
}

// CreateMultipleApplicationsForJob creates multiple applications for a job
func (f *ApplicationFactory) CreateMultipleApplicationsForJob(jobID int64, count int) []*ApplicationBuilder {
	applications := make([]*ApplicationBuilder, count)
	for i := 0; i < count; i++ {
		applications[i] = f.CreateApplicationForJob(jobID)
	}
	return applications
}

// CreateMultipleApplicationsByUser creates multiple applications by a user
func (f *ApplicationFactory) CreateMultipleApplicationsByUser(userID int64, count int) []*ApplicationBuilder {
	applications := make([]*ApplicationBuilder, count)
	for i := 0; i < count; i++ {
		applications[i] = f.CreateApplicationByUser(userID)
	}
	return applications
}

// ApplicationStage represents an application stage
type ApplicationStage struct {
	ID            int64
	ApplicationID int64
	Stage         string
	Status        string
	ScheduledAt   *time.Time
	CompletedAt   *time.Time
	Notes         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ApplicationDocument represents an application document
type ApplicationDocument struct {
	ID            int64
	ApplicationID int64
	DocumentType  string
	DocumentURL   string
	FileName      string
	FileSize      int64
	UploadedAt    time.Time
	CreatedAt     time.Time
}

// ApplicationNote represents an application note
type ApplicationNote struct {
	ID            int64
	ApplicationID int64
	UserID        int64
	Note          string
	IsInternal    bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// WithStage adds a stage to the application
func (ab *ApplicationBuilder) WithStage(stage ApplicationStage) *ApplicationBuilder {
	// This would be implemented based on your domain model
	return ab
}

// WithDocument adds a document to the application
func (ab *ApplicationBuilder) WithDocument(doc ApplicationDocument) *ApplicationBuilder {
	// This would be implemented based on your domain model
	return ab
}

// WithNote adds a note to the application
func (ab *ApplicationBuilder) WithNote(note ApplicationNote) *ApplicationBuilder {
	// This would be implemented based on your domain model
	return ab
}
