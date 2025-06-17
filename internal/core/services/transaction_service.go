package services

import (
	"fmt"
	"stock-simulation-backend/internal/core/domain"
	"stock-simulation-backend/internal/core/ports/repositories"
	"stock-simulation-backend/internal/core/ports/services"
	"time"
)

type transactionService struct {
	transactionRepo repositories.TransactionRepository
	portfolioRepo   repositories.PortfolioRepository
	stockRepo       repositories.StockRepository
	userRepo        repositories.UserRepository
}

func NewTransactionService(
	transactionRepo repositories.TransactionRepository,
	portfolioRepo repositories.PortfolioRepository,
	stockRepo repositories.StockRepository,
	userRepo repositories.UserRepository,
) services.TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		portfolioRepo:   portfolioRepo,
		stockRepo:       stockRepo,
		userRepo:        userRepo,
	}
}

func (s *transactionService) BuyStock(userID int, req *domain.TransactionRequest) (*domain.TransactionResponse, error) {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get stock
	stock, err := s.stockRepo.GetBySymbol(req.StockSymbol)
	if err != nil {
		return nil, fmt.Errorf("stock not found: %w", err)
	}

	// Calculate total amount
	totalAmount := float64(req.Quantity) * stock.CurrentPrice

	// Check if user has enough balance
	if user.Balance < totalAmount {
		return nil, fmt.Errorf("insufficient balance")
	}

	// Create transaction
	transaction := &domain.Transaction{
		UserID:      userID,
		StockSymbol: req.StockSymbol,
		Type:        domain.TransactionTypeBuy,
		Quantity:    req.Quantity,
		Price:       stock.CurrentPrice,
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
	}

	err = s.transactionRepo.Create(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update user balance
	newBalance := user.Balance - totalAmount
	err = s.userRepo.UpdateBalance(userID, newBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}

	// Update portfolio
	err = s.updatePortfolioAfterBuy(userID, req.StockSymbol, req.Quantity, stock.CurrentPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to update portfolio: %w", err)
	}

	response := &domain.TransactionResponse{
		Transaction: transaction,
		Message:     fmt.Sprintf("Successfully bought %d shares of %s", req.Quantity, req.StockSymbol),
		Balance:     newBalance,
	}

	return response, nil
}

func (s *transactionService) SellStock(userID int, req *domain.TransactionRequest) (*domain.TransactionResponse, error) {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get stock
	stock, err := s.stockRepo.GetBySymbol(req.StockSymbol)
	if err != nil {
		return nil, fmt.Errorf("stock not found: %w", err)
	}

	// Get portfolio item
	portfolioItem, err := s.portfolioRepo.GetByUserIDAndSymbol(userID, req.StockSymbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}
	if portfolioItem == nil || portfolioItem.Quantity < req.Quantity {
		return nil, fmt.Errorf("insufficient shares to sell")
	}

	// Calculate total amount
	totalAmount := float64(req.Quantity) * stock.CurrentPrice

	// Create transaction
	transaction := &domain.Transaction{
		UserID:      userID,
		StockSymbol: req.StockSymbol,
		Type:        domain.TransactionTypeSell,
		Quantity:    req.Quantity,
		Price:       stock.CurrentPrice,
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
	}

	err = s.transactionRepo.Create(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update user balance
	newBalance := user.Balance + totalAmount
	err = s.userRepo.UpdateBalance(userID, newBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}

	// Update portfolio
	err = s.updatePortfolioAfterSell(userID, req.StockSymbol, req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to update portfolio: %w", err)
	}

	// Calculate and update total profit
	profit := totalAmount - (float64(req.Quantity) * portfolioItem.AveragePrice)
	newTotalProfit := user.TotalProfit + profit
	err = s.userRepo.UpdateTotalProfit(userID, newTotalProfit)
	if err != nil {
		return nil, fmt.Errorf("failed to update total profit: %w", err)
	}

	response := &domain.TransactionResponse{
		Transaction: transaction,
		Message:     fmt.Sprintf("Successfully sold %d shares of %s", req.Quantity, req.StockSymbol),
		Balance:     newBalance,
	}

	return response, nil
}

func (s *transactionService) GetUserTransactions(userID int, limit, offset int) ([]domain.Transaction, error) {
	transactions, err := s.transactionRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user transactions: %w", err)
	}

	return transactions, nil
}

func (s *transactionService) GetTransactionHistory(userID int, stockSymbol, transactionType string, limit int) ([]domain.Transaction, error) {
	transactions, err := s.transactionRepo.GetUserTransactionHistory(userID, stockSymbol, transactionType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction history: %w", err)
	}

	return transactions, nil
}

func (s *transactionService) GetTransactionByID(transactionID int) (*domain.Transaction, error) {
	transaction, err := s.transactionRepo.GetByID(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

func (s *transactionService) updatePortfolioAfterBuy(userID int, stockSymbol string, quantity int, price float64) error {
	portfolioItem, err := s.portfolioRepo.GetByUserIDAndSymbol(userID, stockSymbol)
	if err != nil {
		return err
	}

	if portfolioItem == nil {
		// Create new portfolio item
		newPortfolio := &domain.Portfolio{
			UserID:       userID,
			StockSymbol:  stockSymbol,
			Quantity:     quantity,
			AveragePrice: price,
			TotalCost:    float64(quantity) * price,
			UpdatedAt:    time.Now(),
		}
		return s.portfolioRepo.Create(newPortfolio)
	} else {
		// Update existing portfolio item
		totalCost := portfolioItem.TotalCost + (float64(quantity) * price)
		totalQuantity := portfolioItem.Quantity + quantity
		newAveragePrice := totalCost / float64(totalQuantity)

		portfolioItem.Quantity = totalQuantity
		portfolioItem.AveragePrice = newAveragePrice
		portfolioItem.TotalCost = totalCost
		portfolioItem.UpdatedAt = time.Now()

		return s.portfolioRepo.Update(portfolioItem)
	}
}

func (s *transactionService) updatePortfolioAfterSell(userID int, stockSymbol string, quantity int) error {
	portfolioItem, err := s.portfolioRepo.GetByUserIDAndSymbol(userID, stockSymbol)
	if err != nil {
		return err
	}

	if portfolioItem == nil {
		return fmt.Errorf("portfolio item not found")
	}

	newQuantity := portfolioItem.Quantity - quantity
	if newQuantity <= 0 {
		// Delete portfolio item if no shares left
		return s.portfolioRepo.Delete(userID, stockSymbol)
	} else {
		// Update portfolio item
		portfolioItem.Quantity = newQuantity
		portfolioItem.TotalCost = portfolioItem.AveragePrice * float64(newQuantity)
		portfolioItem.UpdatedAt = time.Now()

		return s.portfolioRepo.Update(portfolioItem)
	}
}