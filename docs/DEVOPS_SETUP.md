# DevOps Setup Guide

Panduan lengkap untuk mengatur CI/CD pipeline menggunakan Jenkins, GitHub, dan Docker untuk Stock Simulation Backend.

## 📋 Prerequisites

### Software Requirements
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (Windows)
- [Git](https://git-scm.com/downloads)
- [Go 1.24+](https://golang.org/dl/)
- [Make](http://gnuwin32.sourceforge.net/packages/make.htm) (untuk Windows)

### GitHub Requirements
- GitHub repository dengan admin access
- Personal Access Token dengan permissions:
  - `repo` (full repository access)
  - `admin:repo_hook` (webhook management)

## 🚀 Quick Start

### 1. Setup Jenkins

```bash
# Jalankan script setup Jenkins
.\scripts\setup-jenkins.bat
```

Script ini akan:
- Membuat Jenkins container dengan Docker support
- Mengatur volume untuk persistent data
- Mengexpose Jenkins di port 8080
- Menampilkan initial admin password

### 2. Konfigurasi Jenkins

1. **Akses Jenkins**: Buka http://localhost:8080
2. **Unlock Jenkins**: Gunakan initial admin password
3. **Install Plugins**: Pilih "Install suggested plugins"
4. **Create Admin User**: Buat user admin
5. **Install Additional Plugins**:
   - Docker Pipeline
   - GitHub Integration
   - Blue Ocean
   - HTML Publisher
   - Go Plugin

### 3. Setup GitHub Webhook

```bash
# Jalankan script setup webhook
.\scripts\setup-github-webhook.bat
```

Script ini akan:
- Membuat webhook di GitHub repository
- Mengkonfigurasi trigger untuk push dan pull request
- Mengarahkan ke Jenkins webhook endpoint

### 4. Create Jenkins Pipeline Job

1. **New Item** → **Pipeline**
2. **Pipeline Configuration**:
   - Definition: Pipeline script from SCM
   - SCM: Git
   - Repository URL: `https://github.com/username/stock-simulation-backend.git`
   - Script Path: `Jenkinsfile`
3. **Build Triggers**:
   - ✅ GitHub hook trigger for GITScm polling
4. **Save**

## 🔧 Manual Setup

### Jenkins Manual Installation

```yaml
# docker-compose.jenkins.yml
version: '3.8'
services:
  jenkins:
    image: jenkins/jenkins:lts
    container_name: jenkins-master
    restart: unless-stopped
    ports:
      - "8080:8080"
      - "50000:50000"
    volumes:
      - ./jenkins_home:/var/jenkins_home
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - JENKINS_OPTS=--httpPort=8080
```

```bash
# Start Jenkins
docker-compose -f docker-compose.jenkins.yml up -d

# Get initial password
docker exec jenkins-master cat /var/jenkins_home/secrets/initialAdminPassword
```

### GitHub Webhook Manual Setup

1. **Repository Settings** → **Webhooks** → **Add webhook**
2. **Payload URL**: `http://your-jenkins-url:8080/github-webhook/`
3. **Content type**: `application/json`
4. **Events**: Push, Pull requests
5. **Active**: ✅

## 🧪 Testing Setup

### Local Testing

```bash
# Run unit tests
make test

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration

# Run all tests
make test-all

# Code quality checks
make code-quality
```

### CI/CD Testing Environment

Testing environment menggunakan:
- **Database**: MySQL (port 3308)
- **Cache**: Redis (port 6380)
- **API**: Go application (port 8081)

```bash
# Setup test environment
make test-setup

# Run CI tests
make test-ci

# Teardown test environment
make test-teardown
```

## 📊 Pipeline Stages

### 1. Checkout
- Clone repository dari GitHub
- Checkout ke branch yang di-trigger

### 2. Environment Setup
- Setup Go environment
- Install dependencies
- Prepare test environment

### 3. Code Quality
- **Linting**: `golangci-lint`
- **Security Scan**: `gosec`
- **Vulnerability Check**: `govulncheck`

### 4. Testing
- **Unit Tests**: dengan race detection
- **Coverage Report**: minimum 80%
- **Integration Tests**: dengan test database

### 5. Build
- Build Docker image
- Tag dengan commit SHA dan branch
- Push ke Docker registry (optional)

### 6. Deploy
- **Staging**: Auto-deploy untuk main branch
- **Production**: Manual approval required

## 🔍 Monitoring & Debugging

### Jenkins Logs

```bash
# View Jenkins logs
docker logs jenkins-master

# Follow logs
docker logs -f jenkins-master
```

### Pipeline Debugging

```bash
# Check container status
docker ps

# Check test results
make test-coverage
open coverage.html

# Check application logs
make logs
```

### Common Issues

#### 1. Docker Permission Issues
```bash
# Add Jenkins user to docker group (Linux/WSL)
sudo usermod -aG docker jenkins

# Restart Jenkins
docker-compose -f docker-compose.jenkins.yml restart
```

#### 2. GitHub Webhook Not Triggering
- Check webhook delivery in GitHub settings
- Verify Jenkins URL is accessible from internet
- Check firewall settings for port 8080

#### 3. Test Database Connection Issues
```bash
# Check test environment
make test-setup
docker-compose -f docker-compose.test.yml ps

# Check database logs
docker-compose -f docker-compose.test.yml logs mysql
```

## 📈 Advanced Configuration

### Multi-branch Pipeline

1. **New Item** → **Multibranch Pipeline**
2. **Branch Sources** → **GitHub**
3. **Repository**: `https://github.com/username/stock-simulation-backend`
4. **Scan Multibranch Pipeline Triggers**: Periodically if not otherwise run

### Blue Ocean Interface

1. Install Blue Ocean plugin
2. Access: http://localhost:8080/blue
3. Visual pipeline editor dan monitoring

### Notifications

#### Slack Integration
```groovy
// Jenkinsfile
post {
    success {
        slackSend(
            channel: '#ci-cd',
            color: 'good',
            message: "✅ Build Success: ${env.JOB_NAME} - ${env.BUILD_NUMBER}"
        )
    }
    failure {
        slackSend(
            channel: '#ci-cd',
            color: 'danger',
            message: "❌ Build Failed: ${env.JOB_NAME} - ${env.BUILD_NUMBER}"
        )
    }
}
```

#### Email Notifications
```groovy
// Jenkinsfile
post {
    always {
        emailext(
            subject: "Build ${currentBuild.result}: ${env.JOB_NAME} - ${env.BUILD_NUMBER}",
            body: "Build ${currentBuild.result}\n\nCheck console output at ${env.BUILD_URL}",
            to: "team@company.com"
        )
    }
}
```

## 🔐 Security Best Practices

### 1. Credentials Management
- Gunakan Jenkins Credentials Store
- Jangan hardcode secrets di Jenkinsfile
- Rotate credentials secara berkala

### 2. Access Control
- Setup role-based access control
- Limit admin access
- Use GitHub OAuth untuk authentication

### 3. Network Security
- Gunakan HTTPS untuk Jenkins
- Restrict webhook access
- Setup VPN untuk production access

## 📚 Additional Resources

- [Jenkins Documentation](https://www.jenkins.io/doc/)
- [Docker Pipeline Plugin](https://plugins.jenkins.io/docker-workflow/)
- [GitHub Integration](https://plugins.jenkins.io/github/)
- [Go Plugin](https://plugins.jenkins.io/golang/)
- [Blue Ocean](https://www.jenkins.io/projects/blueocean/)

## 🆘 Support

Jika mengalami masalah:
1. Check logs: `docker logs jenkins-master`
2. Verify webhook delivery di GitHub
3. Test manual build di Jenkins
4. Check network connectivity
5. Review Jenkins system configuration

---

**Happy DevOps! 🚀**