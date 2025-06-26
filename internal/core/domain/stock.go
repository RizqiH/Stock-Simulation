package domain

import (
    "time"
)

type Stock struct {
    ID           int       `json:"id" db:"id"`
    Symbol       string    `json:"symbol" db:"symbol"`
    Name         string    `json:"name" db:"name"`
    CurrentPrice float64   `json:"current_price" db:"current_price"`
    PreviousClose float64  `json:"previous_close" db:"previous_close"`
    Volume       int64     `json:"volume" db:"volume"`
    MarketCap    int64     `json:"market_cap" db:"market_cap"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type StockPrice struct {
    Symbol    string    `json:"symbol"`
    Price     float64   `json:"price"`
    Change    float64   `json:"change"`
    ChangePct float64   `json:"change_pct"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Historical price data for charting
type HistoricalPrice struct {
    ID        int       `json:"id" db:"id"`
    Symbol    string    `json:"symbol" db:"symbol"`
    Date      time.Time `json:"date" db:"date"`
    Open      float64   `json:"open" db:"open"`
    High      float64   `json:"high" db:"high"`
    Low       float64   `json:"low" db:"low"`
    Close     float64   `json:"close" db:"close"`
    Volume    int64     `json:"volume" db:"volume"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Chart data structure for frontend
type ChartData struct {
    Symbol     string            `json:"symbol"`
    Period     string            `json:"period"`
    Prices     []HistoricalPrice `json:"prices"`
    Indicators ChartIndicators   `json:"indicators"`
}

// Technical indicators
type ChartIndicators struct {
    MA20   []float64 `json:"ma20"`   // Moving Average 20
    MA50   []float64 `json:"ma50"`   // Moving Average 50
    RSI    []float64 `json:"rsi"`    // Relative Strength Index
    Volume []int64   `json:"volume"` // Volume data
}

// Note: Order types, status, and structures are now defined in order.go
// This file focuses on Stock-specific domain models