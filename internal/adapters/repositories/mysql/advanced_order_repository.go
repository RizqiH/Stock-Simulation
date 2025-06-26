package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"stock-simulation-backend/internal/core/domain"
	"stock-simulation-backend/internal/core/ports/repositories"
)

type AdvancedOrderRepository struct {
	db *sql.DB
}

func NewAdvancedOrderRepository(db *sql.DB) repositories.AdvancedOrderRepositoryWithSearch {
	return &AdvancedOrderRepository{db: db}
}

func (r *AdvancedOrderRepository) Create(order *domain.Order) error {
	query := `
		INSERT INTO advanced_orders 
		(user_id, stock_symbol, order_type, side, quantity, price, stop_price, trailing_amount, 
		 trailing_percent, time_in_force, status, market_price, bid_price, ask_price, 
		 commission, fees, expires_at, parent_order_id, linked_order_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		order.UserID, order.StockSymbol, order.OrderType, order.Side, order.Quantity,
		order.Price, order.StopPrice, order.TrailingAmount, order.TrailingPercent,
		order.TimeInForce, order.Status, order.MarketPrice, order.BidPrice, order.AskPrice,
		order.Commission, order.Fees, order.ExpiresAt, order.ParentOrderID, order.LinkedOrderID,
	)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get order ID: %w", err)
	}

	order.ID = int(id)
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	return nil
}

func (r *AdvancedOrderRepository) GetByID(orderID int) (*domain.Order, error) {
	query := `
		SELECT id, user_id, stock_symbol, order_type, side, quantity, price, stop_price,
		       trailing_amount, trailing_percent, time_in_force, status, executed_price,
		       executed_quantity, remaining_quantity, market_price, bid_price, ask_price,
		       commission, fees, spread, executed_at, expires_at, parent_order_id,
		       linked_order_id, created_at, updated_at
		FROM advanced_orders 
		WHERE id = ?
	`

	var order domain.Order
	err := r.db.QueryRow(query, orderID).Scan(
		&order.ID, &order.UserID, &order.StockSymbol, &order.OrderType, &order.Side,
		&order.Quantity, &order.Price, &order.StopPrice, &order.TrailingAmount,
		&order.TrailingPercent, &order.TimeInForce, &order.Status, &order.ExecutedPrice,
		&order.ExecutedQuantity, &order.RemainingQuantity, &order.MarketPrice,
		&order.BidPrice, &order.AskPrice, &order.Commission, &order.Fees, &order.Spread,
		&order.ExecutedAt, &order.ExpiresAt, &order.ParentOrderID, &order.LinkedOrderID,
		&order.CreatedAt, &order.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

func (r *AdvancedOrderRepository) Update(order *domain.Order) error {
	query := `
		UPDATE advanced_orders 
		SET price = ?, stop_price = ?, quantity = ?, time_in_force = ?, status = ?,
		    executed_price = ?, executed_quantity = ?, remaining_quantity = ?,
		    trailing_amount = ?, trailing_percent = ?, expires_at = ?, updated_at = NOW()
		WHERE id = ? AND user_id = ?
	`

	_, err := r.db.Exec(query,
		order.Price, order.StopPrice, order.Quantity, order.TimeInForce, order.Status,
		order.ExecutedPrice, order.ExecutedQuantity, order.RemainingQuantity,
		order.TrailingAmount, order.TrailingPercent, order.ExpiresAt,
		order.ID, order.UserID,
	)

	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	order.UpdatedAt = time.Now()
	return nil
}

func (r *AdvancedOrderRepository) Delete(orderID int) error {
	query := `UPDATE advanced_orders SET status = 'CANCELLED', updated_at = NOW() WHERE id = ?`
	
	result, err := r.db.Exec(query, orderID)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found or already cancelled")
	}

	return nil
}

func (r *AdvancedOrderRepository) Search(userID int, criteria *repositories.OrderSearchCriteria) (*repositories.OrderSearchResult, error) {
	baseQuery := `
		SELECT id, user_id, stock_symbol, order_type, side, quantity, price, stop_price,
		       trailing_amount, trailing_percent, time_in_force, status, executed_price,
		       executed_quantity, remaining_quantity, market_price, bid_price, ask_price,
		       commission, fees, spread, executed_at, expires_at, parent_order_id,
		       linked_order_id, created_at, updated_at
		FROM advanced_orders 
		WHERE user_id = ?
	`

	countQuery := `SELECT COUNT(*) FROM advanced_orders WHERE user_id = ?`
	
	var args []interface{}
	var whereConditions []string
	args = append(args, userID)

	if criteria.Symbol != nil {
		whereConditions = append(whereConditions, "stock_symbol = ?")
		args = append(args, *criteria.Symbol)
	}

	if criteria.OrderType != nil {
		whereConditions = append(whereConditions, "order_type = ?")
		args = append(args, *criteria.OrderType)
	}

	if criteria.Status != nil {
		whereConditions = append(whereConditions, "status = ?")
		args = append(args, *criteria.Status)
	}

	if criteria.StartDate != nil {
		whereConditions = append(whereConditions, "created_at >= ?")
		args = append(args, *criteria.StartDate)
	}

	if criteria.EndDate != nil {
		whereConditions = append(whereConditions, "created_at <= ?")
		args = append(args, *criteria.EndDate)
	}

	if len(whereConditions) > 0 {
		whereClause := " AND " + strings.Join(whereConditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	// Get total count
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Add ordering and pagination
	baseQuery += " ORDER BY created_at DESC"
	if criteria.Limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT %d", criteria.Limit)
	}
	if criteria.Offset > 0 {
		baseQuery += fmt.Sprintf(" OFFSET %d", criteria.Offset)
	}

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search orders: %w", err)
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.StockSymbol, &order.OrderType, &order.Side,
			&order.Quantity, &order.Price, &order.StopPrice, &order.TrailingAmount,
			&order.TrailingPercent, &order.TimeInForce, &order.Status, &order.ExecutedPrice,
			&order.ExecutedQuantity, &order.RemainingQuantity, &order.MarketPrice,
			&order.BidPrice, &order.AskPrice, &order.Commission, &order.Fees, &order.Spread,
			&order.ExecutedAt, &order.ExpiresAt, &order.ParentOrderID, &order.LinkedOrderID,
			&order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate orders: %w", err)
	}

	pageSize := criteria.Limit
	if pageSize == 0 {
		pageSize = 20
	}

	totalPages := (total + pageSize - 1) / pageSize
	page := (criteria.Offset / pageSize) + 1

	return &repositories.OrderSearchResult{
		Orders:     orders,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *AdvancedOrderRepository) GetActiveOrdersByUser(userID int) ([]domain.Order, error) {
	query := `
		SELECT id, user_id, stock_symbol, order_type, side, quantity, price, stop_price,
		       trailing_amount, trailing_percent, time_in_force, status, executed_price,
		       executed_quantity, remaining_quantity, market_price, bid_price, ask_price,
		       commission, fees, spread, executed_at, expires_at, parent_order_id,
		       linked_order_id, created_at, updated_at
		FROM advanced_orders 
		WHERE user_id = ? AND status IN ('PENDING', 'PARTIALLY_FILLED')
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active orders: %w", err)
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.StockSymbol, &order.OrderType, &order.Side,
			&order.Quantity, &order.Price, &order.StopPrice, &order.TrailingAmount,
			&order.TrailingPercent, &order.TimeInForce, &order.Status, &order.ExecutedPrice,
			&order.ExecutedQuantity, &order.RemainingQuantity, &order.MarketPrice,
			&order.BidPrice, &order.AskPrice, &order.Commission, &order.Fees, &order.Spread,
			&order.ExecutedAt, &order.ExpiresAt, &order.ParentOrderID, &order.LinkedOrderID,
			&order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *AdvancedOrderRepository) CancelAllOrdersByUser(userID int, symbol *string) (int, error) {
	query := `UPDATE advanced_orders SET status = 'CANCELLED', updated_at = NOW() 
	          WHERE user_id = ? AND status IN ('PENDING', 'PARTIALLY_FILLED')`
	args := []interface{}{userID}

	if symbol != nil {
		query += " AND stock_symbol = ?"
		args = append(args, *symbol)
	}

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to cancel orders: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}

	return int(rowsAffected), nil
}

