package orderBook

import (
	"context"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/internal/database/models"
	"gorm.io/gorm"
	"time"
)

type OrderBookRepository interface {
	Create(ctx context.Context, snapshot *models.OrderBookSnapshot) error
	GetLatestByTradingPair(ctx context.Context, tradingPairID uuid.UUID) (*models.OrderBookSnapshot, error)
	GetHistory(ctx context.Context, tradingPairID uuid.UUID, limit int) ([]models.OrderBookSnapshot, error)
	DeleteOldSnapshots(ctx context.Context, olderThan time.Time) error
}

type GormOrderBookSnapshotRepository struct {
	db *gorm.DB
}

func NewGormOrderBookSnapshotRepository(db *gorm.DB) *GormOrderBookSnapshotRepository {
	return &GormOrderBookSnapshotRepository{db: db}
}

func (r *GormOrderBookSnapshotRepository) Create(ctx context.Context, snapshot *models.OrderBookSnapshot) error {
	return r.db.WithContext(ctx).Create(snapshot).Error
}

func (r *GormOrderBookSnapshotRepository) GetLatestByTradingPair(ctx context.Context, tradingPairID uuid.UUID) (
	*models.OrderBookSnapshot, error) {
	var snapshot models.OrderBookSnapshot
	err := r.db.WithContext(ctx).
		Preload("TradingPair").
		Where("trading_pair_id = ?", tradingPairID).
		Order("snapshot_time DESC").
		First(&snapshot).Error
	if err != nil {
		return nil, err
	}
	return &snapshot, nil
}

func (r *GormOrderBookSnapshotRepository) GetHistory(ctx context.Context, tradingPairID uuid.UUID, limit int) (
	[]models.OrderBookSnapshot, error) {
	var snapshots []models.OrderBookSnapshot
	err := r.db.WithContext(ctx).
		Where("trading_pair_id = ?", tradingPairID).
		Order("snapshot_time DESC").
		Limit(limit).
		Find(&snapshots).Error
	return snapshots, err
}

func (r *GormOrderBookSnapshotRepository) DeleteOldSnapshots(ctx context.Context, olderThan time.Time) error {
	return r.db.WithContext(ctx).Where("snapshot_time < ?", olderThan).Delete(&models.OrderBookSnapshot{}).Error
}
