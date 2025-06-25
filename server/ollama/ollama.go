package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-ai/config"
	"go-ai/db"
	"io"
	"log"
	"net/http"
	"strings"
)

// ─────────────────────────────────────────────────────────────────────────────
// TYPES
// ─────────────────────────────────────────────────────────────────────────────

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type GraphQueryPlan struct {
	TargetNodes []string          `json:"target_nodes"`
	Filters     []db.FilterClause `json:"filters"`
	RawInput    string            `json:"raw_input"`
}

// ─────────────────────────────────────────────────────────────────────────────
// API WRAPPER
// ─────────────────────────────────────────────────────────────────────────────

// SendPrompt sends a prompt to the local Ollama server and returns the string response.
func SendPrompt(prompt string) (string, error) {
	reqBody, err := json.Marshal(OllamaRequest{
		Model:  "llama3",
		Prompt: prompt,
		Stream: false,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal Ollama request: %w", err)
	}

	resp, err := http.Post(config.GetOllamaURI()+"/api/generate", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("Ollama POST request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read Ollama response body: %w", err)
	}

	var result OllamaResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("Unmarshal failed: %w\nRaw body: %s", err, string(body))
	}

	return result.Response, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// PROMPT GENERATION
// ─────────────────────────────────────────────────────────────────────────────

// BuildGraphPlannerPrompt dynamically creates a schema-aware graph planning prompt.
func BuildGraphPlannerPrompt(schema db.GraphSchema, userQuery string) string {
	nodeSection := strings.Join(schema.NodeLabels, "\n- ")
	relSection := strings.Join(schema.Relationships, "\n- ")

	return fmt.Sprintf(`
You are a graph planner for a chatbot that answers questions about Gabriella's resume and experience.

The knowledge graph includes the following:

NODE TYPES:
- %s

RELATIONSHIPS:
- %s

TASK:
Given a natural language question, return a graph query plan using this format:

{
  "target_nodes": ["Project"],
  "filters": [
    {
      "on": "Tag",
      "value": "Frontend",
      "relation": "HAS_TAG"
    }
  ],
  "reasoning": "The user is asking about frontend work, which relates to projects tagged as 'Frontend'"
}

GUIDELINES:
- Use only valid node and relationship types from the schema above.
- Do not return "Person" unless the user is directly asking about Gabriella herself.
- Do not include filters with "value": null or "*".
- Output only a single JSON object. No markdown, no commentary, no alternatives.
- If the query references something ambiguous (like "Val-T" or "Hyperpad"), include both "Project" and "WorkExperience" with no filters.
- Use Tag filters for implied categories (e.g., "Hackathons", "Frontend") even if not stated as tags.
- For explicit references (e.g., “Tell me about Val-T”), use a name-based filter, not a Tag.
- Prefer these filter relationships:
  - HAS_TAG
  - HAS_SKILL
  - HAS_HOBBY
- Ignore or remap any other relationship types to the above.
- If unsure, return broad results with empty filters.

QUESTION:
%s
`, nodeSection, relSection, userQuery)
}

// ─────────────────────────────────────────────────────────────────────────────
// PLANNER LOGIC
// ─────────────────────────────────────────────────────────────────────────────

// PlanGraphQuery builds a structured graph query plan from the user's input.
func PlanGraphQuery(userInput string) (GraphQueryPlan, error) {
	prompt := BuildGraphPlannerPrompt(db.CachedSchema, userInput)

	rawResp, err := SendPrompt(prompt)
	if err != nil {
		return GraphQueryPlan{}, err
	}

	parsed, err := ParseIntentResponse(rawResp)
	if err != nil {
		return GraphQueryPlan{}, err
	}

	parsed.RawInput = userInput
	return parsed, nil
}

// ParseIntentResponse extracts the JSON graph plan from a raw LLM response.
func ParseIntentResponse(response string) (GraphQueryPlan, error) {
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")

	if start == -1 || end == -1 || end <= start {
		log.Println("OLLAMA RAW RESPONSE (unusable):", response)
		return GraphQueryPlan{}, fmt.Errorf("Ollama returned no valid JSON block: %s", response)
	}

	jsonPart := response[start : end+1]
	log.Println("OLLAMA JSON BLOCK:", jsonPart)

	var parsed GraphQueryPlan
	if err := json.Unmarshal([]byte(jsonPart), &parsed); err != nil {
		return GraphQueryPlan{}, fmt.Errorf("failed to parse Ollama JSON: %w\nRaw JSON: %s", err, jsonPart)
	}
	return parsed, nil
}
