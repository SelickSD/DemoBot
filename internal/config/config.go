package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type BotConfig struct {
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`
	BotName  string `envconfig:"BOT_NAME" required:"true"`
}

type Config struct {
	Debug       bool   `envconfig:"DEBUG" default:"false"`
	ConfigEmail string `envconfig:"CONFIG_EMAIL" required:"true"`
	BotConfig   BotConfig
}

func Load() (*Config, error) {
	cfg := new(Config)

	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("error while parse env config | %w", err)
	}

	return cfg, nil
}
