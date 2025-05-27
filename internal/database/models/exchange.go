package models

type Exchange struct {
	BaseModel
	Name        string `gorm:"size:50;not null;uniqueIndex:ux_exchanges_name_active,where:deleted_at IS NULL" json:"name"`
	DisplayName string `gorm:"size:100;not null" json:"display_name"`
	BaseURL     string `gorm:"size:255;not null" json:"base_url"`
	IsActive    bool   `gorm:"not null;default:true" json:"is_active"`
	RateLimit   int    `gorm:"not null;default:1000" json:"rate_limit"`
	Features    JSONB  `gorm:"type:jsonb" json:"features"`

	// Relationships
	ExchangeCredentials []ExchangeCredential `gorm:"foreignKey:ExchangeID;constraint:OnDelete:CASCADE" json:"exchange_credentials,omitempty"`
	TradingPairs        []TradingPair        `gorm:"foreignKey:ExchangeID;constraint:OnDelete:CASCADE" json:"trading_pairs,omitempty"`
	OrderHistories      []OrderHistory       `gorm:"foreignKey:ExchangeID;constraint:OnDelete:CASCADE" json:"order_histories,omitempty"`
	BalanceSnapshots    []BalanceSnapshot    `gorm:"foreignKey:ExchangeID;constraint:OnDelete:CASCADE" json:"balance_snapshots,omitempty"`
}
