package http

import (
	"strconv"
	"strings"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/helpers"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// CompanyBasicHandler handles basic CRUD operations for companies
// This includes listing, creating, getting, updating, deleting companies,
// and managing company images (logo and banner).
type CompanyBasicHandler struct {
	companyService  company.CompanyService
	industryRepo    master.IndustryRepository
	companySizeRepo master.CompanySizeRepository
	provinceRepo    master.ProvinceRepository
	cityRepo        master.CityRepository
	districtRepo    master.DistrictRepository
}

// NewCompanyBasicHandler creates a new instance of CompanyBasicHandler
func NewCompanyBasicHandler(
	companyService company.CompanyService,
	industryRepo master.IndustryRepository,
	companySizeRepo master.CompanySizeRepository,
	provinceRepo master.ProvinceRepository,
	cityRepo master.CityRepository,
	districtRepo master.DistrictRepository,
) *CompanyBasicHandler {
	return &CompanyBasicHandler{
		companyService:  companyService,
		industryRepo:    industryRepo,
		companySizeRepo: companySizeRepo,
		provinceRepo:    provinceRepo,
		cityRepo:        cityRepo,
		districtRepo:    districtRepo,
	}
}

// ListCompanies godoc
// @Summary List companies
// @Description List companies with filters and pagination (Glints-like)
// @Tags companies
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param industry query string false "Industry"
// @Param company_type query string false "Company type"
// @Param size_category query string false "Size category"
// @Param location query string false "Location (city)"
// @Param is_verified query boolean false "Only verified companies"
// @Param sort_by query string false "Sort by (name, created_at, followers, rating)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} utils.Response{data=response.CompanyListResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies [get]
func (h *CompanyBasicHandler) ListCompanies(c *fiber.Ctx) error {
	ctx := c.Context()

	var q request.CompanySearchRequest
	if err := c.QueryParser(&q); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidQueryParams)
	}
	if err := utils.ValidateStruct(&q); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, ErrValidationFailed, errs)
	}
	q.Page, q.Limit = utils.ValidatePagination(q.Page, q.Limit, 100)

	// Sanitize search inputs
	q.Query = utils.SanitizeString(q.Query)
	q.Location = utils.SanitizeString(q.Location)

	// Build filter using helper
	filter := helpers.BuildCompanyFilter(q)

	var (
		companies []company.Company
		total     int64
		err       error
	)
	if q.Query != "" {
		companies, total, err = h.companyService.SearchCompanies(ctx, q.Query, filter)
	} else {
		companies, total, err = h.companyService.ListCompanies(ctx, filter)
	}
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to list companies", err.Error())
	}

	respList := make([]response.CompanyResponse, 0, len(companies))
	for _, comp := range companies {
		cr := mapper.ToCompanyResponse(&comp)
		if cr != nil {
			respList = append(respList, *cr)
		}
	}

	meta := utils.GetPaginationMeta(q.Page, q.Limit, total)
	payload := response.CompanyListResponse{Companies: respList}
	return utils.SuccessResponseWithMeta(c, MsgFetchedSuccess, payload, meta)
}

