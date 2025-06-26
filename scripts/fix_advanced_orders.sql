-- Fix advanced_orders table structure
USE stock_simulation;

-- Disable foreign key checks
SET FOREIGN_KEY_CHECKS = 0;

-- Drop existing tables
DROP TABLE IF EXISTS oco_orders;
DROP TABLE IF EXISTS advanced_orders;

-- Create advanced_orders table with correct column names matching repository
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

-- Create oco_orders table
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

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;

-- Show table structure to verify
DESCRIBE advanced_orders;

SELECT 'Advanced orders table recreated successfully!' as status; 