package repository

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"keerja-backend/internal/domain/application"
)

// MockApplicationRepository is a mock implementation of application.ApplicationRepository
type MockApplicationRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockApplicationRepository) Create(ctx context.Context, app *application.JobApplication) error {
	args := m.Called(ctx, app)
	return args.Error(0)
}

// FindByID mocks the FindByID method
func (m *MockApplicationRepository) FindByID(ctx context.Context, id int64) (*application.JobApplication, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*application.JobApplication), args.Error(1)
}

// FindByJobAndUser mocks the FindByJobAndUser method
func (m *MockApplicationRepository) FindByJobAndUser(ctx context.Context, jobID, userID int64) (*application.JobApplication, error) {
	args := m.Called(ctx, jobID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*application.JobApplication), args.Error(1)
}

// Update mocks the Update method
func (m *MockApplicationRepository) Update(ctx context.Context, app *application.JobApplication) error {
	args := m.Called(ctx, app)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockApplicationRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// List mocks the List method
func (m *MockApplicationRepository) List(ctx context.Context, filter application.ApplicationFilter, page, limit int) ([]application.JobApplication, int64, error) {
	args := m.Called(ctx, filter, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]application.JobApplication), args.Get(1).(int64), args.Error(2)
}

// ListByUser mocks the ListByUser method
func (m *MockApplicationRepository) ListByUser(ctx context.Context, userID int64, filter application.ApplicationFilter, page, limit int) ([]application.JobApplication, int64, error) {
	args := m.Called(ctx, userID, filter, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]application.JobApplication), args.Get(1).(int64), args.Error(2)
}

// ListByJob mocks the ListByJob method
func (m *MockApplicationRepository) ListByJob(ctx context.Context, jobID int64, filter application.ApplicationFilter, page, limit int) ([]application.JobApplication, int64, error) {
	args := m.Called(ctx, jobID, filter, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]application.JobApplication), args.Get(1).(int64), args.Error(2)
}

// ListByCompany mocks the ListByCompany method
func (m *MockApplicationRepository) ListByCompany(ctx context.Context, companyID int64, filter application.ApplicationFilter, page, limit int) ([]application.JobApplication, int64, error) {
	args := m.Called(ctx, companyID, filter, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]application.JobApplication), args.Get(1).(int64), args.Error(2)
}

// UpdateStatus mocks the UpdateStatus method
func (m *MockApplicationRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

// BulkUpdateStatus mocks the BulkUpdateStatus method
func (m *MockApplicationRepository) BulkUpdateStatus(ctx context.Context, ids []int64, status string) error {
	args := m.Called(ctx, ids, status)
	return args.Error(0)
}

// GetApplicationsByStatus mocks the GetApplicationsByStatus method
func (m *MockApplicationRepository) GetApplicationsByStatus(ctx context.Context, jobID int64, status string) ([]application.JobApplication, error) {
	args := m.Called(ctx, jobID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]application.JobApplication), args.Error(1)
}

// MarkAsViewed mocks the MarkAsViewed method
func (m *MockApplicationRepository) MarkAsViewed(ctx context.Context, id int64, viewed bool) error {
	args := m.Called(ctx, id, viewed)
	return args.Error(0)
}

// ToggleBookmark mocks the ToggleBookmark method
func (m *MockApplicationRepository) ToggleBookmark(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetBookmarkedApplications mocks the GetBookmarkedApplications method
func (m *MockApplicationRepository) GetBookmarkedApplications(ctx context.Context, companyID int64, page, limit int) ([]application.JobApplication, int64, error) {
	args := m.Called(ctx, companyID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]application.JobApplication), args.Get(1).(int64), args.Error(2)
}

// SearchApplications mocks the SearchApplications method
func (m *MockApplicationRepository) SearchApplications(ctx context.Context, query string, filter application.ApplicationFilter, page, limit int) ([]application.JobApplication, int64, error) {
	args := m.Called(ctx, query, filter, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]application.JobApplication), args.Get(1).(int64), args.Error(2)
}

// GetApplicationsWithHighScore mocks the GetApplicationsWithHighScore method
func (m *MockApplicationRepository) GetApplicationsWithHighScore(ctx context.Context, jobID int64, minScore float64, limit int) ([]application.JobApplication, error) {
	args := m.Called(ctx, jobID, minScore, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]application.JobApplication), args.Error(1)
}

// CountApplications mocks the CountApplications method
func (m *MockApplicationRepository) CountApplications(ctx context.Context, filter application.ApplicationFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

// CountByJob mocks the CountByJob method
func (m *MockApplicationRepository) CountByJob(ctx context.Context, jobID int64, status string) (int64, error) {
	args := m.Called(ctx, jobID, status)
	return args.Get(0).(int64), args.Error(1)
}

// CountByUser mocks the CountByUser method
func (m *MockApplicationRepository) CountByUser(ctx context.Context, userID int64, status string) (int64, error) {
	args := m.Called(ctx, userID, status)
	return args.Get(0).(int64), args.Error(1)
}

// CountByCompany mocks the CountByCompany method
func (m *MockApplicationRepository) CountByCompany(ctx context.Context, companyID int64, status string) (int64, error) {
	args := m.Called(ctx, companyID, status)
	return args.Get(0).(int64), args.Error(1)
}

// GetUserApplicationStats mocks the GetUserApplicationStats method
func (m *MockApplicationRepository) GetUserApplicationStats(ctx context.Context, userID int64) (*application.UserApplicationStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*application.UserApplicationStats), args.Error(1)
}

