package companyhandler

import (
	"strconv"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/handler/http"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// CompanyReviewHandler handles company review operations
// This includes adding, updating, deleting reviews, and fetching company reviews with ratings.
type CompanyReviewHandler struct {
	companyService company.CompanyService
}

// NewCompanyReviewHandler creates a new instance of CompanyReviewHandler
func NewCompanyReviewHandler(companyService company.CompanyService) *CompanyReviewHandler {
	return &CompanyReviewHandler{companyService: companyService}
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
func (h *CompanyReviewHandler) AddReview(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, http.ErrInvalidID)
	}

	var req request.AddReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, http.ErrInvalidRequest)
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, http.ErrValidationFailed, errs)
	}

	// CRITICAL: Sanitize user-generated content to prevent XSS
	if req.Pros != nil {
		sanitized := utils.SanitizeHTML(*req.Pros)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, http.ErrPotentialXSS)
		}
		req.Pros = &sanitized
	}
	if req.Cons != nil {
		sanitized := utils.SanitizeHTML(*req.Cons)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, http.ErrPotentialXSS)
		}
		req.Cons = &sanitized
	}
	if req.AdviceToManagement != nil {
		sanitized := utils.SanitizeHTML(*req.AdviceToManagement)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, http.ErrPotentialXSS)
		}
		req.AdviceToManagement = &sanitized
	}
	if req.PositionTitle != nil {
		sanitized := utils.SanitizeString(*req.PositionTitle)
		req.PositionTitle = &sanitized
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
		return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
	}
	resp := mapper.ToCompanyReviewResponse(rev)
	return utils.CreatedResponse(c, http.MsgCreatedSuccess, resp)
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
func (h *CompanyReviewHandler) UpdateReview(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	reviewID, err := strconv.Atoi(c.Params("id"))
	if err != nil || reviewID <= 0 {
		return utils.BadRequestResponse(c, http.ErrInvalidID)
	}

	var req request.UpdateReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, http.ErrInvalidRequest)
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, http.ErrValidationFailed, errs)
	}

	// CRITICAL: Sanitize user-generated content
	if req.Pros != nil {
		sanitized := utils.SanitizeHTML(*req.Pros)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, http.ErrPotentialXSS)
		}
		req.Pros = &sanitized
	}
	if req.Cons != nil {
		sanitized := utils.SanitizeHTML(*req.Cons)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, http.ErrPotentialXSS)
		}
		req.Cons = &sanitized
	}
	if req.AdviceToManagement != nil {
		sanitized := utils.SanitizeHTML(*req.AdviceToManagement)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, http.ErrPotentialXSS)
		}
		req.AdviceToManagement = &sanitized
	}
	if req.PositionTitle != nil {
		sanitized := utils.SanitizeString(*req.PositionTitle)
		req.PositionTitle = &sanitized
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
		return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, http.MsgUpdatedSuccess, nil)
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
func (h *CompanyReviewHandler) DeleteReview(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	reviewID, err := strconv.Atoi(c.Params("id"))
	if err != nil || reviewID <= 0 {
		return utils.BadRequestResponse(c, http.ErrInvalidID)
	}

	if err := h.companyService.DeleteReview(ctx, int64(reviewID), userID); err != nil {
		return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, http.MsgDeletedSuccess, nil)
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
func (h *CompanyReviewHandler) GetCompanyReviews(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, http.ErrInvalidID)
	}
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	page, limit = utils.ValidatePagination(page, limit, 100)

	// Optional filters via query - sanitize filter values
	var filt company.ReviewFilter
	if v := c.Query("status"); v != "" {
		s := utils.SanitizeString(v)
		filt.Status = &s
	}
	if v := c.Query("reviewer_type"); v != "" {
		s := utils.SanitizeString(v)
		filt.ReviewerType = &s
	}
	if v := c.Query("min_rating"); v != "" {
		if r, err := strconv.ParseFloat(v, 64); err == nil {
			filt.MinRating = &r
		}
	}
	if v := c.Query("max_rating"); v != "" {
		if r, err := strconv.ParseFloat(v, 64); err == nil {
			filt.MaxRating = &r
		}
	}
	filt.Page = page
	filt.Limit = limit
	filt.SortBy = utils.SanitizeString(c.Query("sort_by"))
	filt.SortOrder = utils.SanitizeString(c.Query("sort_order"))

	reviews, total, err := h.companyService.GetCompanyReviews(ctx, int64(companyID), &filt)
	if err != nil {
		return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
	}

	resp := make([]response.CompanyReviewResponse, 0, len(reviews))
	for _, r := range reviews {
		rr := mapper.ToCompanyReviewResponse(&r)
		if rr != nil {
			resp = append(resp, *rr)
		}
	}
	meta := utils.GetPaginationMeta(page, limit, total)
	return utils.SuccessResponseWithMeta(c, http.MsgFetchedSuccess, fiber.Map{"reviews": resp}, meta)
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
func (h *CompanyReviewHandler) GetAverageRatings(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, http.ErrInvalidID)
	}

	ratings, err := h.companyService.GetAverageRatings(ctx, int64(companyID))
	if err != nil {
		return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, http.MsgFetchedSuccess, ratings)
}
