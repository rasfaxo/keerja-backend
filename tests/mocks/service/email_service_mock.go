package service

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockEmailService is a mock implementation of email service
type MockEmailService struct {
	mock.Mock
}

// SendVerificationEmail mocks sending verification email
func (m *MockEmailService) SendVerificationEmail(ctx context.Context, email, token string) error {
	args := m.Called(ctx, email, token)
	return args.Error(0)
}

// SendWelcomeEmail mocks sending welcome email
func (m *MockEmailService) SendWelcomeEmail(ctx context.Context, email, name string) error {
	args := m.Called(ctx, email, name)
	return args.Error(0)
}

// SendPasswordResetEmail mocks sending password reset email
func (m *MockEmailService) SendPasswordResetEmail(ctx context.Context, email, token string) error {
	args := m.Called(ctx, email, token)
	return args.Error(0)
}

// SendPasswordChangedEmail mocks sending password changed email
func (m *MockEmailService) SendPasswordChangedEmail(ctx context.Context, email, name string) error {
	args := m.Called(ctx, email, name)
	return args.Error(0)
}

// SendJobApplicationEmail mocks sending job application email
func (m *MockEmailService) SendJobApplicationEmail(ctx context.Context, email, name, jobTitle, companyName string) error {
	args := m.Called(ctx, email, name, jobTitle, companyName)
	return args.Error(0)
}

// SendJobStatusUpdateEmail mocks sending job status update email
func (m *MockEmailService) SendJobStatusUpdateEmail(ctx context.Context, email, name, jobTitle, status string) error {
	args := m.Called(ctx, email, name, jobTitle, status)
	return args.Error(0)
}

// SendInterviewInvitationEmail mocks sending interview invitation email
func (m *MockEmailService) SendInterviewInvitationEmail(ctx context.Context, email, name, jobTitle, interviewDate string) error {
	args := m.Called(ctx, email, name, jobTitle, interviewDate)
	return args.Error(0)
}

// SendInterviewReminderEmail mocks sending interview reminder email
func (m *MockEmailService) SendInterviewReminderEmail(ctx context.Context, email, name, jobTitle, interviewDate string) error {
	args := m.Called(ctx, email, name, jobTitle, interviewDate)
	return args.Error(0)
}

// SendJobPostedEmail mocks sending job posted email
func (m *MockEmailService) SendJobPostedEmail(ctx context.Context, email, name, jobTitle string) error {
	args := m.Called(ctx, email, name, jobTitle)
	return args.Error(0)
}

// SendApplicationReceivedEmail mocks sending application received email
func (m *MockEmailService) SendApplicationReceivedEmail(ctx context.Context, email, name, applicantName, jobTitle string) error {
	args := m.Called(ctx, email, name, applicantName, jobTitle)
	return args.Error(0)
}

// SendOfferEmail mocks sending offer email
func (m *MockEmailService) SendOfferEmail(ctx context.Context, email, name, jobTitle, companyName string) error {
	args := m.Called(ctx, email, name, jobTitle, companyName)
	return args.Error(0)
}

// SendRejectionEmail mocks sending rejection email
func (m *MockEmailService) SendRejectionEmail(ctx context.Context, email, name, jobTitle, reason string) error {
	args := m.Called(ctx, email, name, jobTitle, reason)
	return args.Error(0)
}

// SendBulkEmail mocks sending bulk email
func (m *MockEmailService) SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error {
	args := m.Called(ctx, recipients, subject, body)
	return args.Error(0)
}

// SendTemplateEmail mocks sending template email
func (m *MockEmailService) SendTemplateEmail(ctx context.Context, email, templateName string, data map[string]interface{}) error {
	args := m.Called(ctx, email, templateName, data)
	return args.Error(0)
}
