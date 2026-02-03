// Package builtin provides built-in tools for the Gogent framework.
package builtin

import (
	"context"
	"fmt"
)

// Calculator is a tool for performing basic arithmetic operations.
type Calculator struct{}

// NewCalculator creates a new calculator tool.
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Name returns the tool's name.
func (c *Calculator) Name() string {
	return "calculator"
}

// Description returns the tool's description.
func (c *Calculator) Description() string {
	return "Performs basic arithmetic operations. Supports operations: add, subtract, multiply, divide, power, modulo."
}

// Parameters returns the tool's parameter schema as a JSON schema.
func (c *Calculator) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"operation": map[string]interface{}{
				"type":        "string",
				"description": "The operation to perform",
				"enum":        []string{"add", "subtract", "multiply", "divide", "power", "modulo"},
			},
			"a": map[string]interface{}{
				"type":        "number",
				"description": "The first operand",
			},
			"b": map[string]interface{}{
				"type":        "number",
				"description": "The second operand",
			},
		},
		"required": []string{"operation", "a", "b"},
	}
}

// Execute performs the calculation.
func (c *Calculator) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	// Extract parameters
	operation, ok := input["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation must be a string")
	}
	
	// Extract operands - handle both float64 and int
	var a, b float64
	
	switch v := input["a"].(type) {
	case float64:
		a = v
	case int:
		a = float64(v)
	case int64:
		a = float64(v)
	default:
		return nil, fmt.Errorf("operand 'a' must be a number")
	}
	
	switch v := input["b"].(type) {
	case float64:
		b = v
	case int:
		b = float64(v)
	case int64:
		b = float64(v)
	default:
		return nil, fmt.Errorf("operand 'b' must be a number")
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
	case "power":
		// Use math.Pow for accurate power calculations
		result = 1
		for i := 0; i < int(b); i++ {
			result *= a
		}
	case "modulo":
		if b == 0 {
			return nil, fmt.Errorf("modulo by zero")
		}
		// Truncate to integers for modulo operation
		result = float64(int(a) % int(b))
	default:
		return nil, fmt.Errorf("unsupported operation: %s", operation)
	}
	
	return result, nil
}
