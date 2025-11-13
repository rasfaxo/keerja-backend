package notification

import (
	"context"
	"time"
)

// NotificationService defines the interface for notification operations
type NotificationService interface {
	// SendNotification sends a notification to a user
	SendNotification(ctx context.Context, req *SendNotificationRequest) (*Notification, error)

	// SendBulkNotification sends notifications to multiple users
	SendBulkNotification(ctx context.Context, userIDs []int64, req *SendNotificationRequest) error

	// GetUserNotifications retrieves user notifications
	GetUserNotifications(ctx context.Context, userID int64, filter NotificationFilter, page, limit int) ([]Notification, int64, error)

	// GetUnreadNotifications retrieves unread notifications
	GetUnreadNotifications(ctx context.Context, userID int64, limit int) ([]Notification, error)

	// GetNotificationByID retrieves notification by ID
	GetNotificationByID(ctx context.Context, id, userID int64) (*Notification, error)

	// MarkAsRead marks notification as read
	MarkAsRead(ctx context.Context, id, userID int64) error

	// MarkAsUnread marks notification as unread
	MarkAsUnread(ctx context.Context, id, userID int64) error

	// MarkAllAsRead marks all notifications as read
	MarkAllAsRead(ctx context.Context, userID int64) error

	// DeleteNotification deletes a notification
	DeleteNotification(ctx context.Context, id, userID int64) error

	// DeleteAllNotifications deletes all notifications for user
	DeleteAllNotifications(ctx context.Context, userID int64) error

	// GetUnreadCount retrieves unread notification count
	GetUnreadCount(ctx context.Context, userID int64) (int64, error)

	// GetNotificationStats retrieves notification statistics
	GetNotificationStats(ctx context.Context, userID int64) (*NotificationStats, error)

	// NotifyJobApplication sends job application notification
	NotifyJobApplication(ctx context.Context, userID, jobID, applicationID int64) error

	// NotifyInterviewScheduled sends interview scheduled notification
	NotifyInterviewScheduled(ctx context.Context, userID, interviewID int64, interviewDate time.Time) error

	// NotifyStatusUpdate sends status update notification
	NotifyStatusUpdate(ctx context.Context, userID, applicationID int64, oldStatus, newStatus string) error

	// NotifyJobRecommendation sends job recommendation notification
	NotifyJobRecommendation(ctx context.Context, userID, jobID int64) error

	// NotifyCompanyUpdate sends company update notification
	NotifyCompanyUpdate(ctx context.Context, userIDs []int64, companyID int64, updateType string) error

	// GetNotificationPreferences retrieves user notification preferences
	GetNotificationPreferences(ctx context.Context, userID int64) (*NotificationPreference, error)

	// UpdateNotificationPreferences updates user notification preferences
	UpdateNotificationPreferences(ctx context.Context, userID int64, prefs *NotificationPreference) error

	// SendPushNotification sends push notification
	SendPushNotification(ctx context.Context, userID int64, notification *Notification) error

	// SendEmailNotification sends email notification
	SendEmailNotification(ctx context.Context, userID int64, notification *Notification) error

	// CleanupExpiredNotifications removes expired notifications
	CleanupExpiredNotifications(ctx context.Context) error
}

// SendNotificationRequest represents notification sending request
type SendNotificationRequest struct {
	UserID      int64
	Type        string
	Title       string
	Message     string
	Data        map[string]interface{}
	Priority    string
	Category    string
	ActionURL   string
	Icon        string
	SenderID    *int64
	RelatedID   *int64
	RelatedType string
	ExpiresAt   *time.Time
	Channel     string // in_app, email, push, sms
}

// NotificationFilter defines filters for notifications
type NotificationFilter struct {
	Type     string
	Category string
	IsRead   *bool
	Priority string
	DateFrom *time.Time
	DateTo   *time.Time
}

// NotificationStats represents notification statistics
type NotificationStats struct {
	TotalCount        int64            `json:"total_count"`
	UnreadCount       int64            `json:"unread_count"`
	ReadCount         int64            `json:"read_count"`
	TodayCount        int64            `json:"today_count"`
	ThisWeekCount     int64            `json:"this_week_count"`
	HighPriorityCount int64            `json:"high_priority_count"`
	CategoryBreakdown map[string]int64 `json:"category_breakdown"`
}

// NotificationRepository defines the interface for notification data operations
type NotificationRepository interface {
	// Create creates a new notification
	Create(ctx context.Context, notification *Notification) error

	// FindByID finds notification by ID
	FindByID(ctx context.Context, id int64) (*Notification, error)

	// Update updates notification
	Update(ctx context.Context, notification *Notification) error

	// Delete deletes notification
	Delete(ctx context.Context, id int64) error

	// ListByUser lists notifications for user
	ListByUser(ctx context.Context, userID int64, filter NotificationFilter, page, limit int) ([]Notification, int64, error)

	// GetUnreadByUser retrieves unread notifications for user
	GetUnreadByUser(ctx context.Context, userID int64, limit int) ([]Notification, error)

	// CountUnreadByUser counts unread notifications for user
	CountUnreadByUser(ctx context.Context, userID int64) (int64, error)

	// MarkAsRead marks notification as read
	MarkAsRead(ctx context.Context, id int64) error

	// MarkAllAsRead marks all notifications as read for user
	MarkAllAsRead(ctx context.Context, userID int64) error

	// DeleteByUser deletes all notifications for user
	DeleteByUser(ctx context.Context, userID int64) error

	// GetExpiredNotifications retrieves expired notifications
	GetExpiredNotifications(ctx context.Context, limit int) ([]Notification, error)

	// GetStats retrieves notification statistics
	GetStats(ctx context.Context, userID int64) (*NotificationStats, error)

	// BulkCreate creates multiple notifications
	BulkCreate(ctx context.Context, notifications []Notification) error

	// FindPreferenceByUser finds notification preferences for user
	FindPreferenceByUser(ctx context.Context, userID int64) (*NotificationPreference, error)

	// CreatePreference creates notification preferences
	CreatePreference(ctx context.Context, preference *NotificationPreference) error

	// UpdatePreference updates notification preferences
	UpdatePreference(ctx context.Context, preference *NotificationPreference) error
}

