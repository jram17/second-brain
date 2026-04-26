package embedder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jram17/second-brain/services/content/pkg/breaker"
)

func getOllamaURL() string {
	if u := os.Getenv("OLLAMA_URL"); u != "" {
		return u + "/api/embeddings"
	}
	return "http://localhost:11434/api/embeddings"
}

type embeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}
type embeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

var cb = breaker.New("ollama-embedder")

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

	req, err := http.NewRequestWithContext(ctx, "POST", getOllamaURL(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	result, err := cb.Execute(func() (interface{}, error) {
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
		return res.Embedding,nil
	})

	if err!=nil{
		return nil,fmt.Errorf("embedding failed (circuit breaker): %w",err)
	}
	//ipo convert the float 64 to 32 cause qdrant expects 32
	raw := result.([]float64)
	embeddings := make([]float32,len(raw))
	for i,v := range raw{
		embeddings[i] = float32(v)
	}
	return embeddings,nil
}
