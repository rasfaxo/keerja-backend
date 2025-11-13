package routes

import (
	"keerja-backend/internal/handler/http"
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupPushNotificationRoutes configures all push notification related routes
// Routes: /api/v1/push/*
//
// Endpoints (5):
//   - POST   /send/device           Send to specific device token
//   - POST   /send/user/:id         Send to all user's devices
//   - POST   /send/batch            Send to multiple users (batch)
//   - POST   /send/topic            Send to topic subscribers
//   - POST   /test                  Send test notification
//
// Total: 5 endpoints
func SetupPushNotificationRoutes(api fiber.Router, handler *http.PushNotificationHandler, authMw *middleware.AuthMiddleware) {
	// Push Notification routes group
	push := api.Group("/push")

	// Apply authentication middleware to all push notification routes
	push.Use(authMw.AuthRequired())

	// Push Notification Sending Endpoints
	pushSend := push.Group("/send")

	// POST /api/v1/push/send/device - Send to specific device token
	// Body: { token, title, body, data, image_url, sound, priority, badge }
	// Rate limit: 100 requests/minute
	pushSend.Post("/device",
		middleware.ApplicationRateLimiter(),
		handler.SendPushToDevice,
	)

	// POST /api/v1/push/send/user/:id - Send to all user's devices
	// Body: { title, body, data, image_url, sound, priority, badge }
	// Rate limit: 100 requests/minute
	pushSend.Post("/user/:id",
		middleware.ApplicationRateLimiter(),
		handler.SendPushToUser,
	)

	// POST /api/v1/push/send/batch - Send to multiple users (batch)
	// Body: { user_ids, title, body, data, image_url, sound, priority, badge }
	// Rate limit: 30 requests/minute (bulk operation)
	pushSend.Post("/batch",
		middleware.SearchRateLimiter(),
		handler.SendPushToMultipleUsers,
	)

	// POST /api/v1/push/send/topic - Send to topic subscribers
	// Body: { topic, title, body, data, image_url, sound, priority, badge }
	// Rate limit: 30 requests/minute
	pushSend.Post("/topic",
		middleware.SearchRateLimiter(),
		handler.SendPushToTopic,
	)

	// POST /api/v1/push/test - Send test notification
	// Body: { token }
	// Rate limit: 100 requests/minute
	push.Post("/test",
		middleware.ApplicationRateLimiter(),
		handler.SendTestNotification,
	)
}
