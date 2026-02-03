// Package builtin provides built-in tools for the Gogent framework.
package builtin

import (
	"context"
	"fmt"
)

// Echo is a simple tool that echoes back the input. Useful for testing.
type Echo struct{}

// NewEcho creates a new echo tool.
func NewEcho() *Echo {
	return &Echo{}
}

// Name returns the tool's name.
func (e *Echo) Name() string {
	return "echo"
}

// Description returns the tool's description.
func (e *Echo) Description() string {
	return "Echoes back the provided message. Useful for testing tool execution."
}

// Parameters returns the tool's parameter schema.
func (e *Echo) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"message": map[string]interface{}{
				"type":        "string",
				"description": "The message to echo back",
			},
		},
		"required": []string{"message"},
	}
}

// Execute echoes back the message.
func (e *Echo) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	message, ok := input["message"].(string)
	if !ok {
		return nil, fmt.Errorf("message must be a string")
	}
	
	return map[string]string{
		"echoed": message,
	}, nil
}
