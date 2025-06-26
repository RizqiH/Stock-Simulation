package domain

import (
    "time"
    "encoding/json"
)

// WebSocketMessageType represents different types of WebSocket messages
type WebSocketMessageType string

const (
    // Market data messages
    MessageTypePriceUpdate     WebSocketMessageType = "PRICE_UPDATE"
    MessageTypeVolumeUpdate    WebSocketMessageType = "VOLUME_UPDATE"
    MessageTypeMarketStatus    WebSocketMessageType = "MARKET_STATUS"
    MessageTypeOrderBookUpdate WebSocketMessageType = "ORDER_BOOK_UPDATE"
    
    // Order and trade messages
    MessageTypeOrderUpdate     WebSocketMessageType = "ORDER_UPDATE"
    MessageTypeTradeExecution  WebSocketMessageType = "TRADE_EXECUTION"
    MessageTypePositionUpdate  WebSocketMessageType = "POSITION_UPDATE"
    MessageTypeBalanceUpdate   WebSocketMessageType = "BALANCE_UPDATE"
    
    // Portfolio messages
    MessageTypePortfolioUpdate WebSocketMessageType = "PORTFOLIO_UPDATE"
    MessageTypePnLUpdate       WebSocketMessageType = "PNL_UPDATE"
    
    // News and alerts
    MessageTypeNewsAlert       WebSocketMessageType = "NEWS_ALERT"
    MessageTypePriceAlert      WebSocketMessageType = "PRICE_ALERT"
    MessageTypeSystemAlert     WebSocketMessageType = "SYSTEM_ALERT"
    
    // Social trading
    MessageTypeSocialUpdate    WebSocketMessageType = "SOCIAL_UPDATE"
    MessageTypeFollowerUpdate  WebSocketMessageType = "FOLLOWER_UPDATE"
    
    // System messages
    MessageTypeHeartbeat       WebSocketMessageType = "HEARTBEAT"
    MessageTypeSubscription    WebSocketMessageType = "SUBSCRIPTION"
    MessageTypeUnsubscription  WebSocketMessageType = "UNSUBSCRIPTION"
    MessageTypeError           WebSocketMessageType = "ERROR"
)

// WebSocketMessage represents a generic WebSocket message
type WebSocketMessage struct {
    Type      WebSocketMessageType `json:"type"`
    Symbol    *string              `json:"symbol,omitempty"`
    UserID    *int                 `json:"user_id,omitempty"`
    Timestamp time.Time            `json:"timestamp"`
    Data      json.RawMessage      `json:"data"`
    Sequence  int64                `json:"sequence"`
    MessageID string               `json:"message_id"`
}

// PriceUpdateMessage represents a real-time price update
type PriceUpdateMessage struct {
    Symbol        string    `json:"symbol"`
    Price         float64   `json:"price"`
    Change        float64   `json:"change"`
    ChangePercent float64   `json:"change_percent"`
    Volume        int64     `json:"volume"`
    High          float64   `json:"high"`
    Low           float64   `json:"low"`
    Open          float64   `json:"open"`
    PreviousClose float64   `json:"previous_close"`
    BidPrice      *float64  `json:"bid_price,omitempty"`
    AskPrice      *float64  `json:"ask_price,omitempty"`
    BidSize       *int      `json:"bid_size,omitempty"`
    AskSize       *int      `json:"ask_size,omitempty"`
    LastTradeTime time.Time `json:"last_trade_time"`
    MarketCap     *int64    `json:"market_cap,omitempty"`
}

// OrderBookLevel represents a single level in the order book
type OrderBookLevel struct {
    Price    float64 `json:"price"`
    Quantity int     `json:"quantity"`
    Orders   int     `json:"orders"`
}

