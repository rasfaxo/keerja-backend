package repository

import (
	"context"

	"github.com/stretchr/testify/mock"

	"keerja-backend/internal/domain/company"
)

// MockCompanyRepository is a mock implementation of company.CompanyRepository
type MockCompanyRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockCompanyRepository) Create(ctx context.Context, c *company.Company) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

// FindByID mocks the FindByID method
func (m *MockCompanyRepository) FindByID(ctx context.Context, id int64) (*company.Company, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.Company), args.Error(1)
}

// FindByUUID mocks the FindByUUID method
func (m *MockCompanyRepository) FindByUUID(ctx context.Context, uuid string) (*company.Company, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.Company), args.Error(1)
}

// FindBySlug mocks the FindBySlug method
func (m *MockCompanyRepository) FindBySlug(ctx context.Context, slug string) (*company.Company, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.Company), args.Error(1)
}

// Update mocks the Update method
func (m *MockCompanyRepository) Update(ctx context.Context, c *company.Company) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockCompanyRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// List mocks the List method
func (m *MockCompanyRepository) List(ctx context.Context, filter company.CompanyFilter, page, limit int) ([]company.Company, int64, error) {
	args := m.Called(ctx, filter, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]company.Company), args.Get(1).(int64), args.Error(2)
}

// SearchCompanies mocks the SearchCompanies method
func (m *MockCompanyRepository) SearchCompanies(ctx context.Context, query string, page, limit int) ([]company.Company, int64, error) {
	args := m.Called(ctx, query, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]company.Company), args.Get(1).(int64), args.Error(2)
}

// GetVerifiedCompanies mocks the GetVerifiedCompanies method
func (m *MockCompanyRepository) GetVerifiedCompanies(ctx context.Context, limit int) ([]company.Company, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]company.Company), args.Error(1)
}

// GetTopRatedCompanies mocks the GetTopRatedCompanies method
func (m *MockCompanyRepository) GetTopRatedCompanies(ctx context.Context, limit int) ([]company.Company, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]company.Company), args.Error(1)
}

// CreateProfile mocks the CreateProfile method
func (m *MockCompanyRepository) CreateProfile(ctx context.Context, profile *company.CompanyProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

// FindProfileByCompanyID mocks the FindProfileByCompanyID method
func (m *MockCompanyRepository) FindProfileByCompanyID(ctx context.Context, companyID int64) (*company.CompanyProfile, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.CompanyProfile), args.Error(1)
}

// UpdateProfile mocks the UpdateProfile method
func (m *MockCompanyRepository) UpdateProfile(ctx context.Context, profile *company.CompanyProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

// CreateEmployerUser mocks the CreateEmployerUser method
func (m *MockCompanyRepository) CreateEmployerUser(ctx context.Context, employerUser *company.EmployerUser) error {
	args := m.Called(ctx, employerUser)
	return args.Error(0)
}

// FindEmployerUserByID mocks the FindEmployerUserByID method
func (m *MockCompanyRepository) FindEmployerUserByID(ctx context.Context, id int64) (*company.EmployerUser, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.EmployerUser), args.Error(1)
}

// FindEmployerUserByUserAndCompany mocks the FindEmployerUserByUserAndCompany method
func (m *MockCompanyRepository) FindEmployerUserByUserAndCompany(ctx context.Context, userID, companyID int64) (*company.EmployerUser, error) {
	args := m.Called(ctx, userID, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.EmployerUser), args.Error(1)
}

// UpdateEmployerUser mocks the UpdateEmployerUser method
func (m *MockCompanyRepository) UpdateEmployerUser(ctx context.Context, employerUser *company.EmployerUser) error {
	args := m.Called(ctx, employerUser)
	return args.Error(0)
}

// DeleteEmployerUser mocks the DeleteEmployerUser method
func (m *MockCompanyRepository) DeleteEmployerUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetEmployerUsersByCompanyID mocks the GetEmployerUsersByCompanyID method
func (m *MockCompanyRepository) GetEmployerUsersByCompanyID(ctx context.Context, companyID int64) ([]company.EmployerUser, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]company.EmployerUser), args.Error(1)
}

// GetCompaniesByUserID mocks the GetCompaniesByUserID method
func (m *MockCompanyRepository) GetCompaniesByUserID(ctx context.Context, userID int64) ([]company.Company, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]company.Company), args.Error(1)
}

// CreateReview mocks the CreateReview method
func (m *MockCompanyRepository) CreateReview(ctx context.Context, review *company.CompanyReview) error {
	args := m.Called(ctx, review)
	return args.Error(0)
}

