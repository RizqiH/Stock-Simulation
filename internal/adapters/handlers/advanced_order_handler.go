package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"stock-simulation-backend/internal/core/domain"
	"stock-simulation-backend/internal/core/ports/repositories"
	"stock-simulation-backend/internal/core/ports/services"
)

type AdvancedOrderHandler struct {
	orderService      services.AdvancedOrderService
	commissionService services.CommissionService
	marketService     services.MarketService
	realTimeService   services.RealTimeService
}

func NewAdvancedOrderHandler(
	orderService services.AdvancedOrderService,
	commissionService services.CommissionService,
	marketService services.MarketService,
	realTimeService services.RealTimeService,
) *AdvancedOrderHandler {
	return &AdvancedOrderHandler{
		orderService:      orderService,
		commissionService: commissionService,
		marketService:     marketService,
		realTimeService:   realTimeService,
	}
}

// NewSimplifiedAdvancedOrderHandler creates a simplified handler for development
func NewSimplifiedAdvancedOrderHandler() *SimplifiedAdvancedOrderHandler {
	return &SimplifiedAdvancedOrderHandler{}
}

// SimplifiedAdvancedOrderHandler provides placeholder responses for advanced order endpoints
type SimplifiedAdvancedOrderHandler struct{}

func (h *SimplifiedAdvancedOrderHandler) CreateOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Advanced order creation coming soon",
		"status":  "not_implemented",
	})
}

func (h *SimplifiedAdvancedOrderHandler) CreateOCOOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OCO order creation coming soon",
		"status":  "not_implemented",
	})
}

func (h *SimplifiedAdvancedOrderHandler) ModifyOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Order modification coming soon",
		"status":  "not_implemented",
	})
}

func (h *SimplifiedAdvancedOrderHandler) CancelOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Order cancellation coming soon",
		"status":  "not_implemented",
	})
}

func (h *SimplifiedAdvancedOrderHandler) CancelAllOrders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"cancelled_count": 0,
		"message":         "Cancel all orders coming soon",
		"status":          "not_implemented",
	})
}

func (h *SimplifiedAdvancedOrderHandler) GetUserOrders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"orders":      []interface{}{},
		"total":       0,
		"page":        1,
		"page_size":   20,
		"total_pages": 0,
		"message":     "Order retrieval coming soon",
		"status":      "not_implemented",
	})
}

func (h *SimplifiedAdvancedOrderHandler) GetOrderByID(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error":   "Order not found",
		"message": "Order details coming soon",
		"status":  "not_implemented",
	})
}

func (h *SimplifiedAdvancedOrderHandler) GetActiveOrders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"orders": []interface{}{},
		"total":  0,
		"message": "Active orders coming soon",
		"status":  "not_implemented",
	})
}

func (h *SimplifiedAdvancedOrderHandler) GetOrderStatistics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"statistics": gin.H{
			"total_orders":            0,
			"pending_orders":          0,
			"executed_orders":         0,
			"cancelled_orders":        0,
			"partially_filled":        0,
			"success_rate":           0.0,
			"average_execution_time": 0.0,
			"total_commission":       0.0,
			"total_fees":            0.0,
		},
		"message": "Order statistics coming soon",
		"status":  "not_implemented",
	})
}

func (h *SimplifiedAdvancedOrderHandler) GetExecutionMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"metrics": gin.H{
			"user_id":                getUserIDFromContext(c),
			"timeframe":             c.DefaultQuery("timeframe", "month"),
			"total_orders":          0,
			"executed_orders":       0,
			"cancelled_orders":      0,
			"partially_filled":      0,
			"average_execution_time": 0.0,
			"fill_rate":            0.0,
			"average_slippage":     0.0,
			"best_execution":       0.0,
			"worst_execution":      0.0,
			"total_commission":     0.0,
			"total_fees":          0.0,
		},
		"message": "Execution metrics coming soon",
		"status":  "not_implemented",
	})
}

func (h *SimplifiedAdvancedOrderHandler) GetSlippageAnalysis(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"analysis": gin.H{
			"symbol":                  c.Param("symbol"),
			"user_id":                getUserIDFromContext(c),
			"total_trades":           0,
			"average_slippage":       0.0,
			"median_slippage":        0.0,
			"slippage_std_dev":       0.0,
			"best_execution":         0.0,
			"worst_slippage":         0.0,
			"market_order_slippage":  0.0,
			"limit_order_slippage":   0.0,
		},
		"message": "Slippage analysis coming soon",
		"status":  "not_implemented",
	})
}

