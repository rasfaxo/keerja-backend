package http

import (
	"context"
	"keerja-backend/internal/domain/auth"
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/service"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService         *service.AuthService
	oauthService        *service.OAuthService
	registrationService *service.RegistrationService
	refreshTokenService *service.RefreshTokenService
	userRepo            user.UserRepository
	companyRepo         company.CompanyRepository
}

func NewAuthHandler(
	authService *service.AuthService,
	oauthService *service.OAuthService,
	registrationService *service.RegistrationService,
	refreshTokenService *service.RefreshTokenService,
	userRepo user.UserRepository,
	companyRepo company.CompanyRepository,
) *AuthHandler {
	return &AuthHandler{
		authService:         authService,
		oauthService:        oauthService,
		registrationService: registrationService,
		refreshTokenService: refreshTokenService,
		userRepo:            userRepo,
		companyRepo:         companyRepo,
	}
}

// buildAuthResponse builds auth response with company info for employer
func (h *AuthHandler) buildAuthResponse(ctx context.Context, usr *user.User, accessToken, refreshToken string) *response.AuthResponse {
	// If user is employer, include company info
	if usr.UserType == "employer" {
		companies, err := h.companyRepo.GetCompaniesByUserID(ctx, usr.ID)
		if err != nil {
			// Log error but continue without company
			println("DEBUG: Error getting companies for user", usr.ID, ":", err.Error())
		} else if len(companies) > 0 {
			// Get first company (primary company)
			println("DEBUG: Found", len(companies), "companies for user", usr.ID)
			return mapper.ToAuthResponseWithCompany(usr, &companies[0], accessToken, refreshToken)
		} else {
			println("DEBUG: No companies found for employer user", usr.ID)
		}
	}

	// Jobseeker or employer without company
	return mapper.ToAuthResponse(usr, accessToken, refreshToken)
}

// Register godoc
// @Summary Register new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "Register request"
// @Success 201 {object} utils.Response{data=response.AuthResponse}
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Sanitize text input fields
	req.FullName = utils.SanitizeString(req.FullName)
	req.Email = utils.SanitizeString(req.Email)
	if req.Phone != "" {
		req.Phone = utils.SanitizeString(req.Phone)
	}

	// Convert to domain request
	phone := &req.Phone
	if req.Phone == "" {
		phone = nil
	}

	domainReq := &user.RegisterRequest{
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    phone,
		Password: req.Password,
		UserType: req.UserType,
	}

	// Register user
	usr, verificationToken, err := h.authService.Register(ctx, domainReq)
	if err != nil {
		if err == service.ErrEmailAlreadyExists {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Email already exists", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to register user", err.Error())
	}

	// Convert to response DTO
	response := mapper.ToAuthResponse(usr, "", verificationToken)

	return utils.CreatedResponse(c, "User registered successfully. Please check your email for verification.", response)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and get access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login request"
// @Success 200 {object} utils.Response{data=response.AuthResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Sanitize input
	req.Email = utils.SanitizeString(req.Email)

	// Login user
	usr, accessToken, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password", err.Error())
		}
		if err == service.ErrEmailNotVerified {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Email not verified", "Please verify your email before logging in")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to login", err.Error())
	}

	// Build auth response with company info if employer
	authResponse := h.buildAuthResponse(ctx, usr, accessToken, "")

	return utils.SuccessResponse(c, "Login successful", authResponse)
}

// VerifyEmail godoc
// @Summary Verify email
// @Description Verify user's email address with token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.VerifyEmailRequest true "Verify email request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Sanitize input
	req.Token = utils.SanitizeString(req.Token)

	// Verify email
	if err := h.authService.VerifyEmail(ctx, req.Token); err != nil {
		if err == service.ErrInvalidVerificationToken {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid verification token", err.Error())
		}
		if err == service.ErrTokenExpired {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Verification token expired", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to verify email", err.Error())
	}

	return utils.SuccessResponse(c, "Email verified successfully", nil)
}

