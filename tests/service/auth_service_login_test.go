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

// TestAuthService_Login_Success tests successful login scenarios
func TestAuthService_Login_Success(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		userType    string
		isVerified  bool
		status      string
		description string
	}{
		{
			name:        "successful jobseeker login",
			email:       "jobseeker@example.com",
			password:    "ValidPassword123!",
			userType:    "jobseeker",
			isVerified:  true,
			status:      "active",
			description: "should successfully login as job seeker",
		},
		{
			name:        "successful employer login",
			email:       "employer@company.com",
			password:    "EmployerPass456!",
			userType:    "employer",
			isVerified:  true,
			status:      "active",
			description: "should successfully login as employer",
		},
		{
			name:        "successful admin login",
			email:       "admin@keerja.com",
			password:    "AdminSecure789!",
			userType:    "admin",
			isVerified:  true,
			status:      "active",
			description: "should successfully login as admin",
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
					JWTSecret:   "test-jwt-secret-key",
					JWTDuration: 24 * time.Hour,
				},
			)

			ctx := context.Background()

			// Create user with hashed password
			hashedPassword, _ := utils.HashPassword(tt.password)
			existingUser := &user.User{
				ID:           1,
				Email:        tt.email,
				FullName:     "Test User",
				PasswordHash: hashedPassword,
				UserType:     tt.userType,
				IsVerified:   tt.isVerified,
				Status:       tt.status,
			}

			// Mock: Find user by email
			mockUserRepo.On("FindByEmail", ctx, tt.email).
				Return(existingUser, nil).Once()

			// Mock: Update last login
			mockUserRepo.On("Update", ctx, mock.MatchedBy(func(u *user.User) bool {
				return u.ID == existingUser.ID && u.LastLogin != nil
			})).Return(nil).Once()

			// Act
			loggedInUser, token, err := authService.Login(ctx, tt.email, tt.password)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, loggedInUser)
			assert.NotEmpty(t, token)

			// Verify user properties
			assert.Equal(t, existingUser.ID, loggedInUser.ID)
			assert.Equal(t, existingUser.Email, loggedInUser.Email)
			assert.Equal(t, existingUser.UserType, loggedInUser.UserType)
			assert.NotNil(t, loggedInUser.LastLogin)

			// Verify JWT token format
			assert.Greater(t, len(token), 50, "JWT token should be reasonably long")
			assert.Contains(t, token, ".", "JWT should contain dots")

			mockUserRepo.AssertExpectations(t)
		})
	}
}

// TestAuthService_Login_InvalidCredentials tests invalid login scenarios
func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*mockRepo.MockUserRepository, string, string)
		email         string
		password      string
		expectedError error
		description   string
	}{
		{
			name: "non-existent user",
			setupMock: func(repo *mockRepo.MockUserRepository, email, password string) {
				repo.On("FindByEmail", mock.Anything, email).
					Return(nil, errors.New("not found")).Once()
			},
			email:         "nonexistent@example.com",
			password:      "AnyPassword123!",
			expectedError: service.ErrInvalidCredentials,
			description:   "should return invalid credentials for non-existent user",
		},
		{
			name: "incorrect password",
			setupMock: func(repo *mockRepo.MockUserRepository, email, password string) {
				// Create user with different password
				hashedPassword, _ := utils.HashPassword("CorrectPassword123!")
				existingUser := &user.User{
					ID:           1,
					Email:        email,
					PasswordHash: hashedPassword,
					IsVerified:   true,
					Status:       "active",
				}
				repo.On("FindByEmail", mock.Anything, email).
					Return(existingUser, nil).Once()
			},
			email:         "user@example.com",
			password:      "WrongPassword456!",
			expectedError: service.ErrInvalidCredentials,
			description:   "should return invalid credentials for wrong password",
		},
		{
			name: "null user returned",
			setupMock: func(repo *mockRepo.MockUserRepository, email, password string) {
				repo.On("FindByEmail", mock.Anything, email).
					Return(nil, nil).Once()
			},
			email:         "null@example.com",
			password:      "Password123!",
			expectedError: service.ErrInvalidCredentials,
			description:   "should return invalid credentials when user is nil",
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
					JWTSecret:   "test-jwt-secret-key",
					JWTDuration: 24 * time.Hour,
				},
			)

			ctx := context.Background()
			tt.setupMock(mockUserRepo, tt.email, tt.password)

			// Act
			loggedInUser, token, err := authService.Login(ctx, tt.email, tt.password)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
			assert.Nil(t, loggedInUser)
			assert.Empty(t, token)

			mockUserRepo.AssertExpectations(t)
			// Update should not be called on failed login
			mockUserRepo.AssertNotCalled(t, "Update")
		})
	}
}

