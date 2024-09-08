package main

import (
	"context"
	"github.com/barancanatbas/messaging/internal/delivery"
	"github.com/barancanatbas/messaging/pkg/cache"
	"github.com/barancanatbas/messaging/pkg/httpclient"
	"github.com/barancanatbas/messaging/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/barancanatbas/messaging/docs"
	"github.com/barancanatbas/messaging/internal/message"
	"github.com/barancanatbas/messaging/pkg/mysql"
	"github.com/barancanatbas/messaging/pkg/redis"
	"github.com/barancanatbas/messaging/pkg/viper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	swagger "github.com/swaggo/fiber-swagger"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := viper.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	db, err := mysql.NewMysqlClient(cfg.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MySQL")
	}

	redisClient, err := redis.New(cfg.Redis)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}

	validatorService := validator.NewValidator()

	cacheService := cache.NewCacheService(redisClient)
	httpClient := httpclient.NewHttpClient(cfg.HttpClient)

	deliveryService := delivery.NewDeliveryService(httpClient)

	messageRepository := message.NewMessageRepository(db)
	messageService := message.NewMessageService(cfg.App, messageRepository, cacheService, deliveryService)

	app := fiber.New()
	app.Get("/swagger/*", swagger.WrapHandler)

	message.NewMessageHandler(app, messageService, validatorService)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(cfg.App.Port); err != nil {
			log.Error().Err(err).Msg("Fiber server failed")
		}
	}()

	<-quit
	log.Info().Msg("Graceful shutdown initiated...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error().Err(err).Msg("Server failed to shutdown gracefully")
	} else {
		log.Info().Msg("Server shutdown completed")
	}
}
