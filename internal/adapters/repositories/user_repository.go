package repositories

import (
    "stock-simulation-backend/internal/core/domain"
)

type UserRepository interface {
    Create(user *domain.User) error
    GetByID(id int) (*domain.User, error)
    GetByEmail(email string) (*domain.User, error)
    GetByUsername(username string) (*domain.User, error)
    UpdateBalance(userID int, newBalance float64) error
    UpdateTotalProfit(userID int, totalProfit float64) error
    GetLeaderboard(limit int) ([]domain.UserProfile, error)
    EmailExists(email string) (bool, error)
    UsernameExists(username string) (bool, error)
}