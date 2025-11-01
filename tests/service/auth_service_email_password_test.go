package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/service"
	"keerja-backend/internal/utils"
	mockRepo "keerja-backend/tests/mocks/repository"
	mockSvc "keerja-backend/tests/mocks/service"
)

// TestAuthService_VerifyEmail_Success tests successful email verification
func TestAuthService_VerifyEmail_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   "test-secret",
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	email := "user@example.com"
	token := "valid-verification-token-12345"

	// Save token to store
	mockTokenStore.SaveVerificationToken(email, token, time.Now().Add(24*time.Hour))

	// Create unverified user
	unverifiedUser := &user.User{
		ID:         1,
		Email:      email,
		FullName:   "Test User",
		IsVerified: false,
		Status:     "inactive",
	}

	// Mock: Find user by email
	mockUserRepo.On("FindByEmail", ctx, email).
		Return(unverifiedUser, nil).Once()

	// Mock: Update user as verified
	mockUserRepo.On("Update", ctx, mock.MatchedBy(func(u *user.User) bool {
		return u.ID == 1 && u.IsVerified == true && u.Status == "active"
	})).Return(nil).Once()

	// Mock: Send welcome email
	mockEmailService.On("SendWelcomeEmail", ctx, email, "Test User").
		Return(nil).Once()

	// Act
	err := authService.VerifyEmail(ctx, token)

	// Assert
	assert.NoError(t, err)

	// Verify token was deleted from store
	_, tokenErr := mockTokenStore.GetVerificationToken(token)
	assert.Error(t, tokenErr, "token should be deleted after verification")

	mockUserRepo.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

// TestAuthService_VerifyEmail_InvalidToken tests invalid token scenarios
func TestAuthService_VerifyEmail_InvalidToken(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		setupToken    func(service.TokenStore)
		expectedError error
		description   string
	}{
		{
			name:  "non-existent token",
			token: "invalid-token-that-doesnt-exist",
			setupToken: func(ts service.TokenStore) {
				// Don't save any token
			},
			expectedError: service.ErrInvalidVerificationToken,
			description:   "should return error for non-existent token",
		},
		{
			name:  "expired token",
			token: "expired-token-12345",
			setupToken: func(ts service.TokenStore) {
				// Save token with past expiry
				ts.SaveVerificationToken("user@example.com", "expired-token-12345", time.Now().Add(-1*time.Hour))
			},
			expectedError: service.ErrTokenExpired,
			description:   "should return error for expired token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserRepo := new(mockRepo.MockUserRepository)
			mockEmailService := new(mockSvc.MockEmailService)
			mockTokenStore := service.NewInMemoryTokenStore()

			authService := service.NewAuthService(
				mockUserRepo,
				mockEmailService,
				mockTokenStore,
				service.AuthServiceConfig{
					JWTSecret:   "test-secret",
					JWTDuration: 24 * time.Hour,
				},
			)

			tt.setupToken(mockTokenStore)
			ctx := context.Background()

			// Act
			err := authService.VerifyEmail(ctx, tt.token)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)

			// User repo should not be called
			mockUserRepo.AssertNotCalled(t, "FindByEmail")
			mockUserRepo.AssertNotCalled(t, "Update")
			mockEmailService.AssertNotCalled(t, "SendWelcomeEmail")
		})
	}
}

// TestAuthService_VerifyEmail_UserNotFound tests user not found scenario
func TestAuthService_VerifyEmail_UserNotFound(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   "test-secret",
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	email := "nonexistent@example.com"
	token := "valid-token-12345"

	// Save token
	mockTokenStore.SaveVerificationToken(email, token, time.Now().Add(24*time.Hour))

	// Mock: User not found
	mockUserRepo.On("FindByEmail", ctx, email).
		Return(nil, errors.New("user not found")).Once()

	// Act
	err := authService.VerifyEmail(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrUserNotFound, err)

	mockUserRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "Update")
	mockEmailService.AssertNotCalled(t, "SendWelcomeEmail")
}

