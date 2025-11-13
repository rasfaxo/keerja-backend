package postgres

import (
	"context"
	"strings"
	"time"

	"keerja-backend/internal/domain/application"

	"gorm.io/gorm"
)

// applicationRepository implements the application.ApplicationRepository interface
type applicationRepository struct {
	db *gorm.DB
}

// NewApplicationRepository creates a new application repository instance
func NewApplicationRepository(db *gorm.DB) application.ApplicationRepository {
	return &applicationRepository{db: db}
}

// ============================================================================
// JobApplication CRUD Operations
// ============================================================================

// Create creates a new job application
func (r *applicationRepository) Create(ctx context.Context, app *application.JobApplication) error {
	return r.db.WithContext(ctx).Create(app).Error
}

// FindByID finds an application by ID with relationships preloaded
func (r *applicationRepository) FindByID(ctx context.Context, id int64) (*application.JobApplication, error) {
	var app application.JobApplication
	err := r.db.WithContext(ctx).
		Preload("Stages", func(db *gorm.DB) *gorm.DB {
			return db.Order("started_at DESC")
		}).
		Preload("Documents").
		Preload("ApplicationNotes", func(db *gorm.DB) *gorm.DB {
			return db.Order("is_pinned DESC, created_at DESC")
		}).
		Preload("Interviews", func(db *gorm.DB) *gorm.DB {
			return db.Order("scheduled_at DESC")
		}).
		First(&app, id).Error

	if err != nil {
		return nil, err
	}
	return &app, nil
}

// FindByJobAndUser finds an application by job ID and user ID
func (r *applicationRepository) FindByJobAndUser(ctx context.Context, jobID, userID int64) (*application.JobApplication, error) {
	var app application.JobApplication
	err := r.db.WithContext(ctx).
		Where("job_id = ? AND user_id = ?", jobID, userID).
		Preload("Stages", func(db *gorm.DB) *gorm.DB {
			return db.Order("started_at DESC")
		}).
		Preload("Documents").
		Preload("ApplicationNotes", func(db *gorm.DB) *gorm.DB {
			return db.Order("is_pinned DESC, created_at DESC")
		}).
		Preload("Interviews", func(db *gorm.DB) *gorm.DB {
			return db.Order("scheduled_at DESC")
		}).
		First(&app).Error

	if err != nil {
		return nil, err
	}
	return &app, nil
}

// Update updates an existing application
func (r *applicationRepository) Update(ctx context.Context, app *application.JobApplication) error {
	return r.db.WithContext(ctx).Save(app).Error
}

// Delete deletes an application by ID
func (r *applicationRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&application.JobApplication{}, id).Error
}

// ============================================================================
// Application Listing and Filtering
// ============================================================================

// List lists applications with filtering and pagination
func (r *applicationRepository) List(ctx context.Context, filter application.ApplicationFilter, page, limit int) ([]application.JobApplication, int64, error) {
	var apps []application.JobApplication
	var total int64

	query := r.db.WithContext(ctx).Model(&application.JobApplication{})
	query = r.applyApplicationFilter(query, filter)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and sorting
	offset := (page - 1) * limit
	query = r.applySorting(query, filter.SortBy)

	err := query.
		Preload("Stages", func(db *gorm.DB) *gorm.DB {
			return db.Order("started_at DESC").Limit(1)
		}).
		Offset(offset).
		Limit(limit).
		Find(&apps).Error

	return apps, total, err
}

// ListByUser lists applications by user ID
func (r *applicationRepository) ListByUser(ctx context.Context, userID int64, filter application.ApplicationFilter, page, limit int) ([]application.JobApplication, int64, error) {
	filter.UserID = userID
	return r.List(ctx, filter, page, limit)
}

// ListByJob lists applications by job ID
func (r *applicationRepository) ListByJob(ctx context.Context, jobID int64, filter application.ApplicationFilter, page, limit int) ([]application.JobApplication, int64, error) {
	filter.JobID = jobID
	return r.List(ctx, filter, page, limit)
}

// ListByCompany lists applications by company ID
func (r *applicationRepository) ListByCompany(ctx context.Context, companyID int64, filter application.ApplicationFilter, page, limit int) ([]application.JobApplication, int64, error) {
	filter.CompanyID = companyID
	return r.List(ctx, filter, page, limit)
}

// ============================================================================
// Application Status Operations
// ============================================================================

