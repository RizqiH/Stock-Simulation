package domain

import (
    "time"
    "fmt"
)

// CommissionType represents different commission calculation methods
type CommissionType string

const (
    CommissionTypeFlat       CommissionType = "FLAT"       // Fixed amount per trade
    CommissionTypePercentage CommissionType = "PERCENTAGE" // Percentage of trade value
    CommissionTypePerShare   CommissionType = "PER_SHARE"  // Amount per share
    CommissionTypeTiered     CommissionType = "TIERED"     // Different rates based on volume
)

// CommissionTier represents a tier in the tiered commission structure
type CommissionTier struct {
    ID          int     `json:"id" db:"id"`
    MinVolume   float64 `json:"min_volume" db:"min_volume"`     // Minimum monthly volume
    MaxVolume   *float64 `json:"max_volume,omitempty" db:"max_volume"` // Maximum monthly volume (null for unlimited)
    Rate        float64 `json:"rate" db:"rate"`                 // Commission rate for this tier
    MinFee      float64 `json:"min_fee" db:"min_fee"`           // Minimum fee per trade
    MaxFee      *float64 `json:"max_fee,omitempty" db:"max_fee"` // Maximum fee per trade
}

// CommissionStructure represents the commission configuration
type CommissionStructure struct {
    ID                 int              `json:"id" db:"id"`
    Name               string           `json:"name" db:"name"`
    Type               CommissionType   `json:"type" db:"type"`
    BaseRate           float64          `json:"base_rate" db:"base_rate"`           // Base commission rate
    MinimumFee         float64          `json:"minimum_fee" db:"minimum_fee"`       // Minimum commission per trade
    MaximumFee         *float64         `json:"maximum_fee,omitempty" db:"maximum_fee"` // Maximum commission per trade
    
    // Additional fees
    RegulatoryFee      float64          `json:"regulatory_fee" db:"regulatory_fee"`       // SEC/FINRA fees
    ClearingFee        float64          `json:"clearing_fee" db:"clearing_fee"`           // Clearing house fees
    PlatformFee        float64          `json:"platform_fee" db:"platform_fee"`           // Platform usage fee
    
    // Market data fees
    MarketDataFee      float64          `json:"market_data_fee" db:"market_data_fee"`     // Monthly market data fee
    
    // Account fees
    InactivityFee      float64          `json:"inactivity_fee" db:"inactivity_fee"`       // Monthly fee for inactive accounts
    InactivityPeriod   int              `json:"inactivity_period" db:"inactivity_period"` // Days of inactivity before fee
    
    // Special rates
    OptionsRate        float64          `json:"options_rate" db:"options_rate"`           // Commission for options
    ForexRate          float64          `json:"forex_rate" db:"forex_rate"`               // Commission for forex
    CryptoRate         float64          `json:"crypto_rate" db:"crypto_rate"`             // Commission for crypto
    
    // Tiered structure
    Tiers              []CommissionTier `json:"tiers,omitempty"`
    
    IsActive           bool             `json:"is_active" db:"is_active"`
    CreatedAt          time.Time        `json:"created_at" db:"created_at"`
    UpdatedAt          time.Time        `json:"updated_at" db:"updated_at"`
}

