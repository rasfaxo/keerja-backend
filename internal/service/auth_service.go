package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/utils"
)

var (
	ErrInvalidCredentials       = errors.New("invalid email or password")
	ErrEmailAlreadyExists       = errors.New("email already exists")
	ErrInvalidVerificationToken = errors.New("invalid verification token")
	ErrTokenExpired             = errors.New("token has expired")
	ErrUserNotFound             = errors.New("user not found")
	ErrInvalidResetToken        = errors.New("invalid reset token")
	ErrEmailNotVerified         = errors.New("email not verified")
)

type TokenStore interface {
	SaveVerificationToken(email, token string, expiry time.Time) error
	GetVerificationToken(token string) (email string, err error)
	DeleteVerificationToken(token string) error
	SaveResetToken(email, token string, expiry time.Time) error
	GetResetToken(token string) (email string, err error)
	DeleteResetToken(token string) error
}

// inMemoryTokenStore is a simple in-memory token store (should be replaced with Redis in production)
type inMemoryTokenStore struct {
	mu                 sync.RWMutex
	verificationTokens map[string]tokenData
	resetTokens        map[string]tokenData
}

type tokenData struct {
	Email  string
	Expiry time.Time
}

// NewInMemoryTokenStore creates a new in-memory token store
func NewInMemoryTokenStore() TokenStore {
	return &inMemoryTokenStore{
		verificationTokens: make(map[string]tokenData),
		resetTokens:        make(map[string]tokenData),
	}
}

func (s *inMemoryTokenStore) SaveVerificationToken(email, token string, expiry time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.verificationTokens[token] = tokenData{Email: email, Expiry: expiry}
	return nil
}

func (s *inMemoryTokenStore) GetVerificationToken(token string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.verificationTokens[token]
	if !exists {
		return "", ErrInvalidVerificationToken
	}

	if time.Now().After(data.Expiry) {
		return "", ErrTokenExpired
	}

	return data.Email, nil
}

func (s *inMemoryTokenStore) DeleteVerificationToken(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.verificationTokens, token)
	return nil
}

func (s *inMemoryTokenStore) SaveResetToken(email, token string, expiry time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resetTokens[token] = tokenData{Email: email, Expiry: expiry}
	return nil
}

func (s *inMemoryTokenStore) GetResetToken(token string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.resetTokens[token]
	if !exists {
		return "", ErrInvalidResetToken
	}

	if time.Now().After(data.Expiry) {
		return "", ErrTokenExpired
	}

	return data.Email, nil
}

func (s *inMemoryTokenStore) DeleteResetToken(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.resetTokens, token)
	return nil
}

// EmailService defines interface for email operations
type EmailService interface {
	SendVerificationEmail(ctx context.Context, email, token string) error
	SendPasswordResetEmail(ctx context.Context, email, token string) error
	SendWelcomeEmail(ctx context.Context, email, name string) error
}

// AuthService handles authentication business logic
type AuthService struct {
	userRepo     user.UserRepository
	emailService EmailService
	tokenStore   TokenStore
	jwtSecret    string
	jwtDuration  time.Duration
}

// AuthServiceConfig holds auth service configuration
type AuthServiceConfig struct {
	JWTSecret   string
	JWTDuration time.Duration
}

// NewAuthService creates a new auth service instance
func NewAuthService(userRepo user.UserRepository, emailService EmailService, tokenStore TokenStore, cfg AuthServiceConfig) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		emailService: emailService,
		tokenStore:   tokenStore,
		jwtSecret:    cfg.JWTSecret,
		jwtDuration:  cfg.JWTDuration,
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req *user.RegisterRequest) (*user.User, string, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, "", ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate verification token
	verificationToken, err := generateSecureToken(32)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Create user
	newUser := &user.User{
		FullName:     req.FullName,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		UserType:     req.UserType,
		IsVerified:   false,
		Status:       "inactive", // Will be set to 'active' after email verification
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	// Create user profile with slug
	slug := utils.GenerateSlug(req.FullName)
	slug = s.ensureUniqueSlug(ctx, slug)

	profile := &user.UserProfile{
		UserID: newUser.ID,
		Slug:   &slug,
	}
	if err := s.userRepo.CreateProfile(ctx, profile); err != nil {
		// Log error but don't fail registration
		fmt.Printf("Failed to create profile: %v\n", err)
	}

	// Save verification token to store
	if err := s.tokenStore.SaveVerificationToken(newUser.Email, verificationToken, time.Now().Add(24*time.Hour)); err != nil {
		return nil, "", fmt.Errorf("failed to save verification token: %w", err)
	}

	// Send verification email
	if err := s.emailService.SendVerificationEmail(ctx, newUser.Email, verificationToken); err != nil {
		// Log error but don't fail registration
		fmt.Printf("Failed to send verification email: %v\n", err)
	}

	return newUser, verificationToken, nil
}

// Login authenticates a user and returns JWT token
func (s *AuthService) Login(ctx context.Context, email, password string) (*user.User, string, error) {
	// Find user by email
	usr, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || usr == nil {
		return nil, "", ErrInvalidCredentials
	}

	// Verify password
	if !utils.VerifyPassword(password, usr.PasswordHash) {
		return nil, "", ErrInvalidCredentials
	}

	// Check if email is verified
	if !usr.IsVerified {
		return nil, "", ErrEmailNotVerified
	}

	// Check account status
	if usr.Status == "suspended" {
		return nil, "", errors.New("account is suspended")
	}
	if usr.Status == "deactivated" {
		return nil, "", errors.New("account is deactivated")
	}

	// Generate JWT token
	token, err := utils.GenerateAccessToken(usr.ID, usr.Email, usr.UserType, s.jwtSecret, s.jwtDuration)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Update last login
	now := time.Now()
	usr.LastLogin = &now
	if err := s.userRepo.Update(ctx, usr); err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	return usr, token, nil
}

// VerifyEmail verifies user's email with token
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	// Get email from token store
	email, err := s.tokenStore.GetVerificationToken(token)
	if err != nil {
		return err
	}

	// Find user by email
	usr, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || usr == nil {
		return ErrUserNotFound
	}

	// Check if already verified
	if usr.IsVerified {
		return nil // Already verified, no error
	}

	// Update user as verified
	usr.IsVerified = true
	usr.Status = "active"

	if err := s.userRepo.Update(ctx, usr); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Delete verification token
	s.tokenStore.DeleteVerificationToken(token)

	// Send welcome email
	if err := s.emailService.SendWelcomeEmail(ctx, usr.Email, usr.FullName); err != nil {
		// Log error but don't fail verification
		fmt.Printf("Failed to send welcome email: %v\n", err)
	}

	return nil
}

