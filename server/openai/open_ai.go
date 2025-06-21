package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-ai/config"
	"go-ai/db"
	"go-ai/resume"
	"io"
	"net/http"
	"text/template"
)

type OpenAIChatRequest struct {
	Model    string            `json:"model"`
	Messages []db.ChatMessage  `json:"messages"`
}

type OpenAIChatResponse struct {
	Choices []struct {
		Message db.ChatMessage `json:"message"`
	} `json:"choices"`
}
// Helper Functions
func BuildSystemPrompt(p resume.PersonaContext) (string, error) {
	tmpl, err := template.New("prompt").Parse(p.PromptTemplate)
	if err != nil {
		return "", fmt.Errorf("template parsing failed: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, p); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	return buf.String(), nil
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

// SmartQuery performs a 2-step GPT process: filter → generate
func SmartQuery(resumeData *resume.ResumeData, userInput string) (string, error) {
	// Step 1: extract relevant resume sections
	relevantInfo, err := ExtractRelevantResumeInfo(resumeData, userInput)
	if err != nil {
		return "", fmt.Errorf("step 1 (filter) failed: %w", err)
	}

	// Step 2: generate the final answer using persona tone
	finalAnswer, err := GeneratePersonalAnswer(resumeData, relevantInfo, userInput)
	if err != nil {
		return "", fmt.Errorf("step 2 (generate) failed: %w", err)
	}

	return finalAnswer, nil
}

// ExtractRelevantResumeInfo uses GPT to identify the most relevant resume entries
func ExtractRelevantResumeInfo(resumeData *resume.ResumeData, userInput string) (string, error) {
	const filterTemplate = `
You are {{.Identity.Name}}'s assistant. You have her full resume below.

User's Question:
{{.UserInput}}

Resume:
{{.FlatResume}}

Please extract the most relevant work experiences, projects, education, or personal traits that would help answer the question. Return a concise bullet-point summary of those items.
`

	tmpl, err := template.New("filter").Parse(filterTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	flatResume := resumeData.FlattenResume()
	err = tmpl.Execute(&buf, map[string]interface{}{
		"Identity":   resumeData.PersonaContext.Identity,
		"FlatResume": flatResume,
		"UserInput":  userInput,
	})
	if err != nil {
		return "", err
	}

	return CallOpenAI([]db.ChatMessage{
		{Role: "system", Content: "You are a helpful assistant."},
		{Role: "user", Content: buf.String()},
	}, "gpt-3.5-turbo")
}

// GeneratePersonalAnswer crafts the final GPT response in Gabriella's tone
func GeneratePersonalAnswer(resumeData *resume.ResumeData, relevantInfo, userInput string) (string, error) {
	systemPrompt, err := BuildSystemPrompt(resumeData.PersonaContext)
	if err != nil {
		return "", err
	}

	userPrompt := fmt.Sprintf(`Relevant Resume Info:
%s

User Question:
%s`, relevantInfo, userInput)

	return CallOpenAI([]db.ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}, "gpt-4")
}
