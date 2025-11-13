package http

import (
	"strconv"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/job"
	"keerja-backend/internal/domain/master"
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
	jobService        job.JobService
	companyService    company.CompanyService
	jobOptionsService master.JobOptionsService
	skillsService     master.SkillsMasterService
}

// GetJobsGroupedByStatus handles GET /api/v1/jobs/grouped-by-status
// @Summary Get jobs grouped by status for mobile tab UI
// @Description Returns jobs grouped by status (active, inactive, draft, in_review)
// @Tags Jobs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/jobs/grouped-by-status [get]
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

	// Map domain jobs to response DTOs using mapper
	groupedResp := map[string][]response.JobResponse{
		"active":    {},
		"inactive":  {},
		"draft":     {},
		"in_review": {},
	}
	for status, jobs := range groupedJobs {
		for _, j := range jobs {
			jr := mapper.ToJobResponse(&j)
			if jr != nil {
				groupedResp[status] = append(groupedResp[status], *jr)
			}
		}
	}

	return utils.SuccessResponse(c, "Jobs grouped by status retrieved successfully", groupedResp)

}

func NewJobHandler(jobService job.JobService, companyService company.CompanyService, jobOptionsService master.JobOptionsService, skillService master.SkillsMasterService) *JobHandler {
	return &JobHandler{
		jobService:        jobService,
		companyService:    companyService,
		jobOptionsService: jobOptionsService,
		skillsService:     skillService,
	}
}

// GetJobTypesOptions handles GET /api/v1/jobs/job-types
// @Summary Get job types options for mobile
// @Description Get job types, work policies, company addresses, and salary ranges for job posting
// @Tags Jobs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/jobs/job-types [get]
func (h *JobHandler) GetJobTypesOptions(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get authenticated user ID from JWT
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", "User ID not found in context")
	}

	// Get job types
	jobTypes, err := h.jobOptionsService.GetJobTypes(ctx)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job types")
	}

	// Get work policies
	workPolicies, err := h.jobOptionsService.GetWorkPolicies(ctx)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve work policies")
	}

	// Get user's company addresses
	companies, err := h.companyService.GetUserCompanies(ctx, userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve company addresses")
	}

	// Extract addresses from all companies
	var addresses []response.WorkAddressOption
	for _, comp := range companies {
		// only include if company has address
		if comp.FullAddress == "" {
			continue
		}

		address := response.WorkAddressOption{
			ID:          comp.ID, // Using company ID as address ID
			CompanyID:   comp.ID,
			CompanyName: comp.CompanyName,
			FullAddress: comp.FullAddress,
			Latitude:    comp.Latitude,
			Longitude:   comp.Longitude,
		}

		// Add location details if available
		if comp.DistrictRelation != nil || comp.CityRelation != nil || comp.ProvinceRelation != nil {
			location := &response.LocationDetail{}
			if comp.DistrictRelation != nil {
				location.District = comp.DistrictRelation.Name
			}
			if comp.CityRelation != nil {
				location.City = comp.CityRelation.Name
			}
			if comp.ProvinceRelation != nil {
				location.Province = comp.ProvinceRelation.Name
			}
			address.Location = location
		}
		addresses = append(addresses, address)
	}
	// Prepare salary ranges (in rupiah)
	salaryRanges := []response.SalaryRangeOption{
		{Label: "< 1 jt", MinValue: 0, MaxValue: 1000000},
		{Label: "1 jt", MinValue: 1000000, MaxValue: 1000000},
		{Label: "2 jt", MinValue: 2000000, MaxValue: 2000000},
		{Label: "3 jt", MinValue: 3000000, MaxValue: 3000000},
		{Label: "4 jt", MinValue: 4000000, MaxValue: 4000000},
		{Label: "5 jt", MinValue: 5000000, MaxValue: 5000000},
		{Label: "6 jt", MinValue: 6000000, MaxValue: 6000000},
		{Label: "7 jt", MinValue: 7000000, MaxValue: 7000000},
		{Label: "8 jt", MinValue: 8000000, MaxValue: 8000000},
		{Label: "9 jt", MinValue: 9000000, MaxValue: 9000000},
		{Label: "10 jt", MinValue: 10000000, MaxValue: 10000000},
		{Label: "15 jt", MinValue: 15000000, MaxValue: 15000000},
		{Label: "20 jt", MinValue: 20000000, MaxValue: 20000000},
		{Label: "25 jt", MinValue: 25000000, MaxValue: 25000000},
		{Label: "30 jt", MinValue: 30000000, MaxValue: 30000000},
		{Label: "40 jt", MinValue: 40000000, MaxValue: 40000000},
		{Label: "50 jt", MinValue: 50000000, MaxValue: 50000000},
		{Label: "> 50 jt", MinValue: 50000001, MaxValue: 0},
	}
	// Salary display options
	salaryDisplayOptions := []response.SalaryDisplayOption{
		{Value: "range", Label: "Tampilkan rentang gaji", Description: "Contoh: Rp 5.000.000 - Rp 10.000.000"},
		{Value: "starting_from", Label: "Mulai dari", Description: "Contoh: Mulai dari Rp 5.000.000"},
		{Value: "up_to", Label: "Hingga", Description: "Contoh: Hingga Rp 10.000.000"},
		{Value: "hidden", Label: "Sembunyikan gaji", Description: "Gaji tidak ditampilkan di lowongan"},
	}
	resp := response.JobTypesOptionsResponse{
		JobTypes:             jobTypes,
		WorkPolicies:         workPolicies,
		WorkAddresses:        addresses,
		SalaryRanges:         salaryRanges,
		SalaryDisplayOptions: salaryDisplayOptions,
		SalaryInfo:           "Pilih rentang gaji 'Mulai Dari' dan 'Sampai' untuk menentukan range. Gunakan opsi display untuk mengatur cara tampil gaji di lowongan.",
	}
	return utils.SuccessResponse(c, "Job types options retrieved successfully", resp)
}

