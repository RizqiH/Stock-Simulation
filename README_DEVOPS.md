# 🚀 DevOps Setup - Stock Simulation Backend

Panduan cepat untuk setup CI/CD pipeline dengan Jenkins, GitHub, dan Docker.

## ⚡ Quick Start (5 menit)

### 1. Setup Jenkins
```bash
make jenkins-setup
```

### 2. Konfigurasi Jenkins
- Buka http://localhost:8080
- Gunakan password yang ditampilkan
- Install suggested plugins
- Buat admin user

### 3. Setup GitHub Webhook
```bash
make github-webhook
```

### 4. Create Pipeline Job
- New Item → Pipeline
- Repository: `https://github.com/username/stock-simulation-backend.git`
- Script Path: `Jenkinsfile`
- Enable "GitHub hook trigger"

## 🛠️ Available Commands

### Jenkins Management
```bash
make jenkins-setup     # Setup Jenkins dengan Docker
make jenkins-start     # Start Jenkins
make jenkins-stop      # Stop Jenkins
make jenkins-logs      # View logs
make jenkins-password  # Get admin password
```

### CI/CD Operations
```bash
make ci-setup          # Setup complete CI/CD
make ci-local          # Run CI pipeline locally
make github-webhook    # Setup GitHub webhook
make devops-status     # Check services status
make devops-clean      # Clean environment
```

### Testing Commands
```bash
make test              # Unit tests
make test-coverage     # Tests with coverage
make test-integration  # Integration tests
make test-all          # All tests
make code-quality      # Linting + security
```

## 📋 Prerequisites

- ✅ Docker Desktop (running)
- ✅ Git
- ✅ Make
- ✅ GitHub Personal Access Token

## 🔧 Pipeline Features

### Automatic Triggers
- ✅ Push ke repository
- ✅ Pull request
- ✅ Manual trigger

### Quality Checks
- ✅ Go linting (`golangci-lint`)
- ✅ Security scan (`gosec`)
- ✅ Vulnerability check (`govulncheck`)
- ✅ Unit tests dengan coverage
- ✅ Integration tests

### Build & Deploy
- ✅ Docker image build
- ✅ Multi-environment deploy
- ✅ Rollback capability

## 📊 Monitoring

### Jenkins Dashboard
- **URL**: http://localhost:8080
- **Blue Ocean**: http://localhost:8080/blue

### Application
- **API**: http://localhost:8080 (dev)
- **Health**: http://localhost:8080/health
- **Metrics**: http://localhost:8080/metrics

## 🚨 Troubleshooting

### Jenkins tidak start
```bash
# Check Docker
docker ps

# Check logs
make jenkins-logs

# Restart
make jenkins-stop
make jenkins-start
```

### Webhook tidak trigger
```bash
# Check webhook di GitHub settings
# Verify Jenkins URL accessible
# Test manual build
```

### Tests gagal
```bash
# Run locally
make test-all

# Check test environment
make test-setup
make devops-status
```

## 📚 Documentation

- 📖 [Detailed Setup Guide](docs/DEVOPS_SETUP.md)
- 🔧 [Jenkinsfile Configuration](Jenkinsfile)
- 🐳 [Docker Compose Files](docker-compose.*.yml)
- 🧪 [Testing Guide](docs/TESTING.md)

## 🎯 Next Steps

1. **Setup Notifications**
   - Slack integration
   - Email alerts
   - GitHub status checks

2. **Advanced Features**
   - Multi-branch pipeline
   - Parallel testing
   - Performance testing

3. **Production Ready**
   - HTTPS setup
   - Backup strategy
   - Monitoring & alerting

---

**Need help?** Check [DEVOPS_SETUP.md](docs/DEVOPS_SETUP.md) untuk panduan lengkap! 🚀