func (h *SimplifiedAdvancedOrderHandler) CalculateCommission(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"calculation": gin.H{
			"base_commission":   5.00,
			"regulatory_fees":   0.50,
			"clearing_fees":     0.25,
			"platform_fees":     1.00,
			"total_commission":  6.75,
			"effective_rate":    0.001,
		},
		"message": "Commission calculation coming soon",
		"status":  "not_implemented",
	})
}

// @Summary Create a new advanced order
// @Description Create market, limit, stop-loss, take-profit, or trailing stop orders
// @Tags orders
// @Accept json
// @Produce json
// @Param order body domain.OrderRequest true "Order details"
// @Success 201 {object} OrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 422 {object} ValidationErrorResponse
// @Router /orders [post]
func (h *AdvancedOrderHandler) CreateOrder(c *gin.Context) {
	fmt.Printf("🔥 CreateOrder called - Method: %s, Path: %s\n", c.Request.Method, c.Request.URL.Path)
	
	userID := getUserIDFromContext(c)
	fmt.Printf("🔥 UserID from context: %d\n", userID)
	
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User authentication required",
		})
		return
	}

	var request domain.OrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Printf("🔥 JSON binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, ValidationErrorResponse{
			Error:   "Invalid request",
			Message: "Please check your order parameters",
			Details: parseValidationErrors(err),
		})
		return
	}
	
	fmt.Printf("🔥 Order request received: %+v\n", request)

	// Validate order
	if err := h.orderService.ValidateOrder(userID, &request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Error:   "Order validation failed",
			Message: err.Error(),
		})
		return
	}

	// Create order
	fmt.Printf("🔥 Calling orderService.CreateOrder...\n")
	order, err := h.orderService.CreateOrder(userID, &request)
	if err != nil {
		fmt.Printf("🔥 Order creation failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create order",
			Message: err.Error(),
		})
		return
	}
	
	fmt.Printf("🔥 Order created successfully: %+v\n", order)

	// Calculate commission estimate
	var commissionCalc *domain.CommissionCalculation
	if h.commissionService != nil {
		commissionCalc, _ = h.commissionService.CalculateCommission(
			userID,
			float64(order.Quantity)*order.MarketPrice,
			order.OrderType,
			"stock",
		)
	}

	response := OrderResponse{
		Order:              *order,
		CommissionEstimate: commissionCalc,
		Message:            "Order created successfully",
	}

	fmt.Printf("🔥 Sending response: %+v\n", response)
	c.JSON(http.StatusCreated, response)
}

// @Summary Create OCO (One-Cancels-Other) order
// @Description Create a pair of linked orders where execution of one cancels the other
// @Tags orders
// @Accept json
// @Produce json
// @Param oco body OCOOrderRequest true "OCO order details"
// @Success 201 {object} OCOOrderResponse
// @Failure 400 {object} ErrorResponse
// @Router /orders/oco [post]
func (h *AdvancedOrderHandler) CreateOCOOrder(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	var request OCOOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ValidationErrorResponse{
			Error:   "Invalid request",
			Message: "Please check your OCO order parameters",
			Details: parseValidationErrors(err),
		})
		return
	}

	// Create OCO orders
	parentOrder, linkedOrder, err := h.orderService.CreateOCOOrder(
		userID,
		&request.ParentOrder,
		&request.LinkedOrder,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create OCO order",
			Message: err.Error(),
		})
		return
	}

	response := OCOOrderResponse{
		ParentOrder: *parentOrder,
		LinkedOrder: *linkedOrder,
		Message:     "OCO order created successfully",
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary Modify an existing order
// @Description Modify price, quantity, or other parameters of an active order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param modifications body services.OrderModificationRequest true "Order modifications"
// @Success 200 {object} OrderResponse
// @Failure 400 {object} ErrorResponse
// @Router /orders/{id} [put]
func (h *AdvancedOrderHandler) ModifyOrder(c *gin.Context) {
	userID := getUserIDFromContext(c)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid order ID",
			Message: "Order ID must be a valid number",
		})
		return
	}

	var request services.OrderModificationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ValidationErrorResponse{
			Error:   "Invalid request",
			Details: parseValidationErrors(err),
		})
		return
	}

	modifiedOrder, err := h.orderService.ModifyOrder(userID, orderID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to modify order",
			Message: err.Error(),
		})
		return
	}

	response := OrderResponse{
		Order:   *modifiedOrder,
		Message: "Order modified successfully",
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Cancel an order
// @Description Cancel an active order
// @Tags orders
// @Param id path int true "Order ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Router /orders/{id} [delete]
func (h *AdvancedOrderHandler) CancelOrder(c *gin.Context) {
	userID := getUserIDFromContext(c)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid order ID",
		})
		return
	}

	err = h.orderService.CancelOrder(userID, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to cancel order",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Order cancelled successfully",
	})
}

