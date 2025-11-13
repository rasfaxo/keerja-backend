package application

import (
	"context"
	"time"
)

// ApplicationService defines the interface for application business logic
type ApplicationService interface {
	// Application submission and management (Job Seeker)
	ApplyForJob(ctx context.Context, req *ApplyJobRequest) (*JobApplication, error)
	WithdrawApplication(ctx context.Context, applicationID, userID int64) error
	GetMyApplications(ctx context.Context, userID int64, filter ApplicationFilter, page, limit int) (*ApplicationListResponse, error)
	GetApplicationDetail(ctx context.Context, applicationID, userID int64) (*ApplicationDetailResponse, error)
	GetMyApplicationStats(ctx context.Context, userID int64) (*UserApplicationStats, error)

	// Application review and management (Employer)
	GetJobApplications(ctx context.Context, jobID int64, filter ApplicationFilter, page, limit int) (*ApplicationListResponse, error)
	GetCompanyApplications(ctx context.Context, companyID int64, filter ApplicationFilter, page, limit int) (*ApplicationListResponse, error)
	GetApplicationForReview(ctx context.Context, applicationID, employerUserID int64) (*ApplicationDetailResponse, error)
	MarkAsViewed(ctx context.Context, applicationID, employerUserID int64) error
	ToggleBookmark(ctx context.Context, applicationID, employerUserID int64) error
	GetBookmarkedApplications(ctx context.Context, companyID int64, page, limit int) (*ApplicationListResponse, error)

	// Application status workflow (Employer)
	MoveToScreening(ctx context.Context, applicationID, handledBy int64, notes string) error
	MoveToShortlist(ctx context.Context, applicationID, handledBy int64, notes string) error
	MoveToInterview(ctx context.Context, applicationID, handledBy int64, notes string) error
	MakeOffer(ctx context.Context, applicationID, handledBy int64, notes string) error
	MarkAsHired(ctx context.Context, applicationID, handledBy int64, notes string) error
	RejectApplication(ctx context.Context, applicationID, handledBy int64, reason string) error
	BulkUpdateStatus(ctx context.Context, applicationIDs []int64, status string, handledBy int64) error

	// Stage management
	GetApplicationStages(ctx context.Context, applicationID int64) ([]JobApplicationStage, error)
	GetCurrentStage(ctx context.Context, applicationID int64) (*JobApplicationStage, error)
	GetStageHistory(ctx context.Context, applicationID int64) ([]JobApplicationStage, error)
	CompleteStage(ctx context.Context, stageID, handledBy int64, notes string) error

	// Document management
	UploadApplicationDocument(ctx context.Context, req *UploadDocumentRequest) (*ApplicationDocument, error)
	UpdateDocument(ctx context.Context, documentID int64, req *UpdateDocumentRequest) (*ApplicationDocument, error)
	DeleteDocument(ctx context.Context, documentID, userID int64) error
	GetApplicationDocuments(ctx context.Context, applicationID int64) ([]ApplicationDocument, error)
	GetDocumentsByType(ctx context.Context, applicationID int64, docType string) ([]ApplicationDocument, error)
	VerifyDocument(ctx context.Context, documentID, verifiedBy int64, notes string) error
	GetUnverifiedDocuments(ctx context.Context, page, limit int) ([]ApplicationDocument, int64, error)

	// Notes management (Employer)
	AddNote(ctx context.Context, req *AddNoteRequest) (*ApplicationNote, error)
	UpdateNote(ctx context.Context, noteID int64, req *UpdateNoteRequest) (*ApplicationNote, error)
	DeleteNote(ctx context.Context, noteID, authorID int64) error
	GetApplicationNotes(ctx context.Context, applicationID int64, visibility string) ([]ApplicationNote, error)
	GetStageNotes(ctx context.Context, stageID int64) ([]ApplicationNote, error)
	PinNote(ctx context.Context, noteID, employerUserID int64) error
	UnpinNote(ctx context.Context, noteID, employerUserID int64) error
	GetPinnedNotes(ctx context.Context, applicationID int64) ([]ApplicationNote, error)

	// Interview scheduling and management
	ScheduleInterview(ctx context.Context, req *ScheduleInterviewRequest) (*Interview, error)
	RescheduleInterview(ctx context.Context, interviewID int64, req *RescheduleInterviewRequest) (*Interview, error)
	CancelInterview(ctx context.Context, interviewID int64, cancelledBy int64, reason string) error
	CompleteInterview(ctx context.Context, interviewID int64, req *CompleteInterviewRequest) (*Interview, error)
	MarkInterviewNoShow(ctx context.Context, interviewID int64, markedBy int64) error
	GetApplicationInterviews(ctx context.Context, applicationID int64) ([]Interview, error)
	GetInterviewDetail(ctx context.Context, interviewID int64) (*Interview, error)
	GetUpcomingInterviews(ctx context.Context, employerUserID int64, days int) ([]Interview, error)
	GetInterviewsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]Interview, error)
	SendInterviewReminder(ctx context.Context, interviewID int64) error

	// Search and filtering
	SearchApplications(ctx context.Context, filter ApplicationSearchFilter, page, limit int) (*ApplicationListResponse, error)
	GetHighScoreApplications(ctx context.Context, companyID int64, minScore float64, limit int) ([]JobApplication, error)
	GetRecentApplications(ctx context.Context, companyID int64, hours int, limit int) ([]JobApplication, error)

	// Analytics and reporting
	GetApplicationAnalytics(ctx context.Context, applicationID int64) (*ApplicationAnalytics, error)
	GetJobApplicationAnalytics(ctx context.Context, jobID int64, startDate, endDate time.Time) (*JobApplicationAnalytics, error)
	GetCompanyApplicationAnalytics(ctx context.Context, companyID int64, startDate, endDate time.Time) (*CompanyApplicationAnalytics, error)
	GetConversionFunnel(ctx context.Context, jobID int64) (*ConversionFunnel, error)
	GetApplicationTrends(ctx context.Context, companyID int64, startDate, endDate time.Time) ([]ApplicationTrend, error)
	GetAverageTimePerStage(ctx context.Context, companyID int64) ([]StageTimeStats, error)
	GetTopApplicants(ctx context.Context, jobID int64, limit int) ([]JobApplication, error)
	GetApplicationSourceAnalytics(ctx context.Context, companyID int64) ([]SourceStats, error)

	// Notifications
	NotifyApplicationReceived(ctx context.Context, applicationID int64) error
	NotifyStatusUpdate(ctx context.Context, applicationID int64, newStatus string) error
	NotifyInterviewScheduled(ctx context.Context, interviewID int64) error
	NotifyInterviewReminder(ctx context.Context, interviewID int64) error

	// Validation and permissions
	ValidateApplication(ctx context.Context, application *JobApplication) error
	CheckApplicationOwnership(ctx context.Context, applicationID, userID int64) error
	CheckEmployerAccess(ctx context.Context, applicationID, employerUserID int64) error
	CanApplyForJob(ctx context.Context, jobID, userID int64) error

	// Bulk operations
	BulkRejectApplications(ctx context.Context, applicationIDs []int64, rejectedBy int64, reason string) error
	BulkMoveToStage(ctx context.Context, applicationIDs []int64, stage string, handledBy int64) error
	ExportApplications(ctx context.Context, companyID int64, filter ApplicationFilter) ([]byte, error)
}