// UpdateStatus updates the status of an application
func (r *applicationRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// BulkUpdateStatus updates status for multiple applications
func (r *applicationRepository) BulkUpdateStatus(ctx context.Context, ids []int64, status string) error {
	return r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Where("id IN ?", ids).
		Update("status", status).Error
}

// GetApplicationsByStatus gets applications by status
func (r *applicationRepository) GetApplicationsByStatus(ctx context.Context, status string, page, limit int) ([]application.JobApplication, int64, error) {
	filter := application.ApplicationFilter{
		Status: status,
	}
	return r.List(ctx, filter, page, limit)
}

// ============================================================================
// Application Tracking
// ============================================================================

// MarkAsViewed marks an application as viewed by employer
func (r *applicationRepository) MarkAsViewed(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Where("id = ?", id).
		Update("viewed_by_employer", true).Error
}

// ToggleBookmark toggles bookmark status of an application
func (r *applicationRepository) ToggleBookmark(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Where("id = ?", id).
		Update("is_bookmarked", gorm.Expr("NOT is_bookmarked")).Error
}

// GetBookmarkedApplications gets bookmarked applications for a company
func (r *applicationRepository) GetBookmarkedApplications(ctx context.Context, companyID int64, page, limit int) ([]application.JobApplication, int64, error) {
	filter := application.ApplicationFilter{
		CompanyID:      companyID,
		BookmarkedOnly: boolPtr(true),
	}
	return r.List(ctx, filter, page, limit)
}

// ============================================================================
// Application Search
// ============================================================================

// SearchApplications performs advanced search on applications
func (r *applicationRepository) SearchApplications(ctx context.Context, filter application.ApplicationSearchFilter, page, limit int) ([]application.JobApplication, int64, error) {
	var apps []application.JobApplication
	var total int64

	query := r.db.WithContext(ctx).Model(&application.JobApplication{})

	// Apply keyword search
	if filter.Keyword != "" {
		keyword := "%" + strings.ToLower(filter.Keyword) + "%"
		query = query.Where("LOWER(notes) LIKE ?", keyword)
	}

	// Filter by job IDs
	if len(filter.JobIDs) > 0 {
		query = query.Where("job_id IN ?", filter.JobIDs)
	}

	// Filter by company IDs
	if len(filter.CompanyIDs) > 0 {
		query = query.Where("company_id IN ?", filter.CompanyIDs)
	}

	// Filter by statuses
	if len(filter.Statuses) > 0 {
		query = query.Where("status IN ?", filter.Statuses)
	}

	// Filter by minimum score
	if filter.MinScore != nil {
		query = query.Where("match_score >= ?", *filter.MinScore)
	}

	// Filter by sources
	if len(filter.Sources) > 0 {
		query = query.Where("source IN ?", filter.Sources)
	}

	// Filter by applied within days
	if filter.AppliedWithin != nil {
		cutoffDate := time.Now().AddDate(0, 0, -*filter.AppliedWithin)
		query = query.Where("applied_at >= ?", cutoffDate)
	}

	// Filter by has documents
	if filter.HasDocuments != nil && *filter.HasDocuments {
		query = query.Joins("INNER JOIN application_documents ON application_documents.application_id = job_applications.id")
	}

	// Filter by has interviews
	if filter.HasInterviews != nil && *filter.HasInterviews {
		query = query.Joins("INNER JOIN interviews ON interviews.application_id = job_applications.id")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit

	err := query.
		Preload("Stages", func(db *gorm.DB) *gorm.DB {
			return db.Order("started_at DESC").Limit(1)
		}).
		Order("applied_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&apps).Error

	return apps, total, err
}

// GetApplicationsWithHighScore gets applications with match score above threshold
func (r *applicationRepository) GetApplicationsWithHighScore(ctx context.Context, minScore float64, limit int) ([]application.JobApplication, error) {
	var apps []application.JobApplication
	err := r.db.WithContext(ctx).
		Where("match_score >= ?", minScore).
		Order("match_score DESC").
		Limit(limit).
		Find(&apps).Error

	return apps, err
}

// ============================================================================
// Application Statistics
// ============================================================================

// GetApplicationStats gets detailed statistics for an application
func (r *applicationRepository) GetApplicationStats(ctx context.Context, applicationID int64) (*application.ApplicationStats, error) {
	var stats application.ApplicationStats
	stats.ApplicationID = applicationID

	// Get application details
	var app application.JobApplication
	if err := r.db.WithContext(ctx).First(&app, applicationID).Error; err != nil {
		return nil, err
	}

	// Count total stages
	r.db.WithContext(ctx).
		Model(&application.JobApplicationStage{}).
		Where("application_id = ?", applicationID).
		Count(&stats.TotalStages)

	// Count completed stages
	r.db.WithContext(ctx).
		Model(&application.JobApplicationStage{}).
		Where("application_id = ? AND completed_at IS NOT NULL", applicationID).
		Count(&stats.CompletedStages)

	// Get current stage
	var currentStage application.JobApplicationStage
	if err := r.db.WithContext(ctx).
		Where("application_id = ? AND completed_at IS NULL", applicationID).
		Order("started_at DESC").
		First(&currentStage).Error; err == nil {
		stats.CurrentStage = currentStage.StageName
	}

	// Count documents
	r.db.WithContext(ctx).
		Model(&application.ApplicationDocument{}).
		Where("application_id = ?", applicationID).
		Count(&stats.TotalDocuments)

	// Count verified documents
	r.db.WithContext(ctx).
		Model(&application.ApplicationDocument{}).
		Where("application_id = ? AND is_verified = ?", applicationID, true).
		Count(&stats.VerifiedDocuments)

	// Count notes
	r.db.WithContext(ctx).
		Model(&application.ApplicationNote{}).
		Where("application_id = ?", applicationID).
		Count(&stats.TotalNotes)

	// Count interviews
	r.db.WithContext(ctx).
		Model(&application.Interview{}).
		Where("application_id = ?", applicationID).
		Count(&stats.TotalInterviews)

	// Count completed interviews
	r.db.WithContext(ctx).
		Model(&application.Interview{}).
		Where("application_id = ? AND status = ?", applicationID, "completed").
		Count(&stats.CompletedInterviews)

	// Calculate average interview score
	var avgScore struct {
		Avg float64
	}
	r.db.WithContext(ctx).
		Model(&application.Interview{}).
		Select("COALESCE(AVG(overall_score), 0) as avg").
		Where("application_id = ? AND overall_score IS NOT NULL", applicationID).
		Scan(&avgScore)
	stats.AverageInterviewScore = avgScore.Avg

	// Calculate days since applied
	stats.DaysSinceApplied = int(time.Since(app.AppliedAt).Hours() / 24)

	// Get last activity
	stats.LastActivity = app.UpdatedAt

	return &stats, nil
}

// GetUserApplicationStats gets application statistics for a user
func (r *applicationRepository) GetUserApplicationStats(ctx context.Context, userID int64) (*application.UserApplicationStats, error) {
	var stats application.UserApplicationStats
	stats.UserID = userID

	// Count total applications
	r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalApplications)

	// Count by status
	statusCounts := []struct {
		Status string
		Count  int64
	}{}
	r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Select("status, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Scan(&statusCounts)

	for _, sc := range statusCounts {
		switch sc.Status {
		case "applied":
			stats.AppliedCount = sc.Count
		case "screening":
			stats.ScreeningCount = sc.Count
		case "shortlisted":
			stats.ShortlistedCount = sc.Count
		case "interview":
			stats.InterviewCount = sc.Count
		case "offered":
			stats.OfferedCount = sc.Count
		case "hired":
			stats.HiredCount = sc.Count
		case "rejected":
			stats.RejectedCount = sc.Count
		case "withdrawn":
			stats.WithdrawnCount = sc.Count
		}
	}

	// Calculate average match score
	var avgScore struct {
		Avg float64
	}
	r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Select("COALESCE(AVG(match_score), 0) as avg").
		Where("user_id = ?", userID).
		Scan(&avgScore)
	stats.AverageMatchScore = avgScore.Avg

	// Calculate success rate (hired / total)
	if stats.TotalApplications > 0 {
		stats.SuccessRate = float64(stats.HiredCount) / float64(stats.TotalApplications) * 100
	}

	// Calculate average response time (time from applied to first status change)
	var avgResponseTime struct {
		Avg float64
	}
	r.db.WithContext(ctx).
		Table("job_applications").
		Select("COALESCE(AVG(EXTRACT(EPOCH FROM (updated_at - applied_at)) / 86400), 0) as avg").
		Where("user_id = ? AND status != ?", userID, "applied").
		Scan(&avgResponseTime)
	stats.AverageResponseTime = avgResponseTime.Avg

	return &stats, nil
}

// GetJobApplicationStats gets application statistics for a job
func (r *applicationRepository) GetJobApplicationStats(ctx context.Context, jobID int64) (*application.JobApplicationStats, error) {
	var stats application.JobApplicationStats
	stats.JobID = jobID

	// Count total applications
	r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Where("job_id = ?", jobID).
		Count(&stats.TotalApplications)

	// Count by status
	statusCounts := []struct {
		Status string
		Count  int64
	}{}
	r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Select("status, COUNT(*) as count").
		Where("job_id = ?", jobID).
		Group("status").
		Scan(&statusCounts)

	for _, sc := range statusCounts {
		switch sc.Status {
		case "applied":
			stats.AppliedCount = sc.Count
		case "screening":
			stats.ScreeningCount = sc.Count
		case "shortlisted":
			stats.ShortlistedCount = sc.Count
		case "interview":
			stats.InterviewCount = sc.Count
		case "offered":
			stats.OfferedCount = sc.Count
		case "hired":
			stats.HiredCount = sc.Count
		case "rejected":
			stats.RejectedCount = sc.Count
		case "withdrawn":
			stats.WithdrawnCount = sc.Count
		}
	}

	// Calculate average match score
	var avgScore struct {
		Avg float64
	}
	r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Select("COALESCE(AVG(match_score), 0) as avg").
		Where("job_id = ?", jobID).
		Scan(&avgScore)
	stats.AverageMatchScore = avgScore.Avg

	// Calculate conversion rate (hired / total)
	if stats.TotalApplications > 0 {
		stats.ConversionRate = float64(stats.HiredCount) / float64(stats.TotalApplications) * 100
	}

	// Calculate average time to hire
	var avgTimeToHire struct {
		Avg float64
	}
	r.db.WithContext(ctx).
		Table("job_applications").
		Select("COALESCE(AVG(EXTRACT(EPOCH FROM (updated_at - applied_at)) / 86400), 0) as avg").
		Where("job_id = ? AND status = ?", jobID, "hired").
		Scan(&avgTimeToHire)
	stats.AverageTimeToHire = avgTimeToHire.Avg

	// Get top sources
	r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Select("source, COUNT(*) as count").
		Where("job_id = ?", jobID).
		Group("source").
		Order("count DESC").
		Limit(5).
		Scan(&stats.TopSources)

	return &stats, nil
}

// GetCompanyApplicationStats gets application statistics for a company
func (r *applicationRepository) GetCompanyApplicationStats(ctx context.Context, companyID int64) (*application.CompanyApplicationStats, error) {
	var stats application.CompanyApplicationStats
	stats.CompanyID = companyID

	// Count total applications
	r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Where("company_id = ?", companyID).
		Count(&stats.TotalApplications)

	// Count viewed applications
	r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Where("company_id = ? AND viewed_by_employer = ?", companyID, true).
		Count(&stats.ViewedApplications)

	// Count bookmarked applications
	r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Where("company_id = ? AND is_bookmarked = ?", companyID, true).
		Count(&stats.BookmarkedApplications)

	// Count total hires
	r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Where("company_id = ? AND status = ?", companyID, "hired").
		Count(&stats.TotalHires)

	// Calculate average match score
	var avgScore struct {
		Avg float64
	}
	r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Select("COALESCE(AVG(match_score), 0) as avg").
		Where("company_id = ?", companyID).
		Scan(&avgScore)
	stats.AverageMatchScore = avgScore.Avg

	// Calculate average time to hire
	var avgTimeToHire struct {
		Avg float64
	}
	r.db.WithContext(ctx).
		Table("job_applications").
		Select("COALESCE(AVG(EXTRACT(EPOCH FROM (updated_at - applied_at)) / 86400), 0) as avg").
		Where("company_id = ? AND status = ?", companyID, "hired").
		Scan(&avgTimeToHire)
	stats.AverageTimeToHire = avgTimeToHire.Avg

	// Calculate conversion rate
	if stats.TotalApplications > 0 {
		stats.ConversionRate = float64(stats.TotalHires) / float64(stats.TotalApplications) * 100
	}

	// Get top performing jobs
	r.db.WithContext(ctx).
		Table("job_applications").
		Select("job_id, COUNT(*) as application_count, SUM(CASE WHEN status = 'hired' THEN 1 ELSE 0 END) as hired_count, COALESCE(AVG(match_score), 0) as average_match_score").
		Where("company_id = ?", companyID).
		Group("job_id").
		Order("hired_count DESC").
		Limit(10).
		Scan(&stats.TopPerformingJobs)

	// Calculate conversion rates for top performing jobs
	for i := range stats.TopPerformingJobs {
		if stats.TopPerformingJobs[i].ApplicationCount > 0 {
			stats.TopPerformingJobs[i].ConversionRate = float64(stats.TopPerformingJobs[i].HiredCount) / float64(stats.TopPerformingJobs[i].ApplicationCount) * 100
		}
	}

	// Get applications by month (last 12 months)
	r.db.WithContext(ctx).
		Table("job_applications").
		Select("TO_CHAR(applied_at, 'YYYY-MM') as month, COUNT(*) as count").
		Where("company_id = ? AND applied_at >= ?", companyID, time.Now().AddDate(0, -12, 0)).
		Group("month").
		Order("month DESC").
		Scan(&stats.ApplicationsByMonth)

	return &stats, nil
}

// ============================================================================
// JobApplicationStage Operations
// ============================================================================

// CreateStage creates a new application stage
func (r *applicationRepository) CreateStage(ctx context.Context, stage *application.JobApplicationStage) error {
	return r.db.WithContext(ctx).Create(stage).Error
}

// FindStageByID finds a stage by ID
func (r *applicationRepository) FindStageByID(ctx context.Context, id int64) (*application.JobApplicationStage, error) {
	var stage application.JobApplicationStage
	err := r.db.WithContext(ctx).
		Preload("Application").
		Preload("StageNotes").
		Preload("Interviews").
		First(&stage, id).Error

	if err != nil {
		return nil, err
	}
	return &stage, nil
}

// UpdateStage updates an existing stage
func (r *applicationRepository) UpdateStage(ctx context.Context, stage *application.JobApplicationStage) error {
	return r.db.WithContext(ctx).Save(stage).Error
}

// CompleteStage marks a stage as completed
func (r *applicationRepository) CompleteStage(ctx context.Context, id int64, notes string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&application.JobApplicationStage{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"completed_at": now,
			"notes":        notes,
		}).Error
}

// ListStagesByApplication lists all stages for an application
func (r *applicationRepository) ListStagesByApplication(ctx context.Context, applicationID int64) ([]application.JobApplicationStage, error) {
	var stages []application.JobApplicationStage
	err := r.db.WithContext(ctx).
		Where("application_id = ?", applicationID).
		Order("started_at ASC").
		Find(&stages).Error

	return stages, err
}

// GetCurrentStage gets the current (incomplete) stage for an application
func (r *applicationRepository) GetCurrentStage(ctx context.Context, applicationID int64) (*application.JobApplicationStage, error) {
	var stage application.JobApplicationStage
	err := r.db.WithContext(ctx).
		Where("application_id = ? AND completed_at IS NULL", applicationID).
		Order("started_at DESC").
		First(&stage).Error

	if err != nil {
		return nil, err
	}
	return &stage, nil
}

// GetStageHistory gets the stage history for an application
func (r *applicationRepository) GetStageHistory(ctx context.Context, applicationID int64) ([]application.JobApplicationStage, error) {
	var stages []application.JobApplicationStage
	err := r.db.WithContext(ctx).
		Where("application_id = ? AND completed_at IS NOT NULL", applicationID).
		Order("completed_at DESC").
		Find(&stages).Error

	return stages, err
}

// ============================================================================
// ApplicationDocument Operations
// ============================================================================

// CreateDocument creates a new application document
func (r *applicationRepository) CreateDocument(ctx context.Context, document *application.ApplicationDocument) error {
	return r.db.WithContext(ctx).Create(document).Error
}

// FindDocumentByID finds a document by ID
func (r *applicationRepository) FindDocumentByID(ctx context.Context, id int64) (*application.ApplicationDocument, error) {
	var document application.ApplicationDocument
	err := r.db.WithContext(ctx).
		Preload("Application").
		First(&document, id).Error

	if err != nil {
		return nil, err
	}
	return &document, nil
}

// UpdateDocument updates an existing document
func (r *applicationRepository) UpdateDocument(ctx context.Context, document *application.ApplicationDocument) error {
	return r.db.WithContext(ctx).Save(document).Error
}

// DeleteDocument deletes a document by ID
func (r *applicationRepository) DeleteDocument(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&application.ApplicationDocument{}, id).Error
}

// ListDocumentsByApplication lists all documents for an application
func (r *applicationRepository) ListDocumentsByApplication(ctx context.Context, applicationID int64) ([]application.ApplicationDocument, error) {
	var documents []application.ApplicationDocument
	err := r.db.WithContext(ctx).
		Where("application_id = ?", applicationID).
		Order("uploaded_at DESC").
		Find(&documents).Error

	return documents, err
}

// ListDocumentsByUser lists documents for a user filtered by type
func (r *applicationRepository) ListDocumentsByUser(ctx context.Context, userID int64, docType string) ([]application.ApplicationDocument, error) {
	var documents []application.ApplicationDocument
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if docType != "" {
		query = query.Where("document_type = ?", docType)
	}

	err := query.Order("uploaded_at DESC").Find(&documents).Error
	return documents, err
}

// GetDocumentsByType gets documents by application ID and type
func (r *applicationRepository) GetDocumentsByType(ctx context.Context, applicationID int64, docType string) ([]application.ApplicationDocument, error) {
	var documents []application.ApplicationDocument
	err := r.db.WithContext(ctx).
		Where("application_id = ? AND document_type = ?", applicationID, docType).
		Order("uploaded_at DESC").
		Find(&documents).Error

	return documents, err
}

// VerifyDocument marks a document as verified
func (r *applicationRepository) VerifyDocument(ctx context.Context, id int64, verifiedBy int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&application.ApplicationDocument{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_verified": true,
			"verified_by": verifiedBy,
			"verified_at": now,
		}).Error
}

