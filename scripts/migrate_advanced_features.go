package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Database connection
	db, err := sql.Open("mysql", "stockuser:stockpassword@tcp(localhost:3307)/stock_simulation?parseTime=true")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("🚀 Starting Advanced Trading Features Migration...")

	// Run migrations
	err = runAdvancedFeaturesMigration(db)
	if err != nil {
		log.Fatal("Failed to run migration:", err)
	}

	fmt.Println("✅ Advanced Trading Features Migration completed successfully!")
}

func runAdvancedFeaturesMigration(db *sql.DB) error {
	// Split migration into manageable chunks
	migrations := []string{
		// Commission structures table
		`CREATE TABLE IF NOT EXISTS commission_structures (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			type ENUM('FLAT', 'PERCENTAGE', 'PER_SHARE', 'TIERED') NOT NULL,
			base_rate DECIMAL(8,4) NOT NULL DEFAULT 0.0000,
			minimum_fee DECIMAL(8,2) NOT NULL DEFAULT 0.00,
			maximum_fee DECIMAL(8,2) NULL,
			regulatory_fee DECIMAL(8,4) DEFAULT 0.0000,
			clearing_fee DECIMAL(8,2) DEFAULT 0.00,
			platform_fee DECIMAL(8,2) DEFAULT 0.00,
			market_data_fee DECIMAL(8,2) DEFAULT 0.00,
			inactivity_fee DECIMAL(8,2) DEFAULT 0.00,
			inactivity_period INT DEFAULT 90,
			options_rate DECIMAL(8,4) DEFAULT 0.0000,
			forex_rate DECIMAL(8,4) DEFAULT 0.0000,
			crypto_rate DECIMAL(8,4) DEFAULT 0.0000,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_name (name),
			INDEX idx_type (type),
			INDEX idx_active (is_active)
		)`,

		// Commission tiers table
		`CREATE TABLE IF NOT EXISTS commission_tiers (
			id INT AUTO_INCREMENT PRIMARY KEY,
			commission_structure_id INT NOT NULL,
			min_volume DECIMAL(15,2) NOT NULL,
			max_volume DECIMAL(15,2) NULL,
			rate DECIMAL(8,4) NOT NULL,
			min_fee DECIMAL(8,2) NOT NULL,
			max_fee DECIMAL(8,2) NULL,
			FOREIGN KEY (commission_structure_id) REFERENCES commission_structures(id) ON DELETE CASCADE,
			INDEX idx_structure (commission_structure_id),
			INDEX idx_volume (min_volume, max_volume)
		)`,

		// User commission profiles table
		`CREATE TABLE IF NOT EXISTS user_commission_profiles (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL UNIQUE,
			commission_structure_id INT NOT NULL,
			monthly_volume DECIMAL(15,2) DEFAULT 0.00,
			yearly_volume DECIMAL(15,2) DEFAULT 0.00,
			total_trades INT DEFAULT 0,
			vip_level INT DEFAULT 0,
			last_trade_date TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (commission_structure_id) REFERENCES commission_structures(id),
			INDEX idx_user (user_id),
			INDEX idx_structure (commission_structure_id),
			INDEX idx_vip_level (vip_level),
			INDEX idx_volume (monthly_volume, yearly_volume)
		)`,

		// Advanced orders table
		`CREATE TABLE IF NOT EXISTS advanced_orders (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			stock_symbol VARCHAR(10) NOT NULL,
			order_type ENUM('MARKET', 'LIMIT', 'STOP_LOSS', 'TAKE_PROFIT', 'TRAILING_STOP', 'OCO') NOT NULL,
			side ENUM('BUY', 'SELL', 'SHORT', 'COVER') NOT NULL,
			quantity INT NOT NULL,
			price DECIMAL(10,2) NULL,
			stop_price DECIMAL(10,2) NULL,
			trailing_amount DECIMAL(10,2) NULL,
			trailing_percent DECIMAL(5,2) NULL,
			time_in_force ENUM('GTC', 'IOC', 'FOK', 'DAY') DEFAULT 'GTC',
			status ENUM('PENDING', 'EXECUTED', 'CANCELLED', 'EXPIRED', 'PARTIALLY_FILLED') DEFAULT 'PENDING',
			executed_price DECIMAL(10,2) NULL,
			executed_quantity INT DEFAULT 0,
			remaining_quantity INT NOT NULL,
			executed_at TIMESTAMP NULL,
			expires_at TIMESTAMP NULL,
			parent_order_id INT NULL,
			linked_order_id INT NULL,
			commission DECIMAL(10,2) DEFAULT 0.00,
			fees DECIMAL(10,2) DEFAULT 0.00,
			market_price DECIMAL(10,2) NOT NULL,
			bid_price DECIMAL(10,2) NULL,
			ask_price DECIMAL(10,2) NULL,
			spread DECIMAL(10,2) NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (parent_order_id) REFERENCES advanced_orders(id) ON DELETE SET NULL,
			FOREIGN KEY (linked_order_id) REFERENCES advanced_orders(id) ON DELETE SET NULL,
			INDEX idx_user_id (user_id),
			INDEX idx_symbol (stock_symbol),
			INDEX idx_status (status),
			INDEX idx_order_type (order_type),
			INDEX idx_created_at (created_at),
			INDEX idx_expires_at (expires_at),
			INDEX idx_parent_order (parent_order_id),
			INDEX idx_linked_order (linked_order_id)
		)`,

		// Markets table
		`CREATE TABLE IF NOT EXISTS markets (
			id INT AUTO_INCREMENT PRIMARY KEY,
			code VARCHAR(10) NOT NULL UNIQUE,
			name VARCHAR(100) NOT NULL,
			type ENUM('STOCK', 'FOREX', 'CRYPTO', 'OPTION', 'FUTURE') NOT NULL,
			timezone VARCHAR(50) NOT NULL,
			currency VARCHAR(3) NOT NULL,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_code (code),
			INDEX idx_type (type),
			INDEX idx_active (is_active)
		)`,

		// Trading sessions table
		`CREATE TABLE IF NOT EXISTS trading_sessions (
			id INT AUTO_INCREMENT PRIMARY KEY,
			market_id INT NOT NULL,
			type ENUM('PRE_MARKET', 'REGULAR', 'AFTER_HOURS', 'OVERNIGHT') NOT NULL,
			start_time TIME NOT NULL,
			end_time TIME NOT NULL,
			days_of_week JSON NOT NULL,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (market_id) REFERENCES markets(id) ON DELETE CASCADE,
			INDEX idx_market (market_id),
			INDEX idx_type (type),
			INDEX idx_active (is_active)
		)`,

		// Price alerts table
		`CREATE TABLE IF NOT EXISTS price_alerts (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			symbol VARCHAR(10) NOT NULL,
			alert_type ENUM('ABOVE', 'BELOW', 'CHANGE_PERCENT') NOT NULL,
			trigger_price DECIMAL(10,2) NULL,
			trigger_percent DECIMAL(5,2) NULL,
			current_price DECIMAL(10,2) NOT NULL,
			is_triggered BOOLEAN DEFAULT FALSE,
			triggered_at TIMESTAMP NULL,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			INDEX idx_user (user_id),
			INDEX idx_symbol (symbol),
			INDEX idx_active (is_active),
			INDEX idx_triggered (is_triggered)
		)`,

		// User trading statistics table
		`CREATE TABLE IF NOT EXISTS user_trading_stats (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL UNIQUE,
			total_trades INT DEFAULT 0,
			winning_trades INT DEFAULT 0,
			losing_trades INT DEFAULT 0,
			total_volume DECIMAL(15,2) DEFAULT 0.00,
			total_commission DECIMAL(10,2) DEFAULT 0.00,
			total_fees DECIMAL(10,2) DEFAULT 0.00,
			average_holding_period INT DEFAULT 0,
			largest_win DECIMAL(10,2) DEFAULT 0.00,
			largest_loss DECIMAL(10,2) DEFAULT 0.00,
			max_drawdown DECIMAL(10,2) DEFAULT 0.00,
			sharpe_ratio DECIMAL(8,4) DEFAULT 0.0000,
			win_rate DECIMAL(5,2) DEFAULT 0.00,
			profit_factor DECIMAL(8,4) DEFAULT 0.0000,
			last_calculated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			INDEX idx_user (user_id),
			INDEX idx_win_rate (win_rate),
			INDEX idx_sharpe (sharpe_ratio),
			INDEX idx_calculated (last_calculated)
		)`,
	}

	// Execute each migration
	for i, migration := range migrations {
		fmt.Printf("Executing migration %d/%d...\n", i+1, len(migrations))
		_, err := db.Exec(migration)
		if err != nil {
			return fmt.Errorf("failed to execute migration %d: %v", i+1, err)
		}
	}

	// Insert default data
	fmt.Println("Inserting default data...")
	err := insertDefaultData(db)
	if err != nil {
		return fmt.Errorf("failed to insert default data: %v", err)
	}

	return nil
}

