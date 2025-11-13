package companyhandler

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/handler/http"
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
		return utils.BadRequestResponse(c, http.ErrInvalidQueryParams)
	}
	if err := utils.ValidateStruct(&q); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, http.ErrValidationFailed, errs)
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
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, http.ErrFailedOperation, err.Error())
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
	return utils.SuccessResponseWithMeta(c, http.MsgFetchedSuccess, payload, meta)
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
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, http.ErrUnauthorized, "userID not found in context")
	}

	// Check if user already owns a company (1 user = 1 company)
	existingCompanies, err := h.companyService.GetUserCompanies(ctx, userID)
	if err == nil && len(existingCompanies) > 0 {
		return utils.ErrorResponse(c, fiber.StatusForbidden,
			http.ErrForbidden,
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
			return utils.BadRequestResponse(c, http.ErrInvalidBody)
		}
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, http.ErrValidationFailed, errors)
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
			return utils.BadRequestResponse(c, http.ErrPotentialXSS)
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
			return utils.BadRequestResponse(c, http.ErrPotentialXSS)
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
		return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
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
	return utils.CreatedResponse(c, http.MsgCreatedSuccess, response)
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
		return utils.BadRequestResponse(c, "Invalid company ID")
	}

	// Get the company from the service
	companyData, err := h.companyService.GetCompany(ctx, int64(companyID))
	if err != nil {
		return utils.NotFoundResponse(c, http.ErrNotFound)
	}

	// Map to response DTO
	response := mapper.ToCompanyResponse(companyData)
	return utils.SuccessResponse(c, http.MsgFetchedSuccess, response)
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
		return utils.BadRequestResponse(c, http.ErrInvalidRequest)
	}

	companyData, err := h.companyService.GetCompanyBySlug(ctx, slug)
	if err != nil {
		return utils.NotFoundResponse(c, http.ErrNotFound)
	}
	responseDTO := mapper.ToCompanyResponse(companyData)
	return utils.SuccessResponse(c, http.MsgFetchedSuccess, responseDTO)
}

