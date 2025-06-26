package services

import (
	"fmt"
	"time"

	"stock-simulation-backend/internal/core/domain"
	"stock-simulation-backend/internal/core/ports/repositories"
	"stock-simulation-backend/internal/core/ports/services"
)

type AdvancedOrderService struct {
	orderRepo          repositories.AdvancedOrderRepositoryWithSearch
	stockRepo          repositories.StockRepository
	portfolioRepo      repositories.PortfolioRepository
	userRepo           repositories.UserRepository
	transactionService services.TransactionService
}

func NewAdvancedOrderService(
	orderRepo repositories.AdvancedOrderRepositoryWithSearch,
	stockRepo repositories.StockRepository,
	portfolioRepo repositories.PortfolioRepository,
	userRepo repositories.UserRepository,
	transactionService services.TransactionService,
) services.AdvancedOrderService {
	return &AdvancedOrderService{
		orderRepo:          orderRepo,
		stockRepo:          stockRepo,
		portfolioRepo:      portfolioRepo,
		userRepo:           userRepo,
		transactionService: transactionService,
	}
}

func (s *AdvancedOrderService) ValidateOrder(userID int, request *domain.OrderRequest) error {
	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Validate stock exists
	_, err = s.stockRepo.GetBySymbol(request.StockSymbol)
	if err != nil {
		return fmt.Errorf("stock not found: %s", request.StockSymbol)
	}

	// Validate quantity
	if request.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	// Validate price for limit orders
	if request.OrderType == "LIMIT" && (request.Price == nil || *request.Price <= 0) {
		return fmt.Errorf("limit orders require a valid price")
	}

	// Validate stop price for stop orders
	if (request.OrderType == "STOP_LOSS" || request.OrderType == "TAKE_PROFIT") && 
	   (request.StopPrice == nil || *request.StopPrice <= 0) {
		return fmt.Errorf("stop orders require a valid stop price")
	}

	// Validate trailing stop parameters
	if request.OrderType == "TRAILING_STOP" {
		if request.TrailingAmount == nil && request.TrailingPercent == nil {
			return fmt.Errorf("trailing stop orders require either trailing amount or percentage")
		}
		if request.TrailingAmount != nil && *request.TrailingAmount <= 0 {
			return fmt.Errorf("trailing amount must be positive")
		}
		if request.TrailingPercent != nil && (*request.TrailingPercent <= 0 || *request.TrailingPercent >= 100) {
			return fmt.Errorf("trailing percentage must be between 0 and 100")
		}
	}

	return nil
}

