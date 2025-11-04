package mapper

import (
	"encoding/json"
	"strconv"
	"time"

	"keerja-backend/internal/domain/notification"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
)

// ToPushMessage converts SendPushNotificationRequest to PushMessage entity
func ToPushMessage(req *request.SendPushNotificationRequest) *notification.PushMessage {
	msg := &notification.PushMessage{
		Title: req.Title,
		Body:  req.Body,
		Data:  make(map[string]string),
	}

	// Parse JSON data if provided
	if req.Data != nil && len(req.Data) > 0 {
		var dataMap map[string]interface{}
		if err := json.Unmarshal(req.Data, &dataMap); err == nil {
			// Convert map[string]interface{} to map[string]string
			for key, value := range dataMap {
				if strValue, ok := value.(string); ok {
					msg.Data[key] = strValue
				}
			}
		}
	}

	// Set image URL
	if req.ImageURL != nil {
		msg.ImageURL = *req.ImageURL
	}

	// Set sound (default: "default")
	if req.Sound != nil {
		msg.Sound = *req.Sound
	} else {
		msg.Sound = "default"
	}

	// Set priority (default: "high")
	if req.Priority != nil {
		msg.Priority = *req.Priority
	} else {
		msg.Priority = "high"
	}

	// Set badge
	if req.Badge != nil {
		badge := *req.Badge
		msg.Badge = &badge
	}

	return msg
}

// ToPushMessageFromNotification converts Notification entity to PushMessage
func ToPushMessageFromNotification(notif *notification.Notification) *notification.PushMessage {
	msg := &notification.PushMessage{
		Title:    notif.Title,
		Body:     notif.Message,
		Data:     make(map[string]string),
		Sound:    "default",
		Priority: "high",
	}

	// Add notification metadata to data (convert to string)
	msg.Data["notification_id"] = strconv.FormatInt(notif.ID, 10)
	msg.Data["type"] = notif.Type
	msg.Data["category"] = notif.Category

	if notif.ActionURL != "" {
		msg.Data["action_url"] = notif.ActionURL
	}

	if notif.RelatedID != nil {
		msg.Data["related_id"] = strconv.FormatInt(*notif.RelatedID, 10)
	}

	if notif.RelatedType != "" {
		msg.Data["related_type"] = notif.RelatedType
	}

	// Set priority based on notification priority
	switch notif.Priority {
	case "urgent":
		msg.Priority = "high"
	case "high":
		msg.Priority = "high"
	case "normal":
		msg.Priority = "normal"
	case "low":
		msg.Priority = "low"
	default:
		msg.Priority = "normal"
	}

	// Set image if available
	if notif.Icon != "" {
		msg.ImageURL = notif.Icon
	}

	return msg
}

// ToPushNotificationResponse converts PushResult to response DTO
func ToPushNotificationResponse(result *notification.PushResult) *response.PushNotificationResponse {
	if result == nil {
		return nil
	}

	return &response.PushNotificationResponse{
		MessageID:    result.MessageID,
		Success:      result.Success,
		ErrorCode:    result.ErrorCode,
		ErrorMessage: result.ErrorMessage,
		SentAt:       time.Now(),
	}
}

// ToBatchPushNotificationResponse converts multiple PushResult to batch response
func ToBatchPushNotificationResponse(results []notification.PushResult, totalUsers int) *response.BatchPushNotificationResponse {
	successCount := 0
	failureCount := 0

	responses := make([]response.PushNotificationResponse, len(results))
	for i, result := range results {
		if result.Success {
			successCount++
		} else {
			failureCount++
		}
		responses[i] = *ToPushNotificationResponse(&result)
	}

	return &response.BatchPushNotificationResponse{
		TotalSent:    len(results),
		TotalFailed:  failureCount,
		TotalUsers:   totalUsers,
		SuccessCount: successCount,
		FailureCount: failureCount,
		Results:      responses,
		SentAt:       time.Now(),
	}
}

// ToBatchPushNotificationResponseFromMap converts map of user results to batch response
func ToBatchPushNotificationResponseFromMap(resultsMap map[int64][]notification.PushResult, totalUsers int) *response.BatchPushNotificationResponse {
	successCount := 0
	failureCount := 0
	totalSent := 0

	var responses []response.PushNotificationResponse
	for _, userResults := range resultsMap {
		for _, result := range userResults {
			totalSent++
			if result.Success {
				successCount++
			} else {
				failureCount++
			}
			responses = append(responses, *ToPushNotificationResponse(&result))
		}
	}

	return &response.BatchPushNotificationResponse{
		TotalSent:    totalSent,
		TotalFailed:  failureCount,
		TotalUsers:   totalUsers,
		SuccessCount: successCount,
		FailureCount: failureCount,
		Results:      responses,
		SentAt:       time.Now(),
	}
}

// ToPushValidationResponse converts validation result to response
func ToPushValidationResponse(token string, isValid bool, platform string, message string) *response.PushValidationResponse {
	return &response.PushValidationResponse{
		Token:    token,
		IsValid:  isValid,
		Platform: platform,
		Message:  message,
	}
}

// ToPushDeliveryStatsResponse converts delivery stats to response
func ToPushDeliveryStatsResponse(totalSent, totalDelivered, totalFailed int) *response.PushDeliveryStatsResponse {
	deliveryRate := 0.0
	failureRate := 0.0

	if totalSent > 0 {
		deliveryRate = float64(totalDelivered) / float64(totalSent) * 100
		failureRate = float64(totalFailed) / float64(totalSent) * 100
	}

	return &response.PushDeliveryStatsResponse{
		TotalSent:      totalSent,
		TotalDelivered: totalDelivered,
		TotalFailed:    totalFailed,
		DeliveryRate:   deliveryRate,
		FailureRate:    failureRate,
	}
}

// ToTopicSubscriptionResponse converts topic subscription result to response
func ToTopicSubscriptionResponse(topic string, success bool, tokensCount int, errorMessage string) *response.TopicSubscriptionResponse {
	return &response.TopicSubscriptionResponse{
		Topic:        topic,
		Success:      success,
		TokensCount:  tokensCount,
		ErrorMessage: errorMessage,
	}
}
