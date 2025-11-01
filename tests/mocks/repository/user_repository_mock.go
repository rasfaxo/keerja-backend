package repository

import (
	"context"

	"github.com/stretchr/testify/mock"

	"keerja-backend/internal/domain/user"
)

// MockUserRepository is a mock implementation of user.UserRepository
type MockUserRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockUserRepository) Create(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

// FindByID mocks the FindByID method
func (m *MockUserRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

// FindByEmail mocks the FindByEmail method
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

// FindByPhone mocks the FindByPhone method
func (m *MockUserRepository) FindByPhone(ctx context.Context, phone string) (*user.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

// FindByUUID mocks the FindByUUID method
func (m *MockUserRepository) FindByUUID(ctx context.Context, uuid string) (*user.User, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

// Update mocks the Update method
func (m *MockUserRepository) Update(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// SoftDelete mocks the SoftDelete method
func (m *MockUserRepository) SoftDelete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// List mocks the List method
func (m *MockUserRepository) List(ctx context.Context, filter *user.UserFilter) ([]user.User, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]user.User), args.Get(1).(int64), args.Error(2)
}

// Search mocks the Search method
func (m *MockUserRepository) Search(ctx context.Context, query string, userType string, page, limit int) ([]user.User, int64, error) {
	args := m.Called(ctx, query, userType, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]user.User), args.Get(1).(int64), args.Error(2)
}

// UpdateEmail mocks the UpdateEmail method
func (m *MockUserRepository) UpdateEmail(ctx context.Context, id int64, email string) error {
	args := m.Called(ctx, id, email)
	return args.Error(0)
}

// UpdatePhone mocks the UpdatePhone method
func (m *MockUserRepository) UpdatePhone(ctx context.Context, id int64, phone string) error {
	args := m.Called(ctx, id, phone)
	return args.Error(0)
}

// UpdatePassword mocks the UpdatePassword method
func (m *MockUserRepository) UpdatePassword(ctx context.Context, id int64, hashedPassword string) error {
	args := m.Called(ctx, id, hashedPassword)
	return args.Error(0)
}

// UpdateLastLogin mocks the UpdateLastLogin method
func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// UpdateEmailVerification mocks the UpdateEmailVerification method
func (m *MockUserRepository) UpdateEmailVerification(ctx context.Context, id int64, verified bool) error {
	args := m.Called(ctx, id, verified)
	return args.Error(0)
}

// UpdatePhoneVerification mocks the UpdatePhoneVerification method
func (m *MockUserRepository) UpdatePhoneVerification(ctx context.Context, id int64, verified bool) error {
	args := m.Called(ctx, id, verified)
	return args.Error(0)
}

// CountUsers mocks the CountUsers method
func (m *MockUserRepository) CountUsers(ctx context.Context, userType string) (int64, error) {
	args := m.Called(ctx, userType)
	return args.Get(0).(int64), args.Error(1)
}

// GetActiveUsers mocks the GetActiveUsers method
func (m *MockUserRepository) GetActiveUsers(ctx context.Context, limit int) ([]user.User, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.User), args.Error(1)
}

// GetNewUsers mocks the GetNewUsers method
func (m *MockUserRepository) GetNewUsers(ctx context.Context, days int) ([]user.User, error) {
	args := m.Called(ctx, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.User), args.Error(1)
}

// CreateProfile mocks the CreateProfile method
func (m *MockUserRepository) CreateProfile(ctx context.Context, profile *user.UserProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

// FindProfileByUserID mocks the FindProfileByUserID method
func (m *MockUserRepository) FindProfileByUserID(ctx context.Context, userID int64) (*user.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.UserProfile), args.Error(1)
}

// FindProfileBySlug mocks the FindProfileBySlug method
func (m *MockUserRepository) FindProfileBySlug(ctx context.Context, slug string) (*user.UserProfile, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.UserProfile), args.Error(1)
}

// UpdateProfile mocks the UpdateProfile method
func (m *MockUserRepository) UpdateProfile(ctx context.Context, profile *user.UserProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

// CreatePreference mocks the CreatePreference method
func (m *MockUserRepository) CreatePreference(ctx context.Context, preference *user.UserPreference) error {
	args := m.Called(ctx, preference)
	return args.Error(0)
}

// FindPreferenceByUserID mocks the FindPreferenceByUserID method
func (m *MockUserRepository) FindPreferenceByUserID(ctx context.Context, userID int64) (*user.UserPreference, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.UserPreference), args.Error(1)
}

// UpdatePreference mocks the UpdatePreference method
func (m *MockUserRepository) UpdatePreference(ctx context.Context, preference *user.UserPreference) error {
	args := m.Called(ctx, preference)
	return args.Error(0)
}

// GetEducationsByUserID mocks the GetEducationsByUserID method
func (m *MockUserRepository) GetEducationsByUserID(ctx context.Context, userID int64) ([]user.UserEducation, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.UserEducation), args.Error(1)
}

// GetExperiencesByUserID mocks the GetExperiencesByUserID method
func (m *MockUserRepository) GetExperiencesByUserID(ctx context.Context, userID int64) ([]user.UserExperience, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.UserExperience), args.Error(1)
}

