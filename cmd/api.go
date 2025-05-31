package main

import (
	"context"
	"fmt"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/rzabhd80/eye-on/api/bitpin"
	"github.com/rzabhd80/eye-on/api/nobitex"
	userService "github.com/rzabhd80/eye-on/api/user"
	"github.com/rzabhd80/eye-on/domain/balance"
	"github.com/rzabhd80/eye-on/domain/exchange"
	bitpinEntity "github.com/rzabhd80/eye-on/domain/exchange/bitpin"
	nobitexEntity "github.com/rzabhd80/eye-on/domain/exchange/nobitex"
	"github.com/rzabhd80/eye-on/domain/exchange/registry"
	"github.com/rzabhd80/eye-on/domain/exchangeCredentials"
	"github.com/rzabhd80/eye-on/domain/order"
	"github.com/rzabhd80/eye-on/domain/orderBook"
	"github.com/rzabhd80/eye-on/domain/traidingPair"
	"github.com/rzabhd80/eye-on/domain/user"
	db "github.com/rzabhd80/eye-on/internal/database"
	"github.com/rzabhd80/eye-on/internal/envConfig"
	"github.com/rzabhd80/eye-on/internal/helpers"
	"github.com/rzabhd80/eye-on/internal/redis"
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
	redisConn := redis.RedisConnection{EnvConf: devConf}
	appRedisClient := redisConn.NewRedisClient()
	jwtParser := helpers.JWTParser{EnvConf: devConf}
	request := helpers.NewRequest(10 * time.Second)

	defer func(redisCLient *redis2.Client) {
		err := redisCLient.Close()
		if err != nil {

		}

	}(appRedisClient)
	if err != nil {
		return err
	}
	// Test connection
	ctxRedis := context.Background()
	pong, err := appRedisClient.Ping(ctxRedis).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Printf("Connected to Redis: %s\n", pong)
	ctxRedis.Done()
	err = psqlDb.Migrate()
	if err != nil {
		return err
	}
	exchangeRepo := exchange.NewExchangeRepository(psqlDb.GormDb)

	tradingPairRepo := traidingPair.TradingPairRepository{DB: psqlDb.GormDb}

	exchangeCredRepo := exchangeCredentials.NewExchangeCredentialRepository(psqlDb.GormDb, devConf)
	orderRepo := order.NewOrderHistoryRepository(psqlDb.GormDb)
	orderBookRepo := orderBook.NewOrderBookSnapshotRepository(psqlDb.GormDb)
	balanceRepo := balance.NewBalanceSnapshotRepository(psqlDb.GormDb)

	exchangeRegistery := registry.NewRegistry(exchangeRepo, &tradingPairRepo, exchangeCredRepo, psqlDb.GormDb)

	registry.SetDefaultRegistry(exchangeRegistery)

	bitpinSymbolRegistry := bitpinEntity.BitpinSymbolRegistry{}
	NobitexSymbolRegistry := nobitexEntity.NobitexSymbolRegistry{}
	bitpinExchange, err := registry.GetOrCreateExchange(ctx, registry.ExchangeConfig{
		Name:          "bitpin",
		DisplayName:   "bitpin",
		BaseURL:       "https://api.bitpin.ir",
		RateLimit:     0,
		Features:      nil,
		SymbolFactory: &bitpinSymbolRegistry,
	})

	nobitexExchange, err := registry.GetOrCreateExchange(ctx, registry.ExchangeConfig{
		Name:          "nobitex",
		DisplayName:   "nobitex",
		BaseURL:       "https://api.nobitex.ir",
		RateLimit:     0,
		Features:      nil,
		SymbolFactory: &NobitexSymbolRegistry,
	})
	if err != nil {
		return err
	}

	userRepo := user.NewUserRepository(psqlDb.GormDb)

	app := fiber.New()

	userRouter := userService.Router{
		Service: &userService.UserAuthService{User: &user.User{
			UserRepo:         userRepo,
			ExchangeRepo:     exchangeRepo,
			ExchangeCredRepo: exchangeCredRepo,
			JwtParser:        &jwtParser,
			EnvConf:          devConf,
		}},
		Parser: &jwtParser,
	}

	nobitexRouter := nobitex.Router{
		Service: &nobitex.NobitexService{
			Exchange: &nobitexEntity.NobitexExchange{
				NobitexExchangeModel:   nobitexExchange.Exchange,
				ExchangeRepo:           exchangeRepo,
				ExchangeCredentialRepo: exchangeCredRepo,
				UserRepo:               userRepo,
				TradingPairRepo:        &tradingPairRepo,
				OrderRepo:              orderRepo,
				OrderBookRepo:          orderBookRepo,
				BalanceRepo:            balanceRepo,
				Request:                request,
			},
		},
		Parser: &jwtParser,
	}
	bitpinRouter := bitpin.Router{
		Service: &bitpin.BitpinService{
			Exchange: &bitpinEntity.BitpinExchange{
				BitpinExchangeModel:    bitpinExchange.Exchange,
				ExchangeRepo:           exchangeRepo,
				ExchangeCredentialRepo: exchangeCredRepo,
				UserRepo:               userRepo,
				TradingPairRepo:        &tradingPairRepo,
				OrderRepo:              orderRepo,
				OrderBookRepo:          orderBookRepo,
				BalanceRepo:            balanceRepo,
				Request:                request,
			},
		},
		Parser: &jwtParser,
	}

	//Register your routes here
	userRouter.SetUserRouter(app)
	bitpinRouter.SetUserRouter(app)
	nobitexRouter.SetUserRouter(app)

	ctx, stp := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)

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
	err = psqlDb.Close()
	if err != nil {
		return err
	}
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Fatal("Server shutdown error: %v", zap.String("error", err.Error()))
		return err
	}

	logger.Info("Server shutdown complete")
	return nil
}
