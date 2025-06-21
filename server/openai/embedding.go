package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-ai/config"
	"go-ai/resume"
	"io"
	"log"
	"net/http"
	"slices"
	"sort"
	"strings"
)

type EmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

type EmbeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}
type EmbeddedChunk struct {
	Text      string
	Embedding []float64
}

var resumeChunks []EmbeddedChunk

func SearchRelevantChunks(query string) ([]string, []float64, error) {
	log.Println("🔍 Generating embedding for user query:", query)
	queryVec, err := GenerateEmbedding(query)
	if err != nil {
		return nil, nil, err
	}

	type scoredChunk struct {
		Text  string
		Score float64
	}
	var results []scoredChunk

	log.Printf("📦 Comparing against %d resume chunks...", len(resumeChunks))
	for _, chunk := range resumeChunks {
		score := CosineSim(queryVec, chunk.Embedding)
		results = append(results, scoredChunk{chunk.Text, score})
	}

	if len(results) == 0 {
		log.Println("⚠️ No chunks to compare.")
		return nil, nil, nil
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	top := 3
	if len(results) < 3 {
		top = len(results)
	}

	var topChunks []string
	var topScores []float64
	for i := 0; i < top; i++ {
		topChunks = append(topChunks, results[i].Text)
		topScores = append(topScores, results[i].Score)
	}

	log.Println("📥 Retrieved Chunks:", topChunks)
	log.Println("📊 Chunk Scores:", topScores)
	return topChunks, topScores, nil
}


// LoadAndEmbedResumeChunks takes ResumeData and sends each chunk for embedding
func LoadAndEmbedResumeChunks(data *resume.ResumeData) error {
	var chunks []string

	// 🔹 Persona Summary
	if data.PersonaContext.Summary != "" {
		chunks = append(chunks, fmt.Sprintf("Persona Summary: %s", data.PersonaContext.Summary))
	}

	// 🔹 Persona Values
	if len(data.PersonaContext.Values) > 0 {
		chunks = append(chunks, fmt.Sprintf("Core Values: %s", strings.Join(data.PersonaContext.Values, ", ")))
	}

	// 🔹 Education
	for _, edu := range data.Education {
		text := fmt.Sprintf("Education: %s at %s — %s", edu.Name, edu.Institution, edu.Description)
		if len(edu.Skills) > 0 {
			text += fmt.Sprintf(" Skills gained: %s.", strings.Join(edu.Skills, ", "))
		}
		if leadership, ok := any(edu.Leadership).([]string); ok && len(leadership) > 0 {
			text += fmt.Sprintf(" Leadership: %s.", strings.Join(leadership, ", "))
		}
		if slices.Contains(edu.Tags, "Award") {
			text += " This entry includes a notable award or recognition."
		}
		chunks = append(chunks, text)
	}

	// 🔹 Work Experience
	for _, exp := range data.WorkExperience {
		text := fmt.Sprintf("Experience: %s at %s — %s", exp.Name, exp.Institution, exp.Description)
		if len(exp.Skills) > 0 {
			text += fmt.Sprintf(" Skills: %s.", strings.Join(exp.Skills, ", "))
		}
		for _, proj := range exp.Projects {
			text += fmt.Sprintf("\nProject: %s — %s", proj.Title, proj.Description)
		}
		chunks = append(chunks, text)
	}

	// 🔹 Projects
	for _, proj := range data.Projects {
		label := "Project"
		if contains(proj.Tags, "Hackathon") {
			label = "Hackathon Project"
		}
		if proj.Featured {
			label = "Featured " + label
		}
		text := fmt.Sprintf("%s: %s — %s", label, proj.Name, proj.Description)
		if len(proj.Skills) > 0 {
			text += fmt.Sprintf(" Technologies: %s.", strings.Join(proj.Skills, ", "))
		}
		if slices.Contains(proj.Tags, "Award") && slices.Contains(proj.Tags, "Hackathon") {
			text += " This project won an award at the hackathon."
		}
		chunks = append(chunks, text)
	}

	// 🔹 Hobbies
	for _, h := range data.Hobbies {
		text := fmt.Sprintf("Hobby: %s — %s", h.Name, h.Description)
		if len(h.Tags) > 0 {
			text += fmt.Sprintf(" Tags: %s.", strings.Join(h.Tags, ", "))
		}
		chunks = append(chunks, text)

		// Optional: add highlights as separate chunks
		for _, highlight := range h.Highlights {
			chunks = append(chunks, fmt.Sprintf("Hobby Highlight: %s", highlight))
		}
	}

	// 🔹 Embed and store all chunks
	for _, chunk := range chunks {
		emb, err := GenerateEmbedding(chunk)
		if err != nil {
			return fmt.Errorf("embedding failed for chunk: %v", err)
		}
		resumeChunks = append(resumeChunks, EmbeddedChunk{
			Text:      chunk,
			Embedding: emb,
		})
	}

	log.Printf("✅ Embedded %d resume chunks", len(resumeChunks))
	return nil
}
// GenerateEmbedding sends text to OpenAI's embedding endpoint
func GenerateEmbedding(input string) ([]float64, error) {
	apiKey := config.GetOpenAIKey()
	reqBody, err := json.Marshal(EmbeddingRequest{
		Input: input,
		Model: "text-embedding-ada-002",
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", config.GetOpenAIEmbeddingURL(), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenAI error: %s", body)
	}

	var result EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return result.Data[0].Embedding, nil
}

// Utility: case-insensitive match
func contains(tags []string, match string) bool {
	for _, t := range tags {
		if strings.EqualFold(t, match) {
			return true
		}
	}
	return false
}
