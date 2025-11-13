package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

// Token types
const (
	TokenTypeVerification  = "verification"
	TokenTypePasswordReset = "password_reset"
)

// TokenData represents token information
type TokenData struct {
	Token     string
	Type      string
	UserID    int64
	ExpiresAt time.Time
	CreatedAt time.Time
}

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateVerificationToken generates a token for email verification
func GenerateVerificationToken() (string, error) {
	return GenerateSecureToken(32)
}

// GeneratePasswordResetToken generates a token for password reset
func GeneratePasswordResetToken() (string, error) {
	return GenerateSecureToken(32)
}

// GenerateAPIKey generates an API key
func GenerateAPIKey() (string, error) {
	token, err := GenerateSecureToken(32)
	if err != nil {
		return "", err
	}
	return "keerja_" + token, nil
}

// IsExpired checks if a timestamp has expired
func IsExpired(expiresAt time.Time) bool {
	return time.Now().After(expiresAt)
}

// GetVerificationTokenExpiry returns expiration time for verification tokens (24 hours)
func GetVerificationTokenExpiry() time.Time {
	return time.Now().Add(24 * time.Hour)
}

// GetPasswordResetTokenExpiry returns expiration time for password reset tokens (1 hour)
func GetPasswordResetTokenExpiry() time.Time {
	return time.Now().Add(1 * time.Hour)
}

// GetRefreshTokenExpiry returns expiration time for refresh tokens (7 days)
func GetRefreshTokenExpiry() time.Time {
	return time.Now().Add(7 * 24 * time.Hour)
}

// GenerateRandomToken generates a random token of specified length (in bytes)
func GenerateRandomToken(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}
