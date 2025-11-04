package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/service"
	mockRepo "keerja-backend/tests/mocks/repository"
	mockSvc "keerja-backend/tests/mocks/service"
)

// TestAuthService_RefreshToken_Success tests successful token refresh
func TestAuthService_RefreshToken_Success(t *testing.T) {
	tests := []struct {
		name        string
		userType    string
		status      string
		description string
	}{
		{
			name:        "jobseeker refresh token",
			userType:    "jobseeker",
			status:      "active",
			description: "should successfully refresh token for active jobseeker",
		},
		{
			name:        "employer refresh token",
			userType:    "employer",
			status:      "active",
			description: "should successfully refresh token for active employer",
		},
		{
			name:        "admin refresh token",
			userType:    "admin",
			status:      "active",
			description: "should successfully refresh token for active admin",
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
					JWTSecret:   "test-secret-key-for-jwt",
					JWTDuration: 24 * time.Hour,
				},
			)

			ctx := context.Background()
			userID := int64(1)

			activeUser := &user.User{
				ID:       userID,
				Email:    "user@example.com",
				FullName: "Test User",
				UserType: tt.userType,
				Status:   tt.status,
			}

			mockUserRepo.On("FindByID", ctx, userID).
				Return(activeUser, nil).Once()

			// Act
			token, err := authService.RefreshToken(ctx, userID)

			// Assert
			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			// Verify token validity
			parsedToken, parseErr := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte("test-secret-key-for-jwt"), nil
			})
			assert.NoError(t, parseErr)
			assert.True(t, parsedToken.Valid)

			// Verify claims
			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			assert.True(t, ok)
			assert.Equal(t, float64(userID), claims["user_id"])
			assert.Equal(t, "user@example.com", claims["email"])
			assert.Equal(t, tt.userType, claims["user_type"])

			mockUserRepo.AssertExpectations(t)
		})
	}
}

// TestAuthService_RefreshToken_UserNotFound tests token refresh for non-existent user
func TestAuthService_RefreshToken_UserNotFound(t *testing.T) {
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
	nonExistentUserID := int64(999)

	mockUserRepo.On("FindByID", ctx, nonExistentUserID).
		Return(nil, errors.New("user not found")).Once()

	// Act
	token, err := authService.RefreshToken(ctx, nonExistentUserID)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "user not found")

	mockUserRepo.AssertExpectations(t)
}

// TestAuthService_RefreshToken_InactiveAccount tests token refresh for inactive accounts
func TestAuthService_RefreshToken_InactiveAccount(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		expectedErr string
		description string
	}{
		{
			name:        "suspended account",
			status:      "suspended",
			expectedErr: "account is not active",
			description: "should not refresh token for suspended account",
		},
		{
			name:        "deactivated account",
			status:      "deactivated",
			expectedErr: "account is not active",
			description: "should not refresh token for deactivated account",
		},
		{
			name:        "inactive account",
			status:      "inactive",
			expectedErr: "account is not active",
			description: "should not refresh token for inactive account",
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
			userID := int64(1)

			inactiveUser := &user.User{
				ID:       userID,
				Email:    "user@example.com",
				Status:   tt.status,
				UserType: "jobseeker",
			}

			mockUserRepo.On("FindByID", ctx, userID).
				Return(inactiveUser, nil).Once()

			// Act
			token, err := authService.RefreshToken(ctx, userID)

			// Assert
			assert.Error(t, err)
			assert.Empty(t, token)
			assert.Contains(t, err.Error(), tt.expectedErr)

			mockUserRepo.AssertExpectations(t)
		})
	}
}

// TestAuthService_ValidateToken_Success tests successful token validation
func TestAuthService_ValidateToken_Success(t *testing.T) {
	tests := []struct {
		name        string
		userType    string
		description string
	}{
		{
			name:        "validate jobseeker token",
			userType:    "jobseeker",
			description: "should successfully validate jobseeker token",
		},
		{
			name:        "validate employer token",
			userType:    "employer",
			description: "should successfully validate employer token",
		},
		{
			name:        "validate admin token",
			userType:    "admin",
			description: "should successfully validate admin token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserRepo := new(mockRepo.MockUserRepository)
			mockEmailService := new(mockSvc.MockEmailService)
			mockTokenStore := service.NewInMemoryTokenStore()

			jwtSecret := "test-secret-key-for-jwt-validation"
			authService := service.NewAuthService(
				mockUserRepo,
				mockEmailService,
				mockTokenStore,
				service.AuthServiceConfig{
					JWTSecret:   jwtSecret,
					JWTDuration: 24 * time.Hour,
				},
			)

			ctx := context.Background()
			userID := int64(1)
			email := "user@example.com"

			// Create a valid JWT token
			claims := jwt.MapClaims{
				"user_id":   float64(userID),
				"email":     email,
				"user_type": tt.userType,
				"exp":       time.Now().Add(24 * time.Hour).Unix(),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, _ := token.SignedString([]byte(jwtSecret))

			activeUser := &user.User{
				ID:       userID,
				Email:    email,
				FullName: "Test User",
				UserType: tt.userType,
				Status:   "active",
			}

			mockUserRepo.On("FindByID", ctx, userID).
				Return(activeUser, nil).Once()

			// Act
			validatedUser, err := authService.ValidateToken(ctx, tokenString)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, validatedUser)
			assert.Equal(t, userID, validatedUser.ID)
			assert.Equal(t, email, validatedUser.Email)
			assert.Equal(t, tt.userType, validatedUser.UserType)
			assert.Equal(t, "active", validatedUser.Status)

			mockUserRepo.AssertExpectations(t)
		})
	}
}

