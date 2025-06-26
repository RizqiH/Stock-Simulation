package services

import (
	"stock-simulation-backend/internal/core/domain"
	"stock-simulation-backend/internal/core/ports/repositories"
	"time"
)

type AdvancedOrderService interface {
	// Order creation and management
	CreateOrder(userID int, request *domain.OrderRequest) (*domain.Order, error)
	CreateOCOOrder(userID int, parentRequest *domain.OrderRequest, linkedRequest *domain.OrderRequest) (*domain.Order, *domain.Order, error)
	ModifyOrder(userID int, orderID int, modifications *OrderModificationRequest) (*domain.Order, error)
	CancelOrder(userID int, orderID int) error
	CancelAllOrders(userID int, symbol *string) (int, error)
	
	// Order execution
	ExecuteOrder(orderID int, marketPrice float64) (*domain.OrderExecution, error)
	ExecuteMarketOrders(symbol string, currentPrice float64) ([]domain.OrderExecution, error)
	ExecuteLimitOrders(symbol string, currentPrice float64) ([]domain.OrderExecution, error)
	ExecuteStopOrders(symbol string, currentPrice float64) ([]domain.OrderExecution, error)
	ExecuteTrailingStops(priceUpdates map[string]float64) ([]domain.OrderExecution, error)
	
	// Order validation
	ValidateOrder(userID int, request *domain.OrderRequest) error
	ValidateBuyingPower(userID int, order *domain.Order) error
	ValidatePosition(userID int, order *domain.Order) error
	ValidateMarketHours(order *domain.Order) error
	ValidateOrderLimits(userID int, order *domain.Order) error
	
	// Order queries
	GetUserOrders(userID int, status *domain.OrderStatus, limit, offset int) ([]domain.Order, error)
	GetOrderByID(userID int, orderID int) (*domain.Order, error)
	GetActiveOrders(userID int) ([]domain.Order, error)
	GetOrderHistory(userID int, startDate, endDate *time.Time, limit, offset int) ([]domain.Order, error)
	SearchOrders(userID int, criteria *repositories.OrderSearchCriteria) (*repositories.OrderSearchResult, error)
	
	// Order monitoring and management
	MonitorOrders() error // Background service to monitor and execute orders
	ExpireOrders() error  // Background service to expire day orders and time-based orders
	UpdateTrailingStops(priceUpdates map[string]float64) error
	ProcessMarketClose(marketCode string) error
	ProcessMarketOpen(marketCode string) error
	
	// Order statistics and analytics
	GetOrderStatistics(userID int) (*domain.OrderStats, error)
	GetExecutionMetrics(userID int, timeframe string) (*repositories.OrderExecutionMetrics, error)
	GetSlippageAnalysis(userID int, symbol string) (*repositories.SlippageAnalysis, error)
	
	// Risk management
	CheckPositionLimits(userID int, order *domain.Order) error
	CheckDailyLimits(userID int) error
	CalculateMarginRequirement(userID int, order *domain.Order) (float64, error)
	ValidateRiskParameters(userID int, order *domain.Order) error
	
	// Integration with other services
	NotifyOrderUpdate(order *domain.Order, updateType OrderUpdateType) error
	UpdatePortfolioOnExecution(execution *domain.OrderExecution) error
	UpdateBalanceOnExecution(execution *domain.OrderExecution) error
}

type CommissionService interface {
	// Commission calculation
	CalculateCommission(userID int, tradeValue float64, orderType domain.OrderType, assetType string) (*domain.CommissionCalculation, error)
	CalculateSlippage(symbol string, quantity int, orderType domain.OrderType) (*domain.Slippage, error)
	
	// Commission structure management
	CreateCommissionStructure(structure *domain.CommissionStructure) error
	GetCommissionStructures() ([]domain.CommissionStructure, error)
	UpdateCommissionStructure(structure *domain.CommissionStructure) error
	
	// User commission management
	GetUserCommissionProfile(userID int) (*domain.UserCommissionProfile, error)
	UpdateUserCommissionProfile(userID int, updates *CommissionProfileUpdates) error
	UpdateUserVolume(userID int, tradeValue float64) error
	EvaluateVIPUpgrade(userID int) error
	
	// Analytics
	GetCommissionStatistics(userID int) (*repositories.CommissionStatistics, error)
	GetCommissionHistory(userID int, startDate, endDate time.Time) ([]domain.CommissionCalculation, error)
	GetSystemCommissionRevenue(startDate, endDate time.Time) (*repositories.CommissionRevenue, error)
	
	// Monthly operations
	ResetMonthlyVolumes() error
	ProcessMonthlyCommissionReport() error
}

