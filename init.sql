-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    balance DECIMAL(15,2) DEFAULT 10000.00,
    total_profit DECIMAL(15,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_username (username)
);

-- Create stocks table
CREATE TABLE IF NOT EXISTS stocks (
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
CREATE TABLE IF NOT EXISTS transactions (
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
CREATE TABLE IF NOT EXISTS portfolio (
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

-- Insert sample stocks
INSERT IGNORE INTO stocks (symbol, name, current_price, previous_close, market_cap, volume) VALUES
('AAPL', 'Apple Inc.', 150.00, 148.50, 2500000000000, 50000000),
('GOOGL', 'Alphabet Inc.', 2800.00, 2750.00, 1800000000000, 25000000),
('MSFT', 'Microsoft Corporation', 300.00, 295.00, 2200000000000, 30000000),
('AMZN', 'Amazon.com Inc.', 3200.00, 3150.00, 1600000000000, 20000000),
('TSLA', 'Tesla Inc.', 800.00, 780.00, 800000000000, 40000000),
('META', 'Meta Platforms Inc.', 320.00, 315.00, 850000000000, 35000000),
('NVDA', 'NVIDIA Corporation', 450.00, 440.00, 1100000000000, 45000000),
('NFLX', 'Netflix Inc.', 400.00, 395.00, 180000000000, 15000000),
('BABA', 'Alibaba Group', 90.00, 88.50, 240000000000, 25000000),
('V', 'Visa Inc.', 220.00, 218.00, 480000000000, 10000000);