func (s *AdvancedOrderService) CreateOrder(userID int, request *domain.OrderRequest) (*domain.Order, error) {
	// Validate the order first
	if err := s.ValidateOrder(userID, request); err != nil {
		return nil, err
	}

	// Get current stock price
	stock, err := s.stockRepo.GetBySymbol(request.StockSymbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock price: %w", err)
	}

	// Create order object
	order := &domain.Order{
		UserID:           userID,
		StockSymbol:      request.StockSymbol,
		OrderType:        domain.OrderType(request.OrderType),
		Side:             domain.OrderSide(request.Side),
		Quantity:         request.Quantity,
		Price:            request.Price,
		StopPrice:        request.StopPrice,
		TrailingAmount:   request.TrailingAmount,
		TrailingPercent:  request.TrailingPercent,
		TimeInForce:      domain.TimeInForce(request.TimeInForce),
		Status:           "PENDING",
		RemainingQuantity: request.Quantity,
		MarketPrice:      stock.CurrentPrice,
		Commission:       s.calculateCommission(request.Quantity, stock.CurrentPrice),
		Fees:            s.calculateFees(request.Quantity, stock.CurrentPrice),
	}

	if request.ExpiresAt != nil {
		order.ExpiresAt = request.ExpiresAt
	}

	// Save to database first
	err = s.orderRepo.Create(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Execute order based on type
	switch request.OrderType {
	case "MARKET":
		// Market orders execute immediately at current price
		err = s.executeOrderTransaction(order, stock.CurrentPrice)
		if err != nil {
			// If execution fails, cancel the order
			s.orderRepo.Delete(order.ID)
			return nil, fmt.Errorf("failed to execute market order: %w", err)
		}
		
	case "LIMIT":
		// Check if limit order can be executed immediately
		fmt.Printf("🔍 Checking LIMIT order execution: Symbol=%s, Side=%s, LimitPrice=%.2f, CurrentPrice=%.2f\n", 
			order.StockSymbol, order.Side, *order.Price, stock.CurrentPrice)
			
		if s.canExecuteLimitOrder(order, stock.CurrentPrice) {
			fmt.Printf("✅ LIMIT order conditions met, executing immediately...\n")
			err = s.executeOrderTransaction(order, *order.Price)
			if err != nil {
				// If execution fails, keep order pending
				fmt.Printf("⚠️ Limit order execution failed, keeping PENDING: %v\n", err)
			} else {
				fmt.Printf("🎯 LIMIT order successfully executed and should appear in transactions\n")
			}
		} else {
			// Keep order pending for future execution
			fmt.Printf("📝 Limit order created and kept PENDING: ID=%d, LimitPrice=%.2f, CurrentPrice=%.2f\n", 
				order.ID, *order.Price, stock.CurrentPrice)
		}
		
	case "STOP_LOSS", "TAKE_PROFIT":
		// Stop orders are kept pending until triggered
		fmt.Printf("📝 Stop order created and kept PENDING: ID=%d, Type=%s, StopPrice=%.2f\n", 
			order.ID, order.OrderType, *order.StopPrice)
			
	default:
		// Other order types (TRAILING_STOP, OCO) are kept pending
		fmt.Printf("📝 Advanced order created and kept PENDING: ID=%d, Type=%s\n", 
			order.ID, order.OrderType)
	}

	return order, nil
}

// Helper function to check if limit order can be executed immediately
func (s *AdvancedOrderService) canExecuteLimitOrder(order *domain.Order, currentPrice float64) bool {
	if order.Price == nil {
		fmt.Printf("❌ Limit order has no price set\n")
		return false
	}
	
	limitPrice := *order.Price
	
	switch order.Side {
	case "BUY":
		// Buy limit order executes if current price <= limit price
		canExecute := currentPrice <= limitPrice
		fmt.Printf("🔍 BUY LIMIT check: current=%.2f <= limit=%.2f ? %v\n", currentPrice, limitPrice, canExecute)
		return canExecute
	case "SELL":
		// Sell limit order executes if current price >= limit price
		canExecute := currentPrice >= limitPrice
		fmt.Printf("🔍 SELL LIMIT check: current=%.2f >= limit=%.2f ? %v\n", currentPrice, limitPrice, canExecute)
		return canExecute
	default:
		fmt.Printf("❌ Unknown order side: %s\n", order.Side)
		return false
	}
}

func (s *AdvancedOrderService) CreateOCOOrder(userID int, parentRequest, linkedRequest *domain.OrderRequest) (*domain.Order, *domain.Order, error) {
	// Validate both orders
	if err := s.ValidateOrder(userID, parentRequest); err != nil {
		return nil, nil, fmt.Errorf("parent order validation failed: %w", err)
	}
	if err := s.ValidateOrder(userID, linkedRequest); err != nil {
		return nil, nil, fmt.Errorf("linked order validation failed: %w", err)
	}

	// Create parent order first
	parentOrder, err := s.CreateOrder(userID, parentRequest)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create parent order: %w", err)
	}

	// Create linked order
	linkedOrder, err := s.CreateOrder(userID, linkedRequest)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create linked order: %w", err)
	}

	// Update both orders to reference each other
	parentOrder.LinkedOrderID = &linkedOrder.ID
	linkedOrder.ParentOrderID = &parentOrder.ID

	// Update in database
	err = s.orderRepo.Update(parentOrder)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update parent order: %w", err)
	}

	err = s.orderRepo.Update(linkedOrder)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update linked order: %w", err)
	}

	return parentOrder, linkedOrder, nil
}

