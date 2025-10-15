package service

import (
	"context"
	"errors"
	"fmt"

	"keerja-backend/internal/domain/email"
)

// emailService implements email.EmailService interface
type emailService struct {
	emailRepo email.EmailRepository
	// In production, you would inject actual email provider (SMTP, SendGrid, AWS SES, etc.)
	// provider EmailProvider
}

// NewEmailService creates a new email service instance
func NewEmailService(emailRepo email.EmailRepository) email.EmailService {
	return &emailService{
		emailRepo: emailRepo,
	}
}

// SendEmail sends a plain email
func (s *emailService) SendEmail(ctx context.Context, to, subject, body string) error {
	// Create email log
	log := &email.EmailLog{
		Recipient: to,
		Subject:   subject,
		Body:      body,
		Status:    "pending",
		Provider:  "smtp", // Default provider
	}

	// Save to database
	if err := s.emailRepo.Create(ctx, log); err != nil {
		return fmt.Errorf("failed to create email log: %w", err)
	}

	// Send email using provider (simulated here)
	if err := s.sendEmailViaProvider(ctx, to, subject, body); err != nil {
		log.MarkAsFailed(err.Error())
		s.emailRepo.Update(ctx, log)
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Mark as sent
	log.MarkAsSent()
	if err := s.emailRepo.Update(ctx, log); err != nil {
		return fmt.Errorf("failed to update email log: %w", err)
	}

	return nil
}

// SendTemplateEmail sends an email using a template
func (s *emailService) SendTemplateEmail(ctx context.Context, to, template string, data map[string]interface{}) error {
	// Load template and render with data
	subject, body, err := s.renderTemplate(template, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Create email log
	log := &email.EmailLog{
		Recipient: to,
		Subject:   subject,
		Body:      body,
		Template:  template,
		Status:    "pending",
		Provider:  "smtp",
	}

	// Save to database
	if err := s.emailRepo.Create(ctx, log); err != nil {
		return fmt.Errorf("failed to create email log: %w", err)
	}

	// Send email
	if err := s.sendEmailViaProvider(ctx, to, subject, body); err != nil {
		log.MarkAsFailed(err.Error())
		s.emailRepo.Update(ctx, log)
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Mark as sent
	log.MarkAsSent()
	if err := s.emailRepo.Update(ctx, log); err != nil {
		return fmt.Errorf("failed to update email log: %w", err)
	}

	return nil
}

// SendVerificationEmail sends account verification email
func (s *emailService) SendVerificationEmail(ctx context.Context, to, token string) error {
	subject := "Verify Your Keerja Account"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Welcome to Keerja!</h2>
			<p>Please verify your email address by clicking the link below:</p>
			<p><a href="https://keerja.com/verify?token=%s">Verify Email</a></p>
			<p>This link will expire in 24 hours.</p>
			<p>If you didn't create an account, please ignore this email.</p>
		</body>
		</html>
	`, token)

	return s.SendEmail(ctx, to, subject, body)
}

// SendPasswordResetEmail sends password reset email
func (s *emailService) SendPasswordResetEmail(ctx context.Context, to, token string) error {
	subject := "Reset Your Keerja Password"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Password Reset Request</h2>
			<p>You requested to reset your password. Click the link below to proceed:</p>
			<p><a href="https://keerja.com/reset-password?token=%s">Reset Password</a></p>
			<p>This link will expire in 1 hour.</p>
			<p>If you didn't request this, please ignore this email.</p>
		</body>
		</html>
	`, token)

	return s.SendEmail(ctx, to, subject, body)
}

// SendWelcomeEmail sends welcome email to new users
func (s *emailService) SendWelcomeEmail(ctx context.Context, to, name string) error {
	subject := "Welcome to Keerja!"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Welcome to Keerja, %s!</h2>
			<p>We're excited to have you on board. Here's what you can do:</p>
			<ul>
				<li>Complete your profile</li>
				<li>Upload your resume</li>
				<li>Start applying for jobs</li>
				<li>Get personalized job recommendations</li>
			</ul>
			<p><a href="https://keerja.com/dashboard">Go to Dashboard</a></p>
		</body>
		</html>
	`, name)

	return s.SendEmail(ctx, to, subject, body)
}

// SendJobApplicationEmail sends job application confirmation
func (s *emailService) SendJobApplicationEmail(ctx context.Context, to, jobTitle, companyName string) error {
	subject := fmt.Sprintf("Application Received: %s at %s", jobTitle, companyName)
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Application Received</h2>
			<p>Your application for <strong>%s</strong> at <strong>%s</strong> has been successfully submitted.</p>
			<p>The employer will review your application and contact you if they're interested.</p>
			<p><a href="https://keerja.com/applications">View Your Applications</a></p>
		</body>
		</html>
	`, jobTitle, companyName)

	return s.SendEmail(ctx, to, subject, body)
}

// SendInterviewInvitationEmail sends interview invitation
func (s *emailService) SendInterviewInvitationEmail(ctx context.Context, to, jobTitle string, interviewDate string) error {
	subject := fmt.Sprintf("Interview Invitation: %s", jobTitle)
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Interview Invitation</h2>
			<p>Congratulations! You've been invited for an interview for the position of <strong>%s</strong>.</p>
			<p><strong>Interview Date:</strong> %s</p>
			<p>Please confirm your attendance and prepare accordingly.</p>
			<p><a href="https://keerja.com/interviews">View Interview Details</a></p>
		</body>
		</html>
	`, jobTitle, interviewDate)

	return s.SendEmail(ctx, to, subject, body)
}

// SendJobStatusUpdateEmail sends job status update notification
func (s *emailService) SendJobStatusUpdateEmail(ctx context.Context, to, jobTitle, status string) error {
	subject := fmt.Sprintf("Application Status Update: %s", jobTitle)

	statusMessage := map[string]string{
		"screening":   "Your application is being reviewed by the employer.",
		"shortlisted": "Congratulations! You've been shortlisted for an interview.",
		"interview":   "You've been invited for an interview.",
		"offered":     "Congratulations! You've received a job offer.",
		"hired":       "Congratulations! You've been hired for this position.",
		"rejected":    "Unfortunately, your application was not successful this time.",
	}

	message := statusMessage[status]
	if message == "" {
		message = "Your application status has been updated."
	}

	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Application Status Update</h2>
			<p>Your application for <strong>%s</strong> has been updated.</p>
			<p><strong>New Status:</strong> %s</p>
			<p>%s</p>
			<p><a href="https://keerja.com/applications">View Application Details</a></p>
		</body>
		</html>
	`, jobTitle, status, message)

	return s.SendEmail(ctx, to, subject, body)
}

