package mapper

import (
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/job"
	"keerja-backend/internal/dto/response"
)

// Job Entity to Response Mappers

// ToJobResponse maps Job entity to JobResponse DTO (master data only)
func ToJobResponse(j *job.Job) *response.JobResponse {
	if j == nil {
		return nil
	}

	var daysRemaining *int
	if j.ExpiredAt != nil && j.IsActive() {
		days := int(j.ExpiredAt.Sub(j.CreatedAt).Hours() / 24)
		daysRemaining = &days
	}

	resp := &response.JobResponse{
		ID:                j.ID,
		UUID:              j.UUID.String(),
		CompanyID:         j.CompanyID,
		Title:             j.Title,
		Slug:              j.Slug,
		SalaryMin:         j.SalaryMin,
		SalaryMax:         j.SalaryMax,
		SalaryDisplay:     j.SalaryDisplay,
		MinAge:            j.MinAge,
		MaxAge:            j.MaxAge,
		Currency:          j.Currency,
		Status:            j.Status,
		ViewsCount:        j.ViewsCount,
		ApplicationsCount: j.ApplicationsCount,
		PublishedAt:       j.PublishedAt,
		ExpiredAt:         j.ExpiredAt,
		CreatedAt:         j.CreatedAt,
		IsExpired:         j.IsExpired(),
		DaysRemaining:     daysRemaining,
	}

	// Map master data details
	if j.JobType != nil {
		resp.JobType = &response.JobMasterDataItem{
			ID:          j.JobType.ID,
			Code:        j.JobType.Code,
			Name:        j.JobType.Name,
			Description: "",
		}
	}

	if j.WorkPolicy != nil {
		resp.WorkPolicy = &response.JobMasterDataItem{
			ID:          j.WorkPolicy.ID,
			Code:        j.WorkPolicy.Code,
			Name:        j.WorkPolicy.Name,
			Description: "",
		}
	}

	// Include additional master-data if available
	if j.JobTitle != nil {
		resp.JobTitle = &response.JobMasterDataItem{
			ID:   j.JobTitle.ID,
			Code: j.JobTitle.NormalizedName,
			Name: j.JobTitle.Name,
		}
	}

	if j.EducationLevelM != nil {
		resp.EducationLevel = &response.JobMasterDataItem{
			ID:   j.EducationLevelM.ID,
			Code: j.EducationLevelM.Code,
			Name: j.EducationLevelM.Name,
		}
	}

	if j.ExperienceLevelM != nil {
		resp.ExperienceLevel = &response.JobMasterDataItem{
			ID:   j.ExperienceLevelM.ID,
			Code: j.ExperienceLevelM.Code,
			Name: j.ExperienceLevelM.Name,
		}
	}

	if j.GenderPreference != nil {
		resp.GenderPreference = &response.JobMasterDataItem{
			ID:   j.GenderPreference.ID,
			Code: j.GenderPreference.Code,
			Name: j.GenderPreference.Name,
		}
	}

	// Map category/subcategory objects
	if j.Category != nil {
		resp.JobCategory = &response.JobCategoryResponse{
			ID:           j.Category.ID,
			CategoryName: j.Category.Name,
			Description:  j.Category.Description,
		}
	}
	if j.JobSubcategory != nil {
		resp.JobSubcategory = &response.JobSubcategoryResponse{
			ID:              j.JobSubcategory.ID,
			SubcategoryName: j.JobSubcategory.Name,
			Description:     j.JobSubcategory.Description,
		}
	}

	return resp
}

// ToJobResponseWithCompany maps Job entity with company info to JobResponse DTO
func ToJobResponseWithCompany(j *job.Job, comp *company.Company) *response.JobResponse {
	if j == nil {
		return nil
	}

	resp := ToJobResponse(j)
	if resp == nil {
		return nil
	}

	// Add company info if available
	if comp != nil {
		resp.CompanyName = comp.CompanyName
		if comp.LogoURL != nil {
			resp.CompanyLogoURL = *comp.LogoURL
		}
		resp.CompanyVerified = comp.IsVerified()
	}

	return resp
}