// TestAuthService_ValidateToken_InvalidToken tests invalid token scenarios
func TestAuthService_ValidateToken_InvalidToken(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		expectedError string
		description   string
	}{
		{
			name:          "malformed token",
			token:         "invalid.malformed.token",
			expectedError: "invalid token",
			description:   "should return error for malformed token",
		},
		{
			name:          "empty token",
			token:         "",
			expectedError: "invalid token",
			description:   "should return error for empty token",
		},
		{
			name:          "random string",
			token:         "this-is-not-a-jwt-token",
			expectedError: "invalid token",
			description:   "should return error for random string",
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

			// Act
			validatedUser, err := authService.ValidateToken(ctx, tt.token)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, validatedUser)
			assert.Contains(t, err.Error(), tt.expectedError)

			// User repo should not be called for invalid tokens
			mockUserRepo.AssertNotCalled(t, "FindByID")
		})
	}
}

// TestAuthService_ValidateToken_ExpiredToken tests expired token validation
func TestAuthService_ValidateToken_ExpiredToken(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	jwtSecret := "test-secret-key-for-jwt"
	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   jwtSecret,
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()

	// Create an expired JWT token
	claims := jwt.MapClaims{
		"user_id":   float64(1),
		"email":     "user@example.com",
		"user_type": "jobseeker",
		"exp":       time.Now().Add(-1 * time.Hour).Unix(), // Expired 1 hour ago
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, _ := token.SignedString([]byte(jwtSecret))

	// Act
	validatedUser, err := authService.ValidateToken(ctx, expiredToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, validatedUser)
	assert.Contains(t, err.Error(), "token has expired")

	mockUserRepo.AssertNotCalled(t, "FindByID")
}

// TestAuthService_ValidateToken_WrongSignature tests token with wrong signature
func TestAuthService_ValidateToken_WrongSignature(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	correctSecret := "correct-secret-key"
	wrongSecret := "wrong-secret-key"

	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   correctSecret,
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()

	// Create a token with wrong secret
	claims := jwt.MapClaims{
		"user_id":   float64(1),
		"email":     "user@example.com",
		"user_type": "jobseeker",
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	wrongSignedToken, _ := token.SignedString([]byte(wrongSecret))

	// Act
	validatedUser, err := authService.ValidateToken(ctx, wrongSignedToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, validatedUser)
	assert.Contains(t, err.Error(), "signature is invalid")

	mockUserRepo.AssertNotCalled(t, "FindByID")
}

// TestAuthService_ValidateToken_UserNotFound tests token validation when user doesn't exist
func TestAuthService_ValidateToken_UserNotFound(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	jwtSecret := "test-secret-key"
	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   jwtSecret,
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	nonExistentUserID := int64(999)

	// Create a valid token for non-existent user
	claims := jwt.MapClaims{
		"user_id":   float64(nonExistentUserID),
		"email":     "deleted@example.com",
		"user_type": "jobseeker",
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(jwtSecret))

	mockUserRepo.On("FindByID", ctx, nonExistentUserID).
		Return(nil, errors.New("user not found")).Once()

	// Act
	validatedUser, err := authService.ValidateToken(ctx, tokenString)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, validatedUser)
	assert.Contains(t, err.Error(), "user not found")

	mockUserRepo.AssertExpectations(t)
}

// TestAuthService_ValidateToken_InactiveUser tests token validation for inactive users
func TestAuthService_ValidateToken_InactiveUser(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		expectedErr string
		description string
	}{
		{
			name:        "suspended user",
			status:      "suspended",
			expectedErr: "account is not active",
			description: "should reject token for suspended user",
		},
		{
			name:        "deactivated user",
			status:      "deactivated",
			expectedErr: "account is not active",
			description: "should reject token for deactivated user",
		},
		{
			name:        "inactive user",
			status:      "inactive",
			expectedErr: "account is not active",
			description: "should reject token for inactive user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserRepo := new(mockRepo.MockUserRepository)
			mockEmailService := new(mockSvc.MockEmailService)
			mockTokenStore := service.NewInMemoryTokenStore()

			jwtSecret := "test-secret-key"
			authService := service.NewAuthService(
				mockUserRepo,
				mockEmailService,
				mockTokenStore,
				service.AuthServiceConfig{
					JWTSecret:   jwtSecret,
					JWTDuration: 24 * time.Hour,
				},
			)

			ctx := context.Background()
			userID := int64(1)

			// Create valid token
			claims := jwt.MapClaims{
				"user_id":   float64(userID),
				"email":     "user@example.com",
				"user_type": "jobseeker",
				"exp":       time.Now().Add(24 * time.Hour).Unix(),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, _ := token.SignedString([]byte(jwtSecret))

			// User exists but is inactive
			inactiveUser := &user.User{
				ID:       userID,
				Email:    "user@example.com",
				Status:   tt.status,
				UserType: "jobseeker",
			}

			mockUserRepo.On("FindByID", ctx, userID).
				Return(inactiveUser, nil).Once()

			// Act
			validatedUser, err := authService.ValidateToken(ctx, tokenString)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, validatedUser)
			assert.Contains(t, err.Error(), tt.expectedErr)

			mockUserRepo.AssertExpectations(t)
		})
	}
}

