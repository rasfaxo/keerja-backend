package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	mathRand "math/rand"
	"time"

	"keerja-backend/internal/domain/auth"
	"keerja-backend/internal/domain/email"
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

// OTP Configuration constants for registration
const (
	OTPCodeLength            = 6
	OTPCodeExpiryMinutes     = 5
	OTPMaxVerifyAttempts     = 5
	OTPResendWindowSeconds   = 60 // user must wait 60s before resend
	OTPMaxOTPRequestsPerHour = 3  // max 3 OTP requests per hour per user
)

var (
	ErrInvalidOTPCode      = errors.New("invalid OTP code")
	ErrOTPCodeExpired      = errors.New("OTP code has expired")
	ErrOTPCodeAlreadyUsed  = errors.New("OTP code has already been used")
	ErrTooManyOTPAttempts  = errors.New("too many failed OTP verification attempts")
	ErrTooManyOTPRequests  = errors.New("too many OTP requests, please try again later")
	ErrOTPCodeNotFound     = errors.New("no OTP found for this user")
	ErrResendTooSoon       = errors.New("please wait before requesting a new OTP")
	ErrUserAlreadyVerified = errors.New("user email already verified")
)

// RegistrationService handles user registration with OTP verification
type RegistrationService struct {
	userRepo     user.UserRepository
	otpCodeRepo  auth.OTPCodeRepository
	emailService email.EmailService
	jwtSecret    string
	jwtDuration  time.Duration
}

// NewRegistrationService creates a new registration service
func NewRegistrationService(
	userRepo user.UserRepository,
	otpCodeRepo auth.OTPCodeRepository,
	emailService email.EmailService,
	jwtSecret string,
	jwtDuration time.Duration,
) *RegistrationService {
	return &RegistrationService{
		userRepo:     userRepo,
		otpCodeRepo:  otpCodeRepo,
		emailService: emailService,
		jwtSecret:    jwtSecret,
		jwtDuration:  jwtDuration,
	}
}

// generateOTPCode generates a random 6-digit OTP code
func (s *RegistrationService) generateOTPCode() string {
	rnd := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	code := rnd.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

// hashOTPCode creates a SHA256 hash of the OTP code
func (s *RegistrationService) hashOTPCode(email, otpCode string) string {
	hash := sha256.Sum256([]byte(email + "|" + otpCode))
	return hex.EncodeToString(hash[:])
}

// RegisterUser creates a new user with is_verified = false and sends OTP
func (s *RegistrationService) RegisterUser(ctx context.Context, fullName, email, password, phone, userType string) error {
	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return ErrEmailAlreadyExists
	}

	// Hash password
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user with is_verified = false
	newUser := &user.User{
		FullName:     fullName,
		Email:        email,
		Phone:        &phone,
		PasswordHash: passwordHash,
		UserType:     userType,
		IsVerified:   false,
		Status:       "active",
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Create empty user profile for the new user
	userProfile := &user.UserProfile{
		UserID: newUser.ID,
	}
	if err := s.userRepo.CreateProfile(ctx, userProfile); err != nil {
		return fmt.Errorf("failed to create user profile: %w", err)
	}

	// Generate and send OTP
	if err := s.sendOTPToUser(ctx, newUser.ID, email); err != nil {
		// Rollback user creation if OTP sending fails (optional)
		// return error but user is already created
		return fmt.Errorf("user created but failed to send OTP: %w", err)
	}

	return nil
}

// sendOTPToUser generates OTP, saves to DB, and sends via email
func (s *RegistrationService) sendOTPToUser(ctx context.Context, userID int64, email string) error {
	// Check rate limiting: max 3 OTP requests per hour
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	recentCount, err := s.otpCodeRepo.CountRecentByUserID(ctx, userID, oneHourAgo, "email_verification")
	if err != nil {
		return fmt.Errorf("failed to check rate limit: %w", err)
	}
	if recentCount >= OTPMaxOTPRequestsPerHour {
		return ErrTooManyOTPRequests
	}

	// Check resend window: user must wait 60s before new OTP
	latestOTP, err := s.otpCodeRepo.FindByUserIDAndType(ctx, userID, "email_verification")
	if err != nil {
		return fmt.Errorf("failed to check latest OTP: %w", err)
	}
	if latestOTP != nil && !latestOTP.IsUsed {
		timeSinceCreation := time.Since(latestOTP.CreatedAt)
		if timeSinceCreation < OTPResendWindowSeconds*time.Second {
			return ErrResendTooSoon
		}
	}

	// Generate OTP code
	otpCode := s.generateOTPCode()
	otpHash := s.hashOTPCode(email, otpCode)

	// Save OTP to database
	otpRecord := &auth.OTPCode{
		UserID:    userID,
		OTPHash:   otpHash,
		Type:      "email_verification",
		ExpiredAt: time.Now().Add(OTPCodeExpiryMinutes * time.Minute),
		IsUsed:    false,
		Attempts:  0,
	}

	if err := s.otpCodeRepo.Create(ctx, otpRecord); err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}

	// Send OTP via email
	if err := s.emailService.SendOTPRegistrationEmail(ctx, email, "", otpCode); err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	return nil
}

