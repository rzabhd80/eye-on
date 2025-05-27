package order

import "github.com/rzabhd80/eye-on/internal/database/models"

type Order struct {
	orderHistory *models.OrderHistory
	orderEvents  []*models.OrderEvent
}
