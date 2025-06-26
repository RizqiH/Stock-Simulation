package main

import (
	"database/sql"
	"fmt"
	"log"
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

	// Create tables first
	err = createTables(db)
	if err != nil {
		log.Printf("Warning: Could not create tables (may already exist): %v", err)
	}

	// Insert sample historical data
	err = insertHistoricalData(db)
	if err != nil {
		log.Fatal("Failed to insert historical data:", err)
	}

	fmt.Println("✅ Historical data inserted successfully!")
}

func createTables(db *sql.DB) error {
	// Create historical_prices table
	historicalPricesSQL := `
	CREATE TABLE IF NOT EXISTS historical_prices (
		id INT AUTO_INCREMENT PRIMARY KEY,
		symbol VARCHAR(10) NOT NULL,
		date DATE NOT NULL,
		open DECIMAL(10,2) NOT NULL,
		high DECIMAL(10,2) NOT NULL,
		low DECIMAL(10,2) NOT NULL,
		close DECIMAL(10,2) NOT NULL,
		volume BIGINT NOT NULL DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_symbol_date (symbol, date),
		INDEX idx_date (date),
		UNIQUE KEY unique_symbol_date (symbol, date)
	);`

	_, err := db.Exec(historicalPricesSQL)
	if err != nil {
		return err
	}

	// Create orders table
	ordersSQL := `
	CREATE TABLE IF NOT EXISTS orders (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		stock_symbol VARCHAR(10) NOT NULL,
		order_type ENUM('MARKET', 'LIMIT', 'STOP_LOSS', 'TAKE_PROFIT') NOT NULL,
		side ENUM('BUY', 'SELL') NOT NULL,
		quantity INT NOT NULL,
		price DECIMAL(10,2) DEFAULT NULL,
		stop_price DECIMAL(10,2) DEFAULT NULL,
		status ENUM('PENDING', 'EXECUTED', 'CANCELLED', 'EXPIRED') DEFAULT 'PENDING',
		executed_price DECIMAL(10,2) DEFAULT NULL,
		executed_at TIMESTAMP NULL,
		expires_at TIMESTAMP NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_user_id (user_id),
		INDEX idx_symbol (stock_symbol),
		INDEX idx_status (status),
		INDEX idx_created_at (created_at)
	);`

	_, err = db.Exec(ordersSQL)
	return err
}

