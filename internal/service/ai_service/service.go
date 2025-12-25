package ai_service

import "github.com/SelickSD/DemoBot.git/internal/config"

type AiRepo interface {
	PostNewMassage(massage string) string
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

func (s *Service) SendMessage(massage string) string {
	return s.aiRepo.PostNewMassage(massage)
}
