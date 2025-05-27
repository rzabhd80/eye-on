package exchangeCredentials

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type IExchangeCredentialRepository interface {
	Create(ctx context.Context, cred *ExchangeCredential) error
	GetByID(ctx context.Context, id uuid.UUID) (*ExchangeCredential, error)
	GetByUserAndExchange(ctx context.Context, userID, exchangeID uuid.UUID) ([]ExchangeCredential, error)
	Update(ctx context.Context, cred *ExchangeCredential) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastUsed(ctx context.Context, id uuid.UUID) error
}

type ExchangeCredentialRepository struct {
	db *gorm.DB
}

func NewGormExchangeCredentialRepository(db *gorm.DB) IExchangeCredentialRepository {
	return &ExchangeCredentialRepository{db: db}
}

func (r *ExchangeCredentialRepository) Create(ctx context.Context, cred *ExchangeCredential) error {
	return r.db.WithContext(ctx).Create(cred).Error
}

func (r *ExchangeCredentialRepository) GetByID(ctx context.Context, id uuid.UUID) (*ExchangeCredential, error) {
	var cred ExchangeCredential
	err := r.db.WithContext(ctx).Preload("User").Preload("Exchange").First(&cred, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

func (r *ExchangeCredentialRepository) GetByUserAndExchange(ctx context.Context, userID, exchangeID uuid.UUID) ([]ExchangeCredential, error) {
	var creds []ExchangeCredential
	err := r.db.WithContext(ctx).
		Preload("Exchange").
		Where("user_id = ? AND exchange_id = ?", userID, exchangeID).
		Find(&creds).Error
	return creds, err
}

func (r *ExchangeCredentialRepository) Update(ctx context.Context, cred *ExchangeCredential) error {
	return r.db.WithContext(ctx).Save(cred).Error
}

func (r *ExchangeCredentialRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&ExchangeCredential{}, id).Error
}

func (r *ExchangeCredentialRepository) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&ExchangeCredential{}).
		Where("id = ?", id).
		Update("last_used", now).Error
}
