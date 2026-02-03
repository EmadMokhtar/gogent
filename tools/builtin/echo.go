package builtin

import (
	"context"
	"fmt"

	"github.com/EmadMokhtar/gogent/tools"
)

// Echo is a tool that echoes back the input (useful for testing)
type Echo struct{}

// NewEcho creates a new echo tool
func NewEcho() *Echo {
	return &Echo{}
}

// Name returns the name of the tool
func (e *Echo) Name() string {
	return "echo"
}

// Description returns a description of what the tool does
func (e *Echo) Description() string {
	return "Echoes back the provided message (useful for testing)"
}

// Parameters returns the list of parameters the tool accepts
func (e *Echo) Parameters() []tools.Parameter {
	return []tools.Parameter{
		{
			Name:        "message",
			Description: "The message to echo back",
			Type:        "string",
			Required:    true,
		},
	}
}

// Execute runs the tool with the given input
func (e *Echo) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	msgVal, ok := input["message"]
	if !ok {
		return nil, fmt.Errorf("missing required parameter: message")
	}
	
	message, ok := msgVal.(string)
	if !ok {
		return nil, fmt.Errorf("message must be a string")
	}
	
	return message, nil
}
