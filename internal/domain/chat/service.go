package chat

import "context"

// CreateConversationRequest represents the request to create a conversation
type CreateConversationRequest struct {
	InitiatorID int64
	RecipientID int64
}

// SendMessageRequest represents the request to send a message
type SendMessageRequest struct {
	ConversationID int64
	SenderID       int64
	Content        string
}

// ChatService defines the interface for chat business logic
type ChatService interface {
	// CreateConversation creates a new conversation between two users
	CreateConversation(ctx context.Context, req *CreateConversationRequest) (*Conversation, error)

	// GetUserConversations retrieves all conversations for a user with pagination
	GetUserConversations(ctx context.Context, userID int64, page, limit int) ([]Conversation, int64, error)

	// GetConversationByID retrieves a conversation by its ID
	GetConversationByID(ctx context.Context, conversationID, userID int64) (*Conversation, error)

	// SendMessage sends a message in a conversation
	SendMessage(ctx context.Context, req *SendMessageRequest) (*Message, error)

	// GetConversationMessages retrieves messages for a conversation with pagination
	GetConversationMessages(ctx context.Context, conversationID, userID int64, page, limit int) ([]Message, int64, error)

	// MarkAsRead marks a message as read
	MarkAsRead(ctx context.Context, messageID, userID int64) error

	// ArchiveConversation archives a conversation for a user
	ArchiveConversation(ctx context.Context, conversationID, userID int64) error
}
