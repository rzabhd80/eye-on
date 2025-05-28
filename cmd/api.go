package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/rzabhd80/eye-on/domain/exchange"
	"github.com/rzabhd80/eye-on/domain/exchange/registry"
	"github.com/rzabhd80/eye-on/domain/exchangeCredentials"
	"github.com/rzabhd80/eye-on/domain/traidingPair"
	db "github.com/rzabhd80/eye-on/internal/database"
	"github.com/rzabhd80/eye-on/internal/envConfig"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func apiService(cntx *cli.Context, logger *zap.Logger) error {

	ctx, cancel := context.WithCancel(cntx.Context)
	defer cancel()
	devConf, err := envCofig.LoadConfig()
	if err != nil {
		return err
	}

	psqlDb, err := db.NewDatabase(devConf)
	if err != nil {
		return err
	}
	err = psqlDb.Migrate()
	if err != nil {
		return err
	}
	exchangeRepo := exchange.ExchangeRepository{
		Db: psqlDb.GormDb,
	}

	tradingPairRepo := traidingPair.TradingPairRepository{DB: psqlDb.GormDb}

	exchangeCredRepo := exchangeCredentials.ExchangeCredentialRepository{
		Db: psqlDb.GormDb,
	}

	exchangeResitery := registry.NewRegistry(&exchangeRepo, &tradingPairRepo, &exchangeCredRepo)

	//TODO add exchanges
	
	exchangeResitery.Register("", func(cfg registry.ExchangeConfig) (registry.IExchange, error) { return nil, nil })

	ctx, stp := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	app := fiber.New()

	go func() {
		logger.Info("Starting server on port %s", zap.String("port", devConf.PORT))
		if err := app.Listen(":" + devConf.PORT); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()
	<-ctx.Done()
	stp()
	logger.Info("Shutting down server...")

	// Shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	//err = psqlDb.Close()
	//if err != nil {
	//	return err
	//}
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Fatal("Server shutdown error: %v", zap.String("error", err.Error()))
		return err
	}

	logger.Info("Server shutdown complete")
	return nil
}
