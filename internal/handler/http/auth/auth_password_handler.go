package authhandler

import (
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/service"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	req.Email = utils.SanitizeString(req.Email)

	if err := h.authService.ForgotPassword(ctx, req.Email); err != nil {
		return utils.SuccessResponse(c, "If the email exists, a password reset link has been sent", nil)
	}

	return utils.SuccessResponse(c, "Password reset link sent to your email", nil)
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	req.Token = utils.SanitizeString(req.Token)

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

func (h *AuthHandler) ForgotPasswordOTP(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.ForgotPasswordOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	req.Email = utils.SanitizeString(req.Email)

	if err := h.registrationService.RequestPasswordResetOTP(ctx, req.Email); err != nil {
		if err == service.ErrTooManyOTPRequests {
			return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "Too many password reset requests. Please try again later.", err.Error())
		}
		if err == service.ErrEmailNotVerified {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email not verified", err.Error())
		}
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

func (h *AuthHandler) ResetPasswordOTP(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.ResetPasswordOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	req.Email = utils.SanitizeString(req.Email)
	req.OTPCode = utils.SanitizeString(req.OTPCode)

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
