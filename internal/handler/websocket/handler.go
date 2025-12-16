package websocket

import (
	"strconv"
	"time"

	"keerja-backend/internal/config"
	"keerja-backend/internal/domain/chat"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog/log"
)

// Handler handles WebSocket connections
type Handler struct {
	hub              *Hub
	conversationRepo chat.ConversationRepository
	config           *config.Config
}

// NewHandler creates a new WebSocket handler
func NewHandler(
	hub *Hub,
	conversationRepo chat.ConversationRepository,
	config *config.Config,
) *Handler {
	return &Handler{
		hub:              hub,
		conversationRepo: conversationRepo,
		config:           config,
	}
}

// HandleWebSocket handles WebSocket upgrade and connection
func (h *Handler) HandleWebSocket(c *fiber.Ctx) error {
	// Check if it's a WebSocket upgrade request
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	// Extract conversation ID from params
	conversationID, err := strconv.ParseInt(c.Params("conversationId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid conversation ID",
		})
	}

	// Extract token from query parameter
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Authentication token required",
		})
	}

	// Validate token
	claims, err := utils.ValidateToken(token, h.config.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid or expired token",
		})
	}

	userID := claims.UserID

	// Verify user is participant in the conversation
	ctx := c.Context()
	isParticipant, err := h.conversationRepo.IsUserParticipant(ctx, conversationID, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check participant")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to verify participant",
		})
	}
	if !isParticipant {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "You are not a participant in this conversation",
		})
	}

	// Upgrade to WebSocket
	return websocket.New(func(conn *websocket.Conn) {
		// Create client
		client := NewClient(h.hub, conn, conversationID, userID)

		// Register client
		h.hub.register <- client

		// Send welcome message
		client.send <- map[string]interface{}{
			"type": "connected",
			"data": map[string]interface{}{
				"conversation_id": conversationID,
				"user_id":         userID,
			},
			"timestamp": time.Now(),
		}

		// Start client read/write pumps
		client.Start()
	})(c)
}
