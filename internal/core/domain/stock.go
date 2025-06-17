package domain

import (
    "time"
)

type Stock struct {
    ID          int       `json:"id" db:"id"`
    Symbol      string    `json:"symbol" db:"symbol"`
    Name        string    `json:"name" db:"name"`
    CurrentPrice float64  `json:"current_price" db:"current_price"`
    OpenPrice   float64   `json:"open_price" db:"open_price"`
    HighPrice   float64   `json:"high_price" db:"high_price"`
    LowPrice    float64   `json:"low_price" db:"low_price"`
    Volume      int64     `json:"volume" db:"volume"`
    MarketCap   float64   `json:"market_cap" db:"market_cap"`
    Sector      string    `json:"sector" db:"sector"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type StockPrice struct {
    Symbol    string    `json:"symbol"`
    Price     float64   `json:"price"`
    Change    float64   `json:"change"`
    ChangePct float64   `json:"change_pct"`
    UpdatedAt time.Time `json:"updated_at"`
}