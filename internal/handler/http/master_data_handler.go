package http

import (
	"strconv"

	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// MasterDataHandler handles master data HTTP requests
type MasterDataHandler struct {
	jobTitleService   master.JobTitleService
	jobOptionsService master.JobOptionsService
}

// NewMasterDataHandler creates a new master data handler
func NewMasterDataHandler(
	jobTitleService master.JobTitleService,
	jobOptionsService master.JobOptionsService,
) *MasterDataHandler {
	return &MasterDataHandler{
		jobTitleService:   jobTitleService,
		jobOptionsService: jobOptionsService,
	}
}

// GetJobTitles handles GET /api/v1/master/job-titles
// @Summary Get job titles with smart search
// @Description Search job titles with fuzzy matching and category recommendations
// @Tags Master Data
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param limit query int false "Results limit (default: 20, max: 100)"
// @Success 200 {object} utils.Response{data=[]master.JobTitleResponse}
// @Failure 400 {object} utils.Response
// @Failure 429 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/master/job-titles [get]
func (h *MasterDataHandler) GetJobTitles(c *fiber.Ctx) error {
	// Parse query parameters
	query := c.Query("q", "")
	limitStr := c.Query("limit", "20")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	// Get job titles
	jobTitles, err := h.jobTitleService.SearchJobTitles(c.Context(), query, limit)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job titles")
	}

	return utils.SuccessResponse(c, "Job titles retrieved successfully", jobTitles)
}

// GetJobOptions handles GET /api/v1/master/job-options
// @Summary Get all job posting options (DEPRECATED - use /job-types instead)
// @Description Get job types, work policies, education levels, experience levels, and gender preferences
// @Tags Master Data
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=master.JobOptionsResponse}
// @Failure 429 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/master/job-options [get]
// // @Deprecated
func (h *MasterDataHandler) GetJobOptions(c *fiber.Ctx) error {
	// Get all job options (heavily cached)
	options, err := h.jobOptionsService.GetJobOptions(c.Context())
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job options")
	}

	return utils.SuccessResponse(c, "Job options retrieved successfully", options)
}

// GetJobTypes handles GET /api/v1/master/job-types-options
// @Summary Get job types and work policies for mobile
// @Description Get job types, work policies, company addresses, and salary ranges for job posting
// @Tags Master Data
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/master/job-types-options [get]
func (h *MasterDataHandler) GetJobTypesOptions(c *fiber.Ctx) error {
	ctx := c.Context()

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

	// Prepare salary ranges (in millions IDR)
	salaryRanges := []map[string]interface{}{
		{"label": "< 1 jt", "min": 0, "max": 1000000},
		{"label": "1 jt", "min": 1000000, "max": 1000000},
		{"label": "2 jt", "min": 2000000, "max": 2000000},
		{"label": "3 jt", "min": 3000000, "max": 3000000},
		{"label": "4 jt", "min": 4000000, "max": 4000000},
		{"label": "5 jt", "min": 5000000, "max": 5000000},
		{"label": "6 jt", "min": 6000000, "max": 6000000},
		{"label": "7 jt", "min": 7000000, "max": 7000000},
		{"label": "8 jt", "min": 8000000, "max": 8000000},
		{"label": "9 jt", "min": 9000000, "max": 9000000},
		{"label": "10 jt", "min": 10000000, "max": 10000000},
		{"label": "15 jt", "min": 15000000, "max": 15000000},
		{"label": "20 jt", "min": 20000000, "max": 20000000},
		{"label": "25 jt", "min": 25000000, "max": 25000000},
		{"label": "30 jt", "min": 30000000, "max": 30000000},
		{"label": "40 jt", "min": 40000000, "max": 40000000},
		{"label": "50 jt", "min": 50000000, "max": 50000000},
		{"label": "> 50 jt", "min": 50000000, "max": 0},
	}

	response := fiber.Map{
		"job_types":     jobTypes,
		"work_policies": workPolicies,
		"salary_ranges": salaryRanges,
		"salary_units":  []string{"Rp/bulan", "Rp/hari", "Rp/jam", "Rp/proyek"},
		"salary_info":   "Pilih rentang gaji dari 'Mulai Dari' hingga 'Sampai'. Kosongkan jika tidak ingin menampilkan gaji.",
	}

	return utils.SuccessResponse(c, "Job types options retrieved successfully", response)
}

