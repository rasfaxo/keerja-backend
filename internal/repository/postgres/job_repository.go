package postgres

import (
	"context"
	"strings"
	"time"

	"keerja-backend/internal/domain/job"

	"gorm.io/gorm"
)

// jobRepository implements job.JobRepository
type jobRepository struct {
	db *gorm.DB
}

// NewJobRepository creates a new instance of JobRepository
func NewJobRepository(db *gorm.DB) job.JobRepository {
	return &jobRepository{db: db}
}

// ===========================================
// JOB CRUD OPERATIONS
// ===========================================

// Create creates a new job
func (r *jobRepository) Create(ctx context.Context, j *job.Job) error {
	return r.db.WithContext(ctx).Create(j).Error
}

// FindByID finds a job by ID
func (r *jobRepository) FindByID(ctx context.Context, id int64) (*job.Job, error) {
	var j job.Job
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("JobSubcategory").
		Preload("Locations").
		Preload("Benefits").
		Preload("Skills.Skill").
		Preload("JobRequirements").
		// Preload Job Master Data Relations (New)
		Preload("JobTitle").
		Preload("JobType").
		Preload("WorkPolicy").
		Preload("EducationLevelM").
		Preload("ExperienceLevelM").
		Preload("GenderPreference").
		First(&j, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &j, nil
}

// FindByUUID finds a job by UUID
func (r *jobRepository) FindByUUID(ctx context.Context, uuid string) (*job.Job, error) {
	var j job.Job
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("JobSubcategory").
		Preload("Locations").
		Preload("Benefits").
		Preload("Skills.Skill").
		Preload("JobRequirements").
		// Preload Job Master Data Relations (New)
		Preload("JobTitle").
		Preload("JobType").
		Preload("WorkPolicy").
		Preload("EducationLevelM").
		Preload("ExperienceLevelM").
		Preload("GenderPreference").
		Where("uuid = ?", uuid).
		First(&j).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &j, nil
}

// FindBySlug finds a job by slug
func (r *jobRepository) FindBySlug(ctx context.Context, slug string) (*job.Job, error) {
	var j job.Job
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("JobSubcategory").
		Preload("Locations").
		Preload("Benefits").
		Preload("Skills.Skill").
		Preload("JobRequirements").
		Where("slug = ?", slug).
		First(&j).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &j, nil
}

// Update updates a job
func (r *jobRepository) Update(ctx context.Context, j *job.Job) error {
	// Use Updates with Select to avoid GORM overwriting our pointer values
	// This ensures we only update the specific fields we want (all fields except CompanyID)
	return r.db.WithContext(ctx).Model(j).Select(
		// Basic fields
		"Title", "Slug", "Description",
		"JobLevel", "EmploymentType",
		"Requirements", "Responsibilities",

		// Master Data IDs
		"JobTitleID", "JobTypeID", "WorkPolicyID",
		"EducationLevelID", "ExperienceLevelID", "GenderPreferenceID",
		"CategoryID",

		// Location fields
		"Location", "City", "Province", "RemoteOption",

		// Salary fields
		"SalaryMin", "SalaryMax", "SalaryDisplay", "MinAge", "MaxAge", "Currency",

		// Experience and Education (legacy fields)
		"ExperienceMin", "ExperienceMax", "EducationLevel",

		// Job metadata
		"TotalHires", "Status",
		"ViewsCount", "ApplicationsCount",

		// Dates
		"PublishedAt", "ExpiredAt", "UpdatedAt",
	).Updates(j).Error
}

// Delete hard deletes a job
func (r *jobRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&job.Job{}, id).Error
}

// SoftDelete soft deletes a job
func (r *jobRepository) SoftDelete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&job.Job{}, id).Error
}

// ===========================================
// JOB LISTING AND SEARCH
// ===========================================

// List retrieves jobs with filtering and pagination
func (r *jobRepository) List(ctx context.Context, filter job.JobFilter, page, limit int) ([]job.Job, int64, error) {
	var jobs []job.Job
	var total int64

	query := r.db.WithContext(ctx).Model(&job.Job{})
	query = r.applyJobFilter(query, filter)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Apply sorting
	query = r.applySorting(query, filter.SortBy)

	// Execute query
	err := query.
		Preload("Category").
		Preload("JobSubcategory").
		Preload("Locations").
		Preload("Benefits").
		Preload("Skills.Skill").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error

	return jobs, total, err
}

// ListByCompany retrieves jobs by company ID
func (r *jobRepository) ListByCompany(ctx context.Context, companyID int64, filter job.JobFilter, page, limit int) ([]job.Job, int64, error) {
	filter.CompanyID = companyID
	return r.List(ctx, filter, page, limit)
}

