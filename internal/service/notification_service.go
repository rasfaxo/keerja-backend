package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"keerja-backend/internal/domain/notification"
)

// notificationService implements notification.NotificationService interface
type notificationService struct {
	notifRepo notification.NotificationRepository
	// In production, inject push notification service (Firebase, OneSignal, etc.)
	// pushService PushNotificationService
}

// NewNotificationService creates a new notification service instance
func NewNotificationService(notifRepo notification.NotificationRepository) notification.NotificationService {
	return &notificationService{
		notifRepo: notifRepo,
	}
}

// SendNotification sends a notification to a user
func (s *notificationService) SendNotification(ctx context.Context, req *notification.SendNotificationRequest) (*notification.Notification, error) {
	// Get user preferences
	prefs, _ := s.GetNotificationPreferences(ctx, req.UserID)
	if prefs != nil && !prefs.CanSendNotification(req.Type) {
		return nil, errors.New("user has disabled this type of notification")
	}

	// Convert metadata to JSON
	var dataJSON string
	if req.Data != nil {
		data, err := json.Marshal(req.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal notification data: %w", err)
		}
		dataJSON = string(data)
	}

	// Create notification
	notif := &notification.Notification{
		UserID:      req.UserID,
		Type:        req.Type,
		Title:       req.Title,
		Message:     req.Message,
		Data:        dataJSON,
		Priority:    req.Priority,
		Category:    req.Category,
		ActionURL:   req.ActionURL,
		Icon:        req.Icon,
		SenderID:    req.SenderID,
		RelatedID:   req.RelatedID,
		RelatedType: req.RelatedType,
		ExpiresAt:   req.ExpiresAt,
		Channel:     req.Channel,
	}

	// Set defaults
	if notif.Priority == "" {
		notif.Priority = "normal"
	}
	if notif.Channel == "" {
		notif.Channel = "in_app"
	}

	// Save to database
	if err := s.notifRepo.Create(ctx, notif); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Send to appropriate channels
	go s.sendToChannels(ctx, notif, prefs)

	return notif, nil
}

// SendBulkNotification sends notifications to multiple users
func (s *notificationService) SendBulkNotification(ctx context.Context, userIDs []int64, req *notification.SendNotificationRequest) error {
	notifications := make([]notification.Notification, 0, len(userIDs))

	// Convert metadata to JSON
	var dataJSON string
	if req.Data != nil {
		data, err := json.Marshal(req.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal notification data: %w", err)
		}
		dataJSON = string(data)
	}

	for _, userID := range userIDs {
		// Check user preferences
		prefs, _ := s.GetNotificationPreferences(ctx, userID)
		if prefs != nil && !prefs.CanSendNotification(req.Type) {
			continue // Skip if user has disabled this type
		}

		notif := notification.Notification{
			UserID:      userID,
			Type:        req.Type,
			Title:       req.Title,
			Message:     req.Message,
			Data:        dataJSON,
			Priority:    req.Priority,
			Category:    req.Category,
			ActionURL:   req.ActionURL,
			Icon:        req.Icon,
			SenderID:    req.SenderID,
			RelatedID:   req.RelatedID,
			RelatedType: req.RelatedType,
			ExpiresAt:   req.ExpiresAt,
			Channel:     req.Channel,
		}

		// Set defaults
		if notif.Priority == "" {
			notif.Priority = "normal"
		}
		if notif.Channel == "" {
			notif.Channel = "in_app"
		}

		notifications = append(notifications, notif)
	}

	// Bulk create notifications
	if err := s.notifRepo.BulkCreate(ctx, notifications); err != nil {
		return fmt.Errorf("failed to create bulk notifications: %w", err)
	}

	return nil
}

// GetUserNotifications retrieves user notifications
func (s *notificationService) GetUserNotifications(ctx context.Context, userID int64, filter notification.NotificationFilter, page, limit int) ([]notification.Notification, int64, error) {
	return s.notifRepo.ListByUser(ctx, userID, filter, page, limit)
}

