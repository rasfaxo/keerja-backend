package response

// AuthResponse represents authentication response with tokens
type AuthResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	TokenType    string     `json:"token_type"`
	ExpiresIn    int64      `json:"expires_in"`
	User         *UserBasic `json:"user"`
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
