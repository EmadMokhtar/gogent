package tools

import (
	"context"
	"fmt"
)

// Parameter represents a tool parameter
type Parameter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"` // "string", "number", "integer", "boolean", "object", "array"
	Required    bool   `json:"required"`
}

// Tool is the interface that all tools must implement
type Tool interface {
	// Name returns the name of the tool
	Name() string
	
	// Description returns a description of what the tool does
	Description() string
	
	// Parameters returns the list of parameters the tool accepts
	Parameters() []Parameter
	
	// Execute runs the tool with the given input
	Execute(ctx context.Context, input map[string]interface{}) (interface{}, error)
}

// ToolError represents an error from tool execution
type ToolError struct {
	Tool    string
	Message string
	Err     error
}

func (e *ToolError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("tool %s error: %s: %v", e.Tool, e.Message, e.Err)
	}
	return fmt.Sprintf("tool %s error: %s", e.Tool, e.Message)
}

func (e *ToolError) Unwrap() error {
	return e.Err
}