// CreateCompany godoc
// @Summary Create a new company profile
// @Description Create a new company profile with basic details and optional logo
// @Tags companies
// @Accept json,multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param request body request.RegisterCompanyRequest false "Register company request (for JSON)"
// @Param company_name formData string false "Company name (for multipart)"
// @Param industry_id formData int false "Industry ID"
// @Param company_size_id formData int false "Company size ID"
// @Param district_id formData int false "District ID"
// @Param full_address formData string false "Full address"
// @Param description formData string false "Description"
// @Param logo formData file false "Company logo image"
// @Success 201 {object} utils.Response{data=response.CompanyDetailResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies [post]
func (h *CompanyBasicHandler) CreateCompany(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get authenticated user ID from JWT using middleware helper
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated", "userID not found in context")
	}

	// Check if user already owns a company (1 user = 1 company)
	existingCompanies, err := h.companyService.GetUserCompanies(ctx, userID)
	if err == nil && len(existingCompanies) > 0 {
		return utils.ErrorResponse(c, fiber.StatusForbidden,
			"User already owns a company",
			"Business rule violation: Each user can only register one company. You already own a company.")
	}

	// Parse request body (support both JSON and multipart/form-data)
	var req request.RegisterCompanyRequest

	// Check content type
	contentType := string(c.Request().Header.ContentType())
	isMultipart := strings.Contains(contentType, "multipart/form-data")

	if isMultipart {
		// Parse form data
		req.CompanyName = c.FormValue("company_name")

		// Parse optional integer fields (ID-based - backward compatibility)
		if industryID := c.FormValue("industry_id"); industryID != "" {
			if id, err := strconv.ParseInt(industryID, 10, 64); err == nil {
				req.IndustryID = &id
			}
		}
		if companySizeID := c.FormValue("company_size_id"); companySizeID != "" {
			if id, err := strconv.ParseInt(companySizeID, 10, 64); err == nil {
				req.CompanySizeID = &id
			}
		}
		if districtID := c.FormValue("district_id"); districtID != "" {
			if id, err := strconv.ParseInt(districtID, 10, 64); err == nil {
				req.DistrictID = &id
			}
		}

		// Parse location names (for mobile app dropdown)
		if industryName := c.FormValue("industry_name"); industryName != "" {
			req.IndustryName = &industryName
		}
		if companySizeName := c.FormValue("company_size_name"); companySizeName != "" {
			req.CompanySizeName = &companySizeName
		}
		if provinceName := c.FormValue("province_name"); provinceName != "" {
			req.ProvinceName = &provinceName
		}
		if cityName := c.FormValue("city_name"); cityName != "" {
			req.CityName = &cityName
		}
		if districtName := c.FormValue("district_name"); districtName != "" {
			req.DistrictName = &districtName
		}

		// Parse optional string fields
		if fullAddress := c.FormValue("full_address"); fullAddress != "" {
			req.FullAddress = fullAddress
		}
		if description := c.FormValue("description"); description != "" {
			req.Description = &description
		}
		if legalName := c.FormValue("legal_name"); legalName != "" {
			req.LegalName = &legalName
		}
		if registrationNumber := c.FormValue("registration_number"); registrationNumber != "" {
			req.RegistrationNumber = &registrationNumber
		}
		if websiteURL := c.FormValue("website_url"); websiteURL != "" {
			req.WebsiteURL = &websiteURL
		}
		if phone := c.FormValue("phone"); phone != "" {
			req.Phone = &phone
		}
		if about := c.FormValue("about"); about != "" {
			req.About = &about
		}
	} else {
		// Parse JSON body
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, ErrInvalidRequest)
		}
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, ErrValidationFailed, errors)
	}

	// Sanitize required fields
	req.CompanyName = utils.SanitizeString(req.CompanyName)

	// Sanitize optional pointer fields
	if req.LegalName != nil {
		sanitized := utils.SanitizeString(*req.LegalName)
		req.LegalName = &sanitized
	}
	if req.About != nil {
		sanitized := utils.SanitizeHTML(*req.About)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, ErrPotentialXSS)
		}
		req.About = &sanitized
	}

	// Sanitize FullAddress
	if req.FullAddress != "" {
		req.FullAddress = utils.SanitizeString(req.FullAddress)
	}

	// Sanitize Description
	if req.Description != nil {
		sanitized := utils.SanitizeHTML(*req.Description)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, ErrPotentialXSS)
		}
		req.Description = &sanitized
	}

	// Sanitize legacy address fields (for backward compatibility)
	if req.Address != nil {
		sanitized := utils.SanitizeString(*req.Address)
		req.Address = &sanitized
	}
	if req.City != nil {
		sanitized := utils.SanitizeString(*req.City)
		req.City = &sanitized
	}
	if req.Province != nil {
		sanitized := utils.SanitizeString(*req.Province)
		req.Province = &sanitized
	}
	if req.Country != nil {
		sanitized := utils.SanitizeString(*req.Country)
		req.Country = &sanitized
	}

	// Resolve location names to IDs (if name-based input is provided)
	// This allows mobile app to send location names instead of IDs
	if req.IndustryName != nil && *req.IndustryName != "" {
		industry, err := h.industryRepo.GetByName(ctx, *req.IndustryName)
		if err != nil {
			return utils.BadRequestResponse(c, "Industry not found: "+*req.IndustryName)
		}
		req.IndustryID = &industry.ID
	}

	if req.CompanySizeName != nil && *req.CompanySizeName != "" {
		companySize, err := h.companySizeRepo.GetByCategory(ctx, *req.CompanySizeName)
		if err != nil {
			return utils.BadRequestResponse(c, "Company size not found: "+*req.CompanySizeName)
		}
		req.CompanySizeID = &companySize.ID
	}

	// Resolve location hierarchy: Province -> City -> District
	if req.ProvinceName != nil && *req.ProvinceName != "" {
		province, err := h.provinceRepo.GetByName(ctx, *req.ProvinceName)
		if err != nil {
			return utils.BadRequestResponse(c, "Province not found: "+*req.ProvinceName)
		}

		// If city name provided, find city within the province
		if req.CityName != nil && *req.CityName != "" {
			city, err := h.cityRepo.GetByNameAndProvinceID(ctx, *req.CityName, province.ID)
			if err != nil {
				return utils.BadRequestResponse(c, "City not found: "+*req.CityName+" in province: "+*req.ProvinceName)
			}

			// If district name provided, find district within the city
			if req.DistrictName != nil && *req.DistrictName != "" {
				district, err := h.districtRepo.GetByNameAndCityID(ctx, *req.DistrictName, city.ID)
				if err != nil {
					return utils.BadRequestResponse(c, "District not found: "+*req.DistrictName+" in city: "+*req.CityName)
				}
				req.DistrictID = &district.ID
			}
		}
	}

	// Convert to domain request
	domainReq := &company.RegisterCompanyRequest{
		CompanyName:        req.CompanyName,
		LegalName:          req.LegalName,
		RegistrationNumber: req.RegistrationNumber,

		// Master Data Relations
		IndustryID:    req.IndustryID,
		CompanySizeID: req.CompanySizeID,
		DistrictID:    req.DistrictID,
		FullAddress:   req.FullAddress,
		Description:   req.Description,

		// Legacy Fields (for backward compatibility)
		Industry:     req.Industry,
		CompanyType:  req.CompanyType,
		SizeCategory: req.SizeCategory,
		Address:      req.Address,
		City:         req.City,
		Province:     req.Province,

		// Other Fields
		WebsiteURL:  req.WebsiteURL,
		EmailDomain: req.EmailDomain,
		Phone:       req.Phone,
		Country:     req.Country,
		PostalCode:  req.PostalCode,
		About:       req.About,
	}

	// Create the company with user as owner
	createdCompany, err := h.companyService.RegisterCompany(ctx, domainReq, userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, ErrFailedOperation)
	}

	// Handle logo upload if provided (multipart/form-data only)
	if isMultipart {
		if logoFile, err := c.FormFile("logo"); err == nil && logoFile != nil {
			// Upload logo (this will update the database)
			logoURL, uploadErr := h.companyService.UploadLogo(ctx, createdCompany.ID, logoFile)
			if uploadErr == nil && logoURL != "" {
				// Update the company object with logo URL for response
				createdCompany.LogoURL = &logoURL
			}
			// Don't fail the whole request if logo upload fails
		}
	}

	// Map to response DTO with master data
	response := mapper.ToCompanyDetailResponse(createdCompany)
	return utils.CreatedResponse(c, MsgCreatedSuccess, response)
}

