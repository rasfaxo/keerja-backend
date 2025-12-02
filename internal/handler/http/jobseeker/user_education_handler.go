package userhandler

import (
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/helpers"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// UserEducationHandler handles user education operations
type UserEducationHandler struct {
	userService user.UserService
}

// NewUserEducationHandler creates a new instance of UserEducationHandler
func NewUserEducationHandler(userService user.UserService) *UserEducationHandler {
	return &UserEducationHandler{
		userService: userService,
	}
}

func (h *UserEducationHandler) GetEducations(c *fiber.Ctx) error {
	usr, err := helpers.GetProfile(c, h.userService)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, err.Error())
	}

	educations := mapper.MapEntities(usr.Educations, mapper.ToUserEducationResponse)
	return utils.SuccessResponse(c, common.MsgOperationSuccess, educations)
}

func (h *UserEducationHandler) AddEducation(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.AddEducationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrInvalidRequest, err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errors)
	}

	req.InstitutionName = utils.SanitizeIfNonEmpty(req.InstitutionName)
	req.Major = utils.SanitizePtr(req.Major)
	req.DegreeLevel = utils.SanitizePtr(req.DegreeLevel)
	req.Activities = utils.SanitizePtr(req.Activities)
	req.Description = utils.SanitizePtr(req.Description)

	domainReq := &user.AddEducationRequest{
		InstitutionName: req.InstitutionName,
		Major:           req.Major,
		DegreeLevel:     req.DegreeLevel,
		StartYear:       req.StartYear,
		EndYear:         req.EndYear,
		GPA:             req.GPA,
		Activities:      req.Activities,
		Description:     req.Description,
		IsCurrent:       req.IsCurrent,
	}

	education, err := h.userService.AddEducation(ctx, userID, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, err.Error())
	}

	response := mapper.ToEducationResponse(education)
	return utils.CreatedResponse(c, common.MsgOperationSuccess, response)
}

func (h *UserEducationHandler) UpdateEducation(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	educationID, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrInvalidRequest, err.Error())
	}

	var req request.UpdateEducationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrInvalidRequest, err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errors)
	}

	req.InstitutionName = utils.SanitizePtr(req.InstitutionName)
	req.Major = utils.SanitizePtr(req.Major)
	req.DegreeLevel = utils.SanitizePtr(req.DegreeLevel)
	req.Activities = utils.SanitizePtr(req.Activities)
	req.Description = utils.SanitizePtr(req.Description)

	domainReq := &user.UpdateEducationRequest{
		InstitutionName: req.InstitutionName,
		Major:           req.Major,
		DegreeLevel:     req.DegreeLevel,
		StartYear:       req.StartYear,
		EndYear:         req.EndYear,
		GPA:             req.GPA,
		Activities:      req.Activities,
		Description:     req.Description,
		IsCurrent:       req.IsCurrent,
	}

	if err := h.userService.UpdateEducation(ctx, userID, educationID, domainReq); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update education", err.Error())
	}

	return utils.SuccessResponse(c, "Education updated successfully", nil)
}

func (h *UserEducationHandler) DeleteEducation(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	educationID, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid education ID", err.Error())
	}

	if err := h.userService.DeleteEducation(ctx, userID, educationID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete education", err.Error())
	}

	return utils.SuccessResponse(c, "Education deleted successfully", nil)
}