// VerifyEmailOTP verifies OTP code and marks user as verified
func (s *RegistrationService) VerifyEmailOTP(ctx context.Context, email, otpCode string) (string, *user.User, error) {
	// Find user by email
	usr, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, fmt.Errorf("failed to find user: %w", err)
	}
	if usr == nil {
		return "", nil, ErrUserNotFound
	}

	// Check if already verified
	if usr.IsVerified {
		return "", nil, ErrUserAlreadyVerified
	}

	// Find latest OTP for this user
	latestOTP, err := s.otpCodeRepo.FindByUserIDAndType(ctx, usr.ID, "email_verification")
	if err != nil {
		return "", nil, fmt.Errorf("failed to find OTP: %w", err)
	}
	if latestOTP == nil {
		return "", nil, ErrOTPCodeNotFound
	}

	// Check if OTP is expired
	if latestOTP.IsExpired() {
		return "", nil, ErrOTPCodeExpired
	}

	// Check if OTP is already used
	if latestOTP.IsUsed {
		return "", nil, ErrOTPCodeAlreadyUsed
	}

	// Check max attempts
	if !latestOTP.CanAttemptVerification(OTPMaxVerifyAttempts) {
		return "", nil, ErrTooManyOTPAttempts
	}

	// Verify OTP hash
	inputHash := s.hashOTPCode(email, otpCode)
	if inputHash != latestOTP.OTPHash {
		// Increment attempts on failed verification
		if err := s.otpCodeRepo.IncrementAttempts(ctx, latestOTP.ID); err != nil {
			return "", nil, fmt.Errorf("failed to increment attempts: %w", err)
		}
		return "", nil, ErrInvalidOTPCode
	}

	// Mark OTP as used
	if err := s.otpCodeRepo.MarkAsUsed(ctx, latestOTP.ID); err != nil {
		return "", nil, fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	// Update user: set is_verified = true
	usr.IsVerified = true
	if err := s.userRepo.Update(ctx, usr); err != nil {
		return "", nil, fmt.Errorf("failed to update user verification status: %w", err)
	}

	// Generate JWT token for auto-login
	token, err := utils.GenerateAccessToken(usr.ID, usr.Email, usr.UserType, s.jwtSecret, s.jwtDuration)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, usr, nil
}

// ResendOTP resends OTP to user email
func (s *RegistrationService) ResendOTP(ctx context.Context, email string) error {
	// Find user by email
	usr, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if usr == nil {
		return ErrUserNotFound
	}

	// Check if already verified
	if usr.IsVerified {
		return ErrUserAlreadyVerified
	}

	// Send OTP (rate limiting is handled inside)
	if err := s.sendOTPToUser(ctx, usr.ID, email); err != nil {
		return err
	}

	return nil
}

// CleanupExpiredOTPs removes expired OTP codes (should be called periodically)
func (s *RegistrationService) CleanupExpiredOTPs(ctx context.Context) error {
	return s.otpCodeRepo.DeleteExpired(ctx)
}

// ===========================================
// Forgot Password with OTP
// ===========================================

