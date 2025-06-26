# 🚀 Backend Deployment to Railway - Complete Guide

## Overview
This guide will help you deploy your Go backend to Railway, the same platform where your database is hosted. This ensures optimal performance and simplified configuration.

## 🎯 Why Railway for Go Backend?

### ✅ **Perfect for Go Applications**
- **Native Go Support**: Built-in support for Go applications
- **Long-running Services**: Unlike serverless, perfect for persistent connections
- **WebSocket Support**: Full support for real-time features
- **Database Integration**: Internal networking with your MySQL service

### ✅ **Benefits Over Other Platforms**
- **Vercel**: ❌ Serverless functions, no WebSocket, execution limits
- **Railway**: ✅ Full server, WebSocket, unlimited execution time
- **Same Platform as DB**: ✅ Low latency, no external network costs

## 📋 Prerequisites

### 1. Install Railway CLI
```bash
# Using npm
npm install -g @railway/cli

# Using yarn
yarn global add @railway/cli

# Verify installation
railway --version
```

### 2. Verify Database is Running
- Your MySQL database should already be deployed to Railway
- Note the project name where your database is hosted

## 🚀 Automated Deployment

### Option 1: Use Deployment Script (Recommended)
```bash
# Run the automated deployment script
cd stock-simulation-backend/scripts
./deploy_backend_to_railway.bat
```

The script will:
- ✅ Check prerequisites
- ✅ Login to Railway
- ✅ Link to your database project
- ✅ Set environment variables
- ✅ Deploy your backend
- ✅ Provide deployment URL

## 🔧 Manual Deployment (Step by Step)

### Step 1: Login and Setup
```bash
# Login to Railway
railway login

# Link to existing project (where your MySQL is)
railway link

# Or create new project
railway new
```

### Step 2: Environment Variables
```bash
# Application settings
railway variables set PORT=8080
railway variables set ENV=production
railway variables set JWT_SECRET=your-secure-jwt-secret

# Database connection (linking to MySQL service)
railway variables set DB_HOST=${{MySQL.MYSQL_HOST}}
railway variables set DB_PORT=${{MySQL.MYSQL_PORT}}
railway variables set DB_USER=${{MySQL.MYSQL_USER}}
railway variables set DB_PASSWORD=${{MySQL.MYSQL_PASSWORD}}
railway variables set DB_NAME=${{MySQL.MYSQL_DATABASE}}

# CORS settings
railway variables set CORS_ORIGINS=https://your-frontend-domain.com,http://localhost:3000
```

### Step 3: Deploy
```bash
# Deploy to Railway
railway up
```

### Step 4: Get Deployment URL
```bash
# Get your backend URL
railway domain
```

## 📊 Configuration Details

### Environment Variables Reference
```env
# Required - Application
PORT=8080
ENV=production
JWT_SECRET=your-secure-jwt-secret-128-characters

# Required - Database (auto-linked to MySQL service)
DB_HOST=${{MySQL.MYSQL_HOST}}
DB_PORT=${{MySQL.MYSQL_PORT}}
DB_USER=${{MySQL.MYSQL_USER}}
DB_PASSWORD=${{MySQL.MYSQL_PASSWORD}}
DB_NAME=${{MySQL.MYSQL_DATABASE}}

# Required - CORS
CORS_ORIGINS=https://your-frontend.vercel.app,http://localhost:3000

# Optional - Additional features
REDIS_URL=redis://your-redis-url (if using Redis)
```

### Railway Configuration (railway.toml)
```toml
[build]
builder = "DOCKERFILE"
dockerfilePath = "Dockerfile"

[deploy]
startCommand = "./main"
healthcheckPath = "/health"
healthcheckTimeout = 300
restartPolicyType = "ON_FAILURE"
restartPolicyMaxRetries = 10
```

## 🧪 Testing Your Deployment

### Health Check
```bash
# Test health endpoint
curl https://your-backend-url.railway.app/health

# Expected response:
{
  "status": "healthy",
  "service": "stock-simulation-api",
  "environment": "production",
  "version": "1.0.0"
}
```

