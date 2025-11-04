package repository

import (
	"context"

	"github.com/stretchr/testify/mock"

	"keerja-backend/internal/domain/job"
)

// MockJobRepository is a mock implementation of job.JobRepository
type MockJobRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockJobRepository) Create(ctx context.Context, j *job.Job) error {
	args := m.Called(ctx, j)
	return args.Error(0)
}

// FindByID mocks the FindByID method
func (m *MockJobRepository) FindByID(ctx context.Context, id int64) (*job.Job, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*job.Job), args.Error(1)
}

// FindByUUID mocks the FindByUUID method
func (m *MockJobRepository) FindByUUID(ctx context.Context, uuid string) (*job.Job, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*job.Job), args.Error(1)
}

// FindBySlug mocks the FindBySlug method
func (m *MockJobRepository) FindBySlug(ctx context.Context, slug string) (*job.Job, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*job.Job), args.Error(1)
}

// Update mocks the Update method
func (m *MockJobRepository) Update(ctx context.Context, j *job.Job) error {
	args := m.Called(ctx, j)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockJobRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// SoftDelete mocks the SoftDelete method
func (m *MockJobRepository) SoftDelete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// List mocks the List method
func (m *MockJobRepository) List(ctx context.Context, filter job.JobFilter, page, limit int) ([]job.Job, int64, error) {
	args := m.Called(ctx, filter, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]job.Job), args.Get(1).(int64), args.Error(2)
}

// ListByCompany mocks the ListByCompany method
func (m *MockJobRepository) ListByCompany(ctx context.Context, companyID int64, filter job.JobFilter, page, limit int) ([]job.Job, int64, error) {
	args := m.Called(ctx, companyID, filter, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]job.Job), args.Get(1).(int64), args.Error(2)
}

// ListByEmployer mocks the ListByEmployer method
func (m *MockJobRepository) ListByEmployer(ctx context.Context, employerUserID int64, filter job.JobFilter, page, limit int) ([]job.Job, int64, error) {
	args := m.Called(ctx, employerUserID, filter, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]job.Job), args.Get(1).(int64), args.Error(2)
}

// Search mocks the Search method
func (m *MockJobRepository) Search(ctx context.Context, query string, filter job.JobFilter, page, limit int) ([]job.Job, int64, error) {
	args := m.Called(ctx, query, filter, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]job.Job), args.Get(1).(int64), args.Error(2)
}

// GetFeaturedJobs mocks the GetFeaturedJobs method
func (m *MockJobRepository) GetFeaturedJobs(ctx context.Context, limit int) ([]job.Job, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]job.Job), args.Error(1)
}

// GetRecommendedJobs mocks the GetRecommendedJobs method
func (m *MockJobRepository) GetRecommendedJobs(ctx context.Context, userID int64, limit int) ([]job.Job, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]job.Job), args.Error(1)
}

// GetSimilarJobs mocks the GetSimilarJobs method
func (m *MockJobRepository) GetSimilarJobs(ctx context.Context, jobID int64, limit int) ([]job.Job, error) {
	args := m.Called(ctx, jobID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]job.Job), args.Error(1)
}

