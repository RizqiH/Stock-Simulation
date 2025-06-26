@echo off
echo ================================================================
echo          RAILWAY CREDENTIALS HELPER
echo ================================================================
echo.
echo This script will help you get the CORRECT credentials from Railway
echo.
echo IMPORTANT: Railway has TWO types of credentials:
echo   1. INTERNAL (mysql.railway.internal) - Only works INSIDE Railway
echo   2. PUBLIC (containers-us-west-xxx.railway.app) - For EXTERNAL connections
echo.
echo For deploying from your local computer, you need PUBLIC credentials.
echo.
echo ================================================================
echo                     STEP-BY-STEP GUIDE
echo ================================================================
echo.
echo 1. Go to your Railway dashboard
echo 2. Click on your MySQL service
echo 3. Click on "Connect" tab (NOT "Variables" tab)
echo 4. Look for "Public Networking" section
echo 5. You will see something like:
echo.
echo    Host: containers-us-west-123.railway.app
echo    Port: 6543 (or another 4-digit number)
echo    Username: root
echo    Password: [your password from Variables tab]
echo    Database: railway
echo.
echo ================================================================
echo                        EXAMPLE
echo ================================================================
echo.
echo ✅ CORRECT (Public - for external connections):
echo    Host: containers-us-west-123.railway.app
echo    Port: 6543
echo    User: root
echo    Password: GqmhoSbTUbGSeBgTQXwXICUDXKTGITbX
echo    Database: railway
echo.
echo ❌ WRONG (Internal - only works inside Railway):
echo    Host: mysql.railway.internal
echo    Port: 3306
echo.
echo ================================================================
echo                    ALTERNATIVE METHOD
echo ================================================================
echo.
echo If you can't find "Connect" tab, you can extract from MYSQL_PUBLIC_URL:
echo.
echo From Variables tab, look for: MYSQL_PUBLIC_URL
echo Format: mysql://root:PASSWORD@HOST:PORT/DATABASE
echo.
echo Example:
echo MYSQL_PUBLIC_URL = mysql://root:GqmhoSbTUbGSeBgTQXwXICUDXKTGITbX@containers-us-west-123.railway.app:6543/railway
echo                          ↑                                    ↑                                   ↑      ↑
echo                     Password                               Host                                Port  Database
echo.
echo ================================================================
echo.
echo Now run the deployment script again with the CORRECT public credentials!
echo.
pause 