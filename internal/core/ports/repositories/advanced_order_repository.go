package repositories

import (
	"stock-simulation-backend/internal/core/domain"
	"time"
)

type AdvancedOrderRepository interface {
	// Basic CRUD operations
	Create(order *domain.Order) error
	GetByID(id int) (*domain.Order, error)
	Update(order *domain.Order) error
	Delete(id int) error
	
	// Query methods
	GetByUserID(userID int, limit, offset int) ([]domain.Order, error)
	GetByUserIDAndStatus(userID int, status domain.OrderStatus) ([]domain.Order, error)
	GetBySymbol(symbol string) ([]domain.Order, error)
	GetByMarketAndStatus(marketCode string, status domain.OrderStatus) ([]domain.Order, error)
	
	// Active orders management
	GetActiveOrders() ([]domain.Order, error)
	GetPendingOrders() ([]domain.Order, error)
	GetExpiredOrders() ([]domain.Order, error)
	GetOrdersForExecution(symbol string, currentPrice float64) ([]domain.Order, error)
	
	// Order type specific queries
	GetLimitOrders(symbol string, side domain.OrderSide) ([]domain.Order, error)
	GetStopOrders(symbol string) ([]domain.Order, error)
	GetTrailingStopOrders() ([]domain.Order, error)
	GetOCOOrders(userID int) ([]domain.Order, error)
	
	// Order execution
	ExecuteOrder(orderID int, executedPrice float64, executedQuantity int) error
	PartialFillOrder(orderID int, filledQuantity int, filledPrice float64) error
	CancelOrder(orderID int, reason string) error
	ExpireOrder(orderID int) error
	
	// Order status updates
	UpdateStatus(orderID int, status domain.OrderStatus) error
	UpdateExecutionDetails(orderID int, execution *domain.OrderExecution) error
	
	// OCO order management
	CreateOCOOrders(parentOrder *domain.Order, linkedOrder *domain.Order) error
	CancelLinkedOrders(orderID int) error
	GetLinkedOrders(orderID int) ([]domain.Order, error)
	
	// Statistics and analytics
	GetUserOrderStats(userID int) (*domain.OrderStats, error)
	GetOrderStatsByDateRange(userID int, startDate, endDate time.Time) (*domain.OrderStats, error)
	GetOrderCountByType(userID int) (map[domain.OrderType]int, error)
	GetAverageExecutionTime(userID int) (float64, error)
	
	// Reporting and history
	GetOrderHistory(userID int, limit, offset int) ([]domain.Order, error)
	GetOrdersByDateRange(userID int, startDate, endDate time.Time) ([]domain.Order, error)
	GetOrdersForSymbolAndDateRange(userID int, symbol string, startDate, endDate time.Time) ([]domain.Order, error)
	
	// Validation and constraints
	ValidateOrderConstraints(order *domain.Order) error
	CheckDailyOrderLimit(userID int) (bool, error)
	CheckOrderSizeLimit(userID int, symbol string, quantity int) (bool, error)
	
	// Trailing stop specific
	UpdateTrailingStopPrice(orderID int, newStopPrice float64) error
	GetTrailingStopsToUpdate(priceUpdates map[string]float64) ([]domain.Order, error)
	
	// Commission and fees calculation
	CalculateOrderCommission(userID int, order *domain.Order) (float64, error)
	UpdateOrderCommission(orderID int, commission, fees float64) error
	
	// Market hours integration
	GetOrdersAwaitingMarketOpen(marketCode string) ([]domain.Order, error)
	GetOrdersToExpireAtMarketClose(marketCode string) ([]domain.Order, error)
}

