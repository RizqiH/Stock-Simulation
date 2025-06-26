CREATE DATABASE IF NOT EXISTS stock_simulation;
CREATE USER IF NOT EXISTS 'stockuser'@'%' IDENTIFIED BY 'stockpassword';
GRANT ALL PRIVILEGES ON stock_simulation.* TO 'stockuser'@'%';
FLUSH PRIVILEGES; 