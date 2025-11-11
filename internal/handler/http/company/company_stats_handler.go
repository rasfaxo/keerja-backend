package companyhandler

import (
	"strconv"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/utils"
	"keerja-backend/internal/handler/http"
	"github.com/gofiber/fiber/v2"
)

// CompanyStatsHandler handles company statistics and special queries
// This includes verified companies, top-rated companies, and company statistics.
type CompanyStatsHandler struct {
	companyService company.CompanyService
}

// NewCompanyStatsHandler creates a new instance of CompanyStatsHandler
func NewCompanyStatsHandler(companyService company.CompanyService) *CompanyStatsHandler {
	return &CompanyStatsHandler{companyService: companyService}
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
func (h *CompanyStatsHandler) GetVerifiedCompanies(c *fiber.Ctx) error {
	ctx := c.Context()
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	page, limit = utils.ValidatePagination(page, limit, 100)

	companies, total, err := h.companyService.GetVerifiedCompanies(ctx, page, limit)
	if err != nil {
		return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
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
	return utils.SuccessResponseWithMeta(c, http.MsgFetchedSuccess, payload, meta)
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
func (h *CompanyStatsHandler) GetTopRatedCompanies(c *fiber.Ctx) error {
	ctx := c.Context()
	limit := c.QueryInt("limit", 10)
	if limit <= 0 {
		limit = 10
	}

	companies, err := h.companyService.GetTopRatedCompanies(ctx, limit)
	if err != nil {
		return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
	}

	respList := make([]response.CompanyResponse, 0, len(companies))
	for _, comp := range companies {
		cr := mapper.ToCompanyResponse(&comp)
		if cr != nil {
			respList = append(respList, *cr)
		}
	}
	payload := response.CompanyListResponse{Companies: respList}
	return utils.SuccessResponse(c, http.MsgFetchedSuccess, payload)
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
func (h *CompanyStatsHandler) GetCompanyStats(c *fiber.Ctx) error {
	ctx := c.Context()
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, http.ErrInvalidID)
	}

	stats, err := h.companyService.GetCompanyStats(ctx, int64(companyID))
	if err != nil {
		return utils.InternalServerErrorResponse(c, http.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, http.MsgFetchedSuccess, stats)
}
