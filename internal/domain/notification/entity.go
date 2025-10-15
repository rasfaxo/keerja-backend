package notification

import "time"

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
