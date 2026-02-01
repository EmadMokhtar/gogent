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
	
	// Use a channel to collect results
	type resultWithIndex struct {
		index  int
		result *ExecutionResult
		err    error
	}
	
	resultChan := make(chan resultWithIndex, len(requests))
	
	// Execute tools in parallel
	for i, req := range requests {
		go func(idx int, r struct {
			ToolName string
			Input    map[string]interface{}
		}) {
			result, err := e.Execute(ctx, r.ToolName, r.Input)
			resultChan <- resultWithIndex{
				index:  idx,
				result: result,
				err:    err,
			}
		}(i, req)
	}
	
	// Collect results
	var firstError error
	for i := 0; i < len(requests); i++ {
		res := <-resultChan
		results[res.index] = res.result
		if res.err != nil && firstError == nil {
			firstError = res.err
		}
	}
	
	return results, firstError
}
