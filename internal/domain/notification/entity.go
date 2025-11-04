package notification

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Notification represents a notification
type Notification struct {
	ID          int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      int64      `json:"user_id" gorm:"not null;index"`
	Type        string     `json:"type" gorm:"type:varchar(50);not null;index"` // job_application, interview_scheduled, status_update, etc.
	Title       string     `json:"title" gorm:"type:varchar(255);not null"`
	Message     string     `json:"message" gorm:"type:text;not null"`
	Data        string     `json:"data" gorm:"type:json"` // Additional data as JSON
	IsRead      bool       `json:"is_read" gorm:"default:false;index"`
	ReadAt      *time.Time `json:"read_at"`
	Priority    string     `json:"priority" gorm:"type:varchar(20);default:'normal'"` // low, normal, high, urgent
	Category    string     `json:"category" gorm:"type:varchar(50);not null;index"`   // application, job, account, system
	ActionURL   string     `json:"action_url" gorm:"type:varchar(500)"`
	Icon        string     `json:"icon" gorm:"type:varchar(100)"`
	SenderID    *int64     `json:"sender_id" gorm:"index"`
	RelatedID   *int64     `json:"related_id" gorm:"index"`              // ID of related entity (job_id, application_id, etc.)
	RelatedType string     `json:"related_type" gorm:"type:varchar(50)"` // job, application, interview, etc.
	ExpiresAt   *time.Time `json:"expires_at"`
	IsSent      bool       `json:"is_sent" gorm:"default:false"`
	SentAt      *time.Time `json:"sent_at"`
	Channel     string     `json:"channel" gorm:"type:varchar(50)"` // in_app, email, push, sms
	CreatedAt   time.Time  `json:"created_at" gorm:"type:timestamp;default:now()"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"type:timestamp;default:now()"`
}

// TableName specifies the table name
func (Notification) TableName() string {
	return "notifications"
}

// IsRead checks if notification is read
func (n *Notification) IsReadStatus() bool {
	return n.IsRead
}

// IsUnread checks if notification is unread
func (n *Notification) IsUnread() bool {
	return !n.IsRead
}

// IsExpired checks if notification has expired
func (n *Notification) IsExpired() bool {
	if n.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*n.ExpiresAt)
}

// IsHighPriority checks if notification is high priority
func (n *Notification) IsHighPriority() bool {
	return n.Priority == "high" || n.Priority == "urgent"
}

// MarkAsRead marks notification as read
func (n *Notification) MarkAsRead() {
	n.IsRead = true
	now := time.Now()
	n.ReadAt = &now
}

// MarkAsUnread marks notification as unread
func (n *Notification) MarkAsUnread() {
	n.IsRead = false
	n.ReadAt = nil
}

// MarkAsSent marks notification as sent
func (n *Notification) MarkAsSent() {
	n.IsSent = true
	now := time.Now()
	n.SentAt = &now
}

// NotificationPreference represents user notification preferences
type NotificationPreference struct {
	ID                        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID                    int64     `json:"user_id" gorm:"not null;uniqueIndex"`
	EmailEnabled              bool      `json:"email_enabled" gorm:"default:true"`
	PushEnabled               bool      `json:"push_enabled" gorm:"default:true"`
	SMSEnabled                bool      `json:"sms_enabled" gorm:"default:false"`
	JobApplicationsEnabled    bool      `json:"job_applications_enabled" gorm:"default:true"`
	InterviewEnabled          bool      `json:"interview_enabled" gorm:"default:true"`
	StatusUpdatesEnabled      bool      `json:"status_updates_enabled" gorm:"default:true"`
	JobRecommendationsEnabled bool      `json:"job_recommendations_enabled" gorm:"default:true"`
	CompanyUpdatesEnabled     bool      `json:"company_updates_enabled" gorm:"default:true"`
	MarketingEnabled          bool      `json:"marketing_enabled" gorm:"default:false"`
	WeeklyDigestEnabled       bool      `json:"weekly_digest_enabled" gorm:"default:true"`
	CreatedAt                 time.Time `json:"created_at" gorm:"type:timestamp;default:now()"`
	UpdatedAt                 time.Time `json:"updated_at" gorm:"type:timestamp;default:now()"`
}

// TableName specifies the table name
func (NotificationPreference) TableName() string {
	return "notification_preferences"
}

// IsEmailEnabled checks if email notifications are enabled
func (np *NotificationPreference) IsEmailEnabled() bool {
	return np.EmailEnabled
}

// IsPushEnabled checks if push notifications are enabled
func (np *NotificationPreference) IsPushEnabled() bool {
	return np.PushEnabled
}

// CanSendNotification checks if notification can be sent for given type
func (np *NotificationPreference) CanSendNotification(notificationType string) bool {
	switch notificationType {
	case "job_application":
		return np.JobApplicationsEnabled
	case "interview":
		return np.InterviewEnabled
	case "status_update":
		return np.StatusUpdatesEnabled
	case "job_recommendation":
		return np.JobRecommendationsEnabled
	case "company_update":
		return np.CompanyUpdatesEnabled
	case "marketing":
		return np.MarketingEnabled
	default:
		return true
	}
}

// ============================================================================
// FCM Push Notification Entities
// Reference: https://firebase.google.com/docs/cloud-messaging
// ============================================================================

// Platform represents device platform types
type Platform string

const (
	PlatformAndroid Platform = "android"
	PlatformIOS     Platform = "ios"
	PlatformWeb     Platform = "web"
)

// DeviceInfo represents device metadata stored as JSONB
type DeviceInfo struct {
	Model      string `json:"model,omitempty"`       // Device model (e.g., "iPhone 14 Pro", "Pixel 7")
	OSVersion  string `json:"os_version,omitempty"`  // OS version (e.g., "iOS 16.3", "Android 13")
	AppVersion string `json:"app_version,omitempty"` // App version (e.g., "1.2.3")
	Language   string `json:"language,omitempty"`    // Device language
}