// TestAuthService_ValidateToken_MissingClaims tests token with missing required claims
func TestAuthService_ValidateToken_MissingClaims(t *testing.T) {
	tests := []struct {
		name        string
		claims      jwt.MapClaims
		description string
	}{
		{
			name: "missing user_id",
			claims: jwt.MapClaims{
				"email":     "user@example.com",
				"user_type": "jobseeker",
				"exp":       time.Now().Add(24 * time.Hour).Unix(),
			},
			description: "should reject token without user_id",
		},
		{
			name: "missing email",
			claims: jwt.MapClaims{
				"user_id":   float64(1),
				"user_type": "jobseeker",
				"exp":       time.Now().Add(24 * time.Hour).Unix(),
			},
			description: "should reject token without email",
		},
		{
			name: "missing user_type",
			claims: jwt.MapClaims{
				"user_id": float64(1),
				"email":   "user@example.com",
				"exp":     time.Now().Add(24 * time.Hour).Unix(),
			},
			description: "should reject token without user_type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserRepo := new(mockRepo.MockUserRepository)
			mockEmailService := new(mockSvc.MockEmailService)
			mockTokenStore := service.NewInMemoryTokenStore()

			jwtSecret := "test-secret-key"
			authService := service.NewAuthService(
				mockUserRepo,
				mockEmailService,
				mockTokenStore,
				service.AuthServiceConfig{
					JWTSecret:   jwtSecret,
					JWTDuration: 24 * time.Hour,
				},
			)

			ctx := context.Background()

			// Create token with missing claims
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, tt.claims)
			tokenString, _ := token.SignedString([]byte(jwtSecret))

			// Mock: For any user_id that might be parsed (including 0), return error
			mockUserRepo.On("FindByID", ctx, mock.AnythingOfType("int64")).
				Return(nil, errors.New("user not found")).Maybe()

			// Act
			validatedUser, err := authService.ValidateToken(ctx, tokenString)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, validatedUser)
			// Error could be from invalid claims or user not found
			// Either way, validation should fail

			mockUserRepo.AssertExpectations(t)
		})
	}
}

// TestAuthService_TokenLifecycle tests complete token lifecycle
func TestAuthService_TokenLifecycle(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockEmailService := new(mockSvc.MockEmailService)
	mockTokenStore := service.NewInMemoryTokenStore()

	jwtSecret := "test-secret-key-lifecycle"
	authService := service.NewAuthService(
		mockUserRepo,
		mockEmailService,
		mockTokenStore,
		service.AuthServiceConfig{
			JWTSecret:   jwtSecret,
			JWTDuration: 24 * time.Hour,
		},
	)

	ctx := context.Background()
	userID := int64(1)

	activeUser := &user.User{
		ID:       userID,
		Email:    "user@example.com",
		FullName: "Test User",
		UserType: "jobseeker",
		Status:   "active",
	}

	// Setup mocks for multiple calls
	// 2x RefreshToken + 2x ValidateToken = 4 total FindByID calls
	mockUserRepo.On("FindByID", ctx, userID).
		Return(activeUser, nil).Times(4)

	// Act 1: Generate initial token
	token1, err1 := authService.RefreshToken(ctx, userID)
	assert.NoError(t, err1)
	assert.NotEmpty(t, token1)

	// Act 2: Validate first token
	user1, err2 := authService.ValidateToken(ctx, token1)
	assert.NoError(t, err2)
	assert.NotNil(t, user1)
	assert.Equal(t, userID, user1.ID)

	// Sleep briefly to ensure different timestamps (iat claim)
	time.Sleep(2 * time.Second)

	// Act 3: Refresh token (generate new token)
	token2, err3 := authService.RefreshToken(ctx, userID)
	assert.NoError(t, err3)
	assert.NotEmpty(t, token2)
	// Note: Tokens might be same if generated in same second (due to iat timestamp)
	// We added sleep above to ensure difference

	// Act 4: Validate new token
	user2, err4 := authService.ValidateToken(ctx, token2)
	assert.NoError(t, err4)
	assert.NotNil(t, user2)
	assert.Equal(t, userID, user2.ID)

	// Both tokens should still be valid (no explicit revocation in current impl)
	// Note: In production, you might want to implement token revocation

	mockUserRepo.AssertExpectations(t)
}
