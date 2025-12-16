package websocket

import (
	"sync"
	"time"
)

// Hub manages WebSocket connections and broadcasts
type Hub struct {
	// Registered clients per conversation
	// Key: conversationID, Value: map of clients
	conversations map[int64]map[*Client]bool

	// Client ID cache for duplicate prevention (5 min TTL)
	clientIDCache map[string]time.Time

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast messages to specific conversation
	broadcast chan *BroadcastMessage

	// Mutex for thread-safe operations
	mu sync.RWMutex

	// Cleanup interval for expired client IDs
	cleanupInterval time.Duration
}

// BroadcastMessage represents a message to broadcast to a conversation
type BroadcastMessage struct {
	ConversationID int64
	Data           interface{}
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		conversations:   make(map[int64]map[*Client]bool),
		clientIDCache:   make(map[string]time.Time),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		broadcast:       make(chan *BroadcastMessage, 256),
		cleanupInterval: 1 * time.Minute,
	}
}

// Run starts the hub's main event loop
func (h *Hub) Run() {
	// Start cleanup goroutine
	go h.cleanupExpiredClientIDs()

	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastToConversation(message)
		}
	}
}

// registerClient adds a client to the hub
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.conversations[client.conversationID] == nil {
		h.conversations[client.conversationID] = make(map[*Client]bool)
	}
	h.conversations[client.conversationID][client] = true
}

// unregisterClient removes a client from the hub
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.conversations[client.conversationID]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.send)

			// Remove conversation map if empty
			if len(clients) == 0 {
				delete(h.conversations, client.conversationID)
			}
		}
	}
}

// broadcastToConversation sends a message to all clients in a conversation
func (h *Hub) broadcastToConversation(message *BroadcastMessage) {
	h.mu.RLock()
	clients := h.conversations[message.ConversationID]
	h.mu.RUnlock()

	for client := range clients {
		select {
		case client.send <- message.Data:
		default:
			// Client's send channel is full, close and unregister
			h.unregisterClient(client)
		}
	}
}

// BroadcastToConversation queues a message for broadcast to a conversation
func (h *Hub) BroadcastToConversation(conversationID int64, data interface{}) {
	h.broadcast <- &BroadcastMessage{
		ConversationID: conversationID,
		Data:           data,
	}
}

// IsClientIDProcessed checks if a client ID has been processed recently
func (h *Hub) IsClientIDProcessed(clientID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if expireTime, exists := h.clientIDCache[clientID]; exists {
		// Check if still valid (within 5 minutes)
		if time.Now().Before(expireTime) {
			return true
		}
	}
	return false
}

// MarkClientIDProcessed marks a client ID as processed
func (h *Hub) MarkClientIDProcessed(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Set expiration to 5 minutes from now
	h.clientIDCache[clientID] = time.Now().Add(5 * time.Minute)
}

// cleanupExpiredClientIDs removes expired client IDs from cache
func (h *Hub) cleanupExpiredClientIDs() {
	ticker := time.NewTicker(h.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		h.mu.Lock()
		now := time.Now()
		for clientID, expireTime := range h.clientIDCache {
			if now.After(expireTime) {
				delete(h.clientIDCache, clientID)
			}
		}
		h.mu.Unlock()
	}
}

// GetConnectionCount returns the number of active connections for a conversation
func (h *Hub) GetConnectionCount(conversationID int64) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.conversations[conversationID]; ok {
		return len(clients)
	}
	return 0
}
