package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	DBSource         string `env:"DB_SOURCE,required"`
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
	//	MongoURI         string `env:"MONGODB_URI,required"`
	//
	// MongoDatabase    string `env:"MONGODB_DATABASE,required"`
	// MongoUsername    string `env:"MONGODB_USER,required"`
	// MongoPassword    string `env:"MONGODB_PASS,required"`
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing environment variables: %w", err)
	}

	return &cfg, nil
}