// GetUnreadNotifications retrieves unread notifications
func (s *notificationService) GetUnreadNotifications(ctx context.Context, userID int64, limit int) ([]notification.Notification, error) {
	return s.notifRepo.GetUnreadByUser(ctx, userID, limit)
}

// GetNotificationByID retrieves notification by ID
func (s *notificationService) GetNotificationByID(ctx context.Context, id, userID int64) (*notification.Notification, error) {
	notif, err := s.notifRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("notification not found: %w", err)
	}

	// Check ownership
	if notif.UserID != userID {
		return nil, errors.New("you do not have access to this notification")
	}

	return notif, nil
}

// MarkAsRead marks notification as read
func (s *notificationService) MarkAsRead(ctx context.Context, id, userID int64) error {
	// Check ownership
	if _, err := s.GetNotificationByID(ctx, id, userID); err != nil {
		return err
	}

	return s.notifRepo.MarkAsRead(ctx, id)
}

// MarkAsUnread marks notification as unread
func (s *notificationService) MarkAsUnread(ctx context.Context, id, userID int64) error {
	// Check ownership
	notif, err := s.GetNotificationByID(ctx, id, userID)
	if err != nil {
		return err
	}

	// Update notification
	notif.MarkAsUnread()
	return s.notifRepo.Update(ctx, notif)
}

// MarkAllAsRead marks all notifications as read
func (s *notificationService) MarkAllAsRead(ctx context.Context, userID int64) error {
	return s.notifRepo.MarkAllAsRead(ctx, userID)
}

// DeleteNotification deletes a notification
func (s *notificationService) DeleteNotification(ctx context.Context, id, userID int64) error {
	// Check ownership
	if _, err := s.GetNotificationByID(ctx, id, userID); err != nil {
		return err
	}

	return s.notifRepo.Delete(ctx, id)
}

// DeleteAllNotifications deletes all notifications for user
func (s *notificationService) DeleteAllNotifications(ctx context.Context, userID int64) error {
	return s.notifRepo.DeleteByUser(ctx, userID)
}

// GetUnreadCount retrieves unread notification count
func (s *notificationService) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	return s.notifRepo.CountUnreadByUser(ctx, userID)
}

// GetNotificationStats retrieves notification statistics
func (s *notificationService) GetNotificationStats(ctx context.Context, userID int64) (*notification.NotificationStats, error) {
	return s.notifRepo.GetStats(ctx, userID)
}

// ===== Specific Notification Types =====

// NotifyJobApplication sends job application notification
func (s *notificationService) NotifyJobApplication(ctx context.Context, userID, jobID, applicationID int64) error {
	req := &notification.SendNotificationRequest{
		UserID:      userID,
		Type:        "job_application",
		Title:       "Application Submitted",
		Message:     "Your job application has been submitted successfully",
		Category:    "application",
		Priority:    "normal",
		Icon:        "check-circle",
		RelatedID:   &applicationID,
		RelatedType: "application",
		ActionURL:   fmt.Sprintf("/applications/%d", applicationID),
		Data: map[string]interface{}{
			"job_id":         jobID,
			"application_id": applicationID,
		},
	}

	_, err := s.SendNotification(ctx, req)
	return err
}

// NotifyInterviewScheduled sends interview scheduled notification
func (s *notificationService) NotifyInterviewScheduled(ctx context.Context, userID, interviewID int64, interviewDate time.Time) error {
	req := &notification.SendNotificationRequest{
		UserID:      userID,
		Type:        "interview",
		Title:       "Interview Scheduled",
		Message:     fmt.Sprintf("Your interview is scheduled for %s", interviewDate.Format("Jan 02, 2006 at 3:04 PM")),
		Category:    "interview",
		Priority:    "high",
		Icon:        "calendar",
		RelatedID:   &interviewID,
		RelatedType: "interview",
		ActionURL:   fmt.Sprintf("/interviews/%d", interviewID),
		Data: map[string]interface{}{
			"interview_id":   interviewID,
			"interview_date": interviewDate,
		},
	}

	_, err := s.SendNotification(ctx, req)
	return err
}

