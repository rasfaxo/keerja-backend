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

// UserExperienceHandler handles user work experience operations
type UserExperienceHandler struct {
	userService user.UserService
}

// NewUserExperienceHandler creates a new instance of UserExperienceHandler
func NewUserExperienceHandler(userService user.UserService) *UserExperienceHandler {
	return &UserExperienceHandler{
		userService: userService,
	}
}

func (h *UserExperienceHandler) GetExperiences(c *fiber.Ctx) error {
	usr, err := helpers.GetProfile(c, h.userService)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, err.Error())
	}

	experiences := mapper.MapEntities(usr.Experiences, mapper.ToUserExperienceResponse)
	return utils.SuccessResponse(c, common.MsgOperationSuccess, experiences)
}

func (h *UserExperienceHandler) AddExperience(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.AddExperienceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	req.CompanyName = utils.SanitizeIfNonEmpty(req.CompanyName)
	req.PositionTitle = utils.SanitizeIfNonEmpty(req.PositionTitle)
	req.Industry = utils.SanitizePtr(req.Industry)
	req.EmploymentType = utils.SanitizePtr(req.EmploymentType)
	req.Description = utils.SanitizePtr(req.Description)
	req.Achievements = utils.SanitizePtr(req.Achievements)
	req.LocationCity = utils.SanitizePtr(req.LocationCity)
	req.LocationCountry = utils.SanitizePtr(req.LocationCountry)

	domainReq := &user.AddExperienceRequest{
		CompanyName:     req.CompanyName,
		PositionTitle:   req.PositionTitle,
		Industry:        req.Industry,
		EmploymentType:  req.EmploymentType,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		IsCurrent:       req.IsCurrent,
		Description:     req.Description,
		Achievements:    req.Achievements,
		LocationCity:    req.LocationCity,
		LocationCountry: req.LocationCountry,
	}

	experience, err := h.userService.AddExperience(ctx, userID, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to add experience", err.Error())
	}

	response := mapper.ToExperienceResponse(experience)
	return utils.CreatedResponse(c, "Experience added successfully", response)
}

func (h *UserExperienceHandler) UpdateExperience(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	experienceID, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid experience ID", err.Error())
	}

	var req request.UpdateExperienceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	req.CompanyName = utils.SanitizePtr(req.CompanyName)
	req.PositionTitle = utils.SanitizePtr(req.PositionTitle)
	req.Industry = utils.SanitizePtr(req.Industry)
	req.EmploymentType = utils.SanitizePtr(req.EmploymentType)
	req.Description = utils.SanitizePtr(req.Description)
	req.Achievements = utils.SanitizePtr(req.Achievements)
	req.LocationCity = utils.SanitizePtr(req.LocationCity)
	req.LocationCountry = utils.SanitizePtr(req.LocationCountry)

	domainReq := &user.UpdateExperienceRequest{
		CompanyName:     req.CompanyName,
		PositionTitle:   req.PositionTitle,
		Industry:        req.Industry,
		EmploymentType:  req.EmploymentType,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		IsCurrent:       req.IsCurrent,
		Description:     req.Description,
		Achievements:    req.Achievements,
		LocationCity:    req.LocationCity,
		LocationCountry: req.LocationCountry,
	}

	if err := h.userService.UpdateExperience(ctx, userID, experienceID, domainReq); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update experience", err.Error())
	}

	return utils.SuccessResponse(c, "Experience updated successfully", nil)
}

func (h *UserExperienceHandler) DeleteExperience(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	experienceID, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid experience ID", err.Error())
	}

	if err := h.userService.DeleteExperience(ctx, userID, experienceID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete experience", err.Error())
	}

	return utils.SuccessResponse(c, "Experience deleted successfully", nil)
}
