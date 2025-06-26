# Stock Simulation - Database Migration to Railway Script
# This script helps migrate local database to Railway MySQL

Write-Host "🗄️ Stock Simulation - Database Migration to Railway" -ForegroundColor Green
Write-Host "====================================================" -ForegroundColor Green

# Check if mysqldump is available
Write-Host "📋 Checking MySQL tools..." -ForegroundColor Yellow
try {
    $mysqlVersion = mysqldump --version
    Write-Host "✅ MySQL tools found" -ForegroundColor Green
} catch {
    Write-Host "❌ MySQL tools not found!" -ForegroundColor Red
    Write-Host "Please install MySQL client tools first" -ForegroundColor Yellow
    exit 1
}

# Check Railway CLI
try {
    railway --version | Out-Null
    Write-Host "✅ Railway CLI found" -ForegroundColor Green
} catch {
    Write-Host "❌ Railway CLI not found!" -ForegroundColor Red
    Write-Host "Please install: npm install -g @railway/cli" -ForegroundColor Yellow
    exit 1
}

# Get local database credentials
Write-Host "📊 Local Database Configuration" -ForegroundColor Yellow
Write-Host "================================" -ForegroundColor Yellow

$localHost = Read-Host "Local DB Host (default: localhost)"
if ([string]::IsNullOrWhiteSpace($localHost)) { $localHost = "localhost" }

$localPort = Read-Host "Local DB Port (default: 3307)"
if ([string]::IsNullOrWhiteSpace($localPort)) { $localPort = "3307" }

$localUser = Read-Host "Local DB User (default: root)"
if ([string]::IsNullOrWhiteSpace($localUser)) { $localUser = "root" }

$localPassword = Read-Host "Local DB Password (default: root)" -AsSecureString
if ($localPassword.Length -eq 0) { 
    $localPassword = ConvertTo-SecureString "root" -AsPlainText -Force 
}
$localPasswordPlain = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($localPassword))

$localDbName = Read-Host "Local DB Name (default: stock_simulation)"
if ([string]::IsNullOrWhiteSpace($localDbName)) { $localDbName = "stock_simulation" }

# Test local connection
Write-Host "🔍 Testing local database connection..." -ForegroundColor Yellow
try {
    $testResult = mysql -h $localHost -P $localPort -u $localUser -p$localPasswordPlain -e "SELECT 1;" $localDbName 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Local database connection successful" -ForegroundColor Green
    } else {
        throw "Connection failed"
    }
} catch {
    Write-Host "❌ Cannot connect to local database" -ForegroundColor Red
    Write-Host "Please check your credentials and ensure database is running" -ForegroundColor Yellow
    exit 1
}

# Export local database
$backupFile = "stocksim_backup_$(Get-Date -Format 'yyyyMMdd_HHmmss').sql"
Write-Host "📤 Exporting local database to $backupFile..." -ForegroundColor Yellow

try {
    mysqldump -h $localHost -P $localPort -u $localUser -p$localPasswordPlain --single-transaction --routines --triggers $localDbName > $backupFile
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Database exported successfully" -ForegroundColor Green
        $backupSize = (Get-Item $backupFile).Length / 1KB
        Write-Host "📁 Backup size: $([math]::Round($backupSize, 2)) KB" -ForegroundColor Cyan
    } else {
        throw "Export failed"
    }
} catch {
    Write-Host "❌ Failed to export database" -ForegroundColor Red
    exit 1
}

# Get Railway database credentials
Write-Host "🚂 Railway Database Configuration" -ForegroundColor Yellow
Write-Host "=================================" -ForegroundColor Yellow

Write-Host "Getting Railway database credentials..." -ForegroundColor Cyan
try {
    # Get Railway database URL
    $databaseUrl = railway variables get DATABASE_URL
    if ([string]::IsNullOrWhiteSpace($databaseUrl)) {
        throw "DATABASE_URL not found"
    }
    
    # Parse DATABASE_URL
    # Format: mysql://user:password@host:port/database
    if ($databaseUrl -match "mysql://([^:]+):([^@]+)@([^:]+):(\d+)/(.+)") {
        $railwayUser = $matches[1]
        $railwayPassword = $matches[2]
        $railwayHost = $matches[3]
        $railwayPort = $matches[4]
        $railwayDbName = $matches[5]
        
        Write-Host "✅ Railway credentials obtained" -ForegroundColor Green
        Write-Host "🔹 Host: $railwayHost" -ForegroundColor Cyan
        Write-Host "🔹 Port: $railwayPort" -ForegroundColor Cyan
        Write-Host "🔹 Database: $railwayDbName" -ForegroundColor Cyan
    } else {
        throw "Could not parse DATABASE_URL"
    }
} catch {
    Write-Host "❌ Could not get Railway database credentials" -ForegroundColor Red
    Write-Host "Please ensure you're linked to the correct Railway project" -ForegroundColor Yellow
    Write-Host "Run: railway link" -ForegroundColor Cyan
    exit 1
}