// IncrementViews mocks the IncrementViews method
func (m *MockJobRepository) IncrementViews(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// IncrementApplications mocks the IncrementApplications method
func (m *MockJobRepository) IncrementApplications(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// DecrementApplications mocks the DecrementApplications method
func (m *MockJobRepository) DecrementApplications(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// UpdateStatus mocks the UpdateStatus method
func (m *MockJobRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

// CloseExpiredJobs mocks the CloseExpiredJobs method
func (m *MockJobRepository) CloseExpiredJobs(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// CountJobs mocks the CountJobs method
func (m *MockJobRepository) CountJobs(ctx context.Context, filter job.JobFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

// CountByCompany mocks the CountByCompany method
func (m *MockJobRepository) CountByCompany(ctx context.Context, companyID int64, status string) (int64, error) {
	args := m.Called(ctx, companyID, status)
	return args.Get(0).(int64), args.Error(1)
}

// GetJobStats mocks the GetJobStats method
func (m *MockJobRepository) GetJobStats(ctx context.Context, jobID int64) (*job.JobStats, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*job.JobStats), args.Error(1)
}

// GetCompanyJobStats mocks the GetCompanyJobStats method
func (m *MockJobRepository) GetCompanyJobStats(ctx context.Context, companyID int64) (*job.CompanyJobStats, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*job.CompanyJobStats), args.Error(1)
}

// NOTE: methods yg dikomentari karena merujuk pada tipe yang belum didefinisikan 
// yang belum ada di domain pekerjaan. Hapus komentar saat tipe-tipe ini sudah didefinisikan di domain pekerjaan.:
// - GetPopularLocations() - requires job.LocationStats
// - GetTopCompanies() - requires job.CompanyStats (may exist, check)
// - CreateJobView() - requires job.JobView
// - GetJobViewsByUser() - uses job.JobView
// - GetJobTrends() - requires job.JobTrend

// GetTopCompanies mocks the GetTopCompanies method
// Keeping this as CompanyStats might exist in domain
func (m *MockJobRepository) GetTopCompanies(ctx context.Context, limit int) ([]struct{}, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]struct{}), args.Error(1)
}

// CreateSkill mocks the CreateSkill method
func (m *MockJobRepository) CreateSkill(ctx context.Context, skill *job.JobSkill) error {
	args := m.Called(ctx, skill)
	return args.Error(0)
}

// DeleteSkill mocks the DeleteSkill method
func (m *MockJobRepository) DeleteSkill(ctx context.Context, jobID, skillID int64) error {
	args := m.Called(ctx, jobID, skillID)
	return args.Error(0)
}

// CreateBenefit mocks the CreateBenefit method
func (m *MockJobRepository) CreateBenefit(ctx context.Context, benefit *job.JobBenefit) error {
	args := m.Called(ctx, benefit)
	return args.Error(0)
}

// DeleteBenefit mocks the DeleteBenefit method
func (m *MockJobRepository) DeleteBenefit(ctx context.Context, jobID, benefitID int64) error {
	args := m.Called(ctx, jobID, benefitID)
	return args.Error(0)
}

// SaveJob mocks the SaveJob method
func (m *MockJobRepository) SaveJob(ctx context.Context, userID, jobID int64) error {
	args := m.Called(ctx, userID, jobID)
	return args.Error(0)
}

// UnsaveJob mocks the UnsaveJob method
func (m *MockJobRepository) UnsaveJob(ctx context.Context, userID, jobID int64) error {
	args := m.Called(ctx, userID, jobID)
	return args.Error(0)
}

// GetSavedJobs mocks the GetSavedJobs method
func (m *MockJobRepository) GetSavedJobs(ctx context.Context, userID int64, page, limit int) ([]job.Job, int64, error) {
	args := m.Called(ctx, userID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]job.Job), args.Get(1).(int64), args.Error(2)
}

// IsJobSaved mocks the IsJobSaved method
func (m *MockJobRepository) IsJobSaved(ctx context.Context, userID, jobID int64) (bool, error) {
	args := m.Called(ctx, userID, jobID)
	return args.Get(0).(bool), args.Error(1)
}

// NOTE: Commented out methods that reference undefined types:
// CreateJobView, GetJobViewsByUser, GetJobTrends
// Uncomment when job.JobView and job.JobTrend types are defined in domain

// GetCategoryStats mocks the GetCategoryStats method
func (m *MockJobRepository) GetCategoryStats(ctx context.Context) ([]job.CategoryStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]job.CategoryStats), args.Error(1)
}

// PublishJob mocks the PublishJob method
func (m *MockJobRepository) PublishJob(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// CloseJob mocks the CloseJob method
func (m *MockJobRepository) CloseJob(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
