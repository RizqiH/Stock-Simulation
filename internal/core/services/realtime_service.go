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
	clients      map[string]*websocket.Conn
	clientsMu    sync.RWMutex
	upgrader     websocket.Upgrader
	broadcast    chan domain.PriceUpdateMessage
	register     chan *websocket.Conn
	unregister   chan *websocket.Conn
	redisService *RedisService // Redis service for pub/sub
}

func NewRealTimeService(redisService *RedisService) *RealTimeService {
	return &RealTimeService{
		clients:      make(map[string]*websocket.Conn),
		redisService: redisService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow connections from any origin in development
			},
		},
		broadcast:  make(chan domain.PriceUpdateMessage, 100),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// Start the WebSocket hub and Redis subscription
func (s *RealTimeService) Start() {
	log.Println("🔌 Starting WebSocket service for real-time updates...")
	
	// Start WebSocket hub
	go s.runWebSocketHub()
	
	// Start Redis subscription if available
	if s.redisService != nil {
		go s.subscribeToRedis()
	} else {
		log.Println("⚠️ Redis service not available, using WebSocket only")
	}
}

// runWebSocketHub manages WebSocket connections and broadcasting
func (s *RealTimeService) runWebSocketHub() {
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
}

// subscribeToRedis subscribes to Redis price updates
func (s *RealTimeService) subscribeToRedis() {
	log.Println("🔌 Starting Redis subscription for price updates...")
	
	priceUpdateChan, err := s.redisService.SubscribeToPriceUpdates()
	if err != nil {
		log.Printf("❌ Failed to subscribe to Redis: %v", err)
		return
	}
	
	// Forward Redis messages to WebSocket broadcast
	for priceUpdate := range priceUpdateChan {
		select {
		case s.broadcast <- priceUpdate:
			// Successfully queued for broadcast
		default:
			// Channel full, drop message
			log.Println("⚠️ WebSocket broadcast channel full, dropping Redis message")
		}
	}
	
	log.Println("📡 Redis subscription ended")
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
		if err := conn.Close(); err != nil {
			log.Printf("⚠️ Failed to close WebSocket connection: %v", err)
		}
	}()
	
	// Set read deadline and pong handler for keepalive
	if err := conn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
		log.Printf("⚠️ Failed to set read deadline: %v", err)
	}
	conn.SetPongHandler(func(string) error {
		if err := conn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
			log.Printf("⚠️ Failed to set read deadline in pong handler: %v", err)
		}
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
		if err := s.sendToClient(conn, pongResponse); err != nil {
			log.Printf("⚠️ Failed to send pong response: %v", err)
		}
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
		"type":       "welcome",
		"message":    "Connected to real-time price updates",
		"client_id":  clientID,
		"timestamp":  time.Now().Unix(),
		"redis_enabled": s.redisService != nil,
	}
	if err := s.sendToClient(conn, welcomeMsg); err != nil {
		log.Printf("⚠️ Failed to send welcome message: %v", err)
	}
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

// BroadcastPriceUpdate publishes price update to Redis and broadcasts locally
func (s *RealTimeService) BroadcastPriceUpdate(update domain.PriceUpdateMessage) {
	// Publish to Redis if available (for scaling across multiple instances)
	if s.redisService != nil {
		if err := s.redisService.PublishPriceUpdate(update); err != nil {
			log.Printf("⚠️ Failed to publish to Redis: %v", err)
		}
	}
	
	// Also broadcast locally for direct WebSocket connections
	select {
	case s.broadcast <- update:
		// Message queued for broadcast
	default:
		log.Println("⚠️ Local broadcast channel full, dropping price update")
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
		"source":    "redis", // Indicates this came from Redis pub/sub
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

// BroadcastMarketStatus broadcasts market status updates
func (s *RealTimeService) BroadcastMarketStatus(status string) {
	if s.redisService != nil {
		if err := s.redisService.PublishMarketStatus(status); err != nil {
			log.Printf("⚠️ Failed to publish market status to Redis: %v", err)
		}
	}
	
	// Also broadcast locally
	message := map[string]interface{}{
		"type":      "market_status",
		"status":    status,
		"timestamp": time.Now().Unix(),
	}
	
	s.clientsMu.RLock()
	for _, conn := range s.clients {
		if err := s.sendToClient(conn, message); err != nil {
			log.Printf("⚠️ Failed to send market status: %v", err)
		}
	}
	s.clientsMu.RUnlock()
}

// BroadcastTradingAlert broadcasts trading alerts
func (s *RealTimeService) BroadcastTradingAlert(alert domain.TradingAlert) {
	if s.redisService != nil {
		if err := s.redisService.PublishTradingAlert(alert); err != nil {
			log.Printf("⚠️ Failed to publish alert to Redis: %v", err)
		}
	}
	
	// Also broadcast locally
	message := map[string]interface{}{
		"type":      "trading_alert",
		"data":      alert,
		"timestamp": time.Now().Unix(),
	}
	
	s.clientsMu.RLock()
	for _, conn := range s.clients {
		if err := s.sendToClient(conn, message); err != nil {
			log.Printf("⚠️ Failed to send trading alert: %v", err)
		}
	}
	s.clientsMu.RUnlock()
}

// Send message to specific client
func (s *RealTimeService) sendToClient(conn *websocket.Conn, message interface{}) error {
	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		log.Printf("⚠️ Failed to set write deadline: %v", err)
		return err
	}
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
		"type":         "heartbeat",
		"timestamp":    time.Now().Unix(),
		"clients":      len(s.clients),
		"redis_status": s.redisService != nil && s.redisService.GetConnectionStatus(),
	}
	
	for _, conn := range s.clients {
		if err := s.sendToClient(conn, heartbeat); err != nil {
			log.Printf("⚠️ Failed to send heartbeat: %v", err)
		}
	}
}

// GetServiceStatus returns the service status including Redis
func (s *RealTimeService) GetServiceStatus() map[string]interface{} {
	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()
	
	status := map[string]interface{}{
		"websocket_clients": len(s.clients),
		"redis_enabled":     s.redisService != nil,
		"redis_connected":   false,
	}
	
	if s.redisService != nil {
		status["redis_connected"] = s.redisService.GetConnectionStatus()
	}
	
	return status
} 