// ForgotPassword godoc
// @Summary Forgot password
// @Description Request password reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Sanitize input
	req.Email = utils.SanitizeString(req.Email)

	// Request password reset
	if err := h.authService.ForgotPassword(ctx, req.Email); err != nil {
		// Don't expose whether user exists or not
		return utils.SuccessResponse(c, "If the email exists, a password reset link has been sent", nil)
	}

	return utils.SuccessResponse(c, "Password reset link sent to your email", nil)
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset user password with token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Sanitize input
	req.Token = utils.SanitizeString(req.Token)

	// Reset password
	if err := h.authService.ResetPassword(ctx, req.Token, req.NewPassword); err != nil {
		if err == service.ErrInvalidResetToken {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid reset token", err.Error())
		}
		if err == service.ErrTokenExpired {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Reset token expired", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to reset password", err.Error())
	}

	return utils.SuccessResponse(c, "Password reset successfully", nil)
}

// RefreshToken godoc
// @Summary Refresh token
// @Description Get new access token using current token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=response.AuthResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", "")
	}

	// Refresh token
	accessToken, err := h.authService.RefreshToken(ctx, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to refresh token", err.Error())
	}

	// Get user info for response
	usr, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get user info", err.Error())
	}

	// Build auth response with company info if employer
	authResponse := h.buildAuthResponse(ctx, usr, accessToken, "")

	return utils.SuccessResponse(c, "Token refreshed successfully", authResponse)
}

// Logout godoc
// @Summary Logout user
// @Description Logout user and invalidate token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	// Logout user
	if err := h.authService.Logout(ctx, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to logout", err.Error())
	}

	return utils.SuccessResponse(c, "Logout successful", nil)
}

// ResendVerification godoc
// @Summary Resend verification email
// @Description Resend email verification link
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.ResendVerificationRequest true "Resend verification request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/resend-verification [post]
func (h *AuthHandler) ResendVerification(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.ResendVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Sanitize input
	req.Email = utils.SanitizeString(req.Email)

	// Resend verification email
	if err := h.authService.ResendVerificationEmail(ctx, req.Email); err != nil {
		// Don't expose whether user exists or not
		return utils.SuccessResponse(c, "If the email exists and is not verified, a verification link has been sent", nil)
	}

	return utils.SuccessResponse(c, "Verification email sent successfully", nil)
}

// ===========================================
// OTP Registration Handlers
// ===========================================

// RegisterWithOTP godoc
// @Summary Register with OTP verification
// @Description Register new user and send OTP code for email verification
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "Register request"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/register-otp [post]
func (h *AuthHandler) RegisterWithOTP(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Sanitize input
	req.FullName = utils.SanitizeString(req.FullName)
	req.Email = utils.SanitizeString(req.Email)
	if req.Phone != "" {
		req.Phone = utils.SanitizeString(req.Phone)
	}

	// Register user with OTP (service expects individual params)
	if err := h.registrationService.RegisterUser(ctx, req.FullName, req.Email, req.Password, req.Phone, req.UserType); err != nil {
		if err == service.ErrEmailAlreadyExists {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Email already exists", err.Error())
		}
		if err == service.ErrTooManyOTPRequests {
			return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "Too many OTP requests. Please try again later.", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to register user", err.Error())
	}

	return utils.CreatedResponse(c, "Registration successful. Please check your email for OTP verification code.", fiber.Map{
		"email": req.Email,
		"note":  "OTP code is valid for 5 minutes.",
	})
}

// VerifyEmailOTP godoc
// @Summary Verify email with OTP code
// @Description Verify user's email address using OTP code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.VerifyOTPRequest true "Verify OTP request"
// @Success 200 {object} utils.Response{data=response.AuthResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/verify-email-otp [post]
func (h *AuthHandler) VerifyEmailOTP(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.VerifyEmailOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Sanitize input
	req.Email = utils.SanitizeString(req.Email)
	req.OTPCode = utils.SanitizeString(req.OTPCode)

	// Verify OTP (returns accessToken, user, error)
	accessToken, usr, err := h.registrationService.VerifyEmailOTP(ctx, req.Email, req.OTPCode)
	if err != nil {
		if err == service.ErrInvalidOTPCode {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid OTP code", err.Error())
		}
		if err == service.ErrOTPCodeExpired {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "OTP code has expired", err.Error())
		}
		if err == service.ErrOTPCodeAlreadyUsed {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "OTP code has already been used", err.Error())
		}
		if err == service.ErrTooManyOTPAttempts {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Too many failed attempts. Please request a new OTP.", err.Error())
		}
		if err == service.ErrUserAlreadyVerified {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email already verified", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to verify OTP", err.Error())
	}

	// Build auth response with company info if employer
	authResponse := h.buildAuthResponse(ctx, usr, accessToken, "")

	return utils.SuccessResponse(c, "Email verified successfully", authResponse)
}

