package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"stock-simulation-backend/internal/core/domain"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type RealTimeService struct {
	clients    map[string]*websocket.Conn
	clientsMu  sync.RWMutex
	upgrader   websocket.Upgrader
	broadcast  chan domain.PriceUpdateMessage
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

func NewRealTimeService() *RealTimeService {
	return &RealTimeService{
		clients:    make(map[string]*websocket.Conn),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow connections from any origin in development
			},
		},
		broadcast:  make(chan domain.PriceUpdateMessage),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// Start the WebSocket hub
func (s *RealTimeService) Start() {
	log.Println("🔌 Starting WebSocket service for real-time updates...")
	
	go func() {
		for {
			select {
			case conn := <-s.register:
				s.handleNewConnection(conn)
			
			case conn := <-s.unregister:
				s.handleDisconnection(conn)
			
			case priceUpdate := <-s.broadcast:
				s.broadcastToAllClients(priceUpdate)
			}
		}
	}()
}

// Handle WebSocket connection upgrade
func (s *RealTimeService) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ WebSocket upgrade failed: %v", err)
		return
	}
	
	log.Printf("🔌 New WebSocket connection established from %s", conn.RemoteAddr())
	
	// Register the new connection
	s.register <- conn
	
	// Handle incoming messages and connection cleanup
	go s.handleClient(conn)
}

// Handle individual client connection
func (s *RealTimeService) handleClient(conn *websocket.Conn) {
	defer func() {
		s.unregister <- conn
		conn.Close()
	}()
	
	// Set read deadline and pong handler for keepalive
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	
	// Listen for messages from client
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("❌ WebSocket error: %v", err)
			}
			break
		}
		
		// Handle different message types
		if messageType == websocket.TextMessage {
			s.handleClientMessage(conn, message)
		}
	}
}

// Handle messages from clients
func (s *RealTimeService) handleClientMessage(conn *websocket.Conn, message []byte) {
	var clientMsg map[string]interface{}
	if err := json.Unmarshal(message, &clientMsg); err != nil {
		log.Printf("⚠️ Invalid client message: %v", err)
		return
	}
	
	// Handle ping messages
	if msgType, exists := clientMsg["type"]; exists && msgType == "ping" {
		pongResponse := map[string]interface{}{
			"type":      "pong",
			"timestamp": time.Now().Unix(),
		}
		s.sendToClient(conn, pongResponse)
	}
}

// Handle new connection registration
func (s *RealTimeService) handleNewConnection(conn *websocket.Conn) {
	s.clientsMu.Lock()
	clientID := fmt.Sprintf("%s-%d", conn.RemoteAddr().String(), time.Now().UnixNano())
	s.clients[clientID] = conn
	s.clientsMu.Unlock()
	
	log.Printf("✅ Client registered: %s (Total: %d)", clientID, len(s.clients))
	
	// Send welcome message
	welcomeMsg := map[string]interface{}{
		"type":    "welcome",
		"message": "Connected to real-time price updates",
		"client_id": clientID,
		"timestamp": time.Now().Unix(),
	}
	s.sendToClient(conn, welcomeMsg)
}

// Handle client disconnection
func (s *RealTimeService) handleDisconnection(conn *websocket.Conn) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	
	// Find and remove the client
	for clientID, client := range s.clients {
		if client == conn {
			delete(s.clients, clientID)
			log.Printf("🔌 Client disconnected: %s (Remaining: %d)", clientID, len(s.clients))
			break
		}
	}
}

// Broadcast price update to all connected clients
func (s *RealTimeService) BroadcastPriceUpdate(update domain.PriceUpdateMessage) {
	select {
	case s.broadcast <- update:
		// Message queued for broadcast
	default:
		log.Println("⚠️ Broadcast channel full, dropping price update")
	}
}

// Send message to all clients
func (s *RealTimeService) broadcastToAllClients(priceUpdate domain.PriceUpdateMessage) {
	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()
	
	if len(s.clients) == 0 {
		return
	}
	
	message := map[string]interface{}{
		"type":      "price_update",
		"data":      priceUpdate,
		"timestamp": time.Now().Unix(),
	}
	
	var disconnectedClients []string
	
	for clientID, conn := range s.clients {
		if err := s.sendToClient(conn, message); err != nil {
			log.Printf("⚠️ Failed to send to client %s: %v", clientID, err)
			disconnectedClients = append(disconnectedClients, clientID)
		}
	}
	
	// Remove disconnected clients
	for _, clientID := range disconnectedClients {
		delete(s.clients, clientID)
	}
	
	if len(s.clients) > 0 {
		fmt.Printf("📡 Broadcasted %s price update to %d clients\n", priceUpdate.Symbol, len(s.clients))
	}
}

// Send message to specific client
func (s *RealTimeService) sendToClient(conn *websocket.Conn, message interface{}) error {
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return conn.WriteJSON(message)
}

// Get connected clients count
func (s *RealTimeService) GetConnectedClientsCount() int {
	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()
	return len(s.clients)
}

// Send heartbeat to all clients
func (s *RealTimeService) SendHeartbeat() {
	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()
	
	heartbeat := map[string]interface{}{
		"type":      "heartbeat",
		"timestamp": time.Now().Unix(),
		"clients":   len(s.clients),
	}
	
	for _, conn := range s.clients {
		s.sendToClient(conn, heartbeat)
	}
} 