package services

import "stock-simulation-backend/internal/core/domain"

type StockService interface {
	GetAllStocks() ([]domain.Stock, error)
	GetStockBySymbol(symbol string) (*domain.Stock, error)
	GetTopStocks(limit int) ([]domain.Stock, error)
	CreateStock(stock *domain.Stock) error
	UpdateStock(stock *domain.Stock) error
	DeleteStock(symbol string) error
	UpdateStockPrice(symbol string, price float64) error
	SimulateMarketMovement() error
	GetStockPrice(symbol string) (*domain.StockPrice, error)
}