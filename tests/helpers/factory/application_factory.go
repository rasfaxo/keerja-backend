package factory

import (
	"time"

	"github.com/google/uuid"
)

// ApplicationFactory creates test application instances
type ApplicationFactory struct {
	ID          string
	JobID       string
	UserID      string
	Status      string
	CoverLetter string
	ResumeURL   string
	AppliedAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewApplicationFactory creates a new application factory with default values
func NewApplicationFactory(jobID, userID string) *ApplicationFactory {
	now := time.Now().UTC()

	return &ApplicationFactory{
		ID:          uuid.New().String(),
		JobID:       jobID,
		UserID:      userID,
		Status:      "pending",
		CoverLetter: "I am very interested in this position and believe my skills are a great fit.",
		ResumeURL:   "https://storage.example.com/resumes/" + uuid.New().String() + ".pdf",
		AppliedAt:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// WithID sets custom ID
func (f *ApplicationFactory) WithID(id string) *ApplicationFactory {
	f.ID = id
	return f
}

// WithJobID sets custom job ID
func (f *ApplicationFactory) WithJobID(jobID string) *ApplicationFactory {
	f.JobID = jobID
	return f
}

// WithUserID sets custom user ID
func (f *ApplicationFactory) WithUserID(userID string) *ApplicationFactory {
	f.UserID = userID
	return f
}

// WithStatus sets custom status
func (f *ApplicationFactory) WithStatus(status string) *ApplicationFactory {
	f.Status = status
	return f
}

// WithCoverLetter sets custom cover letter
func (f *ApplicationFactory) WithCoverLetter(coverLetter string) *ApplicationFactory {
	f.CoverLetter = coverLetter
	return f
}

// WithResumeURL sets custom resume URL
func (f *ApplicationFactory) WithResumeURL(resumeURL string) *ApplicationFactory {
	f.ResumeURL = resumeURL
	return f
}

// AsPending sets status to pending
func (f *ApplicationFactory) AsPending() *ApplicationFactory {
	f.Status = "pending"
	return f
}

// AsReviewing sets status to reviewing
func (f *ApplicationFactory) AsReviewing() *ApplicationFactory {
	f.Status = "reviewing"
	return f
}

// AsShortlisted sets status to shortlisted
func (f *ApplicationFactory) AsShortlisted() *ApplicationFactory {
	f.Status = "shortlisted"
	return f
}

// AsInterviewing sets status to interviewing
func (f *ApplicationFactory) AsInterviewing() *ApplicationFactory {
	f.Status = "interviewing"
	return f
}

// AsAccepted sets status to accepted
func (f *ApplicationFactory) AsAccepted() *ApplicationFactory {
	f.Status = "accepted"
	return f
}

// AsRejected sets status to rejected
func (f *ApplicationFactory) AsRejected() *ApplicationFactory {
	f.Status = "rejected"
	return f
}

// AsWithdrawn sets status to withdrawn
func (f *ApplicationFactory) AsWithdrawn() *ApplicationFactory {
	f.Status = "withdrawn"
	return f
}

// WithoutResume removes resume URL
func (f *ApplicationFactory) WithoutResume() *ApplicationFactory {
	f.ResumeURL = ""
	return f
}

// Build returns the application factory as map for testing
func (f *ApplicationFactory) Build() map[string]interface{} {
	return map[string]interface{}{
		"id":           f.ID,
		"job_id":       f.JobID,
		"user_id":      f.UserID,
		"status":       f.Status,
		"cover_letter": f.CoverLetter,
		"resume_url":   f.ResumeURL,
		"applied_at":   f.AppliedAt,
		"created_at":   f.CreatedAt,
		"updated_at":   f.UpdatedAt,
	}
}

// BuildMultiple creates multiple applications
func BuildMultipleApplications(count int, jobID, userID string) []*ApplicationFactory {
	applications := make([]*ApplicationFactory, count)
	for i := 0; i < count; i++ {
		applications[i] = NewApplicationFactory(jobID, userID)
	}
	return applications
}
