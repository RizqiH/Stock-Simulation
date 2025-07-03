package services

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"stock-simulation-backend/internal/core/domain"
	"stock-simulation-backend/internal/core/ports/repositories"
	"sync"
	"time"
)

type PriceSimulatorService struct {
	stockRepo           repositories.StockRepository
	historicalPriceRepo repositories.HistoricalPriceRepository
	realTimeService     *RealTimeService
	redisService        *RedisService
	running             bool
	stopChan            chan bool
	mu                  sync.RWMutex
	
	// Configuration
	updateInterval time.Duration
	volatility     float64
	maxChange      float64
}

func NewPriceSimulatorService(
	stockRepo repositories.StockRepository, 
	historicalPriceRepo repositories.HistoricalPriceRepository,
	realTimeService *RealTimeService,
	redisService *RedisService,
) *PriceSimulatorService {
	return &PriceSimulatorService{
		stockRepo:           stockRepo,
		historicalPriceRepo: historicalPriceRepo,
		realTimeService:     realTimeService,
		redisService:        redisService,
		running:             false,
		stopChan:            make(chan bool),
		updateInterval:      5 * time.Second,  // Update every 5 seconds
		volatility:          2.0,               // 2% max normal change
		maxChange:           5.0,               // 5% max extreme change
	}
}

// Start begins the automatic price simulation
func (s *PriceSimulatorService) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.running {
		log.Println("⚠️ Price simulator already running")
		return
	}
	
	s.running = true
	s.stopChan = make(chan bool)
	
	go s.runSimulation()
	
	if s.redisService != nil {
		log.Println("🚀 Price simulator started with Redis pub/sub - updating every", s.updateInterval)
	} else {
		log.Println("🚀 Price simulator started with WebSocket only - updating every", s.updateInterval)
	}
}

// Stop halts the price simulation
func (s *PriceSimulatorService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if !s.running {
		log.Println("⚠️ Price simulator not running")
		return
	}
	
	s.running = false
	close(s.stopChan)
	log.Println("⏹️ Price simulator stopped")
}

// IsRunning returns whether the simulator is currently active
func (s *PriceSimulatorService) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// runSimulation is the main simulation loop
func (s *PriceSimulatorService) runSimulation() {
	ticker := time.NewTicker(s.updateInterval)
	defer ticker.Stop()
	
	// Create random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	if s.redisService != nil {
		log.Println("📈 Starting automatic price updates with Redis pub/sub and real-time broadcasting...")
	} else {
		log.Println("📈 Starting automatic price updates with WebSocket broadcasting...")
	}
	
	for {
		select {
		case <-ticker.C:
			s.updateAllPrices(rng)
		case <-s.stopChan:
			log.Println("📉 Price simulation stopped")
			return
		}
	}
}

// updateAllPrices updates all stock prices with realistic movements
func (s *PriceSimulatorService) updateAllPrices(rng *rand.Rand) {
	stocks, err := s.stockRepo.GetAll()
	if err != nil {
		log.Printf("❌ Failed to get stocks for price update: %v", err)
		return
	}
	
	if len(stocks) == 0 {
		return
	}
	
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("\n📊 [%s] Updating %d stock prices...\n", timestamp, len(stocks))
	
	updatedCount := 0
	for _, stock := range stocks {
		oldPrice := stock.CurrentPrice
		newPrice := s.generateRealisticPrice(stock, rng)
		
		// Update stock price in database
		err := s.stockRepo.UpdatePrice(stock.Symbol, newPrice)
		if err != nil {
			log.Printf("⚠️ Failed to update %s: %v", stock.Symbol, err)
			continue
		}
		
		// Calculate and display change
		change := newPrice - oldPrice
		changePercent := (change / oldPrice) * 100
		
		indicator := s.getPriceIndicator(change)
		fmt.Printf("%s %s: $%.2f → $%.2f (%+.2f%%)\n",
			indicator, stock.Symbol, oldPrice, newPrice, changePercent)
		
		// Create real-time price update message
		priceUpdate := domain.PriceUpdateMessage{
			Symbol:        stock.Symbol,
			Price:         newPrice,
			Change:        change,
			ChangePercent: changePercent,
			Volume:        stock.Volume,
			High:          math.Max(oldPrice, newPrice), // Simplified for demo
			Low:           math.Min(oldPrice, newPrice), // Simplified for demo
			Open:          oldPrice,                     // Using previous price as open
			PreviousClose: stock.PreviousClose,
			LastTradeTime: time.Now(),
			MarketCap:     &stock.MarketCap,
		}
		
		// Publish to Redis first (this will be picked up by all instances)
		if s.redisService != nil {
			if err := s.redisService.PublishPriceUpdate(priceUpdate); err != nil {
				log.Printf("⚠️ Failed to publish %s to Redis: %v", stock.Symbol, err)
			}
			
			// Also cache the price in Redis for quick access
			if err := s.redisService.CacheStockPrice(stock.Symbol, newPrice, 30*time.Second); err != nil {
				log.Printf("⚠️ Failed to cache %s price in Redis: %v", stock.Symbol, err)
			}
		}
		
		// Broadcast to local WebSocket clients (fallback or direct)
		if s.realTimeService != nil {
			s.realTimeService.BroadcastPriceUpdate(priceUpdate)
		}
		
		// Save historical price data every 30 seconds (6 updates)
		if updatedCount%6 == 0 {
			s.saveHistoricalPrice(stock.Symbol, oldPrice, newPrice, float64(stock.Volume))
		}
		
		updatedCount++
	}
	
	if updatedCount > 0 {
		fmt.Printf("✅ Updated %d/%d stock prices", updatedCount, len(stocks))
		
		// Show broadcasting status
		if s.redisService != nil {
			fmt.Print(" 📡 Published to Redis")
		}
		
		if s.realTimeService != nil {
			clientCount := s.realTimeService.GetConnectedClientsCount()
			if clientCount > 0 {
				fmt.Printf(" 🔌 WebSocket clients: %d", clientCount)
			}
		}
		fmt.Println()
	}
}

