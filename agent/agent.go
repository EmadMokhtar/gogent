// Package agent provides the core agent implementation for the Gogent framework.
package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/EmadMokhtar/gogent/llm"
	"github.com/EmadMokhtar/gogent/tools"
	"github.com/EmadMokhtar/gogent/types"
)

// Agent represents an AI agent that can execute tasks using an LLM and tools.
type Agent struct {
	llmProvider   llm.Provider
	toolRegistry  *tools.Registry
	toolExecutor  *tools.Executor
	systemPrompt  string
	maxIterations int
	timeout       time.Duration
}

// Option is a functional option for configuring an agent.
type Option func(*Agent)

// New creates a new agent with the given options.
func New(opts ...Option) *Agent {
	registry := tools.NewRegistry()
	
	a := &Agent{
		toolRegistry:  registry,
		toolExecutor:  tools.NewExecutor(registry),
		systemPrompt:  "You are a helpful AI assistant.",
		maxIterations: 5,
		timeout:       60 * time.Second,
	}
	
	for _, opt := range opts {
		opt(a)
	}
	
	return a
}

// WithLLM sets the LLM provider for the agent.
func WithLLM(provider llm.Provider) Option {
	return func(a *Agent) {
		a.llmProvider = provider
	}
}

// WithTools registers tools with the agent.
func WithTools(toolList ...tools.Tool) Option {
	return func(a *Agent) {
		for _, tool := range toolList {
			a.toolRegistry.Register(tool)
		}
	}
}

// WithSystemPrompt sets the system prompt for the agent.
func WithSystemPrompt(prompt string) Option {
	return func(a *Agent) {
		a.systemPrompt = prompt
	}
}

// WithMaxIterations sets the maximum number of iterations for tool calling.
func WithMaxIterations(max int) Option {
	return func(a *Agent) {
		a.maxIterations = max
	}
}

// WithTimeout sets the timeout for agent execution.
func WithTimeout(timeout time.Duration) Option {
	return func(a *Agent) {
		a.timeout = timeout
	}
}

// Run executes the agent with the given user query.
func (a *Agent) Run(ctx context.Context, query string) (*types.Response, error) {
	if a.llmProvider == nil {
		return nil, fmt.Errorf("LLM provider not configured")
	}
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()
	
	// Create run context
	runCtx := NewRunContext(ctx, a.systemPrompt, a.maxIterations)
	
	// Add user message
	runCtx.AddMessage(types.Message{
		Role:    "user",
		Content: query,
	})
	
	// Execute agent loop
	return a.executeLoop(runCtx)
}

// RunWithMessages executes the agent with a custom message history.
func (a *Agent) RunWithMessages(ctx context.Context, messages []types.Message) (*types.Response, error) {
	if a.llmProvider == nil {
		return nil, fmt.Errorf("LLM provider not configured")
	}
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()
	
	// Create run context
	runCtx := NewRunContext(ctx, a.systemPrompt, a.maxIterations)
	
	// Add messages
	for _, msg := range messages {
		runCtx.AddMessage(msg)
	}
	
	// Execute agent loop
	return a.executeLoop(runCtx)
}

// executeLoop executes the main agent loop with tool calling.
func (a *Agent) executeLoop(runCtx *RunContext) (*types.Response, error) {
	for {
		// Check if we can continue iterating
		if !runCtx.CanIterate() {
			return nil, fmt.Errorf("max iterations reached")
		}
		
		// Build LLM options with tools
		llmOpts := a.buildLLMOptions()
		
		// Call LLM
		response, err := a.llmProvider.Generate(runCtx.Context(), runCtx.Messages(), llmOpts...)
		if err != nil {
			return nil, fmt.Errorf("LLM generation failed: %w", err)
		}
		
		// Add assistant message to context
		if len(response.Messages) > 0 {
			runCtx.AddMessage(response.Messages[0])
		}
		
		// Check if there are tool calls
		if len(response.ToolCalls) == 0 {
			// No tool calls, return final response
			response.Messages = runCtx.messages
			return response, nil
		}
		
		// Execute tool calls
		toolMessages, err := a.toolExecutor.ExecuteToolCalls(runCtx.Context(), response.ToolCalls)
		if err != nil {
			return nil, fmt.Errorf("tool execution failed: %w", err)
		}
		
		// Add tool results to context
		for _, msg := range toolMessages {
			runCtx.AddMessage(msg)
		}
		
		// Increment iteration
		runCtx.IncrementIteration()
		
		// Continue loop to get LLM's response to tool results
	}
}

// buildLLMOptions builds LLM options including tool definitions.
func (a *Agent) buildLLMOptions() []llm.Option {
	opts := []llm.Option{}
	
	// Add tools if any are registered
	toolList := a.toolRegistry.List()
	if len(toolList) > 0 {
		toolDefs := make([]llm.ToolDef, len(toolList))
		for i, tool := range toolList {
			toolDefs[i] = llm.ToolDef{
				Type: "function",
				Function: llm.FunctionDef{
					Name:        tool.Name(),
					Description: tool.Description(),
					Parameters:  tool.Parameters(),
				},
			}
		}
		opts = append(opts, llm.WithTools(toolDefs))
	}
	
	return opts
}

// Stream executes the agent with streaming responses.
func (a *Agent) Stream(ctx context.Context, query string) (<-chan llm.Token, error) {
	if a.llmProvider == nil {
		return nil, fmt.Errorf("LLM provider not configured")
	}
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	
	// Create run context
	runCtx := NewRunContext(ctx, a.systemPrompt, a.maxIterations)
	
	// Add user message
	runCtx.AddMessage(types.Message{
		Role:    "user",
		Content: query,
	})
	
	// Build LLM options
	llmOpts := a.buildLLMOptions()
	
	// Start streaming
	tokenChan, err := a.llmProvider.Stream(runCtx.Context(), runCtx.Messages(), llmOpts...)
	if err != nil {
		cancel()
		return nil, err
	}
	
	// Wrap the channel to handle cleanup
	wrappedChan := make(chan llm.Token)
	go func() {
		defer close(wrappedChan)
		defer cancel()
		
		for token := range tokenChan {
			wrappedChan <- token
		}
	}()
	
	return wrappedChan, nil
}
