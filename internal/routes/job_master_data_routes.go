package routes

import (
	"time"

	"keerja-backend/internal/handler/http"
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupJobMasterDataRoutes sets up routes for job master data (Phase 1-4)
func SetupJobMasterDataRoutes(api fiber.Router, handler *http.MasterDataHandler, authMw *middleware.AuthMiddleware) {
	master := api.Group("/master")

	// Phase 1: Job Titles endpoint (public, with rate limiting)
	// GET /api/v1/master/job-titles?q=search&limit=20
	master.Get("/job-titles",
		middleware.RateLimitByIP(30, 1*time.Minute), // 30 requests per minute per IP
		handler.GetJobTitles,
	)

	// Phase 3: Job Options endpoint (public, heavily cached)
	// GET /api/v1/master/job-options
	master.Get("/job-options",
		middleware.RateLimitByIP(60, 1*time.Minute), // 60 requests per minute (high because it's cached)
		handler.GetJobOptions,
	)

	// Phase 4: Job posting form options for mobile (public)
	// GET /api/v1/master/job-posting-form-options
	master.Get("/job-posting-form-options",
		middleware.RateLimitByIP(60, 1*time.Minute),
		handler.GetJobPostingFormOptions,
	)

	// Admin-only routes for managing job titles
	admin := api.Group("/admin/master")
	admin.Use(authMw.AuthRequired())
	admin.Use(authMw.AdminOnly())

	// POST /api/v1/admin/master/job-titles
	admin.Post("/job-titles", handler.CreateJobTitle)

	// GET /api/v1/admin/master/job-titles/:id
	admin.Get("/job-titles/:id", handler.GetJobTitleByID)

	// PUT /api/v1/admin/master/job-titles/:id
	admin.Put("/job-titles/:id", handler.UpdateJobTitle)

	// DELETE /api/v1/admin/master/job-titles/:id
	admin.Delete("/job-titles/:id", handler.DeleteJobTitle)
}
