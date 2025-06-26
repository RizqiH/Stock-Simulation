package main

import (
    "database/sql"
    "log"
    "net/http"
    "stock-simulation-backend/internal/adapters/handlers"
    "stock-simulation-backend/internal/adapters/middleware"
    mysqlRepo "stock-simulation-backend/internal/adapters/repositories/mysql"
    "stock-simulation-backend/internal/config"
    "stock-simulation-backend/internal/core/services"

    "github.com/gin-gonic/gin"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()
    
    // Set Gin mode based on environment
    gin.SetMode(cfg.Server.GinMode)
    
    log.Printf("🚀 Starting StockSim API Server...")
    log.Printf("📋 Environment: %s", cfg.Server.ENV)
    log.Printf("🌐 Server: %s", cfg.GetServerAddress())

    // Initialize database
    dsn := cfg.GetDSN()
    log.Printf("🔌 Connecting to database...")
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("❌ Failed to connect to database:", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal("❌ Failed to ping database:", err)
    }
    
    log.Printf("✅ Database connected successfully")

    // Initialize repositories
    log.Printf("🏗️ Initializing repositories...")
    userRepo := mysqlRepo.NewUserRepository(db)
    stockRepo := mysqlRepo.NewStockRepository(db)
    transactionRepo := mysqlRepo.NewTransactionRepository(db)
    portfolioRepo := mysqlRepo.NewPortfolioRepository(db)
    historicalPriceRepo := mysqlRepo.NewHistoricalPriceRepository(db)
    advancedOrderRepo := mysqlRepo.NewAdvancedOrderRepository(db)

    // Initialize services
    log.Printf("⚙️ Initializing services...")
    userService := services.NewUserService(userRepo)
    stockService := services.NewStockService(stockRepo)
    transactionService := services.NewTransactionService(transactionRepo, portfolioRepo, stockRepo, userRepo)
    portfolioService := services.NewPortfolioService(portfolioRepo, stockRepo)
    chartService := services.NewChartService(historicalPriceRepo)
    advancedOrderService := services.NewAdvancedOrderService(advancedOrderRepo, stockRepo, portfolioRepo, userRepo, transactionService)
    commissionService := services.NewCommissionService()
    
    // Initialize real-time service for WebSocket support
    log.Printf("🔄 Initializing real-time services...")
    realTimeService := services.NewRealTimeService()
    realTimeService.Start()
    
    // Initialize price simulator service with WebSocket support
    priceSimulator := services.NewPriceSimulatorService(stockRepo, historicalPriceRepo, realTimeService)
    
    // Start automatic price simulation (only in development or if explicitly enabled)
    if cfg.IsDevelopment() {
        log.Printf("📈 Starting automatic price simulation (development mode)...")
        priceSimulator.Start()
    } else {
        log.Printf("📈 Price simulation disabled in production (start manually if needed)")
    }
    
    // Ensure services stop when application exits
    defer priceSimulator.Stop()

    // Initialize handlers
    log.Printf("🎛️ Initializing handlers...")
    userHandler := handlers.NewUserHandler(userService)
    stockHandler := handlers.NewStockHandler(stockService)
    transactionHandler := handlers.NewTransactionHandler(transactionService)
    portfolioHandler := handlers.NewPortfolioHandler(portfolioService)
    chartHandler := handlers.NewChartHandler(chartService)
    
    // Use the full advanced order handler
    advancedOrderHandler := handlers.NewAdvancedOrderHandler(
        advancedOrderService,
        commissionService,
        nil, // marketService - will implement later if needed
        nil, // realTimeService - will implement later if needed
    )

    // Setup router
    router := gin.Default()

    // Apply middleware
    log.Printf("🛡️ Setting up middleware...")
    
    // CORS middleware with config
    router.Use(func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        
        // Check if origin is allowed
        allowed := false
        for _, allowedOrigin := range cfg.CORS.AllowedOrigins {
            if allowedOrigin == "*" || allowedOrigin == origin {
                allowed = true
                break
            }
        }
        
        if allowed {
            c.Header("Access-Control-Allow-Origin", origin)
        }
        
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
        c.Header("Access-Control-Allow-Credentials", "true")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    })
    
    router.Use(middleware.Logging())

    // Health check endpoint
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "healthy", 
            "service": "stock-simulation-api",
            "environment": cfg.Server.ENV,
            "version": "1.0.0",
        })
    })
    
    // Root endpoint
    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "StockSim Pro API",
            "version": "1.0.0",
            "environment": cfg.Server.ENV,
            "endpoints": gin.H{
                "health": "/health",
                "api": "/api/v1",
                "websocket": "/api/v1/ws",
            },
        })
    })

    // WebSocket endpoint (before other routes)
    router.GET("/api/v1/ws", gin.WrapH(http.HandlerFunc(realTimeService.HandleWebSocket)))

    // Public routes
    public := router.Group("/api/v1")
    {
        public.POST("/auth/register", userHandler.Register)
        public.POST("/auth/login", userHandler.Login)
        public.GET("/stocks", stockHandler.GetAllStocks)
        public.GET("/stocks/top", stockHandler.GetTopStocks)
        public.GET("/stocks/:symbol", stockHandler.GetStockBySymbol)
        public.GET("/leaderboard", userHandler.GetLeaderboard)
        
        // Chart routes (public)
        public.GET("/charts/:symbol", chartHandler.GetChartData)
        public.GET("/charts/:symbol/history", chartHandler.GetHistoricalPrices)
        public.GET("/charts/symbols", chartHandler.GetAvailableSymbols)
        
        // Price simulation routes (public for testing)
        public.PUT("/stocks/:symbol/price", stockHandler.UpdateStockPrice)
        public.POST("/stocks/simulate", stockHandler.SimulateMarketMovement)
        
        // Price simulator control endpoints
        public.GET("/simulator/status", func(c *gin.Context) {
            status := priceSimulator.GetStatus()
            c.JSON(200, gin.H{"simulator": status})
        })
        
        // WebSocket status endpoint
        public.GET("/ws/status", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "websocket_enabled": true,
                "connected_clients": realTimeService.GetConnectedClientsCount(),
                "endpoint": "/api/v1/ws",
            })
        })
        
        // Only allow simulator control in development
        if cfg.IsDevelopment() {
            public.POST("/simulator/start", func(c *gin.Context) {
                priceSimulator.Start()
                c.JSON(200, gin.H{"message": "Price simulator started"})
            })
            public.POST("/simulator/stop", func(c *gin.Context) {
                priceSimulator.Stop()
                c.JSON(200, gin.H{"message": "Price simulator stopped"})
            })
        }
    }

    // Protected routes
    protected := router.Group("/api/v1")
    protected.Use(middleware.Auth())
    {
        // User routes
        protected.GET("/profile", userHandler.GetProfile)
        protected.PUT("/profile", userHandler.UpdateProfile)

        // Transaction routes
        protected.POST("/transactions/buy", transactionHandler.BuyStock)
        protected.POST("/transactions/sell", transactionHandler.SellStock)
        protected.GET("/transactions", transactionHandler.GetUserTransactions)

        // Portfolio routes
        protected.GET("/portfolio", portfolioHandler.GetUserPortfolio)
        protected.GET("/portfolio/performance", portfolioHandler.GetPortfolioPerformance)
        protected.GET("/portfolio/value", portfolioHandler.GetPortfolioValue)
        protected.GET("/portfolio/summary", portfolioHandler.GetPortfolioSummary)
        
        // Advanced Order routes
        protected.POST("/orders", advancedOrderHandler.CreateOrder)
        protected.POST("/orders/oco", advancedOrderHandler.CreateOCOOrder)
        protected.GET("/orders", advancedOrderHandler.GetUserOrders)
        protected.GET("/orders/active", advancedOrderHandler.GetActiveOrders)
        protected.GET("/orders/:id", advancedOrderHandler.GetOrderByID)
        protected.PUT("/orders/:id", advancedOrderHandler.ModifyOrder)
        protected.DELETE("/orders/:id", advancedOrderHandler.CancelOrder)
        protected.POST("/orders/cancel-all", advancedOrderHandler.CancelAllOrders)
        protected.GET("/orders/statistics", advancedOrderHandler.GetOrderStatistics)
        protected.GET("/orders/execution-metrics", advancedOrderHandler.GetExecutionMetrics)
        protected.GET("/orders/slippage/:symbol", advancedOrderHandler.GetSlippageAnalysis)
        protected.POST("/commission/calculate", advancedOrderHandler.CalculateCommission)
    }

    log.Printf("🎯 All systems initialized successfully!")
    log.Printf("📡 Server starting on %s", cfg.GetServerAddress())
    log.Printf("🌐 API endpoints available at: http://%s/api/v1", cfg.GetServerAddress())
    log.Printf("💬 WebSocket available at: ws://%s/api/v1/ws", cfg.GetServerAddress())
    
    if cfg.IsProduction() {
        log.Printf("🔒 Running in PRODUCTION mode")
    } else {
        log.Printf("🔧 Running in DEVELOPMENT mode")
    }
    
    // Start server
    log.Fatal(router.Run(":" + cfg.Server.Port))
}
