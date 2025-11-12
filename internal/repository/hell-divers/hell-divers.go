package hell_divers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/SelickSD/DemoBot.git/internal/config"
	"github.com/SelickSD/DemoBot.git/internal/logger"
)

type NewsFeed struct {
	Id        int    `json:"id"`
	Published int    `json:"published"`
	Type      int    `json:"type"`
	Message   string `json:"message"`
}

func GetNews(config config.Config) ([]NewsFeed, error) {
	baseURL := "https://api.helldivers2.dev/raw/api/NewsFeed/801"

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	var allItems []NewsFeed

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		logger.Error.Printf("Ошибка создания запроса: %v", err)
		return nil, err
	}

	req.Header.Add("X-Super-Client", config.BotName)
	req.Header.Add("X-Super-Contact", config.ConfigEmail)
	req.Header.Add("Accept-Language", "ru-RU")

	maxRetries := 5
	retryDelay := time.Second * 2
	var response []NewsFeed

	for retries := 0; retries < maxRetries; retries++ {
		logger.Info.Printf("Попытка запроса %d из %d", retries+1, maxRetries)

		resp, err := client.Do(req)
		if err != nil {
			logger.Error.Printf("Ошибка выполнения запроса: %v", err)
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error.Printf("Ошибка чтения ответа: %v", err)
			return nil, err
		}

		if resp.StatusCode == http.StatusServiceUnavailable {
			logger.Info.Printf("Сервер вернул 503. Попытка %d из %d. Жду %s...", retries+1, maxRetries, retryDelay)
			time.Sleep(retryDelay)
			retryDelay *= 2 // Экспоненциальная задержка
			continue
		}

		if resp.StatusCode != http.StatusOK {
			logger.Error.Printf("Ошибка: статус %d\nОтвет: %s", resp.StatusCode, string(body))
			return nil, fmt.Errorf("API вернул статус: %d", resp.StatusCode)
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			logger.Error.Printf("Ошибка парсинга JSON: %v", err)
			return nil, err
		}

		logger.Info.Printf("Успешно получено %d новостей", len(response))
		break
	}

	allItems = append(allItems, response...)
	return allItems, nil
}
