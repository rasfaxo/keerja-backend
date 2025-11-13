package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"keerja-backend/internal/config"
	"keerja-backend/internal/domain/notification"

	"firebase.google.com/go/v4/messaging"
)

// FCMPushService implements notification.PushNotificationService
type FCMPushService struct {
	deviceTokenRepo notification.DeviceTokenRepository
	cfg             *config.Config
}

// NewFCMPushService creates a new FCM push notification service
func NewFCMPushService(deviceTokenRepo notification.DeviceTokenRepository, cfg *config.Config) notification.PushNotificationService {
	return &FCMPushService{
		deviceTokenRepo: deviceTokenRepo,
		cfg:             cfg,
	}
}

// SendToDevice sends push notification to a specific device token
func (s *FCMPushService) SendToDevice(ctx context.Context, token string, message *notification.PushMessage) (*notification.PushResult, error) {
	if !config.IsFCMEnabled() {
		return &notification.PushResult{
			Success:      false,
			ErrorCode:    "FCM_DISABLED",
			ErrorMessage: "FCM is not enabled",
		}, nil
	}

	fcmClient, err := config.GetFCMClient()
	if err != nil {
		return &notification.PushResult{
			Success:      false,
			ErrorCode:    "FCM_CLIENT_ERROR",
			ErrorMessage: err.Error(),
		}, err
	}

	// Build FCM message
	fcmMessage := s.buildFCMMessage(token, message)

	// Send with timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, s.cfg.FCMTimeout)
	defer cancel()

	// Send message
	messageID, err := fcmClient.Send(ctxTimeout, fcmMessage)
	if err != nil {
		return s.handleFCMError(ctx, token, err)
	}

	// Update token last used timestamp
	if deviceToken, err := s.deviceTokenRepo.FindByToken(ctx, token); err == nil {
		deviceToken.MarkAsUsed()
		deviceToken.ResetFailures()
		_ = s.deviceTokenRepo.Update(ctx, deviceToken)
	}

	return &notification.PushResult{
		Success:   true,
		MessageID: messageID,
	}, nil
}

// SendToUser sends push notification to all user's devices
func (s *FCMPushService) SendToUser(ctx context.Context, userID int64, message *notification.PushMessage) ([]notification.PushResult, error) {
	// Get all active tokens for user
	tokens, err := s.deviceTokenRepo.FindByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user device tokens: %w", err)
	}

	if len(tokens) == 0 {
		log.Printf("No active device tokens found for user %d", userID)
		return []notification.PushResult{}, nil
	}

	// Extract token strings
	tokenStrings := make([]string, len(tokens))
	for i, token := range tokens {
		tokenStrings[i] = token.Token
	}

	// Send to all tokens
	return s.sendToMultipleTokens(ctx, tokenStrings, message)
}

// SendToMultipleUsers sends push notification to multiple users
func (s *FCMPushService) SendToMultipleUsers(ctx context.Context, userIDs []int64, message *notification.PushMessage) (map[int64][]notification.PushResult, error) {
	if len(userIDs) == 0 {
		return map[int64][]notification.PushResult{}, nil
	}

	// Get all tokens for users
	tokens, err := s.deviceTokenRepo.FindByUserIDs(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get device tokens for users: %w", err)
	}

	// Group tokens by user
	userTokens := make(map[int64][]string)
	for _, token := range tokens {
		userTokens[token.UserID] = append(userTokens[token.UserID], token.Token)
	}

	// Send to each user's tokens
	results := make(map[int64][]notification.PushResult)
	for userID, userTokenStrings := range userTokens {
		userResults, err := s.sendToMultipleTokens(ctx, userTokenStrings, message)
		if err != nil {
			log.Printf("Error sending to user %d: %v", userID, err)
		}
		results[userID] = userResults
	}

	return results, nil
}