// TestAuthService_Login_EmailNotVerified tests unverified email scenario
func TestAuthService_Login_EmailNotVerified(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   "test-jwt-secret-key",
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	email := "unverified@example.com"
	password := "Password123!"

	// Create unverified user
	hashedPassword, _ := utils.HashPassword(password)
	unverifiedUser := &user.User{
		ID:           1,
		Email:        email,
		PasswordHash: hashedPassword,
		IsVerified:   false, // Not verified
		Status:       "inactive",
	}

	// Mock: Find user by email
	mockUserRepo.On("FindByEmail", ctx, email).
		Return(unverifiedUser, nil).Once()

	// Act
	loggedInUser, token, err := authService.Login(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrEmailNotVerified, err)
	assert.Nil(t, loggedInUser)
	assert.Empty(t, token)

	mockUserRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "Update")
}

// TestAuthService_Login_AccountStatus tests different account statuses
func TestAuthService_Login_AccountStatus(t *testing.T) {
	tests := []struct {
		name             string
		accountStatus    string
		expectedErrorMsg string
		description      string
	}{
		{
			name:             "suspended account",
			accountStatus:    "suspended",
			expectedErrorMsg: "account is suspended",
			description:      "should not allow login for suspended account",
		},
		{
			name:             "deactivated account",
			accountStatus:    "deactivated",
			expectedErrorMsg: "account is deactivated",
			description:      "should not allow login for deactivated account",
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
					JWTSecret:   "test-jwt-secret-key",
					JWTDuration: 24 * time.Hour,
				},
			)

			ctx := context.Background()
			email := "user@example.com"
			password := "Password123!"

			// Create user with specific status
			hashedPassword, _ := utils.HashPassword(password)
			statusUser := &user.User{
				ID:           1,
				Email:        email,
				PasswordHash: hashedPassword,
				IsVerified:   true,
				Status:       tt.accountStatus,
			}

			// Mock: Find user by email
			mockUserRepo.On("FindByEmail", ctx, email).
				Return(statusUser, nil).Once()

			// Act
			loggedInUser, token, err := authService.Login(ctx, email, password)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErrorMsg)
			assert.Nil(t, loggedInUser)
			assert.Empty(t, token)

			mockUserRepo.AssertExpectations(t)
			mockUserRepo.AssertNotCalled(t, "Update")
		})
	}
}

// TestAuthService_Login_JWTTokenGeneration tests JWT token generation
func TestAuthService_Login_JWTTokenGeneration(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	jwtSecret := "super-secret-jwt-key-for-testing"
	jwtDuration := 2 * time.Hour

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   jwtSecret,
			JWTDuration: jwtDuration,
		},
	)

	ctx := context.Background()
	email := "user@example.com"
	password := "Password123!"

	// Create verified active user
	hashedPassword, _ := utils.HashPassword(password)
	activeUser := &user.User{
		ID:           123,
		Email:        email,
		FullName:     "Test User",
		PasswordHash: hashedPassword,
		UserType:     "jobseeker",
		IsVerified:   true,
		Status:       "active",
	}

	mockUserRepo.On("FindByEmail", ctx, email).Return(activeUser, nil).Once()
	mockUserRepo.On("Update", ctx, mock.Anything).Return(nil).Once()

	// Act
	loggedInUser, token, err := authService.Login(ctx, email, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, loggedInUser)
	assert.NotEmpty(t, token)

	// Validate JWT token can be parsed
	claims, parseErr := utils.ValidateToken(token, jwtSecret)
	assert.NoError(t, parseErr)
	assert.NotNil(t, claims)
	assert.Equal(t, activeUser.ID, claims.UserID)
	assert.Equal(t, activeUser.Email, claims.Email)
	assert.Equal(t, activeUser.UserType, claims.UserType)

	mockUserRepo.AssertExpectations(t)
}

