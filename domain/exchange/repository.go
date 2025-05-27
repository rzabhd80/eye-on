package exchange

import (
	"context"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/domain/exchange/registry"
	"gorm.io/gorm"
)

type ExchangeRepository interface {
	Create(ctx context.Context, exchange *Exchange) error
	GetByID(ctx context.Context, id uuid.UUID) (*Exchange, error)
	GetByName(ctx context.Context, name string) (*Exchange, error)
	Update(ctx context.Context, exchange *Exchange) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, activeOnly bool) ([]Exchange, error)
}

type exchangeRepository struct {
	db *gorm.DB
}

func NewExchangeRepository(db *gorm.DB) ExchangeRepository {
	return &exchangeRepository{db: db}
}

func (r *exchangeRepository) GetByName(ctx context.Context, name string) (*registry.ExchangeConfig, error) {
	var config domain.ExchangeConfig
	err := r.db.WithContext(ctx).Where("name = ? AND is_active = ?", name, true).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *exchangeRepository) GetAll(ctx context.Context) ([]domain.ExchangeConfig, error) {
	var configs []domain.ExchangeConfig
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&configs).Error
	return configs, err
}

func (r *exchangeRepository) Create(ctx context.Context, config *domain.ExchangeConfig) error {
	return r.db.WithContext(ctx).Create(config).Error
}

func (r *exchangeRepository) Update(ctx context.Context, config *domain.ExchangeConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}

func (r *exchangeRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.ExchangeConfig{}, "id = ?", id).Error
}

func (r *exchangeRepository) List(ctx context.Context, activeOnly bool) ([]Exchange, error) {
	var exchanges []Exchange
	query := r.db.WithContext(ctx)
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	err := query.Find(&exchanges).Error
	return exchanges, err
}
