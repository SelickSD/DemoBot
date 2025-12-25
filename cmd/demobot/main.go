package main

import (
	"log"

	"github.com/SelickSD/DemoBot.git/internal/app/bot"
	"github.com/SelickSD/DemoBot.git/internal/config"
	diversRepo "github.com/SelickSD/DemoBot.git/internal/repository/hell-divers"
	polzaApi "github.com/SelickSD/DemoBot.git/internal/repository/polza-ai-api"
	"github.com/SelickSD/DemoBot.git/internal/service/ai_service"
	"github.com/SelickSD/DemoBot.git/internal/service/helldivers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.Load()

	botApiClient, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Panic(err)
	}

	divers := diversRepo.NewRepository()
	diversService := helldivers.NewService(cfg, divers)

	aiApiClient := polzaApi.NewAIApiKey(*cfg)
	aiService := ai_service.NewService(cfg, aiApiClient)

	demobot := bot.NewBot(cfg, botApiClient, diversService, aiService)

	demobot.Run()
}
