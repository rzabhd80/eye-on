package models

import (
	"github.com/google/uuid"
	"time"
)

type ExchangeCredential struct {
	BaseModel
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index:idx_user_exchange" json:"user_id"`
	ExchangeID uuid.UUID  `gorm:"type:uuid;not null;index:idx_user_exchange" json:"exchange_id"`
	Label      string     `gorm:"size:100;not null;default:'Default'" json:"label"`
	APIKey     string     `gorm:"type:text;not null" json:"api_key"`
	SecretKey  string     `gorm:"type:text;" json:"secret_key"`
	RefreshKey string     `gorm:"type:text;" json:"refresh_key"`
	AccessKey  string     `gorm:"type:text" json:"access_key,omitempty"`
	IsActive   bool       `gorm:"not null;default:true" json:"is_active"`
	IsTestnet  bool       `gorm:"not null;default:false" json:"is_testnet"`
	LastUsed   *time.Time `json:"last_used,omitempty"`

	// Relationships
	User           User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Exchange       Exchange       `gorm:"foreignKey:ExchangeID;constraint:OnDelete:CASCADE" json:"exchange,omitempty"`
	OrderHistories []OrderHistory `gorm:"foreignKey:ExchangeCredentialID;constraint:OnDelete:CASCADE" json:"order_histories,omitempty"`

	// Unique constraint
	_ struct{} `gorm:"uniqueIndex:ux_exchange_credentials_user_exchange_label,unique"`
}
