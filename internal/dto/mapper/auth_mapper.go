package mapper

import (
	"keerja-backend/internal/domain/auth"
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/response"
)

// ToAuthResponse converts User entity to AuthResponse DTO
func ToAuthResponse(u *user.User, accessToken string, refreshToken string) *response.AuthResponse {
	if u == nil {
		return nil
	}

	return &response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour in seconds
		User:         ToUserBasic(u),
		Company:      nil, // Will be set by handler if user is employer
	}
}

// ToAuthResponseWithCompany converts User entity and Company to AuthResponse DTO (for employer)
func ToAuthResponseWithCompany(u *user.User, c *company.Company, accessToken string, refreshToken string) *response.AuthResponse {
	if u == nil {
		return nil
	}

	authResp := &response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour in seconds
		User:         ToUserBasic(u),
		Company:      ToCompanyBasic(c),
	}

	return authResp
}

// ToUserBasic converts User entity to UserBasic DTO
func ToUserBasic(u *user.User) *response.UserBasic {
	if u == nil {
		return nil
	}

	phone := ""
	if u.Phone != nil {
		phone = *u.Phone
	}

	return &response.UserBasic{
		ID:         u.ID,
		UUID:       u.UUID.String(),
		FullName:   u.FullName,
		Email:      u.Email,
		Phone:      phone,
		UserType:   u.UserType,
		IsVerified: u.IsVerified,
		Status:     u.Status,
	}
}

// ToCompanyBasic converts Company entity to CompanyBasic DTO
func ToCompanyBasic(c *company.Company) *response.CompanyBasic {
	if c == nil {
		return nil
	}

	logoURL := ""
	if c.LogoURL != nil {
		logoURL = *c.LogoURL
	}

	// Default verification values
	status := "not_requested"
	badgeGranted := false
	npwpNumber := ""
	nibNumber := ""

	// If verification record exists, get details
	if c.Verification != nil {
		status = c.Verification.Status
		badgeGranted = c.Verification.BadgeGranted
		npwpNumber = c.Verification.NPWPNumber
		if c.Verification.NIBNumber != nil {
			nibNumber = *c.Verification.NIBNumber
		}
	}

	return &response.CompanyBasic{
		ID:           c.ID,
		UUID:         c.UUID.String(),
		CompanyName:  c.CompanyName,
		Slug:         c.Slug,
		LogoURL:      logoURL,
		IsVerified:   c.Verified,
		Status:       status,
		BadgeGranted: badgeGranted,
		NPWPNumber:   npwpNumber,
		NIBNumber:    nibNumber,
	}
}

// ToTokenResponse converts tokens to TokenResponse DTO
func ToTokenResponse(accessToken string, refreshToken string) *response.TokenResponse {
	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}
}

// ToDeviceListResponse converts refresh tokens to DeviceListResponse
func ToDeviceListResponse(tokens []*auth.RefreshToken) *response.DeviceListResponse {
	devices := make([]response.DeviceInfo, 0, len(tokens))

	for _, token := range tokens {
		device := response.DeviceInfo{
			ID:         token.ID,
			DeviceName: token.DeviceName,
			DeviceType: token.DeviceType,
			IPAddress:  token.IPAddress,
			LastUsedAt: token.LastUsedAt.Format("2006-01-02T15:04:05Z07:00"),
			CreatedAt:  token.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			IsCurrent:  false, // Can be set by handler based on current request
		}
		devices = append(devices, device)
	}

	return &response.DeviceListResponse{
		Devices: devices,
		Total:   len(devices),
	}
}

// ToOAuthProviderResponse converts OAuthProvider to response DTO
func ToOAuthProviderResponse(provider *auth.OAuthProvider) *response.OAuthProviderResponse {
	if provider == nil {
		return nil
	}

	// Note: OAuthProvider doesn't have LastLoginAt in entity
	// Using UpdatedAt as approximation
	var lastLoginAt *string
	if !provider.UpdatedAt.IsZero() && !provider.UpdatedAt.Equal(provider.CreatedAt) {
		t := provider.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
		lastLoginAt = &t
	}

	return &response.OAuthProviderResponse{
		ID:          provider.ID,
		Provider:    provider.Provider,
		ProviderID:  provider.ProviderUserID, // Using ProviderUserID from entity
		Email:       provider.Email,
		DisplayName: provider.Name, // Using Name from entity
		ConnectedAt: provider.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		LastLoginAt: lastLoginAt,
	}
}

// ToOAuthProviderListResponse converts slice of OAuthProvider to list response
func ToOAuthProviderListResponse(providers []*auth.OAuthProvider) []response.OAuthProviderResponse {
	result := make([]response.OAuthProviderResponse, 0, len(providers))

	for _, provider := range providers {
		if resp := ToOAuthProviderResponse(provider); resp != nil {
			result = append(result, *resp)
		}
	}

	return result
}
