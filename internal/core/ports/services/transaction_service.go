package services

import "stock-simulation-backend/internal/core/domain"

type TransactionService interface {
	BuyStock(userID int, req *domain.TransactionRequest) (*domain.TransactionResponse, error)
	SellStock(userID int, req *domain.TransactionRequest) (*domain.TransactionResponse, error)
	GetUserTransactions(userID int, limit, offset int) ([]domain.Transaction, error)
	GetTransactionHistory(userID int, stockSymbol, transactionType string, limit int) ([]domain.Transaction, error)
	GetTransactionByID(transactionID int) (*domain.Transaction, error)
}