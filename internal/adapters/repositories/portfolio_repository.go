package repositories

import (
    "stock-simulation-backend/internal/core/domain"
)

type PortfolioRepository interface {
    GetByUserID(userID int) ([]domain.Portfolio, error)
    GetByUserIDAndSymbol(userID, symbol string) (*domain.Portfolio, error)
    Create(portfolio *domain.Portfolio) error
    Update(portfolio *domain.Portfolio) error
    Delete(userID int, symbol string) error
    GetPortfolioValue(userID int) (float64, error)
}