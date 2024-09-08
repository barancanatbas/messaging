package viper

import (
	"github.com/barancanatbas/messaging/config"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadConfig() (*config.Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	viper.AutomaticEnv()

	cfg := &config.Config{
		DB: config.DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
		},
		Redis: config.RedisConfig{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		App: config.AppConfig{
			MessageInterval: viper.GetInt("MESSAGE_INTERVAL"),
			Port:            viper.GetString("APP_PORT"),
		},
		HttpClient: config.HttpClientConfig{
			BaseURL: viper.GetString("HTTPCLIENT_BASE_URL"),
			AuthKey: viper.GetString("HTTPCLIENT_AUTH_KEY"),
		},
	}

	return cfg, nil
}
