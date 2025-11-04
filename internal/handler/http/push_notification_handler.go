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

// PushNotificationHandler handles push notification-related HTTP requests
type PushNotificationHandler struct {
	pushService notification.PushNotificationService
	logger      *logrus.Logger
}

// NewPushNotificationHandler creates a new push notification handler
func NewPushNotificationHandler(
	pushService notification.PushNotificationService,
	logger *logrus.Logger,
) *PushNotificationHandler {
	return &PushNotificationHandler{
		pushService: pushService,
		logger:      logger,
	}
}

// SendPushToDevice godoc
// @Summary Send push notification to device
// @Description Send a push notification to a specific device token
// @Tags Push Notifications
// @Accept json
// @Produce json
// @Param request body request.SendPushToDeviceRequest true "Push notification request"
// @Success 200 {object} utils.Response "Push notification sent successfully"
// @Failure 400 {object} utils.Response "Invalid request or failed to send"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/push/send/device [post]
func (h *PushNotificationHandler) SendPushToDevice(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	// Parse request body
	var req request.SendPushToDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.WithError(err).Error("Failed to parse send push to device request")
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		h.logger.WithError(err).Error("Validation failed for send push to device")
		return utils.BadRequestResponse(c, err.Error())
	}

	// Convert to PushMessage
	message := mapper.ToPushMessage(&req.SendPushNotificationRequest)

	// Send notification
	result, err := h.pushService.SendToDevice(ctx, req.Token, message)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send push notification to device")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send push notification", err.Error())
	}

	// Convert to response
	response := mapper.ToPushNotificationResponse(result)

	h.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"success": result.Success,
		"token":   req.Token,
	}).Info("Push notification sent to device")

	if result.Success {
		return utils.SuccessResponse(c, "Push notification sent successfully", response)
	} else {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to send push notification", result.ErrorMessage)
	}
}

// SendPushToUser godoc
// @Summary Send push notification to user
// @Description Send a push notification to all devices of a specific user
// @Tags Push Notifications
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body request.SendPushToUserRequest true "Push notification request"
// @Success 200 {object} utils.Response "Push notification sent to user"
// @Failure 400 {object} utils.Response "Invalid request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/push/send/user/{id} [post]
func (h *PushNotificationHandler) SendPushToUser(c *fiber.Ctx) error {
	ctx := c.Context()
	senderID := middleware.GetUserID(c)

	// Parse user ID from URL
	targetUserID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid user ID")
	}

	// Parse request body
	var req request.SendPushToUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.WithError(err).Error("Failed to parse send push to user request")
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		h.logger.WithError(err).Error("Validation failed for send push to user")
		return utils.BadRequestResponse(c, err.Error())
	}

	// Override user ID from URL
	req.UserID = targetUserID

	// Convert to PushMessage
	message := mapper.ToPushMessage(&req.SendPushNotificationRequest)

	// Send notification to all user's devices
	results, err := h.pushService.SendToUser(ctx, req.UserID, message)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send push notification to user")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send push notification", err.Error())
	}

	// Convert to response
	response := mapper.ToBatchPushNotificationResponse(results, 1)

	h.logger.WithFields(logrus.Fields{
		"sender_id":   senderID,
		"target_user": targetUserID,
		"total_sent":  response.TotalSent,
		"success":     response.SuccessCount,
		"failed":      response.FailureCount,
	}).Info("Push notification sent to user")

	return utils.SuccessResponse(c, "Push notification sent to user", response)
}

