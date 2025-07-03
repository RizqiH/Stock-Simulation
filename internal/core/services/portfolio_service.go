package services

import (
	"fmt"
	"math/rand"
	"time"
	"stock-simulation-backend/internal/core/domain"
	"stock-simulation-backend/internal/core/ports/repositories"
	"stock-simulation-backend/internal/core/ports/services"
)

type portfolioService struct {
	portfolioRepo repositories.PortfolioRepository
	stockRepo     repositories.StockRepository
}

func NewPortfolioService(
	portfolioRepo repositories.PortfolioRepository,
	stockRepo repositories.StockRepository,
) services.PortfolioService {
	return &portfolioService{
		portfolioRepo: portfolioRepo,
		stockRepo:     stockRepo,
	}
}

func (s *portfolioService) GetUserPortfolio(userID int) (*domain.PortfolioSummary, error) {
	portfolios, err := s.portfolioRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user portfolio: %w", err)
	}

	var portfolioItems []domain.PortfolioItem
	var totalValue, totalCost, totalProfit float64

	for _, portfolio := range portfolios {
		// Get current stock price
		stock, err := s.stockRepo.GetBySymbol(portfolio.StockSymbol)
		if err != nil {
			continue // Skip if stock not found
		}

		currentValue := float64(portfolio.Quantity) * stock.CurrentPrice
		profit := currentValue - portfolio.TotalCost
		profitPct := float64(0)
		if portfolio.TotalCost > 0 {
			profitPct = (profit / portfolio.TotalCost) * 100
		}

		portfolioItem := domain.PortfolioItem{
			StockSymbol:   portfolio.StockSymbol,
			StockName:     stock.Name,
			Quantity:      portfolio.Quantity,
			AveragePrice:  portfolio.AveragePrice,
			CurrentPrice:  stock.CurrentPrice,
			TotalCost:     portfolio.TotalCost,
			CurrentValue:  currentValue,
			ProfitLoss:    profit,
			ProfitLossPct: profitPct,
		}

		portfolioItems = append(portfolioItems, portfolioItem)
		totalValue += currentValue
		totalCost += portfolio.TotalCost
		totalProfit += profit
	}

	totalProfitPct := float64(0)
	if totalCost > 0 {
		totalProfitPct = (totalProfit / totalCost) * 100
	}

	summary := &domain.PortfolioSummary{
		TotalValue:     totalValue,
		TotalCost:      totalCost,
		TotalProfit:    totalProfit,
		TotalProfitPct: totalProfitPct,
		Holdings:       portfolioItems,
	}

	return summary, nil
}

func (s *portfolioService) GetPortfolioPerformance(userID int, period string) (*domain.PortfolioPerformance, error) {
	portfolios, err := s.portfolioRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user portfolio: %w", err)
	}

	var totalCost, totalCurrentValue float64
	var totalProfit float64

	for _, portfolio := range portfolios {
		// Get current stock price
		stock, err := s.stockRepo.GetBySymbol(portfolio.StockSymbol)
		if err != nil {
			continue // Skip if stock not found
		}

		currentValue := float64(portfolio.Quantity) * stock.CurrentPrice
		totalCost += portfolio.TotalCost
		totalCurrentValue += currentValue
		totalProfit += (currentValue - portfolio.TotalCost)
	}

	totalProfitPct := float64(0)
	if totalCost > 0 {
		totalProfitPct = (totalProfit / totalCost) * 100
	}

	performance := &domain.PortfolioPerformance{
		Period:       period,
		StartValue:   totalCost,
		EndValue:     totalCurrentValue,
		Profit:       totalProfit,
		ProfitPct:    totalProfitPct,
		Transactions: len(portfolios),
	}

	return performance, nil
}

// New method for time-based portfolio performance with historical data points
func (s *portfolioService) GetPortfolioPerformanceHistory(userID int, startDate, endDate time.Time) ([]domain.PortfolioDataPoint, error) {
	// Get current portfolio to understand holdings
	portfolios, err := s.portfolioRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user portfolio: %w", err)
	}

	// For now, generate realistic mock data based on current portfolio
	// In a real implementation, this would query historical portfolio values from database
	
	var dataPoints []domain.PortfolioDataPoint
	var totalCost float64

	// Calculate current total cost
	for _, portfolio := range portfolios {
		totalCost += portfolio.TotalCost
	}

	// If no portfolio exists, use default values
	if totalCost == 0 {
		totalCost = 50000 // Default starting value
	}

	// Generate daily data points from start to end date
	currentDate := startDate
	initialValue := totalCost
	currentValue := initialValue
	
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		// Simulate realistic portfolio movement (±2% daily volatility)
		// #nosec G404 -- Using math/rand for portfolio simulation, not cryptographic purposes
		dailyChange := (rand.Float64() - 0.5) * 0.04 // ±2% daily change
		marketTrend := 0.0002 // Small positive trend (about 7% annually)
		
		currentValue *= (1 + dailyChange + marketTrend)
		
		// Ensure value doesn't go negative
		if currentValue < totalCost * 0.5 {
			currentValue = totalCost * 0.5
		}
		
		profitLoss := currentValue - totalCost
		profitLossPct := (profitLoss / totalCost) * 100
		
		dataPoint := domain.PortfolioDataPoint{
			Date:            currentDate,
			TotalValue:      currentValue,
			TotalCost:       totalCost,
			ProfitLoss:      profitLoss,
			ProfitLossPct:   profitLossPct,
			CashValue:       5000,  // Assuming some cash
			InvestmentValue: currentValue - 5000,
		}
		
		dataPoints = append(dataPoints, dataPoint)
		
		// Move to next day
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return dataPoints, nil
}

func (s *portfolioService) GetPortfolioValue(userID int) (float64, error) {
	value, err := s.portfolioRepo.GetPortfolioValue(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get portfolio value: %w", err)
	}

	return value, nil
}

func (s *portfolioService) GetPortfolioSummary(userID int) (*domain.PortfolioSummary, error) {
	summary, err := s.portfolioRepo.GetPortfolioSummary(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio summary: %w", err)
	}

	return summary, nil
}

func (s *portfolioService) UpdatePortfolio(userID int, stockSymbol string, quantity int, averagePrice float64) error {
	portfolio, err := s.portfolioRepo.GetByUserIDAndSymbol(userID, stockSymbol)
	if err != nil {
		return fmt.Errorf("failed to get portfolio: %w", err)
	}

	if portfolio == nil {
		return fmt.Errorf("portfolio item not found")
	}

	portfolio.Quantity = quantity
	portfolio.AveragePrice = averagePrice
	portfolio.TotalCost = float64(quantity) * averagePrice

	err = s.portfolioRepo.Update(portfolio)
	if err != nil {
		return fmt.Errorf("failed to update portfolio: %w", err)
	}

	return nil
}

func (s *portfolioService) GetPortfolioItem(userID int, stockSymbol string) (*domain.Portfolio, error) {
	portfolio, err := s.portfolioRepo.GetByUserIDAndSymbol(userID, stockSymbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio item: %w", err)
	}

	if portfolio == nil {
		return nil, fmt.Errorf("portfolio item not found")
	}

	return portfolio, nil
}