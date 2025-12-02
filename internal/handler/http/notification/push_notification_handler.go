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

type PushNotificationHandler struct {
	pushService notification.PushNotificationService
	logger      *logrus.Logger
}

func NewPushNotificationHandler(
	pushService notification.PushNotificationService,
	logger *logrus.Logger,
) *PushNotificationHandler {
	return &PushNotificationHandler{
		pushService: pushService,
		logger:      logger,
	}
}

func (h *PushNotificationHandler) SendPushToDevice(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.SendPushToDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.WithError(err).Error("Failed to parse send push to device request")
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		h.logger.WithError(err).Error("Validation failed for send push to device")
		return utils.BadRequestResponse(c, err.Error())
	}

	message := mapper.ToPushMessage(&req.SendPushNotificationRequest)

	result, err := h.pushService.SendToDevice(ctx, req.Token, message)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send push notification to device")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send push notification", err.Error())
	}

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

func (h *PushNotificationHandler) SendPushToUser(c *fiber.Ctx) error {
	ctx := c.Context()
	senderID := middleware.GetUserID(c)

	targetUserID, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid user ID")
	}

	var req request.SendPushToUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.WithError(err).Error("Failed to parse send push to user request")
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		h.logger.WithError(err).Error("Validation failed for send push to user")
		return utils.BadRequestResponse(c, err.Error())
	}

	req.UserID = targetUserID

	message := mapper.ToPushMessage(&req.SendPushNotificationRequest)

	results, err := h.pushService.SendToUser(ctx, req.UserID, message)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send push notification to user")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send push notification", err.Error())
	}

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

func (h *PushNotificationHandler) SendPushToMultipleUsers(c *fiber.Ctx) error {
	ctx := c.Context()
	senderID := middleware.GetUserID(c)

	var req request.SendPushToMultipleUsersRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.WithError(err).Error("Failed to parse send push to multiple users request")
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		h.logger.WithError(err).Error("Validation failed for send push to multiple users")
		return utils.BadRequestResponse(c, err.Error())
	}

	message := mapper.ToPushMessage(&req.SendPushNotificationRequest)

	results, err := h.pushService.SendToMultipleUsers(ctx, req.UserIDs, message)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send push notifications to multiple users")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send push notifications", err.Error())
	}

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

func (h *PushNotificationHandler) SendPushToTopic(c *fiber.Ctx) error {
	ctx := c.Context()
	senderID := middleware.GetUserID(c)

	var req request.SendPushToTopicRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.WithError(err).Error("Failed to parse send push to topic request")
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		h.logger.WithError(err).Error("Validation failed for send push to topic")
		return utils.BadRequestResponse(c, err.Error())
	}

	message := mapper.ToPushMessage(&req.SendPushNotificationRequest)

	result, err := h.pushService.SendToTopic(ctx, req.Topic, message)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send push notification to topic")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send push notification", err.Error())
	}

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

func (h *PushNotificationHandler) SendTestNotification(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.TestPushNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.WithError(err).Error("Failed to parse test notification request")
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		h.logger.WithError(err).Error("Validation failed for test notification")
		return utils.BadRequestResponse(c, err.Error())
	}

	message := &notification.PushMessage{
		Title:    "Test Notification",
		Body:     "This is a test notification from Keerja Backend",
		Data:     map[string]string{"type": "test", "timestamp": strconv.FormatInt(c.Context().Time().Unix(), 10)},
		Priority: "high",
		Sound:    "default",
	}

	result, err := h.pushService.SendToDevice(ctx, req.Token, message)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send test notification")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send test notification", err.Error())
	}

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
