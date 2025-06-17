# Stock Simulation Backend

Backend API untuk aplikasi simulasi trading saham yang dibangun dengan Go, menggunakan arsitektur Clean Architecture.

## Fitur

- Manajemen pengguna (registrasi, login)
- Data saham real-time
- Simulasi trading (beli/jual saham)
- Portfolio management
- Tracking performa investasi
- Riwayat transaksi

## Teknologi

- **Go** - Bahasa pemrograman
- **Gin** - Web framework
- **GORM** - ORM untuk database
- **MySQL** - Database
- **JWT** - Authentication
- **Clean Architecture** - Arsitektur aplikasi

## Struktur Project

```
stock-simulation-backend/
├── cmd/api/                 # Entry point aplikasi
├── internal/
│   ├── config/             # Konfigurasi aplikasi
│   ├── core/
│   │   ├── domain/         # Domain models
│   │   ├── ports/          # Interfaces
│   │   └── services/       # Business logic
│   ├── handlers/           # HTTP handlers
│   ├── infrastructure/     # Database, external services
│   └── middleware/         # HTTP middleware
├── migrations/             # Database migrations
├── scripts/               # Setup scripts
└── go.mod                 # Go modules
```

## Prerequisites

- Go 1.19 atau lebih baru
- MySQL 8.0 atau lebih baru
- Git

## Setup Options

### Option 1: Docker Setup (Recommended)

Docker setup adalah cara termudah untuk menjalankan aplikasi dengan semua dependensinya.

#### Prerequisites untuk Docker
- Docker Desktop
- Docker Compose

#### Quick Start dengan Docker

```bash
# Clone repository
git clone <repository-url>
cd stock-simulation-backend

# Quick start development environment
make quick-start
```

Atau manual:

```bash
# Development environment
cd scripts
.\docker-dev.bat

# Atau menggunakan make
make dev
```

#### Docker Commands

```bash
# Development
make dev              # Start development environment
make dev-logs         # View logs
make dev-stop         # Stop services
make dev-restart      # Restart services

# Production
make prod             # Start production environment
make prod-logs        # View production logs
make prod-scale       # Scale API service

# Database
make db-shell         # Access MySQL shell
make redis-shell      # Access Redis CLI

# Utilities
make health           # Check service health
make clean            # Clean up Docker resources
make help             # Show all available commands
```

### Option 2: Manual Setup

#### 1. Install Dependencies

Pastikan MySQL dan Redis sudah terinstall dan berjalan di sistem Anda.

#### 2. Setup Database (Windows)

Jalankan script setup otomatis:

```bash
cd scripts
.\setup_database.bat
```

Script ini akan:
- Membuat database `stock_simulation`
- Membuat tabel yang diperlukan
- Menambahkan data sample saham Indonesia
- Membuat user test

#### 3. Setup Manual

Jika ingin setup manual:

```bash
# Login ke MySQL
mysql -u root -p

# Jalankan migration files
source migrations/001_create_database.sql
source migrations/002_insert_sample_data.sql
```

## Setup Aplikasi

### 1. Clone Repository

```bash
git clone <repository-url>
cd stock-simulation-backend
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Konfigurasi Environment

Copy file `.env.example` ke `.env` dan sesuaikan konfigurasi:

```bash
copy .env.example .env
```

Edit file `.env`:

```env
PORT=8080
DATABASE_URL=root:your_mysql_password@tcp(localhost:3306)/stock_simulation?parseTime=true
JWT_SECRET=your-super-secret-jwt-key
```

### 4. Jalankan Aplikasi

```bash
go run ./cmd/api
```

Atau build terlebih dahulu:

```bash
go build -o bin/api ./cmd/api
./bin/api
```

Aplikasi akan berjalan di `http://localhost:8080`

## Docker Services

Ketika menggunakan Docker, aplikasi akan berjalan dengan services berikut:

### Development Environment
- **API Server**: `http://localhost:8080`
- **MySQL Database**: `localhost:3306`
- **Redis Cache**: `localhost:6379`
- **Adminer** (Database UI): `http://localhost:8081`

### Production Environment
- **Load Balancer**: `http://localhost:80`
- **API Server**: `http://localhost:8080`
- **MySQL Database**: `localhost:3306`
- **Redis Cache**: `localhost:6379`
- **Nginx**: Load balancing dan reverse proxy

### Environment Files

- `.env.dev` - Development configuration
- `.env.prod` - Production configuration
- `.env.example` - Template file

**Penting**: Sebelum deploy production, pastikan mengubah semua nilai `CHANGE_THIS` di `.env.prod`

## API Endpoints

### Authentication
- `POST /api/auth/register` - Registrasi user baru
- `POST /api/auth/login` - Login user

