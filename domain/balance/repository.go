package balance

import (
	"context"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/internal/database/models"
	"gorm.io/gorm"
	"time"
)

type BalanceSnapshotRepository interface {
	Create(ctx context.Context, snapshot *Balance) error
	GetLatestByUserAndExchange(ctx context.Context, userID, exchangeID uuid.UUID) ([]*models.BalanceSnapshot, error)
	GetHistory(ctx context.Context, userID, exchangeID uuid.UUID, currency string, limit int) ([]*models.BalanceSnapshot, error)
	DeleteOldSnapshots(ctx context.Context, olderThan time.Time) error
}
type GormBalanceSnapshotRepository struct {
	db *gorm.DB
}

func NewGormBalanceSnapshotRepository(db *gorm.DB) *GormBalanceSnapshotRepository {
	return &GormBalanceSnapshotRepository{db: db}
}

func (r *GormBalanceSnapshotRepository) Create(ctx context.Context, snapshot *Balance) error {
	return r.db.WithContext(ctx).Create(snapshot.balance).Error
}

func (r *GormBalanceSnapshotRepository) GetLatestByUserAndExchange(ctx context.Context, userID, exchangeID uuid.UUID) (
	[]models.BalanceSnapshot, error) {
	var snapshots []models.BalanceSnapshot
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND exchange_id = ?", userID, exchangeID).
		Order("snapshot_time DESC").
		Limit(50). // Get latest 50 currency snapshots
		Find(&snapshots).Error
	return snapshots, err
}

func (r *GormBalanceSnapshotRepository) GetHistory(ctx context.Context, userID, exchangeID uuid.UUID, currency string, limit int) (
	[]models.BalanceSnapshot, error) {
	var snapshots []models.BalanceSnapshot
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND exchange_id = ? AND currency = ?", userID, exchangeID, currency).
		Order("snapshot_time DESC").
		Limit(limit).
		Find(&snapshots).Error
	return snapshots, err
}

func (r *GormBalanceSnapshotRepository) DeleteOldSnapshots(ctx context.Context, olderThan time.Time) error {
	return r.db.WithContext(ctx).Where("snapshot_time < ?", olderThan).Delete(&models.BalanceSnapshot{}).Error
}
