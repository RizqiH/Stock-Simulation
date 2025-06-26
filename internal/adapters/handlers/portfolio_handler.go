package handlers

import (
	"net/http"
	"stock-simulation-backend/internal/core/ports/services"
	"time"
	"log"

	"github.com/gin-gonic/gin"
)

type PortfolioHandler struct {
	portfolioService services.PortfolioService
}

func NewPortfolioHandler(portfolioService services.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{
		portfolioService: portfolioService,
	}
}

func (h *PortfolioHandler) GetUserPortfolio(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	portfolio, err := h.portfolioService.GetUserPortfolio(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"portfolio": portfolio})
}

func (h *PortfolioHandler) GetPortfolioPerformance(c *gin.Context) {
	userID := c.GetInt("userID")  // Fix: use "userID" instead of "user_id"
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get query parameters
	period := c.DefaultQuery("period", "1M")
	
	// Calculate date range based on period
	var startDate time.Time
	now := time.Now()
	
	switch period {
	case "1D":
		startDate = now.AddDate(0, 0, -1)
	case "1W":
		startDate = now.AddDate(0, 0, -7)
	case "1M":
		startDate = now.AddDate(0, -1, 0)
	case "3M":
		startDate = now.AddDate(0, -3, 0)
	case "6M":
		startDate = now.AddDate(0, -6, 0)
	case "1Y":
		startDate = now.AddDate(-1, 0, 0)
	case "ALL":
		startDate = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC) // Far back date
	default:
		startDate = now.AddDate(0, -1, 0) // Default to 1 month
	}

	// Get portfolio performance data
	performanceData, err := h.portfolioService.GetPortfolioPerformanceHistory(userID, startDate, now)
	if err != nil {
		log.Printf("Error getting portfolio performance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get portfolio performance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": performanceData,
		"period": period,
		"start_date": startDate,
		"end_date": now,
	})
}

func (h *PortfolioHandler) GetPortfolioValue(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	value, err := h.portfolioService.GetPortfolioValue(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_value": value})
}

func (h *PortfolioHandler) GetPortfolioSummary(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	summary, err := h.portfolioService.GetPortfolioSummary(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}