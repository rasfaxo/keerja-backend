package userhandler

import (
	"fmt"

	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/helpers"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// UserSkillHandler handles user skill operations
type UserSkillHandler struct {
	userService user.UserService
}

// NewUserSkillHandler creates a new instance of UserSkillHandler
func NewUserSkillHandler(userService user.UserService) *UserSkillHandler {
	return &UserSkillHandler{
		userService: userService,
	}
}

func (h *UserSkillHandler) GetSkills(c *fiber.Ctx) error {
	usr, err := helpers.GetProfile(c, h.userService)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, err.Error())
	}

	skills := mapper.MapEntities(usr.Skills, mapper.ToUserSkillResponse)
	return utils.SuccessResponse(c, common.MsgOperationSuccess, skills)
}

func (h *UserSkillHandler) AddSkills(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.AddUserSkillsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	domainSkills := make([]user.AddUserSkillRequest, len(req.Skills))
	for i, skill := range req.Skills {
		if skill.SkillID == nil && skill.SkillName == "" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed",
				fmt.Sprintf("Skill #%d: Either skill_id or skill_name must be provided", i+1))
		}

		domainSkills[i] = user.AddUserSkillRequest{
			SkillID:           skill.SkillID,
			SkillName:         skill.SkillName,
			ProficiencyLevel:  skill.ProficiencyLevel,
			YearsOfExperience: skill.YearsOfExperience,
		}
	}

	domainReq := &user.AddUserSkillsRequest{
		Skills: domainSkills,
	}

	addedSkills, err := h.userService.AddSkills(ctx, userID, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to add skills", err.Error())
	}

	skillsResponse := make([]map[string]any, len(addedSkills))
	for i, skill := range addedSkills {
		skillsResponse[i] = map[string]any{
			"id":               skill.ID,
			"skill_name":       skill.SkillName,
			"skill_level":      skill.SkillLevel,
			"years_experience": skill.YearsExperience,
		}
	}

	return utils.CreatedResponse(c,
		fmt.Sprintf("Successfully added %d skills", len(addedSkills)),
		map[string]any{
			"skills": skillsResponse,
			"total":  len(addedSkills),
		},
	)
}

func (h *UserSkillHandler) DeleteSkill(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	skillID, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid skill ID", err.Error())
	}

	if err := h.userService.DeleteSkill(ctx, userID, skillID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete skill", err.Error())
	}

	return utils.SuccessResponse(c, "Skill deleted successfully", nil)
}
