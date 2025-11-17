package master

import (
	"strconv"
	"strings"

	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// IndustryHandler handles HTTP requests for industry master data
type IndustryHandler struct {
	service master.IndustryService
}

// NewIndustryHandler creates a new instance of IndustryHandler
func NewIndustryHandler(service master.IndustryService) *IndustryHandler {
	return &IndustryHandler{
		service: service,
	}
}

// GetAllIndustries godoc
// @Summary Get all industries
// @Description Retrieve all industries with optional filtering by active status and search query
// @Tags Master Data - Industries
// @Accept json
// @Produce json
// @Param active query boolean false "Filter by active status (true/false)"
// @Param search query string false "Search industries by name (min 2 characters)"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/meta/industries [get]
func (h *IndustryHandler) GetAllIndustries(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse query parameters
	activeParam := strings.TrimSpace(c.Query("active"))
	search := strings.TrimSpace(c.Query("search"))

	// Determine if we should filter by active status
	var industries []master.IndustryResponse
	var err error

	if activeParam == "true" {
		// Get only active industries
		industries, err = h.service.GetActive(ctx, search)
	} else {
		// Get all industries (including inactive)
		industries, err = h.service.GetAll(ctx, search)
	}

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve industries", err.Error())
	}

	return utils.SuccessResponse(c, "Industries retrieved successfully", industries)
}

// GetIndustryByID godoc
// @Summary Get industry by ID
// @Description Retrieve a single industry by its ID
// @Tags Master Data - Industries
// @Accept json
// @Produce json
// @Param id path int true "Industry ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/meta/industries/{id} [get]
func (h *IndustryHandler) GetIndustryByID(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse ID from path parameter
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid industry ID")
	}

	// Get industry by ID
	industry, err := h.service.GetByID(ctx, id)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == "industry not found" || strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "Industry not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve industry", err.Error())
	}

	return utils.SuccessResponse(c, "Industry retrieved successfully", industry)
}
