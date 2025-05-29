package exchangeCredentials

import (
	"context"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/internal/database/models"
	"gorm.io/gorm"
	"time"
)

type IExchangeCredentialRepository interface {
	Create(ctx context.Context, cred *models.ExchangeCredential) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ExchangeCredential, error)
	GetByUserAndExchange(ctx context.Context, userID, exchangeID uuid.UUID) (*models.ExchangeCredential, error)
	Update(ctx context.Context, cred *models.ExchangeCredential) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastUsed(ctx context.Context, id uuid.UUID) error
}

type ExchangeCredentialRepository struct {
	Db *gorm.DB
}

func NewExchangeCredentialRepository(db *gorm.DB) *ExchangeCredentialRepository {
	return &ExchangeCredentialRepository{Db: db}
}

func (r *ExchangeCredentialRepository) Create(ctx context.Context, cred *models.ExchangeCredential) error {
	return r.Db.WithContext(ctx).Create(cred).Error
}

func (r *ExchangeCredentialRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ExchangeCredential, error) {
	var cred models.ExchangeCredential
	err := r.Db.WithContext(ctx).Preload("User").Preload("Exchange").First(&cred, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

func (r *ExchangeCredentialRepository) GetByUserAndExchange(ctx context.Context, userID, exchangeID uuid.UUID) (
	*models.ExchangeCredential, error) {
	var creds *models.ExchangeCredential
	err := r.Db.WithContext(ctx).
		Preload("Exchange").
		Where("user_id = ? AND exchange_id = ? AND is_active = ?", userID, exchangeID, true).
		Order("created_at DESC").
		Find(&creds).Error
	return creds, err
}

func (r *ExchangeCredentialRepository) Update(ctx context.Context, cred *models.ExchangeCredential) error {
	return r.Db.WithContext(ctx).Save(cred).Error
}

func (r *ExchangeCredentialRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.Db.WithContext(ctx).Delete(&models.ExchangeCredential{}, id).Error
}

func (r *ExchangeCredentialRepository) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.Db.WithContext(ctx).Model(&models.ExchangeCredential{}).
		Where("id = ?", id).
		Update("last_used", now).Error
}
