package routes

import (
	"keerja-backend/internal/handler/http"
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupDeviceTokenRoutes configures all device token related routes
// Routes: /api/v1/device-tokens/*
//
// Endpoints (6):
//   - POST   /                     Register device token
//   - GET    /                     Get user's device tokens (with pagination)
//   - GET    /stats                Get device token statistics
//   - GET    /:id                  Get specific device token by ID
//   - DELETE /:token               Unregister device token
//   - POST   /validate             Validate device token with FCM
//
// Total: 6 endpoints
func SetupDeviceTokenRoutes(api fiber.Router, handler *http.DeviceTokenHandler, authMw *middleware.AuthMiddleware) {
	// Device Token routes group
	deviceTokens := api.Group("/device-tokens")

	// Apply authentication middleware to all device token routes
	deviceTokens.Use(authMw.AuthRequired())

	// POST /api/v1/device-tokens - Register new device token
	// Body: { token, platform, device_info }
	// Rate limit: 100 requests/minute
	deviceTokens.Post("/",
		middleware.ApplicationRateLimiter(),
		handler.RegisterDeviceToken,
	)

	// GET /api/v1/device-tokens - Get user's device tokens
	// Query params: page, page_size, platform, is_active
	// Rate limit: 30 requests/minute
	deviceTokens.Get("/",
		middleware.SearchRateLimiter(),
		handler.GetUserDevices,
	)

	// GET /api/v1/device-tokens/stats - Get device token statistics
	// Returns: total_tokens, active_tokens, inactive_tokens, by platform
	// Rate limit: 30 requests/minute
	deviceTokens.Get("/stats",
		middleware.SearchRateLimiter(),
		handler.GetDeviceTokenStats,
	)

	// GET /api/v1/device-tokens/:id - Get specific device token by ID
	// Rate limit: 30 requests/minute
	deviceTokens.Get("/:id",
		middleware.SearchRateLimiter(),
		handler.GetDeviceTokenByID,
	)

	// DELETE /api/v1/device-tokens/:token - Unregister device token
	// Rate limit: 100 requests/minute
	deviceTokens.Delete("/:token",
		middleware.ApplicationRateLimiter(),
		handler.UnregisterDeviceToken,
	)

	// POST /api/v1/device-tokens/validate - Validate device token with FCM
	// Body: { token }
	// Rate limit: 30 requests/minute
	deviceTokens.Post("/validate",
		middleware.SearchRateLimiter(),
		handler.ValidateDeviceToken,
	)
}
