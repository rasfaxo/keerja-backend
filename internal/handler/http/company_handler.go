package http

import (
	"strconv"
	"strings"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type CompanyHandler struct {
	companyService company.CompanyService
}

func NewCompanyHandler(companyService company.CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: companyService}
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
func (h *CompanyHandler) ListCompanies(c *fiber.Ctx) error {
	ctx := c.Context()

	var q request.CompanySearchRequest
	if err := c.QueryParser(&q); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters", err.Error())
	}
	if err := utils.ValidateStruct(&q); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}
	q.Page, q.Limit = utils.ValidatePagination(q.Page, q.Limit, 100)

	filter := &company.CompanyFilter{
		Verified:    q.IsVerified,
		Page:        q.Page,
		Limit:       q.Limit,
		SortBy:      q.SortBy,
		SortOrder:   q.SortOrder,
	}
	if q.Industry != "" {
		v := q.Industry; filter.Industry = &v
	}
	if q.CompanyType != "" {
		v := q.CompanyType; filter.CompanyType = &v
	}
	if q.SizeCategory != "" {
		v := q.SizeCategory; filter.SizeCategory = &v
	}
	if q.Location != "" {
		v := q.Location; filter.City = &v
	}
	if q.Query != "" {
		v := q.Query; filter.SearchQuery = &v
	}

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
	return utils.SuccessResponseWithMeta(c, "Companies retrieved successfully", payload, meta)
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
func (h *CompanyHandler) CreateCompany(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req request.RegisterCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
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
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create company", err.Error())
	}

	// Map to response DTO
	response := mapper.ToCompanyResponse(createdCompany)
	return utils.CreatedResponse(c, "Company created successfully", response)
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
func (h *CompanyHandler) GetCompany(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get company ID from the URL
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}

	// Get the company from the service
	companyData, err := h.companyService.GetCompany(ctx, int64(companyID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Company not found", err.Error())
	}

	// Map to response DTO
	response := mapper.ToCompanyResponse(companyData)
	return utils.SuccessResponse(c, "Company retrieved successfully", response)
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
func (h *CompanyHandler) GetCompanyBySlug(c *fiber.Ctx) error {
	ctx := c.Context()
	slug := strings.TrimSpace(c.Params("slug"))
	if slug == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company slug", "")
	}

	companyData, err := h.companyService.GetCompanyBySlug(ctx, slug)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Company not found", err.Error())
	}
	responseDTO := mapper.ToCompanyResponse(companyData)
	return utils.SuccessResponse(c, "Company retrieved successfully", responseDTO)
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
func (h *CompanyHandler) UpdateCompany(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get company ID from the URL
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}

	// Parse request body
	var req request.UpdateCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
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
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update company", err.Error())
	}

	return utils.SuccessResponse(c, "Company updated successfully", nil)
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
func (h *CompanyHandler) DeleteCompany(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get company ID from the URL
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}

	// Delete the company profile
	if err := h.companyService.DeleteCompany(ctx, int64(companyID)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete company", err.Error())
	}

	return utils.SuccessResponse(c, "Company deleted successfully", nil)
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
func (h *CompanyHandler) UploadLogo(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}

	file := middleware.GetUploadedFile(c)
	if file == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No file uploaded", "")
	}

	url, err := h.companyService.UploadLogo(ctx, int64(companyID), file)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload logo", err.Error())
	}
	return utils.CreatedResponse(c, "Logo uploaded successfully", fiber.Map{"logo_url": url})
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
func (h *CompanyHandler) UploadBanner(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}

	file := middleware.GetUploadedFile(c)
	if file == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No file uploaded", "")
	}

	url, err := h.companyService.UploadBanner(ctx, int64(companyID), file)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload banner", err.Error())
	}
	return utils.CreatedResponse(c, "Banner uploaded successfully", fiber.Map{"banner_url": url})
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
func (h *CompanyHandler) DeleteLogo(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}

	if err := h.companyService.DeleteLogo(ctx, int64(companyID)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete logo", err.Error())
	}
	return utils.SuccessResponse(c, "Logo deleted successfully", nil)
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
func (h *CompanyHandler) DeleteBanner(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}

	if err := h.companyService.DeleteBanner(ctx, int64(companyID)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete banner", err.Error())
	}
	return utils.SuccessResponse(c, "Banner deleted successfully", nil)
}

