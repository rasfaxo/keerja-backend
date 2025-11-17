package http

import (
	"strconv"

	"keerja-backend/internal/domain/admin"
	"keerja-backend/internal/domain/job"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// AdminHandler handles admin operations like job approval/rejection
type AdminHandler struct {
	adminJobService admin.AdminJobService
}

func NewAdminHandler(adminJobService admin.AdminJobService) *AdminHandler {
	return &AdminHandler{adminJobService: adminJobService}
}

// ApproveJob godoc
// @Summary Approve a pending job posting
// @Description Approve a pending job posting (admin only)
// @Tags admin-jobs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Job ID"
// @Param request body job.ApproveJobRequest true "Approval request"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/jobs/{id}/approve [patch]
func (h *AdminHandler) ApproveJob(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	var req job.ApproveJobRequest
	_ = c.BodyParser(&req) // optional fields
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, ErrValidationFailed, errs)
	}

	approvedJob, err := h.adminJobService.ApproveJob(ctx, id)
	if err != nil {
		// Handle specific error cases
		errMsg := err.Error()

		// 404 Not Found: Job not found
		if errMsg == "job not found: record not found" {
			return utils.NotFoundResponse(c, "Job not found")
		}

		// 409 Conflict: Job not in pending_review status
		if errMsg == "only jobs with pending_review status can be approved (current status: draft)" ||
			errMsg == "only jobs with pending_review status can be approved (current status: published)" ||
			errMsg == "only jobs with pending_review status can be approved (current status: closed)" ||
			errMsg == "only jobs with pending_review status can be approved (current status: expired)" {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Job status conflict", errMsg)
		}

		// 500 Internal Server Error: Other errors
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to approve job", err.Error())
	}

	resp := mapper.ToJobDetailResponse(approvedJob.(*job.Job))
	return utils.SuccessResponse(c, "Job approved successfully", resp)
}

// RejectJob godoc
// @Summary Reject a pending job posting
// @Description Reject a pending job posting (admin only)
// @Tags admin-jobs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Job ID"
// @Param request body job.RejectJobRequest true "Rejection request"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/jobs/{id}/reject [patch]
func (h *AdminHandler) RejectJob(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	var req job.RejectJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, ErrValidationFailed, errs)
	}

	rejectedJob, err := h.adminJobService.RejectJob(ctx, id, req.Reason)
	if err != nil {
		// Handle specific error cases
		errMsg := err.Error()

		// 404 Not Found: Job not found
		if errMsg == "job not found: record not found" {
			return utils.NotFoundResponse(c, "Job not found")
		}

		// 409 Conflict: Job not in pending_review status
		if errMsg == "only jobs with pending_review status can be rejected (current status: draft)" ||
			errMsg == "only jobs with pending_review status can be rejected (current status: published)" ||
			errMsg == "only jobs with pending_review status can be rejected (current status: closed)" ||
			errMsg == "only jobs with pending_review status can be rejected (current status: expired)" {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Job status conflict", errMsg)
		}

		// 500 Internal Server Error: Other errors
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to reject job", err.Error())
	}

	resp := mapper.ToJobDetailResponse(rejectedJob.(*job.Job))
	return utils.SuccessResponse(c, "Job rejected and reverted to draft status", resp)
}
