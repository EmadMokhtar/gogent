// Package tools provides the tool system for agent tool execution.
package tools

import (
	"context"
)

// Tool defines the interface for agent tools.
type Tool interface {
	// Name returns the tool's name.
	Name() string
	
	// Description returns a description of what the tool does.
	Description() string
	
	// Parameters returns the tool's parameter schema.
	Parameters() map[string]interface{}
	
	// Execute runs the tool with the given input.
	Execute(ctx context.Context, input map[string]interface{}) (interface{}, error)
}

// Parameter represents a tool parameter definition.
type Parameter struct {
	Name        string      // Parameter name
	Type        string      // Parameter type (string, number, boolean, etc.)
	Description string      // Parameter description
	Required    bool        // Whether the parameter is required
	Default     interface{} // Default value if not provided
}
