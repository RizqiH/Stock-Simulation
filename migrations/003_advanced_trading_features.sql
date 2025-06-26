-- Migration: Advanced Trading Features
-- Created: 2024
-- Description: Adds comprehensive trading features including advanced orders, commissions, market hours, and real-time updates

-- =====================================================
-- ADVANCED ORDERS TABLES
-- =====================================================

-- Enhanced orders table with advanced order types
DROP TABLE IF EXISTS advanced_orders;
CREATE TABLE advanced_orders (
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
    
    -- OCO Related
    parent_order_id INT NULL,
    linked_order_id INT NULL,
    
    -- Commission and fees
    commission DECIMAL(10,2) DEFAULT 0.00,
    fees DECIMAL(10,2) DEFAULT 0.00,
    
    -- Market data at time of order
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
);

-- =====================================================
-- COMMISSION & FEES TABLES
-- =====================================================

-- Commission structures table
CREATE TABLE commission_structures (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type ENUM('FLAT', 'PERCENTAGE', 'PER_SHARE', 'TIERED') NOT NULL,
    base_rate DECIMAL(8,4) NOT NULL DEFAULT 0.0000,
    minimum_fee DECIMAL(8,2) NOT NULL DEFAULT 0.00,
    maximum_fee DECIMAL(8,2) NULL,
    
    -- Additional fees
    regulatory_fee DECIMAL(8,4) DEFAULT 0.0000,
    clearing_fee DECIMAL(8,2) DEFAULT 0.00,
    platform_fee DECIMAL(8,2) DEFAULT 0.00,
    
    -- Market data fees
    market_data_fee DECIMAL(8,2) DEFAULT 0.00,
    
    -- Account fees
    inactivity_fee DECIMAL(8,2) DEFAULT 0.00,
    inactivity_period INT DEFAULT 90,
    
    -- Special rates
    options_rate DECIMAL(8,4) DEFAULT 0.0000,
    forex_rate DECIMAL(8,4) DEFAULT 0.0000,
    crypto_rate DECIMAL(8,4) DEFAULT 0.0000,
    
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_name (name),
    INDEX idx_type (type),
    INDEX idx_active (is_active)
);

-- Commission tiers table
CREATE TABLE commission_tiers (
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
);

-- User commission profiles table
CREATE TABLE user_commission_profiles (
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
);

-- =====================================================
-- MARKET HOURS & TRADING SESSIONS TABLES
-- =====================================================

-- Markets table
CREATE TABLE markets (
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
);

-- Trading sessions table
CREATE TABLE trading_sessions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    market_id INT NOT NULL,
    type ENUM('PRE_MARKET', 'REGULAR', 'AFTER_HOURS', 'OVERNIGHT') NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    days_of_week JSON NOT NULL, -- [1,2,3,4,5] for Mon-Fri
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (market_id) REFERENCES markets(id) ON DELETE CASCADE,
    INDEX idx_market (market_id),
    INDEX idx_type (type),
    INDEX idx_active (is_active)
);

-- Market holidays table
CREATE TABLE market_holidays (
    id INT AUTO_INCREMENT PRIMARY KEY,
    market_id INT NOT NULL,
    date DATE NOT NULL,
    name VARCHAR(100) NOT NULL,
    type ENUM('FULL_CLOSE', 'EARLY_CLOSE') NOT NULL,
    early_close_time TIME NULL,
    is_recurring BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (market_id) REFERENCES markets(id) ON DELETE CASCADE,
    INDEX idx_market (market_id),
    INDEX idx_date (date),
    INDEX idx_type (type),
    UNIQUE KEY unique_market_date (market_id, date)
);

-- Market data permissions table
CREATE TABLE market_data_permissions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    market_code VARCHAR(10) NOT NULL,
    data_type ENUM('REAL_TIME', 'DELAYED', 'SNAPSHOT') NOT NULL,
    permission_level ENUM('BASIC', 'PREMIUM', 'PROFESSIONAL') NOT NULL,
    expires_at TIMESTAMP NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user (user_id),
    INDEX idx_market (market_code),
    INDEX idx_expires (expires_at),
    INDEX idx_active (is_active),
    UNIQUE KEY unique_user_market_type (user_id, market_code, data_type)
);

-- Trading restrictions table
CREATE TABLE trading_restrictions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    market_code VARCHAR(10) NOT NULL,
    symbol VARCHAR(10) NULL, -- NULL for market-wide restrictions
    restriction_type ENUM('HALT', 'SUSPENSION', 'LIMIT_UP_DOWN', 'CIRCUIT_BREAKER') NOT NULL,
    reason TEXT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_market (market_code),
    INDEX idx_symbol (symbol),
    INDEX idx_type (restriction_type),
    INDEX idx_active (is_active),
    INDEX idx_time_range (start_time, end_time)
);