# Test Railway connection
Write-Host "🔍 Testing Railway database connection..." -ForegroundColor Yellow
try {
    $testResult = mysql -h $railwayHost -P $railwayPort -u $railwayUser -p$railwayPassword -e "SELECT 1;" $railwayDbName 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Railway database connection successful" -ForegroundColor Green
    } else {
        throw "Connection failed: $testResult"
    }
} catch {
    Write-Host "❌ Cannot connect to Railway database" -ForegroundColor Red
    Write-Host "Error: $_" -ForegroundColor Yellow
    exit 1
}

# Import to Railway
Write-Host "📥 Importing database to Railway..." -ForegroundColor Yellow
$proceed = Read-Host "⚠️  This will overwrite the Railway database. Continue? (y/n)"

if ($proceed -eq "y" -or $proceed -eq "Y") {
    try {
        mysql -h $railwayHost -P $railwayPort -u $railwayUser -p$railwayPassword $railwayDbName < $backupFile
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Database imported successfully to Railway" -ForegroundColor Green
        } else {
            throw "Import failed"
        }
    } catch {
        Write-Host "❌ Failed to import database to Railway" -ForegroundColor Red
        Write-Host "Error: $_" -ForegroundColor Yellow
        exit 1
    }
} else {
    Write-Host "ℹ️ Import cancelled" -ForegroundColor Yellow
    exit 0
}

# Verify import
Write-Host "🔍 Verifying import..." -ForegroundColor Yellow
try {
    # Check if main tables exist
    $tables = mysql -h $railwayHost -P $railwayPort -u $railwayUser -p$railwayPassword -e "SHOW TABLES;" $railwayDbName 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Import verification successful" -ForegroundColor Green
        
        # Count records in main tables
        $userCount = mysql -h $railwayHost -P $railwayPort -u $railwayUser -p$railwayPassword -e "SELECT COUNT(*) FROM users;" $railwayDbName 2>&1 | Select-String -Pattern '\d+'
        $stockCount = mysql -h $railwayHost -P $railwayPort -u $railwayUser -p$railwayPassword -e "SELECT COUNT(*) FROM stocks;" $railwayDbName 2>&1 | Select-String -Pattern '\d+'
        
        Write-Host "📊 Database Statistics:" -ForegroundColor Cyan
        Write-Host "🔹 Users: $userCount" -ForegroundColor Cyan
        Write-Host "🔹 Stocks: $stockCount" -ForegroundColor Cyan
    } else {
        throw "Verification failed"
    }
} catch {
    Write-Host "⚠️ Could not verify import, but it may have succeeded" -ForegroundColor Yellow
}

# Cleanup
$cleanup = Read-Host "🧹 Delete backup file? (y/n)"
if ($cleanup -eq "y" -or $cleanup -eq "Y") {
    Remove-Item $backupFile
    Write-Host "✅ Backup file deleted" -ForegroundColor Green
} else {
    Write-Host "📁 Backup file kept: $backupFile" -ForegroundColor Cyan
}

# Final summary
Write-Host ""
Write-Host "🎉 Database Migration Complete!" -ForegroundColor Green
Write-Host "===============================" -ForegroundColor Green
Write-Host ""
Write-Host "📋 Next Steps:" -ForegroundColor Yellow
Write-Host "1. Test your Railway application endpoints"
Write-Host "2. Check that all data is accessible"
Write-Host "3. Update your frontend to use Railway API URL"
Write-Host "4. Run some API tests to verify functionality"
Write-Host ""
Write-Host "🔧 Useful Commands:" -ForegroundColor Yellow
Write-Host "- railway logs                     # Check application logs"
Write-Host "- railway connect MySQL            # Direct MySQL connection"
Write-Host "- railway variables get DATABASE_URL  # Get database URL"
Write-Host ""
Write-Host "🚀 Your database is now on Railway!" -ForegroundColor Green 