package response

import "time"

// PushNotificationResponse represents a single push notification delivery result
type PushNotificationResponse struct {
	MessageID    string    `json:"message_id,omitempty"`
	Token        string    `json:"token,omitempty"`
	Success      bool      `json:"success"`
	ErrorCode    string    `json:"error_code,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
	SentAt       time.Time `json:"sent_at"`
}

// BatchPushNotificationResponse represents results from batch push notification sending
type BatchPushNotificationResponse struct {
	TotalSent    int                        `json:"total_sent"`
	TotalFailed  int                        `json:"total_failed"`
	TotalUsers   int                        `json:"total_users"`
	SuccessCount int                        `json:"success_count"`
	FailureCount int                        `json:"failure_count"`
	Results      []PushNotificationResponse `json:"results,omitempty"`
	SentAt       time.Time                  `json:"sent_at"`
}

// PushValidationResponse represents the result of device token validation
type PushValidationResponse struct {
	Token    string `json:"token"`
	IsValid  bool   `json:"is_valid"`
	Platform string `json:"platform,omitempty"`
	Message  string `json:"message,omitempty"`
}

// PushDeliveryStatsResponse represents push notification delivery statistics
type PushDeliveryStatsResponse struct {
	TotalSent      int     `json:"total_sent"`
	TotalDelivered int     `json:"total_delivered"`
	TotalFailed    int     `json:"total_failed"`
	DeliveryRate   float64 `json:"delivery_rate"`
	FailureRate    float64 `json:"failure_rate"`
}

// TopicSubscriptionResponse represents the result of topic subscription
type TopicSubscriptionResponse struct {
	Topic        string `json:"topic"`
	Success      bool   `json:"success"`
	TokensCount  int    `json:"tokens_count"`
	ErrorMessage string `json:"error_message,omitempty"`
}
