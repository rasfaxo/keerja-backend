package service

import (
	"context"
	"fmt"
	"keerja-backend/internal/domain/job"

	"gorm.io/gorm"
)

// AdminJobService handles admin operations on jobs
type adminJobService struct {
	jobRepo job.Repository
}

// NewAdminJobService creates a new AdminJobService
func NewAdminJobService(jobRepo job.Repository) job.AdminJobService {
	return &adminJobService{
		jobRepo: jobRepo,
	}
}

// ApproveJob approves a job posting
func (s *adminJobService) ApproveJob(ctx context.Context, jobID int64) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("job not found")
		}
		return err
	}

	if job.Status != "pending" {
		return fmt.Errorf("job is not in pending status")
	}

	job.Status = "approved"
	if err := s.jobRepo.Update(ctx, job); err != nil {
		return fmt.Errorf("failed to update job status: %v", err)
	}

	return nil
}

// RejectJob rejects a job posting with a reason
func (s *adminJobService) RejectJob(ctx context.Context, jobID int64, reason string) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("job not found")
		}
		return err
	}

	if job.Status != "pending" {
		return fmt.Errorf("job is not in pending status")
	}

	job.Status = "rejected"
	job.RejectionReason = &reason
	if err := s.jobRepo.Update(ctx, job); err != nil {
		return fmt.Errorf("failed to update job status: %v", err)
	}

	return nil
}