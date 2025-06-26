package repositories

import (
	"stock-simulation-backend/internal/core/domain"
	"time"
)

type HistoricalPriceRepository interface {
	// Create new historical price record
	Create(price *domain.HistoricalPrice) error
	
	// Get historical prices for a symbol within date range
	GetBySymbolAndDateRange(symbol string, startDate, endDate time.Time) ([]domain.HistoricalPrice, error)
	
	// Get historical prices for a symbol with limit (most recent first)
	GetBySymbolWithLimit(symbol string, limit int) ([]domain.HistoricalPrice, error)
	
	// Get historical prices for charting with period
	GetChartData(symbol string, period string) (*domain.ChartData, error)
	
	// Update historical price record
	Update(price *domain.HistoricalPrice) error
	
	// Delete historical price records older than specified date
	DeleteOlderThan(date time.Time) error
	
	// Batch insert historical prices
	BatchInsert(prices []domain.HistoricalPrice) error
	
	// Get latest price for symbol
	GetLatestPrice(symbol string) (*domain.HistoricalPrice, error)
	
	// Get available symbols with historical data
	GetAvailableSymbols() ([]string, error)
} 