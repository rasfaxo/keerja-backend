package application

import (
	"context"
	"time"
)

// ApplicationRepository defines the interface for application data access
type ApplicationRepository interface {
	// JobApplication CRUD operations
	Create(ctx context.Context, application *JobApplication) error
	FindByID(ctx context.Context, id int64) (*JobApplication, error)
	FindByJobAndUser(ctx context.Context, jobID, userID int64) (*JobApplication, error)
	Update(ctx context.Context, application *JobApplication) error
	Delete(ctx context.Context, id int64) error
	
	// Application listing and filtering
	List(ctx context.Context, filter ApplicationFilter, page, limit int) ([]JobApplication, int64, error)
	ListByUser(ctx context.Context, userID int64, filter ApplicationFilter, page, limit int) ([]JobApplication, int64, error)
	ListByJob(ctx context.Context, jobID int64, filter ApplicationFilter, page, limit int) ([]JobApplication, int64, error)
	ListByCompany(ctx context.Context, companyID int64, filter ApplicationFilter, page, limit int) ([]JobApplication, int64, error)
	
	// Application status operations
	UpdateStatus(ctx context.Context, id int64, status string) error
	BulkUpdateStatus(ctx context.Context, ids []int64, status string) error
	GetApplicationsByStatus(ctx context.Context, status string, page, limit int) ([]JobApplication, int64, error)
	
	// Application tracking
	MarkAsViewed(ctx context.Context, id int64) error
	ToggleBookmark(ctx context.Context, id int64) error
	GetBookmarkedApplications(ctx context.Context, companyID int64, page, limit int) ([]JobApplication, int64, error)
	
	// Application search
	SearchApplications(ctx context.Context, filter ApplicationSearchFilter, page, limit int) ([]JobApplication, int64, error)
	GetApplicationsWithHighScore(ctx context.Context, minScore float64, limit int) ([]JobApplication, error)
	
	// Application statistics
	GetApplicationStats(ctx context.Context, applicationID int64) (*ApplicationStats, error)
	GetUserApplicationStats(ctx context.Context, userID int64) (*UserApplicationStats, error)
	GetJobApplicationStats(ctx context.Context, jobID int64) (*JobApplicationStats, error)
	GetCompanyApplicationStats(ctx context.Context, companyID int64) (*CompanyApplicationStats, error)
	
	// JobApplicationStage operations
	CreateStage(ctx context.Context, stage *JobApplicationStage) error
	FindStageByID(ctx context.Context, id int64) (*JobApplicationStage, error)
	UpdateStage(ctx context.Context, stage *JobApplicationStage) error
	CompleteStage(ctx context.Context, id int64, notes string) error
	ListStagesByApplication(ctx context.Context, applicationID int64) ([]JobApplicationStage, error)
	GetCurrentStage(ctx context.Context, applicationID int64) (*JobApplicationStage, error)
	GetStageHistory(ctx context.Context, applicationID int64) ([]JobApplicationStage, error)
	
	// ApplicationDocument operations
	CreateDocument(ctx context.Context, document *ApplicationDocument) error
	FindDocumentByID(ctx context.Context, id int64) (*ApplicationDocument, error)
	UpdateDocument(ctx context.Context, document *ApplicationDocument) error
	DeleteDocument(ctx context.Context, id int64) error
	ListDocumentsByApplication(ctx context.Context, applicationID int64) ([]ApplicationDocument, error)
	ListDocumentsByUser(ctx context.Context, userID int64, docType string) ([]ApplicationDocument, error)
	GetDocumentsByType(ctx context.Context, applicationID int64, docType string) ([]ApplicationDocument, error)
	VerifyDocument(ctx context.Context, id int64, verifiedBy int64) error
	GetUnverifiedDocuments(ctx context.Context, page, limit int) ([]ApplicationDocument, int64, error)
	
	// ApplicationNote operations
	CreateNote(ctx context.Context, note *ApplicationNote) error
	FindNoteByID(ctx context.Context, id int64) (*ApplicationNote, error)
	UpdateNote(ctx context.Context, note *ApplicationNote) error
	DeleteNote(ctx context.Context, id int64) error
	ListNotesByApplication(ctx context.Context, applicationID int64) ([]ApplicationNote, error)
	ListNotesByStage(ctx context.Context, stageID int64) ([]ApplicationNote, error)
	ListNotesByAuthor(ctx context.Context, authorID int64, page, limit int) ([]ApplicationNote, int64, error)
	GetPinnedNotes(ctx context.Context, applicationID int64) ([]ApplicationNote, error)
	PinNote(ctx context.Context, id int64) error
	UnpinNote(ctx context.Context, id int64) error
	
	// Interview operations
	CreateInterview(ctx context.Context, interview *Interview) error
	FindInterviewByID(ctx context.Context, id int64) (*Interview, error)
	UpdateInterview(ctx context.Context, interview *Interview) error
	DeleteInterview(ctx context.Context, id int64) error
	ListInterviewsByApplication(ctx context.Context, applicationID int64) ([]Interview, error)
	ListInterviewsByInterviewer(ctx context.Context, interviewerID int64, filter InterviewFilter) ([]Interview, error)
	GetUpcomingInterviews(ctx context.Context, date time.Time, limit int) ([]Interview, error)
	GetInterviewsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]Interview, error)
	UpdateInterviewStatus(ctx context.Context, id int64, status string) error
	CompleteInterview(ctx context.Context, id int64, scores InterviewScores, feedback string) error
	RescheduleInterview(ctx context.Context, id int64, newSchedule time.Time) error
	CancelInterview(ctx context.Context, id int64) error
	
	// Analytics and reporting
	GetApplicationTrends(ctx context.Context, startDate, endDate time.Time) ([]ApplicationTrend, error)
	GetConversionFunnel(ctx context.Context, jobID int64) (*ConversionFunnel, error)
	GetAverageTimePerStage(ctx context.Context, companyID int64) ([]StageTimeStats, error)
	GetTopApplicants(ctx context.Context, jobID int64, limit int) ([]JobApplication, error)
	GetApplicationSourceStats(ctx context.Context, companyID int64) ([]SourceStats, error)
	
	// Bulk operations
	BulkCreateApplications(ctx context.Context, applications []JobApplication) error
	BulkDeleteApplications(ctx context.Context, ids []int64) error
}

