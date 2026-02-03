package builtin

import (
	"context"
	"fmt"

	"github.com/EmadMokhtar/gogent/tools"
)

// Calculator is a tool that performs basic arithmetic operations
type Calculator struct{}

// NewCalculator creates a new calculator tool
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Name returns the name of the tool
func (c *Calculator) Name() string {
	return "calculator"
}

// Description returns a description of what the tool does
func (c *Calculator) Description() string {
	return "Performs basic arithmetic operations (add, subtract, multiply, divide) on two numbers"
}

// Parameters returns the list of parameters the tool accepts
func (c *Calculator) Parameters() []tools.Parameter {
	return []tools.Parameter{
		{
			Name:        "operation",
			Description: "The operation to perform: 'add', 'subtract', 'multiply', or 'divide'",
			Type:        "string",
			Required:    true,
		},
		{
			Name:        "a",
			Description: "The first number",
			Type:        "number",
			Required:    true,
		},
		{
			Name:        "b",
			Description: "The second number",
			Type:        "number",
			Required:    true,
		},
	}
}

// Execute runs the tool with the given input
func (c *Calculator) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	// Extract operation
	opVal, ok := input["operation"]
	if !ok {
		return nil, fmt.Errorf("missing required parameter: operation")
	}
	operation, ok := opVal.(string)
	if !ok {
		return nil, fmt.Errorf("operation must be a string")
	}
	
	// Extract a
	aVal, ok := input["a"]
	if !ok {
		return nil, fmt.Errorf("missing required parameter: a")
	}
	var a float64
	switch v := aVal.(type) {
	case float64:
		a = v
	case int:
		a = float64(v)
	case int64:
		a = float64(v)
	default:
		return nil, fmt.Errorf("parameter 'a' must be a number")
	}
	
	// Extract b
	bVal, ok := input["b"]
	if !ok {
		return nil, fmt.Errorf("missing required parameter: b")
	}
	var b float64
	switch v := bVal.(type) {
	case float64:
		b = v
	case int:
		b = float64(v)
	case int64:
		b = float64(v)
	default:
		return nil, fmt.Errorf("parameter 'b' must be a number")
	}
	
	// Perform operation
	var result float64
	switch operation {
	case "add":
		result = a + b
	case "subtract":
		result = a - b
	case "multiply":
		result = a * b
	case "divide":
		if b == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		result = a / b
	default:
		return nil, fmt.Errorf("unknown operation: %s (must be add, subtract, multiply, or divide)", operation)
	}
	
	return result, nil
}