// ListByEmployer retrieves jobs by employer user ID
func (r *jobRepository) ListByEmployer(ctx context.Context, employerUserID int64, filter job.JobFilter, page, limit int) ([]job.Job, int64, error) {
	var jobs []job.Job
	var total int64

	query := r.db.WithContext(ctx).Model(&job.Job{}).
		Where("employer_user_id = ?", employerUserID)
	query = r.applyJobFilter(query, filter)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Apply sorting
	query = r.applySorting(query, filter.SortBy)

	err := query.
		Preload("Category").
		Preload("JobSubcategory").
		Preload("JobTitle").
		Preload("JobType").
		Preload("WorkPolicy").
		Preload("EducationLevelM").
		Preload("ExperienceLevelM").
		Preload("GenderPreference").
		Preload("Locations").
		Preload("Benefits").
		Preload("Skills.Skill").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error

	return jobs, total, err
}

// SearchJobs performs advanced job search
func (r *jobRepository) SearchJobs(ctx context.Context, filter job.JobSearchFilter, page, limit int) ([]job.Job, int64, error) {
	var jobs []job.Job
	var total int64

	query := r.db.WithContext(ctx).Model(&job.Job{})

	// Keyword search on title and description
	if filter.Keyword != "" {
		kw := "%" + strings.ToLower(filter.Keyword) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", kw, kw)
	}

	// Location (match city, province or location text)
	if filter.Location != "" {
		loc := "%" + strings.ToLower(filter.Location) + "%"
		query = query.Where("LOWER(city) = ? OR LOWER(province) = ? OR LOWER(location) LIKE ?", strings.ToLower(filter.Location), strings.ToLower(filter.Location), loc)
	}

	// Category filter
	if len(filter.CategoryIDs) > 0 {
		query = query.Where("category_id IN ?", filter.CategoryIDs)
	}

	// Job levels filter
	if len(filter.JobLevels) > 0 {
		query = query.Where("job_level IN ?", filter.JobLevels)
	}

	// Employment types filter
	if len(filter.EmploymentTypes) > 0 {
		query = query.Where("employment_type IN ?", filter.EmploymentTypes)
	}

	// Remote only filter
	if filter.RemoteOnly {
		query = query.Where("remote_option = ?", true)
	}

	// Salary range filter
	if filter.MinSalary != nil {
		query = query.Where("salary_max >= ? OR salary_max IS NULL", *filter.MinSalary)
	}
	if filter.MaxSalary != nil {
		query = query.Where("salary_min <= ? OR salary_min IS NULL", *filter.MaxSalary)
	}

	// Experience filter
	if filter.MinExperience != nil {
		query = query.Where("experience_max >= ? OR experience_max IS NULL", *filter.MinExperience)
	}
	if filter.MaxExperience != nil {
		query = query.Where("experience_min <= ? OR experience_min IS NULL", *filter.MaxExperience)
	}

	// Education levels filter
	if len(filter.EducationLevels) > 0 {
		query = query.Where("education_level IN ?", filter.EducationLevels)
	}

	// Company IDs filter
	if len(filter.CompanyIDs) > 0 {
		query = query.Where("company_id IN ?", filter.CompanyIDs)
	}

	// Posted within filter (days)
	if filter.PostedWithin != nil && *filter.PostedWithin > 0 {
		daysAgo := time.Now().AddDate(0, 0, -*filter.PostedWithin)
		query = query.Where("published_at >= ?", daysAgo)
	}

	// Skills filter (ensure jobs contain all requested skills)
	if len(filter.SkillIDs) > 0 {
		query = query.Joins("INNER JOIN job_skills ON job_skills.job_id = jobs.id").
			Where("job_skills.skill_id IN ?", filter.SkillIDs).
			Group("jobs.id").
			Having("COUNT(DISTINCT job_skills.skill_id) = ?", len(filter.SkillIDs))
	}

	// Only active (published and not expired)
	query = query.Where("status = ?", "published").
		Where("(expired_at IS NULL OR expired_at > ?)", time.Now())

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Execute final query
	err := query.
		Preload("Category").
		Preload("JobSubcategory").
		Preload("Locations").
		Preload("Benefits").
		Preload("Skills.Skill").
		Order("published_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error

	return jobs, total, err
}

// GetJobsGroupedByStatus returns jobs grouped by status for a user
func (r *jobRepository) GetJobsGroupedByStatus(ctx context.Context, userID int64) (map[string][]job.Job, error) {
	var jobs []job.Job
	err := r.db.WithContext(ctx).
		Model(&job.Job{}).
		Joins("JOIN employer_users eu ON eu.id = jobs.employer_user_id").
		Where("eu.user_id = ?", userID).
		Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	grouped := map[string][]job.Job{
		"active":    {},
		"inactive":  {},
		"draft":     {},
		"in_review": {},
	}
	for _, j := range jobs {
		switch j.Status {
		case "published":
			// published jobs are considered 'active' in the mobile UI
			grouped["active"] = append(grouped["active"], j)
		case "inactive":
			grouped["inactive"] = append(grouped["inactive"], j)
		case "draft":
			grouped["draft"] = append(grouped["draft"], j)
		case "in_review", "pending_review":
			// both in_review and pending_review map to the in_review tab
			grouped["in_review"] = append(grouped["in_review"], j)
		}
	}
	return grouped, nil
}

// ===========================================
// JOB STATUS OPERATIONS
// ===========================================

// UpdateStatusByCompany updates all jobs for a company from one status to another
func (r *jobRepository) UpdateStatusByCompany(ctx context.Context, companyID int64, fromStatus, toStatus string) error {
	return r.db.WithContext(ctx).
		Model(&job.Job{}).
		Where("company_id = ? AND status = ?", companyID, fromStatus).
		Update("status", toStatus).Error
}

// UpdateStatus updates job status
func (r *jobRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).
		Model(&job.Job{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// UpdateStatusWithExpiry updates job status and optionally sets published_at and expired_at
func (r *jobRepository) UpdateStatusWithExpiry(ctx context.Context, id int64, status string, publishedAt *time.Time, expiredAt *time.Time) error {
	updates := map[string]interface{}{}
	updates["status"] = status
	if publishedAt != nil {
		updates["published_at"] = *publishedAt
	}
	if expiredAt != nil {
		updates["expired_at"] = *expiredAt
	}

	return r.db.WithContext(ctx).
		Model(&job.Job{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// PublishJob publishes a job
func (r *jobRepository) PublishJob(ctx context.Context, id int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&job.Job{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       "published",
			"published_at": now,
		}).Error
}

// CloseJob closes a job
func (r *jobRepository) CloseJob(ctx context.Context, id int64) error {
	return r.UpdateStatus(ctx, id, "closed")
}

// ExpireJob marks a job as expired
func (r *jobRepository) ExpireJob(ctx context.Context, id int64) error {
	return r.UpdateStatus(ctx, id, "expired")
}

// SuspendJob suspends a job
func (r *jobRepository) SuspendJob(ctx context.Context, id int64) error {
	return r.UpdateStatus(ctx, id, "suspended")
}

// GetExpiredJobs retrieves jobs that have expired
func (r *jobRepository) GetExpiredJobs(ctx context.Context) ([]job.Job, error) {
	var jobs []job.Job
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("status = ? AND expired_at <= ?", "published", now).
		Find(&jobs).Error
	return jobs, err
}

// GetExpiringJobs retrieves jobs expiring within specified days
func (r *jobRepository) GetExpiringJobs(ctx context.Context, days int) ([]job.Job, error) {
	var jobs []job.Job
	now := time.Now()
	futureDate := now.AddDate(0, 0, days)
	err := r.db.WithContext(ctx).
		Where("status = ? AND expired_at > ? AND expired_at <= ?", "published", now, futureDate).
		Find(&jobs).Error
	return jobs, err
}

// ===========================================
// JOB STATISTICS
// ===========================================

// IncrementViews increments job view count
func (r *jobRepository) IncrementViews(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&job.Job{}).
		Where("id = ?", id).
		UpdateColumn("views_count", gorm.Expr("views_count + ?", 1)).Error
}

// IncrementApplications increments job application count
func (r *jobRepository) IncrementApplications(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&job.Job{}).
		Where("id = ?", id).
		UpdateColumn("applications_count", gorm.Expr("applications_count + ?", 1)).Error
}

// GetJobStats retrieves job statistics
func (r *jobRepository) GetJobStats(ctx context.Context, jobID int64) (*job.JobStats, error) {
	var j job.Job
	if err := r.db.WithContext(ctx).First(&j, jobID).Error; err != nil {
		return nil, err
	}

	stats := &job.JobStats{
		JobID:             j.ID,
		ViewsCount:        j.ViewsCount,
		ApplicationsCount: j.ApplicationsCount,
	}

	// Calculate conversion rate
	if j.ViewsCount > 0 {
		stats.ConversionRate = float64(j.ApplicationsCount) / float64(j.ViewsCount) * 100
	}

	return stats, nil
}

// GetCompanyJobStats retrieves company's job statistics
func (r *jobRepository) GetCompanyJobStats(ctx context.Context, companyID int64) (*job.CompanyJobStats, error) {
	var stats job.CompanyJobStats
	stats.CompanyID = companyID

	// Total jobs
	r.db.WithContext(ctx).Model(&job.Job{}).
		Where("company_id = ?", companyID).
		Count(&stats.TotalJobs)

	// Active jobs
	r.db.WithContext(ctx).Model(&job.Job{}).
		Where("company_id = ? AND status = ?", companyID, "published").
		Where("(expired_at IS NULL OR expired_at > ?)", time.Now()).
		Count(&stats.ActiveJobs)

	// Draft jobs
	r.db.WithContext(ctx).Model(&job.Job{}).
		Where("company_id = ? AND status = ?", companyID, "draft").
		Count(&stats.DraftJobs)

	// Closed jobs
	r.db.WithContext(ctx).Model(&job.Job{}).
		Where("company_id = ? AND status = ?", companyID, "closed").
		Count(&stats.ClosedJobs)

	// Expired jobs
	r.db.WithContext(ctx).Model(&job.Job{}).
		Where("company_id = ? AND status = ?", companyID, "expired").
		Count(&stats.ExpiredJobs)

	// Total views and applications
	var result struct {
		TotalViews        int64
		TotalApplications int64
	}
	r.db.WithContext(ctx).Model(&job.Job{}).
		Where("company_id = ?", companyID).
		Select("COALESCE(SUM(views_count), 0) as total_views, COALESCE(SUM(applications_count), 0) as total_applications").
		Scan(&result)

	stats.TotalViews = result.TotalViews
	stats.TotalApplications = result.TotalApplications

	// Calculate averages
	if stats.TotalJobs > 0 {
		stats.AverageViewsPerJob = float64(stats.TotalViews) / float64(stats.TotalJobs)
		stats.AverageApplicationsPerJob = float64(stats.TotalApplications) / float64(stats.TotalJobs)
	}

	return &stats, nil
}

// ===========================================
// RECOMMENDATION AND MATCHING
// ===========================================

// GetRecommendedJobs retrieves recommended jobs for a user
func (r *jobRepository) GetRecommendedJobs(ctx context.Context, userID int64, limit int) ([]job.Job, error) {
	var jobs []job.Job
	// Simplified recommendation: get latest published jobs
	// In production, implement ML-based recommendation
	err := r.db.WithContext(ctx).
		Where("status = ?", "published").
		Where("(expired_at IS NULL OR expired_at > ?)", time.Now()).
		Order("published_at DESC").
		Limit(limit).
		Preload("Category").
		Preload("Locations").
		Preload("Benefits").
		Find(&jobs).Error
	return jobs, err
}

// GetSimilarJobs retrieves similar jobs to a given job
func (r *jobRepository) GetSimilarJobs(ctx context.Context, jobID int64, limit int) ([]job.Job, error) {
	var targetJob job.Job
	if err := r.db.WithContext(ctx).First(&targetJob, jobID).Error; err != nil {
		return nil, err
	}

	var jobs []job.Job
	query := r.db.WithContext(ctx).
		Where("id != ?", jobID).
		Where("status = ?", "published").
		Where("(expired_at IS NULL OR expired_at > ?)", time.Now())

	// Match by category or location
	if targetJob.CategoryID != nil {
		query = query.Where("category_id = ? OR city = ?", *targetJob.CategoryID, targetJob.City)
	} else {
		query = query.Where("city = ?", targetJob.City)
	}

	err := query.
		Order("published_at DESC").
		Limit(limit).
		Preload("Category").
		Preload("Locations").
		Find(&jobs).Error

	return jobs, err
}

// GetMatchingJobs retrieves jobs matching user profile
func (r *jobRepository) GetMatchingJobs(ctx context.Context, userID int64, filter job.JobFilter, page, limit int) ([]job.Job, int64, error) {
	// Simplified matching: return published jobs
	// In production, match against user skills, experience, preferences
	filter.Status = "published"
	isActive := true
	filter.IsActive = &isActive
	return r.List(ctx, filter, page, limit)
}

// ===========================================
// ADVANCED SEARCH
// ===========================================

// SearchByLocation searches jobs near a location
func (r *jobRepository) SearchByLocation(ctx context.Context, latitude, longitude, radius float64, filter job.JobFilter, page, limit int) ([]job.Job, int64, error) {
	var jobs []job.Job
	var total int64

	// Use Haversine formula for distance calculation
	query := r.db.WithContext(ctx).
		Table("jobs").
		Select("jobs.*").
		Joins("INNER JOIN job_locations ON job_locations.job_id = jobs.id").
		Where(`(
			6371 * acos(
				cos(radians(?)) * cos(radians(job_locations.latitude)) *
				cos(radians(job_locations.longitude) - radians(?)) +
				sin(radians(?)) * sin(radians(job_locations.latitude))
			)
		) <= ?`, latitude, longitude, latitude, radius)

	query = r.applyJobFilter(query, filter)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	err := query.
		Preload("Category").
		Preload("Locations").
		Preload("Benefits").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error

	return jobs, total, err
}

// SearchBySkills searches jobs by skill requirements
func (r *jobRepository) SearchBySkills(ctx context.Context, skillIDs []int64, filter job.JobFilter, page, limit int) ([]job.Job, int64, error) {
	var jobs []job.Job
	var total int64

	query := r.db.WithContext(ctx).
		Model(&job.Job{}).
		Joins("INNER JOIN job_skills ON job_skills.job_id = jobs.id").
		Where("job_skills.skill_id IN ?", skillIDs)

	query = r.applyJobFilter(query, filter)

	query = query.Group("jobs.id")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	err := query.
		Preload("Category").
		Preload("Locations").
		Preload("Benefits").
		Preload("Skills.Skill").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error

	return jobs, total, err
}

// SearchBySalaryRange searches jobs within salary range
func (r *jobRepository) SearchBySalaryRange(ctx context.Context, minSalary, maxSalary float64, filter job.JobFilter, page, limit int) ([]job.Job, int64, error) {
	filter.MinSalary = &minSalary
	filter.MaxSalary = &maxSalary
	return r.List(ctx, filter, page, limit)
}

// ===========================================
// JOB CATEGORY OPERATIONS
// ===========================================

// CreateCategory creates a job category
func (r *jobRepository) CreateCategory(ctx context.Context, category *job.JobCategory) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// FindCategoryByID finds a category by ID
func (r *jobRepository) FindCategoryByID(ctx context.Context, id int64) (*job.JobCategory, error) {
	var category job.JobCategory
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Preload("Subcategories").
		First(&category, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// FindCategoryByCode finds a category by code
func (r *jobRepository) FindCategoryByCode(ctx context.Context, code string) (*job.JobCategory, error) {
	var category job.JobCategory
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("code = ?", code).
		First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// UpdateCategory updates a job category
func (r *jobRepository) UpdateCategory(ctx context.Context, category *job.JobCategory) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// DeleteCategory soft deletes a category
func (r *jobRepository) DeleteCategory(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&job.JobCategory{}, id).Error
}

// ListCategories retrieves categories with filtering
func (r *jobRepository) ListCategories(ctx context.Context, filter job.CategoryFilter, page, limit int) ([]job.JobCategory, int64, error) {
	var categories []job.JobCategory
	var total int64

	query := r.db.WithContext(ctx).Model(&job.JobCategory{})

	// Apply filters
	if filter.ParentID != nil {
		query = query.Where("parent_id = ?", *filter.ParentID)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.Keyword != "" {
		keyword := "%" + strings.ToLower(filter.Keyword) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ?", keyword, keyword)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	err := query.
		Preload("Parent").
		Preload("Children").
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Find(&categories).Error

	return categories, total, err
}

// GetCategoryTree retrieves hierarchical category tree
func (r *jobRepository) GetCategoryTree(ctx context.Context) ([]job.JobCategory, error) {
	var categories []job.JobCategory
	err := r.db.WithContext(ctx).
		Preload("Children").
		Preload("Subcategories").
		Preload("Children.Subcategories").
		Where("parent_id IS NULL AND is_active = ?", true).
		Order("name ASC").
		Find(&categories).Error
	return categories, err
}

// GetActiveCategories retrieves all active categories
func (r *jobRepository) GetActiveCategories(ctx context.Context) ([]job.JobCategory, error) {
	var categories []job.JobCategory
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&categories).Error
	return categories, err
}

// ===========================================
// JOB SUBCATEGORY OPERATIONS
// ===========================================

// CreateSubcategory creates a job subcategory
func (r *jobRepository) CreateSubcategory(ctx context.Context, subcategory *job.JobSubcategory) error {
	return r.db.WithContext(ctx).Create(subcategory).Error
}

// FindSubcategoryByID finds a subcategory by ID
func (r *jobRepository) FindSubcategoryByID(ctx context.Context, id int64) (*job.JobSubcategory, error) {
	var subcategory job.JobSubcategory
	err := r.db.WithContext(ctx).
		Preload("Category").
		First(&subcategory, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &subcategory, nil
}

// FindSubcategoryByCode finds a subcategory by code
func (r *jobRepository) FindSubcategoryByCode(ctx context.Context, code string) (*job.JobSubcategory, error) {
	var subcategory job.JobSubcategory
	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("code = ?", code).
		First(&subcategory).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &subcategory, nil
}

// UpdateSubcategory updates a job subcategory
func (r *jobRepository) UpdateSubcategory(ctx context.Context, subcategory *job.JobSubcategory) error {
	return r.db.WithContext(ctx).Save(subcategory).Error
}

// DeleteSubcategory soft deletes a subcategory
func (r *jobRepository) DeleteSubcategory(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&job.JobSubcategory{}, id).Error
}

// ListSubcategories retrieves subcategories by category ID
func (r *jobRepository) ListSubcategories(ctx context.Context, categoryID int64) ([]job.JobSubcategory, error) {
	var subcategories []job.JobSubcategory
	err := r.db.WithContext(ctx).
		Where("category_id = ?", categoryID).
		Order("name ASC").
		Find(&subcategories).Error
	return subcategories, err
}

// GetActiveSubcategories retrieves active subcategories by category ID
func (r *jobRepository) GetActiveSubcategories(ctx context.Context, categoryID int64) ([]job.JobSubcategory, error) {
	var subcategories []job.JobSubcategory
	err := r.db.WithContext(ctx).
		Where("category_id = ? AND is_active = ?", categoryID, true).
		Order("name ASC").
		Find(&subcategories).Error
	return subcategories, err
}

// ===========================================
// JOB LOCATION OPERATIONS
// ===========================================

// CreateLocation creates a job location
func (r *jobRepository) CreateLocation(ctx context.Context, location *job.JobLocation) error {
	return r.db.WithContext(ctx).Create(location).Error
}

// FindLocationByID finds a location by ID
func (r *jobRepository) FindLocationByID(ctx context.Context, id int64) (*job.JobLocation, error) {
	var location job.JobLocation
	err := r.db.WithContext(ctx).First(&location, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &location, nil
}

// UpdateLocation updates a job location
func (r *jobRepository) UpdateLocation(ctx context.Context, location *job.JobLocation) error {
	return r.db.WithContext(ctx).Save(location).Error
}

// DeleteLocation deletes a job location
func (r *jobRepository) DeleteLocation(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&job.JobLocation{}, id).Error
}

// ListLocationsByJob retrieves locations for a job
func (r *jobRepository) ListLocationsByJob(ctx context.Context, jobID int64) ([]job.JobLocation, error) {
	var locations []job.JobLocation
	err := r.db.WithContext(ctx).
		Where("job_id = ?", jobID).
		Order("is_primary DESC, created_at ASC").
		Find(&locations).Error
	return locations, err
}

// GetPrimaryLocation retrieves primary location for a job
func (r *jobRepository) GetPrimaryLocation(ctx context.Context, jobID int64) (*job.JobLocation, error) {
	var location job.JobLocation
	err := r.db.WithContext(ctx).
		Where("job_id = ? AND is_primary = ?", jobID, true).
		First(&location).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &location, nil
}

// SetPrimaryLocation sets a location as primary
func (r *jobRepository) SetPrimaryLocation(ctx context.Context, jobID, locationID int64) error {
	// First, unset all primary flags for this job
	if err := r.db.WithContext(ctx).
		Model(&job.JobLocation{}).
		Where("job_id = ?", jobID).
		Update("is_primary", false).Error; err != nil {
		return err
	}

	// Set the specified location as primary
	return r.db.WithContext(ctx).
		Model(&job.JobLocation{}).
		Where("id = ? AND job_id = ?", locationID, jobID).
		Update("is_primary", true).Error
}

// ===========================================
// JOB BENEFIT OPERATIONS
// ===========================================

// CreateBenefit creates a job benefit
func (r *jobRepository) CreateBenefit(ctx context.Context, benefit *job.JobBenefit) error {
	return r.db.WithContext(ctx).Create(benefit).Error
}

// FindBenefitByID finds a benefit by ID
func (r *jobRepository) FindBenefitByID(ctx context.Context, id int64) (*job.JobBenefit, error) {
	var benefit job.JobBenefit
	err := r.db.WithContext(ctx).First(&benefit, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &benefit, nil
}

// UpdateBenefit updates a job benefit
func (r *jobRepository) UpdateBenefit(ctx context.Context, benefit *job.JobBenefit) error {
	return r.db.WithContext(ctx).Save(benefit).Error
}

// DeleteBenefit deletes a job benefit
func (r *jobRepository) DeleteBenefit(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&job.JobBenefit{}, id).Error
}

// ListBenefitsByJob retrieves benefits for a job
func (r *jobRepository) ListBenefitsByJob(ctx context.Context, jobID int64) ([]job.JobBenefit, error) {
	var benefits []job.JobBenefit
	err := r.db.WithContext(ctx).
		Where("job_id = ?", jobID).
		Order("is_highlight DESC, benefit_name ASC").
		Find(&benefits).Error
	return benefits, err
}

// GetHighlightedBenefits retrieves highlighted benefits
func (r *jobRepository) GetHighlightedBenefits(ctx context.Context, jobID int64) ([]job.JobBenefit, error) {
	var benefits []job.JobBenefit
	err := r.db.WithContext(ctx).
		Where("job_id = ? AND is_highlight = ?", jobID, true).
		Order("benefit_name ASC").
		Find(&benefits).Error
	return benefits, err
}

// BulkCreateBenefits creates multiple benefits
func (r *jobRepository) BulkCreateBenefits(ctx context.Context, benefits []job.JobBenefit) error {
	return r.db.WithContext(ctx).Create(&benefits).Error
}

// BulkDeleteBenefits deletes all benefits for a job
func (r *jobRepository) BulkDeleteBenefits(ctx context.Context, jobID int64) error {
	return r.db.WithContext(ctx).
		Where("job_id = ?", jobID).
		Delete(&job.JobBenefit{}).Error
}

// ===========================================
// JOB SKILL OPERATIONS
// ===========================================

// CreateSkill creates a job skill
func (r *jobRepository) CreateSkill(ctx context.Context, skill *job.JobSkill) error {
	return r.db.WithContext(ctx).Create(skill).Error
}

// FindSkillByID finds a skill by ID
func (r *jobRepository) FindSkillByID(ctx context.Context, id int64) (*job.JobSkill, error) {
	var skill job.JobSkill
	err := r.db.WithContext(ctx).First(&skill, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &skill, nil
}

// UpdateSkill updates a job skill
func (r *jobRepository) UpdateSkill(ctx context.Context, skill *job.JobSkill) error {
	return r.db.WithContext(ctx).Save(skill).Error
}

// DeleteSkill deletes a job skill
func (r *jobRepository) DeleteSkill(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&job.JobSkill{}, id).Error
}

// ListSkillsByJob retrieves skills for a job
func (r *jobRepository) ListSkillsByJob(ctx context.Context, jobID int64) ([]job.JobSkill, error) {
	var skills []job.JobSkill
	err := r.db.WithContext(ctx).
		Where("job_id = ?", jobID).
		Order("importance_level ASC, weight DESC").
		Find(&skills).Error
	return skills, err
}

// GetRequiredSkills retrieves required skills
func (r *jobRepository) GetRequiredSkills(ctx context.Context, jobID int64) ([]job.JobSkill, error) {
	var skills []job.JobSkill
	err := r.db.WithContext(ctx).
		Where("job_id = ? AND importance_level = ?", jobID, "required").
		Order("weight DESC").
		Find(&skills).Error
	return skills, err
}

// GetPreferredSkills retrieves preferred skills
func (r *jobRepository) GetPreferredSkills(ctx context.Context, jobID int64) ([]job.JobSkill, error) {
	var skills []job.JobSkill
	err := r.db.WithContext(ctx).
		Where("job_id = ? AND importance_level = ?", jobID, "preferred").
		Order("weight DESC").
		Find(&skills).Error
	return skills, err
}

// BulkCreateSkills creates multiple skills
func (r *jobRepository) BulkCreateSkills(ctx context.Context, skills []job.JobSkill) error {
	return r.db.WithContext(ctx).Create(&skills).Error
}

// BulkDeleteSkills deletes all skills for a job
func (r *jobRepository) BulkDeleteSkills(ctx context.Context, jobID int64) error {
	return r.db.WithContext(ctx).
		Where("job_id = ?", jobID).
		Delete(&job.JobSkill{}).Error
}

// ===========================================
// JOB REQUIREMENT OPERATIONS
// ===========================================

// CreateRequirement creates a job requirement
func (r *jobRepository) CreateRequirement(ctx context.Context, requirement *job.JobRequirement) error {
	return r.db.WithContext(ctx).Create(requirement).Error
}

// FindRequirementByID finds a requirement by ID
func (r *jobRepository) FindRequirementByID(ctx context.Context, id int64) (*job.JobRequirement, error) {
	var requirement job.JobRequirement
	err := r.db.WithContext(ctx).First(&requirement, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &requirement, nil
}

// UpdateRequirement updates a job requirement
func (r *jobRepository) UpdateRequirement(ctx context.Context, requirement *job.JobRequirement) error {
	return r.db.WithContext(ctx).Save(requirement).Error
}

// DeleteRequirement deletes a job requirement
func (r *jobRepository) DeleteRequirement(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&job.JobRequirement{}, id).Error
}

// ListRequirementsByJob retrieves requirements for a job
func (r *jobRepository) ListRequirementsByJob(ctx context.Context, jobID int64) ([]job.JobRequirement, error) {
	var requirements []job.JobRequirement
	err := r.db.WithContext(ctx).
		Where("job_id = ?", jobID).
		Order("is_mandatory DESC, priority ASC").
		Find(&requirements).Error
	return requirements, err
}

// GetMandatoryRequirements retrieves mandatory requirements
func (r *jobRepository) GetMandatoryRequirements(ctx context.Context, jobID int64) ([]job.JobRequirement, error) {
	var requirements []job.JobRequirement
	err := r.db.WithContext(ctx).
		Where("job_id = ? AND is_mandatory = ?", jobID, true).
		Order("priority ASC").
		Find(&requirements).Error
	return requirements, err
}

// BulkCreateRequirements creates multiple requirements
func (r *jobRepository) BulkCreateRequirements(ctx context.Context, requirements []job.JobRequirement) error {
	return r.db.WithContext(ctx).Create(&requirements).Error
}

// BulkDeleteRequirements deletes all requirements for a job
func (r *jobRepository) BulkDeleteRequirements(ctx context.Context, jobID int64) error {
	return r.db.WithContext(ctx).
		Where("job_id = ?", jobID).
		Delete(&job.JobRequirement{}).Error
}

// ===========================================
// ANALYTICS
// ===========================================

// GetTrendingJobs retrieves trending jobs
func (r *jobRepository) GetTrendingJobs(ctx context.Context, limit int) ([]job.Job, error) {
	var jobs []job.Job
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	err := r.db.WithContext(ctx).
		Where("status = ?", "published").
		Where("published_at >= ?", sevenDaysAgo).
		Order("views_count DESC, applications_count DESC").
		Limit(limit).
		Preload("Category").
		Preload("Locations").
		Preload("Benefits").
		Find(&jobs).Error
	return jobs, err
}

// GetPopularCategories retrieves popular categories
func (r *jobRepository) GetPopularCategories(ctx context.Context, limit int) ([]job.CategoryStats, error) {
	var stats []job.CategoryStats
	err := r.db.WithContext(ctx).
		Model(&job.Job{}).
		Select(`
			job_categories.id as category_id,
			job_categories.name as category_name,
			COUNT(jobs.id) as job_count,
			COUNT(CASE WHEN jobs.status = 'published' THEN 1 END) as active_job_count,
			COALESCE(SUM(jobs.views_count), 0) as total_views,
			COALESCE(SUM(jobs.applications_count), 0) as total_applications
		`).
		Joins("INNER JOIN job_categories ON job_categories.id = jobs.category_id").
		Group("job_categories.id, job_categories.name").
		Order("job_count DESC").
		Limit(limit).
		Scan(&stats).Error
	return stats, err
}

// GetJobsByDateRange retrieves jobs within date range
func (r *jobRepository) GetJobsByDateRange(ctx context.Context, startDate, endDate time.Time, filter job.JobFilter) ([]job.Job, error) {
	var jobs []job.Job
	query := r.db.WithContext(ctx).Model(&job.Job{}).
		Where("published_at BETWEEN ? AND ?", startDate, endDate)

	query = r.applyJobFilter(query, filter)

	err := query.
		Preload("Category").
		Preload("Locations").
		Order("published_at DESC").
		Find(&jobs).Error
	return jobs, err
}

// ===========================================
// MASTER DATA PRELOAD
// ===========================================

// PreloadMasterData preloads all master data relations for a job
func (r *jobRepository) PreloadMasterData(ctx context.Context, j *job.Job) error {
	if j == nil {
		return nil
	}

	return r.db.WithContext(ctx).First(j, j.ID).Error
}

// FindByIDWithMasterData retrieves job with all master data preloaded
func (r *jobRepository) FindByIDWithMasterData(ctx context.Context, id int64) (*job.Job, error) {
	var j job.Job

	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Skills.Skill").
		Preload("Benefits").
		Preload("Locations").
		Preload("JobRequirements").
		Where("id = ?", id).
		First(&j).Error

	if err != nil {
		return nil, err
	}

	return &j, nil
}

// ===========================================
// HELPER FUNCTIONS
// ===========================================

// applyJobFilter applies job filter to query
func (r *jobRepository) applyJobFilter(query *gorm.DB, filter job.JobFilter) *gorm.DB {
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.CompanyID > 0 {
		query = query.Where("company_id = ?", filter.CompanyID)
	}
	if filter.CategoryID > 0 {
		query = query.Where("category_id = ?", filter.CategoryID)
	}
	if filter.City != "" {
		query = query.Where("city = ?", filter.City)
	}
	if filter.Province != "" {
		query = query.Where("province = ?", filter.Province)
	}
	if filter.JobLevel != "" {
		query = query.Where("job_level = ?", filter.JobLevel)
	}
	if filter.EmploymentType != "" {
		query = query.Where("employment_type = ?", filter.EmploymentType)
	}
	if filter.RemoteOption != nil {
		query = query.Where("remote_option = ?", *filter.RemoteOption)
	}
	if filter.MinSalary != nil {
		query = query.Where("salary_max >= ? OR salary_max IS NULL", *filter.MinSalary)
	}
	if filter.MaxSalary != nil {
		query = query.Where("salary_min <= ? OR salary_min IS NULL", *filter.MaxSalary)
	}
	if filter.MinExperience != nil {
		query = query.Where("experience_max >= ? OR experience_max IS NULL", *filter.MinExperience)
	}
	if filter.MaxExperience != nil {
		query = query.Where("experience_min <= ? OR experience_min IS NULL", *filter.MaxExperience)
	}
	if filter.EducationLevel != "" {
		query = query.Where("education_level = ?", filter.EducationLevel)
	}
	if filter.IsActive != nil && *filter.IsActive {
		query = query.Where("status = ?", "published").
			Where("(expired_at IS NULL OR expired_at > ?)", time.Now())
	}
	if filter.PublishedAfter != nil {
		query = query.Where("published_at >= ?", *filter.PublishedAfter)
	}
	return query
}

// applySorting applies sorting to query
func (r *jobRepository) applySorting(query *gorm.DB, sortBy string) *gorm.DB {
	switch sortBy {
	case "latest":
		return query.Order("published_at DESC")
	case "salary_asc":
		return query.Order("salary_min ASC NULLS LAST")
	case "salary_desc":
		return query.Order("salary_max DESC NULLS LAST")
	case "views":
		return query.Order("views_count DESC")
	case "applications":
		return query.Order("applications_count DESC")
	default:
		return query.Order("created_at DESC")
	}
}