func insertHistoricalData(db *sql.DB) error {
	// Sample data for AAPL (last 30 days)
	aaplData := [][]interface{}{
		{"AAPL", time.Now().AddDate(0, 0, -30), 148.50, 152.30, 147.80, 151.20, 82345600},
		{"AAPL", time.Now().AddDate(0, 0, -29), 151.20, 153.40, 149.90, 150.80, 78234500},
		{"AAPL", time.Now().AddDate(0, 0, -28), 150.80, 154.20, 149.30, 153.10, 89765400},
		{"AAPL", time.Now().AddDate(0, 0, -27), 153.10, 155.60, 151.40, 154.90, 92187300},
		{"AAPL", time.Now().AddDate(0, 0, -26), 154.90, 157.30, 153.20, 156.70, 87654200},
		{"AAPL", time.Now().AddDate(0, 0, -25), 156.70, 158.90, 155.10, 157.40, 95432100},
		{"AAPL", time.Now().AddDate(0, 0, -24), 157.40, 159.80, 156.30, 158.20, 88976500},
		{"AAPL", time.Now().AddDate(0, 0, -23), 158.20, 160.50, 157.10, 159.60, 91234700},
		{"AAPL", time.Now().AddDate(0, 0, -22), 159.60, 161.30, 158.40, 160.90, 86543200},
		{"AAPL", time.Now().AddDate(0, 0, -21), 160.90, 162.70, 159.80, 161.50, 94876300},
		{"AAPL", time.Now().AddDate(0, 0, -20), 161.50, 163.20, 160.20, 162.80, 87654300},
		{"AAPL", time.Now().AddDate(0, 0, -19), 162.80, 164.50, 161.90, 163.40, 92187400},
		{"AAPL", time.Now().AddDate(0, 0, -18), 163.40, 165.10, 162.30, 164.70, 89543600},
		{"AAPL", time.Now().AddDate(0, 0, -17), 164.70, 166.20, 163.80, 165.30, 91876500},
		{"AAPL", time.Now().AddDate(0, 0, -16), 165.30, 167.40, 164.50, 166.80, 88765400},
		{"AAPL", time.Now().AddDate(0, 0, -15), 166.80, 168.90, 165.60, 167.50, 94321700},
		{"AAPL", time.Now().AddDate(0, 0, -14), 167.50, 169.30, 166.20, 168.40, 87654800},
		{"AAPL", time.Now().AddDate(0, 0, -13), 168.40, 170.60, 167.10, 169.20, 92456300},
		{"AAPL", time.Now().AddDate(0, 0, -12), 169.20, 171.80, 168.30, 170.50, 89876500},
		{"AAPL", time.Now().AddDate(0, 0, -11), 170.50, 172.40, 169.60, 171.30, 91234800},
		{"AAPL", time.Now().AddDate(0, 0, -10), 171.30, 173.70, 170.40, 172.90, 88543600},
		{"AAPL", time.Now().AddDate(0, 0, -9), 172.90, 174.50, 171.80, 173.60, 94567200},
		{"AAPL", time.Now().AddDate(0, 0, -8), 173.60, 175.30, 172.40, 174.80, 87234500},
		{"AAPL", time.Now().AddDate(0, 0, -7), 174.80, 176.90, 173.90, 175.40, 92876300},
		{"AAPL", time.Now().AddDate(0, 0, -6), 175.40, 177.20, 174.50, 176.10, 89654700},
		{"AAPL", time.Now().AddDate(0, 0, -5), 176.10, 178.60, 175.30, 177.80, 91432800},
		{"AAPL", time.Now().AddDate(0, 0, -4), 177.80, 179.40, 176.90, 178.50, 88765400},
		{"AAPL", time.Now().AddDate(0, 0, -3), 178.50, 180.20, 177.60, 179.30, 94321600},
		{"AAPL", time.Now().AddDate(0, 0, -2), 179.30, 181.80, 178.40, 180.90, 87654200},
		{"AAPL", time.Now().AddDate(0, 0, -1), 180.90, 182.50, 179.80, 181.60, 92456700},
	}

	// Insert AAPL data
	query := `INSERT IGNORE INTO historical_prices (symbol, date, open, high, low, close, volume) VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	fmt.Println("Inserting AAPL historical data...")
	for _, data := range aaplData {
		_, err := db.Exec(query, data...)
		if err != nil {
			return fmt.Errorf("failed to insert AAPL data: %v", err)
		}
	}

	// Sample data for GOOGL
	googlData := [][]interface{}{
		{"GOOGL", time.Now().AddDate(0, 0, -30), 2420.50, 2458.30, 2398.80, 2441.20, 1234560},
		{"GOOGL", time.Now().AddDate(0, 0, -29), 2441.20, 2467.40, 2422.90, 2453.80, 1187450},
		{"GOOGL", time.Now().AddDate(0, 0, -28), 2453.80, 2489.20, 2431.30, 2476.10, 1298740},
		{"GOOGL", time.Now().AddDate(0, 0, -27), 2476.10, 2502.60, 2454.40, 2489.90, 1356730},
		{"GOOGL", time.Now().AddDate(0, 0, -26), 2489.90, 2518.30, 2467.20, 2506.70, 1421850},
		{"GOOGL", time.Now().AddDate(0, 0, -25), 2506.70, 2534.90, 2485.10, 2521.40, 1187640},
		{"GOOGL", time.Now().AddDate(0, 0, -24), 2521.40, 2548.80, 2499.30, 2537.20, 1345290},
		{"GOOGL", time.Now().AddDate(0, 0, -23), 2537.20, 2565.50, 2516.10, 2552.60, 1256730},
		{"GOOGL", time.Now().AddDate(0, 0, -22), 2552.60, 2578.30, 2531.40, 2567.90, 1423180},
		{"GOOGL", time.Now().AddDate(0, 0, -21), 2567.90, 2592.70, 2546.80, 2581.50, 1189540},
	}

	fmt.Println("Inserting GOOGL historical data...")
	for _, data := range googlData {
		_, err := db.Exec(query, data...)
		if err != nil {
			return fmt.Errorf("failed to insert GOOGL data: %v", err)
		}
	}

	// Sample data for MSFT
	msftData := [][]interface{}{
		{"MSFT", time.Now().AddDate(0, 0, -30), 248.50, 255.30, 246.80, 252.20, 28345600},
		{"MSFT", time.Now().AddDate(0, 0, -29), 252.20, 258.40, 249.90, 256.80, 31234500},
		{"MSFT", time.Now().AddDate(0, 0, -28), 256.80, 262.20, 254.30, 260.10, 29765400},
		{"MSFT", time.Now().AddDate(0, 0, -27), 260.10, 266.60, 257.40, 264.90, 32187300},
		{"MSFT", time.Now().AddDate(0, 0, -26), 264.90, 270.30, 262.20, 268.70, 27654200},
		{"MSFT", time.Now().AddDate(0, 0, -25), 268.70, 274.90, 266.10, 272.40, 35432100},
		{"MSFT", time.Now().AddDate(0, 0, -24), 272.40, 278.80, 270.30, 276.20, 28976500},
		{"MSFT", time.Now().AddDate(0, 0, -23), 276.20, 282.50, 274.10, 280.60, 31234700},
		{"MSFT", time.Now().AddDate(0, 0, -22), 280.60, 286.30, 278.40, 284.90, 26543200},
		{"MSFT", time.Now().AddDate(0, 0, -21), 284.90, 290.70, 282.80, 288.50, 34876300},
	}

	fmt.Println("Inserting MSFT historical data...")
	for _, data := range msftData {
		_, err := db.Exec(query, data...)
		if err != nil {
			return fmt.Errorf("failed to insert MSFT data: %v", err)
		}
	}

	return nil
} 