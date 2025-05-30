package exchange

import (
	"context"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/internal/database/models"
	"gorm.io/gorm"
)

type IExchangeRepository interface {
	Create(ctx context.Context, exchange *models.Exchange) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Exchange, error)
	GetByName(ctx context.Context, name string) (*models.Exchange, error)
	Update(ctx context.Context, exchange *models.Exchange) error
	Delete(ctx context.Context, exchange models.Exchange) error
	List(ctx context.Context, activeOnly bool) ([]models.Exchange, error)
}

type ExchangeRepository struct {
	db *gorm.DB
}

func (r *ExchangeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Exchange, error) {
	var config models.Exchange
	err := r.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).Find(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func NewExchangeRepository(db *gorm.DB) *ExchangeRepository {
	return &ExchangeRepository{db: db}
}

func (r *ExchangeRepository) GetByName(ctx context.Context, name string) (*models.Exchange, error) {
	var config models.Exchange
	err := r.db.WithContext(ctx).Where("name = ? AND is_active = ?", name, true).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *ExchangeRepository) GetAll(ctx context.Context) (models.Exchange, error) {
	var configs models.Exchange
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&configs).Error
	return configs, err
}

func (r *ExchangeRepository) Create(ctx context.Context, config *models.Exchange) error {
	return r.db.WithContext(ctx).Create(&config).Error
}

func (r *ExchangeRepository) Update(ctx context.Context, config *models.Exchange) error {
	return r.db.WithContext(ctx).Save(&config).Error
}

func (r *ExchangeRepository) Delete(ctx context.Context, exchange models.Exchange) error {
	return r.db.WithContext(ctx).Delete(&models.Exchange{}, "id = ?", exchange.ID).Error
}

func (r *ExchangeRepository) List(ctx context.Context, activeOnly bool) ([]models.Exchange, error) {
	var exchanges []models.Exchange
	query := r.db.WithContext(ctx)
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	err := query.Find(&exchanges).Error
	return exchanges, err
}
