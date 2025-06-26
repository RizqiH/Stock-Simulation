-- Random Price Movement Simulator
-- This script applies random price changes to all stocks

-- Update AAPL (randomly between $145-155)
UPDATE stocks SET 
    current_price = 150.0 + (RAND() - 0.5) * 10, 
    updated_at = NOW() 
WHERE symbol = 'AAPL';

-- Update GOOGL (randomly between $2350-2450)
UPDATE stocks SET 
    current_price = 2400.0 + (RAND() - 0.5) * 100, 
    updated_at = NOW() 
WHERE symbol = 'GOOGL';

-- Update TSLA (randomly between $145-155)
UPDATE stocks SET 
    current_price = 150.0 + (RAND() - 0.5) * 10, 
    updated_at = NOW() 
WHERE symbol = 'TSLA';

-- Update MSFT (randomly between $320-340)
UPDATE stocks SET 
    current_price = 330.0 + (RAND() - 0.5) * 20, 
    updated_at = NOW() 
WHERE symbol = 'MSFT';

-- Update META (randomly between $290-310)
UPDATE stocks SET 
    current_price = 300.0 + (RAND() - 0.5) * 20, 
    updated_at = NOW() 
WHERE symbol = 'META';

-- Show all updated prices
SELECT 
    symbol, 
    ROUND(current_price, 2) as current_price, 
    updated_at,
    CONCAT(
        CASE 
            WHEN current_price > previous_close THEN '📈 +'
            WHEN current_price < previous_close THEN '📉 '
            ELSE '➡️ '
        END,
        ROUND(((current_price - previous_close) / previous_close * 100), 2), 
        '%'
    ) as change
FROM stocks 
ORDER BY symbol; 