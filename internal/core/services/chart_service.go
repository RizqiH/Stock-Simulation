package services

import (
	"fmt"
	"stock-simulation-backend/internal/core/domain"
	"stock-simulation-backend/internal/core/ports/repositories"
	"time"
)

type chartService struct {
	historicalPriceRepo repositories.HistoricalPriceRepository
}

func NewChartService(historicalPriceRepo repositories.HistoricalPriceRepository) *chartService {
	return &chartService{
		historicalPriceRepo: historicalPriceRepo,
	}
}

func (s *chartService) GetChartData(symbol string, period string) (*domain.ChartData, error) {
	chartData, err := s.historicalPriceRepo.GetChartData(symbol, period)
	if err != nil {
		return nil, fmt.Errorf("failed to get chart data: %w", err)
	}

	return chartData, nil
}

func (s *chartService) GetHistoricalPrices(symbol string, limit int) ([]domain.HistoricalPrice, error) {
	if limit <= 0 {
		limit = 30
	}
	if limit > 365 {
		limit = 365 // Maximum 1 year
	}

	prices, err := s.historicalPriceRepo.GetBySymbolWithLimit(symbol, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical prices: %w", err)
	}

	return prices, nil
}

func (s *chartService) GetAvailableSymbols() ([]string, error) {
	symbols, err := s.historicalPriceRepo.GetAvailableSymbols()
	if err != nil {
		return nil, fmt.Errorf("failed to get available symbols: %w", err)
	}

	return symbols, nil
}

func (s *chartService) AddHistoricalPrice(price *domain.HistoricalPrice) error {
	if price == nil {
		return fmt.Errorf("price data is required")
	}

	if price.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}

	if price.Open <= 0 || price.High <= 0 || price.Low <= 0 || price.Close <= 0 {
		return fmt.Errorf("invalid price data: prices must be positive")
	}

	if price.High < price.Low {
		return fmt.Errorf("invalid price data: high price cannot be less than low price")
	}

	if price.Open > price.High || price.Open < price.Low {
		return fmt.Errorf("invalid price data: open price must be between high and low")
	}

	if price.Close > price.High || price.Close < price.Low {
		return fmt.Errorf("invalid price data: close price must be between high and low")
	}

	err := s.historicalPriceRepo.Create(price)
	if err != nil {
		return fmt.Errorf("failed to add historical price: %w", err)
	}

	return nil
}

func (s *chartService) BatchAddHistoricalPrices(prices []domain.HistoricalPrice) error {
	if len(prices) == 0 {
		return fmt.Errorf("no price data provided")
	}

	// Validate all prices before batch insert
	for i, price := range prices {
		if price.Symbol == "" {
			return fmt.Errorf("price at index %d: symbol is required", i)
		}

		if price.Open <= 0 || price.High <= 0 || price.Low <= 0 || price.Close <= 0 {
			return fmt.Errorf("price at index %d: prices must be positive", i)
		}

		if price.High < price.Low {
			return fmt.Errorf("price at index %d: high price cannot be less than low price", i)
		}

		if price.Open > price.High || price.Open < price.Low {
			return fmt.Errorf("price at index %d: open price must be between high and low", i)
		}

		if price.Close > price.High || price.Close < price.Low {
			return fmt.Errorf("price at index %d: close price must be between high and low", i)
		}
	}

	err := s.historicalPriceRepo.BatchInsert(prices)
	if err != nil {
		return fmt.Errorf("failed to batch add historical prices: %w", err)
	}

	return nil
}

func (s *chartService) UpdateHistoricalPrice(price *domain.HistoricalPrice) error {
	if price == nil {
		return fmt.Errorf("price data is required")
	}

	if price.ID <= 0 {
		return fmt.Errorf("price ID is required for update")
	}

	if price.Open <= 0 || price.High <= 0 || price.Low <= 0 || price.Close <= 0 {
		return fmt.Errorf("invalid price data: prices must be positive")
	}

	if price.High < price.Low {
		return fmt.Errorf("invalid price data: high price cannot be less than low price")
	}

	if price.Open > price.High || price.Open < price.Low {
		return fmt.Errorf("invalid price data: open price must be between high and low")
	}

	if price.Close > price.High || price.Close < price.Low {
		return fmt.Errorf("invalid price data: close price must be between high and low")
	}

	err := s.historicalPriceRepo.Update(price)
	if err != nil {
		return fmt.Errorf("failed to update historical price: %w", err)
	}

	return nil
}

func (s *chartService) GetLatestPrice(symbol string) (*domain.HistoricalPrice, error) {
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	price, err := s.historicalPriceRepo.GetLatestPrice(symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest price: %w", err)
	}

	return price, nil
}

func (s *chartService) CleanOldData(daysToKeep int) error {
	if daysToKeep <= 0 {
		return fmt.Errorf("days to keep must be positive")
	}

	cutoffDate := time.Now().AddDate(0, 0, -daysToKeep)
	
	err := s.historicalPriceRepo.DeleteOlderThan(cutoffDate)
	if err != nil {
		return fmt.Errorf("failed to clean old data: %w", err)
	}

	return nil
} 