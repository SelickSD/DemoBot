package main

import (
	"context"
	"log"

	"github.com/SelickSD/DemoBot.git/internal/app/bot"
	"github.com/SelickSD/DemoBot.git/internal/config"
	db2 "github.com/SelickSD/DemoBot.git/internal/db"
	diversRepo "github.com/SelickSD/DemoBot.git/internal/repository/hell-divers"
	msInfoRepo "github.com/SelickSD/DemoBot.git/internal/repository/messageinfo"
	polzaApi "github.com/SelickSD/DemoBot.git/internal/repository/polza-ai-api"
	"github.com/SelickSD/DemoBot.git/internal/service/ai_service"
	"github.com/SelickSD/DemoBot.git/internal/service/helldivers"
	msInfoSvc "github.com/SelickSD/DemoBot.git/internal/service/message_info"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	if err := db2.Migrate(); err != nil {
		log.Fatal(err)
	}

	if err := db2.Init(ctx); err != nil {
		log.Fatal(err)
	}

	botApiClient, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Panic(err)
	}

	divers := diversRepo.NewRepository()
	diversService := helldivers.NewService(cfg, divers)

	aiApiClient := polzaApi.NewAIApiKey(*cfg)
	aiService := ai_service.NewService(cfg, aiApiClient)

	messageInfoRepo := msInfoRepo.NewRepository()
	messageInfoService := msInfoSvc.NewService(messageInfoRepo)

	demobot := bot.NewBot(
		cfg,
		botApiClient,
		diversService,
		aiService,
		messageInfoService,
	)

	//Run
	demobot.Run()
}
