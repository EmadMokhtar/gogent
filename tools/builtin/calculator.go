package builtin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/EmadMokhtar/gogent/tools"
)

// Calculator is a tool for performing basic calculations
type Calculator struct{}

// NewCalculator creates a new calculator tool
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Name returns the tool name
func (c *Calculator) Name() string {
	return "calculator"
}

// Description returns the tool description
func (c *Calculator) Description() string {
	return "Performs basic arithmetic operations (add, subtract, multiply, divide)"
}

// Parameters returns the tool parameters
func (c *Calculator) Parameters() []tools.Parameter {
	return []tools.Parameter{
		{
			Name:        "operation",
			Type:        tools.TypeString,
			Description: "The operation to perform: add, subtract, multiply, divide",
			Required:    true,
			Enum:        []string{"add", "subtract", "multiply", "divide"},
		},
		{
			Name:        "a",
			Type:        tools.TypeNumber,
			Description: "First number",
			Required:    true,
		},
		{
			Name:        "b",
			Type:        tools.TypeNumber,
			Description: "Second number",
			Required:    true,
		},
	}
}

// Execute runs the calculator tool
func (c *Calculator) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	operation, ok := input["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation must be a string")
	}
	
	a, err := toFloat64(input["a"])
	if err != nil {
		return nil, fmt.Errorf("invalid value for a: %w", err)
	}
	
	b, err := toFloat64(input["b"])
	if err != nil {
		return nil, fmt.Errorf("invalid value for b: %w", err)
	}
	
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
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}
	
	return result, nil
}

// toFloat64 converts an interface{} to float64
func toFloat64(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}