// ApplicationFilter defines filter criteria for application listing
type ApplicationFilter struct {
	Status         string
	JobID          int64
	UserID         int64
	CompanyID      int64
	MinScore       *float64
	MaxScore       *float64
	ViewedOnly     *bool
	BookmarkedOnly *bool
	Source         string
	AppliedAfter   *time.Time
	AppliedBefore  *time.Time
	SortBy         string // "latest", "score_desc", "score_asc"
}

// ApplicationSearchFilter defines advanced search criteria
type ApplicationSearchFilter struct {
	Keyword        string
	JobIDs         []int64
	CompanyIDs     []int64
	Statuses       []string
	MinScore       *float64
	Sources        []string
	AppliedWithin  *int // days
	HasDocuments   *bool
	HasInterviews  *bool
}

// InterviewFilter defines filter criteria for interview listing
type InterviewFilter struct {
	Status        string
	InterviewType string
	ScheduledFrom *time.Time
	ScheduledTo   *time.Time
	CompletedOnly *bool
}

// ApplicationStats represents application statistics
type ApplicationStats struct {
	ApplicationID     int64
	TotalStages       int64
	CompletedStages   int64
	CurrentStage      string
	TotalDocuments    int64
	VerifiedDocuments int64
	TotalNotes        int64
	TotalInterviews   int64
	CompletedInterviews int64
	AverageInterviewScore float64
	DaysSinceApplied  int
	LastActivity      time.Time
}

// UserApplicationStats represents user's application statistics
type UserApplicationStats struct {
	UserID              int64
	TotalApplications   int64
	AppliedCount        int64
	ScreeningCount      int64
	ShortlistedCount    int64
	InterviewCount      int64
	OfferedCount        int64
	HiredCount          int64
	RejectedCount       int64
	WithdrawnCount      int64
	AverageMatchScore   float64
	SuccessRate         float64
	AverageResponseTime float64 // in days
}

// JobApplicationStats represents job's application statistics
type JobApplicationStats struct {
	JobID               int64
	TotalApplications   int64
	AppliedCount        int64
	ScreeningCount      int64
	ShortlistedCount    int64
	InterviewCount      int64
	OfferedCount        int64
	HiredCount          int64
	RejectedCount       int64
	WithdrawnCount      int64
	AverageMatchScore   float64
	ConversionRate      float64
	AverageTimeToHire   float64 // in days
	TopSources          []SourceCount
}

// CompanyApplicationStats represents company's application statistics
type CompanyApplicationStats struct {
	CompanyID           int64
	TotalApplications   int64
	ViewedApplications  int64
	BookmarkedApplications int64
	TotalHires          int64
	AverageMatchScore   float64
	AverageTimeToHire   float64 // in days
	ConversionRate      float64
	TopPerformingJobs   []JobPerformance
	ApplicationsByMonth []MonthlyCount
}

// ApplicationTrend represents application trend data
type ApplicationTrend struct {
	Date             time.Time
	TotalApplications int64
	HiredCount       int64
	RejectedCount    int64
	AverageMatchScore float64
}

// ConversionFunnel represents hiring funnel metrics
type ConversionFunnel struct {
	JobID             int64
	AppliedCount      int64
	ScreeningCount    int64
	ShortlistedCount  int64
	InterviewCount    int64
	OfferedCount      int64
	HiredCount        int64
	ConversionRates   map[string]float64 // stage -> rate
}

// StageTimeStats represents average time spent in each stage
type StageTimeStats struct {
	StageName   string
	AverageDays float64
	MinDays     int
	MaxDays     int
	Count       int64
}

// SourceStats represents application source statistics
type SourceStats struct {
	Source            string
	Count             int64
	HiredCount        int64
	ConversionRate    float64
	AverageMatchScore float64
}

// SourceCount represents source count
type SourceCount struct {
	Source string
	Count  int64
}

// JobPerformance represents job performance metrics
type JobPerformance struct {
	JobID             int64
	JobTitle          string
	ApplicationCount  int64
	HiredCount        int64
	ConversionRate    float64
	AverageMatchScore float64
}

// MonthlyCount represents monthly application count
type MonthlyCount struct {
	Month string
	Count int64
}

// InterviewScores represents interview evaluation scores
type InterviewScores struct {
	OverallScore       *float64
	TechnicalScore     *float64
	CommunicationScore *float64
	PersonalityScore   *float64
}
