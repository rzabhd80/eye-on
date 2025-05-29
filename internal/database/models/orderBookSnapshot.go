package models

import (
	"github.com/google/uuid"
	"time"
)

type OrderBookSnapshot struct {
	BaseModel
	ExchangeID    uuid.UUID `gorm:"type:uuid;not null" json:"exchange_id"`
	TradingPairID uuid.UUID `gorm:"type:uuid;not null" json:"trading_pair_id"`
	Symbol        string    `gorm:"size:20;not null" json:"symbol"`  // For faster queries without joins
	Bids          JSONB     `gorm:"type:jsonb;not null" json:"bids"` // [[price, qty], ...]
	Asks          JSONB     `gorm:"type:jsonb;not null" json:"asks"`
	SnapshotTime  time.Time `gorm:"not null;default:now()" json:"snapshot_time"`

	// Relationships
	Exchange    Exchange    `gorm:"foreignKey:ExchangeID;constraint:OnDelete:CASCADE" json:"exchange,omitempty"`
	TradingPair TradingPair `gorm:"foreignKey:TradingPairID;constraint:OnDelete:CASCADE" json:"trading_pair,omitempty"`

	_ struct{} `gorm:"index:idx_ob_snapshots_exchange_pair_time,composite:exchange_id,trading_pair_id,snapshot_time"`
	_ struct{} `gorm:"index:idx_ob_snapshots_exchange_symbol_time,composite:exchange_id,symbol,snapshot_time"`
	_ struct{} `gorm:"uniqueIndex:ux_ob_snapshots_exchange_pair_time,composite:exchange_id,trading_pair_id,snapshot_time"`
}