func (s *AdvancedOrderService) ModifyOrder(userID, orderID int, modifications *services.OrderModificationRequest) (*domain.Order, error) {
	// Get existing order
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	// Check if order can be modified
	if order.Status != "PENDING" {
		return nil, fmt.Errorf("cannot modify order with status: %s", order.Status)
	}

	// Apply modifications
	if modifications.Price != nil {
		order.Price = modifications.Price
	}
	if modifications.StopPrice != nil {
		order.StopPrice = modifications.StopPrice
	}
	if modifications.Quantity != nil {
		order.Quantity = *modifications.Quantity
		order.RemainingQuantity = *modifications.Quantity
	}
	if modifications.TimeInForce != nil {
		order.TimeInForce = domain.TimeInForce(*modifications.TimeInForce)
	}
	if modifications.TrailingAmount != nil {
		order.TrailingAmount = modifications.TrailingAmount
	}
	if modifications.TrailingPercent != nil {
		order.TrailingPercent = modifications.TrailingPercent
	}

	// Update in database
	err = s.orderRepo.Update(order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return order, nil
}

func (s *AdvancedOrderService) CancelOrder(userID, orderID int) error {
	// Get order to verify ownership and status
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}

	if order.Status != "PENDING" && order.Status != "PARTIALLY_FILLED" {
		return fmt.Errorf("cannot cancel order with status: %s", order.Status)
	}

	// Verify ownership
	if order.UserID != userID {
		return fmt.Errorf("order does not belong to user")
	}

	// Cancel the order
	err = s.orderRepo.Delete(orderID)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	// If this is part of an OCO order, cancel the linked order too
	if order.LinkedOrderID != nil {
		err = s.orderRepo.Delete(*order.LinkedOrderID)
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: failed to cancel linked order %d: %v\n", *order.LinkedOrderID, err)
		}
	}

	return nil
}

func (s *AdvancedOrderService) CancelAllOrders(userID int, symbol *string) (int, error) {
	return s.orderRepo.CancelAllOrdersByUser(userID, symbol)
}

func (s *AdvancedOrderService) GetOrderByID(userID, orderID int) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	
	// Verify ownership
	if order.UserID != userID {
		return nil, fmt.Errorf("order not found")
	}
	
	return order, nil
}

func (s *AdvancedOrderService) SearchOrders(userID int, criteria *repositories.OrderSearchCriteria) (*repositories.OrderSearchResult, error) {
	return s.orderRepo.Search(userID, criteria)
}

func (s *AdvancedOrderService) GetActiveOrders(userID int) ([]domain.Order, error) {
	return s.orderRepo.GetActiveOrdersByUser(userID)
}

func (s *AdvancedOrderService) GetOrderStatistics(userID int) (*domain.OrderStats, error) {
	return s.orderRepo.GetOrderStatistics(userID)
}

func (s *AdvancedOrderService) GetExecutionMetrics(userID int, timeframe string) (*repositories.OrderExecutionMetrics, error) {
	return s.orderRepo.GetExecutionMetrics(userID, timeframe)
}

func (s *AdvancedOrderService) GetSlippageAnalysis(userID int, symbol string) (*repositories.SlippageAnalysis, error) {
	return s.orderRepo.GetSlippageAnalysis(userID, symbol)
}

// Helper functions
func (s *AdvancedOrderService) calculateCommission(quantity int, price float64) float64 {
	// Basic commission calculation: $0.005 per share, minimum $1
	commission := float64(quantity) * 0.005
	if commission < 1.0 {
		commission = 1.0
	}
	return commission
}

func (s *AdvancedOrderService) calculateFees(quantity int, price float64) float64 {
	// Basic fee calculation: 0.1% of trade value, minimum $0.50
	tradeValue := float64(quantity) * price
	fees := tradeValue * 0.001
	if fees < 0.50 {
		fees = 0.50
	}
	return fees
}