// ============================================================================
// FCM Push Notification Interfaces
// Reference: https://firebase.google.com/docs/cloud-messaging
// ============================================================================

// DeviceTokenRepository defines the interface for device token data operations
type DeviceTokenRepository interface {
	// Create creates a new device token
	Create(ctx context.Context, token *DeviceToken) error

	// FindByID finds device token by ID
	FindByID(ctx context.Context, id int64) (*DeviceToken, error)

	// FindByToken finds device token by FCM token string
	FindByToken(ctx context.Context, token string) (*DeviceToken, error)

	// FindByUser finds all active device tokens for a user
	FindByUser(ctx context.Context, userID int64) ([]DeviceToken, error)

	// FindByUserAndPlatform finds device tokens for a user on specific platform
	FindByUserAndPlatform(ctx context.Context, userID int64, platform Platform) ([]DeviceToken, error)

	// Update updates device token
	Update(ctx context.Context, token *DeviceToken) error

	// Delete deletes device token
	Delete(ctx context.Context, id int64) error

	// DeleteByToken deletes device token by token string
	DeleteByToken(ctx context.Context, token string) error

	// Deactivate deactivates a device token
	Deactivate(ctx context.Context, token string) error

	// FindInactiveTokens finds inactive tokens that haven't been used for specified days
	FindInactiveTokens(ctx context.Context, inactiveDays int, limit int) ([]DeviceToken, error)

	// FindInactive finds tokens inactive since the cutoff date
	FindInactive(ctx context.Context, cutoffDate time.Time) ([]DeviceToken, error)

	// FindByFailureCount finds tokens with failure count >= minFailures
	FindByFailureCount(ctx context.Context, minFailures int) ([]DeviceToken, error)

	// CountByUser counts active device tokens for a user
	CountByUser(ctx context.Context, userID int64) (int64, error)

	// BatchUpdate updates multiple device tokens
	BatchUpdate(ctx context.Context, tokens []DeviceToken) error

	// FindByUserIDs finds all active device tokens for multiple users (for batch sending)
	FindByUserIDs(ctx context.Context, userIDs []int64) ([]DeviceToken, error)

	// FindTokensWithHighFailureCount finds tokens with high failure counts for cleanup
	FindTokensWithHighFailureCount(ctx context.Context, minFailures int, limit int) ([]DeviceToken, error)
}

// PushNotificationService defines the interface for FCM push notification operations
type PushNotificationService interface {
	// SendToDevice sends push notification to a specific device token
	SendToDevice(ctx context.Context, token string, message *PushMessage) (*PushResult, error)

	// SendToUser sends push notification to all user's devices
	SendToUser(ctx context.Context, userID int64, message *PushMessage) ([]PushResult, error)

	// SendToMultipleUsers sends push notification to multiple users
	SendToMultipleUsers(ctx context.Context, userIDs []int64, message *PushMessage) (map[int64][]PushResult, error)

	// SendToTopic sends push notification to a topic
	SendToTopic(ctx context.Context, topic string, message *PushMessage) (*PushResult, error)

	// RegisterDeviceToken registers a new device token for a user
	RegisterDeviceToken(ctx context.Context, userID int64, token string, platform Platform, deviceInfo *DeviceInfo) error

	// UnregisterDeviceToken removes a device token
	UnregisterDeviceToken(ctx context.Context, userID int64, token string) error

	// RefreshDeviceToken updates an existing device token
	RefreshDeviceToken(ctx context.Context, oldToken, newToken string) error

	// GetUserDevices retrieves all registered devices for a user
	GetUserDevices(ctx context.Context, userID int64) ([]DeviceToken, error)

	// ValidateToken validates if a token is still valid with FCM
	ValidateToken(ctx context.Context, token string) (bool, error)

	// CleanupInactiveTokens removes inactive device tokens
	CleanupInactiveTokens(ctx context.Context, inactiveDays int) error
}

// PushNotificationLogRepository defines the interface for push notification log operations
type PushNotificationLogRepository interface {
	// Create creates a new push notification log
	Create(ctx context.Context, log *PushNotificationLog) error

	// FindByID finds log by ID
	FindByID(ctx context.Context, id int64) (*PushNotificationLog, error)

	// FindByNotificationID finds logs by notification ID
	FindByNotificationID(ctx context.Context, notificationID int64) ([]PushNotificationLog, error)

	// FindByUserID finds logs by user ID
	FindByUserID(ctx context.Context, userID int64, page, limit int) ([]PushNotificationLog, int64, error)

	// FindByFCMMessageID finds log by FCM message ID
	FindByFCMMessageID(ctx context.Context, messageID string) (*PushNotificationLog, error)

	// Update updates push notification log
	Update(ctx context.Context, log *PushNotificationLog) error

	// UpdateStatus updates log status
	UpdateStatus(ctx context.Context, id int64, status string) error

	// GetDeliveryStats retrieves delivery statistics
	GetDeliveryStats(ctx context.Context, userID int64, from, to time.Time) (map[string]int64, error)

	// CleanupOldLogs removes logs older than specified days
	CleanupOldLogs(ctx context.Context, retentionDays int) error
}
