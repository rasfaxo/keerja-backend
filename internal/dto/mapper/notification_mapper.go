package mapper

import (
	"encoding/json"
	"math"

	"keerja-backend/internal/domain/notification"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
)

// ============================================================================
// Notification Mappers
// ============================================================================

// ToNotificationResponse maps a notification entity to response
func ToNotificationResponse(notif *notification.Notification) response.NotificationResponse {
	var data map[string]interface{}
	if notif.Data != "" {
		json.Unmarshal([]byte(notif.Data), &data)
	}

	return response.NotificationResponse{
		ID:          notif.ID,
		UserID:      notif.UserID,
		Type:        notif.Type,
		Title:       notif.Title,
		Message:     notif.Message,
		Data:        data,
		IsRead:      notif.IsRead,
		ReadAt:      notif.ReadAt,
		Priority:    notif.Priority,
		Category:    notif.Category,
		ActionURL:   notif.ActionURL,
		Icon:        notif.Icon,
		SenderID:    notif.SenderID,
		RelatedID:   notif.RelatedID,
		RelatedType: notif.RelatedType,
		ExpiresAt:   notif.ExpiresAt,
		IsSent:      notif.IsSent,
		SentAt:      notif.SentAt,
		Channel:     notif.Channel,
		CreatedAt:   notif.CreatedAt,
		UpdatedAt:   notif.UpdatedAt,
	}
}

// ToNotificationListResponse maps notifications to paginated response
func ToNotificationListResponse(notifs []notification.Notification, total int64, page, limit int) response.NotificationListResponse {
	notifications := make([]response.NotificationResponse, len(notifs))
	for i, notif := range notifs {
		notifications[i] = ToNotificationResponse(&notif)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return response.NotificationListResponse{
		Notifications: notifications,
		Total:         total,
		Page:          page,
		Limit:         limit,
		TotalPages:    totalPages,
	}
}

// ToNotificationPreferenceResponse maps notification preferences to response
func ToNotificationPreferenceResponse(pref *notification.NotificationPreference) response.NotificationPreferenceResponse {
	return response.NotificationPreferenceResponse{
		ID:                        pref.ID,
		UserID:                    pref.UserID,
		EmailEnabled:              pref.EmailEnabled,
		PushEnabled:               pref.PushEnabled,
		SMSEnabled:                pref.SMSEnabled,
		JobApplicationsEnabled:    pref.JobApplicationsEnabled,
		InterviewEnabled:          pref.InterviewEnabled,
		StatusUpdatesEnabled:      pref.StatusUpdatesEnabled,
		JobRecommendationsEnabled: pref.JobRecommendationsEnabled,
		CompanyUpdatesEnabled:     pref.CompanyUpdatesEnabled,
		MarketingEnabled:          pref.MarketingEnabled,
		WeeklyDigestEnabled:       pref.WeeklyDigestEnabled,
		CreatedAt:                 pref.CreatedAt,
		UpdatedAt:                 pref.UpdatedAt,
	}
}

// ToNotificationStatsResponse maps notification stats to response
func ToNotificationStatsResponse(stats *notification.NotificationStats) response.NotificationStatsResponse {
	return response.NotificationStatsResponse{
		TotalCount:        stats.TotalCount,
		UnreadCount:       stats.UnreadCount,
		ReadCount:         stats.ReadCount,
		TodayCount:        stats.TodayCount,
		ThisWeekCount:     stats.ThisWeekCount,
		HighPriorityCount: stats.HighPriorityCount,
		CategoryBreakdown: stats.CategoryBreakdown,
	}
}

// ============================================================================
// Request to Domain Mappers
// ============================================================================

// ToSendNotificationRequest converts request DTO to domain request
func ToSendNotificationRequest(req *request.SendNotificationRequest) *notification.SendNotificationRequest {
	return &notification.SendNotificationRequest{
		UserID:      req.UserID,
		Type:        req.Type,
		Title:       req.Title,
		Message:     req.Message,
		Data:        req.Data,
		Priority:    req.Priority,
		Category:    req.Category,
		ActionURL:   req.ActionURL,
		Icon:        req.Icon,
		SenderID:    req.SenderID,
		RelatedID:   req.RelatedID,
		RelatedType: req.RelatedType,
		ExpiresAt:   req.ExpiresAt,
		Channel:     req.Channel,
	}
}

// ToNotificationFilter converts request filter to domain filter
func ToNotificationFilter(req *request.NotificationFilterRequest) notification.NotificationFilter {
	return notification.NotificationFilter{
		Type:     req.Type,
		Category: req.Category,
		IsRead:   req.IsRead,
		Priority: req.Priority,
		DateFrom: req.DateFrom,
		DateTo:   req.DateTo,
	}
}

// ToNotificationPreferenceFromRequest updates notification preference entity from request
func ToNotificationPreferenceFromRequest(req *request.UpdateNotificationPreferencesRequest, existing *notification.NotificationPreference) *notification.NotificationPreference {
	if existing == nil {
		existing = &notification.NotificationPreference{}
	}

	if req.EmailEnabled != nil {
		existing.EmailEnabled = *req.EmailEnabled
	}
	if req.PushEnabled != nil {
		existing.PushEnabled = *req.PushEnabled
	}
	if req.SMSEnabled != nil {
		existing.SMSEnabled = *req.SMSEnabled
	}
	if req.JobApplicationsEnabled != nil {
		existing.JobApplicationsEnabled = *req.JobApplicationsEnabled
	}
	if req.InterviewEnabled != nil {
		existing.InterviewEnabled = *req.InterviewEnabled
	}
	if req.StatusUpdatesEnabled != nil {
		existing.StatusUpdatesEnabled = *req.StatusUpdatesEnabled
	}
	if req.JobRecommendationsEnabled != nil {
		existing.JobRecommendationsEnabled = *req.JobRecommendationsEnabled
	}
	if req.CompanyUpdatesEnabled != nil {
		existing.CompanyUpdatesEnabled = *req.CompanyUpdatesEnabled
	}
	if req.MarketingEnabled != nil {
		existing.MarketingEnabled = *req.MarketingEnabled
	}
	if req.WeeklyDigestEnabled != nil {
		existing.WeeklyDigestEnabled = *req.WeeklyDigestEnabled
	}

	return existing
}
