package domain

import (
    "time"
    "fmt"
)

// OrderType represents different types of orders
type OrderType string

const (
    OrderTypeMarket     OrderType = "MARKET"
    OrderTypeLimit      OrderType = "LIMIT"
    OrderTypeStopLoss   OrderType = "STOP_LOSS"
    OrderTypeTakeProfit OrderType = "TAKE_PROFIT"
    OrderTypeTrailingStop OrderType = "TRAILING_STOP"
    OrderTypeOCO        OrderType = "OCO" // One-Cancels-Other
)

// OrderStatus represents the current status of an order
type OrderStatus string

const (
    OrderStatusPending   OrderStatus = "PENDING"
    OrderStatusExecuted  OrderStatus = "EXECUTED"
    OrderStatusCancelled OrderStatus = "CANCELLED"
    OrderStatusExpired   OrderStatus = "EXPIRED"
    OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
)

// OrderSide represents buy or sell
type OrderSide string

const (
    OrderSideBuy  OrderSide = "BUY"
    OrderSideSell OrderSide = "SELL"
    OrderSideShort OrderSide = "SHORT" // For short selling
    OrderSideCover OrderSide = "COVER" // To cover short position
)

// TimeInForce represents how long an order stays active
type TimeInForce string

const (
    TimeInForceGTC TimeInForce = "GTC" // Good Till Cancelled
    TimeInForceIOC TimeInForce = "IOC" // Immediate or Cancel
    TimeInForceFOK TimeInForce = "FOK" // Fill or Kill
    TimeInForceDAY TimeInForce = "DAY" // Day order
)

// Order represents a trading order with advanced features
type Order struct {
    ID                 int         `json:"id" db:"id"`
    UserID             int         `json:"user_id" db:"user_id"`
    StockSymbol        string      `json:"stock_symbol" db:"stock_symbol"`
    OrderType          OrderType   `json:"order_type" db:"order_type"`
    Side               OrderSide   `json:"side" db:"side"`
    Quantity           int         `json:"quantity" db:"quantity"`
    Price              *float64    `json:"price,omitempty" db:"price"` // For limit orders
    StopPrice          *float64    `json:"stop_price,omitempty" db:"stop_price"` // For stop orders
    TrailingAmount     *float64    `json:"trailing_amount,omitempty" db:"trailing_amount"` // For trailing stops
    TrailingPercent    *float64    `json:"trailing_percent,omitempty" db:"trailing_percent"` // For trailing stops
    TimeInForce        TimeInForce `json:"time_in_force" db:"time_in_force"`
    Status             OrderStatus `json:"status" db:"status"`
    ExecutedPrice      *float64    `json:"executed_price,omitempty" db:"executed_price"`
    ExecutedQuantity   int         `json:"executed_quantity" db:"executed_quantity"`
    RemainingQuantity  int         `json:"remaining_quantity" db:"remaining_quantity"`
    ExecutedAt         *time.Time  `json:"executed_at,omitempty" db:"executed_at"`
    ExpiresAt          *time.Time  `json:"expires_at,omitempty" db:"expires_at"`
    CreatedAt          time.Time   `json:"created_at" db:"created_at"`
    UpdatedAt          time.Time   `json:"updated_at" db:"updated_at"`
    
    // OCO Related
    ParentOrderID      *int        `json:"parent_order_id,omitempty" db:"parent_order_id"`
    LinkedOrderID      *int        `json:"linked_order_id,omitempty" db:"linked_order_id"`
    
    // Commission and fees
    Commission         float64     `json:"commission" db:"commission"`
    Fees               float64     `json:"fees" db:"fees"`
    
    // Market data at time of order
    MarketPrice        float64     `json:"market_price" db:"market_price"`
    BidPrice           *float64    `json:"bid_price,omitempty" db:"bid_price"`
    AskPrice           *float64    `json:"ask_price,omitempty" db:"ask_price"`
    Spread             *float64    `json:"spread,omitempty" db:"spread"`
}