// UpdateCompany godoc
// @Summary Update a company profile
// @Description Update company profile with banner, logo, and editable information. Company name, location (country/province/city), employee count, and industry are read-only (set during company creation). Full address is fetched from company creation but can be edited here.
// @Tags companies
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Param banner formData file false "Company banner image (800x300px, jpg/jpeg/png)"
// @Param logo formData file false "Company logo image (120x120px, jpg/jpeg/png)"
// @Param full_address formData string false "Full address (dari data company saat create, bisa di-edit)"
// @Param short_description formData string true "Deskripsi Singkat - Visi dan Misi Perusahaan (max 1000 chars)"
// @Param website_url formData string false "Website URL"
// @Param instagram_url formData string false "Instagram URL"
// @Param facebook_url formData string false "Facebook URL"
// @Param linkedin_url formData string false "LinkedIn URL"
// @Param twitter_url formData string false "Twitter URL"
// @Param company_description formData string true "Company description (Deskripsi Perusahaan)"
// @Param company_culture formData string false "Company culture (Budaya Perusahaan)"
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
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, http.ErrUnauthorized, "userID not found in context")
	}

	// Get company ID from the URL
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, "Invalid company ID")
	}

	// Check if user has permission (owner or admin)
	hasPermission, err := h.companyService.CheckEmployerPermission(ctx, userID, int64(companyID), "admin")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check user permission", err.Error())
	}
	if !hasPermission {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to update this company. Only company owner or admin company can perform this action.", "")
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		// If multipart form parsing fails, try to parse as JSON (backward compatibility)
		var req request.UpdateCompanyRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, http.ErrInvalidRequest)
		}

		// Validate request
		if err := utils.ValidateStruct(&req); err != nil {
			errors := utils.FormatValidationErrors(err)
			return utils.ValidationErrorResponse(c, http.ErrValidationFailed, errors)
		}

		// Convert request to domain request
		domainReq := &company.UpdateCompanyRequest{
			FullAddress:        req.FullAddress,
			ShortDescription:   req.ShortDescription,
			WebsiteURL:         req.WebsiteURL,
			InstagramURL:       req.InstagramURL,
			FacebookURL:        req.FacebookURL,
			LinkedinURL:        req.LinkedinURL,
			TwitterURL:         req.TwitterURL,
			CompanyDescription: req.CompanyDescription,
			CompanyCulture:     req.CompanyCulture,
		}

		// Call service without files (backward compatibility)
		if err := h.companyService.UpdateCompany(ctx, int64(companyID), domainReq, nil, nil); err != nil {
			return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
		}

		return utils.SuccessResponse(c, http.MsgUpdatedSuccess, nil)
	}

	// Get form values
	// NOTE: company_name, country, province, city, employee_count, industry
	// tidak perlu di-parse karena read-only (sudah ada saat create company)
	fullAddress := c.FormValue("full_address")
	shortDescription := c.FormValue("short_description")
	websiteURL := c.FormValue("website_url")
	instagramURL := c.FormValue("instagram_url")
	facebookURL := c.FormValue("facebook_url")
	linkedinURL := c.FormValue("linkedin_url")
	twitterURL := c.FormValue("twitter_url")
	companyDescription := c.FormValue("company_description")
	companyCulture := c.FormValue("company_culture")

	// Build request
	req := &request.UpdateCompanyRequest{}

	if fullAddress != "" {
		req.FullAddress = &fullAddress
	}
	if shortDescription != "" {
		req.ShortDescription = &shortDescription
	}
	if websiteURL != "" {
		req.WebsiteURL = &websiteURL
	}
	if instagramURL != "" {
		req.InstagramURL = &instagramURL
	}
	if facebookURL != "" {
		req.FacebookURL = &facebookURL
	}
	if linkedinURL != "" {
		req.LinkedinURL = &linkedinURL
	}
	if twitterURL != "" {
		req.TwitterURL = &twitterURL
	}
	if companyDescription != "" {
		req.CompanyDescription = &companyDescription
	}
	if companyCulture != "" {
		req.CompanyCulture = &companyCulture
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, http.ErrValidationFailed, errors)
	}

	// Sanitize fields
	if req.FullAddress != nil {
		sanitized := utils.SanitizeString(*req.FullAddress)
		req.FullAddress = &sanitized
	}
	if req.ShortDescription != nil {
		sanitized := utils.SanitizeString(*req.ShortDescription)
		req.ShortDescription = &sanitized
	}
	if req.CompanyDescription != nil {
		sanitized := utils.SanitizeHTML(*req.CompanyDescription)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, http.ErrPotentialXSS)
		}
		req.CompanyDescription = &sanitized
	}
	if req.CompanyCulture != nil {
		sanitized := utils.SanitizeHTML(*req.CompanyCulture)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, http.ErrPotentialXSS)
		}
		req.CompanyCulture = &sanitized
	}

	// Get banner file (optional)
	var bannerFile *multipart.FileHeader
	if files := form.File["banner"]; len(files) > 0 {
		bannerFile = files[0]
	}

	// Get logo file (optional)
	var logoFile *multipart.FileHeader
	if files := form.File["logo"]; len(files) > 0 {
		logoFile = files[0]
	}

	// Convert request to domain request
	domainReq := &company.UpdateCompanyRequest{
		FullAddress:        req.FullAddress,
		ShortDescription:   req.ShortDescription,
		WebsiteURL:         req.WebsiteURL,
		InstagramURL:       req.InstagramURL,
		FacebookURL:        req.FacebookURL,
		LinkedinURL:        req.LinkedinURL,
		TwitterURL:         req.TwitterURL,
		CompanyDescription: req.CompanyDescription,
		CompanyCulture:     req.CompanyCulture,
	}

	// Update the company profile with files
	if err := h.companyService.UpdateCompany(ctx, int64(companyID), domainReq, bannerFile, logoFile); err != nil {
		return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
	}

	return utils.SuccessResponse(c, http.MsgUpdatedSuccess, nil)
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
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, http.ErrUnauthorized, "userID not found in context")
	}

	// Get company ID from the URL
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, http.ErrInvalidCompanyID)
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
		return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
	}

	return utils.SuccessResponse(c, http.MsgDeletedSuccess, nil)
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
		return utils.BadRequestResponse(c, http.ErrInvalidCompanyID)
	}

	file := middleware.GetUploadedFile(c)
	if file == nil {
		return utils.BadRequestResponse(c, http.ErrNoFileUploaded)
	}

	url, err := h.companyService.UploadLogo(ctx, int64(companyID), file)
	if err != nil {
		return utils.InternalServerErrorResponse(c, http.ErrFileUploadFailed)
	}
	return utils.CreatedResponse(c, http.MsgUploadSuccess, fiber.Map{"logo_url": url})
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
		return utils.BadRequestResponse(c, http.ErrInvalidCompanyID)
	}

	file := middleware.GetUploadedFile(c)
	if file == nil {
		return utils.BadRequestResponse(c, http.ErrNoFileUploaded)
	}

	url, err := h.companyService.UploadBanner(ctx, int64(companyID), file)
	if err != nil {
		return utils.InternalServerErrorResponse(c, http.ErrFileUploadFailed)
	}
	return utils.CreatedResponse(c, http.MsgUploadSuccess, fiber.Map{"banner_url": url})
}

