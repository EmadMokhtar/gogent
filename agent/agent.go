package agent

import (
	"context"

	"github.com/EmadMokhtar/gogent/llm"
	"github.com/EmadMokhtar/gogent/memory"
)

// Agent represents an AI agent with memory capabilities.
type Agent struct {
	llm          llm.LLM
	memory       memory.Memory
	systemPrompt string
}

// Response represents the agent's response.
type Response struct {
	Content  string
	Role     string
	Metadata map[string]interface{}
}

// New creates a new agent with the given options.
func New(opts ...Option) *Agent {
	a := &Agent{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Run executes the agent with the given input.
func (a *Agent) Run(ctx context.Context, input string) (*Response, error) {
	// Build message history
	var messages []memory.Message

	// Add system prompt if configured
	if a.systemPrompt != "" {
		messages = append(messages, memory.Message{
			Role:    "system",
			Content: a.systemPrompt,
		})
	}

	// Retrieve context from memory if available
	if a.memory != nil {
		recentMessages, err := a.memory.Get(ctx, 10)
		if err == nil {
			messages = append(messages, recentMessages...)
		}
	}

	// Add the current user input
	userMessage := memory.Message{
		Role:    "user",
		Content: input,
	}
	messages = append(messages, userMessage)

	// Store user message in memory
	if a.memory != nil {
		a.memory.Add(ctx, userMessage)
	}

	// Call LLM if available
	var response string
	if a.llm != nil {
		resp, err := a.llm.Generate(ctx, messages)
		if err != nil {
			return nil, err
		}
		response = resp
	} else {
		// If no LLM, just echo back
		response = "Echo: " + input
	}

	// Create assistant response
	assistantMessage := memory.Message{
		Role:    "assistant",
		Content: response,
	}

	// Store assistant response in memory
	if a.memory != nil {
		a.memory.Add(ctx, assistantMessage)
	}

	return &Response{
		Content: response,
		Role:    "assistant",
	}, nil
}
