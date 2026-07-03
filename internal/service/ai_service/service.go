package ai_service

import (
	"github.com/SelickSD/DemoBot.git/internal/config"
	"github.com/SelickSD/DemoBot.git/internal/repository/polza-ai-api/dto"
)

type AiRepo interface {
	PostNewMassage(massage []dto.Message) string
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

func (s *Service) SendMessage(massage []dto.Message) string {
	return s.aiRepo.PostNewMassage(massage)
}
