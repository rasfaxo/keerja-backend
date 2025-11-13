package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"keerja-backend/internal/domain/auth"
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/utils"
)

// OAuth Configuration
const (
	GoogleTokenURL    = "https://oauth2.googleapis.com/token"
	GoogleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	StateTokenLength  = 32
)

var (
	ErrInvalidProvider     = errors.New("invalid OAuth provider")
	ErrInvalidState        = errors.New("invalid state token")
	ErrOAuthExchangeFailed = errors.New("failed to exchange OAuth code")
	ErrOAuthUserInfoFailed = errors.New("failed to get user info from OAuth provider")
)

// OAuthConfig holds OAuth provider configuration
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
}

// GoogleUserInfo represents user data from Google OAuth
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// GoogleTokenResponse represents Google OAuth token response
type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
}

// OAuthService handles OAuth authentication business logic
type OAuthService struct {
	oauthRepo    auth.OAuthRepository
	userRepo     user.UserRepository
	googleConfig OAuthConfig
	jwtSecret    string
	jwtDuration  time.Duration
	stateStore   map[string]time.Time // In-memory state storage (should use Redis in production)
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(
	oauthRepo auth.OAuthRepository,
	userRepo user.UserRepository,
	googleConfig OAuthConfig,
	jwtSecret string,
	jwtDuration time.Duration,
) *OAuthService {
	return &OAuthService{
		oauthRepo:    oauthRepo,
		userRepo:     userRepo,
		googleConfig: googleConfig,
		jwtSecret:    jwtSecret,
		jwtDuration:  jwtDuration,
		stateStore:   make(map[string]time.Time),
	}
}

// generateStateToken generates a random state token for OAuth
func (s *OAuthService) generateStateToken() (string, error) {
	b := make([]byte, StateTokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	state := base64.URLEncoding.EncodeToString(b)

	// Store state with expiration (5 minutes)
	s.stateStore[state] = time.Now().Add(5 * time.Minute)

	return state, nil
}

// validateState checks if state token is valid and not expired
func (s *OAuthService) validateState(state string) bool {
	expiry, exists := s.stateStore[state]
	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		delete(s.stateStore, state)
		return false
	}

	delete(s.stateStore, state) // Use once
	return true
}

// GetGoogleAuthURL generates Google OAuth authorization URL
func (s *OAuthService) GetGoogleAuthURL(ctx context.Context) (string, error) {
	state, err := s.generateStateToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate state token: %w", err)
	}

	params := url.Values{
		"client_id":     {s.googleConfig.ClientID},
		"redirect_uri":  {s.googleConfig.RedirectURI},
		"response_type": {"code"},
		"scope":         {"openid email profile"},
		"state":         {state},
		"access_type":   {"offline"}, // Get refresh token
		"prompt":        {"consent"},
	}

	authURL := "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()
	return authURL, nil
}

