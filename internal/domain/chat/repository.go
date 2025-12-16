package chat

import "context"

// ConversationRepository defines the interface for conversation data access
type ConversationRepository interface {
	// CreateConversation creates a new conversation
	CreateConversation(ctx context.Context, conversation *Conversation) error

	// GetConversationByID retrieves a conversation by its ID
	GetConversationByID(ctx context.Context, id int64) (*Conversation, error)

	// GetUserConversations retrieves all conversations for a user with pagination
	GetUserConversations(ctx context.Context, userID int64, page, limit int) ([]Conversation, int64, error)

	// FindOrCreateConversation finds an existing conversation or creates a new one
	FindOrCreateConversation(ctx context.Context, user1ID, user2ID int64) (*Conversation, error)

	// ArchiveConversation archives a conversation for a specific user
	ArchiveConversation(ctx context.Context, conversationID, userID int64) error

	// IsUserParticipant checks if a user is a participant in a conversation
	IsUserParticipant(ctx context.Context, conversationID, userID int64) (bool, error)

	// GetUnreadCount returns the count of unread messages for a user in a conversation
	GetUnreadCount(ctx context.Context, conversationID, userID int64) (int64, error)

	// UpdateLastMessageAt updates the last_message_at timestamp
	UpdateLastMessageAt(ctx context.Context, conversationID int64) error
}

// MessageRepository defines the interface for message data access
type MessageRepository interface {
	// CreateMessage creates a new message
	CreateMessage(ctx context.Context, message *Message) error

	// GetConversationMessages retrieves messages for a conversation with pagination
	GetConversationMessages(ctx context.Context, conversationID int64, page, limit int) ([]Message, int64, error)

	// MarkAsRead marks a message as read
	MarkAsRead(ctx context.Context, messageID, userID int64) error

	// GetMessageByID retrieves a message by its ID
	GetMessageByID(ctx context.Context, id int64) (*Message, error)
}
