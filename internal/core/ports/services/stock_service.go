package services

import "stock-simulation-backend/internal/core/domain"

type StockService interface {
	GetAllStocks() ([]domain.Stock, error)
	GetStockBySymbol(symbol string) (*domain.Stock, error)
	GetTopStocks(limit int) ([]domain.Stock, error)
	UpdateStockPrice(symbol string, price float64) error
	CreateStock(stock *domain.Stock) error
	GetStockPrice(symbol string) (*domain.StockPrice, error)
}