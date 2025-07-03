# Redis Real-Time Price Broadcasting Implementation

## 🎯 **Overview**
Implementasi Redis Pub/Sub untuk Real-Time Price Broadcasting yang memungkinkan aplikasi Stock Simulation untuk:
- Publish price updates ke Redis channels
- Subscribe ke price updates dari Redis
- Scale aplikasi dengan multiple instances
- Cache stock prices untuk performa yang lebih baik

## 🏗️ **Architecture**

### **Before (WebSocket Only)**
```
Price Simulator → RealTimeService → WebSocket Clients
```

### **After (Redis + WebSocket)**
```
Price Simulator → Redis Pub/Sub → Multiple App Instances → WebSocket Clients
                      ↓
                 Price Cache
```

## 📂 **Files Modified/Created**

### **New Files:**
- `internal/core/services/redis_service.go` - Redis service dengan pub/sub functionality
- `REDIS_IMPLEMENTATION.md` - Dokumentasi ini

### **Modified Files:**
- `internal/config/config.go` - Added Redis client initialization
- `internal/core/services/realtime_service.go` - Integrated Redis pub/sub
- `internal/core/services/price_simulator_service.go` - Added Redis publishing
- `internal/core/domain/stock.go` - Added TradingAlert structure
- `cmd/api/main.go` - Integrated Redis service in dependency injection
- `go.mod` - Added go-redis dependency

## 🔧 **Setup & Configuration**

### **1. Start Redis Container**
```bash
# Start Redis with Docker Compose
docker-compose up redis -d

# Or manually with Docker
docker run -d --name redis -p 6379:6379 redis:7-alpine
```

### **2. Environment Variables**
```env
# Redis Configuration
REDIS_URL=redis://localhost:6379
REDIS_PORT=6379
```

### **3. Start Application**
```bash
# Development mode (with price simulation)
go run ./cmd/api/main.go

# Production mode
ENV=production go run ./cmd/api/main.go
```

## 📡 **Redis Channels**

### **Price Updates**
- **Channel**: `stock:price_updates`
- **Format**: JSON PriceUpdateMessage
- **Purpose**: Real-time stock price broadcasts

### **Market Status**
- **Channel**: `stock:market_status`
- **Format**: JSON MarketStatusMessage
- **Purpose**: Market open/close notifications

### **Trading Alerts**
- **Channel**: `stock:trading_alerts`
- **Format**: JSON TradingAlert
- **Purpose**: User-specific trading notifications

## 🔍 **Testing Redis Implementation**

### **1. Check Health Status**
```bash
curl http://localhost:8080/health
```
Response:
```json
{
    "status": "healthy",
    "redis": "connected",
    "redis_available": true
}
```

### **2. Check WebSocket + Redis Status**
```bash
curl http://localhost:8080/api/v1/ws/status
```
Response:
```json
{
    "websocket_enabled": true,
    "connected_clients": 0,
    "redis_enabled": true,
    "redis_connected": true
}
```

### **3. Test Redis Pub/Sub (Development Only)**
```bash
curl -X POST http://localhost:8080/api/v1/dev/redis/publish-test
```

### **4. Monitor Redis Channels**
```bash
# Connect to Redis CLI
redis-cli

# Subscribe to price updates
SUBSCRIBE stock:price_updates

# Subscribe to market status
SUBSCRIBE stock:market_status
```

## 💻 **Code Examples**

### **Publishing Price Updates**
```go
// In Price Simulator Service
priceUpdate := domain.PriceUpdateMessage{
    Symbol:        "AAPL",
    Price:         150.25,
    Change:        2.15,
    ChangePercent: 1.45,
    // ... other fields
}

// Publish to Redis (distributed to all instances)
redisService.PublishPriceUpdate(priceUpdate)

// Also broadcast locally (WebSocket fallback)
realTimeService.BroadcastPriceUpdate(priceUpdate)
```

