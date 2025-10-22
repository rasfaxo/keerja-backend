package http

import (
	"strconv"
	"strings"

	"keerja-backend/internal/domain/company"
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
	companyService company.CompanyService
}

// NewCompanyBasicHandler creates a new instance of CompanyBasicHandler
func NewCompanyBasicHandler(companyService company.CompanyService) *CompanyBasicHandler {
	return &CompanyBasicHandler{companyService: companyService}
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
// @Description Create a new company profile with basic details
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.RegisterCompanyRequest true "Register company request"
// @Success 201 {object} utils.Response{data=response.CompanyResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies [post]
func (h *CompanyBasicHandler) CreateCompany(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.RegisterCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
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

	// Convert to domain request
	domainReq := &company.RegisterCompanyRequest{
		CompanyName:        req.CompanyName,
		LegalName:          req.LegalName,
		RegistrationNumber: req.RegistrationNumber,
		Industry:           req.Industry,
		CompanyType:        req.CompanyType,
		SizeCategory:       req.SizeCategory,
		WebsiteURL:         req.WebsiteURL,
		EmailDomain:        req.EmailDomain,
		Phone:              req.Phone,
		Address:            req.Address,
		City:               req.City,
		Province:           req.Province,
		Country:            req.Country,
		PostalCode:         req.PostalCode,
		About:              req.About,
	}

	// Create the company
	createdCompany, err := h.companyService.RegisterCompany(ctx, domainReq)
	if err != nil {
		return utils.InternalServerErrorResponse(c, ErrFailedOperation)
	}

	// Map to response DTO
	response := mapper.ToCompanyResponse(createdCompany)
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
// @Description Update the details of a company profile
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Param request body request.UpdateCompanyRequest true "Update company request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id} [put]
func (h *CompanyBasicHandler) UpdateCompany(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get company ID from the URL
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, ErrInvalidID)
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
		Industry:           req.Industry,
		CompanyType:        req.CompanyType,
		SizeCategory:       req.SizeCategory,
		WebsiteURL:         req.WebsiteURL,
		EmailDomain:        req.EmailDomain,
		Phone:              req.Phone,
		Address:            req.Address,
		City:               req.City,
		Province:           req.Province,
		PostalCode:         req.PostalCode,
		Latitude:           req.Latitude,
		Longitude:          req.Longitude,
		About:              req.About,
		Culture:            req.Culture,
		Benefits:           req.Benefits,
	}

	// Update the company profile
	if err := h.companyService.UpdateCompany(ctx, int64(companyID), domainReq); err != nil {
		return utils.InternalServerErrorResponse(c, ErrFailedOperation)
	}

	return utils.SuccessResponse(c, MsgUpdatedSuccess, nil)
}

// DeleteCompany godoc
// @Summary Delete a company profile
// @Description Delete a company profile by ID
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id} [delete]
func (h *CompanyBasicHandler) DeleteCompany(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get company ID from the URL
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, ErrInvalidID)
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
