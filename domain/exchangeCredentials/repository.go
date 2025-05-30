package exchangeCredentials

import (
	"context"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/internal/database/models"
	envCofig "github.com/rzabhd80/eye-on/internal/envConfig"
	"github.com/rzabhd80/eye-on/internal/helpers"
	"gorm.io/gorm"
	"time"
)

type IExchangeCredentialRepository interface {
	Create(ctx context.Context, cred *models.ExchangeCredential) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ExchangeCredential, error)
	GetByUserAndExchange(ctx context.Context, userID, exchangeID uuid.UUID) (*models.ExchangeCredential, error)
	Update(ctx context.Context, cred *models.ExchangeCredential) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastUsed(ctx context.Context, id uuid.UUID) error
}

type ExchangeCredentialRepository struct {
	Db      *gorm.DB
	EnvConf *envCofig.AppConfig
}

func NewExchangeCredentialRepository(db *gorm.DB, envConf *envCofig.AppConfig) *ExchangeCredentialRepository {
	return &ExchangeCredentialRepository{Db: db, EnvConf: envConf}
}

func (r *ExchangeCredentialRepository) Create(ctx context.Context, cred *models.ExchangeCredential) error {
	return r.Db.WithContext(ctx).Create(cred).Error
}

func (r *ExchangeCredentialRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ExchangeCredential, error) {
	var cred models.ExchangeCredential
	err := r.Db.WithContext(ctx).Preload("User").Preload("Exchange").First(&cred, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	key := r.EnvConf.EncryptionKey
	apiKeyDecr, err := helpers.DecryptAPIKey(cred.APIKey, key)
	if err != nil {
		return nil, err
	}

	var accKeyDec string
	if cred.AccessKey != "" {
		accKeyDec, err = helpers.DecryptAPIKey(cred.AccessKey, key)
		if err != nil {
			return nil, err
		}
	}
	cred.AccessKey = accKeyDec
	cred.APIKey = apiKeyDecr
	return &cred, nil
}

func (r *ExchangeCredentialRepository) GetByUserAndExchange(ctx context.Context, userID, exchangeID uuid.UUID) (
	*models.ExchangeCredential, error) {
	var creds *models.ExchangeCredential
	err := r.Db.WithContext(ctx).
		Preload("Exchange").
		Where("user_id = ? AND exchange_id = ? AND is_active = ?", userID, exchangeID, true).
		Order("created_at DESC").
		First(&creds).Error
	key := r.EnvConf.EncryptionKey
	apiKeyDecr, err := helpers.DecryptAPIKey(creds.APIKey, key)
	if err != nil {
		return nil, err
	}

	var accKeyDec string
	if creds.AccessKey != "" {
		accKeyDec, err = helpers.DecryptAPIKey(creds.AccessKey, key)
		if err != nil {
			return nil, err
		}
	}
	creds.AccessKey = accKeyDec
	creds.APIKey = apiKeyDecr
	return creds, err

}

func (r *ExchangeCredentialRepository) Update(ctx context.Context, cred *models.ExchangeCredential) error {
	return r.Db.WithContext(ctx).Save(cred).Error
}

func (r *ExchangeCredentialRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.Db.WithContext(ctx).Delete(&models.ExchangeCredential{}, id).Error
}

func (r *ExchangeCredentialRepository) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.Db.WithContext(ctx).Model(&models.ExchangeCredential{}).
		Where("id = ?", id).
		Update("last_used", now).Error
}
