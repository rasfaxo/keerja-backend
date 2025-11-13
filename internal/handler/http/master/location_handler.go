package master

import (
	"strconv"
	"strings"

	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// LocationHandler handles HTTP requests for location master data (provinces, cities, districts)
type LocationHandler struct {
	provinceService master.ProvinceService
	cityService     master.CityService
	districtService master.DistrictService
}

// NewLocationHandler creates a new instance of LocationHandler
func NewLocationHandler(
	provinceService master.ProvinceService,
	cityService master.CityService,
	districtService master.DistrictService,
) *LocationHandler {
	return &LocationHandler{
		provinceService: provinceService,
		cityService:     cityService,
		districtService: districtService,
	}
}

// ========================================
// PROVINCE ENDPOINTS
// ========================================

// GetAllProvinces godoc
// @Summary Get all provinces
// @Description Retrieve all provinces with optional filtering by active status and search query
// @Tags Master Data - Locations
// @Accept json
// @Produce json
// @Param active query boolean false "Filter by active status (true/false)"
// @Param search query string false "Search provinces by name (min 2 characters)"
// @Success 200 {object} utils.Response{data=[]master.ProvinceResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/meta/locations/provinces [get]
func (h *LocationHandler) GetAllProvinces(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse query parameters
	activeParam := strings.TrimSpace(c.Query("active"))
	search := strings.TrimSpace(c.Query("search"))

	// Determine if we should filter by active status
	var provinces []master.ProvinceResponse
	var err error

	if activeParam == "true" {
		// Get only active provinces
		provinces, err = h.provinceService.GetActive(ctx, search)
	} else {
		// Get all provinces (including inactive)
		provinces, err = h.provinceService.GetAll(ctx, search)
	}

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve provinces", err.Error())
	}

	return utils.SuccessResponse(c, "Provinces retrieved successfully", provinces)
}

// GetProvinceByID godoc
// @Summary Get province by ID
// @Description Retrieve a single province by its ID
// @Tags Master Data - Locations
// @Accept json
// @Produce json
// @Param id path int true "Province ID"
// @Success 200 {object} utils.Response{data=master.ProvinceResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/meta/locations/provinces/{id} [get]
func (h *LocationHandler) GetProvinceByID(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse ID from path parameter
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid province ID")
	}

	// Get province by ID
	province, err := h.provinceService.GetByID(ctx, id)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == "province not found" || strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "Province not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve province", err.Error())
	}

	return utils.SuccessResponse(c, "Province retrieved successfully", province)
}

// ========================================
// CITY ENDPOINTS
// ========================================

// GetCities godoc
// @Summary Get cities by province
// @Description Retrieve cities with optional filtering by province, active status, and search query
// @Tags Master Data - Locations
// @Accept json
// @Produce json
// @Param province_id query int true "Province ID to filter cities"
// @Param active query boolean false "Filter by active status (true/false)"
// @Param search query string false "Search cities by name (min 2 characters)"
// @Success 200 {object} utils.Response{data=[]master.CityResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/meta/locations/cities [get]
func (h *LocationHandler) GetCities(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse query parameters
	provinceIDParam := strings.TrimSpace(c.Query("province_id"))
	if provinceIDParam == "" {
		return utils.BadRequestResponse(c, "province_id query parameter is required")
	}

	provinceID, err := strconv.ParseInt(provinceIDParam, 10, 64)
	if err != nil || provinceID <= 0 {
		return utils.BadRequestResponse(c, "Invalid province_id parameter")
	}

	activeParam := strings.TrimSpace(c.Query("active"))
	search := strings.TrimSpace(c.Query("search"))

	// Determine if we should filter by active status
	var cities []master.CityResponse

	if activeParam == "true" {
		// Get only active cities in province
		cities, err = h.cityService.GetActiveByProvinceID(ctx, provinceID, search)
	} else {
		// Get all cities in province (including inactive)
		cities, err = h.cityService.GetByProvinceID(ctx, provinceID, search)
	}

	if err != nil {
		// Check for specific errors
		errMsg := err.Error()
		if strings.Contains(errMsg, "province not found") {
			return utils.NotFoundResponse(c, "Province not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve cities", errMsg)
	}

	return utils.SuccessResponse(c, "Cities retrieved successfully", cities)
}

// GetCityByID godoc
// @Summary Get city by ID
// @Description Retrieve a single city by its ID with province information
// @Tags Master Data - Locations
// @Accept json
// @Produce json
// @Param id path int true "City ID"
// @Success 200 {object} utils.Response{data=master.CityResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/meta/locations/cities/{id} [get]
func (h *LocationHandler) GetCityByID(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse ID from path parameter
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid city ID")
	}

	// Get city by ID
	city, err := h.cityService.GetByID(ctx, id)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == "city not found" || strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "City not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve city", err.Error())
	}

	return utils.SuccessResponse(c, "City retrieved successfully", city)
}

// ========================================
// DISTRICT ENDPOINTS
// ========================================

// GetDistricts godoc
// @Summary Get districts by city
// @Description Retrieve districts with optional filtering by city, active status, and search query
// @Tags Master Data - Locations
// @Accept json
// @Produce json
// @Param city_id query int true "City ID to filter districts"
// @Param active query boolean false "Filter by active status (true/false)"
// @Param search query string false "Search districts by name (min 2 characters)"
// @Success 200 {object} utils.Response{data=[]master.DistrictResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/meta/locations/districts [get]
func (h *LocationHandler) GetDistricts(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse query parameters
	cityIDParam := strings.TrimSpace(c.Query("city_id"))
	if cityIDParam == "" {
		return utils.BadRequestResponse(c, "city_id query parameter is required")
	}

	cityID, err := strconv.ParseInt(cityIDParam, 10, 64)
	if err != nil || cityID <= 0 {
		return utils.BadRequestResponse(c, "Invalid city_id parameter")
	}

	activeParam := strings.TrimSpace(c.Query("active"))
	search := strings.TrimSpace(c.Query("search"))

	// Determine if we should filter by active status
	var districts []master.DistrictResponse

	if activeParam == "true" {
		// Get only active districts in city
		districts, err = h.districtService.GetActiveByCityID(ctx, cityID, search)
	} else {
		// Get all districts in city (including inactive)
		districts, err = h.districtService.GetByCityID(ctx, cityID, search)
	}

	if err != nil {
		// Check for specific errors
		errMsg := err.Error()
		if strings.Contains(errMsg, "city not found") {
			return utils.NotFoundResponse(c, "City not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve districts", errMsg)
	}

	return utils.SuccessResponse(c, "Districts retrieved successfully", districts)
}

// GetDistrictByID godoc
// @Summary Get district by ID
// @Description Retrieve a single district by its ID with full location hierarchy (city and province)
// @Tags Master Data - Locations
// @Accept json
// @Produce json
// @Param id path int true "District ID"
// @Success 200 {object} utils.Response{data=master.DistrictResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/meta/locations/districts/{id} [get]
func (h *LocationHandler) GetDistrictByID(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse ID from path parameter
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid district ID")
	}

	// Get district by ID with full location hierarchy
	district, err := h.districtService.GetByID(ctx, id)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == "district not found" || strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "District not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve district", err.Error())
	}

	return utils.SuccessResponse(c, "District retrieved successfully", district)
}