// OrderBookUpdate represents order book changes
type OrderBookUpdate struct {
    Symbol    string           `json:"symbol"`
    Bids      []OrderBookLevel `json:"bids"`      // Best bids (highest prices first)
    Asks      []OrderBookLevel `json:"asks"`      // Best asks (lowest prices first)
    Timestamp time.Time        `json:"timestamp"`
    Sequence  int64            `json:"sequence"`
}

// OrderUpdateMessage represents order status changes
type OrderUpdateMessage struct {
    OrderID           int         `json:"order_id"`
    Status            OrderStatus `json:"status"`
    ExecutedQuantity  int         `json:"executed_quantity"`
    RemainingQuantity int         `json:"remaining_quantity"`
    ExecutedPrice     *float64    `json:"executed_price,omitempty"`
    Commission        float64     `json:"commission"`
    Fees              float64     `json:"fees"`
    Message           string      `json:"message"`
}

// TradeExecutionMessage represents a completed trade
type TradeExecutionMessage struct {
    TradeID      int       `json:"trade_id"`
    OrderID      int       `json:"order_id"`
    Symbol       string    `json:"symbol"`
    Side         OrderSide `json:"side"`
    Quantity     int       `json:"quantity"`
    Price        float64   `json:"price"`
    TotalAmount  float64   `json:"total_amount"`
    Commission   float64   `json:"commission"`
    Fees         float64   `json:"fees"`
    ExecutedAt   time.Time `json:"executed_at"`
}

// PositionUpdateMessage represents position changes
type PositionUpdateMessage struct {
    Symbol           string  `json:"symbol"`
    Quantity         int     `json:"quantity"`
    AveragePrice     float64 `json:"average_price"`
    CurrentPrice     float64 `json:"current_price"`
    UnrealizedPnL    float64 `json:"unrealized_pnl"`
    RealizedPnL      float64 `json:"realized_pnl"`
    TotalPnL         float64 `json:"total_pnl"`
    TotalValue       float64 `json:"total_value"`
    DayChange        float64 `json:"day_change"`
    DayChangePercent float64 `json:"day_change_percent"`
}

// BalanceUpdateMessage represents balance changes
type BalanceUpdateMessage struct {
    CashBalance      float64 `json:"cash_balance"`
    MarginBalance    float64 `json:"margin_balance"`
    BuyingPower      float64 `json:"buying_power"`
    PortfolioValue   float64 `json:"portfolio_value"`
    TotalPnL         float64 `json:"total_pnl"`
    DayPnL           float64 `json:"day_pnl"`
    MaintenanceMargin float64 `json:"maintenance_margin"`
}

// NewsAlertMessage represents news or market alerts
type NewsAlertMessage struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Summary     string    `json:"summary"`
    Category    string    `json:"category"`    // EARNINGS, FDA, MERGER, etc.
    Severity    string    `json:"severity"`    // LOW, MEDIUM, HIGH, CRITICAL
    Symbols     []string  `json:"symbols"`     // Affected symbols
    Source      string    `json:"source"`
    PublishedAt time.Time `json:"published_at"`
    URL         *string   `json:"url,omitempty"`
}

// PriceAlertMessage represents price alert notifications
type PriceAlertMessage struct {
    AlertID       int     `json:"alert_id"`
    Symbol        string  `json:"symbol"`
    AlertType     string  `json:"alert_type"`     // ABOVE, BELOW, CHANGE_PERCENT
    TriggerPrice  float64 `json:"trigger_price"`
    CurrentPrice  float64 `json:"current_price"`
    Message       string  `json:"message"`
    TriggeredAt   time.Time `json:"triggered_at"`
}

// SystemAlertMessage represents system-wide alerts
type SystemAlertMessage struct {
    ID          int       `json:"id"`
    Type        string    `json:"type"`        // MAINTENANCE, OUTAGE, UPDATE
    Title       string    `json:"title"`
    Message     string    `json:"message"`
    Severity    string    `json:"severity"`    // INFO, WARNING, ERROR
    StartTime   *time.Time `json:"start_time,omitempty"`
    EndTime     *time.Time `json:"end_time,omitempty"`
    AffectedSystems []string `json:"affected_systems,omitempty"`
}