// ResendOTP godoc
// @Summary Resend OTP code
// @Description Resend OTP verification code to email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.ResendOTPRequest true "Resend OTP request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 429 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/resend-otp [post]
func (h *AuthHandler) ResendOTP(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.ResendOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Sanitize input
	req.Email = utils.SanitizeString(req.Email)

	// Resend OTP
	if err := h.registrationService.ResendOTP(ctx, req.Email); err != nil {
		if err == service.ErrUserAlreadyVerified {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email already verified", err.Error())
		}
		if err == service.ErrResendTooSoon {
			return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "Please wait before requesting a new OTP", err.Error())
		}
		if err == service.ErrTooManyOTPRequests {
			return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "Too many OTP requests. Please try again later.", err.Error())
		}
		// Don't expose whether user exists
		return utils.SuccessResponse(c, "If the email exists and is not verified, a new OTP has been sent", nil)
	}

	return utils.SuccessResponse(c, "OTP code has been resent to your email.", fiber.Map{
		"email": req.Email,
		"note":  "OTP code is valid for 5 minutes.",
	})
}

// ===========================================
// Forgot Password with OTP Handlers
// ===========================================

// ForgotPasswordOTP godoc
// @Summary Request password reset OTP
// @Description Send OTP code to email for password reset
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.ForgotPasswordOTPRequest true "Forgot password OTP request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 429 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/forgot-password-otp [post]
func (h *AuthHandler) ForgotPasswordOTP(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.ForgotPasswordOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Sanitize input
	req.Email = utils.SanitizeString(req.Email)

	// Request password reset OTP
	if err := h.registrationService.RequestPasswordResetOTP(ctx, req.Email); err != nil {
		if err == service.ErrTooManyOTPRequests {
			return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "Too many password reset requests. Please try again later.", err.Error())
		}
		if err == service.ErrEmailNotVerified {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email not verified", err.Error())
		}
		// Don't expose specific errors for security (silent success)
		return utils.SuccessResponse(c, "If the email exists, a password reset OTP has been sent.", fiber.Map{
			"email": req.Email,
			"note":  "OTP code is valid for 5 minutes.",
		})
	}

	return utils.SuccessResponse(c, "Password reset OTP has been sent to your email.", fiber.Map{
		"email": req.Email,
		"note":  "OTP code is valid for 5 minutes.",
	})
}

// ResetPasswordOTP godoc
// @Summary Reset password with OTP
// @Description Reset password using OTP verification
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.ResetPasswordOTPRequest true "Reset password OTP request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/reset-password-otp [post]
func (h *AuthHandler) ResetPasswordOTP(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.ResetPasswordOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Sanitize input
	req.Email = utils.SanitizeString(req.Email)
	req.OTPCode = utils.SanitizeString(req.OTPCode)

	// Reset password with OTP
	if err := h.registrationService.ResetPasswordWithOTP(ctx, req.Email, req.OTPCode, req.NewPassword); err != nil {
		if err == service.ErrInvalidOTPCode {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid OTP code", err.Error())
		}
		if err == service.ErrOTPCodeExpired {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "OTP code has expired", err.Error())
		}
		if err == service.ErrTooManyOTPAttempts {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Too many failed attempts. Please request a new OTP.", err.Error())
		}
		if err == service.ErrInvalidCredentials {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or OTP", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to reset password", err.Error())
	}

	return utils.SuccessResponse(c, "Password has been reset successfully. You can now login with your new password.", nil)
}

