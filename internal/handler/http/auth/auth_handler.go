package authhandler

import (
	"context"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/service"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService         *service.AuthService
	oauthService        *service.OAuthService
	registrationService *service.RegistrationService
	refreshTokenService *service.RefreshTokenService
	userRepo            user.UserRepository
	companyRepo         company.CompanyRepository
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(
	authService *service.AuthService,
	oauthService *service.OAuthService,
	registrationService *service.RegistrationService,
	refreshTokenService *service.RefreshTokenService,
	userRepo user.UserRepository,
	companyRepo company.CompanyRepository,
) *AuthHandler {
	return &AuthHandler{
		authService:         authService,
		oauthService:        oauthService,
		registrationService: registrationService,
		refreshTokenService: refreshTokenService,
		userRepo:            userRepo,
		companyRepo:         companyRepo,
	}
}

func (h *AuthHandler) buildAuthResponse(ctx context.Context, usr *user.User, accessToken, refreshToken string) *response.AuthResponse {
	if usr.UserType == "employer" {
		companies, err := h.companyRepo.GetCompaniesByUserID(ctx, usr.ID)
		if err == nil && len(companies) > 0 {
			return mapper.ToAuthResponseWithCompany(usr, &companies[0], accessToken, refreshToken)
		}
	}
	return mapper.ToAuthResponse(usr, accessToken, refreshToken)
}
