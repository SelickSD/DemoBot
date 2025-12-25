package bot

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/SelickSD/DemoBot.git/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HellDiversService interface {
	GetLatestNews() (string, error)
}

type AiService interface {
	SendMessage(massage string) string
}

type Bot struct {
	cfg           *config.Config
	botApiClient  *tgbotapi.BotAPI
	diversService HellDiversService
	aiService     AiService
}

func NewBot(
	cfg *config.Config,
	botApiClient *tgbotapi.BotAPI,
	diversService HellDiversService,
	aiService AiService,
) *Bot {
	return &Bot{
		cfg:           cfg,
		botApiClient:  botApiClient,
		diversService: diversService,
		aiService:     aiService,
	}
}

func (b *Bot) Run() {
	b.botApiClient.Debug = b.cfg.Debug
	log.Printf("Authorized on account %s", b.botApiClient.Self.UserName)

	// Удаляем активный webhook
	_, err := b.botApiClient.Request(tgbotapi.DeleteWebhookConfig{})
	if err != nil {
		log.Panic("failed to delete webhook:", err)
	}

	// Обработка graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go b.start()

	<-sigChan
	log.Println("Shutting down bot...")
}

func (b *Bot) start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.botApiClient.GetUpdatesChan(u)

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
		if isBotCommand(strings.ToLower(update.Message.Text)) {
			actualMessage := extractBotMessage(strings.ToLower(update.Message.Text))
			response = b.aiService.SendMessage(actualMessage)
		}
		break
	}

	if err != nil {
		log.Printf("Error handling command: %v", err)
		response = "Произошла ошибка при обработке запроса. Попробуйте позже."
	}

	if response != "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := b.botApiClient.Send(msg); err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}

func (b *Bot) handleDemocracyCommand() (string, error) {
	news, err := b.diversService.GetLatestNews()
	if err != nil {
		return "", err
	}

	return news, nil
}

func (b *Bot) handleHelpCommand() string {
	return `Доступные команды:
• "За демократию!" или /democracy - получить последние новости с фронта
• /help - показать это сообщение

За свободу! За управляемую демократию!`
}

func isBotCommand(message string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(message)), "bot")
}

func extractBotMessage(message string) string {
	message = strings.TrimSpace(message)
	if strings.HasPrefix(strings.ToLower(message), "bot") {
		// Убираем "bot" и лишние запятые/пробелы
		cleaned := strings.TrimSpace(message[3:])
		cleaned = strings.TrimPrefix(cleaned, ",")
		cleaned = strings.TrimPrefix(cleaned, " ")
		return strings.TrimSpace(cleaned)
	}
	return ""
}
