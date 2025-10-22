package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"
	"time"

	"keerja-backend/internal/config"
	"keerja-backend/internal/domain/email"

	"gopkg.in/gomail.v2"
)

// emailService implements email.EmailService interface
type emailService struct {
	emailRepo email.EmailRepository
	config    *config.Config
	dialer    *gomail.Dialer
}

// NewEmailService creates a new email service instance
func NewEmailService(emailRepo email.EmailRepository, cfg *config.Config) email.EmailService {
	// Parse SMTP port
	smtpPort, err := strconv.Atoi(cfg.SMTPPort)
	if err != nil {
		smtpPort = 587 // Default port
	}

	// Create SMTP dialer
	dialer := gomail.NewDialer(cfg.SMTPHost, smtpPort, cfg.SMTPUsername, cfg.SMTPPassword)

	// For MailHog or development, disable SSL/TLS verification
	if cfg.IsDevelopment() || cfg.SMTPHost == "localhost" || cfg.SMTPHost == "mailhog" {
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return &emailService{
		emailRepo: emailRepo,
		config:    cfg,
		dialer:    dialer,
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
		Provider:  "smtp",
	}

	// Save to database
	if err := s.emailRepo.Create(ctx, log); err != nil {
		return fmt.Errorf("failed to create email log: %w", err)
	}

	// Send email using SMTP
	if err := s.sendViaSMTP(to, subject, body); err != nil {
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
func (s *emailService) SendTemplateEmail(ctx context.Context, to, templateName string, data map[string]interface{}) error {
	// Convert template name to EmailTemplate type
	templateType := email.EmailTemplate(templateName)

	// Prepare template data
	templateData := s.mapToTemplateData(data)

	// Render template
	body, err := email.RenderTemplate(templateType, templateData)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Get subject
	subject := email.GetSubject(templateType)

	// Create email log
	log := &email.EmailLog{
		Recipient: to,
		Subject:   subject,
		Body:      body,
		Template:  templateName,
		Status:    "pending",
		Provider:  "smtp",
	}

	// Save to database
	if err := s.emailRepo.Create(ctx, log); err != nil {
		return fmt.Errorf("failed to create email log: %w", err)
	}

	// Debug: Log before sending
	fmt.Printf("[DEBUG] About to send email to: %s\n", to)
	fmt.Printf("[DEBUG] SMTP Host: %s:%s\n", s.config.SMTPHost, s.config.SMTPPort)

	// Send email
	if err := s.sendViaSMTP(to, subject, body); err != nil {
		fmt.Printf("[ERROR] Failed to send email: %v\n", err)
		log.MarkAsFailed(err.Error())
		s.emailRepo.Update(ctx, log)
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("[DEBUG] Email sent successfully\n")

	// Mark as sent
	log.MarkAsSent()
	if err := s.emailRepo.Update(ctx, log); err != nil {
		return fmt.Errorf("failed to update email log: %w", err)
	}

	return nil
}

// SendVerificationEmail sends account verification email
func (s *emailService) SendVerificationEmail(ctx context.Context, to, token string) error {
	verifyURL := fmt.Sprintf("%s?token=%s", s.config.VerifyEmailURL, token)

	data := map[string]interface{}{
		"Name":         to,
		"Token":        token,
		"VerifyURL":    verifyURL,
		"SupportEmail": s.config.SupportEmail,
		"Year":         time.Now().Year(),
	}

	return s.SendTemplateEmail(ctx, to, string(email.TemplateVerification), data)
}

// SendPasswordResetEmail sends password reset email
func (s *emailService) SendPasswordResetEmail(ctx context.Context, to, token string) error {
	resetURL := fmt.Sprintf("%s?token=%s", s.config.ResetPasswordURL, token)

	data := map[string]interface{}{
		"Name":         to,
		"Token":        token,
		"ResetURL":     resetURL,
		"SupportEmail": s.config.SupportEmail,
		"Year":         time.Now().Year(),
	}

	return s.SendTemplateEmail(ctx, to, string(email.TemplateForgotPassword), data)
}

// SendWelcomeEmail sends welcome email to new users
func (s *emailService) SendWelcomeEmail(ctx context.Context, to, name string) error {
	data := map[string]interface{}{
		"Name":         name,
		"DashboardURL": s.config.DashboardURL,
		"SupportEmail": s.config.SupportEmail,
		"Year":         time.Now().Year(),
	}

	return s.SendTemplateEmail(ctx, to, string(email.TemplateWelcome), data)
}

// SendJobApplicationEmail sends job application confirmation
func (s *emailService) SendJobApplicationEmail(ctx context.Context, to, jobTitle, companyName string) error {
	data := map[string]interface{}{
		"Name":         to,
		"JobTitle":     jobTitle,
		"CompanyName":  companyName,
		"DashboardURL": s.config.DashboardURL,
		"SupportEmail": s.config.SupportEmail,
		"Year":         time.Now().Year(),
	}

	return s.SendTemplateEmail(ctx, to, string(email.TemplateApplicationUpdate), data)
}

// SendInterviewInvitationEmail sends interview invitation
func (s *emailService) SendInterviewInvitationEmail(ctx context.Context, to, jobTitle string, interviewDate string) error {
	data := map[string]interface{}{
		"Name":          to,
		"JobTitle":      jobTitle,
		"InterviewDate": interviewDate,
		"InterviewTime": "10:00 AM",
		"InterviewURL":  fmt.Sprintf("%s/interviews", s.config.DashboardURL),
		"SupportEmail":  s.config.SupportEmail,
		"Year":          time.Now().Year(),
	}

	return s.SendTemplateEmail(ctx, to, string(email.TemplateInterviewInvite), data)
}

// SendJobStatusUpdateEmail sends job status update notification
func (s *emailService) SendJobStatusUpdateEmail(ctx context.Context, to, jobTitle, status string) error {
	data := map[string]interface{}{
		"Name":          to,
		"JobTitle":      jobTitle,
		"Status":        status,
		"ApplicationID": "APP-12345",
		"Message":       "Your application has been updated.",
		"DashboardURL":  s.config.DashboardURL,
		"SupportEmail":  s.config.SupportEmail,
		"Year":          time.Now().Year(),
	}

	return s.SendTemplateEmail(ctx, to, string(email.TemplateApplicationUpdate), data)
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
	if err := s.sendViaSMTP(log.Recipient, log.Subject, log.Body); err != nil {
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

// sendViaSMTP sends email using SMTP
func (s *emailService) sendViaSMTP(to, subject, body string) error {
	// Debug: Log connection attempt
	fmt.Printf("🔌 [DEBUG] Attempting SMTP connection to %s:%s\n", s.config.SMTPHost, s.config.SMTPPort)
	fmt.Printf("🔌 [DEBUG] SMTP Username: '%s' (empty=%v)\n", s.config.SMTPUsername, s.config.SMTPUsername == "")
	fmt.Printf("🔌 [DEBUG] SMTP From: %s\n", s.config.SMTPFrom)

	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Send email
	fmt.Printf("[DEBUG] Calling dialer.DialAndSend()...\n")
	if err := s.dialer.DialAndSend(m); err != nil {
		fmt.Printf("[ERROR] SMTP Error: %v\n", err)
		return fmt.Errorf("failed to send email via SMTP: %w", err)
	}

	// Log for development
	if s.config.IsDevelopment() {
		fmt.Printf("\n====== EMAIL SENT VIA SMTP ======\n")
		fmt.Printf("To: %s\n", to)
		fmt.Printf("Subject: %s\n", subject)
		fmt.Printf("SMTP Host: %s:%s\n", s.config.SMTPHost, s.config.SMTPPort)
		fmt.Printf("================================\n\n")
	}

	return nil
}

// mapToTemplateData converts map to TemplateData
func (s *emailService) mapToTemplateData(data map[string]interface{}) email.TemplateData {
	templateData := email.TemplateData{
		SupportEmail: s.config.SupportEmail,
		DashboardURL: s.config.DashboardURL,
		LoginURL:     s.config.FrontendURL + "/login",
		Year:         time.Now().Year(),
	}

	// Map data to struct fields
	if v, ok := data["Name"].(string); ok {
		templateData.Name = v
	}
	if v, ok := data["Email"].(string); ok {
		templateData.Email = v
	}
	if v, ok := data["Token"].(string); ok {
		templateData.Token = v
	}
	if v, ok := data["VerifyURL"].(string); ok {
		templateData.VerifyURL = v
	}
	if v, ok := data["ResetURL"].(string); ok {
		templateData.ResetURL = v
	}
	if v, ok := data["CompanyName"].(string); ok {
		templateData.CompanyName = v
	}
	if v, ok := data["JobTitle"].(string); ok {
		templateData.JobTitle = v
	}
	if v, ok := data["InterviewDate"].(string); ok {
		templateData.InterviewDate = v
	}
	if v, ok := data["InterviewTime"].(string); ok {
		templateData.InterviewTime = v
	}
	if v, ok := data["InterviewURL"].(string); ok {
		templateData.InterviewURL = v
	}
	if v, ok := data["ApplicationID"].(string); ok {
		templateData.ApplicationID = v
	}
	if v, ok := data["Status"].(string); ok {
		templateData.Status = v
	}
	if v, ok := data["Message"].(string); ok {
		templateData.Message = v
	}

	return templateData
}
