@echo off
echo Starting Stock Simulation Production Environment...
echo.

:: Check if Docker is running
docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Docker is not running or not installed
    echo Please start Docker and try again
    pause
    exit /b 1
)

:: Check if .env.prod exists
if not exist ".env.prod" (
    echo Error: .env.prod file not found
    echo Please copy .env.example to .env.prod and configure it with production values
    pause
    exit /b 1
)

:: Validate production environment
echo Validating production configuration...
findstr "CHANGE_THIS" .env.prod >nul
if %errorlevel% equ 0 (
    echo WARNING: Found default values in .env.prod
    echo Please update all CHANGE_THIS values before deploying to production
    echo.
    set /p continue=Continue anyway? (y/N): 
    if /i not "%continue%"=="y" (
        echo Deployment cancelled
        pause
        exit /b 1
    )
)

:: Copy production environment file
copy .env.prod .env >nul 2>&1

echo Building production images...
echo.

:: Build production images
docker-compose -f docker-compose.prod.yml build --no-cache

if %errorlevel% neq 0 (
    echo Error: Failed to build production images
    pause
    exit /b 1
)

echo.
echo Starting production services...
echo.

:: Start production services
docker-compose -f docker-compose.prod.yml up -d

if %errorlevel% neq 0 (
    echo Error: Failed to start production environment
    echo Check the logs with: docker-compose -f docker-compose.prod.yml logs
    pause
    exit /b 1
)

:: Wait for services to be ready
echo Waiting for services to be ready...
timeout /t 30 /nobreak >nul

:: Check service health
echo Checking service health...
docker-compose -f docker-compose.prod.yml ps

echo.
echo Production environment started successfully!
echo.
echo Services:
echo - API Server: http://localhost:8080
echo - Load Balancer: http://localhost:80
echo - Database: localhost:3306
echo - Redis: localhost:6379
echo.
echo Production commands:
echo - View logs: docker-compose -f docker-compose.prod.yml logs -f
echo - Stop services: docker-compose -f docker-compose.prod.yml down
echo - Scale API: docker-compose -f docker-compose.prod.yml up -d --scale api=3
echo - Update service: docker-compose -f docker-compose.prod.yml up -d --no-deps api
echo.
echo IMPORTANT: Monitor logs and system resources in production!
echo.
pause