// TestAuthService_VerifyEmail_AlreadyVerified tests already verified user
func TestAuthService_VerifyEmail_AlreadyVerified(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   "test-secret",
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	email := "verified@example.com"
	token := "valid-token-12345"

	// Save token
	mockTokenStore.SaveVerificationToken(email, token, time.Now().Add(24*time.Hour))

	// Create already verified user
	verifiedUser := &user.User{
		ID:         1,
		Email:      email,
		IsVerified: true, // Already verified
		Status:     "active",
	}

	// Mock: Find user
	mockUserRepo.On("FindByEmail", ctx, email).
		Return(verifiedUser, nil).Once()

	// Act
	err := authService.VerifyEmail(ctx, token)

	// Assert
	assert.NoError(t, err, "should not error if already verified")

	// Update should not be called for already verified user
	mockUserRepo.AssertNotCalled(t, "Update")
	mockEmailService.AssertNotCalled(t, "SendWelcomeEmail")
}

// TestAuthService_VerifyEmail_UpdateFailure tests database update failure
func TestAuthService_VerifyEmail_UpdateFailure(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   "test-secret",
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	email := "user@example.com"
	token := "valid-token-12345"

	mockTokenStore.SaveVerificationToken(email, token, time.Now().Add(24*time.Hour))

	unverifiedUser := &user.User{
		ID:         1,
		Email:      email,
		IsVerified: false,
		Status:     "inactive",
	}

	mockUserRepo.On("FindByEmail", ctx, email).Return(unverifiedUser, nil).Once()
	mockUserRepo.On("Update", ctx, mock.Anything).
		Return(errors.New("database update failed")).Once()

	// Act
	err := authService.VerifyEmail(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update user")

	mockUserRepo.AssertExpectations(t)
	mockEmailService.AssertNotCalled(t, "SendWelcomeEmail")
}

// TestAuthService_VerifyEmail_WelcomeEmailFailure tests welcome email failure
func TestAuthService_VerifyEmail_WelcomeEmailFailure(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   "test-secret",
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	email := "user@example.com"
	token := "valid-token-12345"

	mockTokenStore.SaveVerificationToken(email, token, time.Now().Add(24*time.Hour))

	unverifiedUser := &user.User{
		ID:         1,
		Email:      email,
		FullName:   "Test User",
		IsVerified: false,
		Status:     "inactive",
	}

	mockUserRepo.On("FindByEmail", ctx, email).Return(unverifiedUser, nil).Once()
	mockUserRepo.On("Update", ctx, mock.Anything).Return(nil).Once()
	// Mock: Welcome email fails
	mockEmailService.On("SendWelcomeEmail", ctx, email, "Test User").
		Return(errors.New("SMTP error")).Once()

	// Act
	err := authService.VerifyEmail(ctx, token)

	// Assert
	// Verification should still succeed even if welcome email fails
	assert.NoError(t, err, "verification should succeed even if welcome email fails")

	mockUserRepo.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

// TestAuthService_ResendVerificationEmail_Success tests successful resend
func TestAuthService_ResendVerificationEmail_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   "test-secret",
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	email := "user@example.com"

	unverifiedUser := &user.User{
		ID:         1,
		Email:      email,
		FullName:   "Test User",
		IsVerified: false,
		Status:     "inactive",
	}

	mockUserRepo.On("FindByEmail", ctx, email).Return(unverifiedUser, nil).Once()
	mockEmailService.On("SendVerificationEmail", ctx, email, mock.AnythingOfType("string")).
		Return(nil).Once()

	// Act
	err := authService.ResendVerificationEmail(ctx, email)

	// Assert
	assert.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

// TestAuthService_ResendVerificationEmail_Errors tests error scenarios
func TestAuthService_ResendVerificationEmail_Errors(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*mockRepo.MockUserRepository, *mockSvc.MockEmailService)
		email         string
		expectedError string
		description   string
	}{
		{
			name: "user not found",
			setupMock: func(repo *mockRepo.MockUserRepository, email *mockSvc.MockEmailService) {
				repo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, errors.New("not found")).Once()
			},
			email:         "nonexistent@example.com",
			expectedError: "user not found",
			description:   "should return error when user not found",
		},
		{
			name: "already verified",
			setupMock: func(repo *mockRepo.MockUserRepository, emailSvc *mockSvc.MockEmailService) {
				verifiedUser := &user.User{
					ID:         1,
					Email:      "verified@example.com",
					IsVerified: true,
				}
				repo.On("FindByEmail", mock.Anything, "verified@example.com").
					Return(verifiedUser, nil).Once()
			},
			email:         "verified@example.com",
			expectedError: "email already verified",
			description:   "should return error when email already verified",
		},
		{
			name: "email service failure",
			setupMock: func(repo *mockRepo.MockUserRepository, emailSvc *mockSvc.MockEmailService) {
				unverifiedUser := &user.User{
					ID:         1,
					Email:      "user@example.com",
					IsVerified: false,
				}
				repo.On("FindByEmail", mock.Anything, "user@example.com").
					Return(unverifiedUser, nil).Once()
				emailSvc.On("SendVerificationEmail", mock.Anything, "user@example.com", mock.AnythingOfType("string")).
					Return(errors.New("SMTP connection failed")).Once()
			},
			email:         "user@example.com",
			expectedError: "failed to send verification email",
			description:   "should return error when email service fails",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserRepo := new(mockRepo.MockUserRepository)
			mockEmailService := new(mockSvc.MockEmailService)
			mockTokenStore := service.NewInMemoryTokenStore()

			authService := service.NewAuthService(
				mockUserRepo,
				mockEmailService,
				mockTokenStore,
				service.AuthServiceConfig{
					JWTSecret:   "test-secret",
					JWTDuration: 24 * time.Hour,
				},
			)

			ctx := context.Background()
			tt.setupMock(mockUserRepo, mockEmailService)

			// Act
			err := authService.ResendVerificationEmail(ctx, tt.email)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)

			mockUserRepo.AssertExpectations(t)
			mockEmailService.AssertExpectations(t)
		})
	}
}

