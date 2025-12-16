package websocket

import (
	"encoding/json"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog/log"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 8192
)

// Client represents a WebSocket client connection
type Client struct {
	hub            *Hub
	conn           *websocket.Conn
	send           chan interface{}
	conversationID int64
	userID         int64
}

// NewClient creates a new WebSocket client
func NewClient(hub *Hub, conn *websocket.Conn, conversationID, userID int64) *Client {
	return &Client{
		hub:            hub,
		conn:           conn,
		send:           make(chan interface{}, 256),
		conversationID: conversationID,
		userID:         userID,
	}
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var message map[string]interface{}
		err := c.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Err(err).Msg("WebSocket read error")
			}
			break
		}

		// Handle incoming message
		c.handleMessage(message)
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				log.Error().Err(err).Msg("WebSocket write error")
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming WebSocket messages
func (c *Client) handleMessage(message map[string]interface{}) {
	msgType, ok := message["type"].(string)
	if !ok {
		c.sendError("Invalid message format: missing type")
		return
	}

	switch msgType {
	case "ping":
		// Respond to ping
		c.send <- map[string]interface{}{
			"type":      "pong",
			"timestamp": time.Now(),
		}

	case "typing_start":
		// Broadcast typing indicator
		c.hub.BroadcastToConversation(c.conversationID, map[string]interface{}{
			"type": "user_typing",
			"data": map[string]interface{}{
				"user_id":         c.userID,
				"conversation_id": c.conversationID,
				"is_typing":       true,
			},
			"timestamp": time.Now(),
		})

	case "typing_stop":
		// Broadcast typing stop
		c.hub.BroadcastToConversation(c.conversationID, map[string]interface{}{
			"type": "user_typing",
			"data": map[string]interface{}{
				"user_id":         c.userID,
				"conversation_id": c.conversationID,
				"is_typing":       false,
			},
			"timestamp": time.Now(),
		})

	default:
		log.Warn().Str("type", msgType).Msg("Unknown WebSocket message type")
	}
}

// sendError sends an error message to the client
func (c *Client) sendError(errorMsg string) {
	c.send <- map[string]interface{}{
		"type": "error",
		"data": map[string]interface{}{
			"message": errorMsg,
		},
		"timestamp": time.Now(),
	}
}

// Start begins reading and writing for the client
func (c *Client) Start() {
	go c.writePump()
	go c.readPump()
}

// WSMessage represents an incoming WebSocket message
type WSMessage struct {
	Type     string                 `json:"type"`
	Data     map[string]interface{} `json:"data,omitempty"`
	ClientID string                 `json:"client_id,omitempty"`
}

// WSEvent represents an outgoing WebSocket event
type WSEvent struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// MarshalJSON implements json.Marshaler for WSEvent
func (e *WSEvent) MarshalJSON() ([]byte, error) {
	type Alias WSEvent
	return json.Marshal(&struct {
		*Alias
		Timestamp string `json:"timestamp"`
	}{
		Alias:     (*Alias)(e),
		Timestamp: e.Timestamp.Format(time.RFC3339),
	})
}
