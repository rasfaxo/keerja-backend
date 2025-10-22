package http

import (
	"strconv"

	"keerja-backend/internal/domain/application"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// ApplicationHandler handles application-related HTTP requests
type ApplicationHandler struct {
	appService application.ApplicationService
}

func NewApplicationHandler(appService application.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		appService: appService,
	}
}

// POST /applications - Submit job application
func (h *ApplicationHandler) Apply(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req application.ApplyJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}

	// Set user ID
	req.UserID = userID

	// Sanitize inputs
	req.CoverLetter = utils.SanitizeHTML(req.CoverLetter)
	req.Source = utils.SanitizeString(req.Source)

	// Apply for job
	app, err := h.appService.ApplyForJob(ctx, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, ErrApplicationNotFound, err.Error())
	}

	return utils.CreatedResponse(c, MsgApplicationSubmit, app)
}

// POST /applications/jobs/:job_id/apply - Apply to specific job
func (h *ApplicationHandler) ApplyToJob(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	jobID, err := strconv.ParseInt(c.Params("job_id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	var req application.ApplyJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}

	// Set job ID and user ID
	req.JobID = jobID
	req.UserID = userID

	// Sanitize inputs
	req.CoverLetter = utils.SanitizeHTML(req.CoverLetter)
	req.Source = utils.SanitizeString(req.Source)

	// Apply for job
	app, err := h.appService.ApplyForJob(ctx, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, ErrAlreadyApplied, err.Error())
	}

	return utils.CreatedResponse(c, MsgApplicationSubmit, app)
}

// GET /applications/my-applications - List candidate's applications
func (h *ApplicationHandler) GetMyApplications(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	// Parse pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// Parse filters
	filter := application.ApplicationFilter{}
	if status := c.Query("status"); status != "" {
		filter.Status = status
	}
	if jobIDStr := c.Query("job_id"); jobIDStr != "" {
		if jobID, err := strconv.ParseInt(jobIDStr, 10, 64); err == nil {
			filter.JobID = jobID
		}
	}

	// Get applications
	response, err := h.appService.GetMyApplications(ctx, userID, filter, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, ErrInternalServer, err.Error())
	}

	return utils.SuccessResponse(c, MsgFetchedSuccess, response)
}

// GET /applications/:id - Get application details
func (h *ApplicationHandler) GetApplication(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	// Get application details
	response, err := h.appService.GetApplicationDetail(ctx, appID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, ErrApplicationNotFound, err.Error())
	}

	return utils.SuccessResponse(c, MsgFetchedSuccess, response)
}

// DELETE /applications/:id/withdraw - Withdraw application
func (h *ApplicationHandler) Withdraw(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	// Withdraw application
	if err := h.appService.WithdrawApplication(ctx, appID, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, ErrCannotWithdraw, err.Error())
	}

	return utils.SuccessResponse(c, MsgOperationSuccess, nil)
}

// POST /applications/:id/documents - Upload application document
func (h *ApplicationHandler) UploadDocument(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	var req application.UploadDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}

	// Set application ID and user ID
	req.ApplicationID = appID
	req.UserID = userID

	// Upload document
	doc, err := h.appService.UploadApplicationDocument(ctx, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, ErrFileUploadFailed, err.Error())
	}

	return utils.CreatedResponse(c, MsgUploadSuccess, doc)
}

// POST /applications/:id/rate - Rate application experience
func (h *ApplicationHandler) RateExperience(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	// Parse rating and comment from body
	type RateRequest struct {
		Rating  int    `json:"rating" validate:"required,min=1,max=5"`
		Comment string `json:"comment"`
	}

	var req RateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}

	// TODO: Implement rate experience in domain service
	// For now, return not implemented
	_ = userID
	_ = appID
	_ = req

	return utils.ErrorResponse(c, fiber.StatusNotImplemented, ErrFailedOperation, "Rating feature not yet implemented")
}

// ============================================================================
// EMPLOYER ENDPOINTS
// ============================================================================

// GET /applications/job/:job_id - List applications for job
func (h *ApplicationHandler) ListByJob(c *fiber.Ctx) error {
	ctx := c.Context()
	// employerID := middleware.GetUserID(c)  // TODO: Add authorization check

	jobID, err := strconv.ParseInt(c.Params("job_id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// Parse filters
	filter := application.ApplicationFilter{}
	if status := c.Query("status"); status != "" {
		filter.Status = status
	}

	// Get applications
	response, err := h.appService.GetJobApplications(ctx, jobID, filter, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, ErrInternalServer, err.Error())
	}

	return utils.SuccessResponse(c, MsgFetchedSuccess, response)
}

// POST /applications/search - Search applications
func (h *ApplicationHandler) SearchApplications(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse search filter
	var filter application.ApplicationSearchFilter
	if err := c.BodyParser(&filter); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// Search applications
	response, err := h.appService.SearchApplications(ctx, filter, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, ErrInternalServer, err.Error())
	}

	return utils.SuccessResponse(c, MsgFetchedSuccess, response)
}

