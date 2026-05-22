package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

type Hub struct {
	Clients    map[string]*Client
	Rooms      map[string]map[string]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *WSMessage
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Rooms:      make(map[string]map[string]bool),
		Register:   make(chan *Client, 100),
		Unregister: make(chan *Client, 100),
		Broadcast:  make(chan *WSMessage, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ClientID] = client
			if client.Room != "" {
				if _, exists := h.Rooms[client.Room]; !exists {
					h.Rooms[client.Room] = make(map[string]bool)
				}
				h.Rooms[client.Room][client.ClientID] = true
			}
			h.mu.Unlock()
			log.Printf("Client registered: id=%s, room=%s, role=%s", client.ClientID, client.Room, client.Role)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ClientID]; ok {
				delete(h.Clients, client.ClientID)
				close(client.Send)
				if client.Room != "" {
					if room, exists := h.Rooms[client.Room]; exists {
						delete(room, client.ClientID)
						if len(room) == 0 {
							delete(h.Rooms, client.Room)
						}
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client unregistered: id=%s", client.ClientID)

		case message := <-h.Broadcast:
			h.mu.RLock()
			if message.Room != "" {
				h.broadcastToRoom(message)
			} else {
				h.broadcastAll(message)
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) broadcastAll(message *WSMessage) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal broadcast message: %v", err)
		return
	}

	for _, client := range h.Clients {
		select {
		case client.Send <- data:
		default:
			log.Printf("Send buffer full for client: %s, dropping message", client.ClientID)
		}
	}
}

func (h *Hub) broadcastToRoom(message *WSMessage) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal room message: %v", err)
		return
	}

	if roomClients, exists := h.Rooms[message.Room]; exists {
		for clientID := range roomClients {
			if client, ok := h.Clients[clientID]; ok {
				select {
				case client.Send <- data:
				default:
					log.Printf("Send buffer full for client: %s, dropping room message", clientID)
				}
			}
		}
	}
}

func (h *Hub) SendToClient(clientID string, message *WSMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.Clients[clientID]; ok {
		data, err := json.Marshal(message)
		if err != nil {
			return
		}
		select {
		case client.Send <- data:
		default:
		}
	}
}

func (h *Hub) SendToRoom(room string, message *WSMessage) {
	message.Room = room
	select {
	case h.Broadcast <- message:
	default:
		log.Printf("Hub broadcast buffer full, dropping message for room: %s", room)
	}
}

func (h *Hub) HandleClientMessage(client *Client, msg *WSMessage) {
	switch msg.Type {
	case WSHeartbeat:
		h.SendToClient(client.ClientID, &WSMessage{
			Type:      WSHeartbeat,
			Timestamp: time.Now(),
			Payload:   map[string]string{"status": "ok"},
		})
	case WSProbeEvent:
		h.SendToRoom(msg.Room, msg)
	default:
		log.Printf("Unhandled message type from client %s: %s", client.ClientID, msg.Type)
	}
}
