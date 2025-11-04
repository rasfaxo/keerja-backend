package postgres

import (
	"context"
	"time"

	"keerja-backend/internal/domain/auth"

	"gorm.io/gorm"
)

// ===========================================
// OAuth Repository Implementation
// ===========================================

// oauthRepository implements auth.OAuthRepository
type oauthRepository struct {
	db *gorm.DB
}

// NewOAuthRepository creates a new OAuth repository
func NewOAuthRepository(db *gorm.DB) auth.OAuthRepository {
	return &oauthRepository{db: db}
}

// Create creates a new OAuth provider connection
func (r *oauthRepository) Create(ctx context.Context, provider *auth.OAuthProvider) error {
	return r.db.WithContext(ctx).Create(provider).Error
}

// FindByProviderAndUserID finds OAuth connection by provider and provider user ID
func (r *oauthRepository) FindByProviderAndUserID(ctx context.Context, provider, providerUserID string) (*auth.OAuthProvider, error) {
	var oauth auth.OAuthProvider
	err := r.db.WithContext(ctx).
		Where("provider = ? AND provider_user_id = ?", provider, providerUserID).
		First(&oauth).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &oauth, nil
}

// FindByUserID finds all OAuth connections for a user
func (r *oauthRepository) FindByUserID(ctx context.Context, userID int64) ([]auth.OAuthProvider, error) {
	var providers []auth.OAuthProvider
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Find(&providers).Error
	return providers, err
}

// Update updates OAuth provider connection
func (r *oauthRepository) Update(ctx context.Context, provider *auth.OAuthProvider) error {
	return r.db.WithContext(ctx).Save(provider).Error
}

// Delete deletes OAuth provider connection
func (r *oauthRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&auth.OAuthProvider{}, id).Error
}

// FindByUserAndProvider finds OAuth connection for specific user and provider
func (r *oauthRepository) FindByUserAndProvider(ctx context.Context, userID int64, provider string) (*auth.OAuthProvider, error) {
	var oauth auth.OAuthProvider
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND provider = ? AND is_active = ?", userID, provider, true).
		First(&oauth).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &oauth, nil
}

// ===========================================
// OTPCode Repository Implementation
// ===========================================

// otpCodeRepository implements auth.OTPCodeRepository
type otpCodeRepository struct {
	db *gorm.DB
}

// NewOTPCodeRepository creates a new OTP code repository
func NewOTPCodeRepository(db *gorm.DB) auth.OTPCodeRepository {
	return &otpCodeRepository{db: db}
}

// Create creates a new OTP code
func (r *otpCodeRepository) Create(ctx context.Context, otp *auth.OTPCode) error {
	return r.db.WithContext(ctx).Create(otp).Error
}

// FindByUserIDAndType finds the latest OTP code by user ID and type
func (r *otpCodeRepository) FindByUserIDAndType(ctx context.Context, userID int64, otpType string) (*auth.OTPCode, error) {
	var otp auth.OTPCode
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, otpType).
		Order("created_at DESC").
		First(&otp).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &otp, nil
}

// FindAllByUserIDAndType finds all OTP codes by user ID and type
func (r *otpCodeRepository) FindAllByUserIDAndType(ctx context.Context, userID int64, otpType string) ([]*auth.OTPCode, error) {
	var otps []*auth.OTPCode
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, otpType).
		Order("created_at DESC").
		Find(&otps).Error

	if err != nil {
		return nil, err
	}
	return otps, nil
}

// MarkAsUsed marks an OTP code as used
func (r *otpCodeRepository) MarkAsUsed(ctx context.Context, id int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&auth.OTPCode{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_used": true,
			"used_at": now,
		}).Error
}

// Update updates an OTP code
func (r *otpCodeRepository) Update(ctx context.Context, otp *auth.OTPCode) error {
	return r.db.WithContext(ctx).Save(otp).Error
}

// IncrementAttempts increments failed verification attempts
func (r *otpCodeRepository) IncrementAttempts(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&auth.OTPCode{}).
		Where("id = ?", id).
		UpdateColumn("attempts", gorm.Expr("attempts + ?", 1)).Error
}

