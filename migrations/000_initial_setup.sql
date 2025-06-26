-- Initial Database Setup Migration
-- This will recreate the complete database structure with default data

-- Drop existing tables if they exist (for clean setup)
DROP TABLE IF EXISTS advanced_orders;
DROP TABLE IF EXISTS oco_orders;
DROP TABLE IF EXISTS user_commission_profiles;
DROP TABLE IF EXISTS commission_structures;
DROP TABLE IF EXISTS markets;
DROP TABLE IF EXISTS historical_prices;
DROP TABLE IF EXISTS portfolio;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS stocks;
DROP TABLE IF EXISTS users;

-- Create users table
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL DEFAULT '',
    balance DECIMAL(15,2) DEFAULT 10000.00,
    total_profit DECIMAL(15,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_username (username)
);

-- Create stocks table
CREATE TABLE stocks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    current_price DECIMAL(10,2) NOT NULL,
    previous_close DECIMAL(10,2) NOT NULL,
    market_cap BIGINT,
    volume BIGINT DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_symbol (symbol)
);

-- Create transactions table
CREATE TABLE transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    stock_symbol VARCHAR(10) NOT NULL,
    transaction_type ENUM('buy', 'sell') NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    total_amount DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_stock_symbol (stock_symbol),
    INDEX idx_created_at (created_at)
);

-- Create portfolio table
CREATE TABLE portfolio (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    stock_symbol VARCHAR(10) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    average_price DECIMAL(10,2) NOT NULL,
    total_cost DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_stock (user_id, stock_symbol),
    INDEX idx_user_id (user_id)
);

-- Create historical_prices table
CREATE TABLE historical_prices (
    id INT AUTO_INCREMENT PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    volume BIGINT DEFAULT 0,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (symbol) REFERENCES stocks(symbol) ON DELETE CASCADE,
    INDEX idx_symbol_timestamp (symbol, timestamp),
    INDEX idx_timestamp (timestamp)
);

-- Create advanced_orders table for advanced trading features
CREATE TABLE advanced_orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    stock_symbol VARCHAR(10) NOT NULL,
    order_type ENUM('MARKET', 'LIMIT', 'STOP_LOSS', 'TAKE_PROFIT', 'TRAILING_STOP') NOT NULL,
    side ENUM('BUY', 'SELL', 'SHORT', 'COVER') NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10,2),
    stop_price DECIMAL(10,2),
    trailing_amount DECIMAL(10,2),
    trailing_percent DECIMAL(5,2),
    time_in_force ENUM('DAY', 'GTC', 'IOC', 'FOK') DEFAULT 'DAY',
    status ENUM('PENDING', 'PARTIALLY_FILLED', 'FILLED', 'CANCELLED', 'REJECTED', 'EXPIRED') DEFAULT 'PENDING',
    executed_price DECIMAL(10,2),
    executed_quantity INT DEFAULT 0,
    remaining_quantity INT,
    market_price DECIMAL(10,2),
    bid_price DECIMAL(10,2),
    ask_price DECIMAL(10,2),
    commission DECIMAL(10,2) DEFAULT 0.00,
    fees DECIMAL(10,2) DEFAULT 0.00,
    spread DECIMAL(10,4) DEFAULT 0.00,
    executed_at TIMESTAMP NULL,
    expires_at TIMESTAMP NULL,
    parent_order_id INT,
    linked_order_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_stock_symbol (stock_symbol),
    INDEX idx_status (status),
    INDEX idx_order_type (order_type),
    INDEX idx_side (side),
    INDEX idx_created_at (created_at)
);

-- Create oco_orders table for One-Cancels-Other orders
CREATE TABLE oco_orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    group_id VARCHAR(50) NOT NULL UNIQUE,
    user_id INT NOT NULL,
    primary_order_id INT NOT NULL,
    secondary_order_id INT NOT NULL,
    status ENUM('ACTIVE', 'TRIGGERED', 'CANCELLED') DEFAULT 'ACTIVE',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (primary_order_id) REFERENCES advanced_orders(id) ON DELETE CASCADE,
    FOREIGN KEY (secondary_order_id) REFERENCES advanced_orders(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_group_id (group_id),
    INDEX idx_status (status)
);

