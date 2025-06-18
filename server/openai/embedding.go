package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-ai/config"
	"go-ai/resume"
	"io"
	"log"
	"math"
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
func contains(tags []string, match string) bool {
	for _, t := range tags {
		if strings.EqualFold(t, match) {
			return true
		}
	}
	return false
}

var resumeChunks []EmbeddedChunk

// LoadAndEmbedResumeChunks takes ResumeData and turns it into embedded chunks for RAG.
func LoadAndEmbedResumeChunks(data *resume.ResumeData) error {
	var chunks []string

	// Persona summary
	if data.PersonaContext.Summary != "" {
		chunks = append(chunks, fmt.Sprintf("Persona Summary: %s", data.PersonaContext.Summary))
	}

	// Education
	for _, edu := range data.Education {
		text := fmt.Sprintf("Education: %s at %s — %s", edu.Name, edu.Institution, edu.Description)
		if slices.Contains(edu.Tags, "Award") {
			text += " This entry includes a notable award or recognition."
		}
		chunks = append(chunks, text)
	}

	// Work Experience
	for _, exp := range data.WorkExperience {
		text := fmt.Sprintf("Experience: %s at %s — %s", exp.Name, exp.Institution, exp.Description)
		for _, proj := range exp.Projects {
			text += fmt.Sprintf("\nProject: %s — %s", proj.Title, proj.Description)
		}
		chunks = append(chunks, text)
	}

	for _, proj := range data.Projects {
		label := "Project"
		if contains(proj.Tags, "Hackathon") {
			label = "Hackathon Project"
		}
		if proj.Featured {
			label = "Featured " + label
		}
		text := fmt.Sprintf("%s: %s — %s", label, proj.Name, proj.Description)
		if slices.Contains(proj.Tags, "Award") && slices.Contains(proj.Tags, "Hackathon") {
		text += " This project won an award at the hackathon."
	}
		chunks = append(chunks, text)
	}



	// Optional: Hobbies (useful if persona-related queries arise)
	for _, h := range data.Hobbies {
		text := fmt.Sprintf("Hobby: %s — %s", h.Name, h.Description)
		chunks = append(chunks, text)
	}

	// Embed each chunk
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

// generateEmbedding sends resume text to OpenAI and returns the vector
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
		return nil, fmt.Errorf("No embedding returned")
	}

	return result.Data[0].Embedding, nil
}


// Similarity Logic

// dot product
func dot(a, b []float64) float64 {
	var sum float64
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

func norm(a []float64) float64 {
	var sum float64
	for _, v := range a {
		sum += v * v
	}
	return math.Sqrt(sum)
}

func CosineSim(a, b []float64) float64 {
	return dot(a, b) / (norm(a) * norm(b))
}


// SearchRelevantChunks returns top 3 most similar chunks for the user query
func SearchRelevantChunks(query string) ([]string, error){
	log.Println("🔍 Generating embedding for user query:", query)
	queryVec, err := GenerateEmbedding(query)
	if err != nil {
		return  nil, err
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
		return nil, nil
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	top := 3
	if len(results) < 3 {
		top = len(results)
	}
	var rawTopChunks []string
	for i := 0; i < top; i++ {
		rawTopChunks = append(rawTopChunks, results[i].Text)
	}

	log.Println("📥 Retrieved Chunks", rawTopChunks)
	return rawTopChunks, nil
}