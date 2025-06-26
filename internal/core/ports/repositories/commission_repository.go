package repositories

import (
	"stock-simulation-backend/internal/core/domain"
	"time"
)

type CommissionRepository interface {
	// Commission structure management
	CreateCommissionStructure(structure *domain.CommissionStructure) error
	GetCommissionStructureByID(id int) (*domain.CommissionStructure, error)
	GetCommissionStructureByName(name string) (*domain.CommissionStructure, error)
	GetActiveCommissionStructures() ([]domain.CommissionStructure, error)
	UpdateCommissionStructure(structure *domain.CommissionStructure) error
	DeleteCommissionStructure(id int) error
	
	// Commission tier management
	CreateCommissionTier(tier *domain.CommissionTier) error
	GetCommissionTiersByStructureID(structureID int) ([]domain.CommissionTier, error)
	UpdateCommissionTier(tier *domain.CommissionTier) error
	DeleteCommissionTier(id int) error
	
	// User commission profile management
	CreateUserCommissionProfile(profile *domain.UserCommissionProfile) error
	GetUserCommissionProfile(userID int) (*domain.UserCommissionProfile, error)
	UpdateUserCommissionProfile(profile *domain.UserCommissionProfile) error
	UpdateUserVolume(userID int, tradeValue float64) error
	UpdateUserVIPLevel(userID int, vipLevel int) error
	ResetMonthlyVolume() error // Called monthly to reset volumes
	
	// Commission calculation
	CalculateCommission(userID int, tradeValue float64, orderType domain.OrderType, assetType string) (*domain.CommissionCalculation, error)
	GetApplicableTier(userID int, monthlyVolume float64) (*domain.CommissionTier, error)
	
	// Slippage calculation
	CalculateSlippage(symbol string, quantity int, orderType domain.OrderType, marketConditions string) (*domain.Slippage, error)
	GetHistoricalSlippage(symbol string, days int) ([]domain.Slippage, error)
	
	// Analytics and reporting
	GetUserCommissionHistory(userID int, startDate, endDate time.Time) ([]domain.CommissionCalculation, error)
	GetCommissionStatistics(userID int) (*CommissionStatistics, error)
	GetSystemCommissionRevenue(startDate, endDate time.Time) (*CommissionRevenue, error)
}

// CommissionStatistics represents user commission statistics
type CommissionStatistics struct {
	UserID               int     `json:"user_id"`
	TotalCommissionPaid  float64 `json:"total_commission_paid"`
	TotalFeesPaid        float64 `json:"total_fees_paid"`
	AverageCommissionPerTrade float64 `json:"average_commission_per_trade"`
	MonthlyCommission    float64 `json:"monthly_commission"`
	YearlyCommission     float64 `json:"yearly_commission"`
	CurrentTier          *domain.CommissionTier `json:"current_tier,omitempty"`
	NextTierRequirement  *TierRequirement `json:"next_tier_requirement,omitempty"`
	EstimatedSavings     float64 `json:"estimated_savings"` // Savings compared to standard rate
}

// TierRequirement represents requirements to reach next tier
type TierRequirement struct {
	NextTierID          int     `json:"next_tier_id"`
	RequiredVolume      float64 `json:"required_volume"`
	RemainingVolume     float64 `json:"remaining_volume"`
	EstimatedSavings    float64 `json:"estimated_savings"`
	DaysToNextEvaluation int    `json:"days_to_next_evaluation"`
}

// CommissionRevenue represents system-wide commission revenue
type CommissionRevenue struct {
	Period              string  `json:"period"`
	TotalCommission     float64 `json:"total_commission"`
	TotalFees           float64 `json:"total_fees"`
	TotalRevenue        float64 `json:"total_revenue"`
	NumberOfTrades      int     `json:"number_of_trades"`
	NumberOfUsers       int     `json:"number_of_users"`
	AverageCommissionPerTrade float64 `json:"average_commission_per_trade"`
	AverageCommissionPerUser  float64 `json:"average_commission_per_user"`
	RevenueByStructure  map[string]float64 `json:"revenue_by_structure"`
}

