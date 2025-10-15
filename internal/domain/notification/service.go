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
	Type      string
	Category  string
	IsRead    *bool
	Priority  string
	DateFrom  *time.Time
	DateTo    *time.Time
}

// NotificationStats represents notification statistics
type NotificationStats struct {
	TotalCount          int64 `json:"total_count"`
	UnreadCount         int64 `json:"unread_count"`
	ReadCount           int64 `json:"read_count"`
	TodayCount          int64 `json:"today_count"`
	ThisWeekCount       int64 `json:"this_week_count"`
	HighPriorityCount   int64 `json:"high_priority_count"`
	CategoryBreakdown   map[string]int64 `json:"category_breakdown"`
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
