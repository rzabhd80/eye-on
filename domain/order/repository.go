package order

import (
	"context"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/internal/database/models"
	"gorm.io/gorm"
)

type IOrderRepository interface {
	Create(ctx context.Context, order *models.OrderHistory) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.OrderHistory, error)
	GetByClientOrderID(ctx context.Context, clientOrderID string) (*models.OrderHistory, error)
	GetByExchangeOrderID(ctx context.Context, exchangeOrderID string) (*models.OrderHistory, error)
	GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.OrderHistory, error)
	GetOpenOrders(ctx context.Context, userID, exchangeCredentialID uuid.UUID) ([]models.OrderHistory, error)
	Update(ctx context.Context, order *models.OrderHistory) error
	CreateEvent(ctx context.Context, event *models.OrderEvent) error
	GetByOrderEventID(ctx context.Context, orderHistID uuid.UUID) ([]models.OrderEvent, error)
	EventList(ctx context.Context, limit, offset int) ([]models.OrderEvent, error)
	//UpdateStatusWithEvent(ctx context.Context, orderID uuid.UUID,
	//	status string, executedQty, executedPrice, commission float64,
	//	eventType string, eventTime time.Time) error
	//
	//AddEvent(ctx context.Context, orderID uuid.UUID, event *models.OrderEvent) error
	//GetOrderWithEvents(ctx context.Context, orderID uuid.UUID) (*models.OrderEvent, error)
	//GetEventsByOrder(ctx context.Context, orderID uuid.UUID) ([]OrderEvent, error)

}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderHistoryRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.OrderHistory) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.OrderHistory, error) {
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

func (r *OrderRepository) GetByClientOrderID(ctx context.Context, clientOrderID string) (*models.OrderHistory,
	error) {
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

func (r *OrderRepository) GetByExchangeOrderID(ctx context.Context, exchangeOrderID string) (*models.OrderHistory,
	error) {
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
func (r *OrderRepository) GetOrderHistoryWithTradingPair(ctx context.Context, orderHistoryID uuid.UUID) (*models.OrderHistory, error) {
	var orderHistory models.OrderHistory

	err := r.db.WithContext(ctx).Preload("TradingPair").
		Where("id = ?", orderHistoryID).
		First(&orderHistory).Error

	if err != nil {
		return nil, err
	}

	return &orderHistory, nil
}

func (r *OrderRepository) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.OrderHistory,
	error) {
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

func (r *OrderRepository) GetOpenOrders(ctx context.Context, userID, exchangeCredentialID uuid.UUID) (
	[]models.OrderHistory,
	error) {
	var orders []models.OrderHistory
	err := r.db.WithContext(ctx).
		Preload("Exchange").
		Preload("TradingPair").
		Where("user_id = ? AND exchange_credential_id = ? AND status IN ?",
			userID, exchangeCredentialID, []string{"NEW", "PARTIALLY_FILLED"}).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) Update(ctx context.Context, order *models.OrderHistory) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, executedQty, executedPrice,
	commission float64) error {
	return r.db.WithContext(ctx).Model(&models.OrderHistory{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":         status,
			"executed_qty":   executedQty,
			"executed_price": executedPrice,
			"commission":     commission,
		}).Error
}

func (r *OrderRepository) CreateEvent(ctx context.Context, event *models.OrderEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *OrderRepository) GetByOrderEventID(ctx context.Context, orderHistID uuid.UUID) ([]models.OrderEvent, error) {
	var events []models.OrderEvent
	err := r.db.WithContext(ctx).Where("order_hist_id = ?", orderHistID).Order("event_time ASC").Find(&events).Error
	return events, err
}

func (r *OrderRepository) EventList(ctx context.Context, limit, offset int) ([]models.OrderEvent, error) {
	var events []models.OrderEvent
	err := r.db.WithContext(ctx).Preload("OrderHistory").Order("recorded_at DESC").Limit(limit).Offset(offset).
		Find(&events).Error
	return events, err
}
