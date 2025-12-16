package chat

import (
	"time"

	"github.com/google/uuid"
)

// Conversation represents a chat conversation between two users
type Conversation struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex" json:"uuid"`
	LastMessageAt *time.Time `gorm:"type:timestamp" json:"last_message_at,omitempty"`
	CreatedAt     time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"type:timestamp;index" json:"deleted_at,omitempty"`

	// Relationships
	Participants []ChatParticipant `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE" json:"participants,omitempty"`
	Messages     []Message         `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE" json:"messages,omitempty"`
}

// TableName specifies the table name for Conversation
func (Conversation) TableName() string {
	return "conversations"
}

// ChatParticipant represents a user participating in a conversation
type ChatParticipant struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ConversationID int64     `gorm:"not null;index:idx_conversation_user,unique" json:"conversation_id"`
	UserID         int64     `gorm:"not null;index:idx_conversation_user,unique;index:idx_participant_user" json:"user_id"`
	IsArchived     bool      `gorm:"default:false" json:"is_archived"`
	CreatedAt      time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`

	// Relationships
	Conversation *Conversation `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
}

// TableName specifies the table name for ChatParticipant
func (ChatParticipant) TableName() string {
	return "chat_participants"
}

// Message represents a chat message in a conversation
type Message struct {
	ID             int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex" json:"uuid"`
	ConversationID int64      `gorm:"not null;index:idx_message_conversation" json:"conversation_id"`
	SenderID       int64      `gorm:"not null;index:idx_message_sender" json:"sender_id"`
	Content        string     `gorm:"type:varchar(5000);not null" json:"content"`
	IsRead         bool       `gorm:"default:false" json:"is_read"`
	CreatedAt      time.Time  `gorm:"type:timestamp;default:now();index:idx_message_created" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"type:timestamp;index" json:"deleted_at,omitempty"`

	// Relationships
	Conversation *Conversation `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
}

// TableName specifies the table name for Message
func (Message) TableName() string {
	return "messages"
}

// IsValidParticipantPair validates if two user roles are allowed to chat
func IsValidParticipantPair(role1, role2 string) bool {
	validPairs := map[string][]string{
		"candidate": {"employer", "recruiter", "company"},
		"employer":  {"candidate"},
		"recruiter": {"candidate"},
		"company":   {"candidate"},
	}

	allowed, exists := validPairs[role1]
	if !exists {
		return false
	}

	for _, allowedRole := range allowed {
		if allowedRole == role2 {
			return true
		}
	}

	return false
}
