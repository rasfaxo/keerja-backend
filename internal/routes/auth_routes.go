package routes

import (
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes configures authentication routes
// Routes: /api/v1/auth/*
func SetupAuthRoutes(api fiber.Router, deps *Dependencies, authMw *middleware.AuthMiddleware) {
	auth := api.Group("/auth")

	// Public routes
	auth.Post("/register",
		middleware.RegistrationRateLimiter(),
		deps.AuthHandler.Register,
	)

	auth.Post("/login",
		middleware.AuthRateLimiter(),
		deps.AuthHandler.Login,
	)

	auth.Post("/verify-email", deps.AuthHandler.VerifyEmail)

	auth.Post("/forgot-password",
		middleware.EmailRateLimiter(),
		deps.AuthHandler.ForgotPassword,
	)

	auth.Post("/reset-password", deps.AuthHandler.ResetPassword)

	auth.Post("/resend-verification",
		middleware.EmailRateLimiter(),
		deps.AuthHandler.ResendVerification,
	)

	// Protected routes
	auth.Post("/refresh-token",
		authMw.AuthRequired(),
		deps.AuthHandler.RefreshToken,
	)

	auth.Post("/logout",
		authMw.AuthRequired(),
		deps.AuthHandler.Logout,
	)
}
