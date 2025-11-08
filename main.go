package main

import (
	"github.com/SelickSD/DemoBot.git/internal/config"
	hell_divers "github.com/SelickSD/DemoBot.git/internal/repository/hell-divers"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	//cfg, err := config.LoadConfig(".cicd/deploy/helm-envs/config.yaml")
	//if err != nil {
	//	log.Fatalf("Ошибка загрузки конфига: %v", err)
	//}
	
	bot, err := tgbotapi.NewBotAPI(config.Bot_token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Удаляем активный webhook
	_, err = bot.Request(tgbotapi.DeleteWebhookConfig{})
	if err != nil {
		log.Panic("failed to delete webhook:", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.Text == "За демократию!" {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			news, err := hell_divers.GetNews()
			if err != nil {
				log.Panic(err)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, createMassages(news))
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

func createMassages(news []hell_divers.NewsFeed) string {
	count := len(news)
	if news[count-1].Message != "" {
		return news[count-1].Message
	}
	return ""
}