// GetJobRequirements handles GET /api/v1/jobs/job-requirements
// @Summary Get job requirements options for mobile
// @Description Get gender preferences, age ranges, skills, education levels, and experience levels for job requirements
// @Tags Jobs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/jobs/job-requirements [get]
func (h *JobHandler) GetJobRequirements(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get authenticated user ID from JWT
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", "User ID not found in context")
	}

	// Get gender preferences
	genders, err := h.jobOptionsService.GetGenderPreferences(ctx)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve gender preferences")
	}

	// Get education levels
	educationLevels, err := h.jobOptionsService.GetEducationLevels(ctx)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve education levels")
	}

	// Get experience levels
	experienceLevels, err := h.jobOptionsService.GetExperienceLevels(ctx)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve experience levels")
	}

	// Get recommended skills (popular skills)
	skillsFilter := &master.SkillsFilter{
		IsActive:  utils.BoolPtr(true),
		Page:      1,
		PageSize:  50, // Get top 50 popular skills
		SortBy:    "popularity_score",
		SortOrder: "DESC",
	}
	skillsResult, err := h.skillsService.GetSkills(ctx, skillsFilter)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve skills")
	}

	// Prepare age ranges
	ageRanges := []response.AgeRangeOption{
		{Label: "Tidak ada batasan umur", Min: nil, Max: nil},
		{Label: "18 tahun", Min: utils.IntPtr(18), Max: utils.IntPtr(18)},
		{Label: "19 tahun", Min: utils.IntPtr(19), Max: utils.IntPtr(19)},
		{Label: "20 tahun", Min: utils.IntPtr(20), Max: utils.IntPtr(20)},
		{Label: "21 tahun", Min: utils.IntPtr(21), Max: utils.IntPtr(21)},
		{Label: "22 tahun", Min: utils.IntPtr(22), Max: utils.IntPtr(22)},
		{Label: "23 tahun", Min: utils.IntPtr(23), Max: utils.IntPtr(23)},
		{Label: "24 tahun", Min: utils.IntPtr(24), Max: utils.IntPtr(24)},
		{Label: "25 tahun", Min: utils.IntPtr(25), Max: utils.IntPtr(25)},
		{Label: "26 tahun", Min: utils.IntPtr(26), Max: utils.IntPtr(26)},
		{Label: "27 tahun", Min: utils.IntPtr(27), Max: utils.IntPtr(27)},
		{Label: "28 tahun", Min: utils.IntPtr(28), Max: utils.IntPtr(28)},
		{Label: "29 tahun", Min: utils.IntPtr(29), Max: utils.IntPtr(29)},
		{Label: "30 tahun", Min: utils.IntPtr(30), Max: utils.IntPtr(30)},
		{Label: "35 tahun", Min: utils.IntPtr(35), Max: utils.IntPtr(35)},
		{Label: "40 tahun", Min: utils.IntPtr(40), Max: utils.IntPtr(40)},
		{Label: "45 tahun", Min: utils.IntPtr(45), Max: utils.IntPtr(45)},
		{Label: "50 tahun", Min: utils.IntPtr(50), Max: utils.IntPtr(50)},
		{Label: "55 tahun", Min: utils.IntPtr(55), Max: utils.IntPtr(55)},
		{Label: "60 tahun", Min: utils.IntPtr(60), Max: utils.IntPtr(60)},
	}

	// Map skills to DTO
	var skillOptions []response.SkillOption
	for _, skill := range skillsResult.Skills {
		skillOptions = append(skillOptions, response.SkillOption{
			ID:              skill.ID,
			Code:            skill.Code,
			Name:            skill.Name,
			NormalizedName:  skill.NormalizedName,
			SkillType:       skill.SkillType,
			DifficultyLevel: skill.DifficultyLevel,
			PopularityScore: skill.PopularityScore,
			Description:     skill.Description,
		})
	}

	resp := response.JobRequirementsOptionsResponse{
		Genders:          genders,
		AgeRanges:        ageRanges,
		EducationLevels:  educationLevels,
		ExperienceLevels: experienceLevels,
		Skills:           skillOptions,
		SkillsTotal:      skillsResult.Total,
		AgeInfo:          "Pilih rentang umur 'Min.' dan 'Max.' untuk menentukan batasan umur. Centang 'Tidak ada batasan umur' untuk tidak membatasi.",
		SkillsNote:       "User bisa pilih maksimal 17 skill. Gunakan search untuk mencari skill lain: GET /api/v1/skills/search?q={keyword}",
	}

	return utils.SuccessResponse(c, "Job requirements options retrieved successfully", resp)
}