// UpdateReview mocks the UpdateReview method
func (m *MockCompanyRepository) UpdateReview(ctx context.Context, review *company.CompanyReview) error {
	args := m.Called(ctx, review)
	return args.Error(0)
}

// DeleteReview mocks the DeleteReview method
func (m *MockCompanyRepository) DeleteReview(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetCompanyReviews mocks the GetCompanyReviews method
func (m *MockCompanyRepository) GetCompanyReviews(ctx context.Context, companyID int64, filter company.ReviewFilter, page, limit int) ([]company.CompanyReview, int64, error) {
	args := m.Called(ctx, companyID, filter, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]company.CompanyReview), args.Get(1).(int64), args.Error(2)
}

// GetAverageRatings mocks the GetAverageRatings method
func (m *MockCompanyRepository) GetAverageRatings(ctx context.Context, companyID int64) (*company.AverageRatings, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.AverageRatings), args.Error(1)
}

// CreateDocument mocks the CreateDocument method
func (m *MockCompanyRepository) CreateDocument(ctx context.Context, doc *company.CompanyDocument) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

// UpdateDocument mocks the UpdateDocument method
func (m *MockCompanyRepository) UpdateDocument(ctx context.Context, doc *company.CompanyDocument) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

// DeleteDocument mocks the DeleteDocument method
func (m *MockCompanyRepository) DeleteDocument(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetCompanyDocuments mocks the GetCompanyDocuments method
func (m *MockCompanyRepository) GetCompanyDocuments(ctx context.Context, companyID int64) ([]company.CompanyDocument, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]company.CompanyDocument), args.Error(1)
}

// CreateVerification mocks the CreateVerification method
func (m *MockCompanyRepository) CreateVerification(ctx context.Context, verification *company.CompanyVerification) error {
	args := m.Called(ctx, verification)
	return args.Error(0)
}

// UpdateVerification mocks the UpdateVerification method
func (m *MockCompanyRepository) UpdateVerification(ctx context.Context, verification *company.CompanyVerification) error {
	args := m.Called(ctx, verification)
	return args.Error(0)
}

// FindVerificationByCompanyID mocks the FindVerificationByCompanyID method
func (m *MockCompanyRepository) FindVerificationByCompanyID(ctx context.Context, companyID int64) (*company.CompanyVerification, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.CompanyVerification), args.Error(1)
}

// GetPendingVerifications mocks the GetPendingVerifications method
func (m *MockCompanyRepository) GetPendingVerifications(ctx context.Context) ([]company.CompanyVerification, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]company.CompanyVerification), args.Error(1)
}

// GetCompaniesNeedingVerificationRenewal mocks the GetCompaniesNeedingVerificationRenewal method
func (m *MockCompanyRepository) GetCompaniesNeedingVerificationRenewal(ctx context.Context, daysBeforeExpiry int) ([]company.Company, error) {
	args := m.Called(ctx, daysBeforeExpiry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]company.Company), args.Error(1)
}

// CreateIndustry mocks the CreateIndustry method
func (m *MockCompanyRepository) CreateIndustry(ctx context.Context, industry *company.CompanyIndustry) error {
	args := m.Called(ctx, industry)
	return args.Error(0)
}

// FindIndustryByID mocks the FindIndustryByID method
func (m *MockCompanyRepository) FindIndustryByID(ctx context.Context, id int64) (*company.CompanyIndustry, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.CompanyIndustry), args.Error(1)
}

// FindIndustryByCode mocks the FindIndustryByCode method
func (m *MockCompanyRepository) FindIndustryByCode(ctx context.Context, code string) (*company.CompanyIndustry, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.CompanyIndustry), args.Error(1)
}

// UpdateIndustry mocks the UpdateIndustry method
func (m *MockCompanyRepository) UpdateIndustry(ctx context.Context, industry *company.CompanyIndustry) error {
	args := m.Called(ctx, industry)
	return args.Error(0)
}

// DeleteIndustry mocks the DeleteIndustry method
func (m *MockCompanyRepository) DeleteIndustry(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetAllIndustries mocks the GetAllIndustries method
func (m *MockCompanyRepository) GetAllIndustries(ctx context.Context, activeOnly bool) ([]company.CompanyIndustry, error) {
	args := m.Called(ctx, activeOnly)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]company.CompanyIndustry), args.Error(1)
}

// GetIndustryTree mocks the GetIndustryTree method
func (m *MockCompanyRepository) GetIndustryTree(ctx context.Context) ([]company.CompanyIndustry, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]company.CompanyIndustry), args.Error(1)
}
