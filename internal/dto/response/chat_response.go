package response

import "time"

// ConversationResponse represents a conversation in the response
type ConversationResponse struct {
	ID            int64                 `json:"id"`
	UUID          string                `json:"uuid"`
	Participants  []ParticipantResponse `json:"participants"`
	LastMessage   *MessageResponse      `json:"last_message,omitempty"`
	LastMessageAt *time.Time            `json:"last_message_at,omitempty"`
	UnreadCount   int64                 `json:"unread_count"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

// ParticipantResponse represents a conversation participant
type ParticipantResponse struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	FullName   string    `json:"full_name,omitempty"`
	Email      string    `json:"email,omitempty"`
	UserType   string    `json:"user_type,omitempty"`
	IsArchived bool      `json:"is_archived"`
	CreatedAt  time.Time `json:"created_at"`
}

// MessageResponse represents a message in the response
type MessageResponse struct {
	ID             int64     `json:"id"`
	UUID           string    `json:"uuid"`
	ConversationID int64     `json:"conversation_id"`
	SenderID       int64     `json:"sender_id"`
	Content        string    `json:"content"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ConversationListResponse represents paginated list of conversations
type ConversationListResponse struct {
	Conversations []ConversationResponse `json:"conversations"`
	Total         int64                  `json:"total"`
	Page          int                    `json:"page"`
	Limit         int                    `json:"limit"`
	TotalPages    int                    `json:"total_pages"`
}

// MessageListResponse represents paginated list of messages
type MessageListResponse struct {
	Messages   []MessageResponse `json:"messages"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}
