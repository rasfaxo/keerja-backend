package response

import (
	"time"
)

// NotificationResponse represents a notification response
type NotificationResponse struct {
	ID          int64                  `json:"id"`
	UserID      int64                  `json:"user_id"`
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Data        map[string]interface{} `json:"data,omitempty"`
	IsRead      bool                   `json:"is_read"`
	ReadAt      *time.Time             `json:"read_at,omitempty"`
	Priority    string                 `json:"priority"`
	Category    string                 `json:"category"`
	ActionURL   string                 `json:"action_url,omitempty"`
	Icon        string                 `json:"icon,omitempty"`
	SenderID    *int64                 `json:"sender_id,omitempty"`
	RelatedID   *int64                 `json:"related_id,omitempty"`
	RelatedType string                 `json:"related_type,omitempty"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
	IsSent      bool                   `json:"is_sent"`
	SentAt      *time.Time             `json:"sent_at,omitempty"`
	Channel     string                 `json:"channel"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// NotificationListResponse represents a paginated list of notifications
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Total         int64                  `json:"total"`
	Page          int                    `json:"page"`
	Limit         int                    `json:"limit"`
	TotalPages    int                    `json:"total_pages"`
}

// NotificationPreferenceResponse represents notification preferences
type NotificationPreferenceResponse struct {
	ID                        int64     `json:"id"`
	UserID                    int64     `json:"user_id"`
	EmailEnabled              bool      `json:"email_enabled"`
	PushEnabled               bool      `json:"push_enabled"`
	SMSEnabled                bool      `json:"sms_enabled"`
	JobApplicationsEnabled    bool      `json:"job_applications_enabled"`
	InterviewEnabled          bool      `json:"interview_enabled"`
	StatusUpdatesEnabled      bool      `json:"status_updates_enabled"`
	JobRecommendationsEnabled bool      `json:"job_recommendations_enabled"`
	CompanyUpdatesEnabled     bool      `json:"company_updates_enabled"`
	MarketingEnabled          bool      `json:"marketing_enabled"`
	WeeklyDigestEnabled       bool      `json:"weekly_digest_enabled"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

// NotificationStatsResponse represents notification statistics
type NotificationStatsResponse struct {
	TotalCount        int64            `json:"total_count"`
	UnreadCount       int64            `json:"unread_count"`
	ReadCount         int64            `json:"read_count"`
	TodayCount        int64            `json:"today_count"`
	ThisWeekCount     int64            `json:"this_week_count"`
	HighPriorityCount int64            `json:"high_priority_count"`
	CategoryBreakdown map[string]int64 `json:"category_breakdown"`
}

// UnreadCountResponse represents unread notification count
type UnreadCountResponse struct {
	UnreadCount int64 `json:"unread_count"`
}

// NotificationActionResponse represents a simple success response
type NotificationActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
