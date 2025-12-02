package companyhandler

import (
	"strconv"
	"strings"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// CompanyProfileHandler handles company profile and social features
// This includes profile management, publishing/unpublishing, and company following.
type CompanyProfileHandler struct {
	companyService company.CompanyService
}

// NewCompanyProfileHandler creates a new instance of CompanyProfileHandler
func NewCompanyProfileHandler(companyService company.CompanyService) *CompanyProfileHandler {
	return &CompanyProfileHandler{companyService: companyService}
}

func (h *CompanyProfileHandler) GetProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	profile, err := h.companyService.GetProfile(ctx, int64(companyID))
	if err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}
	resp := mapper.ToCompanyProfileResponse(profile)
	return utils.SuccessResponse(c, common.MsgFetchedSuccess, resp)
}

func (h *CompanyProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	var req request.UpdateCompanyProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errs)
	}

	// Sanitize HTML fields
	if req.Description != nil {
		sanitized := utils.SanitizeHTML(*req.Description)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, common.ErrPotentialXSS)
		}
		req.Description = &sanitized
	}
	if req.Mission != nil {
		sanitized := utils.SanitizeHTML(*req.Mission)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, common.ErrPotentialXSS)
		}
		req.Mission = &sanitized
	}
	if req.Vision != nil {
		sanitized := utils.SanitizeHTML(*req.Vision)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, common.ErrPotentialXSS)
		}
		req.Vision = &sanitized
	}
	if req.CoreValues != nil {
		sanitized := utils.SanitizeString(*req.CoreValues)
		req.CoreValues = &sanitized
	}

	// Sanitize social media URLs
	if req.FacebookURL != nil {
		sanitized := utils.SanitizeString(*req.FacebookURL)
		req.FacebookURL = &sanitized
	}
	if req.TwitterURL != nil {
		sanitized := utils.SanitizeString(*req.TwitterURL)
		req.TwitterURL = &sanitized
	}
	if req.LinkedinURL != nil {
		sanitized := utils.SanitizeString(*req.LinkedinURL)
		req.LinkedinURL = &sanitized
	}
	if req.InstagramURL != nil {
		sanitized := utils.SanitizeString(*req.InstagramURL)
		req.InstagramURL = &sanitized
	}
	if req.YoutubeURL != nil {
		sanitized := utils.SanitizeString(*req.YoutubeURL)
		req.YoutubeURL = &sanitized
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
	if req.FacebookURL != nil && *req.FacebookURL != "" {
		social["facebook"] = *req.FacebookURL
	}
	if req.TwitterURL != nil && *req.TwitterURL != "" {
		social["twitter"] = *req.TwitterURL
	}
	if req.LinkedinURL != nil && *req.LinkedinURL != "" {
		social["linkedin"] = *req.LinkedinURL
	}
	if req.InstagramURL != nil && *req.InstagramURL != "" {
		social["instagram"] = *req.InstagramURL
	}
	if req.YoutubeURL != nil && *req.YoutubeURL != "" {
		social["youtube"] = *req.YoutubeURL
	}
	if len(social) == 0 {
		social = nil
	}

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
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, common.MsgUpdatedSuccess, nil)
}

func (h *CompanyProfileHandler) PublishProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}
	if err := h.companyService.PublishProfile(ctx, int64(companyID)); err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, common.MsgUpdatedSuccess, fiber.Map{"published": true})
}

func (h *CompanyProfileHandler) UnpublishProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}
	if err := h.companyService.UnpublishProfile(ctx, int64(companyID)); err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, common.MsgUpdatedSuccess, fiber.Map{"published": false})
}

func (h *CompanyProfileHandler) FollowCompany(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	if err := h.companyService.FollowCompany(ctx, int64(companyID), userID); err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, common.MsgOperationSuccess, fiber.Map{"followed": true})
}

func (h *CompanyProfileHandler) UnfollowCompany(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	if err := h.companyService.UnfollowCompany(ctx, int64(companyID), userID); err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, common.MsgOperationSuccess, fiber.Map{"followed": false})
}

func (h *CompanyProfileHandler) GetFollowers(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	page, limit = utils.ValidatePagination(page, limit, 100)

	followers, total, err := h.companyService.GetFollowers(ctx, int64(companyID), page, limit)
	if err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}

	resp := mapper.MapEntities[company.CompanyFollower, response.CompanyFollowerResponse](followers, func(f *company.CompanyFollower) *response.CompanyFollowerResponse {
		return mapper.ToCompanyFollowerResponse(f)
	})
	meta := utils.GetPaginationMeta(page, limit, total)
	return utils.SuccessResponseWithMeta(c, common.MsgFetchedSuccess, fiber.Map{"followers": resp}, meta)
}

func (h *CompanyProfileHandler) GetFollowedCompanies(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	page, limit = utils.ValidatePagination(page, limit, 100)

	companies, total, err := h.companyService.GetFollowedCompanies(ctx, userID, page, limit)
	if err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}
	respList := mapper.MapEntities[company.Company, response.CompanyResponse](companies, func(comp *company.Company) *response.CompanyResponse {
		return mapper.ToCompanyResponse(comp)
	})
	meta := utils.GetPaginationMeta(page, limit, total)
	payload := response.CompanyListResponse{Companies: respList}
	return utils.SuccessResponseWithMeta(c, common.MsgFetchedSuccess, payload, meta)
}