// SendToTopic sends push notification to a topic
func (s *FCMPushService) SendToTopic(ctx context.Context, topic string, message *notification.PushMessage) (*notification.PushResult, error) {
	if !config.IsFCMEnabled() {
		return &notification.PushResult{
			Success:      false,
			ErrorCode:    "FCM_DISABLED",
			ErrorMessage: "FCM is not enabled",
		}, nil
	}

	fcmClient, err := config.GetFCMClient()
	if err != nil {
		return &notification.PushResult{
			Success:      false,
			ErrorCode:    "FCM_CLIENT_ERROR",
			ErrorMessage: err.Error(),
		}, err
	}

	// Build FCM message for topic
	fcmMessage := &messaging.Message{
		Topic: topic,
		Notification: &messaging.Notification{
			Title:    message.Title,
			Body:     message.Body,
			ImageURL: message.ImageURL,
		},
		Data: message.Data,
	}

	// Apply platform-specific configs
	s.applyPlatformConfigs(fcmMessage, message)

	// Send with timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, s.cfg.FCMTimeout)
	defer cancel()

	messageID, err := fcmClient.Send(ctxTimeout, fcmMessage)
	if err != nil {
		return &notification.PushResult{
			Success:      false,
			ErrorCode:    "SEND_ERROR",
			ErrorMessage: err.Error(),
		}, err
	}

	return &notification.PushResult{
		Success:   true,
		MessageID: messageID,
	}, nil
}

// RegisterDeviceToken registers a new device token for a user
func (s *FCMPushService) RegisterDeviceToken(ctx context.Context, userID int64, token string, platform notification.Platform, deviceInfo *notification.DeviceInfo) error {
	// Check if token already exists
	existingToken, err := s.deviceTokenRepo.FindByToken(ctx, token)
	if err == nil {
		// Token exists, update it
		existingToken.UserID = userID
		existingToken.Platform = platform
		if deviceInfo != nil {
			existingToken.DeviceInfo = *deviceInfo
		}
		existingToken.Activate()
		return s.deviceTokenRepo.Update(ctx, existingToken)
	}

	// Create new token
	deviceToken := &notification.DeviceToken{
		UserID:   userID,
		Token:    token,
		Platform: platform,
		IsActive: true,
	}

	if deviceInfo != nil {
		deviceToken.DeviceInfo = *deviceInfo
	}

	return s.deviceTokenRepo.Create(ctx, deviceToken)
}

// UnregisterDeviceToken removes a device token
func (s *FCMPushService) UnregisterDeviceToken(ctx context.Context, userID int64, token string) error {
	// Verify token belongs to user
	existingToken, err := s.deviceTokenRepo.FindByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("device token not found")
	}

	if existingToken.UserID != userID {
		return fmt.Errorf("unauthorized: token does not belong to user")
	}

	return s.deviceTokenRepo.DeleteByToken(ctx, token)
}

// RefreshDeviceToken updates an existing device token
func (s *FCMPushService) RefreshDeviceToken(ctx context.Context, oldToken, newToken string) error {
	// Find old token
	existingToken, err := s.deviceTokenRepo.FindByToken(ctx, oldToken)
	if err != nil {
		return fmt.Errorf("old token not found: %w", err)
	}

	// Update token string
	existingToken.Token = newToken
	existingToken.Activate()

	return s.deviceTokenRepo.Update(ctx, existingToken)
}

// GetUserDevices retrieves all registered devices for a user
func (s *FCMPushService) GetUserDevices(ctx context.Context, userID int64) ([]notification.DeviceToken, error) {
	return s.deviceTokenRepo.FindByUser(ctx, userID)
}

// ValidateToken validates if a token is still valid with FCM
func (s *FCMPushService) ValidateToken(ctx context.Context, token string) (bool, error) {
	// Send a dry-run message to validate token
	result, err := s.SendToDevice(ctx, token, &notification.PushMessage{
		Title: "Validation",
		Body:  "Token validation",
	})

	if err != nil {
		return false, err
	}

	return result.Success, nil
}

