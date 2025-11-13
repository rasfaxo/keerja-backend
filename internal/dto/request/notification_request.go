package request

import (
	"time"
)

// SendNotificationRequest represents a request to send a notification
type SendNotificationRequest struct {
	UserID      int64                  `json:"user_id" validate:"required"`
	Type        string                 `json:"type" validate:"required,max=50"`
	Title       string                 `json:"title" validate:"required,max=255"`
	Message     string                 `json:"message" validate:"required"`
	Data        map[string]interface{} `json:"data" validate:"omitempty"`
	Priority    string                 `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
	Category    string                 `json:"category" validate:"required,oneof=application job account system company"`
	ActionURL   string                 `json:"action_url" validate:"omitempty,max=500,url"`
	Icon        string                 `json:"icon" validate:"omitempty,max=100"`
	SenderID    *int64                 `json:"sender_id" validate:"omitempty"`
	RelatedID   *int64                 `json:"related_id" validate:"omitempty"`
	RelatedType string                 `json:"related_type" validate:"omitempty,max=50"`
	ExpiresAt   *time.Time             `json:"expires_at" validate:"omitempty"`
	Channel     string                 `json:"channel" validate:"omitempty,oneof=in_app email push sms"`
}

// SendBulkNotificationRequest represents a request to send notifications to multiple users
type SendBulkNotificationRequest struct {
	UserIDs     []int64                `json:"user_ids" validate:"required,min=1"`
	Type        string                 `json:"type" validate:"required,max=50"`
	Title       string                 `json:"title" validate:"required,max=255"`
	Message     string                 `json:"message" validate:"required"`
	Data        map[string]interface{} `json:"data" validate:"omitempty"`
	Priority    string                 `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
	Category    string                 `json:"category" validate:"required,oneof=application job account system company"`
	ActionURL   string                 `json:"action_url" validate:"omitempty,max=500,url"`
	Icon        string                 `json:"icon" validate:"omitempty,max=100"`
	SenderID    *int64                 `json:"sender_id" validate:"omitempty"`
	RelatedID   *int64                 `json:"related_id" validate:"omitempty"`
	RelatedType string                 `json:"related_type" validate:"omitempty,max=50"`
	ExpiresAt   *time.Time             `json:"expires_at" validate:"omitempty"`
	Channel     string                 `json:"channel" validate:"omitempty,oneof=in_app email push sms"`
}

// UpdateNotificationPreferencesRequest represents a request to update notification preferences
type UpdateNotificationPreferencesRequest struct {
	EmailEnabled              *bool `json:"email_enabled" validate:"omitempty"`
	PushEnabled               *bool `json:"push_enabled" validate:"omitempty"`
	SMSEnabled                *bool `json:"sms_enabled" validate:"omitempty"`
	JobApplicationsEnabled    *bool `json:"job_applications_enabled" validate:"omitempty"`
	InterviewEnabled          *bool `json:"interview_enabled" validate:"omitempty"`
	StatusUpdatesEnabled      *bool `json:"status_updates_enabled" validate:"omitempty"`
	JobRecommendationsEnabled *bool `json:"job_recommendations_enabled" validate:"omitempty"`
	CompanyUpdatesEnabled     *bool `json:"company_updates_enabled" validate:"omitempty"`
	MarketingEnabled          *bool `json:"marketing_enabled" validate:"omitempty"`
	WeeklyDigestEnabled       *bool `json:"weekly_digest_enabled" validate:"omitempty"`
}

// NotificationFilterRequest represents filters for notification queries
type NotificationFilterRequest struct {
	Type     string     `json:"type" validate:"omitempty,max=50"`
	Category string     `json:"category" validate:"omitempty,oneof=application job account system company"`
	IsRead   *bool      `json:"is_read" validate:"omitempty"`
	Priority string     `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
	DateFrom *time.Time `json:"date_from" validate:"omitempty"`
	DateTo   *time.Time `json:"date_to" validate:"omitempty"`
}

// MarkNotificationsAsReadRequest represents a request to mark notifications as read
type MarkNotificationsAsReadRequest struct {
	NotificationIDs []int64 `json:"notification_ids" validate:"required,min=1"`
}
