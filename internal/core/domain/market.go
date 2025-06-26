package domain

import (
    "time"
    "fmt"
)

// MarketType represents different market types
type MarketType string

const (
    MarketTypeStock  MarketType = "STOCK"
    MarketTypeForex  MarketType = "FOREX"
    MarketTypeCrypto MarketType = "CRYPTO"
    MarketTypeOption MarketType = "OPTION"
    MarketTypeFuture MarketType = "FUTURE"
)

// Market represents a trading market
type Market struct {
    ID                  int        `json:"id" db:"id"`
    Code                string     `json:"code" db:"code"`                 // NYSE, NASDAQ, LSE, etc.
    Name                string     `json:"name" db:"name"`
    Type                MarketType `json:"type" db:"type"`
    TimeZone            string     `json:"timezone" db:"timezone"`         // America/New_York, Europe/London, etc.
    Currency            string     `json:"currency" db:"currency"`         // USD, EUR, GBP, etc.
    IsActive            bool       `json:"is_active" db:"is_active"`
    CreatedAt           time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}

// TradingSessionType represents different trading session types
type TradingSessionType string

const (
    SessionTypePreMarket  TradingSessionType = "PRE_MARKET"
    SessionTypeRegular    TradingSessionType = "REGULAR"
    SessionTypeAfterHours TradingSessionType = "AFTER_HOURS"
    SessionTypeOvernight  TradingSessionType = "OVERNIGHT"
)

