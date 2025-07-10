package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddress  string `env:"REDIS_ADDRESS,required"`
	RedisPassword string `env:"REDIS_PASSWORD,required"`
	RedisDatabase string `env:"REDIS_DB,required"`

	DBSource string `env:"DB_SOURCE,required"`
	// WebAppURL string `env:"WEB_APP_URL,required"`
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("error loading .env file: %v", err)
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("error parsing environment variables: %w", err)
	}

	return &cfg, nil
}
