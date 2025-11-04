package repository

import (
	"context"

	"github.com/stretchr/testify/mock"

	"keerja-backend/internal/domain/notification"
)

// MockNotificationRepository is a mock implementation of notification.NotificationRepository
type MockNotificationRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockNotificationRepository) Create(ctx context.Context, n *notification.Notification) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

// FindByID mocks the FindByID method
func (m *MockNotificationRepository) FindByID(ctx context.Context, id int64) (*notification.Notification, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notification.Notification), args.Error(1)
}

// Update mocks the Update method
func (m *MockNotificationRepository) Update(ctx context.Context, n *notification.Notification) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockNotificationRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ListByUser mocks the ListByUser method
func (m *MockNotificationRepository) ListByUser(ctx context.Context, userID int64, filter notification.NotificationFilter, page, limit int) ([]notification.Notification, int64, error) {
	args := m.Called(ctx, userID, filter, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]notification.Notification), args.Get(1).(int64), args.Error(2)
}

// GetUnreadByUser mocks the GetUnreadByUser method
func (m *MockNotificationRepository) GetUnreadByUser(ctx context.Context, userID int64, limit int) ([]notification.Notification, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]notification.Notification), args.Error(1)
}

// CountUnreadByUser mocks the CountUnreadByUser method
func (m *MockNotificationRepository) CountUnreadByUser(ctx context.Context, userID int64) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// MarkAsRead mocks the MarkAsRead method
func (m *MockNotificationRepository) MarkAsRead(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MarkAllAsRead mocks the MarkAllAsRead method
func (m *MockNotificationRepository) MarkAllAsRead(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// DeleteByUser mocks the DeleteByUser method
func (m *MockNotificationRepository) DeleteByUser(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// GetExpiredNotifications mocks the GetExpiredNotifications method
func (m *MockNotificationRepository) GetExpiredNotifications(ctx context.Context, limit int) ([]notification.Notification, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]notification.Notification), args.Error(1)
}

// GetStats mocks the GetStats method
func (m *MockNotificationRepository) GetStats(ctx context.Context, userID int64) (*notification.NotificationStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notification.NotificationStats), args.Error(1)
}

// BulkCreate mocks the BulkCreate method
func (m *MockNotificationRepository) BulkCreate(ctx context.Context, notifications []notification.Notification) error {
	args := m.Called(ctx, notifications)
	return args.Error(0)
}

// FindPreferenceByUser mocks the FindPreferenceByUser method
func (m *MockNotificationRepository) FindPreferenceByUser(ctx context.Context, userID int64) (*notification.NotificationPreference, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notification.NotificationPreference), args.Error(1)
}

// CreatePreference mocks the CreatePreference method
func (m *MockNotificationRepository) CreatePreference(ctx context.Context, preference *notification.NotificationPreference) error {
	args := m.Called(ctx, preference)
	return args.Error(0)
}

// UpdatePreference mocks the UpdatePreference method
func (m *MockNotificationRepository) UpdatePreference(ctx context.Context, preference *notification.NotificationPreference) error {
	args := m.Called(ctx, preference)
	return args.Error(0)
}
