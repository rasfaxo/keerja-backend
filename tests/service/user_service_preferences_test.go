package service_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "keerja-backend/internal/domain/user"
    "keerja-backend/internal/service"
    mockRepo "keerja-backend/tests/mocks/repository"
)

func ptrStr(s string) *string { return &s }
func ptrFloat64(f float64) *float64 { return &f }

func TestUserService_UpdatePreferences_CreateIfMissing(t *testing.T) {
    ctx := context.Background()
    userID := int64(1)

    mockUserRepo := new(mockRepo.MockUserRepository)
    // preference not found initially
    mockUserRepo.On("FindPreferenceByUserID", ctx, userID).Return(nil, nil).Once()
    // Expect CreatePreference to be called to create a default row
    mockUserRepo.On("CreatePreference", ctx, mock.MatchedBy(func(p *user.UserPreference) bool {
        return p != nil && p.UserID == userID
    })).Return(nil).Once()
    // Expect UpdatePreference to be called afterwards with updated fields
    mockUserRepo.On("UpdatePreference", ctx, mock.MatchedBy(func(p *user.UserPreference) bool {
        return p.UserID == userID && p.LanguagePreference != nil && *p.LanguagePreference == "en"
    })).Return(nil).Once()

    svc := service.NewUserService(mockUserRepo, nil, nil)

    req := &user.UpdatePreferenceRequest{
        LanguagePreference: ptrStr("en"),
        PreferredSalaryMin: ptrFloat64(500.0),
    }

    err := svc.UpdatePreferences(ctx, userID, req)
    assert.NoError(t, err)
    mockUserRepo.AssertExpectations(t)
}

func TestUserService_UpdatePreferences_UpdateExisting(t *testing.T) {
    ctx := context.Background()
    userID := int64(2)

    existing := &user.UserPreference{ID: 10, UserID: userID}

    mockUserRepo := new(mockRepo.MockUserRepository)
    mockUserRepo.On("FindPreferenceByUserID", ctx, userID).Return(existing, nil).Once()
    mockUserRepo.On("UpdatePreference", ctx, mock.MatchedBy(func(p *user.UserPreference) bool {
        return p.ID == existing.ID && p.UserID == userID && p.LanguagePreference != nil && *p.LanguagePreference == "fr"
    })).Return(nil).Once()

    svc := service.NewUserService(mockUserRepo, nil, nil)

    req := &user.UpdatePreferenceRequest{
        LanguagePreference: ptrStr("fr"),
    }

    err := svc.UpdatePreferences(ctx, userID, req)
    assert.NoError(t, err)
    mockUserRepo.AssertExpectations(t)
}
