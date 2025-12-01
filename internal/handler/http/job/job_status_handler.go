package jobhandler

import (
	"strings"
	"time"

	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *JobHandler) PublishJob(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	id, err := utils.ParseIDParam(c, "id")
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	var req request.PublishJobRequest
	_ = c.BodyParser(&req)
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errs)
	}

	var expiredAtPtr *time.Time
	if req.ExpiredAt != nil && *req.ExpiredAt != "" {
		expiredAt, err := utils.ParseOptionalDateTime(req.ExpiredAt)
		if err != nil {
			return utils.BadRequestResponse(c, common.ErrInvalidDateFormat)
		}
		if expiredAt != nil {
			if err := utils.MustBeFutureTime(*expiredAt); err != nil {
				return utils.BadRequestResponse(c, common.ErrFutureDateRequired)
			}
			expiredAtPtr = expiredAt
		}
	}

	if err := h.jobService.PublishJob(ctx, id, employerID, expiredAtPtr); err != nil {
		errMsg := err.Error()

		if errMsg == "company is not verified yet" {
			return utils.ForbiddenResponse(c, common.ErrCompanyNotVerified)
		}

		if errMsg == "job not found: record not found" || errMsg == "not authorized to manage this job" {
			return utils.NotFoundResponse(c, common.ErrJobNotFound)
		}

		if errMsg == "job is already pending review" || errMsg == "job is already published" {
			return utils.ErrorResponse(c, fiber.StatusConflict, common.ErrConflict, errMsg)
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrInternalServer, err.Error())
	}

	jobObj, err := h.jobService.GetJob(ctx, id)
	if err != nil {
		return utils.SuccessResponse(c, "Job publish initiated", fiber.Map{"status": "unknown"})
	}

	switch jobObj.Status {
	case "published":
		return utils.SuccessResponse(c, "Job published successfully", fiber.Map{
			"status":  "published",
			"message": "Your job is now live",
		})
	case "pending_review", "in_review":
		return utils.SuccessResponse(c, "Job submitted for review successfully", fiber.Map{
			"status":  "pending_review",
			"message": "Your job has been submitted and is waiting for admin approval",
		})
	default:
		return utils.SuccessResponse(c, common.MsgStatusUpdated, fiber.Map{"status": jobObj.Status})
	}
}

func (h *JobHandler) CloseJob(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	id, err := utils.ParseIDParam(c, "id")
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	var req request.CloseJobRequest
	_ = c.BodyParser(&req)
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errs)
	}

	if err := h.jobService.CloseJob(ctx, id, employerID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrInternalServer, err.Error())
	}
	return utils.SuccessResponse(c, common.MsgUpdatedSuccess, fiber.Map{"closed": true})
}

func (h *JobHandler) InactivateJob(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	id, err := utils.ParseIDParam(c, "id")
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	if employerID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorized, "userID not found in context")
	}

	if err := h.jobService.InactivateJob(ctx, id, employerID); err != nil {
		if err.Error() == "you do not have permission to modify this job" {
			return utils.ForbiddenResponse(c, common.ErrNotJobOwner)
		}
		if strings.Contains(err.Error(), "job not found") {
			return utils.NotFoundResponse(c, common.ErrJobNotFound)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrInternalServer, err.Error())
	}

	return utils.SuccessResponse(c, common.MsgOperationSuccess, fiber.Map{"id": id})
}
