package applicationhandler

import (
	"strconv"

	"keerja-backend/internal/domain/application"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *ApplicationHandler) ListByJob(c *fiber.Ctx) error {
	ctx := c.Context()

	jobID, err := strconv.ParseInt(c.Params("job_id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	filter := application.ApplicationFilter{}
	if status := c.Query("status"); status != "" {
		filter.Status = status
	}

	response, err := h.appService.GetJobApplications(ctx, jobID, filter, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrInternalServer, err.Error())
	}

	return utils.SuccessResponse(c, common.MsgFetchedSuccess, response)
}

func (h *ApplicationHandler) SearchApplications(c *fiber.Ctx) error {
	ctx := c.Context()

	var filter application.ApplicationSearchFilter
	if err := c.BodyParser(&filter); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	response, err := h.appService.SearchApplications(ctx, filter, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrInternalServer, err.Error())
	}

	return utils.SuccessResponse(c, common.MsgFetchedSuccess, response)
}

func (h *ApplicationHandler) UpdateStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	type UpdateStatusRequest struct {
		Status string `json:"status" validate:"required"`
		Notes  string `json:"notes"`
	}

	var req UpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}

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
		return utils.BadRequestResponse(c, common.ErrInvalidApplicationStage)
	}

	if updateErr != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrFailedOperation, updateErr.Error())
	}

	return utils.SuccessResponse(c, common.MsgStatusUpdated, nil)
}

func (h *ApplicationHandler) BulkUpdateStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	type BulkUpdateRequest struct {
		ApplicationIDs []int64 `json:"application_ids" validate:"required,min=1"`
		Status         string  `json:"status" validate:"required"`
	}

	var req BulkUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}

	if err := h.appService.BulkUpdateStatus(ctx, req.ApplicationIDs, req.Status, employerID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, common.MsgOperationSuccess, nil)
}

func (h *ApplicationHandler) UpdateStage(c *fiber.Ctx) error {
	return h.UpdateStatus(c)
}

func (h *ApplicationHandler) Bookmark(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	if err := h.appService.ToggleBookmark(ctx, appID, employerID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, common.MsgOperationSuccess, nil)
}

func (h *ApplicationHandler) MarkViewed(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	if err := h.appService.MarkAsViewed(ctx, appID, employerID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, common.MsgOperationSuccess, nil)
}
