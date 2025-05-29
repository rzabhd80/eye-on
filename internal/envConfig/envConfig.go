package envCofig

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"os"
)

func CheckFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

type AppConfig struct {
	DatabaseConfig
	RedisConfig
	AppName       string `env:"APP_NAME" envDefault:"eye on"`
	AppVersion    string `env:"APP_VERSION" envDefault:"0.0.1"`
	HOST          string `env:"HOST" envDefault:"0.0.0.0"`
	PORT          string `env:"PORT" envDefault:"8080"`
	EncryptionKey string `env:"ENCRYPTION_KEY"`
	JWTKey        string `env:"JWT_KEY"`
}

type RedisConfig struct {
	RedisHost     string ` env:"REDIS_HOST" json:"redis-host"`
	RedisPort     int    `env:"REDIS_PORT" json:"redis-port"`
	RedisPassword string `env:"REDIS_PASSWORD" json:"redis_password"`
	RedisDB       int    `env:"REDIS_DB" json:"redis-db"`
}

type DatabaseConfig struct {
	DbHost     string `env:"DB_HOST" envDefault:"postgres"`
	DbPort     string `env:"DB_PORT" envDefault:"5432"`
	DbUser     string `env:"DB_USER" envDefault:"postgres"`
	DbPassword string `env:"DB_PASSWORD" envDefault:"postgres"`
	DbName     string `env:"DB_NAME" envDefault:"postgres"`
}

func LoadConfig() (*AppConfig, error) {
	var envFile string = ".env"
	if CheckFileExists(envFile) {
		if err := godotenv.Load(envFile); err != nil {
			return nil, err
		}
	}
	devEnv := &AppConfig{}
	if err := env.Parse(devEnv); err != nil {
		return nil, err
	}
	return devEnv, nil
}
