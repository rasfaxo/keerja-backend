package http

import (
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/service"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
	userRepo    user.UserRepository
}

func NewAuthHandler(authService *service.AuthService, userRepo user.UserRepository) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userRepo:    userRepo,
	}
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

	// Convert to response DTO
	response := mapper.ToAuthResponse(usr, accessToken, "")

	return utils.SuccessResponse(c, "Login successful", response)
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

	// Convert to response DTO
	response := mapper.ToAuthResponse(usr, accessToken, "")

	return utils.SuccessResponse(c, "Token refreshed successfully", response)
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