// GetCompanyVerificationStatus godoc
// @Summary Get company verification status
// @Description Get company verification status (verified or pending)
// @Tags companies
// @Accept json
// @Produce json
// @Param id path int true "Company ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/verification-status [get]
func (h *CompanyBasicHandler) GetCompanyVerificationStatus(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, http.ErrInvalidID)
	}

	// Get company basic info
	comp, err := h.companyService.GetCompany(ctx, companyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, http.ErrCompanyNotFound, err.Error())
	}

	// Get verification details
	verification, err := h.companyService.GetVerificationStatus(ctx, companyID)

	resp := fiber.Map{
		"id":           comp.ID,
		"company_name": comp.CompanyName,
		"verified":     comp.Verified,
		"verified_at":  comp.VerifiedAt,
	}

	// Add verification details if exists
	if err == nil && verification != nil {
		resp["status"] = verification.Status
		resp["verification_score"] = verification.VerificationScore
		resp["verification_notes"] = verification.VerificationNotes
		resp["npwp_number"] = verification.NPWPNumber
		resp["nib_number"] = verification.NIBNumber
		resp["reviewed_at"] = verification.ReviewedAt
		resp["verification_expiry"] = verification.VerificationExpiry
		resp["badge_granted"] = verification.BadgeGranted
		resp["rejection_reason"] = verification.RejectionReason
	} else {
		// No verification record yet
		resp["status"] = "not_requested"
	}

	return utils.SuccessResponse(c, "Verification status retrieved successfully", resp)
}

// GetMyCompanyVerificationStatus godoc
// @Summary Get my company verification status
// @Description Get verification status of the authenticated user's company (automatic, no ID needed)
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/me/verification-status [get]
func (h *CompanyBasicHandler) GetMyCompanyVerificationStatus(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get authenticated user ID from JWT
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, http.ErrUnauthorized, "userID not found in context")
	}

	// Get user's companies
	companies, err := h.companyService.GetUserCompanies(ctx, userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
	}

	if len(companies) == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, http.ErrCompanyNotFound, "You don't have any company registered")
	}

	// Get first company (users can only have 1 company based on business rule)
	comp := companies[0]

	// Get verification details
	verification, err := h.companyService.GetVerificationStatus(ctx, comp.ID)

	resp := fiber.Map{
		"id":           comp.ID,
		"company_name": comp.CompanyName,
		"verified":     comp.Verified,
		"verified_at":  comp.VerifiedAt,
	}

	// Add verification details if exists
	if err == nil && verification != nil {
		resp["status"] = verification.Status
		resp["verification_score"] = verification.VerificationScore
		resp["verification_notes"] = verification.VerificationNotes
		resp["npwp_number"] = verification.NPWPNumber
		resp["nib_number"] = verification.NIBNumber
		resp["reviewed_at"] = verification.ReviewedAt
		resp["verification_expiry"] = verification.VerificationExpiry
		resp["badge_granted"] = verification.BadgeGranted
		resp["rejection_reason"] = verification.RejectionReason
	} else {
		// No verification record yet
		resp["status"] = "not_requested"
	}

	return utils.SuccessResponse(c, "Verification status retrieved successfully", resp)
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
		return utils.BadRequestResponse(c, http.ErrInvalidCompanyID)
	}

	if err := h.companyService.DeleteLogo(ctx, int64(companyID)); err != nil {
		return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, http.MsgDeletedSuccess, nil)
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
		return utils.BadRequestResponse(c, http.ErrInvalidCompanyID)
	}

	if err := h.companyService.DeleteBanner(ctx, int64(companyID)); err != nil {
		return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, http.MsgDeletedSuccess, nil)
}

// GetMyCompanies godoc
// @Summary Get my companies
// @Description Get all companies where the authenticated user is a member with full company details
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]response.CompanyDetailResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/my-companies [get]
func (h *CompanyBasicHandler) GetMyCompanies(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get authenticated user ID from JWT using middleware helper
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, http.ErrUnauthorized, "userID not found in context")
	}

	// Get companies where user is a member
	companies, err := h.companyService.GetUserCompanies(ctx, userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
	}

	// Map to response DTOs with full company details
	responses := make([]response.CompanyDetailResponse, 0, len(companies))
	for _, comp := range companies {
		if detail := mapper.ToCompanyDetailResponse(&comp); detail != nil {
			responses = append(responses, *detail)
		}
	}

	return utils.SuccessResponse(c, http.MsgFetchedSuccess, responses)
}

