package job

import "context"

// AdminJobService defines admin operations on jobs
type AdminJobService interface {
	// ApproveJob approves a job posting
	ApproveJob(ctx context.Context, jobID int64) error

	// RejectJob rejects a job posting with a reason
	RejectJob(ctx context.Context, jobID int64, reason string) error
}
