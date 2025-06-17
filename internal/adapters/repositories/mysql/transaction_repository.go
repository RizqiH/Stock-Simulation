package mysql

import (
	"database/sql"
	"fmt"
	"stock-simulation-backend/internal/core/domain"
	"stock-simulation-backend/internal/core/ports/repositories"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) repositories.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *domain.Transaction) error {
	query := `
		INSERT INTO transactions (user_id, stock_symbol, type, quantity, price, total_amount, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())
	`
	result, err := r.db.Exec(query, transaction.UserID, transaction.StockSymbol,
		transaction.Type, transaction.Quantity, transaction.Price, transaction.TotalAmount)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get transaction ID: %w", err)
	}

	transaction.ID = int(id)
	return nil
}

func (r *transactionRepository) GetByID(id int) (*domain.Transaction, error) {
	query := `
		SELECT id, user_id, stock_symbol, type, quantity, price, total_amount, created_at
		FROM transactions WHERE id = ?
	`
	var transaction domain.Transaction
	err := r.db.QueryRow(query, id).Scan(
		&transaction.ID, &transaction.UserID, &transaction.StockSymbol,
		&transaction.Type, &transaction.Quantity, &transaction.Price,
		&transaction.TotalAmount, &transaction.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return &transaction, nil
}

func (r *transactionRepository) GetByUserID(userID int, limit, offset int) ([]domain.Transaction, error) {
	query := `
		SELECT id, user_id, stock_symbol, type, quantity, price, total_amount, created_at
		FROM transactions 
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user transactions: %w", err)
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.StockSymbol,
			&transaction.Type, &transaction.Quantity, &transaction.Price,
			&transaction.TotalAmount, &transaction.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *transactionRepository) GetByUserIDAndSymbol(userID int, stockSymbol string, limit int) ([]domain.Transaction, error) {
	query := `
		SELECT id, user_id, stock_symbol, type, quantity, price, total_amount, created_at
		FROM transactions 
		WHERE user_id = ? AND stock_symbol = ?
		ORDER BY created_at DESC
		LIMIT ?
	`
	rows, err := r.db.Query(query, userID, stockSymbol, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user transactions by symbol: %w", err)
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.StockSymbol,
			&transaction.Type, &transaction.Quantity, &transaction.Price,
			&transaction.TotalAmount, &transaction.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *transactionRepository) GetByUserIDAndType(userID int, transactionType string, limit int) ([]domain.Transaction, error) {
	query := `
		SELECT id, user_id, stock_symbol, type, quantity, price, total_amount, created_at
		FROM transactions 
		WHERE user_id = ? AND type = ?
		ORDER BY created_at DESC
		LIMIT ?
	`
	rows, err := r.db.Query(query, userID, transactionType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user transactions by type: %w", err)
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.StockSymbol,
			&transaction.Type, &transaction.Quantity, &transaction.Price,
			&transaction.TotalAmount, &transaction.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *transactionRepository) GetUserTransactionHistory(userID int, stockSymbol, transactionType string, limit int) ([]domain.Transaction, error) {
	query := `
		SELECT id, user_id, stock_symbol, type, quantity, price, total_amount, created_at
		FROM transactions 
		WHERE user_id = ?
	`
	args := []interface{}{userID}

	if stockSymbol != "" {
		query += " AND stock_symbol = ?"
		args = append(args, stockSymbol)
	}

	if transactionType != "" {
		query += " AND type = ?"
		args = append(args, transactionType)
	}

	query += " ORDER BY created_at DESC LIMIT ?"
	args = append(args, limit)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction history: %w", err)
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.StockSymbol,
			&transaction.Type, &transaction.Quantity, &transaction.Price,
			&transaction.TotalAmount, &transaction.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *transactionRepository) GetTotalTransactionsByUser(userID int) (int, error) {
	query := `SELECT COUNT(*) FROM transactions WHERE user_id = ?`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total transactions: %w", err)
	}
	return count, nil
}