// GetCompanyProfile godoc
// @Summary Get company profile
// @Tags companies
// @Produce json
// @Param id path int true "Company ID"
// @Success 200 {object} utils.Response{data=response.CompanyProfileResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/profile [get]
func (h *CompanyHandler) GetProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}

	profile, err := h.companyService.GetProfile(ctx, int64(companyID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get profile", err.Error())
	}
	resp := mapper.ToCompanyProfileResponse(profile)
	return utils.SuccessResponse(c, "Company profile retrieved successfully", resp)
}

// UpdateCompanyProfile godoc
// @Summary Update company profile
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Param request body request.UpdateCompanyProfileRequest true "Update company profile request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/profile [put]
func (h *CompanyHandler) UpdateProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}

	var req request.UpdateCompanyProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	// Map DTO -> domain UpdateProfileRequest
	var values []string
	if req.CoreValues != nil && strings.TrimSpace(*req.CoreValues) != "" {
		for _, v := range strings.Split(*req.CoreValues, ",") {
			v = strings.TrimSpace(v)
			if v != "" {
				values = append(values, v)
			}
		}
	}
	social := map[string]string{}
	if req.FacebookURL != nil && *req.FacebookURL != "" { social["facebook"] = *req.FacebookURL }
	if req.TwitterURL != nil && *req.TwitterURL != "" { social["twitter"] = *req.TwitterURL }
	if req.LinkedinURL != nil && *req.LinkedinURL != "" { social["linkedin"] = *req.LinkedinURL }
	if req.InstagramURL != nil && *req.InstagramURL != "" { social["instagram"] = *req.InstagramURL }
	if req.YoutubeURL != nil && *req.YoutubeURL != "" { social["youtube"] = *req.YoutubeURL }
	if len(social) == 0 { social = nil }

	domainReq := &company.UpdateProfileRequest{
		LongDescription: req.Description,
		Mission:         req.Mission,
		Vision:          req.Vision,
		Values:          values,
		SEOTitle:        nil,
		SEOKeywords:     nil,
		SEODescription:  nil,
		SocialLinks:     social,
	}

	if err := h.companyService.UpdateProfile(ctx, int64(companyID), domainReq); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update profile", err.Error())
	}
	return utils.SuccessResponse(c, "Company profile updated successfully", nil)
}

// PublishCompanyProfile godoc
// @Summary Publish company profile
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/profile/publish [post]
func (h *CompanyHandler) PublishProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}
	if err := h.companyService.PublishProfile(ctx, int64(companyID)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to publish profile", err.Error())
	}
	return utils.SuccessResponse(c, "Company profile published successfully", fiber.Map{"published": true})
}

// UnpublishCompanyProfile godoc
// @Summary Unpublish company profile
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/profile/unpublish [post]
func (h *CompanyHandler) UnpublishProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}
	if err := h.companyService.UnpublishProfile(ctx, int64(companyID)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to unpublish profile", err.Error())
	}
	return utils.SuccessResponse(c, "Company profile unpublished successfully", fiber.Map{"published": false})
}

// FollowCompany godoc
// @Summary Follow a company
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/follow [post]
func (h *CompanyHandler) FollowCompany(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}

	if err := h.companyService.FollowCompany(ctx, int64(companyID), userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to follow company", err.Error())
	}
	return utils.SuccessResponse(c, "Company followed successfully", fiber.Map{"followed": true})
}

// UnfollowCompany godoc
// @Summary Unfollow a company
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/follow [delete]
func (h *CompanyHandler) UnfollowCompany(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}

	if err := h.companyService.UnfollowCompany(ctx, int64(companyID), userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to unfollow company", err.Error())
	}
	return utils.SuccessResponse(c, "Company unfollowed successfully", fiber.Map{"followed": false})
}

// GetFollowers godoc
// @Summary Get company followers
// @Tags companies
// @Produce json
// @Param id path int true "Company ID"
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/followers [get]
func (h *CompanyHandler) GetFollowers(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	page, limit = utils.ValidatePagination(page, limit, 100)

	followers, total, err := h.companyService.GetFollowers(ctx, int64(companyID), page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get followers", err.Error())
	}

	resp := make([]response.CompanyFollowerResponse, 0, len(followers))
	for _, f := range followers {
		fr := mapper.ToCompanyFollowerResponse(&f)
		if fr != nil {
			resp = append(resp, *fr)
		}
	}
	meta := utils.GetPaginationMeta(page, limit, total)
	return utils.SuccessResponseWithMeta(c, "Followers retrieved successfully", fiber.Map{"followers": resp}, meta)
}

