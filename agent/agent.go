package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/EmadMokhtar/gogent/llm"
	"github.com/EmadMokhtar/gogent/tools"
)

// Agent represents an AI agent that can use tools and interact with an LLM
type Agent struct {
	llm           llm.Provider
	tools         *tools.Registry
	systemPrompt  string
	messages      []llm.Message
	maxIterations int
	temperature   *float64
	maxTokens     *int
	mu            sync.RWMutex
}

// Response represents a response from the agent
type Response struct {
	Content      string
	ToolCalls    []llm.ToolCall
	Usage        llm.Usage
	MessageCount int
}

// Option is a functional option for configuring an Agent
type Option func(*Agent)

// New creates a new Agent with the given options
func New(opts ...Option) *Agent {
	a := &Agent{
		maxIterations: 10,
		messages:      []llm.Message{},
		tools:         tools.NewRegistry(),
	}
	
	for _, opt := range opts {
		opt(a)
	}
	
	return a
}

// WithLLM sets the LLM provider for the agent
func WithLLM(provider llm.Provider) Option {
	return func(a *Agent) {
		a.llm = provider
	}
}

// WithTools registers tools with the agent
func WithTools(toolList ...tools.Tool) Option {
	return func(a *Agent) {
		for _, tool := range toolList {
			a.tools.Register(tool)
		}
	}
}

// WithSystemPrompt sets the system prompt for the agent
func WithSystemPrompt(prompt string) Option {
	return func(a *Agent) {
		a.systemPrompt = prompt
	}
}

// WithMaxIterations sets the maximum number of iterations for tool execution
func WithMaxIterations(max int) Option {
	return func(a *Agent) {
		a.maxIterations = max
	}
}

// WithTemperature sets the temperature for LLM calls
func WithTemperature(temperature float64) Option {
	return func(a *Agent) {
		a.temperature = &temperature
	}
}

// WithMaxTokens sets the maximum tokens for LLM calls
func WithMaxTokens(maxTokens int) Option {
	return func(a *Agent) {
		a.maxTokens = &maxTokens
	}
}

// Run executes the agent with the given input
func (a *Agent) Run(ctx context.Context, input string) (*Response, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	if a.llm == nil {
		return nil, fmt.Errorf("no LLM provider configured")
	}
	
	// Add system prompt if this is the first message and a system prompt is set
	if len(a.messages) == 0 && a.systemPrompt != "" {
		a.messages = append(a.messages, llm.Message{
			Role:    "system",
			Content: a.systemPrompt,
		})
	}
	
	// Add user message
	a.messages = append(a.messages, llm.Message{
		Role:    "user",
		Content: input,
	})
	
	// Execute agent loop with tool calls
	iteration := 0
	for iteration < a.maxIterations {
		// Build LLM options
		opts := a.buildLLMOptions()
		
		// Call LLM
		response, err := a.llm.Generate(ctx, a.messages, opts...)
		if err != nil {
			return nil, fmt.Errorf("LLM generation failed: %w", err)
		}
		
		// Add assistant message to history
		assistantMsg := llm.Message{
			Role:      "assistant",
			Content:   response.Content,
			ToolCalls: response.ToolCalls,
		}
		a.messages = append(a.messages, assistantMsg)
		
		// If no tool calls, we're done
		if len(response.ToolCalls) == 0 {
			return &Response{
				Content:      response.Content,
				ToolCalls:    response.ToolCalls,
				Usage:        response.Usage,
				MessageCount: len(a.messages),
			}, nil
		}
		
		// Execute tool calls
		for _, toolCall := range response.ToolCalls {
			result, err := a.tools.Execute(ctx, toolCall)
			var resultStr string
			if err != nil {
				resultStr = fmt.Sprintf("Error: %v", err)
			} else {
				resultStr = fmt.Sprintf("%v", result)
			}
			
			// Add tool result as a message
			a.messages = append(a.messages, llm.Message{
				Role:    "tool",
				Content: resultStr,
				Metadata: map[string]interface{}{
					"tool_call_id": toolCall.ID,
					"tool_name":    toolCall.Name,
				},
			})
		}
		
		iteration++
	}
	
	return nil, fmt.Errorf("max iterations (%d) reached without final response", a.maxIterations)
}

// buildLLMOptions builds LLM options from agent configuration
func (a *Agent) buildLLMOptions() []llm.Option {
	opts := []llm.Option{}
	
	if a.temperature != nil {
		opts = append(opts, llm.WithTemperature(*a.temperature))
	}
	if a.maxTokens != nil {
		opts = append(opts, llm.WithMaxTokens(*a.maxTokens))
	}
	
	// Add tools if any are registered
	toolDefs := a.tools.ToLLMTools()
	if len(toolDefs) > 0 {
		opts = append(opts, llm.WithTools(toolDefs))
	}
	
	return opts
}

// Messages returns a copy of the current message history
func (a *Agent) Messages() []llm.Message {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	messages := make([]llm.Message, len(a.messages))
	copy(messages, a.messages)
	return messages
}

// Reset clears the agent's message history
func (a *Agent) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	a.messages = []llm.Message{}
}
