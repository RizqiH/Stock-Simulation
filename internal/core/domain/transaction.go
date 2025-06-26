package domain

import (
    "time"
)

type TransactionType string

const (
    TransactionTypeBuy  TransactionType = "BUY"
    TransactionTypeSell TransactionType = "SELL"
)

type Transaction struct {
    ID          int             `json:"id" db:"id"`
    UserID      int             `json:"user_id" db:"user_id"`
    StockSymbol string          `json:"stock_symbol" db:"stock_symbol"`
    Type        TransactionType `json:"type" db:"transaction_type"`
    Quantity    int             `json:"quantity" db:"quantity"`
    Price       float64         `json:"price" db:"price"`
    TotalAmount float64         `json:"total_amount" db:"total_amount"`
    CreatedAt   time.Time       `json:"created_at" db:"created_at"`
}

type TransactionRequest struct {
    StockSymbol string `json:"stock_symbol" binding:"required"`
    Quantity    int    `json:"quantity" binding:"required,min=1"`
}

type TransactionResponse struct {
    Transaction *Transaction `json:"transaction"`
    Message     string       `json:"message"`
    Balance     float64      `json:"new_balance"`
}