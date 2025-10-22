package http

import (
	"strconv"

	"keerja-backend/internal/domain/job"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/helpers"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// JobHandler aligns with the style used in user_handler:
// - Uses utils.Response helpers
// - Extracts user ID via middleware
// - Converts DTO <-> domain using mapper
// - Calls domain job.JobService directly

type JobHandler struct {
	jobService job.JobService
}

func NewJobHandler(jobService job.JobService) *JobHandler {
	return &JobHandler{jobService: jobService}
}

// GET /jobs
func (h *JobHandler) ListJobs(c *fiber.Ctx) error {
	ctx := c.Context()

	var q request.JobSearchRequest
	if err := c.QueryParser(&q); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidQueryParams)
	}
	if err := utils.ValidateStruct(&q); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, ErrValidationFailed, errs)
	}
	q.Page, q.Limit = utils.ValidatePagination(q.Page, q.Limit, 100)

	// Build domain filter using helper
	f := helpers.BuildJobFilter(q)

	jobs, total, err := h.jobService.ListJobs(ctx, f, q.Page, q.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to list jobs", err.Error())
	}

	respJobs := make([]response.JobResponse, 0, len(jobs))
	for _, j := range jobs {
		jr := mapper.ToJobResponse(&j)
		if jr != nil {
			respJobs = append(respJobs, *jr)
		}
	}

	meta := utils.GetPaginationMeta(q.Page, q.Limit, total)
	payload := response.JobListResponse{Jobs: respJobs}
	return utils.SuccessResponseWithMeta(c, MsgFetchedSuccess, payload, meta)
}

// GET /jobs/:id
func (h *JobHandler) GetJob(c *fiber.Ctx) error {
	ctx := c.Context()
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	j, err := h.jobService.GetJob(ctx, id)
	if err != nil {
		// domain layer returns not found error when missing
		return utils.NotFoundResponse(c, ErrJobNotFound)
	}

	resp := mapper.ToJobDetailResponse(j)
	return utils.SuccessResponse(c, MsgFetchedSuccess, resp)
}

// POST /jobs/search
func (h *JobHandler) SearchJobs(c *fiber.Ctx) error {
	ctx := c.Context()

	var q request.JobSearchRequest
	if err := c.BodyParser(&q); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}
	if err := utils.ValidateStruct(&q); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, ErrValidationFailed, errs)
	}
	q.Page, q.Limit = utils.ValidatePagination(q.Page, q.Limit, 100)

	// Sanitize search inputs
	q.Query = utils.SanitizeString(q.Query)
	q.Location = utils.SanitizeString(q.Location)

	// Build domain search filter using helper
	f := helpers.BuildJobSearchFilter(q)

	result, err := h.jobService.SearchJobs(ctx, f, q.Page, q.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to search jobs", err.Error())
	}

	respJobs := make([]response.JobResponse, 0, len(result.Jobs))
	for _, j := range result.Jobs {
		jr := mapper.ToJobResponse(&j)
		if jr != nil {
			respJobs = append(respJobs, *jr)
		}
	}
	meta := utils.GetPaginationMeta(result.Page, result.Limit, result.Total)
	payload := response.JobListResponse{Jobs: respJobs}
	return utils.SuccessResponseWithMeta(c, MsgFetchedSuccess, payload, meta)
}