// OrderSearchCriteria represents search criteria for orders
type OrderSearchCriteria struct {
	UserID       *int                   `json:"user_id,omitempty"`
	Symbol       *string                `json:"symbol,omitempty"`
	OrderType    *domain.OrderType      `json:"order_type,omitempty"`
	Status       *domain.OrderStatus    `json:"status,omitempty"`
	Side         *domain.OrderSide      `json:"side,omitempty"`
	StartDate    *time.Time             `json:"start_date,omitempty"`
	EndDate      *time.Time             `json:"end_date,omitempty"`
	MinPrice     *float64               `json:"min_price,omitempty"`
	MaxPrice     *float64               `json:"max_price,omitempty"`
	MinQuantity  *int                   `json:"min_quantity,omitempty"`
	MaxQuantity  *int                   `json:"max_quantity,omitempty"`
	SortBy       string                 `json:"sort_by,omitempty"`       // created_at, price, quantity
	SortOrder    string                 `json:"sort_order,omitempty"`    // ASC, DESC
	Limit        int                    `json:"limit,omitempty"`
	Offset       int                    `json:"offset,omitempty"`
}

// OrderSearchResult represents the result of an order search
type OrderSearchResult struct {
	Orders     []domain.Order `json:"orders"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// Enhanced repository interface with search capabilities
type AdvancedOrderRepositoryWithSearch interface {
	AdvancedOrderRepository
	
	// Advanced search
	SearchOrders(criteria *OrderSearchCriteria) (*OrderSearchResult, error)
	SearchOrdersByUser(userID int, criteria *OrderSearchCriteria) (*OrderSearchResult, error)
	SearchOrdersBySymbol(symbol string, criteria *OrderSearchCriteria) (*OrderSearchResult, error)
	Search(userID int, criteria *OrderSearchCriteria) (*OrderSearchResult, error)
	
	// Bulk operations
	BulkUpdateStatus(orderIDs []int, status domain.OrderStatus) error
	BulkCancelOrders(orderIDs []int, reason string) error
	BulkExecuteOrders(executions []domain.OrderExecution) error
	
	// Performance analytics
	GetExecutionMetrics(userID int, timeframe string) (*OrderExecutionMetrics, error)
	GetSlippageAnalysis(userID int, symbol string) (*SlippageAnalysis, error)
	
	// Additional methods needed by service
	CancelAllOrdersByUser(userID int, symbol *string) (int, error)
	GetActiveOrdersByUser(userID int) ([]domain.Order, error)
	GetOrderStatistics(userID int) (*domain.OrderStats, error)
}

// OrderExecutionMetrics represents order execution performance metrics
type OrderExecutionMetrics struct {
	UserID              int     `json:"user_id"`
	Timeframe           string  `json:"timeframe"`
	TotalOrders         int     `json:"total_orders"`
	ExecutedOrders      int     `json:"executed_orders"`
	CancelledOrders     int     `json:"cancelled_orders"`
	PartiallyFilled     int     `json:"partially_filled"`
	AverageExecutionTime float64 `json:"average_execution_time"` // in seconds
	FillRate            float64 `json:"fill_rate"`              // percentage
	AverageSlippage     float64 `json:"average_slippage"`       // in basis points
	BestExecution       float64 `json:"best_execution"`         // best price improvement
	WorstExecution      float64 `json:"worst_execution"`        // worst slippage
	TotalCommission     float64 `json:"total_commission"`
	TotalFees           float64 `json:"total_fees"`
}

// SlippageAnalysis represents slippage analysis for a symbol
type SlippageAnalysis struct {
	Symbol            string  `json:"symbol"`
	UserID            int     `json:"user_id"`
	TotalTrades       int     `json:"total_trades"`
	AverageSlippage   float64 `json:"average_slippage"`   // in basis points
	MedianSlippage    float64 `json:"median_slippage"`    // in basis points
	SlippageStdDev    float64 `json:"slippage_std_dev"`   // standard deviation
	BestExecution     float64 `json:"best_execution"`     // best price improvement
	WorstSlippage     float64 `json:"worst_slippage"`     // worst slippage experienced
	MarketOrderSlippage float64 `json:"market_order_slippage"` // average for market orders
	LimitOrderSlippage  float64 `json:"limit_order_slippage"`  // average for limit orders
} 