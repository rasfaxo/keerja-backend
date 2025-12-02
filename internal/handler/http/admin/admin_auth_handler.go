package admin

import (
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/service"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AdminAuthHandler struct {
	adminAuthService *service.AdminAuthService
}

func NewAdminAuthHandler(adminAuthService *service.AdminAuthService) *AdminAuthHandler {
	return &AdminAuthHandler{
		adminAuthService: adminAuthService,
	}
}

func (h *AdminAuthHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.AdminLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	req.Email = utils.SanitizeString(req.Email)

	adminUser, accessToken, refreshToken, err := h.adminAuthService.Login(ctx, req.Email, req.Password)
	if err != nil {
		if err == service.ErrAdminInvalidCredentials {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password", err.Error())
		}
		if err == service.ErrAdminNotActive {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Admin account is not active", err.Error())
		}
		if err == service.ErrAdminSuspended {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Admin account is suspended", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to login", err.Error())
	}

	expiresIn := int64(3600)
	response := mapper.ToAdminAuthResponse(adminUser, accessToken, refreshToken, expiresIn)

	return utils.SuccessResponse(c, "Admin login successful", response)
}

func (h *AdminAuthHandler) Logout(c *fiber.Ctx) error {
	ctx := c.Context()

	adminID := middleware.GetAdminID(c)
	if adminID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", "")
	}

	if err := h.adminAuthService.Logout(ctx, adminID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to logout", err.Error())
	}

	return utils.SuccessResponse(c, "Logout successful", nil)
}

func (h *AdminAuthHandler) RefreshToken(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.AdminRefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	newAccessToken, newRefreshToken, err := h.adminAuthService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if err == service.ErrAdminTokenExpired {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Refresh token expired", err.Error())
		}
		if err == service.ErrAdminInvalidToken {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid refresh token", err.Error())
		}
		if err == service.ErrAdminNotActive {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Admin account is not active", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to refresh token", err.Error())
	}

	expiresIn := int64(3600)
	response := mapper.ToAdminTokenResponse(newAccessToken, newRefreshToken, expiresIn)

	return utils.SuccessResponse(c, "Token refreshed successfully", response)
}

func (h *AdminAuthHandler) GetProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	adminID := middleware.GetAdminID(c)
	if adminID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", "")
	}

	adminUser, err := h.adminAuthService.GetCurrentProfile(ctx, adminID)
	if err != nil {
		if err == service.ErrAdminNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Admin not found", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get profile", err.Error())
	}

	response := mapper.ToAdminProfileResponse(adminUser)

	return utils.SuccessResponse(c, "Profile retrieved successfully", response)
}

func (h *AdminAuthHandler) ChangePassword(c *fiber.Ctx) error {
	ctx := c.Context()

	adminID := middleware.GetAdminID(c)
	if adminID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", "")
	}

	var req request.AdminChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	if err := h.adminAuthService.ChangePassword(ctx, adminID, req.CurrentPassword, req.NewPassword); err != nil {
		if err == service.ErrAdminInvalidPassword {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Current password is incorrect", err.Error())
		}
		if err == service.ErrAdminNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Admin not found", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to change password", err.Error())
	}

	return utils.SuccessResponse(c, "Password changed successfully", nil)
}
