package master

import (
	"strconv"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/job"
	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type MasterDataHandler struct {
	jobTitleService   master.JobTitleService
	jobOptionsService master.JobOptionsService
	jobService        job.JobService
	companyService    company.CompanyService
	skillsService     master.SkillsMasterService
}

func NewMasterDataHandler(
	jobTitleService master.JobTitleService,
	jobOptionsService master.JobOptionsService,
	jobService job.JobService,
	companyService company.CompanyService,
	skillsService master.SkillsMasterService,
) *MasterDataHandler {
	return &MasterDataHandler{
		jobTitleService:   jobTitleService,
		jobOptionsService: jobOptionsService,
		jobService:        jobService,
		companyService:    companyService,
		skillsService:     skillsService,
	}
}

// ==================== JOB TITLES ENDPOINT ====================

func (h *MasterDataHandler) GetJobTitles(c *fiber.Ctx) error {
	query := c.Query("q", "")
	limitStr := c.Query("limit", "20")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	jobTitles, err := h.jobTitleService.SearchJobTitles(c.Context(), query, limit)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job titles")
	}

	return utils.SuccessResponse(c, "Job titles retrieved successfully", jobTitles)
}

// ==================== PAGE-BASED FORM OPTIONS ====================

// GetJobDetailsOptions returns master data for Page 1: Job Details
// Returns: job categories (for Job Name/job_category_id) and subcategories (for Job Field/job_subcategory_id)
// GET /api/v1/master/job-form/job-details?q=search
// Query params:
//   - q: search keyword for category name (optional)
func (h *MasterDataHandler) GetJobDetailsOptions(c *fiber.Ctx) error {
	ctx := c.Context()
	query := c.Query("q", "")

	var categories []job.JobCategory
	var err error

	if query != "" {
		// Search categories with keyword filter
		filter := job.CategoryFilter{
			Keyword:  query,
			IsActive: boolPtr(true),
		}
		categories, _, err = h.jobService.ListCategories(ctx, filter, 1, 100)
	} else {
		// Get all categories with nested subcategories (tree structure)
		categories, err = h.jobService.GetCategoryTree(ctx)
	}

	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job categories")
	}

	return utils.SuccessResponse(c, "Job details options retrieved successfully", fiber.Map{
		"job_categories": categories, // For "Job Name" dropdown (job_category_id)
		// Subcategories are nested inside each category for "Job Field" (job_subcategory_id)
		"_page":        1,
		"_page_name":   "job_details",
		"_searched":    query != "",
		"_description": "Select job category (Job Name) first, then select subcategory (Job Field)",
	})
}

// GetJobTypeOptions returns master data for Page 2: Job Type
// Returns: job types, work policies, company addresses, salary ranges
// GET /api/v1/master/job-form/job-type
// Requires auth for company addresses
func (h *MasterDataHandler) GetJobTypeOptions(c *fiber.Ctx) error {
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

	// Salary ranges
	salaryRanges := []fiber.Map{
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
		"job_types":     jobTypes,     // For "Job Type" dropdown (job_type_id)
		"work_policies": workPolicies, // For "Work Policy" dropdown (work_policy_id)
		"salary_ranges": salaryRanges, // For "Salary" selection
		"salary_units":  []string{"Rp/bulan", "Rp/hari", "Rp/jam", "Rp/proyek"},
		"_page":         2,
		"_page_name":    "job_type",
		"_description":  "Select job type, work policy, work address, and salary range",
	}

	// Get company addresses if user is authenticated
	userID := middleware.GetUserID(c)
	if userID > 0 {
		// Get user's companies
		companies, err := h.companyService.GetUserCompanies(ctx, userID)
		if err == nil && len(companies) > 0 {
			var allAddresses []fiber.Map
			for _, comp := range companies {
				addrs, err := h.companyService.GetCompanyAddresses(ctx, comp.ID, false)
				if err == nil {
					for _, addr := range addrs {
						allAddresses = append(allAddresses, fiber.Map{
							"id":           addr.ID,
							"company_id":   addr.CompanyID,
							"company_name": comp.CompanyName,
							"full_address": addr.FullAddress,
							"latitude":     addr.Latitude,
							"longitude":    addr.Longitude,
						})
					}
				}
			}
			response["company_addresses"] = allAddresses // For "Work Address" dropdown (company_address_id)
		}
	}

	return utils.SuccessResponse(c, "Job type options retrieved successfully", response)
}