// GetMyAddresses godoc
// @Summary Get my company addresses
// @Description Get list of addresses for the authenticated employer's company
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]response.CompanyAddressResponse}
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/me/addresses [get]
func (h *CompanyBasicHandler) GetMyAddresses(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get authenticated user ID from JWT
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, http.ErrUnauthorized, "userID not found in context")
	}

	// Get companies where user is a member
	companies, err := h.companyService.GetUserCompanies(ctx, userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
	}

	// If no company found
	if len(companies) == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, http.ErrCompanyNotFound, "User is not affiliated with any company")
	}

	// Get the first company (user's primary company)
	// In future, we can support multiple companies
	company := companies[0]

	// Build address response
	// Currently, company has single embedded address
	addresses := make([]response.CompanyAddressResponse, 0)

	// Only return address if FullAddress is not empty
	if company.FullAddress != "" {
		addressResp := response.CompanyAddressResponse{
			ID:            company.ID, // Using company ID as address ID
			AlamatLengkap: company.FullAddress,
		}

		// Include coordinates if available
		if company.Latitude != nil {
			addressResp.Latitude = *company.Latitude
		}
		if company.Longitude != nil {
			addressResp.Longitude = *company.Longitude
		}

		addresses = append(addresses, addressResp)
	}

	return utils.SuccessResponse(c, http.MsgFetchedSuccess, addresses)
}

// RequestVerification godoc
// @Summary Request company verification
// @Description Submit a verification request for the company with NPWP and optional NIB
// @Tags companies
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Param npwp_number formData string true "Nomor NPWP Perusahaan (Required)"
// @Param nib_number formData string false "13 Digit NIB Perusahaan (Optional)"
// @Param npwp_file formData file true "Upload NPWP Perusahaan (pdf, jpg, jpeg, png, max 10MB)"
// @Param additional_documents formData file false "Dokumen Tambahan (max 5 files, each max 10MB)"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/request-verification [post]
func (h *CompanyBasicHandler) RequestVerification(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get company ID from URL params
	companyID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, http.ErrInvalidCompanyID, err.Error())
	}

	// Get authenticated user from context
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, http.ErrUnauthorized, "User ID not found in context")
	}

	// Verify company exists
	comp, err := h.companyService.GetCompany(ctx, companyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, http.ErrCompanyNotFound, err.Error())
	}

	// Check if company is already verified
	if comp.Verified {
		return utils.ErrorResponse(c, fiber.StatusConflict, "Company already verified", "This company is already verified")
	}

	// Check if verification request already exists
	verificationStatus, err := h.companyService.GetVerificationStatus(ctx, companyID)
	if err == nil && verificationStatus != nil {
		// Verification record exists, check status
		if verificationStatus.Status == "pending" || verificationStatus.Status == "under_review" {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Verification request already submitted",
				"A verification request is already pending for this company")
		}
	}

	// Get employer user ID
	employerUserID, err := h.companyService.GetEmployerUserID(ctx, userID, companyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied",
			http.ErrNotCompanyMember)
	}

	// Parse multipart form
	var req request.RequestVerificationRequest

	// Get NPWP number (required)
	npwpNumber := c.FormValue("npwp_number")
	if npwpNumber == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "NPWP number is required", "npwp_number field is required")
	}
	req.NPWPNumber = &npwpNumber

	// Get NIB number (optional)
	nibNumber := c.FormValue("nib_number")
	if nibNumber != "" {
		req.NIBNumber = &nibNumber
	}

	// Sanitize inputs
	*req.NPWPNumber = utils.SanitizeString(*req.NPWPNumber)
	if req.NIBNumber != nil {
		*req.NIBNumber = utils.SanitizeString(*req.NIBNumber)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, http.ErrValidationFailed, err.Error())
	}

	// Get NPWP file (required)
	npwpFile, err := c.FormFile("npwp_file")
	if err != nil || npwpFile == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "NPWP file is required", "npwp_file must be uploaded")
	}

	// Validate NPWP file size (max 10MB)
	if npwpFile.Size > 10*1024*1024 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "NPWP file too large", "Maximum file size is 10MB")
	}

	// Get additional documents (optional, max 5 files)
	form, err := c.MultipartForm()
	var additionalFiles []*multipart.FileHeader
	if err == nil && form != nil {
		if files, ok := form.File["additional_documents"]; ok {
			// Limit to 5 files
			maxFiles := 5
			if len(files) > maxFiles {
				files = files[:maxFiles]
			}

			// Validate each file size
			for _, file := range files {
				if file.Size > 10*1024*1024 {
					return utils.ErrorResponse(c, fiber.StatusBadRequest, "File too large",
						fmt.Sprintf("File %s exceeds 10MB limit", file.Filename))
				}
				additionalFiles = append(additionalFiles, file)
			}
		}
	}

	// Request verification with documents
	if err := h.companyService.RequestVerification(
		ctx,
		companyID,
		employerUserID,
		*req.NPWPNumber,
		req.NIBNumber,
		npwpFile,
		additionalFiles,
	); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, http.ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, "Verification request submitted successfully", fiber.Map{
		"company_id": companyID,
		"status":     "pending",
		"message":    "Your verification request has been submitted and will be reviewed by our admin team",
	})
}
