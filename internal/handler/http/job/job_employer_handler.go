package jobhandler

import (
	"keerja-backend/internal/domain/job"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *JobHandler) GetJobsGroupedByStatus(c *fiber.Ctx) error {
	ctx := c.UserContext()
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", "User ID not found in context")
	}

	groupedJobs, err := h.companyService.GetJobsGroupedByStatus(ctx, userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve jobs grouped by status")
	}

	groupedResp := map[string][]response.JobResponse{
		"active":    {},
		"inactive":  {},
		"draft":     {},
		"in_review": {},
	}
	for status, jobs := range groupedJobs {
		groupedResp[status] = mapper.MapEntities[job.Job, response.JobResponse](jobs, func(j *job.Job) *response.JobResponse {
			return mapper.ToJobResponse(j)
		})
	}

	return utils.SuccessResponse(c, "Jobs grouped by status retrieved successfully", groupedResp)
}

func (h *JobHandler) GetMyJobs(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var f request.JobFilterRequest
	if err := c.QueryParser(&f); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidQueryParams)
	}
	if err := utils.ValidateStruct(&f); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errs)
	}
	f.Page, f.Limit = utils.ValidatePagination(f.Page, f.Limit, 100)

	df := job.JobFilter{
		Status:   f.Status,
		City:     "",
		Province: "",
		SortBy:   f.SortBy,
	}
	if f.CompanyID != nil {
		df.CompanyID = *f.CompanyID
	}
	if f.CategoryID != nil {
		df.CategoryID = *f.CategoryID
	}

	var employerID int64
	if f.CompanyID != nil {
		resolvedID, err := h.companyService.GetEmployerUserID(ctx, userID, *f.CompanyID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Failed to resolve employer user ID", err.Error())
		}
		employerID = resolvedID
	} else {
		companies, err := h.companyService.GetUserCompanies(ctx, userID)
		if err != nil || len(companies) == 0 {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "No companies found for user", "No companies found")
		}
		resolvedID, err := h.companyService.GetEmployerUserID(ctx, userID, companies[0].ID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Failed to resolve employer user ID", err.Error())
		}
		employerID = resolvedID
	}

	jobs, total, err := h.jobService.GetMyJobs(ctx, employerID, df, f.Page, f.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get my jobs", err.Error())
	}

	respJobs := mapper.MapEntities[job.Job, response.JobDetailResponse](jobs, func(j *job.Job) *response.JobDetailResponse {
		return mapper.ToJobDetailResponse(j)
	})
	for i := range respJobs {
		comp, err := h.companyService.GetCompany(ctx, jobs[i].CompanyID)
		if err == nil && comp != nil {
			respJobs[i].CompanyName = comp.CompanyName
			respJobs[i].CompanySlug = comp.Slug
			respJobs[i].CompanyVerified = comp.IsVerified()
		}
	}

	meta := utils.GetPaginationMeta(f.Page, f.Limit, total)
	payload := struct {
		Jobs []response.JobDetailResponse `json:"jobs"`
	}{Jobs: respJobs}
	return utils.SuccessResponseWithMeta(c, common.MsgFetchedSuccess, payload, meta)
}

func (h *JobHandler) SaveJobDraft(c *fiber.Ctx) error {
	ctx := c.Context()

	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated", "userID not found in context")
	}

	var req request.SaveJobDraftRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errs)
	}

	companyID := req.CompanyID

	domainReq := &job.SaveJobDraftRequest{
		DraftID:          req.DraftID,
		JobTitleID:       req.JobTitleID,
		JobCategoryID:    req.JobCategoryID,
		JobSubcategoryID: req.JobSubcategoryID,
		JobTypeID:        req.JobTypeID,
		WorkPolicyID:     req.WorkPolicyID,
		GajiMin:          req.GajiMin,
		GajiMaks:         req.GajiMaks,
		AdaBonus:         req.AdaBonus,
		GenderPreference: req.GenderPreference,
		UmurMin:          req.UmurMin,
		UmurMaks:         req.UmurMaks,
		SkillIDs:         req.SkillIDs,
		PendidikanID:     req.PendidikanID,
		PengalamanID:     req.PengalamanID,
		Deskripsi:        req.Deskripsi,
	}

	draft, err := h.jobService.SaveJobDraft(ctx, companyID, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save job draft", err.Error())
	}

	comp, _ := h.companyService.GetCompany(ctx, draft.CompanyID)
	resp := mapper.ToJobDetailResponseWithCompany(draft, comp, nil)

	if req.DraftID == nil || *req.DraftID == 0 {
		return utils.CreatedResponse(c, "Job draft created successfully", resp)
	}
	return utils.SuccessResponse(c, "Job draft updated successfully", resp)
}