// PATCH /applications/:id/status - Update application status
func (h *ApplicationHandler) UpdateStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	type UpdateStatusRequest struct {
		Status string `json:"status" validate:"required"`
		Notes  string `json:"notes"`
	}

	var req UpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}

	// Update status based on value
	var updateErr error
	switch req.Status {
	case "screening":
		updateErr = h.appService.MoveToScreening(ctx, appID, employerID, req.Notes)
	case "shortlisted":
		updateErr = h.appService.MoveToShortlist(ctx, appID, employerID, req.Notes)
	case "interview":
		updateErr = h.appService.MoveToInterview(ctx, appID, employerID, req.Notes)
	case "offered":
		updateErr = h.appService.MakeOffer(ctx, appID, employerID, req.Notes)
	case "hired":
		updateErr = h.appService.MarkAsHired(ctx, appID, employerID, req.Notes)
	case "rejected":
		updateErr = h.appService.RejectApplication(ctx, appID, employerID, req.Notes)
	default:
		return utils.BadRequestResponse(c, ErrInvalidApplicationStage)
	}

	if updateErr != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, ErrFailedOperation, updateErr.Error())
	}

	return utils.SuccessResponse(c, MsgStatusUpdated, nil)
}

// POST /applications/bulk-status - Bulk update status
func (h *ApplicationHandler) BulkUpdateStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	type BulkUpdateRequest struct {
		ApplicationIDs []int64 `json:"application_ids" validate:"required,min=1"`
		Status         string  `json:"status" validate:"required"`
	}

	var req BulkUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}

	// Bulk update
	if err := h.appService.BulkUpdateStatus(ctx, req.ApplicationIDs, req.Status, employerID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, MsgOperationSuccess, nil)
}

// PATCH /applications/:id/stage - Update hiring stage
func (h *ApplicationHandler) UpdateStage(c *fiber.Ctx) error {
	// For now, redirect to UpdateStatus since they're similar
	return h.UpdateStatus(c)
}

// POST /applications/:id/notes - Add note
func (h *ApplicationHandler) AddNote(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	var req application.AddNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}

	// Set application and author ID
	req.ApplicationID = appID
	req.AuthorID = employerID

	// Add note
	note, err := h.appService.AddNote(ctx, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, ErrFailedOperation, err.Error())
	}

	return utils.CreatedResponse(c, MsgCreatedSuccess, note)
}

// PUT /applications/:id/notes/:note_id - Update note
func (h *ApplicationHandler) UpdateNote(c *fiber.Ctx) error {
	ctx := c.Context()
	// employerID := middleware.GetUserID(c)

	noteID, err := strconv.ParseInt(c.Params("note_id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	var req application.UpdateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}

	// Update note
	note, err := h.appService.UpdateNote(ctx, noteID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, MsgUpdatedSuccess, note)
}

// POST /applications/:id/interviews - Schedule interview
func (h *ApplicationHandler) ScheduleInterview(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	var req application.ScheduleInterviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}

	// Set application ID and interviewer ID
	req.ApplicationID = appID
	req.InterviewerID = &employerID

	// Schedule interview
	interview, err := h.appService.ScheduleInterview(ctx, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, ErrFailedOperation, err.Error())
	}

	return utils.CreatedResponse(c, MsgCreatedSuccess, interview)
}

// PUT /applications/:id/interviews/:interview_id - Update interview
func (h *ApplicationHandler) UpdateInterview(c *fiber.Ctx) error {
	ctx := c.Context()

	interviewID, err := strconv.ParseInt(c.Params("interview_id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	// Get interview and update fields
	interview, err := h.appService.GetInterviewDetail(ctx, interviewID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, ErrInterviewNotFound, err.Error())
	}

	// Parse update request
	if err := c.BodyParser(interview); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}

	// TODO: Add proper update logic
	return utils.SuccessResponse(c, MsgUpdatedSuccess, interview)
}

// PATCH /applications/:id/interviews/:interview_id/reschedule - Reschedule interview
func (h *ApplicationHandler) RescheduleInterview(c *fiber.Ctx) error {
	ctx := c.Context()

	interviewID, err := strconv.ParseInt(c.Params("interview_id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	var req application.RescheduleInterviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}

	// Reschedule interview
	interview, err := h.appService.RescheduleInterview(ctx, interviewID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, ErrInterviewConflict, err.Error())
	}

	return utils.SuccessResponse(c, MsgUpdatedSuccess, interview)
}

// PATCH /applications/:id/interviews/:interview_id/complete - Complete interview
func (h *ApplicationHandler) CompleteInterview(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	interviewID, err := strconv.ParseInt(c.Params("interview_id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	var req application.CompleteInterviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}

	// Set completed by
	req.CompletedBy = employerID

	// Complete interview
	interview, err := h.appService.CompleteInterview(ctx, interviewID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, MsgOperationSuccess, interview)
}

// DELETE /applications/:id/interviews/:interview_id - Cancel interview
func (h *ApplicationHandler) CancelInterview(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	interviewID, err := strconv.ParseInt(c.Params("interview_id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	reason := c.Query("reason", "")

	// Cancel interview
	if err := h.appService.CancelInterview(ctx, interviewID, employerID, reason); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, MsgOperationSuccess, nil)
}

// PATCH /applications/:id/bookmark - Toggle bookmark
func (h *ApplicationHandler) Bookmark(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	// Toggle bookmark
	if err := h.appService.ToggleBookmark(ctx, appID, employerID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, MsgOperationSuccess, nil)
}

// PATCH /applications/:id/viewed - Mark as viewed
func (h *ApplicationHandler) MarkViewed(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	// Mark as viewed
	if err := h.appService.MarkAsViewed(ctx, appID, employerID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, MsgOperationSuccess, nil)
}
