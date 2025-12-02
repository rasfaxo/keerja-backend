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
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"keerja-backend/internal/domain/auth"
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/utils"
)

// OAuth Configuration
const (
	GoogleTokenURL    = "https://oauth2.googleapis.com/token"
	GoogleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	StateTokenLength  = 32
	defaultStateTTL   = 5 * time.Minute
	oneTimeCodeTTL    = 2 * time.Minute
)

var (
	ErrInvalidProvider       = errors.New("invalid OAuth provider")
	ErrInvalidState          = errors.New("invalid state token")
	ErrInvalidRedirectURI    = errors.New("redirect_uri is not allowed")
	ErrMissingCodeVerifier   = errors.New("code_verifier is required for PKCE exchange")
	ErrOAuthExchangeFailed   = errors.New("failed to exchange OAuth code")
	ErrOAuthUserInfoFailed   = errors.New("failed to get user info from OAuth provider")
	errMobileRedirectMissing = errors.New("mobile redirect URIs not configured")
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

// GoogleAuthURLRequest captures query options for initiating OAuth.
type GoogleAuthURLRequest struct {
	ClientType           string
	RedirectURI          string
	PostLoginRedirectURI string
	CodeChallenge        string
	CodeChallengeMethod  string
}

// GoogleAuthURLResponse returns the generated auth URL and state metadata.
type GoogleAuthURLResponse struct {
	AuthURL   string `json:"auth_url"`
	State     string `json:"state"`
	ExpiresIn int    `json:"expires_in"`
}

// GoogleExchangeRequest represents payload for PKCE mobile exchange.
type GoogleExchangeRequest struct {
	Code         string `json:"code"`
	CodeVerifier string `json:"code_verifier"`
	State        string `json:"state"`
	RedirectURI  string `json:"redirect_uri"`
}

// GoogleExchangeResponse represents the resulting JWT tokens for mobile apps.
type GoogleExchangeResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// OAuthCallbackResult describes the outcome of browser callback handling.
type OAuthCallbackResult struct {
	AccessToken          string
	ClientType           string
	PostLoginRedirectURI string
}

// OAuthService handles OAuth authentication business logic
type OAuthService struct {
	oauthRepo              auth.OAuthRepository
	userRepo               user.UserRepository
	googleConfig           OAuthConfig
	jwtSecret              string
	jwtDuration            time.Duration
	stateStore             OAuthStateStore
	stateTTL               time.Duration
	allowedMobileRedirects map[string]struct{}

	// fallback in-memory one-time-code store (used when Redis not available)
	oneTimeMu    sync.RWMutex
	oneTimeStore map[string]oneTimeEntry
}

type oneTimeEntry struct {
	Token     string
	ExpiresAt time.Time
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(
	oauthRepo auth.OAuthRepository,
	userRepo user.UserRepository,
	googleConfig OAuthConfig,
	jwtSecret string,
	jwtDuration time.Duration,
	stateStore OAuthStateStore,
	allowedMobileRedirects []string,
) *OAuthService {
	if stateStore == nil {
		stateStore = NewInMemoryOAuthStateStore()
	}

	return &OAuthService{
		oauthRepo: oauthRepo,

		userRepo:               userRepo,
		googleConfig:           googleConfig,
		jwtSecret:              jwtSecret,
		jwtDuration:            jwtDuration,
		stateStore:             stateStore,
		stateTTL:               defaultStateTTL,
		allowedMobileRedirects: normalizeRedirects(allowedMobileRedirects),
	}
}

func normalizeRedirects(values []string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, v := range values {
		if trimmed := strings.TrimSpace(v); trimmed != "" {
			result[trimmed] = struct{}{}
		}
	}
	return result
}

// generateStateToken generates a random state token for OAuth
func (s *OAuthService) generateStateToken() (string, error) {
	b := make([]byte, StateTokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *OAuthService) storeState(ctx context.Context, state string, data OAuthStateData) error {
	if s.stateStore == nil {
		return errors.New("state store not configured")
	}
	return s.stateStore.Save(ctx, state, data, s.stateTTL)
}

func (s *OAuthService) consumeState(ctx context.Context, state string) (*OAuthStateData, error) {
	if s.stateStore == nil {
		return nil, errors.New("state store not configured")
	}

	data, err := s.stateStore.Consume(ctx, state)
	if err != nil {
		if errors.Is(err, ErrStateNotFound) {
			return nil, ErrInvalidState
		}
		return nil, err
	}
	return data, nil
}

func (s *OAuthService) ensureMobileRedirectAllowed(uri string) error {
	if strings.TrimSpace(uri) == "" {
		return fmt.Errorf("%w: redirect_uri is required", ErrInvalidRedirectURI)
	}
	if len(s.allowedMobileRedirects) == 0 {
		return errMobileRedirectMissing
	}
	if !s.isAllowedMobileRedirect(uri) {
		return fmt.Errorf("%w: %s", ErrInvalidRedirectURI, uri)
	}
	return nil
}

func (s *OAuthService) isAllowedMobileRedirect(uri string) bool {
	_, ok := s.allowedMobileRedirects[uri]
	return ok
}

// CreateOneTimeCode stores a single-use code in Redis mapping to the provided jwtToken.
// It returns the generated code. Requires stateStore to be backed by Redis.
func (s *OAuthService) CreateOneTimeCode(ctx context.Context, jwtToken string, ttl time.Duration) (string, error) {
	// try redis-backed store first
	if redisStore, ok := s.stateStore.(*RedisOAuthStateStore); ok && redisStore.client != nil {
		// generate a random code (url-safe)
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			return "", fmt.Errorf("failed to generate one-time code: %w", err)
		}
		code := base64.RawURLEncoding.EncodeToString(b)

		key := "oauth:onetime:" + code
		if ttl <= 0 {
			ttl = oneTimeCodeTTL
		}

		if err := redisStore.client.Set(ctx, key, jwtToken, ttl).Err(); err != nil {
			return "", fmt.Errorf("failed to store one-time code in redis: %w", err)
		}

		return code, nil
	}

	// fallback: use in-memory one-time store
	if ttl <= 0 {
		ttl = oneTimeCodeTTL
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate one-time code: %w", err)
	}
	code := base64.RawURLEncoding.EncodeToString(b)

	s.oneTimeMu.Lock()
	if s.oneTimeStore == nil {
		s.oneTimeStore = make(map[string]oneTimeEntry)
	}
	s.oneTimeStore[code] = oneTimeEntry{Token: jwtToken, ExpiresAt: time.Now().Add(ttl)}
	s.oneTimeMu.Unlock()

	return code, nil
}

// ConsumeOneTimeCode atomically gets and deletes the one-time code from Redis and returns the stored jwtToken.
func (s *OAuthService) ConsumeOneTimeCode(ctx context.Context, code string) (string, error) {
	// Redis-backed store
	if redisStore, ok := s.stateStore.(*RedisOAuthStateStore); ok && redisStore.client != nil {
		key := "oauth:onetime:" + code
		value, err := redisStore.client.Get(ctx, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return "", errors.New("one-time code not found or expired")
			}
			return "", fmt.Errorf("failed to get one-time code: %w", err)
		}

		// delete key after reading
		if err := redisStore.client.Del(ctx, key).Err(); err != nil {
			return "", fmt.Errorf("failed to delete one-time code: %w", err)
		}

		return value, nil
	}

	// fallback in-memory store
	s.oneTimeMu.Lock()
	defer s.oneTimeMu.Unlock()

	entry, ok := s.oneTimeStore[code]
	if !ok || time.Now().After(entry.ExpiresAt) {
		delete(s.oneTimeStore, code)
		return "", errors.New("one-time code not found or expired")
	}

	// consume and delete
	delete(s.oneTimeStore, code)
	return entry.Token, nil
}

// JWTExpirySeconds returns configured JWT expiry in seconds (fallback to 3600 if not set)
func (s *OAuthService) JWTExpirySeconds() int {
	secs := int(s.jwtDuration.Seconds())
	if secs <= 0 {
		secs = 3600
	}
	return secs
}

func (s *OAuthService) normalizeClientType(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "mobile":
		return "mobile"
	default:
		return "web"
	}
}

