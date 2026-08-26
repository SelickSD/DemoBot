package ai_service

import (
	"github.com/SelickSD/DemoBot.git/internal/config"
	"github.com/SelickSD/DemoBot.git/internal/repository/polza-ai-api/dto"
)

type AiRepo interface {
	PostNewMassage(massage []dto.Message, prompt string) string
}

type Service struct {
	cfg    *config.Config
	aiRepo AiRepo
}

func NewService(cfg *config.Config, aiRepo AiRepo) *Service {
	return &Service{
		cfg:    cfg,
		aiRepo: aiRepo,
	}
}

func (s *Service) SendMessage(massage []dto.Message, prompt string) string {
	const (
		defaultPrompt = "Общение ведется в Телеграмме, на русском языке, ответы нужно формировать в дружеской форме. Новое сообщение помечено как NewMessage, нужно ответить на него"
	)

	var actualPrompt string

	switch prompt {
	case "HELLDIVERS2":
		actualPrompt = "Prompt: Общение ведется в Телеграмме, на русском языке. " +
			"В рамках данного сообщения тебе переданы последние 10 игровых новостей из игры HellDivers 2. Каждое сообщение помечено Message ID. " +
			"Проведи анализ и сформируй отчет о ситуации на галактической карте, сообщение помеченное LastReport содержит твой предыдущий ответ. " +
			"Сформируй итоговый отчет который должен содержать Главный приказ, если он не выполнен, если выполнен - результаты и актуальные новости. " +
			"По возможности создай ответ в присущей данной игре пафосной и шуточной форме."
	default:
		actualPrompt = defaultPrompt
	}

	return s.aiRepo.PostNewMassage(massage, actualPrompt)
}
