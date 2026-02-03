// Package tools provides the tool system for agent tool execution.
package tools

import (
	"context"
	"fmt"
	"sync"

	"github.com/EmadMokhtar/gogent/types"
)

// Executor executes tools and manages their execution.
type Executor struct {
	registry *Registry
}

// NewExecutor creates a new tool executor.
func NewExecutor(registry *Registry) *Executor {
	return &Executor{
		registry: registry,
	}
}

// ExecuteToolCall executes a single tool call.
func (e *Executor) ExecuteToolCall(ctx context.Context, toolCall types.ToolCall) (interface{}, error) {
	// Get the tool from registry
	tool, err := e.registry.Get(toolCall.Function.Name)
	if err != nil {
		return nil, fmt.Errorf("tool not found: %w", err)
	}
	
	// Parse arguments
	args, err := parseArguments(toolCall.Function.Arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}
	
	// Execute the tool
	result, err := tool.Execute(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}
	
	return result, nil
}

// ExecuteToolCalls executes multiple tool calls in parallel.
func (e *Executor) ExecuteToolCalls(ctx context.Context, toolCalls []types.ToolCall) ([]types.Message, error) {
	if len(toolCalls) == 0 {
		return []types.Message{}, nil
	}
	
	// Execute tools in parallel
	var wg sync.WaitGroup
	results := make([]types.Message, len(toolCalls))
	errors := make([]error, len(toolCalls))
	
	for i, tc := range toolCalls {
		wg.Add(1)
		go func(idx int, toolCall types.ToolCall) {
			defer wg.Done()
			
			result, err := e.ExecuteToolCall(ctx, toolCall)
			if err != nil {
				errors[idx] = err
				results[idx] = types.Message{
					Role:       "tool",
					ToolCallID: toolCall.ID,
					Content:    fmt.Sprintf("Error: %s", err.Error()),
					Name:       toolCall.Function.Name,
				}
				return
			}
			
			// Format result as JSON string
			resultStr, err := formatResult(result)
			if err != nil {
				errors[idx] = err
				results[idx] = types.Message{
					Role:       "tool",
					ToolCallID: toolCall.ID,
					Content:    fmt.Sprintf("Error formatting result: %s", err.Error()),
					Name:       toolCall.Function.Name,
				}
				return
			}
			
			results[idx] = types.Message{
				Role:       "tool",
				ToolCallID: toolCall.ID,
				Content:    resultStr,
				Name:       toolCall.Function.Name,
			}
		}(i, tc)
	}
	
	wg.Wait()
	
	// Check if any errors occurred (non-fatal, we still return tool messages)
	// The LLM will see the error messages and can handle them
	
	return results, nil
}