// TestAuthService_ForgotPassword_Success tests successful forgot password
func TestAuthService_ForgotPassword_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   "test-secret",
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	email := "user@example.com"

	existingUser := &user.User{
		ID:       1,
		Email:    email,
		FullName: "Test User",
	}

	mockUserRepo.On("FindByEmail", ctx, email).Return(existingUser, nil).Once()
	mockEmailService.On("SendPasswordResetEmail", ctx, email, mock.AnythingOfType("string")).
		Return(nil).Once()

	// Act
	err := authService.ForgotPassword(ctx, email)

	// Assert
	assert.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

// TestAuthService_ForgotPassword_NonExistentUser tests forgot password with non-existent user
func TestAuthService_ForgotPassword_NonExistentUser(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   "test-secret",
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	email := "nonexistent@example.com"

	mockUserRepo.On("FindByEmail", ctx, email).
		Return(nil, errors.New("not found")).Once()

	// Act
	err := authService.ForgotPassword(ctx, email)

	// Assert
	// Should not reveal if user exists or not (security best practice)
	assert.NoError(t, err, "should not reveal if user exists")

	mockUserRepo.AssertExpectations(t)
	// Email should not be sent
	mockEmailService.AssertNotCalled(t, "SendPasswordResetEmail")
}

// TestAuthService_ResetPassword_Success tests successful password reset
func TestAuthService_ResetPassword_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   "test-secret",
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	email := "user@example.com"
	token := "valid-reset-token-12345"
	newPassword := "NewSecurePassword123!"

	// Save reset token
	mockTokenStore.SaveResetToken(email, token, time.Now().Add(1*time.Hour))

	oldPasswordHash, _ := utils.HashPassword("OldPassword123!")
	existingUser := &user.User{
		ID:           1,
		Email:        email,
		PasswordHash: oldPasswordHash,
	}

	var capturedUser *user.User
	mockUserRepo.On("FindByEmail", ctx, email).Return(existingUser, nil).Once()
	mockUserRepo.On("Update", ctx, mock.Anything).
		Return(nil).Once().
		Run(func(args mock.Arguments) {
			capturedUser = args.Get(1).(*user.User)
		})

	// Act
	err := authService.ResetPassword(ctx, token, newPassword)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, capturedUser)
	assert.NotEqual(t, oldPasswordHash, capturedUser.PasswordHash, "password should be changed")
	assert.NotEqual(t, newPassword, capturedUser.PasswordHash, "password should be hashed")

	// Verify new password works
	assert.True(t, utils.VerifyPassword(newPassword, capturedUser.PasswordHash))

	// Token should be deleted
	_, tokenErr := mockTokenStore.GetResetToken(token)
	assert.Error(t, tokenErr, "token should be deleted after use")

	mockUserRepo.AssertExpectations(t)
}