// GetUnverifiedDocuments gets unverified documents with pagination
func (r *applicationRepository) GetUnverifiedDocuments(ctx context.Context, page, limit int) ([]application.ApplicationDocument, int64, error) {
	var documents []application.ApplicationDocument
	var total int64

	query := r.db.WithContext(ctx).
		Model(&application.ApplicationDocument{}).
		Where("is_verified = ?", false)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit

	err := query.
		Order("uploaded_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&documents).Error

	return documents, total, err
}

// ============================================================================
// ApplicationNote Operations
// ============================================================================

// CreateNote creates a new application note
func (r *applicationRepository) CreateNote(ctx context.Context, note *application.ApplicationNote) error {
	return r.db.WithContext(ctx).Create(note).Error
}

// FindNoteByID finds a note by ID
func (r *applicationRepository) FindNoteByID(ctx context.Context, id int64) (*application.ApplicationNote, error) {
	var note application.ApplicationNote
	err := r.db.WithContext(ctx).
		Preload("Application").
		Preload("Stage").
		First(&note, id).Error

	if err != nil {
		return nil, err
	}
	return &note, nil
}

// UpdateNote updates an existing note
func (r *applicationRepository) UpdateNote(ctx context.Context, note *application.ApplicationNote) error {
	return r.db.WithContext(ctx).Save(note).Error
}

// DeleteNote deletes a note by ID
func (r *applicationRepository) DeleteNote(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&application.ApplicationNote{}, id).Error
}