// GetJobRequirementsOptions returns master data for Page 3: Job Requirements
// Returns: gender preferences, skills (paginated), education levels, experience levels
// GET /api/v1/master/job-form/job-requirements
func (h *MasterDataHandler) GetJobRequirementsOptions(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get all job options
	options, err := h.jobOptionsService.GetJobOptions(ctx)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job options")
	}

	// Get skills with pagination
	skillsPage := c.QueryInt("skills_page", 1)
	skillsLimit := c.QueryInt("skills_limit", 50)
	skillsQuery := c.Query("skills_q", "")

	skillsPage, skillsLimit = utils.ValidatePagination(skillsPage, skillsLimit, 100)

	filter := &master.SkillsFilter{
		Search:   skillsQuery,
		IsActive: nil,
		Page:     skillsPage,
		PageSize: skillsLimit,
	}

	skillsResp, err := h.skillsService.GetSkills(ctx, filter)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve skills")
	}

	var skills interface{} = skillsResp.Skills
	if skillsResp.Skills == nil {
		skills = []interface{}{}
	}

	// Age limit options
	ageLimits := []fiber.Map{
		{"label": "No Age Limit", "min": 0, "max": 0},
		{"label": "18-25 tahun", "min": 18, "max": 25},
		{"label": "20-30 tahun", "min": 20, "max": 30},
		{"label": "25-35 tahun", "min": 25, "max": 35},
		{"label": "30-40 tahun", "min": 30, "max": 40},
		{"label": "35-45 tahun", "min": 35, "max": 45},
		{"label": "Max 25 tahun", "min": 0, "max": 25},
		{"label": "Max 30 tahun", "min": 0, "max": 30},
		{"label": "Max 35 tahun", "min": 0, "max": 35},
		{"label": "Max 40 tahun", "min": 0, "max": 40},
	}

	skillsMeta := utils.GetPaginationMeta(skillsPage, skillsLimit, int64(skillsResp.Total))

	return utils.SuccessResponse(c, "Job requirements options retrieved successfully", fiber.Map{
		"gender_preferences": options.GenderPreferences, // For "Gender" dropdown
		"age_limits":         ageLimits,                 // For "Age" dropdown
		"skills":             skills,                    // For "Required Skill" multi-select
		"skills_meta":        skillsMeta,                // Pagination info for skills
		"education_levels":   options.EducationLevels,   // For "Minimum Education Required" dropdown
		"experience_levels":  options.ExperienceLevels,  // For "Required Work Experience" dropdown
		"_page":              3,
		"_page_name":         "job_requirements",
		"_description":       "Select gender preference, age limit, required skills, education, and experience",
	})
}

// GetJobDescriptionOptions returns master data for Page 4: Job Description
// Returns: only metadata since this page is mainly text input
// GET /api/v1/master/job-form/job-description
func (h *MasterDataHandler) GetJobDescriptionOptions(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, "Job description options retrieved successfully", fiber.Map{
		"max_length":  5000, 
		"min_length":  100,  
		"placeholder": "Describe the job responsibilities, requirements, and benefits...",
		"tips": []string{
			"Jelaskan tanggung jawab utama pekerjaan",
			"Sebutkan kualifikasi yang dibutuhkan",
			"Tambahkan benefit dan fasilitas yang ditawarkan",
			"Gunakan format bullet points untuk kemudahan membaca",
		},
		"_page":        4,
		"_page_name":   "job_description",
		"_description": "Enter job description text",
	})
}

// ==================== LEGACY/UTILITY ENDPOINTS ====================

// GetJobOptions returns all job options (legacy endpoint)
func (h *MasterDataHandler) GetJobOptions(c *fiber.Ctx) error {
	options, err := h.jobOptionsService.GetJobOptions(c.Context())
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job options")
	}

	return utils.SuccessResponse(c, "Job options retrieved successfully", options)
}

// GetSkillsPaginated returns paginated skills with optional search
// GET /api/v1/master/skills?q=search&page=1&limit=50
func (h *MasterDataHandler) GetSkillsPaginated(c *fiber.Ctx) error {
	ctx := c.Context()
	query := c.Query("q", "")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)

	page, limit = utils.ValidatePagination(page, limit, 100)

	filter := &master.SkillsFilter{
		Search:   query,
		IsActive: nil,
		Page:     page,
		PageSize: limit,
	}

	skillsResp, err := h.skillsService.GetSkills(ctx, filter)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve skills")
	}

	var skills interface{} = skillsResp.Skills
	if skillsResp.Skills == nil {
		skills = []interface{}{}
	}

	meta := utils.GetPaginationMeta(page, limit, int64(skillsResp.Total))
	return utils.SuccessResponseWithMeta(c, "Skills retrieved successfully", fiber.Map{
		"skills": skills,
	}, meta)
}

// ==================== ADMIN ENDPOINTS ====================

func (h *MasterDataHandler) CreateJobTitle(c *fiber.Ctx) error {
	var req master.CreateJobTitleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	jobTitle, err := h.jobTitleService.CreateJobTitle(c.Context(), &req)
	if err != nil {
		if err.Error() == "job title with this name already exists" {
			return utils.ConflictResponse(c, err.Error())
		}
		return utils.InternalServerErrorResponse(c, "Failed to create job title")
	}

	return utils.CreatedResponse(c, "Job title created successfully", jobTitle)
}

func (h *MasterDataHandler) UpdateJobTitle(c *fiber.Ctx) error {
	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid job title ID")
	}

	var req master.UpdateJobTitleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

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

func (h *MasterDataHandler) DeleteJobTitle(c *fiber.Ctx) error {
	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid job title ID")
	}

	if err := h.jobTitleService.DeleteJobTitle(c.Context(), id); err != nil {
		if err.Error() == "job title not found" {
			return utils.NotFoundResponse(c, err.Error())
		}
		return utils.InternalServerErrorResponse(c, "Failed to delete job title")
	}

	return utils.SuccessResponse(c, "Job title deleted successfully", nil)
}

func (h *MasterDataHandler) GetJobTitleByID(c *fiber.Ctx) error {
	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid job title ID")
	}

	jobTitle, err := h.jobTitleService.GetJobTitle(c.Context(), id)
	if err != nil {
		if err.Error() == "job title not found" {
			return utils.NotFoundResponse(c, err.Error())
		}
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job title")
	}

	return utils.SuccessResponse(c, "Job title retrieved successfully", jobTitle)
}

// ==================== HELPER FUNCTIONS ====================

// boolPtr returns a pointer to a bool value
func boolPtr(b bool) *bool {
	return &b
}
