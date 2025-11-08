package routes

import (
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAdminAuthRoutes configures admin authentication routes
// Routes: /api/v1/auth/admin/*
func SetupAdminAuthRoutes(api fiber.Router, deps *Dependencies, adminAuthMw *middleware.AdminAuthMiddleware) {
	adminAuth := api.Group("/auth/admin")

	// ===========================================
	// Public Routes - Admin Authentication
	// ===========================================

	adminAuth.Post("/login",
		middleware.AuthRateLimiter(),
		deps.AdminAuthHandler.Login,
	)

	adminAuth.Post("/refresh-token",
		deps.AdminAuthHandler.RefreshToken,
	)

	// ===========================================
	// Protected Routes - Require Admin Auth
	// ===========================================

	adminAuth.Post("/logout",
		adminAuthMw.AdminAuthRequired(),
		deps.AdminAuthHandler.Logout,
	)

	adminAuth.Get("/me",
		adminAuthMw.AdminAuthRequired(),
		deps.AdminAuthHandler.GetProfile,
	)

	adminAuth.Put("/change-password",
		adminAuthMw.AdminAuthRequired(),
		deps.AdminAuthHandler.ChangePassword,
	)
}