// SocialUpdateMessage represents social trading updates
type SocialUpdateMessage struct {
    UserID       int       `json:"user_id"`
    Username     string    `json:"username"`
    Action       string    `json:"action"`      // TRADE, FOLLOW, UNFOLLOW, POST
    Symbol       *string   `json:"symbol,omitempty"`
    Quantity     *int      `json:"quantity,omitempty"`
    Price        *float64  `json:"price,omitempty"`
    Message      *string   `json:"message,omitempty"`
    Timestamp    time.Time `json:"timestamp"`
    Followers    int       `json:"followers"`
}

// SubscriptionRequest represents a subscription request
type SubscriptionRequest struct {
    Type     string   `json:"type"`      // PRICE, ORDER_BOOK, NEWS, etc.
    Symbols  []string `json:"symbols"`   // Symbols to subscribe to
    UserID   *int     `json:"user_id,omitempty"`   // For user-specific data
    Interval *int     `json:"interval,omitempty"`  // Update interval in milliseconds
}

// SubscriptionResponse represents a subscription response
type SubscriptionResponse struct {
    Success     bool     `json:"success"`
    Message     string   `json:"message"`
    Subscribed  []string `json:"subscribed"`
    Failed      []string `json:"failed"`
    TotalSubs   int      `json:"total_subscriptions"`
}

// WebSocketConnection represents a WebSocket connection
type WebSocketConnection struct {
    ID            string                 `json:"id"`
    UserID        *int                   `json:"user_id,omitempty"`
    ConnectedAt   time.Time              `json:"connected_at"`
    LastHeartbeat time.Time              `json:"last_heartbeat"`
    Subscriptions map[string][]string    `json:"subscriptions"` // type -> symbols
    IsActive      bool                   `json:"is_active"`
    ClientInfo    *WebSocketClientInfo   `json:"client_info,omitempty"`
}

// WebSocketClientInfo represents client information
type WebSocketClientInfo struct {
    UserAgent   string `json:"user_agent"`
    IPAddress   string `json:"ip_address"`
    Platform    string `json:"platform"`
    Version     string `json:"version"`
}

// MarketDataSnapshot represents a complete market data snapshot
type MarketDataSnapshot struct {
    Symbol        string     `json:"symbol"`
    Price         float64    `json:"price"`
    Change        float64    `json:"change"`
    ChangePercent float64    `json:"change_percent"`
    Volume        int64      `json:"volume"`
    High          float64    `json:"high"`
    Low           float64    `json:"low"`
    Open          float64    `json:"open"`
    PreviousClose float64    `json:"previous_close"`
    MarketCap     *int64     `json:"market_cap,omitempty"`
    OrderBook     *OrderBookUpdate `json:"order_book,omitempty"`
    LastUpdate    time.Time  `json:"last_update"`
}

// WebSocketHub manages WebSocket connections and message distribution
type WebSocketHub interface {
    // Connection management
    AddConnection(conn *WebSocketConnection) error
    RemoveConnection(connectionID string) error
    GetConnection(connectionID string) (*WebSocketConnection, error)
    GetUserConnections(userID int) ([]*WebSocketConnection, error)
    
    // Subscription management
    Subscribe(connectionID string, subscriptionType string, symbols []string) error
    Unsubscribe(connectionID string, subscriptionType string, symbols []string) error
    GetSubscriptions(connectionID string) (map[string][]string, error)
    
    // Message broadcasting
    BroadcastToSymbol(symbol string, messageType WebSocketMessageType, data interface{}) error
    BroadcastToUser(userID int, messageType WebSocketMessageType, data interface{}) error
    BroadcastToAll(messageType WebSocketMessageType, data interface{}) error
    SendToConnection(connectionID string, messageType WebSocketMessageType, data interface{}) error
    
    // Heartbeat and health
    SendHeartbeat(connectionID string) error
    CheckConnections() error
    GetConnectionStats() (*ConnectionStats, error)
}