// GET /jobs
func (h *JobHandler) ListJobs(c *fiber.Ctx) error {
	ctx := c.Context()

	var q request.JobFilterRequest
	if err := c.QueryParser(&q); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidQueryParams)
	}
	if err := utils.ValidateStruct(&q); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, ErrValidationFailed, errs)
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

	respJobs := make([]response.JobResponse, 0, len(jobs))
	for _, j := range jobs {
		jr := mapper.ToJobResponse(&j)
		if jr != nil {
			respJobs = append(respJobs, *jr)
		}
	}

	meta := utils.GetPaginationMeta(q.Page, q.Limit, total)
	payload := response.JobListResponse{Jobs: respJobs}
	// Document: /jobs?status=in_review, /jobs?status=draft, /jobs?status=published, /jobs?status=closed, /jobs?status=expired
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

	// Get company info to include in response
	company, _ := h.companyService.GetCompany(ctx, j.CompanyID)

	resp := mapper.ToJobDetailResponseWithCompany(j, company)
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
	userID := middleware.GetUserID(c)

	var req request.CreateJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, ErrValidationFailed, errs)
	}

	// Sanitize text inputs
	req.Description = utils.SanitizeHTML(req.Description)

	// Security validation
	if !utils.ValidateNoXSS(req.Description) {
		return utils.BadRequestResponse(c, ErrPotentialXSS)
	}

	// Map skills
	skills := make([]job.AddSkillRequest, 0, len(req.Skills))
	for _, s := range req.Skills {
		skills = append(skills, job.AddSkillRequest{
			SkillID:         s.SkillID,
			ImportanceLevel: s.ImportanceLevel,
		})
	}

	// Build domain request - master data only
	// Pass user ID (not employer_user ID), service will resolve it
	domainReq := &job.CreateJobRequest{
		CompanyID:          req.CompanyID,
		EmployerUserID:     userID, // This is user ID, service will look up employer_user ID
		Description:        req.Description,
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

	created, err := h.jobService.CreateJob(ctx, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create job", err.Error())
	}

	// Get company info to include in response
	company, _ := h.companyService.GetCompany(ctx, req.CompanyID)

	resp := mapper.ToJobDetailResponseWithCompany(created, company)
	return utils.CreatedResponse(c, MsgCreatedSuccess, resp)
}

