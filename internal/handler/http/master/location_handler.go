package master

import (
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

func (h *LocationHandler) GetProvinceByID(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse ID from path parameter
	id, err := utils.ParseIDParam(c, "id")
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

func (h *LocationHandler) GetCities(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse query parameters
	provinceID := int64(c.QueryInt("province_id", 0))
	if provinceID <= 0 {
		return utils.BadRequestResponse(c, "province_id query parameter is required and must be a positive integer")
	}

	activeParam := strings.TrimSpace(c.Query("active"))
	search := strings.TrimSpace(c.Query("search"))

	// Determine if we should filter by active status
	var cities []master.CityResponse
	var err error

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

func (h *LocationHandler) GetCityByID(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse ID from path parameter
	id, err := utils.ParseIDParam(c, "id")
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

func (h *LocationHandler) GetDistricts(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse query parameters
	cityID := int64(c.QueryInt("city_id", 0))
	if cityID <= 0 {
		return utils.BadRequestResponse(c, "city_id query parameter is required and must be a positive integer")
	}

	activeParam := strings.TrimSpace(c.Query("active"))
	search := strings.TrimSpace(c.Query("search"))

	// Determine if we should filter by active status
	var districts []master.DistrictResponse
	var err error

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

func (h *LocationHandler) GetDistrictByID(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse ID from path parameter
	id, err := utils.ParseIDParam(c, "id")
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
