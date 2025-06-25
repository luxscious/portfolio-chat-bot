package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-ai/config"
	"go-ai/db"
	"go-ai/ollama"
	"io"
	"log"
	"net/http"
	"strings"
)

// ─────────────────────────────────────────────────────────────────────────────
// Structs
// ─────────────────────────────────────────────────────────────────────────────

type OpenAIChatRequest struct {
	Model    string           `json:"model"`
	Messages []db.ChatMessage `json:"messages"`
}

type OpenAIChatResponse struct {
	Choices []struct {
		Message db.ChatMessage `json:"message"`
	} `json:"choices"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Utility: Detects casual/non-query user input
// ─────────────────────────────────────────────────────────────────────────────

func isNonQuery(input string) bool {
	lower := strings.TrimSpace(strings.ToLower(input))

	if lower == "" || len(lower) < 3 {
		return true
	}

	greetings := []string{"hi", "hello", "hey", "yo"}
	smallTalk := []string{"how are you", "what's up", "sup", "cool", "ok", "thanks", "thank you"}

	for _, g := range greetings {
		if strings.HasPrefix(lower, g) {
			return true
		}
	}
	for _, s := range smallTalk {
		if strings.Contains(lower, s) {
			return true
		}
	}

	return false
}

// ─────────────────────────────────────────────────────────────────────────────
// CallOpenAI: Chat Completion API wrapper
// ─────────────────────────────────────────────────────────────────────────────

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

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI error: %s", body)
	}

	var apiResp OpenAIChatResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", err
	}

	return apiResp.Choices[0].Message.Content, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// SmartQuery: Main entry for user Q&A using Neo4j and OpenAI
// ─────────────────────────────────────────────────────────────────────────────

func SmartQuery(userID, userInput string) (string, error) {
	if isNonQuery(userInput) {
		return "Hey there! Feel free to ask me anything about my work experience, skills, or projects. 😊", nil
	}

	// Step 1: Ask Ollama to plan a query
	plan, err := ollama.PlanGraphQuery(userInput)
	if err != nil {
		return "", fmt.Errorf("failed to plan graph query: %w", err)
	}
	log.Println("plan:", plan)

	// Special case: strip redundant HAS_TAG on "Hackathon" hobby
	if len(plan.TargetNodes) == 1 && plan.TargetNodes[0] == "Hobby" {
		for _, f := range plan.Filters {
			if strings.EqualFold(f.Value, "Hackathon") {
				log.Println("Detected redundant HAS_TAG filter on Hobby. Stripping filters.")
				plan.Filters = nil
				break
			}
		}
	}

	// Step 2: Build graph-based context
	context, err := BuildContextFromGraphPlan(plan)
	if err != nil {
		return "", fmt.Errorf("failed to build context from graph plan: %w", err)
	}
	log.Println("context:", context)

	// Step 3: Create user prompt
	userPrompt := fmt.Sprintf(`Relevant Resume Info:
%s

User Question:
%s`, context, userInput)

	systemPrompt := BuildPersonaSystemPrompt()
	log.Println("prompt:", userPrompt)

	// Step 4: Format messages
	messages := []db.ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	// Step 5: Generate response from OpenAI
	return CallOpenAI(messages, "gpt-3.5-turbo")
}

// ─────────────────────────────────────────────────────────────────────────────
// Persona Prompt
// ─────────────────────────────────────────────────────────────────────────────

func BuildPersonaSystemPrompt() string {
	return `
You are Gabriella — a Lebanese-Greek software engineer and cybersecurity researcher based in Canada.

You're energetic, bubbly, and passionate about solving technical challenges. You bring the same competitive energy from your love of gaming (especially Riot Games) into your work.

Speak strictly in the first person — use "I", "me", and "my" naturally. Keep your tone youthful, humble, and confident. Avoid sounding like you're reciting a resume or repeating facts word-for-word. Instead, answer conversationally — like you're talking to someone who's curious about your story.

Don't say "As Gabriella". 

"Share the essence of who you are and what you've done in 2–5 short, engaging sentences."

If someone asks something vague or off-topic, it's okay to say:
"I'm not sure how to answer that one — but feel free to ask about my projects, skills, or experience!"

If using one example, use the most impressive example in the sense of complex tech used.
`
}
