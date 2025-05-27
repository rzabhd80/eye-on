package order

import (
	"context"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/internal/database/models"
	"gorm.io/gorm"
)

type GormOrderHistoryRepository struct {
	db *gorm.DB
}

func NewGormOrderHistoryRepository(db *gorm.DB) *GormOrderHistoryRepository {
	return &GormOrderHistoryRepository{db: db}
}

func (r *GormOrderHistoryRepository) Create(ctx context.Context, order *models.OrderHistory) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *GormOrderHistoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.OrderHistory, error) {
	var order models.OrderHistory
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Exchange").
		Preload("TradingPair").
		Preload("OrderEvents").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *GormOrderHistoryRepository) GetByClientOrderID(ctx context.Context, clientOrderID string) (*models.OrderHistory, error) {
	var order models.OrderHistory
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Exchange").
		Preload("TradingPair").
		First(&order, "client_order_id = ?", clientOrderID).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *GormOrderHistoryRepository) GetByExchangeOrderID(ctx context.Context, exchangeOrderID string) (*models.OrderHistory, error) {
	var order models.OrderHistory
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Exchange").
		Preload("TradingPair").
		First(&order, "exchange_order_id = ?", exchangeOrderID).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *GormOrderHistoryRepository) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.OrderHistory, error) {
	var orders []models.OrderHistory
	err := r.db.WithContext(ctx).
		Preload("Exchange").
		Preload("TradingPair").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, err
}
