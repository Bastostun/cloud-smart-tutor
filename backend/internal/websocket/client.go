package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WSMessageType string

const (
	WSNodeMetrics      WSMessageType = "node_metrics"
	WSProbeEvent       WSMessageType = "probe_event"
	WSErrorCluster     WSMessageType = "error_cluster"
	WSTeachingAlert    WSMessageType = "teaching_alert"
	WSAICompanion      WSMessageType = "ai_companion"
	WSAuditLog         WSMessageType = "audit_log"
	WSCloneStatus      WSMessageType = "clone_status"
	WSHeartbeat        WSMessageType = "heartbeat"
)

type WSMessage struct {
	Type      WSMessageType `json:"type"`
	Timestamp time.Time     `json:"timestamp"`
	Payload   interface{}   `json:"payload"`
	Room      string        `json:"room,omitempty"` // classroom_id or student_id
}

type Client struct {
	Hub         *Hub
	Conn        *websocket.Conn
	Send        chan []byte
	ClientID    string
	Room        string
	Role        string // "student", "teacher", "enterprise", "admin"
	mu          sync.Mutex
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("WebSocket unexpected close: client=%s, error=%v", c.ClientID, err)
			}
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Invalid WebSocket message: client=%s, error=%v", c.ClientID, err)
			continue
		}

		c.Hub.HandleClientMessage(c, &msg)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.mu.Lock()
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			c.mu.Unlock()

			if err != nil {
				log.Printf("WebSocket write error: client=%s, error=%v", c.ClientID, err)
				return
			}

		case <-ticker.C:
			c.mu.Lock()
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.mu.Unlock()
				return
			}
			c.mu.Unlock()
		}
	}
}

func (c *Client) SendJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	select {
	case c.Send <- data:
		return nil
	default:
		return nil
	}
}
