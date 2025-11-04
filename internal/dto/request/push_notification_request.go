package request

import "encoding/json"

// SendPushNotificationRequest represents base request for sending push notification
type SendPushNotificationRequest struct {
	Title    string          `json:"title" validate:"required,min=1,max=100" example:"New Job Alert"`
	Body     string          `json:"body" validate:"required,min=1,max=500" example:"You have a new job match"`
	Data     json.RawMessage `json:"data,omitempty" swaggertype:"string"`
	ImageURL *string         `json:"image_url,omitempty" validate:"omitempty,url"`
	Sound    *string         `json:"sound,omitempty"`
	Priority *string         `json:"priority,omitempty" validate:"omitempty,oneof=high normal low"`
	Badge    *int            `json:"badge,omitempty" validate:"omitempty,min=0"`
}

// SendPushToDeviceRequest represents a request to send push notification to a specific device
type SendPushToDeviceRequest struct {
	SendPushNotificationRequest
	Token string `json:"token" validate:"required,min=10,max=4096" example:"dYW2s3xZR4e..."`
}

// SendPushToUserRequest represents a request to send push notification to a user
type SendPushToUserRequest struct {
	SendPushNotificationRequest
	UserID int64 `json:"user_id" validate:"required,gt=0" example:"123"`
}

// SendPushToMultipleUsersRequest represents a request to send push notification to multiple users
type SendPushToMultipleUsersRequest struct {
	SendPushNotificationRequest
	UserIDs []int64 `json:"user_ids" validate:"required,min=1,max=1000,dive,gt=0"`
}

// SendPushToTopicRequest represents a request to send push notification to a topic
type SendPushToTopicRequest struct {
	SendPushNotificationRequest
	Topic string `json:"topic" validate:"required,min=1,max=100" example:"job_alerts"`
}

// TestPushNotificationRequest represents a request to send test push notification
type TestPushNotificationRequest struct {
	Token string `json:"token" validate:"required,min=10,max=4096"`
}
