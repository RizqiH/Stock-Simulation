# Stock Simulation Backend Makefile
# Requires: Docker, Docker Compose, Go

.PHONY: help build dev prod test clean logs shell db-shell redis-shell migrate health stop restart scale

# Default target
help: ## Show this help message
	@echo "Stock Simulation Backend - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development Commands
dev: ## Start development environment
	@echo "Starting development environment..."
	@copy .env.dev .env >nul 2>&1 || echo "Using existing .env file"
	@docker-compose --env-file .env.dev up --build -d
	@echo "Development environment started!"
	@echo "API: http://localhost:8080"
	@echo "Adminer: http://localhost:8081"

dev-logs: ## Show development logs
	@docker-compose logs -f

dev-stop: ## Stop development environment
	@docker-compose down

dev-restart: ## Restart development environment
	@docker-compose restart

dev-rebuild: ## Rebuild and restart development environment
	@docker-compose down
	@docker-compose --env-file .env.dev up --build -d

# Production Commands
prod: ## Start production environment
	@echo "Starting production environment..."
	@copy .env.prod .env >nul 2>&1 || echo "Using existing .env file"
	@docker-compose -f docker-compose.prod.yml up --build -d
	@echo "Production environment started!"
	@echo "API: http://localhost:8080"
	@echo "Load Balancer: http://localhost:80"

prod-logs: ## Show production logs
	@docker-compose -f docker-compose.prod.yml logs -f

prod-stop: ## Stop production environment
	@docker-compose -f docker-compose.prod.yml down

prod-restart: ## Restart production environment
	@docker-compose -f docker-compose.prod.yml restart

prod-scale: ## Scale API service (usage: make prod-scale REPLICAS=3)
	@docker-compose -f docker-compose.prod.yml up -d --scale api=$(or $(REPLICAS),2)

# Build Commands
build: ## Build application binary
	@echo "Building application..."
	@go build -o bin/api ./cmd/api
	@echo "Build complete: bin/api"

build-docker: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t stock-simulation-api .
	@echo "Docker image built: stock-simulation-api"

# Database Commands
db-setup: ## Setup database with migrations
	@echo "Setting up database..."
	@docker-compose exec mysql mysql -u root -p$(MYSQL_ROOT_PASSWORD) < /docker-entrypoint-initdb.d/001_create_database.sql
	@docker-compose exec mysql mysql -u root -p$(MYSQL_ROOT_PASSWORD) < /docker-entrypoint-initdb.d/002_insert_sample_data.sql
	@echo "Database setup complete!"

db-shell: ## Access MySQL shell
	@docker-compose exec mysql mysql -u stockuser -p stock_simulation

db-root-shell: ## Access MySQL shell as root
	@docker-compose exec mysql mysql -u root -p

db-backup: ## Backup database
	@echo "Creating database backup..."
	@docker-compose exec mysql mysqldump -u stockuser -p stock_simulation > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "Backup created!"

db-restore: ## Restore database (usage: make db-restore FILE=backup.sql)
	@echo "Restoring database from $(FILE)..."
	@docker-compose exec -T mysql mysql -u stockuser -p stock_simulation < $(FILE)
	@echo "Database restored!"

# Redis Commands
redis-shell: ## Access Redis CLI
	@docker-compose exec redis redis-cli

redis-flush: ## Flush Redis cache
	@docker-compose exec redis redis-cli FLUSHALL
	@echo "Redis cache flushed!"

# Application Commands
shell: ## Access API container shell
	@docker-compose exec api sh

logs: ## Show all service logs
	@docker-compose logs -f

api-logs: ## Show API service logs only
	@docker-compose logs -f api

health: ## Check service health
	@echo "Checking service health..."
	@docker-compose ps
	@echo ""
	@curl -s http://localhost:8080/health || echo "API health check failed"

# Test Commands
test: ## Run unit tests
	@echo "Running unit tests..."
	@go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out
	@echo "Coverage report: coverage.html"

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@copy .env.example .env.test >nul 2>&1 || echo "Using existing .env.test"
	@docker-compose -f docker-compose.test.yml --env-file .env.test up -d --build
	@echo "Waiting for services to be ready..."
	@timeout /t 30 /nobreak >nul 2>&1 || ping 127.0.0.1 -n 31 >nul
	@docker-compose -f docker-compose.test.yml exec -T api go test -tags=integration ./tests/integration/... || echo "Integration tests completed"
	@docker-compose -f docker-compose.test.yml down -v

test-all: ## Run all tests (unit + integration)
	@echo "Running all tests..."
	@make test-coverage
	@make test-integration

test-ci: ## Run tests for CI environment
	@echo "Running CI tests..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run --timeout=5m

security-scan: ## Run security scan
	@echo "Running security scan..."
	@gosec ./...

