package email

import "time"

// EmailLog represents email sending history
type EmailLog struct {
	ID            int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Recipient     string     `json:"recipient" gorm:"type:varchar(255);not null;index"`
	Subject       string     `json:"subject" gorm:"type:varchar(500);not null"`
	Body          string     `json:"body" gorm:"type:text"`
	Template      string     `json:"template" gorm:"type:varchar(100)"`
	Status        string     `json:"status" gorm:"type:varchar(50);not null;index;default:'pending'"` // pending, sent, failed
	Provider      string     `json:"provider" gorm:"type:varchar(50)"`
	SentAt        *time.Time `json:"sent_at"`
	FailureReason string     `json:"failure_reason" gorm:"type:text"`
	Metadata      *string    `json:"metadata,omitempty" gorm:"type:jsonb"`
	RetryCount    int        `json:"retry_count" gorm:"default:0"`
	MaxRetries    int        `json:"max_retries" gorm:"default:3"`
	CreatedAt     time.Time  `json:"created_at" gorm:"type:timestamp;default:now()"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"type:timestamp;default:now()"`
}

// TableName specifies the table name
func (EmailLog) TableName() string {
	return "email_logs"
}

// IsSent checks if email was sent successfully
func (e *EmailLog) IsSent() bool {
	return e.Status == "sent"
}

// IsFailed checks if email sending failed
func (e *EmailLog) IsFailed() bool {
	return e.Status == "failed"
}

// IsPending checks if email is pending
func (e *EmailLog) IsPending() bool {
	return e.Status == "pending"
}

// CanRetry checks if email can be retried
func (e *EmailLog) CanRetry() bool {
	return e.IsFailed() && e.RetryCount < e.MaxRetries
}

// MarkAsSent marks email as sent
func (e *EmailLog) MarkAsSent() {
	e.Status = "sent"
	now := time.Now()
	e.SentAt = &now
}

// MarkAsFailed marks email as failed
func (e *EmailLog) MarkAsFailed(reason string) {
	e.Status = "failed"
	e.FailureReason = reason
	e.RetryCount++
}
