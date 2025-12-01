package authhandler

import (
	"fmt"
	"net/url"

	"keerja-backend/internal/domain/auth"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/service"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) InitiateGoogleLogin(c *fiber.Ctx) error {
	ctx := c.Context()

	req := service.GoogleAuthURLRequest{
		ClientType:           c.Query("client"),
		RedirectURI:          c.Query("redirect_uri"),
		PostLoginRedirectURI: c.Query("post_login_redirect_uri"),
		CodeChallenge:        c.Query("code_challenge"),
		CodeChallengeMethod:  c.Query("code_challenge_method"),
	}

	authResp, err := h.oauthService.GetGoogleAuthURL(ctx, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate auth URL", err.Error())
	}

	return utils.SuccessResponse(c, "Google auth URL generated", fiber.Map{
		"auth_url":   authResp.AuthURL,
		"state":      authResp.State,
		"expires_in": authResp.ExpiresIn,
	})
}

func (h *AuthHandler) HandleGoogleCallback(c *fiber.Ctx) error {
	ctx := c.Context()

	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Missing authorization code", "")
	}

	result, err := h.oauthService.HandleGoogleCallback(ctx, code, state)
	if err != nil {
		if err.Error() == "invalid state" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid OAuth state", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to authenticate with Google", err.Error())
	}

	if result.PostLoginRedirectURI != "" {
		if code, err := h.oauthService.CreateOneTimeCode(ctx, result.AccessToken, 0); err == nil {
			deepLink := fmt.Sprintf("%s?code=%s", result.PostLoginRedirectURI, url.QueryEscape(code))
			return c.Redirect(deepLink, fiber.StatusFound)
		}

		deepLink := result.PostLoginRedirectURI
		fragment := url.QueryEscape(result.AccessToken)
		if fragment != "" {
			deepLink = fmt.Sprintf("%s#token=%s", result.PostLoginRedirectURI, fragment)
		}
		return c.Redirect(deepLink, fiber.StatusFound)
	}

	response := mapper.ToTokenResponse(result.AccessToken, "")

	return utils.SuccessResponse(c, "Google authentication successful", response)
}

func (h *AuthHandler) ExchangeGoogleOAuthCode(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.GoogleOAuthExchangeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	exchangeResp, err := h.oauthService.ExchangeGoogleCode(ctx, service.GoogleExchangeRequest{
		Code:         req.Code,
		CodeVerifier: req.CodeVerifier,
		State:        req.State,
		RedirectURI:  req.RedirectURI,
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to exchange authorization code", err.Error())
	}

	return utils.SuccessResponse(c, "Google authentication successful", exchangeResp)
}

func (h *AuthHandler) ExchangeOneTimeCode(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.OneTimeCodeExchangeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	jwtToken, err := h.oauthService.ConsumeOneTimeCode(ctx, req.Code)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid or expired one-time code", err.Error())
	}

	expiresIn := h.oauthService.JWTExpirySeconds()

	return utils.SuccessResponse(c, "Token exchange successful", map[string]interface{}{
		"access_token": jwtToken,
		"token_type":   "Bearer",
		"expires_in":   expiresIn,
	})
}

func (h *AuthHandler) GetConnectedProviders(c *fiber.Ctx) error {
	ctx := c.Context()

	claims := c.Locals("user")
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "No user context found")
	}

	userClaims := claims.(*utils.Claims)

	providers, err := h.oauthService.GetConnectedProviders(ctx, uint(userClaims.UserID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get connected providers", err.Error())
	}

	providerPointers := make([]*auth.OAuthProvider, len(providers))
	for i := range providers {
		providerPointers[i] = &providers[i]
	}

	response := mapper.ToOAuthProviderListResponse(providerPointers)

	return utils.SuccessResponse(c, "Connected providers retrieved successfully", response)
}

func (h *AuthHandler) DisconnectOAuth(c *fiber.Ctx) error {
	ctx := c.Context()

	claims := c.Locals("user")
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "No user context found")
	}

	userClaims := claims.(*utils.Claims)

	provider := c.Params("provider")
	if provider == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Provider name required", "")
	}

	if err := h.oauthService.DisconnectOAuthProvider(ctx, uint(userClaims.UserID), provider); err != nil {
		if err.Error() == "OAuth provider not connected" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Provider not connected", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to disconnect provider", err.Error())
	}

	return utils.SuccessResponse(c, "OAuth provider disconnected successfully", nil)
}
