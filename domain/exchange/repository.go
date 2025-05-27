package exchange

import (
	"context"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/internal/database/models"
	"gorm.io/gorm"
)

type IExchangeRepository interface {
	Create(ctx context.Context, exchange *Exchange) error
	GetByID(ctx context.Context, id uuid.UUID) (*Exchange, error)
	GetByName(ctx context.Context, name string) (*models.Exchange, error)
	Update(ctx context.Context, exchange *Exchange) error
	Delete(ctx context.Context, exchange Exchange) error
	List(ctx context.Context, activeOnly bool) ([]Exchange, error)
}

type ExchangeRepository struct {
	Db *gorm.DB
}

func (r *ExchangeRepository) GetByID(ctx context.Context, id uuid.UUID) (*Exchange, error) {
	var config Exchange
	err := r.Db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).Find(&config.exchange).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func NewExchangeRepository(db *gorm.DB) IExchangeRepository {
	return &ExchangeRepository{Db: db}
}

func (r *ExchangeRepository) GetByName(ctx context.Context, name string) (*models.Exchange, error) {
	var config Exchange
	err := r.Db.WithContext(ctx).Where("name = ? AND is_active = ?", name, true).First(&config.exchange).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *ExchangeRepository) GetAll(ctx context.Context) (Exchange, error) {
	var configs Exchange
	err := r.Db.WithContext(ctx).Where("is_active = ?", true).Find(&configs.exchange).Error
	return configs, err
}

func (r *ExchangeRepository) Create(ctx context.Context, config *Exchange) error {
	return r.Db.WithContext(ctx).Create(config.exchange).Error
}

func (r *ExchangeRepository) Update(ctx context.Context, config *Exchange) error {
	return r.Db.WithContext(ctx).Save(config.exchange).Error
}

func (r *ExchangeRepository) Delete(ctx context.Context, exchange Exchange) error {
	return r.Db.WithContext(ctx).Delete(&Exchange{}, "id = ?", exchange.exchange.ID).Error
}

func (r *ExchangeRepository) List(ctx context.Context, activeOnly bool) ([]Exchange, error) {
	var exchanges []Exchange
	query := r.Db.WithContext(ctx)
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	err := query.Find(&exchanges).Error
	return exchanges, err
}