// GetGoogleAuthURL generates Google OAuth authorization URL
func (s *OAuthService) GetGoogleAuthURL(ctx context.Context, req GoogleAuthURLRequest) (*GoogleAuthURLResponse, error) {
	clientType := s.normalizeClientType(req.ClientType)

	redirectURI := strings.TrimSpace(req.RedirectURI)
	if redirectURI == "" {
		redirectURI = s.googleConfig.RedirectURI
	}

	if clientType == "mobile" {
		if err := s.ensureMobileRedirectAllowed(redirectURI); err != nil {
			return nil, err
		}
	}

	postLoginRedirect := strings.TrimSpace(req.PostLoginRedirectURI)
	if postLoginRedirect != "" {
		if err := s.ensureMobileRedirectAllowed(postLoginRedirect); err != nil {
			return nil, fmt.Errorf("post_login_redirect_uri invalid: %w", err)
		}
	}

	codeChallenge := strings.TrimSpace(req.CodeChallenge)
	codeChallengeMethod := strings.TrimSpace(strings.ToUpper(req.CodeChallengeMethod))
	if codeChallenge != "" && codeChallengeMethod == "" {
		codeChallengeMethod = "S256"
	}

	state, err := s.generateStateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate state token: %w", err)
	}

	stateData := OAuthStateData{
		RedirectURI:          redirectURI,
		PostLoginRedirectURI: postLoginRedirect,
		CodeChallenge:        codeChallenge,
		CodeChallengeMethod:  codeChallengeMethod,
		ClientType:           clientType,
		CreatedAt:            time.Now(),
	}

	if err := s.storeState(ctx, state, stateData); err != nil {
		return nil, fmt.Errorf("failed to persist oauth state: %w", err)
	}

	params := url.Values{
		"client_id":     {s.googleConfig.ClientID},
		"redirect_uri":  {redirectURI},
		"response_type": {"code"},
		"scope":         {"openid email profile"},
		"state":         {state},
		"access_type":   {"offline"}, // Get refresh token
		"prompt":        {"consent"},
	}

	if codeChallenge != "" {
		params.Set("code_challenge", codeChallenge)
		params.Set("code_challenge_method", codeChallengeMethod)
	}

	authURL := "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()

	return &GoogleAuthURLResponse{
		AuthURL:   authURL,
		State:     state,
		ExpiresIn: int(s.stateTTL.Seconds()),
	}, nil
}