-- =====================================================
-- WEBSOCKET & REAL-TIME TABLES
-- =====================================================

-- WebSocket connections table
CREATE TABLE websocket_connections (
    id VARCHAR(36) PRIMARY KEY,
    user_id INT NULL,
    connected_at TIMESTAMP NOT NULL,
    last_heartbeat TIMESTAMP NOT NULL,
    subscriptions JSON NULL,
    is_active BOOLEAN DEFAULT TRUE,
    client_info JSON NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_user (user_id),
    INDEX idx_active (is_active),
    INDEX idx_heartbeat (last_heartbeat)
);

-- Price alerts table
CREATE TABLE price_alerts (
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
);

-- News alerts table
CREATE TABLE news_alerts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    summary TEXT,
    content TEXT,
    category VARCHAR(50) NOT NULL,
    severity ENUM('LOW', 'MEDIUM', 'HIGH', 'CRITICAL') NOT NULL,
    symbols JSON NULL, -- Affected symbols
    source VARCHAR(100) NOT NULL,
    url TEXT NULL,
    published_at TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_category (category),
    INDEX idx_severity (severity),
    INDEX idx_published (published_at),
    INDEX idx_active (is_active),
    FULLTEXT idx_content (title, summary, content)
);

-- =====================================================
-- ANALYTICS & PERFORMANCE TABLES
-- =====================================================

-- User trading statistics table
CREATE TABLE user_trading_stats (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL UNIQUE,
    total_trades INT DEFAULT 0,
    winning_trades INT DEFAULT 0,
    losing_trades INT DEFAULT 0,
    total_volume DECIMAL(15,2) DEFAULT 0.00,
    total_commission DECIMAL(10,2) DEFAULT 0.00,
    total_fees DECIMAL(10,2) DEFAULT 0.00,
    average_holding_period INT DEFAULT 0, -- in hours
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
);

-- Portfolio performance snapshots table
CREATE TABLE portfolio_snapshots (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    snapshot_date DATE NOT NULL,
    total_value DECIMAL(15,2) NOT NULL,
    cash_balance DECIMAL(15,2) NOT NULL,
    positions_value DECIMAL(15,2) NOT NULL,
    day_pnl DECIMAL(10,2) DEFAULT 0.00,
    total_pnl DECIMAL(10,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user (user_id),
    INDEX idx_date (snapshot_date),
    UNIQUE KEY unique_user_date (user_id, snapshot_date)
);

-- =====================================================
-- SOCIAL TRADING TABLES
-- =====================================================

-- User follows table (for social trading)
CREATE TABLE user_follows (
    id INT AUTO_INCREMENT PRIMARY KEY,
    follower_id INT NOT NULL,
    following_id INT NOT NULL,
    notification_enabled BOOLEAN DEFAULT TRUE,
    copy_trading_enabled BOOLEAN DEFAULT FALSE,
    copy_percentage DECIMAL(5,2) DEFAULT 0.00, -- Percentage of portfolio to copy
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_follower (follower_id),
    INDEX idx_following (following_id),
    UNIQUE KEY unique_follow (follower_id, following_id)
);

-- Social trading posts table
CREATE TABLE social_posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    symbols JSON NULL, -- Related symbols
    trade_id INT NULL, -- Related trade if any
    likes_count INT DEFAULT 0,
    comments_count INT DEFAULT 0,
    shares_count INT DEFAULT 0,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user (user_id),
    INDEX idx_public (is_public),
    INDEX idx_created (created_at),
    FULLTEXT idx_content (content)
);

-- =====================================================
-- SYSTEM CONFIGURATION TABLES
-- =====================================================

-- System settings table
CREATE TABLE system_settings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    setting_key VARCHAR(100) NOT NULL UNIQUE,
    setting_value TEXT NOT NULL,
    data_type ENUM('STRING', 'INTEGER', 'FLOAT', 'BOOLEAN', 'JSON') NOT NULL,
    description TEXT,
    is_public BOOLEAN DEFAULT FALSE, -- Whether setting can be read by frontend
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_key (setting_key),
    INDEX idx_public (is_public)
);

-- =====================================================
-- INSERT DEFAULT DATA
-- =====================================================

-- Insert default commission structure
INSERT INTO commission_structures (name, type, base_rate, minimum_fee, maximum_fee, regulatory_fee, clearing_fee, platform_fee) VALUES
('Standard Retail', 'FLAT', 0.00, 0.00, NULL, 0.0021, 0.50, 1.00),
('Premium Retail', 'PERCENTAGE', 0.25, 1.00, 25.00, 0.0021, 0.25, 0.50),
('Professional', 'PER_SHARE', 0.005, 1.00, NULL, 0.0021, 0.10, 0.00),
('VIP', 'TIERED', 0.00, 0.00, NULL, 0.0021, 0.00, 0.00);

