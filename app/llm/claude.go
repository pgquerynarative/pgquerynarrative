package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const claudeBaseURL = "https://api.anthropic.com/v1"
const claudeAPIVersion = "2023-06-01"

// ClaudeClient calls the Anthropic Messages API for text generation.
type ClaudeClient struct {
	apiKey string
	model  string
	client *http.Client
}

// NewClaudeClient returns a client for the Claude API.
// apiKey is the Anthropic API key (from LLM_API_KEY). model is the model name (e.g. claude-3-5-sonnet-20241022, claude-3-haiku-20240307).
func NewClaudeClient(apiKey, model string) *ClaudeClient {
	if model == "" {
		model = "claude-3-5-sonnet-20241022"
	}
	return &ClaudeClient{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Name returns the provider name.
func (c *ClaudeClient) Name() string {
	return "claude"
}

func (c *ClaudeClient) Model() string {
	return c.model
}

const claudeMaxRetries = 3
const claudeRetryDelay = 6 * time.Second

// Generate sends the prompt to Claude Messages API and returns the generated text.
// On 429 (rate limit) it retries up to claudeMaxRetries with backoff.
func (c *ClaudeClient) Generate(ctx context.Context, prompt string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("claude: LLM_API_KEY is required")
	}

	url := claudeBaseURL + "/messages"

	payload := map[string]interface{}{
		"model":      c.model,
		"max_tokens": 2048,
		"messages": []map[string]interface{}{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("claude: marshal request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt < claudeMaxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return "", fmt.Errorf("claude: create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", c.apiKey)
		req.Header.Set("anthropic-version", claudeAPIVersion)

		resp, err := c.client.Do(req)
		if err != nil {
			return "", fmt.Errorf("claude: request: %w", err)
		}

		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var result struct {
				Content []struct {
					Type string `json:"type"`
					Text string `json:"text"`
				} `json:"content"`
			}
			if err := json.Unmarshal(body, &result); err != nil {
				return "", fmt.Errorf("claude: decode response: %w", err)
			}
			var text string
			for _, block := range result.Content {
				if block.Type == "text" {
					text += block.Text
				}
			}
			if text == "" {
				return "", fmt.Errorf("claude: empty response")
			}
			return text, nil
		}

		lastErr = fmt.Errorf("claude API error: %d - %s", resp.StatusCode, string(body))

		if resp.StatusCode == 429 && attempt < claudeMaxRetries-1 {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(claudeRetryDelay):
				// retry
			}
			continue
		}

		return "", lastErr
	}
	return "", lastErr
}
