package mapper

import (
	"keerja-backend/internal/domain/job"
	"keerja-backend/internal/dto/response"
)

// Job Entity to Response Mappers

// ToJobResponse maps Job entity to JobResponse DTO
func ToJobResponse(j *job.Job) *response.JobResponse {
	if j == nil {
		return nil
	}

	var daysRemaining *int
	if j.ExpiredAt != nil && j.IsActive() {
		days := int(j.ExpiredAt.Sub(j.CreatedAt).Hours() / 24)
		daysRemaining = &days
	}

	return &response.JobResponse{
		ID:                j.ID,
		UUID:              j.UUID.String(),
		CompanyID:         j.CompanyID,
		Title:             j.Title,
		Slug:              j.Slug,
		JobLevel:          j.JobLevel,
		EmploymentType:    j.EmploymentType,
		Location:          j.Location,
		City:              j.City,
		Province:          j.Province,
		RemoteOption:      j.RemoteOption,
		SalaryMin:         j.SalaryMin,
		SalaryMax:         j.SalaryMax,
		Currency:          j.Currency,
		ExperienceMin:     j.ExperienceMin,
		ExperienceMax:     j.ExperienceMax,
		EducationLevel:    j.EducationLevel,
		Status:            j.Status,
		ViewsCount:        j.ViewsCount,
		ApplicationsCount: j.ApplicationsCount,
		PublishedAt:       j.PublishedAt,
		ExpiredAt:         j.ExpiredAt,
		CreatedAt:         j.CreatedAt,
		IsExpired:         j.IsExpired(),
		DaysRemaining:     daysRemaining,
	}
}

// ToJobDetailResponse maps Job entity with relations to JobDetailResponse DTO
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
		ID:                j.ID,
		UUID:              j.UUID.String(),
		CompanyID:         j.CompanyID,
		EmployerUserID:    j.EmployerUserID,
		CategoryID:        j.CategoryID,
		Title:             j.Title,
		Slug:              j.Slug,
		JobLevel:          j.JobLevel,
		EmploymentType:    j.EmploymentType,
		Description:       j.Description,
		Requirements:      j.RequirementsText,
		Responsibilities:  j.Responsibilities,
		Location:          j.Location,
		City:              j.City,
		Province:          j.Province,
		RemoteOption:      j.RemoteOption,
		SalaryMin:         j.SalaryMin,
		SalaryMax:         j.SalaryMax,
		Currency:          j.Currency,
		ExperienceMin:     j.ExperienceMin,
		ExperienceMax:     j.ExperienceMax,
		EducationLevel:    j.EducationLevel,
		TotalHires:        j.TotalHires,
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
		resp.Requirements_ = make([]response.JobRequirementResponse, len(j.JobRequirements))
		for i, req := range j.JobRequirements {
			resp.Requirements_[i] = *ToJobRequirementResponse(&req)
		}
	}

	return resp
}

// ToJobSkillResponse maps JobSkill entity to JobSkillResponse DTO
func ToJobSkillResponse(s *job.JobSkill) *response.JobSkillResponse {
	if s == nil {
		return nil
	}

	return &response.JobSkillResponse{
		ID:          s.ID,
		SkillID:     s.SkillID,
		IsRequired:  s.IsRequired(),
		MinYearsExp: nil, // Not available in entity
	}
}

// ToJobBenefitResponse maps JobBenefit entity to JobBenefitResponse DTO
func ToJobBenefitResponse(b *job.JobBenefit) *response.JobBenefitResponse {
	if b == nil {
		return nil
	}

	var benefitID int64
	if b.BenefitID != nil {
		benefitID = *b.BenefitID
	}

	return &response.JobBenefitResponse{
		ID:          b.ID,
		BenefitID:   benefitID,
		BenefitName: b.BenefitName,
		Description: b.Description,
	}
}

// ToJobLocationResponse maps JobLocation entity to JobLocationResponse DTO
func ToJobLocationResponse(l *job.JobLocation) *response.JobLocationResponse {
	if l == nil {
		return nil
	}

	return &response.JobLocationResponse{
		ID:        l.ID,
		City:      l.City,
		Province:  l.Province,
		Address:   l.Address,
		IsRemote:  l.IsRemote(),
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
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
		Description:     r.RequirementText,
		IsRequired:      r.IsMandatory,
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
