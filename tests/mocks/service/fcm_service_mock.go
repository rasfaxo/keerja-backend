package service

import (
	"context"

	"github.com/stretchr/testify/mock"

	"keerja-backend/internal/domain/notification"
)

// MockFCMService is a mock implementation of FCM (Firebase Cloud Messaging) service
// This mock implements notification.PushNotificationService interface
type MockFCMService struct {
	mock.Mock
}

// SendToDevice mocks sending notification to a single device token
func (m *MockFCMService) SendToDevice(ctx context.Context, token string, message *notification.PushMessage) (*notification.PushResult, error) {
	args := m.Called(ctx, token, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notification.PushResult), args.Error(1)
}

// SendToUser mocks sending notification to all devices of a user
func (m *MockFCMService) SendToUser(ctx context.Context, userID int64, message *notification.PushMessage) ([]notification.PushResult, error) {
	args := m.Called(ctx, userID, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]notification.PushResult), args.Error(1)
}

// SendToMultipleUsers mocks sending notification to multiple users
func (m *MockFCMService) SendToMultipleUsers(ctx context.Context, userIDs []int64, message *notification.PushMessage) (map[int64][]notification.PushResult, error) {
	args := m.Called(ctx, userIDs, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[int64][]notification.PushResult), args.Error(1)
}

// SendToTopic mocks sending notification to a topic
func (m *MockFCMService) SendToTopic(ctx context.Context, topic string, message *notification.PushMessage) (*notification.PushResult, error) {
	args := m.Called(ctx, topic, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notification.PushResult), args.Error(1)
}

// RegisterDeviceToken mocks registering a device token
func (m *MockFCMService) RegisterDeviceToken(ctx context.Context, userID int64, token string, platform notification.Platform, deviceInfo *notification.DeviceInfo) error {
	args := m.Called(ctx, userID, token, platform, deviceInfo)
	return args.Error(0)
}

// UnregisterDeviceToken mocks unregistering a device token
func (m *MockFCMService) UnregisterDeviceToken(ctx context.Context, userID int64, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

// RefreshDeviceToken mocks refreshing a device token
func (m *MockFCMService) RefreshDeviceToken(ctx context.Context, oldToken, newToken string) error {
	args := m.Called(ctx, oldToken, newToken)
	return args.Error(0)
}

// GetUserDevices mocks getting user's devices
func (m *MockFCMService) GetUserDevices(ctx context.Context, userID int64) ([]notification.DeviceToken, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]notification.DeviceToken), args.Error(1)
}

// ValidateToken mocks validating a device token
func (m *MockFCMService) ValidateToken(ctx context.Context, token string) (bool, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(bool), args.Error(1)
}

// CleanupInactiveTokens mocks cleaning up inactive device tokens
func (m *MockFCMService) CleanupInactiveTokens(ctx context.Context, inactiveDays int) error {
	args := m.Called(ctx, inactiveDays)
	return args.Error(0)
}
