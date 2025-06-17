package mysql

import (
	"database/sql"
	"fmt"
	"stock-simulation-backend/internal/core/domain"
	"stock-simulation-backend/internal/core/ports/repositories"
)

type portfolioRepository struct {
	db *sql.DB
}

func NewPortfolioRepository(db *sql.DB) repositories.PortfolioRepository {
	return &portfolioRepository{db: db}
}

func (r *portfolioRepository) Create(portfolio *domain.Portfolio) error {
	query := `
		INSERT INTO portfolios (user_id, stock_symbol, quantity, average_price, total_cost, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW())
	`
	result, err := r.db.Exec(query, portfolio.UserID, portfolio.StockSymbol,
		portfolio.Quantity, portfolio.AveragePrice, portfolio.TotalCost)
	if err != nil {
		return fmt.Errorf("failed to create portfolio: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get portfolio ID: %w", err)
	}

	portfolio.ID = int(id)
	return nil
}

func (r *portfolioRepository) GetByUserID(userID int) ([]domain.Portfolio, error) {
	query := `
		SELECT id, user_id, stock_symbol, quantity, average_price, total_cost, updated_at
		FROM portfolios 
		WHERE user_id = ? AND quantity > 0
		ORDER BY stock_symbol
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user portfolio: %w", err)
	}
	defer rows.Close()

	var portfolios []domain.Portfolio
	for rows.Next() {
		var portfolio domain.Portfolio
		err := rows.Scan(&portfolio.ID, &portfolio.UserID, &portfolio.StockSymbol,
			&portfolio.Quantity, &portfolio.AveragePrice, &portfolio.TotalCost,
			&portfolio.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan portfolio: %w", err)
		}
		portfolios = append(portfolios, portfolio)
	}

	return portfolios, nil
}

func (r *portfolioRepository) GetByUserIDAndSymbol(userID int, stockSymbol string) (*domain.Portfolio, error) {
	query := `
		SELECT id, user_id, stock_symbol, quantity, average_price, total_cost, updated_at
		FROM portfolios 
		WHERE user_id = ? AND stock_symbol = ?
	`
	var portfolio domain.Portfolio
	err := r.db.QueryRow(query, userID, stockSymbol).Scan(
		&portfolio.ID, &portfolio.UserID, &portfolio.StockSymbol,
		&portfolio.Quantity, &portfolio.AveragePrice, &portfolio.TotalCost,
		&portfolio.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Portfolio not found, return nil without error
		}
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}
	return &portfolio, nil
}

func (r *portfolioRepository) Update(portfolio *domain.Portfolio) error {
	query := `
		UPDATE portfolios 
		SET quantity = ?, average_price = ?, total_cost = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.Exec(query, portfolio.Quantity, portfolio.AveragePrice,
		portfolio.TotalCost, portfolio.ID)
	if err != nil {
		return fmt.Errorf("failed to update portfolio: %w", err)
	}
	return nil
}

func (r *portfolioRepository) Delete(userID int, stockSymbol string) error {
	query := `DELETE FROM portfolios WHERE user_id = ? AND stock_symbol = ?`
	_, err := r.db.Exec(query, userID, stockSymbol)
	if err != nil {
		return fmt.Errorf("failed to delete portfolio: %w", err)
	}
	return nil
}

func (r *portfolioRepository) GetPortfolioValue(userID int) (float64, error) {
	query := `
		SELECT COALESCE(SUM(p.quantity * s.current_price), 0) as total_value
		FROM portfolios p
		JOIN stocks s ON p.stock_symbol = s.symbol
		WHERE p.user_id = ? AND p.quantity > 0
	`
	var totalValue float64
	err := r.db.QueryRow(query, userID).Scan(&totalValue)
	if err != nil {
		return 0, fmt.Errorf("failed to get portfolio value: %w", err)
	}
	return totalValue, nil
}

func (r *portfolioRepository) GetPortfolioSummary(userID int) (*domain.PortfolioSummary, error) {
	query := `
		SELECT 
			p.stock_symbol,
			s.name as stock_name,
			p.quantity,
			p.average_price,
			s.current_price,
			p.total_cost,
			(p.quantity * s.current_price) as current_value,
			((p.quantity * s.current_price) - p.total_cost) as profit_loss,
			(((p.quantity * s.current_price) - p.total_cost) / p.total_cost * 100) as profit_loss_pct
		FROM portfolios p
		JOIN stocks s ON p.stock_symbol = s.symbol
		WHERE p.user_id = ? AND p.quantity > 0
		ORDER BY p.stock_symbol
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio summary: %w", err)
	}
	defer rows.Close()

	var holdings []domain.PortfolioItem
	var totalValue, totalCost float64

	for rows.Next() {
		var item domain.PortfolioItem
		err := rows.Scan(&item.StockSymbol, &item.StockName, &item.Quantity,
			&item.AveragePrice, &item.CurrentPrice, &item.TotalCost,
			&item.CurrentValue, &item.ProfitLoss, &item.ProfitLossPct)
		if err != nil {
			return nil, fmt.Errorf("failed to scan portfolio item: %w", err)
		}
		holdings = append(holdings, item)
		totalValue += item.CurrentValue
		totalCost += item.TotalCost
	}

	totalProfit := totalValue - totalCost
	totalProfitPct := float64(0)
	if totalCost > 0 {
		totalProfitPct = (totalProfit / totalCost) * 100
	}

	summary := &domain.PortfolioSummary{
		TotalValue:     totalValue,
		TotalCost:      totalCost,
		TotalProfit:    totalProfit,
		TotalProfitPct: totalProfitPct,
		Holdings:       holdings,
	}

	return summary, nil
}