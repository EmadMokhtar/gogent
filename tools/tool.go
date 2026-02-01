package tools

import (
	"context"
)

// ParameterType represents the type of a tool parameter
type ParameterType string

const (
	TypeString  ParameterType = "string"
	TypeNumber  ParameterType = "number"
	TypeBoolean ParameterType = "boolean"
	TypeObject  ParameterType = "object"
	TypeArray   ParameterType = "array"
)

// Parameter represents a tool parameter
type Parameter struct {
	Name        string        `json:"name"`
	Type        ParameterType `json:"type"`
	Description string        `json:"description"`
	Required    bool          `json:"required"`
	Enum        []string      `json:"enum,omitempty"`
	Default     interface{}   `json:"default,omitempty"`
}

// Tool is the interface that all tools must implement
type Tool interface {
	// Name returns the tool name
	Name() string
	
	// Description returns a description of what the tool does
	Description() string
	
	// Parameters returns the list of parameters the tool accepts
	Parameters() []Parameter
	
	// Execute runs the tool with the given input
	Execute(ctx context.Context, input map[string]interface{}) (interface{}, error)
}

// ExecutionResult represents the result of a tool execution
type ExecutionResult struct {
	ToolName string                 `json:"tool_name"`
	Input    map[string]interface{} `json:"input"`
	Output   interface{}            `json:"output"`
	Error    error                  `json:"error,omitempty"`
}
