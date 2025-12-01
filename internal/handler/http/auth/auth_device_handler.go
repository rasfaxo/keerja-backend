package authhandler

import (
	"keerja-backend/internal/domain/auth"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/service"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) LoginWithRememberMe(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.LoginWithRememberMeRequest
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

	deviceInfo := service.DeviceInfo{
		DeviceID:  req.DeviceID,
		UserAgent: string(c.Request().Header.UserAgent()),
		IPAddress: c.IP(),
	}

	refreshToken, err := h.refreshTokenService.CreateRefreshToken(ctx, usr.ID, deviceInfo, req.RememberMe)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create refresh token", err.Error())
	}

	authResponse := h.buildAuthResponse(ctx, usr, accessToken, refreshToken)

	return utils.SuccessResponse(c, "Login successful", authResponse)
}

func (h *AuthHandler) RefreshAccessToken(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.RefreshAccessTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	claims := c.Locals("user")
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "No user context found")
	}

	userClaims := claims.(*utils.Claims)

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

	tokenResponse := mapper.ToTokenResponse(newAccessToken, newRefreshToken)

	return utils.SuccessResponse(c, "Token refreshed successfully", tokenResponse)
}

func (h *AuthHandler) GetActiveDevices(c *fiber.Ctx) error {
	ctx := c.Context()

	claims := c.Locals("user")
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "No user context found")
	}

	userClaims := claims.(*utils.Claims)

	devices, err := h.refreshTokenService.GetUserDevices(ctx, userClaims.UserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get devices", err.Error())
	}

	devicePointers := make([]*auth.RefreshToken, len(devices))
	for i := range devices {
		devicePointers[i] = &devices[i]
	}

	deviceList := mapper.ToDeviceListResponse(devicePointers)

	return utils.SuccessResponse(c, "Active devices retrieved successfully", deviceList)
}

func (h *AuthHandler) RevokeDevice(c *fiber.Ctx) error {
	ctx := c.Context()

	claims := c.Locals("user")
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "No user context found")
	}

	userClaims := claims.(*utils.Claims)

	var req request.RevokeDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	if err := h.refreshTokenService.RevokeDeviceToken(ctx, userClaims.UserID, req.DeviceID, "user_revoked"); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to revoke device", err.Error())
	}

	return utils.SuccessResponse(c, "Device session revoked successfully", nil)
}

func (h *AuthHandler) LogoutAllDevices(c *fiber.Ctx) error {
	ctx := c.Context()

	claims := c.Locals("user")
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "No user context found")
	}

	userClaims := claims.(*utils.Claims)

	if err := h.refreshTokenService.RevokeAllUserTokens(ctx, userClaims.UserID, "logout_all_devices"); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to logout from all devices", err.Error())
	}

	return utils.SuccessResponse(c, "Logged out from all devices successfully", nil)
}
