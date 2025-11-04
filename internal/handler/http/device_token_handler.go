package http

import (
	"strconv"

	"keerja-backend/internal/domain/notification"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// DeviceTokenHandler handles device token-related HTTP requests
type DeviceTokenHandler struct {
	deviceTokenRepo notification.DeviceTokenRepository
	pushService     notification.PushNotificationService
	logger          *logrus.Logger
}

// NewDeviceTokenHandler creates a new device token handler
func NewDeviceTokenHandler(
	deviceTokenRepo notification.DeviceTokenRepository,
	pushService notification.PushNotificationService,
	logger *logrus.Logger,
) *DeviceTokenHandler {
	return &DeviceTokenHandler{
		deviceTokenRepo: deviceTokenRepo,
		pushService:     pushService,
		logger:          logger,
	}
}

// RegisterDeviceToken godoc
// @Summary Register device token
// @Description Register a device token for push notifications
// @Tags Device Tokens
// @Accept json
// @Produce json
// @Param request body request.RegisterDeviceTokenRequest true "Device token registration request"
// @Success 200 {object} utils.Response "Device token registered successfully"
// @Failure 400 {object} utils.Response "Invalid request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/device-tokens [post]
func (h *DeviceTokenHandler) RegisterDeviceToken(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	// Parse request body
	var req request.RegisterDeviceTokenRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.WithError(err).Error("Failed to parse register device token request")
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		h.logger.WithError(err).Error("Validation failed for register device token")
		return utils.BadRequestResponse(c, err.Error())
	}

	// Check if token already exists for this user
	existingToken, err := h.deviceTokenRepo.FindByToken(ctx, req.Token)
	if err == nil && existingToken != nil {
		// Token exists, check if it belongs to this user
		if existingToken.UserID == userID {
			// Update existing token
			if req.DeviceInfo != nil {
				updateReq := &request.UpdateDeviceTokenRequest{
					DeviceInfo: req.DeviceInfo,
				}
				mapper.UpdateDeviceTokenFromRequest(existingToken, updateReq)
			}

			// Reactivate if inactive
			if !existingToken.IsActive {
				existingToken.Activate()
			}

			existingToken.MarkAsUsed()

			if err := h.deviceTokenRepo.Update(ctx, existingToken); err != nil {
				h.logger.WithError(err).Error("Failed to update existing device token")
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update device token", err.Error())
			}

			response := mapper.ToDeviceTokenResponse(existingToken)
			h.logger.WithField("user_id", userID).Info("Device token updated successfully")
			return utils.SuccessResponse(c, "Device token updated successfully", response)
		}

		// Token belongs to different user - this shouldn't happen normally
		// Deactivate the old token and create new one
		existingToken.Deactivate()
		_ = h.deviceTokenRepo.Update(ctx, existingToken)
	}

	// Create new device token
	deviceToken := mapper.ToDeviceToken(userID, &req)
	if err := h.deviceTokenRepo.Create(ctx, deviceToken); err != nil {
		h.logger.WithError(err).Error("Failed to create device token")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to register device token", err.Error())
	}

	response := mapper.ToDeviceTokenResponse(deviceToken)
	h.logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"platform": req.Platform,
	}).Info("Device token registered successfully")

	return utils.CreatedResponse(c, "Device token registered successfully", response)
}

// UnregisterDeviceToken godoc
// @Summary Unregister device token
// @Description Remove a device token from the system
// @Tags Device Tokens
// @Accept json
// @Produce json
// @Param token path string true "Device token"
// @Success 200 {object} utils.Response "Device token unregistered successfully"
// @Failure 400 {object} utils.Response "Invalid request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Token does not belong to user"
// @Failure 404 {object} utils.Response "Device token not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/device-tokens/{token} [delete]
func (h *DeviceTokenHandler) UnregisterDeviceToken(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)
	token := c.Params("token")

	if token == "" {
		return utils.BadRequestResponse(c, "Token parameter is required")
	}

	// Find token
	deviceToken, err := h.deviceTokenRepo.FindByToken(ctx, token)
	if err != nil {
		h.logger.WithError(err).Warn("Device token not found")
		return utils.NotFoundResponse(c, "Device token not found")
	}

	// Check ownership
	if deviceToken.UserID != userID {
		h.logger.WithFields(logrus.Fields{
			"user_id":       userID,
			"token_user_id": deviceToken.UserID,
		}).Warn("Unauthorized attempt to unregister device token")
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to delete this token", "")
	}

	// Delete token
	if err := h.deviceTokenRepo.DeleteByToken(ctx, token); err != nil {
		h.logger.WithError(err).Error("Failed to delete device token")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to unregister device token", err.Error())
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"token":   token,
	}).Info("Device token unregistered successfully")

	return utils.SuccessResponse(c, "Device token unregistered successfully", nil)
}

