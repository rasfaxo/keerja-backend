package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"keerja-backend/internal/domain/chat"
	"keerja-backend/internal/domain/user"
)

// chatService implements chat.ChatService interface
type chatService struct {
	conversationRepo chat.ConversationRepository
	messageRepo      chat.MessageRepository
	userRepo         user.UserRepository
	wsHub            WebSocketHub // Interface for WebSocket hub
}

// WebSocketHub interface for broadcasting messages
type WebSocketHub interface {
	BroadcastToConversation(conversationID int64, message interface{})
}

// NewChatService creates a new chat service instance
func NewChatService(
	conversationRepo chat.ConversationRepository,
	messageRepo chat.MessageRepository,
	userRepo user.UserRepository,
	wsHub WebSocketHub,
) chat.ChatService {
	return &chatService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		userRepo:         userRepo,
		wsHub:            wsHub,
	}
}

// CreateConversation creates a new conversation between two users
func (s *chatService) CreateConversation(ctx context.Context, req *chat.CreateConversationRequest) (*chat.Conversation, error) {
	// Validate users exist and get their roles
	initiator, err := s.userRepo.FindByID(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find initiator: %w", err)
	}
	if initiator == nil {
		return nil, fmt.Errorf("initiator not found")
	}

	recipient, err := s.userRepo.FindByID(ctx, req.RecipientID)
	if err != nil {
		return nil, fmt.Errorf("failed to find recipient: %w", err)
	}
	if recipient == nil {
		return nil, fmt.Errorf("recipient not found")
	}

	// Validate participant roles
	if !chat.IsValidParticipantPair(initiator.UserType, recipient.UserType) {
		return nil, fmt.Errorf("invalid participant combination: %s cannot chat with %s", initiator.UserType, recipient.UserType)
	}

	// Find or create conversation
	conversation, err := s.conversationRepo.FindOrCreateConversation(ctx, req.InitiatorID, req.RecipientID)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create conversation: %w", err)
	}

	return conversation, nil
}

// GetUserConversations retrieves all conversations for a user with pagination
func (s *chatService) GetUserConversations(ctx context.Context, userID int64, page, limit int) ([]chat.Conversation, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	conversations, total, err := s.conversationRepo.GetUserConversations(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get conversations: %w", err)
	}

	return conversations, total, nil
}

// GetConversationByID retrieves a conversation by its ID
func (s *chatService) GetConversationByID(ctx context.Context, conversationID, userID int64) (*chat.Conversation, error) {
	// Check if user is participant
	isParticipant, err := s.conversationRepo.IsUserParticipant(ctx, conversationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check participant: %w", err)
	}
	if !isParticipant {
		return nil, fmt.Errorf("user is not a participant in this conversation")
	}

	conversation, err := s.conversationRepo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	if conversation == nil {
		return nil, fmt.Errorf("conversation not found")
	}

	return conversation, nil
}

// SendMessage sends a message in a conversation
func (s *chatService) SendMessage(ctx context.Context, req *chat.SendMessageRequest) (*chat.Message, error) {
	// Validate conversation exists and user is participant
	isParticipant, err := s.conversationRepo.IsUserParticipant(ctx, req.ConversationID, req.SenderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check participant: %w", err)
	}
	if !isParticipant {
		return nil, fmt.Errorf("user is not a participant in this conversation")
	}

	// Validate content
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, fmt.Errorf("message content cannot be empty")
	}
	if len(content) > 5000 {
		return nil, fmt.Errorf("message content exceeds maximum length of 5000 characters")
	}

	// Create message
	now := time.Now()
	message := &chat.Message{
		ConversationID: req.ConversationID,
		SenderID:       req.SenderID,
		Content:        content,
		IsRead:         false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.messageRepo.CreateMessage(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Broadcast to WebSocket if hub is available
	if s.wsHub != nil {
		s.wsHub.BroadcastToConversation(req.ConversationID, map[string]interface{}{
			"type":      "message_received",
			"data":      message,
			"timestamp": time.Now(),
		})
	}

	return message, nil
}

// GetConversationMessages retrieves messages for a conversation with pagination
func (s *chatService) GetConversationMessages(ctx context.Context, conversationID, userID int64, page, limit int) ([]chat.Message, int64, error) {
	// Check if user is participant
	isParticipant, err := s.conversationRepo.IsUserParticipant(ctx, conversationID, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to check participant: %w", err)
	}
	if !isParticipant {
		return nil, 0, fmt.Errorf("user is not a participant in this conversation")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	messages, total, err := s.messageRepo.GetConversationMessages(ctx, conversationID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get messages: %w", err)
	}

	return messages, total, nil
}

// MarkAsRead marks a message as read
func (s *chatService) MarkAsRead(ctx context.Context, messageID, userID int64) error {
	// Get message to verify it exists and get conversation ID
	message, err := s.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}
	if message == nil {
		return fmt.Errorf("message not found")
	}

	// Check if user is participant in the conversation
	isParticipant, err := s.conversationRepo.IsUserParticipant(ctx, message.ConversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to check participant: %w", err)
	}
	if !isParticipant {
		return fmt.Errorf("user is not a participant in this conversation")
	}

	// Prevent sender from marking their own message as read
	if message.SenderID == userID {
		return fmt.Errorf("cannot mark own message as read")
	}

	if err := s.messageRepo.MarkAsRead(ctx, messageID, userID); err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	// Broadcast read status to WebSocket if hub is available
	if s.wsHub != nil {
		s.wsHub.BroadcastToConversation(message.ConversationID, map[string]interface{}{
			"type": "message_read",
			"data": map[string]interface{}{
				"message_id": messageID,
				"user_id":    userID,
			},
			"timestamp": time.Now(),
		})
	}

	return nil
}

// ArchiveConversation archives a conversation for a user
func (s *chatService) ArchiveConversation(ctx context.Context, conversationID, userID int64) error {
	// Check if user is participant
	isParticipant, err := s.conversationRepo.IsUserParticipant(ctx, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to check participant: %w", err)
	}
	if !isParticipant {
		return fmt.Errorf("user is not a participant in this conversation")
	}

	if err := s.conversationRepo.ArchiveConversation(ctx, conversationID, userID); err != nil {
		return fmt.Errorf("failed to archive conversation: %w", err)
	}

	return nil
}