// OrderRequest represents a request to create an order
type OrderRequest struct {
    StockSymbol        string      `json:"stock_symbol" binding:"required"`
    OrderType          OrderType   `json:"order_type" binding:"required"`
    Side               OrderSide   `json:"side" binding:"required"`
    Quantity           int         `json:"quantity" binding:"required,min=1"`
    Price              *float64    `json:"price,omitempty"`
    StopPrice          *float64    `json:"stop_price,omitempty"`
    TrailingAmount     *float64    `json:"trailing_amount,omitempty"`
    TrailingPercent    *float64    `json:"trailing_percent,omitempty"`
    TimeInForce        TimeInForce `json:"time_in_force"`
    ExpiresAt          *time.Time  `json:"expires_at,omitempty"`
    LinkedOrderRequest *OrderRequest `json:"linked_order,omitempty"` // For OCO orders
}

// OrderExecution represents the result of order execution
type OrderExecution struct {
    OrderID           int       `json:"order_id"`
    ExecutedPrice     float64   `json:"executed_price"`
    ExecutedQuantity  int       `json:"executed_quantity"`
    Commission        float64   `json:"commission"`
    Fees              float64   `json:"fees"`
    Slippage          float64   `json:"slippage"`
    ExecutedAt        time.Time `json:"executed_at"`
    TotalAmount       float64   `json:"total_amount"`
}

// OrderStats represents order statistics for a user
type OrderStats struct {
    TotalOrders        int     `json:"total_orders"`
    PendingOrders      int     `json:"pending_orders"`
    ExecutedOrders     int     `json:"executed_orders"`
    CancelledOrders    int     `json:"cancelled_orders"`
    PartiallyFilled    int     `json:"partially_filled"`
    SuccessRate        float64 `json:"success_rate"`
    AverageExecutionTime float64 `json:"average_execution_time"`
    TotalCommission    float64 `json:"total_commission"`
    TotalFees          float64 `json:"total_fees"`
}

// Validation methods
func (o *Order) Validate() error {
    if o.UserID <= 0 {
        return fmt.Errorf("invalid user ID")
    }
    
    if o.StockSymbol == "" {
        return fmt.Errorf("stock symbol is required")
    }
    
    if o.Quantity <= 0 {
        return fmt.Errorf("quantity must be positive")
    }
    
    // Validate order type specific requirements
    switch o.OrderType {
    case OrderTypeLimit:
        if o.Price == nil || *o.Price <= 0 {
            return fmt.Errorf("limit orders require a valid price")
        }
    case OrderTypeStopLoss, OrderTypeTakeProfit:
        if o.StopPrice == nil || *o.StopPrice <= 0 {
            return fmt.Errorf("stop orders require a valid stop price")
        }
    case OrderTypeTrailingStop:
        if (o.TrailingAmount == nil && o.TrailingPercent == nil) ||
           (o.TrailingAmount != nil && *o.TrailingAmount <= 0) ||
           (o.TrailingPercent != nil && (*o.TrailingPercent <= 0 || *o.TrailingPercent >= 100)) {
            return fmt.Errorf("trailing stop orders require valid trailing amount or percent")
        }
    }
    
    return nil
}

// Helper methods
func (o *Order) IsActive() bool {
    return o.Status == OrderStatusPending || o.Status == OrderStatusPartiallyFilled
}

func (o *Order) IsCompleted() bool {
    return o.Status == OrderStatusExecuted || o.Status == OrderStatusCancelled || o.Status == OrderStatusExpired
}

func (o *Order) CalculateTotalCost() float64 {
    if o.ExecutedPrice == nil {
        return 0
    }
    return float64(o.ExecutedQuantity) * (*o.ExecutedPrice) + o.Commission + o.Fees
}

func (o *Order) CanBeExecuted(currentPrice float64) bool {
    if !o.IsActive() {
        return false
    }
    
    switch o.OrderType {
    case OrderTypeMarket:
        return true
    case OrderTypeLimit:
        if o.Price == nil {
            return false
        }
        if o.Side == OrderSideBuy {
            return currentPrice <= *o.Price
        }
        return currentPrice >= *o.Price
    case OrderTypeStopLoss:
        if o.StopPrice == nil {
            return false
        }
        if o.Side == OrderSideBuy {
            return currentPrice >= *o.StopPrice
        }
        return currentPrice <= *o.StopPrice
    case OrderTypeTakeProfit:
        if o.StopPrice == nil {
            return false
        }
        if o.Side == OrderSideBuy {
            return currentPrice <= *o.StopPrice
        }
        return currentPrice >= *o.StopPrice
    }
    
    return false
} 