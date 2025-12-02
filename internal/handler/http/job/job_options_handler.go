package jobhandler

import (
	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *JobHandler) GetJobTypesOptions(c *fiber.Ctx) error {
	ctx := c.Context()

	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", "User ID not found in context")
	}

	jobTypes, err := h.jobOptionsService.GetJobTypes(ctx)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job types")
	}

	workPolicies, err := h.jobOptionsService.GetWorkPolicies(ctx)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve work policies")
	}

	companies, err := h.companyService.GetUserCompanies(ctx, userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve company addresses")
	}

	var addresses []response.WorkAddressOption
	for _, comp := range companies {
		addrs, err := h.companyService.GetCompanyAddresses(ctx, comp.ID, false)
		if err != nil {
			continue
		}
		for _, a := range addrs {
			address := response.WorkAddressOption{
				ID:          a.ID,
				CompanyID:   a.CompanyID,
				CompanyName: comp.CompanyName,
				FullAddress: a.FullAddress,
				Latitude:    a.Latitude,
				Longitude:   a.Longitude,
			}
			addresses = append(addresses, address)
		}
	}

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

func (h *JobHandler) GetJobRequirements(c *fiber.Ctx) error {
	ctx := c.Context()

	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", "User ID not found in context")
	}

	genders, err := h.jobOptionsService.GetGenderPreferences(ctx)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve gender preferences")
	}

	educationLevels, err := h.jobOptionsService.GetEducationLevels(ctx)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve education levels")
	}

	experienceLevels, err := h.jobOptionsService.GetExperienceLevels(ctx)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve experience levels")
	}

	skillsFilter := &master.SkillsFilter{
		IsActive:  utils.BoolPtr(true),
		Page:      1,
		PageSize:  50,
		SortBy:    "popularity_score",
		SortOrder: "DESC",
	}
	skillsResult, err := h.skillsService.GetSkills(ctx, skillsFilter)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to retrieve skills")
	}

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
