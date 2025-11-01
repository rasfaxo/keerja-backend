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
	mockRepo "keerja-backend/tests/mocks/repository"
	mockSvc "keerja-backend/tests/mocks/service"
)

// TestAuthService_Register_Success tests successful user registration scenarios
func TestAuthService_Register_Success(t *testing.T) {
	tests := []struct {
		name        string
		request     *user.RegisterRequest
		userType    string
		description string
	}{
		{
			name: "successful jobseeker registration",
			request: &user.RegisterRequest{
				FullName: "John Doe",
				Email:    "john.doe@example.com",
				Phone:    stringPtr("+6281234567890"),
				Password: "SecurePass123!",
				UserType: "jobseeker",
			},
			userType:    "jobseeker",
			description: "should successfully register a job seeker",
		},
		{
			name: "successful employer registration",
			request: &user.RegisterRequest{
				FullName: "Jane Smith",
				Email:    "jane.smith@company.com",
				Phone:    stringPtr("+6289876543210"),
				Password: "CompanyPass456!",
				UserType: "employer",
			},
			userType:    "employer",
			description: "should successfully register an employer",
		},
		{
			name: "successful registration without phone",
			request: &user.RegisterRequest{
				FullName: "No Phone User",
				Email:    "nophone@example.com",
				Phone:    nil,
				Password: "ValidPass789!",
				UserType: "jobseeker",
			},
			userType:    "jobseeker",
			description: "should successfully register without phone number",
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

			// Mock: Check email doesn't exist
			mockUserRepo.On("FindByEmail", ctx, tt.request.Email).
				Return(nil, errors.New("not found")).Once()

			// Mock: Check slug doesn't exist (for ensureUniqueSlug)
			mockUserRepo.On("FindProfileBySlug", ctx, mock.AnythingOfType("string")).
				Return(nil, errors.New("not found")).Maybe()

			// Mock: Create user successfully
			mockUserRepo.On("Create", ctx, mock.MatchedBy(func(u *user.User) bool {
				return u.Email == tt.request.Email &&
					u.FullName == tt.request.FullName &&
					u.UserType == tt.request.UserType &&
					u.IsVerified == false &&
					u.Status == "inactive" &&
					len(u.PasswordHash) > 0 // Password should be hashed
			})).Return(nil).Once().Run(func(args mock.Arguments) {
				// Simulate DB assigning ID
				u := args.Get(1).(*user.User)
				u.ID = 1
			})

			// Mock: Create profile successfully
			mockUserRepo.On("CreateProfile", ctx, mock.MatchedBy(func(p *user.UserProfile) bool {
				return p.UserID == 1 && p.Slug != nil
			})).Return(nil).Once()

			// Mock: Send verification email
			mockEmailService.On("SendVerificationEmail", ctx, tt.request.Email, mock.AnythingOfType("string")).
				Return(nil).Once()

			// Act
			registeredUser, verificationToken, err := authService.Register(ctx, tt.request)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, registeredUser)
			assert.NotEmpty(t, verificationToken)

			// Verify user properties
			assert.Equal(t, tt.request.Email, registeredUser.Email)
			assert.Equal(t, tt.request.FullName, registeredUser.FullName)
			assert.Equal(t, tt.request.UserType, registeredUser.UserType)
			assert.False(t, registeredUser.IsVerified)
			assert.Equal(t, "inactive", registeredUser.Status)
			assert.NotEmpty(t, registeredUser.PasswordHash)
			assert.NotEqual(t, tt.request.Password, registeredUser.PasswordHash, "password should be hashed")

			// Verify token is not empty and has correct length
			assert.NotEmpty(t, verificationToken)
			assert.Greater(t, len(verificationToken), 30, "verification token should be secure")

			// Verify all mocks were called as expected
			mockUserRepo.AssertExpectations(t)
			mockEmailService.AssertExpectations(t)
		})
	}
}

// TestAuthService_Register_DuplicateEmail tests duplicate email scenario
func TestAuthService_Register_DuplicateEmail(t *testing.T) {
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
	request := &user.RegisterRequest{
		FullName: "Existing User",
		Email:    "existing@example.com",
		Password: "Password123!",
		UserType: "jobseeker",
	}

	existingUser := &user.User{
		ID:       1,
		Email:    request.Email,
		FullName: "Already Registered",
	}

	// Mock: Email already exists
	mockUserRepo.On("FindByEmail", ctx, request.Email).
		Return(existingUser, nil).Once()

	// Act
	registeredUser, verificationToken, err := authService.Register(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrEmailAlreadyExists, err)
	assert.Nil(t, registeredUser)
	assert.Empty(t, verificationToken)

	mockUserRepo.AssertExpectations(t)
	// Email service should NOT be called
	mockEmailService.AssertNotCalled(t, "SendVerificationEmail")
}

