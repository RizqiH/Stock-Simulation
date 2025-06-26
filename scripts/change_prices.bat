@echo off
echo.
echo ===========================================
echo     STOCK PRICE CHANGER FOR TESTING
echo ===========================================
echo.
echo Choose an option:
echo 1. Lower AAPL to $148.00 (trigger BUY LIMIT orders)
echo 2. Raise AAPL to $153.00 (trigger SELL LIMIT orders)
echo 3. Random price movements (all stocks)
echo 4. Custom price for AAPL
echo 5. Show current prices
echo 6. Exit
echo.
set /p choice="Enter your choice (1-6): "

if "%choice%"=="1" goto lower_aapl
if "%choice%"=="2" goto raise_aapl
if "%choice%"=="3" goto random_prices
if "%choice%"=="4" goto custom_price
if "%choice%"=="5" goto show_prices
if "%choice%"=="6" goto exit
goto menu

:lower_aapl
echo.
echo Lowering AAPL to $148.00...
docker-compose exec mysql mysql -u stockuser -pstockpassword stock_simulation -e "UPDATE stocks SET current_price = 148.00, updated_at = NOW() WHERE symbol = 'AAPL';"
echo ✅ AAPL price updated to $148.00
goto show_aapl

:raise_aapl
echo.
echo Raising AAPL to $153.00...
docker-compose exec mysql mysql -u stockuser -pstockpassword stock_simulation -e "UPDATE stocks SET current_price = 153.00, updated_at = NOW() WHERE symbol = 'AAPL';"
echo ✅ AAPL price updated to $153.00
goto show_aapl

:random_prices
echo.
echo Applying random price movements...
docker-compose exec mysql mysql -u stockuser -pstockpassword stock_simulation -e "UPDATE stocks SET current_price = current_price * (1 + (RAND() - 0.5) * 0.06), updated_at = NOW();"
echo ✅ All stock prices randomly updated
goto show_prices

:custom_price
echo.
set /p custom="Enter new AAPL price: $"
echo Updating AAPL to $%custom%...
docker-compose exec mysql mysql -u stockuser -pstockpassword stock_simulation -e "UPDATE stocks SET current_price = %custom%, updated_at = NOW() WHERE symbol = 'AAPL';"
echo ✅ AAPL price updated to $%custom%
goto show_aapl

:show_aapl
echo.
echo Current AAPL price:
docker-compose exec mysql mysql -u stockuser -pstockpassword stock_simulation -e "SELECT symbol, current_price, updated_at FROM stocks WHERE symbol = 'AAPL';"
echo.
pause
goto menu

:show_prices
echo.
echo Current stock prices:
docker-compose exec mysql mysql -u stockuser -pstockpassword stock_simulation -e "SELECT symbol, ROUND(current_price, 2) as price, updated_at FROM stocks ORDER BY symbol;"
echo.
pause
goto menu

:exit
echo.
echo Goodbye! 👋
exit 