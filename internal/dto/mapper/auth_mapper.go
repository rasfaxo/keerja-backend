package mapper

import (
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
	}
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

// ToTokenResponse converts tokens to TokenResponse DTO
func ToTokenResponse(accessToken string, refreshToken string) *response.TokenResponse {
	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}
}