// POST /jobs
func (h *JobHandler) CreateJob(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	var req request.CreateJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, ErrValidationFailed, errs)
	}

	// Sanitize text inputs
	req.Title = utils.SanitizeString(req.Title)
	req.Description = utils.SanitizeHTML(req.Description)
	req.RequirementsText = utils.SanitizeHTML(req.RequirementsText)
	req.Responsibilities = utils.SanitizeHTML(req.Responsibilities)

	// Security validation
	if !utils.ValidateNoXSS(req.Title) || !utils.ValidateNoXSS(req.Description) {
		return utils.BadRequestResponse(c, ErrPotentialXSS)
	}

	// Parse optional expired_at using datetime helper
	expiredAt, err := utils.ParseOptionalDateTime(req.ExpiredAt)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidDateFormat)
	}

	// Validate future date if provided
	if expiredAt != nil {
		if err := utils.MustBeFutureTime(*expiredAt); err != nil {
			return utils.BadRequestResponse(c, ErrFutureDateRequired)
		}
	}

	// Map locations
	locs := make([]job.AddLocationRequest, 0, len(req.Locations))
	for _, l := range req.Locations {
		locs = append(locs, job.AddLocationRequest{
			LocationType:  l.LocationType,
			Address:       l.Address,
			City:          l.City,
			Province:      l.Province,
			PostalCode:    l.PostalCode,
			Country:       l.Country,
			Latitude:      l.Latitude,
			Longitude:     l.Longitude,
			GooglePlaceID: l.GooglePlaceID,
			MapURL:        l.MapURL,
			IsPrimary:     l.IsPrimary,
		})
	}
	// Map skills
	skills := make([]job.AddSkillRequest, 0, len(req.Skills))
	for _, s := range req.Skills {
		skills = append(skills, job.AddSkillRequest{
			SkillID:         s.SkillID,
			ImportanceLevel: s.ImportanceLevel,
			Weight:          s.Weight,
		})
	}
	// Map benefits
	benefits := make([]job.AddBenefitRequest, 0, len(req.Benefits))
	for _, b := range req.Benefits {
		benefits = append(benefits, job.AddBenefitRequest{
			BenefitID:   b.BenefitID,
			BenefitName: b.BenefitName,
			Description: b.Description,
			IsHighlight: b.IsHighlight,
		})
	}
	// Map requirements
	reqs := make([]job.AddRequirementRequest, 0, len(req.JobRequirements))
	for _, r := range req.JobRequirements {
		reqs = append(reqs, job.AddRequirementRequest{
			RequirementType: r.RequirementType,
			RequirementText: r.RequirementText,
			SkillID:         r.SkillID,
			MinExperience:   r.MinExperience,
			MaxExperience:   r.MaxExperience,
			EducationLevel:  r.EducationLevel,
			Language:        r.Language,
			IsMandatory:     r.IsMandatory,
			Priority:        r.Priority,
		})
	}

	// Build domain request
	domainReq := &job.CreateJobRequest{
		CompanyID:        req.CompanyID,
		EmployerUserID:   employerID,
		CategoryID:       req.CategoryID,
		Title:            req.Title,
		JobLevel:         req.JobLevel,
		EmploymentType:   req.EmploymentType,
		Description:      req.Description,
		RequirementsText: req.RequirementsText,
		Responsibilities: req.Responsibilities,
		RemoteOption:     req.RemoteOption,
		SalaryMin:        req.SalaryMin,
		SalaryMax:        req.SalaryMax,
		Currency:         req.Currency,
		ExperienceMin:    req.ExperienceMin,
		ExperienceMax:    req.ExperienceMax,
		EducationLevel:   req.EducationLevel,
		TotalHires:       req.TotalHires,
		ExpiredAt:        expiredAt,
		Locations:        locs,
		Benefits:         benefits,
		Skills:           skills,
		JobRequirements:  reqs,
	}

	created, err := h.jobService.CreateJob(ctx, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create job", err.Error())
	}

	resp := mapper.ToJobDetailResponse(created)
	return utils.CreatedResponse(c, MsgCreatedSuccess, resp)
}

// PUT /jobs/:id
func (h *JobHandler) UpdateJob(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	var req request.UpdateJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, ErrValidationFailed, errs)
	}

	// CRITICAL FIX: Check ownership before allowing updates
	existingJob, err := h.jobService.GetJob(ctx, id)
	if err != nil {
		return utils.NotFoundResponse(c, ErrJobNotFound)
	}
	if existingJob.EmployerUserID == nil || *existingJob.EmployerUserID != employerID {
		return utils.ForbiddenResponse(c, ErrNotJobOwner)
	}

	// Parse optional expired_at using datetime helper
	expiredAt, err := utils.ParseOptionalDateTime(req.ExpiredAt)
	if err != nil {
		return utils.BadRequestResponse(c, ErrInvalidDateFormat)
	}

	// Validate future date if provided
	if expiredAt != nil {
		if err := utils.MustBeFutureTime(*expiredAt); err != nil {
			return utils.BadRequestResponse(c, ErrFutureDateRequired)
		}
	}

	// Build domain request (only non-nil/nonnull fields will be applied by service)
	domainReq := &job.UpdateJobRequest{}
	if req.CategoryID != nil {
		domainReq.CategoryID = req.CategoryID
	}
	if req.Title != nil {
		sanitized := utils.SanitizeString(*req.Title)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, ErrPotentialXSS)
		}
		domainReq.Title = sanitized
	}
	if req.JobLevel != nil {
		domainReq.JobLevel = *req.JobLevel
	}
	if req.EmploymentType != nil {
		domainReq.EmploymentType = *req.EmploymentType
	}
	if req.Description != nil {
		sanitized := utils.SanitizeHTML(*req.Description)
		domainReq.Description = sanitized
	}
	if req.RequirementsText != nil {
		sanitized := utils.SanitizeHTML(*req.RequirementsText)
		domainReq.RequirementsText = sanitized
	}
	if req.Responsibilities != nil {
		sanitized := utils.SanitizeHTML(*req.Responsibilities)
		domainReq.Responsibilities = sanitized
	}
	if req.RemoteOption != nil {
		domainReq.RemoteOption = req.RemoteOption
	}
	if req.SalaryMin != nil {
		domainReq.SalaryMin = req.SalaryMin
	}
	if req.SalaryMax != nil {
		domainReq.SalaryMax = req.SalaryMax
	}
	if req.Currency != nil {
		domainReq.Currency = *req.Currency
	}
	if req.ExperienceMin != nil {
		domainReq.ExperienceMin = req.ExperienceMin
	}
	if req.ExperienceMax != nil {
		domainReq.ExperienceMax = req.ExperienceMax
	}
	if req.EducationLevel != nil {
		domainReq.EducationLevel = *req.EducationLevel
	}
	if req.TotalHires != nil {
		domainReq.TotalHires = req.TotalHires
	}
	if expiredAt != nil {
		domainReq.ExpiredAt = expiredAt
	}

	updated, err := h.jobService.UpdateJob(ctx, id, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update job", err.Error())
	}

	resp := mapper.ToJobDetailResponse(updated)
	return utils.SuccessResponse(c, MsgUpdatedSuccess, resp)
}

