# 🚀 Railway Deployment - Quick Start Guide

## ⚡ 5-Minute Railway Setup

### 1. Create Railway Account
```bash
1. Visit https://railway.app
2. Sign in with GitHub
3. Verify your account
```

### 2. Deploy Database (Automated Script)
```bash
# Run the automated deployment script
cd stock-simulation-backend/scripts
./deploy_to_railway.bat
```

**The script will:**
- ✅ Create local database backup
- ✅ Guide you through Railway setup
- ✅ Upload data to Railway MySQL
- ✅ Create production config files
- ✅ Test the connection

### 3. Manual Railway Setup (Alternative)

#### Step 3a: Create MySQL Service
1. **Create Project**: 
   - Click "New Project" → "Empty Project"
   - Name: `stocksim-database`

2. **Add MySQL**: 
   - Click "Add Service" → "Database" → "MySQL"
   - Wait for provisioning (2-3 minutes)

3. **Get Credentials**:
   - Click MySQL service → "Connect" tab → "Public Networking"
   - Copy credentials (format: containers-us-west-xxx.railway.app:XXXX)
   - **Important**: Use the public host and port, NOT the internal ones

#### Step 3b: Upload Database
```bash
# Export local database
mysqldump -u root -proot -h localhost -P 3307 stock_simulation > backup.sql

# Import to Railway (replace with your credentials)
mysql -h containers-us-west-xxx.railway.app -P 6543 -u root -pYOUR_PASSWORD railway < backup.sql
```

### 4. Configure Application

#### Create Production Environment File
```bash
# Copy example and edit
cp config.env.example .env.production

# Edit with your Railway credentials
DB_HOST=containers-us-west-xxx.railway.app
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-railway-password
DB_NAME=railway
ENV=production
JWT_SECRET=your-production-secret
```

### 5. Test Connection
```bash
# Test Railway database connection
mysql -h your-railway-host.railway.app -u root -p railway

# Verify tables
SHOW TABLES;
SELECT COUNT(*) FROM users;
```

### 6. Deploy Backend to Railway (Optional)

#### Option A: Deploy Backend Service
1. **Add Backend Service**:
   - In same project: "Add Service" → "GitHub Repo"
   - Connect your repository

2. **Configure Build**:
   ```
   Build Command: go build -o main cmd/api/main.go
   Start Command: ./main
   ```

3. **Environment Variables**:
   ```
   DB_HOST=${{MySQL.MYSQL_HOST}}
   DB_USER=${{MySQL.MYSQL_USER}}
   DB_PASSWORD=${{MySQL.MYSQL_PASSWORD}}
   DB_NAME=${{MySQL.MYSQL_DATABASE}}
   PORT=8080
   ENV=production
   JWT_SECRET=your-production-secret
   ```

#### Option B: Run Local Backend with Railway DB
```bash
# Copy production config to local
cp .env.production .env

# Start backend
go run cmd/api/main.go
```

### 7. Update Frontend

#### Update API Configuration
```typescript
// nuxt.config.ts
export default defineNuxtConfig({
  runtimeConfig: {
    public: {
      apiBase: process.env.API_BASE_URL || 'https://your-backend.railway.app/api/v1'
    }
  }
})
```

## 📊 Cost & Features

### Railway Pricing
- **Starter**: $5/month for MySQL
- **Pro**: Usage-based pricing
- **Free Tier**: $5 credit monthly

### Features Included
- ✅ **Automatic Backups**: Daily snapshots
- ✅ **SSL/TLS**: Encrypted connections
- ✅ **Monitoring**: Built-in metrics
- ✅ **Scaling**: Automatic resource scaling
- ✅ **99.9% Uptime**: High availability

## 🔧 Configuration Reference

### Required Environment Variables
```env
# Database (from Railway)
DB_HOST=containers-us-west-xxx.railway.app
DB_PORT=3306
DB_USER=root
DB_PASSWORD=xxx
DB_NAME=railway

# Application
PORT=8080
ENV=production
JWT_SECRET=your-production-secret

# CORS (update with your domains)
CORS_ORIGINS=https://your-frontend.vercel.app
```

### Optional Variables
```env
# Redis (if using)
REDIS_URL=redis://xxx

# Custom settings
HOST=0.0.0.0
DATABASE_URL=mysql://user:pass@host:3306/db
```

## 🎯 Next Steps

### 1. Frontend Deployment
- Deploy to **Vercel**: `vercel --prod`
- Deploy to **Netlify**: `netlify deploy --prod`
- Update `CORS_ORIGINS` with your frontend URL

### 2. Custom Domain (Optional)
- Add custom domain in Railway dashboard
- Update DNS settings
- SSL certificate auto-generated

### 3. Monitoring
- Check Railway dashboard for metrics
- Monitor database usage
- Set up alerts for downtime

## 🛠️ Troubleshooting

### Common Issues

#### Connection Refused
```bash
# Check if Railway service is running
railway status

# Test connection
telnet your-host.railway.app 3306
```

#### Authentication Failed
```bash
# Verify credentials in Railway dashboard
railway variables
```

#### SSL Issues
```bash
# Add SSL parameter to connection
DATABASE_URL=mysql://user:pass@host:3306/db?tls=true
```

### Debug Commands
```bash
# Railway CLI commands
railway login
railway link your-project
railway logs
railway shell

# Database commands
railway run mysql -h $MYSQL_HOST -u $MYSQL_USER -p$MYSQL_PASSWORD $MYSQL_DATABASE
```

## 📞 Support

- **Railway Docs**: https://docs.railway.app
- **Railway Discord**: https://discord.gg/railway
- **StockSim Issues**: Create GitHub issue

---

## ✅ Success Checklist

- [ ] Railway account created
- [ ] MySQL service deployed
- [ ] Database uploaded successfully
- [ ] Environment variables configured
- [ ] Connection tested
- [ ] Backend deployment working
- [ ] Frontend updated with new API URL
- [ ] CORS configured for frontend domain
- [ ] Production testing completed

**🎉 Congratulations! Your StockSim database is now running on Railway!** 