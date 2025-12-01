package master

import (
	"strconv"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/job"
	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/dto/mapper"
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

func (h *MasterDataHandler) GetJobOptions(c *fiber.Ctx) error {
	options, err := h.jobOptionsService.GetJobOptions(c.Context())
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job options")
	}

	return utils.SuccessResponse(c, "Job options retrieved successfully", options)
}

func (h *MasterDataHandler) GetJobTypesOptions(c *fiber.Ctx) error {
	ctx := c.Context()

	jobTypes, err := h.jobOptionsService.GetJobTypes(ctx)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job types")
	}

	workPolicies, err := h.jobOptionsService.GetWorkPolicies(ctx)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve work policies")
	}

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

func (h *MasterDataHandler) GetJobPostingFormOptions(c *fiber.Ctx) error {
	ctx := c.Context()

	categories, err := h.jobService.GetCategoryTree(ctx)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job categories")
	}

	options, err := h.jobOptionsService.GetJobOptions(ctx)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job options")
	}

	salaryDefaults := map[string]interface{}{
		"min":      8000000,
		"max":      15000000,
		"display":  "range",
		"currency": "IDR",
	}

	ageDefaults := map[string]int{"min_age": 22, "max_age": 35}

	response := fiber.Map{
		"job_categories":     categories,
		"job_types":          options.JobTypes,
		"work_policies":      options.WorkPolicies,
		"education_levels":   options.EducationLevels,
		"experience_levels":  options.ExperienceLevels,
		"gender_preferences": options.GenderPreferences,
		"salary_defaults":    salaryDefaults,
		"salary_units":       []string{"Rp/bulan", "Rp/hari", "Rp/jam", "Rp/proyek"},
		"age_defaults":       ageDefaults,
	}

	userID := middleware.GetUserID(c)
	if userID != 0 {
		companies, err := h.companyService.GetUserCompanies(ctx, userID)
		if err == nil && len(companies) > 0 {
			comp := companies[0]

			companyAddresses := make([]map[string]interface{}, 0)
			addrs, err := h.companyService.GetCompanyAddresses(ctx, comp.ID, false)
			if err == nil && len(addrs) > 0 {
				for _, a := range addrs {
					addr := map[string]interface{}{
						"id":           a.ID,
						"full_address": a.FullAddress,
					}
					if a.Latitude != nil {
						addr["latitude"] = *a.Latitude
					}
					if a.Longitude != nil {
						addr["longitude"] = *a.Longitude
					}
					companyAddresses = append(companyAddresses, addr)
				}
			}
			response["company_addresses"] = companyAddresses

			if cresp := mapper.ToCompanyResponse(&comp); cresp != nil {
				response["company"] = cresp
			}
		}

		if h.skillsService != nil {
			filter := &master.SkillsFilter{Page: 1, PageSize: 100}
			skillsResp, err := h.skillsService.GetSkills(ctx, filter)
			if err == nil && skillsResp != nil {
				response["skills"] = skillsResp.Skills
			}
		}
	}

	var flatSubcategories []map[string]interface{}
	for _, cat := range categories {
		if len(cat.Subcategories) > 0 {
			for _, sc := range cat.Subcategories {
				flatSubcategories = append(flatSubcategories, map[string]interface{}{
					"id":          sc.ID,
					"category_id": sc.CategoryID,
					"code":        sc.Code,
					"name":        sc.Name,
					"is_active":   sc.IsActive,
				})
			}
		}
		if len(cat.Children) > 0 {
			for _, ch := range cat.Children {
				if len(ch.Subcategories) > 0 {
					for _, sc := range ch.Subcategories {
						flatSubcategories = append(flatSubcategories, map[string]interface{}{
							"id":          sc.ID,
							"category_id": sc.CategoryID,
							"code":        sc.Code,
							"name":        sc.Name,
							"is_active":   sc.IsActive,
						})
					}
				}
			}
		}
	}

	response["job_subcategories"] = flatSubcategories

	return utils.SuccessResponse(c, "Job posting form options retrieved successfully", response)
}

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
