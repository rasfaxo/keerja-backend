package admin

import (
	"strconv"
	"strings"

	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// AdminMasterDataHandler handles admin CRUD operations for master data
type AdminMasterDataHandler struct {
	provinceService    master.AdminProvinceService
	cityService        master.AdminCityService
	districtService    master.AdminDistrictService
	industryService    master.AdminIndustryService
	companySizeService master.AdminCompanySizeService
	jobTypeService     master.AdminJobTypeService
}

// NewAdminMasterDataHandler creates a new admin master data handler
func NewAdminMasterDataHandler(
	provinceService master.AdminProvinceService,
	cityService master.AdminCityService,
	districtService master.AdminDistrictService,
	industryService master.AdminIndustryService,
	companySizeService master.AdminCompanySizeService,
	jobTypeService master.AdminJobTypeService,
) *AdminMasterDataHandler {
	return &AdminMasterDataHandler{
		provinceService:    provinceService,
		cityService:        cityService,
		districtService:    districtService,
		industryService:    industryService,
		companySizeService: companySizeService,
		jobTypeService:     jobTypeService,
	}
}

// ========================================
// PROVINCE CRUD ENDPOINTS
// ========================================

// CreateProvince handles POST /api/v1/admin/master/provinces
func (h *AdminMasterDataHandler) CreateProvince(c *fiber.Ctx) error {
	var req master.CreateProvinceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Check duplicate code
	exists, err := h.provinceService.CheckDuplicateCode(c.Context(), req.Code)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to check duplicate province code")
	}
	if exists {
		return utils.ConflictResponse(c, "Province with this code already exists")
	}

	// Create province
	province, err := h.provinceService.Create(c.Context(), req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create province")
	}

	return utils.CreatedResponse(c, "Province created successfully", province)
}

// GetProvinces handles GET /api/v1/admin/master/provinces
func (h *AdminMasterDataHandler) GetProvinces(c *fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search", ""))
	activeParam := strings.TrimSpace(c.Query("active", ""))

	var provinces []master.ProvinceResponse
	var err error

	if activeParam == "true" {
		provinces, err = h.provinceService.GetActive(c.Context(), search)
	} else {
		provinces, err = h.provinceService.GetAll(c.Context(), search)
	}

	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve provinces")
	}

	return utils.SuccessResponse(c, "Provinces retrieved successfully", provinces)
}

// GetProvinceByID handles GET /api/v1/admin/master/provinces/:id
func (h *AdminMasterDataHandler) GetProvinceByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid province ID")
	}

	province, err := h.provinceService.GetByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "Province not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to retrieve province")
	}

	return utils.SuccessResponse(c, "Province retrieved successfully", province)
}

// UpdateProvince handles PUT /api/v1/admin/master/provinces/:id
func (h *AdminMasterDataHandler) UpdateProvince(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid province ID")
	}

	var req master.UpdateProvinceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Check duplicate code if code is being updated
	if req.Code != "" {
		exists, err := h.provinceService.CheckDuplicateCode(c.Context(), req.Code)
		if err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to check duplicate province code")
		}
		if exists {
			// Verify it's not the same record
			existing, err := h.provinceService.GetByID(c.Context(), id)
			if err == nil && existing != nil && existing.Code != req.Code {
				return utils.ConflictResponse(c, "Province with this code already exists")
			}
		}
	}

	// Update province
	province, err := h.provinceService.Update(c.Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "Province not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to update province")
	}

	return utils.SuccessResponse(c, "Province updated successfully", province)
}

// DeleteProvince handles DELETE /api/v1/admin/master/provinces/:id
func (h *AdminMasterDataHandler) DeleteProvince(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid province ID")
	}

	// Check references before deleting
	cities, companies, err := h.provinceService.CountReferences(c.Context(), id)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to check province references")
	}

	if cities > 0 || companies > 0 {
		return utils.ConflictResponse(c,
			"Cannot delete province: it is still referenced by cities or companies")
	}

	// Delete province
	if err := h.provinceService.Delete(c.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "Province not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to delete province")
	}

	return utils.SuccessResponse(c, "Province deleted successfully", nil)
}

// ========================================
// CITY CRUD ENDPOINTS
// ========================================

