// Package tools provides the tool system for agent tool execution.
package tools

import (
	"encoding/json"
	"fmt"
)

// parseArguments parses JSON-encoded tool arguments.
func parseArguments(argsJSON string) (map[string]interface{}, error) {
	if argsJSON == "" {
		return map[string]interface{}{}, nil
	}
	
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}
	
	return args, nil
}

// formatResult formats a tool result as a JSON string.
func formatResult(result interface{}) (string, error) {
	// If result is already a string, return it
	if str, ok := result.(string); ok {
		return str, nil
	}
	
	// Otherwise, encode as JSON
	data, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}
	
	return string(data), nil
}
