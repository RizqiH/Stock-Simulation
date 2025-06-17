@echo off
REM Setup Jenkins for Stock Simulation Backend CI/CD
REM This script sets up Jenkins with Docker support

echo ========================================
echo Setting up Jenkins for CI/CD Pipeline
echo ========================================
echo.

REM Check if Docker is running
docker version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Docker is not running or not installed!
    echo Please install Docker Desktop and make sure it's running.
    pause
    exit /b 1
)

echo ✓ Docker is running
echo.

REM Create Jenkins directory
if not exist "jenkins_home" (
    mkdir jenkins_home
    echo ✓ Created jenkins_home directory
) else (
    echo ✓ jenkins_home directory already exists
)

REM Set permissions for Jenkins home (Windows equivalent)
echo ✓ Setting up Jenkins home directory permissions

REM Create docker-compose file for Jenkins
echo Creating Jenkins docker-compose configuration...
(
echo version: '3.8'
echo.
echo services:
echo   jenkins:
echo     image: jenkins/jenkins:lts
echo     container_name: jenkins-master
echo     restart: unless-stopped
echo     ports:
echo       - "8080:8080"
echo       - "50000:50000"
echo     volumes:
echo       - ./jenkins_home:/var/jenkins_home
echo       - /var/run/docker.sock:/var/run/docker.sock
echo       - /usr/bin/docker:/usr/bin/docker
echo     environment:
echo       - JENKINS_OPTS=--httpPort=8080
echo       - JAVA_OPTS=-Djenkins.install.runSetupWizard=false
echo     networks:
echo       - jenkins-network
echo.
echo   jenkins-agent:
echo     image: jenkins/inbound-agent:latest
echo     container_name: jenkins-agent
echo     restart: unless-stopped
echo     environment:
echo       - JENKINS_URL=http://jenkins:8080
echo       - JENKINS_SECRET=${JENKINS_AGENT_SECRET}
echo       - JENKINS_AGENT_NAME=docker-agent
echo       - JENKINS_AGENT_WORKDIR=/home/jenkins/agent
echo     volumes:
echo       - /var/run/docker.sock:/var/run/docker.sock
echo       - /usr/bin/docker:/usr/bin/docker
echo     depends_on:
echo       - jenkins
echo     networks:
echo       - jenkins-network
echo.
echo networks:
echo   jenkins-network:
echo     driver: bridge
) > docker-compose.jenkins.yml

echo ✓ Created Jenkins docker-compose configuration
echo.

REM Start Jenkins
echo Starting Jenkins...
docker-compose -f docker-compose.jenkins.yml up -d

if %errorlevel% neq 0 (
    echo ERROR: Failed to start Jenkins!
    pause
    exit /b 1
)

echo ✓ Jenkins started successfully
echo.

REM Wait for Jenkins to start
echo Waiting for Jenkins to initialize...
timeout /t 30 /nobreak >nul 2>&1

REM Get initial admin password
echo ========================================
echo Jenkins Setup Information
echo ========================================
echo.
echo Jenkins is starting up at: http://localhost:8080
echo.
echo To get the initial admin password, run:
echo docker exec jenkins-master cat /var/jenkins_home/secrets/initialAdminPassword
echo.
echo Or check the Jenkins logs:
echo docker logs jenkins-master
echo.

REM Try to get the initial password
echo Attempting to retrieve initial admin password...
docker exec jenkins-master cat /var/jenkins_home/secrets/initialAdminPassword 2>nul
if %errorlevel% equ 0 (
    echo.
    echo ✓ Use the password above to unlock Jenkins
) else (
    echo ⚠ Password not ready yet. Please wait a moment and run:
    echo docker exec jenkins-master cat /var/jenkins_home/secrets/initialAdminPassword
)

echo.
echo ========================================
echo Next Steps:
echo ========================================
echo 1. Open http://localhost:8080 in your browser
echo 2. Use the initial admin password to unlock Jenkins
echo 3. Install suggested plugins
echo 4. Create an admin user
echo 5. Install additional plugins:
echo    - Docker Pipeline
echo    - GitHub Integration
echo    - Blue Ocean
echo    - Pipeline Stage View
echo    - HTML Publisher
echo 6. Configure GitHub webhook for automatic builds
echo 7. Create a new Pipeline job using the Jenkinsfile
echo.
echo ========================================
echo Useful Commands:
echo ========================================
echo Start Jenkins:  docker-compose -f docker-compose.jenkins.yml up -d
echo Stop Jenkins:   docker-compose -f docker-compose.jenkins.yml down
echo View logs:      docker logs jenkins-master
echo Get password:   docker exec jenkins-master cat /var/jenkins_home/secrets/initialAdminPassword
echo.

pause