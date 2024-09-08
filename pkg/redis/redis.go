package redis

import (
	"fmt"
	"github.com/barancanatbas/messaging/config"
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

func New(config config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	if err := client.Ping().Err(); err != nil {
		log.Error().Err(err).Msg("Could not connect to Redis")
		return nil, fmt.Errorf("could not connect to Redis: %w", err)
	}

	log.Info().Str("addr", config.Addr).Msg("Successfully connected to Redis")

	return client, nil
}