// saveHistoricalPrice saves price data for charting
func (s *PriceSimulatorService) saveHistoricalPrice(symbol string, oldPrice, newPrice, volume float64) {
	if s.historicalPriceRepo == nil {
		return
	}
	
	now := time.Now()
	
	// Create historical price entry
	historicalPrice := &domain.HistoricalPrice{
		Symbol:    symbol,
		Date:      now,
		Open:      oldPrice,
		High:      math.Max(oldPrice, newPrice),
		Low:       math.Min(oldPrice, newPrice),
		Close:     newPrice,
		Volume:    int64(volume),
		CreatedAt: now,
	}
	
	// Save to database
	err := s.historicalPriceRepo.Create(historicalPrice)
	if err != nil {
		log.Printf("⚠️ Failed to save historical price for %s: %v", symbol, err)
	} else {
		fmt.Printf("💾 Saved historical price for %s\n", symbol)
	}
}

// generateRealisticPrice creates realistic price movements
func (s *PriceSimulatorService) generateRealisticPrice(stock domain.Stock, rng *rand.Rand) float64 {
	currentPrice := stock.CurrentPrice
	
	// Market trends (simulate bull/bear market influences)
	marketTrend := s.getMarketTrend(rng)
	
	// Base volatility
	baseVolatility := s.volatility
	
	// Stock-specific volatility based on price
	if currentPrice < 50 {
		baseVolatility *= 1.5 // Penny stocks more volatile
	} else if currentPrice > 1000 {
		baseVolatility *= 0.8 // High-price stocks less volatile
	}
	
	// Random change with market bias
	randomChange := (rng.Float64() - 0.5) * 2 * baseVolatility // -volatility% to +volatility%
	trendInfluence := marketTrend * 0.3 // Market trend contributes 30%
	
	changePercent := randomChange + trendInfluence
	
	// Apply extreme events (rare large moves)
	if rng.Float64() < 0.02 { // 2% chance
		extremeChange := (rng.Float64() - 0.5) * 2 * s.maxChange
		changePercent = extremeChange
		if math.Abs(extremeChange) > 3 {
			fmt.Printf("💥 EXTREME MOVE: %s %+.1f%%\n", stock.Symbol, extremeChange)
		}
	}
	
	// Apply change
	change := currentPrice * (changePercent / 100)
	newPrice := currentPrice + change
	
	// Ensure price doesn't go below $0.01
	if newPrice < 0.01 {
		newPrice = 0.01
	}
	
	// Round to 2 decimal places
	return math.Round(newPrice*100) / 100
}

// getMarketTrend simulates overall market sentiment
func (s *PriceSimulatorService) getMarketTrend(rng *rand.Rand) float64 {
	// Simulate market cycles
	hour := time.Now().Hour()
	
	// Market opening hours tend to be more volatile
	if hour >= 9 && hour <= 10 {
		return (rng.Float64() - 0.5) * 2 // Higher volatility at open
	}
	
	// Lunch time usually calmer
	if hour >= 12 && hour <= 13 {
		return (rng.Float64() - 0.5) * 0.5 // Lower volatility
	}
	
	// Normal market hours
	return (rng.Float64() - 0.5) * 1.0
}

// getPriceIndicator returns emoji indicator for price movement
func (s *PriceSimulatorService) getPriceIndicator(change float64) string {
	if change > 0.5 {
		return "🚀" // Strong up
	} else if change > 0 {
		return "📈" // Up
	} else if change < -0.5 {
		return "💥" // Strong down
	} else if change < 0 {
		return "📉" // Down
	}
	return "➡️" // No change
}

// SetUpdateInterval allows configuring update frequency
func (s *PriceSimulatorService) SetUpdateInterval(interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.updateInterval = interval
	log.Printf("⚙️ Update interval changed to %v", interval)
}

// SetVolatility allows configuring price volatility
func (s *PriceSimulatorService) SetVolatility(volatility float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.volatility = volatility
	log.Printf("⚙️ Volatility changed to %.1f%%", volatility)
}

// GetStatus returns current simulator status
func (s *PriceSimulatorService) GetStatus() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	status := map[string]interface{}{
		"running":         s.running,
		"update_interval": s.updateInterval.String(),
		"volatility":      s.volatility,
		"max_change":      s.maxChange,
		"redis_enabled":   s.redisService != nil,
		"redis_connected": false,
	}
	
	if s.redisService != nil {
		status["redis_connected"] = s.redisService.GetConnectionStatus()
	}
	
	if s.realTimeService != nil {
		status["connected_clients"] = s.realTimeService.GetConnectedClientsCount()
	}
	
	return status
}

// PublishMarketEvent publishes special market events to Redis
func (s *PriceSimulatorService) PublishMarketEvent(eventType, message string) {
	if s.redisService != nil {
		if err := s.redisService.PublishMarketStatus(eventType + ": " + message); err != nil {
			log.Printf("⚠️ Failed to publish market event: %v", err)
		}
	}
	
	if s.realTimeService != nil {
		s.realTimeService.BroadcastMarketStatus(eventType + ": " + message)
	}
} 