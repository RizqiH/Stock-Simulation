package services

import (
	"time"
	"stock-simulation-backend/internal/core/domain"
)

type PortfolioService interface {
	GetUserPortfolio(userID int) (*domain.PortfolioSummary, error)
	GetPortfolioPerformance(userID int, period string) (*domain.PortfolioPerformance, error)
	GetPortfolioPerformanceHistory(userID int, startDate, endDate time.Time) ([]domain.PortfolioDataPoint, error)
	GetPortfolioValue(userID int) (float64, error)
	GetPortfolioSummary(userID int) (*domain.PortfolioSummary, error)
	UpdatePortfolio(userID int, stockSymbol string, quantity int, averagePrice float64) error
	GetPortfolioItem(userID int, stockSymbol string) (*domain.Portfolio, error)
}