// AddReview godoc
// @Summary Add a company review
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Param request body request.AddReviewRequest true "Add review request"
// @Success 201 {object} utils.Response{data=response.CompanyReviewResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/reviews [post]
func (h *CompanyHandler) AddReview(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}

	var req request.AddReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	domainReq := &company.AddReviewRequest{
		CompanyID:          int64(companyID),
		UserID:             userID,
		ReviewerType:       req.ReviewerType,
		PositionTitle:      req.PositionTitle,
		EmploymentPeriod:   req.EmploymentPeriod,
		RatingOverall:      req.RatingOverall,
		RatingCulture:      req.RatingCulture,
		RatingWorkLife:     req.RatingWorkLife,
		RatingSalary:       req.RatingSalary,
		RatingManagement:   req.RatingManagement,
		Pros:               req.Pros,
		Cons:               req.Cons,
		AdviceToManagement: req.AdviceToManagement,
		IsAnonymous:        req.IsAnonymous,
		RecommendToFriend:  req.RecommendToFriend,
	}

	rev, err := h.companyService.AddReview(ctx, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to add review", err.Error())
	}
	resp := mapper.ToCompanyReviewResponse(rev)
	return utils.CreatedResponse(c, "Review added successfully", resp)
}

// UpdateReview godoc
// @Summary Update a company review
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Review ID"
// @Param request body request.UpdateReviewRequest true "Update review request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/reviews/{id} [put]
func (h *CompanyHandler) UpdateReview(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	reviewID, err := strconv.Atoi(c.Params("id"))
	if err != nil || reviewID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid review ID", "")
	}

	var req request.UpdateReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	domainReq := &company.UpdateReviewRequest{
		ReviewerType:       req.ReviewerType,
		PositionTitle:      req.PositionTitle,
		EmploymentPeriod:   req.EmploymentPeriod,
		RatingOverall:      req.RatingOverall,
		RatingCulture:      req.RatingCulture,
		RatingWorkLife:     req.RatingWorkLife,
		RatingSalary:       req.RatingSalary,
		RatingManagement:   req.RatingManagement,
		Pros:               req.Pros,
		Cons:               req.Cons,
		AdviceToManagement: req.AdviceToManagement,
		RecommendToFriend:  req.RecommendToFriend,
	}

	if err := h.companyService.UpdateReview(ctx, int64(reviewID), userID, domainReq); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update review", err.Error())
	}
	return utils.SuccessResponse(c, "Review updated successfully", nil)
}

// DeleteReview godoc
// @Summary Delete a company review
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param id path int true "Review ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/reviews/{id} [delete]
func (h *CompanyHandler) DeleteReview(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	reviewID, err := strconv.Atoi(c.Params("id"))
	if err != nil || reviewID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid review ID", "")
	}

	if err := h.companyService.DeleteReview(ctx, int64(reviewID), userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete review", err.Error())
	}
	return utils.SuccessResponse(c, "Review deleted successfully", nil)
}

// GetCompanyReviews godoc
// @Summary List company reviews
// @Tags companies
// @Produce json
// @Param id path int true "Company ID"
// @Param status query string false "Review status"
// @Param reviewer_type query string false "Reviewer type"
// @Param min_rating query number false "Minimum rating"
// @Param max_rating query number false "Maximum rating"
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param sort_by query string false "Sort by"
// @Param sort_order query string false "Sort order"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/reviews [get]
func (h *CompanyHandler) GetCompanyReviews(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	page, limit = utils.ValidatePagination(page, limit, 100)

	// Optional filters via query
	var filt company.ReviewFilter
	if v := c.Query("status"); v != "" { s := v; filt.Status = &s }
	if v := c.Query("reviewer_type"); v != "" { s := v; filt.ReviewerType = &s }
	if v := c.Query("min_rating"); v != "" {
		if r, err := strconv.ParseFloat(v, 64); err == nil { filt.MinRating = &r }
	}
	if v := c.Query("max_rating"); v != "" {
		if r, err := strconv.ParseFloat(v, 64); err == nil { filt.MaxRating = &r }
	}
	filt.Page = page
	filt.Limit = limit
	filt.SortBy = c.Query("sort_by")
	filt.SortOrder = c.Query("sort_order")

	reviews, total, err := h.companyService.GetCompanyReviews(ctx, int64(companyID), &filt)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get reviews", err.Error())
	}

	resp := make([]response.CompanyReviewResponse, 0, len(reviews))
	for _, r := range reviews {
		rr := mapper.ToCompanyReviewResponse(&r)
		if rr != nil {
			resp = append(resp, *rr)
		}
	}
	meta := utils.GetPaginationMeta(page, limit, total)
	return utils.SuccessResponseWithMeta(c, "Reviews retrieved successfully", fiber.Map{"reviews": resp}, meta)
}

