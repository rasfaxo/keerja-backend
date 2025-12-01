package companyhandler

import (
	"strconv"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/handler/http/common"
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

func (h *CompanyReviewHandler) AddReview(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	var req request.AddReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errs)
	}

	// CRITICAL: Sanitize user-generated content to prevent XSS
	if req.Pros != nil {
		sanitized := utils.SanitizeHTML(*req.Pros)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, common.ErrPotentialXSS)
		}
		req.Pros = &sanitized
	}
	if req.Cons != nil {
		sanitized := utils.SanitizeHTML(*req.Cons)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, common.ErrPotentialXSS)
		}
		req.Cons = &sanitized
	}
	if req.AdviceToManagement != nil {
		sanitized := utils.SanitizeHTML(*req.AdviceToManagement)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, common.ErrPotentialXSS)
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
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}
	resp := mapper.ToCompanyReviewResponse(rev)
	return utils.CreatedResponse(c, common.MsgCreatedSuccess, resp)
}

func (h *CompanyReviewHandler) UpdateReview(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	reviewID, err := strconv.Atoi(c.Params("id"))
	if err != nil || reviewID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	var req request.UpdateReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errs)
	}

	// CRITICAL: Sanitize user-generated content
	if req.Pros != nil {
		sanitized := utils.SanitizeHTML(*req.Pros)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, common.ErrPotentialXSS)
		}
		req.Pros = &sanitized
	}
	if req.Cons != nil {
		sanitized := utils.SanitizeHTML(*req.Cons)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, common.ErrPotentialXSS)
		}
		req.Cons = &sanitized
	}
	if req.AdviceToManagement != nil {
		sanitized := utils.SanitizeHTML(*req.AdviceToManagement)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, common.ErrPotentialXSS)
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
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, common.MsgUpdatedSuccess, nil)
}

func (h *CompanyReviewHandler) DeleteReview(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	reviewID, err := strconv.Atoi(c.Params("id"))
	if err != nil || reviewID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	if err := h.companyService.DeleteReview(ctx, int64(reviewID), userID); err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, common.MsgDeletedSuccess, nil)
}

func (h *CompanyReviewHandler) GetCompanyReviews(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
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
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}

	resp := mapper.MapEntities[company.CompanyReview, response.CompanyReviewResponse](reviews, func(r *company.CompanyReview) *response.CompanyReviewResponse {
		return mapper.ToCompanyReviewResponse(r)
	})
	meta := utils.GetPaginationMeta(page, limit, total)
	return utils.SuccessResponseWithMeta(c, common.MsgFetchedSuccess, fiber.Map{"reviews": resp}, meta)
}

func (h *CompanyReviewHandler) GetAverageRatings(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	ratings, err := h.companyService.GetAverageRatings(ctx, int64(companyID))
	if err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, common.MsgFetchedSuccess, ratings)
}
