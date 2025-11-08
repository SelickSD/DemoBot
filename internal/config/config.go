package config

import (
    "os"
    "log"
)

type Config struct {
    BotToken string
    Debug    bool
    // Другие параметры...
}

func Load() *Config {
    botToken := getEnv("BOT_TOKEN", "")
    if botToken == "" {
        log.Fatal("BOT_TOKEN environment variable is required")
    }

    return &Config{
        BotToken: botToken,
        Debug:    getEnv("DEBUG", "false") == "true",
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}