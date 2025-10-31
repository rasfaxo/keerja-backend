package routes

import (
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupNotificationRoutes configures notification routes
// Routes: /api/v1/notifications/*
//
// User Endpoints (12):
//   - GET    /                         Get user notifications
//   - GET    /unread                   Get unread notifications
//   - GET    /unread/count             Get unread notification count
//   - GET    /:id                      Get notification by ID
//   - PUT    /:id/read                 Mark notification as read
//   - PUT    /:id/unread               Mark notification as unread
//   - PUT    /read-all                 Mark all notifications as read
//   - DELETE /:id                      Delete notification
//   - DELETE /                         Delete all notifications
//   - GET    /stats                    Get notification statistics
//   - GET    /preferences              Get notification preferences
//   - PUT    /preferences              Update notification preferences
//
// Admin Endpoints (2):
//   - POST   /send                     Send notification (admin only)
//   - POST   /send-bulk                Send bulk notifications (admin only)
//
// Total: 14 endpoints
func SetupNotificationRoutes(api fiber.Router, deps *Dependencies, authMw *middleware.AuthMiddleware) {
	// ============================================
	// USER ROUTES (12 endpoints)
	// All require authentication
	// ============================================
	notifications := api.Group("/notifications")
	notifications.Use(authMw.AuthRequired())

	// GET /api/v1/notifications - Get user notifications
	// Query params: page, limit, type, category, is_read, priority, date_from, date_to
	// Rate limit: 30 requests/minute
	notifications.Get("/",
		middleware.SearchRateLimiter(),
		deps.NotificationHandler.GetNotifications,
	)

	// GET /api/v1/notifications/unread - Get unread notifications
	// Query params: limit (default: 50)
	// Rate limit: 30 requests/minute
	notifications.Get("/unread",
		middleware.SearchRateLimiter(),
		deps.NotificationHandler.GetUnreadNotifications,
	)

	// GET /api/v1/notifications/unread/count - Get unread notification count
	// Returns: { unread_count: number }
	// Rate limit: 30 requests/minute
	notifications.Get("/unread/count",
		middleware.SearchRateLimiter(),
		deps.NotificationHandler.GetUnreadCount,
	)

	// GET /api/v1/notifications/stats - Get notification statistics
	// Returns: total, unread, read, today, this week, high priority, category breakdown
	// Rate limit: 30 requests/minute
	notifications.Get("/stats",
		middleware.SearchRateLimiter(),
		deps.NotificationHandler.GetNotificationStats,
	)

	// GET /api/v1/notifications/preferences - Get notification preferences
	// Returns: email, push, sms preferences, notification type preferences
	// Rate limit: 30 requests/minute
	notifications.Get("/preferences",
		middleware.SearchRateLimiter(),
		deps.NotificationHandler.GetPreferences,
	)

	// PUT /api/v1/notifications/preferences - Update notification preferences
	// Body: { email_enabled, push_enabled, job_applications_enabled, ... }
	// Rate limit: 100 requests/minute
	notifications.Put("/preferences",
		middleware.ApplicationRateLimiter(),
		deps.NotificationHandler.UpdatePreferences,
	)

	// PUT /api/v1/notifications/read-all - Mark all notifications as read
	// Rate limit: 100 requests/minute
	notifications.Put("/read-all",
		middleware.ApplicationRateLimiter(),
		deps.NotificationHandler.MarkAllAsRead,
	)

	// DELETE /api/v1/notifications - Delete all notifications
	// Rate limit: 100 requests/minute
	notifications.Delete("/",
		middleware.ApplicationRateLimiter(),
		deps.NotificationHandler.DeleteAllNotifications,
	)

	// GET /api/v1/notifications/:id - Get notification by ID
	// Returns: Full notification details
	notifications.Get("/:id",
		deps.NotificationHandler.GetNotificationByID,
	)

	// PUT /api/v1/notifications/:id/read - Mark notification as read
	// Rate limit: 100 requests/minute
	notifications.Put("/:id/read",
		middleware.ApplicationRateLimiter(),
		deps.NotificationHandler.MarkAsRead,
	)

	// PUT /api/v1/notifications/:id/unread - Mark notification as unread
	// Rate limit: 100 requests/minute
	notifications.Put("/:id/unread",
		middleware.ApplicationRateLimiter(),
		deps.NotificationHandler.MarkAsUnread,
	)

	// DELETE /api/v1/notifications/:id - Delete notification
	// Rate limit: 100 requests/minute
	notifications.Delete("/:id",
		middleware.ApplicationRateLimiter(),
		deps.NotificationHandler.DeleteNotification,
	)

	// ============================================
	// ADMIN ROUTES (2 endpoints)
	// Require authentication + admin role
	// ============================================
	admin := notifications.Group("")
	admin.Use(authMw.AdminOnly())

	// POST /api/v1/notifications/send - Send notification (admin only)
	// Body: { user_id, type, title, message, data, priority, category, ... }
	// Rate limit: 100 requests/minute
	admin.Post("/send",
		middleware.ApplicationRateLimiter(),
		deps.NotificationHandler.SendNotification,
	)

	// POST /api/v1/notifications/send-bulk - Send bulk notifications (admin only)
	// Body: { user_ids: [], type, title, message, data, priority, category, ... }
	// Rate limit: 30 requests/minute (bulk operation)
	admin.Post("/send-bulk",
		middleware.SearchRateLimiter(),
		deps.NotificationHandler.SendBulkNotification,
	)
}