// ConnectionStats represents WebSocket connection statistics
type ConnectionStats struct {
    TotalConnections    int               `json:"total_connections"`
    AuthenticatedUsers  int               `json:"authenticated_users"`
    AnonymousUsers      int               `json:"anonymous_users"`
    TotalSubscriptions  int               `json:"total_subscriptions"`
    MessagesSent        int64             `json:"messages_sent"`
    MessagesReceived    int64             `json:"messages_received"`
    SubscriptionsByType map[string]int    `json:"subscriptions_by_type"`
    ConnectionsByHour   map[string]int    `json:"connections_by_hour"`
    LastUpdated         time.Time         `json:"last_updated"`
}

// RealTimeDataProvider interface for market data providers
type RealTimeDataProvider interface {
    // Price data
    GetRealTimePrice(symbol string) (*PriceUpdateMessage, error)
    GetOrderBook(symbol string, depth int) (*OrderBookUpdate, error)
    
    // Market status
    GetMarketStatus(marketCode string) (*MarketStatus, error)
    
    // Subscriptions
    SubscribeToSymbol(symbol string, callback func(*PriceUpdateMessage)) error
    UnsubscribeFromSymbol(symbol string) error
    
    // Provider management
    Start() error
    Stop() error
    IsConnected() bool
    GetLastUpdate() time.Time
}

// WebSocket message creation helpers
func NewWebSocketMessage(msgType WebSocketMessageType, data interface{}) (*WebSocketMessage, error) {
    dataBytes, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }
    
    return &WebSocketMessage{
        Type:      msgType,
        Timestamp: time.Now(),
        Data:      dataBytes,
        Sequence:  time.Now().UnixNano(),
        MessageID: generateMessageID(),
    }, nil
}

func (wsm *WebSocketMessage) UnmarshalData(v interface{}) error {
    return json.Unmarshal(wsm.Data, v)
}

func (wsm *WebSocketMessage) SetSymbol(symbol string) {
    wsm.Symbol = &symbol
}

func (wsm *WebSocketMessage) SetUserID(userID int) {
    wsm.UserID = &userID
}

// Helper function to generate unique message IDs
func generateMessageID() string {
    return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
    }
    return string(b)
}

// Connection validation
func (wsc *WebSocketConnection) IsExpired(timeout time.Duration) bool {
    return time.Since(wsc.LastHeartbeat) > timeout
}

func (wsc *WebSocketConnection) UpdateHeartbeat() {
    wsc.LastHeartbeat = time.Now()
}

func (wsc *WebSocketConnection) AddSubscription(subscriptionType string, symbols []string) {
    if wsc.Subscriptions == nil {
        wsc.Subscriptions = make(map[string][]string)
    }
    
    existing := wsc.Subscriptions[subscriptionType]
    for _, symbol := range symbols {
        // Add only if not already subscribed
        found := false
        for _, existingSymbol := range existing {
            if existingSymbol == symbol {
                found = true
                break
            }
        }
        if !found {
            existing = append(existing, symbol)
        }
    }
    wsc.Subscriptions[subscriptionType] = existing
}

func (wsc *WebSocketConnection) RemoveSubscription(subscriptionType string, symbols []string) {
    if wsc.Subscriptions == nil {
        return
    }
    
    existing := wsc.Subscriptions[subscriptionType]
    var filtered []string
    
    for _, existingSymbol := range existing {
        shouldKeep := true
        for _, symbolToRemove := range symbols {
            if existingSymbol == symbolToRemove {
                shouldKeep = false
                break
            }
        }
        if shouldKeep {
            filtered = append(filtered, existingSymbol)
        }
    }
    
    if len(filtered) == 0 {
        delete(wsc.Subscriptions, subscriptionType)
    } else {
        wsc.Subscriptions[subscriptionType] = filtered
    }
} 