// exchangeGoogleCode exchanges authorization code for access token
func (s *OAuthService) exchangeGoogleCode(ctx context.Context, code string) (*GoogleTokenResponse, error) {
	data := url.Values{
		"client_id":     {s.googleConfig.ClientID},
		"client_secret": {s.googleConfig.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {s.googleConfig.RedirectURI},
	}

	resp, err := http.PostForm(GoogleTokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokenResp GoogleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// getGoogleUserInfo fetches user info from Google
func (s *OAuthService) getGoogleUserInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", GoogleUserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s", string(body))
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

// HandleGoogleCallback handles Google OAuth callback
func (s *OAuthService) HandleGoogleCallback(ctx context.Context, code, state string) (string, error) {
	// Validate state
	if !s.validateState(state) {
		return "", ErrInvalidState
	}

	// Exchange code for token
	tokenResp, err := s.exchangeGoogleCode(ctx, code)
	if err != nil {
		return "", fmt.Errorf("exchange code failed: %w", err)
	}

	// Get user info
	userInfo, err := s.getGoogleUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		return "", fmt.Errorf("get user info failed: %w", err)
	}

	// Check if OAuth provider exists
	provider, err := s.oauthRepo.FindByProviderAndUserID(ctx, "google", userInfo.ID)
	if err != nil {
		return "", fmt.Errorf("failed to find OAuth provider: %w", err)
	}

	var usr *user.User

	if provider != nil {
		// Existing OAuth connection - find user
		usr, err = s.userRepo.FindByID(ctx, provider.UserID)
		if err != nil {
			return "", fmt.Errorf("failed to find user: %w", err)
		}

		// Update tokens
		accessToken := tokenResp.AccessToken
		refreshToken := tokenResp.RefreshToken
		expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

		provider.AccessToken = &accessToken
		provider.RefreshToken = &refreshToken
		provider.TokenExpiry = &expiresAt

		if err := s.oauthRepo.Update(ctx, provider); err != nil {
			return "", fmt.Errorf("failed to update OAuth provider: %w", err)
		}
	} else {
		// New OAuth connection - find or create user by email
		usr, err = s.userRepo.FindByEmail(ctx, userInfo.Email)
		if err != nil {
			return "", fmt.Errorf("failed to find user by email: %w", err)
		}

		if usr == nil {
			// Create new user
			usr = &user.User{
				Email:        userInfo.Email,
				FullName:     userInfo.Name,
				UserType:     "jobseeker", // Default type
				IsVerified:   userInfo.VerifiedEmail,
				Status:       "active",
				PasswordHash: "", // OAuth users don't need password
			}

			if err := s.userRepo.Create(ctx, usr); err != nil {
				return "", fmt.Errorf("failed to create user: %w", err)
			}
		}

		// Create OAuth provider record
		rawDataJSON, _ := json.Marshal(userInfo)
		rawDataStr := string(rawDataJSON)
		accessToken := tokenResp.AccessToken
		refreshToken := tokenResp.RefreshToken
		email := userInfo.Email
		name := userInfo.Name
		avatar := userInfo.Picture
		expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

		provider = &auth.OAuthProvider{
			UserID:         usr.ID,
			Provider:       "google",
			ProviderUserID: userInfo.ID,
			Email:          &email,
			Name:           &name,
			AvatarURL:      &avatar,
			AccessToken:    &accessToken,
			RefreshToken:   &refreshToken,
			TokenExpiry:    &expiresAt,
			RawData:        &rawDataStr,
		}

		if err := s.oauthRepo.Create(ctx, provider); err != nil {
			return "", fmt.Errorf("failed to create OAuth provider: %w", err)
		}
	}

	// Update last login
	now := time.Now()
	usr.LastLogin = &now
	_ = s.userRepo.Update(ctx, usr)

	// Generate JWT token
	jwtToken, err := utils.GenerateAccessToken(int64(usr.ID), usr.Email, usr.UserType, s.jwtSecret, s.jwtDuration)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return jwtToken, nil
}

// DisconnectOAuthProvider disconnects an OAuth provider from user
func (s *OAuthService) DisconnectOAuthProvider(ctx context.Context, userID uint, provider string) error {
	// Check if provider exists
	oauthProvider, err := s.oauthRepo.FindByUserAndProvider(ctx, int64(userID), provider)
	if err != nil {
		return fmt.Errorf("failed to find OAuth provider: %w", err)
	}

	if oauthProvider == nil {
		return errors.New("OAuth provider not connected")
	}

	// Delete provider
	if err := s.oauthRepo.Delete(ctx, oauthProvider.ID); err != nil {
		return fmt.Errorf("failed to delete OAuth provider: %w", err)
	}

	return nil
}

// GetConnectedProviders returns list of connected OAuth providers for user
func (s *OAuthService) GetConnectedProviders(ctx context.Context, userID uint) ([]auth.OAuthProvider, error) {
	providers, err := s.oauthRepo.FindByUserID(ctx, int64(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get connected providers: %w", err)
	}

	return providers, nil
}
