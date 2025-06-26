# Railway Database Deployment Guide

## 🚀 Deploy StockSim Database to Railway

### Prerequisites
- Railway account (sign up at https://railway.app)
- GitHub account
- Local database backup

### Step 1: Create Railway Account
1. Go to https://railway.app
2. Sign in with GitHub
3. Verify your account

### Step 2: Create MySQL Database Service
1. **Create New Project**
   ```
   - Click "New Project"
   - Select "Empty Project"
   - Name: "stocksim-database"
   ```

2. **Add MySQL Service**
   ```
   - Click "Add Service"
   - Select "Database"
   - Choose "MySQL"
   - Railway will provision MySQL instance
   ```

3. **Get Database Credentials**
   ```
   - Click on MySQL service
   - Go to "Variables" tab
   - Copy these values:
     * MYSQL_HOST
     * MYSQL_PORT
     * MYSQL_USER
     * MYSQL_PASSWORD
     * MYSQL_DATABASE
     * DATABASE_URL (complete connection string)
   ```

### Step 3: Update Application Configuration

#### Update Go Application
1. **Create environment file** (`stock-simulation-backend/.env.production`):
   ```env
   # Railway MySQL Configuration
   DB_HOST=your-railway-mysql-host.railway.app
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=your-railway-password
   DB_NAME=railway
   DATABASE_URL=mysql://root:password@host:3306/railway
   
   # Application Settings
   PORT=8080
   JWT_SECRET=your-jwt-secret-key
   API_BASE_URL=https://your-backend-domain.railway.app
   ```

2. **Update config.go** to read from environment:
   ```go
   package config
   
   import (
       "os"
       "strconv"
   )
   
   type Config struct {
       Database DatabaseConfig
       Server   ServerConfig
       JWT      JWTConfig
   }
   
   type DatabaseConfig struct {
       Host     string
       Port     int
       User     string
       Password string
       DBName   string
       URL      string
   }
   
   type ServerConfig struct {
       Port string
   }
   
   type JWTConfig struct {
       Secret string
   }
   
   func LoadConfig() *Config {
       port, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
       
       return &Config{
           Database: DatabaseConfig{
               Host:     getEnv("DB_HOST", "localhost"),
               Port:     port,
               User:     getEnv("DB_USER", "root"),
               Password: getEnv("DB_PASSWORD", "root"),
               DBName:   getEnv("DB_NAME", "stock_simulation"),
               URL:      getEnv("DATABASE_URL", ""),
           },
           Server: ServerConfig{
               Port: getEnv("PORT", "8080"),
           },
           JWT: JWTConfig{
               Secret: getEnv("JWT_SECRET", "your-secret-key"),
           },
       }
   }
   
   func getEnv(key, defaultValue string) string {
       if value := os.Getenv(key); value != "" {
           return value
       }
       return defaultValue
   }
   ```

### Step 4: Database Migration

#### Method 1: Direct SQL Import (Recommended)
1. **Export local database**:
   ```bash
   mysqldump -u root -proot -h localhost -P 3307 stock_simulation > stocksim_backup.sql
   ```

2. **Connect to Railway MySQL**:
   ```bash
   mysql -h your-railway-host.railway.app -u root -p railway < stocksim_backup.sql
   ```

#### Method 2: Using Railway CLI
1. **Install Railway CLI**:
   ```bash
   npm install -g @railway/cli
   ```

2. **Login to Railway**:
   ```bash
   railway login
   ```

3. **Connect to your project**:
   ```bash
   railway link
   ```

4. **Upload database**:
   ```bash
   railway run mysql -h $MYSQL_HOST -u $MYSQL_USER -p$MYSQL_PASSWORD $MYSQL_DATABASE < stocksim_backup.sql
   ```

### Step 5: Deploy Backend Application

#### Option A: Deploy Backend to Railway
1. **Create new service for backend**:
   ```
   - In same project, click "Add Service"
   - Select "GitHub Repo"
   - Connect your backend repository
   ```

2. **Configure build settings**:
   ```
   - Build Command: go build -o main cmd/api/main.go
   - Start Command: ./main
   ```

3. **Add environment variables**:
   ```
   DB_HOST=${{MySQL.MYSQL_HOST}}
   DB_PORT=${{MySQL.MYSQL_PORT}}
   DB_USER=${{MySQL.MYSQL_USER}}
   DB_PASSWORD=${{MySQL.MYSQL_PASSWORD}}
   DB_NAME=${{MySQL.MYSQL_DATABASE}}
   DATABASE_URL=${{MySQL.DATABASE_URL}}
   PORT=8080
   JWT_SECRET=your-jwt-secret
   ```

#### Option B: Update Local Backend for Railway DB
1. **Update connection in main.go**:
   ```go
   func main() {
       config := config.LoadConfig()
       
       var dsn string
       if config.Database.URL != "" {
           dsn = config.Database.URL
       } else {
           dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
               config.Database.User,
               config.Database.Password,
               config.Database.Host,
               config.Database.Port,
               config.Database.DBName,
           )
       }
       
       db, err := sql.Open("mysql", dsn)
       if err != nil {
           log.Fatal("Failed to connect to database:", err)
       }
       defer db.Close()
       
       // Rest of your application...
   }
   ```

### Step 6: Update Frontend Configuration

1. **Update API base URL** in `nuxt.config.ts`:
   ```typescript
   export default defineNuxtConfig({
     runtimeConfig: {
       public: {
         apiBase: process.env.API_BASE_URL || 'https://your-backend.railway.app/api'
       }
     }
   })
   ```

2. **Update composables/useApi.ts**:
   ```typescript
   const config = useRuntimeConfig()
   const baseURL = config.public.apiBase
   ```

### Step 7: Test Connection

1. **Test database connection**:
   ```bash
   # From your local machine
   mysql -h your-railway-host.railway.app -u root -p railway
   ```

2. **Verify data**:
   ```sql
   SHOW TABLES;
   SELECT COUNT(*) FROM users;
   SELECT COUNT(*) FROM stocks;
   ```

### Step 8: Production Environment Variables

Create production environment files:

#### Backend `.env.production`:
```env
DB_HOST=containers-us-west-xxx.railway.app
DB_PORT=3306
DB_USER=root
DB_PASSWORD=xxx
DB_NAME=railway
JWT_SECRET=production-jwt-secret-key
PORT=8080
CORS_ORIGINS=https://your-frontend-domain.com
```

#### Frontend `.env.production`:
```env
API_BASE_URL=https://your-backend.railway.app/api
```

### Step 9: Security Configuration

1. **Configure CORS for production**:
   ```go
   func corsMiddleware() gin.HandlerFunc {
       config := cors.DefaultConfig()
       
       if os.Getenv("ENV") == "production" {
           config.AllowOrigins = []string{
               "https://your-frontend-domain.com",
               "https://your-custom-domain.com",
           }
       } else {
           config.AllowAllOrigins = true
       }
       
       return cors.New(config)
   }
   ```

2. **Set secure headers**:
   ```go
   r.Use(func(c *gin.Context) {
       c.Header("X-Content-Type-Options", "nosniff")
       c.Header("X-Frame-Options", "DENY")
       c.Header("X-XSS-Protection", "1; mode=block")
       c.Next()
   })
   ```

### Troubleshooting

#### Common Issues:
1. **Connection timeout**: Check firewall settings
2. **Authentication failed**: Verify credentials
3. **Database not found**: Ensure database name is correct
4. **SSL connection**: Railway requires SSL, add `?tls=true` to connection string

#### Debug Commands:
```bash
# Test connection
telnet your-host.railway.app 3306

# Check environment variables
railway variables

# View logs
railway logs
```

### Cost Optimization
- **Starter Plan**: $5/month for database
- **Pro Plan**: Usage-based pricing
- **Optimize queries**: Use indexes, limit results
- **Connection pooling**: Implement connection pooling in Go

### Backup Strategy
1. **Automated backups**: Railway provides automatic backups
2. **Manual backups**: Export data regularly
3. **Migration scripts**: Keep migration scripts versioned

```bash
# Create backup
mysqldump -h your-host.railway.app -u root -p railway > backup_$(date +%Y%m%d).sql

# Restore backup
mysql -h your-host.railway.app -u root -p railway < backup_20241227.sql
``` 