// TestAuthService_Login_LastLoginUpdate tests last login timestamp update
func TestAuthService_Login_LastLoginUpdate(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   "test-jwt-secret-key",
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	email := "user@example.com"
	password := "Password123!"

	// Create user
	hashedPassword, _ := utils.HashPassword(password)
	testUser := &user.User{
		ID:           1,
		Email:        email,
		PasswordHash: hashedPassword,
		IsVerified:   true,
		Status:       "active",
		LastLogin:    nil, // Initially no last login
	}

	var capturedUser *user.User
	mockUserRepo.On("FindByEmail", ctx, email).Return(testUser, nil).Once()
	mockUserRepo.On("Update", ctx, mock.Anything).
		Return(nil).Once().
		Run(func(args mock.Arguments) {
			capturedUser = args.Get(1).(*user.User)
		})

	// Act
	beforeLogin := time.Now()
	loggedInUser, token, err := authService.Login(ctx, email, password)
	afterLogin := time.Now()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, loggedInUser)
	assert.NotEmpty(t, token)

	// Verify LastLogin was updated
	assert.NotNil(t, capturedUser)
	assert.NotNil(t, capturedUser.LastLogin)

	// LastLogin should be between beforeLogin and afterLogin
	assert.True(t, capturedUser.LastLogin.After(beforeLogin) || capturedUser.LastLogin.Equal(beforeLogin))
	assert.True(t, capturedUser.LastLogin.Before(afterLogin) || capturedUser.LastLogin.Equal(afterLogin))

	mockUserRepo.AssertExpectations(t)
}

// TestAuthService_Login_UpdateLastLoginFailure tests when updating last login fails
func TestAuthService_Login_UpdateLastLoginFailure(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   "test-jwt-secret-key",
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	email := "user@example.com"
	password := "Password123!"

	hashedPassword, _ := utils.HashPassword(password)
	testUser := &user.User{
		ID:           1,
		Email:        email,
		PasswordHash: hashedPassword,
		IsVerified:   true,
		Status:       "active",
	}

	mockUserRepo.On("FindByEmail", ctx, email).Return(testUser, nil).Once()
	// Mock: Update fails but should not affect login success
	mockUserRepo.On("Update", ctx, mock.Anything).
		Return(errors.New("database update failed")).Once()

	// Act
	loggedInUser, token, err := authService.Login(ctx, email, password)

	// Assert
	// Login should still succeed even if last login update fails
	assert.NoError(t, err, "login should succeed even if last login update fails")
	assert.NotNil(t, loggedInUser)
	assert.NotEmpty(t, token)

	mockUserRepo.AssertExpectations(t)
}

// TestAuthService_Login_ConcurrentLogins tests concurrent login attempts
func TestAuthService_Login_ConcurrentLogins(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   "test-jwt-secret-key",
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	email := "user@example.com"
	password := "Password123!"

	hashedPassword, _ := utils.HashPassword(password)
	testUser := &user.User{
		ID:           1,
		Email:        email,
		PasswordHash: hashedPassword,
		IsVerified:   true,
		Status:       "active",
	}

	// Mock for multiple concurrent calls
	mockUserRepo.On("FindByEmail", ctx, email).Return(testUser, nil).Times(5)
	mockUserRepo.On("Update", ctx, mock.Anything).Return(nil).Times(5)

	// Act - Simulate 5 concurrent login attempts
	results := make(chan error, 5)
	for i := 0; i < 5; i++ {
		go func() {
			_, token, err := authService.Login(ctx, email, password)
			if err != nil {
				results <- err
			} else if token == "" {
				results <- errors.New("empty token")
			} else {
				results <- nil
			}
		}()
	}

	// Assert - Collect results
	for i := 0; i < 5; i++ {
		err := <-results
		assert.NoError(t, err, "concurrent login should succeed")
	}

	mockUserRepo.AssertExpectations(t)
}