// ListNotesByApplication lists all notes for an application
func (r *applicationRepository) ListNotesByApplication(ctx context.Context, applicationID int64) ([]application.ApplicationNote, error) {
	var notes []application.ApplicationNote
	err := r.db.WithContext(ctx).
		Where("application_id = ?", applicationID).
		Order("is_pinned DESC, created_at DESC").
		Find(&notes).Error

	return notes, err
}

// ListNotesByStage lists all notes for a stage
func (r *applicationRepository) ListNotesByStage(ctx context.Context, stageID int64) ([]application.ApplicationNote, error) {
	var notes []application.ApplicationNote
	err := r.db.WithContext(ctx).
		Where("stage_id = ?", stageID).
		Order("created_at DESC").
		Find(&notes).Error

	return notes, err
}

// ListNotesByAuthor lists notes by author with pagination
func (r *applicationRepository) ListNotesByAuthor(ctx context.Context, authorID int64, page, limit int) ([]application.ApplicationNote, int64, error) {
	var notes []application.ApplicationNote
	var total int64

	query := r.db.WithContext(ctx).
		Model(&application.ApplicationNote{}).
		Where("author_id = ?", authorID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notes).Error

	return notes, total, err
}

// GetPinnedNotes gets all pinned notes for an application
func (r *applicationRepository) GetPinnedNotes(ctx context.Context, applicationID int64) ([]application.ApplicationNote, error) {
	var notes []application.ApplicationNote
	err := r.db.WithContext(ctx).
		Where("application_id = ? AND is_pinned = ?", applicationID, true).
		Order("created_at DESC").
		Find(&notes).Error

	return notes, err
}