// ToJobDetailResponse maps Job entity with relations to JobDetailResponse DTO (master data only)
func ToJobDetailResponse(j *job.Job) *response.JobDetailResponse {
	if j == nil {
		return nil
	}

	var daysRemaining *int
	if j.ExpiredAt != nil && j.IsActive() {
		days := int(j.ExpiredAt.Sub(j.CreatedAt).Hours() / 24)
		daysRemaining = &days
	}

	resp := &response.JobDetailResponse{
		ID:             j.ID,
		UUID:           j.UUID.String(),
		CompanyID:      j.CompanyID,
		EmployerUserID: j.EmployerUserID,
		Title:          j.Title,
		Slug:           j.Slug,

		Description: j.Description,

		SalaryMin:         j.SalaryMin,
		SalaryMax:         j.SalaryMax,
		SalaryDisplay:     j.SalaryDisplay,
		MinAge:            j.MinAge,
		MaxAge:            j.MaxAge,
		Currency:          j.Currency,
		Status:            j.Status,
		ViewsCount:        j.ViewsCount,
		ApplicationsCount: j.ApplicationsCount,
		PublishedAt:       j.PublishedAt,
		ExpiredAt:         j.ExpiredAt,
		CreatedAt:         j.CreatedAt,
		UpdatedAt:         j.UpdatedAt,
		IsExpired:         j.IsExpired(),
		DaysRemaining:     daysRemaining,
	}

	// Category/Subcategory objects are populated below (no numeric IDs in response)

	// Map Job Master Data Details
	if j.HasJobMasterDataRelations() {
		// Job Title
		if j.JobTitle != nil {
			resp.JobTitle = &response.JobMasterDataItem{
				ID:          j.JobTitle.ID,
				Code:        j.JobTitle.NormalizedName,
				Name:        j.JobTitle.Name,
				Description: "",
			}
		}

		// Job Type
		if j.JobType != nil {
			resp.JobType = &response.JobMasterDataItem{
				ID:          j.JobType.ID,
				Code:        j.JobType.Code,
				Name:        j.JobType.Name,
				Description: "",
			}
		}

		// Work Policy
		if j.WorkPolicy != nil {
			resp.WorkPolicy = &response.JobMasterDataItem{
				ID:          j.WorkPolicy.ID,
				Code:        j.WorkPolicy.Code,
				Name:        j.WorkPolicy.Name,
				Description: "",
			}
		}

		// Education Level
		if j.EducationLevelM != nil {
			resp.EducationLevel = &response.JobMasterDataItem{
				ID:          j.EducationLevelM.ID,
				Code:        j.EducationLevelM.Code,
				Name:        j.EducationLevelM.Name,
				Description: "",
			}
		}

		// Experience Level
		if j.ExperienceLevelM != nil {
			resp.ExperienceLevel = &response.JobMasterDataItem{
				ID:          j.ExperienceLevelM.ID,
				Code:        j.ExperienceLevelM.Code,
				Name:        j.ExperienceLevelM.Name,
				Description: "",
			}
		}

		// Gender Preference
		if j.GenderPreference != nil {
			resp.GenderPreference = &response.JobMasterDataItem{
				ID:          j.GenderPreference.ID,
				Code:        j.GenderPreference.Code,
				Name:        j.GenderPreference.Name,
				Description: "",
			}
		}

		// Map category/subcategory objects
		if j.Category != nil {
			resp.JobCategory = &response.JobCategoryResponse{
				ID:           j.Category.ID,
				CategoryName: j.Category.Name,
				Description:  j.Category.Description,
			}
		}
		if j.JobSubcategory != nil {
			resp.JobSubcategory = &response.JobSubcategoryResponse{
				ID:              j.JobSubcategory.ID,
				SubcategoryName: j.JobSubcategory.Name,
				Description:     j.JobSubcategory.Description,
			}
		}
	}

	// Map skills
	if len(j.Skills) > 0 {
		resp.Skills = make([]response.JobSkillResponse, len(j.Skills))
		for i, skill := range j.Skills {
			resp.Skills[i] = *ToJobSkillResponse(&skill)
		}
	}

	// Map benefits
	if len(j.Benefits) > 0 {
		resp.Benefits = make([]response.JobBenefitResponse, len(j.Benefits))
		for i, benefit := range j.Benefits {
			resp.Benefits[i] = *ToJobBenefitResponse(&benefit)
		}
	}

	// Map locations
	if len(j.Locations) > 0 {
		resp.Locations = make([]response.JobLocationResponse, len(j.Locations))
		for i, location := range j.Locations {
			resp.Locations[i] = *ToJobLocationResponse(&location)
		}
	}

	// Map requirements
	if len(j.JobRequirements) > 0 {
		resp.JobRequirements = make([]response.JobRequirementResponse, len(j.JobRequirements))
		for i, req := range j.JobRequirements {
			resp.JobRequirements[i] = *ToJobRequirementResponse(&req)
		}
	}

	// If job has a preloaded CompanyAddress relation (job-local type), map it
	if j.CompanyAddress != nil {
		resp.CompanyAddress = &response.CompanyAddressResponse{
			ID:          j.CompanyAddress.ID,
			FullAddress: j.CompanyAddress.FullAddress,
		}
		if j.CompanyAddress.Latitude != nil {
			resp.CompanyAddress.Latitude = *j.CompanyAddress.Latitude
		}
		if j.CompanyAddress.Longitude != nil {
			resp.CompanyAddress.Longitude = *j.CompanyAddress.Longitude
		}
		resp.CompanyAddress.ProvinceID = j.CompanyAddress.ProvinceID
		resp.CompanyAddress.CityID = j.CompanyAddress.CityID
		resp.CompanyAddress.DistrictID = j.CompanyAddress.DistrictID
	}

	return resp
}

