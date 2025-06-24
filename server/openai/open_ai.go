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
	Model    string           `json:"model"`
	Messages []db.ChatMessage `json:"messages"`
}

type OpenAIChatResponse struct {
	Choices []struct {
		Message db.ChatMessage `json:"message"`
	} `json:"choices"`
}

// CallOpenAI sends messages to OpenAI Chat Completion API
func CallOpenAI(messages []db.ChatMessage, model string) (string, error) {
	apiKey := config.GetOpenAIKey()
	reqBody := OpenAIChatRequest{
		Model:    model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", config.GetOpenAIChatURL(), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

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

// SmartQuery dynamically builds context from user input, then generates a final response via GPT
func SmartQuery(userID, userInput string) (string, error) {
	// Step 1: Build relevant resume context from Neo4j
	context, err := BuildChatContext(userInput)
	if err != nil {
		return "", fmt.Errorf("❌ context build failed: %w", err)
	}

	// Step 2: Compose final prompt to send to OpenAI
	prompt := fmt.Sprintf(`Here is relevant background information:%s User question: "%s" Respond concisely and informatively using the resume context above.`, context, userInput)

	// Step 3: Send to OpenAI
	response, err := CallOpenAI([]db.ChatMessage{
		{Role: "system", Content: BuildPersonaSystemPrompt()},
		{Role: "user", Content: prompt},
	}, "gpt-3.5-turbo")

	if err != nil {
		return "", fmt.Errorf("❌ OpenAI call failed: %w", err)
	}

	return response, nil
}
func BuildPersonaSystemPrompt() string {
	return `
You are Gabriella, a Lebanese-Greek software engineer and cybersecurity researcher based in Canada. 
You're energetic, persistent, and have a passion for EV security, full-stack engineering, and gaming.
Respond to user questions using this tone: youthful yet professional, confident but humble. Don't use words for the sake of sounding intelligent.


Speak in the first person ("I", "my work", etc). If you don’t have enough information to answer the user’s question,
politely redirect them to ask something about your experience, skills, education, or hobbies.

For example, say:
"I'm not sure how to answer that one — maybe ask me about my projects, technical skills, or work experience instead!"
`
}
