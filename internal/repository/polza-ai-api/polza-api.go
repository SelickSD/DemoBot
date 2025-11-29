package polzaaiapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/SelickSD/DemoBot.git/internal/config"
)

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatRequest struct {
    Model       string    `json:"model"`
    Messages    []Message `json:"messages"`
    MaxTokens   int       `json:"max_tokens,omitempty"`
    Temperature float64   `json:"temperature,omitempty"`
    Stream      bool      `json:"stream,omitempty"`
}

type Usage struct {
    PromptTokens     int     `json:"prompt_tokens"`
    CompletionTokens int     `json:"completion_tokens"`
    TotalTokens      int     `json:"total_tokens"`
    Cost             float64 `json:"cost,omitempty"`
}

type Choice struct {
    Index        int    `json:"index"`
    Message      Message `json:"message,omitempty"`
    FinishReason string `json:"finish_reason"`
}

type ChatResponse struct {
    ID      string   `json:"id"`
    Object  string   `json:"object"`
    Created int64    `json:"created"`
    Model   string   `json:"model"`
    Choices []Choice `json:"choices"`
    Usage   Usage    `json:"usage"`
}

type AiApyClient struct{
	cfg config.Config
}

func NewAIApiKey(cfg config.Config) *AiApyClient {
return &AiApyClient{
	cfg: cfg,
}
}

func createChatCompletion(apiKey string, request ChatRequest) (*ChatResponse, error) {
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

    var chatResp ChatResponse
    err = json.Unmarshal(body, &chatResp)
    if err != nil {
        return nil, err
    }

    return &chatResp, nil
}

// func streamChatCompletion(apiKey string, request ChatRequest) error {
//     request.Stream = true
//     jsonData, err := json.Marshal(request)
//     if err != nil {
//         return err
//     }

//     req, err := http.NewRequest("POST", "https://api.polza.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
//     if err != nil {
//         return err
//     }

//     req.Header.Set("Authorization", "Bearer "+apiKey)
//     req.Header.Set("Content-Type", "application/json")
//     req.Header.Set("Accept", "text/event-stream")

//     client := &http.Client{}
//     resp, err := client.Do(req)
//     if err != nil {
//         return err
//     }
//     defer resp.Body.Close()

//     scanner := bufio.NewScanner(resp.Body)
//     for scanner.Scan() {
//         line := scanner.Text()
//         if strings.HasPrefix(line, "data: ") {
//             data := line[6:]
//             if data == "[DONE]" {
//                 break
//             }

//             var chunk map[string]interface{}
//             if err := json.Unmarshal([]byte(data), &chunk); err == nil {
//                 if choices, ok := chunk["choices"].([]interface{}); ok && len(choices) > 0 {
//                     choice := choices[0].(map[string]interface{})
//                     if delta, ok := choice["delta"].(map[string]interface{}); ok {
//                         if content, ok := delta["content"].(string); ok {
//                             fmt.Print(content)
//                         }
//                     }
//                 }

//                 // Финальная статистика
//                 if usage, ok := chunk["usage"].(map[string]interface{}); ok {
//                     fmt.Printf("\n\nСтатистика:\n")
//                     fmt.Printf("Токены: %.0f\n", usage["total_tokens"])
//                     if cost, ok := usage["cost"]; ok {
//                         fmt.Printf("Стоимость: %.2f руб.\n", cost)
//                     }
//                 }
//             }
//         }
//     }

//     return scanner.Err()
// }

func (c *AiApyClient) PostNewMassage(massage string) string {
    
    // Обычный запрос
    request := ChatRequest{
        Model: "deepseek/deepseek-chat-v3.1",
        Messages: []Message{
            {Role: "user", Content: massage},
        },
        Temperature: 0.7,
        MaxTokens: 1000,
    }

    resp, err := createChatCompletion(c.cfg.AiApiKey, request)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Ответ: %s\n", resp.Choices[0].Message.Content)
    fmt.Printf("Стоимость: %.2f руб.\n", resp.Usage.Cost)

	return resp.Choices[0].Message.Content

    // Streaming запрос
    // fmt.Printf("\n\nStreaming пример:\n")
    // streamRequest := ChatRequest{
    //     Model: "anthropic/claude-3-5-sonnet",
    //     Messages: []Message{
    //         {Role: "user", Content: "Напиши короткую историю про программиста"},
    //     },
    // }

    // err = streamChatCompletion(apiKey, streamRequest)
    // if err != nil {
    //     panic(err)
    // }
}