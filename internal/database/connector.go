package postgresDB

import (
	"database/sql"
	"fmt"
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rzabhd80/eye-on/internal/envConfig"
)

type Database struct {
	db  *sql.DB
	cfg *envCofig.AppConfig
}

func NewDatabase(config *envCofig.AppConfig) (*Database, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	database := &Database{db: db, cfg: config}
	return database, nil
}

func (database *Database) Migrate() error {
	driver, err := postgres.WithInstance(database.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migrate driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", database.cfg.DbName, driver,
	)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

func (db *Database) Close() error {
	err := db.db.Close()
	if err != nil {
		return err
	}
	return nil
}
