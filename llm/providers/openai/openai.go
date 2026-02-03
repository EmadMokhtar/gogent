// Package openai provides an OpenAI LLM provider implementation.
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
	"github.com/EmadMokhtar/gogent/types"
)

// Provider implements the LLM provider interface for OpenAI.
type Provider struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

// Option is a functional option for configuring the OpenAI provider.
type Option func(*Provider)

// New creates a new OpenAI provider.
func New(opts ...Option) *Provider {
	p := &Provider{
		model:   "gpt-4",
		baseURL: "https://api.openai.com/v1",
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
	
	for _, opt := range opts {
		opt(p)
	}
	
	return p
}

// WithAPIKey sets the OpenAI API key.
func WithAPIKey(apiKey string) Option {
	return func(p *Provider) {
		p.apiKey = apiKey
	}
}

// WithModel sets the model to use.
func WithModel(model string) Option {
	return func(p *Provider) {
		p.model = model
	}
}

// WithBaseURL sets the base URL for API requests.
func WithBaseURL(baseURL string) Option {
	return func(p *Provider) {
		p.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(p *Provider) {
		p.httpClient = client
	}
}

// Generate produces a response from OpenAI.
func (p *Provider) Generate(ctx context.Context, messages []types.Message, opts ...llm.Option) (*types.Response, error) {
	// Apply options
	options := llm.DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	
	// Use model override if provided
	model := p.model
	if options.Model != "" {
		model = options.Model
	}
	
	// Build request
	reqBody := buildChatCompletionRequest(messages, model, options)
	
	// Make API call
	respBody, err := p.makeRequest(ctx, "/chat/completions", reqBody)
	if err != nil {
		return nil, err
	}
	
	// Parse response
	var chatResp chatCompletionResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}
	
	// Convert to agent response
	choice := chatResp.Choices[0]
	response := &types.Response{
		Content:  choice.Message.Content,
		Messages: []types.Message{convertToAgentMessage(choice.Message)},
		Metadata: map[string]interface{}{
			"model":       chatResp.Model,
			"usage":       chatResp.Usage,
			"finish_reason": choice.FinishReason,
		},
	}
	
	// Add tool calls if present
	if len(choice.Message.ToolCalls) > 0 {
		response.ToolCalls = convertToAgentToolCalls(choice.Message.ToolCalls)
	}
	
	return response, nil
}

// Stream produces a streaming response from OpenAI.
func (p *Provider) Stream(ctx context.Context, messages []types.Message, opts ...llm.Option) (<-chan llm.Token, error) {
	// Apply options
	options := llm.DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	
	// Use model override if provided
	model := p.model
	if options.Model != "" {
		model = options.Model
	}
	
	// Build request with streaming enabled
	reqBody := buildChatCompletionRequest(messages, model, options)
	reqBody["stream"] = true
	
	// Create channel for tokens
	tokenChan := make(chan llm.Token)
	
	// Start streaming in a goroutine
	go func() {
		defer close(tokenChan)
		
		// Make streaming request
		if err := p.makeStreamingRequest(ctx, "/chat/completions", reqBody, tokenChan); err != nil {
			tokenChan <- llm.Token{Error: err, Done: true}
		}
	}()
	
	return tokenChan, nil
}

// makeRequest makes an HTTP request to the OpenAI API.
func (p *Provider) makeRequest(ctx context.Context, path string, body map[string]interface{}) ([]byte, error) {
	// Encode request body
	reqData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+path, bytes.NewReader(reqData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	
	// Make request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	// Check for errors
	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		if err := json.Unmarshal(respData, &errResp); err == nil {
			return nil, fmt.Errorf("API error: %s", errResp.Error.Message)
		}
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(respData))
	}
	
	return respData, nil
}

// makeStreamingRequest makes a streaming HTTP request to the OpenAI API.
func (p *Provider) makeStreamingRequest(ctx context.Context, path string, body map[string]interface{}, tokenChan chan<- llm.Token) error {
	// Encode request body
	reqData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+path, bytes.NewReader(reqData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Accept", "text/event-stream")
	
	// Make request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Check status
	if resp.StatusCode != http.StatusOK {
		respData, _ := io.ReadAll(resp.Body)
		var errResp errorResponse
		if err := json.Unmarshal(respData, &errResp); err == nil {
			return fmt.Errorf("API error: %s", errResp.Error.Message)
		}
		return fmt.Errorf("API error: status %d", resp.StatusCode)
	}
	
	// Read SSE stream
	return p.readSSEStream(resp.Body, tokenChan)
}

// readSSEStream reads an SSE stream and sends tokens to the channel.
func (p *Provider) readSSEStream(body io.Reader, tokenChan chan<- llm.Token) error {
	// Use bufio.Reader for line-by-line reading
	reader := io.Reader(body)
	buf := make([]byte, 1)
	var line []byte
	var content string
	
	for {
		// Read byte by byte to find line boundaries
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		
		if n == 0 {
			break
		}
		
		// Check for newline
		if buf[0] == '\n' {
			lineStr := string(line)
			line = []byte{} // Reset line buffer
			
			// Parse SSE data
			if len(lineStr) > 6 && lineStr[:6] == "data: " {
				data := lineStr[6:]
				
				// Trim whitespace including trailing newlines
				data = string(bytes.TrimSpace([]byte(data)))
				
				// Check for stream end
				if data == "[DONE]" {
					tokenChan <- llm.Token{Content: content, Done: true}
					return nil
				}
				
				// Parse chunk
				var chunk chatCompletionChunk
				if err := json.Unmarshal([]byte(data), &chunk); err != nil {
					continue // Skip malformed chunks
				}
				
				if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
					delta := chunk.Choices[0].Delta.Content
					content += delta
					tokenChan <- llm.Token{Content: content, Delta: delta, Done: false}
				}
			}
		} else {
			line = append(line, buf[0])
		}
		
		if err == io.EOF {
			break
		}
	}
	
	tokenChan <- llm.Token{Content: content, Done: true}
	return nil
}