// Add missing interface method
func (s *AdvancedOrderService) CalculateMarginRequirement(userID int, order *domain.Order) (float64, error) {
	// Basic margin calculation: 50% of stock value for margin accounts
	tradeValue := float64(order.Quantity) * order.MarketPrice
	marginRequirement := tradeValue * 0.5 // 50% margin requirement
	return marginRequirement, nil
}

// Add missing interface methods with stub implementations
func (s *AdvancedOrderService) ExecuteOrder(orderID int, marketPrice float64) (*domain.OrderExecution, error) {
	return &domain.OrderExecution{}, nil
}

func (s *AdvancedOrderService) ExecuteMarketOrders(symbol string, currentPrice float64) ([]domain.OrderExecution, error) {
	return []domain.OrderExecution{}, nil
}

func (s *AdvancedOrderService) ExecuteLimitOrders(symbol string, currentPrice float64) ([]domain.OrderExecution, error) {
	return []domain.OrderExecution{}, nil
}

func (s *AdvancedOrderService) ExecuteStopOrders(symbol string, currentPrice float64) ([]domain.OrderExecution, error) {
	return []domain.OrderExecution{}, nil
}

func (s *AdvancedOrderService) ExecuteTrailingStops(priceUpdates map[string]float64) ([]domain.OrderExecution, error) {
	return []domain.OrderExecution{}, nil
}

func (s *AdvancedOrderService) ValidateBuyingPower(userID int, order *domain.Order) error {
	return nil
}

func (s *AdvancedOrderService) ValidatePosition(userID int, order *domain.Order) error {
	return nil
}

func (s *AdvancedOrderService) ValidateMarketHours(order *domain.Order) error {
	return nil
}

func (s *AdvancedOrderService) ValidateOrderLimits(userID int, order *domain.Order) error {
	return nil
}

func (s *AdvancedOrderService) GetUserOrders(userID int, status *domain.OrderStatus, limit, offset int) ([]domain.Order, error) {
	if status != nil {
		return s.orderRepo.GetByUserIDAndStatus(userID, *status)
	}
	return s.orderRepo.GetByUserID(userID, limit, offset)
}

func (s *AdvancedOrderService) GetOrderHistory(userID int, startDate, endDate *time.Time, limit, offset int) ([]domain.Order, error) {
	return s.orderRepo.GetOrderHistory(userID, limit, offset)
}

func (s *AdvancedOrderService) MonitorOrders() error {
	return nil
}

func (s *AdvancedOrderService) ExpireOrders() error {
	return nil
}

func (s *AdvancedOrderService) UpdateTrailingStops(priceUpdates map[string]float64) error {
	return nil
}

func (s *AdvancedOrderService) ProcessMarketClose(marketCode string) error {
	return nil
}

func (s *AdvancedOrderService) ProcessMarketOpen(marketCode string) error {
	return nil
}

func (s *AdvancedOrderService) CheckPositionLimits(userID int, order *domain.Order) error {
	return nil
}

func (s *AdvancedOrderService) CheckDailyLimits(userID int) error {
	return nil
}

func (s *AdvancedOrderService) ValidateRiskParameters(userID int, order *domain.Order) error {
	return nil
}

func (s *AdvancedOrderService) NotifyOrderUpdate(order *domain.Order, updateType services.OrderUpdateType) error {
	return nil
}

func (s *AdvancedOrderService) UpdatePortfolioOnExecution(execution *domain.OrderExecution) error {
	return nil
}

