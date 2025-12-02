package applicationhandler

import (
	"strconv"

	"keerja-backend/internal/domain/application"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *ApplicationHandler) AddNote(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	appID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	var req application.AddNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}

	req.ApplicationID = appID
	req.AuthorID = employerID

	note, err := h.appService.AddNote(ctx, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrFailedOperation, err.Error())
	}

	return utils.CreatedResponse(c, common.MsgCreatedSuccess, note)
}

func (h *ApplicationHandler) UpdateNote(c *fiber.Ctx) error {
	ctx := c.Context()

	noteID, err := strconv.ParseInt(c.Params("note_id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	var req application.UpdateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}

	note, err := h.appService.UpdateNote(ctx, noteID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, common.MsgUpdatedSuccess, note)
}
