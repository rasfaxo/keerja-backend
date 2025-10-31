package request

import "encoding/json"

// RegisterDeviceTokenRequest represents a request to register a device token for push notifications
type RegisterDeviceTokenRequest struct {
	Token      string          `json:"token" validate:"required,min=10,max=4096" example:"dYW2s3xZR4e..."`   // FCM device token
	Platform   string          `json:"platform" validate:"required,oneof=android ios web" example:"android"` // Device platform
	DeviceInfo json.RawMessage `json:"device_info,omitempty" swaggertype:"string"`                           // Device information as JSON string
}

// UpdateDeviceTokenRequest represents a request to update device token information
type UpdateDeviceTokenRequest struct {
	DeviceInfo json.RawMessage `json:"device_info,omitempty" swaggertype:"string"` // Device information as JSON string
}

// DeviceTokenFilterRequest represents filters for device token queries
type DeviceTokenFilterRequest struct {
	UserID   *int64  `query:"user_id"`
	Platform *string `query:"platform" validate:"omitempty,oneof=android ios web"`
	IsActive *bool   `query:"is_active"`
	Page     int     `query:"page" validate:"min=1" default:"1"`
	PageSize int     `query:"page_size" validate:"min=1,max=100" default:"20"`
}

// ValidateDeviceTokenRequest represents a request to validate a device token
type ValidateDeviceTokenRequest struct {
	Token string `json:"token" validate:"required,min=10,max=4096" example:"dYW2s3xZR4e..."` // FCM device token to validate
}
