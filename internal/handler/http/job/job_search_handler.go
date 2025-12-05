package jobhandler

import (
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/helpers"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *JobHandler) SearchJobs(c *fiber.Ctx) error {
	ctx := c.Context()

	var q request.JobSearchRequest
	if err := c.BodyParser(&q); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}
	if err := utils.ValidateStruct(&q); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errs)
	}
	q.Page, q.Limit = utils.ValidatePagination(q.Page, q.Limit, 100)

	q.Query = utils.SanitizeIfNonEmpty(q.Query)
	q.Location = utils.SanitizeIfNonEmpty(q.Location)

	f := helpers.BuildJobSearchFilter(q)

	result, err := h.jobService.SearchJobs(ctx, f, q.Page, q.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to search jobs", err.Error())
	}

	// Collect unique company IDs for batch fetching
	companyIDMap := make(map[int64]bool)
	for _, j := range result.Jobs {
		companyIDMap[j.CompanyID] = true
	}

	// Fetch companies in batch to avoid N+1 queries
	companies := make(map[int64]*company.Company)
	for companyID := range companyIDMap {
		comp, err := h.companyService.GetCompany(ctx, companyID)
		if err == nil && comp != nil {
			companies[companyID] = comp
		}
	}

	// Map jobs with company data
	respJobs := make([]response.JobResponse, 0, len(result.Jobs))
	for _, j := range result.Jobs {
		comp := companies[j.CompanyID]
		jobResp := mapper.ToJobResponseWithCompany(&j, comp)
		if jobResp != nil {
			respJobs = append(respJobs, *jobResp)
		}
	}

	meta := utils.GetPaginationMeta(result.Page, result.Limit, result.Total)
	payload := response.JobListResponse{Jobs: respJobs}
	return utils.SuccessResponseWithMeta(c, common.MsgFetchedSuccess, payload, meta)
}
