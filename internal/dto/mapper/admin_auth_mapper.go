package mapper

import (
	"keerja-backend/internal/domain/admin"
	"keerja-backend/internal/dto/response"
)

// ToAdminAuthResponse converts AdminUser to AdminAuthResponse
func ToAdminAuthResponse(adminUser *admin.AdminUser, accessToken, refreshToken string, expiresIn int64) *response.AdminAuthResponse {
	resp := &response.AdminAuthResponse{
		AdminID:         adminUser.ID,
		UUID:            adminUser.UUID,
		FullName:        adminUser.FullName,
		Email:           adminUser.Email,
		Phone:           adminUser.Phone,
		Status:          adminUser.Status,
		ProfileImageURL: adminUser.ProfileImageURL,
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		TokenType:       "Bearer",
		ExpiresIn:       expiresIn,
		LastLogin:       adminUser.LastLogin,
		CreatedAt:       adminUser.CreatedAt,
	}

	// Map role info if exists
	if adminUser.Role != nil {
		resp.Role = &response.AdminRoleInfo{
			RoleID:          adminUser.Role.ID,
			RoleName:        adminUser.Role.RoleName,
			RoleDescription: adminUser.Role.RoleDescription,
			AccessLevel:     adminUser.Role.AccessLevel,
			IsSystemRole:    adminUser.Role.IsSystemRole,
		}
	}

	return resp
}

// ToAdminProfileResponse converts AdminUser to AdminProfileResponse
func ToAdminProfileResponse(adminUser *admin.AdminUser) *response.AdminProfileResponse {
	resp := &response.AdminProfileResponse{
		AdminID:         adminUser.ID,
		UUID:            adminUser.UUID,
		FullName:        adminUser.FullName,
		Email:           adminUser.Email,
		Phone:           adminUser.Phone,
		Status:          adminUser.Status,
		ProfileImageURL: adminUser.ProfileImageURL,
		LastLogin:       adminUser.LastLogin,
		CreatedAt:       adminUser.CreatedAt,
		UpdatedAt:       adminUser.UpdatedAt,
	}

	// Map role info if exists
	if adminUser.Role != nil {
		resp.Role = &response.AdminRoleInfo{
			RoleID:          adminUser.Role.ID,
			RoleName:        adminUser.Role.RoleName,
			RoleDescription: adminUser.Role.RoleDescription,
			AccessLevel:     adminUser.Role.AccessLevel,
			IsSystemRole:    adminUser.Role.IsSystemRole,
		}
	}

	return resp
}

// ToAdminTokenResponse creates token response
func ToAdminTokenResponse(accessToken, refreshToken string, expiresIn int64) *response.AdminTokenResponse {
	return &response.AdminTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}
}