// DELETE /jobs/:id
func (h *JobHandler) DeleteJob(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	if err := h.jobService.DeleteJob(ctx, id, employerID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete job", err.Error())
	}

	return utils.SuccessResponse(c, MsgDeletedSuccess, fiber.Map{"deleted": true})
}

// GET /jobs/my-jobs
func (h *JobHandler) GetMyJobs(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	var f request.JobFilterRequest
	if err := c.QueryParser(&f); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidQueryParams)
	}
	if err := utils.ValidateStruct(&f); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, ErrValidationFailed, errs)
	}
	f.Page, f.Limit = utils.ValidatePagination(f.Page, f.Limit, 100)

	// Build domain filter
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
	if f.IsExpired != nil {
		// translate IsExpired -> Status or PublishedAfter/IsActive if needed; here we rely on repository to interpret
	}

	jobs, total, err := h.jobService.GetMyJobs(ctx, employerID, df, f.Page, f.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get my jobs", err.Error())
	}

	respJobs := make([]response.JobResponse, 0, len(jobs))
	for _, j := range jobs {
		jr := mapper.ToJobResponse(&j)
		if jr != nil {
			respJobs = append(respJobs, *jr)
		}
	}
	meta := utils.GetPaginationMeta(f.Page, f.Limit, total)
	payload := response.JobListResponse{Jobs: respJobs}
	return utils.SuccessResponseWithMeta(c, MsgFetchedSuccess, payload, meta)
}

// POST /jobs/:id/publish
func (h *JobHandler) PublishJob(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	var req request.PublishJobRequest
	_ = c.BodyParser(&req) // optional fields
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, ErrValidationFailed, errs)
	}

	// If ExpiredAt provided, parse and validate using datetime helper
	if req.ExpiredAt != nil && *req.ExpiredAt != "" {
		expiredAt, err := utils.ParseOptionalDateTime(req.ExpiredAt)
		if err != nil {
			return utils.BadRequestResponse(c, ErrInvalidDateFormat)
		}
		if expiredAt != nil {
			if err := utils.MustBeFutureTime(*expiredAt); err != nil {
				return utils.BadRequestResponse(c, ErrFutureDateRequired)
			}
			if err := h.jobService.SetJobExpiry(ctx, id, *expiredAt); err != nil {
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to set job expiry", err.Error())
			}
		}
	}

	if err := h.jobService.PublishJob(ctx, id, employerID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to publish job", err.Error())
	}
	return utils.SuccessResponse(c, MsgUpdatedSuccess, fiber.Map{"published": true})
}

// POST /jobs/:id/close
func (h *JobHandler) CloseJob(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, ErrInvalidID)
	}

	var req request.CloseJobRequest
	_ = c.BodyParser(&req) // reason optional; domain CloseJob doesn't take it
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, ErrValidationFailed, errs)
	}

	if err := h.jobService.CloseJob(ctx, id, employerID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to close job", err.Error())
	}
	return utils.SuccessResponse(c, MsgUpdatedSuccess, fiber.Map{"closed": true})
}
