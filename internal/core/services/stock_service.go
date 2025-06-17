package services

import (
	"fmt"
	"stock-simulation-backend/internal/core/domain"
	"stock-simulation-backend/internal/core/ports/repositories"
	"stock-simulation-backend/internal/core/ports/services"
	"time"
)

type stockService struct {
	stockRepo repositories.StockRepository
}

func NewStockService(stockRepo repositories.StockRepository) services.StockService {
	return &stockService{
		stockRepo: stockRepo,
	}
}

func (s *stockService) GetAllStocks() ([]domain.Stock, error) {
	stocks, err := s.stockRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all stocks: %w", err)
	}

	return stocks, nil
}

func (s *stockService) GetStockBySymbol(symbol string) (*domain.Stock, error) {
	stock, err := s.stockRepo.GetBySymbol(symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock by symbol: %w", err)
	}

	return stock, nil
}

func (s *stockService) GetTopStocks(limit int) ([]domain.Stock, error) {
	stocks, err := s.stockRepo.GetTopStocks(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top stocks: %w", err)
	}

	return stocks, nil
}

func (s *stockService) UpdateStockPrice(symbol string, price float64) error {
	if price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}

	err := s.stockRepo.UpdatePrice(symbol, price)
	if err != nil {
		return fmt.Errorf("failed to update stock price: %w", err)
	}

	return nil
}

func (s *stockService) CreateStock(stock *domain.Stock) error {
	if stock.Symbol == "" {
		return fmt.Errorf("stock symbol is required")
	}
	if stock.Name == "" {
		return fmt.Errorf("stock name is required")
	}
	if stock.CurrentPrice <= 0 {
		return fmt.Errorf("current price must be greater than 0")
	}

	// Check if stock already exists
	existingStock, err := s.stockRepo.GetBySymbol(stock.Symbol)
	if err == nil && existingStock != nil {
		return fmt.Errorf("stock with symbol %s already exists", stock.Symbol)
	}

	stock.UpdatedAt = time.Now()
	err = s.stockRepo.Create(stock)
	if err != nil {
		return fmt.Errorf("failed to create stock: %w", err)
	}

	return nil
}

func (s *stockService) GetStockPrice(symbol string) (*domain.StockPrice, error) {
	stock, err := s.stockRepo.GetBySymbol(symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock: %w", err)
	}

	// Calculate change and change percentage
	change := stock.CurrentPrice - stock.OpenPrice
	changePct := float64(0)
	if stock.OpenPrice > 0 {
		changePct = (change / stock.OpenPrice) * 100
	}

	stockPrice := &domain.StockPrice{
		Symbol:    stock.Symbol,
		Price:     stock.CurrentPrice,
		Change:    change,
		ChangePct: changePct,
		UpdatedAt: stock.UpdatedAt,
	}

	return stockPrice, nil
}