-- Quick Database Setup Script
-- Run this to reset and setup the database with default data

USE stock_simulation;

-- Source the main migration file
SOURCE /docker-entrypoint-initdb.d/000_initial_setup.sql;

-- Verify setup
SELECT 'Setup verification:' as info;
SELECT 'Users created:' as table_name, COUNT(*) as count FROM users;
SELECT 'Stocks loaded:' as table_name, COUNT(*) as count FROM stocks;
SELECT 'Portfolios created:' as table_name, COUNT(*) as count FROM portfolio;
SELECT 'Transactions created:' as table_name, COUNT(*) as count FROM transactions;

-- Show default login credentials
SELECT 'Default Login Credentials:' as info;
SELECT 'Username/Email' as credential, 'Password' as value
UNION ALL
SELECT 'admin@stocksim.com', 'password123'
UNION ALL
SELECT 'john@example.com', 'password123'
UNION ALL
SELECT 'jane@example.com', 'password123'
UNION ALL
SELECT 'demo@stocksim.com', 'password123'; 