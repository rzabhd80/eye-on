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
	"github.com/rzabhd80/eye-on/internal/helpers"
	"gorm.io/gorm"
	"time"
)

type Constructor func(cfg ExchangeConfig) (IExchange, error)

var (
	constructors = make(map[string]Constructor)
)

// Registry manages exchange registration and creation
type Registry struct {
	exchangeRepo            *exchange.ExchangeRepository
	tradingPairRepo         *traidingPair.TradingPairRepository
	exchangeCredentialsRepo *exchangeCredentials.ExchangeCredentialRepository
	constructors            map[string]Constructor
}

func NewRegistry(repo *exchange.ExchangeRepository, tradingRepo *traidingPair.TradingPairRepository,
	exchangeCredentialsRepo *exchangeCredentials.ExchangeCredentialRepository) *Registry {
	return &Registry{
		exchangeRepo:            repo,
		tradingPairRepo:         tradingRepo,
		exchangeCredentialsRepo: exchangeCredentialsRepo,
		constructors:            make(map[string]Constructor),
	}
}

// Register is called by each adapter in its init()
func (r *Registry) Register(name string, constructor Constructor, envConfig envCofig.AppConfig) {
	if _, dup := r.constructors[name]; dup {
		panic("exchange " + name + " already registered")
	}
	r.constructors[name] = constructor
}

// ExchangeResult contains both the database model and runtime instance
type ExchangeResult struct {
	Exchange        *models.Exchange
	Credential      *models.ExchangeCredential
	symbols         []traidingPair.TradingPair
	Instance        IExchange
	IsNewExchange   bool
	IsNewCredential bool
}

// GetOrCreateExchangeConfig creates or retrieves an exchange and its credentials from the database
func (r *Registry) GetOrCreateExchangeConfig(ctx context.Context, cfg ExchangeConfig, envConf envCofig.AppConfig) (
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

	// Check if credential exists for this user and exchangeInstance
	var credential models.ExchangeCredential
	label := cfg.Label
	if label == "" {
		label = "Default"
	}

	err = tx.Where(
		"user_id = ? AND exchange_id = ? AND label = ? AND is_active = ?",
		cfg.UserID, exchangeInstance.ID, label, true,
	).First(&credential).Error

	isNewCredential := false
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			secretKey, err := helpers.Encrypt(cfg.SecretKey, []byte(envConf.EncryptionKey))
			if err != nil {
				return nil, err
			}
			// Create new credential
			credential = models.ExchangeCredential{
				BaseModel: models.BaseModel{
					ID: uuid.New(),
				},
				UserID:      cfg.UserID,
				ExchangeID:  exchangeInstance.ID,
				Label:       label,
				APIKey:      cfg.APIKey,
				SecretKey:   secretKey,
				Passphrase:  cfg.Passphrase,
				IsActive:    true,
				IsTestnet:   cfg.IsTestnet,
				Permissions: models.JSONB(map[string]interface{}{}), // Initialize empty
			}

			if err := tx.Create(&credential).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create exchangeInstance credential: %w", err)
			}
			isNewCredential = true
		} else {
			tx.Rollback()
			return nil, fmt.Errorf("failed to query exchangeInstance credential: %w", err)
		}
	} else {
		// Update existing credential if API keys have changed
		updated := false
		if credential.APIKey != cfg.APIKey {
			credential.APIKey = cfg.APIKey
			updated = true
		}
		if credential.SecretKey != cfg.SecretKey {
			credential.SecretKey = cfg.SecretKey
			updated = true
		}
		if credential.Passphrase != cfg.Passphrase {
			credential.Passphrase = cfg.Passphrase
			updated = true
		}
		if credential.IsTestnet != cfg.IsTestnet {
			credential.IsTestnet = cfg.IsTestnet
			updated = true
		}

		if updated {
			if err := tx.Save(&credential).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to update exchangeInstance credential: %w", err)
			}
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

	// Update last used timestamp
	now := time.Now()
	credential.LastUsed = &now
	r.exchangeRepo.Db.WithContext(ctx).Save(&credential)

	return &ExchangeResult{
		Exchange:        &exchangeInstance,
		Credential:      &credential,
		Instance:        instance,
		IsNewExchange:   isNewExchange,
		IsNewCredential: isNewCredential,
		symbols:         tradingPairs,
	}, nil
}

// GetExchange retrieves an existing exchange instance without creating new records
func (r *Registry) GetExchange(ctx context.Context, userID uuid.UUID, exchangeName, label string) (*ExchangeResult, error) {
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
		Exchange:        &credential.Exchange,
		Credential:      &credential,
		Instance:        instance,
		IsNewExchange:   false,
		IsNewCredential: false,
	}, nil
}

// ListSupportedExchanges returns a list of registered exchange names
func (r *Registry) ListSupportedExchanges() []string {
	exchanges := make([]string, 0, len(r.constructors))
	for name := range r.constructors {
		exchanges = append(exchanges, name)
	}
	return exchanges
}

// DeactivateCredential marks a credential as inactive
func (r *Registry) DeactivateCredential(ctx context.Context, userID uuid.UUID, exchangeName, label string) error {
	if label == "" {
		label = "Default"
	}

	result := r.exchangeRepo.Db.WithContext(ctx).
		Table("exchange_credentials").
		Where("user_id = ? AND label = ?", userID, label).
		Where("exchange_id IN (SELECT id FROM exchanges WHERE name = ?)", exchangeName).
		Update("is_active", false)

	if result.Error != nil {
		return fmt.Errorf("failed to deactivate credential: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("credential not found")
	}

	return nil
}

// Global registry instance (optional, for backward compatibility)
var defaultRegistry *Registry

// SetDefaultRegistry sets the global registry instance
func SetDefaultRegistry(registry *Registry) {
	defaultRegistry = registry
}

// Register registers an exchange constructor in the default registry
func Register(name string, constructor Constructor, envConf envCofig.AppConfig) {
	if defaultRegistry == nil {
		panic("default registry not initialized. Call SetDefaultRegistry first")
	}
	defaultRegistry.Register(name, constructor, envConf)
}

// GetOrCreateExchange uses the default registry
func GetOrCreateExchange(ctx context.Context, cfg ExchangeConfig, envConf envCofig.AppConfig) (*ExchangeResult, error) {
	if defaultRegistry == nil {
		return nil, fmt.Errorf("default registry not initialized")
	}
	return defaultRegistry.GetOrCreateExchangeConfig(ctx, cfg, envConf)
}
