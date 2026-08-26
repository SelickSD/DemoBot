package helldivers

import (
	"context"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/SelickSD/DemoBot.git/internal/config"
	"github.com/SelickSD/DemoBot.git/internal/enum"
	"github.com/SelickSD/DemoBot.git/internal/logger"
	"github.com/SelickSD/DemoBot.git/internal/repository/game_news_reports"
	"github.com/SelickSD/DemoBot.git/internal/repository/hell-divers/dto"
	dtoPolza "github.com/SelickSD/DemoBot.git/internal/repository/polza-ai-api/dto"
)

type DiversRepo interface {
	GetNews(config config.Config) ([]dto.NewsFeed, error)
}

type newsReportRepo interface {
	CreateGameNewsReport(
		ctx context.Context,
		gameID int64,
		report string,
	) error
	GetLatestGameNewsReport(
		ctx context.Context,
		gameID int64,
	) (*game_news_reports.GameNewsReport, error)
}

type aiRepo interface {
	PostNewMassage(massage []dtoPolza.Message, prompt string) string
}

type Service struct {
	cfg            *config.Config
	diversRepo     DiversRepo
	newsReportRepo newsReportRepo
	aiRepo         aiRepo
}

func NewService(
	cfg *config.Config,
	diversRepo DiversRepo,
	newsReportRepo newsReportRepo,
	aiRepo aiRepo,
) *Service {
	return &Service{
		cfg:            cfg,
		diversRepo:     diversRepo,
		newsReportRepo: newsReportRepo,
		aiRepo:         aiRepo,
	}
}

func (s *Service) GetLatestNews(ctx context.Context) (string, error) {
	news, err := s.diversRepo.GetNews(*s.cfg)
	if err != nil {
		return "", err
	}
	return s.createMessages(ctx, news), nil
}

func (s *Service) createMessages(ctx context.Context, news []dto.NewsFeed) string {
	if len(news) == 0 {
		return "Новостей с фронта пока нет. Демократия ждет ваших свершений!"
	}
	latestNews := make([]dtoPolza.Message, 0, 12)
	newsIdArray := make([]int, len(news))

	for i := range news {
		newsIdArray[i] = news[i].Id
	}

	newsIdArray = top10(newsIdArray)

	for _, value := range news {
		if slices.Contains(newsIdArray, value.Id) {

			result := strings.Replace(value.Message, "<i=1>", "", -1)
			result = strings.Replace(result, "</i>", "", -1)
			result = strings.Replace(result, "<i=3>", "", -1)
			result = strings.Replace(result, "<br>", "\n", -1)

			latestNews = append(latestNews, dtoPolza.Message{
				Role:    "user",
				Content: "Message ID: " + strconv.Itoa(value.Id) + " " + result,
			})
		}
	}

	lastReport, err := s.newsReportRepo.GetLatestGameNewsReport(ctx, int64(enum.HELLDIVERS2))
	if err != nil {
		logger.Info.Println("Ошибка получения последнего отчета из репозитория newsReportRepo")
	}

	if lastReport != nil {
		latestNews = append(latestNews, dtoPolza.Message{
			Role:    "user",
			Content: "LastReport: " + lastReport.Report,
		})
	}

	if len(latestNews) == 0 {
		return "Получена пустая новость. Возможно, враги демократии вмешались в коммуникации!"
	}

	latestNews = append(latestNews, dtoPolza.Message{
		Role:    "user",
		Content: "Prompt: " + enum.HELLDIVERSPROMPT,
	})

	result := s.aiRepo.PostNewMassage(latestNews, "HELLDIVERS2")

	if result == "" {
		return "Новостей с фронта пока нет. Демократия ждет ваших свершений!"
	}

	err = s.newsReportRepo.CreateGameNewsReport(ctx, int64(enum.HELLDIVERS2), result)
	if err != nil {
		logger.Info.Println(err)
	}

	return result
}

func top10(values []int) []int {
	sort.Sort(sort.Reverse(sort.IntSlice(values)))

	if len(values) > 10 {
		values = values[:10]
	}

	return values
}
