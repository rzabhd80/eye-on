package models

import (
	"github.com/google/uuid"
	"time"
)

type OrderBookSnapshot struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TradingPairID uuid.UUID `gorm:"type:uuid;not null" json:"trading_pair_id"`
	Bids          JSONB     `gorm:"type:jsonb;not null" json:"bids"` // [[price, qty], ...]
	Asks          JSONB     `gorm:"type:jsonb;not null" json:"asks"`
	SnapshotTime  time.Time `gorm:"not null;default:now()" json:"snapshot_time"`

	// Relationships
	TradingPair TradingPair `gorm:"foreignKey:TradingPairID;constraint:OnDelete:CASCADE" json:"trading_pair,omitempty"`

	// Composite index
	_ struct{} `gorm:"index:idx_ob_snapshots_pair_time,composite:trading_pair_id,snapshot_time"`
}
