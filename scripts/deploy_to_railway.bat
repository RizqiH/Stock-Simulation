@echo off
echo ================================================================
echo             STOCKSIM PRO - RAILWAY DEPLOYMENT
echo ================================================================
echo.
echo This script will help you deploy your database to Railway
echo.
echo Prerequisites:
echo - Railway account (https://railway.app)
echo - MySQL client installed
echo - Railway CLI installed (npm install -g @railway/cli)
echo.
pause

echo.
echo 🚀 Starting Railway deployment process...
echo.

cd /d "%~dp0\.."

echo 📊 Step 1: Creating database backup
echo ----------------------------------------
echo Creating backup of local database...
mysqldump -u root -proot -h localhost -P 3307 stock_simulation > stocksim_railway_backup.sql

if %ERRORLEVEL% EQU 0 (
    echo ✅ Database backup created: stocksim_railway_backup.sql
) else (
    echo ❌ Failed to create database backup
    echo Please check if MySQL is running on localhost:3307
    pause
    exit /b 1
)

echo.
echo 🔧 Step 2: Railway setup instructions
echo ----------------------------------------
echo.
echo Please complete these steps manually:
echo.
echo 1. Go to https://railway.app and sign in
echo 2. Create new project: "stocksim-database"
echo 3. Add MySQL service
echo 4. Copy the database credentials
echo.
echo After completing the above steps, please enter your Railway MySQL credentials:
echo.
echo IMPORTANT: Your Railway MySQL credentials should look like this:
echo   Host: containers-us-west-xxx.railway.app (NOT mysql.railway.internal)
echo   Port: 6543 (or another 4-digit port, NOT 3306)
echo   User: root
echo   Password: [long random string]
echo   Database: railway
echo.
echo You can find these in Railway Dashboard ^> MySQL Service ^> Connect ^> Public Networking
echo.

set /p RAILWAY_HOST="Enter Railway MySQL Host (format: containers-us-west-xxx.railway.app): "
set /p RAILWAY_PORT="Enter Railway MySQL Port (NOT 3306, usually 4 digits): "
set /p RAILWAY_USER="Enter Railway MySQL User (usually 'root'): "
set /p RAILWAY_PASSWORD="Enter Railway MySQL Password: "
set /p RAILWAY_DATABASE="Enter Railway Database Name (usually 'railway'): "

echo.
echo 📤 Step 3: Uploading database to Railway
echo ----------------------------------------
echo Connecting to Railway MySQL and importing data...

mysql -h %RAILWAY_HOST% -P %RAILWAY_PORT% -u %RAILWAY_USER% -p%RAILWAY_PASSWORD% %RAILWAY_DATABASE% < stocksim_railway_backup.sql

if %ERRORLEVEL% EQU 0 (
    echo ✅ Database successfully uploaded to Railway!
) else (
    echo ❌ Failed to upload database to Railway
    echo Please check your credentials and try again
    pause
    exit /b 1
)

echo.
echo 🔧 Step 4: Creating production configuration
echo ----------------------------------------

echo Creating .env.production file...
(
echo # Railway MySQL Configuration
echo DB_HOST=%RAILWAY_HOST%
echo DB_PORT=%RAILWAY_PORT%
echo DB_USER=%RAILWAY_USER%
echo DB_PASSWORD=%RAILWAY_PASSWORD%
echo DB_NAME=%RAILWAY_DATABASE%
echo DATABASE_URL=mysql://%RAILWAY_USER%:%RAILWAY_PASSWORD%@%RAILWAY_HOST%:%RAILWAY_PORT%/%RAILWAY_DATABASE%
echo.
echo # Application Settings
echo PORT=8080
echo JWT_SECRET=your-production-jwt-secret-change-this
echo ENV=production
echo.
echo # CORS Settings
echo CORS_ORIGINS=http://localhost:3000,https://your-frontend-domain.com
) > .env.production

echo ✅ Created .env.production file

echo.
echo 🧪 Step 5: Testing connection
echo ----------------------------------------
echo Testing connection to Railway database...

mysql -h %RAILWAY_HOST% -P %RAILWAY_PORT% -u %RAILWAY_USER% -p%RAILWAY_PASSWORD% %RAILWAY_DATABASE% -e "SHOW TABLES;"

if %ERRORLEVEL% EQU 0 (
    echo ✅ Connection test successful!
    echo.
    echo 📋 Verifying data...
    mysql -h %RAILWAY_HOST% -P %RAILWAY_PORT% -u %RAILWAY_USER% -p%RAILWAY_PASSWORD% %RAILWAY_DATABASE% -e "SELECT 'Users:' as Table_Name, COUNT(*) as Count FROM users UNION SELECT 'Stocks:', COUNT(*) FROM stocks UNION SELECT 'Transactions:', COUNT(*) FROM transactions;"
) else (
    echo ❌ Connection test failed
)

echo.
echo 📝 Step 6: Next steps
echo ----------------------------------------
echo.
echo ✅ Database deployed to Railway successfully!
echo.
echo Next steps to complete deployment:
echo.
echo 1. Update your Go application to use the new database:
echo    - Copy the .env.production file values
echo    - Update your config.go to read environment variables
echo.
echo 2. For local development with Railway database:
echo    - Copy .env.production to .env
echo    - Restart your Go application
echo.
echo 3. Deploy your backend to Railway:
echo    - Push your code to GitHub
echo    - Create a new service in Railway
echo    - Connect your GitHub repository
echo    - Set environment variables
echo.
echo 4. Update frontend configuration:
echo    - Update API base URL in nuxt.config.ts
echo    - Deploy frontend to Vercel/Netlify
echo.
echo 📊 Database Information:
echo Host: %RAILWAY_HOST%
echo Port: %RAILWAY_PORT%
echo Database: %RAILWAY_DATABASE%
echo User: %RAILWAY_USER%
echo.
echo 💰 Cost: ~$5/month for MySQL service
echo 🔒 Security: SSL enabled by default
echo 📈 Scaling: Automatic scaling available
echo.

echo ================================================================
echo              RAILWAY DEPLOYMENT COMPLETED
echo ================================================================
echo.
pause 