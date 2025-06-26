package handlers

import (
	"net/http"
	"stock-simulation-backend/internal/core/ports/services"

	"github.com/gin-gonic/gin"
)

type ChartHandler struct {
	chartService services.ChartService
}

func NewChartHandler(chartService services.ChartService) *ChartHandler {
	return &ChartHandler{
		chartService: chartService,
	}
}

// GetChartData handles GET /api/charts/:symbol
func (h *ChartHandler) GetChartData(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
		return
	}

	period := c.DefaultQuery("period", "30D")
	
	chartData, err := h.chartService.GetChartData(symbol, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": chartData,
	})
}

// GetHistoricalPrices handles GET /api/charts/:symbol/history
func (h *ChartHandler) GetHistoricalPrices(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
		return
	}

	limit := 30
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := parseLimit(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	prices, err := h.chartService.GetHistoricalPrices(symbol, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"symbol": symbol,
			"prices": prices,
		},
	})
}

// GetAvailableSymbols handles GET /api/charts/symbols
func (h *ChartHandler) GetAvailableSymbols(c *gin.Context) {
	symbols, err := h.chartService.GetAvailableSymbols()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"symbols": symbols,
		},
	})
}

// Helper function to parse limit parameter
func parseLimit(limitStr string) (int, error) {
	switch limitStr {
	case "7":
		return 7, nil
	case "30":
		return 30, nil
	case "90":
		return 90, nil
	case "365":
		return 365, nil
	default:
		return 30, nil
	}
} 