// CleanupInactiveTokens removes inactive device tokens
func (s *FCMPushService) CleanupInactiveTokens(ctx context.Context, inactiveDays int) error {
	// Find inactive tokens
	inactiveTokens, err := s.deviceTokenRepo.FindInactiveTokens(ctx, inactiveDays, 1000)
	if err != nil {
		return fmt.Errorf("failed to find inactive tokens: %w", err)
	}

	log.Printf("Found %d inactive tokens to cleanup", len(inactiveTokens))

	// Delete inactive tokens
	for _, token := range inactiveTokens {
		if err := s.deviceTokenRepo.Delete(ctx, token.ID); err != nil {
			log.Printf("Failed to delete inactive token %d: %v", token.ID, err)
		}
	}

	// Also cleanup tokens with high failure count
	highFailureTokens, err := s.deviceTokenRepo.FindTokensWithHighFailureCount(ctx, 5, 1000)
	if err != nil {
		log.Printf("Failed to find high failure tokens: %v", err)
		return nil
	}

	log.Printf("Found %d tokens with high failure count to cleanup", len(highFailureTokens))

	for _, token := range highFailureTokens {
		if err := s.deviceTokenRepo.Delete(ctx, token.ID); err != nil {
			log.Printf("Failed to delete high failure token %d: %v", token.ID, err)
		}
	}

	return nil
}

// ============================================================================
// Private Helper Methods
// ============================================================================

// sendToMultipleTokens sends to multiple tokens with batch support (max 500 per batch)
func (s *FCMPushService) sendToMultipleTokens(ctx context.Context, tokens []string, message *notification.PushMessage) ([]notification.PushResult, error) {
	if !config.IsFCMEnabled() {
		results := make([]notification.PushResult, len(tokens))
		for i := range results {
			results[i] = notification.PushResult{
				Success:      false,
				ErrorCode:    "FCM_DISABLED",
				ErrorMessage: "FCM is not enabled",
			}
		}
		return results, nil
	}

	fcmClient, err := config.GetFCMClient()
	if err != nil {
		return nil, err
	}

	results := make([]notification.PushResult, 0, len(tokens))

	// Split into batches (FCM limit: 500 tokens per batch)
	batchSize := s.cfg.FCMBatchSize
	if batchSize > 500 {
		batchSize = 500
	}

	for i := 0; i < len(tokens); i += batchSize {
		end := i + batchSize
		if end > len(tokens) {
			end = len(tokens)
		}

		batch := tokens[i:end]

		// Build multicast message
		multicastMessage := &messaging.MulticastMessage{
			Tokens: batch,
			Notification: &messaging.Notification{
				Title:    message.Title,
				Body:     message.Body,
				ImageURL: message.ImageURL,
			},
			Data: message.Data,
		}

		// Apply platform configs
		s.applyMulticastPlatformConfigs(multicastMessage, message)

		// Send batch with timeout
		ctxTimeout, cancel := context.WithTimeout(ctx, s.cfg.FCMTimeout)
		batchResponse, err := fcmClient.SendEachForMulticast(ctxTimeout, multicastMessage)
		cancel()

		if err != nil {
			log.Printf("Error sending batch: %v", err)
			// Add error results for this batch
			for range batch {
				results = append(results, notification.PushResult{
					Success:      false,
					ErrorCode:    "BATCH_SEND_ERROR",
					ErrorMessage: err.Error(),
				})
			}
			continue
		}

		// Process batch results
		for idx, response := range batchResponse.Responses {
			tokenStr := batch[idx]

			if response.Success {
				results = append(results, notification.PushResult{
					Success:   true,
					MessageID: response.MessageID,
				})

				// Update token success
				if deviceToken, err := s.deviceTokenRepo.FindByToken(ctx, tokenStr); err == nil {
					deviceToken.MarkAsUsed()
					deviceToken.ResetFailures()
					_ = s.deviceTokenRepo.Update(ctx, deviceToken)
				}
			} else {
				// Handle failure
				result, _ := s.handleFCMError(ctx, tokenStr, response.Error)
				results = append(results, *result)
			}
		}
	}

	return results, nil
}

// buildFCMMessage builds FCM message from PushMessage
func (s *FCMPushService) buildFCMMessage(token string, message *notification.PushMessage) *messaging.Message {
	fcmMessage := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title:    message.Title,
			Body:     message.Body,
			ImageURL: message.ImageURL,
		},
		Data: message.Data,
	}

	s.applyPlatformConfigs(fcmMessage, message)

	return fcmMessage
}

