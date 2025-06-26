package services

import (
	"stock-simulation-backend/internal/core/domain"
)

type ChartService interface {
	// Get chart data with technical indicators for a symbol and period
	GetChartData(symbol string, period string) (*domain.ChartData, error)
	
	// Get historical prices for a symbol with limit
	GetHistoricalPrices(symbol string, limit int) ([]domain.HistoricalPrice, error)
	
	// Get available symbols with historical data
	GetAvailableSymbols() ([]string, error)
	
	// Add new historical price data
	AddHistoricalPrice(price *domain.HistoricalPrice) error
	
	// Batch add historical prices
	BatchAddHistoricalPrices(prices []domain.HistoricalPrice) error
	
	// Update existing historical price
	UpdateHistoricalPrice(price *domain.HistoricalPrice) error
	
	// Get latest price for symbol
	GetLatestPrice(symbol string) (*domain.HistoricalPrice, error)
	
	// Clean old historical data
	CleanOldData(daysToKeep int) error
} 