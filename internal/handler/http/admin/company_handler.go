package admin

import (
	"strconv"

	"keerja-backend/internal/domain/admin"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// CompanyHandler handles admin company management endpoints
type CompanyHandler struct {
	adminCompanyService admin.AdminCompanyService
}

// NewCompanyHandler creates a new admin company handler
func NewCompanyHandler(adminCompanyService admin.AdminCompanyService) *CompanyHandler {
	return &CompanyHandler{
		adminCompanyService: adminCompanyService,
	}
}

// =============================================================================
// Task 2.1: List Companies (Moderation Queue)
// GET /api/v1/admin/companies
// =============================================================================

// ListCompanies retrieves paginated list of companies with filters
// @Summary List companies for moderation
// @Description Get paginated list of companies with search, filter, and sort
// @Tags Admin - Company
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param search query string false "Search by company name, email, legal name"
// @Param status query string false "Filter by verification status"
// @Param verified query boolean false "Filter by verified status"
// @Param is_active query boolean false "Filter by active status"
// @Param industry_id query int false "Filter by industry ID"
// @Param company_size_id query int false "Filter by company size ID"
// @Param province_id query int false "Filter by province ID"
// @Param city_id query int false "Filter by city ID"
// @Param created_from query string false "Filter by created date from (YYYY-MM-DD)"
// @Param created_to query string false "Filter by created date to (YYYY-MM-DD)"
// @Param sort_by query string false "Sort by field" Enums(company_name, created_at, verified_at, updated_at)
// @Param sort_order query string false "Sort order" Enums(asc, desc)
// @Success 200 {object} response.AdminCompaniesListResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/admin/companies [get]
func (h *CompanyHandler) ListCompanies(c *fiber.Ctx) error {
	// Parse query parameters
	var req request.AdminGetCompaniesRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	}

	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	// Convert to service request
	serviceReq := &admin.AdminCompanyListRequest{
		Page:               req.Page,
		Limit:              req.Limit,
		Search:             req.Search,
		Status:             req.Status,
		VerificationStatus: req.VerificationStatus,
		IndustryID:         req.IndustryID,
		CompanySizeID:      req.CompanySizeID,
		ProvinceID:         req.ProvinceID,
		CityID:             req.CityID,
		Verified:           req.Verified,
		IsActive:           req.IsActive,
		CreatedFrom:        req.CreatedFrom,
		CreatedTo:          req.CreatedTo,
		SortBy:             req.SortBy,
		SortOrder:          req.SortOrder,
	}

	// Call service
	result, err := h.adminCompanyService.ListCompanies(c.Context(), serviceReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch companies", err.Error())
	}

	// Map to response
	companies := make([]response.AdminCompanyListItemResponse, 0, len(result.Companies))
	for _, item := range result.Companies {
		companies = append(companies, mapper.ToAdminCompanyListItemResponse(&item))
	}

	resp := response.AdminCompaniesListResponse{
		Companies: companies,
		Meta: response.PaginationMeta{
			CurrentPage: result.Page,
			PerPage:     result.Limit,
			Total:       result.Total,
			TotalPages:  result.TotalPages,
			HasNext:     result.HasNext,
			HasPrev:     result.HasPrev,
		},
	}

	return utils.SuccessResponse(c, "Companies retrieved successfully", resp)
}

// =============================================================================
// Task 2.2: Get Company Detail
// GET /api/v1/admin/companies/:id
// =============================================================================

// GetCompanyDetail retrieves full company details for moderation
// @Summary Get company detail
// @Description Get full company details including verification, documents, and stats
// @Tags Admin - Company
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} response.AdminCompanyDetailResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/admin/companies/{id} [get]
func (h *CompanyHandler) GetCompanyDetail(c *fiber.Ctx) error {
	// Parse company ID
	companyID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", err.Error())
	}

	// Call service
	detail, err := h.adminCompanyService.GetCompanyDetail(c.Context(), companyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Company not found", err.Error())
	}

	// Map to response
	resp := mapper.ToAdminCompanyDetailResponse(detail)

	return utils.SuccessResponse(c, "Company detail retrieved successfully", resp)
}

// =============================================================================
// Task 2.3: Update Company Status (Approve/Reject/Suspend)
// PATCH /api/v1/admin/companies/:id/status
// =============================================================================

// UpdateCompanyStatus updates company verification status
// @Summary Update company status
// @Description Approve, reject, suspend, or blacklist a company
// @Tags Admin - Company
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Param request body request.AdminUpdateCompanyStatusRequest true "Status update request"
// @Success 200 {object} response.AdminCompanyStatusUpdateResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/admin/companies/{id}/status [patch]
func (h *CompanyHandler) UpdateCompanyStatus(c *fiber.Ctx) error {
	// Parse company ID
	companyID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", err.Error())
	}

	// Parse request body
	var req request.AdminUpdateCompanyStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	}

	// Additional validation: rejection_reason required when status is rejected
	if req.Status == "rejected" && (req.RejectionReason == nil || *req.RejectionReason == "") {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", "rejection_reason is required when status is rejected")
	}

	// Get admin ID from context (set by auth middleware)
	adminID := c.Locals("admin_id").(int64)

	// Convert to service request
	serviceReq := &admin.AdminCompanyStatusRequest{
		Status:          req.Status,
		RejectionReason: req.RejectionReason,
		Notes:           req.Notes,
		GrantBadge:      req.GrantBadge,
	}

	// Call service
	err = h.adminCompanyService.UpdateCompanyStatus(c.Context(), companyID, serviceReq, adminID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update company status", err.Error())
	}

	// Note: In real implementation, service should return updated company details
	// For now, we'll return a simple success response
	return utils.SuccessResponse(c, "Company status updated successfully", fiber.Map{
		"company_id": companyID,
		"status":     req.Status,
		"updated_by": adminID,
	})
}

