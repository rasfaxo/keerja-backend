package routes

import (
	"keerja-backend/internal/handler/http/chat"
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupChatRoutes sets up routes for chat operations
func SetupChatRoutes(api fiber.Router, chatHandler *chat.ChatHandler, authMw *middleware.AuthMiddleware) {
	chatGroup := api.Group("/chat")
	chatGroup.Use(authMw.AuthRequired())

	// Conversation routes
	chatGroup.Post("/conversations", chatHandler.CreateConversation)
	chatGroup.Get("/conversations", chatHandler.GetConversations)
	chatGroup.Get("/conversations/:id/messages", chatHandler.GetConversationMessages)

	// Message routes
	chatGroup.Post("/conversations/:id/messages", chatHandler.SendMessage)
	chatGroup.Put("/conversations/:id/messages/:msgId/read", chatHandler.MarkAsRead)

	// Archive route
	chatGroup.Put("/conversations/:id/archive", chatHandler.ArchiveConversation)
}