// SendBulkEmail sends email to multiple recipients
func (s *emailService) SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error {
	var errors []error

	for _, recipient := range recipients {
		if err := s.SendEmail(ctx, recipient, subject, body); err != nil {
			errors = append(errors, fmt.Errorf("failed to send to %s: %w", recipient, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to send %d out of %d emails", len(errors), len(recipients))
	}

	return nil
}

// GetEmailLog retrieves email log by ID
func (s *emailService) GetEmailLog(ctx context.Context, id int64) (*email.EmailLog, error) {
	return s.emailRepo.FindByID(ctx, id)
}

// GetEmailLogs retrieves email logs with pagination
func (s *emailService) GetEmailLogs(ctx context.Context, filter email.EmailFilter, page, limit int) ([]email.EmailLog, int64, error) {
	return s.emailRepo.List(ctx, filter, page, limit)
}

// RetryFailedEmail retries sending a failed email
func (s *emailService) RetryFailedEmail(ctx context.Context, logID int64) error {
	// Get email log
	log, err := s.emailRepo.FindByID(ctx, logID)
	if err != nil {
		return fmt.Errorf("email log not found: %w", err)
	}

	// Check if can retry
	if !log.CanRetry() {
		return errors.New("email cannot be retried (max retries reached or not in failed status)")
	}

	// Reset status to pending
	log.Status = "pending"

	// Attempt to send
	if err := s.sendEmailViaProvider(ctx, log.Recipient, log.Subject, log.Body); err != nil {
		log.MarkAsFailed(err.Error())
		s.emailRepo.Update(ctx, log)
		return fmt.Errorf("failed to resend email: %w", err)
	}

	// Mark as sent
	log.MarkAsSent()
	if err := s.emailRepo.Update(ctx, log); err != nil {
		return fmt.Errorf("failed to update email log: %w", err)
	}

	return nil
}

// GetFailedEmails retrieves all failed emails
func (s *emailService) GetFailedEmails(ctx context.Context, page, limit int) ([]email.EmailLog, int64, error) {
	return s.emailRepo.GetFailedEmails(ctx, page, limit)
}

// ===== Helper Methods =====

// sendEmailViaProvider sends email using the actual provider
// This is a placeholder - in production, integrate with SMTP, SendGrid, AWS SES, etc.
func (s *emailService) sendEmailViaProvider(ctx context.Context, to, subject, body string) error {
	// TODO: Implement actual email sending logic
	// For now, simulate successful sending and log to console

	fmt.Printf("\n====== EMAIL SENT ======\n")
	fmt.Printf("To: %s\n", to)
	fmt.Printf("Subject: %s\n", subject)
	fmt.Printf("========================\n\n")

	// Example with SMTP:
	// return s.smtpClient.Send(to, subject, body)

	// Example with SendGrid:
	// return s.sendgridClient.Send(to, subject, body)

	// Simulate success
	return nil
}

// renderTemplate renders email template with data
func (s *emailService) renderTemplate(template string, data map[string]interface{}) (string, string, error) {
	// TODO: Implement template rendering
	// Load template from file or database
	// Render with data using template engine (html/template, etc.)

	// For now, return placeholder
	subject := fmt.Sprintf("Email from template: %s", template)
	body := fmt.Sprintf("Template: %s, Data: %v", template, data)

	return subject, body, nil
}
