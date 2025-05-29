package main

import (
	"context"
	"fmt"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	user2 "github.com/rzabhd80/eye-on/api/user"
	"github.com/rzabhd80/eye-on/domain/exchange"
	"github.com/rzabhd80/eye-on/domain/exchange/registry"
	"github.com/rzabhd80/eye-on/domain/exchangeCredentials"
	"github.com/rzabhd80/eye-on/domain/traidingPair"
	"github.com/rzabhd80/eye-on/domain/user"
	db "github.com/rzabhd80/eye-on/internal/database"
	"github.com/rzabhd80/eye-on/internal/envConfig"
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
	redisCLient := redisConn.NewRedisClient()

	defer func(redisCLient *redis2.Client) {
		err := redisCLient.Close()
		if err != nil {

		}

	}(redisCLient)
	if err != nil {
		return err
	}
	// Test connection
	ctxRedis := context.Background()
	pong, err := redisCLient.Ping(ctxRedis).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Printf("Connected to Redis: %s\n", pong)
	ctxRedis.Done()

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

	exchangeRegistery := registry.NewRegistry(&exchangeRepo, &tradingPairRepo, &exchangeCredRepo)

	registry.SetDefaultRegistry(exchangeRegistery)

	bitpinExchange, err := registry.GetOrCreateExchange(ctx, registry.ExchangeConfig{
		Name:        "bitpint",
		DisplayName: "bitpin",
		BaseURL:     "https://api.bitpin.ir",
		RateLimit:   0,
		Timeout:     0,
		Features:    nil,
		Label:       "",
	})
	nobitexExchange, err := registry.GetOrCreateExchange(ctx, registry.ExchangeConfig{
		Name:        "nobitex",
		DisplayName: "nobitex",
		BaseURL:     "https://api.nobitex.ir/v3",
		RateLimit:   0,
		Timeout:     0,
		Features:    nil,
		Label:       "",
	})

	if err != nil {
		return err
	}
	userRepo := user.UserRepository{
		Db: psqlDb.GormDb,
	}

	app := fiber.New()

	userRouter := user2.Router{
		Service: &user2.UserAuthService{User: &user.User{
			UserRepo:         &userRepo,
			ExchangeRepo:     &exchangeRepo,
			ExchangeCredRepo: &exchangeCredRepo,
		}},
	}

	//Register your routes here
	userRouter.SetUserRouter(app)

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
