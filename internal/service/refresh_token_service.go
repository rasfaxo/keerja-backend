package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"keerja-backend/internal/domain/auth"
	"keerja-backend/internal/utils"
)

// Refresh token configuration
const (
	RefreshTokenLength     = 64   // 64 bytes = 512 bits
	RefreshTokenExpiryDays = 30   // Default 30 days
	RememberMeExpiryDays   = 90   // Remember me: 90 days
	MaxActiveTokensPerUser = 5    // Max 5 devices
	RefreshTokenRotation   = true // Enable token rotation on refresh
)

var (
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrRefreshTokenExpired  = errors.New("refresh token has expired")
	ErrRefreshTokenRevoked  = errors.New("refresh token has been revoked")
	ErrMaxDevicesExceeded   = errors.New("maximum number of devices exceeded")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

// DeviceInfo contains information about the device making the request
type DeviceInfo struct {
	DeviceName string
	DeviceType string // mobile, desktop, tablet, unknown
	DeviceID   string // unique identifier
	UserAgent  string
	IPAddress  string
}

// RefreshTokenService handles refresh token operations
type RefreshTokenService struct {
	refreshTokenRepo auth.RefreshTokenRepository
	jwtSecret        string
	jwtDuration      time.Duration
}

// NewRefreshTokenService creates a new refresh token service
func NewRefreshTokenService(
	refreshTokenRepo auth.RefreshTokenRepository,
	jwtSecret string,
	jwtDuration time.Duration,
) *RefreshTokenService {
	return &RefreshTokenService{
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        jwtSecret,
		jwtDuration:      jwtDuration,
	}
}

// generateRefreshToken generates a cryptographically secure random token
func (s *RefreshTokenService) generateRefreshToken() (string, error) {
	bytes := make([]byte, RefreshTokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// hashRefreshToken creates SHA256 hash of refresh token
func (s *RefreshTokenService) hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// parseDeviceType determines device type from user agent
func (s *RefreshTokenService) parseDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)

	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") ||
		strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		return "mobile"
	}

	if strings.Contains(ua, "tablet") {
		return "tablet"
	}

	if strings.Contains(ua, "windows") || strings.Contains(ua, "mac") ||
		strings.Contains(ua, "linux") {
		return "desktop"
	}

	return "unknown"
}

// parseDeviceName creates user-friendly device name from user agent
func (s *RefreshTokenService) parseDeviceName(userAgent string) string {
	ua := strings.ToLower(userAgent)

	// Check for specific browsers
	browser := "Unknown Browser"
	if strings.Contains(ua, "chrome") && !strings.Contains(ua, "edg") {
		browser = "Chrome"
	} else if strings.Contains(ua, "firefox") {
		browser = "Firefox"
	} else if strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome") {
		browser = "Safari"
	} else if strings.Contains(ua, "edg") {
		browser = "Edge"
	}

	// Check for OS
	os := "Unknown OS"
	if strings.Contains(ua, "windows") {
		os = "Windows"
	} else if strings.Contains(ua, "mac") {
		os = "macOS"
	} else if strings.Contains(ua, "linux") {
		os = "Linux"
	} else if strings.Contains(ua, "android") {
		os = "Android"
	} else if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		os = "iOS"
	}

	return fmt.Sprintf("%s on %s", browser, os)
}

// CreateRefreshToken creates a new refresh token for a user
func (s *RefreshTokenService) CreateRefreshToken(
	ctx context.Context,
	userID int64,
	deviceInfo DeviceInfo,
	rememberMe bool,
) (string, error) {
	// Check active token limit per user
	activeCount, err := s.refreshTokenRepo.CountActiveByUserID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to count active tokens: %w", err)
	}

	if activeCount >= MaxActiveTokensPerUser {
		// Revoke oldest token to make room
		tokens, err := s.refreshTokenRepo.FindActiveByUserID(ctx, userID)
		if err != nil {
			return "", fmt.Errorf("failed to find active tokens: %w", err)
		}

		if len(tokens) > 0 {
			oldestToken := tokens[len(tokens)-1]
			if err := s.refreshTokenRepo.Revoke(ctx, oldestToken.ID, "max_devices_exceeded"); err != nil {
				return "", fmt.Errorf("failed to revoke oldest token: %w", err)
			}
		}
	}

	// Generate refresh token
	token, err := s.generateRefreshToken()
	if err != nil {
		return "", err
	}

	tokenHash := s.hashRefreshToken(token)

	// Determine expiry based on remember me
	expiryDays := RefreshTokenExpiryDays
	if rememberMe {
		expiryDays = RememberMeExpiryDays
	}
	expiresAt := time.Now().Add(time.Duration(expiryDays) * 24 * time.Hour)

	// Parse device info
	deviceType := s.parseDeviceType(deviceInfo.UserAgent)
	deviceName := s.parseDeviceName(deviceInfo.UserAgent)

	// If device info not provided, use parsed values
	if deviceInfo.DeviceType == "" {
		deviceInfo.DeviceType = deviceType
	}
	if deviceInfo.DeviceName == "" {
		deviceInfo.DeviceName = deviceName
	}

	// Create refresh token record
	refreshToken := &auth.RefreshToken{
		UserID:     userID,
		TokenHash:  tokenHash,
		DeviceName: &deviceInfo.DeviceName,
		DeviceType: &deviceInfo.DeviceType,
		DeviceID:   &deviceInfo.DeviceID,
		UserAgent:  &deviceInfo.UserAgent,
		IPAddress:  &deviceInfo.IPAddress,
		ExpiresAt:  expiresAt,
		Revoked:    false,
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return token, nil
}

