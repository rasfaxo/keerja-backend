package master

import (
	"strconv"
	"strings"

	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// CompanySizeHandler handles HTTP requests for company size master data
type CompanySizeHandler struct {
	service master.CompanySizeService
}

// NewCompanySizeHandler creates a new instance of CompanySizeHandler
func NewCompanySizeHandler(service master.CompanySizeService) *CompanySizeHandler {
	return &CompanySizeHandler{
		service: service,
	}
}

// GetAllCompanySizes godoc
// @Summary Get all company sizes
// @Description Retrieve all company size categories with optional filtering by active status
// @Tags Master Data - Company Sizes
// @Accept json
// @Produce json
// @Param active query boolean false "Filter by active status (true/false)"
// @Success 200 {object} utils.Response{data=[]master.CompanySizeResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/meta/company-sizes [get]
func (h *CompanySizeHandler) GetAllCompanySizes(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse query parameters
	activeParam := strings.TrimSpace(c.Query("active"))

	// Determine if we should filter by active status
	var companySizes []master.CompanySizeResponse
	var err error

	if activeParam == "true" {
		// Get only active company sizes
		companySizes, err = h.service.GetActive(ctx)
	} else {
		// Get all company sizes (including inactive)
		companySizes, err = h.service.GetAll(ctx)
	}

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve company sizes", err.Error())
	}

	return utils.SuccessResponse(c, "Company sizes retrieved successfully", companySizes)
}

// GetCompanySizeByID godoc
// @Summary Get company size by ID
// @Description Retrieve a single company size category by its ID
// @Tags Master Data - Company Sizes
// @Accept json
// @Produce json
// @Param id path int true "Company Size ID"
// @Success 200 {object} utils.Response{data=master.CompanySizeResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/meta/company-sizes/{id} [get]
func (h *CompanySizeHandler) GetCompanySizeByID(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse ID from path parameter
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid company size ID")
	}

	// Get company size by ID
	companySize, err := h.service.GetByID(ctx, id)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == "company size not found" || strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "Company size not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve company size", err.Error())
	}

	return utils.SuccessResponse(c, "Company size retrieved successfully", companySize)
}