// PinNote marks a note as pinned
func (r *applicationRepository) PinNote(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&application.ApplicationNote{}).
		Where("id = ?", id).
		Update("is_pinned", true).Error
}

// UnpinNote removes pinned status from a note
func (r *applicationRepository) UnpinNote(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&application.ApplicationNote{}).
		Where("id = ?", id).
		Update("is_pinned", false).Error
}

// ============================================================================
// Interview Operations
// ============================================================================

// CreateInterview creates a new interview
func (r *applicationRepository) CreateInterview(ctx context.Context, interview *application.Interview) error {
	return r.db.WithContext(ctx).Create(interview).Error
}

// FindInterviewByID finds an interview by ID
func (r *applicationRepository) FindInterviewByID(ctx context.Context, id int64) (*application.Interview, error) {
	var interview application.Interview
	err := r.db.WithContext(ctx).
		Preload("Application").
		Preload("Stage").
		First(&interview, id).Error

	if err != nil {
		return nil, err
	}
	return &interview, nil
}

// UpdateInterview updates an existing interview
func (r *applicationRepository) UpdateInterview(ctx context.Context, interview *application.Interview) error {
	return r.db.WithContext(ctx).Save(interview).Error
}

// DeleteInterview deletes an interview by ID
func (r *applicationRepository) DeleteInterview(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&application.Interview{}, id).Error
}

