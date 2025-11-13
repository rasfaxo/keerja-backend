package response

// AuthResponse represents authentication response with tokens
type AuthResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	TokenType    string        `json:"token_type"`
	ExpiresIn    int64         `json:"expires_in"`
	User         *UserBasic    `json:"user"`
	Company      *CompanyBasic `json:"company"`
}

// UserBasic represents basic user info in auth response
type UserBasic struct {
	ID         int64  `json:"id"`
	UUID       string `json:"uuid"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone,omitempty"`
	UserType   string `json:"user_type"`
	IsVerified bool   `json:"is_verified"`
	Status     string `json:"status"`
}

// CompanyBasic represents basic company info in auth response (for employer only)
type CompanyBasic struct {
	ID           int64  `json:"id"`
	UUID         string `json:"uuid"`
	CompanyName  string `json:"company_name"`
	Slug         string `json:"slug"`
	LogoURL      string `json:"logo_url,omitempty"`
	IsVerified   bool   `json:"verified"`
	Status       string `json:"status"`
	BadgeGranted bool   `json:"badge_granted"`
	NPWPNumber   string `json:"npwp_number,omitempty"`
	NIBNumber    string `json:"nib_number,omitempty"`
}

// TokenResponse represents token-only response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// VerificationResponse represents email verification response
type VerificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// LogoutResponse represents logout response
type LogoutResponse struct {
	Message string `json:"message"`
}

// PasswordResetResponse represents password reset response
type PasswordResetResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// DeviceInfo represents device session information
type DeviceInfo struct {
	ID         int64   `json:"id"`
	DeviceName *string `json:"device_name"`
	DeviceType *string `json:"device_type"`
	IPAddress  *string `json:"ip_address"`
	LastUsedAt string  `json:"last_used_at"`
	CreatedAt  string  `json:"created_at"`
	IsCurrent  bool    `json:"is_current"` // whether this is the current device
}

// DeviceListResponse represents list of active devices
type DeviceListResponse struct {
	Devices []DeviceInfo `json:"devices"`
	Total   int          `json:"total"`
}

// OAuthProviderResponse represents connected OAuth provider
type OAuthProviderResponse struct {
	ID          int64   `json:"id"`
	Provider    string  `json:"provider"`
	ProviderID  string  `json:"provider_id"`
	Email       *string `json:"email,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	ConnectedAt string  `json:"connected_at"`
	LastLoginAt *string `json:"last_login_at,omitempty"`
}
