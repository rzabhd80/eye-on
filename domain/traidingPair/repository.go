package traidingPair

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TradingPairRepository interface {
	Create(ctx context.Context, pair *TradingPair) error
	GetByID(ctx context.Context, id uuid.UUID) (*TradingPair, error)
	GetByExchangeAndSymbol(ctx context.Context, exchangeID uuid.UUID, symbol string) (*TradingPair, error)
	GetByExchange(ctx context.Context, exchangeID uuid.UUID, activeOnly bool) ([]TradingPair, error)
	Update(ctx context.Context, pair *TradingPair) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type GormTradingPairRepository struct {
	db *gorm.DB
}

func NewGormTradingPairRepository(db *gorm.DB) *GormTradingPairRepository {
	return &GormTradingPairRepository{db: db}
}

func (r *GormTradingPairRepository) Create(ctx context.Context, pair *TradingPair) error {
	return r.db.WithContext(ctx).Create(pair).Error
}

func (r *GormTradingPairRepository) GetByID(ctx context.Context, id uuid.UUID) (*TradingPair, error) {
	var pair TradingPair
	err := r.db.WithContext(ctx).Preload("Exchange").First(&pair, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &pair, nil
}

func (r *GormTradingPairRepository) GetByExchangeAndSymbol(ctx context.Context, exchangeID uuid.UUID, symbol string) (*TradingPair, error) {
	var pair TradingPair
	err := r.db.WithContext(ctx).
		Preload("Exchange").
		Where("exchange_id = ? AND symbol = ?", exchangeID, symbol).
		First(&pair).Error
	if err != nil {
		return nil, err
	}
	return &pair, nil
}

func (r *GormTradingPairRepository) GetByExchange(ctx context.Context, exchangeID uuid.UUID, activeOnly bool) ([]TradingPair, error) {
	var pairs []TradingPair
	query := r.db.WithContext(ctx).Where("exchange_id = ?", exchangeID)
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	err := query.Find(&pairs).Error
	return pairs, err
}

func (r *GormTradingPairRepository) Update(ctx context.Context, pair *TradingPair) error {
	return r.db.WithContext(ctx).Save(pair).Error
}

func (r *GormTradingPairRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&TradingPair{}, id).Error
}