// CreateCity handles POST /api/v1/admin/master/cities
func (h *AdminMasterDataHandler) CreateCity(c *fiber.Ctx) error {
	var req master.CreateCityRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Check duplicate name in province
	exists, err := h.cityService.CheckDuplicateNameInProvince(c.Context(), req.Name, req.ProvinceID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to check duplicate city name")
	}
	if exists {
		return utils.ConflictResponse(c, "City with this name already exists in the province")
	}

	// Create city
	city, err := h.cityService.Create(c.Context(), req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create city")
	}

	return utils.CreatedResponse(c, "City created successfully", city)
}

// GetCities handles GET /api/v1/admin/master/cities
func (h *AdminMasterDataHandler) GetCities(c *fiber.Ctx) error {
	provinceIDParam := strings.TrimSpace(c.Query("province_id"))
	search := strings.TrimSpace(c.Query("search", ""))
	activeParam := strings.TrimSpace(c.Query("active", ""))

	if provinceIDParam == "" {
		return utils.BadRequestResponse(c, "province_id query parameter is required")
	}

	provinceID, err := strconv.ParseInt(provinceIDParam, 10, 64)
	if err != nil || provinceID <= 0 {
		return utils.BadRequestResponse(c, "Invalid province_id parameter")
	}

	var cities []master.CityResponse
	if activeParam == "true" {
		cities, err = h.cityService.GetActiveByProvinceID(c.Context(), provinceID, search)
	} else {
		cities, err = h.cityService.GetByProvinceID(c.Context(), provinceID, search)
	}

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "Province not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to retrieve cities")
	}

	return utils.SuccessResponse(c, "Cities retrieved successfully", cities)
}

// GetCityByID handles GET /api/v1/admin/master/cities/:id
func (h *AdminMasterDataHandler) GetCityByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid city ID")
	}

	city, err := h.cityService.GetByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "City not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to retrieve city")
	}

	return utils.SuccessResponse(c, "City retrieved successfully", city)
}

// UpdateCity handles PUT /api/v1/admin/master/cities/:id
func (h *AdminMasterDataHandler) UpdateCity(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid city ID")
	}

	var req master.UpdateCityRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Check duplicate name if name or province is being updated
	if req.Name != "" {
		provinceID := req.ProvinceID
		if provinceID == nil {
			// Get existing city to use its province ID
			existing, err := h.cityService.GetByID(c.Context(), id)
			if err != nil {
				return utils.NotFoundResponse(c, "City not found")
			}
			provinceID = &existing.ProvinceID
		}

		exists, err := h.cityService.CheckDuplicateNameInProvince(c.Context(), req.Name, *provinceID)
		if err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to check duplicate city name")
		}
		if exists {
			// Verify it's not the same record
			existing, err := h.cityService.GetByID(c.Context(), id)
			if err == nil && existing != nil && existing.Name != req.Name {
				return utils.ConflictResponse(c, "City with this name already exists in the province")
			}
		}
	}

	// Update city
	city, err := h.cityService.Update(c.Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "City not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to update city")
	}

	return utils.SuccessResponse(c, "City updated successfully", city)
}

// DeleteCity handles DELETE /api/v1/admin/master/cities/:id
func (h *AdminMasterDataHandler) DeleteCity(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid city ID")
	}

	// Check references before deleting
	districts, companies, err := h.cityService.CountReferences(c.Context(), id)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to check city references")
	}

	if districts > 0 || companies > 0 {
		return utils.ConflictResponse(c,
			"Cannot delete city: it is still referenced by districts or companies")
	}

	// Delete city
	if err := h.cityService.Delete(c.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "City not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to delete city")
	}

	return utils.SuccessResponse(c, "City deleted successfully", nil)
}

// ========================================
// DISTRICT CRUD ENDPOINTS
// ========================================

// CreateDistrict handles POST /api/v1/admin/master/districts
func (h *AdminMasterDataHandler) CreateDistrict(c *fiber.Ctx) error {
	var req master.CreateDistrictRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Check duplicate name in city
	exists, err := h.districtService.CheckDuplicateNameInCity(c.Context(), req.Name, req.CityID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to check duplicate district name")
	}
	if exists {
		return utils.ConflictResponse(c, "District with this name already exists in the city")
	}

	// Create district
	district, err := h.districtService.Create(c.Context(), req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create district")
	}

	return utils.CreatedResponse(c, "District created successfully", district)
}

