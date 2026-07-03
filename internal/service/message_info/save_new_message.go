package message_info

import (
	"context"

	"github.com/SelickSD/DemoBot.git/internal/logger"
	"github.com/SelickSD/DemoBot.git/internal/repository/messageinfo"
)

func (s *Service) SaveNewMessage(ctx context.Context, msg messageinfo.MessageInfo) error {
	logger.Init()
	logger.Info.Printf("SaveNewMessage: start")

	if err := s.msInfoRepo.Save(ctx, msg); err != nil {
		logger.Error.Printf("SaveNewMessage: %v", err)
		return err
	}
	logger.Info.Printf("SaveNewMessage: end")

	return nil
}