-- Create markets table for trading hours and market info
CREATE TABLE markets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    timezone VARCHAR(50) NOT NULL,
    open_time TIME NOT NULL,
    close_time TIME NOT NULL,
    trading_days VARCHAR(20) DEFAULT 'MON-FRI',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_name (name),
    INDEX idx_is_active (is_active)
);

-- Create commission_structures table
CREATE TABLE commission_structures (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    market_id INT NOT NULL,
    tier_name VARCHAR(50) NOT NULL,
    min_volume INT DEFAULT 0,
    max_volume INT,
    fixed_fee DECIMAL(10,2) DEFAULT 0.00,
    percentage_fee DECIMAL(5,4) DEFAULT 0.00,
    min_commission DECIMAL(10,2) DEFAULT 0.00,
    max_commission DECIMAL(10,2),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (market_id) REFERENCES markets(id) ON DELETE CASCADE,
    INDEX idx_market_tier (market_id, tier_name),
    INDEX idx_volume_range (min_volume, max_volume)
);

-- Create user_commission_profiles table
CREATE TABLE user_commission_profiles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    market_id INT NOT NULL,
    tier_name VARCHAR(50) NOT NULL DEFAULT 'STANDARD',
    total_volume INT DEFAULT 0,
    monthly_volume INT DEFAULT 0,
    vip_level ENUM('NONE', 'BRONZE', 'SILVER', 'GOLD', 'PLATINUM') DEFAULT 'NONE',
    discount_percentage DECIMAL(5,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (market_id) REFERENCES markets(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_market (user_id, market_id),
    INDEX idx_user_id (user_id),
    INDEX idx_tier (tier_name)
);

-- Insert default users with hashed passwords
-- Password for all users: "password123"
-- Hash generated using bcrypt with cost 10
INSERT INTO users (username, email, password_hash, full_name, balance, total_profit) VALUES
('admin', 'admin@stocksim.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Administrator', 100000.00, 0.00),
('john_doe', 'john@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'John Doe', 50000.00, 2500.00),
('jane_smith', 'jane@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Jane Smith', 75000.00, 5000.00),
('demo_user', 'demo@stocksim.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Demo User', 25000.00, 1000.00);

-- Insert sample stocks with realistic data
INSERT INTO stocks (symbol, name, current_price, previous_close, market_cap, volume) VALUES
('AAPL', 'Apple Inc.', 150.25, 148.75, 2500000000000, 52000000),
('GOOGL', 'Alphabet Inc.', 2835.50, 2798.20, 1800000000000, 28000000),
('MSFT', 'Microsoft Corporation', 305.80, 298.45, 2200000000000, 35000000),
('AMZN', 'Amazon.com Inc.', 3250.75, 3180.25, 1600000000000, 22000000),
('TSLA', 'Tesla Inc.', 820.40, 795.60, 800000000000, 45000000),
('META', 'Meta Platforms Inc.', 325.90, 320.15, 850000000000, 38000000),
('NVDA', 'NVIDIA Corporation', 465.80, 448.25, 1100000000000, 48000000),
('NFLX', 'Netflix Inc.', 410.50, 398.75, 180000000000, 18000000),
('BABA', 'Alibaba Group', 92.30, 89.85, 240000000000, 28000000),
('V', 'Visa Inc.', 225.60, 221.40, 480000000000, 12000000),
('JPM', 'JPMorgan Chase & Co.', 165.75, 162.30, 485000000000, 15000000),
('JNJ', 'Johnson & Johnson', 175.25, 173.80, 460000000000, 8500000),
('WMT', 'Walmart Inc.', 145.90, 144.25, 400000000000, 9200000),
('PG', 'Procter & Gamble', 155.40, 153.95, 372000000000, 7800000),
('UNH', 'UnitedHealth Group', 520.75, 515.20, 490000000000, 6500000);

-- Insert sample portfolio data for demo users
INSERT INTO portfolio (user_id, stock_symbol, quantity, average_price, total_cost) VALUES
(2, 'AAPL', 50, 145.00, 7250.00),
(2, 'GOOGL', 5, 2750.00, 13750.00),
(2, 'TSLA', 10, 780.00, 7800.00),
(3, 'MSFT', 25, 290.00, 7250.00),
(3, 'NVDA', 15, 420.00, 6300.00),
(3, 'META', 20, 310.00, 6200.00),
(4, 'AAPL', 20, 148.00, 2960.00),
(4, 'AMZN', 2, 3100.00, 6200.00);

-- Insert sample transaction history
INSERT INTO transactions (user_id, stock_symbol, transaction_type, quantity, price, total_amount) VALUES
(2, 'AAPL', 'buy', 50, 145.00, 7250.00),
(2, 'GOOGL', 'buy', 5, 2750.00, 13750.00),
(2, 'TSLA', 'buy', 10, 780.00, 7800.00),
(3, 'MSFT', 'buy', 25, 290.00, 7250.00),
(3, 'NVDA', 'buy', 15, 420.00, 6300.00),
(3, 'META', 'buy', 20, 310.00, 6200.00),
(4, 'AAPL', 'buy', 20, 148.00, 2960.00),
(4, 'AMZN', 'buy', 2, 3100.00, 6200.00);

-- Insert historical price data (sample data for charts)
INSERT INTO historical_prices (symbol, price, volume, timestamp) VALUES
-- AAPL historical data (last 30 days)
('AAPL', 145.50, 50000000, '2025-05-27 09:30:00'),
('AAPL', 147.25, 52000000, '2025-05-28 09:30:00'),
('AAPL', 146.80, 48000000, '2025-05-29 09:30:00'),
('AAPL', 148.75, 55000000, '2025-05-30 09:30:00'),
('AAPL', 150.25, 52000000, '2025-06-26 09:30:00'),
-- GOOGL historical data
('GOOGL', 2750.00, 25000000, '2025-05-27 09:30:00'),
('GOOGL', 2798.20, 28000000, '2025-05-28 09:30:00'),
('GOOGL', 2835.50, 28000000, '2025-06-26 09:30:00'),
-- TSLA historical data
('TSLA', 780.00, 40000000, '2025-05-27 09:30:00'),
('TSLA', 795.60, 42000000, '2025-05-28 09:30:00'),
('TSLA', 820.40, 45000000, '2025-06-26 09:30:00');

-- Insert default markets
INSERT INTO markets (name, timezone, open_time, close_time, trading_days, is_active) VALUES
('NYSE', 'America/New_York', '09:30:00', '16:00:00', 'MON-FRI', TRUE),
('NASDAQ', 'America/New_York', '09:30:00', '16:00:00', 'MON-FRI', TRUE),
('FOREX', 'UTC', '00:00:00', '23:59:59', 'MON-FRI', TRUE),
('CRYPTO', 'UTC', '00:00:00', '23:59:59', 'MON-SUN', TRUE);

-- Insert commission structures
INSERT INTO commission_structures (name, market_id, tier_name, min_volume, max_volume, fixed_fee, percentage_fee, min_commission, max_commission, is_active) VALUES
-- NYSE/NASDAQ Standard Tiers
('US Stocks Standard', 1, 'BASIC', 0, 10000, 0.00, 0.0025, 1.00, NULL, TRUE),
('US Stocks Standard', 1, 'STANDARD', 10001, 100000, 0.00, 0.0020, 1.00, NULL, TRUE),
('US Stocks Standard', 1, 'PREMIUM', 100001, 1000000, 0.00, 0.0015, 1.00, NULL, TRUE),
('US Stocks Standard', 1, 'VIP', 1000001, NULL, 0.00, 0.0010, 1.00, NULL, TRUE),

('US Stocks Standard', 2, 'BASIC', 0, 10000, 0.00, 0.0025, 1.00, NULL, TRUE),
('US Stocks Standard', 2, 'STANDARD', 10001, 100000, 0.00, 0.0020, 1.00, NULL, TRUE),
('US Stocks Standard', 2, 'PREMIUM', 100001, 1000000, 0.00, 0.0015, 1.00, NULL, TRUE),
('US Stocks Standard', 2, 'VIP', 1000001, NULL, 0.00, 0.0010, 1.00, NULL, TRUE),

-- FOREX
('Forex Standard', 3, 'BASIC', 0, 100000, 0.00, 0.0001, 0.10, NULL, TRUE),
('Forex Standard', 3, 'STANDARD', 100001, 1000000, 0.00, 0.00008, 0.10, NULL, TRUE),
('Forex Standard', 3, 'PREMIUM', 1000001, NULL, 0.00, 0.00005, 0.10, NULL, TRUE),

-- CRYPTO
('Crypto Standard', 4, 'BASIC', 0, 10000, 0.00, 0.0050, 0.01, NULL, TRUE),
('Crypto Standard', 4, 'STANDARD', 10001, 100000, 0.00, 0.0040, 0.01, NULL, TRUE),
('Crypto Standard', 4, 'PREMIUM', 100001, NULL, 0.00, 0.0030, 0.01, NULL, TRUE);

-- Insert user commission profiles for existing users
INSERT INTO user_commission_profiles (user_id, market_id, tier_name, total_volume, monthly_volume, vip_level, discount_percentage) VALUES
-- Admin user - VIP access to all markets
(1, 1, 'VIP', 5000000, 500000, 'PLATINUM', 20.00),
(1, 2, 'VIP', 5000000, 500000, 'PLATINUM', 20.00),
(1, 3, 'PREMIUM', 2000000, 200000, 'GOLD', 15.00),
(1, 4, 'PREMIUM', 1000000, 100000, 'GOLD', 15.00),

-- John Doe - Standard user
(2, 1, 'STANDARD', 50000, 5000, 'BRONZE', 5.00),
(2, 2, 'STANDARD', 50000, 5000, 'BRONZE', 5.00),
(2, 3, 'BASIC', 10000, 1000, 'NONE', 0.00),
(2, 4, 'BASIC', 25000, 2500, 'NONE', 0.00),

-- Jane Smith - Premium user
(3, 1, 'PREMIUM', 150000, 15000, 'SILVER', 10.00),
(3, 2, 'PREMIUM', 150000, 15000, 'SILVER', 10.00),
(3, 3, 'STANDARD', 75000, 7500, 'BRONZE', 5.00),
(3, 4, 'STANDARD', 100000, 10000, 'BRONZE', 5.00),

-- Demo user - Basic access
(4, 1, 'BASIC', 25000, 2500, 'NONE', 0.00),
(4, 2, 'BASIC', 25000, 2500, 'NONE', 0.00),
(4, 3, 'BASIC', 5000, 500, 'NONE', 0.00),
(4, 4, 'BASIC', 10000, 1000, 'NONE', 0.00);

-- Create indexes for better performance
CREATE INDEX idx_historical_prices_symbol_date ON historical_prices(symbol, timestamp DESC);
CREATE INDEX idx_transactions_user_date ON transactions(user_id, created_at DESC);
CREATE INDEX idx_portfolio_user_symbol ON portfolio(user_id, stock_symbol);

-- Show setup completion message
SELECT 'Database setup completed successfully!' as status,
       COUNT(*) as total_users FROM users
UNION ALL
SELECT 'Total stocks loaded:' as status,
       COUNT(*) as count FROM stocks
UNION ALL
SELECT 'Sample portfolios created:' as status,
       COUNT(*) as count FROM portfolio; 