// SendPushToMultipleUsers godoc
// @Summary Send push notification to multiple users
// @Description Send a push notification to multiple users (batch send)
// @Tags Push Notifications
// @Accept json
// @Produce json
// @Param request body request.SendPushToMultipleUsersRequest true "Batch push notification request"
// @Success 200 {object} utils.Response "Push notifications sent successfully"
// @Failure 400 {object} utils.Response "Invalid request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/push/send/batch [post]
func (h *PushNotificationHandler) SendPushToMultipleUsers(c *fiber.Ctx) error {
	ctx := c.Context()
	senderID := middleware.GetUserID(c)

	// Parse request body
	var req request.SendPushToMultipleUsersRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.WithError(err).Error("Failed to parse send push to multiple users request")
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		h.logger.WithError(err).Error("Validation failed for send push to multiple users")
		return utils.BadRequestResponse(c, err.Error())
	}

	// Convert to PushMessage
	message := mapper.ToPushMessage(&req.SendPushNotificationRequest)

	// Send notification to multiple users
	results, err := h.pushService.SendToMultipleUsers(ctx, req.UserIDs, message)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send push notifications to multiple users")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send push notifications", err.Error())
	}

	// Convert to response
	response := mapper.ToBatchPushNotificationResponseFromMap(results, len(req.UserIDs))

	h.logger.WithFields(logrus.Fields{
		"sender_id":  senderID,
		"user_count": len(req.UserIDs),
		"total_sent": response.TotalSent,
		"success":    response.SuccessCount,
		"failed":     response.FailureCount,
	}).Info("Push notifications sent to multiple users")

	return utils.SuccessResponse(c, "Push notifications sent successfully", response)
}

// SendPushToTopic godoc
// @Summary Send push notification to topic
// @Description Send a push notification to all subscribers of a topic
// @Tags Push Notifications
// @Accept json
// @Produce json
// @Param request body request.SendPushToTopicRequest true "Topic push notification request"
// @Success 200 {object} utils.Response "Push notification sent to topic successfully"
// @Failure 400 {object} utils.Response "Invalid request or failed to send"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/push/send/topic [post]
func (h *PushNotificationHandler) SendPushToTopic(c *fiber.Ctx) error {
	ctx := c.Context()
	senderID := middleware.GetUserID(c)

	// Parse request body
	var req request.SendPushToTopicRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.WithError(err).Error("Failed to parse send push to topic request")
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		h.logger.WithError(err).Error("Validation failed for send push to topic")
		return utils.BadRequestResponse(c, err.Error())
	}

	// Convert to PushMessage
	message := mapper.ToPushMessage(&req.SendPushNotificationRequest)

	// Send notification to topic
	result, err := h.pushService.SendToTopic(ctx, req.Topic, message)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send push notification to topic")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send push notification", err.Error())
	}

	// Convert to response
	response := mapper.ToPushNotificationResponse(result)

	h.logger.WithFields(logrus.Fields{
		"sender_id": senderID,
		"topic":     req.Topic,
		"success":   result.Success,
	}).Info("Push notification sent to topic")

	if result.Success {
		return utils.SuccessResponse(c, "Push notification sent to topic successfully", response)
	} else {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to send push notification to topic", result.ErrorMessage)
	}
}

// SendTestNotification godoc
// @Summary Send test notification
// @Description Send a test push notification to verify FCM integration
// @Tags Push Notifications
// @Accept json
// @Produce json
// @Param request body request.TestPushNotificationRequest true "Test notification request"
// @Success 200 {object} utils.Response "Test notification sent successfully"
// @Failure 400 {object} utils.Response "Invalid request or failed to send"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/push/test [post]
func (h *PushNotificationHandler) SendTestNotification(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	// Parse request body
	var req request.TestPushNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.WithError(err).Error("Failed to parse test notification request")
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		h.logger.WithError(err).Error("Validation failed for test notification")
		return utils.BadRequestResponse(c, err.Error())
	}

	// Create test message
	message := &notification.PushMessage{
		Title:    "Test Notification",
		Body:     "This is a test notification from Keerja Backend",
		Data:     map[string]string{"type": "test", "timestamp": strconv.FormatInt(c.Context().Time().Unix(), 10)},
		Priority: "high",
		Sound:    "default",
	}

	// Send notification
	result, err := h.pushService.SendToDevice(ctx, req.Token, message)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send test notification")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send test notification", err.Error())
	}

	// Convert to response
	response := mapper.ToPushNotificationResponse(result)

	h.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"success": result.Success,
		"token":   req.Token,
	}).Info("Test push notification sent")

	if result.Success {
		return utils.SuccessResponse(c, "Test notification sent successfully", response)
	} else {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to send test notification", result.ErrorMessage)
	}
}
