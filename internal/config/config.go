package config

import (
	"os"

	"github.com/SelickSD/DemoBot.git/internal/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	BotToken    string
	Debug       bool
	ConfigEmail string
	BotName     string
	AiApiKey string
}

func Load() *Config {
	// Инициализируем логгер
	logger.Init()

	// Загружаем .env файл (игнорируем ошибку если файла нет)
	// Это для локальной разработки, на проде используем переменные окружения
	_ = godotenv.Load()

	botToken := getEnv("BOT_TOKEN", "")
	if botToken == "" {
		logger.Error.Fatal("BOT_TOKEN environment variable is required")
	}

	configEmail := getEnv("CONFIG_EMAIL", "")
	if configEmail == "" {
		logger.Error.Fatal("CONFIG_EMAIL environment variable is required")
	}

	botName := getEnv("BOT_NAME", "")
	if botName == "" {
		logger.Error.Fatal("BOT_NAME environment variable is required")
	}

	aiApiKey := getEnv("AI_API_KEY", "")
	if aiApiKey == "" {
		logger.Error.Fatal("AI_API_KEY environment variable is required")
	}

	logger.Info.Printf("Bot configured with name: %s", botName)
	logger.Info.Printf("Debug mode: %t", getEnv("DEBUG", "false") == "true")

	return &Config{
		BotToken:    botToken,
		Debug:       getEnv("DEBUG", "false") == "true",
		ConfigEmail: configEmail,
		BotName:     botName,
		AiApiKey: aiApiKey,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
