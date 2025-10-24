package email

import "context"

// EmailService defines the interface for email operations
type EmailService interface {
	// Send sends an email
	SendEmail(ctx context.Context, to, subject, body string) error

	// SendTemplateEmail sends an email using a template
	SendTemplateEmail(ctx context.Context, to, template string, data map[string]interface{}) error

	// SendVerificationEmail sends account verification email
	SendVerificationEmail(ctx context.Context, to, token string) error

	// SendPasswordResetEmail sends password reset email
	SendPasswordResetEmail(ctx context.Context, to, token string) error

	// SendWelcomeEmail sends welcome email to new users
	SendWelcomeEmail(ctx context.Context, to, name string) error

	// SendJobApplicationEmail sends job application confirmation
	SendJobApplicationEmail(ctx context.Context, to, jobTitle, companyName string) error

	// SendInterviewInvitationEmail sends interview invitation
	SendInterviewInvitationEmail(ctx context.Context, to, jobTitle string, interviewDate string) error

	// SendJobStatusUpdateEmail sends job status update notification
	SendJobStatusUpdateEmail(ctx context.Context, to, jobTitle, status string) error

	// SendBulkEmail sends email to multiple recipients
	SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error

	// GetEmailLog retrieves email log by ID
	GetEmailLog(ctx context.Context, id int64) (*EmailLog, error)

	// GetEmailLogs retrieves email logs with pagination
	GetEmailLogs(ctx context.Context, filter EmailFilter, page, limit int) ([]EmailLog, int64, error)

	// RetryFailedEmail retries sending a failed email
	RetryFailedEmail(ctx context.Context, logID int64) error

	// GetFailedEmails retrieves all failed emails
	GetFailedEmails(ctx context.Context, page, limit int) ([]EmailLog, int64, error)

	// SendOTPEmail sends OTP code via email
	SendOTPEmail(ctx context.Context, to, code, purpose string) error

	// SendOTPRegistrationEmail sends OTP for email verification during registration
	SendOTPRegistrationEmail(ctx context.Context, to, name, code string) error
}

// EmailFilter defines filters for email logs
type EmailFilter struct {
	Recipient string
	Status    string
	Template  string
	DateFrom  *string
	DateTo    *string
}

// EmailRepository defines the interface for email data operations
type EmailRepository interface {
	// Create creates a new email log
	Create(ctx context.Context, log *EmailLog) error

	// FindByID finds email log by ID
	FindByID(ctx context.Context, id int64) (*EmailLog, error)

	// Update updates email log
	Update(ctx context.Context, log *EmailLog) error

	// List retrieves email logs with pagination
	List(ctx context.Context, filter EmailFilter, page, limit int) ([]EmailLog, int64, error)

	// GetFailedEmails retrieves failed emails
	GetFailedEmails(ctx context.Context, page, limit int) ([]EmailLog, int64, error)

	// GetPendingEmails retrieves pending emails
	GetPendingEmails(ctx context.Context, limit int) ([]EmailLog, error)

	// Delete deletes email log
	Delete(ctx context.Context, id int64) error
}