// NotifyStatusUpdate sends status update notification
func (s *notificationService) NotifyStatusUpdate(ctx context.Context, userID, applicationID int64, oldStatus, newStatus string) error {
	statusMessages := map[string]string{
		"screening":   "Your application is being reviewed",
		"shortlisted": "Congratulations! You've been shortlisted",
		"interview":   "You've been invited for an interview",
		"offered":     "Congratulations! You've received a job offer",
		"hired":       "Congratulations! You've been hired",
		"rejected":    "Your application status has been updated",
	}

	message := statusMessages[newStatus]
	if message == "" {
		message = "Your application status has been updated"
	}

	priority := "normal"
	if newStatus == "offered" || newStatus == "hired" || newStatus == "shortlisted" {
		priority = "high"
	}

	req := &notification.SendNotificationRequest{
		UserID:      userID,
		Type:        "status_update",
		Title:       "Application Status Update",
		Message:     message,
		Category:    "application",
		Priority:    priority,
		Icon:        "bell",
		RelatedID:   &applicationID,
		RelatedType: "application",
		ActionURL:   fmt.Sprintf("/applications/%d", applicationID),
		Data: map[string]interface{}{
			"application_id": applicationID,
			"old_status":     oldStatus,
			"new_status":     newStatus,
		},
	}

	_, err := s.SendNotification(ctx, req)
	return err
}

// NotifyJobRecommendation sends job recommendation notification
func (s *notificationService) NotifyJobRecommendation(ctx context.Context, userID, jobID int64) error {
	req := &notification.SendNotificationRequest{
		UserID:      userID,
		Type:        "job_recommendation",
		Title:       "New Job Match",
		Message:     "We found a job that matches your profile",
		Category:    "job",
		Priority:    "normal",
		Icon:        "briefcase",
		RelatedID:   &jobID,
		RelatedType: "job",
		ActionURL:   fmt.Sprintf("/jobs/%d", jobID),
		Data: map[string]interface{}{
			"job_id": jobID,
		},
	}

	_, err := s.SendNotification(ctx, req)
	return err
}

// NotifyCompanyUpdate sends company update notification
func (s *notificationService) NotifyCompanyUpdate(ctx context.Context, userIDs []int64, companyID int64, updateType string) error {
	req := &notification.SendNotificationRequest{
		Type:        "company_update",
		Title:       "Company Update",
		Message:     "Your followed company has posted an update",
		Category:    "company",
		Priority:    "low",
		Icon:        "building",
		RelatedID:   &companyID,
		RelatedType: "company",
		ActionURL:   fmt.Sprintf("/companies/%d", companyID),
		Data: map[string]interface{}{
			"company_id":  companyID,
			"update_type": updateType,
		},
	}

	return s.SendBulkNotification(ctx, userIDs, req)
}

// ===== Notification Preferences =====

// GetNotificationPreferences retrieves user notification preferences
func (s *notificationService) GetNotificationPreferences(ctx context.Context, userID int64) (*notification.NotificationPreference, error) {
	prefs, err := s.notifRepo.FindPreferenceByUser(ctx, userID)
	if err != nil {
		// If preferences don't exist, create default ones
		prefs = &notification.NotificationPreference{
			UserID:                    userID,
			EmailEnabled:              true,
			PushEnabled:               true,
			SMSEnabled:                false,
			JobApplicationsEnabled:    true,
			InterviewEnabled:          true,
			StatusUpdatesEnabled:      true,
			JobRecommendationsEnabled: true,
			CompanyUpdatesEnabled:     true,
			MarketingEnabled:          false,
			WeeklyDigestEnabled:       true,
		}
		if err := s.notifRepo.CreatePreference(ctx, prefs); err != nil {
			return nil, fmt.Errorf("failed to create default preferences: %w", err)
		}
	}

	return prefs, nil
}

