package routes

import (
	"keerja-backend/internal/handler/websocket"

	"github.com/gofiber/fiber/v2"
)

// SetupWebSocketRoutes sets up WebSocket routes for chat
func SetupWebSocketRoutes(app *fiber.App, wsHandler *websocket.Handler) {
	// WebSocket endpoint - no auth middleware (handled in handler via query param token)
	app.Get("/api/v1/chat/ws/:conversationId", wsHandler.HandleWebSocket)
}
