# Stock Simulation Backend - Railway Deployment Script
# This script helps deploy the backend to Railway with proper configuration

Write-Host "🚀 Stock Simulation - Railway Deployment Script" -ForegroundColor Green
Write-Host "=================================================" -ForegroundColor Green

# Check if Railway CLI is installed
Write-Host "📋 Checking Railway CLI installation..." -ForegroundColor Yellow
try {
    $railwayVersion = railway --version
    Write-Host "✅ Railway CLI found: $railwayVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Railway CLI not found!" -ForegroundColor Red
    Write-Host "Please install Railway CLI first:" -ForegroundColor Yellow
    Write-Host "npm install -g @railway/cli" -ForegroundColor Cyan
    exit 1
}

# Login check
Write-Host "🔐 Checking Railway authentication..." -ForegroundColor Yellow
try {
    railway whoami | Out-Null
    Write-Host "✅ Already logged in to Railway" -ForegroundColor Green
} catch {
    Write-Host "⚠️ Not logged in to Railway" -ForegroundColor Yellow
    Write-Host "Please login first: railway login" -ForegroundColor Cyan
    $login = Read-Host "Do you want to login now? (y/n)"
    if ($login -eq "y" -or $login -eq "Y") {
        railway login
    } else {
        exit 1
    }
}

# Project selection
Write-Host "📁 Railway Project Setup..." -ForegroundColor Yellow
$createNew = Read-Host "Do you want to create a new Railway project? (y/n)"

if ($createNew -eq "y" -or $createNew -eq "Y") {
    $projectName = Read-Host "Enter project name (e.g., stocksim-backend)"
    if ([string]::IsNullOrWhiteSpace($projectName)) {
        $projectName = "stocksim-backend"
    }
    
    Write-Host "Creating new Railway project: $projectName" -ForegroundColor Cyan
    railway project new $projectName
    railway link
} else {
    Write-Host "Linking to existing project..." -ForegroundColor Cyan
    railway link
}

# Add MySQL service
Write-Host "🗄️ Database Setup..." -ForegroundColor Yellow
$addDb = Read-Host "Do you need to add MySQL database service? (y/n)"
if ($addDb -eq "y" -or $addDb -eq "Y") {
    Write-Host "Adding MySQL service..." -ForegroundColor Cyan
    railway add --database mysql
    Write-Host "✅ MySQL service added" -ForegroundColor Green
    Write-Host "⏳ Waiting for MySQL to initialize (30 seconds)..." -ForegroundColor Yellow
    Start-Sleep 30
}

# Set environment variables
Write-Host "⚙️ Setting Environment Variables..." -ForegroundColor Yellow

# JWT Secret
$jwtSecret = Read-Host "Enter JWT secret for production (leave empty for auto-generated)"
if ([string]::IsNullOrWhiteSpace($jwtSecret)) {
    # Generate a random JWT secret
    $jwtSecret = [System.Web.Security.Membership]::GeneratePassword(64, 0)
}
railway variables set JWT_SECRET="$jwtSecret"

# CORS Origins
$corsOrigins = Read-Host "Enter CORS origins (e.g., https://your-frontend.railway.app)"
if (![string]::IsNullOrWhiteSpace($corsOrigins)) {
    railway variables set CORS_ORIGINS="$corsOrigins"
}

# Set other production variables
railway variables set ENV="production"
railway variables set PORT="8080"
railway variables set HOST="0.0.0.0"

Write-Host "✅ Environment variables set" -ForegroundColor Green

# Deploy the application
Write-Host "🚀 Deploying Application..." -ForegroundColor Yellow
railway deploy

# Show deployment info
Write-Host "📊 Deployment Information:" -ForegroundColor Green
Write-Host "=========================" -ForegroundColor Green

# Get the domain
try {
    $domain = railway domain
    Write-Host "🌐 Application URL: https://$domain" -ForegroundColor Cyan
    Write-Host "🏥 Health Check: https://$domain/health" -ForegroundColor Cyan
    Write-Host "📡 API Base: https://$domain/api/v1" -ForegroundColor Cyan
    Write-Host "🔌 WebSocket: wss://$domain/api/v1/ws" -ForegroundColor Cyan
} catch {
    Write-Host "⚠️ Could not retrieve domain. Check Railway dashboard." -ForegroundColor Yellow
}

# Show next steps
Write-Host ""
Write-Host "🎉 Deployment Complete!" -ForegroundColor Green
Write-Host "======================" -ForegroundColor Green
Write-Host ""
Write-Host "📋 Next Steps:" -ForegroundColor Yellow
Write-Host "1. Check Railway dashboard for deployment status"
Write-Host "2. Test the health endpoint"
Write-Host "3. Import your database data (if needed)"
Write-Host "4. Update your frontend API URL"
Write-Host "5. Test all API endpoints"
Write-Host ""
Write-Host "🔧 Useful Railway Commands:" -ForegroundColor Yellow
Write-Host "- railway logs      # View application logs"
Write-Host "- railway status    # Check service status"
Write-Host "- railway variables # View environment variables"
Write-Host "- railway domain    # Get application URL"
Write-Host "- railway connect   # Connect to services"
Write-Host ""
Write-Host "📖 For database migration, see:" -ForegroundColor Yellow
Write-Host "docs/RAILWAY_DEPLOYMENT.md"

Write-Host ""
Write-Host "🚀 Happy deploying!" -ForegroundColor Green 