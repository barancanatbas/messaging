package config

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type AppConfig struct {
	MessageInterval int
	Port            string
}

type HttpClientConfig struct {
	BaseURL string
	AuthKey string
}

type Config struct {
	DB         DatabaseConfig
	Redis      RedisConfig
	App        AppConfig
	HttpClient HttpClientConfig
}