// New method to execute order and update portfolio/balance
func (s *AdvancedOrderService) executeOrderTransaction(order *domain.Order, executionPrice float64) error {
	// Update order status to executed
	order.Status = "EXECUTED"
	order.ExecutedPrice = &executionPrice
	order.ExecutedQuantity = order.Quantity
	order.RemainingQuantity = 0
	now := time.Now()
	order.ExecutedAt = &now

	// Use transaction service to create proper transaction records
	// This ensures orders appear in transaction history
	transactionRequest := &domain.TransactionRequest{
		StockSymbol: order.StockSymbol,
		Quantity:    order.Quantity,
	}

	var err error
	switch order.Side {
	case "BUY":
		fmt.Printf("🔄 Executing BUY order via transaction service: %s x%d @ $%.2f\n", 
			order.StockSymbol, order.Quantity, executionPrice)
		_, err = s.transactionService.BuyStock(order.UserID, transactionRequest)
		if err != nil {
			return fmt.Errorf("failed to execute buy transaction: %w", err)
		}

	case "SELL":
		fmt.Printf("🔄 Executing SELL order via transaction service: %s x%d @ $%.2f\n", 
			order.StockSymbol, order.Quantity, executionPrice)
		_, err = s.transactionService.SellStock(order.UserID, transactionRequest)
		if err != nil {
			return fmt.Errorf("failed to execute sell transaction: %w", err)
		}

	case "SHORT":
		// For SHORT: Add money, create negative position (custom logic)
		err = s.processShortTransaction(order.UserID, order.StockSymbol, order.Quantity, 
			float64(order.Quantity)*executionPrice-order.Commission-order.Fees)
		if err != nil {
			return fmt.Errorf("failed to process short transaction: %w", err)
		}

	case "COVER":
		// For COVER: Deduct money, reduce negative position (custom logic)
		err = s.processCoverTransaction(order.UserID, order.StockSymbol, order.Quantity, 
			float64(order.Quantity)*executionPrice+order.Commission+order.Fees)
		if err != nil {
			return fmt.Errorf("failed to process cover transaction: %w", err)
		}
	}

	// Update order in database
	err = s.orderRepo.Update(order)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	fmt.Printf("✅ Order executed successfully: ID=%d, Type=%s, Side=%s\n", 
		order.ID, order.OrderType, order.Side)
	return nil
}

func (s *AdvancedOrderService) processBuyTransaction(userID int, symbol string, quantity int, totalCost float64) error {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Check if user has enough balance
	if user.Balance < totalCost {
		return fmt.Errorf("insufficient balance: required %.2f, available %.2f", totalCost, user.Balance)
	}

	// Deduct balance
	user.Balance -= totalCost
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update user balance: %w", err)
	}

	// Add to portfolio
	return s.updatePortfolioForBuy(userID, symbol, quantity)
}

func (s *AdvancedOrderService) processSellTransaction(userID int, symbol string, quantity int, proceeds float64) error {
	// Check if user has enough shares
	portfolio, err := s.portfolioRepo.GetByUserIDAndSymbol(userID, symbol)
	if err != nil {
		return fmt.Errorf("stock not found in portfolio: %s", symbol)
	}

	if portfolio.Quantity < quantity {
		return fmt.Errorf("insufficient shares: required %d, available %d", quantity, portfolio.Quantity)
	}

	// Add proceeds to user balance
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	user.Balance += proceeds
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update user balance: %w", err)
	}

	// Remove from portfolio
	return s.updatePortfolioForSell(userID, symbol, quantity)
}

func (s *AdvancedOrderService) processShortTransaction(userID int, symbol string, quantity int, proceeds float64) error {
	// Add proceeds to user balance (minus margin requirement)
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	user.Balance += proceeds
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update user balance: %w", err)
	}

	// Create negative position in portfolio
	return s.updatePortfolioForShort(userID, symbol, quantity)
}

func (s *AdvancedOrderService) processCoverTransaction(userID int, symbol string, quantity int, totalCost float64) error {
	// Check if user has short position
	portfolio, err := s.portfolioRepo.GetByUserIDAndSymbol(userID, symbol)
	if err != nil {
		return fmt.Errorf("short position not found for: %s", symbol)
	}

	if portfolio.Quantity >= 0 || (-portfolio.Quantity) < quantity {
		return fmt.Errorf("insufficient short position: required %d, available %d", quantity, -portfolio.Quantity)
	}

	// Deduct cost from user balance
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if user.Balance < totalCost {
		return fmt.Errorf("insufficient balance for cover: required %.2f, available %.2f", totalCost, user.Balance)
	}

	user.Balance -= totalCost
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update user balance: %w", err)
	}

	// Reduce short position
	return s.updatePortfolioForCover(userID, symbol, quantity)
}

