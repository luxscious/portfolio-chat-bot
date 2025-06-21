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
	const extractionTemplate = `
You are {{.Identity.Name}}'s assistant. You have her full resume below in JSON format.

User's Question:
{{.UserInput}}

Resume (JSON):
{{.FullResumeJSON}}

---

Please extract the most relevant items from the resume that would help answer the user's question.

Each item must match this format:

{
  "id": "exact-id-from-json",
  "startDate": "YYYY-MM-DD", // omit if not applicable
  "endDate": "YYYY-MM-DD",
  "name": "Exact name",
  "description": "Why it's relevant to the question",
  "institution": "If applicable",
  "skills": ["tech1", "tech2"],
  "tags": ["Project", "Education", "Trait", ...],
  "type": "project" | "education" | "experience" | "trait",
  "featured": true|false
}

Return a JSON array of all relevant items. Use exact values from the provided JSON resume.
`

	tmpl, err := template.New("resumeExtraction").Parse(extractionTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	// Serialize the full resume data to JSON
	resumeJSON, err := json.MarshalIndent(resumeData, "", "  ")
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&buf, map[string]interface{}{
		"Identity":       resumeData.PersonaContext.Identity,
		"UserInput":      userInput,
		"FullResumeJSON": string(resumeJSON),
	})
	if err != nil {
		return "", err
	}

	return CallOpenAI([]db.ChatMessage{
		{Role: "system", Content: "You are a helpful assistant."},
		{Role: "user", Content: buf.String()},
	}, "gpt-3.5-turbo")
}


func GeneratePersonalAnswer(resumeData *resume.ResumeData, relevantInfo, userInput string) (string, error) {
	systemPrompt, err := BuildSystemPrompt(resumeData.PersonaContext)
	if err != nil {
		return "", err
	}

	userPrompt := fmt.Sprintf(`Relevant Resume Info: %s User Question: %s`, relevantInfo, userInput)

	return CallOpenAI([]db.ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}, "gpt-4o-mini")
}