// GetCompany godoc
// @Summary Get company profile by ID
// @Description Retrieve a specific company profile by ID
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} utils.Response{data=response.CompanyResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id} [get]
func (h *CompanyBasicHandler) GetCompany(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get company ID from the URL
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	// Get the company from the service
	companyData, err := h.companyService.GetCompany(ctx, int64(companyID))
	if err != nil {
		return utils.NotFoundResponse(c, ErrCompanyNotFound)
	}

	// Map to response DTO
	response := mapper.ToCompanyResponse(companyData)
	return utils.SuccessResponse(c, MsgFetchedSuccess, response)
}

// GetCompanyBySlug godoc
// @Summary Get company profile by slug
// @Tags companies
// @Accept json
// @Produce json
// @Param slug path string true "Company slug"
// @Success 200 {object} utils.Response{data=response.CompanyResponse}
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/slug/{slug} [get]
func (h *CompanyBasicHandler) GetCompanyBySlug(c *fiber.Ctx) error {
	ctx := c.Context()
	slug := utils.SanitizeString(strings.TrimSpace(c.Params("slug")))
	if slug == "" {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}

	companyData, err := h.companyService.GetCompanyBySlug(ctx, slug)
	if err != nil {
		return utils.NotFoundResponse(c, ErrCompanyNotFound)
	}
	responseDTO := mapper.ToCompanyResponse(companyData)
	return utils.SuccessResponse(c, MsgFetchedSuccess, responseDTO)
}

