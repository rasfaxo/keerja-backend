package routes

import (
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupApplicationRoutes configures application routes
// Routes: /api/v1/applications/*
//
// Candidate Endpoints (7):
//   - POST   /                         Submit job application
//   - POST   /jobs/:job_id/apply       Apply to specific job
//   - GET    /my-applications          List my applications
//   - GET    /:id                      Get application details
//   - DELETE /:id/withdraw             Withdraw application
//   - POST   /:id/documents            Upload application document
//   - POST   /:id/rate                 Rate application experience
//
// Employer Endpoints (14):
//   - GET    /job/:job_id              List applications for job
//   - POST   /search                   Search applications
//   - PATCH  /:id/status               Update application status
//   - POST   /bulk-status              Bulk update status
//   - PATCH  /:id/stage                Update hiring stage
//   - POST   /:id/notes                Add note to application
//   - PUT    /:id/notes/:note_id       Update existing note
//   - POST   /:id/interviews           Schedule interview
//   - PUT    /:id/interviews/:int_id   Update interview
//   - PATCH  /:id/interviews/:int_id/reschedule  Reschedule interview
//   - PATCH  /:id/interviews/:int_id/complete    Complete interview
//   - DELETE /:id/interviews/:int_id   Cancel interview
//   - PATCH  /:id/bookmark             Toggle bookmark
//   - PATCH  /:id/viewed               Mark as viewed
//
// Total: 21 endpoints
func SetupApplicationRoutes(api fiber.Router, deps *Dependencies, authMw *middleware.AuthMiddleware) {
	// ============================================
	// CANDIDATE ROUTES (7 endpoints)
	// All require authentication
	// ============================================
	applications := api.Group("/applications")
	applications.Use(authMw.AuthRequired())

	// POST /api/v1/applications - Submit job application
	// Body: { job_id, cover_letter, resume_url, documents, ... }
	// Rate limit: 100 requests/minute
	applications.Post("/",
		authMw.JobSeekerOnly(),
		middleware.ApplicationRateLimiter(),
		deps.ApplicationHandler.Apply,
	)

	// POST /api/v1/applications/jobs/:job_id/apply - Apply to specific job
	// Alternative endpoint for job-specific application
	// Body: { cover_letter, resume_url, documents }
	// Rate limit: 100 requests/minute
	applications.Post("/jobs/:job_id/apply",
		authMw.JobSeekerOnly(),
		middleware.ApplicationRateLimiter(),
		deps.ApplicationHandler.ApplyToJob,
	)

	// GET /api/v1/applications/my-applications - List candidate's applications
	// Query params: page, limit, status, job_id
	// Rate limit: 30 requests/minute
	applications.Get("/my-applications",
		authMw.JobSeekerOnly(),
		middleware.SearchRateLimiter(),
		deps.ApplicationHandler.GetMyApplications,
	)

	// GET /api/v1/applications/:id - Get application details
	// Returns: Full application details with job & company info
	applications.Get("/:id",
		deps.ApplicationHandler.GetApplication,
	)

	// DELETE /api/v1/applications/:id/withdraw - Withdraw application
	// Body: { reason }
	// Only candidate who submitted can withdraw
	applications.Delete("/:id/withdraw",
		authMw.JobSeekerOnly(),
		deps.ApplicationHandler.Withdraw,
	)

	// POST /api/v1/applications/:id/documents - Upload application document
	// Body: multipart/form-data with document file
	// Rate limit: 100 requests/minute
	applications.Post("/:id/documents",
		authMw.JobSeekerOnly(),
		middleware.ApplicationRateLimiter(),
		deps.ApplicationHandler.UploadDocument,
	)

	// POST /api/v1/applications/:id/rate - Rate application experience
	// Body: { rating, comment }
	// Candidate can rate their experience after process completion
	applications.Post("/:id/rate",
		authMw.JobSeekerOnly(),
		deps.ApplicationHandler.RateExperience,
	)

	// ============================================
	// EMPLOYER ROUTES (14 endpoints)
	// All require authentication + employer role
	// ============================================
	employer := applications.Group("")
	employer.Use(authMw.EmployerOnly())

	// GET /api/v1/applications/job/:job_id - List applications for specific job
	// Query params: page, limit, status, stage
	// Rate limit: 30 requests/minute
	employer.Get("/job/:job_id",
		middleware.SearchRateLimiter(),
		deps.ApplicationHandler.ListByJob,
	)

	// POST /api/v1/applications/search - Search applications
	// Body: { query, filters, pagination }
	// Rate limit: 30 requests/minute
	employer.Post("/search",
		middleware.SearchRateLimiter(),
		deps.ApplicationHandler.SearchApplications,
	)

	// PATCH /api/v1/applications/:id/status - Update application status
	// Body: { status, reason }
	// Status: pending, reviewing, shortlisted, rejected, accepted
	employer.Patch("/:id/status",
		deps.ApplicationHandler.UpdateStatus,
	)

	// POST /api/v1/applications/bulk-status - Bulk update application status
	// Body: { application_ids: [], status, reason }
	// Rate limit: 30 requests/minute (bulk operation)
	employer.Post("/bulk-status",
		middleware.SearchRateLimiter(),
		deps.ApplicationHandler.BulkUpdateStatus,
	)

	// PATCH /api/v1/applications/:id/stage - Update hiring stage
	// Body: { stage }
	// Stage: applied, screening, interview, offer, hired
	employer.Patch("/:id/stage",
		deps.ApplicationHandler.UpdateStage,
	)

	// POST /api/v1/applications/:id/notes - Add note to application
	// Body: { content, is_private }
	// Rate limit: 100 requests/minute
	employer.Post("/:id/notes",
		middleware.ApplicationRateLimiter(),
		deps.ApplicationHandler.AddNote,
	)

	// PUT /api/v1/applications/:id/notes/:note_id - Update existing note
	// Body: { content, is_private }
	employer.Put("/:id/notes/:note_id",
		deps.ApplicationHandler.UpdateNote,
	)

	// POST /api/v1/applications/:id/interviews - Schedule interview
	// Body: { type, scheduled_at, duration, location, notes, interviewers }
	// Rate limit: 100 requests/minute
	employer.Post("/:id/interviews",
		middleware.ApplicationRateLimiter(),
		deps.ApplicationHandler.ScheduleInterview,
	)

	// PUT /api/v1/applications/:id/interviews/:interview_id - Update interview
	// Body: { type, scheduled_at, duration, location, notes, status }
	employer.Put("/:id/interviews/:interview_id",
		deps.ApplicationHandler.UpdateInterview,
	)

	// PATCH /api/v1/applications/:id/interviews/:interview_id/reschedule - Reschedule interview
	// Body: { scheduled_at, reason }
	employer.Patch("/:id/interviews/:interview_id/reschedule",
		deps.ApplicationHandler.RescheduleInterview,
	)

	// PATCH /api/v1/applications/:id/interviews/:interview_id/complete - Complete interview
	// Body: { feedback, rating, decision }
	employer.Patch("/:id/interviews/:interview_id/complete",
		deps.ApplicationHandler.CompleteInterview,
	)

	// DELETE /api/v1/applications/:id/interviews/:interview_id - Cancel interview
	// Query param: reason
	employer.Delete("/:id/interviews/:interview_id",
		deps.ApplicationHandler.CancelInterview,
	)

	// PATCH /api/v1/applications/:id/bookmark - Toggle bookmark
	// Body: { bookmarked: true/false }
	employer.Patch("/:id/bookmark",
		deps.ApplicationHandler.Bookmark,
	)

	// PATCH /api/v1/applications/:id/viewed - Mark application as viewed
	// Body: { viewed: true }
	employer.Patch("/:id/viewed",
		deps.ApplicationHandler.MarkViewed,
	)
}
