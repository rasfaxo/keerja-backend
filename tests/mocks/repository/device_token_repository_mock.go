package repository

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"keerja-backend/internal/domain/notification"
)

// MockDeviceTokenRepository is a mock implementation of notification.DeviceTokenRepository
type MockDeviceTokenRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockDeviceTokenRepository) Create(ctx context.Context, token *notification.DeviceToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

// FindByID mocks the FindByID method
func (m *MockDeviceTokenRepository) FindByID(ctx context.Context, id int64) (*notification.DeviceToken, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notification.DeviceToken), args.Error(1)
}

// FindByToken mocks the FindByToken method
func (m *MockDeviceTokenRepository) FindByToken(ctx context.Context, token string) (*notification.DeviceToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notification.DeviceToken), args.Error(1)
}

// FindByUser mocks the FindByUser method
func (m *MockDeviceTokenRepository) FindByUser(ctx context.Context, userID int64) ([]notification.DeviceToken, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]notification.DeviceToken), args.Error(1)
}

// FindByUserAndPlatform mocks the FindByUserAndPlatform method
func (m *MockDeviceTokenRepository) FindByUserAndPlatform(ctx context.Context, userID int64, platform notification.Platform) ([]notification.DeviceToken, error) {
	args := m.Called(ctx, userID, platform)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]notification.DeviceToken), args.Error(1)
}

// Update mocks the Update method
func (m *MockDeviceTokenRepository) Update(ctx context.Context, token *notification.DeviceToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockDeviceTokenRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// DeleteByToken mocks the DeleteByToken method
func (m *MockDeviceTokenRepository) DeleteByToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

// Deactivate mocks the Deactivate method
func (m *MockDeviceTokenRepository) Deactivate(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

// FindInactiveTokens mocks the FindInactiveTokens method
func (m *MockDeviceTokenRepository) FindInactiveTokens(ctx context.Context, inactiveDays int, limit int) ([]notification.DeviceToken, error) {
	args := m.Called(ctx, inactiveDays, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]notification.DeviceToken), args.Error(1)
}

// FindInactive mocks the FindInactive method
func (m *MockDeviceTokenRepository) FindInactive(ctx context.Context, cutoffDate time.Time) ([]notification.DeviceToken, error) {
	args := m.Called(ctx, cutoffDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]notification.DeviceToken), args.Error(1)
}

// FindByFailureCount mocks the FindByFailureCount method
func (m *MockDeviceTokenRepository) FindByFailureCount(ctx context.Context, minFailures int) ([]notification.DeviceToken, error) {
	args := m.Called(ctx, minFailures)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]notification.DeviceToken), args.Error(1)
}

// CountByUser mocks the CountByUser method
func (m *MockDeviceTokenRepository) CountByUser(ctx context.Context, userID int64) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// BatchUpdate mocks the BatchUpdate method
func (m *MockDeviceTokenRepository) BatchUpdate(ctx context.Context, tokens []notification.DeviceToken) error {
	args := m.Called(ctx, tokens)
	return args.Error(0)
}

// FindByUserIDs mocks the FindByUserIDs method
func (m *MockDeviceTokenRepository) FindByUserIDs(ctx context.Context, userIDs []int64) ([]notification.DeviceToken, error) {
	args := m.Called(ctx, userIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]notification.DeviceToken), args.Error(1)
}
