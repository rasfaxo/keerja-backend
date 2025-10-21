package routes

import (
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupJobRoutes configures job routes
// Routes: /api/v1/jobs/*
//
// Public Endpoints (3):
//   - GET    /                   List all jobs with filters & pagination
//   - GET    /:id                Get job details by ID
//   - POST   /search             Advanced job search
//
// Employer Endpoints (6):
//   - POST   /                   Create new job posting
//   - PUT    /:id                Update existing job
//   - DELETE /:id                Delete job posting
//   - GET    /my-jobs            List employer's own jobs
//   - PATCH  /:id/publish        Publish draft job
//   - PATCH  /:id/close          Close active job
//
// Total: 9 endpoints
func SetupJobRoutes(api fiber.Router, deps *Dependencies, authMw *middleware.AuthMiddleware) {
	jobs := api.Group("/jobs")

	// ============================================
	// PUBLIC ROUTES (3 endpoints)
	// ============================================

	// GET /api/v1/jobs - List all jobs
	// Query params: page, limit, company_id, location, type, level, etc.
	// Rate limit: 30 requests/minute (search rate)
	jobs.Get("/",
		middleware.SearchRateLimiter(),
		deps.JobHandler.ListJobs,
	)

	// GET /api/v1/jobs/:id - Get job details
	// Returns: Job details with company info
	jobs.Get("/:id",
		deps.JobHandler.GetJob,
	)

	// POST /api/v1/jobs/search - Advanced job search
	// Body: { query, location, filters, pagination }
	// Rate limit: 30 requests/minute
	jobs.Post("/search",
		middleware.SearchRateLimiter(),
		deps.JobHandler.SearchJobs,
	)

	// ============================================
	// PROTECTED ROUTES - EMPLOYER ONLY (6 endpoints)
	// ============================================
	protected := jobs.Group("")
	protected.Use(authMw.AuthRequired())
	protected.Use(authMw.EmployerOnly())

	// POST /api/v1/jobs - Create new job posting
	// Body: { title, description, requirements, company_id, ... }
	// Rate limit: 100 requests/minute (default)
	protected.Post("/",
		middleware.ApplicationRateLimiter(),
		deps.JobHandler.CreateJob,
	)

	// PUT /api/v1/jobs/:id - Update job posting
	// Body: { title, description, requirements, ... }
	// Rate limit: 100 requests/minute
	protected.Put("/:id",
		middleware.ApplicationRateLimiter(),
		deps.JobHandler.UpdateJob,
	)

	// DELETE /api/v1/jobs/:id - Delete job posting
	// Only owner can delete
	protected.Delete("/:id",
		deps.JobHandler.DeleteJob,
	)

	// GET /api/v1/jobs/my-jobs - List employer's jobs
	// Query params: page, limit, status
	// Rate limit: 30 requests/minute
	protected.Get("/my-jobs",
		middleware.SearchRateLimiter(),
		deps.JobHandler.GetMyJobs,
	)

	// PATCH /api/v1/jobs/:id/publish - Publish draft job
	// Changes status from draft to published
	protected.Patch("/:id/publish",
		deps.JobHandler.PublishJob,
	)

	// PATCH /api/v1/jobs/:id/close - Close active job
	// Changes status from published to closed
	protected.Patch("/:id/close",
		deps.JobHandler.CloseJob,
	)
}
