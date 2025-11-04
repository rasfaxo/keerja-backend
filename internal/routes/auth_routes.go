package routes

import (
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes configures authentication routes
// Routes: /api/v1/auth/*
func SetupAuthRoutes(api fiber.Router, deps *Dependencies, authMw *middleware.AuthMiddleware) {
	auth := api.Group("/auth")

	// ===========================================
	// Public Routes - Traditional Auth
	// ===========================================

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

	// ===========================================
	// Public Routes - OTP Registration
	// ===========================================

	auth.Post("/register-otp",
		middleware.RegistrationRateLimiter(),
		deps.AuthHandler.RegisterWithOTP,
	)

	auth.Post("/verify-email-otp",
		middleware.AuthRateLimiter(),
		deps.AuthHandler.VerifyEmailOTP,
	)

	auth.Post("/resend-otp",
		middleware.EmailRateLimiter(),
		deps.AuthHandler.ResendOTP,
	)

	// ===========================================
	// Public Routes - Forgot Password with OTP
	// ===========================================

	auth.Post("/forgot-password-otp",
		middleware.EmailRateLimiter(),
		deps.AuthHandler.ForgotPasswordOTP,
	)

	auth.Post("/reset-password-otp",
		middleware.AuthRateLimiter(),
		deps.AuthHandler.ResetPasswordOTP,
	)

	// ===========================================
	// Public Routes - Login with Remember Me
	// ===========================================

	auth.Post("/login-remember",
		middleware.AuthRateLimiter(),
		deps.AuthHandler.LoginWithRememberMe,
	)

	auth.Post("/refresh",
		deps.AuthHandler.RefreshAccessToken,
	)

	// ===========================================
	// Public Routes - OAuth
	// ===========================================

	oauth := auth.Group("/oauth")

	oauth.Get("/google",
		deps.AuthHandler.InitiateGoogleLogin,
	)

	oauth.Get("/google/callback",
		deps.AuthHandler.HandleGoogleCallback,
	)

	// ===========================================
	// Protected Routes - Device Management
	// ===========================================

	auth.Get("/devices",
		authMw.AuthRequired(),
		deps.AuthHandler.GetActiveDevices,
	)

	auth.Post("/devices/revoke",
		authMw.AuthRequired(),
		deps.AuthHandler.RevokeDevice,
	)

	auth.Post("/logout-all",
		authMw.AuthRequired(),
		deps.AuthHandler.LogoutAllDevices,
	)

	// ===========================================
	// Protected Routes - OAuth Management
	// ===========================================

	oauth.Get("/connected",
		authMw.AuthRequired(),
		deps.AuthHandler.GetConnectedProviders,
	)

	oauth.Delete("/:provider",
		authMw.AuthRequired(),
		deps.AuthHandler.DisconnectOAuth,
	)

	// ===========================================
	// Protected Routes - Legacy
	// ===========================================

	auth.Post("/refresh-token",
		authMw.AuthRequired(),
		deps.AuthHandler.RefreshToken,
	)

	auth.Post("/logout",
		authMw.AuthRequired(),
		deps.AuthHandler.Logout,
	)
}
