#!/usr/bin/env pwsh
# Database Migration Script for Stock Simulation
# This script will reset and setup the database with default data

Write-Host "🚀 Starting Stock Simulation Database Migration..." -ForegroundColor Green

# Check if Docker is running
Write-Host "📋 Checking Docker status..." -ForegroundColor Yellow
$dockerStatus = docker info 2>$null
if (!$dockerStatus) {
    Write-Host "❌ Docker is not running. Please start Docker first." -ForegroundColor Red
    exit 1
}

# Check if we're in the correct directory
if (!(Test-Path "docker-compose.yml")) {
    Write-Host "❌ Please run this script from the stock-simulation-backend directory" -ForegroundColor Red
    exit 1
}

# Copy migration file to the correct location
Write-Host "📁 Copying migration files..." -ForegroundColor Yellow
Copy-Item "migrations/000_initial_setup.sql" "migrations/000_initial_setup.sql.bak" -ErrorAction SilentlyContinue

# Stop existing containers
Write-Host "🛑 Stopping existing containers..." -ForegroundColor Yellow
docker-compose down --volumes

# Start MySQL container
Write-Host "🐳 Starting MySQL container..." -ForegroundColor Yellow
docker-compose up -d mysql

# Wait for MySQL to be ready
Write-Host "⏳ Waiting for MySQL to be ready..." -ForegroundColor Yellow
$maxAttempts = 30
$attempt = 0
do {
    $attempt++
    Start-Sleep 2
    $mysqlReady = docker-compose exec mysql mysqladmin ping -h localhost --silent 2>$null
    if ($mysqlReady) {
        Write-Host "✅ MySQL is ready!" -ForegroundColor Green
        break
    }
    Write-Host "   Attempt $attempt/$maxAttempts..." -ForegroundColor Gray
} while ($attempt -lt $maxAttempts)

if ($attempt -ge $maxAttempts) {
    Write-Host "❌ MySQL failed to start within the timeout period" -ForegroundColor Red
    exit 1
}

# Run the migration
Write-Host "🔄 Running database migration..." -ForegroundColor Yellow
docker-compose exec mysql mysql -u stockuser -pstockpassword stock_simulation -e "SOURCE /docker-entrypoint-initdb.d/000_initial_setup.sql;"

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Database migration completed successfully!" -ForegroundColor Green
    
    # Start all services
    Write-Host "🚀 Starting all services..." -ForegroundColor Yellow
    docker-compose up -d
    
    Write-Host "📊 System Status:" -ForegroundColor Cyan
    Start-Sleep 3
    docker-compose ps
    
    Write-Host ""
    Write-Host "🎉 Stock Simulation is ready!" -ForegroundColor Green
    Write-Host "🌐 API: http://localhost:8082" -ForegroundColor Cyan
    Write-Host "🗄️  Database: localhost:3307" -ForegroundColor Cyan
    Write-Host "🔧 Adminer: http://localhost:8081 (if enabled)" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Default Login Credentials:" -ForegroundColor Yellow
    Write-Host "   Email: admin@stocksim.com" -ForegroundColor White
    Write-Host "   Email: demo@stocksim.com" -ForegroundColor White
    Write-Host "   Email: john@example.com" -ForegroundColor White
    Write-Host "   Password: password123" -ForegroundColor White
    Write-Host ""
    
} else {
    Write-Host "❌ Database migration failed!" -ForegroundColor Red
    Write-Host "📋 Checking logs..." -ForegroundColor Yellow
    docker-compose logs mysql --tail=20
    exit 1
} 