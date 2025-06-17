@echo off
echo Starting Stock Simulation Development Environment...
echo.

:: Check if Docker is running
docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Docker is not running or not installed
    echo Please start Docker Desktop and try again
    pause
    exit /b 1
)

:: Check if .env.dev exists
if not exist ".env.dev" (
    echo Error: .env.dev file not found
    echo Please copy .env.example to .env.dev and configure it
    pause
    exit /b 1
)

:: Copy development environment file
copy .env.dev .env >nul 2>&1

echo Building and starting development containers...
echo.

:: Build and start services
docker-compose --env-file .env.dev up --build -d

if %errorlevel% neq 0 (
    echo Error: Failed to start development environment
    echo Check the logs with: docker-compose logs
    pause
    exit /b 1
)

echo.
echo Development environment started successfully!
echo.
echo Services:
echo - API Server: http://localhost:8080
echo - Database: localhost:3306 (user: stockuser, password: stockpassword)
echo - Redis: localhost:6379
echo - Adminer: http://localhost:8081 (optional database management)
echo.
echo Useful commands:
echo - View logs: docker-compose logs -f
echo - Stop services: docker-compose down
echo - Restart API: docker-compose restart api
echo - Access API container: docker-compose exec api sh
echo - Access MySQL: docker-compose exec mysql mysql -u stockuser -p stock_simulation
echo.
echo Press any key to view logs...
pause >nul
docker-compose logs -f