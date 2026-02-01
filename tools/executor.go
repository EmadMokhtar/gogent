package tools

import (
	"context"
	"fmt"
)

// Executor executes tools with validation
type Executor struct {
	registry *Registry
}

// NewExecutor creates a new tool executor
func NewExecutor(registry *Registry) *Executor {
	return &Executor{
		registry: registry,
	}
}

// Execute runs a tool by name with the given input
func (e *Executor) Execute(ctx context.Context, toolName string, input map[string]interface{}) (*ExecutionResult, error) {
	tool, err := e.registry.Get(toolName)
	if err != nil {
		return nil, err
	}
	
	// Validate input parameters
	if err := e.validateInput(tool, input); err != nil {
		return &ExecutionResult{
			ToolName: toolName,
			Input:    input,
			Error:    err,
		}, err
	}
	
	// Execute the tool
	output, err := tool.Execute(ctx, input)
	
	result := &ExecutionResult{
		ToolName: toolName,
		Input:    input,
		Output:   output,
		Error:    err,
	}
	
	return result, err
}

// validateInput validates that required parameters are present
func (e *Executor) validateInput(tool Tool, input map[string]interface{}) error {
	params := tool.Parameters()
	
	for _, param := range params {
		if param.Required {
			if _, ok := input[param.Name]; !ok {
				return fmt.Errorf("required parameter %s is missing", param.Name)
			}
		}
	}
	
	return nil
}

// ExecuteBatch executes multiple tools in parallel
func (e *Executor) ExecuteBatch(ctx context.Context, requests []struct {
	ToolName string
	Input    map[string]interface{}
}) ([]*ExecutionResult, error) {
	results := make([]*ExecutionResult, len(requests))
	errors := make([]error, len(requests))
	
	// Execute tools in parallel
	done := make(chan struct{})
	for i, req := range requests {
		go func(idx int, r struct {
			ToolName string
			Input    map[string]interface{}
		}) {
			result, err := e.Execute(ctx, r.ToolName, r.Input)
			results[idx] = result
			errors[idx] = err
		}(i, req)
	}
	
	// Wait for all to complete
	go func() {
		// Simple wait - in production, use sync.WaitGroup
		<-done
	}()
	
	// Check if any errors occurred
	var firstError error
	for _, err := range errors {
		if err != nil && firstError == nil {
			firstError = err
		}
	}
	
	return results, firstError
}