// RefreshAccessToken validates refresh token and issues new access token
func (s *RefreshTokenService) RefreshAccessToken(
	ctx context.Context,
	refreshToken string,
	userID int64,
	email string,
	userType string,
) (string, string, error) {
	// Hash the provided token
	tokenHash := s.hashRefreshToken(refreshToken)

	// Find refresh token in database
	storedToken, err := s.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return "", "", fmt.Errorf("failed to find refresh token: %w", err)
	}

	if storedToken == nil {
		return "", "", ErrRefreshTokenNotFound
	}

	// Validate token
	if storedToken.Revoked {
		return "", "", ErrRefreshTokenRevoked
	}

	if storedToken.IsExpired() {
		return "", "", ErrRefreshTokenExpired
	}

	if storedToken.UserID != userID {
		return "", "", ErrInvalidRefreshToken
	}

	// Update last used timestamp
	if err := s.refreshTokenRepo.UpdateLastUsed(ctx, storedToken.ID); err != nil {
		return "", "", fmt.Errorf("failed to update last used: %w", err)
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(userID, email, userType, s.jwtSecret, s.jwtDuration)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Optional: Token rotation - generate new refresh token
	newRefreshToken := refreshToken
	if RefreshTokenRotation {
		newRefreshToken, err = s.generateRefreshToken()
		if err != nil {
			return "", "", fmt.Errorf("failed to generate new refresh token: %w", err)
		}

		// Update token hash
		storedToken.TokenHash = s.hashRefreshToken(newRefreshToken)
		if err := s.refreshTokenRepo.Update(ctx, storedToken); err != nil {
			return "", "", fmt.Errorf("failed to rotate refresh token: %w", err)
		}
	}

	return accessToken, newRefreshToken, nil
}

// RevokeRefreshToken revokes a specific refresh token
func (s *RefreshTokenService) RevokeRefreshToken(ctx context.Context, refreshToken string, reason string) error {
	tokenHash := s.hashRefreshToken(refreshToken)

	storedToken, err := s.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to find refresh token: %w", err)
	}

	if storedToken == nil {
		return ErrRefreshTokenNotFound
	}

	return s.refreshTokenRepo.Revoke(ctx, storedToken.ID, reason)
}

// RevokeAllUserTokens revokes all refresh tokens for a user (logout all devices)
func (s *RefreshTokenService) RevokeAllUserTokens(ctx context.Context, userID int64, reason string) error {
	return s.refreshTokenRepo.RevokeAllByUserID(ctx, userID, reason)
}

// RevokeDeviceToken revokes refresh token for specific device
func (s *RefreshTokenService) RevokeDeviceToken(ctx context.Context, userID int64, deviceID string, reason string) error {
	return s.refreshTokenRepo.RevokeByDeviceID(ctx, userID, deviceID, reason)
}

// GetUserDevices returns all active devices for a user
func (s *RefreshTokenService) GetUserDevices(ctx context.Context, userID int64) ([]auth.RefreshToken, error) {
	return s.refreshTokenRepo.FindActiveByUserID(ctx, userID)
}

// CleanupExpiredTokens removes expired refresh tokens (should be called periodically)
func (s *RefreshTokenService) CleanupExpiredTokens(ctx context.Context) error {
	return s.refreshTokenRepo.DeleteExpired(ctx)
}

// CleanupRevokedTokens removes old revoked tokens (cleanup after 90 days)
func (s *RefreshTokenService) CleanupRevokedTokens(ctx context.Context) error {
	ninetyDaysAgo := time.Now().Add(-90 * 24 * time.Hour)
	return s.refreshTokenRepo.DeleteRevoked(ctx, ninetyDaysAgo)
}