// @Summary Cancel all orders
// @Description Cancel all active orders, optionally filtered by symbol
// @Tags orders
// @Param symbol query string false "Symbol to filter by"
// @Success 200 {object} CancelAllOrdersResponse
// @Router /orders/cancel-all [post]
func (h *AdvancedOrderHandler) CancelAllOrders(c *gin.Context) {
	userID := getUserIDFromContext(c)
	symbol := c.Query("symbol")
	
	var symbolPtr *string
	if symbol != "" {
		symbolPtr = &symbol
	}

	cancelledCount, err := h.orderService.CancelAllOrders(userID, symbolPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to cancel orders",
			Message: err.Error(),
		})
		return
	}

	response := CancelAllOrdersResponse{
		CancelledCount: cancelledCount,
		Message:        "Orders cancelled successfully",
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get user orders
// @Description Get orders for the authenticated user with filtering and pagination
// @Tags orders
// @Param status query string false "Order status filter"
// @Param symbol query string false "Symbol filter"
// @Param order_type query string false "Order type filter"
// @Param start_date query string false "Start date (RFC3339 format)"
// @Param end_date query string false "End date (RFC3339 format)"
// @Param limit query int false "Number of results per page" default(20)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} OrderListResponse
// @Router /orders [get]
func (h *AdvancedOrderHandler) GetUserOrders(c *gin.Context) {
	fmt.Printf("🟦 GetUserOrders called - Method: %s, Path: %s\n", c.Request.Method, c.Request.URL.Path)
	userID := getUserIDFromContext(c)
	
	// Parse query parameters
	criteria := &repositories.OrderSearchCriteria{
		UserID: &userID,
		Limit:  getIntQuery(c, "limit", 20),
		Offset: getIntQuery(c, "offset", 0),
	}

	if status := c.Query("status"); status != "" {
		orderStatus := domain.OrderStatus(status)
		criteria.Status = &orderStatus
	}

	if symbol := c.Query("symbol"); symbol != "" {
		criteria.Symbol = &symbol
	}

	if orderType := c.Query("order_type"); orderType != "" {
		ot := domain.OrderType(orderType)
		criteria.OrderType = &ot
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse(time.RFC3339, startDate); err == nil {
			criteria.StartDate = &t
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse(time.RFC3339, endDate); err == nil {
			criteria.EndDate = &t
		}
	}

	result, err := h.orderService.SearchOrders(userID, criteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve orders",
			Message: err.Error(),
		})
		return
	}

	response := OrderListResponse{
		Orders:     result.Orders,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get order by ID
// @Description Get detailed information about a specific order
// @Tags orders
// @Param id path int true "Order ID"
// @Success 200 {object} OrderResponse
// @Failure 404 {object} ErrorResponse
// @Router /orders/{id} [get]
func (h *AdvancedOrderHandler) GetOrderByID(c *gin.Context) {
	userID := getUserIDFromContext(c)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid order ID",
		})
		return
	}

	order, err := h.orderService.GetOrderByID(userID, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Order not found",
			Message: err.Error(),
		})
		return
	}

	response := OrderResponse{
		Order: *order,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get active orders
// @Description Get all active orders for the authenticated user
// @Tags orders
// @Success 200 {object} OrderListResponse
// @Router /orders/active [get]
func (h *AdvancedOrderHandler) GetActiveOrders(c *gin.Context) {
	userID := getUserIDFromContext(c)

	orders, err := h.orderService.GetActiveOrders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve active orders",
			Message: err.Error(),
		})
		return
	}

	response := OrderListResponse{
		Orders: orders,
		Total:  len(orders),
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get order statistics
// @Description Get comprehensive order statistics for the user
// @Tags orders
// @Success 200 {object} OrderStatsResponse
// @Router /orders/statistics [get]
func (h *AdvancedOrderHandler) GetOrderStatistics(c *gin.Context) {
	userID := getUserIDFromContext(c)

	stats, err := h.orderService.GetOrderStatistics(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve order statistics",
			Message: err.Error(),
		})
		return
	}

	response := OrderStatsResponse{
		Statistics: *stats,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get execution metrics
// @Description Get detailed execution performance metrics
// @Tags orders
// @Param timeframe query string false "Timeframe (day, week, month, year)" default(month)
// @Success 200 {object} ExecutionMetricsResponse
// @Router /orders/execution-metrics [get]
func (h *AdvancedOrderHandler) GetExecutionMetrics(c *gin.Context) {
	userID := getUserIDFromContext(c)
	timeframe := c.DefaultQuery("timeframe", "month")

	metrics, err := h.orderService.GetExecutionMetrics(userID, timeframe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve execution metrics",
			Message: err.Error(),
		})
		return
	}

	response := ExecutionMetricsResponse{
		Metrics: *metrics,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get slippage analysis
// @Description Get slippage analysis for a specific symbol
// @Tags orders
// @Param symbol path string true "Stock symbol"
// @Success 200 {object} SlippageAnalysisResponse
// @Router /orders/slippage/{symbol} [get]
func (h *AdvancedOrderHandler) GetSlippageAnalysis(c *gin.Context) {
	userID := getUserIDFromContext(c)
	symbol := c.Param("symbol")

	analysis, err := h.orderService.GetSlippageAnalysis(userID, symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve slippage analysis",
			Message: err.Error(),
		})
		return
	}

	response := SlippageAnalysisResponse{
		Analysis: *analysis,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Calculate commission
// @Description Calculate commission for a trade
// @Tags orders
// @Accept json
// @Produce json
// @Param request body CommissionRequest true "Commission calculation request"
// @Success 200 {object} CommissionResponse
// @Router /commission/calculate [post]
func (h *AdvancedOrderHandler) CalculateCommission(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	var request CommissionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: "Please check commission calculation parameters",
		})
		return
	}

	if h.commissionService == nil {
		// Fallback calculation
		c.JSON(http.StatusOK, gin.H{
			"calculation": gin.H{
				"base_commission":   5.00,
				"regulatory_fees":   0.50,
				"clearing_fees":     0.25,
				"platform_fees":     1.00,
				"total_commission":  6.75,
				"effective_rate":    0.001,
			},
		})
		return
	}

	calculation, err := h.commissionService.CalculateCommission(
		userID,
		request.TradeValue,
		domain.OrderType(request.OrderType),
		request.AssetType,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to calculate commission",
			Message: err.Error(),
		})
		return
	}

	response := CommissionResponse{
		Calculation: *calculation,
	}

	c.JSON(http.StatusOK, response)
}