### User
- `GET /api/user/profile` - Get user profile
- `PUT /api/user/profile` - Update user profile

### Stocks
- `GET /api/stocks` - Get daftar saham
- `GET /api/stocks/:symbol` - Get detail saham

### Portfolio
- `GET /api/portfolio` - Get portfolio user
- `GET /api/portfolio/performance` - Get performa portfolio
- `GET /api/portfolio/:symbol` - Get detail holding saham

### Transactions
- `POST /api/transactions/buy` - Beli saham
- `POST /api/transactions/sell` - Jual saham
- `GET /api/transactions` - Get riwayat transaksi

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    balance DECIMAL(15,2) DEFAULT 100000.00,
    total_profit DECIMAL(15,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Stocks Table
```sql
CREATE TABLE stocks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    current_price DECIMAL(10,2) NOT NULL,
    open_price DECIMAL(10,2) NOT NULL,
    high_price DECIMAL(10,2) NOT NULL,
    low_price DECIMAL(10,2) NOT NULL,
    volume BIGINT DEFAULT 0,
    market_cap DECIMAL(20,2) DEFAULT 0.00,
    sector VARCHAR(50)
);
```

### Transactions Table
```sql
CREATE TABLE transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    stock_symbol VARCHAR(10) NOT NULL,
    type ENUM('BUY', 'SELL') NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    total_amount DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Portfolios Table
```sql
CREATE TABLE portfolios (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    stock_symbol VARCHAR(10) NOT NULL,
    quantity INT NOT NULL,
    average_price DECIMAL(10,2) NOT NULL,
    total_cost DECIMAL(15,2) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## Sample Data

Database sudah terisi dengan:
- 20 saham populer Indonesia (BBCA, BBRI, BMRI, TLKM, dll)
- 1 user test:
  - Username: `testuser`
  - Email: `test@example.com`
  - Password: `password` (hash: `$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi`)
  - Balance: Rp 100,000,000

## Development

### Menjalankan Tests

```bash
go test ./...
```

### Build untuk Production

```bash
go build -ldflags="-s -w" -o bin/api ./cmd/api
```

## Troubleshooting

### Docker Issues

#### Docker not running
```bash
# Check Docker status
docker info

# Start Docker Desktop (Windows)
# Atau restart Docker service (Linux)
sudo systemctl restart docker
```

#### Port conflicts
```bash
# Check what's using the port
netstat -ano | findstr :8080

# Stop conflicting services atau ubah port di .env file
```

#### Container build failures
```bash
# Clean Docker cache
make clean-all

# Rebuild from scratch
make dev-rebuild
```

#### Database connection issues
```bash
# Check container logs
make logs

# Access database directly
make db-shell

# Reset database
make dev-stop
docker volume rm stock-simulation-backend_mysql_data
make dev
```

### Manual Setup Issues

#### Database Connection Error

1. Pastikan MySQL berjalan
2. Periksa konfigurasi `DATABASE_URL` di file `.env`
3. Pastikan database `stock_simulation` sudah dibuat
4. Periksa username/password MySQL

#### Port Already in Use

Ubah port di file `.env`:
```env
PORT=8081
```

#### Build Errors

Pastikan Go version minimal 1.19:
```bash
go version
```

Update dependencies:
```bash
go mod tidy
```

### Performance Issues

#### High memory usage
```bash
# Monitor resource usage
make monitor

# Adjust Docker resource limits in docker-compose files
```

#### Slow database queries
```bash
# Check slow query log
make db-shell
SHOW VARIABLES LIKE 'slow_query_log';
```

## Development Workflow

### Using Docker (Recommended)

```bash
# Start development environment
make dev

# Make code changes...

# View logs
make api-logs

# Restart API service after changes
docker-compose restart api

# Run tests
make test

# Stop environment
make dev-stop
```

### Using Make Commands

```bash
# See all available commands
make help

# Quick development setup
make quick-start

# Database operations
make db-shell
make db-backup
make redis-flush

# Monitoring
make health
make monitor
make top
```

### Production Deployment

```bash
# Prepare production environment
cp .env.prod .env
# Edit .env with production values

# Deploy
make prod

# Scale API service
make prod-scale REPLICAS=3

# Monitor
make prod-logs
```

## Testing

```bash
# Unit tests
make test

# Test with coverage
make test-coverage

# Integration tests
make test-integration
```

## Security

```bash
# Security scan
make security-scan

# Check for vulnerabilities
go list -json -m all | nancy sleuth
```

## Contributing

1. Fork repository
2. Buat feature branch
3. Commit changes
4. Push ke branch
5. Buat Pull Request

## License

MIT License#   S t o c k - S i m u l a t i o n  
 