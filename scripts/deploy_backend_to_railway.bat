@echo off
echo ================================================================
echo          STOCKSIM PRO - BACKEND RAILWAY DEPLOYMENT
echo ================================================================
echo.
echo This script will help you deploy your Go backend to Railway
echo.
echo Prerequisites:
echo - Railway CLI installed (npm install -g @railway/cli)
echo - Git repository
echo - Railway database already deployed
echo.
pause

echo.
echo 🚀 Starting backend deployment to Railway...
echo.

cd /d "%~dp0\.."

echo 📋 Step 1: Checking prerequisites
echo ----------------------------------------
echo Checking Railway CLI...
railway --version
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Railway CLI not found!
    echo Please install: npm install -g @railway/cli
    pause
    exit /b 1
)

echo Checking Git...
git --version
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Git not found!
    echo Please install Git first
    pause
    exit /b 1
)

echo ✅ All prerequisites met!

echo.
echo 🔧 Step 2: Railway setup
echo ----------------------------------------
echo.
echo Please complete these steps:
echo 1. Login to Railway: railway login
echo 2. Link to your existing project (where MySQL is): railway link
echo 3. Or create new project if needed: railway new
echo.
echo Do you want to login now? (y/n)
set /p LOGIN_CHOICE="Enter your choice: "

if /i "%LOGIN_CHOICE%"=="y" (
    echo Logging in to Railway...
    railway login
    if %ERRORLEVEL% NEQ 0 (
        echo ❌ Railway login failed
        pause
        exit /b 1
    )
)

echo.
echo Do you want to link to existing project? (y/n)
set /p LINK_CHOICE="Enter your choice: "

if /i "%LINK_CHOICE%"=="y" (
    echo Available projects:
    railway projects
    echo.
    echo Linking to project...
    railway link
    if %ERRORLEVEL% NEQ 0 (
        echo ❌ Railway link failed
        pause
        exit /b 1
    )
)

echo.
echo 📦 Step 3: Preparing for deployment
echo ----------------------------------------
echo.
echo Creating .railwayignore file...
(
echo node_modules/
echo .git/
echo .env*
echo *.log
echo stocksim_railway_backup.sql
echo scripts/
echo docs/
echo jenkins/
echo jenkins_home/
echo .github/
echo *.md
echo .gitignore
echo .dockerignore
echo docker-compose*.yml
echo Makefile
echo Jenkinsfile
) > .railwayignore

echo ✅ .railwayignore created

echo.
echo 🌐 Step 4: Environment variables setup
echo ----------------------------------------
echo.
echo Setting up environment variables...
echo.
echo Please enter your production configuration:
echo.

set /p JWT_SECRET="Enter JWT Secret (or press Enter to use generated one): "
if "%JWT_SECRET%"=="" (
    set "JWT_SECRET=cc5707c0eafe48f74323671eecde11f664cc2362c79488aa8c3ac3658fb35dc4692e432ccd68425f0b7ce1e8785fa9976f83e3839166659d1e11d62d6ef16fbfd"
)

set /p FRONTEND_URL="Enter your frontend URL (e.g., https://your-app.vercel.app): "
if "%FRONTEND_URL%"=="" (
    set "FRONTEND_URL=http://localhost:3000"
)

echo.
echo Setting Railway environment variables...

railway variables set PORT=8080
railway variables set ENV=production
railway variables set JWT_SECRET=%JWT_SECRET%
railway variables set CORS_ORIGINS="%FRONTEND_URL%,http://localhost:3000"

echo Setting database variables (linking to existing MySQL service)...
railway variables set DB_HOST=${{MySQL.MYSQL_HOST}}
railway variables set DB_PORT=${{MySQL.MYSQL_PORT}}
railway variables set DB_USER=${{MySQL.MYSQL_USER}}
railway variables set DB_PASSWORD=${{MySQL.MYSQL_PASSWORD}}
railway variables set DB_NAME=${{MySQL.MYSQL_DATABASE}}

echo ✅ Environment variables set!

echo.
echo 🚀 Step 5: Deploying to Railway
echo ----------------------------------------
echo.
echo Starting deployment...
railway up

if %ERRORLEVEL% EQU 0 (
    echo ✅ Backend deployed successfully!
    echo.
    echo 🌐 Getting deployment URL...
    railway domain
    echo.
    echo 📋 Step 6: Next steps
    echo ----------------------------------------
    echo.
    echo ✅ Backend deployed to Railway successfully!
    echo.
    echo Next steps:
    echo 1. Note your backend URL from above
    echo 2. Update frontend configuration with this URL
    echo 3. Test your API endpoints
    echo 4. Deploy frontend to Vercel/Netlify
    echo.
    echo 🧪 Testing your deployment:
    echo 1. Health check: curl https://your-backend-url.railway.app/health
    echo 2. API check: curl https://your-backend-url.railway.app/api/v1/stocks
    echo.
    echo 📊 Monitoring:
    echo 1. View logs: railway logs
    echo 2. Monitor metrics in Railway dashboard
    echo 3. Check WebSocket: ws://your-backend-url.railway.app/api/v1/ws
) else (
    echo ❌ Deployment failed!
    echo Check the error messages above and try again
    echo.
    echo Common issues:
    echo 1. Make sure you're linked to the correct project
    echo 2. Check that your MySQL service is running
    echo 3. Verify environment variables are set correctly
    echo.
    echo Troubleshooting commands:
    echo - railway logs
    echo - railway variables
    echo - railway status
)

echo.
echo ================================================================
echo            BACKEND RAILWAY DEPLOYMENT COMPLETED
echo ================================================================
echo.
pause 