// UpdateCompany godoc
// @Summary Update a company profile
// @Description Update the details of a company profile (Only owner/admin can update)
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Param request body request.UpdateCompanyRequest true "Update company request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id} [put]
func (h *CompanyBasicHandler) UpdateCompany(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get authenticated user ID from JWT using middleware helper
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated", "userID not found in context")
	}

	// Get company ID from the URL
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	// Check if user has permission (owner or admin)
	hasPermission, err := h.companyService.CheckEmployerPermission(ctx, userID, int64(companyID), "admin")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check user permission", err.Error())
	}
	if !hasPermission {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to update this company. Only company owner or admin can perform this action.", "")
	}

	// Parse request body
	var req request.UpdateCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, ErrValidationFailed, errors)
	}

	// Sanitize all string pointer fields
	if req.CompanyName != nil {
		sanitized := utils.SanitizeString(*req.CompanyName)
		req.CompanyName = &sanitized
	}
	if req.LegalName != nil {
		sanitized := utils.SanitizeString(*req.LegalName)
		req.LegalName = &sanitized
	}
	if req.About != nil {
		sanitized := utils.SanitizeHTML(*req.About)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, ErrPotentialXSS)
		}
		req.About = &sanitized
	}
	if req.Culture != nil {
		sanitized := utils.SanitizeHTML(*req.Culture)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, ErrPotentialXSS)
		}
		req.Culture = &sanitized
	}

	// Sanitize FullAddress
	if req.FullAddress != nil {
		sanitized := utils.SanitizeString(*req.FullAddress)
		req.FullAddress = &sanitized
	}

	// Sanitize Description
	if req.Description != nil {
		sanitized := utils.SanitizeHTML(*req.Description)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, ErrPotentialXSS)
		}
		req.Description = &sanitized
	}

	// Sanitize legacy address fields (for backward compatibility)
	if req.Address != nil {
		sanitized := utils.SanitizeString(*req.Address)
		req.Address = &sanitized
	}
	if req.City != nil {
		sanitized := utils.SanitizeString(*req.City)
		req.City = &sanitized
	}
	if req.Province != nil {
		sanitized := utils.SanitizeString(*req.Province)
		req.Province = &sanitized
	}

	// Convert to domain request
	domainReq := &company.UpdateCompanyRequest{
		CompanyName:        req.CompanyName,
		LegalName:          req.LegalName,
		RegistrationNumber: req.RegistrationNumber,

		// Master Data Relations
		IndustryID:    req.IndustryID,
		CompanySizeID: req.CompanySizeID,
		DistrictID:    req.DistrictID,
		FullAddress:   req.FullAddress,
		Description:   req.Description,

		// Legacy Fields (for backward compatibility)
		Industry:     req.Industry,
		CompanyType:  req.CompanyType,
		SizeCategory: req.SizeCategory,
		Address:      req.Address,
		City:         req.City,
		Province:     req.Province,

		// Location
		Latitude:  req.Latitude,
		Longitude: req.Longitude,

		// Other Fields
		WebsiteURL:  req.WebsiteURL,
		EmailDomain: req.EmailDomain,
		Phone:       req.Phone,
		PostalCode:  req.PostalCode,
		About:       req.About,
		Culture:     req.Culture,
		Benefits:    req.Benefits,
	}

	// Update the company profile
	if err := h.companyService.UpdateCompany(ctx, int64(companyID), domainReq); err != nil {
		return utils.InternalServerErrorResponse(c, ErrFailedOperation)
	}

	return utils.SuccessResponse(c, MsgUpdatedSuccess, nil)
}

