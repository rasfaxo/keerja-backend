package auth

import "time"

// OAuthProvider represents OAuth provider connection entity
type OAuthProvider struct {
	ID             int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         int64      `gorm:"not null" json:"user_id"`
	Provider       string     `gorm:"type:text;not null" json:"provider"` // google, facebook, github
	ProviderUserID string     `gorm:"type:text;not null" json:"provider_user_id"`
	Email          *string    `gorm:"type:text" json:"email,omitempty"`
	Name           *string    `gorm:"type:text" json:"name,omitempty"`
	AvatarURL      *string    `gorm:"type:text" json:"avatar_url,omitempty"`
	AccessToken    *string    `gorm:"type:text" json:"-"` // never expose
	RefreshToken   *string    `gorm:"type:text" json:"-"` // never expose
	TokenExpiry    *time.Time `gorm:"type:timestamptz" json:"token_expiry,omitempty"`
	RawData        *string    `gorm:"type:jsonb" json:"-"` // store full profile
	IsActive       bool       `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time  `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"type:timestamptz;default:now()" json:"updated_at"`
}

// TableName specifies the table name for OAuthProvider
func (OAuthProvider) TableName() string {
	return "oauth_providers"
}

// IsGoogle checks if provider is Google
func (o *OAuthProvider) IsGoogle() bool {
	return o.Provider == "google"
}

// IsTokenExpired checks if access token is expired
func (o *OAuthProvider) IsTokenExpired() bool {
	if o.TokenExpiry == nil {
		return true
	}
	return time.Now().After(*o.TokenExpiry)
}

// OTPCode represents OTP verification codes for user registration/verification
type OTPCode struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64      `gorm:"not null;index" json:"user_id"`
	OTPHash   string     `gorm:"type:text;not null" json:"-"`           // SHA256 hash, never expose
	Type      string     `gorm:"type:varchar(50);not null" json:"type"` // email_verification, password_reset
	ExpiredAt time.Time  `gorm:"type:timestamptz;not null" json:"expired_at"`
	IsUsed    bool       `gorm:"default:false" json:"is_used"`
	UsedAt    *time.Time `gorm:"type:timestamptz" json:"used_at,omitempty"`
	Attempts  int        `gorm:"default:0" json:"attempts"`
	CreatedAt time.Time  `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:timestamptz;default:now()" json:"updated_at"`
}

// TableName specifies the table name for OTPCode
func (OTPCode) TableName() string {
	return "otp_codes"
}

// IsExpired checks if OTP code is expired
func (o *OTPCode) IsExpired() bool {
	return time.Now().After(o.ExpiredAt)
}

// CanAttempt checks if more attempts are allowed
func (o *OTPCode) CanAttemptVerification(maxAttempts int) bool {
	return o.Attempts < maxAttempts
}

// IsValid checks if OTP is valid for verification
func (o *OTPCode) IsValid() bool {
	return !o.IsUsed && !o.IsExpired()
}

// RefreshToken represents a refresh token for persistent device sessions
type RefreshToken struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64      `gorm:"not null;index" json:"user_id"`
	TokenHash     string     `gorm:"type:text;not null;uniqueIndex" json:"-"` // SHA256 hash, never expose
	DeviceName    *string    `gorm:"type:varchar(255)" json:"device_name,omitempty"`
	DeviceType    *string    `gorm:"type:varchar(50)" json:"device_type,omitempty"` // mobile, desktop, tablet, unknown
	DeviceID      *string    `gorm:"type:varchar(255);index" json:"device_id,omitempty"`
	UserAgent     *string    `gorm:"type:text" json:"user_agent,omitempty"`
	IPAddress     *string    `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	LastUsedAt    time.Time  `gorm:"type:timestamptz;default:now()" json:"last_used_at"`
	ExpiresAt     time.Time  `gorm:"type:timestamptz;not null" json:"expires_at"`
	Revoked       bool       `gorm:"default:false" json:"revoked"`
	RevokedAt     *time.Time `gorm:"type:timestamptz" json:"revoked_at,omitempty"`
	RevokedReason *string    `gorm:"type:varchar(255)" json:"revoked_reason,omitempty"`
	CreatedAt     time.Time  `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"type:timestamptz;default:now()" json:"updated_at"`
}

// TableName specifies the table name for RefreshToken
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsExpired checks if refresh token is expired
func (r *RefreshToken) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// IsValid checks if refresh token is valid (not expired and not revoked)
func (r *RefreshToken) IsValid() bool {
	return !r.IsExpired() && !r.Revoked
}

// Revoke marks the refresh token as revoked
func (r *RefreshToken) Revoke(reason string) {
	now := time.Now()
	r.Revoked = true
	r.RevokedAt = &now
	r.RevokedReason = &reason
}

// UpdateLastUsed updates the last used timestamp
func (r *RefreshToken) UpdateLastUsed() {
	r.LastUsedAt = time.Now()
}
