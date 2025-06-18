package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-ai/config"
	"io"
	"net/http"
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

// generateEmbedding sends resume text to OpenAI and returns the vector
func generateEmbedding(input string) ([]float64, error) {
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
