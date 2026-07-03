package message_info

import (
	"context"

	"github.com/SelickSD/DemoBot.git/internal/repository/messageinfo"
)

type MessageInfoRepo interface {
	Save(ctx context.Context, msg messageinfo.MessageInfo) error
	DeleteAll(ctx context.Context) error
}

type Service struct {
	msInfoRepo MessageInfoRepo
}

func NewService(msInfoRepo MessageInfoRepo) *Service {
	return &Service{
		msInfoRepo: msInfoRepo,
	}
}

func (s *Service) DellAll(ctx context.Context) error {
	return s.msInfoRepo.DeleteAll(ctx)
}