func insertDefaultData(db *sql.DB) error {
	// Insert default commission structures
	commissionStructures := [][]interface{}{
		{"Standard Retail", "FLAT", 0.00, 0.00, nil, 0.0021, 0.50, 1.00},
		{"Premium Retail", "PERCENTAGE", 0.25, 1.00, 25.00, 0.0021, 0.25, 0.50},
		{"Professional", "PER_SHARE", 0.005, 1.00, nil, 0.0021, 0.10, 0.00},
		{"VIP", "TIERED", 0.00, 0.00, nil, 0.0021, 0.00, 0.00},
	}

	query := `INSERT IGNORE INTO commission_structures (name, type, base_rate, minimum_fee, maximum_fee, regulatory_fee, clearing_fee, platform_fee) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	for _, data := range commissionStructures {
		_, err := db.Exec(query, data...)
		if err != nil {
			return fmt.Errorf("failed to insert commission structure: %v", err)
		}
	}

	// Insert default markets
	markets := [][]interface{}{
		{"NYSE", "New York Stock Exchange", "STOCK", "America/New_York", "USD"},
		{"NASDAQ", "NASDAQ Stock Market", "STOCK", "America/New_York", "USD"},
		{"FOREX", "Foreign Exchange Market", "FOREX", "UTC", "USD"},
		{"CRYPTO", "Cryptocurrency Market", "CRYPTO", "UTC", "USD"},
	}

	query = `INSERT IGNORE INTO markets (code, name, type, timezone, currency) VALUES (?, ?, ?, ?, ?)`
	for _, data := range markets {
		_, err := db.Exec(query, data...)
		if err != nil {
			return fmt.Errorf("failed to insert market: %v", err)
		}
	}

	// Insert trading sessions
	sessions := [][]interface{}{
		{1, "PRE_MARKET", "04:00:00", "09:30:00", "[1,2,3,4,5]"},
		{1, "REGULAR", "09:30:00", "16:00:00", "[1,2,3,4,5]"},
		{1, "AFTER_HOURS", "16:00:00", "20:00:00", "[1,2,3,4,5]"},
		{2, "PRE_MARKET", "04:00:00", "09:30:00", "[1,2,3,4,5]"},
		{2, "REGULAR", "09:30:00", "16:00:00", "[1,2,3,4,5]"},
		{2, "AFTER_HOURS", "16:00:00", "20:00:00", "[1,2,3,4,5]"},
		{3, "REGULAR", "00:00:00", "23:59:59", "[1,2,3,4,5]"},
		{4, "REGULAR", "00:00:00", "23:59:59", "[0,1,2,3,4,5,6]"},
	}

	query = `INSERT IGNORE INTO trading_sessions (market_id, type, start_time, end_time, days_of_week) VALUES (?, ?, ?, ?, ?)`
	for _, data := range sessions {
		_, err := db.Exec(query, data...)
		if err != nil {
			return fmt.Errorf("failed to insert trading session: %v", err)
		}
	}

	// Initialize commission profiles for existing users
	_, err := db.Exec(`INSERT IGNORE INTO user_commission_profiles (user_id, commission_structure_id) 
		SELECT id, 1 FROM users WHERE id NOT IN (SELECT user_id FROM user_commission_profiles)`)
	if err != nil {
		return fmt.Errorf("failed to initialize user commission profiles: %v", err)
	}

	// Initialize trading stats for existing users
	_, err = db.Exec(`INSERT IGNORE INTO user_trading_stats (user_id) 
		SELECT id FROM users WHERE id NOT IN (SELECT user_id FROM user_trading_stats)`)
	if err != nil {
		return fmt.Errorf("failed to initialize user trading stats: %v", err)
	}

	return nil
} 