// Response types
type OrderResponse struct {
	Order              domain.Order                     `json:"order"`
	CommissionEstimate *domain.CommissionCalculation    `json:"commission_estimate,omitempty"`
	Message            string                           `json:"message,omitempty"`
}

type OCOOrderRequest struct {
	ParentOrder domain.OrderRequest `json:"parent_order" binding:"required"`
	LinkedOrder domain.OrderRequest `json:"linked_order" binding:"required"`
}

type OCOOrderResponse struct {
	ParentOrder domain.Order `json:"parent_order"`
	LinkedOrder domain.Order `json:"linked_order"`
	Message     string       `json:"message"`
}

type OrderListResponse struct {
	Orders     []domain.Order `json:"orders"`
	Total      int           `json:"total"`
	Page       int           `json:"page,omitempty"`
	PageSize   int           `json:"page_size,omitempty"`
	TotalPages int           `json:"total_pages,omitempty"`
}

type CancelAllOrdersResponse struct {
	CancelledCount int    `json:"cancelled_count"`
	Message        string `json:"message"`
}

type OrderStatsResponse struct {
	Statistics domain.OrderStats `json:"statistics"`
}

type ExecutionMetricsResponse struct {
	Metrics repositories.OrderExecutionMetrics `json:"metrics"`
}

type SlippageAnalysisResponse struct {
	Analysis repositories.SlippageAnalysis `json:"analysis"`
}

type CommissionRequest struct {
	TradeValue float64 `json:"trade_value" binding:"required"`
	OrderType  string  `json:"order_type" binding:"required"`
	AssetType  string  `json:"asset_type"`
}

type CommissionResponse struct {
	Calculation domain.CommissionCalculation `json:"calculation"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

// Helper functions
func getUserIDFromContext(c *gin.Context) int {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}

func getIntQuery(c *gin.Context, key string, defaultValue int) int {
	if value := c.Query(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func parseValidationErrors(err error) map[string]string {
	// This would parse validation errors from the binding
	// Implementation depends on your validation library
	return map[string]string{
		"general": err.Error(),
	}
} 