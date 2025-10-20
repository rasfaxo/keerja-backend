package http

import (
	"strconv"
	"strings"
	"time"

	"keerja-backend/internal/domain/job"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
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
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query params", err.Error())
	}
	if err := utils.ValidateStruct(&q); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}
	q.Page, q.Limit = utils.ValidatePagination(q.Page, q.Limit, 100)

	// Build domain filter for public listing (published jobs by default)
	f := job.JobFilter{
		Status:         "published",
		City:           q.City,
		Province:       q.Province,
		JobLevel:       q.JobLevel,
		EmploymentType: q.EmploymentType,
		EducationLevel: q.EducationLevel,
	}
	// Sorting mapping
	switch q.SortBy {
	case "posted_date":
		f.SortBy = "latest"
	case "salary":
		if strings.ToLower(q.SortOrder) == "asc" {
			f.SortBy = "salary_asc"
		} else {
			f.SortBy = "salary_desc"
		}
	case "views":
		f.SortBy = "views"
	case "applications":
		f.SortBy = "applications"
	}
	if q.RemoteOnly {
		b := true
		f.RemoteOption = &b
	}
	if q.SalaryMin != nil {
		f.MinSalary = q.SalaryMin
	}
	if q.SalaryMax != nil {
		f.MaxSalary = q.SalaryMax
	}
	if q.ExperienceMin != nil {
		f.MinExperience = q.ExperienceMin
	}
	if q.ExperienceMax != nil {
		f.MaxExperience = q.ExperienceMax
	}
	if q.CategoryID != nil {
		f.CategoryID = *q.CategoryID
	}
	if q.CompanyID != nil {
		f.CompanyID = *q.CompanyID
	}

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
	return utils.SuccessResponseWithMeta(c, "Jobs retrieved successfully", payload, meta)
}

// GET /jobs/:id
func (h *JobHandler) GetJob(c *fiber.Ctx) error {
	ctx := c.Context()
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid job ID")
	}

	j, err := h.jobService.GetJob(ctx, id)
	if err != nil {
		// domain layer returns not found error when missing
		return utils.NotFoundResponse(c, "Job not found")
	}

	resp := mapper.ToJobDetailResponse(j)
	return utils.SuccessResponse(c, "Job retrieved successfully", resp)
}

// POST /jobs/search
func (h *JobHandler) SearchJobs(c *fiber.Ctx) error {
	ctx := c.Context()

	var q request.JobSearchRequest
	if err := c.BodyParser(&q); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := utils.ValidateStruct(&q); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}
	q.Page, q.Limit = utils.ValidatePagination(q.Page, q.Limit, 100)

	// Build domain search filter
	f := job.JobSearchFilter{
		Keyword:      q.Query,
		Location:     q.Location,
		RemoteOnly:   q.RemoteOnly,
		MinSalary:    q.SalaryMin,
		MaxSalary:    q.SalaryMax,
		MinExperience: q.ExperienceMin,
		MaxExperience: q.ExperienceMax,
		PostedWithin: q.PostedWithin,
	}
	if q.CategoryID != nil {
		f.CategoryIDs = []int64{*q.CategoryID}
	}
	if len(q.SkillIDs) > 0 {
		f.SkillIDs = q.SkillIDs
	}
	if q.EmploymentType != "" {
		f.EmploymentTypes = []string{q.EmploymentType}
	}
	if q.JobLevel != "" {
		f.JobLevels = []string{q.JobLevel}
	}
	if q.EducationLevel != "" {
		f.EducationLevels = []string{q.EducationLevel}
	}
	if q.CompanyID != nil {
		f.CompanyIDs = []int64{*q.CompanyID}
	}

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
	return utils.SuccessResponseWithMeta(c, "Jobs searched successfully", payload, meta)
}

