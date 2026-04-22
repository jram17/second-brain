package embedder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const ollamaURL = "http://localhost:11434/api/embeddings"

type embeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}
type embeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

func Embed(text string) ([]float32, error) {
	reqBody := embeddingRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", ollamaURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var res embeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	embedding := make([]float32, len(res.Embedding))
	for i, v := range res.Embedding {
		embedding[i] = float32(v)
	}
	return embedding, nil
}