// TestAuthService_Register_DatabaseError tests database failure scenarios
func TestAuthService_Register_DatabaseError(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*mockRepo.MockUserRepository)
		expectedError string
	}{
		{
			name: "user creation fails",
			setupMock: func(repo *mockRepo.MockUserRepository) {
				repo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, errors.New("not found")).Once()
				repo.On("Create", mock.Anything, mock.Anything).
					Return(errors.New("database connection failed")).Once()
			},
			expectedError: "failed to create user",
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

			tt.setupMock(mockUserRepo)

			ctx := context.Background()
			request := &user.RegisterRequest{
				FullName: "Test User",
				Email:    "test@example.com",
				Password: "Password123!",
				UserType: "jobseeker",
			}

			// Act
			registeredUser, verificationToken, err := authService.Register(ctx, request)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
			assert.Nil(t, registeredUser)
			assert.Empty(t, verificationToken)

			mockUserRepo.AssertExpectations(t)
		})
	}
}

// TestAuthService_Register_EmailServiceError tests when email service fails
func TestAuthService_Register_EmailServiceError(t *testing.T) {
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
	request := &user.RegisterRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "Password123!",
		UserType: "jobseeker",
	}

	// Mock: User creation succeeds
	mockUserRepo.On("FindByEmail", ctx, request.Email).
		Return(nil, errors.New("not found")).Once()
	mockUserRepo.On("FindProfileBySlug", ctx, mock.AnythingOfType("string")).
		Return(nil, errors.New("not found")).Maybe()
	mockUserRepo.On("Create", ctx, mock.Anything).
		Return(nil).Once().Run(func(args mock.Arguments) {
		u := args.Get(1).(*user.User)
		u.ID = 1
	})
	mockUserRepo.On("CreateProfile", ctx, mock.Anything).Return(nil).Once()

	// Mock: Email service fails
	mockEmailService.On("SendVerificationEmail", ctx, request.Email, mock.AnythingOfType("string")).
		Return(errors.New("SMTP connection failed")).Once()

	// Act
	registeredUser, verificationToken, err := authService.Register(ctx, request)

	// Assert
	// Registration should still succeed even if email fails
	assert.NoError(t, err, "registration should succeed even if email sending fails")
	assert.NotNil(t, registeredUser)
	assert.NotEmpty(t, verificationToken)

	mockUserRepo.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

// TestAuthService_Register_ProfileCreationError tests profile creation failure
func TestAuthService_Register_ProfileCreationError(t *testing.T) {
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
	request := &user.RegisterRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "Password123!",
		UserType: "jobseeker",
	}

	// Mock: User creation succeeds
	mockUserRepo.On("FindByEmail", ctx, request.Email).
		Return(nil, errors.New("not found")).Once()
	mockUserRepo.On("FindProfileBySlug", ctx, mock.AnythingOfType("string")).
		Return(nil, errors.New("not found")).Maybe()
	mockUserRepo.On("Create", ctx, mock.Anything).
		Return(nil).Once().Run(func(args mock.Arguments) {
		u := args.Get(1).(*user.User)
		u.ID = 1
	})

	// Mock: Profile creation fails
	mockUserRepo.On("CreateProfile", ctx, mock.Anything).
		Return(errors.New("profile creation failed")).Once()

	// Mock: Email service
	mockEmailService.On("SendVerificationEmail", ctx, request.Email, mock.AnythingOfType("string")).
		Return(nil).Once()

	// Act
	registeredUser, verificationToken, err := authService.Register(ctx, request)

	// Assert
	// Registration should still succeed even if profile creation fails
	assert.NoError(t, err, "registration should succeed even if profile creation fails")
	assert.NotNil(t, registeredUser)
	assert.NotEmpty(t, verificationToken)

	mockUserRepo.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

// TestAuthService_Register_PasswordHashing tests password is properly hashed
func TestAuthService_Register_PasswordHashing(t *testing.T) {
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
	plainPassword := "MySecurePassword123!"
	request := &user.RegisterRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: plainPassword,
		UserType: "jobseeker",
	}

	var capturedPasswordHash string

	// Mock setup
	mockUserRepo.On("FindByEmail", ctx, request.Email).Return(nil, errors.New("not found")).Once()
	mockUserRepo.On("FindProfileBySlug", ctx, mock.AnythingOfType("string")).
		Return(nil, errors.New("not found")).Maybe()
	mockUserRepo.On("Create", ctx, mock.Anything).
		Return(nil).Once().Run(func(args mock.Arguments) {
		u := args.Get(1).(*user.User)
		u.ID = 1
		capturedPasswordHash = u.PasswordHash
	})
	mockUserRepo.On("CreateProfile", ctx, mock.Anything).Return(nil).Once()
	mockEmailService.On("SendVerificationEmail", ctx, request.Email, mock.AnythingOfType("string")).
		Return(nil).Once()

	// Act
	registeredUser, _, err := authService.Register(ctx, request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, registeredUser)

	// Verify password is hashed
	assert.NotEmpty(t, capturedPasswordHash)
	assert.NotEqual(t, plainPassword, capturedPasswordHash, "password should be hashed")
	assert.Greater(t, len(capturedPasswordHash), 50, "hashed password should be reasonably long")
	assert.Contains(t, capturedPasswordHash, "$", "bcrypt hash should contain $")

	mockUserRepo.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
