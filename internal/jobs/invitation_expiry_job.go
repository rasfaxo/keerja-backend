package jobs

import (
	"context"
	"fmt"

	"keerja-backend/internal/domain/company"
)

// InvitationExpiryJob expires old company invitations
type InvitationExpiryJob struct {
	companyService company.CompanyService
}

// NewInvitationExpiryJob creates a new invitation expiry job
func NewInvitationExpiryJob(companyService company.CompanyService) *InvitationExpiryJob {
	return &InvitationExpiryJob{
		companyService: companyService,
	}
}

// Name returns the job name
func (j *InvitationExpiryJob) Name() string {
	return "invitation_expiry"
}

// Schedule returns the cron schedule (every hour)
func (j *InvitationExpiryJob) Schedule() string {
	return "0 0 * * * *" // Every hour at minute 0
}

// Run executes the job
func (j *InvitationExpiryJob) Run(ctx context.Context) error {
	fmt.Println("Running invitation expiry job...")

	// Expire old invitations
	count, err := j.companyService.ExpireOldInvitations(ctx)
	if err != nil {
		return fmt.Errorf("failed to expire invitations: %w", err)
	}

	if count > 0 {
		fmt.Printf("Expired %d old invitations\n", count)
	} else {
		fmt.Println("No invitations to expire")
	}

	return nil
}
