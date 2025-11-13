package postgres

import (
	"context"
	"fmt"
	"time"

	"keerja-backend/internal/domain/notification"

	"gorm.io/gorm"
)

// DeviceTokenRepository implements notification.DeviceTokenRepository
type DeviceTokenRepository struct {
	db *gorm.DB
}

// NewDeviceTokenRepository creates a new device token repository
func NewDeviceTokenRepository(db *gorm.DB) notification.DeviceTokenRepository {
	return &DeviceTokenRepository{db: db}
}

// Create creates a new device token
func (r *DeviceTokenRepository) Create(ctx context.Context, token *notification.DeviceToken) error {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("failed to create device token: %w", err)
	}
	return nil
}

// FindByID finds device token by ID
func (r *DeviceTokenRepository) FindByID(ctx context.Context, id int64) (*notification.DeviceToken, error) {
	var token notification.DeviceToken
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("device token not found")
		}
		return nil, fmt.Errorf("failed to find device token: %w", err)
	}
	return &token, nil
}

// FindByToken finds device token by FCM token string
func (r *DeviceTokenRepository) FindByToken(ctx context.Context, token string) (*notification.DeviceToken, error) {
	var deviceToken notification.DeviceToken
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&deviceToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("device token not found")
		}
		return nil, fmt.Errorf("failed to find device token: %w", err)
	}
	return &deviceToken, nil
}

// FindByUser finds all active device tokens for a user
func (r *DeviceTokenRepository) FindByUser(ctx context.Context, userID int64) ([]notification.DeviceToken, error) {
	var tokens []notification.DeviceToken
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("created_at DESC").
		Find(&tokens).Error; err != nil {
		return nil, fmt.Errorf("failed to find user device tokens: %w", err)
	}
	return tokens, nil
}

// FindByUserAndPlatform finds device tokens for a user on specific platform
func (r *DeviceTokenRepository) FindByUserAndPlatform(ctx context.Context, userID int64, platform notification.Platform) ([]notification.DeviceToken, error) {
	var tokens []notification.DeviceToken
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND platform = ? AND is_active = ?", userID, platform, true).
		Order("created_at DESC").
		Find(&tokens).Error; err != nil {
		return nil, fmt.Errorf("failed to find user device tokens by platform: %w", err)
	}
	return tokens, nil
}

// Update updates device token
func (r *DeviceTokenRepository) Update(ctx context.Context, token *notification.DeviceToken) error {
	if err := r.db.WithContext(ctx).Save(token).Error; err != nil {
		return fmt.Errorf("failed to update device token: %w", err)
	}
	return nil
}

// Delete deletes device token
func (r *DeviceTokenRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&notification.DeviceToken{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete device token: %w", err)
	}
	return nil
}

// DeleteByToken deletes device token by token string
func (r *DeviceTokenRepository) DeleteByToken(ctx context.Context, token string) error {
	if err := r.db.WithContext(ctx).Where("token = ?", token).Delete(&notification.DeviceToken{}).Error; err != nil {
		return fmt.Errorf("failed to delete device token by token: %w", err)
	}
	return nil
}

// Deactivate deactivates a device token
func (r *DeviceTokenRepository) Deactivate(ctx context.Context, token string) error {
	if err := r.db.WithContext(ctx).
		Model(&notification.DeviceToken{}).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"is_active":  false,
			"updated_at": time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("failed to deactivate device token: %w", err)
	}
	return nil
}

// FindInactiveTokens finds inactive tokens that haven't been used for specified days
func (r *DeviceTokenRepository) FindInactiveTokens(ctx context.Context, inactiveDays int, limit int) ([]notification.DeviceToken, error) {
	var tokens []notification.DeviceToken
	cutoffDate := time.Now().AddDate(0, 0, -inactiveDays)

	query := r.db.WithContext(ctx).
		Where("(is_active = ? OR last_used_at < ?)", false, cutoffDate).
		Order("last_used_at ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&tokens).Error; err != nil {
		return nil, fmt.Errorf("failed to find inactive device tokens: %w", err)
	}
	return tokens, nil
}

// FindInactive finds tokens inactive since the cutoff date (for cleanup job)
func (r *DeviceTokenRepository) FindInactive(ctx context.Context, cutoffDate time.Time) ([]notification.DeviceToken, error) {
	var tokens []notification.DeviceToken

	if err := r.db.WithContext(ctx).
		Where("is_active = ? AND (last_used_at IS NULL OR last_used_at < ?)", true, cutoffDate).
		Order("last_used_at ASC").
		Find(&tokens).Error; err != nil {
		return nil, fmt.Errorf("failed to find inactive device tokens: %w", err)
	}
	return tokens, nil
}

// FindByFailureCount finds tokens with failure count >= minFailures (for cleanup job)
func (r *DeviceTokenRepository) FindByFailureCount(ctx context.Context, minFailures int) ([]notification.DeviceToken, error) {
	var tokens []notification.DeviceToken

	if err := r.db.WithContext(ctx).
		Where("failure_count >= ?", minFailures).
		Order("failure_count DESC, last_failure_at DESC").
		Find(&tokens).Error; err != nil {
		return nil, fmt.Errorf("failed to find tokens with high failure count: %w", err)
	}
	return tokens, nil
}

// CountByUser counts active device tokens for a user
func (r *DeviceTokenRepository) CountByUser(ctx context.Context, userID int64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&notification.DeviceToken{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count user device tokens: %w", err)
	}
	return count, nil
}

// BatchUpdate updates multiple device tokens
func (r *DeviceTokenRepository) BatchUpdate(ctx context.Context, tokens []notification.DeviceToken) error {
	if len(tokens) == 0 {
		return nil
	}

	// Use transaction for batch update
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, token := range tokens {
			if err := tx.Save(&token).Error; err != nil {
				return fmt.Errorf("failed to batch update device token: %w", err)
			}
		}
		return nil
	})
}

// FindByUserIDs finds all active device tokens for multiple users (for batch sending)
func (r *DeviceTokenRepository) FindByUserIDs(ctx context.Context, userIDs []int64) ([]notification.DeviceToken, error) {
	if len(userIDs) == 0 {
		return []notification.DeviceToken{}, nil
	}

	var tokens []notification.DeviceToken
	if err := r.db.WithContext(ctx).
		Where("user_id IN ? AND is_active = ?", userIDs, true).
		Order("user_id ASC, created_at DESC").
		Find(&tokens).Error; err != nil {
		return nil, fmt.Errorf("failed to find device tokens for multiple users: %w", err)
	}
	return tokens, nil
}

// CountActiveTokens counts total active device tokens in the system
func (r *DeviceTokenRepository) CountActiveTokens(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&notification.DeviceToken{}).
		Where("is_active = ?", true).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count active device tokens: %w", err)
	}
	return count, nil
}

// FindTokensWithHighFailureCount finds tokens with high failure counts for cleanup
func (r *DeviceTokenRepository) FindTokensWithHighFailureCount(ctx context.Context, minFailures int, limit int) ([]notification.DeviceToken, error) {
	var tokens []notification.DeviceToken

	query := r.db.WithContext(ctx).
		Where("failure_count >= ?", minFailures).
		Order("failure_count DESC, last_failure_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&tokens).Error; err != nil {
		return nil, fmt.Errorf("failed to find tokens with high failure count: %w", err)
	}
	return tokens, nil
}
