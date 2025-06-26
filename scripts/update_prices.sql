-- Manual Stock Price Updates for Testing Limit Orders
-- Run this script to change stock prices and test limit order execution

-- Scenario 1: Lower AAPL price to trigger BUY LIMIT orders
UPDATE stocks SET current_price = 148.75, updated_at = NOW() WHERE symbol = 'AAPL';

-- Scenario 2: Raise GOOGL price to trigger SELL LIMIT orders  
UPDATE stocks SET current_price = 2425.30, updated_at = NOW() WHERE symbol = 'GOOGL';

-- Scenario 3: Mix of price changes
UPDATE stocks SET current_price = 151.80, updated_at = NOW() WHERE symbol = 'TSLA';
UPDATE stocks SET current_price = 331.45, updated_at = NOW() WHERE symbol = 'MSFT';
UPDATE stocks SET current_price = 298.60, updated_at = NOW() WHERE symbol = 'META';

-- Show updated prices
SELECT symbol, current_price, updated_at FROM stocks ORDER BY symbol; 