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
	JobHandler         interface{} // TODO: Change to *http.JobHandler when implemented
	ApplicationHandler interface{} // TODO: Change to *http.ApplicationHandler when implemented
	CompanyHandler     interface{} // TODO: Change to *http.CompanyHandler when implemented
	AdminHandler       interface{} // TODO: Change to *http.AdminHandler when implemented
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