// RequestPasswordResetOTP sends OTP for password reset
func (s *RegistrationService) RequestPasswordResetOTP(ctx context.Context, email string) error {
	// Find user by email
	usr, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Don't reveal if user doesn't exist (security best practice)
	if usr == nil {
		return nil // Silent success
	}

	// Check if user is verified
	if !usr.IsVerified {
		return ErrEmailNotVerified
	}

	// Check rate limiting - max 3 password reset requests per hour
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	recentOTPs, err := s.otpCodeRepo.FindAllByUserIDAndType(ctx, usr.ID, "password_reset")
	if err != nil {
		return fmt.Errorf("failed to check recent OTPs: %w", err)
	}

	count := 0
	for _, otp := range recentOTPs {
		if otp.CreatedAt.After(oneHourAgo) {
			count++
		}
	}

	if count >= 3 {
		return ErrTooManyOTPRequests
	}

	// Revoke all existing password reset OTPs for this user
	for _, otp := range recentOTPs {
		if !otp.IsUsed {
			otp.IsUsed = true
			now := time.Now()
			otp.UsedAt = &now
			if err := s.otpCodeRepo.Update(ctx, otp); err != nil {
				// Log error but continue
				continue
			}
		}
	}

	// Generate OTP code
	otpCode := s.generateOTPCode()
	otpHash := s.hashOTPCode(email, otpCode)

	// Create OTP record
	otp := &auth.OTPCode{
		UserID:    usr.ID,
		OTPHash:   otpHash,
		Type:      "password_reset",
		ExpiredAt: time.Now().Add(5 * time.Minute), // 5 minutes expiry
		IsUsed:    false,
		Attempts:  0,
	}

	if err := s.otpCodeRepo.Create(ctx, otp); err != nil {
		return fmt.Errorf("failed to create OTP: %w", err)
	}

	// Send OTP email
	subject := "Password Reset - Keerja"
	body := fmt.Sprintf(`
		<h2>Password Reset Request</h2>
		<p>Hello %s,</p>
		<p>You have requested to reset your password. Please use the following OTP code:</p>
		<h1 style="color: #4F46E5; letter-spacing: 5px;">%s</h1>
		<p>This code will expire in <strong>5 minutes</strong>.</p>
		<p>If you didn't request this, please ignore this email.</p>
		<hr>
		<p style="color: #666; font-size: 12px;">Keerja - Job Portal Platform</p>
	`, usr.FullName, otpCode)

	if err := s.emailService.SendEmail(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	return nil
}

// ResetPasswordWithOTP resets password using OTP verification
func (s *RegistrationService) ResetPasswordWithOTP(ctx context.Context, email, otpCode, newPassword string) error {
	// Find user by email
	usr, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	if usr == nil {
		return ErrInvalidCredentials
	}

	// Hash the provided OTP
	otpHash := s.hashOTPCode(email, otpCode)

	// Find active password reset OTP
	otps, err := s.otpCodeRepo.FindAllByUserIDAndType(ctx, usr.ID, "password_reset")
	if err != nil {
		return fmt.Errorf("failed to find OTP: %w", err)
	}

	var validOTP *auth.OTPCode
	for _, otp := range otps {
		if !otp.IsUsed && otp.OTPHash == otpHash {
			validOTP = otp
			break
		}
	}

	if validOTP == nil {
		return ErrInvalidOTPCode
	}

	// Check if OTP is expired
	if time.Now().After(validOTP.ExpiredAt) {
		return ErrOTPCodeExpired
	}

	// Check attempts (max 5)
	if validOTP.Attempts >= 5 {
		return ErrTooManyOTPAttempts
	}

	// Mark OTP as used
	validOTP.IsUsed = true
	now := time.Now()
	validOTP.UsedAt = &now

	if err := s.otpCodeRepo.Update(ctx, validOTP); err != nil {
		return fmt.Errorf("failed to update OTP: %w", err)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	usr.PasswordHash = string(hashedPassword)
	if err := s.userRepo.Update(ctx, usr); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Send confirmation email
	subject := "Password Changed Successfully - Keerja"
	body := fmt.Sprintf(`
		<h2>Password Changed</h2>
		<p>Hello %s,</p>
		<p>Your password has been successfully changed.</p>
		<p>If you didn't make this change, please contact support immediately.</p>
		<hr>
		<p style="color: #666; font-size: 12px;">Keerja - Job Portal Platform</p>
	`, usr.FullName)

	// Send email asynchronously (ignore error)
	go s.emailService.SendEmail(context.Background(), usr.Email, subject, body)

	return nil
}
