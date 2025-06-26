package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Stock struct {
	Symbol       string  `json:"symbol"`
	CurrentPrice float64 `json:"current_price"`
}

func main() {
	// Database connection
	db, err := sql.Open("mysql", "stockuser:stockpassword@tcp(localhost:3307)/stock_simulation")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	fmt.Println("🚀 Starting Stock Price Simulator...")
	fmt.Println("📈 This will randomly update stock prices every 3-10 seconds")
	fmt.Println("⏹️  Press Ctrl+C to stop")

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Run simulation
	for {
		updateStockPrices(db)
		
		// Random delay between 3-10 seconds
		delay := 3 + rand.Intn(8)
		fmt.Printf("⏰ Next price update in %d seconds...\n\n", delay)
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

func updateStockPrices(db *sql.DB) {
	// Get all stocks
	query := "SELECT symbol, current_price FROM stocks"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error fetching stocks: %v", err)
		return
	}
	defer rows.Close()

	var stocks []Stock
	for rows.Next() {
		var stock Stock
		err := rows.Scan(&stock.Symbol, &stock.CurrentPrice)
		if err != nil {
			log.Printf("Error scanning stock: %v", err)
			continue
		}
		stocks = append(stocks, stock)
	}

	if len(stocks) == 0 {
		fmt.Println("⚠️  No stocks found in database")
		return
	}

	fmt.Printf("📊 Updating prices for %d stocks at %s\n", len(stocks), time.Now().Format("15:04:05"))

	// Update each stock price
	for _, stock := range stocks {
		newPrice := generateNewPrice(stock.CurrentPrice)
		
		// Update in database
		updateQuery := "UPDATE stocks SET current_price = ?, updated_at = NOW() WHERE symbol = ?"
		_, err := db.Exec(updateQuery, newPrice, stock.Symbol)
		if err != nil {
			log.Printf("Error updating %s: %v", stock.Symbol, err)
			continue
		}

		// Calculate change
		change := newPrice - stock.CurrentPrice
		changePercent := (change / stock.CurrentPrice) * 100
		
		// Display update with colors
		var indicator string
		if change > 0 {
			indicator = "📈"
		} else if change < 0 {
			indicator = "📉"
		} else {
			indicator = "➡️"
		}

		fmt.Printf("%s %s: $%.2f → $%.2f (%.2f%%) [%+.2f]\n",
			indicator, stock.Symbol, stock.CurrentPrice, newPrice, changePercent, change)
	}
}

func generateNewPrice(currentPrice float64) float64 {
	// Generate realistic price movement (-3% to +3%)
	maxChangePercent := 3.0
	changePercent := (rand.Float64() - 0.5) * 2 * maxChangePercent // -3% to +3%
	
	// Apply some volatility patterns
	volatilityFactor := 1.0
	if rand.Float64() < 0.1 { // 10% chance of high volatility
		volatilityFactor = 2.0
	} else if rand.Float64() < 0.05 { // 5% chance of very high volatility
		volatilityFactor = 3.0
	}
	
	changePercent *= volatilityFactor
	
	// Calculate new price
	change := currentPrice * (changePercent / 100)
	newPrice := currentPrice + change
	
	// Ensure price doesn't go below $0.01
	if newPrice < 0.01 {
		newPrice = 0.01
	}
	
	// Round to 2 decimal places
	return math.Round(newPrice*100) / 100
} 