// TradingSession represents a trading session for a market
type TradingSession struct {
    ID          int                `json:"id" db:"id"`
    MarketID    int                `json:"market_id" db:"market_id"`
    Type        TradingSessionType `json:"type" db:"type"`
    StartTime   string             `json:"start_time" db:"start_time"`   // HH:MM format
    EndTime     string             `json:"end_time" db:"end_time"`       // HH:MM format
    DaysOfWeek  string             `json:"days_of_week" db:"days_of_week"` // JSON array: [1,2,3,4,5] for Mon-Fri
    IsActive    bool               `json:"is_active" db:"is_active"`
    CreatedAt   time.Time          `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time          `json:"updated_at" db:"updated_at"`
}

// MarketHoliday represents market holidays
type MarketHoliday struct {
    ID          int       `json:"id" db:"id"`
    MarketID    int       `json:"market_id" db:"market_id"`
    Date        time.Time `json:"date" db:"date"`
    Name        string    `json:"name" db:"name"`
    Type        string    `json:"type" db:"type"`          // FULL_CLOSE, EARLY_CLOSE
    EarlyCloseTime *string `json:"early_close_time,omitempty" db:"early_close_time"` // HH:MM format
    IsRecurring bool      `json:"is_recurring" db:"is_recurring"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// MarketStatus represents the current status of a market
type MarketStatus struct {
    MarketCode        string             `json:"market_code"`
    IsOpen            bool               `json:"is_open"`
    CurrentSession    *TradingSessionType `json:"current_session,omitempty"`
    NextOpenTime      *time.Time         `json:"next_open_time,omitempty"`
    NextCloseTime     *time.Time         `json:"next_close_time,omitempty"`
    TimeToOpen        *time.Duration     `json:"time_to_open,omitempty"`
    TimeToClose       *time.Duration     `json:"time_to_close,omitempty"`
    LocalTime         time.Time          `json:"local_time"`
    Message           string             `json:"message"`
    LastUpdated       time.Time          `json:"last_updated"`
}

// MarketCalendar represents market calendar information
type MarketCalendar struct {
    Date           time.Time           `json:"date"`
    MarketCode     string              `json:"market_code"`
    IsMarketDay    bool                `json:"is_market_day"`
    Sessions       []SessionSchedule   `json:"sessions"`
    Holidays       []MarketHoliday     `json:"holidays"`
    SpecialHours   *SpecialHours       `json:"special_hours,omitempty"`
}

// SessionSchedule represents the schedule for a specific session
type SessionSchedule struct {
    Type      TradingSessionType `json:"type"`
    StartTime time.Time          `json:"start_time"`
    EndTime   time.Time          `json:"end_time"`
    IsActive  bool               `json:"is_active"`
}

// SpecialHours represents special trading hours (early close, etc.)
type SpecialHours struct {
    Reason    string    `json:"reason"`
    CloseTime time.Time `json:"close_time"`
}

// MarketDataPermission represents user's market data permissions
type MarketDataPermission struct {
    ID               int       `json:"id" db:"id"`
    UserID           int       `json:"user_id" db:"user_id"`
    MarketCode       string    `json:"market_code" db:"market_code"`
    DataType         string    `json:"data_type" db:"data_type"`         // REAL_TIME, DELAYED, SNAPSHOT
    PermissionLevel  string    `json:"permission_level" db:"permission_level"` // BASIC, PREMIUM, PROFESSIONAL
    ExpiresAt        *time.Time `json:"expires_at,omitempty" db:"expires_at"`
    IsActive         bool      `json:"is_active" db:"is_active"`
    CreatedAt        time.Time `json:"created_at" db:"created_at"`
    UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// MarketConditions represents current market conditions
type MarketConditions struct {
    MarketCode    string          `json:"market_code"`
    Volatility    float64         `json:"volatility"`        // VIX-like volatility index
    Volume        int64           `json:"volume"`            // Current day volume
    AverageVolume int64           `json:"average_volume"`    // 30-day average volume
    Sentiment     string          `json:"sentiment"`         // BULLISH, BEARISH, NEUTRAL
    Trend         string          `json:"trend"`             // UP, DOWN, SIDEWAYS
    Liquidity     string          `json:"liquidity"`         // HIGH, MEDIUM, LOW
    LastUpdated   time.Time       `json:"last_updated"`
}

// TradingRestriction represents trading restrictions
type TradingRestriction struct {
    ID              int       `json:"id" db:"id"`
    MarketCode      string    `json:"market_code" db:"market_code"`
    Symbol          *string   `json:"symbol,omitempty" db:"symbol"`         // Null for market-wide restrictions
    RestrictionType string    `json:"restriction_type" db:"restriction_type"` // HALT, SUSPENSION, LIMIT_UP_DOWN
    Reason          string    `json:"reason" db:"reason"`
    StartTime       time.Time `json:"start_time" db:"start_time"`
    EndTime         *time.Time `json:"end_time,omitempty" db:"end_time"`
    IsActive        bool      `json:"is_active" db:"is_active"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Market hours calculation methods
func (m *Market) GetTimeZone() (*time.Location, error) {
    return time.LoadLocation(m.TimeZone)
}

func (m *Market) GetCurrentTime() (time.Time, error) {
    tz, err := m.GetTimeZone()
    if err != nil {
        return time.Time{}, err
    }
    return time.Now().In(tz), nil
}

func (ts *TradingSession) ParseStartTime(date time.Time, timezone *time.Location) (time.Time, error) {
    startTime, err := time.Parse("15:04", ts.StartTime)
    if err != nil {
        return time.Time{}, err
    }
    
    return time.Date(
        date.Year(), date.Month(), date.Day(),
        startTime.Hour(), startTime.Minute(), 0, 0,
        timezone,
    ), nil
}

func (ts *TradingSession) ParseEndTime(date time.Time, timezone *time.Location) (time.Time, error) {
    endTime, err := time.Parse("15:04", ts.EndTime)
    if err != nil {
        return time.Time{}, err
    }
    
    // Handle sessions that end after midnight
    endDate := date
    if endTime.Hour() < 12 && ts.StartTime > "12:00" {
        endDate = date.AddDate(0, 0, 1)
    }
    
    return time.Date(
        endDate.Year(), endDate.Month(), endDate.Day(),
        endTime.Hour(), endTime.Minute(), 0, 0,
        timezone,
    ), nil
}

func (ts *TradingSession) IsActiveOnDay(weekday time.Weekday) bool {
    // Simple implementation - assumes DaysOfWeek contains weekday numbers
    // In real implementation, this would parse the JSON array
    switch weekday {
    case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
        return true
    default:
        return ts.Type == SessionTypeOvernight || ts.Type == SessionTypeAfterHours
    }
}

// Market status calculation methods
func (ms *MarketStatus) CalculateTimeToOpen() {
    if ms.NextOpenTime != nil {
        duration := time.Until(*ms.NextOpenTime)
        if duration > 0 {
            ms.TimeToOpen = &duration
        }
    }
}

func (ms *MarketStatus) CalculateTimeToClose() {
    if ms.NextCloseTime != nil && ms.IsOpen {
        duration := time.Until(*ms.NextCloseTime)
        if duration > 0 {
            ms.TimeToClose = &duration
        }
    }
}

func (ms *MarketStatus) UpdateMessage() {
    if ms.IsOpen {
        if ms.CurrentSession != nil {
            switch *ms.CurrentSession {
            case SessionTypePreMarket:
                ms.Message = "Pre-market trading is open"
            case SessionTypeRegular:
                ms.Message = "Market is open for regular trading"
            case SessionTypeAfterHours:
                ms.Message = "After-hours trading is open"
            case SessionTypeOvernight:
                ms.Message = "Overnight trading is open"
            }
        } else {
            ms.Message = "Market is open"
        }
    } else {
        if ms.TimeToOpen != nil {
            hours := int(ms.TimeToOpen.Hours())
            minutes := int(ms.TimeToOpen.Minutes()) % 60
            ms.Message = fmt.Sprintf("Market closed. Opens in %dh %dm", hours, minutes)
        } else {
            ms.Message = "Market is closed"
        }
    }
}

// Validation methods
func (m *Market) Validate() error {
    if m.Code == "" {
        return fmt.Errorf("market code is required")
    }
    
    if m.Name == "" {
        return fmt.Errorf("market name is required")
    }
    
    if m.TimeZone == "" {
        return fmt.Errorf("timezone is required")
    }
    
    // Validate timezone
    _, err := time.LoadLocation(m.TimeZone)
    if err != nil {
        return fmt.Errorf("invalid timezone: %v", err)
    }
    
    return nil
}

func (ts *TradingSession) Validate() error {
    if ts.StartTime == "" || ts.EndTime == "" {
        return fmt.Errorf("start time and end time are required")
    }
    
    // Validate time format
    _, err := time.Parse("15:04", ts.StartTime)
    if err != nil {
        return fmt.Errorf("invalid start time format: %v", err)
    }
    
    _, err = time.Parse("15:04", ts.EndTime)
    if err != nil {
        return fmt.Errorf("invalid end time format: %v", err)
    }
    
    return nil
}

// Market service interface
type MarketService interface {
    GetMarketStatus(marketCode string) (*MarketStatus, error)
    GetMarketCalendar(marketCode string, date time.Time) (*MarketCalendar, error)
    IsMarketOpen(marketCode string) (bool, error)
    GetNextMarketOpen(marketCode string) (*time.Time, error)
    GetNextMarketClose(marketCode string) (*time.Time, error)
    GetMarketConditions(marketCode string) (*MarketConditions, error)
    GetTradingRestrictions(marketCode string, symbol *string) ([]TradingRestriction, error)
    CanTrade(marketCode string, symbol string, orderType OrderType) (bool, string, error)
}

// Helper functions for common market operations
func IsWeekend(t time.Time) bool {
    weekday := t.Weekday()
    return weekday == time.Saturday || weekday == time.Sunday
}

func IsUSHoliday(t time.Time) bool {
    // Simplified US holiday check - in real implementation, this would check against a database
    
    // New Year's Day
    if t.Month() == time.January && t.Day() == 1 {
        return true
    }
    
    // Independence Day
    if t.Month() == time.July && t.Day() == 4 {
        return true
    }
    
    // Christmas Day
    if t.Month() == time.December && t.Day() == 25 {
        return true
    }
    
    // Add more holidays as needed
    return false
}

func GetMarketTimeZone(marketCode string) string {
    switch marketCode {
    case "NYSE", "NASDAQ":
        return "America/New_York"
    case "LSE":
        return "Europe/London"
    case "TSE":
        return "Asia/Tokyo"
    case "HKEX":
        return "Asia/Hong_Kong"
    default:
        return "UTC"
    }
} 