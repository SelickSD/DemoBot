package config

import (
    "os"
    "log"
)

type Config struct {
    BotToken string
    Debug    bool
	ConfigEmail string
	BotName string
    // Другие параметры...
}

func Load() *Config {
    botToken := getEnv("BOT_TOKEN", "")
    if botToken == "" {
        log.Fatal("BOT_TOKEN environment variable is required")
    }

	configEmail := getEnv("CONFIG_EMAIL", "")
    if botToken == "" {
        log.Fatal("CONFIG_EMAIL environment variable is required")
    }

	botName := getEnv("BOT_NAME", "")
    if botToken == "" {
        log.Fatal("BOT_NAME environment variable is required")
    }

    return &Config{
        BotToken: botToken,
        Debug:    getEnv("DEBUG", "false") == "true",
		ConfigEmail: configEmail,
		BotName: botName,
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}