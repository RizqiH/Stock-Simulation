package handlers

import (
	"net/http"
	"strconv"
	"stock-simulation-backend/internal/core/ports/services"

	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	stockService services.StockService
}

func NewStockHandler(stockService services.StockService) *StockHandler {
	return &StockHandler{
		stockService: stockService,
	}
}

func (h *StockHandler) GetAllStocks(c *gin.Context) {
	stocks, err := h.stockService.GetAllStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stocks": stocks})
}

func (h *StockHandler) GetStockBySymbol(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
		return
	}

	stock, err := h.stockService.GetStockBySymbol(symbol)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stock": stock})
}

func (h *StockHandler) GetTopStocks(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	stocks, err := h.stockService.GetTopStocks(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stocks": stocks})
}

func (h *StockHandler) UpdateStockPrice(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
		return
	}

	var request struct {
		Price float64 `json:"price" binding:"required,min=0.01"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.stockService.UpdateStockPrice(symbol, request.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get updated stock
	stock, err := h.stockService.GetStockBySymbol(symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stock price updated successfully",
		"stock":   stock,
	})
}

// SimulateMarketMovement randomly updates all stock prices for testing
func (h *StockHandler) SimulateMarketMovement(c *gin.Context) {
	err := h.stockService.SimulateMarketMovement()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get all updated stocks
	stocks, err := h.stockService.GetAllStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Market movement simulated successfully",
		"stocks":  stocks,
	})
}