// GetDistricts handles GET /api/v1/admin/master/districts
func (h *AdminMasterDataHandler) GetDistricts(c *fiber.Ctx) error {
	cityIDParam := strings.TrimSpace(c.Query("city_id"))
	search := strings.TrimSpace(c.Query("search", ""))
	activeParam := strings.TrimSpace(c.Query("active", ""))

	if cityIDParam == "" {
		return utils.BadRequestResponse(c, "city_id query parameter is required")
	}

	cityID, err := strconv.ParseInt(cityIDParam, 10, 64)
	if err != nil || cityID <= 0 {
		return utils.BadRequestResponse(c, "Invalid city_id parameter")
	}

	var districts []master.DistrictResponse
	if activeParam == "true" {
		districts, err = h.districtService.GetActiveByCityID(c.Context(), cityID, search)
	} else {
		districts, err = h.districtService.GetByCityID(c.Context(), cityID, search)
	}

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "City not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to retrieve districts")
	}

	return utils.SuccessResponse(c, "Districts retrieved successfully", districts)
}

// GetDistrictByID handles GET /api/v1/admin/master/districts/:id
func (h *AdminMasterDataHandler) GetDistrictByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid district ID")
	}

	district, err := h.districtService.GetByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "District not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to retrieve district")
	}

	return utils.SuccessResponse(c, "District retrieved successfully", district)
}

// UpdateDistrict handles PUT /api/v1/admin/master/districts/:id
func (h *AdminMasterDataHandler) UpdateDistrict(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid district ID")
	}

	var req master.UpdateDistrictRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Check duplicate name if name or city is being updated
	if req.Name != "" {
		cityID := req.CityID
		if cityID == nil {
			// Get existing district to use its city ID
			existing, err := h.districtService.GetByID(c.Context(), id)
			if err != nil {
				return utils.NotFoundResponse(c, "District not found")
			}
			cityID = &existing.CityID
		}

		exists, err := h.districtService.CheckDuplicateNameInCity(c.Context(), req.Name, *cityID)
		if err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to check duplicate district name")
		}
		if exists {
			// Verify it's not the same record
			existing, err := h.districtService.GetByID(c.Context(), id)
			if err == nil && existing != nil && existing.Name != req.Name {
				return utils.ConflictResponse(c, "District with this name already exists in the city")
			}
		}
	}

	// Update district
	district, err := h.districtService.Update(c.Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "District not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to update district")
	}

	return utils.SuccessResponse(c, "District updated successfully", district)
}

// DeleteDistrict handles DELETE /api/v1/admin/master/districts/:id
func (h *AdminMasterDataHandler) DeleteDistrict(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid district ID")
	}

	// Check references before deleting
	companies, err := h.districtService.CountCompanyReferences(c.Context(), id)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to check district references")
	}

	if companies > 0 {
		return utils.ConflictResponse(c,
			"Cannot delete district: it is still referenced by companies")
	}

	// Delete district
	if err := h.districtService.Delete(c.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "District not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to delete district")
	}

	return utils.SuccessResponse(c, "District deleted successfully", nil)
}

// ========================================
// INDUSTRY CRUD ENDPOINTS
// ========================================

// CreateIndustry handles POST /api/v1/admin/master/industries
func (h *AdminMasterDataHandler) CreateIndustry(c *fiber.Ctx) error {
	var req master.CreateIndustryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Check duplicate name
	exists, err := h.industryService.CheckDuplicateName(c.Context(), req.Name)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to check duplicate industry name")
	}
	if exists {
		return utils.ConflictResponse(c, "Industry with this name already exists")
	}

	// Create industry
	industry, err := h.industryService.Create(c.Context(), req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create industry")
	}

	return utils.CreatedResponse(c, "Industry created successfully", industry)
}

// GetIndustries handles GET /api/v1/admin/master/industries
func (h *AdminMasterDataHandler) GetIndustries(c *fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search", ""))
	activeParam := strings.TrimSpace(c.Query("active", ""))

	var industries []master.IndustryResponse
	var err error

	if activeParam == "true" {
		industries, err = h.industryService.GetActive(c.Context(), search)
	} else {
		industries, err = h.industryService.GetAll(c.Context(), search)
	}

	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve industries")
	}

	return utils.SuccessResponse(c, "Industries retrieved successfully", industries)
}

