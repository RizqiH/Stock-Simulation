package repositories

import "stock-simulation-backend/internal/core/domain"

type PortfolioRepository interface {
	Create(portfolio *domain.Portfolio) error
	GetByUserID(userID int) ([]domain.Portfolio, error)
	GetByUserIDAndSymbol(userID int, stockSymbol string) (*domain.Portfolio, error)
	Update(portfolio *domain.Portfolio) error
	Delete(userID int, stockSymbol string) error
	GetPortfolioValue(userID int) (float64, error)
	GetPortfolioSummary(userID int) (*domain.PortfolioSummary, error)
}