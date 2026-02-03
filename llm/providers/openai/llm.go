package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/EmadMokhtar/gogent/llm"
)

// LLM is an OpenAI language model provider
type LLM struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

// LLMOption defines configuration options for LLM
type LLMOption func(*LLM)

// WithLLMAPIKey sets the OpenAI API key
func WithLLMAPIKey(apiKey string) LLMOption {
	return func(l *LLM) {
		l.apiKey = apiKey
	}
}

// WithLLMModel sets the model name
func WithLLMModel(model string) LLMOption {
	return func(l *LLM) {
		l.model = model
	}
}

// New creates a new OpenAI LLM provider
func New(opts ...LLMOption) *LLM {
	llmProvider := &LLM{
		model:      "gpt-3.5-turbo",
		baseURL:    "https://api.openai.com/v1",
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}

	for _, opt := range opts {
		opt(llmProvider)
	}

	return llmProvider
}

// Generate generates a response for the given messages
func (l *LLM) Generate(ctx context.Context, messages []llm.Message) (*llm.Response, error) {
	if l.apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	// Convert messages to OpenAI format
	var apiMessages []map[string]string
	for _, msg := range messages {
		apiMessages = append(apiMessages, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	reqBody := map[string]interface{}{
		"model":    l.model,
		"messages": apiMessages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := l.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+l.apiKey)

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var apiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(apiResponse.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	response := &llm.Response{
		Content: apiResponse.Choices[0].Message.Content,
		Model:   l.model,
		Usage: &llm.Usage{
			PromptTokens:     apiResponse.Usage.PromptTokens,
			CompletionTokens: apiResponse.Usage.CompletionTokens,
			TotalTokens:      apiResponse.Usage.TotalTokens,
		},
	}

	return response, nil
}

// Model returns the model name
func (l *LLM) Model() string {
	return l.model
}
