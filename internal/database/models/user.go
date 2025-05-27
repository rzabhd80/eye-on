package models

type User struct {
	BaseModel
	Username string `gorm:"size:50;not null;uniqueIndex:ux_users_username_active,where:deleted_at IS NULL" json:"username"`
	Email    string `gorm:"size:255;not null;uniqueIndex:ux_users_email_active,where:deleted_at IS NULL" json:"email"`
	IsActive bool   `gorm:"not null;default:true" json:"is_active"`

	ExchangeCredentials []ExchangeCredential `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"exchange_credentials,omitempty"`
	OrderHistories      []OrderHistory       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"order_histories,omitempty"`
	BalanceSnapshots    []BalanceSnapshot    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"balance_snapshots,omitempty"`
}