// ResendVerificationEmail resends verification email
func (s *AuthService) ResendVerificationEmail(ctx context.Context, email string) error {
	// Find user by email
	usr, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || usr == nil {
		return ErrUserNotFound
	}

	// Check if already verified
	if usr.IsVerified {
		return errors.New("email already verified")
	}

	// Generate new verification token
	verificationToken, err := generateSecureToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Save to token store
	if err := s.tokenStore.SaveVerificationToken(usr.Email, verificationToken, time.Now().Add(24*time.Hour)); err != nil {
		return fmt.Errorf("failed to save verification token: %w", err)
	}

	// Send verification email
	if err := s.emailService.SendVerificationEmail(ctx, usr.Email, verificationToken); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

// ForgotPassword initiates password reset process
func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	// Find user by email
	usr, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || usr == nil {
		// Don't reveal if user exists or not
		return nil
	}

	// Generate reset token
	resetToken, err := generateSecureToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Save to token store
	if err := s.tokenStore.SaveResetToken(usr.Email, resetToken, time.Now().Add(1*time.Hour)); err != nil {
		return fmt.Errorf("failed to save reset token: %w", err)
	}

	// Send password reset email
	if err := s.emailService.SendPasswordResetEmail(ctx, usr.Email, resetToken); err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}

// ResetPassword resets user password with token
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Get email from token store
	email, err := s.tokenStore.GetResetToken(token)
	if err != nil {
		return err
	}

	// Find user by email
	usr, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || usr == nil {
		return ErrUserNotFound
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	usr.PasswordHash = hashedPassword

	if err := s.userRepo.Update(ctx, usr); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Delete reset token
	s.tokenStore.DeleteResetToken(token)

	return nil
}

// ChangePassword changes user password (requires current password)
func (s *AuthService) ChangePassword(ctx context.Context, userID int64, currentPassword, newPassword string) error {
	// Find user
	usr, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || usr == nil {
		return ErrUserNotFound
	}

	// Verify current password
	if !utils.VerifyPassword(currentPassword, usr.PasswordHash) {
		return errors.New("invalid current password")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	usr.PasswordHash = hashedPassword

	if err := s.userRepo.Update(ctx, usr); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// RefreshToken generates a new access token
func (s *AuthService) RefreshToken(ctx context.Context, userID int64) (string, error) {
	// Find user
	usr, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || usr == nil {
		return "", ErrUserNotFound
	}

	// Check account status
	if usr.Status != "active" {
		return "", errors.New("account is not active")
	}

	// Generate new token
	token, err := utils.GenerateAccessToken(usr.ID, usr.Email, usr.UserType, s.jwtSecret, s.jwtDuration)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// ValidateToken validates JWT token and returns user
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*user.User, error) {
	// Validate and parse token
	claims, err := utils.ValidateToken(tokenString, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Find user
	usr, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil || usr == nil {
		return nil, ErrUserNotFound
	}

	// Check account status
	if usr.Status != "active" {
		return nil, errors.New("account is not active")
	}

	return usr, nil
}

// Logout handles user logout (in stateless JWT, this is mainly for cleanup)
func (s *AuthService) Logout(ctx context.Context, userID int64) error {
	// In JWT-based auth, logout is typically handled client-side
	// nanti disini bisa implement token blacklist pake redis

	// For now, just verify user exists
	usr, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || usr == nil {
		return ErrUserNotFound
	}

	// bisa implement token blacklist disini kalo diperlukan
	// e.g., redis.Set(tokenHash, "blacklisted", expirationTime)

	return nil
}

// Helper functions

// ensureUniqueSlug ensures slug is unique by appending number if needed
func (s *AuthService) ensureUniqueSlug(ctx context.Context, baseSlug string) string {
	slug := baseSlug
	counter := 1

	for {
		// Check if slug exists
		existingProfile, err := s.userRepo.FindProfileBySlug(ctx, slug)
		if err != nil || existingProfile == nil {
			// Slug is unique
			return slug
		}

		// Slug exists, try with counter
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++

		// Safety limit
		if counter > 1000 {
			// Append random string
			randomStr, _ := generateSecureToken(4)
			return fmt.Sprintf("%s-%s", baseSlug, randomStr)
		}
	}
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// timePtr returns a pointer to time.Time
func timePtr(t time.Time) *time.Time {
	return &t
}