type MarketService interface {
	// Market status
	GetMarketStatus(marketCode string) (*domain.MarketStatus, error)
	IsMarketOpen(marketCode string) (bool, error)
	GetMarketCalendar(marketCode string, date time.Time) (*domain.MarketCalendar, error)
	GetNextMarketOpen(marketCode string) (*time.Time, error)
	GetNextMarketClose(marketCode string) (*time.Time, error)
	
	// Trading validation
	CanTrade(marketCode string, symbol string, orderType domain.OrderType) (bool, string, error)
	ValidateMarketHours(marketCode string, orderTime time.Time) (bool, error)
	GetCurrentTradingSession(marketCode string) (*domain.TradingSessionType, error)
	
	// Market conditions
	GetMarketConditions(marketCode string) (*domain.MarketConditions, error)
	UpdateMarketConditions(marketCode string, conditions *domain.MarketConditions) error
	
	// Trading restrictions
	GetTradingRestrictions(marketCode string, symbol *string) ([]domain.TradingRestriction, error)
	CreateTradingRestriction(restriction *domain.TradingRestriction) error
	RemoveTradingRestriction(restrictionID int) error
	
	// Market data permissions
	ValidateMarketDataAccess(userID int, marketCode string, dataType string) (bool, error)
	GrantMarketDataAccess(userID int, marketCode string, dataType string, level string) error
	RevokeMarketDataAccess(userID int, marketCode string, dataType string) error
	
	// Market management
	CreateMarket(market *domain.Market) error
	UpdateMarket(market *domain.Market) error
	CreateMarketHoliday(holiday *domain.MarketHoliday) error
	
	// Background services
	UpdateMarketStatuses() error
	ProcessMarketOpenEvents() error
	ProcessMarketCloseEvents() error
}

type RealTimeService interface {
	// WebSocket connection management
	HandleConnection(connectionID string, userID *int, clientInfo *domain.WebSocketClientInfo) error
	HandleDisconnection(connectionID string) error
	HandleSubscription(connectionID string, request *domain.SubscriptionRequest) (*domain.SubscriptionResponse, error)
	HandleUnsubscription(connectionID string, subscriptionType string, symbols []string) error
	
	// Message broadcasting
	BroadcastPriceUpdate(update *domain.PriceUpdateMessage) error
	BroadcastOrderUpdate(userID int, update *domain.OrderUpdateMessage) error
	BroadcastTradeExecution(userID int, execution *domain.TradeExecutionMessage) error
	BroadcastPositionUpdate(userID int, update *domain.PositionUpdateMessage) error
	BroadcastBalanceUpdate(userID int, update *domain.BalanceUpdateMessage) error
	BroadcastNewsAlert(alert *domain.NewsAlertMessage) error
	BroadcastMarketStatus(status *domain.MarketStatus) error
	
	// Price alerts
	CreatePriceAlert(userID int, alert *PriceAlertRequest) error
	UpdatePriceAlert(userID int, alertID int, updates *PriceAlertUpdates) error
	DeletePriceAlert(userID int, alertID int) error
	GetUserPriceAlerts(userID int) ([]PriceAlert, error)
	CheckPriceAlerts(priceUpdates map[string]float64) error
	
	// Market data services
	GetRealTimePrice(symbol string) (*domain.PriceUpdateMessage, error)
	GetOrderBook(symbol string, depth int) (*domain.OrderBookUpdate, error)
	GetMarketDataSnapshot(symbol string) (*domain.MarketDataSnapshot, error)
	
	// Connection health
	SendHeartbeat(connectionID string) error
	CheckConnections() error
	CleanupExpiredConnections() error
	GetConnectionStats() (*domain.ConnectionStats, error)
	
	// Data provider integration
	StartDataProvider() error
	StopDataProvider() error
	IsDataProviderConnected() bool
}

// Supporting types and requests

