package domain

import (
    "time"
)

type User struct {
    ID           int       `json:"id" db:"id"`
    Username     string    `json:"username" db:"username"`
    Email        string    `json:"email" db:"email"`
    PasswordHash string    `json:"-" db:"password_hash"`
    Balance      float64   `json:"balance" db:"balance"`
    TotalProfit  float64   `json:"total_profit" db:"total_profit"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserRegistration struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type UserLogin struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type UserProfile struct {
    ID          int     `json:"id"`
    Username    string  `json:"username"`
    Email       string  `json:"email"`
    Balance     float64 `json:"balance"`
    TotalProfit float64 `json:"total_profit"`
    Rank        int     `json:"rank,omitempty"`
}