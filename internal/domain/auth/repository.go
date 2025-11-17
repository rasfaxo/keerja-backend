package auth

import (
	"context"
	"time"
)

// OAuthRepository defines data access methods for OAuth providers
type OAuthRepository interface {
	// Create creates a new OAuth provider connection
	Create(ctx context.Context, provider *OAuthProvider) error

	// FindByProviderAndUserID finds OAuth connection by provider and provider user ID
	FindByProviderAndUserID(ctx context.Context, provider, providerUserID string) (*OAuthProvider, error)

	// FindByUserID finds all OAuth connections for a user
	FindByUserID(ctx context.Context, userID int64) ([]OAuthProvider, error)

	// Update updates OAuth provider connection
	Update(ctx context.Context, provider *OAuthProvider) error

	// Delete deletes OAuth provider connection
	Delete(ctx context.Context, id int64) error

	// FindByUserAndProvider finds OAuth connection for specific user and provider
	FindByUserAndProvider(ctx context.Context, userID int64, provider string) (*OAuthProvider, error)
}

// OTPCodeRepository defines data access methods for OTP verification codes
type OTPCodeRepository interface {
	// Create creates a new OTP code
	Create(ctx context.Context, otp *OTPCode) error

	// FindByUserIDAndType finds the latest OTP code by user ID and type
	FindByUserIDAndType(ctx context.Context, userID int64, otpType string) (*OTPCode, error)

	// FindAllByUserIDAndType finds all OTP codes by user ID and type
	FindAllByUserIDAndType(ctx context.Context, userID int64, otpType string) ([]*OTPCode, error)

	// MarkAsUsed marks an OTP code as used
	MarkAsUsed(ctx context.Context, id int64) error

	// IncrementAttempts increments failed verification attempts
	IncrementAttempts(ctx context.Context, id int64) error

	// DeleteExpired deletes all expired OTP codes
	DeleteExpired(ctx context.Context) error

	// CountRecentByUserID counts recent OTP requests by user
	CountRecentByUserID(ctx context.Context, userID int64, since time.Time, otpType string) (int64, error)

	// Update updates an OTP code
	Update(ctx context.Context, otp *OTPCode) error
}

// RefreshTokenRepository defines data access methods for refresh tokens
type RefreshTokenRepository interface {
	// Create creates a new refresh token
	Create(ctx context.Context, token *RefreshToken) error

	// FindByTokenHash finds refresh token by hash
	FindByTokenHash(ctx context.Context, tokenHash string) (*RefreshToken, error)

	// FindByUserID finds all refresh tokens for a user
	FindByUserID(ctx context.Context, userID int64) ([]RefreshToken, error)

	// FindActiveByUserID finds all active (non-revoked, non-expired) tokens for user
	FindActiveByUserID(ctx context.Context, userID int64) ([]RefreshToken, error)

	// FindByUserAndDevice finds refresh token by user and device
	FindByUserAndDevice(ctx context.Context, userID int64, deviceID string) (*RefreshToken, error)

	// Update updates refresh token
	Update(ctx context.Context, token *RefreshToken) error

	// UpdateLastUsed updates last used timestamp
	UpdateLastUsed(ctx context.Context, id int64) error

	// Revoke revokes a refresh token
	Revoke(ctx context.Context, id int64, reason string) error

	// RevokeAllByUserID revokes all refresh tokens for a user
	RevokeAllByUserID(ctx context.Context, userID int64, reason string) error

	// RevokeByDeviceID revokes refresh token by device ID
	RevokeByDeviceID(ctx context.Context, userID int64, deviceID string, reason string) error

	// DeleteExpired deletes all expired refresh tokens
	DeleteExpired(ctx context.Context) error

	// DeleteRevoked deletes old revoked tokens (cleanup)
	DeleteRevoked(ctx context.Context, olderThan time.Time) error

	// CountActiveByUserID counts active tokens for user
	CountActiveByUserID(ctx context.Context, userID int64) (int64, error)
}
