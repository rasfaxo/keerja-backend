package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"keerja-backend/internal/domain/job"
)

// adminJobService implements admin.AdminJobService interface for job moderation
type adminJobService struct {
	jobRepo job.JobRepository
}

// NewAdminJobService creates a new admin job service instance
func NewAdminJobService(jobRepo job.JobRepository) AdminJobService {
	return &adminJobService{
		jobRepo: jobRepo,
	}
}

// AdminJobService interface (implementation of admin.AdminJobService)
type AdminJobService interface {
	// Job approval/rejection (moderation)
	ApproveJob(ctx context.Context, jobID int64) (interface{}, error)
	RejectJob(ctx context.Context, jobID int64, reason string) (interface{}, error)

	// Job list for approval
	GetPendingJobs(ctx context.Context, page, limit int) ([]interface{}, int64, error)
	GetJobsForReview(ctx context.Context, status string, page, limit int) ([]interface{}, int64, error)
}

// ApproveJob approves a pending job posting (admin only)
// Changes status from pending_review to published
func (s *adminJobService) ApproveJob(ctx context.Context, jobID int64) (interface{}, error) {
	// Get job
	j, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	// Verify job is pending review
	if j.Status != "pending_review" {
		return nil, fmt.Errorf("only jobs with pending_review status can be approved (current status: %s)", j.Status)
	}

	// Update status to published
	if err := s.jobRepo.UpdateStatus(ctx, jobID, "published"); err != nil {
		return nil, fmt.Errorf("failed to update job status: %w", err)
	}

	// Set published_at timestamp
	now := time.Now()
	j.PublishedAt = &now
	if err := s.jobRepo.Update(ctx, j); err != nil {
		return nil, fmt.Errorf("failed to update published_at: %w", err)
	}

	// Reload and return updated job
	return s.jobRepo.FindByID(ctx, jobID)
}

// RejectJob rejects a pending job posting (admin only)
// Changes status from pending_review back to draft so employer can fix and resubmit
func (s *adminJobService) RejectJob(ctx context.Context, jobID int64, reason string) (interface{}, error) {
	// Get job
	j, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	// Verify job is pending review
	if j.Status != "pending_review" {
		return nil, fmt.Errorf("only jobs with pending_review status can be rejected (current status: %s)", j.Status)
	}

	// Update status to draft (so employer can fix and resubmit)
	if err := s.jobRepo.UpdateStatus(ctx, jobID, "draft"); err != nil {
		return nil, fmt.Errorf("failed to update job status: %w", err)
	}

	// TODO: Store rejection reason in audit table or notification system
	// For now, log it for reference
	fmt.Printf("Job %d rejected with reason: %s\n", jobID, reason)

	// Reload and return updated job
	return s.jobRepo.FindByID(ctx, jobID)
}

// GetPendingJobs retrieves all jobs pending review (status = pending_review)
func (s *adminJobService) GetPendingJobs(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Build filter for pending_review status
	filter := job.JobFilter{
		Status: "pending_review",
	}

	jobs, total, err := s.jobRepo.List(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get pending jobs: %w", err)
	}

	// Convert to interface{} slice
	result := make([]interface{}, len(jobs))
	for i, j := range jobs {
		result[i] = j
	}

	return result, total, nil
}

// GetJobsForReview retrieves jobs for review with specific status
func (s *adminJobService) GetJobsForReview(ctx context.Context, status string, page, limit int) ([]interface{}, int64, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Validate status - only allow review-related statuses
	validStatuses := map[string]bool{
		"pending_review": true,
		"published":      true,
		"draft":          true,
		"suspended":      true,
	}

	if !validStatuses[status] {
		return nil, 0, fmt.Errorf("invalid status for review: %s", status)
	}

	// Build filter
	filter := job.JobFilter{
		Status: status,
	}

	jobs, total, err := s.jobRepo.List(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get jobs for review: %w", err)
	}

	// Convert to interface{} slice
	result := make([]interface{}, len(jobs))
	for i, j := range jobs {
		result[i] = j
	}

	// Calculate pagination info
	totalPages := int64(math.Ceil(float64(total) / float64(limit)))
	fmt.Printf("Retrieved %d jobs with status '%s' (page %d of %d)\n", len(jobs), status, page, totalPages)

	return result, total, nil
}
