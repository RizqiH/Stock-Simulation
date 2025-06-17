package services

import "stock-simulation-backend/internal/core/domain"

type UserService interface {
	Register(req *domain.UserRegistration) (*domain.UserProfile, error)
	Login(req *domain.UserLogin) (string, *domain.UserProfile, error)
	GetProfile(userID int) (*domain.UserProfile, error)
	UpdateProfile(profile *domain.UserProfile) error
	GetLeaderboard(limit int) ([]domain.UserProfile, error)
	UpdateBalance(userID int, newBalance float64) error
	GetByID(userID int) (*domain.User, error)
}