// ===== Request DTOs =====

// ApplyJobRequest represents request to apply for a job
type ApplyJobRequest struct {
	JobID       int64                   `json:"job_id" validate:"required"`
	UserID      int64                   `json:"user_id" validate:"required"`
	ResumeURL   string                  `json:"resume_url,omitempty"`
	CoverLetter string                  `json:"cover_letter,omitempty"`
	Source      string                  `json:"source,omitempty"`
	Documents   []UploadDocumentRequest `json:"documents,omitempty"`
}

// UploadDocumentRequest represents request to upload application document
type UploadDocumentRequest struct {
	ApplicationID int64  `json:"application_id" validate:"required"`
	UserID        int64  `json:"user_id" validate:"required"`
	DocumentType  string `json:"document_type" validate:"required,oneof='cv' 'cover_letter' 'portfolio' 'certificate' 'transcript' 'other'"`
	FileName      string `json:"file_name,omitempty"`
	FileURL       string `json:"file_url" validate:"required"`
	FileType      string `json:"file_type,omitempty"`
	FileSize      int64  `json:"file_size,omitempty"`
	Notes         string `json:"notes,omitempty"`
}

// UpdateDocumentRequest represents request to update document
type UpdateDocumentRequest struct {
	FileName string `json:"file_name,omitempty"`
	FileURL  string `json:"file_url,omitempty"`
	Notes    string `json:"notes,omitempty"`
}

// AddNoteRequest represents request to add application note
type AddNoteRequest struct {
	ApplicationID int64  `json:"application_id" validate:"required"`
	StageID       *int64 `json:"stage_id,omitempty"`
	AuthorID      int64  `json:"author_id" validate:"required"`
	NoteType      string `json:"note_type" validate:"omitempty,oneof='evaluation' 'feedback' 'reminder' 'internal'"`
	NoteText      string `json:"note_text" validate:"required"`
	Visibility    string `json:"visibility" validate:"omitempty,oneof='internal' 'public'"`
	Sentiment     string `json:"sentiment" validate:"omitempty,oneof='positive' 'neutral' 'negative'"`
	IsPinned      bool   `json:"is_pinned"`
}