// GetAverageRatings godoc
// @Summary Get company's average ratings
// @Tags companies
// @Produce json
// @Param id path int true "Company ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/ratings [get]
func (h *CompanyHandler) GetAverageRatings(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}

	ratings, err := h.companyService.GetAverageRatings(ctx, int64(companyID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get average ratings", err.Error())
	}
	return utils.SuccessResponse(c, "Average ratings retrieved successfully", ratings)
}

// GetVerifiedCompanies godoc
// @Summary Get verified companies
// @Tags companies
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} utils.Response{data=response.CompanyListResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/verified [get]
func (h *CompanyHandler) GetVerifiedCompanies(c *fiber.Ctx) error {
	ctx := c.Context()
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	page, limit = utils.ValidatePagination(page, limit, 100)

	companies, total, err := h.companyService.GetVerifiedCompanies(ctx, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get verified companies", err.Error())
	}

	respList := make([]response.CompanyResponse, 0, len(companies))
	for _, comp := range companies {
		cr := mapper.ToCompanyResponse(&comp)
		if cr != nil {
			respList = append(respList, *cr)
		}
	}
	meta := utils.GetPaginationMeta(page, limit, total)
	payload := response.CompanyListResponse{Companies: respList}
	return utils.SuccessResponseWithMeta(c, "Verified companies retrieved successfully", payload, meta)
}

// GetTopRatedCompanies godoc
// @Summary Get top-rated companies
// @Tags companies
// @Produce json
// @Param limit query int false "Max companies to return"
// @Success 200 {object} utils.Response{data=response.CompanyListResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/top-rated [get]
func (h *CompanyHandler) GetTopRatedCompanies(c *fiber.Ctx) error {
	ctx := c.Context()
	limit := c.QueryInt("limit", 10)
	if limit <= 0 {
		limit = 10
	}

	companies, err := h.companyService.GetTopRatedCompanies(ctx, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get top-rated companies", err.Error())
	}

	respList := make([]response.CompanyResponse, 0, len(companies))
	for _, comp := range companies {
		cr := mapper.ToCompanyResponse(&comp)
		if cr != nil {
			respList = append(respList, *cr)
		}
	}
	payload := response.CompanyListResponse{Companies: respList}
	return utils.SuccessResponse(c, "Top-rated companies retrieved successfully", payload)
}

// GetCompanyStats godoc
// @Summary Get company statistics (jobs, applications, followers, reviews, etc.)
// @Tags companies
// @Produce json
// @Param id path int true "Company ID"
// @Success 200 {object} utils.Response{data=response.CompanyStatsResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/stats [get]
func (h *CompanyHandler) GetCompanyStats(c *fiber.Ctx) error {
	ctx := c.Context()
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", "")
	}

	stats, err := h.companyService.GetCompanyStats(ctx, int64(companyID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get company stats", err.Error())
	}
	return utils.SuccessResponse(c, "Company stats retrieved successfully", stats)
}

// GetFollowedCompanies godoc
// @Summary Get companies followed by current user
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} utils.Response{data=response.CompanyListResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/companies/following [get]
func (h *CompanyHandler) GetFollowedCompanies(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	page, limit = utils.ValidatePagination(page, limit, 100)

	companies, total, err := h.companyService.GetFollowedCompanies(ctx, userID, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get followed companies", err.Error())
	}
	respList := make([]response.CompanyResponse, 0, len(companies))
	for _, comp := range companies {
		cr := mapper.ToCompanyResponse(&comp)
		if cr != nil {
			respList = append(respList, *cr)
		}
	}
	meta := utils.GetPaginationMeta(page, limit, total)
	payload := response.CompanyListResponse{Companies: respList}
	return utils.SuccessResponseWithMeta(c, "Followed companies retrieved successfully", payload, meta)
}
