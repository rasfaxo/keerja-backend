package helpers

import (
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/job"
	"keerja-backend/internal/dto/request"
	"strings"
)

// ============================================================================
// Job Filter Builders
// ============================================================================

// BuildJobFilter builds job.JobFilter from request.JobSearchRequest
func BuildJobFilter(q request.JobSearchRequest) job.JobFilter {
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

	// Optional boolean fields
	if q.RemoteOnly {
		b := true
		f.RemoteOption = &b
	}

	// Optional numeric fields
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

	// Optional ID fields
	if q.CategoryID != nil {
		f.CategoryID = *q.CategoryID
	}
	if q.CompanyID != nil {
		f.CompanyID = *q.CompanyID
	}

	return f
}

// BuildJobSearchFilter builds job.JobSearchFilter from request.JobSearchRequest
func BuildJobSearchFilter(q request.JobSearchRequest) job.JobSearchFilter {
	f := job.JobSearchFilter{
		Keyword:           q.Query,
		Location:          q.Location,
		RemoteOnly:        q.RemoteOnly,
		MinSalary:         q.SalaryMin,
		MaxSalary:         q.SalaryMax,
		MinExperience:     q.ExperienceMin,
		MaxExperience:     q.ExperienceMax,
		PostedWithin:      q.PostedWithin,
		EducationLevelID:  q.EducationLevelID,
		ExperienceLevelID: q.ExperienceLevelID,
	}

	// Optional ID fields
	if q.CategoryID != nil {
		f.CategoryIDs = []int64{*q.CategoryID}
	}
	if q.CompanyID != nil {
		f.CompanyIDs = []int64{*q.CompanyID}
	}

	// Optional array fields
	if len(q.SkillIDs) > 0 {
		f.SkillIDs = q.SkillIDs
	}

	// Master Data ID array filters (UI: Job Type & Work Policy chips)
	if len(q.JobTypeIDs) > 0 {
		f.JobTypeIDs = q.JobTypeIDs
	}
	if len(q.WorkPolicyIDs) > 0 {
		f.WorkPolicyIDs = q.WorkPolicyIDs
	}

	// Optional string array fields (legacy support)
	if q.EmploymentType != "" {
		f.EmploymentTypes = []string{q.EmploymentType}
	}
	if q.JobLevel != "" {
		f.JobLevels = []string{q.JobLevel}
	}
	if q.EducationLevel != "" {
		f.EducationLevels = []string{q.EducationLevel}
	}

	return f
}

// ============================================================================
// Company Filter Builders
// ============================================================================

// BuildCompanyFilter builds company.CompanyFilter from request.CompanySearchRequest
func BuildCompanyFilter(q request.CompanySearchRequest) *company.CompanyFilter {
	filter := &company.CompanyFilter{
		Verified:  q.IsVerified,
		Page:      q.Page,
		Limit:     q.Limit,
		SortBy:    q.SortBy,
		SortOrder: q.SortOrder,
	}

	// Optional string filters
	if q.Industry != "" {
		v := q.Industry
		filter.Industry = &v
	}
	if q.CompanyType != "" {
		v := q.CompanyType
		filter.CompanyType = &v
	}
	if q.SizeCategory != "" {
		v := q.SizeCategory
		filter.SizeCategory = &v
	}
	if q.Location != "" {
		v := q.Location
		filter.City = &v
	}
	if q.Query != "" {
		v := q.Query
		filter.SearchQuery = &v
	}

	return filter
}
