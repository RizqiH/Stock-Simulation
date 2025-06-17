@echo off
echo Setting up Stock Simulation Database...
echo.

:: Check if MySQL is installed and accessible
mysql --version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: MySQL is not installed or not in PATH
    echo Please install MySQL and add it to your PATH environment variable
    pause
    exit /b 1
)

:: Prompt for MySQL root password
set /p MYSQL_PASSWORD=Enter MySQL root password: 

echo.
echo Creating database and tables...
mysql -u root -p%MYSQL_PASSWORD% < ../migrations/001_create_database.sql
if %errorlevel% neq 0 (
    echo Error: Failed to create database and tables
    pause
    exit /b 1
)

echo.
echo Inserting sample data...
mysql -u root -p%MYSQL_PASSWORD% < ../migrations/002_insert_sample_data.sql
if %errorlevel% neq 0 (
    echo Error: Failed to insert sample data
    pause
    exit /b 1
)

echo.
echo Database setup completed successfully!
echo.
echo Database: stock_simulation
echo Tables created:
echo - users
echo - stocks
echo - transactions
echo - portfolios
echo.
echo Sample data inserted:
echo - 20 Indonesian stocks
echo - 1 test user (username: testuser, email: test@example.com)
echo.
echo You can now run the application with: go run ./cmd/api
pause