// UpdateNotificationPreferences updates user notification preferences
func (s *notificationService) UpdateNotificationPreferences(ctx context.Context, userID int64, prefs *notification.NotificationPreference) error {
	// Ensure user ID matches
	prefs.UserID = userID

	// Check if preferences exist
	existingPrefs, _ := s.notifRepo.FindPreferenceByUser(ctx, userID)
	if existingPrefs == nil {
		// Create new preferences
		return s.notifRepo.CreatePreference(ctx, prefs)
	}

	// Update existing preferences
	prefs.ID = existingPrefs.ID
	return s.notifRepo.UpdatePreference(ctx, prefs)
}

// ===== Push and Email Notifications =====

// SendPushNotification sends push notification
func (s *notificationService) SendPushNotification(ctx context.Context, userID int64, notif *notification.Notification) error {
	// Check if push enabled
	prefs, _ := s.GetNotificationPreferences(ctx, userID)
	if prefs != nil && !prefs.IsPushEnabled() {
		return nil // Skip if push disabled
	}

	// TODO: Implement push notification
	// Integrate with Firebase Cloud Messaging, OneSignal, etc.
	// Send push notification to user's devices

	fmt.Printf("\n====== PUSH NOTIFICATION ======\n")
	fmt.Printf("User ID: %d\n", userID)
	fmt.Printf("Title: %s\n", notif.Title)
	fmt.Printf("Message: %s\n", notif.Message)
	fmt.Printf("===============================\n\n")

	return nil
}

// SendEmailNotification sends email notification
func (s *notificationService) SendEmailNotification(ctx context.Context, userID int64, notif *notification.Notification) error {
	// Check if email enabled
	prefs, _ := s.GetNotificationPreferences(ctx, userID)
	if prefs != nil && !prefs.IsEmailEnabled() {
		return nil // Skip if email disabled
	}

	// TODO: Integrate with email service
	// Send email notification to user

	fmt.Printf("\n====== EMAIL NOTIFICATION ======\n")
	fmt.Printf("User ID: %d\n", userID)
	fmt.Printf("Subject: %s\n", notif.Title)
	fmt.Printf("Body: %s\n", notif.Message)
	fmt.Printf("================================\n\n")

	return nil
}

// CleanupExpiredNotifications removes expired notifications
func (s *notificationService) CleanupExpiredNotifications(ctx context.Context) error {
	// Get expired notifications
	expiredNotifs, err := s.notifRepo.GetExpiredNotifications(ctx, 100)
	if err != nil {
		return fmt.Errorf("failed to get expired notifications: %w", err)
	}

	// Delete expired notifications
	for _, notif := range expiredNotifs {
		if err := s.notifRepo.Delete(ctx, notif.ID); err != nil {
			// Log error but continue
			fmt.Printf("failed to delete expired notification %d: %v\n", notif.ID, err)
		}
	}

	return nil
}

// ===== Helper Methods =====

// sendToChannels sends notification to appropriate channels based on preferences
func (s *notificationService) sendToChannels(ctx context.Context, notif *notification.Notification, prefs *notification.NotificationPreference) {
	// Send in-app notification (already saved to database)
	notif.MarkAsSent()
	s.notifRepo.Update(ctx, notif)

	if prefs == nil {
		return
	}

	// Send push notification if enabled
	if prefs.PushEnabled && (notif.Channel == "push" || notif.Channel == "in_app") {
		s.SendPushNotification(ctx, notif.UserID, notif)
	}

	// Send email notification if enabled and high priority
	if prefs.EmailEnabled && notif.IsHighPriority() {
		s.SendEmailNotification(ctx, notif.UserID, notif)
	}
}