### **Subscribing to Updates**
```go
// In RealTime Service
priceUpdateChan, err := redisService.SubscribeToPriceUpdates()
if err != nil {
    log.Printf("Failed to subscribe: %v", err)
    return
}

// Forward Redis messages to WebSocket clients
for priceUpdate := range priceUpdateChan {
    s.broadcastToWebSocketClients(priceUpdate)
}
```

### **Caching Stock Prices**
```go
// Cache price for 30 seconds
redisService.CacheStockPrice("AAPL", 150.25, 30*time.Second)

// Retrieve cached price
price, err := redisService.GetCachedStockPrice("AAPL")
```

## 🚀 **Benefits**

### **Scalability**
- Multiple app instances dapat subscribe ke same Redis channels
- Horizontal scaling untuk handle lebih banyak WebSocket connections
- Load balancing antara multiple servers

### **Performance**
- Redis caching untuk stock prices
- Reduced database queries untuk real-time data
- Faster price lookups

### **Reliability**
- Fallback ke WebSocket-only mode kalau Redis tidak available
- Graceful degradation
- No single point of failure

### **Real-Time Broadcasting**
- Instant price updates ke semua connected clients
- Cross-instance message distribution
- Consistent data across all app instances

## 📊 **Monitoring & Debugging**

### **Application Logs**
```
✅ Redis connected successfully
🔌 Subscribed to Redis channel: stock:price_updates
📡 Published price update for AAPL to Redis
📡 Broadcasted AAPL price update to 5 clients
```

### **Redis Monitoring Commands**
```bash
# Check active channels
redis-cli PUBSUB CHANNELS

# Monitor all commands
redis-cli MONITOR

# Check cached prices
redis-cli KEYS "stock:price:*"

# Get specific cached price
redis-cli GET "stock:price:AAPL"
```

### **API Endpoints for Monitoring**
- `GET /health` - Overall system health including Redis
- `GET /api/v1/ws/status` - WebSocket and Redis status
- `GET /api/v1/redis/status` - Redis-specific status
- `GET /api/v1/simulator/status` - Price simulator status

## 🔄 **Fallback Behavior**

Jika Redis tidak tersedia:
1. ✅ Application tetap berjalan normal
2. ✅ WebSocket connections tetap bekerja
3. ✅ Price updates tetap di-broadcast secara lokal
4. ❌ Cross-instance broadcasting disabled
5. ❌ Price caching disabled

## 🎮 **Usage in Frontend**

WebSocket client tidak perlu berubah. Message format tetap sama:

```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/ws');

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    
    if (data.type === 'price_update') {
        // Handle price update
        console.log('Price update:', data.data);
        console.log('Source:', data.source); // "redis" or "local"
    }
};
```

## 📝 **Future Enhancements**

1. **Redis Cluster** - Multi-node Redis setup untuk high availability
2. **Message Persistence** - Store critical messages dengan TTL
3. **User-Specific Channels** - Personal notifications per user
4. **Rate Limiting** - Redis-based rate limiting untuk API calls
5. **Session Storage** - Store WebSocket sessions in Redis
6. **Analytics** - Real-time analytics dengan Redis Streams

## ⚠️ **Important Notes**

1. **Development vs Production**: Redis pub/sub enabled di semua environment, tapi price simulation hanya auto-start di development
2. **Memory Usage**: Monitor Redis memory usage, set appropriate `maxmemory` policy
3. **Network**: Pastikan Redis port (6379) accessible antar app instances
4. **Security**: Gunakan Redis AUTH di production environment
5. **Backup**: Setup Redis persistence (RDB + AOF) untuk data penting

## 🎯 **Conclusion**

Implementasi Redis pub/sub ini memberikan foundation yang kuat untuk real-time features di Stock Simulation aplikasi. Dengan architecture yang scalable dan fallback mechanism yang reliable, aplikasi siap untuk handle multiple users dan real-time trading scenarios. 