// ListInterviewsByApplication lists all interviews for an application
func (r *applicationRepository) ListInterviewsByApplication(ctx context.Context, applicationID int64) ([]application.Interview, error) {
	var interviews []application.Interview
	err := r.db.WithContext(ctx).
		Where("application_id = ?", applicationID).
		Order("scheduled_at DESC").
		Find(&interviews).Error

	return interviews, err
}

// ListInterviewsByInterviewer lists interviews by interviewer with filtering
func (r *applicationRepository) ListInterviewsByInterviewer(ctx context.Context, interviewerID int64, filter application.InterviewFilter) ([]application.Interview, error) {
	var interviews []application.Interview
	query := r.db.WithContext(ctx).Where("interviewer_id = ?", interviewerID)

	// Apply filters
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.InterviewType != "" {
		query = query.Where("interview_type = ?", filter.InterviewType)
	}

	if filter.ScheduledFrom != nil {
		query = query.Where("scheduled_at >= ?", filter.ScheduledFrom)
	}

	if filter.ScheduledTo != nil {
		query = query.Where("scheduled_at <= ?", filter.ScheduledTo)
	}

	if filter.CompletedOnly != nil && *filter.CompletedOnly {
		query = query.Where("status = ?", "completed")
	}

	err := query.Order("scheduled_at DESC").Find(&interviews).Error
	return interviews, err
}

