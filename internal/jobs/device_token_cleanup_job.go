package jobs

import (
	"context"
	"fmt"
	"time"

	"keerja-backend/internal/domain/notification"

	"github.com/sirupsen/logrus"
)

// DeviceTokenCleanupJob cleans up inactive and invalid device tokens
type DeviceTokenCleanupJob struct {
	deviceTokenRepo notification.DeviceTokenRepository
	logger          *logrus.Logger
	config          CleanupConfig
}

// CleanupConfig holds configuration for the cleanup job
type CleanupConfig struct {
	InactiveDays        int // Days of inactivity before token is considered inactive
	MaxFailureCount     int // Maximum failure count before token is removed
	BatchSize           int // Number of tokens to process per batch
	EnableInactiveClean bool
	EnableFailureClean  bool
}

// NewDeviceTokenCleanupJob creates a new device token cleanup job
func NewDeviceTokenCleanupJob(
	deviceTokenRepo notification.DeviceTokenRepository,
	logger *logrus.Logger,
	config CleanupConfig,
) *DeviceTokenCleanupJob {
	// Set defaults if not provided
	if config.InactiveDays == 0 {
		config.InactiveDays = 90 // Default: 90 days
	}
	if config.MaxFailureCount == 0 {
		config.MaxFailureCount = 10 // Default: 10 failures
	}
	if config.BatchSize == 0 {
		config.BatchSize = 100 // Default: 100 tokens per batch
	}
	// Enable both by default
	if !config.EnableInactiveClean && !config.EnableFailureClean {
		config.EnableInactiveClean = true
		config.EnableFailureClean = true
	}

	return &DeviceTokenCleanupJob{
		deviceTokenRepo: deviceTokenRepo,
		logger:          logger,
		config:          config,
	}
}

// Name returns the job name
func (j *DeviceTokenCleanupJob) Name() string {
	return "device_token_cleanup"
}

// Schedule returns the cron schedule (daily at midnight)
// Format: second minute hour day month weekday
func (j *DeviceTokenCleanupJob) Schedule() string {
	return "0 0 0 * * *" // Every day at midnight (00:00:00)
}

// Run executes the cleanup job
func (j *DeviceTokenCleanupJob) Run(ctx context.Context) error {
	startTime := time.Now()

	j.logger.Info("🧹 Starting device token cleanup job...")

	var totalCleaned int64
	var errors []error

	// Cleanup 1: Remove inactive tokens (not used for X days)
	if j.config.EnableInactiveClean {
		inactiveCleaned, err := j.cleanupInactiveTokens(ctx)
		if err != nil {
			j.logger.WithError(err).Error("Failed to cleanup inactive tokens")
			errors = append(errors, fmt.Errorf("inactive cleanup: %w", err))
		} else {
			totalCleaned += inactiveCleaned
		}
	}

	// Cleanup 2: Remove tokens with high failure counts
	if j.config.EnableFailureClean {
		failureCleaned, err := j.cleanupFailedTokens(ctx)
		if err != nil {
			j.logger.WithError(err).Error("Failed to cleanup failed tokens")
			errors = append(errors, fmt.Errorf("failure cleanup: %w", err))
		} else {
			totalCleaned += failureCleaned
		}
	}

	duration := time.Since(startTime)

	// Log results
	j.logger.WithFields(logrus.Fields{
		"total_cleaned": totalCleaned,
		"duration_ms":   duration.Milliseconds(),
		"errors":        len(errors),
	}).Info("✅ Device token cleanup job completed")

	// Return error if any cleanup failed
	if len(errors) > 0 {
		return fmt.Errorf("cleanup job completed with %d errors: %v", len(errors), errors)
	}

	return nil
}

