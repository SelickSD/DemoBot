package config

import (
	"os"

	"github.com/SelickSD/DemoBot.git/internal/logger"
)

type Config struct {
	BotToken    string
	Debug       bool
	ConfigEmail string
	BotName     string
}

func Load() *Config {
	// Инициализируем логгер
	logger.Init()

	botToken := getEnv("BOT_TOKEN", "")
	if botToken == "" {
		logger.Error.Fatal("BOT_TOKEN environment variable is required")
	}

	configEmail := getEnv("CONFIG_EMAIL", "")
	if configEmail == "" { // Была ошибка - проверяли botToken вместо configEmail
		logger.Error.Fatal("CONFIG_EMAIL environment variable is required")
	}

	botName := getEnv("BOT_NAME", "")
	if botName == "" { // Была ошибка - проверяли botToken вместо botName
		logger.Error.Fatal("BOT_NAME environment variable is required")
	}

	return &Config{
		BotToken:    botToken,
		Debug:       getEnv("DEBUG", "false") == "true",
		ConfigEmail: configEmail,
		BotName:     botName,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
