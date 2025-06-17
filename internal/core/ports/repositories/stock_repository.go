package repositories

import "stock-simulation-backend/internal/core/domain"

type StockRepository interface {
	GetAll() ([]domain.Stock, error)
	GetBySymbol(symbol string) (*domain.Stock, error)
	UpdatePrice(symbol string, price float64) error
	Create(stock *domain.Stock) error
	GetTopStocks(limit int) ([]domain.Stock, error)
	Update(stock *domain.Stock) error
	Delete(symbol string) error
}