// UserCommissionProfile represents a user's commission profile
type UserCommissionProfile struct {
    ID                    int     `json:"id" db:"id"`
    UserID                int     `json:"user_id" db:"user_id"`
    CommissionStructureID int     `json:"commission_structure_id" db:"commission_structure_id"`
    MonthlyVolume         float64 `json:"monthly_volume" db:"monthly_volume"`
    YearlyVolume          float64 `json:"yearly_volume" db:"yearly_volume"`
    TotalTrades           int     `json:"total_trades" db:"total_trades"`
    VIPLevel              int     `json:"vip_level" db:"vip_level"`              // 0=Regular, 1=Bronze, 2=Silver, 3=Gold, 4=Platinum
    LastTradeDate         *time.Time `json:"last_trade_date,omitempty" db:"last_trade_date"`
    CreatedAt             time.Time `json:"created_at" db:"created_at"`
    UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

// CommissionCalculation represents the result of commission calculation
type CommissionCalculation struct {
    BaseCommission    float64 `json:"base_commission"`
    RegulatoryFees    float64 `json:"regulatory_fees"`
    ClearingFees      float64 `json:"clearing_fees"`
    PlatformFees      float64 `json:"platform_fees"`
    TotalCommission   float64 `json:"total_commission"`
    EffectiveRate     float64 `json:"effective_rate"`     // Total commission as percentage of trade value
    TierApplied       *CommissionTier `json:"tier_applied,omitempty"`
}

// Slippage represents market slippage simulation
type Slippage struct {
    BaseSlippage      float64 `json:"base_slippage"`      // Base slippage percentage
    VolumeImpact      float64 `json:"volume_impact"`      // Additional slippage based on volume
    MarketImpact      float64 `json:"market_impact"`      // Market conditions impact
    TotalSlippage     float64 `json:"total_slippage"`     // Total slippage applied
    SlippageAmount    float64 `json:"slippage_amount"`    // Dollar amount of slippage
}

// CommissionService interface defines commission calculation methods
type CommissionCalculator interface {
    CalculateCommission(userID int, tradeValue float64, orderType OrderType, assetType string) (*CommissionCalculation, error)
    CalculateSlippage(symbol string, quantity int, orderType OrderType, marketConditions string) (*Slippage, error)
    GetUserCommissionProfile(userID int) (*UserCommissionProfile, error)
    UpdateUserVolume(userID int, tradeValue float64) error
}

// Commission calculation methods
func (cs *CommissionStructure) CalculateBasicCommission(tradeValue float64) float64 {
    var commission float64
    
    switch cs.Type {
    case CommissionTypeFlat:
        commission = cs.BaseRate
    case CommissionTypePercentage:
        commission = tradeValue * (cs.BaseRate / 100)
    case CommissionTypePerShare:
        // This would need quantity parameter, simplified for now
        commission = cs.BaseRate
    }
    
    // Apply minimum fee
    if commission < cs.MinimumFee {
        commission = cs.MinimumFee
    }
    
    // Apply maximum fee if set
    if cs.MaximumFee != nil && commission > *cs.MaximumFee {
        commission = *cs.MaximumFee
    }
    
    return commission
}

func (cs *CommissionStructure) CalculateTieredCommission(tradeValue, monthlyVolume float64) float64 {
    if len(cs.Tiers) == 0 {
        return cs.CalculateBasicCommission(tradeValue)
    }
    
    // Find applicable tier
    var applicableTier *CommissionTier
    for i := range cs.Tiers {
        tier := &cs.Tiers[i]
        if monthlyVolume >= tier.MinVolume {
            if tier.MaxVolume == nil || monthlyVolume <= *tier.MaxVolume {
                applicableTier = tier
                break
            }
        }
    }
    
    if applicableTier == nil {
        return cs.CalculateBasicCommission(tradeValue)
    }
    
    commission := tradeValue * (applicableTier.Rate / 100)
    
    // Apply tier-specific minimum and maximum
    if commission < applicableTier.MinFee {
        commission = applicableTier.MinFee
    }
    
    if applicableTier.MaxFee != nil && commission > *applicableTier.MaxFee {
        commission = *applicableTier.MaxFee
    }
    
    return commission
}

func (cs *CommissionStructure) GetSpecialRate(assetType string) float64 {
    switch assetType {
    case "options":
        return cs.OptionsRate
    case "forex":
        return cs.ForexRate
    case "crypto":
        return cs.CryptoRate
    default:
        return cs.BaseRate
    }
}

// Validation methods
func (cs *CommissionStructure) Validate() error {
    if cs.Name == "" {
        return fmt.Errorf("commission structure name is required")
    }
    
    if cs.BaseRate < 0 {
        return fmt.Errorf("base rate cannot be negative")
    }
    
    if cs.MinimumFee < 0 {
        return fmt.Errorf("minimum fee cannot be negative")
    }
    
    if cs.MaximumFee != nil && *cs.MaximumFee < cs.MinimumFee {
        return fmt.Errorf("maximum fee cannot be less than minimum fee")
    }
    
    // Validate tiers
    for i, tier := range cs.Tiers {
        if tier.MinVolume < 0 {
            return fmt.Errorf("tier %d: minimum volume cannot be negative", i)
        }
        
        if tier.MaxVolume != nil && *tier.MaxVolume <= tier.MinVolume {
            return fmt.Errorf("tier %d: maximum volume must be greater than minimum volume", i)
        }
        
        if tier.Rate < 0 {
            return fmt.Errorf("tier %d: rate cannot be negative", i)
        }
    }
    
    return nil
}

// Helper methods for user profile
func (ucp *UserCommissionProfile) IsVIP() bool {
    return ucp.VIPLevel > 0
}

func (ucp *UserCommissionProfile) GetVIPLevelName() string {
    switch ucp.VIPLevel {
    case 1:
        return "Bronze"
    case 2:
        return "Silver"
    case 3:
        return "Gold"
    case 4:
        return "Platinum"
    default:
        return "Regular"
    }
}

func (ucp *UserCommissionProfile) IsActive() bool {
    if ucp.LastTradeDate == nil {
        return false
    }
    
    daysSinceLastTrade := time.Since(*ucp.LastTradeDate).Hours() / 24
    return daysSinceLastTrade <= 30 // Consider active if traded in last 30 days
}

// Market hours and conditions
type MarketCondition string

const (
    MarketConditionNormal     MarketCondition = "NORMAL"
    MarketConditionVolatile   MarketCondition = "VOLATILE"
    MarketConditionIlliquid   MarketCondition = "ILLIQUID"
    MarketConditionHighVolume MarketCondition = "HIGH_VOLUME"
)

// Slippage calculation helpers
func CalculateMarketImpactSlippage(quantity int, averageVolume int64, baseSlippage float64) float64 {
    if averageVolume == 0 {
        return baseSlippage * 2 // High slippage for unknown volume
    }
    
    volumeRatio := float64(quantity) / float64(averageVolume)
    
    // Increase slippage based on order size relative to average volume
    if volumeRatio > 0.1 {
        return baseSlippage * (1 + volumeRatio)
    }
    
    return baseSlippage
} 