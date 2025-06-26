-- Add total_cost column to portfolio table
ALTER TABLE portfolio ADD COLUMN IF NOT EXISTS total_cost DECIMAL(15,2) NOT NULL DEFAULT 0.00;

-- Update existing records to calculate total_cost based on quantity * average_price
UPDATE portfolio SET total_cost = quantity * average_price WHERE total_cost = 0.00;