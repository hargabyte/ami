package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type OllamaClient struct {
	BaseURL string
	Model   string
	HTTP    *http.Client
}

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func NewOllamaClient(baseURL, model string) *OllamaClient {
	return &OllamaClient{
		BaseURL: baseURL,
		Model:   model,
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Generate sends a request to Ollama with a circuit breaker pattern
func (c *OllamaClient) Generate(ctx context.Context, prompt string) (string, error) {
	reqBody := GenerateRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Simple Retry Logic (3 attempts)
	var lastErr error
	for i := 0; i < 3; i++ {
		req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/generate", c.BaseURL), bytes.NewBuffer(jsonData))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.HTTP.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * time.Second) // Exponential-ish backoff
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("ollama returned status: %d", resp.StatusCode)
			continue
		}

		var genResp GenerateResponse
		if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
			return "", err
		}

		return genResp.Response, nil
	}

	return "", fmt.Errorf("ollama request failed after 3 attempts: %w", lastErr)
}
