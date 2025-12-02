package applicationhandler

import (
	"strconv"

	"keerja-backend/internal/domain/application"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *ApplicationHandler) Apply(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req application.ApplyJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}

	req.UserID = userID
	req.CoverLetter = utils.SanitizeHTML(req.CoverLetter)
	req.Source = utils.SanitizeString(req.Source)

	app, err := h.appService.ApplyForJob(ctx, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrApplicationNotFound, err.Error())
	}

	return utils.CreatedResponse(c, common.MsgApplicationSubmit, app)
}

func (h *ApplicationHandler) ApplyToJob(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	jobID, err := strconv.ParseInt(c.Params("job_id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	var req application.ApplyJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}

	req.JobID = jobID
	req.UserID = userID
	req.CoverLetter = utils.SanitizeHTML(req.CoverLetter)
	req.Source = utils.SanitizeString(req.Source)

	app, err := h.appService.ApplyForJob(ctx, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrAlreadyApplied, err.Error())
	}

	return utils.CreatedResponse(c, common.MsgApplicationSubmit, app)
}

func (h *ApplicationHandler) GetMyApplications(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	filter := application.ApplicationFilter{}
	if status := c.Query("status"); status != "" {
		filter.Status = status
	}
	if jobIDStr := c.Query("job_id"); jobIDStr != "" {
		if jobID, err := strconv.ParseInt(jobIDStr, 10, 64); err == nil {
			filter.JobID = jobID
		}
	}

	response, err := h.appService.GetMyApplications(ctx, userID, filter, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrInternalServer, err.Error())
	}

	return utils.SuccessResponse(c, common.MsgFetchedSuccess, response)
}

func (h *ApplicationHandler) GetApplication(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	response, err := h.appService.GetApplicationDetail(ctx, appID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, common.ErrApplicationNotFound, err.Error())
	}

	return utils.SuccessResponse(c, common.MsgFetchedSuccess, response)
}

func (h *ApplicationHandler) Withdraw(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	if err := h.appService.WithdrawApplication(ctx, appID, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrCannotWithdraw, err.Error())
	}

	return utils.SuccessResponse(c, common.MsgOperationSuccess, nil)
}

func (h *ApplicationHandler) UploadDocument(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	var req application.UploadDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}

	req.ApplicationID = appID
	req.UserID = userID

	doc, err := h.appService.UploadApplicationDocument(ctx, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrFileUploadFailed, err.Error())
	}

	return utils.CreatedResponse(c, common.MsgUploadSuccess, doc)
}

func (h *ApplicationHandler) RateExperience(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	type RateRequest struct {
		Rating  int    `json:"rating" validate:"required,min=1,max=5"`
		Comment string `json:"comment"`
	}

	var req RateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}

	// TODO: Implement rate experience in domain service
	_ = userID
	_ = appID
	_ = req

	return utils.ErrorResponse(c, fiber.StatusNotImplemented, common.ErrFailedOperation, "Rating feature not yet implemented")
}
