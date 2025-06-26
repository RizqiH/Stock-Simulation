@echo off
echo ========================================
echo 🚀 Advanced Trading Features Migration
echo ========================================

echo Checking if backend services are running...
docker-compose ps

echo.
echo Running migration script...
go run scripts/migrate_advanced_features.go

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ✅ Migration completed successfully!
    echo.
    echo ========================================
    echo 🎯 Next Steps:
    echo ========================================
    echo 1. Check database tables
    echo 2. Restart backend services
    echo 3. Test advanced features
    echo ========================================
) else (
    echo.
    echo ❌ Migration failed!
    echo Please check the error messages above.
)

pause 