// DeleteExpired deletes all expired OTP codes
func (r *otpCodeRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expired_at < ?", time.Now()).
		Delete(&auth.OTPCode{}).Error
}

// CountRecentByUserID counts recent OTP requests by user
func (r *otpCodeRepository) CountRecentByUserID(ctx context.Context, userID int64, since time.Time, otpType string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&auth.OTPCode{}).
		Where("user_id = ? AND type = ? AND created_at >= ?", userID, otpType, since).
		Count(&count).Error
	return count, err
}

// ===========================================
// RefreshToken Repository Implementation
// ===========================================

// refreshTokenRepository implements auth.RefreshTokenRepository
type refreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *gorm.DB) auth.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Create creates a new refresh token
func (r *refreshTokenRepository) Create(ctx context.Context, token *auth.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// FindByTokenHash finds refresh token by hash
func (r *refreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*auth.RefreshToken, error) {
	var token auth.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		First(&token).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// FindByUserID finds all refresh tokens for a user
func (r *refreshTokenRepository) FindByUserID(ctx context.Context, userID int64) ([]auth.RefreshToken, error) {
	var tokens []auth.RefreshToken
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tokens).Error
	return tokens, err
}

// FindActiveByUserID finds all active (non-revoked, non-expired) tokens for user
func (r *refreshTokenRepository) FindActiveByUserID(ctx context.Context, userID int64) ([]auth.RefreshToken, error) {
	var tokens []auth.RefreshToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND revoked = ? AND expires_at > ?", userID, false, time.Now()).
		Order("last_used_at DESC").
		Find(&tokens).Error
	return tokens, err
}

// FindByUserAndDevice finds refresh token by user and device
func (r *refreshTokenRepository) FindByUserAndDevice(ctx context.Context, userID int64, deviceID string) (*auth.RefreshToken, error) {
	var token auth.RefreshToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND device_id = ? AND revoked = ? AND expires_at > ?", userID, deviceID, false, time.Now()).
		Order("created_at DESC").
		First(&token).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// Update updates refresh token
func (r *refreshTokenRepository) Update(ctx context.Context, token *auth.RefreshToken) error {
	return r.db.WithContext(ctx).Save(token).Error
}

// UpdateLastUsed updates last used timestamp
func (r *refreshTokenRepository) UpdateLastUsed(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&auth.RefreshToken{}).
		Where("id = ?", id).
		Update("last_used_at", time.Now()).Error
}

// Revoke revokes a refresh token
func (r *refreshTokenRepository) Revoke(ctx context.Context, id int64, reason string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&auth.RefreshToken{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"revoked":        true,
			"revoked_at":     now,
			"revoked_reason": reason,
		}).Error
}

// RevokeAllByUserID revokes all refresh tokens for a user
func (r *refreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID int64, reason string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&auth.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Updates(map[string]interface{}{
			"revoked":        true,
			"revoked_at":     now,
			"revoked_reason": reason,
		}).Error
}

// RevokeByDeviceID revokes refresh token by device ID
func (r *refreshTokenRepository) RevokeByDeviceID(ctx context.Context, userID int64, deviceID string, reason string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&auth.RefreshToken{}).
		Where("user_id = ? AND device_id = ? AND revoked = ?", userID, deviceID, false).
		Updates(map[string]interface{}{
			"revoked":        true,
			"revoked_at":     now,
			"revoked_reason": reason,
		}).Error
}

// DeleteExpired deletes all expired refresh tokens
func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&auth.RefreshToken{}).Error
}

// DeleteRevoked deletes old revoked tokens (cleanup)
func (r *refreshTokenRepository) DeleteRevoked(ctx context.Context, olderThan time.Time) error {
	return r.db.WithContext(ctx).
		Where("revoked = ? AND revoked_at < ?", true, olderThan).
		Delete(&auth.RefreshToken{}).Error
}

// CountActiveByUserID counts active tokens for user
func (r *refreshTokenRepository) CountActiveByUserID(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&auth.RefreshToken{}).
		Where("user_id = ? AND revoked = ? AND expires_at > ?", userID, false, time.Now()).
		Count(&count).Error
	return count, err
}