type MarketRepository interface {
	// Market management
	CreateMarket(market *domain.Market) error
	GetMarketByCode(code string) (*domain.Market, error)
	GetAllMarkets() ([]domain.Market, error)
	GetMarketsByType(marketType domain.MarketType) ([]domain.Market, error)
	UpdateMarket(market *domain.Market) error
	DeleteMarket(code string) error
	
	// Trading session management
	CreateTradingSession(session *domain.TradingSession) error
	GetTradingSessionsByMarket(marketID int) ([]domain.TradingSession, error)
	GetTradingSessionByType(marketID int, sessionType domain.TradingSessionType) (*domain.TradingSession, error)
	UpdateTradingSession(session *domain.TradingSession) error
	DeleteTradingSession(id int) error
	
	// Market holiday management
	CreateMarketHoliday(holiday *domain.MarketHoliday) error
	GetMarketHolidays(marketID int, year int) ([]domain.MarketHoliday, error)
	GetMarketHolidayByDate(marketID int, date time.Time) (*domain.MarketHoliday, error)
	UpdateMarketHoliday(holiday *domain.MarketHoliday) error
	DeleteMarketHoliday(id int) error
	
	// Market status calculation
	GetMarketStatus(marketCode string) (*domain.MarketStatus, error)
	IsMarketOpen(marketCode string) (bool, error)
	GetNextMarketOpen(marketCode string) (*time.Time, error)
	GetNextMarketClose(marketCode string) (*time.Time, error)
	GetCurrentTradingSession(marketCode string) (*domain.TradingSessionType, error)
	
	// Market calendar
	GetMarketCalendar(marketCode string, date time.Time) (*domain.MarketCalendar, error)
	GetMarketCalendarRange(marketCode string, startDate, endDate time.Time) ([]domain.MarketCalendar, error)
	
	// Market data permissions
	CreateMarketDataPermission(permission *domain.MarketDataPermission) error
	GetUserMarketDataPermissions(userID int) ([]domain.MarketDataPermission, error)
	GetMarketDataPermission(userID int, marketCode string, dataType string) (*domain.MarketDataPermission, error)
	UpdateMarketDataPermission(permission *domain.MarketDataPermission) error
	RevokeMarketDataPermission(userID int, marketCode string, dataType string) error
	
	// Trading restrictions
	CreateTradingRestriction(restriction *domain.TradingRestriction) error
	GetTradingRestrictions(marketCode string, symbol *string) ([]domain.TradingRestriction, error)
	GetActiveTradingRestrictions(marketCode string) ([]domain.TradingRestriction, error)
	UpdateTradingRestriction(restriction *domain.TradingRestriction) error
	RemoveTradingRestriction(id int) error
	
	// Market conditions
	GetMarketConditions(marketCode string) (*domain.MarketConditions, error)
	UpdateMarketConditions(conditions *domain.MarketConditions) error
	
	// Validation and business logic
	CanTrade(marketCode string, symbol string, orderType domain.OrderType) (bool, string, error)
	ValidateMarketHours(marketCode string, orderTime time.Time) (bool, error)
	GetMarketTimeZone(marketCode string) (*time.Location, error)
}

type WebSocketRepository interface {
	// Connection management
	CreateConnection(connection *domain.WebSocketConnection) error
	GetConnection(connectionID string) (*domain.WebSocketConnection, error)
	GetUserConnections(userID int) ([]*domain.WebSocketConnection, error)
	GetAllActiveConnections() ([]*domain.WebSocketConnection, error)
	UpdateConnection(connection *domain.WebSocketConnection) error
	DeleteConnection(connectionID string) error
	
	// Subscription management
	AddSubscription(connectionID string, subscriptionType string, symbols []string) error
	RemoveSubscription(connectionID string, subscriptionType string, symbols []string) error
	GetSubscriptions(connectionID string) (map[string][]string, error)
	GetConnectionsBySubscription(subscriptionType string, symbol string) ([]string, error)
	
	// Connection health
	UpdateHeartbeat(connectionID string) error
	GetExpiredConnections(timeout time.Duration) ([]string, error)
	CleanupExpiredConnections(timeout time.Duration) error
	
	// Statistics
	GetConnectionStats() (*domain.ConnectionStats, error)
	GetConnectionCountByUser(userID int) (int, error)
	GetSubscriptionStats() (map[string]int, error)
	
	// Message history (optional)
	StoreMessage(connectionID string, message *domain.WebSocketMessage) error
	GetMessageHistory(connectionID string, limit int) ([]domain.WebSocketMessage, error)
}

type RealTimeDataRepository interface {
	// Price data storage
	StorePriceUpdate(update *domain.PriceUpdateMessage) error
	GetLatestPrice(symbol string) (*domain.PriceUpdateMessage, error)
	GetPriceHistory(symbol string, minutes int) ([]domain.PriceUpdateMessage, error)
	
	// Order book data
	StoreOrderBookUpdate(update *domain.OrderBookUpdate) error
	GetLatestOrderBook(symbol string) (*domain.OrderBookUpdate, error)
	GetOrderBookHistory(symbol string, minutes int) ([]domain.OrderBookUpdate, error)
	
	// Market data snapshots
	StoreMarketSnapshot(snapshot *domain.MarketDataSnapshot) error
	GetMarketSnapshot(symbol string) (*domain.MarketDataSnapshot, error)
	GetMarketSnapshots(symbols []string) (map[string]*domain.MarketDataSnapshot, error)
	
	// Data cleanup
	CleanupOldPriceData(olderThanHours int) error
	CleanupOldOrderBookData(olderThanHours int) error
	
	// Analytics
	GetPriceVolatility(symbol string, minutes int) (float64, error)
	GetAverageSpread(symbol string, minutes int) (float64, error)
	GetTradingVolume(symbol string, timeframe string) (int64, error)
}

// Additional helper types for market operations
type MarketSessionInfo struct {
	MarketCode      string                     `json:"market_code"`
	CurrentSession  *domain.TradingSessionType `json:"current_session,omitempty"`
	IsOpen          bool                       `json:"is_open"`
	NextOpen        *time.Time                 `json:"next_open,omitempty"`
	NextClose       *time.Time                 `json:"next_close,omitempty"`
	TimeZone        string                     `json:"timezone"`
	LocalTime       time.Time                  `json:"local_time"`
}

type TradingHours struct {
	MarketCode   string                 `json:"market_code"`
	Date         time.Time              `json:"date"`
	Sessions     []SessionInfo          `json:"sessions"`
	IsHoliday    bool                   `json:"is_holiday"`
	SpecialHours *domain.SpecialHours   `json:"special_hours,omitempty"`
}

type SessionInfo struct {
	Type      domain.TradingSessionType `json:"type"`
	StartTime time.Time                 `json:"start_time"`
	EndTime   time.Time                 `json:"end_time"`
	IsActive  bool                      `json:"is_active"`
} 