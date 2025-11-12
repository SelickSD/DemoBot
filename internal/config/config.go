package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type BotConfig struct {
	Token string `envconfig:"TOKEN" required:"true"`
	Name  string `envconfig:"NAME" required:"true"`
}

type Config struct {
	Debug       bool   `envconfig:"DEBUG" default:"false"`
	ConfigEmail string `envconfig:"CONFIG_EMAIL" required:"true"`
	Bot         BotConfig
}

func Load() (*Config, error) {
	// Подгружаем .env только если есть (для локальной разработки)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	}

	cfg := new(Config)

	// Основные переменные
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("failed to load base config: %w", err)
	}

	// Переменные, начинающиеся с BOT_
	if err := envconfig.Process("BOT", &cfg.Bot); err != nil {
		return nil, fmt.Errorf("failed to load bot config: %w", err)
	}

	return cfg, nil
}
