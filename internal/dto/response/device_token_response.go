package response

import (
	"encoding/json"
	"time"
)

// DeviceTokenResponse represents a device token in API responses
type DeviceTokenResponse struct {
	ID            int64           `json:"id" example:"1"`
	UserID        int64           `json:"user_id" example:"123"`
	Token         string          `json:"token" example:"dYW2s3xZR4e..."`
	Platform      string          `json:"platform" example:"android"`
	DeviceInfo    json.RawMessage `json:"device_info,omitempty" swaggertype:"string"`
	IsActive      bool            `json:"is_active" example:"true"`
	LastUsedAt    *time.Time      `json:"last_used_at,omitempty" example:"2025-11-01T10:00:00Z"`
	FailureCount  int             `json:"failure_count" example:"0"`
	LastFailureAt *time.Time      `json:"last_failure_at,omitempty"`
	FailureReason *string         `json:"failure_reason,omitempty"`
	CreatedAt     time.Time       `json:"created_at" example:"2025-10-01T10:00:00Z"`
	UpdatedAt     time.Time       `json:"updated_at" example:"2025-11-01T10:00:00Z"`
}

// DeviceTokenListResponse represents a paginated list of device tokens
type DeviceTokenListResponse struct {
	Data       []DeviceTokenResponse `json:"data"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

// DeviceTokenStatsResponse represents device token statistics
type DeviceTokenStatsResponse struct {
	TotalTokens    int `json:"total_tokens"`
	ActiveTokens   int `json:"active_tokens"`
	InactiveTokens int `json:"inactive_tokens"`
	AndroidTokens  int `json:"android_tokens"`
	IOSTokens      int `json:"ios_tokens"`
	WebTokens      int `json:"web_tokens"`
}

// DeviceTokenCountResponse represents device token count for a user
type DeviceTokenCountResponse struct {
	UserID       int64 `json:"user_id"`
	TotalTokens  int   `json:"total_tokens"`
	ActiveTokens int   `json:"active_tokens"`
}