// derefInt64 safely dereferences an int64 pointer, returning 0 if nil
func derefInt64(ptr *int64) int64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// ToJobDetailResponseWithCompany maps Job entity with company info to JobDetailResponse DTO
func ToJobDetailResponseWithCompany(j *job.Job, comp *company.Company, addr *company.CompanyAddress) *response.JobDetailResponse {
	if j == nil {
		return nil
	}

	resp := ToJobDetailResponse(j)
	if resp == nil {
		return nil
	}

	// Add company info if available
	if comp != nil {
		resp.CompanyName = comp.CompanyName
		if comp.LogoURL != nil {
			resp.CompanyLogoURL = *comp.LogoURL
		}
		resp.CompanyVerified = comp.IsVerified()
		resp.CompanySlug = comp.Slug
	}

	// Add selected company address if provided
	if addr != nil {
		resp.CompanyAddress = &response.CompanyAddressResponse{
			ID:          addr.ID,
			FullAddress: addr.FullAddress,
		}
		if addr.Latitude != nil {
			resp.CompanyAddress.Latitude = *addr.Latitude
		}
		if addr.Longitude != nil {
			resp.CompanyAddress.Longitude = *addr.Longitude
		}
		resp.CompanyAddress.ProvinceID = addr.ProvinceID
		resp.CompanyAddress.CityID = addr.CityID
		resp.CompanyAddress.DistrictID = addr.DistrictID
	}

	return resp
}

// ToJobSkillResponse maps JobSkill entity to JobSkillResponse DTO
func ToJobSkillResponse(s *job.JobSkill) *response.JobSkillResponse {
	if s == nil {
		return nil
	}

	skillName := ""
	if s.Skill != nil {
		skillName = s.Skill.Name
	}

	return &response.JobSkillResponse{
		ID:              s.ID,
		SkillID:         s.SkillID,
		SkillName:       skillName,
		ImportanceLevel: s.ImportanceLevel,
		Weight:          s.Weight,
	}
}

// ToJobBenefitResponse maps JobBenefit entity to JobBenefitResponse DTO
func ToJobBenefitResponse(b *job.JobBenefit) *response.JobBenefitResponse {
	if b == nil {
		return nil
	}

	return &response.JobBenefitResponse{
		ID:          b.ID,
		BenefitID:   b.BenefitID,
		BenefitName: b.BenefitName,
		Description: b.Description,
		IsHighlight: b.IsHighlight,
	}
}

// ToJobLocationResponse maps JobLocation entity to JobLocationResponse DTO
func ToJobLocationResponse(l *job.JobLocation) *response.JobLocationResponse {
	if l == nil {
		return nil
	}

	return &response.JobLocationResponse{
		ID:            l.ID,
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
	}
}

// ToJobRequirementResponse maps JobRequirement entity to JobRequirementResponse DTO
func ToJobRequirementResponse(r *job.JobRequirement) *response.JobRequirementResponse {
	if r == nil {
		return nil
	}

	return &response.JobRequirementResponse{
		ID:              r.ID,
		RequirementType: r.RequirementType,
		RequirementText: r.RequirementText,
		SkillID:         r.SkillID,
		MinExperience:   r.MinExperience,
		MaxExperience:   r.MaxExperience,
		EducationLevel:  r.EducationLevel,
		Language:        r.Language,
		IsMandatory:     r.IsMandatory,
		Priority:        r.Priority,
	}
}

// ToJobCategoryResponse maps JobCategory entity to JobCategoryResponse DTO
func ToJobCategoryResponse(c *job.JobCategory) *response.JobCategoryResponse {
	if c == nil {
		return nil
	}

	return &response.JobCategoryResponse{
		ID:           c.ID,
		CategoryName: c.Name,
		Description:  c.Description,
	}
}
