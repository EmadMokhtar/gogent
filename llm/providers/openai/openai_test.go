package openai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EmadMokhtar/gogent/llm"
)

func TestNew(t *testing.T) {
	provider := New()
	if provider == nil {
		t.Fatal("New() returned nil")
	}
	if provider.baseURL != defaultBaseURL {
		t.Errorf("New() baseURL = %s, want %s", provider.baseURL, defaultBaseURL)
	}
	if provider.model != defaultModel {
		t.Errorf("New() model = %s, want %s", provider.model, defaultModel)
	}
}

func TestWithAPIKey(t *testing.T) {
	apiKey := "test-api-key"
	provider := New(WithAPIKey(apiKey))
	if provider.apiKey != apiKey {
		t.Errorf("WithAPIKey() apiKey = %s, want %s", provider.apiKey, apiKey)
	}
}

func TestWithModel(t *testing.T) {
	model := "gpt-3.5-turbo"
	provider := New(WithModel(model))
	if provider.model != model {
		t.Errorf("WithModel() model = %s, want %s", provider.model, model)
	}
}

func TestWithBaseURL(t *testing.T) {
	baseURL := "https://custom.api.com/v1"
	provider := New(WithBaseURL(baseURL))
	if provider.baseURL != baseURL {
		t.Errorf("WithBaseURL() baseURL = %s, want %s", provider.baseURL, baseURL)
	}
}

func TestWithHTTPClient(t *testing.T) {
	client := &http.Client{}
	provider := New(WithHTTPClient(client))
	if provider.httpClient != client {
		t.Error("WithHTTPClient() did not set HTTP client")
	}
}

func TestProvider_Generate(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "test-id",
			"object": "chat.completion",
			"created": 1234567890,
			"model": "gpt-4",
			"choices": [
				{
					"index": 0,
					"message": {
						"role": "assistant",
						"content": "Hello! I'm here to help."
					},
					"finish_reason": "stop"
				}
			],
			"usage": {
				"prompt_tokens": 10,
				"completion_tokens": 8,
				"total_tokens": 18
			}
		}`))
	}))
	defer server.Close()

	provider := New(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	ctx := context.Background()
	messages := []llm.Message{
		{Role: "user", Content: "Hello"},
	}

	response, err := provider.Generate(ctx, messages)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if response.Content != "Hello! I'm here to help." {
		t.Errorf("Generate() content = %s, want %s", response.Content, "Hello! I'm here to help.")
	}
	if response.Usage.TotalTokens != 18 {
		t.Errorf("Generate() usage.TotalTokens = %d, want 18", response.Usage.TotalTokens)
	}
}

func TestProvider_Generate_WithOptions(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "test-id",
			"object": "chat.completion",
			"created": 1234567890,
			"model": "gpt-4",
			"choices": [
				{
					"index": 0,
					"message": {
						"role": "assistant",
						"content": "Response"
					},
					"finish_reason": "stop"
				}
			],
			"usage": {
				"prompt_tokens": 10,
				"completion_tokens": 5,
				"total_tokens": 15
			}
		}`))
	}))
	defer server.Close()

	provider := New(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	ctx := context.Background()
	messages := []llm.Message{
		{Role: "user", Content: "Hello"},
	}

	temp := 0.7
	maxTokens := 100
	topP := 0.9
	response, err := provider.Generate(ctx, messages,
		llm.WithTemperature(temp),
		llm.WithMaxTokens(maxTokens),
		llm.WithTopP(topP),
		llm.WithStop([]string{"END"}),
	)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if response.Content != "Response" {
		t.Errorf("Generate() content = %s, want %s", response.Content, "Response")
	}
}

func TestProvider_Generate_WithToolCalls(t *testing.T) {
	// Create a mock server that returns a response with tool calls
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "test-id",
			"object": "chat.completion",
			"created": 1234567890,
			"model": "gpt-4",
			"choices": [
				{
					"index": 0,
					"message": {
						"role": "assistant",
						"content": "",
						"tool_calls": [
							{
								"id": "call_123",
								"type": "function",
								"function": {
									"name": "calculator",
									"arguments": "{\"operation\":\"add\",\"a\":5,\"b\":3}"
								}
							}
						]
					},
					"finish_reason": "tool_calls"
				}
			],
			"usage": {
				"prompt_tokens": 20,
				"completion_tokens": 15,
				"total_tokens": 35
			}
		}`))
	}))
	defer server.Close()

	provider := New(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	ctx := context.Background()
	messages := []llm.Message{
		{Role: "user", Content: "What is 5 + 3?"},
	}

	tools := []llm.ToolDefinition{
		{
			Type: "function",
			Function: llm.FunctionDefinition{
				Name:        "calculator",
				Description: "Perform calculations",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"operation": map[string]interface{}{"type": "string"},
						"a":         map[string]interface{}{"type": "number"},
						"b":         map[string]interface{}{"type": "number"},
					},
				},
			},
		},
	}

	response, err := provider.Generate(ctx, messages, llm.WithTools(tools))
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if len(response.ToolCalls) != 1 {
		t.Fatalf("Generate() returned %d tool calls, want 1", len(response.ToolCalls))
	}

	toolCall := response.ToolCalls[0]
	if toolCall.Name != "calculator" {
		t.Errorf("Generate() tool call name = %s, want calculator", toolCall.Name)
	}
	if toolCall.ID != "call_123" {
		t.Errorf("Generate() tool call ID = %s, want call_123", toolCall.ID)
	}
}

func TestProvider_Generate_WithMessagesWithToolCalls(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "test-id",
			"object": "chat.completion",
			"created": 1234567890,
			"model": "gpt-4",
			"choices": [
				{
					"index": 0,
					"message": {
						"role": "assistant",
						"content": "The result is 8"
					},
					"finish_reason": "stop"
				}
			],
			"usage": {
				"prompt_tokens": 20,
				"completion_tokens": 5,
				"total_tokens": 25
			}
		}`))
	}))
	defer server.Close()

	provider := New(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	ctx := context.Background()
	messages := []llm.Message{
		{Role: "user", Content: "What is 5 + 3?"},
		{
			Role:    "assistant",
			Content: "",
			ToolCalls: []llm.ToolCall{
				{
					ID:   "call_123",
					Name: "calculator",
					Input: map[string]interface{}{
						"operation": "add",
						"a":         5,
						"b":         3,
					},
				},
			},
		},
		{Role: "tool", Content: "8"},
	}

	response, err := provider.Generate(ctx, messages)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if response.Content != "The result is 8" {
		t.Errorf("Generate() content = %s, want %s", response.Content, "The result is 8")
	}
}

func TestProvider_Generate_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": {
				"message": "Invalid request",
				"type": "invalid_request_error",
				"code": "invalid_request"
			}
		}`))
	}))
	defer server.Close()

	provider := New(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	ctx := context.Background()
	messages := []llm.Message{
		{Role: "user", Content: "Hello"},
	}

	_, err := provider.Generate(ctx, messages)
	if err == nil {
		t.Error("Generate() should return error for bad request")
	}
}

