package http

import (
	"context"
	"strconv"
	"strings"
	"time"

	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// ============================================================================
// Service contract untuk Application 
// ============================================================================
type ApplicationService interface {
	// Candidate
	Apply(ctx context.Context, userID int64, req request.ApplyJobRequest) (*response.ApplicationDetailResponse, error)
	ApplyToJob(ctx context.Context, userID, jobID int64, req request.ApplyJobRequest) (*response.ApplicationDetailResponse, error)
	GetDetail(ctx context.Context, appID int64, requesterID int64) (*response.ApplicationDetailResponse, error)
	Withdraw(ctx context.Context, appID int64, userID int64, req request.WithdrawApplicationRequest) error
	ListMy(ctx context.Context, userID int64, f request.ApplicationFilterRequest) (*response.ApplicationListResponse, any, error)
	UploadDocument(ctx context.Context, actorID int64, appID int64, req request.UploadApplicationDocumentRequest) (*response.ApplicationDocumentResponse, error)
	RateExperience(ctx context.Context, userID int64, appID int64, req request.RateApplicationExperienceRequest) error

	// Employer
	ListByJob(ctx context.Context, employerID int64, jobID int64, f request.ApplicationFilterRequest) (*response.ApplicationListResponse, any, error)
	Search(ctx context.Context, employerID int64, q request.ApplicationSearchRequest) (*response.ApplicationListResponse, any, error)

	UpdateStatus(ctx context.Context, employerID int64, appID int64, req request.UpdateApplicationStatusRequest) error
	BulkUpdateStatus(ctx context.Context, employerID int64, req request.BulkUpdateApplicationsRequest) error

	UpdateStage(ctx context.Context, employerID int64, appID int64, req request.UpdateApplicationStageRequest) error

	AddNote(ctx context.Context, employerID int64, appID int64, req request.AddApplicationNoteRequest) (*response.ApplicationNoteResponse, error)
	UpdateNote(ctx context.Context, employerID int64, appID int64, noteID int64, req request.UpdateApplicationNoteRequest) (*response.ApplicationNoteResponse, error)

	ScheduleInterview(ctx context.Context, employerID int64, appID int64, req request.ScheduleInterviewRequest) (*response.InterviewResponse, error)
	UpdateInterview(ctx context.Context, employerID int64, appID int64, interviewID int64, req request.UpdateInterviewRequest) (*response.InterviewResponse, error)
	RescheduleInterview(ctx context.Context, employerID int64, appID int64, interviewID int64, req request.RescheduleInterviewRequest) (*response.InterviewResponse, error)
	CompleteInterview(ctx context.Context, employerID int64, appID int64, interviewID int64, req request.CompleteInterviewRequest) (*response.InterviewResponse, error)
	CancelInterview(ctx context.Context, employerID int64, appID int64, interviewID int64) error

	Bookmark(ctx context.Context, employerID int64, appID int64, bookmark bool) error
	MarkViewed(ctx context.Context, employerID int64, appID int64, viewed bool) error
}

// ============================================================================
// Handler (diselaraskan dengan pola user_handler: utils.Response + middleware)
// ============================================================================
type ApplicationHandler struct {
	svc ApplicationService
}

func NewApplicationHandler(svc ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{svc: svc}
}

// ------------------------------ Candidate ------------------------------

// POST /applications
func (h *ApplicationHandler) Apply(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.ApplyJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	out, err := h.svc.Apply(ctx, userID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to apply for job", err.Error())
	}
	return utils.CreatedResponse(c, "Application created successfully", out)
}

// POST /jobs/:job_id/apply
func (h *ApplicationHandler) ApplyToJob(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	jobID, err := strconv.ParseInt(c.Params("job_id"), 10, 64)
	if err != nil || jobID <= 0 {
		return utils.BadRequestResponse(c, "Invalid job ID")
	}

	var req request.ApplyJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	// Paksa gunakan job_id dari path
	req.JobID = jobID

	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	out, err := h.svc.ApplyToJob(ctx, userID, jobID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to apply to job", err.Error())
	}
	return utils.CreatedResponse(c, "Application created successfully", out)
}

// GET /applications/my
func (h *ApplicationHandler) GetMyApplications(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var f request.ApplicationFilterRequest
	if err := c.QueryParser(&f); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query params", err.Error())
	}
	if err := utils.ValidateStruct(&f); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}
	f.Page, f.Limit = utils.ValidatePagination(f.Page, f.Limit, 100)

	data, meta, err := h.svc.ListMy(ctx, userID, f)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get applications", err.Error())
	}
	return utils.SuccessResponseWithMeta(c, "Applications retrieved successfully", data, meta)
}

