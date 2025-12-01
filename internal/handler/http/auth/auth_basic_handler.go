package authhandler

import (
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/service"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	req.FullName = utils.SanitizeString(req.FullName)
	req.Email = utils.SanitizeString(req.Email)
	if req.Phone != "" {
		req.Phone = utils.SanitizeString(req.Phone)
	}

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

	usr, verificationToken, err := h.authService.Register(ctx, domainReq)
	if err != nil {
		if err == service.ErrEmailAlreadyExists {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Email already exists", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to register user", err.Error())
	}

	response := mapper.ToAuthResponse(usr, "", verificationToken)

	return utils.CreatedResponse(c, "User registered successfully. Please check your email for verification.", response)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	req.Email = utils.SanitizeString(req.Email)

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

	authResponse := h.buildAuthResponse(ctx, usr, accessToken, "")

	return utils.SuccessResponse(c, "Login successful", authResponse)
}

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	req.Token = utils.SanitizeString(req.Token)

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

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", "")
	}

	accessToken, err := h.authService.RefreshToken(ctx, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to refresh token", err.Error())
	}

	usr, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get user info", err.Error())
	}

	authResponse := h.buildAuthResponse(ctx, usr, accessToken, "")

	return utils.SuccessResponse(c, "Token refreshed successfully", authResponse)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	if err := h.authService.Logout(ctx, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to logout", err.Error())
	}

	return utils.SuccessResponse(c, "Logout successful", nil)
}