### API Endpoints
```bash
# Test stocks API
curl https://your-backend-url.railway.app/api/v1/stocks

# Test WebSocket (using wscat)
wscat -c wss://your-backend-url.railway.app/api/v1/ws
```

### Database Connection
```bash
# Check Railway logs for database connection
railway logs

# Look for:
# ✅ Database connected successfully
# 🚀 Starting StockSim API Server...
# 📡 Server starting on 0.0.0.0:8080
```

## 📈 Monitoring and Management

### View Logs
```bash
# Real-time logs
railway logs

# Specific service logs
railway logs --service backend
```

### Monitor Metrics
- Go to Railway dashboard
- Click on your backend service
- View metrics: CPU, Memory, Network
- Monitor request rates and response times

### Environment Management
```bash
# List all variables
railway variables

# Update variable
railway variables set KEY=value

# Delete variable
railway variables delete KEY
```

## 🔄 CI/CD Setup (Optional)

### GitHub Integration
1. **Connect Repository**: Link your GitHub repo to Railway
2. **Auto-deploy**: Enable automatic deployments on push
3. **Branch Deploy**: Set up preview deployments for PRs

### Environment Setup
```yaml
# .github/workflows/railway.yml
name: Deploy to Railway
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: railway/action@v1
        with:
          token: ${{ secrets.RAILWAY_TOKEN }}
```

## 🛠️ Troubleshooting

### Common Issues

#### 1. Build Failures
```bash
# Check build logs
railway logs --build

# Common fixes:
# - Ensure Dockerfile is in root directory
# - Check Go version compatibility
# - Verify dependencies are available
```

#### 2. Database Connection Issues
```bash
# Check database service status
railway status

# Verify environment variables
railway variables | grep DB_

# Test connection manually
railway shell
```

#### 3. Port Issues
```bash
# Railway automatically sets PORT
# Make sure your app uses os.Getenv("PORT")
# Default: 8080
```

#### 4. Memory/CPU Limits
```bash
# Check resource usage
railway metrics

# Upgrade plan if needed
railway upgrade
```

### Debug Commands
```bash
# Service status
railway status

# Resource usage
railway metrics

# Environment check
railway variables

# Connect to container
railway shell

# Restart service
railway restart
```

## 💰 Cost Optimization

### Railway Pricing
- **Starter**: $5/month per service
- **Pro**: Usage-based pricing
- **Free Tier**: $5 credit monthly

### Cost Tips
1. **Same Project**: Keep backend and database in same project
2. **Resource Monitoring**: Monitor CPU/Memory usage
3. **Sleep Policy**: Use sleep for development environments
4. **Efficient Code**: Optimize database queries and memory usage

## 🔗 Next Steps

### 1. Frontend Configuration
Update your frontend configuration to use the Railway backend URL:

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

### 2. Domain Setup (Optional)
```bash
# Add custom domain
railway domain add yourdomain.com

# Set up DNS
# Add CNAME record: api.yourdomain.com -> your-backend.railway.app
```

### 3. SSL/HTTPS
Railway automatically provides SSL certificates for all deployments.

### 4. Monitoring Setup
- Set up health check monitoring
- Configure alerts for downtime
- Monitor database performance
- Track API response times

## 📞 Support

- **Railway Docs**: https://docs.railway.app
- **Railway Discord**: https://discord.gg/railway
- **Railway GitHub**: https://github.com/railwayapp/railway

---

## ✅ Deployment Checklist

- [ ] Railway CLI installed
- [ ] Database deployed and running
- [ ] Environment variables configured
- [ ] Backend deployed successfully
- [ ] Health check passing
- [ ] API endpoints working
- [ ] WebSocket connection working
- [ ] Frontend configured with backend URL
- [ ] CORS properly configured
- [ ] Monitoring set up

**🎉 Congratulations! Your StockSim backend is now running on Railway!** 