// TestAuthService_ResetPassword_InvalidToken tests invalid reset token scenarios
func TestAuthService_ResetPassword_InvalidToken(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		setupToken    func(service.TokenStore)
		expectedError error
		description   string
	}{
		{
			name:  "non-existent token",
			token: "invalid-token",
			setupToken: func(ts service.TokenStore) {
				// Don't save any token
			},
			expectedError: service.ErrInvalidResetToken,
			description:   "should return error for non-existent token",
		},
		{
			name:  "expired token",
			token: "expired-token-12345",
			setupToken: func(ts service.TokenStore) {
				ts.SaveResetToken("user@example.com", "expired-token-12345", time.Now().Add(-1*time.Hour))
			},
			expectedError: service.ErrTokenExpired,
			description:   "should return error for expired token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserRepo := new(mockRepo.MockUserRepository)
			mockEmailService := new(mockSvc.MockEmailService)
			mockTokenStore := service.NewInMemoryTokenStore()

			authService := service.NewAuthService(
				mockUserRepo,
				mockEmailService,
				mockTokenStore,
				service.AuthServiceConfig{
					JWTSecret:   "test-secret",
					JWTDuration: 24 * time.Hour,
				},
			)

			tt.setupToken(mockTokenStore)
			ctx := context.Background()

			// Act
			err := authService.ResetPassword(ctx, tt.token, "NewPassword123!")

			// Assert
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)

			mockUserRepo.AssertNotCalled(t, "FindByEmail")
			mockUserRepo.AssertNotCalled(t, "Update")
		})
	}
}

// TestAuthService_ChangePassword_Success tests successful password change
func TestAuthService_ChangePassword_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   "test-secret",
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	userID := int64(1)
	currentPassword := "CurrentPassword123!"
	newPassword := "NewSecurePassword456!"

	currentPasswordHash, _ := utils.HashPassword(currentPassword)
	existingUser := &user.User{
		ID:           userID,
		Email:        "user@example.com",
		PasswordHash: currentPasswordHash,
	}

	var capturedUser *user.User
	mockUserRepo.On("FindByID", ctx, userID).Return(existingUser, nil).Once()
	mockUserRepo.On("Update", ctx, mock.Anything).
		Return(nil).Once().
		Run(func(args mock.Arguments) {
			capturedUser = args.Get(1).(*user.User)
		})

	// Act
	err := authService.ChangePassword(ctx, userID, currentPassword, newPassword)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, capturedUser)
	assert.NotEqual(t, currentPasswordHash, capturedUser.PasswordHash)
	assert.True(t, utils.VerifyPassword(newPassword, capturedUser.PasswordHash))

	mockUserRepo.AssertExpectations(t)
}

// TestAuthService_ChangePassword_InvalidCurrentPassword tests invalid current password
func TestAuthService_ChangePassword_InvalidCurrentPassword(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   "test-secret",
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	userID := int64(1)
	wrongCurrentPassword := "WrongPassword123!"
	newPassword := "NewPassword456!"

	correctPasswordHash, _ := utils.HashPassword("CorrectPassword123!")
	existingUser := &user.User{
		ID:           userID,
		PasswordHash: correctPasswordHash,
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(existingUser, nil).Once()

	// Act
	err := authService.ChangePassword(ctx, userID, wrongCurrentPassword, newPassword)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid current password")

	mockUserRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "Update")
}