// GetFullProfile mocks the GetFullProfile method
func (m *MockUserRepository) GetFullProfile(ctx context.Context, userID int64) (*user.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

// AddExperience mocks AddExperience method
func (m *MockUserRepository) AddExperience(ctx context.Context, experience *user.UserExperience) error {
	args := m.Called(ctx, experience)
	return args.Error(0)
}

// UpdateExperience mocks UpdateExperience method
func (m *MockUserRepository) UpdateExperience(ctx context.Context, experience *user.UserExperience) error {
	args := m.Called(ctx, experience)
	return args.Error(0)
}

// DeleteExperience mocks the DeleteExperience method
func (m *MockUserRepository) DeleteExperience(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// AddEducation mocks AddEducation method
func (m *MockUserRepository) AddEducation(ctx context.Context, education *user.UserEducation) error {
	args := m.Called(ctx, education)
	return args.Error(0)
}

// UpdateEducation mocks UpdateEducation method
func (m *MockUserRepository) UpdateEducation(ctx context.Context, education *user.UserEducation) error {
	args := m.Called(ctx, education)
	return args.Error(0)
}

// DeleteEducation mocks the DeleteEducation method
func (m *MockUserRepository) DeleteEducation(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// AddSkill mocks the AddSkill method
func (m *MockUserRepository) AddSkill(ctx context.Context, skill *user.UserSkill) error {
	args := m.Called(ctx, skill)
	return args.Error(0)
}

// CreateSkill mocks the CreateSkill method (deprecated, use AddSkill)
func (m *MockUserRepository) CreateSkill(ctx context.Context, skill *user.UserSkill) error {
	args := m.Called(ctx, skill)
	return args.Error(0)
}

// DeleteSkill mocks the DeleteSkill method
func (m *MockUserRepository) DeleteSkill(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// AddCertification mocks the AddCertification method
func (m *MockUserRepository) AddCertification(ctx context.Context, cert *user.UserCertification) error {
	args := m.Called(ctx, cert)
	return args.Error(0)
}

// AddDocument mocks the AddDocument method
func (m *MockUserRepository) AddDocument(ctx context.Context, doc *user.UserDocument) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

// AddLanguage mocks the AddLanguage method
func (m *MockUserRepository) AddLanguage(ctx context.Context, lang *user.UserLanguage) error {
	args := m.Called(ctx, lang)
	return args.Error(0)
}

// UpdateLanguage mocks the UpdateLanguage method
func (m *MockUserRepository) UpdateLanguage(ctx context.Context, lang *user.UserLanguage) error {
	args := m.Called(ctx, lang)
	return args.Error(0)
}

// DeleteLanguage mocks the DeleteLanguage method
func (m *MockUserRepository) DeleteLanguage(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetLanguagesByUserID mocks the GetLanguagesByUserID method
func (m *MockUserRepository) GetLanguagesByUserID(ctx context.Context, userID int64) ([]user.UserLanguage, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.UserLanguage), args.Error(1)
}

// AddProject mocks the AddProject method
func (m *MockUserRepository) AddProject(ctx context.Context, project *user.UserProject) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

// UpdateProject mocks the UpdateProject method
func (m *MockUserRepository) UpdateProject(ctx context.Context, project *user.UserProject) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

// DeleteProject mocks the DeleteProject method
func (m *MockUserRepository) DeleteProject(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetProjectsByUserID mocks the GetProjectsByUserID method
func (m *MockUserRepository) GetProjectsByUserID(ctx context.Context, userID int64) ([]user.UserProject, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.UserProject), args.Error(1)
}

// UpdateDocument mocks the UpdateDocument method
func (m *MockUserRepository) UpdateDocument(ctx context.Context, doc *user.UserDocument) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

// DeleteDocument mocks the DeleteDocument method
func (m *MockUserRepository) DeleteDocument(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetDocumentsByUserID mocks the GetDocumentsByUserID method
func (m *MockUserRepository) GetDocumentsByUserID(ctx context.Context, userID int64) ([]user.UserDocument, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.UserDocument), args.Error(1)
}

// UpdateCertification mocks the UpdateCertification method
func (m *MockUserRepository) UpdateCertification(ctx context.Context, cert *user.UserCertification) error {
	args := m.Called(ctx, cert)
	return args.Error(0)
}

// DeleteCertification mocks the DeleteCertification method
func (m *MockUserRepository) DeleteCertification(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetCertificationsByUserID mocks the GetCertificationsByUserID method
func (m *MockUserRepository) GetCertificationsByUserID(ctx context.Context, userID int64) ([]user.UserCertification, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.UserCertification), args.Error(1)
}

// UpdateSkill mocks the UpdateSkill method
func (m *MockUserRepository) UpdateSkill(ctx context.Context, skill *user.UserSkill) error {
	args := m.Called(ctx, skill)
	return args.Error(0)
}

// GetSkillsByUserID mocks the GetSkillsByUserID method
func (m *MockUserRepository) GetSkillsByUserID(ctx context.Context, userID int64) ([]user.UserSkill, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.UserSkill), args.Error(1)
}
