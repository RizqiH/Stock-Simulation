package repositories

import (
	"stock-simulation-backend/internal/core/domain"
)

type OrderRepository interface {
	// Create new order
	Create(order *domain.Order) error
	
	// Get order by ID
	GetByID(id int) (*domain.Order, error)
	
	// Get orders by user ID
	GetByUserID(userID int, limit, offset int) ([]domain.Order, error)
	
	// Get orders by user ID and status
	GetByUserIDAndStatus(userID int, status domain.OrderStatus) ([]domain.Order, error)
	
	// Get orders by symbol
	GetBySymbol(symbol string) ([]domain.Order, error)
	
	// Get pending orders by symbol and type
	GetPendingOrdersBySymbolAndType(symbol string, orderType domain.OrderType) ([]domain.Order, error)
	
	// Get all pending orders (for order processing)
	GetAllPendingOrders() ([]domain.Order, error)
	
	// Update order
	Update(order *domain.Order) error
	
	// Update order status
	UpdateStatus(id int, status domain.OrderStatus) error
	
	// Execute order (update status, executed price, executed time)
	ExecuteOrder(id int, executedPrice float64) error
	
	// Cancel order
	CancelOrder(id int) error
	
	// Delete order
	Delete(id int) error
	
	// Get order statistics for user
	GetUserOrderStats(userID int) (*domain.OrderStats, error)
	
	// Get expired orders
	GetExpiredOrders() ([]domain.Order, error)
} 