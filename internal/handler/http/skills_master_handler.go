package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/utils"
)

// SkillsMasterHandler handles HTTP requests for skills master data
type SkillsMasterHandler struct {
	service master.SkillsMasterService
}

// NewSkillsMasterHandler creates a new skills master handler
func NewSkillsMasterHandler(service master.SkillsMasterService) *SkillsMasterHandler {
	return &SkillsMasterHandler{
		service: service,
	}
}

// GetAllSkills retrieves all skills with optional filters
// GET /api/v1/skills
func (h *SkillsMasterHandler) GetAllSkills(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse query parameters for filtering
	filter := &master.SkillsFilter{
		Search:          c.Query("search", ""),
		SkillType:       c.Query("skill_type", ""),
		DifficultyLevel: c.Query("difficulty_level", ""),
		Page:            1,
		PageSize:        20,
		SortBy:          c.Query("sort_by", "name"),
		SortOrder:       c.Query("sort_order", "ASC"),
	}

	// Parse pagination (use QueryInt with defaults)
	filter.Page = c.QueryInt("page", 1)
	filter.PageSize = c.QueryInt("page_size", 20)

	// Parse category_id if provided (as query param)
	if v := c.QueryInt("category_id", 0); v > 0 {
		cid := int64(v)
		filter.CategoryID = &cid
	}

	// Parse is_active if provided
	if isActive := c.Query("is_active"); isActive != "" {
		if ia, err := strconv.ParseBool(isActive); err == nil {
			filter.IsActive = &ia
		}
	}

	// Parse popularity range
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

	// Get skills from service
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

// SearchSkills searches skills by query string
// GET /api/v1/skills/search
func (h *SkillsMasterHandler) SearchSkills(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get search query
	query := c.Query("q", "")
	if query == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest,
			"Search query is required",
			"Please provide a search query using the 'q' parameter",
		)
	}

	// Parse pagination
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)

	// Search skills
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

// GetSkillsByType retrieves skills by type
// GET /api/v1/skills/type/:type
func (h *SkillsMasterHandler) GetSkillsByType(c *fiber.Ctx) error {
	ctx := c.Context()

	skillType := c.Params("type")
	if skillType == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest,
			"Skill type is required",
			"Please provide a valid skill type",
		)
	}

	// Validate skill type
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

// GetSkillByID retrieves a single skill by ID
// GET /api/v1/skills/:id
func (h *SkillsMasterHandler) GetSkillByID(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse skill ID
	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest,
			"Invalid skill ID",
			"Skill ID must be a valid number",
		)
	}

	// Get skill from service
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

// GetSkillsByIDs retrieves multiple skills by IDs
// POST /api/v1/skills/by-ids
// Body: { "ids": [1, 2, 3] }
func (h *SkillsMasterHandler) GetSkillsByIDs(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
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

	// Limit to max 50 IDs to prevent abuse
	if len(req.IDs) > 50 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest,
			"Too many IDs",
			"Maximum 50 skill IDs allowed per request",
		)
	}

	// Get skills from service
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

// GetRecommendedSkills retrieves recommended skills based on context
// GET /api/v1/skills/recommended
func (h *SkillsMasterHandler) GetRecommendedSkills(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse optional filters
	skillType := c.Query("skill_type", "")
	limit := c.QueryInt("limit", 10)

	// Build filter for popular/recommended skills
	filter := &master.SkillsFilter{
		SkillType: skillType,
		IsActive:  utils.BoolPtr(true),
		Page:      1,
		PageSize:  limit,
		SortBy:    "popularity_score",
		SortOrder: "DESC",
	}

	// Get skills from service
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
