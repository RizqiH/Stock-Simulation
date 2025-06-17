package repositories

import "stock-simulation-backend/internal/core/domain"

type TransactionRepository interface {
	Create(transaction *domain.Transaction) error
	GetByID(id int) (*domain.Transaction, error)
	GetByUserID(userID int, limit, offset int) ([]domain.Transaction, error)
	GetByUserIDAndSymbol(userID int, stockSymbol string, limit int) ([]domain.Transaction, error)
	GetByUserIDAndType(userID int, transactionType string, limit int) ([]domain.Transaction, error)
	GetUserTransactionHistory(userID int, stockSymbol, transactionType string, limit int) ([]domain.Transaction, error)
	GetTotalTransactionsByUser(userID int) (int, error)
}