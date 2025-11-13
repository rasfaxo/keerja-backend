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
// Employer Endpoints (7):
//   - POST   /                   Create new job posting
//   - POST   /draft              Save job draft (Phase 6)
//   - PUT    /:id                Update existing job
//   - DELETE /:id                Delete job posting
//   - GET    /my-jobs            List employer's own jobs
//   - PATCH  /:id/publish        Publish draft job
//   - PATCH  /:id/close          Close active job
//
// Total: 10 endpoints
func SetupJobRoutes(api fiber.Router, deps *Dependencies, authMw *middleware.AuthMiddleware) {
	jobs := api.Group("/jobs")

	// ============================================
	// PUBLIC ROUTES (5 endpoints)
	// ============================================

	// GET /api/v1/jobs/job-types - Get job types options for mobile
	// Returns: job_types, work_policies, work_addresses, salary_ranges
	// Auth required to get user's company addresses
	jobs.Get("/job-types",
		authMw.AuthRequired(),
		deps.JobHandler.GetJobTypesOptions,
	)

	// GET /api/v1/jobs/job-requirements - Get job requirements options for mobile
	// Returns: genders, age_ranges, skills_info, education_levels, experience_levels
	// Auth required
	jobs.Get("/job-requirements",
		authMw.AuthRequired(),
		deps.JobHandler.GetJobRequirements,
	)

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
	// PROTECTED ROUTES - EMPLOYER ONLY (7 endpoints)
	// ============================================
	protected := jobs.Group("")
	protected.Use(authMw.AuthRequired())
	protected.Use(authMw.EmployerOnly())

	// POST /api/v1/jobs/draft - Save job draft (Phase 6)
	// Body: SaveJobDraftRequest with all job posting data
	// Rate limit: 100 requests/minute
	// Response: 201 Created (new draft) or 200 OK (update)
	protected.Post("/draft",
		middleware.ApplicationRateLimiter(),
		deps.JobHandler.SaveJobDraft,
	)

	// POST /api/v1/jobs - Create new job posting
	// Body: { title, description, requirements, company_id, ... }
	// Rate limit: 100 requests/minute (default)
	protected.Post("/",
		middleware.ApplicationRateLimiter(),
		deps.JobHandler.CreateJob,
	)

	// GET /api/v1/jobs/my-jobs - List employer's jobs
	// Query params: page, limit, status
	// Rate limit: 30 requests/minute
	// IMPORTANT: This must be defined BEFORE /:id routes to avoid conflicts
	protected.Get("/my-jobs",
		middleware.SearchRateLimiter(),
		deps.JobHandler.GetMyJobs,
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

	// PATCH /api/v1/jobs/:id/publish - Publish draft job (Phase 7)
	// Changes status from draft to pending_review
	// Requires company to be verified
	// Triggers admin notification for review
	protected.Patch("/:id/publish",
		deps.JobHandler.PublishJob,
	)

	// PATCH /api/v1/jobs/:id/close - Close active job
	// Changes status from published to closed
	protected.Patch("/:id/close",
		deps.JobHandler.CloseJob,
	)
}
