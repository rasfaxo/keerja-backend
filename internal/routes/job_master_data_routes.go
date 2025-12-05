package routes

import (
	"time"

	"keerja-backend/internal/handler/http/master"
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupJobMasterDataRoutes sets up routes for job master data (Phase 1-4)
func SetupJobMasterDataRoutes(api fiber.Router, handler *master.MasterDataHandler, authMw *middleware.AuthMiddleware) {
	master := api.Group("/master")

	// Phase 1: Job Titles endpoint (public, with rate limiting)
	// GET /api/v1/master/job-titles?q=search&limit=20
	master.Get("/job-titles",
		middleware.RateLimitByIP(30, 1*time.Minute), // 30 requests per minute per IP
		handler.GetJobTitles,
	)

	// ==================== PAGE-BASED JOB FORM OPTIONS ====================
	// Endpoints organized by Create Job form pages (UI-driven)

	jobForm := master.Group("/job-form")

	// Page 1: Job Details - Categories & Subcategories
	// GET /api/v1/master/job-form/job-details
	// Returns: job_categories (with nested subcategories)
	// Used for: Job Name (job_category_id), Job Field (job_subcategory_id)
	jobForm.Get("/job-details",
		middleware.RateLimitByIP(60, 1*time.Minute),
		handler.GetJobDetailsOptions,
	)

	// Page 2: Job Type - Types, Policies, Addresses, Salary
	// GET /api/v1/master/job-form/job-type
	// Returns: job_types, work_policies, company_addresses (if auth), salary_ranges
	// Used for: Job Type (job_type_id), Work Policy (work_policy_id), Work Address (company_address_id), Salary
	jobForm.Get("/job-type",
		middleware.RateLimitByIP(60, 1*time.Minute),
		authMw.OptionalAuth(), // Optional auth for company addresses
		handler.GetJobTypeOptions,
	)

	// Page 3: Job Requirements - Gender, Age, Skills, Education, Experience
	// GET /api/v1/master/job-form/job-requirements?skills_page=1&skills_limit=50&skills_q=search
	// Returns: gender_preferences, age_limits, skills (paginated), education_levels, experience_levels
	// Used for: Gender, Age, Required Skill, Minimum Education Required, Required Work Experience
	jobForm.Get("/job-requirements",
		middleware.RateLimitByIP(60, 1*time.Minute),
		handler.GetJobRequirementsOptions,
	)

	// Page 4: Job Description - Metadata only (text input page)
	// GET /api/v1/master/job-form/job-description
	// Returns: max_length, min_length, placeholder, tips
	// Used for: Job description text input guidance
	jobForm.Get("/job-description",
		middleware.RateLimitByIP(60, 1*time.Minute),
		handler.GetJobDescriptionOptions,
	)

	// ==================== UTILITY ENDPOINTS ====================

	// Legacy job options endpoint (all options in one response)
	// GET /api/v1/master/job-options
	master.Get("/job-options",
		middleware.RateLimitByIP(60, 1*time.Minute),
		handler.GetJobOptions,
	)

	// Paginated skills endpoint (standalone)
	// GET /api/v1/master/skills?q=search&page=1&limit=50
	master.Get("/skills",
		middleware.RateLimitByIP(60, 1*time.Minute),
		handler.GetSkillsPaginated,
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
