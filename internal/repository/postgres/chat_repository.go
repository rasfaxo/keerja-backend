package postgres

import (
	"context"
	"fmt"
	"time"

	"keerja-backend/internal/domain/chat"

	"gorm.io/gorm"
)

// chatRepository implements chat.ConversationRepository and chat.MessageRepository
type chatRepository struct {
	db *gorm.DB
}

// NewChatRepository creates a new chat repository instance
func NewChatRepository(db *gorm.DB) interface {
	chat.ConversationRepository
	chat.MessageRepository
} {
	return &chatRepository{db: db}
}

// ==================== Conversation Repository Methods ====================

// CreateConversation creates a new conversation
func (r *chatRepository) CreateConversation(ctx context.Context, conversation *chat.Conversation) error {
	return r.db.WithContext(ctx).Create(conversation).Error
}

// GetConversationByID retrieves a conversation by its ID
func (r *chatRepository) GetConversationByID(ctx context.Context, id int64) (*chat.Conversation, error) {
	var conversation chat.Conversation
	err := r.db.WithContext(ctx).
		Preload("Participants").
		First(&conversation, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &conversation, nil
}

// GetUserConversations retrieves all conversations for a user with pagination
func (r *chatRepository) GetUserConversations(ctx context.Context, userID int64, page, limit int) ([]chat.Conversation, int64, error) {
	var conversations []chat.Conversation
	var total int64

	// Subquery to get conversation IDs for the user
	subQuery := r.db.WithContext(ctx).
		Model(&chat.ChatParticipant{}).
		Select("conversation_id").
		Where("user_id = ? AND is_archived = ?", userID, false)

	// Main query
	query := r.db.WithContext(ctx).
		Model(&chat.Conversation{}).
		Where("id IN (?)", subQuery).
		Where("deleted_at IS NULL")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (page - 1) * limit
	err := query.
		Preload("Participants").
		Order("COALESCE(last_message_at, created_at) DESC").
		Offset(offset).
		Limit(limit).
		Find(&conversations).Error

	if err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// FindOrCreateConversation finds an existing conversation or creates a new one
func (r *chatRepository) FindOrCreateConversation(ctx context.Context, user1ID, user2ID int64) (*chat.Conversation, error) {
	// Try to find existing conversation with both participants
	var conversation chat.Conversation

	err := r.db.WithContext(ctx).
		Raw(`
			SELECT DISTINCT c.* 
			FROM conversations c
			INNER JOIN chat_participants cp1 ON c.id = cp1.conversation_id
			INNER JOIN chat_participants cp2 ON c.id = cp2.conversation_id
			WHERE cp1.user_id = ? 
			  AND cp2.user_id = ? 
			  AND c.deleted_at IS NULL
			LIMIT 1
		`, user1ID, user2ID).
		Scan(&conversation).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// If conversation exists, load participants
	if conversation.ID != 0 {
		if err := r.db.WithContext(ctx).
			Preload("Participants").
			First(&conversation, conversation.ID).Error; err != nil {
			return nil, err
		}
		return &conversation, nil
	}

	// Create new conversation with participants
	now := time.Now()
	conversation = chat.Conversation{
		CreatedAt: now,
		UpdatedAt: now,
		Participants: []chat.ChatParticipant{
			{UserID: user1ID, CreatedAt: now},
			{UserID: user2ID, CreatedAt: now},
		},
	}

	if err := r.db.WithContext(ctx).Create(&conversation).Error; err != nil {
		return nil, err
	}

	// Reload with participants
	if err := r.db.WithContext(ctx).
		Preload("Participants").
		First(&conversation, conversation.ID).Error; err != nil {
		return nil, err
	}

	return &conversation, nil
}

// ArchiveConversation archives a conversation for a specific user
func (r *chatRepository) ArchiveConversation(ctx context.Context, conversationID, userID int64) error {
	result := r.db.WithContext(ctx).
		Model(&chat.ChatParticipant{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Update("is_archived", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("participant not found")
	}

	return nil
}

// IsUserParticipant checks if a user is a participant in a conversation
func (r *chatRepository) IsUserParticipant(ctx context.Context, conversationID, userID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&chat.ChatParticipant{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetUnreadCount returns the count of unread messages for a user in a conversation
func (r *chatRepository) GetUnreadCount(ctx context.Context, conversationID, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&chat.Message{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = ? AND deleted_at IS NULL", conversationID, userID, false).
		Count(&count).Error

	return count, err
}

// UpdateLastMessageAt updates the last_message_at timestamp
func (r *chatRepository) UpdateLastMessageAt(ctx context.Context, conversationID int64) error {
	return r.db.WithContext(ctx).
		Model(&chat.Conversation{}).
		Where("id = ?", conversationID).
		Update("last_message_at", time.Now()).Error
}

// ==================== Message Repository Methods ====================

// CreateMessage creates a new message
func (r *chatRepository) CreateMessage(ctx context.Context, message *chat.Message) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the message
		if err := tx.Create(message).Error; err != nil {
			return err
		}

		// Update conversation's last_message_at
		if err := tx.Model(&chat.Conversation{}).
			Where("id = ?", message.ConversationID).
			Update("last_message_at", message.CreatedAt).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetConversationMessages retrieves messages for a conversation with pagination
func (r *chatRepository) GetConversationMessages(ctx context.Context, conversationID int64, page, limit int) ([]chat.Message, int64, error) {
	var messages []chat.Message
	var total int64

	query := r.db.WithContext(ctx).
		Model(&chat.Message{}).
		Where("conversation_id = ? AND deleted_at IS NULL", conversationID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination - order by created_at DESC (newest first)
	offset := (page - 1) * limit
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&messages).Error

	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// MarkAsRead marks a message as read
func (r *chatRepository) MarkAsRead(ctx context.Context, messageID, userID int64) error {
	// Verify the user is not the sender before marking as read
	result := r.db.WithContext(ctx).
		Model(&chat.Message{}).
		Where("id = ? AND sender_id != ?", messageID, userID).
		Update("is_read", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("message not found or user is sender")
	}

	return nil
}

// GetMessageByID retrieves a message by its ID
func (r *chatRepository) GetMessageByID(ctx context.Context, id int64) (*chat.Message, error) {
	var message chat.Message
	err := r.db.WithContext(ctx).
		First(&message, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &message, nil
}
