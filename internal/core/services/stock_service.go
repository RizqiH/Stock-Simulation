package services

import (
	"fmt"
	"math"
	"math/rand"
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

func (s *stockService) UpdateStock(stock *domain.Stock) error {
	return s.stockRepo.Update(stock)
}

func (s *stockService) DeleteStock(symbol string) error {
	return s.stockRepo.Delete(symbol)
}

// UpdateStockPrice updates a specific stock's price
func (s *stockService) UpdateStockPrice(symbol string, price float64) error {
	// Validate stock exists
	_, err := s.stockRepo.GetBySymbol(symbol)
	if err != nil {
		return fmt.Errorf("stock not found: %s", symbol)
	}

	// Update price
	err = s.stockRepo.UpdatePrice(symbol, price)
	if err != nil {
		return fmt.Errorf("failed to update stock price: %w", err)
	}

	fmt.Printf("📈 Stock price updated: %s → $%.2f\n", symbol, price)
	return nil
}

// SimulateMarketMovement randomly updates all stock prices
func (s *stockService) SimulateMarketMovement() error {
	// Get all stocks
	stocks, err := s.stockRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get stocks: %w", err)
	}

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	fmt.Println("📊 Simulating market movement...")

	// Update each stock with random price movement
	for _, stock := range stocks {
		newPrice := s.generateRandomPrice(stock.CurrentPrice)
		
		err = s.stockRepo.UpdatePrice(stock.Symbol, newPrice)
		if err != nil {
			fmt.Printf("⚠️ Failed to update %s: %v\n", stock.Symbol, err)
			continue
		}

		// Calculate change
		change := newPrice - stock.CurrentPrice
		changePercent := (change / stock.CurrentPrice) * 100
		
		// Display update
		var indicator string
		if change > 0 {
			indicator = "📈"
		} else if change < 0 {
			indicator = "📉"
		} else {
			indicator = "➡️"
		}

		fmt.Printf("%s %s: $%.2f → $%.2f (%.2f%%)\n",
			indicator, stock.Symbol, stock.CurrentPrice, newPrice, changePercent)
	}

	return nil
}

// generateRandomPrice creates a new price with realistic market movement
func (s *stockService) generateRandomPrice(currentPrice float64) float64 {
	// Generate price movement between -3% to +3%
	maxChangePercent := 3.0
	changePercent := (rand.Float64() - 0.5) * 2 * maxChangePercent

	// Apply volatility factor
	volatilityFactor := 1.0
	if rand.Float64() < 0.1 { // 10% chance of high volatility
		volatilityFactor = 2.0
	} else if rand.Float64() < 0.05 { // 5% chance of very high volatility
		volatilityFactor = 3.0
	}

	changePercent *= volatilityFactor

	// Calculate new price
	change := currentPrice * (changePercent / 100)
	newPrice := currentPrice + change

	// Ensure price doesn't go below $0.01
	if newPrice < 0.01 {
		newPrice = 0.01
	}

	// Round to 2 decimal places
	return math.Round(newPrice*100) / 100
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
	change := stock.CurrentPrice - stock.PreviousClose
	changePct := float64(0)
	if stock.PreviousClose > 0 {
		changePct = (change / stock.PreviousClose) * 100
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