// GetUpcomingInterviews gets upcoming interviews from a specific date
func (r *applicationRepository) GetUpcomingInterviews(ctx context.Context, date time.Time, limit int) ([]application.Interview, error) {
	var interviews []application.Interview
	err := r.db.WithContext(ctx).
		Where("scheduled_at >= ? AND status = ?", date, "scheduled").
		Order("scheduled_at ASC").
		Limit(limit).
		Find(&interviews).Error

	return interviews, err
}

// GetInterviewsByDateRange gets interviews within date range
func (r *applicationRepository) GetInterviewsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]application.Interview, error) {
	var interviews []application.Interview
	err := r.db.WithContext(ctx).
		Where("scheduled_at BETWEEN ? AND ?", startDate, endDate).
		Order("scheduled_at ASC").
		Find(&interviews).Error

	return interviews, err
}

// UpdateInterviewStatus updates interview status
func (r *applicationRepository) UpdateInterviewStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).
		Model(&application.Interview{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// CompleteInterview marks interview as completed with scores and feedback
func (r *applicationRepository) CompleteInterview(ctx context.Context, id int64, scores application.InterviewScores, feedback string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":           "completed",
		"ended_at":         now,
		"feedback_summary": feedback,
	}

	if scores.OverallScore != nil {
		updates["overall_score"] = *scores.OverallScore
	}
	if scores.TechnicalScore != nil {
		updates["technical_score"] = *scores.TechnicalScore
	}
	if scores.CommunicationScore != nil {
		updates["communication_score"] = *scores.CommunicationScore
	}
	if scores.PersonalityScore != nil {
		updates["personality_score"] = *scores.PersonalityScore
	}

	return r.db.WithContext(ctx).
		Model(&application.Interview{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// RescheduleInterview reschedules an interview to a new time
func (r *applicationRepository) RescheduleInterview(ctx context.Context, id int64, newSchedule time.Time) error {
	return r.db.WithContext(ctx).
		Model(&application.Interview{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"scheduled_at": newSchedule,
			"status":       "rescheduled",
		}).Error
}

// CancelInterview cancels an interview
func (r *applicationRepository) CancelInterview(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&application.Interview{}).
		Where("id = ?", id).
		Update("status", "cancelled").Error
}

// ============================================================================
// Analytics and Reporting
// ============================================================================

// GetApplicationTrends gets application trends over time
func (r *applicationRepository) GetApplicationTrends(ctx context.Context, startDate, endDate time.Time) ([]application.ApplicationTrend, error) {
	var trends []application.ApplicationTrend

	err := r.db.WithContext(ctx).
		Table("job_applications").
		Select(`
			DATE(applied_at) as date,
			COUNT(*) as total_applications,
			SUM(CASE WHEN status = 'hired' THEN 1 ELSE 0 END) as hired_count,
			SUM(CASE WHEN status = 'rejected' THEN 1 ELSE 0 END) as rejected_count,
			COALESCE(AVG(match_score), 0) as average_match_score
		`).
		Where("applied_at BETWEEN ? AND ?", startDate, endDate).
		Group("DATE(applied_at)").
		Order("date DESC").
		Scan(&trends).Error

	return trends, err
}

// GetConversionFunnel gets conversion funnel metrics for a job
func (r *applicationRepository) GetConversionFunnel(ctx context.Context, jobID int64) (*application.ConversionFunnel, error) {
	var funnel application.ConversionFunnel
	funnel.JobID = jobID

	// Count applications by stage
	statusCounts := []struct {
		Status string
		Count  int64
	}{}
	r.db.WithContext(ctx).
		Model(&application.JobApplication{}).
		Select("status, COUNT(*) as count").
		Where("job_id = ?", jobID).
		Group("status").
		Scan(&statusCounts)

	for _, sc := range statusCounts {
		switch sc.Status {
		case "applied":
			funnel.AppliedCount = sc.Count
		case "screening":
			funnel.ScreeningCount = sc.Count
		case "shortlisted":
			funnel.ShortlistedCount = sc.Count
		case "interview":
			funnel.InterviewCount = sc.Count
		case "offered":
			funnel.OfferedCount = sc.Count
		case "hired":
			funnel.HiredCount = sc.Count
		}
	}

	// Calculate conversion rates
	funnel.ConversionRates = make(map[string]float64)

	if funnel.AppliedCount > 0 {
		funnel.ConversionRates["screening"] = float64(funnel.ScreeningCount) / float64(funnel.AppliedCount) * 100
		funnel.ConversionRates["shortlisted"] = float64(funnel.ShortlistedCount) / float64(funnel.AppliedCount) * 100
		funnel.ConversionRates["interview"] = float64(funnel.InterviewCount) / float64(funnel.AppliedCount) * 100
		funnel.ConversionRates["offered"] = float64(funnel.OfferedCount) / float64(funnel.AppliedCount) * 100
		funnel.ConversionRates["hired"] = float64(funnel.HiredCount) / float64(funnel.AppliedCount) * 100
	}

	return &funnel, nil
}

// GetAverageTimePerStage gets average time spent in each stage
func (r *applicationRepository) GetAverageTimePerStage(ctx context.Context, companyID int64) ([]application.StageTimeStats, error) {
	var stats []application.StageTimeStats

	err := r.db.WithContext(ctx).
		Table("job_application_stages").
		Select(`
			stage_name,
			COALESCE(AVG(EXTRACT(EPOCH FROM (completed_at - started_at)) / 86400), 0) as average_days,
			COALESCE(MIN(EXTRACT(EPOCH FROM (completed_at - started_at)) / 86400), 0) as min_days,
			COALESCE(MAX(EXTRACT(EPOCH FROM (completed_at - started_at)) / 86400), 0) as max_days,
			COUNT(*) as count
		`).
		Joins("INNER JOIN job_applications ON job_applications.id = job_application_stages.application_id").
		Where("job_applications.company_id = ? AND job_application_stages.completed_at IS NOT NULL", companyID).
		Group("stage_name").
		Order("stage_name ASC").
		Scan(&stats).Error

	return stats, err
}

// GetTopApplicants gets top applicants for a job based on match score
func (r *applicationRepository) GetTopApplicants(ctx context.Context, jobID int64, limit int) ([]application.JobApplication, error) {
	var apps []application.JobApplication
	err := r.db.WithContext(ctx).
		Where("job_id = ?", jobID).
		Order("match_score DESC").
		Limit(limit).
		Find(&apps).Error

	return apps, err
}

// GetApplicationSourceStats gets statistics by application source
func (r *applicationRepository) GetApplicationSourceStats(ctx context.Context, companyID int64) ([]application.SourceStats, error) {
	var stats []application.SourceStats

	err := r.db.WithContext(ctx).
		Table("job_applications").
		Select(`
			source,
			COUNT(*) as count,
			SUM(CASE WHEN status = 'hired' THEN 1 ELSE 0 END) as hired_count,
			CASE WHEN COUNT(*) > 0 THEN (SUM(CASE WHEN status = 'hired' THEN 1 ELSE 0 END)::float / COUNT(*) * 100) ELSE 0 END as conversion_rate,
			COALESCE(AVG(match_score), 0) as average_match_score
		`).
		Where("company_id = ?", companyID).
		Group("source").
		Order("count DESC").
		Scan(&stats).Error

	return stats, err
}

// ============================================================================
// Bulk Operations
// ============================================================================

// BulkCreateApplications creates multiple applications at once
func (r *applicationRepository) BulkCreateApplications(ctx context.Context, applications []application.JobApplication) error {
	return r.db.WithContext(ctx).Create(&applications).Error
}

// BulkDeleteApplications deletes multiple applications by IDs
func (r *applicationRepository) BulkDeleteApplications(ctx context.Context, ids []int64) error {
	return r.db.WithContext(ctx).Delete(&application.JobApplication{}, ids).Error
}

// ============================================================================
// Helper Functions
// ============================================================================

// applyApplicationFilter applies filter criteria to query
func (r *applicationRepository) applyApplicationFilter(query *gorm.DB, filter application.ApplicationFilter) *gorm.DB {
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.JobID > 0 {
		query = query.Where("job_id = ?", filter.JobID)
	}

	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if filter.CompanyID > 0 {
		query = query.Where("company_id = ?", filter.CompanyID)
	}

	if filter.MinScore != nil {
		query = query.Where("match_score >= ?", *filter.MinScore)
	}

	if filter.MaxScore != nil {
		query = query.Where("match_score <= ?", *filter.MaxScore)
	}

	if filter.ViewedOnly != nil && *filter.ViewedOnly {
		query = query.Where("viewed_by_employer = ?", true)
	}

	if filter.BookmarkedOnly != nil && *filter.BookmarkedOnly {
		query = query.Where("is_bookmarked = ?", true)
	}

	if filter.Source != "" {
		query = query.Where("source = ?", filter.Source)
	}

	if filter.AppliedAfter != nil {
		query = query.Where("applied_at >= ?", filter.AppliedAfter)
	}

	if filter.AppliedBefore != nil {
		query = query.Where("applied_at <= ?", filter.AppliedBefore)
	}

	return query
}

// applySorting applies sorting to query
func (r *applicationRepository) applySorting(query *gorm.DB, sortBy string) *gorm.DB {
	switch sortBy {
	case "score_desc":
		return query.Order("match_score DESC")
	case "score_asc":
		return query.Order("match_score ASC")
	case "latest":
		return query.Order("applied_at DESC")
	default:
		return query.Order("applied_at DESC")
	}
}

// boolPtr returns a pointer to a boolean value
func boolPtr(b bool) *bool {
	return &b
}