// ===========================================
// Refresh Token Handlers (Remember Me)
// ===========================================

// LoginWithRememberMe godoc
// @Summary Login with remember me option
// @Description Login and get refresh token for persistent sessions
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.LoginWithRememberMeRequest true "Login with remember me"
// @Success 200 {object} utils.Response{data=response.AuthResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/login-remember [post]
func (h *AuthHandler) LoginWithRememberMe(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.LoginWithRememberMeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Sanitize input
	req.Email = utils.SanitizeString(req.Email)

	// Login user
	usr, accessToken, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password", err.Error())
		}
		if err == service.ErrEmailNotVerified {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Email not verified", "Please verify your email before logging in")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to login", err.Error())
	}

	// Create device info
	deviceInfo := service.DeviceInfo{
		DeviceID:  req.DeviceID,
		UserAgent: string(c.Request().Header.UserAgent()),
		IPAddress: c.IP(),
	}

	// Create refresh token
	refreshToken, err := h.refreshTokenService.CreateRefreshToken(ctx, usr.ID, deviceInfo, req.RememberMe)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create refresh token", err.Error())
	}

	// Build auth response with company info if employer
	authResponse := h.buildAuthResponse(ctx, usr, accessToken, refreshToken)

	return utils.SuccessResponse(c, "Login successful", authResponse)
}

// RefreshAccessToken godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.RefreshAccessTokenRequest true "Refresh token request"
// @Success 200 {object} utils.Response{data=response.TokenResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshAccessToken(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.RefreshAccessTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Get user from context (if available from middleware)
	claims := c.Locals("user")
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "No user context found")
	}

	userClaims := claims.(*utils.Claims)

	// Refresh access token
	newAccessToken, newRefreshToken, err := h.refreshTokenService.RefreshAccessToken(
		ctx,
		req.RefreshToken,
		userClaims.UserID,
		userClaims.Email,
		userClaims.UserType,
	)
	if err != nil {
		switch err {
		case service.ErrRefreshTokenNotFound:
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid refresh token", err.Error())
		case service.ErrRefreshTokenExpired:
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Refresh token expired", err.Error())
		case service.ErrRefreshTokenRevoked:
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Refresh token has been revoked", err.Error())
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to refresh token", err.Error())
		}
	}

	// Return tokens
	tokenResponse := mapper.ToTokenResponse(newAccessToken, newRefreshToken)

	return utils.SuccessResponse(c, "Token refreshed successfully", tokenResponse)
}

// GetActiveDevices godoc
// @Summary Get user's active devices
// @Description Get list of all active device sessions
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=response.DeviceListResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/devices [get]
func (h *AuthHandler) GetActiveDevices(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user from context
	claims := c.Locals("user")
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "No user context found")
	}

	userClaims := claims.(*utils.Claims)

	// Get active devices
	devices, err := h.refreshTokenService.GetUserDevices(ctx, userClaims.UserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get devices", err.Error())
	}

	// Convert []auth.RefreshToken to []*auth.RefreshToken for mapper
	devicePointers := make([]*auth.RefreshToken, len(devices))
	for i := range devices {
		devicePointers[i] = &devices[i]
	}

	// Convert to response
	deviceList := mapper.ToDeviceListResponse(devicePointers)

	return utils.SuccessResponse(c, "Active devices retrieved successfully", deviceList)
}

// RevokeDevice godoc
// @Summary Revoke specific device session
// @Description Revoke refresh token for a specific device
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.RevokeDeviceRequest true "Device ID to revoke"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/devices/revoke [post]
func (h *AuthHandler) RevokeDevice(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user from context
	claims := c.Locals("user")
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "No user context found")
	}

	userClaims := claims.(*utils.Claims)

	// Parse request
	var req request.RevokeDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Revoke device
	if err := h.refreshTokenService.RevokeDeviceToken(ctx, userClaims.UserID, req.DeviceID, "user_revoked"); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to revoke device", err.Error())
	}

	return utils.SuccessResponse(c, "Device session revoked successfully", nil)
}

