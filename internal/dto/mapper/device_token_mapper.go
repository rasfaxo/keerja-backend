package mapper

import (
	"encoding/json"

	"keerja-backend/internal/domain/notification"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
)

// ToDeviceToken converts RegisterDeviceTokenRequest to DeviceToken entity
func ToDeviceToken(userID int64, req *request.RegisterDeviceTokenRequest) *notification.DeviceToken {
	var deviceInfo notification.DeviceInfo

	// Parse device info JSON if provided
	if req.DeviceInfo != nil && len(req.DeviceInfo) > 0 {
		_ = json.Unmarshal(req.DeviceInfo, &deviceInfo)
	}

	return &notification.DeviceToken{
		UserID:     userID,
		Token:      req.Token,
		Platform:   notification.Platform(req.Platform),
		DeviceInfo: deviceInfo,
		IsActive:   true,
	}
}

// UpdateDeviceTokenFromRequest updates DeviceToken entity from UpdateDeviceTokenRequest
func UpdateDeviceTokenFromRequest(token *notification.DeviceToken, req *request.UpdateDeviceTokenRequest) {
	if req.DeviceInfo != nil && len(req.DeviceInfo) > 0 {
		var deviceInfo notification.DeviceInfo
		if err := json.Unmarshal(req.DeviceInfo, &deviceInfo); err == nil {
			token.DeviceInfo = deviceInfo
		}
	}
}

// ToDeviceTokenResponse converts DeviceToken entity to response DTO
func ToDeviceTokenResponse(token *notification.DeviceToken) *response.DeviceTokenResponse {
	if token == nil {
		return nil
	}

	// Marshal DeviceInfo to JSON
	deviceInfoJSON, _ := json.Marshal(token.DeviceInfo)

	// Convert FailureReason string to pointer
	var failureReason *string
	if token.FailureReason != "" {
		failureReason = &token.FailureReason
	}

	return &response.DeviceTokenResponse{
		ID:            token.ID,
		UserID:        token.UserID,
		Token:         token.Token,
		Platform:      string(token.Platform),
		DeviceInfo:    deviceInfoJSON,
		IsActive:      token.IsActive,
		LastUsedAt:    token.LastUsedAt,
		FailureCount:  token.FailureCount,
		LastFailureAt: token.LastFailureAt,
		FailureReason: failureReason,
		CreatedAt:     token.CreatedAt,
		UpdatedAt:     token.UpdatedAt,
	}
}

// ToDeviceTokenListResponse converts slice of DeviceToken entities to list response
func ToDeviceTokenListResponse(tokens []notification.DeviceToken, total int64, page, pageSize int) *response.DeviceTokenListResponse {
	data := make([]response.DeviceTokenResponse, len(tokens))
	for i, token := range tokens {
		data[i] = *ToDeviceTokenResponse(&token)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &response.DeviceTokenListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// ToDeviceTokenStatsResponse converts token counts to statistics response
func ToDeviceTokenStatsResponse(total, active, inactive, android, ios, web int) *response.DeviceTokenStatsResponse {
	return &response.DeviceTokenStatsResponse{
		TotalTokens:    total,
		ActiveTokens:   active,
		InactiveTokens: inactive,
		AndroidTokens:  android,
		IOSTokens:      ios,
		WebTokens:      web,
	}
}

// ToDeviceTokenCountResponse converts user token count to response
func ToDeviceTokenCountResponse(userID int64, total, active int) *response.DeviceTokenCountResponse {
	return &response.DeviceTokenCountResponse{
		UserID:       userID,
		TotalTokens:  total,
		ActiveTokens: active,
	}
}
