package service

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"keerja-backend/internal/domain/notification"
)

// MockNotificationService is a mock implementation of notification.NotificationService
type MockNotificationService struct {
	mock.Mock
}

// SendNotification mocks the SendNotification method
func (m *MockNotificationService) SendNotification(ctx context.Context, req *notification.SendNotificationRequest) (*notification.Notification, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notification.Notification), args.Error(1)
}

// GetUserNotifications mocks the GetUserNotifications method
func (m *MockNotificationService) GetUserNotifications(ctx context.Context, userID int64, filter notification.NotificationFilter, page, limit int) ([]notification.Notification, int64, error) {
	args := m.Called(ctx, userID, filter, page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]notification.Notification), args.Get(1).(int64), args.Error(2)
}

// GetUnreadNotifications mocks the GetUnreadNotifications method
func (m *MockNotificationService) GetUnreadNotifications(ctx context.Context, userID int64, limit int) ([]notification.Notification, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]notification.Notification), args.Error(1)
}

// GetNotificationByID mocks the GetNotificationByID method
func (m *MockNotificationService) GetNotificationByID(ctx context.Context, id, userID int64) (*notification.Notification, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notification.Notification), args.Error(1)
}

// MarkAsRead mocks the MarkAsRead method
func (m *MockNotificationService) MarkAsRead(ctx context.Context, id, userID int64) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

// MarkAsUnread mocks the MarkAsUnread method
func (m *MockNotificationService) MarkAsUnread(ctx context.Context, id, userID int64) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

// MarkAllAsRead mocks the MarkAllAsRead method
func (m *MockNotificationService) MarkAllAsRead(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// DeleteNotification mocks the DeleteNotification method
func (m *MockNotificationService) DeleteNotification(ctx context.Context, id, userID int64) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

// DeleteAllNotifications mocks the DeleteAllNotifications method
func (m *MockNotificationService) DeleteAllNotifications(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// GetUnreadCount mocks the GetUnreadCount method
func (m *MockNotificationService) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// GetNotificationStats mocks the GetNotificationStats method
func (m *MockNotificationService) GetNotificationStats(ctx context.Context, userID int64) (*notification.NotificationStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notification.NotificationStats), args.Error(1)
}

// NotifyJobApplication mocks the NotifyJobApplication method
func (m *MockNotificationService) NotifyJobApplication(ctx context.Context, userID, jobID, applicationID int64) error {
	args := m.Called(ctx, userID, jobID, applicationID)
	return args.Error(0)
}

// NotifyInterviewScheduled mocks the NotifyInterviewScheduled method
func (m *MockNotificationService) NotifyInterviewScheduled(ctx context.Context, userID, interviewID int64, interviewDate time.Time) error {
	args := m.Called(ctx, userID, interviewID, interviewDate)
	return args.Error(0)
}

// NotifyStatusUpdate mocks the NotifyStatusUpdate method
func (m *MockNotificationService) NotifyStatusUpdate(ctx context.Context, userID, applicationID int64, oldStatus, newStatus string) error {
	args := m.Called(ctx, userID, applicationID, oldStatus, newStatus)
	return args.Error(0)
}

// NotifyJobRecommendation mocks the NotifyJobRecommendation method
func (m *MockNotificationService) NotifyJobRecommendation(ctx context.Context, userID, jobID int64) error {
	args := m.Called(ctx, userID, jobID)
	return args.Error(0)
}

// NotifyCompanyUpdate mocks the NotifyCompanyUpdate method
func (m *MockNotificationService) NotifyCompanyUpdate(ctx context.Context, userIDs []int64, companyID int64, updateType string) error {
	args := m.Called(ctx, userIDs, companyID, updateType)
	return args.Error(0)
}

// CleanupExpiredNotifications mocks the CleanupExpiredNotifications method
func (m *MockNotificationService) CleanupExpiredNotifications(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// GetNotificationPreferences mocks the GetNotificationPreferences method
func (m *MockNotificationService) GetNotificationPreferences(ctx context.Context, userID int64) (*notification.NotificationPreference, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notification.NotificationPreference), args.Error(1)
}

// UpdateNotificationPreferences mocks the UpdateNotificationPreferences method
func (m *MockNotificationService) UpdateNotificationPreferences(ctx context.Context, userID int64, prefs *notification.NotificationPreference) error {
	args := m.Called(ctx, userID, prefs)
	return args.Error(0)
}

// SendPushNotification mocks the SendPushNotification method
func (m *MockNotificationService) SendPushNotification(ctx context.Context, userID int64, notification *notification.Notification) error {
	args := m.Called(ctx, userID, notification)
	return args.Error(0)
}

// SendEmailNotification mocks the SendEmailNotification method
func (m *MockNotificationService) SendEmailNotification(ctx context.Context, userID int64, notification *notification.Notification) error {
	args := m.Called(ctx, userID, notification)
	return args.Error(0)
}
