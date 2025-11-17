package job

import (
	"context"
)

// Repository defines job persistence operations
type Repository interface {
	// GetByID retrieves a job by its ID
	GetByID(ctx context.Context, id int64) (*Job, error)

	// Update updates an existing job
	Update(ctx context.Context, job *Job) error
}

// Job represents a job posting entity
type Job struct {
	ID              int64   `json:"id" gorm:"primaryKey"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Status          string  `json:"status"` // pending, approved, rejected
	RejectionReason *string `json:"rejection_reason,omitempty"`
}
