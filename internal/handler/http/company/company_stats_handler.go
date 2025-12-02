package companyhandler

import (
	"strconv"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/utils"

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

func (h *CompanyStatsHandler) GetVerifiedCompanies(c *fiber.Ctx) error {
	ctx := c.Context()
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	page, limit = utils.ValidatePagination(page, limit, 100)

	companies, total, err := h.companyService.GetVerifiedCompanies(ctx, page, limit)
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

func (h *CompanyStatsHandler) GetTopRatedCompanies(c *fiber.Ctx) error {
	ctx := c.Context()
	limit := c.QueryInt("limit", 10)
	if limit <= 0 {
		limit = 10
	}

	companies, err := h.companyService.GetTopRatedCompanies(ctx, limit)
	if err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}

	respList := mapper.MapEntities[company.Company, response.CompanyResponse](companies, func(comp *company.Company) *response.CompanyResponse {
		return mapper.ToCompanyResponse(comp)
	})
	payload := response.CompanyListResponse{Companies: respList}
	return utils.SuccessResponse(c, common.MsgFetchedSuccess, payload)
}

func (h *CompanyStatsHandler) GetCompanyStats(c *fiber.Ctx) error {
	ctx := c.Context()
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	stats, err := h.companyService.GetCompanyStats(ctx, int64(companyID))
	if err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, common.MsgFetchedSuccess, stats)
}