// =============================================================================
// Task 2.4: Update Company Details
// PUT /api/v1/admin/companies/:id
// =============================================================================

// UpdateCompany updates company details (admin support)
// @Summary Update company details
// @Description Admin can edit company information to fix errors or typos
// @Tags Admin - Company
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Param request body request.AdminUpdateCompanyRequest true "Company update request"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/admin/companies/{id} [put]
func (h *CompanyHandler) UpdateCompany(c *fiber.Ctx) error {
	// Parse company ID
	companyID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", err.Error())
	}

	// Parse request body
	var req request.AdminUpdateCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	}

	// Get admin ID from context
	adminID := c.Locals("admin_id").(int64)

	// Convert to service request
	serviceReq := &admin.AdminUpdateCompanyRequest{
		CompanyName:        req.CompanyName,
		LegalName:          req.LegalName,
		RegistrationNumber: req.RegistrationNumber,
		IndustryID:         req.IndustryID,
		CompanySizeID:      req.CompanySizeID,
		DistrictID:         req.DistrictID,
		FullAddress:        req.FullAddress,
		Description:        req.Description,
		WebsiteURL:         req.WebsiteURL,
		EmailDomain:        req.EmailDomain,
		Phone:              req.Phone,
		About:              req.About,
		Culture:            req.Culture,
		IsActive:           req.IsActive,
		Verified:           req.Verified,
	}

	// Call service
	err = h.adminCompanyService.UpdateCompany(c.Context(), companyID, serviceReq, adminID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update company", err.Error())
	}

	return utils.SuccessResponse(c, "Company updated successfully", fiber.Map{
		"company_id": companyID,
		"updated_by": adminID,
	})
}

// =============================================================================
// Task 2.5: Delete Company
// DELETE /api/v1/admin/companies/:id
// =============================================================================

// DeleteCompany deletes a company (with validation)
// @Summary Delete company
// @Description Delete company with validation. Cannot delete if has active jobs unless force flag is used.
// @Tags Admin - Company
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Param force query boolean false "Force delete even with active jobs"
// @Param request body request.AdminDeleteCompanyRequest true "Delete request with reason"
// @Success 200 {object} response.AdminCompanyDeleteResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/admin/companies/{id} [delete]
func (h *CompanyHandler) DeleteCompany(c *fiber.Ctx) error {
	// Parse company ID
	companyID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", err.Error())
	}

	// Parse query for force flag
	force := c.QueryBool("force", false)

	// Parse request body
	var req request.AdminDeleteCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Override force from query if provided
	if force {
		req.Force = force
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	}

	// Get admin ID from context
	adminID := c.Locals("admin_id").(int64)

	// Convert to service request
	serviceReq := &admin.AdminDeleteCompanyRequest{
		Force:  req.Force,
		Reason: req.Reason,
	}

	// Call service
	err = h.adminCompanyService.DeleteCompany(c.Context(), companyID, serviceReq, adminID)
	if err != nil {
		// Check if error is due to active jobs
		if err.Error() == "company has active jobs" && !req.Force {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Cannot delete company with active jobs", "Use force=true to delete anyway")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete company", err.Error())
	}

	return utils.SuccessResponse(c, "Company deleted successfully", fiber.Map{
		"company_id": companyID,
		"deleted_by": adminID,
		"reason":     req.Reason,
	})
}

// =============================================================================
// Additional Endpoints
// =============================================================================

// GetCompanyStats retrieves company statistics
// GET /api/v1/admin/companies/:id/stats
func (h *CompanyHandler) GetCompanyStats(c *fiber.Ctx) error {
	companyID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", err.Error())
	}

	stats, err := h.adminCompanyService.GetCompanyStats(c.Context(), companyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch stats", err.Error())
	}

	return utils.SuccessResponse(c, "Company stats retrieved successfully", stats)
}

// GetDashboardStats retrieves overall dashboard statistics
// GET /api/v1/admin/dashboard/stats
func (h *CompanyHandler) GetDashboardStats(c *fiber.Ctx) error {
	stats, err := h.adminCompanyService.GetDashboardStats(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch dashboard stats", err.Error())
	}

	return utils.SuccessResponse(c, "Dashboard stats retrieved successfully", stats)
}

// GetAuditLogs retrieves company audit logs
// GET /api/v1/admin/companies/:id/audit-logs
func (h *CompanyHandler) GetAuditLogs(c *fiber.Ctx) error {
	companyID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", err.Error())
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	logs, err := h.adminCompanyService.GetAuditLogs(c.Context(), companyID, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch audit logs", err.Error())
	}

	return utils.SuccessResponse(c, "Audit logs retrieved successfully", logs)
}
