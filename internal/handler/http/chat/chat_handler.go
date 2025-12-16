package chat

import (
	"strconv"

	"keerja-backend/internal/domain/chat"
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// ChatHandler handles HTTP requests for chat operations
type ChatHandler struct {
	chatService      chat.ChatService
	userRepo         user.UserRepository
	conversationRepo chat.ConversationRepository
}

// NewChatHandler creates a new chat handler instance
func NewChatHandler(
	chatService chat.ChatService,
	userRepo user.UserRepository,
	conversationRepo chat.ConversationRepository,
) *ChatHandler {
	return &ChatHandler{
		chatService:      chatService,
		userRepo:         userRepo,
		conversationRepo: conversationRepo,
	}
}

// CreateConversation handles POST /conversations
func (h *ChatHandler) CreateConversation(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.CreateConversationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Validation failed", err)
	}

	// Create conversation
	conversation, err := h.chatService.CreateConversation(ctx, &chat.CreateConversationRequest{
		InitiatorID: userID,
		RecipientID: req.RecipientID,
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create conversation", err.Error())
	}

	response := mapper.ToConversationResponse(conversation, h.userRepo, ctx, userID)
	return utils.CreatedResponse(c, "Conversation created successfully", response)
}

// GetConversations handles GET /conversations
func (h *ChatHandler) GetConversations(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var filter request.ConversationFilterRequest
	if err := c.QueryParser(&filter); err != nil {
		return utils.BadRequestResponse(c, "Invalid query parameters")
	}

	// Set defaults
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}

	conversations, total, err := h.chatService.GetUserConversations(ctx, userID, filter.Page, filter.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get conversations", err.Error())
	}

	response := mapper.ToConversationListResponse(
		conversations,
		total,
		filter.Page,
		filter.Limit,
		h.userRepo,
		h.conversationRepo,
		ctx,
		userID,
	)

	return utils.SuccessResponse(c, "Conversations retrieved successfully", response)
}

// GetConversationMessages handles GET /conversations/:id/messages
func (h *ChatHandler) GetConversationMessages(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	conversationID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid conversation ID")
	}

	var filter request.MessageFilterRequest
	if err := c.QueryParser(&filter); err != nil {
		return utils.BadRequestResponse(c, "Invalid query parameters")
	}

	// Set defaults
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 50
	}

	messages, total, err := h.chatService.GetConversationMessages(ctx, conversationID, userID, filter.Page, filter.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Failed to get messages", err.Error())
	}

	response := mapper.ToMessageListResponse(messages, total, filter.Page, filter.Limit)
	return utils.SuccessResponse(c, "Messages retrieved successfully", response)
}

// SendMessage handles POST /conversations/:id/messages
func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	conversationID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid conversation ID")
	}

	var req request.SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Validation failed", err)
	}

	// Send message
	message, err := h.chatService.SendMessage(ctx, &chat.SendMessageRequest{
		ConversationID: conversationID,
		SenderID:       userID,
		Content:        req.Content,
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to send message", err.Error())
	}

	response := mapper.ToMessageResponse(message)
	return utils.CreatedResponse(c, "Message sent successfully", response)
}

// MarkAsRead handles PUT /conversations/:id/messages/:msgId/read
func (h *ChatHandler) MarkAsRead(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	messageID, err := strconv.ParseInt(c.Params("msgId"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid message ID")
	}

	if err := h.chatService.MarkAsRead(ctx, messageID, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to mark message as read", err.Error())
	}

	return utils.SuccessResponse(c, "Message marked as read", nil)
}

// ArchiveConversation handles PUT /conversations/:id/archive
func (h *ChatHandler) ArchiveConversation(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	conversationID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid conversation ID")
	}

	if err := h.chatService.ArchiveConversation(ctx, conversationID, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to archive conversation", err.Error())
	}

	return utils.SuccessResponse(c, "Conversation archived successfully", nil)
}