// GetJobApplicationStats mocks the GetJobApplicationStats method
func (m *MockApplicationRepository) GetJobApplicationStats(ctx context.Context, jobID int64) (*application.JobApplicationStats, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*application.JobApplicationStats), args.Error(1)
}

// GetCompanyApplicationStats mocks the GetCompanyApplicationStats method
func (m *MockApplicationRepository) GetCompanyApplicationStats(ctx context.Context, companyID int64) (*application.CompanyApplicationStats, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*application.CompanyApplicationStats), args.Error(1)
}

// CreateStage mocks the CreateStage method
func (m *MockApplicationRepository) CreateStage(ctx context.Context, stage *application.JobApplicationStage) error {
	args := m.Called(ctx, stage)
	return args.Error(0)
}

// GetCurrentStage mocks the GetCurrentStage method
func (m *MockApplicationRepository) GetCurrentStage(ctx context.Context, applicationID int64) (*application.JobApplicationStage, error) {
	args := m.Called(ctx, applicationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*application.JobApplicationStage), args.Error(1)
}

// GetStageHistory mocks the GetStageHistory method
func (m *MockApplicationRepository) GetStageHistory(ctx context.Context, applicationID int64) ([]application.JobApplicationStage, error) {
	args := m.Called(ctx, applicationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]application.JobApplicationStage), args.Error(1)
}

// CompleteStage mocks the CompleteStage method
func (m *MockApplicationRepository) CompleteStage(ctx context.Context, stageID int64, notes string) error {
	args := m.Called(ctx, stageID, notes)
	return args.Error(0)
}

// GetAverageTimePerStage mocks the GetAverageTimePerStage method
func (m *MockApplicationRepository) GetAverageTimePerStage(ctx context.Context, companyID int64) ([]application.StageTimeStats, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]application.StageTimeStats), args.Error(1)
}

// CreateNote mocks the CreateNote method
func (m *MockApplicationRepository) CreateNote(ctx context.Context, note *application.ApplicationNote) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

// UpdateNote mocks the UpdateNote method
func (m *MockApplicationRepository) UpdateNote(ctx context.Context, note *application.ApplicationNote) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

// DeleteNote mocks the DeleteNote method
func (m *MockApplicationRepository) DeleteNote(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetApplicationNotes mocks the GetApplicationNotes method
func (m *MockApplicationRepository) GetApplicationNotes(ctx context.Context, applicationID int64) ([]application.ApplicationNote, error) {
	args := m.Called(ctx, applicationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]application.ApplicationNote), args.Error(1)
}

// CreateDocument mocks the CreateDocument method
func (m *MockApplicationRepository) CreateDocument(ctx context.Context, doc *application.ApplicationDocument) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

// DeleteDocument mocks the DeleteDocument method
func (m *MockApplicationRepository) DeleteDocument(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetApplicationDocuments mocks the GetApplicationDocuments method
func (m *MockApplicationRepository) GetApplicationDocuments(ctx context.Context, applicationID int64) ([]application.ApplicationDocument, error) {
	args := m.Called(ctx, applicationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]application.ApplicationDocument), args.Error(1)
}

// CreateInterview mocks the CreateInterview method
func (m *MockApplicationRepository) CreateInterview(ctx context.Context, interview *application.Interview) error {
	args := m.Called(ctx, interview)
	return args.Error(0)
}

// UpdateInterview mocks the UpdateInterview method
func (m *MockApplicationRepository) UpdateInterview(ctx context.Context, interview *application.Interview) error {
	args := m.Called(ctx, interview)
	return args.Error(0)
}

// DeleteInterview mocks the DeleteInterview method
func (m *MockApplicationRepository) DeleteInterview(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetApplicationInterviews mocks the GetApplicationInterviews method
func (m *MockApplicationRepository) GetApplicationInterviews(ctx context.Context, applicationID int64) ([]application.Interview, error) {
	args := m.Called(ctx, applicationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]application.Interview), args.Error(1)
}

// GetUpcomingInterviews mocks the GetUpcomingInterviews method
func (m *MockApplicationRepository) GetUpcomingInterviews(ctx context.Context, companyID int64, startDate, endDate time.Time) ([]application.Interview, error) {
	args := m.Called(ctx, companyID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]application.Interview), args.Error(1)
}

// NOTE: The following evaluation methods are commented out because they reference
// undefined type application.InterviewEvaluation that doesn't exist in the domain yet:
// - CreateEvaluation()
// - UpdateEvaluation()
// - GetInterviewEvaluations()
// Uncomment when InterviewEvaluation type is defined

// GetConversionFunnel mocks the GetConversionFunnel method
func (m *MockApplicationRepository) GetConversionFunnel(ctx context.Context, jobID int64) (*application.ConversionFunnel, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*application.ConversionFunnel), args.Error(1)
}

// GetApplicationTrends mocks the GetApplicationTrends method
func (m *MockApplicationRepository) GetApplicationTrends(ctx context.Context, companyID int64, startDate, endDate time.Time) ([]application.ApplicationTrend, error) {
	args := m.Called(ctx, companyID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]application.ApplicationTrend), args.Error(1)
}

// GetSourceAnalytics mocks the GetSourceAnalytics method
func (m *MockApplicationRepository) GetSourceAnalytics(ctx context.Context, companyID int64) ([]application.SourceStats, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]application.SourceStats), args.Error(1)
}
