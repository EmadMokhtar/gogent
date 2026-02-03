package agent

import (
	"context"
	"testing"

	"github.com/EmadMokhtar/gogent/llm"
	"github.com/EmadMokhtar/gogent/tools/builtin"
)

// mockProvider is a mock LLM provider for testing
type mockProvider struct {
	responses []*llm.Response
	callCount int
}

func (m *mockProvider) Generate(ctx context.Context, messages []llm.Message, opts ...llm.Option) (*llm.Response, error) {
	if m.callCount < len(m.responses) {
		resp := m.responses[m.callCount]
		m.callCount++
		return resp, nil
	}
	return &llm.Response{Content: "Default response"}, nil
}

func (m *mockProvider) Stream(ctx context.Context, messages []llm.Message, opts ...llm.Option) (<-chan llm.Token, error) {
	ch := make(chan llm.Token)
	close(ch)
	return ch, nil
}

func TestNew(t *testing.T) {
	agent := New()
	if agent == nil {
		t.Fatal("New() returned nil")
	}
	if agent.maxIterations != 10 {
		t.Errorf("New() maxIterations = %d, want 10", agent.maxIterations)
	}
}

func TestWithLLM(t *testing.T) {
	provider := &mockProvider{}
	agent := New(WithLLM(provider))
	
	if agent.llm == nil {
		t.Error("WithLLM() did not set LLM provider")
	}
}

func TestWithSystemPrompt(t *testing.T) {
	prompt := "You are a helpful assistant"
	agent := New(WithSystemPrompt(prompt))
	
	if agent.systemPrompt != prompt {
		t.Errorf("WithSystemPrompt() systemPrompt = %s, want %s", agent.systemPrompt, prompt)
	}
}

func TestWithMaxIterations(t *testing.T) {
	agent := New(WithMaxIterations(5))
	
	if agent.maxIterations != 5 {
		t.Errorf("WithMaxIterations() maxIterations = %d, want 5", agent.maxIterations)
	}
}

func TestWithTemperature(t *testing.T) {
	temp := 0.7
	agent := New(WithTemperature(temp))
	
	if agent.temperature == nil {
		t.Error("WithTemperature() did not set temperature")
		return
	}
	if *agent.temperature != temp {
		t.Errorf("WithTemperature() = %v, want %v", *agent.temperature, temp)
	}
}

func TestWithMaxTokens(t *testing.T) {
	maxTokens := 1000
	agent := New(WithMaxTokens(maxTokens))
	
	if agent.maxTokens == nil {
		t.Error("WithMaxTokens() did not set max tokens")
		return
	}
	if *agent.maxTokens != maxTokens {
		t.Errorf("WithMaxTokens() = %v, want %v", *agent.maxTokens, maxTokens)
	}
}

func TestWithTools(t *testing.T) {
	agent := New(WithTools(builtin.NewCalculator()))
	
	tools := agent.tools.List()
	if len(tools) != 1 {
		t.Errorf("WithTools() registered %d tools, want 1", len(tools))
	}
}

func TestAgent_Run_NoLLM(t *testing.T) {
	agent := New()
	ctx := context.Background()
	
	_, err := agent.Run(ctx, "test input")
	if err == nil {
		t.Error("Agent.Run() should return error when no LLM is configured")
	}
}

func TestAgent_Run_Simple(t *testing.T) {
	provider := &mockProvider{
		responses: []*llm.Response{
			{
				Content: "Hello! How can I help you?",
				Usage: llm.Usage{
					PromptTokens:     10,
					CompletionTokens: 8,
					TotalTokens:      18,
				},
			},
		},
	}
	
	agent := New(
		WithLLM(provider),
		WithSystemPrompt("You are a helpful assistant"),
	)
	
	ctx := context.Background()
	response, err := agent.Run(ctx, "Hello")
	
	if err != nil {
		t.Fatalf("Agent.Run() error = %v", err)
	}
	
	if response.Content != "Hello! How can I help you?" {
		t.Errorf("Agent.Run() content = %s, want %s", response.Content, "Hello! How can I help you?")
	}
	
	if response.Usage.TotalTokens != 18 {
		t.Errorf("Agent.Run() usage.TotalTokens = %d, want 18", response.Usage.TotalTokens)
	}
}

func TestAgent_Run_WithToolCalls(t *testing.T) {
	// Simulate a tool call workflow
	provider := &mockProvider{
		responses: []*llm.Response{
			{
				Content: "",
				ToolCalls: []llm.ToolCall{
					{
						ID:   "call_1",
						Name: "calculator",
						Input: map[string]interface{}{
							"operation": "add",
							"a":         float64(5),
							"b":         float64(3),
						},
					},
				},
			},
			{
				Content: "The result is 8",
			},
		},
	}
	
	agent := New(
		WithLLM(provider),
		WithTools(builtin.NewCalculator()),
	)
	
	ctx := context.Background()
	response, err := agent.Run(ctx, "What is 5 + 3?")
	
	if err != nil {
		t.Fatalf("Agent.Run() error = %v", err)
	}
	
	if response.Content != "The result is 8" {
		t.Errorf("Agent.Run() content = %s, want %s", response.Content, "The result is 8")
	}
}

func TestAgent_Messages(t *testing.T) {
	provider := &mockProvider{
		responses: []*llm.Response{
			{Content: "Response 1"},
		},
	}
	
	agent := New(WithLLM(provider))
	ctx := context.Background()
	
	agent.Run(ctx, "Test input")
	
	messages := agent.Messages()
	if len(messages) != 2 { // user message + assistant response
		t.Errorf("Agent.Messages() returned %d messages, want 2", len(messages))
	}
}

func TestAgent_Reset(t *testing.T) {
	provider := &mockProvider{
		responses: []*llm.Response{
			{Content: "Response 1"},
		},
	}
	
	agent := New(WithLLM(provider))
	ctx := context.Background()
	
	agent.Run(ctx, "Test input")
	agent.Reset()
	
	messages := agent.Messages()
	if len(messages) != 0 {
		t.Errorf("Agent.Reset() should clear messages, but got %d messages", len(messages))
	}
}