// GET /applications/:id
func (h *ApplicationHandler) GetApplication(c *fiber.Ctx) error {
	ctx := c.Context()
	requesterID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || appID <= 0 {
		return utils.BadRequestResponse(c, "Invalid application ID")
	}

	out, err := h.svc.GetDetail(ctx, appID, requesterID)
	if err != nil {
		low := strings.ToLower(err.Error())
		if strings.Contains(low, "not found") || strings.Contains(low, "record not found") {
			return utils.NotFoundResponse(c, "Application not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get application", err.Error())
	}
	return utils.SuccessResponse(c, "Application retrieved successfully", out)
}

// DELETE /applications/:id
func (h *ApplicationHandler) Withdraw(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || appID <= 0 {
		return utils.BadRequestResponse(c, "Invalid application ID")
	}

	var req request.WithdrawApplicationRequest
	_ = c.BodyParser(&req) // body opsional

	if strings.TrimSpace(req.Reason) != "" {
		if err := utils.ValidateStruct(&req); err != nil {
			errs := utils.FormatValidationErrors(err)
			return utils.ValidationErrorResponse(c, "Validation failed", errs)
		}
	}

	if err := h.svc.Withdraw(ctx, appID, userID, req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to withdraw application", err.Error())
	}
	return utils.SuccessResponse(c, "Application withdrawn successfully", fiber.Map{"withdrawn": true})
}

// POST /applications/:id/documents
func (h *ApplicationHandler) UploadDocument(c *fiber.Ctx) error {
	ctx := c.Context()
	actorID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || appID <= 0 {
		return utils.BadRequestResponse(c, "Invalid application ID")
	}

	var req request.UploadApplicationDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	out, err := h.svc.UploadDocument(ctx, actorID, appID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload document", err.Error())
	}
	return utils.CreatedResponse(c, "Document uploaded successfully", out)
}

// POST /applications/:id/experience
func (h *ApplicationHandler) RateExperience(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || appID <= 0 {
		return utils.BadRequestResponse(c, "Invalid application ID")
	}

	var req request.RateApplicationExperienceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	if err := h.svc.RateExperience(ctx, userID, appID, req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to rate experience", err.Error())
	}
	return utils.SuccessResponse(c, "Experience rated successfully", fiber.Map{"rated": true})
}

// ------------------------------ Employer ------------------------------

// GET /jobs/:id/applications
func (h *ApplicationHandler) ListByJob(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	jobID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || jobID <= 0 {
		return utils.BadRequestResponse(c, "Invalid job ID")
	}

	var f request.ApplicationFilterRequest
	if err := c.QueryParser(&f); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query params", err.Error())
	}
	if err := utils.ValidateStruct(&f); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}
	f.Page, f.Limit = utils.ValidatePagination(f.Page, f.Limit, 100)

	data, meta, err := h.svc.ListByJob(ctx, employerID, jobID, f)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get applications", err.Error())
	}
	return utils.SuccessResponseWithMeta(c, "Applications retrieved successfully", data, meta)
}

// POST /applications/search
func (h *ApplicationHandler) SearchApplications(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	var q request.ApplicationSearchRequest
	if err := c.BodyParser(&q); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 || q.Limit > 100 {
		q.Limit = 10
	}
	if err := utils.ValidateStruct(&q); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	data, meta, err := h.svc.Search(ctx, employerID, q)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to search applications", err.Error())
	}
	return utils.SuccessResponseWithMeta(c, "Applications searched successfully", data, meta)
}

// PATCH /applications/:id/status
func (h *ApplicationHandler) UpdateStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || appID <= 0 {
		return utils.BadRequestResponse(c, "Invalid application ID")
	}

	var req request.UpdateApplicationStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if strings.EqualFold(req.Status, "rejected") && strings.TrimSpace(req.RejectionReason) == "" {
		return utils.BadRequestResponse(c, "rejection_reason is required when status=rejected")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	if err := h.svc.UpdateStatus(ctx, employerID, appID, req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update application status", err.Error())
	}
	return utils.SuccessResponse(c, "Application status updated successfully", fiber.Map{"updated": true})
}

// POST /applications/bulk-status
func (h *ApplicationHandler) BulkUpdateStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	var req request.BulkUpdateApplicationsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	if err := h.svc.BulkUpdateStatus(ctx, employerID, req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to bulk update applications", err.Error())
	}
	return utils.SuccessResponse(c, "Applications updated successfully", fiber.Map{"bulk_updated": true})
}

// PATCH /applications/:id/stage
func (h *ApplicationHandler) UpdateStage(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || appID <= 0 {
		return utils.BadRequestResponse(c, "Invalid application ID")
	}

	var req request.UpdateApplicationStageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	if err := h.svc.UpdateStage(ctx, employerID, appID, req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update stage", err.Error())
	}
	return utils.SuccessResponse(c, "Stage updated successfully", fiber.Map{"stage_updated": true})
}

// POST /applications/:id/notes
func (h *ApplicationHandler) AddNote(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || appID <= 0 {
		return utils.BadRequestResponse(c, "Invalid application ID")
	}

	var req request.AddApplicationNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	out, err := h.svc.AddNote(ctx, employerID, appID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to add note", err.Error())
	}
	return utils.CreatedResponse(c, "Note added successfully", out)
}

