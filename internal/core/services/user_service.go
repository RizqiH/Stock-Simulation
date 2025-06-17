package services

import (
	"fmt"
	"stock-simulation-backend/internal/adapters/middleware"
	"stock-simulation-backend/internal/core/domain"
	"stock-simulation-backend/internal/core/ports/repositories"
	"stock-simulation-backend/internal/core/ports/services"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) services.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Register(req *domain.UserRegistration) (*domain.UserProfile, error) {
	// Check if email already exists
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	// Check if username already exists
	exists, err = s.userRepo.UsernameExists(req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Balance:      100000.0, // Starting balance
		TotalProfit:  0.0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Return user profile
	profile := &domain.UserProfile{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Balance:     user.Balance,
		TotalProfit: user.TotalProfit,
	}

	return profile, nil
}

func (s *userService) Login(req *domain.UserLogin) (string, *domain.UserProfile, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return "", nil, fmt.Errorf("invalid email or password")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return token and user profile
	profile := &domain.UserProfile{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Balance:     user.Balance,
		TotalProfit: user.TotalProfit,
	}

	return token, profile, nil
}

func (s *userService) GetProfile(userID int) (*domain.UserProfile, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	profile := &domain.UserProfile{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Balance:     user.Balance,
		TotalProfit: user.TotalProfit,
	}

	return profile, nil
}

func (s *userService) UpdateProfile(profile *domain.UserProfile) error {
	// Get current user
	user, err := s.userRepo.GetByID(profile.ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields
	user.Username = profile.Username
	user.Email = profile.Email
	user.UpdatedAt = time.Now()

	// Save changes
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *userService) GetLeaderboard(limit int) ([]domain.UserProfile, error) {
	leaderboard, err := s.userRepo.GetLeaderboard(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}

	return leaderboard, nil
}

func (s *userService) UpdateBalance(userID int, newBalance float64) error {
	err := s.userRepo.UpdateBalance(userID, newBalance)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}

func (s *userService) GetByID(userID int) (*domain.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}