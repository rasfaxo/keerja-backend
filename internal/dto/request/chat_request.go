package request

// CreateConversationRequest represents the request to create a conversation
type CreateConversationRequest struct {
	RecipientID int64 `json:"recipient_id" validate:"required,min=1"`
}

// SendMessageRequest represents the request to send a message
type SendMessageRequest struct {
	Content string `json:"content" validate:"required,max=5000"`
}

// MessageFilterRequest represents the filter for fetching messages
type MessageFilterRequest struct {
	Page  int `json:"page" query:"page" validate:"omitempty,min=1"`
	Limit int `json:"limit" query:"limit" validate:"omitempty,min=1,max=100"`
}

// ConversationFilterRequest represents the filter for fetching conversations
type ConversationFilterRequest struct {
	Page  int `json:"page" query:"page" validate:"omitempty,min=1"`
	Limit int `json:"limit" query:"limit" validate:"omitempty,min=1,max=100"`
}
