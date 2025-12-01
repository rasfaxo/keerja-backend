package authhandler

import (
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/service"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) RegisterWithOTP(c *fiber.Ctx) error {
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

func (h *AuthHandler) VerifyEmailOTP(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.VerifyEmailOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	req.Email = utils.SanitizeString(req.Email)
	req.OTPCode = utils.SanitizeString(req.OTPCode)

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

	authResponse := h.buildAuthResponse(ctx, usr, accessToken, "")

	return utils.SuccessResponse(c, "Email verified successfully", authResponse)
}

func (h *AuthHandler) ResendOTP(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.ResendOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	req.Email = utils.SanitizeString(req.Email)

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
		return utils.SuccessResponse(c, "If the email exists and is not verified, a new OTP has been sent", nil)
	}

	return utils.SuccessResponse(c, "OTP code has been resent to your email.", fiber.Map{
		"email": req.Email,
		"note":  "OTP code is valid for 5 minutes.",
	})
}

func (h *AuthHandler) ResendVerification(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.ResendVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	req.Email = utils.SanitizeString(req.Email)

	if err := h.authService.ResendVerificationEmail(ctx, req.Email); err != nil {
		return utils.SuccessResponse(c, "If the email exists and is not verified, a verification link has been sent", nil)
	}

	return utils.SuccessResponse(c, "Verification email sent successfully", nil)
}
