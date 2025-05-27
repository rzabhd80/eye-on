package models

import (
	"github.com/google/uuid"
	"time"
)

type OrderEvent struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderHistID  uuid.UUID `gorm:"type:uuid;not null;index:idx_order_events_order_hist_id" json:"order_hist_id"`
	EventType    string    `gorm:"size:30;not null" json:"event_type"` // new, partial_fill, filled, canceled
	FilledQty    float64   `gorm:"type:decimal(20,8);not null" json:"filled_qty"`
	RemainingQty float64   `gorm:"type:decimal(20,8);not null" json:"remaining_qty"`
	EventTime    time.Time `gorm:"not null" json:"event_time"`
	RecordedAt   time.Time `gorm:"not null;default:now()" json:"recorded_at"`

	// Relationships
	OrderHistory OrderHistory `gorm:"foreignKey:OrderHistID;constraint:OnDelete:CASCADE" json:"order_history,omitempty"`
}
