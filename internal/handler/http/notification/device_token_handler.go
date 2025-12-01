package notification

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

type DeviceTokenHandler struct {
	deviceTokenRepo notification.DeviceTokenRepository
	pushService     notification.PushNotificationService
	logger          *logrus.Logger
}

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

func (h *DeviceTokenHandler) RegisterDeviceToken(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.RegisterDeviceTokenRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.WithError(err).Error("Failed to parse register device token request")
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		h.logger.WithError(err).Error("Validation failed for register device token")
		return utils.BadRequestResponse(c, err.Error())
	}

	existingToken, err := h.deviceTokenRepo.FindByToken(ctx, req.Token)
	if err == nil && existingToken != nil {
		if existingToken.UserID == userID {
			if req.DeviceInfo != nil {
				updateReq := &request.UpdateDeviceTokenRequest{
					DeviceInfo: req.DeviceInfo,
				}
				mapper.UpdateDeviceTokenFromRequest(existingToken, updateReq)
			}

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

		existingToken.Deactivate()
		_ = h.deviceTokenRepo.Update(ctx, existingToken)
	}

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

func (h *DeviceTokenHandler) UnregisterDeviceToken(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)
	token := c.Params("token")

	if token == "" {
		return utils.BadRequestResponse(c, "Token parameter is required")
	}

	deviceToken, err := h.deviceTokenRepo.FindByToken(ctx, token)
	if err != nil {
		h.logger.WithError(err).Warn("Device token not found")
		return utils.NotFoundResponse(c, "Device token not found")
	}

	if deviceToken.UserID != userID {
		h.logger.WithFields(logrus.Fields{
			"user_id":       userID,
			"token_user_id": deviceToken.UserID,
		}).Warn("Unauthorized attempt to unregister device token")
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to delete this token", "")
	}

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

func (h *DeviceTokenHandler) GetUserDevices(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)

	var filterReq request.DeviceTokenFilterRequest
	if err := c.QueryParser(&filterReq); err != nil {
		h.logger.WithError(err).Warn("Failed to parse device token filter")
	}

	filterReq.UserID = &userID

	var tokens []notification.DeviceToken
	var err error

	if filterReq.Platform != nil {
		platform := notification.Platform(*filterReq.Platform)
		tokens, err = h.deviceTokenRepo.FindByUserAndPlatform(ctx, userID, platform)
	} else if filterReq.IsActive != nil && !*filterReq.IsActive {
		tokens, err = h.deviceTokenRepo.FindByUser(ctx, userID)
	} else {
		tokens, err = h.deviceTokenRepo.FindByUser(ctx, userID)
	}

	if err != nil {
		h.logger.WithError(err).Error("Failed to retrieve device tokens")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve device tokens", err.Error())
	}

	total := int64(len(tokens))
	response := mapper.ToDeviceTokenListResponse(tokens, total, page, pageSize)

	return utils.SuccessResponse(c, "Device tokens retrieved successfully", response)
}

func (h *DeviceTokenHandler) ValidateDeviceToken(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.ValidateDeviceTokenRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.WithError(err).Error("Failed to parse validate token request")
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		h.logger.WithError(err).Error("Validation failed for validate token")
		return utils.BadRequestResponse(c, err.Error())
	}

	deviceToken, err := h.deviceTokenRepo.FindByToken(ctx, req.Token)
	if err != nil {
		isValid, err := h.pushService.ValidateToken(ctx, req.Token)
		if err != nil {
			h.logger.WithError(err).Error("Failed to validate token with FCM")
			response := mapper.ToPushValidationResponse(req.Token, false, "", "Failed to validate token")
			return utils.SuccessResponse(c, "Token validation completed", response)
		}

		response := mapper.ToPushValidationResponse(req.Token, isValid, "", "Token not registered")
		return utils.SuccessResponse(c, "Token validation completed", response)
	}

	if deviceToken.UserID != userID {
		response := mapper.ToPushValidationResponse(req.Token, false, "", "Token belongs to different user")
		return utils.SuccessResponse(c, "Token validation completed", response)
	}

	isValid, err := h.pushService.ValidateToken(ctx, req.Token)
	if err != nil {
		h.logger.WithError(err).Error("Failed to validate token with FCM")
		response := mapper.ToPushValidationResponse(req.Token, false, string(deviceToken.Platform), "Failed to validate with FCM")
		return utils.SuccessResponse(c, "Token validation completed", response)
	}

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

func (h *DeviceTokenHandler) GetDeviceTokenStats(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	tokens, err := h.deviceTokenRepo.FindByUser(ctx, userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to retrieve device tokens")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve statistics", err.Error())
	}

	count := int64(len(tokens))

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
		int(count),
		int(count),
		0,
		androidCount,
		iosCount,
		webCount,
	)

	return utils.SuccessResponse(c, "Device token statistics retrieved successfully", response)
}

func (h *DeviceTokenHandler) GetDeviceTokenByID(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid device token ID")
	}

	deviceToken, err := h.deviceTokenRepo.FindByID(ctx, id)
	if err != nil {
		h.logger.WithError(err).Warn("Device token not found")
		return utils.NotFoundResponse(c, "Device token not found")
	}

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
