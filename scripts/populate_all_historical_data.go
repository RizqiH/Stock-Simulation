package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Database connection
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3307)/stock_simulation?parseTime=true")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("🔍 Checking symbols without historical data...")
	
	// Get all stock symbols
	symbols, err := getAllStockSymbols(db)
	if err != nil {
		log.Fatal("Failed to get stock symbols:", err)
	}
	
	// Check which ones need historical data
	missingSymbols := []string{}
	for _, symbol := range symbols {
		count, err := getHistoricalDataCount(db, symbol)
		if err != nil {
			log.Printf("Warning: Could not check data for %s: %v", symbol, err)
			continue
		}
		
		if count == 0 {
			missingSymbols = append(missingSymbols, symbol)
			fmt.Printf("❌ %s: No historical data (%d records)\n", symbol, count)
		} else {
			fmt.Printf("✅ %s: Has historical data (%d records)\n", symbol, count)
		}
	}
	
	if len(missingSymbols) == 0 {
		fmt.Println("🎉 All symbols already have historical data!")
		return
	}
	
	fmt.Printf("\n📊 Adding historical data for %d symbols: %v\n", len(missingSymbols), missingSymbols)
	
	// Add historical data for missing symbols
	for _, symbol := range missingSymbols {
		err := addHistoricalDataForSymbol(db, symbol)
		if err != nil {
			log.Printf("Failed to add data for %s: %v", symbol, err)
		} else {
			fmt.Printf("✅ Added 30 days of historical data for %s\n", symbol)
		}
	}
	
	fmt.Println("\n🎉 Historical data population completed!")
	
	// Final summary
	fmt.Println("\n📋 Final Summary:")
	for _, symbol := range symbols {
		count, _ := getHistoricalDataCount(db, symbol)
		fmt.Printf("  %s: %d records\n", symbol, count)
	}
}

func getAllStockSymbols(db *sql.DB) ([]string, error) {
	query := `SELECT DISTINCT symbol FROM stocks ORDER BY symbol`
	rows, err := db.Query(query)
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

func getHistoricalDataCount(db *sql.DB, symbol string) (int, error) {
	query := `SELECT COUNT(*) FROM historical_prices WHERE symbol = ?`
	var count int
	err := db.QueryRow(query, symbol).Scan(&count)
	return count, err
}

func addHistoricalDataForSymbol(db *sql.DB, symbol string) error {
	// Get current price from stocks table
	var currentPrice float64
	err := db.QueryRow("SELECT current_price FROM stocks WHERE symbol = ?", symbol).Scan(&currentPrice)
	if err != nil {
		// Default fallback prices for each symbol
		priceMap := map[string]float64{
			"AAPL":  150.00,
			"GOOGL": 2800.00,
			"MSFT":  300.00,
			"AMZN":  3200.00,
			"TSLA":  800.00,
			"META":  320.00,
			"NVDA":  450.00,
			"NFLX":  400.00,
			"BABA":  90.00,
			"V":     220.00,
			"JPM":   165.00,
			"JNJ":   175.00,
			"WMT":   145.00,
			"PG":    155.00,
			"UNH":   520.00,
		}
		
		if price, exists := priceMap[symbol]; exists {
			currentPrice = price
		} else {
			currentPrice = 100.00 // Default fallback
		}
	}
	
	// Generate 30 days of historical data
	query := `INSERT IGNORE INTO historical_prices (symbol, date, open, high, low, close, volume) VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	baseVolume := getBaseVolumeForSymbol(symbol)
	
	for i := 30; i >= 1; i-- {
		date := time.Now().AddDate(0, 0, -i)
		
		// Calculate prices with some realistic volatility
		volatility := 0.02 // 2% daily volatility
		change := (rand.Float64() - 0.5) * 2 * volatility
		
		basePrice := currentPrice * (1 + change*float64(i)*0.1) // Slight backward trend
		
		open := basePrice * (1 + (rand.Float64()-0.5)*0.01)
		high := basePrice * (1 + rand.Float64()*0.02)
		low := basePrice * (1 - rand.Float64()*0.02)
		close := basePrice * (1 + (rand.Float64()-0.5)*0.015)
		
		// Ensure high >= low and prices are reasonable
		if high < low {
			high, low = low, high
		}
		if close > high {
			high = close
		}
		if close < low {
			low = close
		}
		if open > high {
			high = open
		}
		if open < low {
			low = open
		}
		
		volume := baseVolume + int(rand.Float64()*float64(baseVolume)*0.5)
		
		_, err := db.Exec(query, symbol, date.Format("2006-01-02"), open, high, low, close, volume)
		if err != nil {
			return fmt.Errorf("failed to insert data for %s on %s: %v", symbol, date.Format("2006-01-02"), err)
		}
	}
	
	return nil
}

func getBaseVolumeForSymbol(symbol string) int {
	volumeMap := map[string]int{
		"AAPL":  50000000,
		"GOOGL": 25000000,
		"MSFT":  30000000,
		"AMZN":  20000000,
		"TSLA":  40000000,
		"META":  35000000,
		"NVDA":  45000000,
		"NFLX":  15000000,
		"BABA":  25000000,
		"V":     10000000,
		"JPM":   15000000,
		"JNJ":   8500000,
		"WMT":   9200000,
		"PG":    7800000,
		"UNH":   6500000,
	}
	
	if volume, exists := volumeMap[symbol]; exists {
		return volume
	}
	return 10000000 // Default volume
} 