package mapper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/mapper"
)

func ptr(s string) *string { return &s }

func TestToAuthResponse_IncludesUserPreferences(t *testing.T) {
	pref := &user.UserPreference{
		ID:                 1,
		LanguagePreference: ptr("en"),
		ThemePreference:    "dark",
		PreferredJobType:   "remote",
	}

	u := &user.User{
		ID:         123,
		FullName:   "Tester",
		Email:      "tester@example.com",
		UserType:   "jobseeker",
		IsVerified: true,
		Status:     "active",
		Preference: pref,
	}

	resp := mapper.ToAuthResponse(u, "tok", "")
	assert.NotNil(t, resp)
	if assert.NotNil(t, resp.User) {
		assert.NotNil(t, resp.User.Preference)
		assert.Equal(t, "en", *resp.User.Preference.LanguagePreference)
	}
}
