package applicationhandler

import (
	"strconv"

	"keerja-backend/internal/domain/application"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *ApplicationHandler) ScheduleInterview(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	var req application.ScheduleInterviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}

	req.ApplicationID = appID
	req.InterviewerID = &employerID

	interview, err := h.appService.ScheduleInterview(ctx, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrFailedOperation, err.Error())
	}

	return utils.CreatedResponse(c, common.MsgCreatedSuccess, interview)
}

func (h *ApplicationHandler) UpdateInterview(c *fiber.Ctx) error {
	ctx := c.Context()

	interviewID, err := strconv.ParseInt(c.Params("interview_id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	interview, err := h.appService.GetInterviewDetail(ctx, interviewID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, common.ErrInterviewNotFound, err.Error())
	}

	if err := c.BodyParser(interview); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}

	// TODO: Add proper update logic
	return utils.SuccessResponse(c, common.MsgUpdatedSuccess, interview)
}

func (h *ApplicationHandler) RescheduleInterview(c *fiber.Ctx) error {
	ctx := c.Context()

	interviewID, err := strconv.ParseInt(c.Params("interview_id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	var req application.RescheduleInterviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}

	interview, err := h.appService.RescheduleInterview(ctx, interviewID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrInterviewConflict, err.Error())
	}

	return utils.SuccessResponse(c, common.MsgUpdatedSuccess, interview)
}

func (h *ApplicationHandler) CompleteInterview(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	interviewID, err := strconv.ParseInt(c.Params("interview_id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	var req application.CompleteInterviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}

	req.CompletedBy = employerID

	interview, err := h.appService.CompleteInterview(ctx, interviewID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, common.MsgOperationSuccess, interview)
}

func (h *ApplicationHandler) CancelInterview(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	interviewID, err := strconv.ParseInt(c.Params("interview_id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	reason := c.Query("reason", "")

	if err := h.appService.CancelInterview(ctx, interviewID, employerID, reason); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, common.MsgOperationSuccess, nil)
}