// GetIndustryByID handles GET /api/v1/admin/master/industries/:id
func (h *AdminMasterDataHandler) GetIndustryByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid industry ID")
	}

	industry, err := h.industryService.GetByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "Industry not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to retrieve industry")
	}

	return utils.SuccessResponse(c, "Industry retrieved successfully", industry)
}

// UpdateIndustry handles PUT /api/v1/admin/master/industries/:id
func (h *AdminMasterDataHandler) UpdateIndustry(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid industry ID")
	}

	var req master.UpdateIndustryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Check duplicate name if name is being updated
	if req.Name != "" {
		exists, err := h.industryService.CheckDuplicateName(c.Context(), req.Name)
		if err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to check duplicate industry name")
		}
		if exists {
			// Verify it's not the same record
			existing, err := h.industryService.GetByID(c.Context(), id)
			if err == nil && existing != nil && existing.Name != req.Name {
				return utils.ConflictResponse(c, "Industry with this name already exists")
			}
		}
	}

	// Update industry
	industry, err := h.industryService.Update(c.Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "Industry not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to update industry")
	}

	return utils.SuccessResponse(c, "Industry updated successfully", industry)
}

// DeleteIndustry handles DELETE /api/v1/admin/master/industries/:id
func (h *AdminMasterDataHandler) DeleteIndustry(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid industry ID")
	}

	// Check references before deleting
	companies, err := h.industryService.CountCompanyReferences(c.Context(), id)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to check industry references")
	}

	if companies > 0 {
		return utils.ConflictResponse(c,
			"Cannot delete industry: it is still referenced by companies")
	}

	// Delete industry
	if err := h.industryService.Delete(c.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "Industry not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to delete industry")
	}

	return utils.SuccessResponse(c, "Industry deleted successfully", nil)
}

// ========================================
// JOB TYPE CRUD ENDPOINTS
// ========================================

// CreateJobType handles POST /api/v1/admin/master/job-types
func (h *AdminMasterDataHandler) CreateJobType(c *fiber.Ctx) error {
	var req master.CreateJobTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Check duplicate code
	exists, err := h.jobTypeService.CheckDuplicateCode(c.Context(), req.Code, nil)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to check duplicate job type code")
	}
	if exists {
		return utils.ConflictResponse(c, "Job type with this code already exists")
	}

	// Create job type
	jobType, err := h.jobTypeService.Create(c.Context(), req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create job type")
	}

	return utils.CreatedResponse(c, "Job type created successfully", jobType)
}

// GetJobTypes handles GET /api/v1/admin/master/job-types
func (h *AdminMasterDataHandler) GetJobTypes(c *fiber.Ctx) error {
	jobTypes, err := h.jobTypeService.GetJobTypes(c.Context())
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job types")
	}

	// Convert to response format
	responses := make([]master.JobTypeResponse, len(jobTypes))
	for i, jt := range jobTypes {
		responses[i] = master.JobTypeResponse{
			ID:    jt.ID,
			Name:  jt.Name,
			Code:  jt.Code,
			Order: jt.Order,
		}
	}

	return utils.SuccessResponse(c, "Job types retrieved successfully", responses)
}

// GetJobTypeByID handles GET /api/v1/admin/master/job-types/:id
func (h *AdminMasterDataHandler) GetJobTypeByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid job type ID")
	}

	jobType, err := h.jobTypeService.GetJobTypeByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "Job type not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job type")
	}

	if jobType == nil {
		return utils.NotFoundResponse(c, "Job type not found")
	}

	response := master.JobTypeResponse{
		ID:    jobType.ID,
		Name:  jobType.Name,
		Code:  jobType.Code,
		Order: jobType.Order,
	}

	return utils.SuccessResponse(c, "Job type retrieved successfully", response)
}