func (s *AdvancedOrderService) updatePortfolioForBuy(userID int, symbol string, quantity int) error {
	portfolio, err := s.portfolioRepo.GetByUserIDAndSymbol(userID, symbol)
	if err != nil {
		// Create new portfolio entry if doesn't exist
		stock, err := s.stockRepo.GetBySymbol(symbol)
		if err != nil {
			return err
		}

		newPortfolio := &domain.Portfolio{
			UserID:       userID,
			StockSymbol:  symbol,
			Quantity:     quantity,
			AveragePrice: stock.CurrentPrice,
			TotalCost:    float64(quantity) * stock.CurrentPrice,
		}
		return s.portfolioRepo.Create(newPortfolio)
	}

	// Update existing portfolio
	stock, err := s.stockRepo.GetBySymbol(symbol)
	if err != nil {
		return err
	}

	newQuantity := portfolio.Quantity + quantity
	newCost := portfolio.TotalCost + (float64(quantity) * stock.CurrentPrice)
	portfolio.Quantity = newQuantity
	portfolio.TotalCost = newCost
	portfolio.AveragePrice = newCost / float64(newQuantity)

	return s.portfolioRepo.Update(portfolio)
}

func (s *AdvancedOrderService) updatePortfolioForSell(userID int, symbol string, quantity int) error {
	portfolio, err := s.portfolioRepo.GetByUserIDAndSymbol(userID, symbol)
	if err != nil {
		return err
	}

	portfolio.Quantity -= quantity
	
	// If all shares sold, delete portfolio entry
	if portfolio.Quantity <= 0 {
		return s.portfolioRepo.Delete(userID, symbol)
	}

	// Update total cost proportionally
	sellRatio := float64(quantity) / float64(portfolio.Quantity + quantity)
	portfolio.TotalCost *= (1 - sellRatio)

	return s.portfolioRepo.Update(portfolio)
}

func (s *AdvancedOrderService) updatePortfolioForShort(userID int, symbol string, quantity int) error {
	portfolio, err := s.portfolioRepo.GetByUserIDAndSymbol(userID, symbol)
	if err != nil {
		// Create new short position
		stock, err := s.stockRepo.GetBySymbol(symbol)
		if err != nil {
			return err
		}

		newPortfolio := &domain.Portfolio{
			UserID:       userID,
			StockSymbol:  symbol,
			Quantity:     -quantity, // Negative for short position
			AveragePrice: stock.CurrentPrice,
			TotalCost:    -float64(quantity) * stock.CurrentPrice,
		}
		return s.portfolioRepo.Create(newPortfolio)
	}

	// Add to existing short position
	stock, err := s.stockRepo.GetBySymbol(symbol)
	if err != nil {
		return err
	}

	portfolio.Quantity -= quantity // More negative
	portfolio.TotalCost -= float64(quantity) * stock.CurrentPrice

	return s.portfolioRepo.Update(portfolio)
}

func (s *AdvancedOrderService) updatePortfolioForCover(userID int, symbol string, quantity int) error {
	portfolio, err := s.portfolioRepo.GetByUserIDAndSymbol(userID, symbol)
	if err != nil {
		return err
	}

	portfolio.Quantity += quantity // Less negative

	// If fully covered, delete portfolio entry
	if portfolio.Quantity >= 0 {
		return s.portfolioRepo.Delete(userID, symbol)
	}

	// Update cost proportionally
	coverRatio := float64(quantity) / float64(-portfolio.Quantity + quantity)
	portfolio.TotalCost *= (1 - coverRatio)

	return s.portfolioRepo.Update(portfolio)
}

func (s *AdvancedOrderService) UpdateBalanceOnExecution(execution *domain.OrderExecution) error {
	return nil
} 