// UpdateNoteRequest represents request to update note
type UpdateNoteRequest struct {
	NoteText   string `json:"note_text,omitempty"`
	Visibility string `json:"visibility,omitempty"`
	Sentiment  string `json:"sentiment,omitempty"`
	IsPinned   *bool  `json:"is_pinned,omitempty"`
}

// ScheduleInterviewRequest represents request to schedule interview
type ScheduleInterviewRequest struct {
	ApplicationID int64     `json:"application_id" validate:"required"`
	StageID       *int64    `json:"stage_id,omitempty"`
	InterviewerID *int64    `json:"interviewer_id,omitempty"`
	ScheduledAt   time.Time `json:"scheduled_at" validate:"required"`
	InterviewType string    `json:"interview_type" validate:"omitempty,oneof='online' 'onsite' 'hybrid'"`
	MeetingLink   string    `json:"meeting_link,omitempty"`
	Location      string    `json:"location,omitempty"`
}

// RescheduleInterviewRequest represents request to reschedule interview
type RescheduleInterviewRequest struct {
	ScheduledAt time.Time `json:"scheduled_at" validate:"required"`
	Reason      string    `json:"reason,omitempty"`
	MeetingLink string    `json:"meeting_link,omitempty"`
	Location    string    `json:"location,omitempty"`
}

// CompleteInterviewRequest represents request to complete interview with evaluation
type CompleteInterviewRequest struct {
	OverallScore       *float64 `json:"overall_score,omitempty" validate:"omitempty,min=0,max=100"`
	TechnicalScore     *float64 `json:"technical_score,omitempty" validate:"omitempty,min=0,max=100"`
	CommunicationScore *float64 `json:"communication_score,omitempty" validate:"omitempty,min=0,max=100"`
	PersonalityScore   *float64 `json:"personality_score,omitempty" validate:"omitempty,min=0,max=100"`
	Remarks            string   `json:"remarks,omitempty"`
	FeedbackSummary    string   `json:"feedback_summary,omitempty"`
	CompletedBy        int64    `json:"completed_by" validate:"required"`
}

// ===== Response DTOs =====

// ApplicationListResponse represents paginated application list
type ApplicationListResponse struct {
	Applications []ApplicationSummary `json:"applications"`
	Total        int64                `json:"total"`
	Page         int                  `json:"page"`
	Limit        int                  `json:"limit"`
	TotalPages   int                  `json:"total_pages"`
	Stats        *ListStats           `json:"stats,omitempty"`
}

// ApplicationSummary represents summary of application for listing
type ApplicationSummary struct {
	ID               int64     `json:"id"`
	JobID            int64     `json:"job_id"`
	JobTitle         string    `json:"job_title"`
	CompanyName      string    `json:"company_name"`
	UserID           int64     `json:"user_id"`
	UserName         string    `json:"user_name"`
	Status           string    `json:"status"`
	MatchScore       float64   `json:"match_score"`
	AppliedAt        time.Time `json:"applied_at"`
	ViewedByEmployer bool      `json:"viewed_by_employer"`
	IsBookmarked     bool      `json:"is_bookmarked"`
	CurrentStage     string    `json:"current_stage"`
	DaysSinceApplied int       `json:"days_since_applied"`
}

// ApplicationDetailResponse represents detailed application information
type ApplicationDetailResponse struct {
	Application JobApplication        `json:"application"`
	Job         JobDetail             `json:"job"`
	Applicant   ApplicantProfile      `json:"applicant"`
	Stages      []JobApplicationStage `json:"stages"`
	Documents   []ApplicationDocument `json:"documents"`
	Notes       []ApplicationNote     `json:"notes"`
	Interviews  []Interview           `json:"interviews"`
	Stats       *ApplicationStats     `json:"stats,omitempty"`
}

// JobDetail represents job information in application context
type JobDetail struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	CompanyID   int64  `json:"company_id"`
	CompanyName string `json:"company_name"`
	Location    string `json:"location"`
	JobLevel    string `json:"job_level"`
	Status      string `json:"status"`
}

// ApplicantProfile represents applicant information
type ApplicantProfile struct {
	UserID      int64    `json:"user_id"`
	FullName    string   `json:"full_name"`
	Email       string   `json:"email"`
	Phone       string   `json:"phone"`
	PhotoURL    string   `json:"photo_url"`
	CurrentRole string   `json:"current_role"`
	Experience  int      `json:"experience_years"`
	Skills      []string `json:"skills"`
	Education   string   `json:"education"`
	ResumeURL   string   `json:"resume_url"`
}

