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
	envCofig "github.com/rzabhd80/eye-on/internal/envConfig"
	"gorm.io/gorm"
	"time"
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
	symbols       []traidingPair.TradingPair
	Instance      IExchange
	IsNewExchange bool
}

// GetOrCreateExchangeConfig creates or retrieves an exchange and its credentials from the database
func (r *ExchangeRegistry) GetOrCreateExchangeConfig(ctx context.Context, cfg ExchangeConfig, envConf envCofig.AppConfig) (
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
	//TODO!!! it should be edited to check each and every on of them
	//setting up the symbols
	var exchangeSymbols []string
	for _, symbol := range cfg.symbols {
		exchangeSymbols = append(exchangeSymbols, symbol.Symbol)
	}
	tradingPairs, err := r.tradingPairRepo.GetSymbolsList(ctx, exchangeInstance.ID, true, exchangeSymbols)
	if err != nil {
		return nil, err
	}
	if len(tradingPairs) != len(cfg.symbols) {
		var newPairs []traidingPair.TradingPair

		for _, symbolPair := range cfg.symbols {
			tradingPair := traidingPair.TradingPair{
				TradingPair: models.TradingPair{
					BaseModel: models.BaseModel{
						ID: uuid.New(),
					},
					ExchangeID: exchangeInstance.ID,
					Symbol:     symbolPair.Symbol,
					BaseAsset:  symbolPair.BaseAsset,
					QuoteAsset: symbolPair.QuoteAsset,
				},
			}
			errPair := r.tradingPairRepo.Create(ctx, &tradingPair)
			if errPair != nil {
				return nil, errPair
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
		symbols:       tradingPairs,
	}, nil
}

// GetExchange retrieves an existing exchange instance without creating new records
func (r *ExchangeRegistry) GetExchange(ctx context.Context, userID uuid.UUID, exchangeName, label string) (*ExchangeResult, error) {
	if label == "" {
		label = "Default"
	}

	constructor, ok := r.constructors[exchangeName]
	if !ok {
		return nil, fmt.Errorf("exchange %s not supported", exchangeName)
	}

	// Query with joins to get both exchange and credential info
	var credential models.ExchangeCredential
	err := r.exchangeRepo.Db.WithContext(ctx).
		Preload("Exchange").
		Where("user_id = ? AND label = ? AND is_active = ?", userID, label, true).
		Joins("JOIN exchanges ON exchanges.id = exchange_credentials.exchange_id").
		Where("exchanges.name = ? AND exchanges.is_active = ?", exchangeName, true).
		First(&credential).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("exchange credential not found for user %s, exchange %s, label %s", userID, exchangeName, label)
		}
		return nil, fmt.Errorf("failed to query exchange credential: %w", err)
	}

	// Create exchange config from stored data
	cfg := ExchangeConfig{
		Name:        credential.Exchange.Name,
		DisplayName: credential.Exchange.DisplayName,
		BaseURL:     credential.Exchange.BaseURL,
		APIKey:      credential.APIKey,
		SecretKey:   credential.SecretKey,
		RefreshKey:  credential.RefreshKey,
		Passphrase:  credential.Passphrase,
		IsTestnet:   credential.IsTestnet,
		RateLimit:   credential.Exchange.RateLimit,
		UserID:      credential.UserID,
		Label:       credential.Label,
	}

	// Create the exchange instance
	instance, err := constructor(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create exchange instance: %w", err)
	}

	// Update last used timestamp
	now := time.Now()
	credential.LastUsed = &now
	r.exchangeRepo.Db.WithContext(ctx).Save(&credential)

	return &ExchangeResult{
		Exchange: &credential.Exchange,
		Instance: instance,
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

// Global registry instance (optional, for backward compatibility)
var defaultRegistry *ExchangeRegistry

// SetDefaultRegistry sets the global registry instance
func SetDefaultRegistry(registry *ExchangeRegistry) {
	defaultRegistry = registry
}

// Register registers an exchange constructor in the default registry
func Register(name string, constructor Constructor, envConf envCofig.AppConfig) {
	if defaultRegistry == nil {
		panic("default registry not initialized. Call SetDefaultRegistry first")
	}
	defaultRegistry.Register(name, constructor)
}

// GetOrCreateExchange uses the default registry
func GetOrCreateExchange(ctx context.Context, cfg ExchangeConfig, envConf envCofig.AppConfig) (*ExchangeResult, error) {
	if defaultRegistry == nil {
		return nil, fmt.Errorf("default registry not initialized")
	}
	return defaultRegistry.GetOrCreateExchangeConfig(ctx, cfg, envConf)
}
