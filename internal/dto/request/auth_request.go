package request

// RegisterRequest represents user registration request
type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=150"`
	Email    string `json:"email" validate:"required,email,max=150"`
	Phone    string `json:"phone" validate:"omitempty,min=10,max=20"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	UserType string `json:"user_type" validate:"required,oneof=jobseeker employer"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// VerifyEmailRequest represents email verification request
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// ForgotPasswordRequest represents forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=72"`
}

// ChangePasswordRequest represents change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=72"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ResendVerificationRequest represents resend verification email request
type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ===========================================
// OTP Authentication Requests
// ===========================================

// RequestOTPRequest represents OTP request
type RequestOTPRequest struct {
	Contact string `json:"contact" validate:"required,email"`                                    // email address
	Purpose string `json:"purpose" validate:"omitempty,oneof=login verify_email reset_password"` // default: login
}

// VerifyOTPRequest represents OTP verification request
type VerifyOTPRequest struct {
	Contact string `json:"contact" validate:"required,email"`
	OTP     string `json:"otp" validate:"required,len=6,numeric"`
	Purpose string `json:"purpose" validate:"omitempty,oneof=login verify_email reset_password"`
}

// ===========================================
// OAuth2 / Social Login Requests
// ===========================================

// OAuthCallbackRequest represents OAuth callback data
type OAuthCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
}

// OAuthLoginRequest represents OAuth login initiation
type OAuthLoginRequest struct {
	Provider string `json:"provider" validate:"required,oneof=google facebook github"`
	UserType string `json:"user_type" validate:"required,oneof=jobseeker employer"` // for new users
}

// ===========================================
// OTP Verification for Registration
// ===========================================

// RegisterWithOTPRequest represents registration request with OTP verification
type RegisterWithOTPRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=150"`
	Email    string `json:"email" validate:"required,email,max=150"`
	Phone    string `json:"phone" validate:"omitempty,min=10,max=20"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	UserType string `json:"user_type" validate:"required,oneof=jobseeker employer"`
}

// VerifyEmailOTPRequest represents email verification with OTP
type VerifyEmailOTPRequest struct {
	Email   string `json:"email" validate:"required,email"`
	OTPCode string `json:"otp_code" validate:"required,len=6,numeric"`
}

// ResendOTPRequest represents resend OTP request
type ResendOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ===========================================
// Forgot Password with OTP
// ===========================================

// ForgotPasswordOTPRequest represents forgot password OTP request
type ForgotPasswordOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordOTPRequest represents reset password with OTP request
type ResetPasswordOTPRequest struct {
	Email       string `json:"email" validate:"required,email"`
	OTPCode     string `json:"otp_code" validate:"required,len=6,numeric"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=72"`
}

// ===========================================
// Refresh Token Requests
// ===========================================

// LoginWithRememberMeRequest represents login with remember me option
type LoginWithRememberMeRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	RememberMe bool   `json:"remember_me"` // optional, default false
	DeviceID   string `json:"device_id"`   // optional, client-generated unique ID
}

// RefreshAccessTokenRequest represents refresh token request
type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RevokeDeviceRequest represents revoke specific device request
type RevokeDeviceRequest struct {
	DeviceID string `json:"device_id" validate:"required"`
}
