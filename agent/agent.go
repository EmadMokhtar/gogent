package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/EmadMokhtar/gogent/llm"
	"github.com/EmadMokhtar/gogent/tools"
)

// Agent represents an autonomous agent
type Agent struct {
	opts     Options
	runner   *Runner
	registry *tools.Registry
	executor *tools.Executor
}

// New creates a new agent with the given options
func New(opts ...Option) *Agent {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}
	
	// Create tool registry and executor
	registry := tools.NewRegistry()
	for _, tool := range options.Tools {
		_ = registry.Register(tool)
	}
	executor := tools.NewExecutor(registry)
	
	agent := &Agent{
		opts:     options,
		registry: registry,
		executor: executor,
	}
	
	agent.runner = NewRunner(agent)
	
	return agent
}

// Run executes the agent with a query
func (a *Agent) Run(ctx context.Context, query string) (*Result, error) {
	if a.opts.LLM == nil {
		return nil, fmt.Errorf("LLM provider not configured")
	}
	
	startTime := time.Now()
	
	// Prepare messages
	messages, err := a.runner.prepareMessages(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare messages: %w", err)
	}
	
	// Get RAG sources if available
	var ragSources []interface{}
	if a.opts.RAG != nil {
		results, err := a.opts.RAG.Query(ctx, query, 3)
		if err == nil {
			// Convert to interface{} slice for Result
			for _, r := range results {
				ragSources = append(ragSources, r)
			}
		}
	}
	
	// Call LLM
	llmOpts := []llm.Option{
		llm.WithMaxTokens(a.opts.MaxTokens),
		llm.WithTemperature(a.opts.Temperature),
	}
	
	response, err := a.opts.LLM.Generate(ctx, messages, llmOpts...)
	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}
	
	// Save to memory
	if err := a.runner.saveToMemory(ctx, query, response.Content); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to save to memory: %v\n", err)
	}
	
	result := &Result{
		Response: response.Content,
		Sources:  ragSources,
		Metadata: make(map[string]interface{}),
		Duration: time.Since(startTime),
	}
	
	return result, nil
}

// RunWithContext executes the agent with a custom context
func (a *Agent) RunWithContext(agentCtx *Context, query string) (*Result, error) {
	return a.Run(agentCtx.Context(), query)
}

// Stream executes the agent with streaming
func (a *Agent) Stream(ctx context.Context, query string) (<-chan StreamChunk, error) {
	if a.opts.LLM == nil {
		return nil, fmt.Errorf("LLM provider not configured")
	}
	
	ch := make(chan StreamChunk, 10)
	
	go func() {
		defer close(ch)
		
		// Prepare messages
		messages, err := a.runner.prepareMessages(ctx, query)
		if err != nil {
			ch <- StreamChunk{Error: fmt.Errorf("failed to prepare messages: %w", err)}
			return
		}
		
		// Call LLM stream
		llmOpts := []llm.Option{
			llm.WithMaxTokens(a.opts.MaxTokens),
			llm.WithTemperature(a.opts.Temperature),
			llm.WithStream(),
		}
		
		tokenChan, err := a.opts.LLM.Stream(ctx, messages, llmOpts...)
		if err != nil {
			ch <- StreamChunk{Error: fmt.Errorf("LLM streaming failed: %w", err)}
			return
		}
		
		fullResponse := ""
		for token := range tokenChan {
			if token.Error != nil {
				ch <- StreamChunk{Error: token.Error}
				return
			}
			
			fullResponse += token.Content
			ch <- StreamChunk{
				Content: token.Content,
				Done:    token.Done,
			}
			
			if token.Done {
				// Save to memory
				if err := a.runner.saveToMemory(ctx, query, fullResponse); err != nil {
					fmt.Printf("Warning: failed to save to memory: %v\n", err)
				}
			}
		}
	}()
	
	return ch, nil
}

// GetRegistry returns the tool registry
func (a *Agent) GetRegistry() *tools.Registry {
	return a.registry
}

// GetExecutor returns the tool executor
func (a *Agent) GetExecutor() *tools.Executor {
	return a.executor
}
