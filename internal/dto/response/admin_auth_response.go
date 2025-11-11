package response

import "time"

// AdminAuthResponse represents admin authentication response
type AdminAuthResponse struct {
	AdminID         int64              `json:"admin_id"`
	UUID            string             `json:"uuid"`
	FullName        string             `json:"full_name"`
	Email           string             `json:"email"`
	Phone           string             `json:"phone,omitempty"`
	Status          string             `json:"status"`
	ProfileImageURL string             `json:"profile_image_url,omitempty"`
	Role            *AdminRoleInfo     `json:"role,omitempty"`
	AccessToken     string             `json:"access_token"`
	RefreshToken    string             `json:"refresh_token,omitempty"`
	TokenType       string             `json:"token_type"`
	ExpiresIn       int64              `json:"expires_in"` // seconds
	LastLogin       *time.Time         `json:"last_login,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
}

// AdminRoleInfo represents role information in auth response
type AdminRoleInfo struct {
	RoleID          int64  `json:"role_id"`
	RoleName        string `json:"role_name"`
	RoleDescription string `json:"role_description,omitempty"`
	AccessLevel     int16  `json:"access_level"`
	IsSystemRole    bool   `json:"is_system_role"`
}

// AdminProfileResponse represents current admin profile
type AdminProfileResponse struct {
	AdminID         int64          `json:"admin_id"`
	UUID            string         `json:"uuid"`
	FullName        string         `json:"full_name"`
	Email           string         `json:"email"`
	Phone           string         `json:"phone,omitempty"`
	Status          string         `json:"status"`
	ProfileImageURL string         `json:"profile_image_url,omitempty"`
	Role            *AdminRoleInfo `json:"role,omitempty"`
	LastLogin       *time.Time     `json:"last_login,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// AdminTokenResponse represents token refresh response
type AdminTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
}
