package mysql

import (
	"database/sql"
	"fmt"
	"stock-simulation-backend/internal/core/domain"
	"stock-simulation-backend/internal/core/ports/repositories"
)

type stockRepository struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) repositories.StockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) GetAll() ([]domain.Stock, error) {
	query := `
		SELECT id, symbol, name, current_price, previous_close, volume, market_cap, updated_at
		FROM stocks
		ORDER BY symbol
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all stocks: %w", err)
	}
	defer rows.Close()

	var stocks []domain.Stock
	for rows.Next() {
		var stock domain.Stock
		err := rows.Scan(&stock.ID, &stock.Symbol, &stock.Name, &stock.CurrentPrice,
			&stock.PreviousClose, &stock.Volume, &stock.MarketCap, &stock.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock: %w", err)
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

func (r *stockRepository) GetBySymbol(symbol string) (*domain.Stock, error) {
	query := `
		SELECT id, symbol, name, current_price, previous_close, volume, market_cap, updated_at
		FROM stocks WHERE symbol = ?
	`
	var stock domain.Stock
	err := r.db.QueryRow(query, symbol).Scan(
		&stock.ID, &stock.Symbol, &stock.Name, &stock.CurrentPrice,
		&stock.PreviousClose, &stock.Volume, &stock.MarketCap, &stock.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("stock not found")
		}
		return nil, fmt.Errorf("failed to get stock: %w", err)
	}
	return &stock, nil
}

func (r *stockRepository) UpdatePrice(symbol string, price float64) error {
	query := `UPDATE stocks SET current_price = ?, updated_at = NOW() WHERE symbol = ?`
	_, err := r.db.Exec(query, price, symbol)
	if err != nil {
		return fmt.Errorf("failed to update stock price: %w", err)
	}
	return nil
}

func (r *stockRepository) Create(stock *domain.Stock) error {
	query := `
		INSERT INTO stocks (symbol, name, current_price, previous_close, volume, market_cap, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())
	`
	result, err := r.db.Exec(query, stock.Symbol, stock.Name, stock.CurrentPrice,
		stock.PreviousClose, stock.Volume, stock.MarketCap)
	if err != nil {
		return fmt.Errorf("failed to create stock: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get stock ID: %w", err)
	}

	stock.ID = int(id)
	return nil
}

func (r *stockRepository) GetTopStocks(limit int) ([]domain.Stock, error) {
	query := `
		SELECT id, symbol, name, current_price, previous_close, volume, market_cap, updated_at
		FROM stocks
		ORDER BY market_cap DESC
		LIMIT ?
	`
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top stocks: %w", err)
	}
	defer rows.Close()

	var stocks []domain.Stock
	for rows.Next() {
		var stock domain.Stock
		err := rows.Scan(&stock.ID, &stock.Symbol, &stock.Name, &stock.CurrentPrice,
			&stock.PreviousClose, &stock.Volume, &stock.MarketCap, &stock.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock: %w", err)
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

func (r *stockRepository) Update(stock *domain.Stock) error {
	query := `
		UPDATE stocks 
		SET name = ?, current_price = ?, previous_close = ?, volume = ?, market_cap = ?, updated_at = NOW()
		WHERE symbol = ?
	`
	_, err := r.db.Exec(query, stock.Name, stock.CurrentPrice, stock.PreviousClose,
		stock.Volume, stock.MarketCap, stock.Symbol)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}
	return nil
}

func (r *stockRepository) Delete(symbol string) error {
	query := `DELETE FROM stocks WHERE symbol = ?`
	_, err := r.db.Exec(query, symbol)
	if err != nil {
		return fmt.Errorf("failed to delete stock: %w", err)
	}
	return nil
}