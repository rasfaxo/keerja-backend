package request

// AdminLoginRequest represents admin login request
type AdminLoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=150"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

// AdminChangePasswordRequest represents admin change password request
type AdminChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=6,max=100"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=100"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// AdminRefreshTokenRequest represents admin refresh token request
type AdminRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