func (r *AdvancedOrderRepository) GetOrderStatistics(userID int) (*domain.OrderStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_orders,
			COUNT(CASE WHEN status = 'PENDING' THEN 1 END) as pending_orders,
			COUNT(CASE WHEN status = 'EXECUTED' THEN 1 END) as executed_orders,
			COUNT(CASE WHEN status = 'CANCELLED' THEN 1 END) as cancelled_orders,
			COUNT(CASE WHEN status = 'PARTIALLY_FILLED' THEN 1 END) as partially_filled,
			COALESCE(SUM(commission), 0) as total_commission,
			COALESCE(SUM(fees), 0) as total_fees
		FROM advanced_orders 
		WHERE user_id = ?
	`

	var stats domain.OrderStats
	var totalOrders, executedOrders int

	err := r.db.QueryRow(query, userID).Scan(
		&totalOrders,
		&stats.PendingOrders,
		&executedOrders,
		&stats.CancelledOrders,
		&stats.PartiallyFilled,
		&stats.TotalCommission,
		&stats.TotalFees,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get order statistics: %w", err)
	}

	stats.TotalOrders = totalOrders
	stats.ExecutedOrders = executedOrders

	if totalOrders > 0 {
		stats.SuccessRate = float64(executedOrders) / float64(totalOrders) * 100
	}

	// Calculate average execution time (handle NULL with COALESCE)
	timeQuery := `
		SELECT COALESCE(AVG(TIMESTAMPDIFF(MICROSECOND, created_at, executed_at)) / 1000000, 0)
		FROM advanced_orders 
		WHERE user_id = ? AND status = 'EXECUTED' AND executed_at IS NOT NULL
	`
	
	err = r.db.QueryRow(timeQuery, userID).Scan(&stats.AverageExecutionTime)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate execution time: %w", err)
	}

	return &stats, nil
}

func (r *AdvancedOrderRepository) GetExecutionMetrics(userID int, timeframe string) (*repositories.OrderExecutionMetrics, error) {
	// For now, return mock data - this can be implemented with complex queries
	return &repositories.OrderExecutionMetrics{
		UserID:               userID,
		Timeframe:           timeframe,
		TotalOrders:         0,
		ExecutedOrders:      0,
		CancelledOrders:     0,
		PartiallyFilled:     0,
		AverageExecutionTime: 0.0,
		FillRate:            0.0,
		AverageSlippage:     0.0,
		BestExecution:       0.0,
		WorstExecution:      0.0,
		TotalCommission:     0.0,
		TotalFees:          0.0,
	}, nil
}

func (r *AdvancedOrderRepository) GetSlippageAnalysis(userID int, symbol string) (*repositories.SlippageAnalysis, error) {
	// For now, return mock data - this can be implemented with complex queries
	return &repositories.SlippageAnalysis{
		Symbol:               symbol,
		UserID:              userID,
		TotalTrades:         0,
		AverageSlippage:     0.0,
		MedianSlippage:      0.0,
		SlippageStdDev:      0.0,
		BestExecution:       0.0,
		WorstSlippage:       0.0,
		MarketOrderSlippage: 0.0,
		LimitOrderSlippage:  0.0,
	}, nil
}

// Additional methods required by interface
func (r *AdvancedOrderRepository) GetByUserID(userID int, limit, offset int) ([]domain.Order, error) {
	query := `
		SELECT id, user_id, stock_symbol, order_type, side, quantity, price, stop_price,
		       trailing_amount, trailing_percent, time_in_force, status, executed_price,
		       executed_quantity, remaining_quantity, market_price, bid_price, ask_price,
		       commission, fees, spread, executed_at, expires_at, parent_order_id,
		       linked_order_id, created_at, updated_at
		FROM advanced_orders 
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by user: %w", err)
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.StockSymbol, &order.OrderType, &order.Side,
			&order.Quantity, &order.Price, &order.StopPrice, &order.TrailingAmount,
			&order.TrailingPercent, &order.TimeInForce, &order.Status, &order.ExecutedPrice,
			&order.ExecutedQuantity, &order.RemainingQuantity, &order.MarketPrice,
			&order.BidPrice, &order.AskPrice, &order.Commission, &order.Fees, &order.Spread,
			&order.ExecutedAt, &order.ExpiresAt, &order.ParentOrderID, &order.LinkedOrderID,
			&order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *AdvancedOrderRepository) GetByUserIDAndStatus(userID int, status domain.OrderStatus) ([]domain.Order, error) {
	query := `
		SELECT id, user_id, stock_symbol, order_type, side, quantity, price, stop_price,
		       trailing_amount, trailing_percent, time_in_force, status, executed_price,
		       executed_quantity, remaining_quantity, market_price, bid_price, ask_price,
		       commission, fees, spread, executed_at, expires_at, parent_order_id,
		       linked_order_id, created_at, updated_at
		FROM advanced_orders 
		WHERE user_id = ? AND status = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by user and status: %w", err)
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.StockSymbol, &order.OrderType, &order.Side,
			&order.Quantity, &order.Price, &order.StopPrice, &order.TrailingAmount,
			&order.TrailingPercent, &order.TimeInForce, &order.Status, &order.ExecutedPrice,
			&order.ExecutedQuantity, &order.RemainingQuantity, &order.MarketPrice,
			&order.BidPrice, &order.AskPrice, &order.Commission, &order.Fees, &order.Spread,
			&order.ExecutedAt, &order.ExpiresAt, &order.ParentOrderID, &order.LinkedOrderID,
			&order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *AdvancedOrderRepository) CalculateOrderCommission(userID int, order *domain.Order) (float64, error) {
	// Simple commission calculation: $5 base + $0.005 per share
	commission := 5.00 + (float64(order.Quantity) * 0.005)
	if commission < 1.0 {
		commission = 1.0
	}
	return commission, nil
}

// Stub implementations for remaining interface methods
func (r *AdvancedOrderRepository) GetBySymbol(symbol string) ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *AdvancedOrderRepository) GetByMarketAndStatus(marketCode string, status domain.OrderStatus) ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *AdvancedOrderRepository) GetActiveOrders() ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *AdvancedOrderRepository) GetPendingOrders() ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *AdvancedOrderRepository) GetExpiredOrders() ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *AdvancedOrderRepository) GetOrdersForExecution(symbol string, currentPrice float64) ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *AdvancedOrderRepository) GetLimitOrders(symbol string, side domain.OrderSide) ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *AdvancedOrderRepository) GetStopOrders(symbol string) ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *AdvancedOrderRepository) GetTrailingStopOrders() ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *AdvancedOrderRepository) GetOCOOrders(userID int) ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *AdvancedOrderRepository) ExecuteOrder(orderID int, executedPrice float64, executedQuantity int) error {
	return nil
}

func (r *AdvancedOrderRepository) PartialFillOrder(orderID int, filledQuantity int, filledPrice float64) error {
	return nil
}

func (r *AdvancedOrderRepository) CancelOrder(orderID int, reason string) error {
	return nil
}

func (r *AdvancedOrderRepository) ExpireOrder(orderID int) error {
	return nil
}

func (r *AdvancedOrderRepository) UpdateStatus(orderID int, status domain.OrderStatus) error {
	return nil
}

func (r *AdvancedOrderRepository) UpdateExecutionDetails(orderID int, execution *domain.OrderExecution) error {
	return nil
}

func (r *AdvancedOrderRepository) CreateOCOOrders(parentOrder *domain.Order, linkedOrder *domain.Order) error {
	return nil
}

func (r *AdvancedOrderRepository) CancelLinkedOrders(orderID int) error {
	return nil
}

func (r *AdvancedOrderRepository) GetLinkedOrders(orderID int) ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *AdvancedOrderRepository) GetUserOrderStats(userID int) (*domain.OrderStats, error) {
	return &domain.OrderStats{}, nil
}

func (r *AdvancedOrderRepository) GetOrderStatsByDateRange(userID int, startDate, endDate time.Time) (*domain.OrderStats, error) {
	return &domain.OrderStats{}, nil
}

func (r *AdvancedOrderRepository) GetOrderCountByType(userID int) (map[domain.OrderType]int, error) {
	return make(map[domain.OrderType]int), nil
}

func (r *AdvancedOrderRepository) GetAverageExecutionTime(userID int) (float64, error) {
	return 0.0, nil
}

func (r *AdvancedOrderRepository) GetOrderHistory(userID int, limit, offset int) ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *AdvancedOrderRepository) GetOrdersByDateRange(userID int, startDate, endDate time.Time) ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *AdvancedOrderRepository) GetOrdersForSymbolAndDateRange(userID int, symbol string, startDate, endDate time.Time) ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *AdvancedOrderRepository) ValidateOrderConstraints(order *domain.Order) error {
	return nil
}

func (r *AdvancedOrderRepository) CheckDailyOrderLimit(userID int) (bool, error) {
	return true, nil
}

func (r *AdvancedOrderRepository) CheckOrderSizeLimit(userID int, symbol string, quantity int) (bool, error) {
	return true, nil
}

func (r *AdvancedOrderRepository) UpdateTrailingStopPrice(orderID int, newStopPrice float64) error {
	return nil
}

func (r *AdvancedOrderRepository) GetTrailingStopsToUpdate(priceUpdates map[string]float64) ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *AdvancedOrderRepository) UpdateOrderCommission(orderID int, commission, fees float64) error {
	return nil
}

func (r *AdvancedOrderRepository) GetOrdersAwaitingMarketOpen(marketCode string) ([]domain.Order, error) {
	return []domain.Order{}, nil
}

func (r *AdvancedOrderRepository) GetOrdersToExpireAtMarketClose(marketCode string) ([]domain.Order, error) {
	return []domain.Order{}, nil
}

// Methods from AdvancedOrderRepositoryWithSearch interface
func (r *AdvancedOrderRepository) SearchOrders(criteria *repositories.OrderSearchCriteria) (*repositories.OrderSearchResult, error) {
	// Implementation already exists as Search method
	return r.Search(0, criteria)
}

func (r *AdvancedOrderRepository) SearchOrdersByUser(userID int, criteria *repositories.OrderSearchCriteria) (*repositories.OrderSearchResult, error) {
	return r.Search(userID, criteria)
}

func (r *AdvancedOrderRepository) SearchOrdersBySymbol(symbol string, criteria *repositories.OrderSearchCriteria) (*repositories.OrderSearchResult, error) {
	criteria.Symbol = &symbol
	return r.Search(0, criteria)
}

func (r *AdvancedOrderRepository) BulkUpdateStatus(orderIDs []int, status domain.OrderStatus) error {
	return nil
}

func (r *AdvancedOrderRepository) BulkCancelOrders(orderIDs []int, reason string) error {
	return nil
}

func (r *AdvancedOrderRepository) BulkExecuteOrders(executions []domain.OrderExecution) error {
	return nil
} 