// ListStats represents statistics for application list
type ListStats struct {
	TotalApplications int64            `json:"total_applications"`
	ViewedCount       int64            `json:"viewed_count"`
	BookmarkedCount   int64            `json:"bookmarked_count"`
	AverageMatchScore float64          `json:"average_match_score"`
	StatusBreakdown   map[string]int64 `json:"status_breakdown"`
}

// ApplicationAnalytics represents detailed application analytics
type ApplicationAnalytics struct {
	ApplicationID  int64              `json:"application_id"`
	Timeline       []TimelineEvent    `json:"timeline"`
	StageProgress  []StageProgress    `json:"stage_progress"`
	DocumentStats  DocumentStats      `json:"document_stats"`
	InterviewStats InterviewStats     `json:"interview_stats"`
	MatchAnalysis  MatchAnalysis      `json:"match_analysis"`
	ActivityLog    []ActivityLogEntry `json:"activity_log"`
}

// TimelineEvent represents timeline event
type TimelineEvent struct {
	Date        time.Time `json:"date"`
	EventType   string    `json:"event_type"`
	Description string    `json:"description"`
	Actor       string    `json:"actor"`
}

// StageProgress represents stage progress
type StageProgress struct {
	StageName   string     `json:"stage_name"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Duration    string     `json:"duration"`
	Status      string     `json:"status"`
}

// DocumentStats represents document statistics
type DocumentStats struct {
	TotalDocuments    int64            `json:"total_documents"`
	VerifiedDocuments int64            `json:"verified_documents"`
	DocumentTypes     map[string]int64 `json:"document_types"`
}

// InterviewStats represents interview statistics
type InterviewStats struct {
	TotalInterviews     int64   `json:"total_interviews"`
	CompletedInterviews int64   `json:"completed_interviews"`
	AverageScore        float64 `json:"average_score"`
	HighestScore        float64 `json:"highest_score"`
}

// MatchAnalysis represents match analysis
type MatchAnalysis struct {
	OverallScore    float64  `json:"overall_score"`
	SkillsMatch     float64  `json:"skills_match"`
	ExperienceMatch float64  `json:"experience_match"`
	EducationMatch  float64  `json:"education_match"`
	MatchedSkills   []string `json:"matched_skills"`
	MissingSkills   []string `json:"missing_skills"`
}

// ActivityLogEntry represents activity log entry
type ActivityLogEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	Action      string    `json:"action"`
	Actor       string    `json:"actor"`
	Description string    `json:"description"`
}

// JobApplicationAnalytics represents job application analytics
type JobApplicationAnalytics struct {
	JobID                int64             `json:"job_id"`
	JobTitle             string            `json:"job_title"`
	Period               string            `json:"period"`
	TotalApplications    int64             `json:"total_applications"`
	StatusBreakdown      map[string]int64  `json:"status_breakdown"`
	AverageMatchScore    float64           `json:"average_match_score"`
	ConversionFunnel     *ConversionFunnel `json:"conversion_funnel"`
	ApplicationsOverTime []TimeSeriesData  `json:"applications_over_time"`
	TopSources           []SourceStats     `json:"top_sources"`
	AverageTimeToHire    float64           `json:"average_time_to_hire"`
}

// CompanyApplicationAnalytics represents company application analytics
type CompanyApplicationAnalytics struct {
	CompanyID            int64            `json:"company_id"`
	CompanyName          string           `json:"company_name"`
	Period               string           `json:"period"`
	TotalApplications    int64            `json:"total_applications"`
	TotalHires           int64            `json:"total_hires"`
	ConversionRate       float64          `json:"conversion_rate"`
	AverageTimeToHire    float64          `json:"average_time_to_hire"`
	ApplicationsByJob    []JobStats       `json:"applications_by_job"`
	ApplicationsOverTime []TimeSeriesData `json:"applications_over_time"`
	SourceBreakdown      []SourceStats    `json:"source_breakdown"`
	StageTimeAnalysis    []StageTimeStats `json:"stage_time_analysis"`
}

// TimeSeriesData represents time-series data point
type TimeSeriesData struct {
	Date  time.Time `json:"date"`
	Value int64     `json:"value"`
}

// JobStats represents job statistics
type JobStats struct {
	JobID            int64   `json:"job_id"`
	JobTitle         string  `json:"job_title"`
	ApplicationCount int64   `json:"application_count"`
	HiredCount       int64   `json:"hired_count"`
	ConversionRate   float64 `json:"conversion_rate"`
}
