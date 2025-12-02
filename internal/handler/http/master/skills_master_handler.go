package master

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/utils"
)

type SkillsMasterHandler struct {
	service master.SkillsMasterService
}

func NewSkillsMasterHandler(service master.SkillsMasterService) *SkillsMasterHandler {
	return &SkillsMasterHandler{
		service: service,
	}
}

func (h *SkillsMasterHandler) GetAllSkills(c *fiber.Ctx) error {
	ctx := c.Context()

	filter := &master.SkillsFilter{
		Search:          c.Query("search", ""),
		SkillType:       c.Query("skill_type", ""),
		DifficultyLevel: c.Query("difficulty_level", ""),
		Page:            1,
		PageSize:        20,
		SortBy:          c.Query("sort_by", "name"),
		SortOrder:       c.Query("sort_order", "ASC"),
	}

	filter.Page = c.QueryInt("page", 1)
	filter.PageSize = c.QueryInt("page_size", 20)

	if v := c.QueryInt("category_id", 0); v > 0 {
		cid := int64(v)
		filter.CategoryID = &cid
	}

	if isActive := c.Query("is_active"); isActive != "" {
		if ia, err := strconv.ParseBool(isActive); err == nil {
			filter.IsActive = &ia
		}
	}

	if minPop := c.Query("min_popularity"); minPop != "" {
		if mp, err := strconv.ParseFloat(minPop, 64); err == nil {
			filter.MinPopularity = &mp
		}
	}

	if maxPop := c.Query("max_popularity"); maxPop != "" {
		if mp, err := strconv.ParseFloat(maxPop, 64); err == nil {
			filter.MaxPopularity = &mp
		}
	}

	result, err := h.service.GetSkills(ctx, filter)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError,
			"Failed to retrieve skills",
			err.Error(),
		)
	}

	return utils.SuccessResponse(c,
		"Skills retrieved successfully",
		result,
	)
}

func (h *SkillsMasterHandler) SearchSkills(c *fiber.Ctx) error {
	ctx := c.Context()

	query := c.Query("q", "")
	if query == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest,
			"Search query is required",
			"Please provide a search query using the 'q' parameter",
		)
	}

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)

	result, err := h.service.SearchSkills(ctx, query, page, pageSize)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError,
			"Failed to search skills",
			err.Error(),
		)
	}

	return utils.SuccessResponse(c,
		"Skills search completed successfully",
		result,
	)
}

func (h *SkillsMasterHandler) GetSkillsByType(c *fiber.Ctx) error {
	ctx := c.Context()

	skillType := c.Params("type")
	if skillType == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest,
			"Skill type is required",
			"Please provide a valid skill type",
		)
	}

	validTypes := map[string]bool{
		"technical": true,
		"soft":      true,
		"language":  true,
		"tool":      true,
	}

	if !validTypes[skillType] {
		return utils.ErrorResponse(c, fiber.StatusBadRequest,
			"Invalid skill type",
			"Skill type must be one of: technical, soft, language, tool",
		)
	}

	skills, err := h.service.GetSkillsByType(ctx, skillType)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError,
			"Failed to retrieve skills by type",
			err.Error(),
		)
	}

	return utils.SuccessResponse(c,
		"Skills retrieved successfully",
		map[string]interface{}{
			"skills": skills,
			"type":   skillType,
			"total":  len(skills),
		},
	)
}

func (h *SkillsMasterHandler) GetSkillByID(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest,
			"Invalid skill ID",
			"Skill ID must be a valid number",
		)
	}

	skill, err := h.service.GetSkill(ctx, id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound,
			"Skill not found",
			err.Error(),
		)
	}

	return utils.SuccessResponse(c,
		"Skill retrieved successfully",
		skill,
	)
}

func (h *SkillsMasterHandler) GetSkillsByIDs(c *fiber.Ctx) error {
	ctx := c.Context()

	var req struct {
		IDs []int64 `json:"ids" validate:"required,min=1"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
	}

	if len(req.IDs) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest,
			"IDs array cannot be empty",
			"Please provide at least one skill ID",
		)
	}

	if len(req.IDs) > 50 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest,
			"Too many IDs",
			"Maximum 50 skill IDs allowed per request",
		)
	}

	var skills []interface{}
	for _, id := range req.IDs {
		skill, err := h.service.GetSkill(ctx, id)
		if err == nil && skill != nil {
			skills = append(skills, skill)
		}
	}

	return utils.SuccessResponse(c,
		"Skills retrieved successfully",
		map[string]interface{}{
			"skills": skills,
			"total":  len(skills),
		},
	)
}

func (h *SkillsMasterHandler) GetRecommendedSkills(c *fiber.Ctx) error {
	ctx := c.Context()

	skillType := c.Query("skill_type", "")
	limit := c.QueryInt("limit", 10)

	filter := &master.SkillsFilter{
		SkillType: skillType,
		IsActive:  utils.BoolPtr(true),
		Page:      1,
		PageSize:  limit,
		SortBy:    "popularity_score",
		SortOrder: "DESC",
	}

	result, err := h.service.GetSkills(ctx, filter)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError,
			"Failed to retrieve recommended skills",
			err.Error(),
		)
	}

	return utils.SuccessResponse(c,
		"Recommended skills retrieved successfully",
		result,
	)
}