type OrderModificationRequest struct {
	Price           *float64  `json:"price,omitempty"`
	StopPrice       *float64  `json:"stop_price,omitempty"`
	Quantity        *int      `json:"quantity,omitempty"`
	TimeInForce     *domain.TimeInForce `json:"time_in_force,omitempty"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	TrailingAmount  *float64  `json:"trailing_amount,omitempty"`
	TrailingPercent *float64  `json:"trailing_percent,omitempty"`
}

type OrderUpdateType string

const (
	OrderUpdateTypeCreated         OrderUpdateType = "CREATED"
	OrderUpdateTypeModified        OrderUpdateType = "MODIFIED"
	OrderUpdateTypeCancelled       OrderUpdateType = "CANCELLED"
	OrderUpdateTypeExecuted        OrderUpdateType = "EXECUTED"
	OrderUpdateTypePartiallyFilled OrderUpdateType = "PARTIALLY_FILLED"
	OrderUpdateTypeExpired         OrderUpdateType = "EXPIRED"
)

type CommissionProfileUpdates struct {
	CommissionStructureID *int `json:"commission_structure_id,omitempty"`
	VIPLevel             *int `json:"vip_level,omitempty"`
}

type PriceAlertRequest struct {
	Symbol          string  `json:"symbol" binding:"required"`
	AlertType       string  `json:"alert_type" binding:"required"` // ABOVE, BELOW, CHANGE_PERCENT
	TriggerPrice    *float64 `json:"trigger_price,omitempty"`
	TriggerPercent  *float64 `json:"trigger_percent,omitempty"`
	Message         string  `json:"message"`
	IsActive        bool    `json:"is_active"`
}

type PriceAlertUpdates struct {
	TriggerPrice   *float64 `json:"trigger_price,omitempty"`
	TriggerPercent *float64 `json:"trigger_percent,omitempty"`
	Message        *string  `json:"message,omitempty"`
	IsActive       *bool    `json:"is_active,omitempty"`
}

type PriceAlert struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	Symbol         string    `json:"symbol"`
	AlertType      string    `json:"alert_type"`
	TriggerPrice   *float64  `json:"trigger_price,omitempty"`
	TriggerPercent *float64  `json:"trigger_percent,omitempty"`
	CurrentPrice   float64   `json:"current_price"`
	Message        string    `json:"message"`
	IsTriggered    bool      `json:"is_triggered"`
	TriggeredAt    *time.Time `json:"triggered_at,omitempty"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Order processing pipeline interface
type OrderProcessor interface {
	ProcessOrder(order *domain.Order) (*domain.OrderExecution, error)
	ValidateOrder(order *domain.Order) error
	CalculateExecution(order *domain.Order, marketPrice float64) (*domain.OrderExecution, error)
	ApplySlippage(execution *domain.OrderExecution) error
	CalculateCommission(execution *domain.OrderExecution) error
	UpdatePortfolio(execution *domain.OrderExecution) error
	NotifyExecution(execution *domain.OrderExecution) error
}

// Risk management interface
type RiskManager interface {
	ValidateOrder(userID int, order *domain.Order) error
	CheckPositionLimits(userID int, order *domain.Order) error
	CheckMarginRequirements(userID int, order *domain.Order) error
	CheckDailyLimits(userID int) error
	CalculateRiskMetrics(userID int) (*RiskMetrics, error)
	MonitorPositions(userID int) (*PositionRisk, error)
}

type RiskMetrics struct {
	UserID              int     `json:"user_id"`
	TotalExposure       float64 `json:"total_exposure"`
	MarginUsed          float64 `json:"margin_used"`
	MarginAvailable     float64 `json:"margin_available"`
	MarginRatio         float64 `json:"margin_ratio"`
	PortfolioVaR        float64 `json:"portfolio_var"`        // Value at Risk
	MaxDrawdown         float64 `json:"max_drawdown"`
	ConcentrationRisk   float64 `json:"concentration_risk"`
	LeverageRatio       float64 `json:"leverage_ratio"`
	RiskScore           int     `json:"risk_score"`           // 1-10 scale
	LastCalculated      time.Time `json:"last_calculated"`
}

type PositionRisk struct {
	UserID             int                    `json:"user_id"`
	Positions          []PositionRiskItem     `json:"positions"`
	TotalRisk          float64                `json:"total_risk"`
	ConcentrationRisk  float64                `json:"concentration_risk"`
	CorrelationRisk    float64                `json:"correlation_risk"`
	LiquidityRisk      float64                `json:"liquidity_risk"`
	Recommendations    []RiskRecommendation   `json:"recommendations"`
}

type PositionRiskItem struct {
	Symbol            string  `json:"symbol"`
	Quantity          int     `json:"quantity"`
	MarketValue       float64 `json:"market_value"`
	PortfolioWeight   float64 `json:"portfolio_weight"`
	Beta              float64 `json:"beta"`
	Volatility        float64 `json:"volatility"`
	VaR               float64 `json:"var"`
	RiskContribution  float64 `json:"risk_contribution"`
	LiquidityScore    int     `json:"liquidity_score"`
}

type RiskRecommendation struct {
	Type        string  `json:"type"`        // REDUCE, DIVERSIFY, HEDGE, etc.
	Symbol      string  `json:"symbol"`
	Action      string  `json:"action"`
	Reason      string  `json:"reason"`
	Priority    string  `json:"priority"`    // HIGH, MEDIUM, LOW
	EstimatedImpact float64 `json:"estimated_impact"`
} 