// POST /jobs
func (h *JobHandler) CreateJob(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	var req request.CreateJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	// Parse optional expired_at
	var expiredAt *time.Time
	if req.ExpiredAt != nil && *req.ExpiredAt != "" {
		if t, err := time.Parse(time.RFC3339, *req.ExpiredAt); err == nil {
			expiredAt = &t
		} else {
			return utils.BadRequestResponse(c, "expired_at must be RFC3339 format")
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
	return utils.CreatedResponse(c, "Job created successfully", resp)
}

// PUT /jobs/:id
func (h *JobHandler) UpdateJob(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid job ID")
	}

	var req request.UpdateJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	var expiredAt *time.Time
	if req.ExpiredAt != nil && *req.ExpiredAt != "" {
		if t, err := time.Parse(time.RFC3339, *req.ExpiredAt); err == nil {
			expiredAt = &t
		} else {
			return utils.BadRequestResponse(c, "expired_at must be RFC3339 format")
		}
	}

	// Build domain request (only non-nil/nonnull fields will be applied by service)
	domainReq := &job.UpdateJobRequest{}
	if req.CategoryID != nil {
		domainReq.CategoryID = req.CategoryID
	}
	if req.Title != nil {
		domainReq.Title = *req.Title
	}
	if req.JobLevel != nil {
		domainReq.JobLevel = *req.JobLevel
	}
	if req.EmploymentType != nil {
		domainReq.EmploymentType = *req.EmploymentType
	}
	if req.Description != nil {
		domainReq.Description = *req.Description
	}
	if req.RequirementsText != nil {
		domainReq.RequirementsText = *req.RequirementsText
	}
	if req.Responsibilities != nil {
		domainReq.Responsibilities = *req.Responsibilities
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

	// ownership is checked in service on specific operations; update doesn't require employerID directly here
	_ = employerID

	resp := mapper.ToJobDetailResponse(updated)
	return utils.SuccessResponse(c, "Job updated successfully", resp)
}

// DELETE /jobs/:id
func (h *JobHandler) DeleteJob(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid job ID")
	}

	if err := h.jobService.DeleteJob(ctx, id, employerID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete job", err.Error())
	}

	return utils.SuccessResponse(c, "Job deleted successfully", fiber.Map{"deleted": true})
}

// GET /jobs/my-jobs
func (h *JobHandler) GetMyJobs(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	var f request.JobFilterRequest
	if err := c.QueryParser(&f); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query params", err.Error())
	}
	if err := utils.ValidateStruct(&f); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}
	f.Page, f.Limit = utils.ValidatePagination(f.Page, f.Limit, 100)

	// Build domain filter
	df := job.JobFilter{
		Status:    f.Status,
		City:      "",
		Province:  "",
		SortBy:    f.SortBy,
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
	return utils.SuccessResponseWithMeta(c, "My jobs retrieved successfully", payload, meta)
}

// POST /jobs/:id/publish
func (h *JobHandler) PublishJob(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid job ID")
	}

	var req request.PublishJobRequest
	_ = c.BodyParser(&req) // optional fields
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	// If ExpiredAt provided, set it before publishing
	if req.ExpiredAt != nil && *req.ExpiredAt != "" {
		if t, err := time.Parse(time.RFC3339, *req.ExpiredAt); err == nil {
			if err := h.jobService.SetJobExpiry(ctx, id, t); err != nil {
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to set job expiry", err.Error())
			}
		} else {
			return utils.BadRequestResponse(c, "expired_at must be RFC3339 format")
		}
	}

	if err := h.jobService.PublishJob(ctx, id, employerID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to publish job", err.Error())
	}
	return utils.SuccessResponse(c, "Job published successfully", fiber.Map{"published": true})
}

// POST /jobs/:id/close
func (h *JobHandler) CloseJob(c *fiber.Ctx) error {
	ctx := c.Context()
	employerID := middleware.GetUserID(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return utils.BadRequestResponse(c, "Invalid job ID")
	}

	var req request.CloseJobRequest
	_ = c.BodyParser(&req) // reason optional; domain CloseJob doesn't take it
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errs)
	}

	if err := h.jobService.CloseJob(ctx, id, employerID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to close job", err.Error())
	}
	return utils.SuccessResponse(c, "Job closed successfully", fiber.Map{"closed": true})
}
