package routes

import (
	"keerja-backend/internal/config"
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
	AdminHandler       interface{}              // TODO: Change to *http.AdminHandler when implemented

	// Company handlers (split by domain for better organization)
	CompanyBasicHandler   *http.CompanyBasicHandler   // CRUD operations (10 endpoints)
	CompanyProfileHandler *http.CompanyProfileHandler // Profile & social features (8 endpoints)
	CompanyReviewHandler  *http.CompanyReviewHandler  // Review system (5 endpoints)
	CompanyStatsHandler   *http.CompanyStatsHandler   // Statistics & queries (3 endpoints)
	CompanyInviteHandler  *http.CompanyInviteHandler  // Employee invitation (1 endpoint)
}

// SetupRoutes configures all application routes
// This is the main entry point for route configuration
func SetupRoutes(app *fiber.App, deps *Dependencies) {
	// Initialize auth middleware
	authMw := middleware.NewAuthMiddleware(deps.Config)

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
	SetupAuthRoutes(api, deps, authMw)        // auth_routes.go
	SetupUserRoutes(api, deps, authMw)        // user_routes.go
	SetupJobRoutes(api, deps, authMw)         // job_routes.go
	SetupApplicationRoutes(api, deps, authMw) // application_routes.go
	SetupCompanyRoutes(api, deps, authMw)     // company_routes.go
	SetupAdminRoutes(api, deps, authMw)       // admin_routes.go
}