// LogoutAllDevices godoc
// @Summary Logout from all devices
// @Description Revoke all refresh tokens (logout from all devices)
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/logout-all [post]
func (h *AuthHandler) LogoutAllDevices(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user from context
	claims := c.Locals("user")
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "No user context found")
	}

	userClaims := claims.(*utils.Claims)

	// Revoke all tokens
	if err := h.refreshTokenService.RevokeAllUserTokens(ctx, userClaims.UserID, "logout_all_devices"); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to logout from all devices", err.Error())
	}

	return utils.SuccessResponse(c, "Logged out from all devices successfully", nil)
}

// ===========================================
// OAuth Handlers
// ===========================================

// InitiateGoogleLogin godoc
// @Summary Initiate Google OAuth login
// @Description Get Google OAuth authorization URL
// @Tags auth
// @Produce json
// @Success 200 {object} utils.Response{data=map[string]string}
// @Failure 500 {object} utils.Response
// @Router /auth/oauth/google [get]
func (h *AuthHandler) InitiateGoogleLogin(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get Google auth URL
	authURL, err := h.oauthService.GetGoogleAuthURL(ctx)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate auth URL", err.Error())
	}

	return utils.SuccessResponse(c, "Google auth URL generated", fiber.Map{
		"auth_url": authURL,
	})
}

// HandleGoogleCallback godoc
// @Summary Handle Google OAuth callback
// @Description Process Google OAuth callback and authenticate user
// @Tags auth
// @Produce json
// @Param code query string true "Authorization code"
// @Param state query string true "State token"
// @Success 200 {object} utils.Response{data=response.AuthResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/oauth/google/callback [get]
func (h *AuthHandler) HandleGoogleCallback(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get code and state from query
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Missing authorization code", "")
	}

	// Authenticate with Google (returns JWT token)
	accessToken, err := h.oauthService.HandleGoogleCallback(ctx, code, state)
	if err != nil {
		// Check for known service errors
		if err.Error() == "invalid state" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid OAuth state", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to authenticate with Google", err.Error())
	}

	// Return token response (no refresh token in basic OAuth flow)
	response := mapper.ToTokenResponse(accessToken, "")

	return utils.SuccessResponse(c, "Google authentication successful", response)
}

// GetConnectedProviders godoc
// @Summary Get connected OAuth providers
// @Description Get list of OAuth providers connected to user account
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]response.OAuthProviderResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/oauth/connected [get]
func (h *AuthHandler) GetConnectedProviders(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user from context
	claims := c.Locals("user")
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "No user context found")
	}

	userClaims := claims.(*utils.Claims)

	// Get connected providers (service expects uint)
	providers, err := h.oauthService.GetConnectedProviders(ctx, uint(userClaims.UserID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get connected providers", err.Error())
	}

	// Convert []auth.OAuthProvider to []*auth.OAuthProvider for mapper
	providerPointers := make([]*auth.OAuthProvider, len(providers))
	for i := range providers {
		providerPointers[i] = &providers[i]
	}

	// Convert to response
	response := mapper.ToOAuthProviderListResponse(providerPointers)

	return utils.SuccessResponse(c, "Connected providers retrieved successfully", response)
}

// DisconnectOAuth godoc
// @Summary Disconnect OAuth provider
// @Description Disconnect an OAuth provider from user account
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Param provider path string true "Provider name (google, facebook, github)"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/oauth/{provider} [delete]
func (h *AuthHandler) DisconnectOAuth(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user from context
	claims := c.Locals("user")
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "No user context found")
	}

	userClaims := claims.(*utils.Claims)

	// Get provider from path
	provider := c.Params("provider")
	if provider == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Provider name required", "")
	}

	// Disconnect provider (service expects uint)
	if err := h.oauthService.DisconnectOAuthProvider(ctx, uint(userClaims.UserID), provider); err != nil {
		// Check if provider not connected
		if err.Error() == "OAuth provider not connected" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Provider not connected", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to disconnect provider", err.Error())
	}

	return utils.SuccessResponse(c, "OAuth provider disconnected successfully", nil)
}
