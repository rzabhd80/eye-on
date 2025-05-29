package models

import "github.com/google/uuid"

type OrderHistory struct {
	BaseModel
	UserID               uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ExchangeCredentialID uuid.UUID `gorm:"type:uuid;not null" json:"exchange_credential_id"`
	ExchangeID           uuid.UUID `gorm:"type:uuid;not null" json:"exchange_id"`
	TradingPairID        uuid.UUID `gorm:"type:uuid;not null" json:"trading_pair_id"`
	ClientOrderID        string    `gorm:"size:100;not null;index:idx_order_histories_client_order_id" json:"client_order_id"`
	ExchangeOrderID      string    `gorm:"size:100;not null;index:idx_order_histories_order_id" json:"exchange_order_id"`
	Side                 string    `gorm:"size:10;not null" json:"side"` // buy/sell
	Type                 string    `gorm:"size:10;not null" json:"type"` // limit/market
	Quantity             float64   `gorm:"type:decimal(20,8);not null" json:"quantity"`
	Price                *float64  `gorm:"type:decimal(20,8)" json:"price,omitempty"`
	Status               string    `gorm:"size:20;not null" json:"status"`
	// Relationships
	User               User               `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	ExchangeCredential ExchangeCredential `gorm:"foreignKey:ExchangeCredentialID;constraint:OnDelete:CASCADE" json:"exchange_credential,omitempty"`
	Exchange           Exchange           `gorm:"foreignKey:ExchangeID;constraint:OnDelete:CASCADE" json:"exchange,omitempty"`
	TradingPair        TradingPair        `gorm:"foreignKey:TradingPairID;constraint:OnDelete:RESTRICT" json:"trading_pair,omitempty"`
	OrderEvents        []OrderEvent       `gorm:"foreignKey:OrderHistID;constraint:OnDelete:CASCADE" json:"order_events,omitempty"`

	// Composite index
	_ struct{} `gorm:"index:idx_orders_user_ex_cred_status,composite:user_id,exchange_credential_id,status"`
}
