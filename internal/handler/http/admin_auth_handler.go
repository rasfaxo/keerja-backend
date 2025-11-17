package http

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

// Login godoc
// @Summary Admin login
// @Description Authenticate admin user and get access token
// @Tags admin-auth
// @Accept json
// @Produce json
// @Param request body request.AdminLoginRequest true "Admin login request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/admin/login [post]
func (h *AdminAuthHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.AdminLoginRequest
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

	// Login admin
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

	// Convert to response DTO
	expiresIn := int64(3600) // 1 hour in seconds
	response := mapper.ToAdminAuthResponse(adminUser, accessToken, refreshToken, expiresIn)

	return utils.SuccessResponse(c, "Admin login successful", response)
}

// Logout godoc
// @Summary Admin logout
// @Description Logout admin and invalidate session
// @Tags admin-auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/admin/logout [post]
func (h *AdminAuthHandler) Logout(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get admin ID from context (set by admin auth middleware)
	adminID := middleware.GetAdminID(c)
	if adminID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", "")
	}

	// Logout admin
	if err := h.adminAuthService.Logout(ctx, adminID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to logout", err.Error())
	}

	return utils.SuccessResponse(c, "Logout successful", nil)
}

// RefreshToken godoc
// @Summary Refresh admin access token
// @Description Get new access token using refresh token
// @Tags admin-auth
// @Accept json
// @Produce json
// @Param request body request.AdminRefreshTokenRequest true "Refresh token request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/admin/refresh-token [post]
func (h *AdminAuthHandler) RefreshToken(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.AdminRefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Refresh token
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

	// Convert to response
	expiresIn := int64(3600) // 1 hour in seconds
	response := mapper.ToAdminTokenResponse(newAccessToken, newRefreshToken, expiresIn)

	return utils.SuccessResponse(c, "Token refreshed successfully", response)
}

// GetProfile godoc
// @Summary Get current admin profile
// @Description Get authenticated admin user profile
// @Tags admin-auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/admin/me [get]
func (h *AdminAuthHandler) GetProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get admin ID from context
	adminID := middleware.GetAdminID(c)
	if adminID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", "")
	}

	// Get admin profile
	adminUser, err := h.adminAuthService.GetCurrentProfile(ctx, adminID)
	if err != nil {
		if err == service.ErrAdminNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Admin not found", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get profile", err.Error())
	}

	// Convert to response
	response := mapper.ToAdminProfileResponse(adminUser)

	return utils.SuccessResponse(c, "Profile retrieved successfully", response)
}

// ChangePassword godoc
// @Summary Change admin password
// @Description Change current admin user password
// @Tags admin-auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.AdminChangePasswordRequest true "Change password request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/admin/change-password [put]
func (h *AdminAuthHandler) ChangePassword(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get admin ID from context
	adminID := middleware.GetAdminID(c)
	if adminID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", "")
	}

	// Parse request body
	var req request.AdminChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Change password
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
