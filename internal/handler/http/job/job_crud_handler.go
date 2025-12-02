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

func (h *JobHandler) ListJobs(c *fiber.Ctx) error {
	ctx := c.Context()

	var q request.JobFilterRequest
	if err := c.QueryParser(&q); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidQueryParams)
	}
	if err := utils.ValidateStruct(&q); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errs)
	}
	q.Page, q.Limit = utils.ValidatePagination(q.Page, q.Limit, 100)

	// Build domain filter
	f := job.JobFilter{
		Status:     q.Status,
		CompanyID:  0,
		CategoryID: 0,
		City:       "",
		Province:   "",
		SortBy:     q.SortBy,
	}
	if q.CompanyID != nil {
		f.CompanyID = *q.CompanyID
	}
	if q.CategoryID != nil {
		f.CategoryID = *q.CategoryID
	}

	jobs, total, err := h.jobService.ListJobs(ctx, f, q.Page, q.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to list jobs", err.Error())
	}

	respJobs := mapper.MapEntities[job.Job, response.JobResponse](jobs, func(j *job.Job) *response.JobResponse {
		return mapper.ToJobResponse(j)
	})

	meta := utils.GetPaginationMeta(q.Page, q.Limit, total)
	payload := response.JobListResponse{Jobs: respJobs}
	return utils.SuccessResponseWithMeta(c, common.MsgFetchedSuccess, payload, meta)
}

func (h *JobHandler) GetJob(c *fiber.Ctx) error {
	ctx := c.Context()
	id, err := utils.ParseIDParam(c, "id")
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	j, err := h.jobService.GetJob(ctx, id)
	if err != nil {
		return utils.NotFoundResponse(c, common.ErrJobNotFound)
	}

	comp, _ := h.companyService.GetCompany(ctx, j.CompanyID)
	resp := mapper.ToJobDetailResponseWithCompany(j, comp, nil)
	return utils.SuccessResponse(c, common.MsgFetchedSuccess, resp)
}

func (h *JobHandler) CreateJob(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.CreateJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errs)
	}

	req.Description = utils.SanitizeHTML(req.Description)

	if !utils.ValidateNoXSS(req.Description) {
		return utils.BadRequestResponse(c, common.ErrPotentialXSS)
	}

	skills := make([]job.AddSkillRequest, 0, len(req.Skills))
	for _, s := range req.Skills {
		skills = append(skills, job.AddSkillRequest{
			SkillID:         s.SkillID,
			ImportanceLevel: s.ImportanceLevel,
		})
	}

	domainReq := &job.CreateJobRequest{
		CompanyID:          req.CompanyID,
		EmployerUserID:     userID,
		Description:        req.Description,
		JobTitleID:         req.JobTitleID,
		JobCategoryID:      req.JobCategoryID,
		JobSubcategoryID:   req.JobSubcategoryID,
		JobTypeID:          req.JobTypeID,
		WorkPolicyID:       req.WorkPolicyID,
		EducationLevelID:   req.EducationLevelID,
		ExperienceLevelID:  req.ExperienceLevelID,
		GenderPreferenceID: req.GenderPreferenceID,
		SalaryMin:          req.SalaryMin,
		SalaryMax:          req.SalaryMax,
		SalaryDisplay:      req.SalaryDisplay,
		MinAge:             req.MinAge,
		MaxAge:             req.MaxAge,
		CompanyAddressID:   req.CompanyAddressID,
		Skills:             skills,
	}

	created, err := h.jobService.CreateJob(ctx, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create job", err.Error())
	}

	comp, _ := h.companyService.GetCompany(ctx, req.CompanyID)
	resp := mapper.ToJobDetailResponseWithCompany(created, comp, nil)
	return utils.CreatedResponse(c, common.MsgCreatedSuccess, resp)
}

func (h *JobHandler) UpdateJob(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	id, err := utils.ParseIDParam(c, "id")
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	var req request.UpdateJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errs)
	}

	existingJob, err := h.jobService.GetJob(ctx, id)
	if err != nil {
		return utils.NotFoundResponse(c, common.ErrJobNotFound)
	}
	if existingJob.EmployerUserID == nil {
		return utils.ForbiddenResponse(c, common.ErrNotJobOwner)
	}

	skills := make([]job.AddSkillRequest, 0, len(req.Skills))
	for _, s := range req.Skills {
		skills = append(skills, job.AddSkillRequest{
			SkillID:         s.SkillID,
			ImportanceLevel: s.ImportanceLevel,
		})
	}

	domainReq := &job.UpdateJobRequest{
		EmployerUserID:     employerID,
		CompanyID:          existingJob.CompanyID,
		JobTitleID:         req.JobTitleID,
		JobTypeID:          req.JobTypeID,
		WorkPolicyID:       req.WorkPolicyID,
		EducationLevelID:   req.EducationLevelID,
		ExperienceLevelID:  req.ExperienceLevelID,
		GenderPreferenceID: req.GenderPreferenceID,
		SalaryMin:          req.SalaryMin,
		SalaryMax:          req.SalaryMax,
		SalaryDisplay:      req.SalaryDisplay,
		MinAge:             req.MinAge,
		MaxAge:             req.MaxAge,
		CompanyAddressID:   req.CompanyAddressID,
		Skills:             skills,
	}

	if req.Description != nil {
		sanitized := utils.SanitizeHTML(*req.Description)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, common.ErrPotentialXSS)
		}
		domainReq.Description = sanitized
	}

	_, err = h.jobService.UpdateJob(ctx, id, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update job", err.Error())
	}

	latestJob, err := h.jobService.GetJob(ctx, id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve updated job", err.Error())
	}

	comp, _ := h.companyService.GetCompany(ctx, latestJob.CompanyID)
	resp := mapper.ToJobDetailResponseWithCompany(latestJob, comp, nil)
	return utils.SuccessResponse(c, common.MsgUpdatedSuccess, resp)
}

func (h *JobHandler) DeleteJob(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	id, err := utils.ParseIDParam(c, "id")
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	if err := h.jobService.DeleteJob(ctx, id, employerID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrInternalServer, err.Error())
	}

	return utils.SuccessResponse(c, common.MsgDeletedSuccess, fiber.Map{"deleted": true})
}
