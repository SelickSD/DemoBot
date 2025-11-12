package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type BotConfig struct {
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`
	BotName  string `envconfig:"BOT_NAME" required:"true"`
}

type Config struct {
	Debug       bool   `envconfig:"DEBUG" default:"false"`
	ConfigEmail string `envconfig:"CONFIG_EMAIL" required:"true"`
	Bot         BotConfig
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️  .env file not found, using system environment variables")
	}

	cfg := new(Config)

	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("error while parse env config | %w", err)
	}

	return cfg, nil
}
