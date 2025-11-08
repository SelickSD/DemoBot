package main

import (
	"log"
	"strings"

	"github.com/SelickSD/DemoBot.git/internal/config"
	hell_divers "github.com/SelickSD/DemoBot.git/internal/repository/hell-divers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	 // Загружаем конфиг из переменных окружения
    cfg := config.Load()

    bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
    if err != nil {
        log.Panic(err)
    }

	bot.Debug = cfg.Debug
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
		
		result := strings.Replace(news[count-1].Message, "<i=1>", "", -1)  
		result = strings.Replace(result, "</i>", "", -1)
		result = strings.Replace(result, "<i=3>", "", -1)   

		return result
	}
	return ""
}
