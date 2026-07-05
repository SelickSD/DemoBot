package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/SelickSD/DemoBot.git/internal/config"
	"github.com/SelickSD/DemoBot.git/internal/logger"
	"github.com/SelickSD/DemoBot.git/internal/repository/messageinfo"
	"github.com/SelickSD/DemoBot.git/internal/repository/polza-ai-api/dto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HellDiversService interface {
	GetLatestNews() (string, error)
}

type AiService interface {
	SendMessage(massage []dto.Message) string
}

type MessageService interface {
	SaveNewMessage(ctx context.Context, msg messageinfo.MessageInfo) error
	GetByChatID(
		ctx context.Context,
		chatID int64,
		limit int,
	) ([]messageinfo.MessageInfo, error)
	DellAll(ctx context.Context) error
}

type Bot struct {
	cfg           *config.Config
	botApiClient  *tgbotapi.BotAPI
	diversService HellDiversService
	aiService     AiService
	msInfoService MessageService
}

func NewBot(
	cfg *config.Config,
	botApiClient *tgbotapi.BotAPI,
	diversService HellDiversService,
	aiService AiService,
	msInfoService MessageService,
) *Bot {
	return &Bot{
		cfg:           cfg,
		botApiClient:  botApiClient,
		diversService: diversService,
		aiService:     aiService,
		msInfoService: msInfoService,
	}
}

func (b *Bot) Run() {
	b.botApiClient.Debug = b.cfg.Debug
	logger.Info.Printf("Authorized on account %s", b.botApiClient.Self.UserName)

	// Удаляем активный webhook
	_, err := b.botApiClient.Request(tgbotapi.DeleteWebhookConfig{})
	if err != nil {
		logger.Error.Panic("failed to delete webhook:", err)
	}

	// Обработка graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go b.start()

	<-sigChan
	logger.Info.Println("Shutting down bot...")
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
			logger.Info.Printf("Recovered from panic in handleUpdate: %v", r)
		}
	}()

	if update.Message == nil {
		return
	}

	logger.Info.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	var response string
	var err error

	switch strings.ToLower(update.Message.Text) {
	case "за демократию!", "/democracy":
		response, err = b.handleDemocracyCommand()
	case "/start", "/help":
		response = b.handleHelpCommand()
	case "delete all":
		err = b.msInfoService.DellAll(context.Background())
		if err != nil {
			logger.Error.Printf("Error deleting all messages: %v", err)
		}
		response = "All messages have been deleted."
	default:
		isReplay := isReplayToBot(update)
		if isBotCommand(strings.ToLower(update.Message.Text)) || isReplay {
			ctx := context.Background()
			messageWithContext := b.prepareNewMassage(ctx, update.Message.Chat.ID)
			err = b.saveNewMassage(update, "")
			if err != nil {
				logger.Error.Printf("Error saving new massage: %v", err)
			}

			actualMessage := update.Message.Text

			if !isReplay {
				actualMessage = extractBotMessage(strings.ToLower(update.Message.Text))
			}

			if len(messageWithContext) == 0 {
				response = b.aiService.SendMessage([]dto.Message{{
					Role:    "user",
					Content: fmt.Sprintf("User: %d, MessageID: %d, NewMessage: %s", update.Message.From.ID, update.Message.MessageID, actualMessage),
				}})
			} else {
				messageWithContext = append(messageWithContext, dto.Message{
					Role:    "user",
					Content: fmt.Sprintf("User: %d, MessageID: %d, NewMessage: %s", update.Message.From.ID, update.Message.MessageID, actualMessage),
				})

				response = b.aiService.SendMessage(messageWithContext)
			}
		}
		break
	}

	if err != nil {
		logger.Info.Printf("Error handling command: %v", err)
		response = "Произошла ошибка при обработке запроса. Попробуйте позже."
	}

	if response != "" {
		// Разбиваем сообщение, если оно слишком длинное
		messages := splitMessage(response, 4096)
		for _, msgText := range messages {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			msg.ReplyToMessageID = update.Message.MessageID

			if msgText != "All messages have been deleted." {
				err = b.saveNewMassage(update, msgText)
				if err != nil {
					logger.Error.Printf("Error saving new massage: %v", err)
				}
			}

			if _, err := b.botApiClient.Send(msg); err != nil {
				logger.Info.Printf("Error sending message part: %v", err)
			}
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

// splitMessage разбивает строку на части заданной максимальной длины
func splitMessage(text string, maxSize int) []string {
	var result []string

	for len(text) > maxSize {
		// Пытаемся разбить по последнему пробелу в пределах maxSize, чтобы не резать слова
		breakIndex := maxSize
		if i := strings.LastIndex(text[:maxSize], " "); i > 0 {
			breakIndex = i
		}
		result = append(result, text[:breakIndex])
		text = strings.TrimSpace(text[breakIndex:])
	}
	if len(text) > 0 {
		result = append(result, text)
	}

	return result
}

func (b *Bot) saveNewMassage(update tgbotapi.Update, replyMessage string) error {
	message := update.Message.Text

	if replyMessage != "" {
		message = fmt.Sprintf("Reply from message id: %d, ChatID: %d, Reply text: %s", update.Message.MessageID, update.Message.Chat.ID, replyMessage)
	}

	ctx := context.Background()
	newMessage := messageinfo.MessageInfo{
		ID:        int64(update.UpdateID),
		MessageID: int64(update.Message.MessageID),
		ChatID:    update.Message.Chat.ID,
		Message:   message,
		UserID:    update.Message.From.ID,
		CreatedAt: time.Now(),
	}

	err := b.msInfoService.SaveNewMessage(ctx, newMessage)
	if err != nil {
		logger.Info.Printf("Error saving new message: %v", err)
		return fmt.Errorf("save new message: %w", err)
	}

	return nil
}

func (b *Bot) prepareNewMassage(ctx context.Context, chatID int64) []dto.Message {
	messageContext, err := b.msInfoService.GetByChatID(ctx, chatID, 10)
	if err != nil {
		return nil
	}

	if len(messageContext) == 0 {
		return nil
	}

	result := make([]dto.Message, 0, len(messageContext)+1) // +1 for the new message

	for _, message := range messageContext {
		var text string

		if strings.Contains(message.Message, "Reply from message id") {
			text = message.Message
		} else {
			text = fmt.Sprintf("User: %d, MessageID: %d, Text: %s", message.UserID, message.MessageID, text)
		}

		result = append(result, dto.Message{
			Role:    "user",
			Content: text,
		})
	}

	logger.Info.Printf("%s", result)

	return result
}

func isReplayToBot(update tgbotapi.Update) bool {
	if update.Message != nil && update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From != nil {
		switch update.Message.ReplyToMessage.From.UserName {
		case "SelickBot", "SuperDemocracyBot":
			return true
		default:
			return false
		}
	}
	return false
}
