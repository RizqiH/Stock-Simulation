@echo off
echo ========================================
echo 🚀 Stock Trading Simulation - Complete System
echo ========================================

echo.
echo 📋 STEP 1: Starting Backend Services...
docker-compose down
docker-compose up -d

echo.
echo 📋 STEP 2: Waiting for services to be ready...
timeout /t 10 /nobreak

echo.
echo 📋 STEP 3: Running Advanced Features Migration...
go run scripts/migrate_advanced_features.go

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ✅ Backend setup completed!
    
    echo.
    echo 📋 STEP 4: Starting Frontend...
    echo Opening new terminal for frontend...
    
    start cmd /k "cd /d ..\stock-simulation-frontend && echo ========================================= && echo 🎨 Frontend Development Server && echo ========================================= && echo Installing dependencies... && npm install && echo. && echo Starting development server... && npm run dev"
    
    echo.
    echo ========================================
    echo 🎯 SYSTEM READY!
    echo ========================================
    echo Backend API: http://localhost:8082
    echo Frontend App: http://localhost:3000
    echo ========================================
    echo.
    echo 📊 Available Features:
    echo - Advanced Orders (Market, Limit, Stop Loss, etc.)
    echo - Real-time WebSocket updates
    echo - Commission & Fees system
    echo - Market hours simulation
    echo - Order analytics & metrics
    echo ========================================
    
) else (
    echo.
    echo ❌ Backend setup failed!
    echo Please check the error messages above.
)

pause 