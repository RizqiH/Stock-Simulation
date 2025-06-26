package mysql

import (
	"database/sql"
	"stock-simulation-backend/internal/core/domain"
	"time"
)

type historicalPriceRepository struct {
	db *sql.DB
}

func NewHistoricalPriceRepository(db *sql.DB) *historicalPriceRepository {
	return &historicalPriceRepository{db: db}
}

func (r *historicalPriceRepository) Create(price *domain.HistoricalPrice) error {
	query := `
		INSERT INTO historical_prices (symbol, date, open, high, low, close, volume, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.Exec(query, price.Symbol, price.Date, price.Open, price.High, 
		price.Low, price.Close, price.Volume, time.Now())
	
	return err
}

func (r *historicalPriceRepository) GetBySymbolAndDateRange(symbol string, startDate, endDate time.Time) ([]domain.HistoricalPrice, error) {
	query := `
		SELECT id, symbol, date, open, high, low, close, volume, created_at 
		FROM historical_prices 
		WHERE symbol = ? AND date BETWEEN ? AND ? 
		ORDER BY date ASC
	`
	
	rows, err := r.db.Query(query, symbol, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var prices []domain.HistoricalPrice
	for rows.Next() {
		var price domain.HistoricalPrice
		err := rows.Scan(&price.ID, &price.Symbol, &price.Date, &price.Open, 
			&price.High, &price.Low, &price.Close, &price.Volume, &price.CreatedAt)
		if err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}
	
	return prices, nil
}

func (r *historicalPriceRepository) GetBySymbolWithLimit(symbol string, limit int) ([]domain.HistoricalPrice, error) {
	query := `
		SELECT id, symbol, date, open, high, low, close, volume, created_at 
		FROM historical_prices 
		WHERE symbol = ? 
		ORDER BY date DESC 
		LIMIT ?
	`
	
	rows, err := r.db.Query(query, symbol, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var prices []domain.HistoricalPrice
	for rows.Next() {
		var price domain.HistoricalPrice
		err := rows.Scan(&price.ID, &price.Symbol, &price.Date, &price.Open, 
			&price.High, &price.Low, &price.Close, &price.Volume, &price.CreatedAt)
		if err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}
	
	// Reverse to get chronological order
	for i, j := 0, len(prices)-1; i < j; i, j = i+1, j-1 {
		prices[i], prices[j] = prices[j], prices[i]
	}
	
	return prices, nil
}

func (r *historicalPriceRepository) GetChartData(symbol string, period string) (*domain.ChartData, error) {
	var limit int
	switch period {
	case "1D":
		limit = 1
	case "7D":
		limit = 7
	case "30D":
		limit = 30
	case "90D":
		limit = 90
	case "1Y":
		limit = 365
	default:
		limit = 30
	}
	
	prices, err := r.GetBySymbolWithLimit(symbol, limit)
	if err != nil {
		return nil, err
	}
	
	if len(prices) == 0 {
		return &domain.ChartData{
			Symbol:     symbol,
			Period:     period,
			Prices:     []domain.HistoricalPrice{},
			Indicators: domain.ChartIndicators{},
		}, nil
	}
	
	// Calculate technical indicators
	indicators := r.calculateIndicators(prices)
	
	return &domain.ChartData{
		Symbol:     symbol,
		Period:     period,
		Prices:     prices,
		Indicators: indicators,
	}, nil
}

func (r *historicalPriceRepository) calculateIndicators(prices []domain.HistoricalPrice) domain.ChartIndicators {
	indicators := domain.ChartIndicators{
		MA20:   []float64{},
		MA50:   []float64{},
		RSI:    []float64{},
		Volume: []int64{},
	}
	
	if len(prices) == 0 {
		return indicators
	}
	
	// Extract closing prices and volumes
	closePrices := make([]float64, len(prices))
	volumes := make([]int64, len(prices))
	
	for i, price := range prices {
		closePrices[i] = price.Close
		volumes[i] = price.Volume
	}
	
	indicators.Volume = volumes
	
	// Calculate Moving Averages
	indicators.MA20 = r.calculateMovingAverage(closePrices, 20)
	indicators.MA50 = r.calculateMovingAverage(closePrices, 50)
	
	// Calculate RSI
	indicators.RSI = r.calculateRSI(closePrices, 14)
	
	return indicators
}

func (r *historicalPriceRepository) calculateMovingAverage(prices []float64, period int) []float64 {
	if len(prices) < period {
		return make([]float64, len(prices))
	}
	
	ma := make([]float64, len(prices))
	
	for i := period - 1; i < len(prices); i++ {
		sum := 0.0
		for j := i - period + 1; j <= i; j++ {
			sum += prices[j]
		}
		ma[i] = sum / float64(period)
	}
	
	return ma
}

func (r *historicalPriceRepository) calculateRSI(prices []float64, period int) []float64 {
	if len(prices) < period+1 {
		return make([]float64, len(prices))
	}
	
	rsi := make([]float64, len(prices))
	gains := make([]float64, len(prices)-1)
	losses := make([]float64, len(prices)-1)
	
	// Calculate gains and losses
	for i := 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains[i-1] = change
		} else {
			losses[i-1] = -change
		}
	}
	
	// Calculate RSI
	for i := period; i < len(prices); i++ {
		avgGain := 0.0
		avgLoss := 0.0
		
		for j := i - period; j < i; j++ {
			avgGain += gains[j]
			avgLoss += losses[j]
		}
		
		avgGain /= float64(period)
		avgLoss /= float64(period)
		
		if avgLoss == 0 {
			rsi[i] = 100
		} else {
			rs := avgGain / avgLoss
			rsi[i] = 100 - (100 / (1 + rs))
		}
	}
	
	return rsi
}

func (r *historicalPriceRepository) Update(price *domain.HistoricalPrice) error {
	query := `
		UPDATE historical_prices 
		SET open = ?, high = ?, low = ?, close = ?, volume = ? 
		WHERE id = ?
	`
	
	_, err := r.db.Exec(query, price.Open, price.High, price.Low, 
		price.Close, price.Volume, price.ID)
	
	return err
}

func (r *historicalPriceRepository) DeleteOlderThan(date time.Time) error {
	query := `DELETE FROM historical_prices WHERE date < ?`
	_, err := r.db.Exec(query, date)
	return err
}

func (r *historicalPriceRepository) BatchInsert(prices []domain.HistoricalPrice) error {
	if len(prices) == 0 {
		return nil
	}
	
	query := `
		INSERT INTO historical_prices (symbol, date, open, high, low, close, volume, created_at) 
		VALUES `
	
	values := []interface{}{}
	for i, price := range prices {
		if i > 0 {
			query += ", "
		}
		query += "(?, ?, ?, ?, ?, ?, ?, ?)"
		values = append(values, price.Symbol, price.Date, price.Open, price.High,
			price.Low, price.Close, price.Volume, time.Now())
	}
	
	_, err := r.db.Exec(query, values...)
	return err
}

func (r *historicalPriceRepository) GetLatestPrice(symbol string) (*domain.HistoricalPrice, error) {
	query := `
		SELECT id, symbol, date, open, high, low, close, volume, created_at 
		FROM historical_prices 
		WHERE symbol = ? 
		ORDER BY date DESC 
		LIMIT 1
	`
	
	var price domain.HistoricalPrice
	err := r.db.QueryRow(query, symbol).Scan(
		&price.ID, &price.Symbol, &price.Date, &price.Open,
		&price.High, &price.Low, &price.Close, &price.Volume, &price.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &price, nil
}

func (r *historicalPriceRepository) GetAvailableSymbols() ([]string, error) {
	query := `SELECT DISTINCT symbol FROM historical_prices ORDER BY symbol`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var symbols []string
	for rows.Next() {
		var symbol string
		if err := rows.Scan(&symbol); err != nil {
			return nil, err
		}
		symbols = append(symbols, symbol)
	}
	
	return symbols, nil
} 