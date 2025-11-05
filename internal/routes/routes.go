package routes

import (
	"keerja-backend/internal/config"
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/handler/http"
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// Dependencies holds all handler dependencies
type Dependencies struct {
	Config             *config.Config
	AuthHandler        *http.AuthHandler
	UserHandler        *http.UserHandler
	JobHandler         *http.JobHandler         // Job management (9 endpoints)
	ApplicationHandler *http.ApplicationHandler // Application management (21 endpoints)
	AdminHandler       *http.AdminHandler       // Admin moderation & job approval

	// Company handlers (split by domain for better organization)
	CompanyBasicHandler   *http.CompanyBasicHandler   // CRUD operations (10 endpoints)
	CompanyProfileHandler *http.CompanyProfileHandler // Profile & social features (8 endpoints)
	CompanyReviewHandler  *http.CompanyReviewHandler  // Review system (5 endpoints)
	CompanyStatsHandler   *http.CompanyStatsHandler   // Statistics & queries (3 endpoints)
	CompanyInviteHandler  *http.CompanyInviteHandler  // Employee invitation (5 endpoints)

	// Master data handlers
	SkillsMasterHandler *http.SkillsMasterHandler // Skills master data (8 endpoints)
	MasterDataHandlers  *MasterDataHandlers       // Industry, company size, location (10 endpoints)
	MasterDataHandler   *http.MasterDataHandler   // Job titles & options (Phase 1-4)

	// FCM Notification handlers (Firebase Cloud Messaging)
	DeviceTokenHandler      *http.DeviceTokenHandler      // Device token management (6 endpoints)
	PushNotificationHandler *http.PushNotificationHandler // Push notifications (5 endpoints)

	// Services (for middlewares)
	CompanyService company.CompanyService
}

// SetupRoutes configures all application routes
// This is the main entry point for route configuration
func SetupRoutes(app *fiber.App, deps *Dependencies) {
	// Initialize auth middleware
	authMw := middleware.NewAuthMiddleware(deps.Config)

	// Initialize permission middleware
	permMw := middleware.NewPermissionMiddleware(deps.CompanyService)

	// API v1 group
	api := app.Group("/api/v1")

	// Health check endpoint
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Keerja API is running",
		})
	})

	// Setup route groups (each in separate file)
	SetupAuthRoutes(api, deps, authMw)               // auth_routes.go
	SetupUserRoutes(api, deps, authMw)               // user_routes.go
	SetupJobRoutes(api, deps, authMw)                // job_routes.go
	SetupApplicationRoutes(api, deps, authMw)        // application_routes.go
	SetupCompanyRoutes(api, deps, authMw, permMw)    // company_routes.go
	SetupAdminRoutes(api, deps, authMw)              // admin_routes.go
	SetupSkillsRoutes(api, deps.SkillsMasterHandler) // skills_routes.go

	// Master data routes (industries, company sizes, locations)
	if deps.MasterDataHandlers != nil {
		SetupMasterDataRoutes(api, deps.MasterDataHandlers) // master_routes.go
	}

	// Job master data routes (job titles & options - Phase 1-4)
	if deps.MasterDataHandler != nil {
		SetupJobMasterDataRoutes(api, deps.MasterDataHandler, authMw) // job_master_data_routes.go
	}

	// FCM Notification routes
	if deps.DeviceTokenHandler != nil {
		SetupDeviceTokenRoutes(api, deps.DeviceTokenHandler, authMw) // device_token_routes.go
	}
	if deps.PushNotificationHandler != nil {
		SetupPushNotificationRoutes(api, deps.PushNotificationHandler, authMw) // push_notification_routes.go
	}
}
