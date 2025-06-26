@echo off
echo ================================================================
echo             STOCKSIM PRO - HISTORICAL DATA SETUP
echo ================================================================
echo.
echo This script will add historical data for all stock symbols
echo that currently don't have any historical data.
echo.
echo Make sure your MySQL server is running on localhost:3307
echo.
pause

echo.
echo 🔍 Starting historical data population...
echo.

cd /d "%~dp0"

echo 📊 Method 1: Running Go script (Recommended)
echo ----------------------------------------
go run populate_all_historical_data.go

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ✅ Historical data setup completed successfully!
    echo.
    echo 📋 What was added:
    echo   - 30 days of OHLCV data for each missing symbol
    echo   - Realistic price movements and volumes
    echo   - Data compatible with chart displays
    echo.
    echo 🎯 You can now:
    echo   - View stock charts for all symbols
    echo   - See portfolio performance charts
    echo   - All technical analysis features will work
    echo.
) else (
    echo.
    echo ❌ Go script failed. Trying SQL method...
    echo.
    echo 📊 Method 2: Running SQL script (Fallback)
    echo ----------------------------------------
    mysql -u root -proot -h localhost -P 3307 stock_simulation < complete_historical_data.sql
    
    if %ERRORLEVEL% EQU 0 (
        echo ✅ SQL script completed successfully!
    ) else (
        echo ❌ Both methods failed. Please check:
        echo   1. MySQL server is running on localhost:3307
        echo   2. Database 'stock_simulation' exists
        echo   3. User 'root' with password 'root' has access
    )
)

echo.
echo ================================================================
echo              HISTORICAL DATA SETUP COMPLETED
echo ================================================================
echo.
pause 