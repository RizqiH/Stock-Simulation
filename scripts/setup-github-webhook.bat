@echo off
REM Setup GitHub Webhook for Jenkins CI/CD
REM This script helps configure GitHub webhook for automatic builds

echo ========================================
echo GitHub Webhook Setup for Jenkins
echo ========================================
echo.

REM Check if required tools are available
where curl >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: curl is not installed or not in PATH!
    echo Please install curl or use Git Bash.
    pause
    exit /b 1
)

echo ✓ curl is available
echo.

REM Get user input
set /p GITHUB_TOKEN="Enter your GitHub Personal Access Token: "
set /p GITHUB_REPO="Enter your GitHub repository (format: username/repo-name): "
set /p JENKINS_URL="Enter your Jenkins URL (default: http://localhost:8080): "

REM Set default Jenkins URL if not provided
if "%JENKINS_URL%"=="" set JENKINS_URL=http://localhost:8080

echo.
echo Configuration:
echo - GitHub Repository: %GITHUB_REPO%
echo - Jenkins URL: %JENKINS_URL%
echo - Webhook URL: %JENKINS_URL%/github-webhook/
echo.

set /p CONFIRM="Continue with webhook setup? (y/n): "
if /i not "%CONFIRM%"=="y" (
    echo Setup cancelled.
    pause
    exit /b 0
)

echo.
echo Setting up GitHub webhook...

REM Create webhook payload
(
echo {
echo   "name": "web",
echo   "active": true,
echo   "events": [
echo     "push",
echo     "pull_request"
echo   ],
echo   "config": {
echo     "url": "%JENKINS_URL%/github-webhook/",
echo     "content_type": "json",
echo     "insecure_ssl": "0"
echo   }
echo }
) > webhook_payload.json

REM Create the webhook
curl -X POST ^
  -H "Authorization: token %GITHUB_TOKEN%" ^
  -H "Accept: application/vnd.github.v3+json" ^
  -H "Content-Type: application/json" ^
  -d @webhook_payload.json ^
  "https://api.github.com/repos/%GITHUB_REPO%/hooks"

if %errorlevel% equ 0 (
    echo.
    echo ✓ GitHub webhook created successfully!
) else (
    echo.
    echo ✗ Failed to create GitHub webhook.
    echo Please check your token permissions and repository access.
)

REM Clean up
del webhook_payload.json >nul 2>&1

echo.
echo ========================================
echo Next Steps:
echo ========================================
echo 1. Verify webhook in GitHub repository settings
echo 2. Go to: https://github.com/%GITHUB_REPO%/settings/hooks
echo 3. Check that webhook URL is: %JENKINS_URL%/github-webhook/
echo 4. Test webhook by making a commit to your repository
echo 5. Monitor Jenkins for automatic build triggers
echo.
echo ========================================
echo Jenkins Pipeline Setup:
echo ========================================
echo 1. Create new Pipeline job in Jenkins
echo 2. Configure Pipeline script from SCM
echo 3. Set Repository URL: https://github.com/%GITHUB_REPO%.git
echo 4. Set Script Path: Jenkinsfile
echo 5. Enable "GitHub hook trigger for GITScm polling"
echo.
echo ========================================
echo Troubleshooting:
echo ========================================
echo - Ensure Jenkins is accessible from GitHub
echo - Check firewall settings for port 8080
echo - Verify GitHub token has repo and admin:repo_hook permissions
echo - Check Jenkins logs: docker logs jenkins-master
echo.

pause