package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/SelickSD/DemoBot.git/internal/config"
	hell_divers "github.com/SelickSD/DemoBot.git/internal/repository/hell-divers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	cfg    *config.Config
	botAPI *tgbotapi.BotAPI
}

func main() {
	// Загружаем конфиг из переменных окружения
	cfg := config.Load()

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Panic(err)
	}

	demobot := &Bot{
		cfg:    cfg,
		botAPI: bot,
	}

	demobot.botAPI.Debug = cfg.Debug
	log.Printf("Authorized on account %s", demobot.botAPI.Self.UserName)

	// Удаляем активный webhook
	_, err = demobot.botAPI.Request(tgbotapi.DeleteWebhookConfig{})
	if err != nil {
		log.Panic("failed to delete webhook:", err)
	}

	// Обработка graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go demobot.start()

	<-sigChan
	log.Println("Shutting down bot...")
}

func (b *Bot) start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.botAPI.GetUpdatesChan(u)

	for update := range updates {
		go b.handleUpdate(update)
	}
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in handleUpdate: %v", r)
		}
	}()

	if update.Message == nil {
		return
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	var response string
	var err error

	switch strings.ToLower(update.Message.Text) {
	case "за демократию!", "/democracy":
		response, err = b.handleDemocracyCommand()
	case "/start", "/help":
		response = b.handleHelpCommand()
	default:
		response = "Неизвестная команда. Напишите /help для списка команд."
	}

	if err != nil {
		log.Printf("Error handling command: %v", err)
		response = "Произошла ошибка при обработке запроса. Попробуйте позже."
	}

	if response != "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := b.botAPI.Send(msg); err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}

func (b *Bot) handleDemocracyCommand() (string, error) {
	news, err := hell_divers.GetNews(*b.cfg)
	if err != nil {
		return "", err
	}

	return createMessages(news), nil
}

func (b *Bot) handleHelpCommand() string {
	return `Доступные команды:
• "За демократию!" или /democracy - получить последние новости с фронта
• /help - показать это сообщение

За свободу! За управляемую демократию!`
}

func createMessages(news []hell_divers.NewsFeed) string {
	if len(news) == 0 {
		return "Новостей с фронта пока нет. Демократия ждет ваших свершений!"
	}

	// Берем последнюю новость
	latestNews := news[len(news)-1]

	if latestNews.Message == "" {
		return "Получена пустая новость. Возможно, враги демократии вмешались в коммуникации!"
	}

	// Очищаем HTML теги
	result := strings.Replace(latestNews.Message, "<i=1>", "", -1)
	result = strings.Replace(result, "</i>", "", -1)
	result = strings.Replace(result, "<i=3>", "", -1)
	result = strings.Replace(result, "<br>", "\n", -1)

	// Добавляем заголовок если есть контент
	if result != "" {
		result = "📢 СВЕЖИЕ НОВОСТИ С ФРОНТА:\n\n" + result
	}

	return result
}
