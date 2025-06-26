package domain

import (
    "time"
)

type Portfolio struct {
    ID           int       `json:"id" db:"id"`
    UserID       int       `json:"user_id" db:"user_id"`
    StockSymbol  string    `json:"stock_symbol" db:"stock_symbol"`
    Quantity     int       `json:"quantity" db:"quantity"`
    AveragePrice float64   `json:"average_price" db:"average_price"`
    TotalCost    float64   `json:"total_cost" db:"total_cost"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type PortfolioItem struct {
    StockSymbol   string  `json:"stock_symbol"`
    StockName     string  `json:"stock_name"`
    Quantity      int     `json:"quantity"`
    AveragePrice  float64 `json:"average_price"`
    CurrentPrice  float64 `json:"current_price"`
    TotalCost     float64 `json:"total_cost"`
    CurrentValue  float64 `json:"current_value"`
    ProfitLoss    float64 `json:"profit_loss"`
    ProfitLossPct float64 `json:"profit_loss_pct"`
}

type PortfolioSummary struct {
    TotalValue     float64         `json:"total_value"`
    TotalCost      float64         `json:"total_cost"`
    TotalProfit    float64         `json:"total_profit"`
    TotalProfitPct float64         `json:"total_profit_pct"`
    Holdings       []PortfolioItem `json:"holdings"`
}

type PortfolioPerformance struct {
    Period      string  `json:"period"`
    StartValue  float64 `json:"start_value"`
    EndValue    float64 `json:"end_value"`
    Profit      float64 `json:"profit"`
    ProfitPct   float64 `json:"profit_pct"`
    Transactions int    `json:"transactions"`
}

type PortfolioDataPoint struct {
    Date            time.Time `json:"date"`
    TotalValue      float64   `json:"total_value"`
    TotalCost       float64   `json:"total_cost"`
    ProfitLoss      float64   `json:"profit_loss"`
    ProfitLossPct   float64   `json:"profit_loss_pct"`
    CashValue       float64   `json:"cash_value"`
    InvestmentValue float64   `json:"investment_value"`
}