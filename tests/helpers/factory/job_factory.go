package factory

import (
	"fmt"
	"time"
)

// JobFactory provides methods to create test jobs
type JobFactory struct {
	sequence int
}

// NewJobFactory creates a new job factory
func NewJobFactory() *JobFactory {
	return &JobFactory{
		sequence: 0,
	}
}

// JobBuilder provides a fluent interface for building jobs
type JobBuilder struct {
	ID               int64
	CompanyID        int64
	CreatedByUserID  int64
	Title            string
	Description      string
	Requirements     string
	Responsibilities string
	Benefits         string
	Location         string
	LocationType     string // remote, onsite, hybrid
	EmploymentType   string
	ExperienceLevel  string
	SalaryMin        *int64
	SalaryMax        *int64
	SalaryCurrency   string
	Skills           []string
	Status           string
	ViewCount        int
	ApplicationCount int
	IsFeatured       bool
	PostedAt         time.Time
	ClosingDate      *time.Time
	ExpiresAt        *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Build builds the job
func (jb *JobBuilder) Build() *JobBuilder {
	// Set default timestamps if not set
	if jb.CreatedAt.IsZero() {
		jb.CreatedAt = time.Now()
	}
	if jb.UpdatedAt.IsZero() {
		jb.UpdatedAt = time.Now()
	}
	if jb.PostedAt.IsZero() {
		jb.PostedAt = time.Now()
	}

	// Set default status
	if jb.Status == "" {
		jb.Status = "open"
	}

	// Set default currency
	if jb.SalaryCurrency == "" {
		jb.SalaryCurrency = "IDR"
	}

	return jb
}

// WithID sets the job ID
func (jb *JobBuilder) WithID(id int64) *JobBuilder {
	jb.ID = id
	return jb
}

// WithCompanyID sets the company ID
func (jb *JobBuilder) WithCompanyID(id int64) *JobBuilder {
	jb.CompanyID = id
	return jb
}

// WithCreatedBy sets the creator user ID
func (jb *JobBuilder) WithCreatedBy(userID int64) *JobBuilder {
	jb.CreatedByUserID = userID
	return jb
}

// WithTitle sets the job title
func (jb *JobBuilder) WithTitle(title string) *JobBuilder {
	jb.Title = title
	return jb
}

// WithDescription sets the description
func (jb *JobBuilder) WithDescription(desc string) *JobBuilder {
	jb.Description = desc
	return jb
}

// WithRequirements sets the requirements
func (jb *JobBuilder) WithRequirements(req string) *JobBuilder {
	jb.Requirements = req
	return jb
}

// WithResponsibilities sets the responsibilities
func (jb *JobBuilder) WithResponsibilities(resp string) *JobBuilder {
	jb.Responsibilities = resp
	return jb
}

// WithBenefits sets the benefits
func (jb *JobBuilder) WithBenefits(benefits string) *JobBuilder {
	jb.Benefits = benefits
	return jb
}

// WithLocation sets the location
func (jb *JobBuilder) WithLocation(location string) *JobBuilder {
	jb.Location = location
	return jb
}

// Remote sets the job as remote
func (jb *JobBuilder) Remote() *JobBuilder {
	jb.LocationType = "remote"
	return jb
}

// Onsite sets the job as onsite
func (jb *JobBuilder) Onsite() *JobBuilder {
	jb.LocationType = "onsite"
	return jb
}

// Hybrid sets the job as hybrid
func (jb *JobBuilder) Hybrid() *JobBuilder {
	jb.LocationType = "hybrid"
	return jb
}

// WithEmploymentType sets the employment type
func (jb *JobBuilder) WithEmploymentType(empType string) *JobBuilder {
	jb.EmploymentType = empType
	return jb
}

// FullTime sets the job as full time
func (jb *JobBuilder) FullTime() *JobBuilder {
	jb.EmploymentType = "full_time"
	return jb
}

// PartTime sets the job as part time
func (jb *JobBuilder) PartTime() *JobBuilder {
	jb.EmploymentType = "part_time"
	return jb
}

// Contract sets the job as contract
func (jb *JobBuilder) Contract() *JobBuilder {
	jb.EmploymentType = "contract"
	return jb
}

// Freelance sets the job as freelance
func (jb *JobBuilder) Freelance() *JobBuilder {
	jb.EmploymentType = "freelance"
	return jb
}

// Internship sets the job as internship
func (jb *JobBuilder) Internship() *JobBuilder {
	jb.EmploymentType = "internship"
	return jb
}

// WithExperienceLevel sets the experience level
func (jb *JobBuilder) WithExperienceLevel(level string) *JobBuilder {
	jb.ExperienceLevel = level
	return jb
}

// EntryLevel sets the experience level as entry level
func (jb *JobBuilder) EntryLevel() *JobBuilder {
	jb.ExperienceLevel = "entry_level"
	return jb
}

// Junior sets the experience level as junior
func (jb *JobBuilder) Junior() *JobBuilder {
	jb.ExperienceLevel = "junior"
	return jb
}

// MidLevel sets the experience level as mid level
func (jb *JobBuilder) MidLevel() *JobBuilder {
	jb.ExperienceLevel = "mid_level"
	return jb
}

// Senior sets the experience level as senior
func (jb *JobBuilder) Senior() *JobBuilder {
	jb.ExperienceLevel = "senior"
	return jb
}

// Lead sets the experience level as lead
func (jb *JobBuilder) Lead() *JobBuilder {
	jb.ExperienceLevel = "lead"
	return jb
}

// Manager sets the experience level as manager
func (jb *JobBuilder) Manager() *JobBuilder {
	jb.ExperienceLevel = "manager"
	return jb
}

// WithSalary sets the salary range
func (jb *JobBuilder) WithSalary(min, max int64, currency string) *JobBuilder {
	jb.SalaryMin = &min
	jb.SalaryMax = &max
	jb.SalaryCurrency = currency
	return jb
}

// WithSalaryIDR sets the salary range in IDR
func (jb *JobBuilder) WithSalaryIDR(min, max int64) *JobBuilder {
	return jb.WithSalary(min, max, "IDR")
}

// WithSalaryUSD sets the salary range in USD
func (jb *JobBuilder) WithSalaryUSD(min, max int64) *JobBuilder {
	return jb.WithSalary(min, max, "USD")
}

// WithSkills sets the required skills
func (jb *JobBuilder) WithSkills(skills []string) *JobBuilder {
	jb.Skills = skills
	return jb
}

// WithStatus sets the job status
func (jb *JobBuilder) WithStatus(status string) *JobBuilder {
	jb.Status = status
	return jb
}

// Open sets the job status as open
func (jb *JobBuilder) Open() *JobBuilder {
	jb.Status = "open"
	return jb
}

// Closed sets the job status as closed
func (jb *JobBuilder) Closed() *JobBuilder {
	jb.Status = "closed"
	return jb
}

// Draft sets the job status as draft
func (jb *JobBuilder) Draft() *JobBuilder {
	jb.Status = "draft"
	return jb
}

// Filled sets the job status as filled
func (jb *JobBuilder) Filled() *JobBuilder {
	jb.Status = "filled"
	return jb
}

// Cancelled sets the job status as cancelled
func (jb *JobBuilder) Cancelled() *JobBuilder {
	jb.Status = "cancelled"
	return jb
}

// Featured marks the job as featured
func (jb *JobBuilder) Featured() *JobBuilder {
	jb.IsFeatured = true
	return jb
}

// WithViewCount sets the view count
func (jb *JobBuilder) WithViewCount(count int) *JobBuilder {
	jb.ViewCount = count
	return jb
}

// WithApplicationCount sets the application count
func (jb *JobBuilder) WithApplicationCount(count int) *JobBuilder {
	jb.ApplicationCount = count
	return jb
}

// WithClosingDate sets the closing date
func (jb *JobBuilder) WithClosingDate(date time.Time) *JobBuilder {
	jb.ClosingDate = &date
	return jb
}

// WithExpiryDate sets the expiry date
func (jb *JobBuilder) WithExpiryDate(date time.Time) *JobBuilder {
	jb.ExpiresAt = &date
	return jb
}

// Expired marks the job as expired
func (jb *JobBuilder) Expired() *JobBuilder {
	past := time.Now().AddDate(0, 0, -1)
	jb.ExpiresAt = &past
	jb.Status = "closed"
	return jb
}

// WithPostedAt sets the posted at time
func (jb *JobBuilder) WithPostedAt(t time.Time) *JobBuilder {
	jb.PostedAt = t
	return jb
}

// WithCreatedAt sets the created at time
func (jb *JobBuilder) WithCreatedAt(t time.Time) *JobBuilder {
	jb.CreatedAt = t
	return jb
}

// WithUpdatedAt sets the updated at time
func (jb *JobBuilder) WithUpdatedAt(t time.Time) *JobBuilder {
	jb.UpdatedAt = t
	return jb
}

// CreateJob creates a job builder with default values
func (f *JobFactory) CreateJob() *JobBuilder {
	f.sequence++
	salaryMin := int64(5000000)
	salaryMax := int64(10000000)

	return &JobBuilder{
		ID:               int64(f.sequence),
		CompanyID:        1,
		CreatedByUserID:  1,
		Title:            fmt.Sprintf("Software Engineer %d", f.sequence),
		Description:      "We are looking for a talented software engineer to join our team.",
		Requirements:     "Bachelor's degree in Computer Science or related field. 2+ years of experience.",
		Responsibilities: "Develop and maintain software applications. Collaborate with team members.",
		Benefits:         "Health insurance, flexible working hours, remote work options.",
		Location:         "Jakarta, Indonesia",
		LocationType:     "onsite",
		EmploymentType:   "full_time",
		ExperienceLevel:  "mid_level",
		SalaryMin:        &salaryMin,
		SalaryMax:        &salaryMax,
		SalaryCurrency:   "IDR",
		Skills:           []string{"Go", "PostgreSQL", "Docker"},
		Status:           "open",
		IsFeatured:       false,
		ViewCount:        0,
		ApplicationCount: 0,
	}
}

// CreateJobForCompany creates a job for a specific company
func (f *JobFactory) CreateJobForCompany(companyID int64) *JobBuilder {
	return f.CreateJob().WithCompanyID(companyID)
}

// CreateJobByUser creates a job created by a specific user
func (f *JobFactory) CreateJobByUser(userID int64) *JobBuilder {
	return f.CreateJob().WithCreatedBy(userID)
}

// CreateRemoteJob creates a remote job
func (f *JobFactory) CreateRemoteJob() *JobBuilder {
	return f.CreateJob().Remote().WithLocation("Remote")
}

// CreateOpenJob creates an open job
func (f *JobFactory) CreateOpenJob() *JobBuilder {
	return f.CreateJob().Open()
}

// CreateClosedJob creates a closed job
func (f *JobFactory) CreateClosedJob() *JobBuilder {
	return f.CreateJob().Closed()
}

// CreateExpiredJob creates an expired job
func (f *JobFactory) CreateExpiredJob() *JobBuilder {
	return f.CreateJob().Expired()
}

// CreateFeaturedJob creates a featured job
func (f *JobFactory) CreateFeaturedJob() *JobBuilder {
	return f.CreateJob().Featured()
}

// CreateMultipleJobs creates multiple jobs
func (f *JobFactory) CreateMultipleJobs(count int) []*JobBuilder {
	jobs := make([]*JobBuilder, count)
	for i := 0; i < count; i++ {
		jobs[i] = f.CreateJob()
	}
	return jobs
}

// CreateMultipleJobsForCompany creates multiple jobs for a company
func (f *JobFactory) CreateMultipleJobsForCompany(companyID int64, count int) []*JobBuilder {
	jobs := make([]*JobBuilder, count)
	for i := 0; i < count; i++ {
		jobs[i] = f.CreateJobForCompany(companyID)
	}
	return jobs
}
