-- Insert sample stock data
USE stock_simulation;

-- Insert popular Indonesian stocks
INSERT INTO stocks (symbol, name, current_price, open_price, high_price, low_price, volume, market_cap, sector) VALUES
('BBCA', 'Bank Central Asia Tbk', 8750.00, 8700.00, 8800.00, 8650.00, 1500000, 1050000000000.00, 'Financial Services'),
('BBRI', 'Bank Rakyat Indonesia Tbk', 4580.00, 4550.00, 4620.00, 4520.00, 2100000, 550000000000.00, 'Financial Services'),
('BMRI', 'Bank Mandiri Tbk', 9175.00, 9150.00, 9200.00, 9100.00, 1800000, 450000000000.00, 'Financial Services'),
('TLKM', 'Telkom Indonesia Tbk', 3640.00, 3620.00, 3680.00, 3600.00, 3200000, 360000000000.00, 'Telecommunications'),
('ASII', 'Astra International Tbk', 5575.00, 5550.00, 5600.00, 5525.00, 1200000, 380000000000.00, 'Consumer Cyclicals'),
('UNVR', 'Unilever Indonesia Tbk', 2630.00, 2620.00, 2650.00, 2610.00, 800000, 190000000000.00, 'Consumer Defensive'),
('ICBP', 'Indofood CBP Sukses Makmur Tbk', 10900.00, 10850.00, 10950.00, 10800.00, 600000, 110000000000.00, 'Consumer Defensive'),
('INDF', 'Indofood Sukses Makmur Tbk', 6525.00, 6500.00, 6575.00, 6475.00, 900000, 57000000000.00, 'Consumer Defensive'),
('KLBF', 'Kalbe Farma Tbk', 1545.00, 1535.00, 1560.00, 1530.00, 2500000, 25000000000.00, 'Healthcare'),
('INTP', 'Indocement Tunggal Prakarsa Tbk', 10025.00, 10000.00, 10100.00, 9950.00, 400000, 110000000000.00, 'Basic Materials'),
('SMGR', 'Semen Indonesia Tbk', 5200.00, 5175.00, 5250.00, 5150.00, 700000, 62000000000.00, 'Basic Materials'),
('ANTM', 'Aneka Tambang Tbk', 1780.00, 1770.00, 1800.00, 1760.00, 1800000, 21000000000.00, 'Basic Materials'),
('PTBA', 'Bukit Asam Tbk', 3910.00, 3890.00, 3950.00, 3870.00, 1100000, 47000000000.00, 'Energy'),
('PGAS', 'Perusahaan Gas Negara Tbk', 1565.00, 1555.00, 1580.00, 1545.00, 1600000, 22000000000.00, 'Energy'),
('ADRO', 'Adaro Energy Tbk', 2630.00, 2620.00, 2650.00, 2610.00, 2200000, 54000000000.00, 'Energy'),
('GGRM', 'Gudang Garam Tbk', 25000.00, 24800.00, 25200.00, 24700.00, 150000, 120000000000.00, 'Consumer Defensive'),
('HMSP', 'HM Sampoerna Tbk', 1385.00, 1380.00, 1395.00, 1375.00, 800000, 163000000000.00, 'Consumer Defensive'),
('EXCL', 'XL Axiata Tbk', 2510.00, 2500.00, 2530.00, 2490.00, 1900000, 60000000000.00, 'Telecommunications'),
('ISAT', 'Indosat Ooredoo Hutchison Tbk', 5800.00, 5775.00, 5850.00, 5750.00, 1300000, 53000000000.00, 'Telecommunications'),
('JSMR', 'Jasa Marga Tbk', 4120.00, 4100.00, 4150.00, 4080.00, 600000, 62000000000.00, 'Industrials');

-- Insert a sample user for testing
INSERT INTO users (username, email, password_hash, balance, total_profit) VALUES
('testuser', 'test@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 100000.00, 0.00);