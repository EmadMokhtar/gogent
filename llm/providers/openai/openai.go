package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/EmadMokhtar/gogent/llm"
)

const (
	defaultBaseURL = "https://api.openai.com/v1"
	defaultModel   = "gpt-4"
)

// Provider is an OpenAI LLM provider
type Provider struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

// Option is a functional option for configuring the OpenAI provider
type Option func(*Provider)

// New creates a new OpenAI provider
func New(opts ...Option) *Provider {
	p := &Provider{
		baseURL:    defaultBaseURL,
		model:      defaultModel,
		httpClient: &http.Client{},
	}
	
	for _, opt := range opts {
		opt(p)
	}
	
	return p
}

// WithAPIKey sets the API key for the OpenAI provider
func WithAPIKey(apiKey string) Option {
	return func(p *Provider) {
		p.apiKey = apiKey
	}
}

// WithModel sets the model for the OpenAI provider
func WithModel(model string) Option {
	return func(p *Provider) {
		p.model = model
	}
}

// WithBaseURL sets the base URL for the OpenAI provider
func WithBaseURL(baseURL string) Option {
	return func(p *Provider) {
		p.baseURL = baseURL
	}
}

// WithHTTPClient sets the HTTP client for the OpenAI provider
func WithHTTPClient(client *http.Client) Option {
	return func(p *Provider) {
		p.httpClient = client
	}
}

// Generate sends messages to OpenAI and returns a response
func (p *Provider) Generate(ctx context.Context, messages []llm.Message, opts ...llm.Option) (*llm.Response, error) {
	options := llm.ApplyOptions(opts...)
	
	// Build request
	reqBody := p.buildChatCompletionRequest(messages, options)
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, p.parseError(resp.StatusCode, body)
	}
	
	var chatResp chatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return p.parseResponse(&chatResp), nil
}

// Stream sends messages to OpenAI and returns a channel of tokens
func (p *Provider) Stream(ctx context.Context, messages []llm.Message, opts ...llm.Option) (<-chan llm.Token, error) {
	options := llm.ApplyOptions(opts...)
	
	// Build request with stream=true
	reqBody := p.buildChatCompletionRequest(messages, options)
	reqBody["stream"] = true
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, p.parseError(resp.StatusCode, body)
	}
	
	tokenChan := make(chan llm.Token)
	
	go func() {
		defer close(tokenChan)
		defer resp.Body.Close()
		
		decoder := json.NewDecoder(resp.Body)
		for {
			var line map[string]interface{}
			if err := decoder.Decode(&line); err != nil {
				if err != io.EOF {
					tokenChan <- llm.Token{Error: err}
				}
				return
			}
			
			// Parse streaming response (simplified)
			if choices, ok := line["choices"].([]interface{}); ok && len(choices) > 0 {
				if choice, ok := choices[0].(map[string]interface{}); ok {
					if delta, ok := choice["delta"].(map[string]interface{}); ok {
						if content, ok := delta["content"].(string); ok {
							tokenChan <- llm.Token{Content: content}
						}
					}
				}
			}
		}
	}()
	
	return tokenChan, nil
}

// buildChatCompletionRequest builds a chat completion request
func (p *Provider) buildChatCompletionRequest(messages []llm.Message, options *llm.Options) map[string]interface{} {
	req := map[string]interface{}{
		"model":    p.model,
		"messages": p.convertMessages(messages),
	}
	
	if options.Temperature != nil {
		req["temperature"] = *options.Temperature
	}
	if options.MaxTokens != nil {
		req["max_tokens"] = *options.MaxTokens
	}
	if options.TopP != nil {
		req["top_p"] = *options.TopP
	}
	if len(options.Stop) > 0 {
		req["stop"] = options.Stop
	}
	if len(options.Tools) > 0 {
		req["tools"] = p.convertTools(options.Tools)
	}
	
	return req
}

// convertMessages converts llm.Message to OpenAI format
func (p *Provider) convertMessages(messages []llm.Message) []map[string]interface{} {
	result := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		m := map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		}
		if len(msg.ToolCalls) > 0 {
			toolCalls := make([]map[string]interface{}, len(msg.ToolCalls))
			for j, tc := range msg.ToolCalls {
				toolCalls[j] = map[string]interface{}{
					"id":   tc.ID,
					"type": "function",
					"function": map[string]interface{}{
						"name":      tc.Name,
						"arguments": toJSONString(tc.Input),
					},
				}
			}
			m["tool_calls"] = toolCalls
		}
		result[i] = m
	}
	return result
}

// convertTools converts tool definitions to OpenAI format
func (p *Provider) convertTools(tools []llm.ToolDefinition) []map[string]interface{} {
	result := make([]map[string]interface{}, len(tools))
	for i, tool := range tools {
		result[i] = map[string]interface{}{
			"type": tool.Type,
			"function": map[string]interface{}{
				"name":        tool.Function.Name,
				"description": tool.Function.Description,
				"parameters":  tool.Function.Parameters,
			},
		}
	}
	return result
}

// parseResponse converts OpenAI response to llm.Response
func (p *Provider) parseResponse(resp *chatCompletionResponse) *llm.Response {
	result := &llm.Response{
		Usage: llm.Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
		Metadata: make(map[string]interface{}),
	}
	
	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]
		result.Content = choice.Message.Content
		
		if len(choice.Message.ToolCalls) > 0 {
			result.ToolCalls = make([]llm.ToolCall, len(choice.Message.ToolCalls))
			for i, tc := range choice.Message.ToolCalls {
				var input map[string]interface{}
				if tc.Function.Arguments != "" {
					json.Unmarshal([]byte(tc.Function.Arguments), &input)
				}
				result.ToolCalls[i] = llm.ToolCall{
					ID:    tc.ID,
					Name:  tc.Function.Name,
					Input: input,
				}
			}
		}
	}
	
	return result
}

// parseError parses an error response from OpenAI
func (p *Provider) parseError(statusCode int, body []byte) error {
	var errResp errorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return fmt.Errorf("OpenAI API error (status %d): %s", statusCode, string(body))
	}
	return fmt.Errorf("OpenAI API error: %s", errResp.Error.Message)
}

// toJSONString converts a map to JSON string
func toJSONString(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

// fromJSONString converts JSON string to map
func fromJSONString(s string) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal([]byte(s), &result)
	return result
}
