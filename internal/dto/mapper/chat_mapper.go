package mapper

import (
	"context"
	"math"

	"keerja-backend/internal/domain/chat"
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/response"
)

// ToConversationResponse converts a conversation domain entity to response DTO
func ToConversationResponse(conv *chat.Conversation, userRepo user.UserRepository, ctx context.Context, currentUserID int64) response.ConversationResponse {
	resp := response.ConversationResponse{
		ID:            conv.ID,
		UUID:          conv.UUID.String(),
		LastMessageAt: conv.LastMessageAt,
		CreatedAt:     conv.CreatedAt,
		UpdatedAt:     conv.UpdatedAt,
	}

	// Map participants
	if len(conv.Participants) > 0 {
		resp.Participants = make([]response.ParticipantResponse, len(conv.Participants))
		for i, p := range conv.Participants {
			resp.Participants[i] = response.ParticipantResponse{
				ID:         p.ID,
				UserID:     p.UserID,
				IsArchived: p.IsArchived,
				CreatedAt:  p.CreatedAt,
			}

			// Fetch user details if userRepo is provided
			if userRepo != nil {
				if usr, err := userRepo.FindByID(ctx, p.UserID); err == nil && usr != nil {
					resp.Participants[i].FullName = usr.FullName
					resp.Participants[i].Email = usr.Email
					resp.Participants[i].UserType = usr.UserType
				}
			}
		}
	}

	// Map last message if exists
	if len(conv.Messages) > 0 {
		lastMsg := conv.Messages[0]
		msgResp := ToMessageResponse(&lastMsg)
		resp.LastMessage = &msgResp
	}

	return resp
}

// ToConversationListResponse converts conversations to paginated response
func ToConversationListResponse(
	conversations []chat.Conversation,
	total int64,
	page, limit int,
	userRepo user.UserRepository,
	conversationRepo chat.ConversationRepository,
	ctx context.Context,
	currentUserID int64,
) response.ConversationListResponse {
	items := make([]response.ConversationResponse, len(conversations))

	for i, conv := range conversations {
		items[i] = ToConversationResponse(&conv, userRepo, ctx, currentUserID)

		// Get unread count if repo is provided
		if conversationRepo != nil {
			if unreadCount, err := conversationRepo.GetUnreadCount(ctx, conv.ID, currentUserID); err == nil {
				items[i].UnreadCount = unreadCount
			}
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return response.ConversationListResponse{
		Conversations: items,
		Total:         total,
		Page:          page,
		Limit:         limit,
		TotalPages:    totalPages,
	}
}

// ToMessageResponse converts a message domain entity to response DTO
func ToMessageResponse(msg *chat.Message) response.MessageResponse {
	return response.MessageResponse{
		ID:             msg.ID,
		UUID:           msg.UUID.String(),
		ConversationID: msg.ConversationID,
		SenderID:       msg.SenderID,
		Content:        msg.Content,
		IsRead:         msg.IsRead,
		CreatedAt:      msg.CreatedAt,
		UpdatedAt:      msg.UpdatedAt,
	}
}

// ToMessageListResponse converts messages to paginated response
func ToMessageListResponse(messages []chat.Message, total int64, page, limit int) response.MessageListResponse {
	items := make([]response.MessageResponse, len(messages))
	for i, msg := range messages {
		items[i] = ToMessageResponse(&msg)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return response.MessageListResponse{
		Messages:   items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}