// Value implements the driver.Valuer interface for GORM JSONB
func (d DeviceInfo) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface for GORM JSONB
func (d *DeviceInfo) Scan(value interface{}) error {
	if value == nil {
		*d = DeviceInfo{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value")
	}

	return json.Unmarshal(bytes, d)
}

// DeviceToken represents a user's FCM device registration token
type DeviceToken struct {
	ID            int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        int64      `json:"user_id" gorm:"not null;index:idx_device_tokens_user_id"`
	Token         string     `json:"token" gorm:"type:varchar(4096);not null;index:idx_device_tokens_token,unique:unique_user_token"`
	Platform      Platform   `json:"platform" gorm:"type:varchar(20);not null;index:idx_device_tokens_platform;check:platform IN ('android', 'ios', 'web')"`
	DeviceInfo    DeviceInfo `json:"device_info" gorm:"type:jsonb;default:'{}'"`
	IsActive      bool       `json:"is_active" gorm:"default:true;not null;index:idx_device_tokens_user_id,idx_device_tokens_token,idx_device_tokens_user_platform"`
	LastUsedAt    *time.Time `json:"last_used_at" gorm:"index:idx_device_tokens_inactive"`
	FailureCount  int        `json:"failure_count" gorm:"default:0;not null;index:idx_device_tokens_failure"`
	LastFailureAt *time.Time `json:"last_failure_at" gorm:"index:idx_device_tokens_failure"`
	FailureReason string     `json:"failure_reason" gorm:"type:text"`
	CreatedAt     time.Time  `json:"created_at" gorm:"type:timestamp;default:now()"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"type:timestamp;default:now()"`
}

// TableName specifies the table name
func (DeviceToken) TableName() string {
	return "device_tokens"
}

// IsValid checks if device token is valid and active
func (dt *DeviceToken) IsValid() bool {
	return dt.IsActive && dt.Token != ""
}

// MarkAsUsed updates the last used timestamp
func (dt *DeviceToken) MarkAsUsed() {
	now := time.Now()
	dt.LastUsedAt = &now
	dt.UpdatedAt = now
}

// RecordFailure increments failure count and records the reason
func (dt *DeviceToken) RecordFailure(reason string) {
	dt.FailureCount++
	now := time.Now()
	dt.LastFailureAt = &now
	dt.FailureReason = reason
	dt.UpdatedAt = now

	// Auto-deactivate after 5 consecutive failures (per FCM best practices)
	if dt.FailureCount >= 5 {
		dt.IsActive = false
	}
}

// ResetFailures clears failure tracking
func (dt *DeviceToken) ResetFailures() {
	dt.FailureCount = 0
	dt.LastFailureAt = nil
	dt.FailureReason = ""
	dt.UpdatedAt = time.Now()
}

// Deactivate marks the token as inactive
func (dt *DeviceToken) Deactivate() {
	dt.IsActive = false
	dt.UpdatedAt = time.Now()
}

// Activate marks the token as active
func (dt *DeviceToken) Activate() {
	dt.IsActive = true
	dt.FailureCount = 0
	dt.LastFailureAt = nil
	dt.FailureReason = ""
	dt.UpdatedAt = time.Now()
}

// PushNotificationLog represents a log entry for push notification delivery
type PushNotificationLog struct {
	ID             int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	NotificationID *int64     `json:"notification_id" gorm:"index:idx_push_logs_notification_id"`
	DeviceTokenID  *int64     `json:"device_token_id" gorm:"index:idx_push_logs_device_token_id"`
	UserID         int64      `json:"user_id" gorm:"not null;index:idx_push_logs_user_id"`
	FCMMessageID   string     `json:"fcm_message_id" gorm:"type:varchar(255);index:idx_push_logs_fcm_message_id"`
	Status         string     `json:"status" gorm:"type:varchar(20);default:'pending';index:idx_push_logs_status"` // pending, sent, delivered, failed, clicked
	ErrorCode      string     `json:"error_code" gorm:"type:varchar(100)"`
	ErrorMessage   string     `json:"error_message" gorm:"type:text"`
	FCMResponse    string     `json:"fcm_response" gorm:"type:jsonb"` // Full response from FCM
	SentAt         *time.Time `json:"sent_at"`
	DeliveredAt    *time.Time `json:"delivered_at"`
	ClickedAt      *time.Time `json:"clicked_at"`
	CreatedAt      time.Time  `json:"created_at" gorm:"type:timestamp;default:now()"`
}

// TableName specifies the table name
func (PushNotificationLog) TableName() string {
	return "push_notification_logs"
}

// PushMessage represents a push notification message payload
type PushMessage struct {
	Title       string            `json:"title"`                  // Notification title
	Body        string            `json:"body"`                   // Notification body
	ImageURL    string            `json:"image_url,omitempty"`    // Optional image URL
	Data        map[string]string `json:"data,omitempty"`         // Additional data payload
	Priority    string            `json:"priority,omitempty"`     // normal, high
	Sound       string            `json:"sound,omitempty"`        // Notification sound
	Badge       *int              `json:"badge,omitempty"`        // Badge count (iOS)
	ClickAction string            `json:"click_action,omitempty"` // Action when clicked
	TTL         int               `json:"ttl,omitempty"`          // Time to live in seconds
}

// PushResult represents the result of sending a push notification
type PushResult struct {
	Success      bool   `json:"success"`
	MessageID    string `json:"message_id,omitempty"`    // FCM message ID
	ErrorCode    string `json:"error_code,omitempty"`    // Error code from FCM
	ErrorMessage string `json:"error_message,omitempty"` // Error message
}