-- Insert commission tiers for VIP structure
INSERT INTO commission_tiers (commission_structure_id, min_volume, max_volume, rate, min_fee, max_fee) VALUES
(4, 0, 10000, 0.50, 2.00, 15.00),
(4, 10000, 50000, 0.25, 1.00, 10.00),
(4, 50000, 100000, 0.15, 0.50, 5.00),
(4, 100000, NULL, 0.10, 0.00, 2.00);

-- Insert default markets
INSERT INTO markets (code, name, type, timezone, currency) VALUES
('NYSE', 'New York Stock Exchange', 'STOCK', 'America/New_York', 'USD'),
('NASDAQ', 'NASDAQ Stock Market', 'STOCK', 'America/New_York', 'USD'),
('LSE', 'London Stock Exchange', 'STOCK', 'Europe/London', 'GBP'),
('TSE', 'Tokyo Stock Exchange', 'STOCK', 'Asia/Tokyo', 'JPY'),
('FOREX', 'Foreign Exchange Market', 'FOREX', 'UTC', 'USD'),
('CRYPTO', 'Cryptocurrency Market', 'CRYPTO', 'UTC', 'USD');

-- Insert trading sessions for NYSE/NASDAQ
INSERT INTO trading_sessions (market_id, type, start_time, end_time, days_of_week) VALUES
(1, 'PRE_MARKET', '04:00:00', '09:30:00', '[1,2,3,4,5]'),
(1, 'REGULAR', '09:30:00', '16:00:00', '[1,2,3,4,5]'),
(1, 'AFTER_HOURS', '16:00:00', '20:00:00', '[1,2,3,4,5]'),
(2, 'PRE_MARKET', '04:00:00', '09:30:00', '[1,2,3,4,5]'),
(2, 'REGULAR', '09:30:00', '16:00:00', '[1,2,3,4,5]'),
(2, 'AFTER_HOURS', '16:00:00', '20:00:00', '[1,2,3,4,5]');

-- Insert forex sessions (24/5 market)
INSERT INTO trading_sessions (market_id, type, start_time, end_time, days_of_week) VALUES
(5, 'REGULAR', '00:00:00', '23:59:59', '[1,2,3,4,5]');

-- Insert crypto sessions (24/7 market)
INSERT INTO trading_sessions (market_id, type, start_time, end_time, days_of_week) VALUES
(6, 'REGULAR', '00:00:00', '23:59:59', '[0,1,2,3,4,5,6]');

-- Insert some major US holidays
INSERT INTO market_holidays (market_id, date, name, type) VALUES
(1, '2024-01-01', 'New Years Day', 'FULL_CLOSE'),
(1, '2024-07-04', 'Independence Day', 'FULL_CLOSE'),
(1, '2024-12-25', 'Christmas Day', 'FULL_CLOSE'),
(1, '2024-11-29', 'Black Friday', 'EARLY_CLOSE'),
(2, '2024-01-01', 'New Years Day', 'FULL_CLOSE'),
(2, '2024-07-04', 'Independence Day', 'FULL_CLOSE'),
(2, '2024-12-25', 'Christmas Day', 'FULL_CLOSE'),
(2, '2024-11-29', 'Black Friday', 'EARLY_CLOSE');

-- Insert default system settings
INSERT INTO system_settings (setting_key, setting_value, data_type, description, is_public) VALUES
('market_open_buffer_minutes', '5', 'INTEGER', 'Buffer time before market opens for order placement', FALSE),
('max_order_quantity', '10000', 'INTEGER', 'Maximum quantity per order', TRUE),
('websocket_heartbeat_interval', '30', 'INTEGER', 'WebSocket heartbeat interval in seconds', FALSE),
('price_update_interval', '1000', 'INTEGER', 'Price update interval in milliseconds', TRUE),
('max_websocket_connections_per_user', '5', 'INTEGER', 'Maximum WebSocket connections per user', FALSE),
('trading_enabled', 'true', 'BOOLEAN', 'Whether trading is enabled system-wide', TRUE),
('maintenance_mode', 'false', 'BOOLEAN', 'Whether system is in maintenance mode', TRUE);

-- Update existing users with default commission profile
INSERT INTO user_commission_profiles (user_id, commission_structure_id) 
SELECT id, 1 FROM users WHERE id NOT IN (SELECT user_id FROM user_commission_profiles);

-- Initialize trading stats for existing users
INSERT INTO user_trading_stats (user_id) 
SELECT id FROM users WHERE id NOT IN (SELECT user_id FROM user_trading_stats); 