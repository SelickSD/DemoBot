package hell_divers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type NewsFeed struct {
	Id        int    `json:"id"`
	Published int    `json:"published"`
	Type      int    `json:"type"`
	Message   string `json:"message"`
}

func GetNews() ([]NewsFeed, error) {
	baseURL := "https://api.helldivers2.dev/raw/api/NewsFeed/801"

	client := &http.Client{}
	var allItems []NewsFeed

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return nil, err
	}

	req.Header.Add("X-Super-Client", "@SuperDemocracyBot")
	req.Header.Add("X-Super-Contact", "selicksd@gmail.com")

	maxRetries := 5
	retryDelay := time.Second
	var response []NewsFeed

	for retries := 0; retries < maxRetries; retries++ {
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Ошибка выполнения запроса:", err)
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Ошибка чтения ответа:", err)
			return nil, err
		}

		if resp.StatusCode == http.StatusServiceUnavailable {
			fmt.Printf("Сервер вернул 503. Попытка %d из %d. Жду %s...\n", retries+1, maxRetries, retryDelay)
			time.Sleep(retryDelay)
			continue // пробуем снова
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Ошибка: статус %d\nОтвет: %s\n", resp.StatusCode, string(body))
			return nil, fmt.Errorf("Status Error")
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		break
	}

	allItems = append(allItems, response...)
	return allItems, nil
}
