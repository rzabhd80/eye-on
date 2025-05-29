package balance

import (
	"context"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/internal/database/models"
	"gorm.io/gorm"
	"time"
)

type IBalanceSnapshotRepository interface {
	Create(ctx context.Context, snapshot *models.BalanceSnapshot) error
	GetLatestByUserAndExchange(ctx context.Context, userID, exchangeID uuid.UUID) (*[]models.BalanceSnapshot, error)
	GetHistory(ctx context.Context, userID, exchangeID uuid.UUID, currency string, limit int) (*[]models.BalanceSnapshot, error)
	DeleteOldSnapshots(ctx context.Context, olderThan time.Time) error
}
type BalanceSnapshotRepository struct {
	db *gorm.DB
}

func NewBalanceSnapshotRepository(db *gorm.DB) *BalanceSnapshotRepository {
	return &BalanceSnapshotRepository{db: db}
}

func (r *BalanceSnapshotRepository) Create(ctx context.Context, snapshot *models.BalanceSnapshot) error {
	return r.db.WithContext(ctx).Create(snapshot).Error
}

func (r *BalanceSnapshotRepository) BulkCreate(ctx context.Context, snapshots *[]models.BalanceSnapshot) error {
	return r.db.WithContext(ctx).Create(snapshots).Error
}
func (r *BalanceSnapshotRepository) GetLatestByUserAndExchange(ctx context.Context, userID, exchangeID uuid.UUID) (
	*[]models.BalanceSnapshot, error) {
	var snapshots *[]models.BalanceSnapshot
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND exchange_id = ?", userID, exchangeID).
		Order("snapshot_time DESC").
		Limit(50). // Get latest 50 currency snapshots
		Find(&snapshots).Error
	return snapshots, err
}

func (r *BalanceSnapshotRepository) GetHistory(ctx context.Context, userID, exchangeID uuid.UUID, currency string, limit int) (
	*[]models.BalanceSnapshot, error) {
	var snapshots *[]models.BalanceSnapshot
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND exchange_id = ? AND currency = ?", userID, exchangeID, currency).
		Order("snapshot_time DESC").
		Limit(limit).
		Find(&snapshots).Error
	return snapshots, err
}

func (r *BalanceSnapshotRepository) DeleteOldSnapshots(ctx context.Context, olderThan time.Time) error {
	return r.db.WithContext(ctx).Where("snapshot_time < ?", olderThan).Delete(&models.BalanceSnapshot{}).Error
}
