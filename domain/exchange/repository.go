package exchange

import (
	"context"
	"github.com/google/uuid"
)

type ExchangeRepository interface {
	Create(ctx context.Context, exchange *Exchange) error
	GetByID(ctx context.Context, id uuid.UUID) (*Exchange, error)
	GetByName(ctx context.Context, name string) (*Exchange, error)
	Update(ctx context.Context, exchange *Exchange) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, activeOnly bool) ([]Exchange, error)
}

type ExchangeCredentialRepository interface {
	Create(ctx context.Context, cred *ExchangeCredential) error
	GetByID(ctx context.Context, id uuid.UUID) (*ExchangeCredential, error)
	GetByUserAndExchange(ctx context.Context, userID, exchangeID uuid.UUID) ([]ExchangeCredential, error)
	Update(ctx context.Context, cred *ExchangeCredential) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastUsed(ctx context.Context, id uuid.UUID) error
}