// PATCH /applications/:id/notes/:note_id
func (h *ApplicationHandler) UpdateNote(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || appID <= 0 {
		return utils.BadRequestResponse(c, "Invalid application ID")
	}
	noteID, err := strconv.ParseInt(c.Params("note_id"), 10, 64)
	if err != nil || noteID <= 0 {
		return utils.BadRequestResponse(c, "Invalid note ID")
	}

	var req request.UpdateApplicationNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	out, err := h.svc.UpdateNote(ctx, employerID, appID, noteID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update note", err.Error())
	}
	return utils.SuccessResponse(c, "Note updated successfully", out)
}

// POST /applications/:id/interviews
func (h *ApplicationHandler) ScheduleInterview(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || appID <= 0 {
		return utils.BadRequestResponse(c, "Invalid application ID")
	}

	var req request.ScheduleInterviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	// Paksa gunakan :id
	req.ApplicationID = appID
	// gtfield=Now tidak berjalan di validator → cek manual
	if req.ScheduledAt.Before(time.Now().Add(-1 * time.Minute)) {
		return utils.BadRequestResponse(c, "scheduled_at must be in the future")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	out, err := h.svc.ScheduleInterview(ctx, employerID, appID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to schedule interview", err.Error())
	}
	return utils.CreatedResponse(c, "Interview scheduled successfully", out)
}

// PATCH /applications/:id/interviews/:interview_id
func (h *ApplicationHandler) UpdateInterview(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || appID <= 0 {
		return utils.BadRequestResponse(c, "Invalid application ID")
	}
	invID, err := strconv.ParseInt(c.Params("interview_id"), 10, 64)
	if err != nil || invID <= 0 {
		return utils.BadRequestResponse(c, "Invalid interview ID")
	}

	var req request.UpdateInterviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if req.ScheduledAt != nil && req.ScheduledAt.Before(time.Now().Add(-1*time.Minute)) {
		return utils.BadRequestResponse(c, "scheduled_at must be in the future")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	out, err := h.svc.UpdateInterview(ctx, employerID, appID, invID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update interview", err.Error())
	}
	return utils.SuccessResponse(c, "Interview updated successfully", out)
}

// POST /applications/:id/interviews/:interview_id/reschedule
func (h *ApplicationHandler) RescheduleInterview(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || appID <= 0 {
		return utils.BadRequestResponse(c, "Invalid application ID")
	}
	invID, err := strconv.ParseInt(c.Params("interview_id"), 10, 64)
	if err != nil || invID <= 0 {
		return utils.BadRequestResponse(c, "Invalid interview ID")
	}

	var req request.RescheduleInterviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if req.ScheduledAt.Before(time.Now().Add(-1 * time.Minute)) {
		return utils.BadRequestResponse(c, "scheduled_at must be in the future")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	out, err := h.svc.RescheduleInterview(ctx, employerID, appID, invID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to reschedule interview", err.Error())
	}
	return utils.SuccessResponse(c, "Interview rescheduled successfully", out)
}

// POST /applications/:id/interviews/:interview_id/complete
func (h *ApplicationHandler) CompleteInterview(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || appID <= 0 {
		return utils.BadRequestResponse(c, "Invalid application ID")
	}
	invID, err := strconv.ParseInt(c.Params("interview_id"), 10, 64)
	if err != nil || invID <= 0 {
		return utils.BadRequestResponse(c, "Invalid interview ID")
	}

	var req request.CompleteInterviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	out, err := h.svc.CompleteInterview(ctx, employerID, appID, invID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to complete interview", err.Error())
	}
	return utils.SuccessResponse(c, "Interview completed successfully", out)
}

// DELETE /applications/:id/interviews/:interview_id
func (h *ApplicationHandler) CancelInterview(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || appID <= 0 {
		return utils.BadRequestResponse(c, "Invalid application ID")
	}
	invID, err := strconv.ParseInt(c.Params("interview_id"), 10, 64)
	if err != nil || invID <= 0 {
		return utils.BadRequestResponse(c, "Invalid interview ID")
	}

	if err := h.svc.CancelInterview(ctx, employerID, appID, invID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to cancel interview", err.Error())
	}
	return utils.SuccessResponse(c, "Interview cancelled successfully", fiber.Map{"cancelled": true})
}

// POST /applications/:id/bookmark?bookmark=false
func (h *ApplicationHandler) Bookmark(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || appID <= 0 {
		return utils.BadRequestResponse(c, "Invalid application ID")
	}

	bookmark := true
	if v := c.Query("bookmark"); v != "" {
		bookmark = strings.ToLower(v) != "false"
	}

	if err := h.svc.Bookmark(ctx, employerID, appID, bookmark); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update bookmark", err.Error())
	}
	return utils.SuccessResponse(c, "Bookmark updated successfully", fiber.Map{"bookmarked": bookmark})
}

// POST /applications/:id/viewed?viewed=false
func (h *ApplicationHandler) MarkViewed(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || appID <= 0 {
		return utils.BadRequestResponse(c, "Invalid application ID")
	}

	viewed := true
	if v := c.Query("viewed"); v != "" {
		viewed = strings.ToLower(v) != "false"
	}

	if err := h.svc.MarkViewed(ctx, employerID, appID, viewed); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update viewed status", err.Error())
	}
	return utils.SuccessResponse(c, "Viewed status updated successfully", fiber.Map{"viewed": viewed})
}