// GetUserDevices godoc
// @Summary Get user device tokens
// @Description Get all device tokens registered for the authenticated user
// @Tags Device Tokens
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param platform query string false "Filter by platform" Enums(android, ios, web)
// @Param is_active query boolean false "Filter by active status"
// @Success 200 {object} utils.Response "Device tokens retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/device-tokens [get]
func (h *DeviceTokenHandler) GetUserDevices(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	// Parse query parameters
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)

	// Parse filters
	var filterReq request.DeviceTokenFilterRequest
	if err := c.QueryParser(&filterReq); err != nil {
		h.logger.WithError(err).Warn("Failed to parse device token filter")
	}

	// Override user_id from authentication
	filterReq.UserID = &userID

	// Get tokens
	var tokens []notification.DeviceToken
	var err error

	if filterReq.Platform != nil {
		// Filter by platform
		platform := notification.Platform(*filterReq.Platform)
		tokens, err = h.deviceTokenRepo.FindByUserAndPlatform(ctx, userID, platform)
	} else if filterReq.IsActive != nil && !*filterReq.IsActive {
		// Get all tokens (including inactive) - no specific method for this
		// We'll use FindByUser which returns active tokens
		tokens, err = h.deviceTokenRepo.FindByUser(ctx, userID)
	} else {
		// Get active tokens (default)
		tokens, err = h.deviceTokenRepo.FindByUser(ctx, userID)
	}

	if err != nil {
		h.logger.WithError(err).Error("Failed to retrieve device tokens")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve device tokens", err.Error())
	}

	// Convert to response
	total := int64(len(tokens))
	response := mapper.ToDeviceTokenListResponse(tokens, total, page, pageSize)

	return utils.SuccessResponse(c, "Device tokens retrieved successfully", response)
}

// ValidateDeviceToken godoc
// @Summary Validate device token
// @Description Validate if a device token is still valid with FCM
// @Tags Device Tokens
// @Accept json
// @Produce json
// @Param request body request.ValidateDeviceTokenRequest true "Validation request"
// @Success 200 {object} utils.Response "Token validation result"
// @Failure 400 {object} utils.Response "Invalid request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/device-tokens/validate [post]
func (h *DeviceTokenHandler) ValidateDeviceToken(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	// Parse request body
	var req request.ValidateDeviceTokenRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.WithError(err).Error("Failed to parse validate token request")
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		h.logger.WithError(err).Error("Validation failed for validate token")
		return utils.BadRequestResponse(c, err.Error())
	}

	// Find token in database
	deviceToken, err := h.deviceTokenRepo.FindByToken(ctx, req.Token)
	if err != nil {
		// Token not in database - try to validate with FCM
		isValid, err := h.pushService.ValidateToken(ctx, req.Token)
		if err != nil {
			h.logger.WithError(err).Error("Failed to validate token with FCM")
			response := mapper.ToPushValidationResponse(req.Token, false, "", "Failed to validate token")
			return utils.SuccessResponse(c, "Token validation completed", response)
		}

		response := mapper.ToPushValidationResponse(req.Token, isValid, "", "Token not registered")
		return utils.SuccessResponse(c, "Token validation completed", response)
	}

	// Check ownership
	if deviceToken.UserID != userID {
		response := mapper.ToPushValidationResponse(req.Token, false, "", "Token belongs to different user")
		return utils.SuccessResponse(c, "Token validation completed", response)
	}

	// Validate with FCM
	isValid, err := h.pushService.ValidateToken(ctx, req.Token)
	if err != nil {
		h.logger.WithError(err).Error("Failed to validate token with FCM")
		response := mapper.ToPushValidationResponse(req.Token, false, string(deviceToken.Platform), "Failed to validate with FCM")
		return utils.SuccessResponse(c, "Token validation completed", response)
	}

	// Update token status if invalid
	if !isValid && deviceToken.IsActive {
		deviceToken.RecordFailure("Token invalid per FCM validation")
		_ = h.deviceTokenRepo.Update(ctx, deviceToken)
	}

	message := "Token is valid"
	if !isValid {
		message = "Token is invalid"
	}

	response := mapper.ToPushValidationResponse(req.Token, isValid, string(deviceToken.Platform), message)
	return utils.SuccessResponse(c, "Token validation completed", response)
}

// GetDeviceTokenStats godoc
// @Summary Get device token statistics
// @Description Get statistics of user's device tokens by platform
// @Tags Device Tokens
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "Device token statistics"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/device-tokens/stats [get]
func (h *DeviceTokenHandler) GetDeviceTokenStats(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	// Get all user's tokens to calculate stats
	tokens, err := h.deviceTokenRepo.FindByUser(ctx, userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to retrieve device tokens")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve statistics", err.Error())
	}

	count := int64(len(tokens))

	// Calculate stats by platform
	var androidCount, iosCount, webCount int
	for _, token := range tokens {
		switch token.Platform {
		case notification.PlatformAndroid:
			androidCount++
		case notification.PlatformIOS:
			iosCount++
		case notification.PlatformWeb:
			webCount++
		}
	}

	response := mapper.ToDeviceTokenStatsResponse(
		int(count),   // total
		int(count),   // active (FindByUser returns only active)
		0,            // inactive
		androidCount, // android
		iosCount,     // ios
		webCount,     // web
	)

	return utils.SuccessResponse(c, "Device token statistics retrieved successfully", response)
}

// GetDeviceTokenByID godoc
// @Summary Get device token by ID
// @Description Get a specific device token by its ID
// @Tags Device Tokens
// @Accept json
// @Produce json
// @Param id path int true "Device token ID"
// @Success 200 {object} utils.Response "Device token retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Token does not belong to user"
// @Failure 404 {object} utils.Response "Device token not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/device-tokens/{id} [get]
func (h *DeviceTokenHandler) GetDeviceTokenByID(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	// Parse ID
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid device token ID")
	}

	// Find token
	deviceToken, err := h.deviceTokenRepo.FindByID(ctx, id)
	if err != nil {
		h.logger.WithError(err).Warn("Device token not found")
		return utils.NotFoundResponse(c, "Device token not found")
	}

	// Check ownership
	if deviceToken.UserID != userID {
		h.logger.WithFields(logrus.Fields{
			"user_id":       userID,
			"token_user_id": deviceToken.UserID,
		}).Warn("Unauthorized attempt to access device token")
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to access this token", "")
	}

	response := mapper.ToDeviceTokenResponse(deviceToken)
	return utils.SuccessResponse(c, "Device token retrieved successfully", response)
}