// exchangeGoogleCode exchanges authorization code for access token
func (s *OAuthService) exchangeGoogleCode(ctx context.Context, code, redirectURI, codeVerifier string) (*GoogleTokenResponse, error) {
	data := url.Values{
		"client_id":     {s.googleConfig.ClientID},
		"client_secret": {s.googleConfig.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {redirectURI},
	}

	if codeVerifier != "" {
		data.Set("code_verifier", codeVerifier)
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
func (s *OAuthService) HandleGoogleCallback(ctx context.Context, code, state string) (*OAuthCallbackResult, error) {
	if strings.TrimSpace(state) == "" {
		return nil, ErrInvalidState
	}

	stateData, err := s.consumeState(ctx, state)
	if err != nil {
		return nil, err
	}

	redirectURI := stateData.RedirectURI
	if redirectURI == "" {
		redirectURI = s.googleConfig.RedirectURI
	}

	tokenResp, err := s.exchangeGoogleCode(ctx, code, redirectURI, "")
	if err != nil {
		return nil, fmt.Errorf("exchange code failed: %w", err)
	}

	jwtToken, err := s.finalizeGoogleLogin(ctx, tokenResp)
	if err != nil {
		return nil, err
	}

	return &OAuthCallbackResult{
		AccessToken:          jwtToken,
		ClientType:           stateData.ClientType,
		PostLoginRedirectURI: stateData.PostLoginRedirectURI,
	}, nil
}

// ExchangeGoogleCode handles PKCE exchange for mobile clients.
func (s *OAuthService) ExchangeGoogleCode(ctx context.Context, req GoogleExchangeRequest) (*GoogleExchangeResponse, error) {
	if strings.TrimSpace(req.Code) == "" || strings.TrimSpace(req.State) == "" {
		return nil, fmt.Errorf("code and state are required")
	}

	stateData, err := s.consumeState(ctx, req.State)
	if err != nil {
		return nil, err
	}

	if stateData.ClientType != "mobile" {
		return nil, errors.New("state was not created for mobile flow")
	}

	if stateData.CodeChallenge != "" && strings.TrimSpace(req.CodeVerifier) == "" {
		return nil, ErrMissingCodeVerifier
	}

	redirectURI := stateData.RedirectURI
	if redirectURI == "" {
		redirectURI = s.googleConfig.RedirectURI
	}

	if strings.TrimSpace(req.RedirectURI) != "" && req.RedirectURI != redirectURI {
		return nil, fmt.Errorf("redirect_uri mismatch")
	}

	tokenResp, err := s.exchangeGoogleCode(ctx, req.Code, redirectURI, strings.TrimSpace(req.CodeVerifier))
	if err != nil {
		return nil, err
	}

	jwtToken, err := s.finalizeGoogleLogin(ctx, tokenResp)
	if err != nil {
		return nil, err
	}

	expiresIn := int(s.jwtDuration.Seconds())
	if expiresIn == 0 {
		expiresIn = int((time.Hour).Seconds())
	}

	return &GoogleExchangeResponse{
		AccessToken: jwtToken,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	}, nil
}

func (s *OAuthService) finalizeGoogleLogin(ctx context.Context, tokenResp *GoogleTokenResponse) (string, error) {
	userInfo, err := s.getGoogleUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		return "", fmt.Errorf("get user info failed: %w", err)
	}

	usr, err := s.upsertGoogleUser(ctx, userInfo, tokenResp)
	if err != nil {
		return "", err
	}

	now := time.Now()
	usr.LastLogin = &now
	_ = s.userRepo.Update(ctx, usr)

	jwtToken, err := utils.GenerateAccessToken(int64(usr.ID), usr.Email, usr.UserType, s.jwtSecret, s.jwtDuration)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return jwtToken, nil
}

func (s *OAuthService) upsertGoogleUser(ctx context.Context, userInfo *GoogleUserInfo, tokenResp *GoogleTokenResponse) (*user.User, error) {
	provider, err := s.oauthRepo.FindByProviderAndUserID(ctx, "google", userInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find OAuth provider: %w", err)
	}

	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	accessToken := tokenResp.AccessToken

	if provider != nil {
		usr, err := s.userRepo.FindByID(ctx, provider.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to find user: %w", err)
		}

		provider.AccessToken = &accessToken
		if tokenResp.RefreshToken != "" {
			refreshToken := tokenResp.RefreshToken
			provider.RefreshToken = &refreshToken
		}
		provider.TokenExpiry = &expiresAt

		if err := s.oauthRepo.Update(ctx, provider); err != nil {
			return nil, fmt.Errorf("failed to update OAuth provider: %w", err)
		}
		return usr, nil
	}

	usr, err := s.userRepo.FindByEmail(ctx, userInfo.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if usr == nil {
		usr = &user.User{
			Email:        userInfo.Email,
			FullName:     userInfo.Name,
			UserType:     "jobseeker",
			IsVerified:   userInfo.VerifiedEmail,
			Status:       "active",
			PasswordHash: "",
		}

		if err := s.userRepo.Create(ctx, usr); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	rawDataJSON, _ := json.Marshal(userInfo)
	rawDataStr := string(rawDataJSON)
	email := userInfo.Email
	name := userInfo.Name
	avatar := userInfo.Picture

	provider = &auth.OAuthProvider{
		UserID:         usr.ID,
		Provider:       "google",
		ProviderUserID: userInfo.ID,
		Email:          &email,
		Name:           &name,
		AvatarURL:      &avatar,
		AccessToken:    &accessToken,
		TokenExpiry:    &expiresAt,
		RawData:        &rawDataStr,
	}

	if tokenResp.RefreshToken != "" {
		refreshToken := tokenResp.RefreshToken
		provider.RefreshToken = &refreshToken
	}

	if err := s.oauthRepo.Create(ctx, provider); err != nil {
		return nil, fmt.Errorf("failed to create OAuth provider: %w", err)
	}

	// TODO: Implement scheduled refresh flow using stored refresh tokens when expiry is reached.

	return usr, nil
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

// RefreshGoogleAccessToken refreshes an expired Google access token using a stored refresh_token.
// If successful it updates the OAuthProvider record with the new access token and expiry.
// This is a helper used by background jobs or on-demand refresh flows.
func (s *OAuthService) RefreshGoogleAccessToken(ctx context.Context, provider *auth.OAuthProvider) (*GoogleTokenResponse, error) {
	if provider == nil || provider.RefreshToken == nil || *provider.RefreshToken == "" {
		return nil, errors.New("no refresh token available")
	}

	data := url.Values{
		"client_id":     {s.googleConfig.ClientID},
		"client_secret": {s.googleConfig.ClientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {*provider.RefreshToken},
	}

	resp, err := http.PostForm(GoogleTokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("refresh token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("refresh request failed: %s", string(body))
	}

	var tokenResp GoogleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode refresh response: %w", err)
	}

	// update provider
	now := time.Now()
	expiresAt := now.Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	provider.AccessToken = &tokenResp.AccessToken
	provider.TokenExpiry = &expiresAt

	if err := s.oauthRepo.Update(ctx, provider); err != nil {
		return nil, fmt.Errorf("failed to update provider after refresh: %w", err)
	}

	return &tokenResp, nil
}

// GetConnectedProviders returns list of connected OAuth providers for user
func (s *OAuthService) GetConnectedProviders(ctx context.Context, userID uint) ([]auth.OAuthProvider, error) {
	providers, err := s.oauthRepo.FindByUserID(ctx, int64(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get connected providers: %w", err)
	}

	return providers, nil
}
