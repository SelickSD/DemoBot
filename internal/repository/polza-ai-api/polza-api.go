package polzaaiapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/SelickSD/DemoBot.git/internal/config"
	"github.com/SelickSD/DemoBot.git/internal/repository/polza-ai-api/dto"
)

type AiApyClient struct {
	cfg config.Config
}

func NewAIApiKey(cfg config.Config) *AiApyClient {
	return &AiApyClient{
		cfg: cfg,
	}
}

func (c *AiApyClient) PostNewMassage(massage string) string {
	request := dto.ChatRequest{
		Model: "deepseek/deepseek-chat-v3.1",
		Messages: []dto.Message{
			{Role: "user", Content: massage},
		},
		Temperature: 0.7,
		MaxTokens:   1500,
	}

	resp, err := createChatCompletion(c.cfg.AiApiKey, request)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Ответ: %s\n", resp.Choices[0].Message.Content)
	fmt.Printf("Стоимость: %.2f руб.\n", resp.Usage.Cost)

	return resp.Choices[0].Message.Content
}

func createChatCompletion(apiKey string, request dto.ChatRequest) (*dto.ChatResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.polza.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var chatResp dto.ChatResponse
	err = json.Unmarshal(body, &chatResp)
	if err != nil {
		return nil, err
	}

	return &chatResp, nil
}