// UpdateJobType handles PUT /api/v1/admin/master/job-types/:id
func (h *AdminMasterDataHandler) UpdateJobType(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid job type ID")
	}

	var req master.UpdateJobTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Check duplicate code if code is being updated
	if req.Code != "" {
		excludeID := &id
		exists, err := h.jobTypeService.CheckDuplicateCode(c.Context(), req.Code, excludeID)
		if err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to check duplicate job type code")
		}
		if exists {
			return utils.ConflictResponse(c, "Job type with this code already exists")
		}
	}

	// Update job type
	jobType, err := h.jobTypeService.Update(c.Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "Job type not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to update job type")
	}

	return utils.SuccessResponse(c, "Job type updated successfully", jobType)
}

// DeleteJobType handles DELETE /api/v1/admin/master/job-types/:id
func (h *AdminMasterDataHandler) DeleteJobType(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid job type ID")
	}

	// Check references before deleting
	jobs, err := h.jobTypeService.CountJobReferences(c.Context(), id)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to check job type references")
	}

	if jobs > 0 {
		return utils.ConflictResponse(c,
			"Cannot delete job type: it is still referenced by jobs")
	}

	// Delete job type
	if err := h.jobTypeService.Delete(c.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "Job type not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to delete job type")
	}

	return utils.SuccessResponse(c, "Job type deleted successfully", nil)
}

// ========================================
// COMPANY SIZE CRUD ENDPOINTS
// ========================================

// CreateCompanySize handles POST /api/v1/admin/meta/company-sizes
func (h *AdminMasterDataHandler) CreateCompanySize(c *fiber.Ctx) error {
	var req master.CreateCompanySizeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Check duplicate label
	exists, err := h.companySizeService.CheckDuplicateCategory(c.Context(), req.Label)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to check duplicate company size label")
	}
	if exists {
		return utils.ConflictResponse(c, "Company size with this label already exists")
	}

	// Create company size
	companySize, err := h.companySizeService.Create(c.Context(), req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create company size")
	}

	return utils.CreatedResponse(c, "Company size created successfully", companySize)
}

// GetCompanySizes handles GET /api/v1/admin/meta/company-sizes
func (h *AdminMasterDataHandler) GetCompanySizes(c *fiber.Ctx) error {
	activeParam := strings.TrimSpace(c.Query("active", ""))

	var companySizes []master.CompanySizeResponse
	var err error

	if activeParam == "true" {
		companySizes, err = h.companySizeService.GetActive(c.Context())
	} else {
		companySizes, err = h.companySizeService.GetAll(c.Context())
	}

	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve company sizes")
	}

	return utils.SuccessResponse(c, "Company sizes retrieved successfully", companySizes)
}

// GetCompanySizeByID handles GET /api/v1/admin/meta/company-sizes/:id
func (h *AdminMasterDataHandler) GetCompanySizeByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid company size ID")
	}

	companySize, err := h.companySizeService.GetByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "Company size not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to retrieve company size")
	}

	return utils.SuccessResponse(c, "Company size retrieved successfully", companySize)
}

// UpdateCompanySize handles PUT /api/v1/admin/meta/company-sizes/:id
func (h *AdminMasterDataHandler) UpdateCompanySize(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid company size ID")
	}

	var req master.UpdateCompanySizeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Check duplicate label if label is being updated
	if req.Label != "" {
		exists, err := h.companySizeService.CheckDuplicateCategory(c.Context(), req.Label)
		if err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to check duplicate company size label")
		}
		if exists {
			// Verify it's not the same record
			existing, err := h.companySizeService.GetByID(c.Context(), id)
			if err == nil && existing != nil && existing.Label != req.Label {
				return utils.ConflictResponse(c, "Company size with this label already exists")
			}
		}
	}

	// Update company size
	companySize, err := h.companySizeService.Update(c.Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "Company size not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to update company size")
	}

	return utils.SuccessResponse(c, "Company size updated successfully", companySize)
}

// DeleteCompanySize handles DELETE /api/v1/admin/meta/company-sizes/:id
func (h *AdminMasterDataHandler) DeleteCompanySize(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid company size ID")
	}

	// Check references before deleting
	companies, err := h.companySizeService.CountCompanyReferences(c.Context(), id)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to check company size references")
	}

	if companies > 0 {
		return utils.ConflictResponse(c,
			"Cannot delete company size: it is still referenced by companies")
	}

	// Delete company size
	if err := h.companySizeService.Delete(c.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFoundResponse(c, "Company size not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to delete company size")
	}

	return utils.SuccessResponse(c, "Company size deleted successfully", nil)
}
