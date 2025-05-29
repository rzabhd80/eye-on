package registry

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/domain/exchange"
	"github.com/rzabhd80/eye-on/domain/exchangeCredentials"
	"github.com/rzabhd80/eye-on/domain/traidingPair"
	"github.com/rzabhd80/eye-on/internal/database/models"
	"gorm.io/gorm"
)

type Constructor func(cfg ExchangeConfig) (IExchange, error)

var (
	constructors = make(map[string]Constructor)
)

// ExchangeRegistry manages exchange registration and creation
type ExchangeRegistry struct {
	exchangeRepo            *exchange.ExchangeRepository
	tradingPairRepo         *traidingPair.TradingPairRepository
	exchangeCredentialsRepo *exchangeCredentials.ExchangeCredentialRepository
	constructors            map[string]Constructor
}

func NewRegistry(repo *exchange.ExchangeRepository, tradingRepo *traidingPair.TradingPairRepository,
	exchangeCredentialsRepo *exchangeCredentials.ExchangeCredentialRepository) *ExchangeRegistry {
	return &ExchangeRegistry{
		exchangeRepo:            repo,
		tradingPairRepo:         tradingRepo,
		exchangeCredentialsRepo: exchangeCredentialsRepo,
		constructors:            make(map[string]Constructor),
	}
}

// Register is called by each adapter in its init()
func (r *ExchangeRegistry) Register(name string, constructor Constructor) {
	if _, dup := r.constructors[name]; dup {
		panic("exchange " + name + " already registered")
	}
	r.constructors[name] = constructor
}

// ExchangeResult contains both the database model and runtime instance
type ExchangeResult struct {
	Exchange      *models.Exchange
	symbols       []models.TradingPair
	Instance      IExchange
	IsNewExchange bool
}

// GetOrCreateExchangeConfig creates or retrieves an exchange
func (r *ExchangeRegistry) GetOrCreateExchangeConfig(ctx context.Context, cfg ExchangeConfig) (
	*ExchangeResult, error) {

	constructor, ok := r.constructors[cfg.Name]
	if !ok {
		return nil, fmt.Errorf("exchangeInstance %s not supported", cfg.Name)
	}

	// Start a transaction
	tx := r.exchangeRepo.Db.WithContext(ctx).Begin()
	defer func() {

		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if exchangeInstance exists
	var exchangeInstance models.Exchange
	err := tx.Where("name = ? AND is_active = ?", cfg.Name, true).First(&exchangeInstance).Error
	isNewExchange := false

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new exchangeInstance
			exchangeInstance = models.Exchange{
				BaseModel: models.BaseModel{
					ID: uuid.New(),
				},
				Name:        cfg.Name,
				DisplayName: cfg.DisplayName,
				BaseURL:     cfg.BaseURL,
				IsActive:    true,
				RateLimit:   cfg.RateLimit,
				Features:    models.JSONB(cfg.Features),
			}

			if err := tx.Create(&exchangeInstance).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create exchangeInstance: %w", err)
			}
			isNewExchange = true
		} else {
			tx.Rollback()
			return nil, fmt.Errorf("failed to query exchangeInstance: %w", err)
		}
	}

	symbols := cfg.SymbolFactory.RegisterExchangeSymbols(&exchangeInstance)
	//setting up the symbols
	var exchangeSymbols []string
	for _, symbol := range *symbols {
		exchangeSymbols = append(exchangeSymbols, symbol.Symbol)
	}

	tradingPairs, err := r.tradingPairRepo.GetSymbolsList(ctx, exchangeInstance.ID, true, exchangeSymbols)

	if err != nil {
		return nil, err
	}
	existingPairsMap := make(map[string]bool, len(*tradingPairs))

	for _, pair := range *tradingPairs {
		existingPairsMap[pair.Symbol] = true
	}

	var newPairs []models.TradingPair
	for _, symbolPair := range *symbols {
		if !existingPairsMap[symbolPair.Symbol] {
			tradingPair := models.TradingPair{
				BaseModel: models.BaseModel{
					ID: uuid.New(),
				},
				ExchangeID: exchangeInstance.ID,
				Symbol:     symbolPair.Symbol,
				BaseAsset:  symbolPair.BaseAsset,
				QuoteAsset: symbolPair.QuoteAsset,
			}
			errPair := r.tradingPairRepo.Create(ctx, &tradingPair)
			if errPair != nil {
				return nil, fmt.Errorf("failed to create tradingPair for symbol %s: %w", symbolPair.Symbol, errPair)
			}
			newPairs = append(newPairs, tradingPair)
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Create the exchangeInstance instance
	instance, err := constructor(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create exchangeInstance instance: %w", err)
	}

	return &ExchangeResult{
		Exchange:      &exchangeInstance,
		Instance:      instance,
		IsNewExchange: isNewExchange,
		symbols:       *tradingPairs,
	}, nil
}

// ListSupportedExchanges returns a list of registered exchange names
func (r *ExchangeRegistry) ListSupportedExchanges() []string {
	exchanges := make([]string, 0, len(r.constructors))
	for name := range r.constructors {
		exchanges = append(exchanges, name)
	}
	return exchanges
}

var defaultRegistry *ExchangeRegistry

func SetDefaultRegistry(registry *ExchangeRegistry) {
	defaultRegistry = registry
}

// Register registers an exchange constructor in the default registry
func Register(name string, constructor Constructor) {
	if defaultRegistry == nil {
		panic("default registry not initialized. Call SetDefaultRegistry first")
	}
	defaultRegistry.Register(name, constructor)
}

// GetOrCreateExchange uses the default registry
func GetOrCreateExchange(ctx context.Context, cfg ExchangeConfig) (*ExchangeResult, error) {
	if defaultRegistry == nil {
		return nil, fmt.Errorf("default registry not initialized")
	}
	return defaultRegistry.GetOrCreateExchangeConfig(ctx, cfg)
}
