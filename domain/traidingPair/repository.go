package traidingPair

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITradingPairRepository interface {
	Create(ctx context.Context, pair *TradingPair) error
	GetByID(ctx context.Context, id uuid.UUID) (*TradingPair, error)
	GetByExchangeAndSymbol(ctx context.Context, exchangeID uuid.UUID, symbol string) (*TradingPair, error)
	GetByExchange(ctx context.Context, exchangeID uuid.UUID, activeOnly bool) ([]TradingPair, error)
	Update(ctx context.Context, pair *TradingPair) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetSymbolsList(ctx context.Context, exchangeID uuid.UUID, activeOnly bool, symbols []string) ([]TradingPair, error)
}

type TradingPairRepository struct {
	DB *gorm.DB
}

func NewTradingPairRepository(db *gorm.DB) ITradingPairRepository {
	return &TradingPairRepository{DB: db}
}

func (r *TradingPairRepository) Create(ctx context.Context, pair *TradingPair) error {
	return r.DB.WithContext(ctx).Create(pair).Error
}

func (r *TradingPairRepository) BulkCreate(ctx context.Context, pair *[]TradingPair) error {
	return r.DB.WithContext(ctx).Create(&pair).Error
}
func (r *TradingPairRepository) GetByID(ctx context.Context, id uuid.UUID) (*TradingPair, error) {
	var pair TradingPair
	err := r.DB.WithContext(ctx).Preload("Exchange").First(&pair, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &pair, nil
}

func (r *TradingPairRepository) GetByExchangeAndSymbol(ctx context.Context, exchangeID uuid.UUID, symbol string) (*TradingPair, error) {
	var pair TradingPair
	err := r.DB.WithContext(ctx).
		Preload("Exchange").
		Where("exchange_id = ? AND symbol = ?", exchangeID, symbol).
		First(&pair).Error
	if err != nil {
		return nil, err
	}
	return &pair, nil
}

func (r *TradingPairRepository) GetByExchange(ctx context.Context, exchangeID uuid.UUID, activeOnly bool) ([]TradingPair, error) {
	var pairs []TradingPair
	query := r.DB.WithContext(ctx).Where("exchange_id = ?", exchangeID)
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	err := query.Find(&pairs).Error
	return pairs, err
}

func (r *TradingPairRepository) GetSymbolsList(ctx context.Context, exchangeID uuid.UUID, activeOnly bool, symbols []string) (
	[]TradingPair, error) {
	var pairs []TradingPair
	query := r.DB.WithContext(ctx).Where("exchange_id = ? AND Symbol IN ?", exchangeID, symbols)
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	err := query.Find(&pairs).Error
	if err != nil {
		return nil, err
	}
	return pairs, err
}
func (r *TradingPairRepository) Update(ctx context.Context, pair *TradingPair) error {
	return r.DB.WithContext(ctx).Save(pair).Error
}

func (r *TradingPairRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&TradingPair{}, id).Error
}