// Admin-only endpoints for managing job titles

// CreateJobTitle handles POST /api/v1/admin/master/job-titles
// @Summary Create a new job title (admin only)
// @Description Create a new job title with category recommendation
// @Tags Admin - Master Data
// @Accept json
// @Produce json
// @Param request body master.CreateJobTitleRequest true "Job title details"
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=master.JobTitle}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/admin/master/job-titles [post]
func (h *MasterDataHandler) CreateJobTitle(c *fiber.Ctx) error {
	var req master.CreateJobTitleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Create job title
	jobTitle, err := h.jobTitleService.CreateJobTitle(c.Context(), &req)
	if err != nil {
		if err.Error() == "job title with this name already exists" {
			return utils.ConflictResponse(c, err.Error())
		}
		return utils.InternalServerErrorResponse(c, "Failed to create job title")
	}

	return utils.CreatedResponse(c, "Job title created successfully", jobTitle)
}

// UpdateJobTitle handles PUT /api/v1/admin/master/job-titles/:id
// @Summary Update a job title (admin only)
// @Description Update an existing job title
// @Tags Admin - Master Data
// @Accept json
// @Produce json
// @Param id path int true "Job Title ID"
// @Param request body master.UpdateJobTitleRequest true "Updated job title details"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=master.JobTitle}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/admin/master/job-titles/{id} [put]
func (h *MasterDataHandler) UpdateJobTitle(c *fiber.Ctx) error {
	// Parse ID
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid job title ID")
	}

	var req master.UpdateJobTitleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Update job title
	jobTitle, err := h.jobTitleService.UpdateJobTitle(c.Context(), id, &req)
	if err != nil {
		if err.Error() == "job title not found" {
			return utils.NotFoundResponse(c, err.Error())
		}
		if err.Error() == "job title with this name already exists" {
			return utils.ConflictResponse(c, err.Error())
		}
		return utils.InternalServerErrorResponse(c, "Failed to update job title")
	}

	return utils.SuccessResponse(c, "Job title updated successfully", jobTitle)
}

// DeleteJobTitle handles DELETE /api/v1/admin/master/job-titles/:id
// @Summary Delete a job title (admin only)
// @Description Delete a job title
// @Tags Admin - Master Data
// @Accept json
// @Produce json
// @Param id path int true "Job Title ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/admin/master/job-titles/{id} [delete]
func (h *MasterDataHandler) DeleteJobTitle(c *fiber.Ctx) error {
	// Parse ID
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid job title ID")
	}

	// Delete job title
	if err := h.jobTitleService.DeleteJobTitle(c.Context(), id); err != nil {
		if err.Error() == "job title not found" {
			return utils.NotFoundResponse(c, err.Error())
		}
		return utils.InternalServerErrorResponse(c, "Failed to delete job title")
	}

	return utils.SuccessResponse(c, "Job title deleted successfully", nil)
}

// GetJobTitleByID handles GET /api/v1/admin/master/job-titles/:id
// @Summary Get a job title by ID (admin only)
// @Description Retrieve a single job title by ID
// @Tags Admin - Master Data
// @Accept json
// @Produce json
// @Param id path int true "Job Title ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=master.JobTitleResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/admin/master/job-titles/{id} [get]
func (h *MasterDataHandler) GetJobTitleByID(c *fiber.Ctx) error {
	// Parse ID
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid job title ID")
	}

	// Get job title
	jobTitle, err := h.jobTitleService.GetJobTitle(c.Context(), id)
	if err != nil {
		if err.Error() == "job title not found" {
			return utils.NotFoundResponse(c, err.Error())
		}
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job title")
	}

	return utils.SuccessResponse(c, "Job title retrieved successfully", jobTitle)
}
