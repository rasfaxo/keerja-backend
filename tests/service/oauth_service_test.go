package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"keerja-backend/internal/service"
)

func TestGetGoogleAuthURL_WithPKCEStateSaved(t *testing.T) {
	stateStore := service.NewInMemoryOAuthStateStore()
	svc := service.NewOAuthService(
		nil,
		nil,
		service.OAuthConfig{
			ClientID:    "test-client",
			RedirectURI: "http://localhost/callback",
		},
		"secret",
		time.Hour,
		stateStore,
		[]string{"myapp://oauth-callback"},
	)

	resp, err := svc.GetGoogleAuthURL(context.Background(), service.GoogleAuthURLRequest{
		ClientType:           "mobile",
		RedirectURI:          "myapp://oauth-callback",
		PostLoginRedirectURI: "myapp://oauth-callback",
		CodeChallenge:        "abc123",
		CodeChallengeMethod:  "S256",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.AuthURL)
	require.NotEmpty(t, resp.State)

	stateData, err := stateStore.Consume(context.Background(), resp.State)
	require.NoError(t, err)
	require.Equal(t, "mobile", stateData.ClientType)
	require.Equal(t, "myapp://oauth-callback", stateData.RedirectURI)
	require.Equal(t, "abc123", stateData.CodeChallenge)
	require.Equal(t, "S256", stateData.CodeChallengeMethod)
}

func TestGetGoogleAuthURL_DisallowsUnknownMobileRedirect(t *testing.T) {
	stateStore := service.NewInMemoryOAuthStateStore()
	svc := service.NewOAuthService(
		nil,
		nil,
		service.OAuthConfig{
			ClientID:    "test-client",
			RedirectURI: "http://localhost/callback",
		},
		"secret",
		time.Hour,
		stateStore,
		[]string{"myapp://oauth-callback"},
	)

	_, err := svc.GetGoogleAuthURL(context.Background(), service.GoogleAuthURLRequest{
		ClientType:  "mobile",
		RedirectURI: "otherapp://callback",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "redirect_uri")
}

func TestExchangeGoogleCode_RequiresCodeVerifierWhenChallengeSet(t *testing.T) {
	stateStore := service.NewInMemoryOAuthStateStore()
	svc := service.NewOAuthService(
		nil,
		nil,
		service.OAuthConfig{
			ClientID:    "test-client",
			RedirectURI: "http://localhost/callback",
		},
		"secret",
		time.Hour,
		stateStore,
		[]string{"myapp://oauth-callback"},
	)

	// Create state with code challenge
	state := "teststate123"
	_ = stateStore.Save(context.Background(), state, service.OAuthStateData{
		RedirectURI:         "myapp://oauth-callback",
		CodeChallenge:       "abc123",
		CodeChallengeMethod: "S256",
		ClientType:          "mobile",
		CreatedAt:           time.Now(),
	}, time.Minute*5)

	_, err := svc.ExchangeGoogleCode(context.Background(), service.GoogleExchangeRequest{
		Code:  "fakecode",
		State: state,
		// no code_verifier provided -> should error
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "code_verifier")
}
