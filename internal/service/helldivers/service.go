package helldivers

import (
	"strings"

	"github.com/SelickSD/DemoBot.git/internal/config"
	"github.com/SelickSD/DemoBot.git/internal/repository/hell-divers/dto"
)

type DiversRepo interface {
	GetNews(config config.Config) ([]dto.NewsFeed, error)
}

type Service struct {
	cfg        *config.Config
	diversRepo DiversRepo
}

func NewService(
	cfg *config.Config,
	diversRepo DiversRepo,
) *Service {
	return &Service{
		cfg:        cfg,
		diversRepo: diversRepo,
	}
}

func (s *Service) GetLatestNews() (string, error) {
	news, err := s.diversRepo.GetNews(*s.cfg)
	if err != nil {
		return "", err
	}
	return createMessages(news), nil
}

func createMessages(news []dto.NewsFeed) string {
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
