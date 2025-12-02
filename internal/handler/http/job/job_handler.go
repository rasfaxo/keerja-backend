package jobhandler

import (
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/job"
	"keerja-backend/internal/domain/master"
)

// JobHandler handles job-related operations
type JobHandler struct {
	jobService        job.JobService
	companyService    company.CompanyService
	jobOptionsService master.JobOptionsService
	skillsService     master.SkillsMasterService
}

// NewJobHandler creates a new instance of JobHandler
func NewJobHandler(
	jobService job.JobService,
	companyService company.CompanyService,
	jobOptionsService master.JobOptionsService,
	skillService master.SkillsMasterService,
) *JobHandler {
	return &JobHandler{
		jobService:        jobService,
		companyService:    companyService,
		jobOptionsService: jobOptionsService,
		skillsService:     skillService,
	}
}