// POST /jobs/draft - Save job draft (Phase 6)
// @Summary Save job draft
// @Description Save job draft with validation (salary, age ranges, XSS sanitization). Support create new or update existing by draft_id
// @Tags jobs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.SaveJobDraftRequest true "Job draft data"
// @Success 201 {object} utils.Response{data=response.JobDetailResponse}
// @Success 200 {object} utils.Response{data=response.JobDetailResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /jobs/draft [post]
func (h *JobHandler) SaveJobDraft(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get authenticated user ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated", "userID not found in context")
	}

	// Parse request body
	var req request.SaveJobDraftRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, ErrInvalidRequest)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, ErrValidationFailed, errs)
	}

	// Get company ID from user context
	// Note: We need to get the user's company. For now, we'll use a helper or assume CompanyID is available
	// In production, you should fetch this from user's employer relationship
	companyID := req.CompanyID

	// Map request DTO to domain request
	domainReq := &job.SaveJobDraftRequest{
		DraftID:          req.DraftID,
		JobTitleID:       req.JobTitleID,
		JobCategoryID:    req.JobCategoryID,
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

	// Call service to save draft
	draft, err := h.jobService.SaveJobDraft(ctx, companyID, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save job draft", err.Error())
	}

	// Get company info to include in response
	company, _ := h.companyService.GetCompany(ctx, draft.CompanyID)

	// Map to response
	resp := mapper.ToJobDetailResponseWithCompany(draft, company)

	// Return 201 Created for new draft, 200 OK for update
	if req.DraftID == nil || *req.DraftID == 0 {
		return utils.CreatedResponse(c, "Job draft created successfully", resp)
	}
	return utils.SuccessResponse(c, "Job draft updated successfully", resp)
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

	// Check ownership before allowing updates
	existingJob, err := h.jobService.GetJob(ctx, id)
	if err != nil {
		return utils.NotFoundResponse(c, ErrJobNotFound)
	}
	if existingJob.EmployerUserID == nil {
		return utils.ForbiddenResponse(c, ErrNotJobOwner)
	}

	// Build domain request - master data only
	domainReq := &job.UpdateJobRequest{
		EmployerUserID:     employerID,            // User ID for verification
		CompanyID:          existingJob.CompanyID, // Company ID for employer lookup
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
	}

	// NOTE: Status is NOT updated by users - it's controlled by workflow
	// Users can only update job details, not status
	// Status changes are handled by:
	// - PublishJob endpoint (draft → pending_approval)
	// - Admin ApproveJob endpoint (pending_approval → published)
	// - Admin RejectJob endpoint (pending_approval → rejected)
	// - System/Admin lifecycle management (published → closed/expired/suspended)

	// Sanitize description if provided
	if req.Description != nil {
		sanitized := utils.SanitizeHTML(*req.Description)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, ErrPotentialXSS)
		}
		domainReq.Description = sanitized
	}

	_, err = h.jobService.UpdateJob(ctx, id, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update job", err.Error())
	}

	// Get latest job data after update to ensure all changes are reflected in response
	latestJob, err := h.jobService.GetJob(ctx, id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve updated job", err.Error())
	}

	// Get company info to include in response
	company, _ := h.companyService.GetCompany(ctx, latestJob.CompanyID)

	resp := mapper.ToJobDetailResponseWithCompany(latestJob, company)
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
	userID := middleware.GetUserID(c)

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

	// Resolve employer_user_id from user_id and company_id
	var employerID int64
	if f.CompanyID != nil {
		resolvedID, err := h.companyService.GetEmployerUserID(ctx, userID, *f.CompanyID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Failed to resolve employer user ID", err.Error())
		}
		employerID = resolvedID
	} else {
		// If no company_id filter, get all companies for user and use the first
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

	respJobs := make([]response.JobDetailResponse, 0, len(jobs))
	for _, j := range jobs {
		jr := mapper.ToJobDetailResponse(&j)
		if jr != nil {
			// Fetch company info for each job to fill company_name
			company, err := h.companyService.GetCompany(ctx, j.CompanyID)
			if err == nil && company != nil {
				jr.CompanyName = company.CompanyName
				jr.CompanySlug = company.Slug
				jr.CompanyVerified = company.IsVerified()
			}
			respJobs = append(respJobs, *jr)
		}
	}
	meta := utils.GetPaginationMeta(f.Page, f.Limit, total)
	payload := struct {
		Jobs []response.JobDetailResponse `json:"jobs"`
	}{Jobs: respJobs}
	return utils.SuccessResponseWithMeta(c, MsgFetchedSuccess, payload, meta)
}

// POST /jobs/:id/publish - Publish job (Phase 7)
// @Summary Publish job for review
// @Description Change job status from draft to pending_review. Requires company to be verified.
// @Tags jobs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Job ID"
// @Param request body request.PublishJobRequest false "Optional publish settings"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response "User not authenticated"
// @Failure 403 {object} utils.Response "Company not verified"
// @Failure 404 {object} utils.Response "Job not found"
// @Failure 409 {object} utils.Response "Job already published or in wrong status"
// @Failure 500 {object} utils.Response
// @Router /jobs/{id}/publish [post]
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
		// Phase 7: Handle specific error cases
		errMsg := err.Error()

		// 403 Forbidden: Company not verified
		if errMsg == "company is not verified yet" {
			return utils.ForbiddenResponse(c, "Company must be verified before publishing jobs")
		}

		// 404 Not Found: Job not found or not owned by employer
		if errMsg == "job not found: record not found" || errMsg == "not authorized to manage this job" {
			return utils.NotFoundResponse(c, "Job not found or you don't have permission")
		}

		// 409 Conflict: Job already published or in wrong status
		if errMsg == "job is already pending review" || errMsg == "job is already published" {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Job status conflict", errMsg)
		}

		// 500 Internal Server Error: Other errors
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to publish job", err.Error())
	}

	// Phase 7: Return 200 OK with status pending_review
	return utils.SuccessResponse(c, "Job submitted for review successfully", fiber.Map{
		"status":  "pending_review",
		"message": "Your job has been submitted and is waiting for admin approval",
	})
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
