package notification

import (
	"keerja-backend/internal/domain/notification"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	notifService notification.NotificationService
}

func NewNotificationHandler(notifService notification.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notifService: notifService,
	}
}

func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	var filterReq request.NotificationFilterRequest
	if err := c.QueryParser(&filterReq); err != nil {
		return utils.BadRequestResponse(c, "Invalid filter parameters")
	}

	filter := mapper.ToNotificationFilter(&filterReq)

	notifications, total, err := h.notifService.GetUserNotifications(ctx, userID, filter, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve notifications", err.Error())
	}

	response := mapper.ToNotificationListResponse(notifications, total, page, limit)
	return utils.SuccessResponse(c, "Notifications retrieved successfully", response)
}

func (h *NotificationHandler) GetUnreadNotifications(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	limit := c.QueryInt("limit", 50)

	notifications, err := h.notifService.GetUnreadNotifications(ctx, userID, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve unread notifications", err.Error())
	}

	responses := make([]interface{}, len(notifications))
	for i, notif := range notifications {
		resp := mapper.ToNotificationResponse(&notif)
		responses[i] = resp
	}

	return utils.SuccessResponse(c, "Unread notifications retrieved successfully", responses)
}

func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	count, err := h.notifService.GetUnreadCount(ctx, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve unread count", err.Error())
	}

	return utils.SuccessResponse(c, "Unread count retrieved successfully", fiber.Map{
		"unread_count": count,
	})
}

func (h *NotificationHandler) GetNotificationByID(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)
	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid notification ID")
	}

	notif, err := h.notifService.GetNotificationByID(ctx, id, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Notification not found", err.Error())
	}

	response := mapper.ToNotificationResponse(notif)
	return utils.SuccessResponse(c, "Notification retrieved successfully", response)
}

func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid notification ID")
	}

	if err := h.notifService.MarkAsRead(ctx, id, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to mark notification as read", err.Error())
	}

	return utils.SuccessResponse(c, "Notification marked as read", nil)
}

func (h *NotificationHandler) MarkAsUnread(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid notification ID")
	}

	if err := h.notifService.MarkAsUnread(ctx, id, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to mark notification as unread", err.Error())
	}

	return utils.SuccessResponse(c, "Notification marked as unread", nil)
}

func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	if err := h.notifService.MarkAllAsRead(ctx, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to mark all notifications as read", err.Error())
	}

	return utils.SuccessResponse(c, "All notifications marked as read", nil)
}

func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid notification ID")
	}

	if err := h.notifService.DeleteNotification(ctx, id, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete notification", err.Error())
	}

	return utils.SuccessResponse(c, "Notification deleted successfully", nil)
}

func (h *NotificationHandler) DeleteAllNotifications(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	if err := h.notifService.DeleteAllNotifications(ctx, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete notifications", err.Error())
	}

	return utils.SuccessResponse(c, "All notifications deleted successfully", nil)
}

func (h *NotificationHandler) GetNotificationStats(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	stats, err := h.notifService.GetNotificationStats(ctx, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve notification statistics", err.Error())
	}

	response := mapper.ToNotificationStatsResponse(stats)
	return utils.SuccessResponse(c, "Notification statistics retrieved successfully", response)
}

func (h *NotificationHandler) GetPreferences(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	prefs, err := h.notifService.GetNotificationPreferences(ctx, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve notification preferences", err.Error())
	}

	response := mapper.ToNotificationPreferenceResponse(prefs)
	return utils.SuccessResponse(c, "Notification preferences retrieved successfully", response)
}

func (h *NotificationHandler) UpdatePreferences(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.UpdateNotificationPreferencesRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	existingPrefs, _ := h.notifService.GetNotificationPreferences(ctx, userID)

	prefs := mapper.ToNotificationPreferenceFromRequest(&req, existingPrefs)
	if err := h.notifService.UpdateNotificationPreferences(ctx, userID, prefs); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update notification preferences", err.Error())
	}

	updatedPrefs, _ := h.notifService.GetNotificationPreferences(ctx, userID)
	response := mapper.ToNotificationPreferenceResponse(updatedPrefs)

	return utils.SuccessResponse(c, "Notification preferences updated successfully", response)
}

func (h *NotificationHandler) SendNotification(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.SendNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	req.Title = utils.SanitizeString(req.Title)
	req.Message = utils.SanitizeString(req.Message)

	domainReq := mapper.ToSendNotificationRequest(&req)
	notif, err := h.notifService.SendNotification(ctx, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to send notification", err.Error())
	}

	response := mapper.ToNotificationResponse(notif)
	return utils.CreatedResponse(c, "Notification sent successfully", response)
}

func (h *NotificationHandler) SendBulkNotification(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.SendBulkNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	req.Title = utils.SanitizeString(req.Title)
	req.Message = utils.SanitizeString(req.Message)

	domainReq := &notification.SendNotificationRequest{
		Type:        req.Type,
		Title:       req.Title,
		Message:     req.Message,
		Data:        req.Data,
		Priority:    req.Priority,
		Category:    req.Category,
		ActionURL:   req.ActionURL,
		Icon:        req.Icon,
		SenderID:    req.SenderID,
		RelatedID:   req.RelatedID,
		RelatedType: req.RelatedType,
		ExpiresAt:   req.ExpiresAt,
		Channel:     req.Channel,
	}

	if err := h.notifService.SendBulkNotification(ctx, req.UserIDs, domainReq); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to send bulk notifications", err.Error())
	}

	return utils.CreatedResponse(c, "Bulk notifications sent successfully", nil)
}
