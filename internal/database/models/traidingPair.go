package models

import "github.com/google/uuid"

type TradingPair struct {
	BaseModel
	ExchangeID  uuid.UUID `gorm:"type:uuid;not null" json:"exchange_id"`
	Symbol      string    `gorm:"size:20;not null;index:idx_trading_pairs_symbol" json:"symbol"`
	BaseAsset   string    `gorm:"size:10;not null" json:"base_asset"`
	QuoteAsset  string    `gorm:"size:10;not null" json:"quote_asset"`
	MinQuantity *float64  `gorm:"type:decimal(20,8)" json:"min_quantity,omitempty"`
	MaxQuantity *float64  `gorm:"type:decimal(20,8)" json:"max_quantity,omitempty"`
	StepSize    *float64  `gorm:"type:decimal(20,8)" json:"step_size,omitempty"`
	TickSize    *float64  `gorm:"type:decimal(20,8)" json:"tick_size,omitempty"`
	IsActive    bool      `gorm:"not null;default:true" json:"is_active"`

	// Relationships
	Exchange           Exchange            `gorm:"foreignKey:ExchangeID;constraint:OnDelete:CASCADE" json:"exchange,omitempty"`
	OrderHistories     []OrderHistory      `gorm:"foreignKey:TradingPairID;constraint:OnDelete:RESTRICT" json:"order_histories,omitempty"`
	OrderBookSnapshots []OrderBookSnapshot `gorm:"foreignKey:TradingPairID;constraint:OnDelete:CASCADE" json:"order_book_snapshots,omitempty"`

	// Unique constraint
	_ struct{} `gorm:"uniqueIndex:ux_trading_pairs_exchange_symbol_active,where:deleted_at IS NULL"`
}
