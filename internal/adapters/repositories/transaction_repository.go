package repositories

import (
    "stock-simulation-backend/internal/core/domain"
)

type TransactionRepository interface {
    Create(transaction *domain.Transaction) error
    GetByUserID(userID int, limit, offset int) ([]domain.Transaction, error)
    GetByUserIDAndSymbol(userID int, symbol string) ([]domain.Transaction, error)
    GetTotalTransactionsByUser(userID int) (int, error)
}