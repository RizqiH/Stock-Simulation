package main

import (
    "database/sql"
    "log"
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
    cfg := config.Load()

    // Initialize database
    db, err := sql.Open("mysql", cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }

    // Initialize repositories
    userRepo := mysqlRepo.NewUserRepository(db)
    stockRepo := mysqlRepo.NewStockRepository(db)
    transactionRepo := mysqlRepo.NewTransactionRepository(db)
    portfolioRepo := mysqlRepo.NewPortfolioRepository(db)

    // Initialize services
    userService := services.NewUserService(userRepo)
    stockService := services.NewStockService(stockRepo)
    transactionService := services.NewTransactionService(transactionRepo, portfolioRepo, stockRepo, userRepo)
    portfolioService := services.NewPortfolioService(portfolioRepo, stockRepo)

    // Initialize handlers
    userHandler := handlers.NewUserHandler(userService)
    stockHandler := handlers.NewStockHandler(stockService)
    transactionHandler := handlers.NewTransactionHandler(transactionService)
    portfolioHandler := handlers.NewPortfolioHandler(portfolioService)

    // Setup router
    router := gin.Default()

    // Apply middleware
    router.Use(middleware.CORS())
    router.Use(middleware.Logging())

    // Health check endpoint
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy", "service": "stock-simulation-api"})
    })

    // Public routes
    public := router.Group("/api/v1")
    {
        public.POST("/auth/register", userHandler.Register)
        public.POST("/auth/login", userHandler.Login)
        public.GET("/stocks", stockHandler.GetAllStocks)
        public.GET("/stocks/:symbol", stockHandler.GetStockBySymbol)
        public.GET("/leaderboard", userHandler.GetLeaderboard)
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
    }

    log.Printf("Server starting on port %s", cfg.Port)
    log.Fatal(router.Run(":" + cfg.Port))
}
