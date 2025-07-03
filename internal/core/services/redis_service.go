package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"stock-simulation-backend/internal/core/domain"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisService(redisClient *redis.Client) *RedisService {
	if redisClient == nil {
		log.Println("⚠️ Redis client is nil, Redis service will be disabled")
		return nil
	}

	return &RedisService{
		client: redisClient,
		ctx:    context.Background(),
	}
}

// PublishPriceUpdate publishes a price update to Redis
func (s *RedisService) PublishPriceUpdate(update domain.PriceUpdateMessage) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("redis service not available")
	}

	// Convert price update to JSON
	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal price update: %w", err)
	}

	// Publish to Redis channel
	channel := "stock:price_updates"
	err = s.client.Publish(s.ctx, channel, data).Err()
	if err != nil {
		return fmt.Errorf("failed to publish to Redis: %w", err)
	}

	log.Printf("📡 Published price update for %s to Redis", update.Symbol)
	return nil
}

// PublishToChannel publishes a message to a specific Redis channel
func (s *RedisService) PublishToChannel(channel string, message interface{}) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("redis service not available")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = s.client.Publish(s.ctx, channel, data).Err()
	if err != nil {
		return fmt.Errorf("failed to publish to channel %s: %w", channel, err)
	}

	return nil
}

// SubscribeToPriceUpdates subscribes to price updates and returns a channel
func (s *RedisService) SubscribeToPriceUpdates() (<-chan domain.PriceUpdateMessage, error) {
	if s == nil || s.client == nil {
		return nil, fmt.Errorf("redis service not available")
	}

	channel := "stock:price_updates"
	pubsub := s.client.Subscribe(s.ctx, channel)

	// Channel to send parsed price updates
	priceUpdateChan := make(chan domain.PriceUpdateMessage, 100)

	go func() {
		defer close(priceUpdateChan)
		defer pubsub.Close()

		log.Printf("🔌 Subscribed to Redis channel: %s", channel)

		// Listen for messages
		for {
			msg, err := pubsub.ReceiveMessage(s.ctx)
			if err != nil {
				log.Printf("❌ Error receiving Redis message: %v", err)
				continue
			}

			// Parse the message
			var priceUpdate domain.PriceUpdateMessage
			if err := json.Unmarshal([]byte(msg.Payload), &priceUpdate); err != nil {
				log.Printf("❌ Error parsing price update: %v", err)
				continue
			}

			// Send to channel
			select {
			case priceUpdateChan <- priceUpdate:
				// Message sent successfully
			default:
				// Channel is full, drop message
				log.Println("⚠️ Price update channel full, dropping message")
			}
		}
	}()

	return priceUpdateChan, nil
}

// SubscribeToChannel subscribes to a specific Redis channel
func (s *RedisService) SubscribeToChannel(channel string) (<-chan *redis.Message, error) {
	if s == nil || s.client == nil {
		return nil, fmt.Errorf("redis service not available")
	}

	pubsub := s.client.Subscribe(s.ctx, channel)
	messageChan := make(chan *redis.Message, 100)

	go func() {
		defer close(messageChan)
		defer pubsub.Close()

		log.Printf("🔌 Subscribed to Redis channel: %s", channel)

		for {
			msg, err := pubsub.ReceiveMessage(s.ctx)
			if err != nil {
				log.Printf("❌ Error receiving message from %s: %v", channel, err)
				continue
			}

			select {
			case messageChan <- msg:
				// Message sent successfully
			default:
				// Channel is full, drop message
				log.Printf("⚠️ Channel %s is full, dropping message", channel)
			}
		}
	}()

	return messageChan, nil
}

// CacheStockPrice caches a stock price in Redis with TTL
func (s *RedisService) CacheStockPrice(symbol string, price float64, ttl time.Duration) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("redis service not available")
	}

	key := fmt.Sprintf("stock:price:%s", symbol)
	err := s.client.Set(s.ctx, key, price, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to cache stock price: %w", err)
	}

	return nil
}

// GetCachedStockPrice retrieves a cached stock price from Redis
func (s *RedisService) GetCachedStockPrice(symbol string) (float64, error) {
	if s == nil || s.client == nil {
		return 0, fmt.Errorf("redis service not available")
	}

	key := fmt.Sprintf("stock:price:%s", symbol)
	result, err := s.client.Get(s.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, fmt.Errorf("price not found in cache")
		}
		return 0, fmt.Errorf("failed to get cached price: %w", err)
	}

	// Parse price from string
	var price float64
	if err := json.Unmarshal([]byte(result), &price); err != nil {
		return 0, fmt.Errorf("failed to parse cached price: %w", err)
	}

	return price, nil
}

// PublishMarketStatus publishes market status updates
func (s *RedisService) PublishMarketStatus(status string) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("redis service not available")
	}

	message := map[string]interface{}{
		"type":      "market_status",
		"status":    status,
		"timestamp": time.Now().Unix(),
	}

	return s.PublishToChannel("stock:market_status", message)
}

// PublishTradingAlert publishes trading alerts
func (s *RedisService) PublishTradingAlert(alert domain.TradingAlert) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("redis service not available")
	}

	return s.PublishToChannel("stock:trading_alerts", alert)
}

// GetConnectionStatus returns the Redis connection status
func (s *RedisService) GetConnectionStatus() bool {
	if s == nil || s.client == nil {
		return false
	}

	_, err := s.client.Ping(s.ctx).Result()
	return err == nil
}

// Close closes the Redis connection
func (s *RedisService) Close() error {
	if s == nil || s.client == nil {
		return nil
	}

	return s.client.Close()
} 