// DeleteCompany godoc
// @Summary Delete a company profile
// @Description Delete a company profile by ID (Only owner can delete)
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id} [delete]
func (h *CompanyBasicHandler) DeleteCompany(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get authenticated user ID from JWT using middleware helper
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated", "userID not found in context")
	}

	// Get company ID from the URL
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	// Check if user is the owner (only owner can delete company)
	isOwner, err := h.companyService.CheckEmployerPermission(ctx, userID, int64(companyID), "owner")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check user permission", err.Error())
	}
	if !isOwner {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to delete this company. Only company owner can perform this action.", "")
	}

	// Delete the company profile
	if err := h.companyService.DeleteCompany(ctx, int64(companyID)); err != nil {
		return utils.InternalServerErrorResponse(c, ErrFailedOperation)
	}

	return utils.SuccessResponse(c, MsgDeletedSuccess, nil)
}

// UploadLogo godoc
// @Summary Upload company logo
// @Description Upload a logo image for the company
// @Tags companies
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Param file formData file true "Logo image file"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/logo [post]
func (h *CompanyBasicHandler) UploadLogo(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	file := middleware.GetUploadedFile(c)
	if file == nil {
		return utils.BadRequestResponse(c, ErrNoFileUploaded)
	}

	url, err := h.companyService.UploadLogo(ctx, int64(companyID), file)
	if err != nil {
		return utils.InternalServerErrorResponse(c, ErrFileUploadFailed)
	}
	return utils.CreatedResponse(c, MsgUploadSuccess, fiber.Map{"logo_url": url})
}

// UploadBanner godoc
// @Summary Upload company banner
// @Description Upload a banner image for the company
// @Tags companies
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Param file formData file true "Banner image file"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/banner [post]
func (h *CompanyBasicHandler) UploadBanner(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	file := middleware.GetUploadedFile(c)
	if file == nil {
		return utils.BadRequestResponse(c, ErrNoFileUploaded)
	}

	url, err := h.companyService.UploadBanner(ctx, int64(companyID), file)
	if err != nil {
		return utils.InternalServerErrorResponse(c, ErrFileUploadFailed)
	}
	return utils.CreatedResponse(c, MsgUploadSuccess, fiber.Map{"banner_url": url})
}

// DeleteLogo godoc
// @Summary Delete company logo
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/logo [delete]
func (h *CompanyBasicHandler) DeleteLogo(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	if err := h.companyService.DeleteLogo(ctx, int64(companyID)); err != nil {
		return utils.InternalServerErrorResponse(c, ErrFailedOperation)
	}
	return utils.SuccessResponse(c, MsgDeletedSuccess, nil)
}

// DeleteBanner godoc
// @Summary Delete company banner
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/banner [delete]
func (h *CompanyBasicHandler) DeleteBanner(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	if err := h.companyService.DeleteBanner(ctx, int64(companyID)); err != nil {
		return utils.InternalServerErrorResponse(c, ErrFailedOperation)
	}
	return utils.SuccessResponse(c, MsgDeletedSuccess, nil)
}

// GetMyCompanies godoc
// @Summary Get my companies
// @Description Get all companies where the authenticated user is a member
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]response.CompanyResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/my-companies [get]
func (h *CompanyBasicHandler) GetMyCompanies(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get authenticated user ID from JWT using middleware helper
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated", "userID not found in context")
	}

	// Get companies where user is a member
	companies, err := h.companyService.GetUserCompanies(ctx, userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, ErrFailedOperation)
	}

	// Map to response DTOs
	responses := make([]response.CompanyResponse, 0, len(companies))
	for _, comp := range companies {
		responses = append(responses, *mapper.ToCompanyResponse(&comp))
	}

	return utils.SuccessResponse(c, MsgFetchedSuccess, responses)
}