vuln-check: ## Check for vulnerabilities
	@echo "Checking for vulnerabilities..."
	@govulncheck ./...

code-quality: ## Run all code quality checks
	@echo "Running code quality checks..."
	@make lint
	@make security-scan
	@make vuln-check

test-setup: ## Setup test environment
	@echo "Setting up test environment..."
	@copy .env.example .env.test >nul 2>&1 || echo "Test env file exists"
	@docker-compose -f docker-compose.test.yml --env-file .env.test up -d mysql redis
	@echo "Waiting for test services..."
	@timeout /t 20 /nobreak >nul 2>&1 || ping 127.0.0.1 -n 21 >nul

test-teardown: ## Teardown test environment
	@echo "Tearing down test environment..."
	@docker-compose -f docker-compose.test.yml down -v
	@docker system prune -f

# Utility Commands
clean: ## Clean up Docker resources
	@echo "Cleaning up Docker resources..."
	@docker-compose down -v --remove-orphans
	@docker system prune -f
	@echo "Cleanup complete!"

clean-all: ## Clean up everything including images
	@echo "Cleaning up all Docker resources..."
	@docker-compose down -v --remove-orphans
	@docker system prune -af
	@echo "Complete cleanup done!"

stop: ## Stop all services
	@docker-compose down
	@docker-compose -f docker-compose.prod.yml down

restart: ## Restart all services
	@docker-compose restart

status: ## Show service status
	@docker-compose ps

env-dev: ## Copy development environment file
	@copy .env.dev .env >nul 2>&1 || echo "File already exists or copied"
	@echo "Development environment file copied!"

env-prod: ## Copy production environment file
	@copy .env.prod .env >nul 2>&1 || echo "File already exists or copied"
	@echo "Production environment file copied!"

# Monitoring Commands
monitor: ## Monitor resource usage
	@docker stats

top: ## Show running processes in containers
	@docker-compose top

# Security Commands
security-scan: ## Run security scan on Docker image
	@echo "Running security scan..."
	@docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy image stock-simulation-api

# Documentation
docs: ## Generate API documentation
	@echo "Generating API documentation..."
	@swag init -g cmd/api/main.go -o docs/
	@echo "Documentation generated in docs/"

# DevOps Commands
jenkins-setup: ## Setup Jenkins for CI/CD
	@echo "Setting up Jenkins..."
	@.\scripts\setup-jenkins.bat

jenkins-start: ## Start Jenkins
	@echo "Starting Jenkins..."
	@docker-compose -f docker-compose.jenkins.yml up -d
	@echo "Jenkins available at http://localhost:8080"

jenkins-stop: ## Stop Jenkins
	@echo "Stopping Jenkins..."
	@docker-compose -f docker-compose.jenkins.yml down

jenkins-logs: ## View Jenkins logs
	@echo "Viewing Jenkins logs..."
	@docker logs -f jenkins-master

jenkins-password: ## Get Jenkins initial admin password
	@echo "Getting Jenkins initial admin password..."
	@docker exec jenkins-master cat /var/jenkins_home/secrets/initialAdminPassword

github-webhook: ## Setup GitHub webhook
	@echo "Setting up GitHub webhook..."
	@.\scripts\setup-github-webhook.bat

ci-local: ## Run CI pipeline locally
	@echo "Running CI pipeline locally..."
	@make code-quality
	@make test-ci
	@echo "Local CI pipeline completed successfully!"

ci-setup: ## Setup complete CI/CD environment
	@echo "Setting up complete CI/CD environment..."
	@make jenkins-setup
	@echo "Please configure Jenkins and then run 'make github-webhook'"
	@echo "See docs/DEVOPS_SETUP.md for detailed instructions"

devops-status: ## Check DevOps services status
	@echo "Checking DevOps services status..."
	@echo "=== Jenkins Status ==="
	@docker ps --filter "name=jenkins" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" || echo "Jenkins not running"
	@echo "\n=== Application Status ==="
	@make status

devops-clean: ## Clean DevOps environment
	@echo "Cleaning DevOps environment..."
	@docker-compose -f docker-compose.jenkins.yml down -v
	@docker system prune -f
	@echo "DevOps environment cleaned"

# Quick start
start: ## Quick start development environment
	@echo "Quick starting Stock Simulation Backend..."
	@make env-dev
	@make dev
	@echo "Waiting for services to be ready..."
	@timeout /t 30 /nobreak >nul 2>&1 || ping 127.0.0.1 -n 31 >nul
	@make health
	@echo ""
	@echo "🚀 Stock Simulation Backend is ready!"
	@echo "📊 API: http://localhost:8080"
	@echo "🗄️  Adminer: http://localhost:8081"
	@echo "📝 Logs: make logs"