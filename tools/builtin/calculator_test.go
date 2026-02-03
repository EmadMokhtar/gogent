// Package builtin provides built-in tools for the Gogent framework.
package builtin

import (
	"context"
	"testing"
)

func TestCalculator_Name(t *testing.T) {
	calc := NewCalculator()
	if calc.Name() != "calculator" {
		t.Errorf("expected name 'calculator', got %s", calc.Name())
	}
}

func TestCalculator_Description(t *testing.T) {
	calc := NewCalculator()
	desc := calc.Description()
	if desc == "" {
		t.Error("description should not be empty")
	}
}

func TestCalculator_Parameters(t *testing.T) {
	calc := NewCalculator()
	params := calc.Parameters()
	if params == nil {
		t.Error("parameters should not be nil")
	}
}

func TestCalculator_Execute_Add(t *testing.T) {
	calc := NewCalculator()
	ctx := context.Background()
	
	result, err := calc.Execute(ctx, map[string]interface{}{
		"operation": "add",
		"a":         float64(10),
		"b":         float64(5),
	})
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if result != float64(15) {
		t.Errorf("expected 15, got %v", result)
	}
}

func TestCalculator_Execute_Subtract(t *testing.T) {
	calc := NewCalculator()
	ctx := context.Background()
	
	result, err := calc.Execute(ctx, map[string]interface{}{
		"operation": "subtract",
		"a":         float64(10),
		"b":         float64(5),
	})
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if result != float64(5) {
		t.Errorf("expected 5, got %v", result)
	}
}

func TestCalculator_Execute_Multiply(t *testing.T) {
	calc := NewCalculator()
	ctx := context.Background()
	
	result, err := calc.Execute(ctx, map[string]interface{}{
		"operation": "multiply",
		"a":         float64(10),
		"b":         float64(5),
	})
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if result != float64(50) {
		t.Errorf("expected 50, got %v", result)
	}
}

func TestCalculator_Execute_Divide(t *testing.T) {
	calc := NewCalculator()
	ctx := context.Background()
	
	result, err := calc.Execute(ctx, map[string]interface{}{
		"operation": "divide",
		"a":         float64(10),
		"b":         float64(5),
	})
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if result != float64(2) {
		t.Errorf("expected 2, got %v", result)
	}
}

func TestCalculator_Execute_DivideByZero(t *testing.T) {
	calc := NewCalculator()
	ctx := context.Background()
	
	_, err := calc.Execute(ctx, map[string]interface{}{
		"operation": "divide",
		"a":         float64(10),
		"b":         float64(0),
	})
	
	if err == nil {
		t.Error("expected error for division by zero")
	}
}

func TestCalculator_Execute_Power(t *testing.T) {
	calc := NewCalculator()
	ctx := context.Background()
	
	result, err := calc.Execute(ctx, map[string]interface{}{
		"operation": "power",
		"a":         float64(2),
		"b":         float64(3),
	})
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if result != float64(8) {
		t.Errorf("expected 8, got %v", result)
	}
}

func TestCalculator_Execute_Modulo(t *testing.T) {
	calc := NewCalculator()
	ctx := context.Background()
	
	result, err := calc.Execute(ctx, map[string]interface{}{
		"operation": "modulo",
		"a":         float64(10),
		"b":         float64(3),
	})
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if result != float64(1) {
		t.Errorf("expected 1, got %v", result)
	}
}

func TestCalculator_Execute_UnsupportedOperation(t *testing.T) {
	calc := NewCalculator()
	ctx := context.Background()
	
	_, err := calc.Execute(ctx, map[string]interface{}{
		"operation": "invalid",
		"a":         float64(10),
		"b":         float64(5),
	})
	
	if err == nil {
		t.Error("expected error for unsupported operation")
	}
}

func TestCalculator_Execute_InvalidInput(t *testing.T) {
	calc := NewCalculator()
	ctx := context.Background()
	
	_, err := calc.Execute(ctx, map[string]interface{}{
		"operation": "add",
		"a":         "not a number",
		"b":         float64(5),
	})
	
	if err == nil {
		t.Error("expected error for invalid input")
	}
}