// cleanupInactiveTokens removes tokens that haven't been used for X days
func (j *DeviceTokenCleanupJob) cleanupInactiveTokens(ctx context.Context) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -j.config.InactiveDays)

	j.logger.WithFields(logrus.Fields{
		"cutoff_date":   cutoffDate.Format("2006-01-02"),
		"inactive_days": j.config.InactiveDays,
	}).Info("🔍 Searching for inactive device tokens...")

	// Find inactive tokens
	tokens, err := j.deviceTokenRepo.FindInactive(ctx, cutoffDate)
	if err != nil {
		return 0, fmt.Errorf("failed to find inactive tokens: %w", err)
	}

	if len(tokens) == 0 {
		j.logger.Info("No inactive tokens found to clean up")
		return 0, nil
	}

	j.logger.WithField("count", len(tokens)).Info("Found inactive tokens to delete")

	// Delete tokens in batches
	var deleted int64
	for i := 0; i < len(tokens); i += j.config.BatchSize {
		end := i + j.config.BatchSize
		if end > len(tokens) {
			end = len(tokens)
		}

		batch := tokens[i:end]
		for _, token := range batch {
			if err := j.deviceTokenRepo.Delete(ctx, token.ID); err != nil {
				j.logger.WithError(err).WithField("token_id", token.ID).
					Error("Failed to delete inactive token")
				continue
			}
			deleted++
		}
	}

	j.logger.WithFields(logrus.Fields{
		"deleted":     deleted,
		"total_found": len(tokens),
	}).Info("✓ Inactive tokens cleanup completed")

	return deleted, nil
}

// cleanupFailedTokens removes tokens with high failure counts
func (j *DeviceTokenCleanupJob) cleanupFailedTokens(ctx context.Context) (int64, error) {
	j.logger.WithField("max_failures", j.config.MaxFailureCount).
		Info("🔍 Searching for failed device tokens...")

	// Find tokens with high failure counts
	tokens, err := j.deviceTokenRepo.FindByFailureCount(ctx, j.config.MaxFailureCount)
	if err != nil {
		return 0, fmt.Errorf("failed to find failed tokens: %w", err)
	}

	if len(tokens) == 0 {
		j.logger.Info("No failed tokens found to clean up")
		return 0, nil
	}

	j.logger.WithField("count", len(tokens)).Info("Found failed tokens to delete")

	// Delete tokens in batches
	var deleted int64
	for i := 0; i < len(tokens); i += j.config.BatchSize {
		end := i + j.config.BatchSize
		if end > len(tokens) {
			end = len(tokens)
		}

		batch := tokens[i:end]
		for _, token := range batch {
			if err := j.deviceTokenRepo.Delete(ctx, token.ID); err != nil {
				j.logger.WithError(err).WithFields(logrus.Fields{
					"token_id":      token.ID,
					"failure_count": token.FailureCount,
				}).Error("Failed to delete failed token")
				continue
			}
			deleted++
		}
	}

	j.logger.WithFields(logrus.Fields{
		"deleted":     deleted,
		"total_found": len(tokens),
	}).Info("✓ Failed tokens cleanup completed")

	return deleted, nil
}

// GetStats returns current statistics (for monitoring/debugging)
func (j *DeviceTokenCleanupJob) GetStats(ctx context.Context) (map[string]interface{}, error) {
	// This could be called from an admin endpoint to check job health
	cutoffDate := time.Now().AddDate(0, 0, -j.config.InactiveDays)

	inactiveTokens, err := j.deviceTokenRepo.FindInactive(ctx, cutoffDate)
	if err != nil {
		return nil, err
	}

	failedTokens, err := j.deviceTokenRepo.FindByFailureCount(ctx, j.config.MaxFailureCount)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"inactive_tokens_pending": len(inactiveTokens),
		"failed_tokens_pending":   len(failedTokens),
		"config": map[string]interface{}{
			"inactive_days":         j.config.InactiveDays,
			"max_failure_count":     j.config.MaxFailureCount,
			"batch_size":            j.config.BatchSize,
			"enable_inactive_clean": j.config.EnableInactiveClean,
			"enable_failure_clean":  j.config.EnableFailureClean,
		},
		"next_run": "daily at 00:00:00",
	}, nil
}
