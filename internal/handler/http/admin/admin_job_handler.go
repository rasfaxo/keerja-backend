package admin

import (
	"strconv"

	"keerja-backend/internal/domain/admin"
	"keerja-backend/internal/domain/job"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AdminJobHandler struct {
	adminJobService admin.AdminJobService
}

func NewAdminJobHandler(adminJobService admin.AdminJobService) *AdminJobHandler {
	return &AdminJobHandler{adminJobService: adminJobService}
}

func (h *AdminJobHandler) ApproveJob(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	var req job.ApproveJobRequest
	_ = c.BodyParser(&req)
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errs)
	}

	approvedJob, err := h.adminJobService.ApproveJob(ctx, id)
	if err != nil {
		errMsg := err.Error()

		if errMsg == "job not found: record not found" {
			return utils.NotFoundResponse(c, "Job not found")
		}

		if errMsg == "only jobs with pending_review status can be approved (current status: draft)" ||
			errMsg == "only jobs with pending_review status can be approved (current status: published)" ||
			errMsg == "only jobs with pending_review status can be approved (current status: closed)" ||
			errMsg == "only jobs with pending_review status can be approved (current status: expired)" {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Job status conflict", errMsg)
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to approve job", err.Error())
	}

	resp := mapper.ToJobDetailResponse(approvedJob.(*job.Job))
	return utils.SuccessResponse(c, "Job approved successfully", resp)
}

func (h *AdminJobHandler) RejectJob(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	var req job.RejectJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errs)
	}

	rejectedJob, err := h.adminJobService.RejectJob(ctx, id, req.Reason)
	if err != nil {
		errMsg := err.Error()

		if errMsg == "job not found: record not found" {
			return utils.NotFoundResponse(c, "Job not found")
		}

		if errMsg == "only jobs with pending_review status can be rejected (current status: draft)" ||
			errMsg == "only jobs with pending_review status can be rejected (current status: published)" ||
			errMsg == "only jobs with pending_review status can be rejected (current status: closed)" ||
			errMsg == "only jobs with pending_review status can be rejected (current status: expired)" {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Job status conflict", errMsg)
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to reject job", err.Error())
	}

	resp := mapper.ToJobDetailResponse(rejectedJob.(*job.Job))
	return utils.SuccessResponse(c, "Job rejected and reverted to draft status", resp)
}
