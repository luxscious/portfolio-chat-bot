package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-ai/config"
	"go-ai/db"
	"io"
	"net/http"
)

type OpenAIChatRequest struct {
	Model    string        `json:"model"`
	Messages []db.ChatMessage `json:"messages"`
}

type OpenAIChatResponse struct {
	Choices []struct {
		Message db.ChatMessage `json:"message"`
	} `json:"choices"`
}


// callOpenAI sends messages to OpenAI Chat Completion API
func CallOpenAI(messages []db.ChatMessage) (string, error) {
	apiKey := config.GetOpenAIKey()
	reqBody := OpenAIChatRequest{
		Model:    "gpt-4", 
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Make request
	req, err := http.NewRequest("POST", config.GetOpenAIChatURL() , bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ apiKey)

	// Send Request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("OpenAI error: %s", body)
	}

	var apiResp OpenAIChatResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", err
	}

	return apiResp.Choices[0].Message.Content, nil
}