// applyPlatformConfigs applies platform-specific configurations
func (s *FCMPushService) applyPlatformConfigs(fcmMessage *messaging.Message, message *notification.PushMessage) {
	// Android config
	fcmMessage.Android = &messaging.AndroidConfig{
		Priority: "high",
		Notification: &messaging.AndroidNotification{
			Sound:       s.getSound(message),
			ClickAction: message.ClickAction,
		},
		TTL: s.getTTL(message),
	}

	// iOS config
	badge := 1
	if message.Badge != nil {
		badge = *message.Badge
	}
	fcmMessage.APNS = &messaging.APNSConfig{
		Payload: &messaging.APNSPayload{
			Aps: &messaging.Aps{
				Badge:    &badge,
				Sound:    s.getSound(message),
				Category: message.ClickAction,
			},
		},
	}

	// Web config
	fcmMessage.Webpush = &messaging.WebpushConfig{
		Notification: &messaging.WebpushNotification{
			Title: message.Title,
			Body:  message.Body,
			Icon:  message.ImageURL,
		},
	}
}

// applyMulticastPlatformConfigs applies configs to multicast messages
func (s *FCMPushService) applyMulticastPlatformConfigs(multicastMessage *messaging.MulticastMessage, message *notification.PushMessage) {
	// Android config
	multicastMessage.Android = &messaging.AndroidConfig{
		Priority: "high",
		Notification: &messaging.AndroidNotification{
			Sound:       s.getSound(message),
			ClickAction: message.ClickAction,
		},
		TTL: s.getTTL(message),
	}

	// iOS config
	badge := 1
	if message.Badge != nil {
		badge = *message.Badge
	}
	multicastMessage.APNS = &messaging.APNSConfig{
		Payload: &messaging.APNSPayload{
			Aps: &messaging.Aps{
				Badge:    &badge,
				Sound:    s.getSound(message),
				Category: message.ClickAction,
			},
		},
	}

	// Web config
	multicastMessage.Webpush = &messaging.WebpushConfig{
		Notification: &messaging.WebpushNotification{
			Title: message.Title,
			Body:  message.Body,
			Icon:  message.ImageURL,
		},
	}
}

// getSound returns notification sound
func (s *FCMPushService) getSound(message *notification.PushMessage) string {
	if message.Sound != "" {
		return message.Sound
	}
	return s.cfg.PushDefaultSound
}

// getTTL returns time-to-live duration
func (s *FCMPushService) getTTL(message *notification.PushMessage) *time.Duration {
	var ttl time.Duration
	if message.TTL > 0 {
		ttl = time.Duration(message.TTL) * time.Second
	} else {
		ttl = time.Duration(s.cfg.PushDefaultTTL) * time.Second
	}
	return &ttl
}

// handleFCMError handles FCM-specific errors and updates token status
func (s *FCMPushService) handleFCMError(ctx context.Context, token string, err error) (*notification.PushResult, error) {
	errorCode := "UNKNOWN_ERROR"
	errorMessage := err.Error()

	// Parse FCM error codes
	// Reference: https://firebase.google.com/docs/cloud-messaging/admin/errors
	if messaging.IsInvalidArgument(err) {
		errorCode = "INVALID_ARGUMENT"
	} else if messaging.IsUnregistered(err) {
		errorCode = "UNREGISTERED"
		// Token is invalid, deactivate it
		_ = s.deviceTokenRepo.Deactivate(ctx, token)
		log.Printf("Deactivated invalid token: %s", token)
	} else if messaging.IsInternal(err) || messaging.IsUnavailable(err) {
		errorCode = "SERVER_ERROR"
	} else if messaging.IsQuotaExceeded(err) {
		errorCode = "QUOTA_EXCEEDED"
	}

	// Record failure
	if deviceToken, err := s.deviceTokenRepo.FindByToken(ctx, token); err == nil {
		deviceToken.RecordFailure(errorCode + ": " + errorMessage)
		_ = s.deviceTokenRepo.Update(ctx, deviceToken)
	}

	return &notification.PushResult{
		Success:      false,
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}, nil
}
