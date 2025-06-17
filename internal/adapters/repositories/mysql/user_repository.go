package mysql

import (
	"database/sql"
	"fmt"
	"stock-simulation-backend/internal/core/domain"
	"stock-simulation-backend/internal/core/ports/repositories"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, balance, total_profit, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`
	result, err := r.db.Exec(query, user.Username, user.Email, user.PasswordHash, user.Balance, user.TotalProfit)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get user ID: %w", err)
	}

	user.ID = int(id)
	return nil
}

func (r *userRepository) GetByID(id int) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, balance, total_profit, created_at, updated_at
		FROM users WHERE id = ?
	`
	var user domain.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Balance, &user.TotalProfit, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, balance, total_profit, created_at, updated_at
		FROM users WHERE email = ?
	`
	var user domain.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Balance, &user.TotalProfit, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, balance, total_profit, created_at, updated_at
		FROM users WHERE username = ?
	`
	var user domain.User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Balance, &user.TotalProfit, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) UpdateBalance(userID int, newBalance float64) error {
	query := `UPDATE users SET balance = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.Exec(query, newBalance, userID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}
	return nil
}

func (r *userRepository) UpdateTotalProfit(userID int, totalProfit float64) error {
	query := `UPDATE users SET total_profit = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.Exec(query, totalProfit, userID)
	if err != nil {
		return fmt.Errorf("failed to update total profit: %w", err)
	}
	return nil
}

func (r *userRepository) GetLeaderboard(limit int) ([]domain.UserProfile, error) {
	query := `
		SELECT id, username, email, balance, total_profit,
		       ROW_NUMBER() OVER (ORDER BY total_profit DESC) as rank
		FROM users
		ORDER BY total_profit DESC
		LIMIT ?
	`
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	defer rows.Close()

	var profiles []domain.UserProfile
	for rows.Next() {
		var profile domain.UserProfile
		err := rows.Scan(&profile.ID, &profile.Username, &profile.Email,
			&profile.Balance, &profile.TotalProfit, &profile.Rank)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user profile: %w", err)
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (r *userRepository) EmailExists(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`
	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}

func (r *userRepository) UsernameExists(username string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE username = ?`
	var count int
	err := r.db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}
	return count > 0, nil
}

func (r *userRepository) Update(user *domain.User) error {
	query := `
		UPDATE users 
		SET username = ?, email = ?, balance = ?, total_profit = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.Exec(query, user.Username, user.Email, user.Balance, user.TotalProfit, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}