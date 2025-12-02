package request

// GoogleOAuthExchangeRequest represents the payload for the mobile OAuth exchange endpoint.
type GoogleOAuthExchangeRequest struct {
	Code         string `json:"code" validate:"required"`
	CodeVerifier string `json:"code_verifier"`
	State        string `json:"state" validate:"required"`
	RedirectURI  string `json:"redirect_uri"`
}

// OneTimeCodeExchangeRequest is used by mobile clients to exchange a one-time code
// received via deep-link into an application JWT.
type OneTimeCodeExchangeRequest struct {
	Code string `json:"code" validate:"required"`
}
