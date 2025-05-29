package models

import (
	"github.com/google/uuid"
	"time"
)

type BalanceSnapshot struct {
	BaseModel
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ExchangeID   uuid.UUID `gorm:"type:uuid;not null" json:"exchange_id"`
	Currency     string    `gorm:"size:10;not null" json:"currency"`
	Total        float64   `gorm:"type:numeric(30,10);not null" json:"total"`
	Available    float64   `gorm:"type:numeric(30,10);not null" json:"available"`
	SnapshotTime time.Time `gorm:"not null;default:now()" json:"snapshot_time"`

	// Relationships
	User     User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Exchange Exchange `gorm:"foreignKey:ExchangeID;constraint:OnDelete:CASCADE" json:"exchange,omitempty"`

	// Composite index
	_ struct{} `gorm:"index:idx_balance_snapshots_user_time,composite:user_id,snapshot_time"`
}
